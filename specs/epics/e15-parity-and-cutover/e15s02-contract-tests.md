# e15s02 — WS / stream / error-catalog contract tests

## 1. Metadata
- **Story ID:** e15s02
- **Epic:** e15 — Parity harness & cutover
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** medium

## 2. User Story
As **the migration operator**, I want to **assert protocol-level parity beyond plain request/response**, so that **realtime, streaming, and error contracts match legacy so the SPA keeps working**.

## 3. Business Value
Covers the non-REST deviations (DEV-001..008): WS realtime envelope, assistant raw text stream (not SSE), error type/code catalog, response .total shape.

## 4. Context & Background
Part of epic **e15 — Parity harness & cutover** (Parallel Run harness, contract tests (response goldens, WS, stream, error catalog), and the per-context flip/canary/decommission runbook.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- WS realtime envelope tests
- assistant raw-text stream framing (not SSE)
- error type/code catalog
- .total shape parity
- Date header, console headers

## 6. Out of Scope
- Multi-region (DEV-009, deferred)

## 7. Dependencies
- e15s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- (test infra)

## 11. Interface / API Surface
- `WS/stream/error contract suite`

## 12. Business Rules Referenced
- (no specific BR — mechanical/CRUD parity)

## 13. Paradigm Notes (PHP → Go)
Contract parity, per screen_deviation_log DEV-001..008.

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
Feature: WS / stream / error-catalog contract tests
  Scenario: ws envelope
    Given a realtime event is emitted
    Then the WS frame matches the legacy envelope shape

  Scenario: assistant stream
    Given the assistant endpoint streams a response
    Then it emits raw concatenated text (streamProtocol:text), not SSEframes

  Scenario: error catalog
    Given a known error occurs
    Then the JSON type+code match the legacy error catalog

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
