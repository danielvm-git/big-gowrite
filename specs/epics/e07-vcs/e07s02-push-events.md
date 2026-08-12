# e07s02 — repository.pushed event publishing via jobs

## 1. Metadata
- **Story ID:** e07s02
- **Epic:** e07 — VCS (installations & repo events)
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **the VCS webhook handler**, I want to **publish a repository.pushed event on push**, so that **Compute consumes it to trigger builds, with no direct cross-context write**.

## 3. Business Value
Decouples VCS→Compute (red-team B2): VCS publishes an event to the jobs table instead of writing deployments directly.

## 4. Context & Background
Part of epic **e07 — VCS (installations & repo events)** (GitHub/GitLab installations, repo linkage, and repository.pushed event publishing.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- provider webhook receipt
- publish event:repository.pushed to jobs
- signature verification

## 6. Out of Scope
- Deployment creation (compute consumes the event)

## 7. Dependencies
- e07s01
- e01s05

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- jobs

## 11. Interface / API Surface
- `VCS webhook endpoint`

## 12. Business Rules Referenced
- BR-MIGRAR-014
- BR-MIGRAR-015

## 13. Paradigm Notes (PHP → Go)
Direct deployments write → jobs table event; Compute is the consumer (no shared table write).

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
Feature: repository.pushed event publishing via jobs
  Scenario: publish on push
    Given a valid push webhook arrives
    Then an event:repository.pushed job is enqueued with the repo/commit payload

  Scenario: bad signature
    Given a push webhook with an invalid signature arrives
    Then it is rejected and no job is enqueued

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
