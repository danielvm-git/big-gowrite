# Contracts: Storage

> Unit: `Storage`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/storage`. Scopes: `buckets.*`, `files.*`.

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `/buckets` | `bucketId`, `name`, `permissions`, `fileSecurity`, `enabled`, `compression`, `encryption`, `antivirus`, `transformations` | 201 bucket | 400 |
| PUT | `/buckets/:bucketId` | same as create | 200 bucket | 404 |
| DELETE | `/buckets/:bucketId` | — | 204 | soft delete + enqueue |
| POST | `/buckets/:bucketId/files` | `fileId`, file, `Content-Range` | 201 file | 429 lock |
| PUT | `/buckets/:bucketId/files/:fileId` | `name`, `permissions` | 200 file | 404 |
| DELETE | `/buckets/:bucketId/files/:fileId` | — | 204 | 404 |
| GET | `/buckets/:bucketId/files/:fileId/download` | `Range` | 200/206 attachment | 416 |
| GET | `/buckets/:bucketId/files/:fileId/view` | `Range` | 200/206 inline | 416 |
| GET | `/buckets/:bucketId/files/:fileId/preview` | `width`, `height`, `gravity`, `quality`, `borderRadius`, `output`, `crop` | 200 image | 500 no imagick |
| GET | `/buckets/:bucketId/files/:fileId/push` | JWT in auth | 200 file | 401 |

## Upload contract

| Feature | Detail |
|---------|--------|
| Chunking | `Content-Range` header; `chunksTotal`/`chunksUploaded` tracking |
| Lock | `storage:file:{projectId}:{bucketId}:{fileId}` |
| Compression | GZIP or ZSTD at finalization |
| Encryption | AES-128-GCM at finalization |
| Size gates | `sizeOriginal <= APP_STORAGE_READ_BUFFER` for compress/encrypt; `<= APP_LIMIT_ANTIVIRUS` for ClamAV |

## Push JWT contract

| Field | Detail |
|-------|--------|
| Algorithm | HS256 |
| TTL | 1 hour |
| `internal` flag | `true` → dbForPlatform; `false` → dbForProject |

## Queue contract

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForDeletes` | Bucket/file delete → worker-deletes | Physical deletion + cache invalidation |

## Event contract

- `buckets.[bucketId].{create,update,delete}`
- `buckets.[bucketId].files.[fileId].{create,update,delete}`

## Cache-Control contract

Computed by `StorageCacheControl` based on bucket caching policy, file state (signed/expired), and token presence.
