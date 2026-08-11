# Design: Users (Admin)

> Unit: `Users`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Users module is a **large HTTP-only administrative service** over the project database (`users`, `identities`, `targets`, `sessions`, `tokens`, `authenticators`, `memberships`). It owns no workers; asynchronous cleanup is delegated to the delete worker via `publisherForDeletes`. All 47 endpoints are Admin/API-key restricted.

```
Client (admin console / API key)
  │ scopes: users.* (+ sessions.write)
  ▼
app/http.php → Users\Http\** (47 actions)
  │
  ├── dbForProject (users, identities, targets, sessions, tokens, ...)
  ├── dbForPlatform (console project context)
  ├── queueForEvents (webhooks, function triggers)
  ├── publisherForDeletes (DELETE_TYPE_DOCUMENT)
  ├── hooks (passwordValidator)
  ├── plan (email policy gates)
  ├── store / proofForToken (session + token secrets)
  └── usage / locale / geoRecord (metrics, session metadata)
```

## 2. Route-to-action map

### 2.1 User CRUD + hash imports

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/users` | `createUser` | users.write |
| POST | `/v1/users/argon2` | `createArgon2User` | users.write |
| POST | `/v1/users/bcrypt` | `createBcryptUser` | users.write |
| POST | `/v1/users/md5` | `createMD5User` | users.write |
| POST | `/v1/users/sha` | `createSHAUser` | users.write |
| POST | `/v1/users/phpass` | `createPHPassUser` | users.write |
| POST | `/v1/users/scrypt` | `createScryptUser` | users.write |
| POST | `/v1/users/scrypt-modified` | `createScryptModifiedUser` | users.write |
| GET | `/v1/users/:userId` | `getUser` | users.read |
| GET | `/v1/users` | `listUsers` | users.read |
| DELETE | `/v1/users/:userId` | `deleteUser` | users.write |

### 2.2 Profile updates

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| PATCH | `/v1/users/:userId/status` | `updateUserStatus` | users.write |
| PUT | `/v1/users/:userId/labels` | `updateUserLabels` | users.write |
| PATCH | `/v1/users/:userId/impersonation` | `updateUserImpersonation` | users.write |
| PATCH | `/v1/users/:userId/name` | `updateUserName` | users.write |
| PATCH | `/v1/users/:userId/password` | `updateUserPassword` | users.write |
| PATCH | `/v1/users/:userId/email` | `updateUserEmail` | users.write |
| PATCH | `/v1/users/:userId/phone` | `updateUserPhone` | users.write |
| PATCH | `/v1/users/:userId/verification` | `updateUserVerification` | users.write |
| PATCH | `/v1/users/:userId/verification/phone` | `updateUserPhoneVerification` | users.write |
| GET | `/v1/users/:userId/prefs` | `getUserPrefs` | users.read |
| PATCH | `/v1/users/:userId/prefs` | `updateUserPrefs` | users.write |

### 2.3 Sub-resources: targets, sessions, tokens, JWT, memberships, identities, MFA

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/users/:userId/targets` | `createUserTarget` | users.write |
| GET | `/v1/users/:userId/targets/:targetId` | `getUserTarget` | users.read |
| GET | `/v1/users/:userId/targets` | `listUserTargets` | users.read |
| PATCH | `/v1/users/:userId/targets/:targetId` | `updateUserTarget` | users.write |
| DELETE | `/v1/users/:userId/targets/:targetId` | `deleteUserTarget` | users.write |
| POST | `/v1/users/:userId/sessions` | `createUserSession` | users.write, sessions.write |
| GET | `/v1/users/:userId/sessions` | `listUserSessions` | users.read |
| DELETE | `/v1/users/:userId/sessions/:sessionId` | `deleteUserSession` | users.write |
| DELETE | `/v1/users/:userId/sessions` | `deleteUserSessions` | users.write |
| POST | `/v1/users/:userId/tokens` | `createUserToken` | users.write |
| POST | `/v1/users/:userId/jwts` | `createUserJWT` | users.write |
| GET | `/v1/users/:userId/memberships` | `listUserMemberships` | users.read |
| GET | `/v1/users/:userId/identities` | `listUserIdentities` | users.read |
| DELETE | `/v1/users/:userId/identities/:identityId` | `deleteUserIdentity` | users.write |
| PATCH | `/v1/users/:userId/mfa` | `updateUserMFA` | users.write |
| GET | `/v1/users/:userId/mfa/factors` | `listUserMFAFactors` | users.read |
| GET | `/v1/users/:userId/mfa/recovery-codes` | `getUserMFARecoveryCodes` | users.read |
| POST | `/v1/users/:userId/mfa/recovery-codes` | `createUserMFARecoveryCodes` | users.write |
| PUT | `/v1/users/:userId/mfa/recovery-codes` | `updateUserMFARecoveryCodes` | users.write |
| DELETE | `/v1/users/:userId/mfa/authenticators/:type` | `deleteUserMFAAuthenticator` | users.write |

