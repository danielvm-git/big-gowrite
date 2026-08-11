# Design: Health

> Unit: `Health`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Health module is a **read-only diagnostics service** (no workers, no tasks, no events). Each endpoint probes one infrastructure component through injected adapters (pools, cache, devices, publishers) and returns `{ status, ping }`.

```
Admin / SRE
  │ health.read (or public / root)
  ▼
app/http.php → Health\Http\** (25 actions)
  │
  ├── $pools (DB, PubSub) — ping()
  ├── $cache (Redis) — ping()
  ├── $deviceFor{Files,Functions,Sites,Builds} — storage probe
  ├── $publisherFor* (13 queues) — queue size
  ├── external: time.google.com (UDP), ClamAV, geo service
  └── $register — raw Redis info for stats
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| GET | `/v1/health` | `getHealth` | health.read |
| GET | `/v1/health/version` | `getHealthVersion` | public |
| GET | `/v1/health/db` | `getHealthDB` | health.read |
| GET | `/v1/health/cache` | `getHealthCache` | health.read |
| GET | `/v1/health/pubsub` | `getHealthPubSub` | health.read |
| GET | `/v1/health/time` | `getHealthTime` | health.read |
| GET | `/v1/health/certificate` | `getHealthCertificate` | health.read |
| GET | `/v1/health/storage` | `getHealthStorage` | health.read |
| GET | `/v1/health/storage/local` | `getHealthStorageLocal` | health.read |
| GET | `/v1/health/anti-virus` | `getHealthAntiVirus` | health.read |
| GET | `/v1/health/geo` | `getHealthGeo` | health.read |
| GET | `/v1/health/stats` | `getHealthStats` | root |
| GET | `/v1/health/queue/audits` | `getQueueAudits` | health.read |
| GET | `/v1/health/queue/webhooks` | `getQueueWebhooks` | health.read |
| GET | `/v1/health/queue/certificates` | `getQueueCertificates` | health.read |
| GET | `/v1/health/queue/builds` | `getQueueBuilds` | health.read |
| GET | `/v1/health/queue/databases` | `getQueueDatabases` | health.read |
| GET | `/v1/health/queue/deletes` | `getQueueDeletes` | health.read |
| GET | `/v1/health/queue/mails` | `getQueueMails` | health.read |
| GET | `/v1/health/queue/messaging` | `getQueueMessaging` | health.read |
| GET | `/v1/health/queue/migrations` | `getQueueMigrations` | health.read |
| GET | `/v1/health/queue/notifications` | `getQueueNotifications` | health.read |
| GET | `/v1/health/queue/functions` | `getQueueFunctions` | health.read |
| GET | `/v1/health/queue/stats-resources` | `getQueueStatsResources` | health.read |
| GET | `/v1/health/queue/stats-usage` | `getQueueUsage` | health.read |
| GET | `/v1/health/queue/failed/:name` | `getFailedJobs` | health.read |

## 3. Data flow — queue check (shared)

```
1. Resolve queue size from publisher (Queue)
2. assertQueueThreshold(size, threshold=5000)
3. size >= threshold → throw HEALTH_QUEUE_SIZE_EXCEEDED
4. Else return { size, threshold }
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `$response` | HTTP response |
| `$pools` | DB + PubSub pool groups |
| `$cache` | Redis cache adapter |
| `$authorization` | DB auth context for pool ping |
| `$register` | Registry (raw Redis info for stats) |
| `$deviceForFiles`, `$deviceForFunctions`, `$deviceForSites`, `$deviceForBuilds` | Storage devices |
| `$publisherForAudits` … `$publisherForUsage` (13) | Queue publishers |
| `$queueForWebhooks` | Webhook queue (non-publisher pattern) |

## 5. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `HEALTH_QUEUE_SIZE_EXCEEDED` | 400/503 | queue size ≥ threshold |
| `HEALTH_CERTIFICATE_EXPIRED` | 400/503 | expired TLS cert |
| `GENERAL_SERVER_ERROR` | 500 | geo endpoint unset; infra probe failures |
| Unauthorized | 401/403 | missing health.read / root |

## 6. Security

- `health.read` scope on most endpoints; `root` on stats; `public` on version.
- Stats intentionally excluded from SDK docs (`docs: false`).
