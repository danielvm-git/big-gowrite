# Requirements: Presences

> Identifier: `PRS-REALTIME`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Presences module tracks real-time user presence state within an Appwrite project. Each record captures status (free-form string), optional metadata, source (`rest`/`graphql`), and expiry timestamp, stored in the `presenceLogs` common collection. It supports upsert (PUT), partial update (PATCH), get, list (with TTL-backed in-memory cache), and delete. Presence records are soft-expired — filtered out on read after `expiresAt` passes, not auto-deleted.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Presences` | Module purpose | 🟢 |
| `specs/Presences/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/domain.md#Presence` | Presence model, expiry | 🟢 |
| `specs/architecture.md#Realtime` | Realtime integration | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App user | Show presence | Set status; others list active users |
| App developer | High-frequency reads | List with TTL cache |
| App user | Leave room | Delete presence / expire |

## 4. Business rules

1. **BR-01:** `userId` in write ops only permitted for API key/privileged users; session callers operate on own identity. 🟢
2. **BR-02:** API keys must explicitly provide `userId`. 🟢
3. **BR-03:** Default expiry 15 minutes when not specified. 🟢
4. **BR-04:** Maximum expiry 30 days (DatetimeValidator bounds). 🟢
5. **BR-05:** Expired presences excluded from GET/LIST (`expiresAt < now()`); treated as `PRESENCE_NOT_FOUND`. 🟢
6. **BR-06:** Permissions default to owner full access; override via `permissions` param. 🟢
7. **BR-07:** `listPresences` supports `ttl` (0–86400 s); when >0 results cached per project/roles/query via `PresenceState`; response carries `X-Appwrite-Cache: hit|miss`. 🟢
8. **BR-08:** List query always appends `expiresAt > now()` server-side (outside user query, not polluting cache key). 🟢
9. **BR-09:** update/delete accept `purge` flag → `PresenceState::purgeListCache`. 🟢
10. **BR-10:** `source` field records `rest`/`graphql` via `x-appwrite-source` header. 🟢
11. **BR-11:** `presenceLogs` name suggests future event-log expansion. 🟡

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Upsert presence | Must | Create/replace; expiry default 15min, max 30d; source detection | 🟢 |
| RF-02 | Get presence | Must | Expired → PRESENCE_NOT_FOUND | 🟢 |
| RF-03 | List presences with TTL cache | Must | Cache hit/miss header; expiresAt>now appended | 🟢 |
| RF-04 | Partial update presence | Must | Sparse update; purge option | 🟢 |
| RF-05 | Delete presence | Must | Delete + purge cache + metric decrement | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Performance | TTL in-memory list cache | `PresenceState` cache fields + X-Appwrite-Cache | 🟢 |
| Security | userId scoping for sessions vs API keys | role detection | 🟢 |
| Observability | Usage metrics on create/delete | METRIC_USERS_PRESENCE | 🟢 |
| Lifecycle | Soft-expiry without background cleanup | expiresAt filtering | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Session user writes own presence
  Given a session-authenticated user
  When PUT /v1/presences/:presenceId with a different userId is called
  Then the session identity is used and the request is rejected/ignored for userId

Scenario: Expired presence hidden
  Given a presence with expiresAt in the past
  When GET /v1/presences/:presenceId is called
  Then PRESENCE_NOT_FOUND is returned

Scenario: List cache hit
  Given a cached list with ttl > 0
  When the same list request is made
  Then the response carries X-Appwrite-Cache: hit
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..05 (CRUD) | Must | Core presence |
| TTL cache (RNF) | Should | Read performance |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] `presenceLogs` future event-log expansion.
- 🔴 [GAP] Cache eviction semantics of PresenceState not fully captured.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
