# Tasks: Webhooks

> Unit: `Webhooks`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-WHK-01 — Module bootstrap + Init hook + route registration
- **Behavior source:** `Modules/Webhooks/Module.php`, `Services/Http.php`, `Http/Init.php`
- **Description:** Register HTTP service with 6 endpoints + Init hook blocking the `console` project.
- **Done when:** Console block verified; route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-WHK-02 — Create webhook
- **Behavior source:** `Http/Webhooks/Create.php`
- **Description:** PublicDomain URL validation; secret auto-gen (128-char hex); projectInternalId scoping; cache purge.
- **Done when:** Creation verified (with and without explicit secret).
- **Priority:** Must · Confidence: 🟢

### T-WHK-03 — Get / List webhook
- **Behavior source:** `Http/Webhooks/{Get,XList}.php`
- **Description:** Project-scoped reads; `signatureKey` stripped from responses.
- **Done when:** Field-stripping verified.
- **Priority:** Must · Confidence: 🟢

### T-WHK-04 — Update webhook (full replace)
- **Behavior source:** `Http/Webhooks/Update.php`
- **Description:** Full replace; re-enable resets `attempts`; cache purge.
- **Done when:** Replace + attempts reset verified.
- **Priority:** Must · Confidence: 🟢

### T-WHK-05 — Delete webhook
- **Behavior source:** `Http/Webhooks/Delete.php`
- **Description:** Removes configuration; cache purge.
- **Done when:** Deletion verified.
- **Priority:** Must · Confidence: 🟢

### T-WHK-06 — Rotate secret
- **Behavior source:** `Http/Webhooks/Secret/Update.php`
- **Description:** Explicit or auto-generated secret; returned once; cache purge.
- **Done when:** Rotation verified.
- **Priority:** Must · Confidence: 🟢

### T-WHK-07 — Legacy path aliases
- **Behavior source:** `Services/Http.php`
- **Description:** `/v1/projects/:projectId/webhooks[...]` aliases via `httpAlias`.
- **Done when:** Aliases verified.
- **Priority:** Should · Confidence: 🟢

### T-WHK-08 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Webhooks/Http/**`
- **Description:** OpenAPI 3.1 for webhooks endpoints.
- **Done when:** `openapi/webhooks.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-WHK-09 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Webhooks/**`
- **Description:** Cover CRUD, secret rotation + stripping, console-project block, legacy aliases.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
