# e09s01 — Proxy rules & domain verification

## 1. Metadata
- **Story ID:** e09s01
- **Epic:** e09 — Networking (proxy rules, domains, TLS)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **a project admin**, I want to **add custom domains and verify ownership**, so that **domain verification and routing rules match legacy**.

## 3. Business Value
Preserves the rule/proxy state machine (BR-MIGRAR-010) for custom domains.

## 4. Context & Background
Part of epic **e09 — Networking (proxy rules, domains, TLS)** (Proxy rules, domain verification, and TLS certificate lifecycle.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- rule CRUD
- domain verification state machine

## 6. Out of Scope
- TLS issuance (own story)

## 7. Dependencies
- e04s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- rules

## 11. Interface / API Surface
- `POST /v1/proxy/rules`
- `domain verification`

## 12. Business Rules Referenced
- BR-MIGRAR-010

## 13. Paradigm Notes (PHP → Go)
Rule document + state → typed state machine.

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
Feature: Proxy rules & domain verification
  Scenario: verify domain
    Given a domain's DNS is correctly configured and verification runs
    Then the rule transitions to verified

  Scenario: unverified
    Given DNS is missing
    Then verification fails and the rule stays unverified

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
