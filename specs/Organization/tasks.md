# Tasks: Organization

> Unit: `Organization`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-ORG-01 — Module bootstrap + Init hook
- **Behavior source:** `Modules/Organization/Module.php`, `Services/Http.php`, `Http/Init.php`
- **Description:** Register HTTP service + Init hook validating `$team` presence.
- **Done when:** Empty team blocked on all endpoints.
- **Priority:** Must · Confidence: 🟢

### T-ORG-02 — Create project with provisioning
- **Behavior source:** `Http/Projects/Create.php`
- **Description:** Reserved id rejection; region filtering; DSN selection (override/shared/random); shared-table join vs per-project namespace + audit + collections bootstrap; auth seeding; afterProjectCreation hook.
- **Done when:** All provisioning modes verified.
- **Priority:** Must · Confidence: 🟢

### T-ORG-03 — List team projects
- **Behavior source:** `Http/Projects/XList.php`
- **Description:** Team-scoped list with select queries + skipFilters optimization.
- **Done when:** List returns only team projects; select works.
- **Priority:** Must · Confidence: 🟢

### T-ORG-04 — Get/Update project
- **Behavior source:** `Http/Projects/{Get,Update}.php`
- **Description:** Ownership check (teamInternalId) before acting; sparse rename.
- **Done when:** Cross-team → PROJECT_NOT_FOUND.
- **Priority:** Must · Confidence: 🟢

### T-ORG-05 — Delete project + async cleanup
- **Behavior source:** `Http/Projects/Delete.php`
- **Description:** Ownership check; delete with authorization skip; enqueue DELETE_TYPE_DOCUMENT.
- **Done when:** Cleanup message enqueued.
- **Priority:** Must · Confidence: 🟢

### T-ORG-06 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Organization/Http/**`
- **Description:** OpenAPI 3.1 for organization endpoints.
- **Done when:** `openapi/organization.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-ORG-07 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Organization/**`
- **Description:** Cover create/list/get/update/delete, reserved id, cross-team, shared tables.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
