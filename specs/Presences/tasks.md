# Tasks: Presences

> Unit: `Presences`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-PRS-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Presences/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 5 endpoints and presences scopes.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-PRS-02 — Upsert presence
- **Behavior source:** `HTTP/Upsert.php`
- **Description:** Identity resolution (session vs API key); expiry defaults/bounds; source detection; permissions; event + metric.
- **Done when:** Upsert verified across identity modes.
- **Priority:** Must · Confidence: 🟢

### T-PRS-03 — Get presence
- **Behavior source:** `HTTP/Get.php`
- **Description:** Fetch with expiry filtering → PRESENCE_NOT_FOUND.
- **Done when:** Expired hidden.
- **Priority:** Must · Confidence: 🟢

### T-PRS-04 — List presences with TTL cache
- **Behavior source:** `HTTP/XList.php`
- **Description:** Cache key per project/roles/query; X-Appwrite-Cache header; server-side expiresAt filter.
- **Done when:** Cache hit/miss verified.
- **Priority:** Must · Confidence: 🟢

### T-PRS-05 — Update presence
- **Behavior source:** `HTTP/Update.php`
- **Description:** Sparse update; purge flag → cache invalidation; event.
- **Done when:** Update + purge verified.
- **Priority:** Must · Confidence: 🟢

### T-PRS-06 — Delete presence
- **Behavior source:** `HTTP/Delete.php`
- **Description:** Delete + purge + metric decrement + event.
- **Done when:** Delete path verified.
- **Priority:** Must · Confidence: 🟢

### T-PRS-07 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Presences/HTTP/**`
- **Description:** OpenAPI 3.1 for presences endpoints.
- **Done when:** `openapi/presences.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-PRS-08 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Presences/**`
- **Description:** Cover upsert/get/list/update/delete, expiry, cache header, identity scoping.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
