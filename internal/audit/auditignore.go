package audit

import (
	"os"
	"path/filepath"
	"strings"
)

// defaultIgnore are patterns excluded from /audit out of the box: obvious
// high-churn-by-design files that are not code smells.
var defaultIgnore = []string{
	"*.lock",
	"package-lock.json",
	"pnpm-lock.yaml",
	"yarn.lock",
	"dist/",
	"build/",
	"node_modules/",
}

// loadAuditIgnore returns the default patterns plus any from a local
// .auditignore file (one glob per line, # comments allowed).
func loadAuditIgnore() []string {
	patterns := append([]string{}, defaultIgnore...)
	data, err := os.ReadFile(".auditignore")
	if err != nil {
		return patterns
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns
}

// isIgnored reports whether path matches any ignore pattern. Directory
// patterns end in "/" and match any path under them.
func isIgnored(path string, patterns []string) bool {
	for _, p := range patterns {
		if strings.HasSuffix(p, "/") {
			if strings.HasPrefix(path, p) || strings.Contains(path, "/"+p) {
				return true
			}
			continue
		}
		if ok, _ := filepath.Match(p, filepath.Base(path)); ok {
			return true
		}
	}
	return false
}
