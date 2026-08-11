# Requirements: Databases

> Identifier: `DBS-DATA-PLANE`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Databases module is Appwrite's primary data persistence layer. It exposes four database product types sharing one class hierarchy via inheritance: **legacy** (`/v1/databases`, collections-documents), **TablesDB** (`/v1/tablesdb`, SQL-oriented tables/rows/columns), **DocumentsDB** (`/v1/documentsdb`, NoSQL documents), and **VectorsDB** (`/v1/vectorsdb`, vector-native with similarity search). It also provides **Embeddings** (`POST /v1/embeddings/text`). A background **Databases worker** processes async schema operations (attribute/index create/delete) against the backing store, updating status flags. A **transaction system** provides optimistic-locking, ACID-like staging across collections with a single atomic commit.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Databases` | Module purpose | 🟢 |
| `specs/Databases/legacy-mapping.md` | Endpoints, worker, business rules, DI | 🟢 |
| `specs/domain.md#Databases` | Four product types, transactions | 🟢 |
| `specs/architecture.md#DataPlane` | DSN pool, dedicated DB | 🟢 |
| `specs/workers.md#Databases` | Async schema worker | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App developer | Store structured data | Create database/collection, CRUD documents |
| Data analyst | SQL-style modeling | TablesDB tables/rows/columns |
| AI developer | Store and search vectors | VectorsDB + similarity query |
| App developer | Multi-step writes | Open transaction, stage ops, commit |
| Schema admin | Evolve schema | Async create attribute/index |

## 4. Business rules

1. **BR-01:** Database types are distinct products — own URL prefix, SDK namespace, terminology (collection/table, document/row, attribute/column); classes reused via `getDatabaseType()` override. 🟢
2. **BR-02:** Attribute/index creation is asynchronous — HTTP creates metadata doc with `status: processing` + enqueues `DatabaseMessage`; worker performs DDL and flips to `available`/`failed`. 🟢
3. **BR-03:** Collection create resolves data-plane DSN before writing metadata — no orphaned metadata on provisioning failure. 🟢
4. **BR-04:** Inline attributes/indexes on `createCollection`; relationship attributes explicitly rejected from inline path. 🟢
5. **BR-05:** Bulk operations are API-key/admin only; session-user path allows single-document ops. 🟢
6. **BR-06:** Transaction staging — write endpoints accept optional `transactionId`; staged ops return mock response; real writes only on `commit=true`. 🟢
7. **BR-07:** Transaction optimistic concurrency — commit compares `$updatedAt` of each doc against staged timestamp; conflict → `ConflictException`. 🟢
8. **BR-08:** Permission aggregation via `Permission::aggregate()` for shorthand roles. 🟢
9. **BR-09:** String attributes may set `encrypt=true`; server stores `encrypt` filter; attribute not queryable; min size `APP_DATABASE_ENCRYPT_SIZE_MIN` + paid plan. 🟢
10. **BR-10:** Relationship attributes use separate creation path — mirrored attribute via worker; `twoWayKey` set to available after creation. 🟢
11. **BR-11:** Index limit per collection — checked against adapter `getLimitForIndexes()` (≤60, count up to 61). 🟢
12. **BR-12:** VectorsDB backend defaults to PostgreSQL; DocumentsDB to MongoDB. 🟡
13. **BR-13:** Attribute deletion cascades to indexes (worker trims referencing indexes). 🟡
14. **BR-14:** Database deletion recursive — worker deletes every collection via `deleteByGroup` then drops metadata collection. 🟡

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Database CRUD (all 4 types) | Must | create/get/update/delete/list per type | 🟢 |
| RF-02 | Collection/table CRUD + inline attrs/indexes | Must | Create with inline arrays; relationship rejected inline | 🟢 |
| RF-03 | Document/row CRUD (single + bulk) | Must | Single via session; bulk API-key/admin only | 🟢 |
| RF-04 | Typed attributes (16 types) async create/update | Must | Metadata processing → worker DDL → available/failed | 🟢 |
| RF-05 | Index CRUD async | Must | createIndex returns 202; worker builds; limit enforced | 🟢 |
| RF-06 | Transactions (open/list/get/commit/rollback/delete + stage ops) | Must | Optimistic-locking atomic commit | 🟢 |
| RF-07 | VectorsDB similarity query | Must | `createVectorsDBDocumentsQuery` with optional TTL cache | 🟢 |
| RF-08 | Embeddings generation | Must | `POST /v1/embeddings/text` via embedding agent | 🟢 |
| RF-09 | Encrypted string attributes | Should | encrypt=true, non-queryable, size/plan gates | 🟢 |
| RF-10 | Relationship attributes | Must | Mirrored attribute via worker | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Consistency | ACID-like transactions via optimistic locking | `$updatedAt` comparison, `ConflictException` | 🟢 |
| Performance | Async schema ops off the request path | DatabaseMessage → worker | 🟢 |
| Performance | TTL cache on vector queries | VectorsDB query options | 🟢 |
| Security | Encrypted attributes non-queryable | `encrypt` filter | 🟢 |
| Multi-tenancy | Per-project metadata + data plane | `$dbForProject`, `$getDatabasesDB` | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Async attribute creation
  Given a collection
  When POST .../attributes/string is called
  Then a processing attribute doc is created and a DatabaseMessage is enqueued
  And the worker eventually sets status to available

Scenario: Transaction commit with conflict
  Given a transaction with staged ops
  When commit=true is sent and a doc changed since staging
  Then a ConflictException is raised and no writes are applied

Scenario: Bulk update via session user
  Given a session-authenticated user
  When PATCH .../documents (bulk) is called
  Then the request is rejected (API-key/admin only)

Scenario: Vector similarity search
  Given documents with embeddings in a vectorsdb collection
  When POST .../documents/query is called
  Then ranked results by similarity are returned
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01/02/03 (core CRUD) | Must | Primary API |
| RF-04/05 (schema async) | Must | Schema evolution |
| RF-06 (transactions) | Should | Advanced consistency |
| RF-07 (vector query) | Should | VectorsDB differentiator |
| RF-08 (embeddings) | Should | AI features |
| RF-09 (encryption) | Should | Paid plan feature |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] VectorsDB PostgreSQL / DocumentsDB MongoDB defaults from `Pool::dsn()`.
- 🟡 [INFERRED] Attribute-deletion index cascade details.
- 🟡 [INFERRED] Recursive database-deletion flow details.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |

