# e10s01 — Message lifecycle & providers

## 1. Metadata
- **Story ID:** e10s01
- **Epic:** e10 — Messaging (email/SMS/push)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **a project admin**, I want to **compose and schedule messages to targets**, so that **message states and provider selection match legacy**.

## 3. Business Value
Preserves the message state machine (BR-MIGRAR-013): draft→scheduled→processing→sent|failed across providers.

## 4. Context & Background
Part of epic **e10 — Messaging (email/SMS/push)** (Message lifecycle and provider delivery.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- message + provider + topic/target CRUD
- state machine
- scheduling via jobs

## 6. Out of Scope
- Actual provider network calls (delivery story)

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
- messages
- providers
- topics
- targets

## 11. Interface / API Surface
- `POST /v1/messaging/messages`
- `providers/topics/targets`

## 12. Business Rules Referenced
- BR-MIGRAR-013

## 13. Paradigm Notes (PHP → Go)
Message document + state → typed state machine; schedule via jobs.

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
Feature: Message lifecycle & providers
  Scenario: schedule
    Given a message is scheduled for a future time
    Then it sits in scheduled until the time, then moves to processing

  Scenario: draft edit
    Given a draft message is edited
    Then it remains editable until scheduled/processing

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
