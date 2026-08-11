# Design: Account (MFA)

> Unit: `Account` · Module: MFA subsystem
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Account MFA module is a **self-contained HTTP service** within the Appwrite modular monolith. It serves as a REST API layer over the user sessions and authenticators stored in the project database. It does **not** own workers or tasks — all async behavior (email, SMS) is delegated to external publishers (`publisherForMails`, `publisherForMessaging`).

```
Browser/Client
  │
  ▼
Traefik (reverse proxy)
  │
  ▼
app/http.php → Account\Http\* (10 actions)
  │
  ├── dbForProject (MongoDB/MariaDB)
  │     ├── users collection (mfa, mfaRecoveryCodes)
  │     ├── sessions collection (factors, mfaUpdatedAt)
  │     └── authenticators collection (TOTP secrets)
  ├── queueForEvents → Redis → worker-events
  ├── publisherForMails → Redis → worker-mails
  ├── publisherForMessaging → Redis → worker-messaging
  └── timelimit (rate limiter)
```

## 2. Route-to-action map

| Method | Path | Action class | Action ID | Groups | Scope |

## 3. Directory structure

```
Modules/Account/
├── Module.php                        # register(Http::class)
├── Services/Http.php                 # 10 route registrations
└── Http/Account/MFA/
    ├── Update.php                    # enable/disable MFA
    ├── Factors/XList.php             # list factors
    ├── Authenticators/
    │   ├── Create.php                # create TOTP authenticator
    │   ├── Update.php                # confirm TOTP authenticator
    │   └── Delete.php                # delete TOTP authenticator
    ├── Challenges/
    │   ├── Create.php                # initiate challenge
    │   └── Update.php                # complete challenge
    └── RecoveryCodes/
        ├── Create.php                # generate initial codes
        ├── Get.php                   # retrieve codes
        └── Update.php                # regenerate codes
```

## 4. Data flow — key actions

### 4.1 Enable/disable MFA (`MFA/Update`)

```
1. Receive { mfa: boolean }
2. Validate input
3. Read user document from dbForProject
4. Set user.mfa = value
5. Recalculate session.factors (check TOTP, email, phone verified)
6. Save user document
7. Emit users.[userId].update.mfa event
8. Return updated user
```

### 4.2 Create TOTP authenticator (`Authenticators/Create`)

```
1. Validate :type ∈ ['totp']
2. Read existing authenticators for user
3. If verified authenticator exists → 409 error
4. If unverified exists → delete it
5. Generate TOTP secret
6. Build provisioning URI (otpauth://totp/{issuer}:{email}?secret=...)
7. Create authenticator document in dbForProject
8. Emit users.[userId].update.mfa event
9. Return { type: "totp", secret, uri }
```

### 4.3 Confirm TOTP authenticator (`Authenticators/Update`)

```
1. Read unverified authenticator from dbForProject
2. Receive { otp } from request
3. Verify OTP against stored secret
4. If invalid → 401 error
5. Set verified = true on authenticator
6. Append "totp" to session.factors (deduplicate)
7. Emit users.[userId].update.mfa event
8. Return updated authenticator
```

|--------|------|-------------|-----------|--------|-------|
| PATCH | `/v1/account/mfa` | `Account/MFA/Update` | `updateMFA` | `api`, `account` | `account` |
| GET | `/v1/account/mfa/factors` | `Account/MFA/Factors/XList` | `listMFAFactors` | `api`, `account` | `account` |
| POST | `/v1/account/mfa/authenticators/:type` | `Account/MFA/Authenticators/Create` | `createMFAAuthenticator` | `api`, `account` | `account` |
| PUT | `/v1/account/mfa/authenticators/:type` | `Account/MFA/Authenticators/Update` | `updateMFAAuthenticator` | `api`, `account` | `account` |
| DELETE | `/v1/account/mfa/authenticators/:type` | `Account/MFA/Authenticators/Delete` | `deleteMFAAuthenticator` | `api`, `account`, `mfaProtected` | `account` |
| POST | `/v1/account/mfa/challenges` | `Account/MFA/Challenges/Create` | `createMFAChallenge` | `api`, `account` | `account` |
| PUT | `/v1/account/mfa/challenges` | `Account/MFA/Challenges/Update` | `updateMFAChallenge` | `api`, `account` | `account` |
| POST | `/v1/account/mfa/recovery-codes` | `Account/MFA/RecoveryCodes/Create` | `createMFARecoveryCodes` | `api`, `account` | `account` |
| GET | `/v1/account/mfa/recovery-codes` | `Account/MFA/RecoveryCodes/Get` | `getMFARecoveryCodes` | `api`, `account`, `mfaProtected` | `account` |
| PATCH | `/v1/account/mfa/recovery-codes` | `Account/MFA/RecoveryCodes/Update` | `updateMFARecoveryCodes` | `api`, `account`, `mfaProtected` | `account` |

