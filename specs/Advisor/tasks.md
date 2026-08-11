# Tasks: Advisor

> Unit: `Advisor`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-ADV-01 — Register Advisor module boilerplate
- **Behavior source:** `Modules/Advisor/Module.php`
- **Description:** Create `Module.php` registering the Advisor HTTP service.
- **Done when:** Module loads and registers Http service.
- **Priority:** Must · Confidence: 🟢

### T-ADV-02 — Register the 5 HTTP routes + enums
- **Behavior source:** `Modules/Advisor/Services/Http.php`, `Enums/*.php`
- **Description:** Register listReports, getReport, deleteReport, listInsights, getInsight with correct scopes (`reports.read`, `reports.write`, `insights.read`). Define ReportType, InsightType, InsightSeverity, InsightStatus, InsightCTAMethod, InsightCTAService enums.
- **Done when:** Route table matches mapping; enums compile.
- **Priority:** Must · Confidence: 🟢

### T-ADV-03 — List reports with bulk insight sub-population
- **Behavior source:** `Modules/Advisor/Http/Reports/XList.php`
- **Description:** Implement project-scoped pagination with cursor validation; skip insights sub-query on report fetch; bulk-fetch insights grouped by reportInternalId; support `total` toggle.
- **Done when:** No N+1 (single grouped insights query); pagination + project scoping verified by E2E.
- **Priority:** Must · Confidence: 🟢

### T-ADV-04 — Get report with embedded insights
- **Behavior source:** `Modules/Advisor/Http/Reports/Get.php`
- **Description:** Fetch report scoped to project (skip insights sub-query), then re-fetch insights explicitly.
- **Done when:** Report returns with nested insights; wrong project → 404.
- **Priority:** Must · Confidence: 🟢

### T-ADV-05 — Delete report with async insight cleanup
- **Behavior source:** `Modules/Advisor/Http/Reports/Delete.php`
- **Description:** Sync delete of report doc; enqueue DeleteMessage `DELETE_TYPE_REPORT` via publisherForDeletes; emit `reports.[reportId].delete`.
- **Done when:** Report removed, delete message enqueued, event emitted; insights cleaned asynchronously.
- **Priority:** Must · Confidence: 🟢

### T-ADV-06 — List insights with cursor pagination
- **Behavior source:** `Modules/Advisor/Http/Insights/XList.php`
- **Description:** Verify parent report exists in project scope; paginate insights by reportInternalId.
- **Done when:** Insights paged correctly, scoped to report and project.
- **Priority:** Must · Confidence: 🟢

### T-ADV-07 — Get single insight
- **Behavior source:** `Modules/Advisor/Http/Insights/Get.php`
- **Description:** Fetch insight by insightId + reportInternalId + projectInternalId.
- **Done when:** Returns insight; cross-report/cross-project access → 404.
- **Priority:** Must · Confidence: 🟢

### T-ADV-08 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Advisor/Http/**`
- **Description:** OpenAPI 3.1 definitions for 5 endpoints; SDK labels `list/get/delete`.
- **Done when:** `openapi/advisor.yaml` matches action `getName()`.
- **Priority:** Should · Confidence: 🟡

### T-ADV-09 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Advisor/**`
- **Description:** Cover list/get/delete report, list/get insight, cross-project isolation, invalid cursor, async deletion.
- **Done when:** E2E suite passes with ADMIN/KEY auth.
- **Priority:** Must · Confidence: 🟢
