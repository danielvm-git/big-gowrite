# Requirements: Avatars

> Identifier: `AVT-IMAGES`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Avatars module provides image-generation and image-proxy endpoints: static avatar assets (browser icons, credit card logos, country flags), dynamic images (user initials, QR codes), remote image/favicon fetching with resizing, webpage screenshots via a headless browser service, and Appwrite Cloud membership cards. All responses are binary image data; nothing writes to a database. Responses are cached aggressively (30–45 days via `Cache-Control`) and grouped under `avatars.read` scope.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Avatars` | Module purpose | 🟢 |
| `specs/Avatars/legacy-mapping.md` | 11 endpoints, business rules, DI | 🟢 |
| `specs/domain.md#Avatars` | Avatar domain, URL validation | 🟢 |
| `specs/architecture.md#Integrations` | Browser microservice integration | 🟢 |
| `specs/security.md#SSRF` | SSRF protections (allowlist, redirect limits) | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App user | Personalize profile | Generate initials avatar or QR code |
| App developer | Proxy external images | Fetch and crop a remote image safely |
| Console admin | Preview websites | Capture a screenshot of a site |
| Marketing | Membership card | Render Cloud membership card |

## 4. Business rules

1. **BR-01:** Static avatar sets loaded from `Config::getParam('avatar-{type}')`, validated via `WhiteList`; unknown codes throw `AVATAR_SET_NOT_FOUND` or `AVATAR_NOT_FOUND`. 🟢
2. **BR-02:** Static avatar rendering requires Imagick PHP extension; missing → `GENERAL_SERVER_ERROR`. 🟢
3. **BR-03:** Remote URLs (favicon, image, screenshot) validated against `PublicHostname` allowlist; only `http`/`https`; private/loopback hosts rejected. 🟢
4. **BR-04:** Favicon fetching follows max 5 redirects manually (`MAX_REDIRECTS = 5`), each hop re-validated (SSRF). 🟢
5. **BR-05:** Favicon formats: SVG sanitized via `SvgSanitizer`; ICO returned as-is; PNG/JPG cropped to 56×56. 🟢
6. **BR-06:** Initials background color deterministic from name char codes (5 themes, indexed mod 5); custom hex `background` overrides. 🟢
7. **BR-07:** Screenshots proxy to internal browser microservice `_APP_BROWSER_HOST` (default `http://appwrite-browser:3000/v1`); options: viewport, theme, geolocation, locale, timezone, permissions, user-agent, sleep. 🟢
8. **BR-08:** Cloud card tiers: `isEmployee` > `isHero` > `isContributor` (mutually exclusive); every 100th user `isPlatinum`; golden overrides platinum. 🟢
9. **BR-09:** Cloud cards resolve GitHub identity via first `provider == 'github'` session; token refreshed if expired (retry up to 10 attempts, 500 ms apart). 🟡
10. **BR-10:** Screenshots usage-metered via `METRIC_AVATARS_SCREENSHOTS_GENERATED`. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Get browser icon PNG by code | Must | Valid code → PNG; unknown → AVATAR_NOT_FOUND | 🟢 |
| RF-02 | Get credit card logo PNG by code | Must | Valid code → PNG; unknown → error | 🟢 |
| RF-03 | Get country flag PNG by ISO Alpha-2 | Must | Valid code → PNG; unknown → error | 🟢 |
| RF-04 | Fetch and return favicon for URL | Must | Safe URL → sanitized favicon; private host rejected; ≤5 redirects | 🟢 |
| RF-05 | Fetch and resize remote image | Must | Safe URL → cropped/resized image; SSRF guards | 🟢 |
| RF-06 | Generate initials PNG avatar | Must | Deterministic color; custom background override | 🟢 |
| RF-07 | Generate QR code PNG from text | Must | Text → QR PNG | 🟢 |
| RF-08 | Capture webpage screenshot | Must | Proxy to browser service; options forwarded; usage metered | 🟢 |
| RF-09 | Render Cloud card front | Must | Tier-based PNG (employee/hero/contributor/platinum/golden) | 🟢 |
| RF-10 | Render Cloud card back | Must | Back PNG rendered | 🟢 |
| RF-11 | Render Cloud card OG image | Must | OG PNG rendered | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | SSRF: hostname allowlist, scheme whitelist, ≤5 redirects | `PublicHostname` + `MAX_REDIRECTS=5` | 🟢 |
| Security | SVG sanitization before serving | `SvgSanitizer` | 🟢 |
| Performance | Aggressive caching 30–45 days Cache-Control | HTTP headers in action base | 🟢 |
| Reliability | Screenshot proxy options incl. retry/sleep | `_APP_BROWSER_HOST` + options | 🟢 |
| Observability | Screenshot usage metered | `METRIC_AVATARS_SCREENSHOTS_GENERATED` | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Get browser icon
  Given a valid browser code
  When GET /v1/avatars/browsers/:code is called
  Then a PNG with cache headers is returned

Scenario: Unknown avatar code
  Given an unknown code
  When GET /v1/avatars/browsers/:code is called
  Then AVATAR_NOT_FOUND error is returned

Scenario: SSRF blocked on private host
  Given a URL resolving to 127.0.0.1
  When GET /v1/avatars/favicon?url=... is called
  Then the request is rejected with an error

Scenario: Redirect chain over limit
  Given a URL with more than 5 redirects
  When GET /v1/avatars/favicon?url=... is called
  Then the request fails safely without reaching an internal host
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..03 (static avatars) | Must | Core feature set |
| RF-04 (favicon) | Must | SSRF-guarded fetch |
| RF-05 (image proxy) | Must | Remote image resizing |
| RF-06 (initials) | Must | Deterministic generation |
| RF-07 (QR) | Must | Common use |
| RF-08 (screenshots) | Should | Depends on browser service |
| RF-09..11 (cloud cards) | Should | Cloud-specific, undocumented |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Cloud card tier arrays (heroes, contributors, employees) source not verified.
- 🟡 [INFERRED] GitHub token refresh retry loop inferred from code comments.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |

