# Domain Model & Tech Stack — gowrite (Appwrite → Go migration)

> Produced via the `model-domain` grilling session (2026-08-11), grounded in the extracted reversa specs.
> Invariants and state machines are the hard gate — here they reference the already-extracted authoritative sources.

## Target stack (revised post red-team 2026-08-11)
- **Language**: Go (transformational paradigm — see `specs/migration/paradigm_decision.md`).
- **Database**: PostgreSQL only (relational + JSONB + pgvector). Consolidates legacy Mongo+MariaDB+Postgres **and Redis** — jobs via `jobs` table + `SKIP LOCKED`, realtime via `LISTEN/NOTIFY`, rate-limit via in-memory token bucket (red-team O3/O5).
- **Data access**: `sqlc` (typed SQL) + explicit `crypto.Encrypt()` — no reimplemented query engine (red-team O4).
- **Topology**: Hybrid — 11 bounded contexts + Console BFF; **flat idiomatic Go packages** (no hexagonal ports/adapters — red-team O1); `platform/` = leaf packages only (O2).
- **Strategy**: Strangler Fig + Parallel Run, flip-unit = context, topological order (`specs/migration/migration_strategy.md`).

## Bounded contexts (revised 2026-08-11 after red-team — 11 contexts + BFF)

> The initial 9-context map was red-teamed (`specs/migration/redteam_findings.md`). Fixes applied: promoted **Migrations** and **Advisor** to real contexts (they have aggregates + state machines — B5/B6); moved **Webhooks** into Projects (it's a `projects` aggregate child — B4); reclassified **Console** as a BFF above the contexts, not inside Projects (B3). Only genuinely stateless modules remain in `platform/`.

| Context | Legacy modules folded in | Responsibility | Core invariant(s) | Flip unit |
|---------|--------------------------|----------------|-------------------|-----------|
| **Identity** | Account, Users | end-user auth, sessions, MFA, identities, targets | email unique across `users`+`identities`; verified TOTP can't be recreated (409) | **co-flips with Organizations** (S2) |
| **Organizations** | Teams, Organization | teams, memberships | ≥1 owner (`MEMBERSHIP_DOWNGRADE_PROHIBITED`); status derived from timestamps | **co-flips with Identity** (S2) |
| **Projects** | Project, Projects, **Webhooks** | project lifecycle, keys, platforms, OAuth2, onboarding, schedules, webhooks | project = tenant root; **Projects is the sole owner of the `projects` write-model** (B1) | own |
| **Data** | Databases, VectorsDB | collections, attributes, indexes, transactions, embeddings | optimistic-lock transactions (self-contained — S6) | own |
| **Storage** | Storage, Tokens | buckets, files, resource tokens | chunked-upload lock; resource token = `bucketId:fileId` (B8 ✅) | own |
| **Compute** | Functions, Sites, Compute | deployments, executions, builds | deployment/execution SMs; build dedup; `ready` before execute | **after/with VCS** (S4) |
| **VCS** | VCS | installations, repos, webhook deploys | webhook signature (403); HMAC state | **before/with Compute** (S4); publishes `repository.pushed` via outbox — **no direct writes into Compute** (B2) |
| **Networking** | Proxy | rules, DNS, TLS | rule SM; CAA gate | own (read-port into Compute readiness — B7) |
| **Messaging** | Notifications (+ messaging worker) | messages, targets, providers, topics | message SM; terminal idempotent | own |
| **Migrations** | Migrations | provider imports (Firebase/Supabase/…), CSV/JSON | `pending→processing→completed\|failed`; retry only from failed | own (touches Data on import) |
| **Advisor** | Advisor | reports → insights | Report has-many Insight; Insight `active→dismissed` + cascade | own |

### Console = BFF (edge layer above the contexts) — not a bounded context (B3)
`Console` (Redirects, Scopes, Assistant, Init guards, Email templates, Variables) composes reads/writes across Identity + Organizations + Projects. Modeled as a **backend-for-frontend** that delegates into the contexts, never owning their aggregates.

### `platform/` (cross-cutting leaf packages — NOT a bounded context, NOT a god-package — O2)
Small independent leaf packages: **auth, permissions** (`app/config/roles.php` 6-role enum), **crypto** (25 encrypted fields), **errors** (`AppError`), **ratelimit** (route labels). DB/queue access is direct `pgx`/stdlib usage inside each context, not an abstraction layer. Genuinely stateless services: **Health, Avatars, Presences** (no aggregate, no state machine).

## Invariants (hard gate — authoritative sources)
- **State machines**: `specs/state-machines.md` — Deployment, Execution, Rule/Proxy, Migration, Membership, Insight, Webhook, **Message** (8 machines, all with legal transitions + terminal states).
- **Domain rules**: `specs/domain.md` — consolidated auth/teams/data/storage/compute invariants.
- **Permissions**: `specs/permissions.md` — role enum, scopes, ACL, rate limits.
- **Data invariants**: `specs/database/` — 48 collections, 25 encrypted fields, 13 unique constraints.

## Concurrency (paradigm note)
Legacy in-process concurrency = Swoole coroutines; target = goroutines + channels (CSP). The mandatory-terminal-state rule for queue jobs (ADR-008) MUST be preserved: every worker message ends in a terminal state; idempotency keyed by CloudEvent id. Shared mutable state (rate-limit windows, presence cache) moves to Postgres/Redis with explicit synchronization — no in-process globals.

## Domain events (design note for Designer)
Target is Go-native (CSP), not full event-sourcing. Recommendation: formalize an **outbox** for cross-context side effects (webhooks, function triggers, realtime) to preserve the legacy event-bus semantics with at-least-once + idempotency, rather than implicit queue publishes. Deferred to `target_domain_model.md`.
