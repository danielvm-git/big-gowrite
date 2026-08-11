# Requirements: Account (MFA)

> Identifier: `ACC-MFA`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Account MFA module delivers multi-factor authentication for end-user accounts. It enables users to register TOTP authenticators, verify them via one-time passwords, manage backup recovery codes, and complete MFA challenges via email, SMS, TOTP, or recovery code. Every challenge completion elevates the current session's factor set, granting access to MFA-protected operations. Broader account features (registration, login, profile, sessions) are handled elsewhere; this module is scoped exclusively to `/v1/account/mfa/**`.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Account` | Module purpose and dependency injections | 🟢 |
| `specs/Account/legacy-mapping.md` | Full endpoint inventory, business rules, events | 🟢 |
| `specs/domain.md#Authentication` | Authentication domain, factor elevation flow | 🟢 |
| `specs/permissions.md#account` | `account` scope grants factor-elevation access | 🟢 |
| `specs/architecture.md#Security` | Argon2id hashing, TOTP secret storage | 🟢 |
| `specs/state-machines.md#Session` | Session factor state transitions | 🟡 |
| `specs/dependencies.md#Account` | Injected services (db, queue, publishers, rate limiter) | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| End user | Secure account with a second factor | Enable MFA, scan TOTP QR code, confirm with OTP |
| End user | Recover access when authenticator is lost | Use a one-time recovery code to complete a challenge |
| End user | Replace authenticator app | Delete old TOTP, register a new one |

## 4. Business rules

1. **BR-01:** MFA can be enabled or disabled on the current user account via a simple boolean toggle. 🟢 — Origin: `PATCH /v1/account/mfa` (Update.php)

2. **BR-02:** Only `totp` is a supported authenticator type; the `:type` param is whitelisted to `['totp']`. 🟢 — Origin: `Authenticators/Create.php` param validation

3. **BR-03:** Creating a TOTP authenticator when one already verified exists returns `USER_AUTHENTICATOR_ALREADY_VERIFIED` (409). If unverified authenticator exists, it is deleted and recreated. 🟢 — Origin: `Authenticators/Create.php`

4. **BR-04:** Confirming a TOTP authenticator requires a live OTP; success sets `verified = true` and appends `totp` to session `factors`. 🟢 — Origin: `Authenticators/Update.php`

5. **BR-05:** Deleting an authenticator purges the user's cached document (`purgeCachedDocument`). 🟢 — Origin: `Authenticators/Delete.php`

6. **BR-06:** MFA challenge factors are `email`, `phone`, `totp`, `recovery_code`. Email requires SMTP + verified email; Phone requires SMS provider + verified phone. 🟢 — Origin: `Challenges/Create.php`

7. **BR-07:** Recovery-code challenges consume the code on success (removed from `mfaRecoveryCodes` array). 🟢 — Origin: `Challenges/Update.php`

8. **BR-08:** Completing any challenge appends the verified factor type to `session.factors` (deduplicated) and sets `session.mfaUpdatedAt`. 🟢 — Origin: `Challenges/Update.php`

9. **BR-09:** Recovery codes are one-time sets; creating when codes exist throws `USER_RECOVERY_CODES_ALREADY_EXISTS` — use PATCH to regenerate. 🟢 — Origin: `RecoveryCodes/Create.php`, `RecoveryCodes/Update.php`

10. **BR-10:** Challenge creation rate-limited to 10 per `(url, userId)`; challenge completion limited to 10 per `(url, challengeId)`. 🟢 — Origin: `Challenges/Create.php`, `Challenges/Update.php`

11. **BR-11:** `mfaProtected` group applies to: get recovery codes, update recovery codes, delete authenticator. 🟢 — Origin: `Services/Http.php` groups


## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Enable or disable MFA on the current user account | Must | Toggle `mfa` boolean; recalculate session factors; emit `users.[userId].update.mfa` | 🟢 |
| RF-02 | List available and verified MFA factors | Must | Return factors array validated against email/phone verification + TOTP registration | 🟢 |
| RF-03 | Create a TOTP authenticator; return secret + provisioning URI | Must | Generate secret, return QR URI; handle verified/unverified guard | 🟢 |
| RF-04 | Confirm a TOTP authenticator by verifying a live OTP | Must | Validate OTP against secret; mark verified; append totp to session factors | 🟢 |
| RF-05 | Delete a TOTP authenticator | Must | Remove authenticator; purge cache; emit delete event | 🟢 |
| RF-06 | Initiate MFA challenge (email, phone, totp, recovery_code) | Must | Dispatch email/SMS OTP or prepare TOTP/recovery-code; emit create event | 🟢 |
| RF-07 | Complete MFA challenge and elevate session factors | Must | Validate challenge code; append factor to session.factors; set mfaUpdatedAt | 🟢 |
| RF-08 | Generate initial set of backup recovery codes | Must | Create 10 hashed one-time codes; emit update event | 🟢 |
| RF-09 | Retrieve current recovery codes | Must | Return plaintext codes from user document | 🟢 |
| RF-10 | Regenerate (replace) recovery codes | Must | Overwrite existing codes; emit update event | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | TOTP secrets stored in project database, protected by session + account scope | `Authenticators/Create.php` in `dbForProject` | 🟢 |
| Security | Rate limiting: 10/(url,userId) challenge create; 10/(url,challengeId) complete | `Challenges/*` `timelimit` calls | 🟢 |
| Performance | Cached user document purged on authenticator change | `Authenticators/Delete.php` → `purgeCachedDocument` | 🟢 |
| Observability | Usage metrics for phone auth (country code) | `usage` injection + phone challenge logging | 🟡 |
| Observability | Events emitted for all state mutations | All actions emit to `queueForEvents` | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Enable MFA with no authenticators
  Given a user without any registered authenticators
  When the user enables MFA via PATCH /v1/account/mfa
  Then MFA is enabled on the account and session factors are recalculated

Scenario: Create TOTP authenticator
  Given a user with MFA enabled and no existing verified authenticator
  When the user creates a TOTP authenticator via POST /v1/account/mfa/authenticators/totp
  Then a TOTP secret and provisioning URI are returned

Scenario: Confirm TOTP authenticator with valid OTP
  Given a user with an unverified TOTP authenticator
  When the user confirms with a valid OTP via PUT /v1/account/mfa/authenticators/totp
  Then the authenticator is marked verified and "totp" is added to session factors

Scenario: Duplicate authenticator creation fails
  Given a user with an already-verified TOTP authenticator
  When the user tries to create another via POST /v1/account/mfa/authenticators/totp
  Then a 409 USER_AUTHENTICATOR_ALREADY_VERIFIED error is returned

Scenario: Complete MFA challenge with recovery code
  Given a user with existing recovery codes and an active challenge
  When the user submits a valid recovery code via PUT /v1/account/mfa/challenges
  Then the session is elevated with recovery_code factor and the code is consumed

Scenario: Rate limit blocks excessive challenges
  Given a user who has initiated 10 challenges from the same URL
  When the user attempts an 11th within the time window
  Then the request is rejected with a rate limit error
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01 (Enable/disable MFA) | Must | Core toggle; all flows depend on it |
| RF-02 (List factors) | Must | Clients need to know which factors to prompt |
| RF-03 (Create authenticator) | Must | Only supported authenticator type |
| RF-04 (Confirm authenticator) | Must | TOTP is useless without verification |
| RF-05 (Delete authenticator) | Must | User needs to replace lost authenticators |
| RF-06 (Initiate challenge) | Must | Core MFA flow |
| RF-07 (Complete challenge) | Must | Session elevation is the security benefit |
| RF-08 (Generate recovery codes) | Must | Backup for lost authenticator |
| RF-09 (Get recovery codes) | Should | Convenience; client can store locally |
| RF-10 (Regenerate recovery codes) | Must | Codes may be lost or compromised |
| Rate limiting (RNF) | Must | Brute-force protection |
| Cached doc purge (RNF) | Should | Prevents stale reads |

## 9. Clarifications

> No clarification session recorded. Run `/reversa-clarify` when `[GAP]` pending.

## 10. Gaps

- 🔴 [GAP] Broader account features (registration, login, profile, sessions, verifications) are outside this module — where are they registered?
- 🔴 [GAP] SMS provider abstraction (`_APP_SMS_PROVIDER`) referenced but code not visible in this module.
- 🟡 [INFERRED] Multiple SDK `Method` entries per action for deprecation (camelCase `Mfa` vs `MFA`). Full table not exhaustively verified.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |

