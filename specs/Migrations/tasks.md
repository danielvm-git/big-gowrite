# Tasks: Migrations

> Unit: `Migrations`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-MIG-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Migrations/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 16 endpoints; all creation endpoints admin-only.
- **Done when:** Route table matches mapping; admin gate enforced.
- **Priority:** Must · Confidence: 🟢

### T-MIG-02 — Migration CRUD + retry
- **Behavior source:** `Http/Migrations/{Get,XList,Update,Delete}.php`
- **Description:** get/list/delete; retry only when failed; re-enqueue + emit retry event.
- **Done when:** Retry logic verified.
- **Priority:** Must · Confidence: 🟢

### T-MIG-03 — Appwrite migration + report
- **Behavior source:** `Http/Migrations/Appwrite/{Create,Report/Get}.php`
- **Description:** Appwrite-to-Appwrite with onDuplicate policy (fail/skip/overwrite); dry-run report.
- **Done when:** Policy honored; report returns preview.
- **Priority:** Must · Confidence: 🟢

### T-MIG-04 — Firebase migration + report
- **Behavior source:** `Http/Migrations/Firebase/{Create,Report/Get}.php`
- **Description:** Service account JSON validation (project_id, client_email, private_key); report.
- **Done when:** Invalid service account rejected.
- **Priority:** Must · Confidence: 🟢

### T-MIG-05 — Supabase migration + report
- **Behavior source:** `Http/Migrations/Supabase/{Create,Report/Get}.php`
- **Description:** DB credentials stored in credentials field; report.
- **Done when:** Credentials persisted securely.
- **Priority:** Must · Confidence: 🟢

### T-MIG-06 — NHost migration + report
- **Behavior source:** `Http/Migrations/NHost/{Create,Report/Get}.php`
- **Description:** Same as Supabase pattern.
- **Done when:** Migration + report work.
- **Priority:** Must · Confidence: 🟢

### T-MIG-07 — CSV export/import
- **Behavior source:** `Http/Migrations/CSV/{Exports,Imports}/Create.php`
- **Description:** Export with delimiter/enclosure/escape/header options to internal bucket; import with decrypt/decompress + CSV_ALLOWED_DATABASE_TYPES gate + backward-compat alias.
- **Done when:** Exports/imports verified; type gate enforced.
- **Priority:** Must · Confidence: 🟢

### T-MIG-08 — JSON export/import
- **Behavior source:** `Http/Migrations/JSON/{Exports,Imports}/Create.php`
- **Description:** All DB types supported.
- **Done when:** JSON round-trip verified.
- **Priority:** Must · Confidence: 🟢

### T-MIG-09 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Migrations/Http/**`
- **Description:** OpenAPI 3.1 for migrations endpoints.
- **Done when:** `openapi/migrations.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-MIG-10 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Migrations/**`
- **Description:** Cover provider migrations, retries, reports, CSV/JSON ops.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
