# e08s03 — Sites & VCS deploy fan-out + preview rules

## 1. Metadata
- **Story ID:** e08s03
- **Epic:** e08 — Compute (functions, sites, executions)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **a developer**, I want to **deploy sites and get preview URLs from Git pushes**, so that **push-triggered deploys and preview rules match legacy**.

## 3. Business Value
Completes compute: sites plus VCS-triggered deploy fan-out and preview rule creation (BR-MIGRAR-015).

## 4. Context & Background
Part of epic **e08 — Compute (functions, sites, executions)** (Functions and sites with deployment/execution state machines and VCS deploy fan-out.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- site CRUD
- deploy fan-out on repository.pushed
- preview rule creation

## 6. Out of Scope
- Proxy rule/domain lifecycle (networking epic)

## 7. Dependencies
- e08s01
- e07s02

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- sites
- deployments
- rules

## 11. Interface / API Surface
- `POST /v1/sites`
- `consumes repository.pushed`

## 12. Business Rules Referenced
- BR-MIGRAR-015

## 13. Paradigm Notes (PHP → Go)
Coroutine fan-out → worker consuming repository.pushed and creating deployments + preview rules.

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
Feature: Sites & VCS deploy fan-out + preview rules
  Scenario: push deploy
    Given a linked repo receives a push
    Then a site deployment is created and a preview rule is generated

  Scenario: preview url
    Given a preview deployment reaches ready
    Then a preview URL rule points to it

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
