# Tasks: Notifications

> Unit: `Notifications`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-NTF-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Notifications/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 3 endpoints and scopes (account, account, public).
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-NTF-02 — List console notifications
- **Behavior source:** `Http/Notifications/XList.php`
- **Description:** Console-only guard; server-side filters (channel, resourceType, resourceId); pagination.
- **Done when:** Non-console project → USER_UNAUTHORIZED; filters applied.
- **Priority:** Must · Confidence: 🟢

### T-NTF-03 — Update notification (mark read)
- **Behavior source:** `Http/Notifications/Update.php`
- **Description:** Sparse read update; 404 vs unauthorized distinction via skip query.
- **Done when:** Read state updated; errors correct.
- **Priority:** Must · Confidence: 🟢

### T-NTF-04 — Logo + email open tracking
- **Behavior source:** `Http/Notifications/Logos/Appwrite/Get.php`
- **Description:** SVG by theme with no-store; JWT decode/validate; read + firstSeen/lastSeen update; audit enqueue; silent error handling.
- **Done when:** Tracking updates notifications; invalid JWT still returns SVG.
- **Priority:** Must · Confidence: 🟢

### T-NTF-05 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Notifications/Http/**`
- **Description:** OpenAPI 3.1 for notifications endpoints.
- **Done when:** `openapi/notifications.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-NTF-06 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Notifications/**`
- **Description:** Cover list/update console gating, tracking JWT valid/invalid, theme variants.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
