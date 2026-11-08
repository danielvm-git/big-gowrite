# e02s01 — Users & accounts with password-hash import parity

## 1. Metadata
- **Story ID:** e02s01
- **Epic:** e02 — Identity (accounts, sessions, MFA, tokens)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** high

## 2. User Story
As **an API client**, I want to **create and manage users/accounts and authenticate with imported password hashes**, so that **existing credentials (Bcrypt/MD5/Argon2/SHA/PHPass/Scrypt variants) keep working**.

## 3. Business Value
Core identity CRUD plus multi-algorithm hash import (BR-MIGRAR-016) — required so existing users authenticate unchanged after the flip.

## 4. Context & Background
Part of epic **e02 — Identity (accounts, sessions, MFA, tokens)** (Accounts, users, sessions, MFA/TOTP, tokens, and password-hash import parity.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- user/account CRUD
- password verify across all legacy hash algorithms
- hash import
- email/phone uniqueness

## 6. Out of Scope
- Sessions (own story)
- MFA (own story)

## 7. Dependencies
- e01s02
- e01s03

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- users

## 11. Interface / API Surface
- `POST /v1/account`
- `GET/PATCH /v1/users/:id`
- `auth.VerifyPassword(hash, plain)`

## 12. Business Rules Referenced
- BR-MIGRAR-016

## 13. Paradigm Notes (PHP → Go)
Active-record-ish user document → typed User struct + sqlc; hash verify dispatched by stored algorithm tag.

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
Feature: Users & accounts with password-hash import parity
  Scenario: create user
    Given a client POSTs a new account with email+password
    Then the user is stored with an encrypted/hashed password and returned without the secret

  Scenario: legacy hash login
    Given a user imported with a Bcrypt hash logs in with the correct password
    Then verification succeeds; an Argon2/MD5/SHA/PHPass/Scrypt user behaves identically

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
