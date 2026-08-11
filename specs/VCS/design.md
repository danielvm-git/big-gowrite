# Design: VCS

> Unit: `VCS`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The VCS module is an **HTTP integration service** bridging Appwrite and external VCS providers. It stores installations/repositories in the **platform database** and, on webhook events, fans out to per-project databases to create Function/Site deployments. No dedicated workers — deployment creation and cleanup are triggered inline via `deploymentsFactory` and `publisherForDeletes`.

```
Provider (GitHub / GitLab / Gitea)
  │  OAuth2 / GitHub App        webhooks (push, PR, installation)
  ▼                              ▼
Authorize/Callback ────────── Events/Create (secret-validated)
  │                              │
  ▼                              ▼
Installations (dbForPlatform)  Deployment trait (fan-out)
  │                              │
  ├── vcsFactory (adapter)       ├── getProjectDB (functions/sites)
  ├── installationTokens         └── deploymentsFactory → Builds
  └── getRepository* endpoints
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| GET | `/v1/vcs/github/authorize` | `getVCSGitHubAuthorize` | vcs.read |
| GET | `/v1/vcs/github/callback` | `getVCSGitHubCallback` | public |
| GET | `/v1/vcs/gitea/authorize` | `getVCSGiteaAuthorize` | vcs.read |
| GET | `/v1/vcs/gitea/callback` | `getVCSGiteaCallback` | public |
| GET | `/v1/vcs/gitlab/authorize` | `getVCSGitlabAuthorize` | vcs.read |
| GET | `/v1/vcs/gitlab/callback` | `getVCSGitlabCallback` | public |
| POST | `/v1/vcs/github/events` | `createVCSGitHubEvent` | public |
| POST | `/v1/vcs/gitea/events` | `createVCSGiteaEvent` | public |
| POST | `/v1/vcs/gitlab/events` | `createVCSGitlabEvent` | public |
| PATCH | `/v1/vcs/github/installations/:installationId/repositories/:repositoryId` | `updateExternalDeployment` | vcs.write |
| GET | `/v1/vcs/installations/:installationId` | `getInstallation` | vcs.read |
| GET | `/v1/vcs/installations` | `listInstallations` | vcs.read |
| DELETE | `/v1/vcs/installations/:installationId` | `deleteInstallation` | vcs.write |
| GET | `/v1/vcs/installations/:installationId/namespaces` | `listNamespaces` | vcs.read |
| POST | `/v1/vcs/github/installations/:installationId/providerRepositories` | `createRepository` | vcs.write |
| GET | `/v1/vcs/github/installations/:installationId/providerRepositories/:providerRepositoryId` | `getRepository` | vcs.read |
| GET | `/v1/vcs/github/installations/:installationId/providerRepositories` | `listRepositories` | vcs.read |
| GET | `/v1/vcs/github/installations/:installationId/providerRepositories/:providerRepositoryId/branches` | `listRepositoryBranches` | vcs.read |
| GET | `/v1/vcs/github/installations/:installationId/providerRepositories/:providerRepositoryId/contents` | `getRepositoryContents` | vcs.read |
| POST | `/v1/vcs/github/installations/:installationId/detections` | `createRepositoryDetection` | vcs.write |

## 3. Data flows — key actions

### 3.1 Authorize + callback

```
1. Authorize: build provider URL; GitHub App installation page / OAuth2 auth URL
2. Callback: verify HMAC state (GitHub public endpoint), exchange code
3. Upsert installation doc in dbForPlatform (projectInternalId scoped)
4. Return redirect to console
```

### 3.2 Webhook ingestion (push / PR / installation)

```
1. Validate provider signature → 403 on mismatch
2. Ignore bot-authored commits (APP_VCS_*_EMAIL)
3. Deployment trait fan-out:
   a. Resolve repositories by providerRepositoryId
   b. Apply skip patterns / branch globs / path globs
   c. Fork PRs without auth → comment + pending status (no deployment)
   d. Authorized → deploymentsFactory → BuildMessage
4. Site deployments → up to 3 preview rules (deploymentId / branch / commit prefix)
5. PR comments → optimistic lock (vcsCommentLocks, 9 retries)
```

### 3.3 Repository browsing

```
1. vcsFactory builds adapter for the installation
2. installationTokens refreshes OAuth2 token (Gitea/GitLab)
3. Swoole coroutines for listRepositories / detection (batch, go, WaitGroup)
```

### 3.4 Installation delete / App uninstall

```
1. Delete installation doc
2. Enqueue cleanup via publisherForDeletes
3. installation.deleted → region-scoped cleanup of repositories +
   installationId/repositoryId refs on functions and sites
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForPlatform` | installations, repositories, vcsComments, vcsCommentLocks, rules, projects |
| `getProjectDB` | Per-project DB (functions, sites, deployments) |
| `vcsFactory` | Provider adapter builder (GitHub, Gitea, GitLab) |
| `installationTokens` | OAuth2 token refresh (Gitea/GitLab) |
| `vcsWebhookSecret` | Provider webhook secret callable |
| `deploymentsFactory` | Deployment document creation |
| `publisherForDeletes` | Installation cleanup |
| `authorization` | Permission bypass in webhook processing |
| `redirectValidator`, `vcsConfigured`, `vcsProviders` | OAuth2/redirect + provider config |

## 5. Events emitted

- No Appwrite-internal `queueForEvents` events — deployments trigger downstream pipeline events via the Deployments service.
- External events posted back to the provider: commit status updates (pending/success/failure) and PR comments via the VCS adapter.

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| Invalid webhook signature | 403 | Provider secret mismatch |
| Installation not found | 404 | Unknown `:installationId` |
| Provider API failure | 4xx/5xx passthrough | Adapter call errors |

## 7. Security

- HMAC-signed OAuth2 `state` on public callback endpoints.
- Provider webhook secrets validated on every event.
- Bot-email filtering prevents deployment loops.
- Region-scoped uninstall cleanup.
- Optimistic locking on PR comments.

