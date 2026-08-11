# Tasks: Console

> Unit: `Console`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-CSL-01 — Register Console module + init hooks
- **Behavior source:** `Modules/Console/Module.php`, `Http/Init/{API,Web}.php`
- **Description:** Register HTTP service; implement API (console-project guard) and Web (security headers) init hooks with correct group assignments.
- **Done when:** Console group blocked outside console project; web group gets headers.
- **Priority:** Must · Confidence: 🟢

### T-CSL-02 — Register 15 HTTP routes
- **Behavior source:** `Modules/Console/Services/Http.php`
- **Description:** Register 7 `/v1/console/**` routes + 8 legacy redirect routes with correct scopes and action IDs.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-CSL-03 — getVariables
- **Behavior source:** `Modules/Console/Http/Variables/Get.php`
- **Description:** Collect env config + DB adapter capabilities + feature flags.
- **Done when:** Payload matches legacy fields; UI adaptation data present.
- **Priority:** Must · Confidence: 🟢

### T-CSL-04 — getConsoleEmailTemplate
- **Behavior source:** `Modules/Console/Http/Templates/Email/Get.php`
- **Description:** Render built-in default template by type+locale with English fallback.
- **Done when:** Built-in template returned regardless of project customization.
- **Priority:** Must · Confidence: 🟢

### T-CSL-05 — listOAuth2Providers
- **Behavior source:** `Modules/Console/Http/OAuth2Providers/XList.php`
- **Description:** List configured enabled providers; exclude mock/disabled; include sandbox.
- **Done when:** Provider list matches config.
- **Priority:** Must · Confidence: 🟢

### T-CSL-06 — Scope lists (project + organization)
- **Behavior source:** `Modules/Console/Http/Scopes/{Project,Organization}/XList.php`
- **Description:** Static scope lists from config params.
- **Done when:** Lists match `projectScopes`/`organizationScopes`.
- **Priority:** Must · Confidence: 🟢

### T-CSL-07 — Assistant SSE proxy
- **Behavior source:** `Modules/Console/Http/Assistant/Create.php`
- **Description:** Rate limit 15/user; proxy to appwrite-assistant; stream SSE via chunk.
- **Done when:** Streamed responses; rate limit enforced.
- **Priority:** Should · Confidence: 🟢

### T-CSL-08 — getResource domain availability
- **Behavior source:** `Modules/Console/Http/Resources/Get.php`
- **Description:** Structure validation, deny list, reserved hostnames, uniqueness in dbForPlatform.rules.
- **Done when:** All 4 checks enforced.
- **Priority:** Must · Confidence: 🟢

### T-CSL-09 — Legacy redirects
- **Behavior source:** `Modules/Console/Http/Redirects/**`
- **Description:** 8 redirect actions preserving path/query/fragment to /console.
- **Done when:** All redirects verified by E2E.
- **Priority:** Must · Confidence: 🟢

### T-CSL-10 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Console/Http/**`
- **Description:** OpenAPI 3.1 for console endpoints.
- **Done when:** `openapi/console.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-CSL-11 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Console/**`
- **Description:** Cover variables, email template, providers, scopes, assistant, resources, redirects, console guard, web headers.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
