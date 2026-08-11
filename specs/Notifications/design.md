# Design: Notifications

> Unit: `Notifications`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Notifications module is a **thin HTTP service** over the platform `notifications` collection. It serves console-user notifications and a public tracking pixel backed by signed JWTs. Notification creation happens outside this module.

```
Client (console project)
  │ account / public scopes
  ▼
app/http.php → Notifications\Http\** (3 actions)
  │
  ├── dbForPlatform (notifications collection)
  ├── publisherForAudits (notification.view)
  └── JWT decode (HS256, _APP_NOTIFICATIONS_TRACKING_SECRET)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| GET | `/v1/notifications` | `listNotifications` | account |
| PATCH | `/v1/notifications/:notificationId` | `updateNotification` | account |
| GET | `/v1/notifications/logos/appwrite` | `getNotificationLogo` | public |

## 3. Data flows — key actions

### 3.1 List notifications

```
1. Guard: project id == console, else USER_UNAUTHORIZED
2. Filters: channel=console, resourceType=users, resourceId=$user->getId()
3. Paginate; return list
```

### 3.2 Update notification

```
1. Guard console-only
2. getDocument → empty → privileged skip query (404 vs 401)
3. Sparse update: { read: bool }
```

### 3.3 Logo + email open tracking

```
1. Serve SVG by theme (system/light/dark)
2. If jwt param present:
   a. Decode HS256 with _APP_NOTIFICATIONS_TRACKING_SECRET
   b. Validate messageId, recipientHash, channel, projectId, purpose=notification_track
   c. Update notification: read=true, firstSeen (first only), lastSeen
   d. Enqueue AuditMessage (notification.view) via publisherForAudits
3. Any tracking error → swallow, still return SVG
4. Cache-Control: no-store
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForPlatform` | notifications collection |
| `response` | HTTP response |
| `request` | IP/user-agent/hostname for audits |
| `authorization` | Permission checks + skip query |
| `project` | Console-only guard |
| `user` | Current user |
| `publisherForAudits` | notification.view audit events |

## 5. Events / audit

| Type | Detail |
|------|--------|
| Audit | `notification.view` (tracking pixel hit) |

## 6. Security

- Console-only guard on list/update.
- Public logo endpoint; tracking JWT validated with dedicated secret.
- Tracking failures never surface; no-store on pixel.
