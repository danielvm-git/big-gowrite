# Design: Webhooks

> Unit: `Webhooks`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Webhooks module is an **HTTP-only management service** for project webhook configurations. It reads/writes the `webhooks` collection in the **platform database** (scoped by `projectInternalId`) and emits lifecycle events via `queueForEvents`. The `console` project is blocked at the Init hook.

```
Client (admin / API key)
  │ scopes: webhooks.*
  ▼
app/http.php → Webhooks\Http\** (6 actions + Init hook)
  │
  ├── Init.php — blocks console project → GENERAL_ACCESS_FORBIDDEN
  ├── dbForPlatform (webhooks collection, projects cache purge)
  ├── queueForEvents (webhooks.[webhookId].* events)
  └── authorization (skip ACL on platform records)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/webhooks` | `createWebhook` | webhooks.write |
| GET | `/v1/webhooks` | `listWebhooks` | webhooks.read |
| GET | `/v1/webhooks/:webhookId` | `getWebhook` | webhooks.read |
| PUT | `/v1/webhooks/:webhookId` | `updateWebhook` | webhooks.write |
| DELETE | `/v1/webhooks/:webhookId` | `deleteWebhook` | webhooks.write |
| PATCH | `/v1/webhooks/:webhookId/secret` | `updateWebhookSecret` | webhooks.write |

Legacy aliases (via `httpAlias`): `/v1/projects/:projectId/webhooks[...]`.

## 3. Data flows — key actions

### 3.1 Create webhook

```
1. Init hook: block console project
2. Validate URL via PublicDomain validator
3. Secret = provided or bin2hex(random_bytes(64))
4. Store: httpUser/httpPass (from authUsername/authPassword), security (from tls)
5. Scope by projectInternalId
6. Purge cached projects document
7. Emit webhooks.[webhookId].create (returns signatureKey once)
```

### 3.2 Get / List

```
1. Scope query by projectInternalId
2. Response model filter strips signatureKey
```

### 3.3 Update (full replace)

```
1. Replace configuration
2. enabled=true → reset attempts to 0
3. Purge cached projects document
4. Emit webhooks.[webhookId].update
```

### 3.4 Rotate secret

```
1. Secret = provided or auto-generated 128-char hex
2. Store; return signatureKey once
3. Purge cached projects document
```

### 3.5 Delete

```
1. Delete document
2. Purge cached projects document
3. Emit webhooks.[webhookId].delete
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForPlatform` | webhooks collection + projects cache purge |
| `queueForEvents` | Lifecycle events |
| `project` | projectInternalId scoping |
| `authorization` | Skip ACL on platform records |
| `response` | HTTP response builder |

## 5. Events emitted

- `webhooks.[webhookId].create`
- `webhooks.[webhookId].update`
- `webhooks.[webhookId].delete`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `GENERAL_ACCESS_FORBIDDEN` | 403 | Console project (Init hook) |
| URL validation | 400 | Non-public or non-HTTP(S) URL |
| `WEBHOOK_ALREADY_EXISTS` | 409 | Duplicate webhook ID |
| Webhook not found | 404 | Unknown `:webhookId` |

## 7. Security

- Console project blocked at Init hook.
- `signatureKey` only returned on create/secret-rotation.
- Basic-auth credentials stored as `httpUser`/`httpPass`; `tls` stored as `security` (remapped on read).
