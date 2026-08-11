# Contracts: Databases

> Unit: `Databases`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract — legacy surface (representative)

Base: `/v1/databases`. Scopes: `databases.*`, `collections.*`, `documents.*`.

| Method | Path | Key params | Success | Notes |
|--------|------|-----------|---------|-------|
| POST | `/databases` | `databaseId`, `name`, `enabled` | 201 database | |
| GET | `/databases` | `queries`, `search`, `cursor` | 200 databaseList | |
| GET | `/databases/:databaseId` | — | 200 database | |
| PUT | `/databases/:databaseId` | `name`, `enabled` | 200 database | |
| DELETE | `/databases/:databaseId` | — | 204 | async |
| POST | `/databases/:databaseId/collections` | `collectionId`, `name`, `permissions`, `attributes[]`, `indexes[]` | 201 collection | |
| POST | `.../documents` | `documentId`, `data`, `permissions`, `transactionId` | 201 document | bulk via `documents[]` (admin/key) |
| PATCH | `.../documents/:documentId` | `data`, `permissions`, `transactionId` | 200 document | |
| PUT | `.../documents/:documentId` | `data`, `permissions` | 200 document | upsert |
| PATCH | `.../documents` | `queries`, `data` | 200 | bulk update (admin/key) |
| PATCH | `.../documents/:documentId/:attribute/increment` | `value` | 200 document | also decrement |
| POST | `.../attributes/string` | `key`, `size`, `required`, `default`, `array`, `encrypt` | 202 attribute | status processing; 16 types |
| POST | `.../indexes` | `key`, `type`, `attributes[]`, `orders[]` | 202 index | async |
| POST | `/databases/transactions` | — | 201 transaction token | TTL-bounded |
| PATCH | `/databases/transactions/:transactionId` | `commit`/`rollback` | 200 | 409 on conflict |
| POST | `/databases/transactions/:transactionId/operations` | `operations[]` | 201 mock | staged |

**Mirrored:** TablesDB (`/v1/tablesdb`, tables/columns/rows), DocumentsDB (`/v1/documentsdb`), VectorsDB (`/v1/vectorsdb` + `POST .../documents/query`).

## Embeddings contract

| Method | Path | Params | Success |
|--------|------|--------|---------|
| POST | `/v1/embeddings/text` | `text[]` | 200 `{ embeddings: [[float]] }` |

## Queue contract (Databases worker)

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForDatabase` | Http → worker-databases | `DatabaseMessage`: type (DELETE_DATABASE, DELETE_COLLECTION, CREATE_ATTRIBUTE, DELETE_ATTRIBUTE, CREATE_INDEX, DELETE_INDEX), databaseId, collectionId, attributeId/indexId, projectInternalId |
| `publisherForDeletes` | Transactions → worker-deletes | Cleanup of committed/rolled-back transaction docs |

## Event contract

- `databases.[databaseId].{create,update,delete}`
- `databases.[databaseId].collections.[collectionId].{create,update,delete}`
- `...attributes.[attributeId].{create,update,delete}` (update fired on worker completion)
- `...indexes.[indexId].{create,update,delete}` (update fired on worker completion)
- `...documents.[documentId].{create,update,delete}`

## Data-plane contract (backends)

| Type | Default backend | Notes |
|------|-----------------|-------|
| legacy | MongoDB (config) | collections/documents |
| tablesdb | PostgreSQL (config) | tables/rows/columns |
| documentsdb | MongoDB (default) | documents |
| vectorsdb | PostgreSQL (default) | vector search |

DSN resolved by `Pool::dsn(type, region)` from env-var pool config.

## Error model

Standard Appwrite error JSON; `409 ConflictException` for transaction optimistic-lock conflicts.
