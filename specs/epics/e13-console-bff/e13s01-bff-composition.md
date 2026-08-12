# e13s01 — Console BFF composition

## 1. Metadata
- **Story ID:** e13s01
- **Epic:** e13 — Console BFF + SPA hosting
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **the console SPA**, I want to **call BFF endpoints that aggregate identity/org/project data**, so that **the console works against Go with the same responses as legacy**.

## 3. Business Value
Composes identity+organizations+projects for the console without owning aggregates (red-team B3).

## 4. Context & Background
Part of epic **e13 — Console BFF + SPA hosting** (BFF composition over identity+orgs+projects and serving the preserved SPA via go:embed.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- BFF endpoints aggregating owned contexts
- console-only headers (DEV-006)
- cookie-session (DEV-005)

## 6. Out of Scope
- Owning any aggregate (delegates only)

## 7. Dependencies
- e02s01
- e03s01
- e04s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- (none — composition layer)

## 11. Interface / API Surface
- `console BFF routes`

## 12. Business Rules Referenced
- BR-MIGRAR-007

## 13. Paradigm Notes (PHP → Go)
Console module → BFF delegating to contexts; no data ownership.

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
Feature: Console BFF composition
  Scenario: aggregate view
    Given the console requests a project overview
    Then the BFF composes identity+org+project data into the legacy response shape

  Scenario: cookie session
    Given the console authenticates via cookie session
    Then the BFF honors the cookie-session contract (DEV-005)

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
