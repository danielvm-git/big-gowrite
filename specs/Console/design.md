# Design: Console

> Unit: `Console`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Console module is a **gateway + utility HTTP service** with two init hooks (group guards) and no workers/tasks. Its routes split into: console API endpoints (`/v1/console/**`), redirects (`/`, `/auth/*`, etc.), and the AI assistant SSE proxy.

```
Browser
  │
  ▼
Traefik → app/http.php
  │
  ├── Init/API.php  → guards `console` group (project == 'console')
  ├── Init/Web.php  → security headers for `web` group
  │
  ├── /v1/console/**  → Variables, Templates, OAuth2, Scopes, Assistant, Resources
  └── /, /auth/*, /login, ... → Redirects/Base.php → /console/...
```

## 2. Route-to-action map

| Method | Path | Action class | Action ID | Scope |
|--------|------|-------------|-----------|-------|
| GET | `/v1/console/variables` | `Variables/Get` | `getVariables` | projects.read |
| GET | `/v1/console/templates/email/:templateId` | `Templates/Email/Get` | `getConsoleEmailTemplate` | public |
| GET | `/v1/console/oauth2-providers` | `OAuth2Providers/XList` | `listOAuth2Providers` | public |
| GET | `/v1/console/scopes/project` | `Scopes/Project/XList` | `listConsoleProjectScopes` | public |
| GET | `/v1/console/scopes/organization` | `Scopes/Organization/XList` | `listConsoleOrganizationScopes` | public |
| POST | `/v1/console/assistant` | `Assistant/Create` | `createAssistantQuery` | assistant.read |
| GET | `/v1/console/resources` | `Resources/Get` | `getResource` | rules.read |
| GET | `/` | `Redirects/Root/Get` | `consoleRedirectRoot` | home |
| GET | `/auth/*` | `Redirects/Auth/Get` | `consoleRedirectAuth` | home |
| GET | `/invite` | `Redirects/Invite/Get` | `consoleRedirectInvite` | home |
| GET | `/login` | `Redirects/Login/Get` | `consoleRedirectLogin` | home |
| GET | `/mfa/*` | `Redirects/MFA/Get` | `consoleRedirectMFA` | home |
| GET | `/card` | `Redirects/Card/Get` | `consoleRedirectCard` | home |
| GET | `/recover` | `Redirects/Recover/Get` | `consoleRedirectRecover` | home |
| GET | `/register` | `Redirects/Register/Get` | `consoleRedirectRegister` | home |

## 3. Init hooks

| Hook | Group | Behavior |
|------|-------|----------|
| `Init/API.php` | `console` | Reject unless `$project->getId() === 'console'` |
| `Init/Web.php` | `web` | Inject X-Frame-Options, X-XSS-Protection, X-UA-Compatible |

## 4. Data flow — key actions

### 4.1 getVariables

```
1. Collect env config (VCS, domains, assistant flags)
2. Detect DB adapter capabilities (relationships, operators, spatial, fulltext)
3. Assemble variables payload
4. Return JSON
```

### 4.2 createAssistantQuery (SSE proxy)

```
1. Rate limit 15/user (abuse-key userId)
2. Forward prompt to appwrite-assistant:3003/v1/models/assistant/prompt
3. Stream SSE response via $response->chunk()
```

### 4.3 getResource

```
1. Validate domain structure (subdomain depth for sites/functions)
2. Check deny list (_APP_CUSTOM_DOMAIN_DENY_LIST)
3. Check reserved platform hostnames
4. Query dbForPlatform.rules for uniqueness (authorization skip)
5. Return { available: bool }
```

## 5. Dependency injection

| Injection | Role | Used by |
|-----------|------|---------|
| `response`, `request` | HTTP layer | All actions |
| `platform` | Hostnames, sites/functions domains, consoleHostname | Variables, Resources |
| `dbForPlatform` | Rules lookup | Resources |
| `dbForProject` | Adapter capability detection | Variables |
| `authorization` | Skip permissions in rule uniqueness | Resources |
| `vcsProviders` | Configured VCS provider keys | Variables |
| `project` | Current project (console guard) | Init/API |
| `localeCodes` | Valid locale codes | Email template |

## 6. Security

- **Console guard:** console-group endpoints only callable under the `console` project.
- **Web headers:** SAMEORIGIN frame policy + XSS protections.
- **Assistant rate limit:** 15/user.
- **Domain validation:** deny list + reserved hostnames + uniqueness.
