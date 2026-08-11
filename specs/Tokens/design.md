# Design: Tokens

> Unit: `Tokens`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Tokens module is a **small HTTP-only service** over the `resourceTokens` project collection. It owns no workers or CLI tasks — all work happens inline in the action layer. Endpoints are guarded by `tokens.*` scopes and restricted to ADMIN/KEY auth types.

```
Client (API key / admin)
  │ scopes: tokens.*
  ▼
app/http.php → Tokens\Http\** (5 actions)
  │
  ├── dbForProject (resourceTokens, buckets, file-bucket collections)
  ├── authorization (permission checks)
  ├── queueForEvents (lifecycle events)
  ├── user (auth context)
  └── usage (metrics)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/tokens/buckets/:bucketId/files/:fileId` | `createFileToken` | tokens.write |
| GET | `/v1/tokens/buckets/:bucketId/files/:fileId` | `listFileTokens` | tokens.read |
| GET | `/v1/tokens/:tokenId` | `getToken` | tokens.read |
| PATCH | `/v1/tokens/:tokenId` | `updateToken` | tokens.write |
| DELETE | `/v1/tokens/:tokenId` | `deleteToken` | tokens.write |

## 3. Data flows — key actions

### 3.1 Create file token

```
1. Resolve bucket + file via getFileAndBucket() (Action.php):
   enabled-state check, read permission, file existence
2. Permission gate: bucket UPDATE or file UPDATE (fileSecurity)
3. Generate 128-byte random secret → encrypt filter at write time
4. Store resourceId = bucketId:fileId, resourceInternalId = compound key
5. Rate limit (write constants)
6. Emit tokens.[tokenId].create
```

### 3.2 List tokens for a file

```
1. Resolve bucket + file (same guard)
2. Filter resourceType = TOKENS_RESOURCE_TYPE_FILES
3. Filter resourceInternalId = compound key
4. Return paginated list
```

### 3.3 Update expiry

```
1. Get token by :tokenId
2. Validate expire: null or future datetime
3. Update sparse document
4. Emit tokens.[tokenId].update
```

### 3.4 Delete token

```
1. Get token by :tokenId
2. Delete document
3. Emit tokens.[tokenId].delete
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | resourceTokens, buckets, file collections |
| `response` | HTTP response builder |
| `user` | Authenticated user context |
| `queueForEvents` | Lifecycle events |
| `authorization` | Permission checks |
| `usage` | Metrics |

## 5. Events emitted

- `tokens.[tokenId].create` — file token created
- `tokens.[tokenId].update` — expiry updated
- `tokens.[tokenId].delete` — token deleted

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `STORAGE_BUCKET_NOT_FOUND` | 404 | Disabled bucket for non-privileged caller |
| Token not found | 404 | Invalid `:tokenId` |
| `GENERAL_ARGUMENT_INVALID` | 400 | `expire` in the past |
| Rate limit | 429 | Write-operation limits |

## 7. Security

- Secrets encrypted at rest; plaintext only in memory during creation.
- Only ADMIN/KEY auth types — no session/JWT access.
- Write operations rate-limited.
- Bucket enabled-state enforced on every operation.