### 4.4 Initiate MFA challenge (`Challenges/Create`)

```
1. Rate limit: 10/(url, userId)
2. Receive { factor: "email"|"phone"|"totp"|"recovery_code" }
3. Validate factor availability:
   - email: SMTP configured + email verified
   - phone: SMS provider + phone verified
   - totp: authenticator exists and verified
   - recovery_code: codes exist on user document
4. Generate challenge token + code (proofForToken, proofForCode)
5. If email: publish to publisherForMails
6. If phone: publish to publisherForMessaging
7. Create challenge document in dbForProject
8. Emit challenges.[challengeId].create event
9. Return { id, token (partial) }
```

### 4.5 Complete MFA challenge (`Challenges/Update`)

```
1. Rate limit: 10/(url, challengeId)
2. Read challenge document from dbForProject
3. Validate { challengeId, otp/code } from request
4. Verify code against challenge
5. If recovery_code: remove consumed code from mfaRecoveryCodes
6. Append factor to session.factors (deduplicate)
7. Set session.mfaUpdatedAt = now
8. Save session in dbForProject
9. Emit sessions.update event
10. Return updated session
```

## 5. Dependency injection

| Injection | Role | Used by |
|-----------|------|---------|
| `response` | HTTP response builder | All actions |
| `request` | HTTP request (user-agent) | All actions |
| `user` | Current authenticated user | All actions |
| `session` | Current session (factors stored) | MFA/Update, Authenticators/Update, Challenges/Update |
| `project` | Project config (issuer, SMTP) | MFA/Update, Authenticators/Create, Challenges/Create |
| `dbForProject` | Project-scoped database | All actions |
| `queueForEvents` | Event queue | All mutating actions |
| `publisherForMails` | Email OTP delivery | Challenges/Create |
| `publisherForMessaging` | SMS OTP delivery | Challenges/Create |
| `locale` | i18n for email/SMS templates | Challenges/Create |
| `platform` | Platform config (branding) | Challenges/Create |
| `timelimit` | Rate limiter callable | Challenges/Create, Challenges/Update |
| `usage` | Usage metrics (phone country) | Challenges/Create |
| `plan` | Project plan config | MFA/Update |
| `proofForToken` | Token generator | Challenges/Create |
| `proofForCode` | OTP code generator | Challenges/Create |
| `requestTimestamp` | Request timestamp | MFA/Update |

## 6. Events emitted

| Event pattern | Trigger | Payload |
|---------------|---------|---------|
| `users.[userId].update.mfa` | MFA toggle, authenticator CRUD, recovery codes CRUD | user document snapshot |
| `users.[userId].delete.mfa` | TOTP authenticator deleted | user document snapshot |
| `users.[userId].challenges.[challengeId].create` | Challenge initiated | challengeId |
| `users.[userId].sessions.[sessionId].create` | Challenge completed (session elevation) | session snapshot |

## 7. Error handling

| Error code | HTTP | Trigger |
|------------|------|---------|
| `USER_AUTHENTICATOR_ALREADY_VERIFIED` | 409 | Creating TOTP when verified exists |
| `USER_AUTHENTICATOR_NOT_FOUND` | 404 | Confirming/deleting non-existent authenticator |
| `USER_CHALLENGE_NOT_FOUND` | 404 | Completing non-existent challenge |
| `USER_RECOVERY_CODES_ALREADY_EXISTS` | 409 | Creating when codes exist |
| `USER_RECOVERY_CODES_NOT_FOUND` | 404 | Getting/regenerating non-existent codes |
| Rate limit exceeded | 429 | Challenge creation/completion limits |
| Invalid OTP | 401 | Wrong TOTP during confirmation |

All errors follow Appwrite model: JSON `{ message, code, type, version }`.

## 8. Security

- **Scope:** All endpoints use `account` scope — requires valid session token (JWT or API key).
- **mfaProtected:** `deleteMFAAuthenticator`, `getMFARecoveryCodes`, `updateMFARecoveryCodes` require MFA-completed session.
- **Secrets:** TOTP secrets in `dbForProject` (project-scoped, not platform DB).
- **Rate limiting:** Challenge create (10/user), complete (10/challenge).
- **Codes:** Recovery codes hashed before storage; plaintext only on creation/regeneration.
- **Elevation:** Factors appended to server-side session, never client-supplied.

