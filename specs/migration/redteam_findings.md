---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: redteam_findings
producedBy: red-team (3 adversarial agents)
hash: "sha256:479d5e9ca21c4cb589966941debe4e1727eb33afbee6d6a0f5821b09d51a5f24"
---

# Red-Team Findings — Topology Validation

> Three adversarial agents attacked the 9-context hybrid topology from distinct lenses. Verdicts:
> **Boundary critic** = SOUND-WITH-FIXES · **Strategy-fit critic** = UNDERMINES (as documented) · **Simplicity/paradigm critic** = OVER-ENGINEERED.
> Convergence: two+ agents independently flagged the `platform/ops` bucket and the routing/coupling granularity — high confidence.

## A. Boundary & ownership (boundary critic)
| # | Sev | Finding | Evidence | Fix |
|---|-----|---------|----------|-----|
| B1 | CRITICAL | `projects` collection written by **two** contexts (Organization + Projects) — split aggregate ownership | `Organization/Http/Projects/Create.php:151` & `Projects/Http/Projects/Create.php:154` | Projects owns the write-model; `/v1/organization/projects/*` is a view delegating into Projects |
| B2 | CRITICAL | VCS↔Compute bidirectional hard coupling (not messaging) | `functions`/`sites` store `installationId`/`repositoryId` FKs (`relationships.md:88-101`); VCS writes `deployments` on webhook (`VCS/Http/GitHub/Deployment.php`) | Merge VCS into Compute, **or** force that edge through outbox + remove VCS's direct writes |
| B3 | HIGH | Console mis-grouped into Projects — it's a BFF spanning Identity/Orgs/Projects | `Console/Http/{Redirects,Scopes,Assistant,Init}` | Model Console as a **BFF/edge layer above** the contexts, not inside Projects |
| B4 | HIGH | Webhooks demoted to `platform/ops` but is a `projects` aggregate child (like Keys/DevKeys/Platforms) + has state machine | `relationships.md:124` (`subQueryWebhooks`); `state-machines.md:117` | Move Webhooks **into Projects** |
| B5 | HIGH | Migrations demoted to platform but has state machine + provider aggregates | `state-machines.md:76-91`; `domain.md:86-88` | Own context (or fold into Data) with a real `domain/` |
| B6 | HIGH | Advisor demoted to platform but has Report→Insight aggregate + state machine | `relationships.md:74-79`; `state-machines.md:106-115` | Own context (or fold into Data) |
| B7 | MED | Networking(Proxy) reads Compute `deployments` at write time | `Proxy/Http/Rules/Function/Create.php:103` | Documented read-port, or denormalize readiness via event |
| B8 | ✅ | Tokens→Storage grouping **validated as sound** | `Tokens/Http/Tokens/Buckets/Files/` only | keep |

## B. Strategy-fit (strategy-fit critic)
| # | Sev | Finding | Evidence | Fix |
|---|-----|---------|----------|-----|
| S1 | CRITICAL | Routing-unit mismatch: strategy flips **22 modules**, topology groups **9 contexts** | `target_architecture.md:45,54` vs `cutover_plan.md:23-32` | Make the **context** the flip unit; co-scheduled modules flip atomically (Account+Users together) |
| S2 | CRITICAL | Cross-context writes mid-request break per-collection ownership (R-D2) | `Teams/Http/Memberships/Status/Update.php:129,139-150` writes `users`+`sessions`; `Memberships/Create.php:213` creates `users` | Co-schedule **Identity+Organizations** as one atomic flip unit |
| S3 | HIGH | `Deletes.php` worker cascades writes across all contexts; no home | `Workers/Deletes.php:685-946,1005,1695` | Cross-cutting workers migrate **last**, after every context they touch |
| S4 | HIGH | Compute scheduled to flip **before** VCS despite synchronous read of `installations` | `Functions/Http/Functions/Create.php:221,354`; `Sites/.../Create.php:162` | Reorder VCS before/with Compute; add "referenced-context migrated" go/no-go gate |
| S5 | MED | R-D2 measured at collection level; real coupling is request/feature-level | `relationships.md` (projectInternalId ubiquitous) | Derive migration order **topologically** from a dependency graph, not learning curve |
| S6 | ✅ | Data optimistic-lock transactions are **self-contained** (no cross-context atomic writes) | `Databases/.../withTransaction()` | design sound in isolation |

## C. Over-engineering / paradigm (simplicity critic)
| # | Sev | Finding | Fix |
|---|-----|---------|-----|
| O1 | CRITICAL | Hexagonal `ports/adapters` × 9 = premature abstraction; single Postgres ⇒ one impl forever (anti-Go-idiom) | Flat `internal/{ctx}/` with concrete structs; interfaces only where substitution is real (SMTP, executor); test against ephemeral Postgres |
| O2 | CRITICAL | `platform/` is a god-package rebuilding the framework we're escaping; `ops/` hides real contexts | Split into small leaf packages (errors/auth/permissions/crypto); use `pgx`/`chi` directly, no home-grown abstraction |
| O3 | HIGH | Footprint asserted, not engineered; outbox worker **+** Redis pub/sub = two overlapping dispatch mechanisms | Drop Redis if measurement allows: Postgres `LISTEN/NOTIFY` + `SKIP LOCKED` + in-memory token-bucket. One fewer daemon = real win |
| O4 | HIGH | BR-HUMANA-001 (reimplement utopia-php/database subset) = multi-month yak-shave; schema is now static | Use `sqlc` + hand-written queries; explicit `crypto.Encrypt()` at the 2-3 write sites; drop the generic filter/relationship engine |
| O5 | MED | Outbox pattern cargo-culted onto single-Postgres (outbox solves dual-write across stores; there's one store) | Single `jobs` table consumed via `SELECT ... FOR UPDATE SKIP LOCKED` |

## Convergence (multiple agents, same flaw)
- **`platform/ops` hides real bounded contexts** — flagged by BOTH boundary (B4/B5/B6) and simplicity (O2). High confidence.
- **Granularity/coupling** — boundary (B1/B2) and strategy (S1/S2/S4) both show module fusion + cross-context writes undermine clean per-unit migration.

---

# Red-Team Findings — Console/Screen Decision (2026-08-11)

> Two adversarial agents validated the "preserve SPA + contract-parity" decision.
> **Hosting-feasibility** = FEASIBLE-WITH-CHANGES · **Contract-parity** = INSUFFICIENT-as-stated (oracle was hollow).
> Outcome: mode (**preserve SPA**) is CORRECT and kept; the **parity oracle** was wrong and is corrected.

## Hosting (FEASIBLE-WITH-CHANGES)
- ✅ Console is a genuine static SPA (`adapter-static`, `ssr=false`, 0 server routes) → `go:embed` works, **no Node runtime**.
- Fixes: serve under base `/console` + `301 /→/console`; `try_files` SPA fallback; br/gz negotiation (`precompress`); 3 cache tiers.
- **Corrected error**: `PUBLIC_APPWRITE_ENDPOINT` is build-time, not Go-runtime — leave unset, SDK derives from `window.location` (`sdk.ts:43`).

## Contract-parity (INSUFFICIENT-as-stated → corrected)
- **CRITICAL-0 (verified)**: `specs/openapi/appwrite.yaml` has **0 `responses:` / 0 `requestBody:`** — routing only, no payloads. Console reads `.total` 273×. "Contract parity" ≠ payload parity.
- Non-REST UI dependencies (→ DEV-002..009): Realtime WS (529 files), assistant **text-stream** (spec wrongly said SSE), error `type`/`code`/`message` control flow, cookie-session + `X-Fallback-Cookies`, console-only headers (impersonation/`X-Appwrite-Mode`, fingerprint, organization), `Date` header, multi-region topology.
- **Fix**: parity oracle = response goldens (from `specs/database/`) + WS protocol + stream framing + error catalog + E2E console suite — NOT the hollow OpenAPI.

## Convergent outcome
Both agents affirm **preserve the SPA**; both expose that the naive assumptions (stock FileServer, OpenAPI-as-oracle) were wrong. Decision doc + deviation log corrected accordingly. **Discovery gap logged**: the Writer's OpenAPI lacks response schemas.
