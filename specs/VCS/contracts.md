# Contracts: VCS

> Unit: `VCS`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/vcs`. Scopes: `vcs.*`; authorize/callback/events endpoints `public`.

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| GET | `/github/authorize` | — | 302 redirect | — |
| GET | `/github/callback` | `code`, `state` (HMAC) | 302 redirect | invalid state |
| GET | `/gitea/authorize` | — | 302 redirect | — |
| GET | `/gitea/callback` | `code`, `state` | 302 redirect | — |
| GET | `/gitlab/authorize` | — | 302 redirect | — |
| GET | `/gitlab/callback` | `code`, `state` | 302 redirect | — |
| POST | `/github/events` | provider webhook JSON | 200 | 403 invalid signature |
| POST | `/gitea/events` | provider webhook JSON | 200 | 403 invalid signature |
| POST | `/gitlab/events` | provider webhook JSON | 200 | 403 invalid signature |
| PATCH | `/github/installations/:installationId/repositories/:repositoryId` | `enabled` | 200 | 404 |
| GET | `/installations/:installationId` | — | 200 installation | 404 |
| GET | `/installations` | `queries` | 200 installationList | — |
| DELETE | `/installations/:installationId` | — | 204 | 404 |
| GET | `/installations/:installationId/namespaces` | — | 200 namespaceList | — |
| POST | `/github/installations/:installationId/providerRepositories` | `name`, `private` | 201 repository | 404 |
| GET | `/github/installations/:installationId/providerRepositories/:providerRepositoryId` | — | 200 repository | 404 |
| GET | `/github/installations/:installationId/providerRepositories` | `search`, `queries` | 200 repositoryList | — |
| GET | `/github/installations/:installationId/providerRepositories/:providerRepositoryId/branches` | — | 200 branchList | — |
| GET | `/github/installations/:installationId/providerRepositories/:providerRepositoryId/contents` | `path` | 200 contentList | 404 |
| POST | `/github/installations/:installationId/detections` | `providerRepositoryId` | 201 detection | 404 |

## Installation document contract

| Field | Detail |
|-------|--------|
| Provider | `github` (App) / `gitea`, `gitlab` (OAuth2 PAT) |
| Scoping | `dbForPlatform`, `projectInternalId` |
| Token handling | Gitea/GitLab refreshed by `installationTokens`; GitHub App via VCS factory |

## Webhook contract

| Aspect | Detail |
|--------|--------|
| Validation | Provider secret (`_APP_VCS_*_WEBHOOK_SECRET`); 403 on mismatch |
| Bot-email ignore | `APP_VCS_*_EMAIL` pushes skipped |
| Skip patterns | `VCS_DEPLOYMENT_SKIP_PATTERNS` in commit message |
| Branch gate | `providerBranches` glob |
| Path gate | `providerPaths` glob |
| Fork PR | Unauthorized → comment + pending status; authorized → deployment |

## Preview rules contract (site deployments)

Up to 3 `rules` docs in `dbForPlatform` per deployment: unique deployment ID, branch, commit-hash prefix.

## Queue contracts

| Channel | Direction | Payload |
|---------|-----------|---------|
| `deploymentsFactory` | Webhook → Deployments service | Deployment docs → BuildMessage |
| `publisherForDeletes` | Install/uninstall → worker-deletes | Installation + repo cleanup, function/site ref removal |

## Event contract (external)

- Commit status updates: pending / success / failure (posted to provider).
- PR comments: authorize URL for external contributors (optimistic-locked via `vcsCommentLocks`).
