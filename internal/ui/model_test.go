package ui

import (
	"errors"
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nicovegasr/nib-git/internal/commit"
	"github.com/nicovegasr/nib-git/internal/git"
)

// fakeRepo is a trivial in-memory Repo for testing model transitions.
type fakeRepo struct {
	status    []git.FileChange
	statusErr error
	stageErr  error
	commitErr error

	stagedArg []string
	commitArg string
}

func (f *fakeRepo) Status() ([]git.FileChange, error) { return f.status, f.statusErr }
func (f *fakeRepo) Stage(files []string) error        { f.stagedArg = files; return f.stageErr }
func (f *fakeRepo) Commit(subject string) error       { f.commitArg = subject; return f.commitErr }

// key builds a tea.KeyMsg whose String() matches the model's switch arms.
func key(s string) tea.KeyMsg {
	switch s {
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "backspace":
		return tea.KeyMsg{Type: tea.KeyBackspace}
	case "space":
		return tea.KeyMsg{Type: tea.KeySpace}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

// send applies a key and returns the resulting Model and command.
func send(m Model, s string) (Model, tea.Cmd) {
	next, cmd := m.Update(key(s))
	return next.(Model), cmd
}

func sampleModel() (Model, *fakeRepo) {
	repo := &fakeRepo{status: []git.FileChange{
		{Path: "a.go", Status: 'M', Staged: false},
		{Path: "b.go", Status: '?', Staged: false},
		{Path: "c.go", Status: 'A', Staged: true},
	}}
	return New(repo), repo
}

func TestNewLoadsFilesAndPreMarksStaged(t *testing.T) {
	m, _ := sampleModel()

	if len(m.files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(m.files))
	}
	if !m.staged["c.go"] {
		t.Error("c.go was staged in git and should be pre-marked")
	}
	if m.staged["a.go"] {
		t.Error("a.go should not be pre-marked")
	}
}

func TestNewPropagatesStatusError(t *testing.T) {
	m := New(&fakeRepo{statusErr: errors.New("boom")})
	if m.err == nil {
		t.Error("expected err set when Status fails")
	}
}

func TestCursorMovesWithinBounds(t *testing.T) {
	m, _ := sampleModel()

	m, _ = send(m, "k") // already at top, stays
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 (clamped at top)", m.cursor)
	}
	m, _ = send(m, "j")
	m, _ = send(m, "j")
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want 2", m.cursor)
	}
	m, _ = send(m, "j") // at bottom, stays
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want 2 (clamped at bottom)", m.cursor)
	}
}

func TestSpaceTogglesStage(t *testing.T) {
	m, _ := sampleModel() // cursor at a.go

	m, _ = send(m, "space")
	if !m.staged["a.go"] {
		t.Error("space should stage a.go")
	}
	m, _ = send(m, "space")
	if m.staged["a.go"] {
		t.Error("space again should unstage a.go")
	}
}

func TestTypeShortcutSelectsType(t *testing.T) {
	m, _ := sampleModel()

	want := -1
	for i, ty := range commit.Types {
		if ty.Shortcut == 'x' { // fix
			want = i
		}
	}
	if want < 0 {
		t.Fatal("no type uses shortcut 'x'")
	}

	m, _ = send(m, "x")
	if m.typeIdx != want {
		t.Errorf("typeIdx = %d, want %d (fix)", m.typeIdx, want)
	}
}

func TestMessageFieldFocusAndEditing(t *testing.T) {
	m, _ := sampleModel()

	m, _ = send(m, "tab")
	if !m.focusMsg {
		t.Fatal("tab should focus the message field")
	}
	// In message mode, "j" is text, not navigation.
	m, _ = send(m, "h")
	m, _ = send(m, "i")
	m, _ = send(m, "j")
	if m.message != "hij" {
		t.Errorf("message = %q, want %q", m.message, "hij")
	}
	if m.cursor != 0 {
		t.Errorf("cursor moved during typing: %d", m.cursor)
	}
	m, _ = send(m, "backspace")
	if m.message != "hi" {
		t.Errorf("message = %q, want %q after backspace", m.message, "hi")
	}
	m, _ = send(m, "esc")
	if m.focusMsg {
		t.Error("esc should leave the message field")
	}
}

func TestCommitStagesSelectedAndFormatsSubject(t *testing.T) {
	m, repo := sampleModel()

	// Select fix, stage a.go (b and c default: c pre-staged), type a message.
	m, _ = send(m, "x")     // fix
	m, _ = send(m, "space") // stage a.go (cursor at a.go)
	m, _ = send(m, "tab")
	for _, ch := range "algo roto" {
		m, _ = send(m, string(ch))
	}

	_, cmd := send(m, "enter")
	if cmd == nil {
		t.Fatal("enter should return a commit command")
	}
	msg := cmd()
	done, ok := msg.(commitDoneMsg)
	if !ok {
		t.Fatalf("expected commitDoneMsg, got %T", msg)
	}
	if done.err != nil {
		t.Fatalf("unexpected commit error: %v", done.err)
	}

	// a.go (just staged) and c.go (pre-staged), in file order.
	want := []string{"a.go", "c.go"}
	if !reflect.DeepEqual(repo.stagedArg, want) {
		t.Errorf("staged = %v, want %v", repo.stagedArg, want)
	}
	if repo.commitArg != "fix: algo roto." {
		t.Errorf("commit subject = %q, want %q", repo.commitArg, "fix: algo roto.")
	}
}

func TestCommitStageErrorIsReported(t *testing.T) {
	m, repo := sampleModel()
	repo.stageErr = errors.New("stage failed")

	_, cmd := send(m, "enter")
	done := cmd().(commitDoneMsg)
	if done.err == nil {
		t.Error("expected stage error to be reported")
	}
	if repo.commitArg != "" {
		t.Error("commit must not run if staging failed")
	}
}

func TestCommitDoneMsgDrivesQuitOrError(t *testing.T) {
	m, _ := sampleModel()

	// Success: marks done and quits.
	next, cmd := m.Update(commitDoneMsg{})
	dm := next.(Model)
	if !dm.done {
		t.Error("successful commit should set done")
	}
	if cmd == nil {
		t.Error("successful commit should return a quit command")
	}

	// Failure: surfaces the error, no quit.
	next, cmd = m.Update(commitDoneMsg{err: errors.New("nope")})
	em := next.(Model)
	if em.err == nil {
		t.Error("failed commit should set err")
	}
	if cmd != nil {
		t.Error("failed commit should not quit")
	}
}

func TestQuitKeys(t *testing.T) {
	m, _ := sampleModel()
	if _, cmd := send(m, "q"); cmd == nil {
		t.Error("q should quit")
	}
}
