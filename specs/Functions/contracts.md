# Contracts: Functions

> Unit: `Functions`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract (representative)

Base: `/v1/functions`. Scopes: `functions.*`, `executions.*`, `execution.*`.

| Method | Path | Key params | Success | Notes |
|--------|------|-----------|---------|-------|
| POST | `/functions` | `functionId`, `name`, `runtime`, `schedule`, `timeout`, `scopes`, `entrypoint`, `providerRepositoryId` (VCS) | 201 function | |
| PUT | `/functions/:functionId` | `name`, `schedule`, `spec` | 200 function | VCS disconnect |
| PATCH | `/functions/:functionId/deployment` | `deploymentId` | 200 function | activate |
| GET | `/functions/runtimes` | — | 200 runtimeList | public |
| POST | `/functions/:functionId/deployments` | gzip file, `content-range` | 201 deployment | chunked |
| PATCH | `/functions/:functionId/deployments/:deploymentId/status` | `status=canceled` | 200 | cancel build |
| POST | `/functions/:functionId/deployments/:deploymentId/duplicate` | — | 201 deployment | |
| POST | `/functions/:functionId/deployments/template` | `templateId`, `providerRepositoryId`, `providerBranch` | 201 deployment | |
| POST | `/functions/:functionId/deployments/vcs` | `providerBranch`, `providerCommitHash` | 201 deployment | |
| GET | `/functions/:functionId/deployments/:deploymentId/download` | — | 200 gzip | |
| POST | `/functions/:functionId/executions` | `body`, `path`, `method`, `headers`, `async` | 200/202 execution | sync waits |
| POST | `/functions/:functionId/variables` | `key`, `value` (secret?) | 201 variable | |
| GET | `/functions/templates` | — | 200 templateList | |
| POST | `/functions/:functionId/jobs/event` | CloudEvent callback | 200 | internal |

## Queue contracts

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForBuilds` | Http → worker-builds | BuildMessage: deploymentId, resourceType=functions |
| `publisherForFunctions` | Http → worker-functions | FunctionMessage: executionId, functionId |
| `publisherForJobs` | Http → worker-jobs | CloudEvent (log/artifact/exit/complete) |
| `publisherForScreenshots` | worker-builds → worker-screenshots | Site deployment screenshot request |
| `publisherForDeletes` | Executions → worker-deletes | Retention pruning (executionsRetentionCount) |
| `publisherForUsage` | Jobs worker → usage | Build/execution usage metrics |

## Executor contract (open-runtimes)

| Call | Purpose |
|------|---------|
| `Executor::createExecution()` | Sync execution (direct) |
| `Executor::deleteRuntime()` | Cold-start on spec change |

Auth: per-execution JWT (HS256, `_APP_OPENSSL_KEY_V1`), header `x-appwrite-key` with projectId + function scopes.

## Jobs-service CloudEvent callbacks

| Type | Meaning |
|------|---------|
| `orchestrator.job.log` | Build log chunk (streamed via Realtime) |
| `orchestrator.job.artifact` | Artifact manifest (sites) |
| `orchestrator.job.exit` | Exit code (0 = success) |
| `orchestrator.job.complete` | Build finished |

## Event contract

- `functions.[functionId].{create,update,delete}`
- `functions.[functionId].deployments.[deploymentId].{create,update}`
- `functions.[functionId].executions.[executionId].create`
- `variables.[variableId].{create,update,delete}`
- `rules.[ruleId].create`

## Configuration contract

| Env var | Effect |
|---------|--------|
| `_APP_FUNCTIONS_CREATION_ABUSE_LIMIT` (50) | Per-project/day creation rate limit |
| `_APP_FUNCTIONS_RUNTIMES` | Runtime allow-list |
| `APP_FUNCTION_LOG_LENGTH_LIMIT` / `APP_FUNCTION_ERROR_LENGTH_LIMIT` | Log/error truncation |
| `_APP_OPENSSL_KEY_V1` | JWT signing key |
| `ENABLE_EXECUTIONS_LIMIT_ON_ROUTE` | Retention pruning flag |
