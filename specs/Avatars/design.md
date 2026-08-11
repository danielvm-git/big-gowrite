# Design: Avatars

> Unit: `Avatars`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Avatars module is a **stateless binary-output HTTP service**. It never writes to a database; it reads static config assets, generates images in-memory (Imagick), proxies remote fetches (favicon, image) with SSRF guards, and proxies screenshots to an internal browser microservice. Cloud cards are the only endpoints touching databases (session/token lookup for GitHub identity).

```
Client
  │ avatars.read scope
  ▼
app/http.php → Avatars\Http\* (11 actions) → binary image response
  │
  ├── Config (avatar sets: browsers, credit-cards, flags)
  ├── Imagick (rendering)
  ├── SvgSanitizer (favicon SVG)
  ├── PublicHostname (SSRF allowlist) + HTTP client (max 5 redirects)
  ├── Browser microservice (_APP_BROWSER_HOST) — screenshots
  ├── dbForProject / dbForPlatform (cloud cards: sessions, users)
  └── usage (screenshot metrics)
```

## 2. Route-to-action map

| Method | Path | Action class | Action ID | Scope |
|--------|------|-------------|-----------|-------|
| GET | `/v1/avatars/browsers/:code` | `Browsers/Get` | `getBrowser` | avatars.read |
| GET | `/v1/avatars/credit-cards/:code` | `CreditCards/Get` | `getCreditCard` | avatars.read |
| GET | `/v1/avatars/flags/:code` | `Flags/Get` | `getFlag` | avatars.read |
| GET | `/v1/avatars/favicon` | `Favicon/Get` | `getFavicon` | avatars.read |
| GET | `/v1/avatars/image` | `Image/Get` | `getImage` | avatars.read |
| GET | `/v1/avatars/initials` | `Initials/Get` | `getInitials` | avatars.read |
| GET | `/v1/avatars/qr` | `QR/Get` | `getQR` | avatars.read |
| GET | `/v1/avatars/screenshots` | `Screenshots/Get` | `getScreenshot` | avatars.read |
| GET | `/v1/cards/cloud` | `Cards/Cloud/Front/Get` | `getCloudCard` | avatars.read |
| GET | `/v1/cards/cloud-back` | `Cards/Cloud/Back/Get` | `getCloudCardBack` | avatars.read |
| GET | `/v1/cards/cloud-og` | `Cards/Cloud/OG/Get` | `getCloudCardOG` | avatars.read |

## 3. Shared base (`Http/Action.php`)

- `avatar(string $key, ...)`: loads static PNG from `Config::getParam('avatar-{type}')`, crops to size, outputs as binary with `Cache-Control`.
- `getUserGitHub(...)`: finds first session with `provider == 'github'`, refreshes token if expired (retry loop 10×500 ms), returns GitHub identity for cloud cards.

## 4. Data flow — key actions

### 4.1 Static avatar (browsers/credit-cards/flags)

```
1. Validate :code against WhiteList from Config avatar-{type}
2. Load PNG asset, crop with Imagick
3. Output binary PNG + cache headers
```

### 4.2 Favicon fetch

```
1. Validate URL scheme (http/https) + PublicHostname allowlist
2. Fetch (max 5 manual redirects, re-validating each hop)
3. Branch by content type:
   - SVG → SvgSanitizer → output
   - ICO → output as-is
   - PNG/JPG/JPEG → crop 56×56
```

### 4.3 Initials generation

```
1. Compute deterministic theme index: sum(char codes of name) % 5
2. Use default theme color OR custom hex background param
3. Render initials text onto PNG canvas via Imagick
```

### 4.4 Screenshot capture

```
1. Validate URL (SSRF guards)
2. Build options (viewport, theme, geolocation, locale, timezone, permissions, user-agent, sleep)
3. POST to _APP_BROWSER_HOST/v1 (screenshot endpoint)
4. Track METRIC_AVATARS_SCREENSHOTS_GENERATED
5. Stream binary response
```

### 4.5 Cloud card render

```
1. Resolve user (dbForPlatform)
2. Classify tier: isEmployee > isHero > isContributor; sequence % 100 == 0 → isPlatinum; golden overrides platinum
3. Resolve GitHub identity (dbForProject sessions, token refresh)
4. Load tier-specific PNG asset
5. Output PNG
```

## 5. Dependency injection

| Injection | Role | Used by |
|-----------|------|---------|
| `response` | Binary output | All actions |
| `user` | Current user | Initials, Cloud cards |
| `project` | Project config | Cloud cards |
| `dbForProject` | Sessions (GitHub token) | Cloud cards |
| `dbForPlatform` | User lookup | Cloud cards |
| `heroes`, `contributors`, `employees` | Tier arrays from config | Cloud cards |
| `logger` | Error logging | Cloud cards |
| `authorization` | Privilege escalation (skip-auth lookups) | Cloud cards |
| `usage` | Usage metrics | Screenshots |

## 6. Security

- **SSRF:** `PublicHostname` allowlist, http/https only, manual redirect chain (≤5) with per-hop re-validation.
- **SVG sanitization:** `SvgSanitizer` before serving foreign SVG favicons.
- **Scope:** all endpoints `avatars.read`.
- **No persistence:** module stateless except cloud-card DB reads.
