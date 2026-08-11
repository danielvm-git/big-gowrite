# Contracts: Webhooks

> Unit: `Webhooks`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/webhooks`. Scopes: `webhooks.*`. The `console` project is blocked (Init hook → `GENERAL_ACCESS_FORBIDDEN`).

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `` | `webhookId`, `name`, `url`, `events: string[]`, `enabled`, `security`, `httpUser`, `httpPass` | 201 webhook (+ signatureKey) | GENERAL_ARGUMENT_INVALID, URL invalid, WEBHOOK_ALREADY_EXISTS |
| GET | `` | `queries` | 200 webhookList (no signatureKey) | — |
| GET | `/:webhookId` | — | 200 webhook (no signatureKey) | 404 |
| PUT | `/:webhookId` | same as create | 200 webhook | 404, URL invalid |
| DELETE | `/:webhookId` | — | 204 | 404 |
| PATCH | `/:webhookId/secret` | `secret` (optional) | 200 webhook (+ signatureKey once) | 404 |

Legacy aliases: `/v1/projects/:projectId/webhooks[...]` (via `httpAlias`).

## Field mapping contract

| API field | Stored field | Notes |
|-----------|--------------|-------|
| `authUsername` | `httpUser` | Mapped back in response model |
| `authPassword` | `httpPass` | Mapped back in response model |
| `tls` | `security` | Remapped in response model |
| `signatureKey` | `signatureKey` | Stripped on get/list; returned only on create + secret rotation |

## Secret contract

| Aspect | Detail |
|--------|--------|
| Auto-generation | `bin2hex(random_bytes(64))` → 128-char hex |
| Exposure | Only on `createWebhook` and `updateWebhookSecret` |

## Lifecycle contract

| Aspect | Detail |
|--------|--------|
| Re-enable | `enabled=true` on update resets `attempts` to 0 |
| Scoping | `projectInternalId` (indexed, not `projectId`) |
| Cache | `purgeCachedDocument('projects', ...)` after every write |

## Event contract

- `webhooks.[webhookId].create`
- `webhooks.[webhookId].update`
- `webhooks.[webhookId].delete`
