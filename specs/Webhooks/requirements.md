# Requirements: Webhooks

> Identifier: `WHK-OUTBOUND`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Webhooks module provides full CRUD management for project-level webhook configurations. Webhooks are HTTP callbacks Appwrite fires when platform events occur (user creation, document update, etc.). The module lets developers register, inspect, modify, delete, and regenerate the signing secret for project webhooks. All endpoints are project-scoped and inaccessible from the `console` project.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Webhooks` | Module purpose | 🟢 |
| `specs/Webhooks/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/architecture.md#Events` | Event → webhook delivery | 🟢 |
| `specs/permissions.md` | Project scoping | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App developer | Integrate externally | Register a webhook URL for project events |
| App developer | Rotate credentials | Regenerate the signing key after a leak |
| App developer | Pause delivery | Disable a webhook; re-enable resets attempts |

## 4. Business rules

1. **BR-01:** The `console` project is blocked at the Init hook — all webhook endpoints return `GENERAL_ACCESS_FORBIDDEN`. 🟢
2. **BR-02:** Webhook URL must be a valid HTTP/HTTPS URL pointing to a public domain (validated by `PublicDomain` validator). 🟢
3. **BR-03:** If `secret` is not provided on creation or rotation, a 128-character hex string is auto-generated via `bin2hex(random_bytes(64))`. 🟢
4. **BR-04:** The `signatureKey` field is stripped from responses in `getWebhook` and `listWebhooks`; it is only returned on `createWebhook` and `updateWebhookSecret`. 🟢
5. **BR-05:** Webhooks are scoped to a project via `projectInternalId` (not `projectId`) for indexed-query performance. 🟢
6. **BR-06:** Re-enabling a webhook (`enabled=true` on update) resets the `attempts` counter to 0. 🟢
7. **BR-07:** After any write operation, the project's cached document is purged (`purgeCachedDocument('projects', ...)`). 🟢
8. **BR-08:** `authUsername`/`authPassword` are stored as `httpUser`/`httpPass` internally and mapped back in the response model filter. 🟢
9. **BR-09:** The `tls` field is stored as `security` internally and remapped in the response model filter. 🟢
10. **BR-10:** Duplicate webhook IDs throw `WEBHOOK_ALREADY_EXISTS`. 🟡
11. **BR-11:** Legacy path aliases (`/v1/projects/:projectId/webhooks`) are registered via `httpAlias` for backwards compatibility. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Create webhook | Must | URL validation, secret auto-gen, project scoping | 🟢 |
| RF-02 | List webhooks | Must | Signature key excluded | 🟢 |
| RF-03 | Get webhook | Must | Signature key excluded | 🟢 |
| RF-04 | Update webhook (replace) | Must | Full replace; re-enable resets attempts | 🟢 |
| RF-05 | Delete webhook | Must | Removes configuration | 🟢 |
| RF-06 | Rotate secret | Must | Auto-gen or explicit; returned once | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Console project blocked | Init hook, GENERAL_ACCESS_FORBIDDEN | 🟢 |
| Security | Key not exposed on reads | signatureKey stripped | 🟢 |
| Reliability | Cache consistency | purgeCachedDocument on writes | 🟢 |
| Performance | Indexed scoping | projectInternalId queries | 🟢 |
| Compat | Legacy path aliases | httpAlias | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Console project access
  Given the console project
  When any webhook endpoint is called
  Then GENERAL_ACCESS_FORBIDDEN is returned

Scenario: Create webhook without secret
  When POST /v1/webhooks is called without a secret
  Then a 128-character hex signatureKey is auto-generated and returned once

Scenario: Get webhook hides secret
  Given a webhook
  When GET /v1/webhooks/:webhookId is called
  Then the response contains no signatureKey field

Scenario: Re-enable resets attempts
  Given a disabled webhook with attempts > 0
  When PUT updates enabled=true
  Then attempts is reset to 0

Scenario: Invalid URL
  Given a non-public or non-HTTP URL
  When creating a webhook
  Then URL validation fails
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..06 (core) | Must | Webhook lifecycle |
| Legacy aliases | Should | Backwards compatibility |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] `WEBHOOK_ALREADY_EXISTS` behavior on concurrent duplicate creation.
- 🔴 [GAP] Exact delivery/retry semantics of the webhook executor (attempts counter lifecycle).

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
