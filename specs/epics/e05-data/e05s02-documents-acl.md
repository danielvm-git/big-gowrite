# e05s02 — Documents with document-level ACL

## 1. Metadata
- **Story ID:** e05s02
- **Epic:** e05 — Data (databases, documents, transactions, vectors)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** high

## 2. User Story
As **an app end-user**, I want to **read/write documents subject to their roles**, so that **ACL filtering matches legacy: unauthorized documents are invisible**.

## 3. Business Value
Document CRUD with per-document ACL (BR-MIGRAR-007) — the most security-sensitive data path.

## 4. Context & Background
Part of epic **e05 — Data (databases, documents, transactions, vectors)** (Databases, collections, documents with ACL, optimistic-lock transactions, and pgvector.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- document CRUD (JSONB)
- per-document _permissions
- query + ACL filter
- documentSecurity flag

## 6. Out of Scope
- Transactions (own story)

## 7. Dependencies
- e05s01
- e01s03

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- documents

## 11. Interface / API Surface
- `POST/GET/PATCH/DELETE /v1/databases/:db/collections/:c/documents`

## 12. Business Rules Referenced
- BR-MIGRAR-007

## 13. Paradigm Notes (PHP → Go)
subQuery permission filter → SQL WHERE built from Principal roles; JSONB payload column.

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
Feature: Documents with document-level ACL
  Scenario: acl read filter
    Given a user queries documents where some lack a matching read role
    Then only permitted documents are returned

  Scenario: acl write
    Given a user without write role updates a document
    Then the write is rejected with legacy-equivalent status

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
