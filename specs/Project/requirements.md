# Requirements: Project

> Identifier: `PRJ-CONFIG`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Project module manages configuration and lifecycle of a single Appwrite project: get/delete project, API keys (standard + ephemeral JWT), project-level env variables, client platform registrations (web, Apple, Android, Windows, Linux), OAuth2 provider credentials (38 providers), email templates, SMTP config, security/auth policies, auth method toggles, protocol/service toggles, mock phones for OTP testing, and project labels. All endpoints operate in the current project context (Init hook) and write to `dbForPlatform` (keys, platforms, OAuth2, SMTP) or `dbForProject` (variables).

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Project` | Module purpose | 🟢 |
| `specs/Project/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/domain.md#Project` | Project document model | 🟢 |
| `specs/security.md#AuthPolicies` | Policies, keys, OAuth2 | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Console admin | Configure project | SMTP, OAuth2, policies, platforms |
| Console admin | Issue API keys | Standard + ephemeral keys |
| App developer | Set env | Project variables (secret flag) |
| Console admin | OTP testing | Mock phones |

## 4. Business rules

1. **BR-01:** Init hook blocks console meta-project and missing projects. 🟢
2. **BR-02:** Standard keys store secret as `API_KEY_STANDARD_` + 256 bytes hex random; ephemeral keys store signed HS256 JWT as `API_KEY_EPHEMERAL_` prefix. 🟢
3. **BR-03:** Ephemeral keys not persisted; in-memory Document; only JWT secret returned; max 3600s. 🟢
4. **BR-04:** Key ownership: `resourceType == 'projects'` + `resourceInternalId == project.getSequence()`. 🟢
5. **BR-05:** After key writes, project document cache purged. 🟢
6. **BR-06:** Project variable create/update marks all functions/sites `live=false` (rebuild signal). 🟢
7. **BR-07:** Project variables use `resourceType = 'project'` (singular); keys use `'projects'`; platforms use `projectInternalId`. 🟢
8. **BR-08:** SMTP update with `enabled=true` attempts real connection (`PHPMailer::SmtpConnect`); throws on failure; `enabled=null` + success → auto-enable. 🟢
9. **BR-09:** Email templates only updatable when custom SMTP enabled. 🟢
10. **BR-10:** OAuth2 creds stored in `oAuthProviders` map: `{providerId}Appid`, `{providerId}Secret`, `{providerId}Enabled`; secret always empty on reads. 🟢
11. **BR-11:** Enabling OAuth2 provider triggers `verifyCredentials()` if present; throws on failure. 🟢
12. **BR-12:** Mock phones stored in `project.auths.mockNumbers`, capped at `APP_LIMIT_COUNT`, duplicates rejected. 🟢
13. **BR-13:** Policies are virtual documents synthesized from `project.auths.*`. 🟢
14. **BR-14:** Platform backward-compat: legacy `type` field mapping in Web platform Create. 🟡
15. **BR-15:** Secret variables cannot be downgraded (`VARIABLE_CANNOT_UNSET_SECRET`). 🟢
16. **BR-16:** Project deletion: delete doc from dbForPlatform + enqueue `DELETE_TYPE_DOCUMENT`. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Get/delete project | Must | Init guard; delete enqueues cleanup | 🟢 |
| RF-02 | Project labels | Must | Replace array | 🟢 |
| RF-03 | Protocol/service/auth-method toggles | Must | Enable/disable by id | 🟢 |
| RF-04 | SMTP config + test email | Must | Connection validated; test mail via publisher | 🟢 |
| RF-05 | Email templates (get/list/update) | Must | Locale-specific; SMTP gate | 🟢 |
| RF-06 | Project variables CRUD | Must | secret flag; live=false propagation | 🟢 |
| RF-07 | API keys CRUD + ephemeral | Must | Ownership; JWT ephemeral ≤3600s; cache purge | 🟢 |
| RF-08 | Platforms CRUD (5 types) | Must | Per-type paths + legacy mapping | 🟢 |
| RF-09 | Mock phones CRUD | Must | Cap + dedupe | 🟢 |
| RF-10 | Policies list/get/update (12 policies) | Must | Virtual docs from auths.* | 🟢 |
| RF-11 | OAuth2 list/get/update (38 providers) | Must | Credential verify; secret masked | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Ephemeral keys non-persisted | in-memory JWT | 🟢 |
| Security | OAuth2 secrets never returned | always empty on read | 🟢 |
| Security | SMTP connection validation | PHPMailer::SmtpConnect | 🟢 |
| Consistency | Cache purge after key writes | purgeCachedDocument('projects') | 🟢 |
| Performance | Rebuild signal via live=false | functions/sites marking | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Ephemeral key creation
  Given an authenticated console admin
  When POST /v1/project/keys/ephemeral with expire > 3600 is called
  Then the request is rejected

Scenario: SMTP validation failure
  Given invalid SMTP credentials with enabled=true
  When PATCH /v1/project/smtp is called
  Then an error is thrown and config is not saved

Scenario: OAuth2 secret masked
  Given a configured provider
  When GET /v1/project/oauth2/:provider is called
  Then the secret field is returned empty

Scenario: Secret variable downgrade
  Given a secret project variable
  When updating with secret=false is attempted
  Then VARIABLE_CANNOT_UNSET_SECRET is returned
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..08 (core) | Must | Project config |
| RF-09 (mock phones) | Should | OTP testing |
| RF-10 (policies) | Must | Security policies |
| RF-11 (OAuth2) | Must | Provider auth |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Legacy platform `type` mapping completeness.
- 🔴 [GAP] Full policy virtual-document synthesis logic.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |

