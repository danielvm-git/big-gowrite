# e03s01 — Teams CRUD

## 1. Metadata
- **Story ID:** e03s01
- **Epic:** e03 — Organizations (teams & memberships)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** low

## 2. User Story
As **an org owner**, I want to **create and manage teams**, so that **projects and memberships have a stable tenancy parent**.

## 3. Business Value
Tenancy root that projects depend on (teamInternalId); straightforward CRUD but load-bearing.

## 4. Context & Background
Part of epic **e03 — Organizations (teams & memberships)** (Teams and memberships with the ≥1-owner invariant and invite lifecycle.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- team create/get/list/update/delete
- prefs

## 6. Out of Scope
- Memberships (own story)

## 7. Dependencies
- e02s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- teams

## 11. Interface / API Surface
- `POST /v1/teams`
- `GET/PATCH/DELETE /v1/teams/:id`

## 12. Business Rules Referenced
- (no specific BR — mechanical/CRUD parity)

## 13. Paradigm Notes (PHP → Go)
Team document → typed Team struct + sqlc.

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
Feature: Teams CRUD
  Scenario: create team
    Given an owner creates a team
    Then the team is stored and the creator is seeded as owner membership

  Scenario: delete cascade guard
    Given a team with projects is deleted
    Then the delete follows legacy cascade semantics (deferred to Deletes worker)

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
