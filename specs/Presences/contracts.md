# Contracts: Presences

> Unit: `Presences`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/presences`. Scopes: `presences.read` / `presences.write`.

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| PUT | `/presences/:presenceId` | `status`, `metadata`, `expiresAt`, `userId` (key only), `permissions` | 200 presence | 400 invalid expiry |
| GET | `/presences/:presenceId` | — | 200 presence | PRESENCE_NOT_FOUND |
| GET | `/presences` | `queries`, `ttl` (0–86400), `cursor`, `limit` | 200 presenceList + `X-Appwrite-Cache` | 400 |
| PATCH | `/presences/:presenceId` | `status`, `metadata`, `expiresAt`, `purge` | 200 presence | PRESENCE_NOT_FOUND |
| DELETE | `/presences/:presenceId` | `purge` | 204 | PRESENCE_NOT_FOUND |

## Headers

| Header | Direction | Meaning |
|--------|-----------|---------|
| `x-appwrite-source` | Request | `rest` or `graphql` → stored in `source` |
| `X-Appwrite-Cache` | Response | `hit` / `miss` on TTL-cached lists |

## Document model (`presenceLogs` in dbForProject)

| Field | Type | Notes |
|-------|------|-------|
| `userId` | string | Owner identity |
| `status` | string | Free-form status |
| `metadata` | object | Arbitrary metadata |
| `source` | string | `rest` / `graphql` |
| `expiresAt` | datetime | Default +15 min; max +30 days |
| `permissions` | array | Default owner full access |

## Event contract

- `presences.[presenceId].upsert`
- `presences.[presenceId].update`
- `presences.[presenceId].delete` (payload: deleted presence doc)

## Metrics

| Metric | Direction |
|--------|-----------|
| `METRIC_USERS_PRESENCE` | +1 on upsert/create, −1 on delete |
