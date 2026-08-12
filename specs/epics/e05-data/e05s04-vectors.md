# e05s04 — VectorsDB via pgvector

## 1. Metadata
- **Story ID:** e05s04
- **Epic:** e05 — Data (databases, documents, transactions, vectors)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P2  ·  **Security surface:** low

## 2. User Story
As **a developer**, I want to **store and query vector embeddings**, so that **similarity search works without a separate vector store**.

## 3. Business Value
Ports VectorsDB (VAR_VECTOR, configurable dimension) onto pgvector, keeping embeddings search in the single datastore.

## 4. Context & Background
Part of epic **e05 — Data (databases, documents, transactions, vectors)** (Databases, collections, documents with ACL, optimistic-lock transactions, and pgvector.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- vector attribute type
- configurable dimension
- similarity query

## 6. Out of Scope
- Embedding generation (external)

## 7. Dependencies
- e05s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- vector_documents

## 11. Interface / API Surface
- `vector attribute + nearest-neighbor query`

## 12. Business Rules Referenced
- BR-MIGRAR-018

## 13. Paradigm Notes (PHP → Go)
VAR_VECTOR filter → pgvector column + ivfflat/hnsw index.

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
Feature: VectorsDB via pgvector
  Scenario: knn query
    Given vectors are stored and a nearest-neighbor query is issued
    Then the closest vectors are returned in distance order

  Scenario: dimension mismatch
    Given a vector of wrong dimension is inserted
    Then the write is rejected

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
