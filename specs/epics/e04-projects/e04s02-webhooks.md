# e04s02 — Webhooks

## 1. Metadata
- **Story ID:** e04s02
- **Epic:** e04 — Projects (+ platforms, keys, webhooks)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **a project admin**, I want to **register webhooks and have events delivered with signatures**, so that **event fan-out to webhooks matches legacy payloads and retry behavior**.

## 3. Business Value
Project-scoped webhook delivery (moved into Projects per red-team B4); async via the jobs table.

## 4. Context & Background
Part of epic **e04 — Projects (+ platforms, keys, webhooks)** (Projects, platforms, keys, and webhooks.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- webhook CRUD
- signature
- delivery via jobs table with retry

## 6. Out of Scope
- The event taxonomy itself (emitted by each context)

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
- webhooks

## 11. Interface / API Surface
- `POST /v1/projects/:id/webhooks`
- `worker: webhook delivery`

## 12. Business Rules Referenced
- BR-MIGRAR-014

## 13. Paradigm Notes (PHP → Go)
Redis webhook queue → jobs table kind=webhook; terminal state + retry.

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
Feature: Webhooks
  Scenario: delivery
    Given an event matches a registered webhook
    Then a signed POST is delivered and marked terminal on success

  Scenario: retry
    Given a webhook endpoint returns 5xx
    Then delivery retries per policy and records a terminal failed state after exhaustion

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
