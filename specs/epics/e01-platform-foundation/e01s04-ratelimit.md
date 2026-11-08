# e01s04 — In-memory token-bucket rate limiting

## 1. Metadata
- **Story ID:** e01s04
- **Epic:** e01 — Platform Foundation
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **the request middleware**, I want to **apply per-route rate limits keyed by ip/userId/url/email**, so that **abusive traffic is throttled with the same limits and 429 responses as legacy**.

## 3. Business Value
Preserves abuse protection (BR-MIGRAR-009) without Redis, honoring the legacy per-route limits and keys.

## 4. Context & Background
Part of epic **e01 — Platform Foundation** (Shared leaf packages (errors, auth, permissions, crypto, ratelimit), base Postgres schema, and the jobs table with SKIP-LOCKED dispatch + LISTEN/NOTIFY realtime.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- token-bucket per instance
- route limit labels (10/15/50/60/100/120; windows 3600s/60s)
- keys ip/userId/url/email
- 429 + Retry-After

## 6. Out of Scope
- Distributed/shared rate limiting across instances (single-VPS scope)

## 7. Dependencies
- e01s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- (in-memory only)

## 11. Interface / API Surface
- `ratelimit.Allow(key string, limit Limit) (bool, retryAfter time.Duration)`

## 12. Business Rules Referenced
- BR-MIGRAR-009

## 13. Paradigm Notes (PHP → Go)
utopia-php/abuse (Redis-backed) → in-memory token bucket per instance; acceptable because target is a single VPS.

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
Feature: In-memory token-bucket rate limiting
  Scenario: under limit
    Given requests below the route limit arrive
    Then all are allowed

  Scenario: over limit
    Given requests exceed the route limit within the window
    Then excess requests get 429 with a Retry-After header

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
