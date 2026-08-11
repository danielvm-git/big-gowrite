# Contracts: Notifications

> Unit: `Notifications`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/notifications`.

| Method | Path | Params | Success | Errors |
|--------|------|--------|---------|--------|
| GET | `/notifications` | `limit`, `cursor` | 200 notificationList | 401 USER_UNAUTHORIZED (non-console) |
| PATCH | `/notifications/:notificationId` | `read` (bool) | 200 notification | 404, 401 |
| GET | `/notifications/logos/appwrite` | `theme` (system/light/dark), `jwt` (optional) | 200 image/svg+xml | — |

## Tracking JWT contract

| Claim | Required | Notes |
|-------|----------|-------|
| `messageId` | yes | Target notification |
| `recipientHash` | yes | Recipient binding |
| `channel` | yes | Channel identifier |
| `projectId` | yes | Project scope |
| `purpose` | yes | Must equal `notification_track` |

Algorithm: `HS256`, key `_APP_NOTIFICATIONS_TRACKING_SECRET`.

## Response headers (logo)

- `Content-Type: image/svg+xml`
- `Cache-Control: no-store`

## Queue contract

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForAudits` | Tracking → audits | `AuditMessage`: event `notification.view`, IP/user-agent/hostname |

## Document model (`notifications` in dbForPlatform)

| Field | Type | Notes |
|-------|------|-------|
| `channel` | string | `console` |
| `resourceType` | string | `users` |
| `resourceId` | string | Target user id |
| `read` | bool | Mark-read state |
| `firstSeen` / `lastSeen` | datetime | Set by tracking |
| `messageId` | string | Email/batch message id |
