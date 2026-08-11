# Design: Functions

> Unit: `Functions`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Functions module implements the **serverless lifecycle** across HTTP (30 actions), three background workers (Builds, Jobs, Screenshots), and integrations (open-runtimes Executor, jobs-service, VCS providers, storage devices).

```
Client
  │ scopes: functions.*, executions.*
  ▼
app/http.php → Functions\Http\** (30 actions)
  │
  ├── dbForProject (functions, deployments, executions, variables)
  ├── dbForPlatform (installations, schedules, repositories, rules)
  ├── Executor (open-runtimes) — sync executions, cold-start
  ├── Deployments service — deployment docs + build trigger
  ├── deviceForFunctions / deviceForBuilds / deviceForFiles
  ├── publisherForBuilds → worker-builds
  ├── publisherForFunctions → worker-functions (executions)
  ├── publisherForJobs → worker-jobs (callbacks)
  ├── publisherForScreenshots → worker-screenshots
  └── publisherForDeletes / publisherForUsage
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/functions` | `createFunction` | functions.write |
| GET | `/v1/functions/:functionId` | `getFunction` | functions.read |
| PUT | `/v1/functions/:functionId` | `updateFunction` | functions.write |
| GET | `/v1/functions` | `listFunctions` | functions.read |
| DELETE | `/v1/functions/:functionId` | `deleteFunction` | functions.write |
| PATCH | `/v1/functions/:functionId/deployment` | `updateFunctionDeployment` | functions.write |
| GET | `/v1/functions/runtimes` | `listRuntimes` | public |
| GET | `/v1/functions/specifications` | `listSpecifications` | functions.read |
| POST | `/v1/functions/:functionId/deployments` | `createDeployment` | functions.write |
| GET | `/v1/functions/:functionId/deployments/:deploymentId` | `getDeployment` | functions.read |
| GET | `/v1/functions/:functionId/deployments` | `listDeployments` | functions.read |
| DELETE | `/v1/functions/:functionId/deployments/:deploymentId` | `deleteDeployment` | functions.write |
| PATCH | `/v1/functions/:functionId/deployments/:deploymentId/status` | `updateDeploymentStatus` | functions.write |
| POST | `/v1/functions/:functionId/deployments/:deploymentId/duplicate` | `createDuplicateDeployment` | functions.write |
| POST | `/v1/functions/:functionId/deployments/template` | `createTemplateDeployment` | functions.write |
| POST | `/v1/functions/:functionId/deployments/vcs` | `createVcsDeployment` | functions.write |
| GET | `/v1/functions/:functionId/deployments/:deploymentId/download` | `downloadDeployment` | functions.read |
| POST | `/v1/functions/:functionId/executions` | `createExecution` | executions.write, execution.write |
| DELETE | `/v1/functions/:functionId/executions/:executionId` | `deleteExecution` | executions.write |
| POST | `/v1/functions/:functionId/variables` | `createVariable` | functions.write |
| GET | `/v1/functions/:functionId/variables/:variableId` | `getVariable` | functions.read |
| GET | `/v1/functions/:functionId/variables` | `listVariables` | functions.read |
| PUT | `/v1/functions/:functionId/variables/:variableId` | `updateVariable` | functions.write |
| DELETE | `/v1/functions/:functionId/variables/:variableId` | `deleteVariable` | functions.write |
| GET | `/v1/functions/templates/:templateId` | `getTemplate` | functions.read |
| GET | `/v1/functions/templates` | `listTemplates` | functions.read |
| POST | `/v1/functions/:functionId/jobs/event` | `createJobEvent` | functions.write |

## 3. Data flows — key actions

### 3.1 Chunked deployment upload

```
1. Validate gzip content + content-range headers
2. Acquire per-deployment distributed lock (locks)
3. Append chunk to deviceForFunctions archive
4. On complete: create deployment doc (status=processing)
5. Trigger build (Deployments service / BuildMessage) + emit events
6. Release lock
```

### 3.2 Execution (sync vs async)

```
Sync (async=false):
  → Executor::createExecution() direct, wait, return (no persistence)

Async:
  → enqueue FunctionMessage → worker-functions
  → execution doc persisted (status=processing)

Scheduled:
  → schedules doc (SCHEDULE_RESOURCE_TYPE_EXECUTION) in dbForPlatform
  → per-execution JWT (HS256, projectId + scopes) → x-appwrite-key
```

### 3.3 Build pipeline (workers)

```
Builds worker:
  clone user repo + template repo → merge → commit/push → Deployments::createFromUrl()

Jobs worker (callback processing):
  CloudEvent callbacks (log, artifact, exit, complete) → per-deployment lock + Redis dedup
  ready() join: exit(0) + complete (+ manifest for sites)
  → state transitions (ready/failed), activation, usage, webhooks

Screenshots worker:
  after site deploy → light/dark screenshots via browser Client → platform files → update docs
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | functions, deployments, executions, variables |
| `dbForPlatform` | installations, schedules, repositories, rules |
| `queueForEvents` / `queueForWebhooks` / `queueForRealtime` | Event fan-out |
| `publisherForBuilds` / `publisherForFunctions` / `publisherForDeletes` / `publisherForScreenshots` / `publisherForUsage` | Worker queues |
| `executor` | open-runtimes executor (sync exec, runtime lifecycle) |
| `deployments` | Deployment doc factory + build trigger |
| `deviceForFunctions` / `deviceForLocal` / `deviceForBuilds` / `deviceForFiles` | Storage devices |
| `vcsFactory` / `repositoryWebhooks` | VCS adapters + webhook lifecycle |
| `authorization` / `locks` / `bus` / `cache` | Auth, distributed locks, domain events, Redis cache |
| `timelimit` | Rate limiting |
| `plan` | Billing capabilities |

## 5. Events emitted

- `functions.[functionId].{create,update,delete}`
- `functions.[functionId].deployments.[deploymentId].{create,update}`
- `functions.[functionId].executions.[executionId].create`
- `variables.[variableId].{create,update,delete}`
- `rules.[ruleId].create` (legacy function creation with domain rule)

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| Rate limit exceeded | 429 | `_APP_FUNCTIONS_CREATION_ABUSE_LIMIT` |
| Runtime not allowed | 400 | Not in `_APP_FUNCTIONS_RUNTIMES` allow-list |
| Build cancelled | 400 | PATCH status on terminal build |
| Execution errors | 500/400 | Executor failures, truncated logs |

