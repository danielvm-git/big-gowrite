# gowrite — Go implementation conventions

Idiomatic Go implementation of the Appwrite backend. Specs live one level up in `../specs/`.

## Paradigm (from specs/migration/paradigm_decision.md)

- **Explicit errors**, not exceptions. Return `(T, error)`. Wrap with `%w`. No panics for control flow.
- **Composition over inheritance.** Concrete structs + small interfaces at the edge only.
- **CSP concurrency.** goroutines + channels; `errgroup` for fan-out. No shared-memory locks unless a resource genuinely needs one.
- **Statically typed schemas.** Go structs + `sqlc`-generated queries. No runtime-array schema, no generic ORM engine.

## Structure (from specs/migration/target_architecture.md — flat, no hexagonal)

```
cmd/{server,worker}/      # binaries
platform/                 # leaf packages, minimal deps
  errors/ auth/ permissions/ crypto/ ratelimit/
internal/                 # bounded contexts (flat: service.go, entity.go, queries.sql, *_test.go)
  jobs/ identity/ organizations/ projects/ ...
db/                       # migrations/*.sql + sqlc queries + generated code
```

- **No `ports/`, no `adapters/`, no DI container.** Wire dependencies explicitly in `cmd/server/main.go`.
- `platform/` packages are leaves — they must not import `internal/`.

## Data access

- **PostgreSQL only** (relational + JSONB + pgvector). No Redis, no Mongo, no MariaDB.
- Data access via `sqlc` + hand-written SQL. Encryption at write sites via `platform/crypto` (never an implicit filter).
- Async dispatch: `jobs` table + `SELECT … FOR UPDATE SKIP LOCKED`. Realtime: `LISTEN/NOTIFY`.
- Every async job drives to a **terminal state** (ADR-008); idempotency via `dedup_key`.

## Testing (develop-tdd)

- Table-driven tests. RED → GREEN → refactor.
- Integration tests use an **ephemeral Postgres** (env `GOWRITE_TEST_DATABASE_URL`); skip cleanly when unset.
- Verify commands in `../specs/epics/*/‹story›-tasks.yaml` must pass before a task flips to `passing`.

## Errors

- One error type: `platform/errors.AppError{Code, HTTPStatus, Err}`.
- Middleware translates `AppError` → the legacy JSON envelope `{type, code, message}`. Unknown errors → 500, no internal leak.

## Naming

- Packages: short, lowercase, singular (`crypto`, `jobs`). Exported identifiers documented with a leading comment.
- Match the legacy REST surface for HTTP paths and response field names (camelCase; `$id`/`$createdAt` system fields preserved).

## Security (non-negotiable)

- 25 encrypted-at-rest fields via AES-GCM (`platform/crypto`). Never log plaintext secrets.
- All access flows through `platform/auth` (Principal) + `platform/permissions` (scopes + document ACL).
