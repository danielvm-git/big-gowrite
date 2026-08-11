# Requirements: Functions

> Identifier: `FNC-SERVERLESS`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Functions module is the core **serverless compute** surface. It manages the full function lifecycle: definition (name, runtime, schedule, timeouts, scopes), code deployments (file upload, VCS/Git-connected, template-based, duplicate), environment variables, triggered executions (sync and async), and jobs-service callbacks. Read-only views expose runtimes, machine specifications, and curated templates. Background workers handle the async build pipeline (Builds), real-time build log streaming + deployment state (Jobs), and site screenshots.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Functions` | Module purpose | 🟢 |
| `specs/Functions/legacy-mapping.md` | Endpoints, workers, business rules, DI | 🟢 |
| `specs/domain.md#Functions` | Function/deployment/execution model | 🟢 |
| `specs/workers.md#{Builds,Jobs,Screenshots}` | Worker pipelines | 🟢 |
| `specs/architecture.md#JobsService` | Build orchestration callbacks | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App developer | Write and deploy code | Create function, upload gzip, activate deployment |
| App developer | Connect to Git | VCS deployment on branch push |
| App developer | Run on schedule/trigger | Configure CRON schedule, async executions |
| App developer | Configure env | Create/update secret variables |

## 4. Business rules

1. **BR-01:** Function creation rate-limited: `_APP_FUNCTIONS_CREATION_ABUSE_LIMIT` (default 50) per project/day. 🟢
2. **BR-02:** Runtime must be in project allow-list `_APP_FUNCTIONS_RUNTIMES` if set. 🟢
3. **BR-03:** VCS connection requires `installationId` + `providerBranch` + `providerRepositoryId`. 🟢
4. **BR-04:** VCS function creation creates a `repositories` doc in dbForPlatform; repo webhooks registered if installation lacks installation-scoped webhooks. 🟢
5. **BR-05:** A `schedules` doc is always created at function creation; active only when CRON schedule + active deployment. 🟢
6. **BR-06:** VCS disconnect deletes all associated `repositories` docs. 🟢
7. **BR-07:** Changing spec on a function with active deployment triggers cold-start via `Executor::deleteRuntime()`. 🟢
8. **BR-08:** Variable create/update sets `function.live = false` (redeploy signal); schedule refreshed. 🟢
9. **BR-09:** Secret variables masked for non-privileged calls; returned normally via admin/key. 🟢
10. **BR-10:** Code deployments chunked multipart via `content-range`; per-deployment distributed lock serializes chunks. 🟢
11. **BR-11:** Sync execution (`async=false`) calls `Executor::createExecution()` directly; execution doc not persisted. 🟢
12. **BR-12:** Async execution enqueues `FunctionMessage`; scheduled executions create `schedules` doc with `SCHEDULE_RESOURCE_TYPE_EXECUTION`. 🟢
13. **BR-13:** Per-execution short-lived JWT (HS256, `_APP_OPENSSL_KEY_V1`) passed as `x-appwrite-key` to executor with project ID + function scopes. 🟢
14. **BR-14:** Execution log/error truncated to `APP_FUNCTION_LOG_LENGTH_LIMIT` / `APP_FUNCTION_ERROR_LENGTH_LIMIT`. 🟢
15. **BR-15:** Jobs worker serializes callbacks per deployment under distributed lock (TTL 30s, timeout 10s), dedupes by CloudEvent `id` (Redis TTL 1h). 🟢
16. **BR-16:** Build success requires `exit(0)` + `complete` callbacks (+ `manifest` for sites); join pattern in `ready()`. 🟢
17. **BR-17:** On successful build with `activate=true`, function `deploymentId` updated; routing rules repointed. 🟢
18. **BR-18:** v2 runtime compat injects legacy `APPWRITE_FUNCTION_TRIGGER`/`APPWRITE_FUNCTION_DATA` env vars. 🟡
19. **BR-19:** Executions retention pruning async via `publisherForDeletes` when `executionsRetentionCount > 0` + `ENABLE_EXECUTIONS_LIMIT_ON_ROUTE`. 🟡

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Function CRUD + VCS linkage | Must | create/get/update/list/delete; VCS connect/disconnect | 🟢 |
| RF-02 | Activate deployment | Must | PATCH .../deployment repoints active deployment | 🟢 |
| RF-03 | List runtimes/specifications | Must | Public runtime list; spec tiers | 🟢 |
| RF-04 | Deployments: upload (chunked), get/list/delete, download | Must | Multipart assembly via content-range + lock | 🟢 |
| RF-05 | Deployment status (cancel build) | Must | PATCH .../status cancels in-progress build | 🟢 |
| RF-06 | Duplicate deployment | Must | Clone existing deployment | 🟢 |
| RF-07 | Template deployment | Must | Deploy from curated template | 🟢 |
| RF-08 | VCS deployment | Must | Trigger branch deployment | 🟢 |
| RF-09 | Executions (sync/async) + delete | Must | Sync via Executor; async via FunctionMessage | 🟢 |
| RF-10 | Variables CRUD with secret masking | Must | Mask for non-privileged; live=false on change | 🟢 |
| RF-11 | Templates list/get | Must | Curated template catalog | 🟢 |
| RF-12 | Jobs-service callback endpoint | Must | POST .../jobs/event (internal) | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Reliability | Chunk upload serialization via distributed lock | `locks` per deployment | 🟢 |
| Reliability | Callback dedup + join pattern | Redis dedup (TTL 1h), `ready()` join | 🟢 |
| Security | Per-execution JWT with scoped claims | HS256 `x-appwrite-key` | 🟢 |
| Performance | Sync vs async execution paths | Executor direct vs queue | 🟢 |
| Governança | Creation rate limit + runtime allow-list | `_APP_FUNCTIONS_CREATION_ABUSE_LIMIT`, `_APP_FUNCTIONS_RUNTIMES` | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Chunked deployment upload
  Given a multipart upload split into chunks
  When POST .../deployments is called with content-range headers
  Then chunks are assembled on the storage device under a per-deployment lock

Scenario: Sync execution
  Given a function with an active deployment
  When POST .../executions with async=false is called
  Then the executor runs directly and the result is returned without a persisted execution doc

Scenario: Variable change forces redeploy
  Given an existing function with a variable
  When the variable is updated
  Then function.live is set to false and the schedule is refreshed

Scenario: Secret masked for non-privileged call
  Given a secret variable
  When a non-privileged client lists variables
  Then the value is masked
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..05, 09, 10 (core) | Must | Serverless lifecycle |
| RF-06..08 (deploy modes) | Must | VCS/template/duplicate |
| RF-11 (templates) | Should | Catalog UX |
| RF-12 (jobs callback) | Must | Build pipeline integration |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] v2 runtime compat env vars exact injection points.
- 🟡 [INFERRED] Executions retention pruning trigger conditions.
- 🟡 [INFERRED] Jobs-service callback CloudEvent schema.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |

