# Requirements: Projects

> Identifier: `PRJS-PLATFORM`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Projects module manages Appwrite projects on the platform layer: creation (database provisioning, collection scaffolding, audit schema), metadata updates, team ownership transfers, developer keys (deprecated/tombstoned), CRON schedules linked to resources, and onboarding stage tracking. Projects are top-level entities stored in `dbForPlatform`.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Projects` | Module purpose | 🟢 |
| `specs/Projects/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/architecture.md#MultiTenant` | Provisioning modes | 🟢 |
| `specs/domain.md#Project` | Project document model | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Platform admin | Create projects | Provision with region + DSN selection |
| Team admin | Manage ownership | Transfer project to another team |
| Admin | Schedule jobs | Create CRON schedule for function |
| Console user | Onboarding | Track/skip onboarding stages |

## 4. Business rules

1. **BR-01:** `projectId = console` reserved → `PROJECT_RESERVED_PROJECT`. 🟢
2. **BR-02:** Owning team must exist → `TEAM_NOT_FOUND`. 🟢
3. **BR-03:** Region must be in `_APP_PROJECT_REGIONS` when list non-empty. 🟢
4. **BR-04:** DSN selected at creation: region-filtered, randomized, overridable via `_APP_DATABASE_OVERRIDE`; stored on project doc. 🟢
5. **BR-05:** Shared-table projects share `appwrite` schema; dedicated projects get own namespace (`_<internalId>`) with collections + audit schema initialized. 🟢
6. **BR-06:** `createDevKey` always throws `DEV_KEY_GONE` (tombstoned for backward-compat messaging). 🟢
7. **BR-07:** Dev key ops verify `projectInternalId` match. 🟢
8. **BR-08:** Schedule creation validates referenced resource exists in project DB before platform DB write. 🟢
9. **BR-09:** `listSchedules` filters by `projectId` + `projectInternalId` (allows null for legacy). 🟢
10. **BR-10:** Onboarding stages from `onboarding` config; completed stage cannot be overwritten as skipped. 🟢
11. **BR-11:** `updateStage` resolves actor type (key-account, key-organization, key-project, admin, user, guest). 🟢
12. **BR-12:** Team transfer propagates permissions to installations, repositories, VCS comments. 🟢
13. **BR-13:** Project status defaults to `PROJECT_STATUS_ACTIVE`; no endpoint changes it here. 🟡
14. **BR-14:** Post-creation hook `afterProjectCreation` for migration mirroring. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Create project (provision DB) | Must | Region/DSN selection; shared vs dedicated; reserved id rejected | 🟢 |
| RF-02 | Update project metadata | Must | Sparse metadata update | 🟢 |
| RF-03 | List projects (team filter) | Must | Pagination + optional team scope | 🟢 |
| RF-04 | Transfer project team | Must | Ownership + permission propagation | 🟢 |
| RF-05 | Dev keys (tombstoned create + CRUD) | Must | Create → DEV_KEY_GONE; others scoped | 🟢 |
| RF-06 | Schedules CRUD | Must | Resource validation; dual-id filtering | 🟢 |
| RF-07 | Onboarding stages list/update | Must | Config-driven; skip guard | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Multi-tenancy | Cross-project dev-key isolation | projectInternalId checks | 🟢 |
| Compatibility | Legacy schedule records | dual-id filter (null allowed) | 🟢 |
| Extensibility | Post-creation hook | afterProjectCreation | 🟢 |
| Backward-compat | Tombstoned dev-key create | DEV_KEY_GONE | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Reserved project id
  Given projectId = "console"
  When POST /v1/projects is called
  Then PROJECT_RESERVED_PROJECT is returned

Scenario: Dev key creation tombstoned
  Given any valid project
  When POST /v1/projects/:projectId/dev-keys is called
  Then DEV_KEY_GONE error is returned

Scenario: Schedule resource validation
  Given a schedule referencing a non-existent function
  When POST /v1/projects/:projectId/schedules is called
  Then the request is rejected

Scenario: Completed stage cannot be skipped
  Given a completed onboarding stage
  When PATCH /v1/projects/:projectId/stages/:stageId is called
  Then the skip is rejected
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..04 (project lifecycle) | Must | Platform core |
| RF-05 (dev keys) | Should | Deprecated surface |
| RF-06 (schedules) | Must | CRON jobs |
| RF-07 (stages) | Should | Onboarding UX |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Project status default/enum semantics.
- 🔴 [GAP] Team-transfer permission propagation exact scope.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
