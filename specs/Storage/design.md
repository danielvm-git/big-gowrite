# Design: Storage

> Unit: `Storage`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Storage module is an **HTTP object-storage service** with two tiers (buckets, files). Uploads assemble chunks under a distributed lock, then finalize with compression/encryption and antivirus. Served via three read paths (download/view/preview) plus a JWT push path. Physical deletion is async via the deletes worker.

```
Client
  │ scopes: buckets.*, files.*
  ▼
app/http.php → Storage\Http\** (14 actions)
  │
  ├── dbForProject (buckets, bucket_{seq} collections)
  ├── deviceForFiles / deviceForLocal
  ├── locks (upload concurrency)
  ├── Imagick (preview)
  ├── publisherForDeletes (physical cleanup)
  └── cacheControlForStorage (Cache-Control)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/storage/buckets` | `createBucket` | buckets.write |
| GET | `/v1/storage/buckets` | `listBuckets` | buckets.read |
| GET | `/v1/storage/buckets/:bucketId` | `getBucket` | buckets.read |
| PUT | `/v1/storage/buckets/:bucketId` | `updateBucket` | buckets.write |
| DELETE | `/v1/storage/buckets/:bucketId` | `deleteBucket` | buckets.write |
| POST | `/v1/storage/buckets/:bucketId/files` | `createFile` | files.write |
| GET | `/v1/storage/buckets/:bucketId/files` | `listFiles` | files.read |
| GET | `/v1/storage/buckets/:bucketId/files/:fileId` | `getFile` | files.read |
| PUT | `/v1/storage/buckets/:bucketId/files/:fileId` | `updateFile` | files.write |
| DELETE | `/v1/storage/buckets/:bucketId/files/:fileId` | `deleteFile` | files.write |
| GET | `/v1/storage/buckets/:bucketId/files/:fileId/download` | `getFileDownload` | files.read |
| GET | `/v1/storage/buckets/:bucketId/files/:fileId/view` | `getFileView` | files.read |
| GET | `/v1/storage/buckets/:bucketId/files/:fileId/preview` | `getFilePreview` | files.read |
| GET | `/v1/storage/buckets/:bucketId/files/:fileId/push` | `getFileForPush` | public |

## 3. Data flows — key actions

### 3.1 Bucket create

```
1. Validate settings (security, compression, encryption, antivirus, permissions)
2. Create bucket_{sequence} collection (files schema)
3. Create bucket doc; emit buckets.[bucketId].create
```

### 3.2 Chunked file upload

```
1. Acquire lock storage:file:{projectId}:{bucketId}:{fileId}
2. Append chunk to deviceForLocal (track chunksTotal/chunksUploaded)
3. On final chunk:
   a. sizeOriginal <= APP_STORAGE_READ_BUFFER → compress (GZIP/ZSTD) + AES-128-GCM encrypt
   b. ClamAV scan if enabled + size <= APP_LIMIT_ANTIVIRUS
   c. Move to deviceForFiles
4. Release lock; emit files.create event
```

### 3.3 Download / View

```
1. Resolve auth: session/API key or resource token
2. fileSecurity ? per-file permissions : bucket permissions
3. Stream with HTTP Range support
4. Cache-Control via StorageCacheControl
```

### 3.4 Preview

```
1. Require imagick extension; bucket.transformations enabled?
2. Apply transforms (crop/resize/rotate/quality...)
3. Stream transformed image
```

### 3.5 Push

```
1. Validate signed JWT (HS256, 1h TTL)
2. internal flag → dbForPlatform else dbForProject
3. Stream file
```

### 3.6 File delete

```
1. Remove physical file (or abort partial upload)
2. Delete DB doc
3. Enqueue cache-invalidation message
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | Bucket + file docs |
| `dbForPlatform` | Push endpoint (internal=true) |
| `queueForEvents` | Lifecycle events |
| `publisherForDeletes` | Physical cleanup |
| `deviceForFiles` | Project file storage |
| `deviceForLocal` | Temp upload storage |
| `authorization` / `user` | RBAC |
| `resourceToken` | Tokenized access |
| `locks` | Upload lock factory |
| `cacheControlForStorage` | Cache-Control computation |

## 5. Events emitted

- `buckets.[bucketId].{create,update,delete}`
- `buckets.[bucketId].files.[fileId].{create,update,delete}`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| Bucket disabled | 403 | End-user on enabled=false |
| Imagick missing | 500 | Preview without extension |
| Invalid JWT | 401 | Push endpoint |
| Range invalid | 416 | Bad Range header |
