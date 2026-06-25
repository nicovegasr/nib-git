# ADR-0002 — Pragmatic TDD with a per-package coverage gate

- **Status:** accepted
- **Date:** 2026-06-25

## Context

We want TDD without the absurdity of writing failing tests for terminal
rendering and mouse plumbing. Coverage must be honest, not gamed.

## Decision

- **Pragmatic TDD:** any function with logic (parsing, hotspot ranking,
  `commit.Format`, model `Update → state` transitions) is written test-first.
  Raw rendering and the thin git shell-out are not.
- **Coverage gate, per package:** domain packages (`commit`, `audit`, `stats`,
  `git`, `settings`) must each stay at **≥ 80%**. `internal/ui` requires **0%**.
- Measured per package (never globally) so a well-covered package cannot mask a
  poorly-covered one. Implemented in `scripts/coverage.sh`.

## Consequences

- Domain shell-outs are made injectable (a `var run = func(...)` seam) so
  parsing is tested without a real repository. This nudged the design toward
  testable seams — see ADR-0005.
- The UI is tested by behaviour (state transitions), not by snapshotting output.
- The coverage gate runs in CI / pre-push, not pre-commit, to keep commits fast.
