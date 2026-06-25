# ADR-0004 — Branch + PR workflow and manual changelog

- **Status:** accepted
- **Date:** 2026-06-25

## Context

`main` should always be green and reviewable, even on a solo project. Versioning
follows SemVer (`VERSION` = 0.1.0).

## Decision

- **Always branch + PR.** No direct commits to `main`. CI (`make ci`) is the
  gate before merge.
- **Changelog:** a manual `CHANGELOG.md` (Keep a Changelog) for now. Later,
  derive it automatically from the conventional commits — a natural fit, since
  nib *is* a conventional-commit tool.

## Consequences

- `main` stays releasable; history is reviewable.
- Slightly more ceremony on a solo repo, accepted for the discipline.
