package git

import (
	"errors"
	"reflect"
	"testing"
)

// stubRun replaces the package runner for the duration of a test.
func stubRun(t *testing.T, fn func(args ...string) (string, error)) {
	t.Helper()
	orig := run
	run = fn
	t.Cleanup(func() { run = orig })
}

func TestStatusParsesPorcelain(t *testing.T) {
	stubRun(t, func(_ ...string) (string, error) {
		return " M internal/git/git.go\nA  main.go\n?? new.txt\nMM both.go\n", nil
	})

	got, err := Status()
	if err != nil {
		t.Fatalf("Status() error: %v", err)
	}

	want := []FileChange{
		{Path: "internal/git/git.go", Status: 'M', Staged: false},
		{Path: "main.go", Status: 'A', Staged: true},
		{Path: "new.txt", Status: '?', Staged: false},
		{Path: "both.go", Status: 'M', Staged: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Status() = %+v, want %+v", got, want)
	}
}

func TestStatusSkipsShortLines(t *testing.T) {
	stubRun(t, func(_ ...string) (string, error) { return "xy\n M ok.go\n", nil })
	got, err := Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Path != "ok.go" {
		t.Errorf("expected only ok.go, got %+v", got)
	}
}

func TestStatusPropagatesError(t *testing.T) {
	stubRun(t, func(_ ...string) (string, error) { return "", errors.New("boom") })
	if _, err := Status(); err == nil {
		t.Error("expected error from Status()")
	}
}

func TestStageNoFilesIsNoop(t *testing.T) {
	called := false
	stubRun(t, func(_ ...string) (string, error) { called = true; return "", nil })
	if err := Stage(nil); err != nil {
		t.Fatal(err)
	}
	if called {
		t.Error("Stage(nil) must not shell out")
	}
}

func TestStagePassesExplicitFiles(t *testing.T) {
	var gotArgs []string
	stubRun(t, func(args ...string) (string, error) { gotArgs = args; return "", nil })

	if err := Stage([]string{"a.go", "b.go"}); err != nil {
		t.Fatal(err)
	}
	want := []string{"add", "--", "a.go", "b.go"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("Stage args = %v, want %v", gotArgs, want)
	}
}

func TestCommitPassesSubject(t *testing.T) {
	var gotArgs []string
	stubRun(t, func(args ...string) (string, error) { gotArgs = args; return "", nil })

	if err := Commit("feat: add thing."); err != nil {
		t.Fatal(err)
	}
	want := []string{"commit", "-m", "feat: add thing."}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Errorf("Commit args = %v, want %v", gotArgs, want)
	}
}

func TestLogBuildsArgsByScope(t *testing.T) {
	cases := []struct {
		name  string
		scope string
		limit int
		want  []string
	}{
		{"current with limit", "current", 10, []string{"log", "--oneline", "--decorate", "--color=never", "-n", "10"}},
		{"all graphs", "all", 0, []string{"log", "--oneline", "--decorate", "--color=never", "--all", "--graph"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotArgs []string
			stubRun(t, func(args ...string) (string, error) { gotArgs = args; return "log output", nil })

			out, err := Log(tc.scope, tc.limit)
			if err != nil {
				t.Fatal(err)
			}
			if out != "log output" {
				t.Errorf("Log output = %q", out)
			}
			if !reflect.DeepEqual(gotArgs, tc.want) {
				t.Errorf("Log args = %v, want %v", gotArgs, tc.want)
			}
		})
	}
}

func TestItoa(t *testing.T) {
	cases := map[int]string{0: "0", 5: "5", 42: "42", 100: "100"}
	for in, want := range cases {
		if got := itoa(in); got != want {
			t.Errorf("itoa(%d) = %q, want %q", in, got, want)
		}
	}
}
