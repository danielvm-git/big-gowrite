# e08s01 — Functions & deployment state machine

## 1. Metadata
- **Story ID:** e08s01
- **Epic:** e08 — Compute (functions, sites, executions)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** medium

## 2. User Story
As **a developer**, I want to **deploy a function and track its build to a terminal state**, so that **deployment states and cancel guards match legacy exactly**.

## 3. Business Value
Preserves the deployment lifecycle (BR-MIGRAR-002): waiting→processing→building→ready|failed|canceled, with build callback idempotency.

## 4. Context & Background
Part of epic **e08 — Compute (functions, sites, executions)** (Functions and sites with deployment/execution state machines and VCS deploy fan-out.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- function CRUD
- deployment state machine
- build via jobs + external executor
- cancel guard on terminal state
- CloudEvent id dedupe

## 6. Out of Scope
- Executions (own story)
- the executor service itself (external)

## 7. Dependencies
- e07s02
- e01s05

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- functions
- deployments

## 11. Interface / API Surface
- `POST /v1/functions`
- `POST .../deployments`
- `build callback`

## 12. Business Rules Referenced
- BR-MIGRAR-002
- BR-MIGRAR-014

## 13. Paradigm Notes (PHP → Go)
Swoole build coroutine → jobs worker driving the state machine; terminal-state + idempotent callback.

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
Feature: Functions & deployment state machine
  Scenario: build to ready
    Given a deployment is created and the executor reports success
    Then the deployment ends in ready via waiting→processing→building

  Scenario: callback dedupe
    Given the same build callback CloudEvent id arrives twice
    Then the state transition applies once

  Scenario: cancel on terminal
    Given a cancel is requested on an already-ready deployment
    Then the request is rejected

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
