# Requirements: Advisor

> Identifier: `ADV-REPORTS`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Advisor module exposes a read-mostly API for surfacing automated analysis results about a project's health, performance, and configuration. It operates on two entities: **Reports** (top-level results of analysis runs — Lighthouse audits, database-analyzer scans) and **Insights** (actionable findings nested inside reports, each with severity, type, status, and optional CTA). Clients list/fetch reports and insights, and can delete reports (which enqueues async deletion of related insights). Reports and insights are never created by this module — they are produced by background workers or external processes writing to the platform database.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Advisor` | Module purpose, entities | 🟢 |
| `specs/Advisor/legacy-mapping.md` | Endpoints, business rules, enums | 🟢 |
| `specs/domain.md#Analysis` | Report/Insight domain model | 🟢 |
| `specs/architecture.md#PlatformDB` | Reports/Insights live in platform DB | 🟢 |
| `specs/workers.md#Deletes` | Async delete worker handling `DELETE_TYPE_REPORT` | 🟡 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Console admin | Review project health | List reports, open a report, inspect critical insights |
| Console admin | Act on a recommendation | Follow an insight CTA to create a database index |
| Console admin | Clean up stale analyses | Delete an old report (async cleanup of insights) |

## 4. Business rules

1. **BR-01:** All queries are scoped to the current project via `projectInternalId == project.getSequence()` — prevents cross-project leakage. 🟢
2. **BR-02:** Insights are additionally scoped to their parent report via `reportInternalId == report.getSequence()`. 🟢
3. **BR-03:** `getReport` and `deleteReport` skip the insights sub-query when fetching the report doc, then re-fetch insights explicitly for the response payload. 🟢
4. **BR-04:** `listReports` bulk-fetches all insights for the page in a single grouped query (avoids N+1). 🟢
5. **BR-05:** Cursor pagination validates the cursor belongs to the current project; invalid cursors throw `GENERAL_CURSOR_NOT_FOUND`. 🟢
6. **BR-06:** `deleteReport` deletes the report doc synchronously, then enqueues `DeleteMessage` (`DELETE_TYPE_REPORT`) for async insight cleanup. 🟢
7. **BR-07:** Authentication requires `ADMIN` or `KEY` auth type — management-level API, no session/JWT. 🟢
8. **BR-08:** `total` param defaults to `true` on list endpoints; `false` skips COUNT query for performance. 🟡

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | List paginated reports for current project | Must | Returns page of reports with nested insights; project-scoped; cursor pagination | 🟢 |
| RF-02 | Get a single report with insights embedded | Must | Returns report + insights; scoped by project and report ID | 🟢 |
| RF-03 | Delete a report | Must | Deletes report sync, enqueues async insight deletion, emits delete event | 🟢 |
| RF-04 | List insights for a report (cursor pagination) | Must | Returns insight page; scoped to parent report | 🟢 |
| RF-05 | Get a single insight scoped to parent report | Must | Returns insight or 404 | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Performance | N+1 avoidance: insights bulk-fetched in one grouped query | `listReports` grouped fetch | 🟢 |
| Performance | `total=false` skips COUNT query | `total` param default `true` | 🟡 |
| Security | Project scoping on every query (platform DB, multi-tenant) | `projectInternalId` filter everywhere | 🟢 |
| Reliability | Async insight cleanup decoupled from sync delete | `publisherForDeletes` enqueue | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: List reports for project
  Given an admin with ADMIN/KEY auth for a project with 2 reports
  When GET /v1/reports is called
  Then paginated reports with embedded insights are returned, scoped to project

Scenario: Invalid cursor rejected
  Given an admin with a cursor from a different project
  When GET /v1/reports with that cursor is called
  Then a GENERAL_CURSOR_NOT_FOUND error is returned

Scenario: Delete report enqueues async cleanup
  Given an admin with a report containing insights
  When DELETE /v1/reports/:reportId is called
  Then the report doc is removed synchronously, a DeleteMessage is enqueued, and a delete event is emitted

Scenario: Insight scoped to wrong report
  Given an admin with an insight belonging to report A
  When GET /v1/reports/:reportB/insights/:insightId is called
  Then a 404 error is returned
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01 (List reports) | Must | Primary read path |
| RF-02 (Get report) | Must | Report detail with insights |
| RF-03 (Delete report) | Must | Management lifecycle |
| RF-04 (List insights) | Must | Insights per report |
| RF-05 (Get insight) | Must | Individual insight detail |
| Async cleanup (RNF) | Must | Avoid orphan insights |

## 9. Clarifications

> No clarification session recorded. Run `/reversa-clarify` when `[GAP]` pending.

## 10. Gaps

- 🔴 [GAP] Which workers/external processes produce reports and insights? Not visible in this module.
- 🟡 [INFERRED] Exact DeleteMessage payload fields for `DELETE_TYPE_REPORT` not fully verified.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
