# ADR-0003 — Makefile as the single quality gate

- **Status:** accepted
- **Date:** 2026-06-25

## Context

Hooks and CI must run *exactly* the same checks. Duplicating commands across a
git hook, a CI YAML and a contributor's terminal guarantees drift.

## Decision

A `Makefile` is the single source of truth for developer tasks:

- `make check` — fast gate: `fmt-check` + `vet` + `lint` + `test`.
- `make cover` — per-package coverage gate (ADR-0002).
- `make ci`    — `check` + `cover`, run in CI.

`lefthook` pre-commit calls `make check`; pre-push calls `make cover`; CI calls
`make ci`. Static analysis is `golangci-lint` with the "recommended" linter set
in `.golangci.yml`.

## Consequences

- Identical behaviour locally and in CI; no drift.
- `misspell` is intentionally excluded: user-facing strings are Spanish by
  design and misspell only understands English.
- New checks are added in one place.
