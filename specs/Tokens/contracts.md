# Contracts: Tokens

> Unit: `Tokens`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/tokens`. Scopes: `tokens.*`; auth types: **ADMIN, KEY** only (no session/JWT).

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `/buckets/:bucketId/files/:fileId` | `expire` (optional, future) | 201 token | STORAGE_BUCKET_NOT_FOUND, GENERAL_ARGUMENT_INVALID, 429 |
| GET | `/buckets/:bucketId/files/:fileId` | `queries` | 200 tokenList | STORAGE_BUCKET_NOT_FOUND |
| GET | `/:tokenId` | — | 200 token | 404 |
| PATCH | `/:tokenId` | `expire` (future or null) | 200 token | GENERAL_ARGUMENT_INVALID, 404, 429 |
| DELETE | `/:tokenId` | — | 204 | 404, 429 |

## Token document contract

| Field | Type | Detail |
|-------|------|--------|
| `resourceType` | string | `files` (`TOKENS_RESOURCE_TYPE_FILES`) |
| `resourceId` | string | `bucketId:fileId` |
| `resourceInternalId` | string | Sequence-based compound key |
| `secret` | bytes | 128-byte random; encrypted at rest (`encrypt` filter) |
| `expire` | datetime \| null | Future datetime or never |

## Usage contract

| Aspect | Detail |
|--------|--------|
| Embedding | Token appended as query-string parameter on file URLs |
| Scoping | Bucket+file pair only — no cross-resource tokens today |
| Revocation | Delete token or let `expire` pass |

## Event contract

- `tokens.[tokenId].{create,update,delete}`
