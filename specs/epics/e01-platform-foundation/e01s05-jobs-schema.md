# e01s05 — Base Postgres schema + jobs table + worker skeleton

## 1. Metadata
- **Story ID:** e01s05
- **Epic:** e01 — Platform Foundation
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** medium

## 2. User Story
As **a context needing async work or realtime**, I want to **enqueue a job and have a worker consume it exactly once**, so that **background work and realtime events run on Postgres alone, with a terminal state guarantee (BR-MIGRAR-014)**.

## 3. Business Value
Establishes the single-datastore dispatch backbone: jobs table with SKIP-LOCKED consumption and LISTEN/NOTIFY realtime — the substitute for Redis + queues + outbox.

## 4. Context & Background
Part of epic **e01 — Platform Foundation** (Shared leaf packages (errors, auth, permissions, crypto, ratelimit), base Postgres schema, and the jobs table with SKIP-LOCKED dispatch + LISTEN/NOTIFY realtime.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- base schema migrations
- jobs table (id, kind, payload jsonb, dedup_key unique, status)
- SELECT … FOR UPDATE SKIP LOCKED consumer
- LISTEN/NOTIFY realtime channel
- cmd/worker skeleton
- sqlc setup

## 6. Out of Scope
- Per-context job kinds (defined in their epics)

## 7. Dependencies
- e01s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- jobs

## 11. Interface / API Surface
- `jobs.Enqueue(ctx, kind, payload, dedupKey) error`
- `worker.Consume(ctx, handlers) error`

## 12. Business Rules Referenced
- BR-MIGRAR-014
- BR-DESCARTAR-004

## 13. Paradigm Notes (PHP → Go)
Swoole coroutines + Redis queues + outbox → single jobs table consumed by a goroutine pool with SKIP LOCKED; idempotency via dedup_key; terminal state per ADR-008.

## 14. Risks
- Risk tier **P0**. Primary risk: behavioral divergence from legacy caught by the parity gate before flip.

## 15. Non-Functional Requirements
- Parity: response shape + protocol match legacy within the agreed divergence budget.
- Security: secrets via platform/crypto; access via platform/auth+permissions.

## 16. Observability
- Tracing spans with canonical keys (project.id, user.id, function.id as applicable).
- Job state transitions recorded to a terminal state (ADR-008) where async.

## 17. Acceptance Criteria (Gherkin)
```gherkin
Feature: Base Postgres schema + jobs table + worker skeleton
  Scenario: exactly-once consumption
    Given two workers poll a single pending job
    Then exactly one worker claims it via SKIP LOCKED and drives it to a terminal state

  Scenario: idempotent enqueue
    Given the same job is enqueued twice with one dedup_key
    Then only one row exists and the duplicate is a no-op

  Scenario: realtime notify
    Given a row change fires NOTIFY on the realtime channel
    Then a LISTEN subscriber receives the event

```

## 18. Test Strategy
- Unit + integration tests against an ephemeral Postgres.
- Parallel Run shadow comparison against legacy (epic e15) before flip.

## 19. Definition of Done
- All tasks pass their `verify:` commands.
- Acceptance scenarios green.
- Parity gate for the parent context passes; no open deviations block the flip.

## 20. Traceability & References
- `specs/product/SCOPE_LATEST.yaml`
- `specs/migration/handoff.md`, `specs/migration/target_business_rules.md`
- `specs/migration/target_architecture.md`, `specs/migration/parity_specs.md`
