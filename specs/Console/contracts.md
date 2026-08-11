# Contracts: Console

> Unit: `Console`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

| Method | Path | Params | Success | Errors |
|--------|------|--------|---------|--------|
| GET | `/v1/console/variables` | — | 200 `{ variables, capabilities, features }` | 401 |
| GET | `/v1/console/templates/email/:templateId` | `locale` | 200 rendered HTML template | 400 unknown type/locale |
| GET | `/v1/console/oauth2-providers` | — | 200 `{ providers: [...] }` | — |
| GET | `/v1/console/scopes/project` | — | 200 `{ scopes: [...] }` | — |
| GET | `/v1/console/scopes/organization` | — | 200 `{ scopes: [...] }` | — |
| POST | `/v1/console/assistant` | `{ message }` | 200 text/event-stream (SSE) | 429 rate limit, 502 assistant down |
| GET | `/v1/console/resources` | `type`, `value` | 200 `{ available: bool }` | 400 invalid |
| GET | `/` | — | 301 → /console | — |
| GET | `/auth/*` | — | 301 → /console/auth/* | — |
| GET | `/invite` | — | 301 → /console/invite | — |
| GET | `/login` | — | 301 → /console/login | — |
| GET | `/mfa/*` | — | 301 → /console/mfa/* | — |
| GET | `/card` | — | 301 → /console/card | — |
| GET | `/recover` | — | 301 → /console/recover | — |
| GET | `/register` | — | 301 → /console/register | — |

## External service contract (assistant)

| Endpoint | Method | Payload |
|----------|--------|---------|
| `http://appwrite-assistant:3003/v1/models/assistant/prompt` | POST | `{ message }` → SSE stream of assistant tokens |

## Configuration contract

| Env var | Effect |
|---------|--------|
| `_APP_CUSTOM_DOMAIN_DENY_LIST` | Domains blocked from custom resource assignment |
| `appwrite-assistant` (host) | Assistant microservice URL (default port 3003) |

## Security contract

- `console` group: only under project id `console`.
- `web` group headers: `X-Frame-Options: SAMEORIGIN`, `X-XSS-Protection`, `X-UA-Compatible`.
- Assistant abuse: `abuse-limit: 15`, `abuse-key: userId:{userId}`.
