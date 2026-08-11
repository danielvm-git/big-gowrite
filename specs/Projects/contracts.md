# Contracts: Projects

> Unit: `Projects`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/projects`. Scopes: `projects.*`, `devKeys.*`, `schedules.*`, `stages.*`.

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `/projects` | `projectId`, `name`, `teamId`, `region` | 201 project | PROJECT_RESERVED_PROJECT, TEAM_NOT_FOUND |
| PATCH | `/projects/:projectId` | `name`, `description`, `logo`, `url` | 200 project | 404 |
| GET | `/projects` | `teamId`, `cursor`, `limit` | 200 projectList | 401 |
| PATCH | `/projects/:projectId/team` | `teamId` | 200 project | TEAM_NOT_FOUND |
| POST | `/projects/:projectId/dev-keys` | `name`, `expire` | — | DEV_KEY_GONE (410) |
| PUT | `/projects/:projectId/dev-keys/:keyId` | `name`, `expire` | 200 devKey | 404 |
| GET | `/projects/:projectId/dev-keys/:keyId` | — | 200 devKey | 404 |
| GET | `/projects/:projectId/dev-keys` | — | 200 devKeyList | 401 |
| DELETE | `/projects/:projectId/dev-keys/:keyId` | — | 204 | 404 |
| POST | `/projects/:projectId/schedules` | `scheduleId`, `cron`, `resourceType`, `resourceId` | 201 schedule | 400 resource |
| GET | `/projects/:projectId/schedules/:scheduleId` | — | 200 schedule | 404 |
| GET | `/projects/:projectId/schedules` | `cursor`, `limit` | 200 scheduleList | 401 |
| GET | `/projects/:projectId/stages` | — | 200 stageList | 401 |
| PATCH | `/projects/:projectId/stages/:stageId` | `skipped` | 200 stage | 400 completed |

## Event contract

- `schedules.[scheduleId].create`

## Configuration contract

| Env var | Effect |
|---------|--------|
| `_APP_PROJECT_REGIONS` | Region allow-list |
| `_APP_DATABASE_KEYS` | DSN pool |
| `_APP_DATABASE_OVERRIDE` | Forced DSN |

## Provisioning contract

| Mode | Schema | Collections |
|------|--------|-------------|
| Shared tables | `appwrite` (shared namespace) | none per-project |
| Dedicated | `_<internalId>` | all project collections + audit schema |

## Hook contract

| Hook | Timing |
|------|--------|
| `afterProjectCreation` | After project doc creation (migration mirroring) |
