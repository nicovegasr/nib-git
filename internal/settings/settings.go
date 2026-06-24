// Package settings holds user-toggleable options. Philosophy: "rich by
// default, everything heavy can be switched off". These map 1:1 to the
// /settings mock in docs/ux-mockups.html.
package settings

// Settings is the persisted configuration. TODO: load/save from
// ~/.config/nib/settings.toml.
type Settings struct {
	ShowLog      bool   // git-log panel visible (off => leaner/faster)
	MouseCapture bool   // click on files; off restores native terminal selection/scroll
	LogScope     string // "current" | "all" | "custom"
	LogCustomCmd string // used when LogScope == "custom"

	AuditDays      int // default window for /audit
	AuditThreshold int // commit count that flags a file as "refactor?"
}

// Default returns the out-of-the-box configuration: rich, with everything on.
func Default() Settings {
	return Settings{
		ShowLog:        true,
		MouseCapture:   true,
		LogScope:       "current",
		AuditDays:      30,
		AuditThreshold: 20,
	}
}

// TODO(load): Load() (Settings, error) reading the TOML file, falling back to
// Default() when absent.
// TODO(save): Save(Settings) error writing the TOML file.