## 3. Data flows — key actions

### 3.1 Create user (plaintext) — shared by all Create actions via `Base.php`

```
1. Normalize + lowercase email; uniqueness check across users + identities
2. personalDataCheck (if enabled): password vs id/email/name/phone
3. Email policy gates (disposable/canonical/free/corporate) plan-gated
4. Hash password (Argon2id default); store hash name
5. Create user document
6. Auto-create targets for email and/or phone
7. Emit users.[userId].create
```

### 3.2 Import with legacy hash

```
1. Same validation as 3.1, but password stored verbatim with hash name
2. Hash name persists for future verification (users.hash)
```

### 3.3 Admin session create

```
1. Skip password verification
2. provider = SESSION_PROVIDER_SERVER
3. secret = proofForToken generated, encoded via store
4. Attach user-agent + geoRecord country via locale
5. Emit users.[userId].sessions.[sessionId].create
```

### 3.4 JWT create

```
1. Sign HS256 with _APP_OPENSSL_KEY_V1
2. Claims scoped to userId + sessionId
```

### 3.5 Delete user

```
1. Hard-delete identities + targets inline
2. Enqueue DELETE_TYPE_DOCUMENT (delete worker cleans sessions,
   memberships, logs, user-owned files)
3. Emit users.[userId].delete
```

### 3.6 List users (performance path)

```
1. Paginated query with search/queries support
2. If no select restriction: bulk-load targets by userInternalId (no N+1)
3. skipFilters() for sessions/tokens/authenticators/challenges/memberships sub-queries
```


## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | users, identities, targets, sessions, tokens, authenticators, memberships |
| `dbForPlatform` | Console project context |
| `queueForEvents` | Webhooks / function triggers |
| `publisherForDeletes` | Async cleanup queue |
| `hooks` | passwordValidator hook |
| `plan` | Email policy capability flags |
| `store`, `proofForToken` | Session/token secrets |
| `locale`, `geoRecord`, `request` | Session metadata (country, user-agent) |
| `usage` | Metrics |
| `response`, `project`, `user` | Standard context |

## 5. Events emitted

- `users.[userId].{create,delete}`
- `users.[userId].update.{status,mfa,name,email,phone,password,verification,verification.phone,prefs,labels}`
- `users.[userId].sessions.[sessionId].{create,delete}` and `users.[userId].sessions.delete` (bulk)
- `users.[userId].targets.[targetId].{create,update,delete}`
- `users.[userId].tokens.[tokenId].create`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| Email uniqueness | 409 | Email already in users or identities |
| Invalid password policy | 400 | personalDataCheck or password history |
| Blocked email domain | 400 | Plan-gated email policy |
| User not found | 404 | Unknown `:userId` |
| MFA authenticator mismatch | 400 | Deletion of non-enrolled type |

All errors follow the Appwrite model: `{ message, code, type, version }`.

## 7. Security

- Admin/API-key only — no end-user sessions on these routes.
- Argon2id default hashing; imports stored verbatim with hash name.
- Password history enforcement via hooks.
- JWT HS256 with `_APP_OPENSSL_KEY_V1`.
- Deletion cascade via queue for asynchronous, non-blocking cleanup.

