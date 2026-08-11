# Tasks: Project

> Unit: `Project`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-PRJ-01 — Module bootstrap + Init hook
- **Behavior source:** `Modules/Project/Module.php`, `Services/Http.php`, `Http/Init.php`
- **Description:** Register HTTP service (60+ actions); Init guard blocks console/missing projects.
- **Done when:** Guard enforced; routes registered.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-02 — Project get/delete + labels
- **Behavior source:** `Http/Project/{Get,Delete}.php`, `Labels/Update.php`
- **Description:** Get project; delete + enqueue DELETE_TYPE_DOCUMENT; labels replace.
- **Done when:** Delete cleanup enqueued.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-03 — Protocol/service/auth-method toggles
- **Behavior source:** `Http/Project/{Protocols,Services,AuthMethods}/Update.php`
- **Description:** Enable/disable by id with events.
- **Done when:** Toggles verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-04 — SMTP config + test email
- **Behavior source:** `Http/Project/SMTP/{Update,Tests/Create}.php`
- **Description:** Connection validation on enable; auto-enable path; test email via publisherForMails.
- **Done when:** Validation + test flow verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-05 — Email templates
- **Behavior source:** `Http/Project/Templates/Email/**`
- **Description:** List/get/update locale-specific templates; SMTP gate on update.
- **Done when:** Template CRUD + gate verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-06 — Project variables CRUD
- **Behavior source:** `Http/Project/Variables/**`
- **Description:** CRUD with secret flag; live=false propagation; secret unset guard.
- **Done when:** Secret rules enforced.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-07 — API keys (standard + ephemeral)
- **Behavior source:** `Http/Project/Keys/**`
- **Description:** CRUD with ownership; ephemeral JWT ≤3600s non-persisted; cache purge.
- **Done when:** Both key types verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-08 — Platforms CRUD (5 types)
- **Behavior source:** `Http/Project/Platforms/**`
- **Description:** Per-type create/update paths; legacy type mapping; delete.
- **Done when:** All platform types verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-09 — Mock phones
- **Behavior source:** `Http/Project/MockPhone/**`
- **Description:** CRUD with APP_LIMIT_COUNT cap + duplicate rejection.
- **Done when:** Cap + dedupe enforced.
- **Priority:** Should · Confidence: 🟢

### T-PRJ-10 — Policies (12)
- **Behavior source:** `Http/Project/Policies/**`
- **Description:** Virtual docs from auths.*; list/get/update each policy.
- **Done when:** All 12 policies verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-11 — OAuth2 providers (38)
- **Behavior source:** `Http/Project/OAuth2/**`
- **Description:** List/get/update with credential verify; secret masked.
- **Done when:** Provider updates verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJ-12 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Project/Http/**`
- **Description:** OpenAPI 3.1 for project endpoints.
- **Done when:** `openapi/project.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-PRJ-13 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Project/**`
- **Description:** Cover project config, keys, platforms, policies, OAuth2, SMTP, variables.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
