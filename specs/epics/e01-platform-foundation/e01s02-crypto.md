# e01s02 — AES-GCM encryption for 25 at-rest fields

## 1. Metadata
- **Story ID:** e01s02
- **Epic:** e01 — Platform Foundation
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P0  ·  **Security surface:** high

## 2. User Story
As **a context persisting secrets**, I want to **call crypto.Encrypt/Decrypt on the 25 designated fields**, so that **OAuth tokens, passwords, and credentials are never stored in plaintext, matching legacy parity**.

## 3. Business Value
Encryption at rest is a non-negotiable security rule (BR-MIGRAR-006); getting the typed crypto layer right unblocks every context that persists secrets.

## 4. Context & Background
Part of epic **e01 — Platform Foundation** (Shared leaf packages (errors, auth, permissions, crypto, ratelimit), base Postgres schema, and the jobs table with SKIP-LOCKED dispatch + LISTEN/NOTIFY realtime.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- platform/crypto AES-GCM encrypt/decrypt
- key derivation from env
- the 25-field registry
- re-encrypt helper for data migration (R-D3)

## 6. Out of Scope
- Password hashing (identity epic, import parity)
- TLS certs (networking epic)

## 7. Dependencies
- None (foundational)

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- 25 encrypted columns across users/projects/keys/installations/etc.

## 11. Interface / API Surface
- `crypto.Encrypt(plaintext []byte) (cipher []byte, error)`
- `crypto.Decrypt(cipher []byte) ([]byte, error)`

## 12. Business Rules Referenced
- BR-MIGRAR-006

## 13. Paradigm Notes (PHP → Go)
utopia-php `encrypt` document filter → explicit crypto.Encrypt() call at each write site (no implicit filter engine).

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
Feature: AES-GCM encryption for 25 at-rest fields
  Scenario: round-trip
    Given a value is encrypted with crypto.Encrypt then decrypted with crypto.Decrypt
    Then the decrypted value byte-equals the original

  Scenario: ciphertext is not plaintext
    Given a secret is encrypted
    Then the stored bytes do not contain the plaintext and include a GCM nonce+tag

  Scenario: legacy decrypt parity
    Given a value encrypted by the legacy PHP filter is read
    Then crypto.Decrypt recovers the original value using the shared key

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
