# e02s02 — Sessions & MFA/TOTP lifecycle

## 1. Metadata
- **Story ID:** e02s02
- **Epic:** e02 — Identity (accounts, sessions, MFA, tokens)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** high

## 2. User Story
As **an authenticated user**, I want to **create sessions and enroll/verify TOTP factors**, so that **MFA challenges, factor dedupe, and session lifecycle match legacy exactly**.

## 3. Business Value
Preserves the MFA/TOTP state machine and session factor recomputation (BR-MIGRAR-001) — auth correctness under the new paradigm.

## 4. Context & Background
Part of epic **e02 — Identity (accounts, sessions, MFA, tokens)** (Accounts, users, sessions, MFA/TOTP, tokens, and password-hash import parity.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- session create/list/delete
- TOTP secret→OTP confirm→persist authenticator
- factor recompute/dedupe on session
- challenge rate limits

## 6. Out of Scope
- OAuth token exchange (tokens story)

## 7. Dependencies
- e02s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- sessions
- authenticators
- challenges

## 11. Interface / API Surface
- `POST /v1/account/sessions`
- `POST /v1/account/mfa/authenticators/totp`
- `PATCH .../sessions/:id`

## 12. Business Rules Referenced
- BR-MIGRAR-001

## 13. Paradigm Notes (PHP → Go)
Challenge documents + state → typed challenge state machine; 409 recreate-blocked becomes explicit Go error.

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
Feature: Sessions & MFA/TOTP lifecycle
  Scenario: totp enroll
    Given a user posts a TOTP secret then confirms with a valid OTP
    Then the authenticator is persisted and the session factors include totp

  Scenario: recreate blocked
    Given a user re-enrolls an already-confirmed TOTP factor
    Then the request fails with 409, matching legacy

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
