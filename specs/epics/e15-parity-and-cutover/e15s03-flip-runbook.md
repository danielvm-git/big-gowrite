# e15s03 — Per-context flip, canary & PHP decommission

## 1. Metadata
- **Story ID:** e15s03
- **Epic:** e15 — Parity harness & cutover
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P0  ·  **Security surface:** high

## 2. User Story
As **the migration operator**, I want to **flip a context to Go behind a canary with instant rollback**, so that **each context goes live safely and PHP is retired incrementally**.

## 3. Business Value
The operational runbook that turns a passed gate into a live flip: Traefik reroute, canary, rollback-by-route, and decommission.

## 4. Context & Background
Part of epic **e15 — Parity harness & cutover** (Parallel Run harness, contract tests (response goldens, WS, stream, error catalog), and the per-context flip/canary/decommission runbook.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- Traefik route flip by context
- canary %
- rollback = re-route
- PHP decommission checklist
- Identity+Orgs atomic co-flip

## 6. Out of Scope
- Big-bang cutover (rejected)

## 7. Dependencies
- e15s01
- e15s02

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- (infra/runbook)

## 11. Interface / API Surface
- `Traefik dynamic config`

## 12. Business Rules Referenced
- (no specific BR — mechanical/CRUD parity)

## 13. Paradigm Notes (PHP → Go)
Strangler Fig flip by context; rollback is a route change.

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
Feature: Per-context flip, canary & PHP decommission
  Scenario: canary flip
    Given a context passes its gate and is flipped at 10% canary
    Then 10% of traffic hits Go and metrics are compared before full flip

  Scenario: rollback
    Given a regression appears post-flip
    Then the route reverts to PHP with no data migration

  Scenario: atomic co-flip
    Given Identity and Organizations are flipped
    Then they flip together as one unit (team-invite spans both)

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
