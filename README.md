# nib

Fast, terminal-first commit tool for conventional commits. A lighter "git IDE"
you keep open in your VSCode terminal to commit on the fly — much faster than
the Source Control panel or lazygit.

> Sibling project of [`conventional-stats`](https://github.com/nicovegasr/conventional-stats)
> (the zsh terminal toolkit). The two are **independent** — no submodule. The
> commit-type list is duplicated by design.

## Status

`v0.1.0` — **scaffold**. Structure, commit-type model and git/audit/stats
parsers are in place; the Bubble Tea UI is a skeleton. See
[`docs/ux-mockups.html`](docs/ux-mockups.html) for the target UX.

## Design

- **Minimal flow** (irreducible): one line — type dropdown + message — over the
  changed-files list. Mark files, pick type, type message, `enter`.
- **Mouse-first**, keyboard available. Mouse capture is toggleable (it hijacks
  native terminal selection/scroll).
- **Rich by default, everything heavy can be switched off** in `/settings`.
- **Slash-commands** for the heavy features:
  - `/audit` — change-frequency hotspots (refactor candidates). Sorted by
    **number of commits** touching a file; lines are context only.
  - `/stats` — commits by conventional-commit type over a window.
  - `/settings` — toggles (git-log panel, mouse capture, audit window/threshold,
    `.auditignore`).

## Build

```sh
brew install go        # not yet installed on this machine
go mod tidy            # populate go.sum + transitive deps
go build -o nib .
./nib
```

## Layout

```
main.go               entrypoint (Bubble Tea program)
internal/commit       conventional-commit types + message formatting
internal/git          git CLI wrappers (status, stage, commit, log)
internal/audit        /audit hotspots + .auditignore
internal/stats        /stats commit-type breakdown
internal/settings     toggleable options
internal/ui           Bubble Tea model (skeleton)
docs/ux-mockups.html  UX reference (all states)
```
