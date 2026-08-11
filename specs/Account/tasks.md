# Tasks: Account (MFA)

> Unit: `Account`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP
> Each task cites its legacy source file. Done criteria always present.

## Task list

### T-ACC-01 — Register Account module boilerplate
- **Behavior source:** `Modules/Account/Module.php`
- **Description:** Create `Module.php` that registers the Account HTTP service on the platform App instance.
- **Done when:** `Module.php` exists and registers `Http::class`; module loads without error on bootstrap.
- **Priority:** Must · Confidence: 🟢

### T-ACC-02 — Register the 10 MFA HTTP routes
- **Behavior source:** `Modules/Account/Services/Http.php`
- **Description:** Register all 10 `/v1/account/mfa/**` routes with correct verbs, paths, action IDs, groups (`api`, `account`, `mfaProtected`) and `account` scope labels.
- **Done when:** Route table matches legacy mapping: updateMFA → PATCH, listMFAFactors → GET, createMFAAuthenticator → POST, updateMFAAuthenticator → PUT, deleteMFAAuthenticator → DELETE, createMFAChallenge → POST, updateMFAChallenge → PUT, createMFARecoveryCodes → POST, getMFARecoveryCodes → GET, updateMFARecoveryCodes → PATCH.
- **Priority:** Must · Confidence: 🟢

### T-ACC-03 — Enable/disable MFA
- **Behavior source:** `Modules/Account/Http/Account/MFA/Update.php`
- **Description:** Implement `updateMFA` to accept `{ mfa: boolean }`, persist to user document, recalculate session factors (TOTP verified, email verified, phone verified), set `mfaUpdatedAt`, emit `users.[userId].update.mfa`.
- **Done when:** Enabling/disabling MFA mutates user doc, recalcs factors, emits event; E2E test green.
- **Priority:** Must · Confidence: 🟢

### T-ACC-04 — List MFA factors
- **Behavior source:** `Modules/Account/Http/Account/MFA/Factors/XList.php`
- **Description:** Implement `listMFAFactors` returning factors derived from TOTP registration, email verification, phone verification state.
- **Done when:** Returns accurate factor list matching user state.
- **Priority:** Must · Confidence: 🟢

### T-ACC-05 — Create TOTP authenticator
- **Behavior source:** `Modules/Account/Http/Account/MFA/Authenticators/Create.php`
- **Description:** Implement `createMFAAuthenticator` supporting only `totp`. Guard 409 `USER_AUTHENTICATOR_ALREADY_VERIFIED` if verified exists; delete unverified stale authenticator then recreate. Generate secret + provisioning URI; persist to dbForProject; emit update event.
- **Done when:** Secret + `otpauth://` URI returned; duplicate guard verified by E2E.
- **Priority:** Must · Confidence: 🟢

### T-ACC-06 — Confirm TOTP authenticator
- **Behavior source:** `Modules/Account/Http/Account/MFA/Authenticators/Update.php`
- **Description:** Implement `updateMFAAuthenticator` verifying submitted OTP against stored secret; on success set `verified=true`, append `totp` to session factors.
- **Done when:** Valid OTP verifies authenticator; invalid OTP returns 401.
- **Priority:** Must · Confidence: 🟢


### T-ACC-08 — Initiate MFA challenge
- **Behavior source:** `Modules/Account/Http/Account/MFA/Challenges/Create.php`
- **Description:** Implement `createMFAChallenge` with rate limit 10/(url,userId). Validate factor availability (email: SMTP + verified email; phone: SMS provider + verified phone; totp: verified authenticator; recovery_code: codes exist). Generate token+code via proofForToken/proofForCode; publish to mails/messaging publishers; create challenge doc; emit create event.
- **Done when:** All four factors initiate correctly; email/SMS dispatched; rate limit enforced at 10.
- **Priority:** Must · Confidence: 🟢

### T-ACC-09 — Complete MFA challenge
- **Behavior source:** `Modules/Account/Http/Account/MFA/Challenges/Update.php`
- **Description:** Implement `updateMFAChallenge` with rate limit 10/(url,challengeId). Validate challenge and code; consume recovery code if used; append factor to session.factors; set `mfaUpdatedAt`; emit session event.
- **Done when:** Session elevated with correct factor; recovery code consumed; rate limit enforced.
- **Priority:** Must · Confidence: 🟢

### T-ACC-10 — Generate initial recovery codes
- **Behavior source:** `Modules/Account/Http/Account/MFA/RecoveryCodes/Create.php`
- **Description:** Implement `createMFARecoveryCodes` generating ~10 one-time backup codes, hashed before storage; guard 409 `USER_RECOVERY_CODES_ALREADY_EXISTS` when codes exist.
- **Done when:** Initial codes created and hashed; duplicate attempt returns 409.
- **Priority:** Must · Confidence: 🟢

### T-ACC-11 — Get recovery codes
- **Behavior source:** `Modules/Account/Http/Account/MFA/RecoveryCodes/Get.php`
- **Description:** Implement `getMFARecoveryCodes` returning existing codes (plaintext) from user doc. Requires `mfaProtected`.
- **Done when:** Returns codes for users with codes; empty/404 for users without.
- **Priority:** Must · Confidence: 🟢

### T-ACC-12 — Regenerate recovery codes
- **Behavior source:** `Modules/Account/Http/Account/MFA/RecoveryCodes/Update.php`
- **Description:** Implement `updateMFARecoveryCodes` replacing current code set with a fresh set. Requires `mfaProtected`.
- **Done when:** Old codes invalidated and new set returned.
- **Priority:** Must · Confidence: 🟢

### T-ACC-13 — Contract / OpenAPI surface
- **Behavior source:** `Modules/Account/Http/Account/MFA/**` + `Services/Http.php`
- **Description:** Generate OpenAPI 3.1 definitions + SDK Method labels for the 10 endpoints, honoring deprecated camelCase `Mfa` variants replaced by uppercase `MFA`.
- **Done when:** `openapi/account-mfa.yaml` conforms to action `getName()` and SDK labels `create/get/list/update/delete`.
- **Priority:** Should · Confidence: 🟡

### T-ACC-14 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Account/**`
- **Description:** Cover MFA lifecycle: enable, list factors, create TOTP, confirm valid/invalid OTP, delete, challenge via email/phone/totp/recovery, complete challenge, recovery codes create/get/regenerate, mfaProtected enforcement.
- **Done when:** `docker compose exec appwrite test tests/e2e/Services/Account --filter=MFA` passes.
- **Priority:** Must · Confidence: 🟢
### T-ACC-07 — Delete TOTP authenticator
- **Behavior source:** `Modules/Account/Http/Account/MFA/Authenticators/Delete.php`
- **Description:** Implement `deleteMFAAuthenticator`: remove authenticator doc, purge user cached document, emit `users.[userId].delete.mfa`. Requires `mfaProtected` group.
- **Done when:** Deletion removes authenticator and purges cache.
- **Priority:** Must · Confidence: 🟢
