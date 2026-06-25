# ADR-0001 — Quality harness in three layers

- **Status:** accepted
- **Date:** 2026-06-25

## Context

nib-git must stay a professional, maintainable codebase — not "vibe coding".
Written rules alone rely on discipline; we want guarantees that survive any
contributor (human or agent).

## Decision

Enforce quality in three complementary layers:

1. **`CLAUDE.md`** — the source of truth for rules and conventions.
2. **Automated hooks** — `lefthook` git hooks run the gate on every commit; a
   Claude Code `PostToolUse` hook runs `gofmt` after each edit.
3. **Project skills** — only where a repeatable flow truly exists (e.g. a future
   `/docs`). We do not add skills speculatively.

## Consequences

- The gate runs regardless of who commits; rules are not just aspirational.
- One Makefile is the single source of truth shared by hooks and CI (ADR-0003).
- Slight setup cost: contributors run `make tools` once.
