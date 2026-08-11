# Tasks: Sites

> Unit: `Sites`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-SIT-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Sites/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 26 endpoints.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-SIT-02 — Site CRUD + VCS linkage
- **Behavior source:** `Http/Sites/{Create,Get,XList,Update,Delete}.php`
- **Description:** Framework/adapter validation; VCS connect/disconnect (repositories + webhooks); live=false logic; cascading delete.
- **Done when:** VCS flows + live flag verified.
- **Priority:** Must · Confidence: 🟢

### T-SIT-03 — Deployment activation
- **Behavior source:** `Http/Sites/Deployment/Update.php`
- **Description:** ready gate; repoint rules.
- **Done when:** Activation verified.
- **Priority:** Must · Confidence: 🟢

### T-SIT-04 — Chunked deployment upload + CRUD + download
- **Behavior source:** `Http/Deployments/{Create,Get,XList,Delete}.php`, `Download/Get.php`
- **Description:** Locked chunk assembly; preview rule creation; delete pointer resets; token-auth download.
- **Done when:** Upload serialization + download auth verified.
- **Priority:** Must · Confidence: 🟢

### T-SIT-05 — Template/VCS/Duplicate deployments
- **Behavior source:** `Http/Deployments/{Template,Vcs,Duplicate}/Create.php`
- **Description:** Template from public GitHub; VCS branch/commit/tag; duplicate rebuild.
- **Done when:** All three modes verified.
- **Priority:** Must · Confidence: 🟢

### T-SIT-06 — Deployment status (cancel)
- **Behavior source:** `Http/Deployments/Status/Update.php`
- **Description:** Terminal-state guard; per-deployment lock write.
- **Done when:** Cancel flow verified.
- **Priority:** Must · Confidence: 🟢

### T-SIT-07 — Variables CRUD
- **Behavior source:** `Http/Variables/**`
- **Description:** CRUD with secret rules; live=false on change.
- **Done when:** Secret + live flag verified.
- **Priority:** Must · Confidence: 🟢

### T-SIT-08 — Catalog endpoints
- **Behavior source:** `Http/{Frameworks,Templates,Specifications}/XList.php`, `Http/Templates/Get.php`
- **Description:** Allowlist-filtered frameworks, templates with filters, plan-gated specs.
- **Done when:** Catalog filtering verified.
- **Priority:** Should · Confidence: 🟢

### T-SIT-09 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Sites/Http/**`
- **Description:** OpenAPI 3.1 for sites endpoints.
- **Done when:** `openapi/sites.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-SIT-10 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Sites/**`
- **Description:** Cover CRUD, all deployment modes, activation, cancel, variables, download, catalog.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
