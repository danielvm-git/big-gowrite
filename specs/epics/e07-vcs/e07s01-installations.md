# e07s01 — VCS installations & repository linkage

## 1. Metadata
- **Story ID:** e07s01
- **Epic:** e07 — VCS (installations & repo events)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** high

## 2. User Story
As **a project developer**, I want to **connect a Git provider and link repositories**, so that **functions/sites can reference a linked repo, as in legacy**.

## 3. Business Value
Stores GitHub/GitLab installations (encrypted credentials) that Compute reads to build from source.

## 4. Context & Background
Part of epic **e07 — VCS (installations & repo events)** (GitHub/GitLab installations, repo linkage, and repository.pushed event publishing.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- provider OAuth/app installation
- installation storage (encrypted)
- repo listing/linkage

## 6. Out of Scope
- Build execution (compute epic)

## 7. Dependencies
- e04s01
- e01s02

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- installations
- repositories

## 11. Interface / API Surface
- `VCS provider callbacks`
- `GET repos`

## 12. Business Rules Referenced
- BR-MIGRAR-006

## 13. Paradigm Notes (PHP → Go)
Installation document → typed struct; provider tokens via crypto.Encrypt.

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
Feature: VCS installations & repository linkage
  Scenario: install
    Given a developer installs the GitHub app
    Then an installation with encrypted credentials is stored

  Scenario: link repo
    Given a developer links a repository to a project
    Then the linkage is persisted and listable

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
