# Design: Databases

> Unit: `Databases`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Databases module implements a **multi-paradigm data plane** over metadata + dedicated backend stores. Four product types share the action class hierarchy via `getDatabaseType()` overrides and per-type SDK namespaces. Requests write metadata to `$dbForProject`; async schema operations are enqueued as `DatabaseMessage` and executed by the `databases` worker against the actual backend (DSN resolved via `Pool`). Transactions stage operations in `transactionLogs` and commit atomically with optimistic locking.

```
Client
  │ scopes: databases.*, collections.*, documents.*
  ▼
app/http.php → Databases\Http\** (per type)
  │
  ├── $dbForProject (metadata: databases, attributes, indexes, transactions, transactionLogs)
  ├── $getDatabasesDB → Pool::dsn(type, region) → backend (MongoDB/PostgreSQL/...)
  ├── publisherForDatabase → worker-databases (DatabaseMessage)
  ├── publisherForDeletes (transaction cleanup)
  ├── queueForEvents / queueForRealtime / queueForWebhooks
  └── embeddingAgent (embeddings endpoint)
```

## 2. Route-to-action map (legacy surface)

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/databases` | `createDatabase` | databases.write |
| GET | `/v1/databases` | `listDatabases` | databases.read |
| GET | `/v1/databases/:databaseId` | `getDatabase` | databases.read |
| PUT | `/v1/databases/:databaseId` | `updateDatabase` | databases.write |
| DELETE | `/v1/databases/:databaseId` | `deleteDatabase` | databases.write |
| POST | `/v1/databases/:databaseId/collections` | `createCollection` | collections.write |
| GET | `/v1/databases/:databaseId/collections` | `listCollections` | collections.read |
| GET | `/v1/databases/:databaseId/collections/:collectionId` | `getCollection` | collections.read |
| PATCH | `/v1/databases/:databaseId/collections/:collectionId` | `updateCollection` | collections.write |
| DELETE | `/v1/databases/:databaseId/collections/:collectionId` | `deleteCollection` | collections.write |
| POST | `.../collections/:collectionId/documents` | `createDocument` | documents.write |
| GET | `.../documents` | `listDocuments` | documents.read |
| GET | `.../documents/:documentId` | `getDocument` | documents.read |
| PATCH | `.../documents/:documentId` | `updateDocument` | documents.write |
| PUT | `.../documents/:documentId` | `upsertDocument` | documents.write |
| DELETE | `.../documents/:documentId` | `deleteDocument` | documents.write |
| PATCH | `.../documents` | `updateDocuments` | documents.write |
| PUT | `.../documents` | `upsertDocuments` | documents.write |
| DELETE | `.../documents` | `deleteDocuments` | documents.write |
| PATCH | `.../documents/:documentId/:attribute/increment` | `incrementDocumentAttribute` | documents.write |
| PATCH | `.../documents/:documentId/:attribute/decrement` | `decrementDocumentAttribute` | documents.write |
| POST | `.../attributes/{16 types}` | `create{Type}Attribute` | collections.write |
| PATCH | `.../attributes/{type}/:attributeId` | `update{Type}Attribute` | collections.write |
| GET | `.../attributes` | `listAttributes` | collections.read |
| GET | `.../attributes/:attributeId` | `getAttribute` | collections.read |
| DELETE | `.../attributes/:attributeId` | `deleteAttribute` | collections.write |
| POST | `.../indexes` | `createIndex` | collections.write |
| GET | `.../indexes` | `listIndexes` | collections.read |
| GET | `.../indexes/:indexId` | `getIndex` | collections.read |
| DELETE | `.../indexes/:indexId` | `deleteIndex` | collections.write |
| POST | `/v1/databases/transactions` | `createDatabasesTransaction` | documents.write |
| GET | `/v1/databases/transactions` | `listDatabasesTransactions` | documents.read |
| GET | `/v1/databases/transactions/:transactionId` | `getDatabasesTransaction` | documents.read |
| PATCH | `/v1/databases/transactions/:transactionId` | `updateDatabasesTransaction` | documents.write |
| DELETE | `/v1/databases/transactions/:transactionId` | `deleteDatabasesTransaction` | documents.write |
| POST | `/v1/databases/transactions/:transactionId/operations` | `createDatabasesTransactionOperations` | documents.write |

**Mirrored surfaces:** TablesDB (`/v1/tablesdb`, tables/columns/rows), DocumentsDB (`/v1/documentsdb`, no column tree), VectorsDB (`/v1/vectorsdb` + `POST .../documents/query` → `createVectorsDBDocumentsQuery`).

**Embeddings:** `POST /v1/embeddings/text` → `createTextEmbeddings` (documents.write).

## 3. Data flows — key actions

### 3.1 Async attribute creation (all types)

```
1. Validate type + plan limits
2. Create metadata doc (status=processing) in $dbForProject
3. Enqueue DatabaseMessage (DATABASE_TYPE_CREATE_ATTRIBUTE)
4. Return 202 with processing attribute
5. Worker: DDL on backend → status available/failed → realtime event
```

### 3.2 Collection create with DSN pre-resolution

```
1. Validate collection payload (+ inline attributes/indexes)
2. Resolve data-plane DSN via Pool (provisioning)
3. On success: write metadata doc + enqueue ops
4. On failure: abort, no orphaned metadata
```

### 3.3 Transaction flow

```
1. POST /transactions → open (TTL-bounded token)
2. POST /transactions/:id/operations → stage batch into transactionLogs (mock response)
3. PATCH /transactions/:id {commit:true} → for each staged op:
   a. Compare doc $updatedAt vs staged timestamp
   b. Conflict → ConflictException (abort all)
   c. Apply write
