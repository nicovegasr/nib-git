package audit

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRankSortsByCommitFrequency(t *testing.T) {
	out := "a.go\nb.go\na.go\na.go\nb.go\nc.go\n"
	got := rank(out, nil)

	if len(got) != 3 {
		t.Fatalf("expected 3 hotspots, got %+v", got)
	}
	if got[0].Path != "a.go" || got[0].Commits != 3 {
		t.Errorf("top hotspot = %+v, want a.go:3", got[0])
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].Commits < got[i].Commits {
			t.Errorf("not sorted by frequency: %+v", got)
		}
	}
}

func TestRankAppliesIgnore(t *testing.T) {
	out := "src/app.go\nnode_modules/x.js\nyarn.lock\nsrc/app.go\n"
	got := rank(out, defaultIgnore)

	if len(got) != 1 || got[0].Path != "src/app.go" {
		t.Errorf("ignore not applied, got %+v", got)
	}
}

func TestIsIgnored(t *testing.T) {
	patterns := []string{"*.lock", "node_modules/", "dist/"}
	cases := map[string]bool{
		"yarn.lock":             true,
		"deep/path/pnpm.lock":   true,
		"node_modules/react.js": true,
		"a/node_modules/x.js":   true,
		"dist/bundle.js":        true,
		"src/app.go":            false,
		"distinct/keep.go":      false, // "dist/" must not match "distinct/"
	}
	for path, want := range cases {
		if got := isIgnored(path, patterns); got != want {
			t.Errorf("isIgnored(%q) = %v, want %v", path, got, want)
		}
	}
}

func TestLoadAuditIgnoreReadsFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".auditignore"),
		[]byte("# comment\n\nvendor/\n*.gen.go\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	chdir(t, dir)

	got := loadAuditIgnore()
	if !contains(got, "vendor/") || !contains(got, "*.gen.go") {
		t.Errorf("custom patterns missing: %v", got)
	}
	if !contains(got, "node_modules/") {
		t.Errorf("default patterns missing: %v", got)
	}
}

func TestLoadAuditIgnoreFallsBackToDefaults(t *testing.T) {
	chdir(t, t.TempDir()) // no .auditignore present
	got := loadAuditIgnore()
	if len(got) != len(defaultIgnore) {
		t.Errorf("expected only defaults, got %v", got)
	}
}

func TestItoaDays(t *testing.T) {
	cases := map[int]string{0: "30 days ago", -1: "30 days ago", 14: "14 days ago"}
	for in, want := range cases {
		if got := itoaDays(in); got != want {
			t.Errorf("itoaDays(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestItoa(t *testing.T) {
	cases := map[int]string{0: "0", 7: "7", 30: "30"}
	for in, want := range cases {
		if got := itoa(in); got != want {
			t.Errorf("itoa(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestRunUsesStubbedLog(t *testing.T) {
	orig := gitLogNames
	t.Cleanup(func() { gitLogNames = orig })

	gitLogNames = func(_ int) (string, error) { return "x.go\nx.go\n", nil }
	chdir(t, t.TempDir())
	got, err := Run(30, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Commits != 2 {
		t.Errorf("Run = %+v, want x.go:2", got)
	}

	gitLogNames = func(_ int) (string, error) { return "", errors.New("boom") }
	if _, err := Run(30, nil); err == nil {
		t.Error("expected error to propagate")
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
