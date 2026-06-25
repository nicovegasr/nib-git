# ADR-0006 — Three test tiers

- **Status:** accepted
- **Date:** 2026-06-25

## Context

Unit tests stub the git shell-out (ADR-0002) for speed and hermeticity. That
leaves the real `git` invocations unexercised. We want confidence on the happy
paths end-to-end without weakening the unit discipline.

## Decision

Three tiers:

1. **Unit** — logic with git stubbed. `make test`. Fast, hermetic. Counts toward
   the ≥80% domain coverage gate.
2. **Integration** — real `git` in a throwaway repo, happy paths. Behind the
   `integration` build tag. `make test-integration`. **Excluded** from the
   coverage gate so the shell-out seam stays "uncovered" by unit metrics and the
   discipline holds.
3. **UI e2e** — drive the Bubble Tea model with `teatest` against a real repo.
   Same `integration` tag. Added when the model is wired.

CI runs all three (`make ci`).

## Consequences

- Real git behaviour is covered without inflating unit coverage.
- The fast loop (`make test`) stays hermetic and quick.
- Integration tests require `git` on the runner (always present in CI).
