# e04s01 — Projects, platforms & API keys

## 1. Metadata
- **Story ID:** e04s01
- **Epic:** e04 — Projects (+ platforms, keys, webhooks)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **an org owner**, I want to **create projects and manage platforms and API keys**, so that **SDKs and clients authenticate against a project exactly as in legacy**.

## 3. Business Value
Projects is the sole writer of the projects collection; platforms and keys gate SDK access — the tenancy layer data/storage/compute build on.

## 4. Context & Background
Part of epic **e04 — Projects (+ platforms, keys, webhooks)** (Projects, platforms, keys, and webhooks.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- project CRUD (sole owner)
- platforms
- API keys (scoped, encrypted secret)
- /organization/projects becomes a read view

## 6. Out of Scope
- Webhooks (own story)

## 7. Dependencies
- e03s01
- e01s02
- e01s03

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- projects
- platforms
- keys

## 11. Interface / API Surface
- `POST /v1/projects`
- `POST /v1/projects/:id/keys`
- `POST /v1/projects/:id/platforms`

## 12. Business Rules Referenced
- BR-MIGRAR-007
- BR-MIGRAR-006

## 13. Paradigm Notes (PHP → Go)
Dual-ownership of projects removed → Projects is sole owner; Organizations reads a view.

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
Feature: Projects, platforms & API keys
  Scenario: create project
    Given an owner creates a project under a team
    Then the project is stored with Projects as sole owner

  Scenario: scoped key
    Given a request uses an API key limited to documents.read
    Then write routes are rejected, read routes allowed

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
