# e08s02 — Execution state machine

## 1. Metadata
- **Story ID:** e08s02
- **Epic:** e08 — Compute (functions, sites, executions)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **a developer**, I want to **invoke a function and read its execution result**, so that **execution status derivation (≥500→failed, sync-timeout→failed) matches legacy**.

## 3. Business Value
Preserves execution lifecycle (BR-MIGRAR-003): waiting/processing/scheduled→completed|failed|canceled, status derived from HTTP code.

## 4. Context & Background
Part of epic **e08 — Compute (functions, sites, executions)** (Functions and sites with deployment/execution state machines and VCS deploy fan-out.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- execution create (sync/async/scheduled)
- status from HTTP code
- sync-timeout→failed (no timeout state)

## 6. Out of Scope
- Scheduling infrastructure specifics (jobs table)

## 7. Dependencies
- e08s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- executions

## 11. Interface / API Surface
- `POST /v1/functions/:id/executions`

## 12. Business Rules Referenced
- BR-MIGRAR-003

## 13. Paradigm Notes (PHP → Go)
Execution document + coroutine → typed struct; status derived explicitly, no persisted timeout.

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
Feature: Execution state machine
  Scenario: sync success
    Given a sync execution returns HTTP 200
    Then the execution is completed

  Scenario: server error
    Given the function returns HTTP 500
    Then the execution is failed

  Scenario: sync timeout
    Given a sync execution exceeds the timeout
    Then the execution is failed (no timeout state persisted)

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
