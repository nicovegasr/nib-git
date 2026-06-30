// Package ui holds the Bubble Tea model for the minimal commit flow:
// mark files → pick type → type message → enter. It depends on a narrow Repo
// port (ADR-0005) so the state transitions are testable without real git.
//
// Phase 2 adds three focus zones (Files → Type → Message), an openable type
// dropdown, and mouse clicks. Picking a type auto-advances to the message so
// the keyboard hot path stays a single `tab` + letter. The type shortcut only
// fires while the Type zone is focused, keeping it predictable.
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nicovegasr/nib-git/internal/commit"
	"github.com/nicovegasr/nib-git/internal/git"
)

// Repo is the narrow slice of git the commit flow needs. main.go injects a real
// adapter; tests inject a fake.
type Repo interface {
	Status() ([]git.FileChange, error)
	Stage(files []string) error
	Commit(subject string) error
}

// zone is the part of the UI currently driving the keyboard.
type zone int

const (
	zoneFiles zone = iota
	zoneType
	zoneMessage
)

// commitDoneMsg reports the outcome of a commit attempt.
type commitDoneMsg struct{ err error }

// Model is the root application state.
type Model struct {
	repo     Repo
	files    []git.FileChange
	staged   map[string]bool // path -> staged in this session
	cursor   int
	typeIdx  int    // index into commit.Types
	message  string // commit message being typed
	focus    zone   // which zone owns the keyboard
	dropOpen bool   // type dropdown expanded
	dropIdx  int    // highlighted item while the dropdown is open
	flash    string // transient success feedback, cleared on the next keystroke
	err      error

	width   int
	height  int
	showLog bool // git-log panel enabled (Phase 4 /settings will toggle it)

	shortcuts map[string]int // type shortcut key -> index into commit.Types
}

// New builds the initial model and loads the working tree through repo.
func New(repo Repo) Model {
	m := Model{
		repo:      repo,
		staged:    map[string]bool{},
		shortcuts: typeShortcuts(),
	}
	m.showLog = true
	files, err := repo.Status()
	if err != nil {
		m.err = err
		return m
	}
	m.files = files
	for _, f := range files {
		if f.Staged {
			m.staged[f.Path] = true
		}
	}
	return m
}

// typeShortcuts maps each commit type's shortcut rune to its index.
func typeShortcuts() map[string]int {
	m := make(map[string]int, len(commit.Types))
	for i, t := range commit.Types {
		m[string(t.Shortcut)] = i
	}
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case commitDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		// Persistent session (ADR-0007): refresh state, never quit.
		m.message = ""
		m.flash = "✓ commit creado"
		files, err := m.repo.Status()
		if err != nil {
			m.err = err
			return m, nil
		}
		m.files = files
		m.staged = map[string]bool{}
		for _, f := range files {
			if f.Staged {
				m.staged[f.Path] = true
			}
		}
		if m.cursor >= len(files) {
			m.cursor = 0
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		return m.updateKey(msg)
	case tea.MouseMsg:
		return m.updateMouse(msg)
	}
	return m, nil
}

func (m Model) updateKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.flash = "" // any keystroke dismisses the success flash
	if key.String() == "ctrl+c" {
		return m, tea.Quit
	}
	switch m.focus {
	case zoneFiles:
		return m.updateFiles(key)
	case zoneType:
		return m.updateType(key)
	default:
		return m.updateMessage(key)
	}
}

func (m Model) updateFiles(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "q":
		return m, tea.Quit
	case "j", "down":
		if m.cursor < len(m.files)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case " ": // toggle stage on the file under the cursor
		if len(m.files) > 0 {
			p := m.files[m.cursor].Path
			m.staged[p] = !m.staged[p]
		}
	case "tab":
		m.focus = zoneType
	case "shift+tab":
		m.focus = zoneMessage
	case "esc":
		return m, tea.Quit
	case "enter":
		return m, m.commit()
	}
	return m, nil
}

