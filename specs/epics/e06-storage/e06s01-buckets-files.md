# e06s01 — Buckets & files

## 1. Metadata
- **Story ID:** e06s01
- **Epic:** e06 — Storage (buckets, files, chunked upload)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **an app developer**, I want to **create buckets and manage files**, so that **file metadata and access control match legacy**.

## 3. Business Value
Baseline storage CRUD and signed access, project-scoped.

## 4. Context & Background
Part of epic **e06 — Storage (buckets, files, chunked upload)** (Buckets and files with concurrent-safe chunked upload.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- bucket/file CRUD
- signed URLs / token access
- per-file permissions

## 6. Out of Scope
- Chunked upload finalize (own story)

## 7. Dependencies
- e04s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- buckets
- files

## 11. Interface / API Surface
- `POST /v1/storage/buckets`
- `GET /v1/storage/buckets/:b/files/:f/view`

## 12. Business Rules Referenced
- BR-MIGRAR-007

## 13. Paradigm Notes (PHP → Go)
Storage documents → typed structs; object bytes on FS/object storage.

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
Feature: Buckets & files
  Scenario: create bucket
    Given a developer creates a bucket
    Then the bucket is stored with its limits

  Scenario: permissioned view
    Given a user without read role views a file
    Then access is denied per legacy

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
