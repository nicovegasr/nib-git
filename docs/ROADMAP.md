# nib — roadmap & design decisions

Decisions locked in a design session (2026-06-25). This file is the source of
truth so a fresh Claude session (or you, in a month) doesn't re-debate them.

## Why nib exists

A commit flow you keep open in the VSCode terminal and use 50× a day, faster
than the Source Control panel or lazygit/gitui (both tried, neither fit). It is
**not** a full git client — it's a fast committer + a few inspection commands.

## Decisions (locked)

| Topic | Decision |
|---|---|
| Repo | Separate from `conventional-stats`. No submodule. Independent. |
| Language | Go + Bubble Tea (TUI). |
| Versioning | SemVer, starts `0.1.0`. |
| Minimal flow | One line: type dropdown + message, over the file list. |
| Interaction | Mouse-first; keyboard always available. |
| Mouse caveat | Mouse capture hijacks terminal selection/scroll → toggle in `/settings`. |
| Staging | `git add` takes an explicit file list, never `git add .`. |
| Heavy features | Slash-commands: `/audit`, `/stats`, `/settings`. Not panels in the hot path. |
| Philosophy | Rich by default; everything heavy can be switched off. |
| `/audit` metric | Sort by **commit count** per file (Tornhill hotspots). Lines = context only. |
| Commit types | Duplicated from conventional-stats. No sync. `tests` → `test:`. |

## Phases

- [x] **Phase 0 — scaffold.** Structure, commit/git/audit/stats/settings
  packages, skeleton UI, mockups, CLAUDE.md. (current)
- [ ] **Phase 1 — minimal commit flow (the core).** Wire `internal/ui` to a
  real render of the file list + inline commit line. Keyboard path first:
  `j/k` move, `space` stage, `tab` to message, type shortcuts, `enter` commit.
  Ship a working `nib` that commits. *This alone is the MVP.*
- [ ] **Phase 2 — type dropdown + mouse.** Open/close dropdown, select type;
  mouse clicks on files and dropdown. `tea.MouseMsg` handling.
- [ ] **Phase 3 — git-log panel.** Optional right-side log, `LogScope`
  current/all/custom. Off-by-toggle.
- [ ] **Phase 4 — `/settings`.** Command palette + persisted TOML
  (`~/.config/nib/settings.toml`), all toggles live.
- [ ] **Phase 5 — `/audit`.** Render hotspots, `--days`/`--exclude`, threshold
  coloring from settings.
- [ ] **Phase 6 — `/stats`.** Commit-type breakdown view; later timeline/authors.
- [ ] **Phase 7 — release.** `gh repo create nib-git`, push, tag `v0.1.0`,
  build/install instructions, maybe `brew`/`go install`.

## Working style

Phase by phase, a commit at the end of each phase (visible progress in git
history). Conventional Commits in English, trailing period in the subject.

## Open questions

- `/stats` timeline granularity (per day / per week)?
- Distribution: `go install`, Homebrew tap, or prebuilt binaries?
- Does `/audit` v0.2 (commits × file size) need file LOC, or git blame lines?
