# ADR-0005 — UI depends on narrow injected interfaces

- **Status:** accepted
- **Date:** 2026-06-25

## Context

The Bubble Tea model needs git, audit and stats. If it calls those packages
directly, its tests must run real git — slow, flaky, and against the TDD goal
(ADR-0002).

## Decision

Lightweight ports & adapters:

- The model receives **narrow interfaces** (only the methods it uses), not whole
  packages.
- `main.go` injects the real adapters (wrappers over `internal/git`, `audit`,
  `stats`).
- Tests inject trivial fakes and assert `Update → state` transitions.

Abstractions follow the **rule of three**: we introduce an interface when a real
need (testability, 3+ repetitions) justifies it — never speculatively.

## Consequences

- The model is testable without a repository.
- `model.go`'s current direct `git.Commit(subject)` call (a scaffold TODO) must
  be replaced by an injected port when the UI is wired.
- Slightly more wiring in `main.go`, paid back in test speed and clarity.
