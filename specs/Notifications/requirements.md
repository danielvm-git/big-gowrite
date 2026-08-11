# Requirements: Notifications

> Identifier: `NTF-CONSOLE`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Notifications module handles in-app/console notifications for Appwrite platform users. Notifications live in the platform-level `notifications` collection (not per-project) and currently target console-channel users (`NOTIFICATION_TYPE_CONSOLE`, `RESOURCE_TYPE_USERS`). It provides list and mark-read endpoints plus a public tracking-pixel endpoint (`/v1/notifications/logos/appwrite`) that serves the Appwrite SVG logo and records email open tracking by decoding a signed JWT.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Notifications` | Module purpose | 🟢 |
| `specs/Notifications/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/security.md#JWT` | Signed tracking tokens | 🟢 |
| `specs/audit.md` | notification.view audit | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Console user | Review notifications | GET /v1/notifications |
| Console user | Manage read state | PATCH notification read=true |
| Email client | Open tracked email | Load logo pixel → open recorded |

## 4. Business rules

1. **BR-01:** List/update are console-only: non-console project → `USER_UNAUTHORIZED`. 🟢
2. **BR-02:** List always appends server-side filters: `channel = NOTIFICATION_TYPE_CONSOLE`, `resourceType = RESOURCE_TYPE_USERS`, `resourceId = $user->getId()`. 🟢
3. **BR-03:** Update only modifies `read` boolean via sparse document update. 🟢
4. **BR-04:** Update auth double-checked: empty getDocument → privileged skip query distinguishes not-found vs unauthorized. 🟢
5. **BR-05:** Logo endpoint scope `public`. 🟢
6. **BR-06:** Email tracking: `jwt` param decoded HS256 with `_APP_NOTIFICATIONS_TRACKING_SECRET`; payload must contain `messageId`, `recipientHash`, `channel`, `projectId`, `purpose=notification_track`. 🟢
7. **BR-07:** Tracking sets `read=true`, records `firstSeen` (first open only) + `lastSeen`. 🟢
8. **BR-08:** Tracking errors silently swallowed; SVG always returned. 🟢
9. **BR-09:** SVG supports `theme` (`system`/`light`/`dark`); system uses CSS media query. 🟢
10. **BR-10:** Logo response `Cache-Control: no-store` — each open fires a fresh request. 🟢
11. **BR-11:** Notifications created by external process (messaging/email workers), not this module. 🟡

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | List console notifications | Must | Server-side filters applied; console-only | 🟢 |
| RF-02 | Update notification read state | Must | Sparse read update; 404 vs 401 distinction | 🟢 |
| RF-03 | Serve Appwrite logo (public) | Must | SVG by theme; no-store cache | 🟢 |
| RF-04 | Email open tracking | Must | JWT validated; read/firstSeen/lastSeen updated; errors swallowed | 🟢 |
| RF-05 | Audit notification views | Must | AuditMessage `notification.view` | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Console-only gating | project id check + USER_UNAUTHORIZED | 🟢 |
| Security | Signed tracking JWT | HS256 with tracking secret | 🟢 |
| Privacy | No-store on tracking pixel | Cache-Control: no-store | 🟢 |
| Reliability | Tracking never breaks logo serving | errors silently swallowed | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Non-console project blocked
  Given a request targeting a non-console project
  When GET /v1/notifications is called
  Then USER_UNAUTHORIZED is returned

Scenario: Valid tracking JWT
  Given a JWT with messageId, recipientHash, channel, projectId, purpose=notification_track
  When GET /v1/notifications/logos/appwrite?jwt=... is called
  Then the notification is marked read with firstSeen/lastSeen and the SVG is returned

Scenario: Invalid tracking JWT
  Given a malformed JWT
  When GET /v1/notifications/logos/appwrite?jwt=... is called
  Then the SVG is still returned and no error surfaces
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01 (List) | Must | Core UX |
| RF-02 (Update) | Must | Read state |
| RF-03 (Logo) | Must | Public asset |
| RF-04 (Tracking) | Should | Email analytics |
| RF-05 (Audit) | Should | Observability |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Notification creation pipeline (external workers) not visible in module.
- 🔴 [GAP] Exact SVG asset source not captured.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
