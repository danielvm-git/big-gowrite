# Design: Migrations

> Unit: `Migrations`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Migrations module is an **admin-only HTTP service** that creates migration documents and enqueues async processing to a migration worker via `publisherForMigrations`. Provider reports connect to sources synchronously at request time. CSV/JSON import/export uses storage devices and the internal platform bucket.

```
Admin
  │ AuthType::ADMIN
  ▼
app/http.php → Migrations\Http\** (16 actions)
  │
  ├── dbForProject (migrations collection)
  ├── publisherForMigrations → worker-migrations
  ├── deviceForFiles (import source) / deviceForMigrations (decrypted temp)
  ├── dbForPlatform (internal bucket lookup)
  └── providers: Appwrite API, Firebase, Supabase, NHost (reports)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| GET | `/v1/migrations` | `listMigrations` | migrations.read |
| GET | `/v1/migrations/:migrationId` | `getMigration` | migrations.read |
| PATCH | `/v1/migrations/:migrationId` | `retryMigration` | migrations.write |
| DELETE | `/v1/migrations/:migrationId` | `deleteMigration` | migrations.write |
| POST | `/v1/migrations/appwrite` | `createAppwriteMigration` | migrations.write |
| GET | `/v1/migrations/appwrite/report` | `getAppwriteReport` | migrations.write |
| POST | `/v1/migrations/firebase` | `createFirebaseMigration` | migrations.write |
| GET | `/v1/migrations/firebase/report` | `getFirebaseReport` | migrations.write |
| POST | `/v1/migrations/supabase` | `createSupabaseMigration` | migrations.write |
| GET | `/v1/migrations/supabase/report` | `getSupabaseReport` | migrations.write |
| POST | `/v1/migrations/nhost` | `createNHostMigration` | migrations.write |
| GET | `/v1/migrations/nhost/report` | `getNHostReport` | migrations.write |
| POST | `/v1/migrations/csv/exports` | `createCSVExport` | migrations.write |
| POST | `/v1/migrations/csv/imports` | `createCSVImport` | migrations.write |
| POST | `/v1/migrations/json/exports` | `createJSONExport` | migrations.write |
| POST | `/v1/migrations/json/imports` | `createJSONImport` | migrations.write |

## 3. Data flows — key actions

### 3.1 Create provider migration

```
1. Validate admin auth + provider credentials
2. (Firebase) validate service account JSON fields
3. Create migrations doc (status=pending, credentials)
4. Enqueue MigrationMessage via publisherForMigrations
5. Emit migrations.[migrationId].create
```

### 3.2 Retry migration

```
1. Fetch migration
2. status != failed → throw MIGRATION_IN_PROGRESS
3. Reset to pending, re-enqueue, emit retry event
```

### 3.3 Provider report (dry-run)

```
1. Connect to source provider at request time
2. Enumerate services/resources to migrate
3. Return preview payload
4. Connection failure → MIGRATION_PROVIDER_ERROR
```

### 3.4 CSV/JSON export

```
1. Validate database type (CSV_ALLOWED_DATABASE_TYPES for CSV)
2. Build resources via Transfer::extractServices()
3. Export docs to internal bucket (default), options for CSV
4. notify=true → email user on completion
```

### 3.5 CSV/JSON import

```
1. Read source file (deviceForFiles); decrypt/decompress if needed → deviceForMigrations
2. Validate DB type
3. Import docs into collection
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | Migrations collection |
| `dbForPlatform` | Internal bucket lookup |
| `queueForEvents` | Lifecycle events |
| `publisherForMigrations` | Worker enqueue |
| `project` / `platform` | Context for MigrationMessage |
| `authorization` | RBAC visibility checks |
| `deviceForFiles` | Import source reads |
| `deviceForMigrations` | Decrypted temp copies |
| `getDatabasesDB` | Appwrite report connection |
| `user` | userInternalId for export notify |

## 5. Events emitted

- `migrations.[migrationId].create`
- `migrations.[migrationId].retry`
- `migrations.[migrationId].delete`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `MIGRATION_IN_PROGRESS` | 400 | Retry non-failed migration |
| `MIGRATION_PROVIDER_ERROR` | 400/502 | Provider report connection failure |
| Unauthorized | 401/403 | Non-admin on create endpoints |
