# e13s02 — Serve SvelteKit SPA via go:embed

## 1. Metadata
- **Story ID:** e13s02
- **Epic:** e13 — Console BFF + SPA hosting
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 3
- **Risk:** P1  ·  **Security surface:** medium

## 2. User Story
As **a console user**, I want to **load the console from the Go server at /console**, so that **the SPA behaves identically, deriving the API endpoint from window.location**.

## 3. Business Value
Serves the preserved console SPA from Go with the 5 hosting requirements — satisfying UI parity (BR-HUMANA-004) without a rewrite.

## 4. Context & Background
Part of epic **e13 — Console BFF + SPA hosting** (BFF composition over identity+orgs+projects and serving the preserved SPA via go:embed.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- go:embed of built SPA
- base /console
- SPA fallback (try_files)
- br/gz negotiation
- cache tiers
- endpoint via window.location (leave PUBLIC_APPWRITE_ENDPOINT unset)

## 6. Out of Scope
- Multi-region hosting (DEV-009, deferred)

## 7. Dependencies
- e13s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- (static assets)

## 11. Interface / API Surface
- `GET /console/*`

## 12. Business Rules Referenced
- (no specific BR — mechanical/CRUD parity)

## 13. Paradigm Notes (PHP → Go)
Nginx/SSR hosting → go:embed static serving with SPA fallback + compression negotiation.

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
Feature: Serve SvelteKit SPA via go:embed
  Scenario: spa fallback
    Given a deep console route is requested directly
    Then index.html is served and the SPA routes client-side

  Scenario: same-origin endpoint
    Given the SPA loads with PUBLIC_APPWRITE_ENDPOINT unset
    Then the SDK derives the API base from window.location

  Scenario: compression
    Given a client advertising br is served
    Then a brotli-encoded asset is returned with correct cache headers

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
