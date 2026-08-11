# Design: Presences

> Unit: `Presences`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Presences module is a **lightweight HTTP CRUD service** over the `presenceLogs` collection in the project DB, with an in-memory TTL list cache for read hot paths. Records soft-expire via `expiresAt`; no background cleanup.

```
Client
  │ scopes: presences.read/write
  ▼
app/http.php → Presences\HTTP\** (5 actions)
  │
  ├── dbForProject (presenceLogs)
  ├── PresenceState (TTL list cache + purge)
  ├── queueForEvents (presence events)
  └── usage (METRIC_USERS_PRESENCE)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| PUT | `/v1/presences/:presenceId` | `upsertPresence` | presences.write |
| GET | `/v1/presences/:presenceId` | `getPresence` | presences.read |
| GET | `/v1/presences` | `listPresences` | presences.read |
| PATCH | `/v1/presences/:presenceId` | `updatePresence` | presences.write |
| DELETE | `/v1/presences/:presenceId` | `deletePresence` | presences.write |

## 3. Data flows — key actions

### 3.1 Upsert

```
1. Resolve identity: session → own userId; API key → explicit userId required
2. Validate status, metadata, expiresAt (default +15min, max 30d)
3. Detect source via x-appwrite-source (rest|graphql)
4. Upsert into presenceLogs (owner permissions default)
5. Emit presences.[presenceId].upsert; increment METRIC_USERS_PRESENCE
```

### 3.2 List with TTL cache

```
1. Build user query
2. Append server-side filter: expiresAt > now() (not part of cache key)
3. If ttl > 0: check PresenceState cache (per project/roles/query)
   - hit → X-Appwrite-Cache: hit
   - miss → query DB, store, X-Appwrite-Cache: miss
4. Return list
```

### 3.3 Update / Delete with purge

```
1. Fetch presence (expired → PRESENCE_NOT_FOUND)
2. Apply sparse update / delete
3. purge=true → PresenceState::purgeListCache
4. Emit update/delete event; decrement metric on delete
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | presenceLogs CRUD + cache |
| `response` | HTTP response |
| `request` | Source header detection |
| `user` | Session identity |
| `authorization` | Permission + role detection |
| `queueForEvents` | Lifecycle events |
| `usage` | METRIC_USERS_PRESENCE |

## 5. Events emitted

- `presences.[presenceId].upsert`
- `presences.[presenceId].update`
- `presences.[presenceId].delete` (payload: deleted doc)

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `PRESENCE_NOT_FOUND` | 404 | Missing or expired presence |
| Invalid expiry | 400 | > 30 days |
| Unauthorized | 401/403 | userId without privileges |

## 7. Security

- Session callers locked to own userId; API keys require explicit userId.
- Server-side expiry filter always applied.
