# Requirements: Organization

> Identifier: `ORG-TEAM-PROJECTS`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Organization module provides team-scoped project management. It exposes CRUD for projects under `/v1/organization/projects`, always bound to the calling team (`$team`). Creation provisions a new database namespace (or joins a shared table), sets up audit schema, bootstraps collection schemas, and seeds auth methods. Deletion enqueues async cleanup. Supports `ORGANIZATION` auth type alongside `ADMIN` and `KEY`.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Organization` | Module purpose | 🟢 |
| `specs/Organization/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/architecture.md#MultiTenant` | DB provisioning modes | 🟢 |
| `specs/domain.md#Project` | Project document model | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Team admin | Create project | Provision new project in region |
| Team admin | Organize projects | List/get/rename/delete team projects |
| Platform ops | Database placement | Region selection, shared tables, overrides |

## 4. Business rules

1. **BR-01:** All endpoints gated by Init hook rejecting empty `$team`. 🟢
2. **BR-02:** `projectId = 'console'` reserved → `PROJECT_RESERVED_PROJECT`. 🟢
3. **BR-03:** Region selection: `_APP_PROJECT_REGIONS` allow-list; non-default region filters DSN candidates from `_APP_DATABASE_KEYS` containing the region string. 🟢
4. **BR-04:** Random DSN from region-filtered pool unless `_APP_DATABASE_OVERRIDE` set. 🟢
5. **BR-05:** Shared table support: DSN in `_APP_DATABASE_SHARED_TABLES` → join shared `appwrite` schema under `_APP_DATABASE_SHARED_NAMESPACE`; no per-project DB. 🟢
6. **BR-06:** Project-table mode: new database namespace (`_` + project sequence), audit schema, all `projects` collections bootstrapped from config. 🟢
7. **BR-07:** All auth methods seeded from `auth` config param. 🟢
8. **BR-08:** Get/Update/Delete verify `project.teamInternalId === team.getSequence()`; mismatch → `PROJECT_NOT_FOUND`. 🟢
9. **BR-09:** Deletion enqueues `DELETE_TYPE_DOCUMENT` via `publisherForDeletes`. 🟢
10. **BR-10:** `hooks->trigger('afterProjectCreation')` fired post-creation (migration mirroring). 🟢
11. **BR-11:** XList supports `select` queries; `skipFilters` for unselected attributes. 🟡
12. **BR-12:** Default auth config: session duration, password history, MFA, session alerts. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Create team-scoped project | Must | Provision DB (shared or per-project), audit, collections, auth; reserved id rejected | 🟢 |
| RF-02 | List team projects | Must | Team scoping + select support | 🟢 |
| RF-03 | Get team project | Must | Ownership verified | 🟢 |
| RF-04 | Update project name | Must | Ownership verified | 🟢 |
| RF-05 | Delete project + async cleanup | Must | Enqueue DELETE_TYPE_DOCUMENT | 🟢 |
| RF-06 | Init hook team validation | Must | Empty team blocked | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Multi-tenancy | Team scoping on all queries | `teamInternalId` checks | 🟢 |
| Performance | Region-filtered DSN selection | `_APP_PROJECT_REGIONS` + pool filter | 🟢 |
| Extensibility | Post-creation hooks | `afterProjectCreation` trigger | 🟢 |
| Compatibility | Shared tables / override modes | `_APP_DATABASE_SHARED_TABLES`, `_APP_DATABASE_OVERRIDE` | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Reserved project id
  Given projectId = "console"
  When POST /v1/organization/projects is called
  Then PROJECT_RESERVED_PROJECT is returned

Scenario: Cross-team access
  Given a project owned by another team
  When GET /v1/organization/projects/:projectId is called
  Then PROJECT_NOT_FOUND is returned

Scenario: Shared table join
  Given a DSN in _APP_DATABASE_SHARED_TABLES
  When a project is created
  Then it joins the shared schema and no per-project database is created
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01 (Create) | Must | Provisioning core |
| RF-02..05 (CRUD) | Must | Project lifecycle |
| RF-06 (Init) | Must | Team gating |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] XList skipFilters optimization details.
- 🔴 [GAP] Exact collection schema bootstrap list not enumerated.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
