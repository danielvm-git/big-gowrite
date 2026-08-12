# e05s01 — Databases, collections & typed schema

## 1. Metadata
- **Story ID:** e05s01
- **Epic:** e05 — Data (databases, documents, transactions, vectors)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **a project developer**, I want to **define databases, collections and attributes**, so that **schema definition and indexing behave as in legacy without the generic PHP engine**.

## 3. Business Value
Ports the 48-collection/606-attribute schema (BR-MIGRAR-018) as typed Go + sqlc — the lean subset resolving BR-HUMANA-001.

## 4. Context & Background
Part of epic **e05 — Data (databases, documents, transactions, vectors)** (Databases, collections, documents with ACL, optimistic-lock transactions, and pgvector.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- database/collection CRUD
- attribute types + indexes (13 unique, 18 fulltext)
- typed schema representation

## 6. Out of Scope
- Documents (own story)
- generic multi-DB engine (out of scope, BR-HUMANA-001)

## 7. Dependencies
- e04s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- databases
- collections
- attributes
- indexes

## 11. Interface / API Surface
- `POST /v1/databases`
- `POST .../collections`
- `POST .../attributes/*`

## 12. Business Rules Referenced
- BR-MIGRAR-018
- BR-DESCARTAR-005

## 13. Paradigm Notes (PHP → Go)
PHP-array runtime schema → static Go structs + Postgres DDL/migrations; sqlc typed access.

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
Feature: Databases, collections & typed schema
  Scenario: create collection
    Given a developer defines a collection with string + integer attributes
    Then the collection and columns exist and reject wrong-typed writes

  Scenario: fulltext index
    Given a fulltext index is declared
    Then text search over that attribute works equivalently to legacy

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
