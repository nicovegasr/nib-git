// Package git wraps the git CLI for the operations nib needs. It shells out to
// `git` rather than using a library to stay small and fast.
package git

import (
	"os/exec"
	"strings"
)

// FileChange is a single entry from `git status`.
type FileChange struct {
	Path   string
	Status rune // M, A, D, R, ? (porcelain XY collapsed to a single hint)
	Staged bool
}

// Status returns the working-tree changes via `git status --porcelain`.
func Status() ([]FileChange, error) {
	out, err := run("status", "--porcelain")
	if err != nil {
		return nil, err
	}
	var changes []FileChange
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if len(line) < 4 {
			continue
		}
		x, y := line[0], line[1]
		path := strings.TrimSpace(line[3:])
		st := rune(x)
		if st == ' ' {
			st = rune(y)
		}
		changes = append(changes, FileChange{
			Path:   path,
			Status: st,
			Staged: x != ' ' && x != '?',
		})
	}
	return changes, nil
}

// Stage stages exactly the given files (replaces the old `git add .`).
func Stage(files []string) error {
	if len(files) == 0 {
		return nil
	}
	_, err := run(append([]string{"add", "--"}, files...)...)
	return err
}

// Commit creates a commit with the given full subject line.
func Commit(subject string) error {
	_, err := run("commit", "-m", subject)
	return err
}

// Log returns a pretty one-line log. scope is "current" or "all".
func Log(scope string, limit int) (string, error) {
	args := []string{"log", "--oneline", "--decorate", "--color=never"}
	if scope == "all" {
		args = append(args, "--all", "--graph")
	}
	if limit > 0 {
		args = append(args, "-n", itoa(limit))
	}
	return run(args...)
}

// run executes git with the given args. It is a package var so tests can stub
// the shell-out and exercise the parsing without a real repository.
var run = func(args ...string) (string, error) {
	out, err := exec.Command("git", args...).Output()
	return string(out), err
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
