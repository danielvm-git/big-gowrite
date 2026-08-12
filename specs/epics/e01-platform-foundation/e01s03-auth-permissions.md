# e01s03 — Auth, console roles + document ACL

## 1. Metadata
- **Story ID:** e01s03
- **Epic:** e01 — Platform Foundation
- **Status:** todo
- **Maturity:** 3 (Countable)
- **BCPs:** 5
- **Risk:** P0  ·  **Security surface:** high

## 2. User Story
As **the request middleware**, I want to **resolve the caller's identity and scopes and enforce document ACL**, so that **only authorized principals reach each resource, identical to legacy behavior**.

## 3. Business Value
Preserves the security model (BR-MIGRAR-007/008): session/JWT verification, console role scopes, and document-level ACL — the gate every request passes through.

## 4. Context & Background
Part of epic **e01 — Platform Foundation** (Shared leaf packages (errors, auth, permissions, crypto, ratelimit), base Postgres schema, and the jobs table with SKIP-LOCKED dispatch + LISTEN/NOTIFY realtime.). Migration follows Strangler Fig + Parallel Run; this story is delivered as a slice of its context and only flips live after the parity gate (epic e15) passes.

## 5. In Scope
- session + JWT verification
- console roles (guests/users/admin/developer/owner/keys) as typed scope sets
- route scope labels {resource}.{read|write}
- document ACL Role::any/users/guests/user/team

## 6. Out of Scope
- Rate limiting (own story)
- MFA challenge flow (identity epic)

## 7. Dependencies
- e01s01

## 8. Assumptions
- The legacy PHP path remains available for Parallel Run comparison until this context flips.
- PostgreSQL (with pgvector) is provisioned per e01.

## 9. Constraints
- Go, flat packages, sqlc + hand SQL; PostgreSQL only (no Redis).
- Solo-dev pace; single-VPS target (not production-critical under load).

## 10. Data & Entities
- users
- sessions
- keys
- document _permissions

## 11. Interface / API Surface
- `auth.Resolve(r *http.Request) (Principal, error)`
- `permissions.Can(p Principal, scope string) bool`
- `permissions.CheckACL(p Principal, doc Document, op Op) error`

## 12. Business Rules Referenced
- BR-MIGRAR-007
- BR-MIGRAR-008

## 13. Paradigm Notes (PHP → Go)
Injected $user/$project + Action scope labels → explicit middleware resolving a Principal struct; roles.php sets become typed scope sets (admin≡developer, owner=member∪admins).

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
Feature: Auth, console roles + document ACL
  Scenario: scope granted
    Given an owner calls a route labelled teams.write
    Then permissions.Can returns true and the request proceeds

  Scenario: scope denied
    Given a member calls a route labelled projects.write reserved for developer
    Then the request is rejected with 401/403 matching legacy

  Scenario: document ACL
    Given a user without a matching read Role queries a document
    Then the document is excluded from results, identical to legacy ACL filtering

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
