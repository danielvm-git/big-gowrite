# e02s03 — Verification, recovery, magic-url & OAuth tokens

## 1. Metadata
- **Story ID:** e02s03
- **Epic:** e02 — Identity (accounts, sessions, MFA, tokens)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** high

## 2. User Story
As **an API client**, I want to **issue and consume verification/recovery/magic-url/OAuth tokens**, so that **account confirmation and recovery flows behave identically to legacy**.

## 3. Business Value
Completes identity surface: the token flows that feed email/SMS verification, recovery, and OAuth — needed for full account parity.

## 4. Context & Background
Part of epic **e02 — Identity (accounts, sessions, MFA, tokens)** (Accounts, users, sessions, MFA/TOTP, tokens, and password-hash import parity.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- token create/consume
- expiry
- OAuth provider token storage (encrypted)

## 6. Out of Scope
- Delivery of the emails/SMS (messaging epic)

## 7. Dependencies
- e02s01
- e01s02

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- tokens

## 11. Interface / API Surface
- `POST /v1/account/verifications`
- `PATCH /v1/account/recovery`
- `OAuth callback handlers`

## 12. Business Rules Referenced
- BR-MIGRAR-006

## 13. Paradigm Notes (PHP → Go)
Token documents → typed Token struct; OAuth secrets stored via crypto.Encrypt.

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
Feature: Verification, recovery, magic-url & OAuth tokens
  Scenario: token consume
    Given a valid, unexpired verification token is consumed
    Then the account is marked verified and the token is invalidated

  Scenario: expired token
    Given an expired token is presented
    Then the request is rejected without side effects

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
