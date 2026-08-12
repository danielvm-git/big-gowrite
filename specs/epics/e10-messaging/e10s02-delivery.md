# e10s02 — Idempotent delivery worker

## 1. Metadata
- **Story ID:** e10s02
- **Epic:** e10 — Messaging (email/SMS/push)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **the messaging worker**, I want to **deliver a message to its provider exactly once**, so that **a message reaches sent or failed and never double-sends**.

## 3. Business Value
Delivers messages via providers with an idempotent terminal state (BR-MIGRAR-013/014).

## 4. Context & Background
Part of epic **e10 — Messaging (email/SMS/push)** (Message lifecycle and provider delivery.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- provider adapters (email/SMS/push)
- delivery worker via jobs
- idempotent terminal

## 6. Out of Scope
- Provider account setup

## 7. Dependencies
- e10s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- messages

## 11. Interface / API Surface
- `worker: message delivery`

## 12. Business Rules Referenced
- BR-MIGRAR-013
- BR-MIGRAR-014

## 13. Paradigm Notes (PHP → Go)
Redis queue delivery → jobs worker; terminal state idempotent.

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
Feature: Idempotent delivery worker
  Scenario: deliver once
    Given a scheduled message is due
    Then the worker delivers it and marks it sent exactly once

  Scenario: provider failure
    Given the provider rejects the message
    Then it is marked failed with no retry storm

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
