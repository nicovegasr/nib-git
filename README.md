# nib

Fast, terminal-first commit tool for conventional commits — a lighter "git IDE"
you keep open in your VSCode terminal to commit on the fly.

> Sibling of [`conventional-stats`](https://github.com/nicovegasr/conventional-stats)
> (zsh toolkit). Independent — no submodule; the commit-type list is duplicated
> by design.

`v0.1.0` — scaffold. Domain packages are in place and tested; the Bubble Tea UI
is being wired.

## Quick start

```sh
go mod tidy
make build
./nib            # mark files → pick type → type message → enter
```

Mouse-first, keyboard always available. Heavy features are slash-commands
(`/audit`, `/stats`, `/settings`) and everything heavy can be switched off.

## Docs

- [Architecture](docs/architecture.md) — package map, data flow, the rules.
- [Contributing](docs/contributing.md) — `make` targets, TDD, code style.
- [Decisions](docs/adr/) — ADRs (the *why*).
- [Roadmap](docs/ROADMAP.md) · [UX mockups](docs/ux-mockups.html)
