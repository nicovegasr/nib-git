// Package audit implements /audit: change-frequency hotspots.
//
// Primary signal is the NUMBER OF COMMITS that touch a file (Adam Tornhill's
// hotspot method, "Your Code as a Crime Scene"). Lines changed are noisy
// (renames, reformats, lockfiles) so they are only context, never the sort
// key. v0.2 idea: commits × file size = the real hotspot.
package audit

import (
	"os/exec"
	"sort"
	"strings"
)

// Hotspot is one file ranked by change frequency.
type Hotspot struct {
	Path    string
	Commits int // primary signal
	// TODO(lines): +/- churn as a secondary context column.
}

// Run returns files ordered by commit frequency over the last `days` days.
// `exclude` are extra inline excludes on top of .auditignore.
func Run(days int, exclude []string) ([]Hotspot, error) {
	ignore := loadAuditIgnore()
	ignore = append(ignore, exclude...)

	out, err := exec.Command("git", "log",
		"--since", itoaDays(days), "--name-only", "--pretty=format:").Output()
	if err != nil {
		return nil, err
	}

	counts := map[string]int{}
	for _, line := range strings.Split(string(out), "\n") {
		path := strings.TrimSpace(line)
		if path == "" || isIgnored(path, ignore) {
			continue
		}
		counts[path]++
	}

	hotspots := make([]Hotspot, 0, len(counts))
	for p, c := range counts {
		hotspots = append(hotspots, Hotspot{Path: p, Commits: c})
	}
	sort.Slice(hotspots, func(i, j int) bool {
		return hotspots[i].Commits > hotspots[j].Commits
	})
	return hotspots, nil
}

func itoaDays(d int) string {
	if d <= 0 {
		d = 30
	}
	return itoa(d) + " days ago" // git accepts "30 days ago"
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
