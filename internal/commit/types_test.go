package commit

import "testing"

func TestFormatAppendsPeriod(t *testing.T) {
	got := Format("feat", "add inline commit line")
	want := "feat: add inline commit line."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatKeepsExistingPeriod(t *testing.T) {
	if got := Format("fix", "crash on empty repo."); got != "fix: crash on empty repo." {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestTestsTypeUsesTestPrefix(t *testing.T) {
	for _, ty := range Types {
		if ty.Name == "tests" && ty.Prefix != "test" {
			t.Fatalf("tests label must map to 'test:' prefix, got %q", ty.Prefix)
		}
	}
}
