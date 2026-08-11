# Contracts: Organization

> Unit: `Organization`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/organization/projects`. Auth: ADMIN | KEY | ORGANIZATION. Scopes: `projects.*`.

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `/projects` | `projectId`, `name`, `region`, `auth` | 201 project | PROJECT_RESERVED_PROJECT, 400 |
| GET | `/projects` | `select`, `cursor`, `limit` | 200 projectList | 401 |
| GET | `/projects/:projectId` | — | 200 project | PROJECT_NOT_FOUND |
| PATCH | `/projects/:projectId` | `name` | 200 project | PROJECT_NOT_FOUND |
| DELETE | `/projects/:projectId` | — | 204 | PROJECT_NOT_FOUND |

## Queue contract

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForDeletes` | Delete action → worker-deletes | `DeleteMessage`: type `DELETE_TYPE_DOCUMENT`, projectId (async cleanup of project resources) |

## Configuration contract

| Env var | Effect |
|---------|--------|
| `_APP_PROJECT_REGIONS` | Region allow-list (empty = all) |
| `_APP_DATABASE_KEYS` | DSN pool for region filtering |
| `_APP_DATABASE_OVERRIDE` | Force a specific DSN |
| `_APP_DATABASE_SHARED_TABLES` | DSNs joining shared schema |
| `_APP_DATABASE_SHARED_NAMESPACE` | Shared namespace name |

## Hook contract

| Hook | Timing |
|------|--------|
| `afterProjectCreation` | Fired after project doc created (migration mirroring) |

## Document model (`projects` in dbForPlatform)

| Field | Type | Notes |
|-------|------|-------|
| `teamInternalId` | int | Ownership (team sequence) |
| `region` | string | Selected region |
| `auth` | array | Seeded auth methods/config |
| `databaseInternalId` | int | Provisioned DB reference |
