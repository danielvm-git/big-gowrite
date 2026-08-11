# Gowrite — Appwrite Migration Specs

This is the **gowrite** specs repository: reverse-engineering documentation and migration planning
for rewriting [Appwrite](https://github.com/appwrite/appwrite) in Go.

## Reference clones (gitignored)

| Directory | Purpose |
|---|---|
| `appwrite/` | Appwrite server source (PHP). Pull to update. |
| `appwrite-console/` | Appwrite console source (SvelteKit). Pull to update. |
| `appwrite-reversa/` | Stale reversa worktree; kept as skill reference. |

## Reversa skills

All reversa pipeline skills live in `.claude/skills/`. Use the flows below from this root:

- `/reversa` — discover and document an existing system
- `/reversa-migrate` — plan the migration of a legacy system
- `/reversa-docs` — generate the visual docs mini-site
- `/reversa-forward` — implement or evolve code from specs
- `/reversa-agents-help` — full agent catalogue

Skills path: `.claude/skills/reversa/SKILL.md` (and siblings).

## Specs layout

```
specs/
├── migration/          ← 6-agent migration pipeline output
│   ├── handoff.md      ← authoritative anchor decisions
│   ├── parity_tests/   ← Gherkin parity scenarios
│   └── ...
├── database/           ← ERD, data dictionary, relationships
├── tech-architecture/  ← target stack decisions
├── adrs/               ← Architecture Decision Records
└── <Module>/           ← per-module specs (22 modules)
```

## Rule

Never modify or delete files in `appwrite/`, `appwrite-console/`, or `appwrite-reversa/`.
The reversa pipeline writes only to `.reversa/`, `specs/`, `_reversa_docs/`, and `_reversa_forward/`.
