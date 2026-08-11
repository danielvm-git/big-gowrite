# Tasks: Databases

> Unit: `Databases`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-DBS-01 — Module bootstrap + registries
- **Behavior source:** `Modules/Databases/Module.php`, `Constants.php`, `Pool.php`, `Services/Registry/*`
- **Description:** Register Http + Workers services; define type constants and resource aliases; implement Pool DSN resolution per type/region; per-type registries (Legacy, TablesDB, DocumentsDB, VectorsDB, Embeddings).
- **Done when:** All four product surfaces registered; DSN pool resolves correctly.
- **Priority:** Must · Confidence: 🟢

### T-DBS-02 — Database CRUD (per type)
- **Behavior source:** `Modules/Databases/Http/Databases/{Create,Get,Update,Delete,XList}.php`
- **Description:** Implement create/get/update/delete/list across all 4 types via shared inheritance + `getDatabaseType()`.
- **Done when:** CRUD verified per type.
- **Priority:** Must · Confidence: 🟢

### T-DBS-03 — Collection/table CRUD + inline schema
- **Behavior source:** `Modules/Databases/Http/Databases/Collections/**`
- **Description:** CRUD collections; support inline `attributes[]`/`indexes[]`; reject relationship attributes inline; pre-resolve DSN before metadata write.
- **Done when:** Inline schema handled; no orphan metadata on provisioning failure.
- **Priority:** Must · Confidence: 🟢

### T-DBS-04 — Document/row CRUD single + bulk
- **Behavior source:** `Modules/Databases/Http/Databases/Collections/Documents/**`
- **Description:** Single CRUD + upsert + increment/decrement; bulk update/upsert/delete restricted to API-key/admin.
- **Done when:** Bulk gating enforced; increments atomic.
- **Priority:** Must · Confidence: 🟢

### T-DBS-05 — Typed attributes (16 types) async create/update
- **Behavior source:** `Modules/Databases/Http/Databases/Collections/Attributes/{Type}/{Create,Update}.php`
- **Description:** Implement all typed attribute create/update paths; enqueue DatabaseMessage; support encrypt for strings with plan gate.
- **Done when:** Attributes create as processing → worker → available.
- **Priority:** Must · Confidence: 🟢

### T-DBS-06 — Relationships
- **Behavior source:** `Modules/Databases/Http/Databases/Collections/Attributes/Relationship/**`
- **Description:** Two-way relationship creation with mirrored attribute via worker; `twoWayKey` → available after creation.
- **Done when:** Mirrored attribute created; status correct.
- **Priority:** Must · Confidence: 🟢

### T-DBS-07 — Indexes async CRUD
- **Behavior source:** `Modules/Databases/Http/Databases/Collections/Indexes/**`
- **Description:** Create (202)/list/get/delete with index limit check vs adapter `getLimitForIndexes()`.
- **Done when:** Limit enforced; async build completes.
- **Priority:** Must · Confidence: 🟢

### T-DBS-08 — Transactions
- **Behavior source:** `Modules/Databases/Http/Databases/Transactions/**`
- **Description:** Open/list/get/commit/rollback/delete + operations staging; optimistic-lock commit via `$transactionState`; cleanup via publisherForDeletes.
- **Done when:** Atomic commit; conflict → 409; rollback discards.
- **Priority:** Should · Confidence: 🟢

### T-DBS-09 — VectorsDB similarity query
- **Behavior source:** `Modules/Databases/Http/VectorsDB/**`
- **Description:** Vector search endpoint with optional TTL cache.
- **Done when:** Ranked similarity results returned.
- **Priority:** Should · Confidence: 🟢

### T-DBS-10 — Embeddings endpoint
- **Behavior source:** `Modules/Databases/Http/Embeddings/Text/Create.php`
- **Description:** Generate text embeddings via `$embeddingAgent`; enforce `databasesMaxEmbeddingTexts`.
- **Done when:** Embeddings returned for batch.
- **Priority:** Should · Confidence: 🟢

### T-DBS-11 — Databases worker
- **Behavior source:** `Modules/Databases/Workers/Databases.php`
- **Description:** Process 6 DatabaseMessage types; update attribute/index status; purge caches; send realtime events; cascade attribute→index and recursive database delete.
- **Done when:** All message types handled; statuses correct.
- **Priority:** Must · Confidence: 🟢

### T-DBS-12 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Databases/Http/**`
- **Description:** OpenAPI 3.1 for all four products + embeddings (4 namespaces).
- **Done when:** `openapi/databases*.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-DBS-13 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Databases/**`
- **Description:** Cover CRUD per type, async schema, transactions, bulk gating, vector search, embeddings.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
