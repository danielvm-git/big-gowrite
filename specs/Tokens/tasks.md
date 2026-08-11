# Tasks: Tokens

> Unit: `Tokens`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-TKN-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Tokens/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 5 endpoint actions.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-TKN-02 — Shared bucket/file resolution guard
- **Behavior source:** `Http/Tokens/Buckets/Files/Action.php`
- **Description:** `getFileAndBucket()` — enabled-state check, read permission, file existence.
- **Done when:** Guard verified for all operations.
- **Priority:** Must · Confidence: 🟢

### T-TKN-03 — Create file token
- **Behavior source:** `Http/Tokens/Buckets/Files/Create.php`
- **Description:** 128-byte secret generation + encryption, `resourceId`/`resourceInternalId` mapping, permission gates (bucket vs file-level with `fileSecurity`), rate limit.
- **Done when:** Creation verified across permission variants; secret stored encrypted.
- **Priority:** Must · Confidence: 🟢

### T-TKN-04 — List tokens for a file
- **Behavior source:** `Http/Tokens/Buckets/Files/XList.php`
- **Description:** Filter `resourceType = files` + compound key.
- **Done when:** Filter verified.
- **Priority:** Must · Confidence: 🟢

### T-TKN-05 — Get / Update / Delete token
- **Behavior source:** `Http/Tokens/{Get,Update,Delete}.php`
- **Description:** Get by ID; expiry future-only update; delete revokes access.
- **Done when:** Lifecycle verified.
- **Priority:** Must · Confidence: 🟢

### T-TKN-06 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Tokens/Http/**`
- **Description:** OpenAPI 3.1 for tokens endpoints.
- **Done when:** `openapi/tokens.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-TKN-07 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Tokens/**`
- **Description:** Cover create/list/get/update/delete, permission variants, expiry validation, disabled-bucket guard.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
