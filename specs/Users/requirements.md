# Requirements: Users (Admin)

> Identifier: `USR-ADMIN`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Users module is the **administrative surface** for managing project users from the server side. It exposes CRUD over user accounts, password-hash imports (Bcrypt, MD5, Argon2, SHA, PHPass, Scrypt, Scrypt Modified), and every sub-resource attached to a user: sessions, targets (push/email/SMS notification endpoints), JWT and one-time login tokens, memberships, OAuth2 identities, preferences, labels, status flags, multi-factor authentication, and MFA recovery codes. All endpoints require `users.read`/`users.write` scope and are restricted to Admin or API Key auth — end users cannot call these routes directly.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Users` | Module purpose | 🟢 |
| `specs/Users/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/domain.md#User` | User model | 🟢 |
| `specs/security.md#Hashing` | Password hashing policies | 🟢 |
| `specs/state-machines.md#User` | User status/verification transitions | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Admin | Onboard users | Create a user or import one with a legacy hash |
| Admin | Support | Reset password, block user, issue admin session |
| Admin | Compliance | MFA recovery codes, OAuth2 identity cleanup |
| App developer | Auth flows | Issue a JWT or magic-link token for a user |

## 4. Business rules

1. **BR-01:** Email is lower-cased and checked for uniqueness across both `users` and `identities` collections before creation. 🟢
2. **BR-02:** If project `auths.personalDataCheck` is enabled, the password must not contain the user's id, email, name, or phone. 🟢
3. **BR-03:** Email policy flags (`disposableEmails`, `canonicalEmails`, `freeEmails`, `corporateEmails`) block certain domains; enforcement gated by `plan` capabilities. 🟢
4. **BR-04:** Plain-text passwords hashed with `Utopia\Auth\Proofs\Password` (Argon2id default); import endpoints store the already-hashed form verbatim with the corresponding hash name. 🟢
5. **BR-05:** Password history maintained if `auths.passwordHistory > 0`; new passwords checked against history via hooks. 🟢
6. **BR-06:** Creating a user with an email or phone auto-creates a `targets` document (providerType `email`/`sms`). 🟢
7. **BR-07:** Admin-created sessions use `SESSION_PROVIDER_SERVER` and skip password verification; the secret is generated via `Token` proof, encoded into the `Store`, and returned as an opaque session token. 🟢
8. **BR-08:** JWT tokens created via `createUserJWT` are signed with `_APP_OPENSSL_KEY_V1` (HS256) and scoped to userId + sessionId. 🟢
9. **BR-09:** User deletion hard-deletes identities and targets inline, then enqueues a `DELETE_TYPE_DOCUMENT` message for the delete worker to clean up sessions, memberships, logs, and user-owned files. 🟢
10. **BR-10:** MFA update stores a boolean `mfa` attribute; MFA authenticator deletion requires verifying the authenticator type matches an enrolled TOTP/email/phone factor. 🟢
11. **BR-11:** `listUsers` optimizes target loading — when no `select` query restricts fields, targets are fetched in a single bulk query by `userInternalId` (avoids N+1). 🟡
12. **BR-12:** Sub-queries for sessions, tokens, authenticators, challenges, and memberships are skipped using `skipFilters()` during list operations for performance. 🟡

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Create user (plaintext) | Must | Hashing + validation + auto-targets | 🟢 |
| RF-02 | Create user (hash import, 7 variants) | Must | Hash stored verbatim with name | 🟢 |
| RF-03 | Get / list / delete user | Must | Pagination, search, hard-delete cascade | 🟢 |
| RF-04 | Profile updates (name/email/phone/password/status/labels/prefs/verification/impersonation) | Must | Sparse updates + events | 🟢 |
| RF-05 | Sessions create / list / delete | Must | Server-provider sessions, bulk delete | 🟢 |
| RF-06 | Targets CRUD | Must | push/email/sms endpoints | 🟢 |
| RF-07 | Tokens + JWT create | Must | Magic-URL token + signed JWT | 🟢 |
| RF-08 | Memberships list / identities list+delete | Must | Scoped reads | 🟢 |
| RF-09 | MFA management | Must | Toggle, factors, recovery codes, authenticator removal | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Password hashing policy | Argon2id default + named imports | 🟢 |
| Security | Personal-data password check | `auths.personalDataCheck` | 🟢 |
| Governança | Email policy plan gates | disposable/canonical/free/corporate | 🟢 |
| Performance | N+1 avoidance on list | bulk target load + `skipFilters()` | 🟡 |
| Auditability | Events for all mutations | `queueForEvents` | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Create user with duplicate email
  Given an existing user with email user@example.com
  When POST /v1/users is called with the same email
  Then creation fails with an email-uniqueness error

Scenario: Import user with legacy hash
  Given a SHA-256 hex hash
  When POST /v1/users/sha is called with that hash
  Then the hash is stored verbatim with hash name sha

Scenario: Admin session without password
  Given an admin with users.write scope
  When POST /v1/users/:userId/sessions is called
  Then a session with SESSION_PROVIDER_SERVER is created, no password verification

Scenario: Delete user cascade
  When DELETE /v1/users/:userId is called
  Then identities and targets are hard-deleted inline and DELETE_TYPE_DOCUMENT is enqueued

Scenario: personalDataCheck violation
  Given auths.personalDataCheck enabled
  When a password containing the user's email is submitted
  Then the password is rejected
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..09 (core) | Must | Admin user lifecycle |
| listUsers optimizations | Should | Performance, not correctness |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Exact error codes for email uniqueness across `identities`.
- 🔴 [GAP] Full 47-endpoint OpenAPI schemas not exhaustively verified.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |

