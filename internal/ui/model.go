// Package ui holds the Bubble Tea model for the minimal commit flow:
// mark files → pick type → type message → enter. It depends on a narrow Repo
// port (ADR-0005) so the state transitions are testable without real git.
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
	focusMsg bool   // true once the user tabbed into the message field
	done     bool   // commit succeeded; ready to quit
	err      error

	shortcuts map[string]int // type shortcut key -> index into commit.Types
}

// New builds the initial model and loads the working tree through repo.
func New(repo Repo) Model {
	m := Model{
		repo:      repo,
		staged:    map[string]bool{},
		shortcuts: typeShortcuts(),
	}
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
		m.done = true
		return m, tea.Quit
	case tea.KeyMsg:
		return m.updateKey(msg)
	}
	// TODO(mouse): handle tea.MouseMsg clicks on files / dropdown (Phase 2).
	return m, nil
}

func (m Model) updateKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.focusMsg {
		return m.updateMessage(key)
	}

	switch key.String() {
	case "q", "ctrl+c":
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
	case "tab": // jump to the message field
		m.focusMsg = true
	case "enter":
		return m, m.commit()
	default:
		if idx, ok := m.shortcuts[key.String()]; ok {
			m.typeIdx = idx
		}
	}
	return m, nil
}

func (m Model) updateMessage(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.focusMsg = false
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

// View implements tea.Model. See docs/ux-mockups.html for the target layout.
func (m Model) View() string {
	if m.err != nil {
		return "nib: " + m.err.Error() + "\n"
	}
	if m.done {
		return "nib: commit creado.\n"
	}

	var b strings.Builder
	b.WriteString("  nib — cambios\n\n")
	for i, f := range m.files {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		mark := "○"
		if m.staged[f.Path] {
			mark = "✓"
		}
		fmt.Fprintf(&b, " %s %s %c %s\n", cursor, mark, f.Status, f.Path)
	}
	t := commit.Types[m.typeIdx]
	msg := m.message
	if msg == "" {
		msg = "(escribe el mensaje)"
	}
	fmt.Fprintf(&b, "\n [ %s ▾ ] %s\n", t.Name, msg)
	b.WriteString("\n j/k mover · espacio stage · letra tipo · tab mensaje · enter commit · q salir\n")
	return b.String()
}
