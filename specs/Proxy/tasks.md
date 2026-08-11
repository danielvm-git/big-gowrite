# Tasks: Proxy

> Unit: `Proxy`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-PRX-01 — Module bootstrap + base action
- **Behavior source:** `Modules/Proxy/Module.php`, `Services/Http.php`, `Action.php`
- **Description:** Register HTTP service (8 actions); implement shared createRule/verifyRule/validateDomainRestrictions/isAppwriteOwned helpers.
- **Done when:** Routes + base helpers functional.
- **Priority:** Must · Confidence: 🟢

### T-PRX-02 — Rule creation (api/function/site)
- **Behavior source:** `Http/Rules/{API,Function,Site}/Create.php`
- **Description:** Domain validation, CAA gate, CNAME/A/AAAA checks, duplicate/orphan handling, cert enqueue, status remap.
- **Done when:** All three types verified.
- **Priority:** Must · Confidence: 🟢

### T-PRX-03 — Redirect rule creation
- **Behavior source:** `Http/Rules/Redirect/Create.php`
- **Description:** Parent resource (site/function) required; status code 301/302/307/308.
- **Done when:** Redirect rules verified.
- **Priority:** Must · Confidence: 🟢

### T-PRX-04 — Get + list rules
- **Behavior source:** `Http/Rules/{Get,XList}.php`
- **Description:** Project-scoped fetch; certificate logs + renewAt hydration.
- **Done when:** Certificate data hydrated.
- **Priority:** Must · Confidence: 🟢

### T-PRX-05 — Delete rule
- **Behavior source:** `Http/Rules/Delete.php`
- **Description:** Ownership; delete + DELETE_TYPE_DOCUMENT enqueue; event.
- **Done when:** Cleanup enqueued.
- **Priority:** Must · Confidence: 🟢

### T-PRX-06 — Retry verification
- **Behavior source:** `Http/Rules/Status/Update.php`
- **Description:** No-op on terminal states; DNS re-validation; TLS trigger on success.
- **Done when:** Retry flow verified.
- **Priority:** Must · Confidence: 🟢

### T-PRX-07 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Proxy/Http/**`
- **Description:** OpenAPI 3.1 for proxy endpoints.
- **Done when:** `openapi/proxy.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-PRX-08 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Proxy/**`
- **Description:** Cover rule types, DNS validation, deny-list, retry, delete cleanup.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
