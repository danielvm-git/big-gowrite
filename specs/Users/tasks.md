# Tasks: Users (Admin)

> Unit: `Users`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-USR-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Users/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 47 endpoint actions.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-USR-02 — Shared user-creation helper
- **Behavior source:** `Modules/Users/Base.php`
- **Description:** `createUser()` — email uniqueness, personalDataCheck, email policy gates, password hashing, auto-target creation.
- **Done when:** All Create actions reuse the helper and pass validation.
- **Priority:** Must · Confidence: 🟢

### T-USR-03 — Plaintext create + 7 hash-import creates
- **Behavior source:** `Http/Users/Create.php`, `Http/Users/{Argon2,Bcrypt,MD5,SHA,PHPass,Scrypt,Scrypt/Modified}/Create.php`
- **Description:** Argon2id default hashing; verbatim import storage with hash name.
- **Done when:** All variants verified.
- **Priority:** Must · Confidence: 🟢

### T-USR-04 — Get / List / Delete user
- **Behavior source:** `Http/Users/{Get,XList,Delete}.php`
- **Description:** Pagination + search; hard-delete cascade (identities/targets inline, DELETE_TYPE_DOCUMENT queued).
- **Done when:** CRUD verified.
- **Priority:** Must · Confidence: 🟢

### T-USR-05 — Profile field updates
- **Behavior source:** `Http/Users/{Status,Labels,Impersonator,Name,Password,Email,Phone}/Update.php`, `Http/Users/Verification/{Update,Phone/Update}.php`, `Http/Users/Prefs/{Get,Update}.php`
- **Description:** Sparse updates with events for status, labels, impersonation, name, password, email, phone, verification (email+phone), prefs.
- **Done when:** All field updates verified.
- **Priority:** Must · Confidence: 🟢

### T-USR-06 — Targets CRUD
- **Behavior source:** `Http/Users/Targets/{Create,Get,XList,Update,Delete}.php`
- **Description:** push/email/sms notification endpoints.
- **Done when:** Target lifecycle verified.
- **Priority:** Must · Confidence: 🟢

### T-USR-07 — Sessions create/list/delete
- **Behavior source:** `Http/Users/Sessions/{Create,XList,Delete,Bulk/Delete}.php`
- **Description:** Server-provider sessions skipping password; bulk delete.
- **Done when:** Session lifecycle verified.
- **Priority:** Must · Confidence: 🟢

### T-USR-08 — Tokens + JWT
- **Behavior source:** `Http/Users/Tokens/Create.php`, `Http/Users/JWTs/Create.php`
- **Description:** Magic-URL/OTP token; HS256 JWT scoped to userId+sessionId.
- **Done when:** Both issuances verified.
- **Priority:** Must · Confidence: 🟢

### T-USR-09 — Memberships/identities
- **Behavior source:** `Http/Users/Memberships/XList.php`, `Http/Users/Identities/{XList,Delete}.php`
- **Description:** Scoped reads + identity removal.
- **Done when:** Verified.
- **Priority:** Must · Confidence: 🟢

### T-USR-10 — MFA management
- **Behavior source:** `Http/Users/MFA/{Update,Factors/XList,RecoveryCodes/Get,RecoveryCodes/Create,RecoveryCodes/Update,Authenticators/Delete}.php`
- **Description:** Toggle, factor list, recovery codes, authenticator removal with enrolled-type check.
- **Done when:** MFA flows verified.
- **Priority:** Must · Confidence: 🟢

### T-USR-11 — listUsers performance path
- **Behavior source:** `Http/Users/XList.php`
- **Description:** Bulk target load by `userInternalId`; `skipFilters()` sub-queries.
- **Done when:** No N+1 in list path.
- **Priority:** Should · Confidence: 🟡

### T-USR-12 — Contract / OpenAPI surface + E2E suite
- **Behavior source:** All `Modules/Users/Http/**`, `tests/e2e/Services/Users/**`
- **Description:** OpenAPI 3.1 for users endpoints; E2E covering CRUD, imports, sessions, targets, MFA, deletion cascade.
- **Done when:** `openapi/users.yaml` matches actions; E2E suite passes.
- **Priority:** Must · Confidence: 🟢
