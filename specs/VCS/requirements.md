# Requirements: VCS

> Identifier: `VCS-DEPLOY`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The VCS module integrates Appwrite with external version control providers (GitHub, GitLab, Gitea). It manages the full lifecycle of VCS installations (GitHub App-based or OAuth2 personal-access-token based), exposes repository browsing/detection endpoints, and receives provider webhooks that automatically trigger deployments for Functions and Sites on push and pull-request events.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#VCS` | Module purpose | 🟢 |
| `specs/VCS/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/architecture.md#Integrations` | External provider integration | 🟢 |
| `specs/domain.md#Installation` | Installation model | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App developer | Connect a repo | OAuth2/App install → browse repositories → create site/function |
| App developer | Auto-deploy | Push to a branch triggers deployment |
| OSS maintainer | Fork PRs | External contributor authorized for deployment |

## 4. Business rules

1. **BR-01:** GitHub uses a GitHub App (App-based installation); GitLab and Gitea use OAuth2 personal access tokens. Authorize/Callback base classes handle Gitea and GitLab; GitHub has its own implementation. 🟢
2. **BR-02:** The GitHub callback endpoint is public; the `state` parameter is HMAC-signed with `_APP_OPENSSL_KEY_V1` and verified before proceeding (prevents token attachment to arbitrary projects). 🟢
3. **BR-03:** Installation documents live in the platform database (`dbForPlatform`), scoped to a project via `projectInternalId`. 🟢
4. **BR-04:** Webhook events are validated against provider-specific secrets (`_APP_VCS_GITHUB_WEBHOOK_SECRET`, `_APP_VCS_GITEA_WEBHOOK_SECRET`, `_APP_VCS_GITLAB_WEBHOOK_SECRET`); invalid signature → 403. 🟢
5. **BR-05:** Push events authored by Appwrite's own bot email (`APP_VCS_GITHUB_EMAIL`, `APP_VCS_GITEA_EMAIL`, `APP_VCS_GITLAB_EMAIL`) are ignored to avoid deployment loops. 🟢
6. **BR-06:** The `Deployment` trait fans out one webhook event to many repositories/projects; each repository lookup uses `providerRepositoryId` as the key. 🟢
7. **BR-07:** Deployments are skipped when the commit message matches `VCS_DEPLOYMENT_SKIP_PATTERNS`, the branch does not match configured `providerBranches` globs, or affected files do not match `providerPaths` globs. 🟢
8. **BR-08:** External (fork) PR deployments require explicit authorization; when not authorized, a comment with an authorize URL is posted and a pending commit status set — no deployment is created. 🟢
9. **BR-09:** Site deployments create up to three preview domain rules in `dbForPlatform.rules`: one per unique deployment ID, one per branch, one per commit-hash prefix. 🟢
10. **BR-10:** PR comments are protected by an optimistic lock (`vcsCommentLocks`, up to 9 retries with 1-second sleep) to avoid concurrent overwrites. 🟢
11. **BR-11:** On GitHub App uninstall (`installation.deleted`), all installation documents, linked repositories, and `installationId`/`repositoryId` references on functions/sites are cleaned up — only for projects whose database is in the current region. 🟢
12. **BR-12:** `InstallationTokens` refreshes OAuth2 tokens for Gitea/GitLab before each API call; GitHub App tokens are obtained via the VCS factory. 🟡
13. **BR-13:** `listRepositories` and `createRepositoryDetection` use Swoole coroutines (`batch`, `go`, `WaitGroup`) for concurrent repository analysis. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Authorize + callback (3 providers) | Must | OAuth2/App flow; HMAC state; installation upsert | 🟢 |
| RF-02 | Installation get/list/delete | Must | Platform-DB scoped; delete enqueues cleanup | 🟢 |
| RF-03 | Repository browse (list, get, branches, contents) | Must | Provider API via adapter; runtime/framework detection | 🟢 |
| RF-04 | Create repository + detection | Must | Provider API create; concurrent detection | 🟢 |
| RF-05 | Webhook ingestion (push, PR, installation) | Must | Secret validation; bot-email ignore; fan-out deployments | 🟢 |
| RF-06 | External (fork) PR authorization | Must | Comment + pending status; authorized → deployment | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Webhook signature validation | provider secrets, 403 on mismatch | 🟢 |
| Security | Anti-loop bot-email filtering | APP_VCS_*_EMAIL | 🟢 |
| Reliability | Comment optimistic locking | vcsCommentLocks, 9 retries | 🟢 |
| Performance | Concurrent repo analysis | Swoole coroutines | 🟢 |
| Governança | Region-scoped uninstall cleanup | current-region check | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Invalid webhook signature
  Given a webhook event with a bad provider secret
  When POST /v1/vcs/github/events is called
  Then 403 is returned

Scenario: Bot-authored push
  Given a push event from Appwrite's own bot email
  When the webhook handler processes it
  Then the event is ignored (no deployment)

Scenario: Branch outside providerBranches glob
  Given a push on a non-matching branch
  When the webhook handler processes it
  Then the deployment is skipped

Scenario: Unauthorized fork PR
  Given an external PR without prior authorization
  When the webhook handler processes it
  Then a comment with an authorize URL and a pending status are posted, no deployment

Scenario: GitHub App uninstall
  Given an installation.deleted event
  When the cleanup runs
  Then installations, repositories, and function/site references are removed for the current region
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..06 (core) | Must | VCS lifecycle + auto-deploy |
| Coroutine optimizations | Should | Performance |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Exact OAuth2 token refresh cadence and error recovery.
- 🔴 [GAP] Provider webhook payload schemas (GitHub/GitLab/Gitea) not captured exhaustively.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |

