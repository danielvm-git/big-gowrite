# e01s01 — Explicit error model + HTTP translation

## 1. Metadata
- **Story ID:** e01s01
- **Epic:** e01 — Platform Foundation
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **a gowrite context author**, I want to **return typed errors that carry an HTTP status and public code**, so that **the API surface stays consistent with the legacy error catalog without per-handler try/catch**.

## 3. Business Value
Replaces PHP try/catch/AppwriteException with idiomatic Go error values, giving every context one error contract and one HTTP translation point.

## 4. Context & Background
Part of epic **e01 — Platform Foundation** (Shared leaf packages (errors, auth, permissions, crypto, ratelimit), base Postgres schema, and the jobs table with SKIP-LOCKED dispatch + LISTEN/NOTIFY realtime.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- platform/errors.AppError{Code, HTTPStatus, Err}
- error-to-HTTP middleware
- error code catalog matching legacy

## 6. Out of Scope
- Per-context domain errors (defined in their own epics)

## 7. Dependencies
- None (foundational)

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- (none — cross-cutting package)

## 11. Interface / API Surface
- `errors.New(code, status, wrap) AppError`
- `middleware.TranslateError(next) http.Handler`

## 12. Business Rules Referenced
- BR-DESCARTAR-001

## 13. Paradigm Notes (PHP → Go)
PHP exceptions + AppwriteException → Go (v, error) returns; middleware translates AppError to the legacy JSON error envelope (type/code/message).

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
Feature: Explicit error model + HTTP translation
  Scenario: handler returns AppError
    Given a handler returns errors.New("user_not_found", 404, err)
    Then the client receives HTTP 404 with JSON {type:"user_not_found", code:404, message:...}

  Scenario: unexpected error is masked
    Given a handler returns a raw non-AppError error
    Then the client receives HTTP 500 with a generic message and no internal detail leaks

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
