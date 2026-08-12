# e05s03 — Optimistic-lock multi-collection transactions

## 1. Metadata
- **Story ID:** e05s03
- **Epic:** e05 — Data (databases, documents, transactions, vectors)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** medium

## 2. User Story
As **an app developer**, I want to **commit changes across collections atomically**, so that **concurrent writers are serialized with version checks, matching legacy semantics**.

## 3. Business Value
Preserves atomic staging + version check + commit/rollback (BR-MIGRAR-004, ADR-005) as idiomatic Go/Postgres transactions.

## 4. Context & Background
Part of epic **e05 — Data (databases, documents, transactions, vectors)** (Databases, collections, documents with ACL, optimistic-lock transactions, and pgvector.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- staging buffer
- version check
- atomic commit
- rollback on conflict

## 6. Out of Scope
- Distributed transactions across services (out of scope)

## 7. Dependencies
- e05s02

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- documents
- _version columns

## 11. Interface / API Surface
- `transaction commit endpoint / batched write API`

## 12. Business Rules Referenced
- BR-MIGRAR-004

## 13. Paradigm Notes (PHP → Go)
utopia-php staged transaction → Postgres transaction with optimistic version columns.

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
Feature: Optimistic-lock multi-collection transactions
  Scenario: atomic commit
    Given a multi-collection write is committed
    Then all changes apply or none do

  Scenario: version conflict
    Given two writers stage conflicting updates to the same doc version
    Then one commits, the other is rejected and rolls back

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
