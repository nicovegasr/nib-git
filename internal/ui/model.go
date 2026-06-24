// Package ui holds the Bubble Tea model. This is a SCAFFOLD: the commit flow
// is wired enough to load files and move the cursor; staging, the type
// dropdown, the git-log panel, and slash-commands are marked TODO and map 1:1
// to the states in docs/ux-mockups.html.
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nicovegasr/nib-git/internal/commit"
	"github.com/nicovegasr/nib-git/internal/git"
)

// Model is the root application state.
type Model struct {
	files    []git.FileChange
	staged   map[string]bool // path -> staged in this session
	cursor   int
	typeIdx  int    // index into commit.Types
	message  string // commit message being typed
	focusMsg bool   // true once the user tabbed into the message field
	err      error
	// TODO(settings): showLog, mouseCapture, logScope — read from settings pkg.
	// TODO(audit): slash-command palette state.
}

// New builds the initial model and loads the working tree.
func New() Model {
	m := Model{staged: map[string]bool{}}
	if files, err := git.Status(); err != nil {
		m.err = err
	} else {
		m.files = files
		for _, f := range files {
			if f.Staged {
				m.staged[f.Path] = true
			}
		}
	}
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		// TODO(mouse): handle tea.MouseMsg clicks on files / dropdown.
		return m, nil
	}

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
	case " ": // space: toggle stage on the file under the cursor
		if len(m.files) > 0 {
			p := m.files[m.cursor].Path
			m.staged[p] = !m.staged[p]
		}
	case "tab": // jump to the message field
		m.focusMsg = true
	case "enter":
		return m, m.commit()
		// TODO: type shortcuts (f/x/r/g…) -> set m.typeIdx from commit.Types.
		// TODO: "/" -> open slash-command palette (/audit, /settings).
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

// commit stages the selected files and creates the commit. TODO: surface
// success/failure in the UI instead of just quitting.
func (m Model) commit() tea.Cmd {
	var files []string
	for p, on := range m.staged {
		if on {
			files = append(files, p)
		}
	}
	return func() tea.Msg {
		if err := git.Stage(files); err != nil {
			return err
		}
		subject := commit.Format(commit.Types[m.typeIdx].Prefix, m.message)
		_ = git.Commit(subject) // TODO: report result
		return tea.Quit()
	}
}

// View implements tea.Model. Placeholder render — see docs/ux-mockups.html for
// the target layout (file list + inline commit line + optional git-log panel).
func (m Model) View() string {
	if m.err != nil {
		return "nib: " + m.err.Error() + "\n"
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
	b.WriteString("\n j/k mover · espacio stage · tab mensaje · enter commit · q salir\n")
	return b.String()
}
