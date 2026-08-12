# e14s02 — Stats worker

## 1. Metadata
- **Story ID:** e14s02
- **Epic:** e14 — Cross-cutting workers (Deletes, Stats)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 2
- **Risk:** P2  ·  **Security surface:** low

## 2. User Story
As **a project owner**, I want to **see usage statistics**, so that **stats aggregation matches legacy numbers**.

## 3. Business Value
Usage/resource stats aggregation (publishers in legacy) as a Go worker.

## 4. Context & Background
Part of epic **e14 — Cross-cutting workers (Deletes, Stats)** (Cascade deletes and stats aggregation workers.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- stats aggregation worker
- usage counters

## 6. Out of Scope
- Real-time billing

## 7. Dependencies
- e01s05

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- stats

## 11. Interface / API Surface
- `worker: stats aggregation`

## 12. Business Rules Referenced
- (no specific BR — mechanical/CRUD parity)

## 13. Paradigm Notes (PHP → Go)
Stats publishers → jobs/interval worker aggregating counters.

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
Feature: Stats worker
  Scenario: aggregate
    Given usage events accumulate over a window
    Then the stats worker aggregates them into the expected totals

  Scenario: idempotent window
    Given an aggregation window is reprocessed
    Then totals are not double-counted

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
