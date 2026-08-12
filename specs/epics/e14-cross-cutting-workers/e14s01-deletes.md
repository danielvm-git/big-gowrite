# e14s01 — Deletes worker (cascade)

## 1. Metadata
- **Story ID:** e14s01
- **Epic:** e14 — Cross-cutting workers (Deletes, Stats)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P0  ·  **Security surface:** medium

## 2. User Story
As **the system**, I want to **cascade-delete a resource and all its dependents**, so that **no orphaned rows remain, matching legacy cascade semantics**.

## 3. Business Value
Cascade delete correctness across all contexts; migrated last to avoid orphaning during coexistence (S3).

## 4. Context & Background
Part of epic **e14 — Cross-cutting workers (Deletes, Stats)** (Cascade deletes and stats aggregation workers.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- cascade delete worker via jobs
- cross-context cleanup ordering

## 6. Out of Scope
- (runs only after all contexts are on Postgres)

## 7. Dependencies
- e02s01
- e03s01
- e04s01
- e05s01
- e06s01
- e08s01
- e01s05

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- all context tables

## 11. Interface / API Surface
- `worker: cascade delete`

## 12. Business Rules Referenced
- BR-MIGRAR-014

## 13. Paradigm Notes (PHP → Go)
Swoole Deletes worker → jobs worker driving ordered cascade in Postgres.

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
Feature: Deletes worker (cascade)
  Scenario: cascade
    Given a team with projects/documents/files is deleted
    Then all dependents are removed with no orphans

  Scenario: idempotent
    Given a delete job is retried after partial completion
    Then it completes without error and reaches a terminal state

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
