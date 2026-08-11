# Design: Advisor

> Unit: `Advisor`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Advisor module is a **read-mostly HTTP service** over the platform database. Reports and Insights live in `dbForPlatform` (multi-tenant control plane). The module is purely reactive — it never creates report/insight documents; producers are external workers. Deletion is split: sync delete of the report + async enqueue of insight cleanup.

```
Console / management client
  │ ADMIN|KEY auth
  ▼
Traefik → app/http.php → Advisor\Http\* (5 actions)
  │
  ├── dbForPlatform (reports, insights collections)
  ├── publisherForDeletes → Redis → worker-deletes
  └── queueForEvents → Redis → worker-events
```

## 2. Route-to-action map

| Method | Path | Action class | Action ID | Groups | Scope |
|--------|------|-------------|-----------|--------|-------|
| GET | `/v1/reports` | `Reports/XList` | `listReports` | `api`, `reports` | `reports.read` |
| GET | `/v1/reports/:reportId` | `Reports/Get` | `getReport` | `api`, `reports` | `reports.read` |
| DELETE | `/v1/reports/:reportId` | `Reports/Delete` | `deleteReport` | `api`, `reports` | `reports.write` |
| GET | `/v1/reports/:reportId/insights` | `Insights/XList` | `listInsights` | `api`, `insights` | `insights.read` |
| GET | `/v1/reports/:reportId/insights/:insightId` | `Insights/Get` | `getInsight` | `api`, `insights` | `insights.read` |

## 3. Domain enums

| Enum | Values |
|------|--------|
| `ReportType` | `lighthouse`, `audit`, `databaseAnalyzer` |
| `InsightType` | `databaseIndex`, `tablesDBIndex`, `documentsDBIndex`, `vectorsDBIndex`, `databasePerformance`, `sitePerformance`, `siteAccessibility`, `siteSeo`, `functionPerformance` |
| `InsightSeverity` | `info`, `warning`, `critical` |
| `InsightStatus` | `active`, `dismissed` |
| `InsightCTAMethod` | `createIndex` |
| `InsightCTAService` | `databases`, `tablesDB`, `documentsDB`, `vectorsDB` |

## 4. Data flow — key actions

### 4.1 List reports (`Reports/XList`)

```
1. Authenticate (ADMIN|KEY), resolve project
2. Build query: { projectInternalId = project.getSequence() }
3. Cursor validation (project-scoped cursor)
4. Fetch report page (skip insights sub-query)
5. Bulk-fetch insights for page IDs (grouped by reportInternalId, single query)
6. Distribute insights to each report
7. Optional COUNT (total=true)
8. Return reportList with nested insights
```

### 4.2 Get report (`Reports/Get`)

```
1. Authenticate, resolve project
2. Fetch report by reportId + projectInternalId (skip insights sub-query)
3. Re-fetch insights for that report explicitly
4. Return report with embedded insights
```

### 4.3 Delete report (`Reports/Delete`)

```
1. Authenticate, resolve project
2. Fetch report scoped to project
3. deleteDocument on reports collection (synchronous)
4. Enqueue DeleteMessage (DELETE_TYPE_REPORT, reportId)
5. Emit reports.[reportId].delete event
6. Return 204
```

### 4.4 List insights (`Insights/XList`)

```
1. Authenticate, resolve project
2. Verify parent report exists (scoped)
3. Query insights: { reportInternalId = report.getSequence() }
4. Cursor pagination (project-scoped validation)
5. Return insightList
```

### 4.5 Get insight (`Insights/Get`)

```
1. Authenticate, resolve project
2. Fetch insight by insightId + reportInternalId + projectInternalId
3. Return insight or 404
```

## 5. Dependency injection

| Injection | Role | Used by |
|-----------|------|---------|
| `response` | HTTP response builder | All actions |
| `project` | Current project (scoping) | All actions |
| `dbForPlatform` | Platform-scoped database | All actions |
| `publisherForDeletes` | Async delete enqueue | Reports/Delete |
| `queueForEvents` | Event queue | Reports/Delete |

## 6. Events emitted

| Event pattern | Trigger |
|---------------|---------|
| `reports.[reportId].delete` | Report deleted |

## 7. Error handling

| Error code | HTTP | Trigger |
|------------|------|---------|
| `REPORT_NOT_FOUND` | 404 | Report ID not found in project scope |
| `INSIGHT_NOT_FOUND` | 404 | Insight ID not found in report scope |
| `GENERAL_CURSOR_NOT_FOUND` | 400 | Cursor from another project or invalid |
| Unauthorized | 401 | Missing/invalid ADMIN or KEY auth |

## 8. Security

- **Auth type:** ADMIN or KEY only — no session/JWT.
- **Scoping:** Every query filtered by `projectInternalId`; insights by `reportInternalId`.
- **Cursor validation:** Cursor must belong to current project before use.
- **Async deletion:** Insights removed by worker, decoupled from API response.
