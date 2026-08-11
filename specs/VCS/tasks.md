# Tasks: VCS

> Unit: `VCS`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-VCS-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/VCS/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with authorize/callback/events/installation/repository endpoints.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-VCS-02 — OAuth2 authorize + callback (Gitea, GitLab)
- **Behavior source:** `Http/Authorize/Base.php`, `Http/Callback/Base.php`, `Http/Gitea/**`, `Http/Gitlab/**`
- **Description:** Shared redirect + callback logic; installation upsert in dbForPlatform.
- **Done when:** Both providers verified end-to-end.
- **Priority:** Must · Confidence: 🟢

### T-VCS-03 — GitHub App authorize + callback
- **Behavior source:** `Http/GitHub/Authorize/Get.php`, `Http/GitHub/Callback/Get.php`
- **Description:** App-install redirect; HMAC state verification on public callback.
- **Done when:** Flow verified with state tamper test.
- **Priority:** Must · Confidence: 🟢

### T-VCS-04 — Webhook ingestion + deployment fan-out
- **Behavior source:** `Http/GitHub/Events/Create.php`, `Http/Gitea/Events/Create.php`, `Http/Gitlab/Events/Create.php`, `Http/GitHub/Deployment.php`
- **Description:** Signature validation; bot-email ignore; skip patterns/globs; external-PR authorization; preview rules; optimistic-lock comments.
- **Done when:** Push, PR, and installation webhooks verified across providers.
- **Priority:** Must · Confidence: 🟢

### T-VCS-05 — Installation get/list/delete
- **Behavior source:** `Http/Installations/{Get,XList,Delete}.php`
- **Description:** Platform-DB CRUD; delete enqueues cleanup; uninstall triggers region-scoped cleanup.
- **Done when:** Lifecycle verified.
- **Priority:** Must · Confidence: 🟢

### T-VCS-06 — Repository browsing + detection
- **Behavior source:** `Http/Installations/Repositories/{Create,Get,XList}.php`, `.../Branches/XList.php`, `.../Contents/Get.php`, `.../Detections/Create.php`
- **Description:** Adapter-based provider calls; token refresh; Swoole coroutines for list/detection.
- **Done when:** All browse endpoints verified.
- **Priority:** Must · Confidence: 🟢

### T-VCS-07 — External deployment authorization
- **Behavior source:** `Http/GitHub/Authorize/External/Update.php`
- **Description:** PATCH grants fork-PR contributor deployment; updates pending status.
- **Done when:** Authorization + deployment verified.
- **Priority:** Must · Confidence: 🟢

### T-VCS-08 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/VCS/Http/**`
- **Description:** OpenAPI 3.1 for VCS endpoints.
- **Done when:** `openapi/vcs.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-VCS-09 — E2E test suite
- **Behavior source:** `tests/e2e/Services/VCS/**`
- **Description:** Cover authorize/callback, webhooks (push/PR/uninstall), skip patterns, fork authorization, repository browsing.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
