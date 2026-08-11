# Design: Projects

> Unit: `Projects`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Projects module is a **platform-layer HTTP service** managing top-level project entities in `dbForPlatform`. It shares provisioning logic patterns with the Organization module but operates on the platform surface (`/v1/projects`).

```
Client
  │ scopes: projects.*, devKeys.*, schedules.*, stages.*
  ▼
app/http.php → Projects\Http\** (14 actions)
  │
  ├── dbForPlatform (projects, devKeys, schedules)
  ├── dbForProject / getProjectDB (schedule resource validation)
  ├── pools + cache (provisioning)
  ├── hooks (afterProjectCreation)
  └── queueForEvents (schedule events)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/projects` | `createProject` | projects.write |
| PATCH | `/v1/projects/:projectId` | `updateProject` | projects.write |
| GET | `/v1/projects` | `listProjects` | projects.read |
| PATCH | `/v1/projects/:projectId/team` | `updateProjectTeam` | projects.write |
| POST | `/v1/projects/:projectId/dev-keys` | `createDevKey` | devKeys.write |
| PUT | `/v1/projects/:projectId/dev-keys/:keyId` | `updateDevKey` | devKeys.write |
| GET | `/v1/projects/:projectId/dev-keys/:keyId` | `getDevKey` | devKeys.read |
| GET | `/v1/projects/:projectId/dev-keys` | `listDevKeys` | devKeys.read |
| DELETE | `/v1/projects/:projectId/dev-keys/:keyId` | `deleteDevKey` | devKeys.write |
| POST | `/v1/projects/:projectId/schedules` | `createSchedule` | schedules.write |
| GET | `/v1/projects/:projectId/schedules/:scheduleId` | `getSchedule` | schedules.read |
| GET | `/v1/projects/:projectId/schedules` | `listSchedules` | schedules.read |
| GET | `/v1/projects/:projectId/stages` | `listStages` | stages.read |
| PATCH | `/v1/projects/:projectId/stages/:stageId` | `updateStage` | stages.write |

## 3. Data flows — key actions

### 3.1 Create project (provisioning)

```
1. projectId != console; team exists
2. Region validation (_APP_PROJECT_REGIONS)
3. DSN selection: override → region-filtered → random
4. Shared tables → appwrite schema; else dedicated namespace `_<internalId>` + collections + audit schema
5. Create projects doc (status=active, teamInternalId, databaseInternalId)
6. Trigger afterProjectCreation hook
```

### 3.2 Schedule create

```
1. Validate CRON + resource ref (function/execution/message)
2. Resource exists in project DB?
3. Create schedules doc in dbForPlatform
4. Emit schedules.[scheduleId].create
```

### 3.3 Team transfer

```
1. Validate target team
2. Update teamInternalId on project
3. Propagate permissions to installations, repositories, VCS comments
```

### 3.4 Stage update

```
1. Load stage from onboarding config
2. Resolve actor type (key-account/key-organization/key-project/admin/user/guest)
3. completed → cannot overwrite as skipped
4. Persist skip with actor tag
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForPlatform` | projects, devKeys, schedules |
| `dbForProject` / `getProjectDB` | Schedule resource validation |
| `cache` / `pools` | Provisioning |
| `hooks` | afterProjectCreation |
| `queueForEvents` | schedule create event |
| `response`, `request` | HTTP layer |
| `team` | List scoping |
| `apiKey`, `user`, `mode` | Stage actor resolution |

## 5. Events emitted

- `schedules.[scheduleId].create`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `PROJECT_RESERVED_PROJECT` | 400 | projectId = console |
| `TEAM_NOT_FOUND` | 404 | Owning team missing |
| `DEV_KEY_GONE` | 410 | createDevKey (tombstone) |
| Schedule resource missing | 400 | Validation failure |
