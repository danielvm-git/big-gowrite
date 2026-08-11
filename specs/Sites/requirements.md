# Requirements: Sites

> Identifier: `SIT-WEB`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Sites module manages the full lifecycle of static and SSR web sites: creation/configuration (framework, build/runtime specs, VCS connection), deployment creation (file upload, public GitHub template, VCS branch/commit, or duplicate), deployment activation, build cancellation, variable management (build and runtime env), and catalogue endpoints (frameworks, templates, compute specs).

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Sites` | Module purpose | 🟢 |
| `specs/Sites/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/domain.md#Site` | Site/deployment model | 🟢 |
| `specs/architecture.md#JobsService` | Build pipeline | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App developer | Host a site | Create site, upload deployment, activate |
| App developer | Connect Git | VCS deployment on branch push |
| App developer | Preview builds | Preview subdomain rules |
| App developer | Configure env | Build/runtime variables |

## 4. Business rules

1. **BR-01:** VCS site creation requires both `installationId` + `providerBranch`; else `GENERAL_ARGUMENT_INVALID`. 🟢
2. **BR-02:** `adapter` (static/ssr) validated against framework's allowed adapters. 🟢
3. **BR-03:** VCS connect creates `repositories` doc + ensures webhook; webhook failure → orphan repo deleted. 🟢
4. **BR-04:** updateSite: null `providerRepositoryId` preserves VCS; empty string disconnects + deletes repositories. 🟢
5. **BR-05:** updateSite sets `live=false` when name, buildCommand, installCommand, startCommand, outputDirectory, providerRootDirectory, or framework change. 🟢
6. **BR-06:** Deployment uploads chunked with per-project/site/deployment distributed lock (600 s timeout); `LockContention` → HTTP 429 `Retry-After: 5`. 🟢
7. **BR-07:** Each new deployment creates preview `rules` doc with unique subdomain under `platform['sitesDomain']`, status `verified` immediately. 🟢
8. **BR-08:** Activate requires deployment `ready`; else `BUILD_NOT_READY`. 🟢
9. **BR-09:** Cancel blocked for `ready`/`failed` → `BUILD_ALREADY_COMPLETED`; written under jobs-worker per-deployment lock. 🟢
10. **BR-10:** Deployment delete resets `latestDeploymentId` to next recent; clears active `deploymentId` if deleted. 🟢
11. **BR-11:** Secret variables cannot be downgraded → `VARIABLE_CANNOT_UNSET_SECRET`. 🟢
12. **BR-12:** Variable create/update sets `site.live=false`. 🟢
13. **BR-13:** Download accepts session/API-key or presigned deployment token; invalid → `USER_UNAUTHORIZED`. 🟢
14. **BR-14:** listFrameworks filters by `_APP_SITES_RUNTIMES` allowlist + enabled build runtimes. 🟢
15. **BR-15:** listSpecifications filters by plan runtime/build specs + `_APP_COMPUTE_CPUS`/`_APP_COMPUTE_MEMORY` caps. 🟢
16. **BR-16:** `_APP_RULES_FORMAT=md5` legacy rule IDs. 🟡

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Site CRUD + VCS linkage | Must | create/get/list/update/delete; VCS connect/disconnect; live=false logic | 🟢 |
| RF-02 | Activate deployment | Must | ready gate; BUILD_NOT_READY otherwise | 🟢 |
| RF-03 | List frameworks/templates/specs | Must | Allowlist filtering; public templates | 🟢 |
| RF-04 | Deployments: upload (chunked), get/list/delete/download | Must | Locking; token download auth | 🟢 |
| RF-05 | Template deployment | Must | Public GitHub template | 🟢 |
| RF-06 | VCS deployment | Must | Branch/commit/tag | 🟢 |
| RF-07 | Duplicate deployment | Must | Rebuild from existing source | 🟢 |
| RF-08 | Deployment status (cancel) | Must | Terminal-state guard | 🟢 |
| RF-09 | Variables CRUD | Must | Secret rules; live=false | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Reliability | Chunk upload serialization | distributed lock 600s | 🟢 |
| Security | Download auth via presigned token | Token::verify | 🟢 |
| Availability | Preview domains instantly verified | rules doc status=verified | 🟢 |
| Governança | Runtime/spec allow-listing | `_APP_SITES_RUNTIMES`, plan caps | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: VCS site missing branch
  Given a VCS site creation without providerBranch
  When POST /v1/sites is called
  Then GENERAL_ARGUMENT_INVALID is returned

Scenario: Activate non-ready deployment
  Given a deployment not in ready state
  When PATCH /v1/sites/:siteId/deployment is called
  Then BUILD_NOT_READY is returned

Scenario: Lock contention on upload
  Given a deployment upload already holding the lock
  When a concurrent chunk is uploaded
  Then HTTP 429 with Retry-After: 5 is returned

Scenario: Secret variable downgrade
  Given a secret site variable
  When updating with secret=false is attempted
  Then VARIABLE_CANNOT_UNSET_SECRET is returned
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..09 (core) | Must | Site lifecycle |
| RF-03 (catalog) | Should | UX conveniences |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] `_APP_RULES_FORMAT=md5` legacy rule IDs (removal planned post-1.7.x).
- 🔴 [GAP] Exact presigned download token payload.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
