# Changelog

All notable changes to nib-git are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/); versioning is
[SemVer](https://semver.org/). (Auto-derivation from conventional commits is a
later goal — ADR-0004.)

## [Unreleased]

### Added

- Minimal commit flow (Phase 1): file list + inline commit line wired to git
  through an injected `Repo` port (ADR-0005). Keyboard path (`j/k`, `space`,
  type shortcuts, `tab`, `enter`), commit result reported, transitions tested.

- Quality harness: `Makefile` gate, `golangci-lint` config, `lefthook` git
  hooks, GitHub Actions CI, and a Claude Code `gofmt` hook (ADR-0001, 0003).
- Per-package coverage gate at 80% for domain packages (`scripts/coverage.sh`,
  ADR-0002).
- Test suites for `commit`, `git`, `audit`, `stats`, `settings`; injectable
  shell-out seam in the domain packages.
- Integration test tier: real git on happy paths behind the `integration`
  build tag (`make test-integration`, ADR-0006).
- Docs: `docs/architecture.md`, `docs/contributing.md`, `docs/adr/`.

## [0.1.0] — 2026-06-25

### Added

- Project scaffold: Go + Bubble Tea TUI skeleton, domain packages
  (`commit`, `git`, `audit`, `stats`, `settings`), `CLAUDE.md`, `ROADMAP`.
