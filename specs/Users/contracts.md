# Contracts: Users (Admin)

> Unit: `Users`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/users`. Scopes: `users.*` (sessions also `sessions.write`). Auth: **ADMIN/KEY** only.

### 2.1 User CRUD + hash imports

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `` | `userId`, `email`, `phone`, `password`, `name` | 201 user | 409 email uniqueness, 400 policy |
| POST | `/argon2` | `userId`, `email`, `phone`, `password` (Argon2 hash) | 201 user | 400/409 |
| POST | `/bcrypt` | same + Bcrypt hash | 201 user | 400/409 |
| POST | `/md5` | same + MD5 hash | 201 user | 400/409 |
| POST | `/sha` | same + SHA hash | 201 user | 400/409 |
| POST | `/phpass` | same + PHPass hash | 201 user | 400/409 |
| POST | `/scrypt` | same + Scrypt hash | 201 user | 400/409 |
| POST | `/scrypt-modified` | same + Modified Scrypt hash | 201 user | 400/409 |
| GET | `/:userId` | — | 200 user | 404 |
| GET | `` | `queries`, `search` | 200 userList | — |
| DELETE | `/:userId` | — | 204 | 404 |

### 2.2 Profile updates

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| PATCH | `/:userId/status` | `status` (blocked boolean) | 200 user | 404 |
| PUT | `/:userId/labels` | `labels: string[]` | 200 user | 404 |
| PATCH | `/:userId/impersonation` | `impersonation` | 200 user | 404 |
| PATCH | `/:userId/name` | `name` | 200 user | 404 |
| PATCH | `/:userId/password` | `password` | 200 user | 400 policy, 404 |
| PATCH | `/:userId/email` | `email` | 200 user | 409 uniqueness, 404 |
| PATCH | `/:userId/phone` | `phone` | 200 user | 404 |
| PATCH | `/:userId/verification` | `verification` | 200 user | 404 |
| PATCH | `/:userId/verification/phone` | `verification` | 200 user | 404 |
| GET | `/:userId/prefs` | — | 200 prefs | 404 |
| PATCH | `/:userId/prefs` | `prefs` (JSON) | 200 prefs | 404 |

### 2.3 Sub-resources: targets, sessions, tokens, JWT, memberships, identities, MFA

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `/:userId/targets` | `targetId`, `providerType` (push/email/sms), `identifier` | 201 target | 404 |
| GET | `/:userId/targets/:targetId` | — | 200 target | 404 |
| GET | `/:userId/targets` | — | 200 targetList | — |
| PATCH | `/:userId/targets/:targetId` | `identifier` | 200 target | 404 |
| DELETE | `/:userId/targets/:targetId` | — | 204 | 404 |
| POST | `/:userId/sessions` | — | 201 session | 404 |
| GET | `/:userId/sessions` | — | 200 sessionList | — |
| DELETE | `/:userId/sessions/:sessionId` | — | 204 | 404 |
| DELETE | `/:userId/sessions` | — | 204 | — |
| POST | `/:userId/tokens` | `type` (magic-url/otp), `length` | 201 token | 404 |
| POST | `/:userId/jwts` | — | 201 jwt | 404 |
| GET | `/:userId/memberships` | — | 200 membershipList | — |
| GET | `/:userId/identities` | — | 200 identityList | — |
| DELETE | `/:userId/identities/:identityId` | — | 204 | 404 |
| PATCH | `/:userId/mfa` | `mfa` (boolean) | 200 user | 404 |
| GET | `/:userId/mfa/factors` | — | 200 factorList | — |
| GET | `/:userId/mfa/recovery-codes` | — | 200 recoveryCodes | — |
| POST | `/:userId/mfa/recovery-codes` | — | 201 recoveryCodes | 409 exists |
| PUT | `/:userId/mfa/recovery-codes` | — | 200 recoveryCodes | 404 |
| DELETE | `/:userId/mfa/authenticators/:type` | `:type` ∈ enrolled factors | 204 | 400 mismatch, 404 |

## Password import contract

| Endpoint | Hash name stored | Notes |
|----------|------------------|-------|
| `/argon2` | argon2 | Stored verbatim |
| `/bcrypt` | bcrypt | Stored verbatim |
| `/md5` | md5 | Stored verbatim |
| `/sha` | sha | Stored verbatim |
| `/phpass` | phpass | Stored verbatim |
| `/scrypt` | scrypt | Stored verbatim |
| `/scrypt-modified` | scrypt-modified | Stored verbatim |
| plain `/users` | password (Argon2id) | Hashed by server via `Utopia\Auth\Proofs\Password` |

## Session & JWT contract

| Aspect | Detail |
|--------|--------|
| Admin session | provider `SESSION_PROVIDER_SERVER`; no password verification; secret via `Token` proof → `Store` → opaque token |
| JWT | HS256, key `_APP_OPENSSL_KEY_V1`, claims scoped to userId + sessionId |

## Delete cascade contract

- Inline: hard-delete identities + targets.
- Queued: `DELETE_TYPE_DOCUMENT` → worker-deletes cleans sessions, memberships, logs, user-owned files.

## Event contract

- `users.[userId].{create,delete}`
- `users.[userId].update.{status,mfa,name,email,phone,password,verification,verification.phone,prefs,labels}`
- `users.[userId].sessions.[sessionId].{create,delete}`; `users.[userId].sessions.delete` (bulk)
- `users.[userId].targets.[targetId].{create,update,delete}`
- `users.[userId].tokens.[tokenId].create`

