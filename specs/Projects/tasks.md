# Tasks: Projects

> Unit: `Projects`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-PRJS-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Projects/Module.php`, `Services/Http.php`, `Http/Projects/Action.php`
- **Description:** Register HTTP service with 14 actions and base action class.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-PRJS-02 — Create project with provisioning
- **Behavior source:** `Http/Projects/Create.php`
- **Description:** Reserved id; team existence; region/DSN selection; shared vs dedicated provisioning; collections + audit schema; afterProjectCreation hook.
- **Done when:** All provisioning modes verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJS-03 — Update + list projects
- **Behavior source:** `Http/Projects/{Update,XList}.php`
- **Description:** Sparse metadata update; paginated list with team filter.
- **Done when:** Update/list verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJS-04 — Team transfer
- **Behavior source:** `Http/Projects/Team/Update.php`
- **Description:** Ownership change + permission propagation to installations/repositories/VCS comments.
- **Done when:** Transfer propagates correctly.
- **Priority:** Must · Confidence: 🟢

### T-PRJS-05 — Dev keys (tombstoned)
- **Behavior source:** `Http/DevKeys/**`
- **Description:** createDevKey throws DEV_KEY_GONE; CRUD scoped by projectInternalId.
- **Done when:** Tombstone + scoping verified.
- **Priority:** Should · Confidence: 🟢

### T-PRJS-06 — Schedules CRUD
- **Behavior source:** `Http/Schedules/**`
- **Description:** Create with resource validation; get/list with dual-id filtering; event emit.
- **Done when:** Schedule lifecycle verified.
- **Priority:** Must · Confidence: 🟢

### T-PRJS-07 — Onboarding stages
- **Behavior source:** `Http/Stages/{XList,Update}.php`
- **Description:** List config-driven stages; update skip with actor resolution + completed guard.
- **Done when:** Stage skip rules verified.
- **Priority:** Should · Confidence: 🟢

### T-PRJS-08 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Projects/Http/**`
- **Description:** OpenAPI 3.1 for projects endpoints.
- **Done when:** `openapi/projects.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-PRJS-09 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Projects/**`
- **Description:** Cover create/update/list/team transfer/dev-keys/schedules/stages.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
