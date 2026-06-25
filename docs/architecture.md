# Architecture

A human-readable map of nib-git. For the *why* behind decisions, see
[`adr/`](adr/). For build/test commands, see [`contributing.md`](contributing.md).

## Big picture

nib is a terminal-first conventional-commit tool: a Bubble Tea TUI over a small
set of typed domain packages that shell out to `git`.

```
main.go                 entrypoint; builds the TUI and injects adapters
internal/commit         conventional-commit types + Format()
internal/git            git CLI wrappers: Status, Stage, Commit, Log
internal/audit          /audit hotspots (by commit frequency) + .auditignore
internal/stats          /stats commits-by-type breakdown
internal/settings       toggleable options (Default() = everything on)
internal/ui/model.go    Bubble Tea model (wiring in progress)
```

## Data flow

```
user (mouse/keys)
      │
      ▼
internal/ui (Bubble Tea model)         ← state machine: Update → state
      │   depends on NARROW interfaces (ADR-0005)
      ▼
internal/{git,audit,stats,commit}      ← domain logic, pure + testable
      │   thin shell-out seam (var run = func ...)
      ▼
git CLI
```

The domain packages split **pure logic** (parsing, ranking, formatting — fully
tested) from a **thin shell-out seam** (`var run`/`var gitLog…`, stubbed in
tests). This is what makes the ≥80% domain coverage real (ADR-0002).

## Non-negotiable rules

- **Minimal flow is one line:** mark files → pick type → type message → enter.
  Never make it pay for a panel that `/settings` can disable.
- **Mouse-first**, keyboard always available.
- `git add` takes an **explicit file list**, never `git add .`.
- `/audit` sorts by **number of commits** touching a file (Tornhill hotspots);
  lines changed are context, never the sort key.
- Commit types are **duplicated on purpose** from conventional-stats; no sync
  mechanism. Update `internal/commit/types.go` by hand.

## Where to start reading

`internal/commit/types.go` (smallest, fully tested) → `internal/git/git.go`
(the shell-out seam pattern) → `internal/ui/model.go` (the work in progress).
