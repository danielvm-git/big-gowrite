# Requirements: Tokens

> Identifier: `TKN-PUBLIC-LINKS`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Tokens module manages short-lived, secret-bearing resource tokens that grant public-URL access to protected Appwrite resources. Currently, the only implemented resource type is storage files (`TOKENS_RESOURCE_TYPE_FILES`). A token carries a 128-byte random secret, an optional expiry datetime, and is associated with a bucket+file pair. Tokens are stored in the `resourceTokens` project-level collection and are designed to be embedded as query-string parameters in file URLs, enabling temporary shareable links without requiring full authentication.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Tokens` | Module purpose | 🟢 |
| `specs/Tokens/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/domain.md#Token` | Token model | 🟢 |
| `specs/data-dictionary.md#resourceTokens` | Token document schema | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App user | Share a file | Create a token and append it to a download URL |
| App developer | Temporary access | Expiring shareable links without authentication |
| App developer | Revoke access | Delete a token to invalidate a link |

## 4. Business rules

1. **BR-01:** Tokens are linked to a bucket+file pair; `resourceId` stores `bucketId:fileId` and `resourceInternalId` stores the sequence-based compound key. 🟢
2. **BR-02:** Token secrets are 128-byte random values stored encrypted in the database (`encrypt` filter on the `secret` field). 🟢
3. **BR-03:** Creating a token requires either bucket-level UPDATE permission or, when `fileSecurity` is enabled on the bucket, file-level UPDATE permission. 🟢
4. **BR-04:** `expire` is optional and must be a future datetime when provided; `null` means the token never expires. 🟢
5. **BR-05:** Listing tokens filters by `resourceType = TOKENS_RESOURCE_TYPE_FILES` and the compound `resourceInternalId`. 🟢
6. **BR-06:** Only `ADMIN` and `KEY` auth types are allowed on all endpoints (no session/JWT access). 🟡
7. **BR-07:** Token access is constrained by the bucket's enabled state; disabled buckets throw `STORAGE_BUCKET_NOT_FOUND` unless the caller is an API key or privileged user. 🟡
8. **BR-08:** Abuse-rate limiting is applied on write operations (create, update, delete) using the default write-rate constants. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Create file token | Must | Bucket/file resolution + permission checks; encrypted secret | 🟢 |
| RF-02 | List tokens for a file | Must | `resourceType` + compound-key filter | 🟢 |
| RF-03 | Get token by ID | Must | Returns token document | 🟢 |
| RF-04 | Update token expiry | Must | Future datetime or null; otherwise invalid | 🟢 |
| RF-05 | Delete token | Must | Revokes the shareable link | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Secret encrypted at rest | `encrypt` filter on `secret` | 🟢 |
| Security | Restricted auth types | ADMIN/KEY only | 🟡 |
| Reliability | Write rate limiting | default write-rate constants | 🟢 |
| Governança | Resource scoping | bucket+file pair, never cross-resource | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Create a file token
  Given bucket UPDATE permission (or file-level when fileSecurity is enabled)
  When POST /v1/tokens/buckets/:bucketId/files/:fileId is called
  Then a token with an encrypted 128-byte random secret is created

Scenario: Update expiry to a past date
  Given an existing token
  When PATCH /v1/tokens/:tokenId sets expire in the past
  Then GENERAL_ARGUMENT_INVALID is returned

Scenario: List tokens for a file
  Given tokens on a bucket+file pair
  When GET /v1/tokens/buckets/:bucketId/files/:fileId is called
  Then only FILES-type tokens for that compound key are returned

Scenario: Disabled bucket access
  Given a disabled bucket and a non-privileged caller
  When any token operation is attempted
  Then STORAGE_BUCKET_NOT_FOUND is returned
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..05 (core) | Must | Token lifecycle |
| Rate limiting (RNF) | Must | Abuse protection on writes |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Exact auth-type enforcement matrix (session/JWT exclusion) not exhaustively verified.
- 🟡 [INFERRED] Downstream consumption of the token query-string by Storage download/view endpoints.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
