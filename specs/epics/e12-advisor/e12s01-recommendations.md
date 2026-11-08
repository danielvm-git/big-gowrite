# e12s01 — Advisor recommendations

## 1. Metadata
- **Story ID:** e12s01
- **Epic:** e12 — Advisor (recommendations)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 2
- **Risk:** P2  ·  **Security surface:** low

## 2. User Story
As **a project owner**, I want to **receive advisory recommendations about my project**, so that **advisor output matches legacy**.

## 3. Business Value
Ports the advisor context that surfaces project recommendations.

## 4. Context & Background
Part of epic **e12 — Advisor (recommendations)** (Advisor recommendations.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- recommendation computation
- advisor endpoints

## 6. Out of Scope
- New recommendation types beyond parity

## 7. Dependencies
- e04s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- (derived/read-mostly)

## 11. Interface / API Surface
- `GET advisor recommendations`

## 12. Business Rules Referenced
- (no specific BR — mechanical/CRUD parity)

## 13. Paradigm Notes (PHP → Go)
Advisor logic → typed Go functions.

## 14. Risks
- Risk tier **P2**. Primary risk: behavioral divergence from legacy caught by the parity gate before flip.

## 15. Non-Functional Requirements
- Parity: response shape + protocol match legacy within the agreed divergence budget.
- Security: secrets via platform/crypto; access via platform/auth+permissions.

## 16. Observability
- Tracing spans with canonical keys (project.id, user.id, function.id as applicable).
- Job state transitions recorded to a terminal state (ADR-008) where async.

## 17. Acceptance Criteria (Gherkin)
```gherkin
Feature: Advisor recommendations
  Scenario: recommendations
    Given a project with a known misconfiguration is evaluated
    Then the matching recommendation is returned as in legacy

  Scenario: clean project
    Given a well-configured project is evaluated
    Then no false recommendations are returned

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
