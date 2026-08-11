# Contracts: Project

> Unit: `Project`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract (representative)

Base: `/v1/project`. Scopes: `project.*`, `keys.*`, `platforms.*`, `policies.*`, `mocks.*`, `templates.*`, `project.oauth2.*`.

| Method | Path | Key params | Success | Notes |
|--------|------|-----------|---------|-------|
| GET | `/project` | — | 200 project | |
| DELETE | `/project` | — | 204 | async cleanup |
| PUT | `/project/labels` | `labels[]` | 200 project | |
| PATCH | `/project/protocols/:protocolId` | `enabled` | 200 | REST/GraphQL/Realtime |
| PATCH | `/project/services/:serviceId` | `enabled` | 200 | |
| PATCH | `/project/auth-methods/:methodId` | `enabled` | 200 | |
| PATCH | `/project/smtp` | `enabled`, `host`, `port`, `username`, `password`, `secure` | 200 | validates connection |
| POST | `/project/smtp/tests` | `sender`, `recipient` | 200 | test email |
| GET | `/project/templates/email` | — | 200 templateList | |
| GET | `/project/templates/email/:templateId` | `locale` | 200 | custom or default |
| PATCH | `/project/templates/email` | `templateId`, `locale`, `subject`, `message` | 200 | SMTP gate |
| POST | `/project/variables` | `key`, `value`, `secret` | 201 variable | |
| PUT | `/project/variables/:variableId` | `key`, `value`, `secret` | 200 | unset-secret blocked |
| POST | `/project/keys` | `name`, `scopes[]`, `expire` | 201 key + secret | |
| POST | `/project/keys/ephemeral` | `name`, `scopes[]`, `expire` (≤3600) | 201 jwt | non-persisted |
| GET | `/project/keys` | — | 200 keyList | |
| POST | `/project/platforms/web` | `name`, `hostname` | 201 platform | per-type |
| POST | `/project/platforms/android` | `name`, `packageName` | 201 platform | |
| POST | `/project/mock-phones` | `number`, `otp` | 201 | |
| GET | `/project/policies` | — | 200 policyList | virtual docs |
| PATCH | `/project/policies/password-history` | `depth` (1–20) | 200 policy | |
| GET | `/project/oauth2` | — | 200 oauth2List | |
| PATCH | `/project/oauth2/{provider}` | `appId`, `secret`, `enabled` | 200 | verify on enable |

## Queue contracts

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForDeletes` | Project delete → worker-deletes | `DeleteMessage`: DELETE_TYPE_DOCUMENT, projectId |
| `publisherForMails` | SMTP test → worker-mails | MailMessage with SMTP test payload |

## Event contract

- `keys.[keyId].{create,update,delete}`
- `variables.[variableId].{create,update}`
- `platforms.[platformId].{create,delete}`
- `mock-phones.[number].{create,update,delete}`
- `oauth2.[providerId].update`
- `authMethod.[methodId].update` / `services.[serviceId].update` / `protocols.[protocolId].update`
- `templates.[templateId].update`
- `projects.[projectId].policies.[policy].update`

## Key storage contract

| Key type | Prefix | Storage |
|----------|--------|---------|
| Standard | `API_KEY_STANDARD_` + 256B hex | Persisted in dbForPlatform |
| Ephemeral | `API_KEY_EPHEMERAL_` + HS256 JWT | In-memory only (≤3600s) |

## Document contract

| Field | Notes |
|-------|-------|
| `oAuthProviders.{providerId}Appid/Secret/Enabled` | Secret always empty on read |
| `auths.mockNumbers[]` | Capped at `APP_LIMIT_COUNT`, no duplicates |
| `auths.*` | Policy virtual docs source |
| Variable `resourceType` | `project` (singular) |
| Key `resourceType` | `projects` (plural) |
