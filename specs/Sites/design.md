# Design: Sites

> Unit: `Sites`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Sites module is an **HTTP service for static/SSR site hosting**. It manages sites, deployments (4 creation modes), and variables, integrating with VCS, storage devices, the jobs-service build pipeline, and the Proxy rules system for preview domains.

```
Client
  │ scopes: sites.*
  ▼
app/http.php → Sites\Http\** (26 actions)
  │
  ├── dbForProject (sites, deployments, variables)
  ├── dbForPlatform (installations, repositories, rules)
  ├── deviceForSites / deviceForBuilds / deviceForLocal
  ├── vcsFactory + repositoryWebhooks
  ├── executor (runtime lifecycle on spec change)
  ├── Deployments service helpers
  ├── locks (upload concurrency)
  └── publisherForDeletes / publisherForBuilds
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/sites` | `createSite` | sites.write |
| GET | `/v1/sites` | `listSites` | sites.read |
| GET | `/v1/sites/:siteId` | `getSite` | sites.read |
| PUT | `/v1/sites/:siteId` | `updateSite` | sites.write |
| DELETE | `/v1/sites/:siteId` | `deleteSite` | sites.write |
| PATCH | `/v1/sites/:siteId/deployment` | `updateSiteDeployment` | sites.write |
| GET | `/v1/sites/frameworks` | `listFrameworks` | public |
| POST | `/v1/sites/:siteId/deployments` | `createDeployment` | sites.write |
| POST | `/v1/sites/:siteId/deployments/template` | `createTemplateDeployment` | sites.write |
| POST | `/v1/sites/:siteId/deployments/vcs` | `createVcsDeployment` | sites.write |
| POST | `/v1/sites/:siteId/deployments/duplicate` | `createDuplicateDeployment` | sites.write |
| GET | `/v1/sites/:siteId/deployments` | `listDeployments` | sites.read |
| GET | `/v1/sites/:siteId/deployments/:deploymentId` | `getDeployment` | sites.read |
| DELETE | `/v1/sites/:siteId/deployments/:deploymentId` | `deleteDeployment` | sites.write |
| PATCH | `/v1/sites/:siteId/deployments/:deploymentId/status` | `updateDeploymentStatus` | sites.write |
| GET | `/v1/sites/:siteId/deployments/:deploymentId/download` | `getDeploymentDownload` | public / sites.read |
| POST | `/v1/sites/:siteId/variables` | `createVariable` | sites.write |
| GET | `/v1/sites/:siteId/variables` | `listVariables` | sites.read |
| GET | `/v1/sites/:siteId/variables/:variableId` | `getVariable` | sites.read |
| PUT | `/v1/sites/:siteId/variables/:variableId` | `updateVariable` | sites.write |
| DELETE | `/v1/sites/:siteId/variables/:variableId` | `deleteVariable` | sites.write |
| GET | `/v1/sites/templates` | `listTemplates` | public |
| GET | `/v1/sites/templates/:templateId` | `getTemplate` | public |
| GET | `/v1/sites/specifications` | `listSpecifications` | sites.read |

## 3. Data flows — key actions

### 3.1 Create site (VCS)

```
1. Validate framework + adapter (static/ssr allowed)
2. VCS: installationId + providerBranch required
3. Create repositories doc + ensure provider webhook (fail → delete orphan repo)
4. Create site doc; emit sites.[siteId].create
```

### 3.2 Chunked deployment upload

```
1. Acquire per-project/site/deployment lock (600s) → 429 Retry-After: 5 on contention
2. Assemble chunks on deviceForSites
3. Create deployment doc; create preview rule (unique subdomain, verified)
4. Trigger build; release lock
```

### 3.3 Activation

```
1. Deployment status == ready? else BUILD_NOT_READY
2. Update site.deploymentId; repoint rules
```

### 3.4 Deployment delete

```
1. Delete deployment + storage artifact
2. Reset latestDeploymentId to next recent
3. Clear active deploymentId if deleted was active
```

### 3.5 Download

```
1. Privileged caller (session/API key) OR presigned deployment token (Token::verify)
2. Stream source or output artifact
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | sites, deployments, variables |
| `dbForPlatform` | installations, repositories, rules |
| `queueForEvents` | Events |
| `publisherForDeletes` | Cascading deletes |
| `publisherForBuilds` | VCS build jobs |
| `deviceForSites` / `deviceForBuilds` / `deviceForLocal` | Storage |
| `vcsFactory` / `repositoryWebhooks` | VCS |
| `executor` | Runtime delete on spec change |
| `authorization` | Auth context |
| `deployments` | Deployment helpers |
| `locks` | Upload lock factory |
| `plan` / `platform` | Plan limits + sitesDomain |

## 5. Events emitted

- `sites.[siteId].{create,update,delete}`
- `sites.[siteId].deployments.[deploymentId].{create,update,delete}`
- `variables.[variableId].{create,update,delete}`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `GENERAL_ARGUMENT_INVALID` | 400 | Missing VCS pair / bad adapter |
| `BUILD_NOT_READY` | 400 | Activate non-ready |
| `BUILD_ALREADY_COMPLETED` | 400 | Cancel ready/failed |
| `VARIABLE_CANNOT_UNSET_SECRET` | 400 | Secret downgrade |
| `LockContention` | 429 | Upload lock busy (Retry-After: 5) |
| `USER_UNAUTHORIZED` | 401 | Download without valid token |
