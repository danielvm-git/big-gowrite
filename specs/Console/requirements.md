# Requirements: Console

> Identifier: `CSL-GATEWAY`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Console module is the gateway and utility layer for the Appwrite web console. It serves the console's own API endpoints (configuration variables, email templates, OAuth2 provider metadata, API key scopes, resource availability, AI assistant proxy) and handles browser-facing concerns: legacy URL redirects to `/console`, security headers for web views, and init hooks that guard console-only and web-only route groups.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Console` | Gateway/utility role | 🟢 |
| `specs/Console/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/domain.md#Console` | Console project, redirects | 🟢 |
| `specs/architecture.md#Integrations` | Assistant microservice | 🟢 |
| `specs/security.md` | Security headers, console guard | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Console user | Access console UI | GET / redirects to /console |
| Console frontend | Adapt UI to server capabilities | GET /v1/console/variables |
| Console user | Configure custom domains | GET /v1/console/resources validates availability |
| Console user | Ask AI assistant for help | POST /v1/console/assistant (SSE stream) |

## 4. Business rules

1. **BR-01:** `API` init action guards the `console` group: blocks requests unless `$project->getId() === 'console'`. 🟢
2. **BR-02:** `Web` init action injects browser security headers: `X-Frame-Options: SAMEORIGIN`, `X-XSS-Protection`, `X-UA-Compatible`. 🟢
3. **BR-03:** All redirect actions share one base class preserving path, query string, and fragment when mapping to /console. 🟢
4. **BR-04:** Assistant endpoint proxies to `http://appwrite-assistant:3003/v1/models/assistant/prompt`, streams SSE via `$response->chunk()`, rate-limited to 15 requests/user (`abuse-limit: 15`, `abuse-key: userId:{userId}`). 🟢
5. **BR-05:** `getVariables` collects env config + DB adapter capabilities (relationship/operator/spatial/fulltext support) + feature flags (VCS, domain forwarding, assistant). 🟢
6. **BR-06:** `getResource` validates domain availability: structure restrictions (subdomain depth), deny list `_APP_CUSTOM_DOMAIN_DENY_LIST`, reserved platform hostnames, uniqueness in `dbForPlatform.rules`. 🟢
7. **BR-07:** `getConsoleEmailTemplate` always returns built-in default template; locale falls back to English. 🟢
8. **BR-08:** `listOAuth2Providers` excludes mock and disabled providers; sandbox included. 🟡
9. **BR-09:** Scope lists are static — read from `projectScopes`/`organizationScopes` config params, no DB. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Return console config variables | Must | Env + adapter capabilities + feature flags returned | 🟢 |
| RF-02 | Render default email template | Must | Built-in template by type+locale; English fallback | 🟢 |
| RF-03 | List enabled OAuth2 providers | Must | Provider metadata; mock/disabled excluded | 🟢 |
| RF-04 | List project scopes | Must | Static list from config | 🟢 |
| RF-05 | List organization scopes | Must | Static list from config | 🟢 |
| RF-06 | Proxy prompt to AI assistant (SSE) | Must | Streamed response; 15/user rate limit | 🟢 |
| RF-07 | Check resource (domain) availability | Must | Enforce structure, deny list, reserved hosts, uniqueness | 🟢 |
| RF-08 | Redirect legacy paths to /console | Must | Preserve path/query/fragment | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Console group blocked outside console project | `API` init guard | 🟢 |
| Security | Web security headers | `Web` init action | 🟢 |
| Performance | Assistant rate limit 15/user | abuse config | 🟢 |
| Performance | getResource DB uniqueness check | dbForPlatform.rules query | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Console API blocked under non-console project
  Given a request targeting a project other than "console"
  When a console-group route is called
  Then access is denied

Scenario: Legacy redirect preserves query
  Given GET /auth?next=/console
  When the redirect action runs
  Then the client is redirected to /console/auth?next=/console

Scenario: Assistant rate limit
  Given a user who sent 15 assistant queries
  When an additional query is sent
  Then the request is rejected (rate limited)

Scenario: Domain already in use
  Given a domain already present in dbForPlatform.rules
  When GET /v1/console/resources is called with that domain
  Then the resource is reported unavailable
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01 (Variables) | Must | Frontend capability adaptation |
| RF-02 (Email template) | Must | Console template preview |
| RF-03 (OAuth2 providers) | Must | Auth setup UX |
| RF-04/05 (Scopes) | Must | Key creation UX |
| RF-06 (Assistant) | Should | Requires assistant service |
| RF-07 (Resource check) | Must | Domain configuration |
| RF-08 (Redirects) | Must | Legacy URL compatibility |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Assistant service response format (SSE schema) not verified in module.
- 🟡 [INFERRED] Full deny list contents of `_APP_CUSTOM_DOMAIN_DENY_LIST` not enumerated.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
