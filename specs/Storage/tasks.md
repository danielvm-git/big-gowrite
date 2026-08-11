# Tasks: Storage

> Unit: `Storage`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-STR-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Storage/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 14 endpoints + scopes.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-STR-02 — Bucket CRUD
- **Behavior source:** `Http/Buckets/{Create,Get,XList,Update,Delete}.php`
- **Description:** Create → bucket_{seq} collection; update full replace; delete soft + enqueue physical deletion.
- **Done when:** Bucket lifecycle verified.
- **Priority:** Must · Confidence: 🟢

### T-STR-03 — Chunked file upload
- **Behavior source:** `Http/Buckets/Files/Create.php`
- **Description:** Content-Range chunking; distributed lock; compression/encryption; antivirus; move to final device.
- **Done when:** Chunk assembly + finalization verified.
- **Priority:** Must · Confidence: 🟢

### T-STR-04 — File metadata CRUD
- **Behavior source:** `Http/Buckets/Files/{Get,XList,Update,Delete}.php`
- **Description:** get/list (searchable)/update name+permissions; delete physical + DB + cache invalidation.
- **Done when:** File CRUD verified.
- **Priority:** Must · Confidence: 🟢

### T-STR-05 — Download + View (Range)
- **Behavior source:** `Http/Buckets/Files/{Download,View}/Get.php`
- **Description:** Attachment vs inline; Range support; permission mode (fileSecurity vs bucket); Cache-Control.
- **Done when:** Range requests + permission modes verified.
- **Priority:** Must · Confidence: 🟢

### T-STR-06 — Preview (image transforms)
- **Behavior source:** `Http/Buckets/Files/Preview/Get.php`
- **Description:** Imagick transforms; per-bucket enable gate.
- **Done when:** Transform pipeline verified.
- **Priority:** Should · Confidence: 🟢

### T-STR-07 — Push delivery (JWT)
- **Behavior source:** `Http/Buckets/Files/Push/Get.php`
- **Description:** HS256 JWT 1h; internal flag selects DB.
- **Done when:** JWT-gated delivery verified.
- **Priority:** Should · Confidence: 🟢

### T-STR-08 — Cache-Control config objects
- **Behavior source:** `Config/CacheControl.php`, `Config/StorageCacheControl.php`
- **Description:** Compute Cache-Control from bucket/file/token state.
- **Done when:** Headers correct across modes.
- **Priority:** Should · Confidence: 🟢

### T-STR-09 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Storage/Http/**`
- **Description:** OpenAPI 3.1 for storage endpoints.
- **Done when:** `openapi/storage.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-STR-10 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Storage/**`
- **Description:** Cover buckets, chunked uploads, download/view/preview/push, Range, encryption.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
