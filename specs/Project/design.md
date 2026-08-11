# Design: Project

> Unit: `Project`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Project module is an **HTTP configuration service** for the current project. An Init hook guards the project group (blocks console/missing projects). Reads/writes split across `dbForPlatform` (projects, keys, platforms, oAuthProviders) and `dbForProject` (variables). Deletion and SMTP tests use async queues.

```
Client
  │ scopes: project.*, keys.*, platforms.*, policies.*, mocks.*, templates.*, oauth2.*
  ▼
app/http.php → Project\Http\** + Init
  │
  ├── dbForPlatform (projects, keys, platforms, oAuthProviders)
  ├── dbForProject (variables, functions, sites)
  ├── publisherForDeletes (project cleanup)
  ├── publisherForMails (SMTP test)
  └── queueForEvents (config events)
```

## 2. Route-to-action map (summary — 60+ actions)

| Method | Path | Action ID |
|--------|------|-----------|
| GET | `/v1/project` | `getProject` |
| DELETE | `/v1/project` | `deleteProject` |
| PUT | `/v1/project/labels` | `updateProjectLabels` |
| PATCH | `/v1/project/protocols/:protocolId` | `updateProjectProtocol` |
| PATCH | `/v1/project/services/:serviceId` | `updateProjectService` |
| PATCH | `/v1/project/auth-methods/:methodId` | `updateProjectAuthMethod` |
| PATCH | `/v1/project/smtp` | `updateProjectSMTP` |
| POST | `/v1/project/smtp/tests` | `createProjectSMTPTest` |
| GET | `/v1/project/templates/email` | `listProjectEmailTemplates` |
| GET | `/v1/project/templates/email/:templateId` | `getProjectEmailTemplate` |
| PATCH | `/v1/project/templates/email` | `updateProjectEmailTemplate` |
| POST | `/v1/project/variables` | `createProjectVariable` |
| GET | `/v1/project/variables` | `listProjectVariables` |
| GET | `/v1/project/variables/:variableId` | `getProjectVariable` |
| PUT | `/v1/project/variables/:variableId` | `updateProjectVariable` |
| DELETE | `/v1/project/variables/:variableId` | `deleteProjectVariable` |
| POST | `/v1/project/keys` | `createProjectKey` |
| POST | `/v1/project/keys/ephemeral` | `createEphemeralProjectKey` |
| GET | `/v1/project/keys` | `listProjectKeys` |
| GET | `/v1/project/keys/:keyId` | `getProjectKey` |
| PUT | `/v1/project/keys/:keyId` | `updateProjectKey` |
| DELETE | `/v1/project/keys/:keyId` | `deleteProjectKey` |
| GET | `/v1/project/platforms` | `listProjectPlatforms` |
| GET | `/v1/project/platforms/:platformId` | `getProjectPlatform` |
| POST | `/v1/project/platforms/{web,apple,android,windows,linux}` | `createProject{Type}Platform` |
| PATCH | `/v1/project/platforms/{type}/:platformId` | `updateProject{Type}Platform` |
| DELETE | `/v1/project/platforms/:platformId` | `deleteProjectPlatform` |
| POST | `/v1/project/mock-phones` | `createProjectMockPhone` |
| GET | `/v1/project/mock-phones` | `listProjectMockPhones` |
| GET | `/v1/project/mock-phones/:number` | `getProjectMockPhone` |
| PATCH | `/v1/project/mock-phones/:number` | `updateProjectMockPhone` |
| DELETE | `/v1/project/mock-phones/:number` | `deleteProjectMockPhone` |
| GET | `/v1/project/policies` | `listPolicies` |
| GET | `/v1/project/policies/:policyId` | `getProjectPolicy` |
| PATCH | `/v1/project/policies/{12 policy paths}` | `updateProject{Name}Policy` |
| GET | `/v1/project/oauth2` | `listProjectOAuth2` |
| GET | `/v1/project/oauth2/:provider` | `getProjectOAuth2` |
| PATCH | `/v1/project/oauth2/{provider}` | `updateProjectOAuth2{Provider}` (38) |


## 3. Data flows — key actions

### 3.1 API key create (standard + ephemeral)

```
Standard:
  1. Validate scopes + optional expiry
  2. Generate secret: API_KEY_STANDARD_ + 256B hex random
  3. Persist in dbForPlatform (resourceType=projects, resourceInternalId=project seq)
  4. Purge project cache; emit keys.[keyId].create

Ephemeral:
  1. Validate expiry ≤ 3600s
  2. Build in-memory Document; sign HS256 JWT (API_KEY_EPHEMERAL_)
  3. Return JWT secret only; nothing persisted
```

### 3.2 Project variable create/update

```
1. Validate key/value/secret
2. Persist in dbForProject (resourceType=project)
3. Mark all functions + sites live=false
4. Emit variables.[variableId].create/update
5. Secret → cannot unset (VARIABLE_CANNOT_UNSET_SECRET)
```

### 3.3 SMTP update

```
1. Validate fields
2. enabled=true → PHPMailer::SmtpConnect() (throw on failure)
   enabled=null + success → auto-enable
3. Persist SMTP config
4. Emit authMethod/services update events
```

### 3.4 OAuth2 provider update

```
1. Persist credentials in oAuthProviders map
2. Enable → verifyCredentials() if exists (throw on failure)
3. Secret masked on all reads
4. Emit oauth2.[providerId].update
```

### 3.5 Project deletion

```
1. Delete project doc (dbForPlatform)
2. Enqueue DELETE_TYPE_DOCUMENT via publisherForDeletes
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForPlatform` | projects, keys, platforms, oAuthProviders |
| `dbForProject` | variables, functions, sites |
| `response`, `request` | HTTP layer |
| `project` | Current project (Init guard) |
| `authorization` | skip() on platform DB writes |
| `queueForEvents` | Config events |
| `publisherForDeletes` | Project cleanup |
| `publisherForMails` | SMTP test email |
| `platform` | Config array (SMTP test) |
| `localeCodes` | Template locale validation |

## 5. Events emitted

- `keys.[keyId].{create,update,delete}`
- `variables.[variableId].{create,update}`
- `platforms.[platformId].{create,delete}`
- `mock-phones.[number].{create,update,delete}`
- `oauth2.[providerId].update`
- `authMethod.[methodId].update`
- `services.[serviceId].update`
- `protocols.[protocolId].update`
- `templates.[templateId].update`
- `projects.[projectId].policies.[policy].update`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `VARIABLE_CANNOT_UNSET_SECRET` | 400 | Downgrade secret variable |
| SMTP connection failure | 400 | enabled=true + connect fail |
| OAuth2 verify failure | 400 | Provider credentials invalid |
| Unauthorized | 401/403 | Init guard / scopes |
