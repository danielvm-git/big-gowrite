# Tasks: Teams

> Unit: `Teams`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-TMS-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Teams/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 13 endpoints.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-TMS-02 — Team CRUD + name update
- **Behavior source:** `Http/Teams/{Create,Get,XList,Delete}.php`, `Http/Teams/Name/Update.php`
- **Description:** Creator ownership in user mode; searchable list; cascade delete with console project handling.
- **Done when:** Team lifecycle verified.
- **Priority:** Must · Confidence: 🟢

### T-TMS-03 — Team preferences
- **Behavior source:** `Http/Preferences/{Get,Update}.php`
- **Description:** Get/replace arbitrary JSON prefs.
- **Done when:** Prefs verified.
- **Priority:** Must · Confidence: 🟢

### T-TMS-04 — Membership create (invite/direct)
- **Behavior source:** `Http/Memberships/Create.php`
- **Description:** userId/email/phone resolution + consistency; stub user creation with email validation; direct vs invite modes; SMTP/SMS gates.
- **Done when:** All invite paths verified.
- **Priority:** Must · Confidence: 🟢

### T-TMS-05 — Membership get/list with privacy
- **Behavior source:** `Http/Memberships/{Get,XList}.php`
- **Description:** Privacy settings on exposed fields; privileged full view.
- **Done when:** Privacy gating verified.
- **Priority:** Must · Confidence: 🟢

### T-TMS-06 — Membership role update
- **Behavior source:** `Http/Memberships/Update.php`
- **Description:** Roles replaced; owner downgrade guard.
- **Done when:** Role updates verified.
- **Priority:** Must · Confidence: 🟢

### T-TMS-07 — Membership delete
- **Behavior source:** `Http/Memberships/Delete.php`
- **Description:** Last-owner guard; primary promotion; total decrement.
- **Done when:** Delete + promotion verified.
- **Priority:** Must · Confidence: 🟢

### T-TMS-08 — Status confirmation
- **Behavior source:** `Http/Memberships/Status/Update.php`
- **Description:** Secret validation; session creation for unauthenticated invitee; total increment.
- **Done when:** Confirm flow verified.
- **Priority:** Must · Confidence: 🟢

### T-TMS-09 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Teams/Http/**`
- **Description:** OpenAPI 3.1 for teams endpoints.
- **Done when:** `openapi/teams.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-TMS-10 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Teams/**`
- **Description:** Cover teams CRUD, invites, confirmation, roles, privacy, owner guards.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
