# Contracts: Health

> Unit: `Health`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/health`. Scopes: `health.read` (default), `public` (version), `root` (stats).

| Method | Path | Params | Success | Errors |
|--------|------|--------|---------|--------|
| GET | `/health` | — | 200 `{ status, ping }` | — |
| GET | `/health/version` | — | 200 `{ version }` | — |
| GET | `/health/db` | — | 200 `{ status, ping }` | 503 on pool failure |
| GET | `/health/cache` | — | 200 `{ status, ping }` | 503 |
| GET | `/health/pubsub` | — | 200 `{ status, ping }` | 503 |
| GET | `/health/time` | — | 200 `{ status, ping }` | 503 drift |
| GET | `/health/certificate` | `url` | 200 `{ status, ping }` | HEALTH_CERTIFICATE_EXPIRED |
| GET | `/health/storage` | — | 200 `{ status, ping }` | 503 |
| GET | `/health/storage/local` | — | 200 `{ status, ping }` | 503 |
| GET | `/health/anti-virus` | — | 200 `{ status: "pass"|"disabled", ping }` | 503 |
| GET | `/health/geo` | — | 200 `{ status, ping }` | GENERAL_SERVER_ERROR |
| GET | `/health/stats` | — | 200 stats | 401 (root required) |
| GET | `/health/queue/audits` | `threshold` (5000) | 200 `{ size, threshold }` | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/webhooks` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/certificates` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/builds` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/databases` | `name` (default database_db_main), `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/deletes` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/mails` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/messaging` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/migrations` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/notifications` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/functions` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/stats-resources` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/stats-usage` | `threshold` | 200 | HEALTH_QUEUE_SIZE_EXCEEDED |
| GET | `/health/queue/failed/:name` | — | 200 `{ count }` | 400 invalid queue name |

## External probes

| Target | Protocol | Constraint |
|--------|----------|------------|
| `time.google.com:123` | raw UDP NTP | 5 s receive timeout; ±60 s drift |
| ClamAV | network ping | skipped when `_APP_STORAGE_ANTIVIRUS=disabled` |
| `_APP_GEO_ENDPOINT` | HTTP probe | missing endpoint → GENERAL_SERVER_ERROR |

## Configuration contract

| Env var | Effect |
|---------|--------|
| `_APP_STORAGE_ANTIVIRUS` | `disabled` → skip ClamAV |
| `_APP_GEO_ENDPOINT` | Geo probe target; empty → error |
| `_APP_*_QUEUE_NAME` + `Event::*_QUEUE_NAME` | Queue name WhiteList for `/queue/failed/:name` |
