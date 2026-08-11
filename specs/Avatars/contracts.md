# Contracts: Avatars

> Unit: `Avatars`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base path: `/v1`. Scope: `avatars.read` on all endpoints. All responses binary image data with long `Cache-Control`.

| Method | Path | Params | Success | Errors |
|--------|------|--------|---------|--------|
| GET | `/avatars/browsers/:code` | — | 200 image/png | 404 AVATAR_NOT_FOUND, 400 AVATAR_SET_NOT_FOUND |
| GET | `/avatars/credit-cards/:code` | — | 200 image/png | 404/400 |
| GET | `/avatars/flags/:code` | — | 200 image/png | 404/400 |
| GET | `/avatars/favicon` | `url`, `width`, `height` | 200 image (png/ico/svg+xml) | 400 invalid URL/SSRF |
| GET | `/avatars/image` | `url`, `width`, `height` | 200 image/* | 400 invalid URL |
| GET | `/avatars/initials` | `name`, `width`, `height`, `background` (hex) | 200 image/png | 400 missing name |
| GET | `/avatars/qr` | `text`, `size`, `margin` | 200 image/png | 400 missing text |
| GET | `/avatars/screenshots` | `url`, `width`, `height`, `theme`, `locale`, `timezone`, `userAgent`, `geolocation`, `permissions`, `sleep` | 200 image/* | 400 invalid URL, 502 browser service |
| GET | `/cards/cloud` | — | 200 image/png | 401 |
| GET | `/cards/cloud-back` | — | 200 image/png | 401 |
| GET | `/cards/cloud-og` | — | 200 image/png | 401 |

## External service contract (browser microservice)

| Endpoint | Method | Payload |
|----------|--------|---------|
| `{_APP_BROWSER_HOST}/v1/screenshot` | POST/GET | `{ url, viewport, theme, geolocation, locale, timezone, permissions, userAgent, sleep }` → binary screenshot |

Default `_APP_BROWSER_HOST` = `http://appwrite-browser:3000/v1`.

## Error model

Standard Appwrite error JSON `{ message, code, type, version }` for 400/401/404/502. Binary endpoints return HTTP 200 + image body on success.
