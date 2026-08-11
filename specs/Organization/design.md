# Design: Organization

> Unit: `Organization`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Organization module is an **HTTP service for team-scoped project management**. It validates `$team` via an Init hook, then provisions project resources (database namespace/shared schema, audit schema, collection schemas, auth methods) at creation. Deletion enqueues async cleanup.

```
Client
  │ scopes: projects.* | auth: ADMIN | KEY | ORGANIZATION
  ▼
app/http.php → Organization\Http\** (5 actions) + Init
  │
  ├── dbForPlatform (projects collection + provisioning)
  ├── pools (connection pool → new Database)
  ├── cache
  ├── hooks (afterProjectCreation)
  └── publisherForDeletes (DELETE_TYPE_DOCUMENT)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/organization/projects` | `createOrganizationProject` | projects.write |
| GET | `/v1/organization/projects` | `listOrganizationProjects` | projects.read |
| GET | `/v1/organization/projects/:projectId` | `getOrganizationProject` | projects.read |
| PATCH | `/v1/organization/projects/:projectId` | `updateOrganizationProject` | projects.write |
| DELETE | `/v1/organization/projects/:projectId` | `deleteOrganizationProject` | projects.write |

## 3. Data flows — key actions

### 3.1 Create project (provisioning)

```
1. Init: $team present
2. projectId != 'console'
3. Region selection:
   - _APP_PROJECT_REGIONS set → validate region in list
   - region != default → filter _APP_DATABASE_KEYS by region string
4. DSN choice: _APP_DATABASE_OVERRIDE | random from pool
5. If DSN in _APP_DATABASE_SHARED_TABLES:
   → join shared schema (appwrite, _APP_DATABASE_SHARED_NAMESPACE)
   Else:
   → create namespace `_` + project sequence
   → setup audit schema
   → bootstrap all projects collections from config
6. Seed auth methods from `auth` config
7. Create projects doc (teamInternalId = team sequence)
8. Trigger afterProjectCreation hook
9. Return project
```

### 3.2 Get/Update/Delete — ownership

```
1. Fetch project
2. project.teamInternalId != team.getSequence() → PROJECT_NOT_FOUND
3. Update: sparse name change
4. Delete: delete doc (authorization skip) + enqueue DELETE_TYPE_DOCUMENT
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForPlatform` | projects collection + provisioning |
| `team` | Owning team (scoping + ownership) |
| `cache` | Passed to new Database instance |
| `pools` | Adapter for project-table provisioning |
| `hooks` | afterProjectCreation trigger |
| `publisherForDeletes` | DELETE_TYPE_DOCUMENT enqueue |
| `authorization` | Skip ACL on delete |
| `response` | HTTP response |

## 5. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `PROJECT_RESERVED_PROJECT` | 400 | projectId = console |
| `PROJECT_NOT_FOUND` | 404 | Cross-team access |
| Unauthorized | 401/403 | Empty team / bad auth |

## 6. Security

- Team ownership enforced on get/update/delete.
- Reserved console project id.
- Region/DSN selection controlled by env config.
