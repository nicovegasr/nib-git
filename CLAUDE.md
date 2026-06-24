# nib-git — Claude Code project guide

`nib` is a fast, **terminal-first** commit tool for conventional commits — a
lighter "git IDE" you keep open in the VSCode terminal to commit on the fly.
Sibling of `conventional-stats` (zsh toolkit), **independent — no submodule**.

Read `docs/ROADMAP.md` for the design decisions and phase plan, and
`docs/ux-mockups.html` for the target UX (all 6 states).

## Build & run

Go is the runtime (Bubble Tea TUI). If `go` is missing: `brew install go`.

```sh
go mod tidy            # populate go.sum + transitive deps (needs network)
go build -o nib .
./nib                  # opens the commit flow
go test ./...          # run tests
go vet ./...           # static checks
gofmt -l .             # list unformatted files (must be empty)
```

## Architecture

```
main.go               entrypoint (tea.NewProgram, AltScreen + mouse)
internal/commit       conventional-commit types + Format() (auto trailing period)
internal/git          git CLI wrappers: Status, Stage([]files), Commit, Log
internal/audit        /audit hotspots (by commit frequency) + .auditignore
internal/stats        /stats commits-by-type breakdown
internal/settings     toggleable options (Default() = everything on)
internal/ui/model.go  Bubble Tea model (currently a SKELETON with TODOs)
```

The UI layer (`internal/ui`) is the unfinished part. Everything it needs from
git/audit/stats/commit already has a typed function — wire the model to those,
do not re-shell-out from the UI.

## Non-negotiable design rules

- **Minimal flow is one line**: type dropdown + message field, over the file
  list. Mark files → pick type → type message → enter. This must stay fast.
- **Mouse-first**, keyboard always available as an alternative.
- **"Rich by default, everything heavy switches off in `/settings`."** Never
  make the minimal flow pay for a panel the user can't disable.
- **`git add` takes an explicit file list**, never `git add .`.
- **`/audit` sorts by NUMBER OF COMMITS** touching a file (Tornhill hotspots).
  Lines changed are context only — never the sort key.
- **Commit types are duplicated** from conventional-stats on purpose. Do not
  build a sync mechanism. If they drift, update `internal/commit/types.go` by
  hand. The `tests` label maps to the `test:` prefix.
- **SemVer**, currently `0.1.0` (see `VERSION`).

## Conventions

- Commits: Conventional Commits in English (`feat:`, `fix:`, `chore:` …),
  trailing period in the subject (matches the tool's own `Format`).
- Prose to the user in Spanish; code/identifiers/commits in English.
- Keep packages small; shell out to `git` rather than pulling a git library.

## Slash-commands (planned)

`/audit` · `/stats` · `/settings` — a command palette inside the running TUI
(Claude-Code-terminal style). See ROADMAP phases.
