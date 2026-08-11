# Contracts: Account (MFA)

> Unit: `Account`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base path: `/v1`. All endpoints require `account` scope (valid session token).

| Method | Path | Request body (JSON) | Success | Errors |
|--------|------|---------------------|---------|--------|
| PATCH | `/account/mfa` | `{ mfa: boolean }` | 200 user | 400 invalid bool |
| GET | `/account/mfa/factors` | — | 200 `{ factors: string[] }` | 401 |
| POST | `/account/mfa/authenticators/:type` | — | 201 `{ type, secret, uri }` | 409 already verified |
| PUT | `/account/mfa/authenticators/:type` | `{ otp: string }` | 200 authenticator | 401 invalid otp, 404 not found |
| DELETE | `/account/mfa/authenticators/:type` | — | 204 | 404, 401 mfaProtected |
| POST | `/account/mfa/challenges` | `{ factor: string }` | 201 `{ id, token }` | 400 unavailable factor, 429 |
| PUT | `/account/mfa/challenges` | `{ challengeId, otp }` | 200 session | 404, 400, 429 |
| POST | `/account/mfa/recovery-codes` | — | 201 `{ recoveryCodes: string[] }` | 409 already exists |
| GET | `/account/mfa/recovery-codes` | — | 200 `{ recoveryCodes: string[] }` | 401 mfaProtected |
| PATCH | `/account/mfa/recovery-codes` | — | 200 `{ recoveryCodes: string[] }` | 401 mfaProtected |

## Path parameters

| Param | Type | Allowed values | Notes |
|-------|------|----------------|-------|
| `:type` | string | `totp` | Whitelisted; other types rejected |

## Request body schemas

### PATCH /account/mfa
```json
{ "mfa": true }
```

### PUT /account/mfa/authenticators/totp
```json
{ "otp": "123456" }
```

### POST /account/mfa/challenges
```json
{ "factor": "email" }
```
`factor` ∈ `email | phone | totp | recovery_code`.

### PUT /account/mfa/challenges
```json
{ "challengeId": "challenge_abc", "otp": "123456" }
```

## Response schemas

### Create TOTP authenticator (201)
```json
{
  "$id": "auth_abc",
  "$createdAt": "...",
  "$updatedAt": "...",
  "type": "totp",
  "secret": "JBSWY3DPEHPK3PXP",
  "uri": "otpauth://totp/Appwrite:user@example.com?secret=...&issuer=Appwrite"
}
```

### Initiate challenge (201)
```json
{
  "$id": "challenge_abc",
  "$createdAt": "...",
  "type": "email",
  "token": "f4f0a2c1..."
}
```

### Session after completed challenge (200)
```json
{
  "$id": "session_abc",
  "factors": ["email", "totp"],
  "mfaUpdatedAt": "..."
}
```

## Queue / messaging contracts

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForMails` | Account → worker-mails | Email OTP: recipient, template id, locale, challenge code, branding context |
| `publisherForMessaging` | Account → worker-messaging | SMS OTP: phone number, message template, provider route, country code |

## Event contract

| Event | Schema |
|-------|--------|
| `users.[userId].update.mfa` | `{ userId, events: [...], timestamp }` |
| `users.[userId].delete.mfa` | `{ userId, events: [...], timestamp }` |
| `users.[userId].challenges.[challengeId].create` | `{ userId, challengeId, timestamp }` |
| `users.[userId].sessions.[sessionId].create` | `{ userId, sessionId, factors, timestamp }` |
