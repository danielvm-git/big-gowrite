# e06s02 — Chunked upload with concurrent-safe finalize

## 1. Metadata
- **Story ID:** e06s02
- **Epic:** e06 — Storage (buckets, files, chunked upload)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** high

## 2. User Story
As **a client uploading a large file**, I want to **send chunks and finalize safely under concurrency**, so that **compression, encryption, and antivirus run once, with 429+Retry-After on contention**.

## 3. Business Value
Preserves chunked upload (BR-MIGRAR-005, ADR-004): Content-Range, finalize with compression/encryption/antivirus, and lock contention → 429.

## 4. Context & Background
Part of epic **e06 — Storage (buckets, files, chunked upload)** (Buckets and files with concurrent-safe chunked upload.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- Content-Range chunk assembly
- finalize: compress + encrypt + ClamAV
- per-resource exclusion → 429 Retry-After

## 6. Out of Scope
- ClamAV service itself (external dependency)

## 7. Dependencies
- e06s01
- e01s02

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- files
- upload state

## 11. Interface / API Surface
- `POST .../files with Content-Range`

## 12. Business Rules Referenced
- BR-MIGRAR-005

## 13. Paradigm Notes (PHP → Go)
Swoole concurrency + distributed lock → Go mutex/advisory-lock per resource; ClamAV over network.

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
Feature: Chunked upload with concurrent-safe finalize
  Scenario: chunked finalize
    Given all chunks of a file are uploaded and finalized
    Then the stored file is compressed, encrypted, and virus-scanned once

  Scenario: lock contention
    Given two finalize requests race on the same file
    Then one proceeds, the other gets 429 with Retry-After

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
