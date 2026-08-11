# Requirements: Migrations

> Identifier: `MIG-IMPORT-EXPORT`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Migrations module lets admin users import project data from external BaaS providers (Appwrite, Firebase, Supabase, NHost) and move data between formats (CSV, JSON) or between Appwrite instances. Each attempt persists as a `migrations` document transitioning `pending` → `processing` → `completed`/`failed`. Provider-specific "report" endpoints preview the migration before committing. Failed migrations can be retried. The module also supports collection-level CSV/JSON export and import.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Migrations` | Module purpose | 🟢 |
| `specs/Migrations/legacy-mapping.md` | 16 endpoints, business rules, DI | 🟢 |
| `specs/workers.md#Migrations` | Migration worker | 🟢 |
| `specs/domain.md#Migrations` | Migration status model | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Admin | Move from Firebase | Create Firebase migration, preview report, commit |
| Admin | Migrate between Appwrite instances | Appwrite-to-Appwrite with onDuplicate policy |
| Admin | Export data for backup | CSV/JSON export to internal bucket |
| Admin | Bulk load data | CSV/JSON import into a collection |

## 4. Business rules

1. **BR-01:** All migration creation endpoints admin-only (`AuthType::ADMIN`). 🟢
2. **BR-02:** Migration docs created in `pending`, then `publisherForMigrations` enqueues processing. 🟢
3. **BR-03:** Retry only works when `status === 'failed'`; otherwise `MIGRATION_IN_PROGRESS`. 🟢
4. **BR-04:** Firebase migration validates service account JSON (`project_id`, `client_email`, `private_key`). 🟢
5. **BR-05:** Supabase/NHost store DB credentials (host, username, password) in `credentials` field. 🟢
6. **BR-06:** Appwrite-to-Appwrite supports `onDuplicate`: `fail` (default), `skip`, `overwrite`. 🟢
7. **BR-07:** Report endpoints connect to source at request time; failures → `MIGRATION_PROVIDER_ERROR`. 🟢
8. **BR-08:** CSV/JSON export writes to internal platform bucket (`default`); `notify` flag emails user on completion. 🟢
9. **BR-09:** CSV/JSON import decrypts/decompresses source file before copying to `deviceForMigrations`. 🟢
10. **BR-10:** CSV export supports `delimiter`, `enclosure`, `escape`, `header` options. 🟢
11. **BR-11:** DB type validation restricts CSV ops to `CSV_ALLOWED_DATABASE_TYPES`; JSON supports all types. 🟢
12. **BR-12:** Import/export builds `resources` via `Transfer::extractServices()` per database type group. 🟢
13. **BR-13:** CSV import registers HTTP alias at `/v1/migrations/csv` for backward compat. 🟡
14. **BR-14:** No `resourceType` label declared (cross-cutting, no RESOURCE_TYPE_ constant). 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Migration CRUD (get/list/retry/delete) | Must | Retry only on failed; delete doc | 🟢 |
| RF-02 | Appwrite migration + report | Must | onDuplicate policy; dry-run report | 🟢 |
| RF-03 | Firebase migration + report | Must | Service account validation | 🟢 |
| RF-04 | Supabase migration + report | Must | Credentials stored | 🟢 |
| RF-05 | NHost migration + report | Must | Credentials stored | 🟢 |
| RF-06 | CSV export | Must | Options; internal bucket; notify | 🟢 |
| RF-07 | CSV import | Must | Type gating; decrypt/decompress | 🟢 |
| RF-08 | JSON export/import | Must | All DB types | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Admin-only creation endpoints | AuthType::ADMIN | 🟢 |
| Security | Credentials stored in migration doc | `credentials` field | 🟢 |
| Reliability | Report preview before commit | provider report endpoints | 🟢 |
| Compatibility | Backward-compat alias for CSV import | `/v1/migrations/csv` | 🟡 |

## 7. Acceptance criteria

```gherkin
Scenario: Retry a failed migration
  Given a migration with status failed
  When PATCH /v1/migrations/:migrationId is called
  Then the migration is re-enqueued

Scenario: Retry non-failed migration
  Given a migration with status processing
  When PATCH /v1/migrations/:migrationId is called
  Then MIGRATION_IN_PROGRESS error is returned

Scenario: Firebase service account validation
  Given a service account JSON missing private_key
  When POST /v1/migrations/firebase is called
  Then the request is rejected

Scenario: CSV import restricted type
  Given a database type not in CSV_ALLOWED_DATABASE_TYPES
  When POST /v1/migrations/csv/imports is called
  Then the request is rejected
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01 (CRUD) | Must | Lifecycle |
| RF-02..05 (providers) | Must | Core imports |
| RF-06..08 (CSV/JSON) | Must | Data portability |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Backward-compat alias details for CSV import.
- 🔴 [GAP] Migration worker processing internals not visible in module.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