4. commit=false/rollback=true → discard staged ops (cleanup via publisherForDeletes)
```

### 3.4 Worker message types

| Message type | Operation |
|--------------|-----------|
| `DATABASE_TYPE_DELETE_DATABASE` | Recursive delete (collections then metadata) |
| `DATABASE_TYPE_DELETE_COLLECTION` | Delete collection + cascade |
| `DATABASE_TYPE_CREATE_ATTRIBUTE` | DDL create column/field |
| `DATABASE_TYPE_DELETE_ATTRIBUTE` | DDL drop + index cascade |
| `DATABASE_TYPE_CREATE_INDEX` | Build index |
| `DATABASE_TYPE_DELETE_INDEX` | Drop index |

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `$dbForProject` | Metadata collections (databases, attributes, indexes, transactions, transactionLogs) |
| `$dbForPlatform` | Project lookup (worker) |
| `$getDatabasesDB` | Data-plane Database factory (DSN via Pool) |
| `$queueForEvents` | Event bus (realtime/webhook/function triggers) |
| `$queueForRealtime` | Realtime pub/sub |
| `$queueForWebhooks` | Webhook delivery |
| `$publisherForFunctions` | Function execution triggers |
| `$publisherForDatabase` | DatabaseMessage enqueue |
| `$publisherForDeletes` | Transaction cleanup |
| `$usage` | Usage metrics |
| `$authorization` | RBAC |
| `$user` | Current user |
| `$transactionState` | Optimistic-lock state during commit |
| `$eventProcessor` | Cached function/webhook event maps |
| `$embeddingAgent` | Text embeddings (Appwrite provider) |
| `$plan` | Plan limits (batchSize, transactionSize, maxEmbeddingTexts, allowEncrypt) |

## 5. Events emitted

- `databases.[databaseId].{create,update,delete}`
- `databases.[databaseId].collections.[collectionId].{create,update,delete}`
- `databases.[databaseId].collections.[collectionId].attributes.[attributeId].{create,update,delete}` (update also fired by worker on async completion)
- `databases.[databaseId].collections.[collectionId].indexes.[indexId].{create,update,delete}` (update fired by worker)
- `databases.[databaseId].collections.[collectionId].documents.[documentId].{create,update,delete}`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `ConflictException` | 409 | Transaction optimistic-lock conflict |
| `GENERAL_..._LIMIT` | 400 | Plan limits exceeded |
| Index limit | 400 | > adapter `getLimitForIndexes()` |
| Attribute not available | 400 | Query on encrypt/processing attribute |

