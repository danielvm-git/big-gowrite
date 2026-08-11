# Contracts: Migrations

> Unit: `Migrations`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/migrations`. Creation/report endpoints require `AuthType::ADMIN`.

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| GET | `/migrations` | `search`, `cursor`, `limit` | 200 migrationList | 401 |
| GET | `/migrations/:migrationId` | — | 200 migration | 404 |
| PATCH | `/migrations/:migrationId` | — | 200 migration | MIGRATION_IN_PROGRESS |
| DELETE | `/migrations/:migrationId` | — | 204 | 404 |
| POST | `/migrations/appwrite` | `resources`, `onDuplicate` (fail/skip/overwrite), source credentials | 202 migration | 400 |
| GET | `/migrations/appwrite/report` | source params | 200 report | MIGRATION_PROVIDER_ERROR |
| POST | `/migrations/firebase` | `serviceAccount` JSON, `resources` | 202 | 400 invalid SA |
| GET | `/migrations/firebase/report` | `serviceAccount` JSON | 200 report | MIGRATION_PROVIDER_ERROR |
| POST | `/migrations/supabase` | `credentials` {host,username,password}, `resources` | 202 | 400 |
| GET | `/migrations/supabase/report` | `credentials` | 200 report | MIGRATION_PROVIDER_ERROR |
| POST | `/migrations/nhost` | `credentials`, `resources` | 202 | 400 |
| GET | `/migrations/nhost/report` | `credentials` | 200 report | MIGRATION_PROVIDER_ERROR |
| POST | `/migrations/csv/exports` | `databaseId`, `collectionId`, `delimiter`, `enclosure`, `escape`, `header`, `notify` | 202 | 400 type |
| POST | `/migrations/csv/imports` | file, `databaseId`, `collectionId` | 202 | 400 type |
| POST | `/migrations/json/exports` | `databaseId`, `collectionId`, `notify` | 202 | 400 |
| POST | `/migrations/json/imports` | file, `databaseId`, `collectionId` | 202 | 400 |

## Migration status model

`pending` → `processing` → `completed` | `failed` (retry from failed → pending).

## Queue contract

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForMigrations` | Http → worker-migrations | `MigrationMessage`: migrationId, provider type, resources, credentials, project/platform context |

## Event contract

- `migrations.[migrationId].create`
- `migrations.[migrationId].retry`
- `migrations.[migrationId].delete`

## Storage contract

| Operation | Device/Bucket |
|-----------|---------------|
| CSV/JSON export | Internal platform bucket (`default`) |
| CSV/JSON import | `deviceForFiles` (source) → decrypt/decompress → `deviceForMigrations` |

## Configuration / constants

| Constant | Effect |
|----------|--------|
| `CSV_ALLOWED_DATABASE_TYPES` | DB types allowed for CSV ops |
| `MIGRATION_PROVIDER_ERROR` | Provider connection failure code |
| `MIGRATION_IN_PROGRESS` | Retry on non-failed status |
