# e11s01 — Migration jobs state machine + retry

## 1. Metadata
- **Story ID:** e11s01
- **Epic:** e11 — Migrations (data import jobs)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** low

## 2. User Story
As **a project owner**, I want to **run and retry data migrations/imports**, so that **migration states and retry rules match legacy**.

## 3. Business Value
Preserves migration lifecycle (BR-MIGRAR-011): pending→processing→completed|failed, retry only from failed.

## 4. Context & Background
Part of epic **e11 — Migrations (data import jobs)** (Migration jobs state machine.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- migration CRUD
- state machine via jobs
- retry only from failed

## 6. Out of Scope
- Specific source connectors beyond parity scope

## 7. Dependencies
- e04s01
- e01s05

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- migrations

## 11. Interface / API Surface
- `POST /v1/migrations`
- `retry endpoint`

## 12. Business Rules Referenced
- BR-MIGRAR-011

## 13. Paradigm Notes (PHP → Go)
Migration document + worker → jobs worker; retry guard explicit.

## 14. Risks
- Risk tier **P1**. Primary risk: behavioral divergence from legacy caught by the parity gate before flip.

## 15. Non-Functional Requirements
- Parity: response shape + protocol match legacy within the agreed divergence budget.
- Security: secrets via platform/crypto; access via platform/auth+permissions.

## 16. Observability
- Tracing spans with canonical keys (project.id, user.id, function.id as applicable).
- Job state transitions recorded to a terminal state (ADR-008) where async.

## 17. Acceptance Criteria (Gherkin)
```gherkin
Feature: Migration jobs state machine + retry
  Scenario: run
    Given a migration is started
    Then it moves pending→processing→completed on success

  Scenario: retry guard
    Given a completed migration is retried
    Then the retry is rejected; only failed migrations retry

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
