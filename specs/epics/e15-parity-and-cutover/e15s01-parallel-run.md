# e15s01 — Parallel Run harness + response goldens

## 1. Metadata
- **Story ID:** e15s01
- **Epic:** e15 — Parity harness & cutover
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** high

## 2. User Story
As **the migration operator**, I want to **mirror requests to Go and PHP and compare responses**, so that **no context flips until it proves parity against captured goldens**.

## 3. Business Value
Implements the parity gate: mirror traffic to Go+PHP, compare outputs, and assert <0.01% divergence — the precondition for every flip.

## 4. Context & Background
Part of epic **e15 — Parity harness & cutover** (Parallel Run harness, contract tests (response goldens, WS, stream, error catalog), and the per-context flip/canary/decommission runbook.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- traffic mirroring (shadow)
- response golden capture from legacy
- divergence report
- acceptance <0.01% over window

## 6. Out of Scope
- The OpenAPI-as-oracle approach (it has 0 response schemas — use goldens)

## 7. Dependencies
- e01s05

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- (test infra)

## 11. Interface / API Surface
- `shadow proxy + comparator`

## 12. Business Rules Referenced
- (no specific BR — mechanical/CRUD parity)

## 13. Paradigm Notes (PHP → Go)
Parity oracle = response goldens + protocol, not the hollow OpenAPI.

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
Feature: Parallel Run harness + response goldens
  Scenario: shadow compare
    Given a request is mirrored to Go and PHP
    Then responses are compared and divergences are logged with field-level diffs

  Scenario: gate pass
    Given divergence stays under 0.01% over the observation window
    Then the context is marked flip-eligible

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
