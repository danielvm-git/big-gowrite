# Contracts: Sites

> Unit: `Sites`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract (representative)

Base: `/v1/sites`. Scopes: `sites.*`; catalog endpoints `public` where noted.

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `/sites` | `siteId`, `name`, `framework`, `adapter`, `installationId`, `providerBranch` (VCS) | 201 site | GENERAL_ARGUMENT_INVALID |
| PUT | `/sites/:siteId` | `name`, build/install/start commands, `outputDirectory`, `providerRepositoryId` | 200 site | 404 |
| DELETE | `/sites/:siteId` | — | 204 | 404 |
| PATCH | `/sites/:siteId/deployment` | `deploymentId` | 200 site | BUILD_NOT_READY |
| GET | `/sites/frameworks` | — | 200 frameworkList | public |
| POST | `/sites/:siteId/deployments` | gzip, `content-range` | 201 deployment | 429 lock (Retry-After: 5) |
| POST | `/sites/:siteId/deployments/template` | `templateId`, `providerRepositoryId`, `providerBranch` | 201 | |
| POST | `/sites/:siteId/deployments/vcs` | `providerBranch`, `providerCommitHash` | 201 | |
| POST | `/sites/:siteId/deployments/duplicate` | — | 201 | |
| GET | `/sites/:siteId/deployments/:deploymentId` | — | 200 deployment | 404 |
| PATCH | `/sites/:siteId/deployments/:deploymentId/status` | `status=canceled` | 200 | BUILD_ALREADY_COMPLETED |
| GET | `/sites/:siteId/deployments/:deploymentId/download` | `token` (presigned) | 200 artifact | USER_UNAUTHORIZED |
| POST | `/sites/:siteId/variables` | `key`, `value`, `secret` | 201 variable | |
| PUT | `/sites/:siteId/variables/:variableId` | `key`, `value`, `secret` | 200 | VARIABLE_CANNOT_UNSET_SECRET |
| GET | `/sites/templates` | `framework`, `useCase` | 200 templateList | public |
| GET | `/sites/specifications` | `runtime`/`build` | 200 specList | 401 |

## Preview domain contract

Each new deployment creates a `rules` doc in dbForPlatform:
- Unique subdomain under `platform['sitesDomain']`
- Status `verified` immediately (no DNS wait)

## Queue contracts

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForBuilds` | Http → worker-builds | BuildMessage (VCS/template deployments) |
| `publisherForDeletes` | Site/deployment delete → worker-deletes | Cascading delete jobs |

## Event contract

- `sites.[siteId].{create,update,delete}`
- `sites.[siteId].deployments.[deploymentId].{create,update,delete}`
- `variables.[variableId].{create,update,delete}`

## Configuration contract

| Env var | Effect |
|---------|--------|
| `_APP_SITES_RUNTIMES` | Framework runtime allow-list |
| `_APP_COMPUTE_CPUS` / `_APP_COMPUTE_MEMORY` | Spec caps |
| `_APP_RULES_FORMAT=md5` | Legacy MD5 rule IDs |
| `platform['sitesDomain']` | Preview subdomain suffix |
