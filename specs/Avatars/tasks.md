# Tasks: Avatars

> Unit: `Avatars`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-AVT-01 — Register Avatars module + base action helpers
- **Behavior source:** `Modules/Avatars/Module.php`, `Modules/Avatars/Http/Action.php`
- **Description:** Register HTTP service; implement shared `avatar()` and `getUserGitHub()` helpers.
- **Done when:** Module loads; helpers functional.
- **Priority:** Must · Confidence: 🟢

### T-AVT-02 — Register 11 HTTP routes
- **Behavior source:** `Modules/Avatars/Services/Http.php`
- **Description:** Register all GET routes with `avatars.read` scope and correct action IDs.
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-AVT-03 — Static avatars (browsers, credit cards, flags)
- **Behavior source:** `Modules/Avatars/Http/{Browsers,CreditCards,Flags}/Get.php`
- **Description:** Load assets from config, whitelist codes, crop via Imagick, output PNG with cache headers.
- **Done when:** Valid codes render; unknown codes → AVATAR_NOT_FOUND.
- **Priority:** Must · Confidence: 🟢

### T-AVT-04 — Favicon proxy
- **Behavior source:** `Modules/Avatars/Http/Favicon/Get.php`
- **Description:** SSRF-safe fetch with ≤5 redirects; SVG sanitize / ICO passthrough / 56×56 crop.
- **Done when:** Private hosts blocked; sanitized SVG served; crop applied.
- **Priority:** Must · Confidence: 🟢

### T-AVT-05 — Remote image proxy
- **Behavior source:** `Modules/Avatars/Http/Image/Get.php`
- **Description:** Fetch + resize remote image with SSRF guards.
- **Done when:** Safe URL resized; private host rejected.
- **Priority:** Must · Confidence: 🟢

### T-AVT-06 — Initials avatar
- **Behavior source:** `Modules/Avatars/Http/Initials/Get.php`
- **Description:** Deterministic theme (char codes % 5) + custom hex background; render via Imagick.
- **Done when:** Same name → same color; override honored.
- **Priority:** Must · Confidence: 🟢

### T-AVT-07 — QR code generation
- **Behavior source:** `Modules/Avatars/Http/QR/Get.php`
- **Description:** Render QR PNG from text input.
- **Done when:** Valid QR PNG produced and scannable.
- **Priority:** Must · Confidence: 🟢

### T-AVT-08 — Screenshot proxy
- **Behavior source:** `Modules/Avatars/Http/Screenshots/Get.php`
- **Description:** Forward options to `_APP_BROWSER_HOST`; meter usage; stream response.
- **Done when:** Screenshot returned; metric incremented.
- **Priority:** Should · Confidence: 🟢

### T-AVT-09 — Cloud cards (front/back/OG)
- **Behavior source:** `Modules/Avatars/Http/Cards/Cloud/{Front,Back,OG}/Get.php`
- **Description:** Tier classification + GitHub identity resolution + tier-specific PNG render.
- **Done when:** Correct tier PNG for employee/hero/contributor/platinum/golden.
- **Priority:** Should · Confidence: 🟢

### T-AVT-10 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Avatars/Http/**`
- **Description:** OpenAPI 3.1 for 11 endpoints with binary responses.
- **Done when:** `openapi/avatars.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-AVT-11 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Avatars/**`
- **Description:** Cover static avatars, favicon SSRF blocks, redirect limits, initials determinism, QR, screenshots, cloud cards.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