// updateType drives the type zone: the dropdown is a discoverable cheat-sheet,
// while the letter shortcut stays the expert path. Selecting a type (letter or
// dropdown) auto-advances to the message field.
func (m Model) updateType(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.dropOpen {
		switch key.String() {
		case "j", "down":
			if m.dropIdx < len(commit.Types)-1 {
				m.dropIdx++
			}
		case "k", "up":
			if m.dropIdx > 0 {
				m.dropIdx--
			}
		case "enter", " ":
			return m.pickType(m.dropIdx), nil
		case "esc":
			m.dropOpen = false
		default:
			if idx, ok := m.shortcuts[key.String()]; ok {
				return m.pickType(idx), nil
			}
		}
		return m, nil
	}

	switch key.String() {
	case "down", "enter", " ":
		m.dropOpen = true
		m.dropIdx = m.typeIdx
	case "tab":
		m.focus = zoneMessage
	case "shift+tab":
		m.focus = zoneFiles
	case "esc":
		m.focus = zoneFiles
	default:
		if idx, ok := m.shortcuts[key.String()]; ok {
			return m.pickType(idx), nil
		}
	}
	return m, nil
}

// pickType selects a commit type, closes the dropdown, and advances focus to
// the message field so the hand never leaves the keyboard.
func (m Model) pickType(idx int) Model {
	m.typeIdx = idx
	m.dropOpen = false
	m.focus = zoneMessage
	return m
}

func (m Model) updateMessage(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.focus = zoneType
	case "tab":
		m.focus = zoneFiles
	case "shift+tab":
		m.focus = zoneType
	case "enter":
		return m, m.commit()
	case "backspace":
		if m.message != "" {
			m.message = m.message[:len(m.message)-1]
		}
	default:
		if len(key.String()) == 1 {
			m.message += key.String()
		}
	}
	return m, nil
}

// updateMouse handles left clicks (mapped through hitTest) and the wheel, which
// scrolls the open dropdown. The mouse-capture toggle lives in /settings
// (Phase 4); here it is always on (main.go enables it).
func (m Model) updateMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.dropOpen && msg.Button == tea.MouseButtonWheelUp {
		if m.dropIdx > 0 {
			m.dropIdx--
		}
		return m, nil
	}
	if m.dropOpen && msg.Button == tea.MouseButtonWheelDown {
		if m.dropIdx < len(commit.Types)-1 {
			m.dropIdx++
		}
		return m, nil
	}
	if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
		return m, nil
	}

	h := m.hitTest(msg.X, msg.Y)
	switch h.kind {
	case hitFile:
		m.cursor = h.idx
		m.focus = zoneFiles
		p := m.files[h.idx].Path
		m.staged[p] = !m.staged[p]
	case hitType:
		m.focus = zoneType
		m.dropOpen = !m.dropOpen
		m.dropIdx = m.typeIdx
	case hitDropItem:
		return m.pickType(h.idx), nil
	case hitMessage:
		m.focus = zoneMessage
	case hitNone:
		m.dropOpen = false
	}
	return m, nil
}

// stagedFiles returns the paths currently marked for commit, in file order so
// the command is deterministic (and testable).
func (m Model) stagedFiles() []string {
	var files []string
	for _, f := range m.files {
		if m.staged[f.Path] {
			files = append(files, f.Path)
		}
	}
	return files
}

// commit stages the selected files and creates the commit, reporting the
// outcome as a commitDoneMsg.
func (m Model) commit() tea.Cmd {
	files := m.stagedFiles()
	subject := commit.Format(commit.Types[m.typeIdx].Prefix, m.message)
	repo := m.repo
	return func() tea.Msg {
		if err := repo.Stage(files); err != nil {
			return commitDoneMsg{err: err}
		}
		return commitDoneMsg{err: repo.Commit(subject)}
	}
}

// --- layout & hit-testing -------------------------------------------------
//
// layout is the SINGLE source of the row geometry: View() draws with it and
// hitTest() reads it, so a click can never drift from what is rendered.

const headerRows = 2 // title + blank line above the file list

const minWidthForLog = 80 // below this the minimal flow keeps the full width

// logVisible reports whether the git-log panel should be drawn: enabled and the
// terminal is wide enough that the minimal flow is not squeezed.
func (m Model) logVisible() bool {
	return m.showLog && m.width >= minWidthForLog
}

