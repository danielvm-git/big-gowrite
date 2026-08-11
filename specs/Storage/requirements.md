# Requirements: Storage

> Identifier: `STR-OBJECT`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Storage module provides object-storage for Appwrite projects: **Buckets** (containers with security, compression, encryption, antivirus, image-transformation policies) and **Files** (binary objects). Files are served as raw download (attachment), inline view, or server-side transformed image preview. A push-notification endpoint retrieves files via a short-lived JWT instead of standard permission checks.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Storage` | Module purpose | 🟢 |
| `specs/Storage/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/security.md#Storage` | Encryption, antivirus, JWT | 🟢 |
| `specs/workers.md#Deletes` | Physical deletion worker | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App developer | Store user files | Create bucket, chunked upload |
| App developer | Serve media | Download/view with Range |
| App developer | Transform images | Preview endpoint (crop/resize) |
| Push service | Deliver notification image | JWT push endpoint |

## 4. Business rules

1. **BR-01:** Creating a bucket creates a dedicated `bucket_{sequence}` collection with files schema from `app/config/collections/buckets.php`. 🟢
2. **BR-02:** `enabled=false` blocks end-user access; API keys/admin can still read/write. 🟢
3. **BR-03:** Uploads support chunking via `Content-Range`; tracked with `chunksTotal`/`chunksUploaded`. 🟢
4. **BR-04:** Concurrency protected by distributed lock `storage:file:{projectId}:{bucketId}:{fileId}`. 🟢
5. **BR-05:** Compression (GZIP/ZSTD) + AES-128-GCM encryption per-file at finalization, only when `sizeOriginal <= APP_STORAGE_READ_BUFFER`. 🟢
6. **BR-06:** ClamAV scan at finalization for local-device files under `APP_LIMIT_ANTIVIRUS` when enabled. 🟢
7. **BR-07:** `fileSecurity=true` → per-file permissions; `false` → bucket-level permissions. 🟢
8. **BR-08:** Permission aggregation converts role groups (`any`) to permission string sets. 🟢
9. **BR-09:** Image preview requires `imagick`; transformations can be disabled per bucket. 🟢
10. **BR-10:** `/push` authenticates via signed HS256 JWT (1h TTL); `internal` flag selects platform DB over project DB. 🟢
11. **BR-11:** Download/view support HTTP Range requests. 🟢
12. **BR-12:** File delete: physical device removal (or abort partial upload) → DB delete → cache-invalidation enqueue. 🟢
13. **BR-13:** `totalSize` always `0` on bucket responses (aggregation out-of-band). 🟡
14. **BR-14:** `resourceType` label `RESOURCE_TYPE_BUCKETS` = `'buckets'`. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Bucket CRUD | Must | Create → collection; update full replace; delete soft + enqueue | 🟢 |
| RF-02 | File create (chunked upload) | Must | Content-Range chunks; lock; compress/encrypt at finalize; AV scan | 🟢 |
| RF-03 | File metadata get/list/update | Must | Searchable list; name/permissions update | 🟢 |
| RF-04 | File delete | Must | Physical + DB + cache invalidation | 🟢 |
| RF-05 | Download (attachment, Range) | Must | Content-Disposition attachment; partial content | 🟢 |
| RF-06 | View (inline, Range) | Must | Inline streaming; partial content | 🟢 |
| RF-07 | Preview (image transforms) | Must | Imagick crop/resize/rotate; bucket gate | 🟢 |
| RF-08 | Push delivery via JWT | Must | HS256 1h token; internal DB option | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Per-file encryption AES-128-GCM | Compression/encryption pipeline | 🟢 |
| Security | Antivirus scan | ClamAV at finalization | 🟢 |
| Performance | Range request support | Download/view partial content | 🟢 |
| Reliability | Locked concurrent uploads | distributed lock key | 🟢 |
| Performance | Cache-Control tuning | StorageCacheControl value object | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Disabled bucket blocks end users
  Given a bucket with enabled=false
  When an end user attempts to read a file
  Then access is denied

Scenario: Chunked upload
  Given a file uploaded in multiple chunks
  When the final Content-Range arrives
  Then the file is finalized with compression/encryption applied

Scenario: JWT push delivery
  Given a valid signed JWT with internal flag
  When GET .../push is called
  Then the file is delivered without standard permission checks

Scenario: Preview without imagick
  Given a server without the imagick extension
  When GET .../preview is called
  Then a server error is returned
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..06 (core) | Must | Object storage |
| RF-07 (preview) | Should | Image transforms |
| RF-08 (push) | Should | Internal delivery |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] totalSize aggregation is out-of-band.
- 🔴 [GAP] Exact JWT push claims beyond `internal` flag.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
