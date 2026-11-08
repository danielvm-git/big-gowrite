# e09s02 — TLS certificate lifecycle

## 1. Metadata
- **Story ID:** e09s02
- **Epic:** e09 — Networking (proxy rules, domains, TLS)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** high

## 2. User Story
As **the system**, I want to **obtain and renew TLS certificates for verified domains**, so that **HTTPS on custom domains works as in legacy**.

## 3. Business Value
Issues/renews certificates for verified domains; certificate material handled securely.

## 4. Context & Background
Part of epic **e09 — Networking (proxy rules, domains, TLS)** (Proxy rules, domain verification, and TLS certificate lifecycle.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- certificate issue/renew via jobs worker
- secure key storage

## 6. Out of Scope
- ACME provider internals (library)

## 7. Dependencies
- e09s01
- e01s05

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- certificates

## 11. Interface / API Surface
- `cert issuance worker`

## 12. Business Rules Referenced
- BR-MIGRAR-010
- BR-MIGRAR-006

## 13. Paradigm Notes (PHP → Go)
Certificate worker → jobs table kind=certificate; keys stored encrypted.

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
Feature: TLS certificate lifecycle
  Scenario: issue
    Given a domain becomes verified
    Then a certificate is issued and stored securely

  Scenario: renew
    Given a certificate nears expiry
    Then a renewal job runs and replaces it

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