// leftColWidth splits the two-column band in half: files left, log right.
//
//nolint:unused // wired into hitTest/View in the next task of this phase.
func (m Model) leftColWidth() int {
	return m.width / 2
}

type layout struct {
	filesStart int // row of the first file
	fileCount  int
	commitRow  int // the "[ type ▾ ] message" line (type and message share it)
	dropStart  int // first dropdown item row (only meaningful when open)
}

func (m Model) layout() layout {
	n := len(m.files)
	commitRow := headerRows + n + 1 // one blank line between files and the commit line
	return layout{
		filesStart: headerRows,
		fileCount:  n,
		commitRow:  commitRow,
		dropStart:  commitRow + 1,
	}
}

type hitKind int

const (
	hitNone hitKind = iota
	hitFile
	hitType
	hitMessage
	hitDropItem
)

type hit struct {
	kind hitKind
	idx  int // file or dropdown-item index, when relevant
}

// hitTest maps a click to a target. Y selects the row; only on the commit line
// does X split the type box from the message field.
func (m Model) hitTest(x, y int) hit {
	l := m.layout()
	if y >= l.filesStart && y < l.filesStart+l.fileCount {
		return hit{hitFile, y - l.filesStart}
	}
	if m.dropOpen && y >= l.dropStart && y < l.dropStart+len(commit.Types) {
		return hit{hitDropItem, y - l.dropStart}
	}
	if y == l.commitRow {
		if x < typeBoxWidth(commit.Types[m.typeIdx].Name) {
			return hit{kind: hitType}
		}
		return hit{kind: hitMessage}
	}
	return hit{kind: hitNone}
}

// typeBoxWidth is the display width of the leading "▸[ name ▾ ] " on the commit
// line, in columns. Used by both View() and hitTest() to agree on where the
// type box ends and the message begins.
func typeBoxWidth(name string) int {
	// "▸" + "[ " + name + " " + "▾" + " ] " = 1 + 2 + len + 1 + 1 + 3
	return len([]rune(name)) + 8
}

// View implements tea.Model. See docs/ux-mockups.html for the target layout.
func (m Model) View() string {
	if m.err != nil {
		return "nib: " + m.err.Error() + "\n"
	}

	var b strings.Builder
	b.WriteString("  nib — cambios\n\n")
	for i, f := range m.files {
		marker := " "
		if m.focus == zoneFiles && i == m.cursor {
			marker = "▸"
		}
		mark := "○"
		if m.staged[f.Path] {
			mark = "✓"
		}
		fmt.Fprintf(&b, "%s %s %c %s\n", marker, mark, f.Status, f.Path)
	}

	b.WriteString("\n")
	b.WriteString(m.commitLine())
	b.WriteString("\n")
	if m.dropOpen {
		b.WriteString(m.dropdownView())
	}
	fmt.Fprintf(&b, "\n %s\n", m.helpLine())
	return b.String()
}

func (m Model) commitLine() string {
	t := commit.Types[m.typeIdx]
	marker := " "
	if m.focus == zoneType {
		marker = "▸"
	}
	arrow := "▾"
	if m.dropOpen {
		arrow = "▴"
	}
	msg := m.message
	if msg == "" {
		msg = "(escribe el mensaje)"
	}
	if m.focus == zoneMessage {
		msg = m.message + "_"
	}
	return fmt.Sprintf("%s[ %s %s ] %s", marker, t.Name, arrow, msg)
}

func (m Model) dropdownView() string {
	var b strings.Builder
	for i, t := range commit.Types {
		sel := " "
		if i == m.dropIdx {
			sel = "▸"
		}
		fmt.Fprintf(&b, "%s %-9s %s\n", sel, t.Name, t.Desc)
	}
	return b.String()
}

func (m Model) helpLine() string {
	switch {
	case m.focus == zoneType && m.dropOpen:
		return "j/k mover · enter elige · esc cierra"
	case m.focus == zoneType:
		return "letra elige tipo · ↓ abre lista · tab mensaje · esc atrás"
	case m.focus == zoneMessage:
		return "escribe · enter commit · esc atrás · tab archivos"
	default:
		return "j/k mover · espacio stage · tab tipo · enter commit · q salir"
	}
}
