# Tasks: Functions

> Unit: `Functions`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-FNC-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Functions/Module.php`, `Services/Http.php`, `Services/Workers.php`
- **Description:** Register HTTP service (27 actions) + Workers (builds, jobs, screenshots).
- **Done when:** All routes + workers registered.
- **Priority:** Must · Confidence: 🟢

### T-FNC-02 — Function CRUD + VCS linkage
- **Behavior source:** `Http/Functions/{Create,Get,Update,XList,Delete}.php`
- **Description:** CRUD with rate limit, runtime allow-list, VCS connect/disconnect (repositories docs, webhooks, schedules doc).
- **Done when:** Creation limits + VCS flows verified.
- **Priority:** Must · Confidence: 🟢

### T-FNC-03 — Deployment activation + spec cold-start
- **Behavior source:** `Http/Functions/Deployment/Update.php`
- **Description:** Activate deployment; spec change triggers Executor::deleteRuntime cold-start.
- **Done when:** Activation repoints deployment; cold-start on spec change.
- **Priority:** Must · Confidence: 🟢

### T-FNC-04 — Runtimes + specifications list
- **Behavior source:** `Http/Runtimes/XList.php`, `Http/Specifications/XList.php`
- **Description:** Read-only runtime and spec tier lists.
- **Done when:** Lists match configured runtimes/specs.
- **Priority:** Must · Confidence: 🟢

### T-FNC-05 — Chunked deployment upload + CRUD + download
- **Behavior source:** `Http/Deployments/{Create,Get,XList,Delete}.php`, `Download/Get.php`
- **Description:** Multipart chunk assembly with content-range + distributed lock; get/list/delete/download.
- **Done when:** Concurrent chunk uploads serialize correctly.
- **Priority:** Must · Confidence: 🟢

### T-FNC-06 — Deployment status (cancel build)
- **Behavior source:** `Http/Deployments/Status/Update.php`
- **Description:** Cancel in-progress build.
- **Done when:** Build cancelled; terminal state handled.
- **Priority:** Must · Confidence: 🟢

### T-FNC-07 — Duplicate deployment
- **Behavior source:** `Http/Deployments/Duplicate/Create.php`
- **Description:** Clone an existing deployment.
- **Done when:** Duplicate produces independent deployment.
- **Priority:** Must · Confidence: 🟢

### T-FNC-08 — Template deployment
- **Behavior source:** `Http/Deployments/Template/Create.php`
- **Description:** Deploy curated template (Builds worker pipeline).
- **Done when:** Template merged and built.
- **Priority:** Must · Confidence: 🟢

### T-FNC-09 — VCS deployment
- **Behavior source:** `Http/Deployments/Vcs/Create.php`
- **Description:** Trigger branch deployment via VCS adapter.
- **Done when:** Branch deployment triggered.
- **Priority:** Must · Confidence: 🟢

### T-FNC-10 — Executions (sync/async/scheduled) + delete
- **Behavior source:** `Http/Executions/{Create,Delete}.php`
- **Description:** Sync via Executor (no persistence); async via FunctionMessage; scheduled via schedules doc; per-execution JWT.
- **Done when:** Both paths + JWT verified.
- **Priority:** Must · Confidence: 🟢

### T-FNC-11 — Variables CRUD + secret masking
- **Behavior source:** `Http/Variables/**`
- **Description:** CRUD; mask secrets for non-privileged; live=false + schedule refresh.
- **Done when:** Masking enforced; redeploy signal set.
- **Priority:** Must · Confidence: 🟢

### T-FNC-12 — Templates list/get
- **Behavior source:** `Http/Templates/{Get,XList}.php`
- **Description:** Curated template catalog read endpoints.
- **Done when:** Catalog served.
- **Priority:** Should · Confidence: 🟢

### T-FNC-13 — Jobs-service callback endpoint
- **Behavior source:** `Http/Jobs/Event/Create.php`
- **Description:** Receive CloudEvent callbacks; enqueue for Jobs worker.
- **Done when:** Callbacks accepted and processed.
- **Priority:** Must · Confidence: 🟢

### T-FNC-14 — Builds worker
- **Behavior source:** `Workers/Builds.php`
- **Description:** Template-into-VCS build pipeline.
- **Done when:** Clone/merge/commit/push + handoff verified.
- **Priority:** Must · Confidence: 🟢

### T-FNC-15 — Jobs worker
- **Behavior source:** `Workers/Jobs.php`
- **Description:** Process log/artifact/exit/complete callbacks; lock + dedup; join pattern; activation; usage; webhooks.
- **Done when:** State transitions and dedup verified.
- **Priority:** Must · Confidence: 🟢

### T-FNC-16 — Screenshots worker
- **Behavior source:** `Workers/Screenshots.php`
- **Description:** Capture light/dark site screenshots; store; update docs.
- **Done when:** Screenshots produced and stored.
- **Priority:** Should · Confidence: 🟢

### T-FNC-17 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Functions/Http/**`
- **Description:** OpenAPI 3.1 for functions endpoints.
- **Done when:** `openapi/functions.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-FNC-18 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Functions/**`
- **Description:** Cover CRUD, deployments (all modes), executions, variables, templates, callbacks.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
