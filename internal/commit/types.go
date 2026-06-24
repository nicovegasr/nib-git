// Package commit defines the conventional-commit types and message formatting.
//
// This list is intentionally DUPLICATED from the conventional-stats zsh
// toolkit (config/git-commits.zsh). The two repos stay independent; there are
// ~13 near-static types, so duplicating them is cheaper than building a sync
// mechanism. Keep this in sync by hand if the zsh list changes.
package commit

import "strings"

// Type is a selectable commit type shown in the dropdown.
type Type struct {
	Name     string // label shown in the dropdown, e.g. "feat"
	Prefix   string // prefix written to the commit, e.g. "feat" -> "feat: ..."
	Desc     string // one-line description
	Shortcut rune   // keyboard shortcut in the dropdown
}

// Types is the ordered list rendered in the type dropdown.
// Mirror of config/git-commits.zsh in conventional-stats.
var Types = []Type{
	{"feat", "feat", "Nueva funcionalidad", 'f'},
	{"fix", "fix", "Corrección de bug", 'x'},
	{"red", "red", "TDD · test roto / WIP", 'r'},
	{"green", "green", "TDD · tests pasan / estado funcional", 'g'},
	{"refactor", "refactor", "Refactor sin cambio de comportamiento", 'e'},
	{"hotfix", "hotfix", "Corrección urgente en producción", 'h'},
	{"docs", "docs", "Documentación", 'd'},
	{"style", "style", "Formato, sin cambio de lógica", 's'},
	{"tests", "test", "Tests añadidos o corregidos", 't'}, // label "tests", prefix "test:"
	{"chore", "chore", "Mantenimiento, dependencias, tooling", 'c'},
	{"perf", "perf", "Mejora de rendimiento", 'p'},
	{"ci", "ci", "Configuración de CI/CD", 'i'},
	{"build", "build", "Sistema de build", 'b'},
}

// Format builds the final commit subject: "<prefix>: <message>" with a
// trailing period auto-appended if missing (mirrors _execute_commit in zsh).
func Format(prefix, message string) string {
	message = strings.TrimSpace(message)
	if message != "" && !strings.HasSuffix(message, ".") {
		message += "."
	}
	return prefix + ": " + message
}
