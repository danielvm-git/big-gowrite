# Contracts: Advisor

> Unit: `Advisor`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base path: `/v1`. Auth: ADMIN or KEY only (no session/JWT).

| Method | Path | Query params | Success | Errors |
|--------|------|--------------|---------|--------|
| GET | `/reports` | `limit`, `cursor`, `total` (default true) | 200 reportList | 400 cursor, 401 |
| GET | `/reports/:reportId` | — | 200 report | 404 |
| DELETE | `/reports/:reportId` | — | 204 | 404, 401 |
| GET | `/reports/:reportId/insights` | `limit`, `cursor`, `total` | 200 insightList | 404 report |
| GET | `/reports/:reportId/insights/:insightId` | — | 200 insight | 404 |

## Domain values (enums)

| Field | Allowed values |
|-------|----------------|
| `report.type` | `lighthouse`, `audit`, `databaseAnalyzer` |
| `insight.type` | `databaseIndex`, `tablesDBIndex`, `documentsDBIndex`, `vectorsDBIndex`, `databasePerformance`, `sitePerformance`, `siteAccessibility`, `siteSeo`, `functionPerformance` |
| `insight.severity` | `info`, `warning`, `critical` |
| `insight.status` | `active`, `dismissed` |
| `insight.cta.method` | `createIndex` |
| `insight.cta.service` | `databases`, `tablesDB`, `documentsDB`, `vectorsDB` |

## Response schemas

### Report (embedded in list/get)
```json
{
  "$id": "report_abc",
  "$createdAt": "...",
  "$updatedAt": "...",
  "type": "lighthouse",
  "status": "...",
  "insights": [ { "type": "sitePerformance", "severity": "critical", "status": "active" } ]
}
```

### Insight
```json
{
  "$id": "insight_abc",
  "$createdAt": "...",
  "$updatedAt": "...",
  "type": "databaseIndex",
  "severity": "warning",
  "status": "active",
  "cta": { "method": "createIndex", "service": "databases" }
}
```

## Queue contract (deletes)

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForDeletes` | Advisor → worker-deletes | `DeleteMessage` with `DELETE_TYPE_REPORT` + reportId (async insight cleanup) |

## Event contract

| Event | Schema |
|-------|--------|
| `reports.[reportId].delete` | `{ reportId, projectInternalId, timestamp }` |
