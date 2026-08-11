# Tasks: Compute

> Unit: `Compute`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-CMP-01 — Implement compute Specification constants
- **Behavior source:** `Modules/Compute/Specification.php`
- **Description:** Define string constants for all compute size identifiers (s-0.5vcpu-512mb … s-8vcpu-8gb).
- **Done when:** Constants match legacy set; no HTTP surface.
- **Priority:** Must · Confidence: 🟢

### T-CMP-02 — Implement Specification validator
- **Behavior source:** `Modules/Compute/Validator/Specification.php`
- **Description:** Filter allowed specifications against plan `runtimeSpecifications` and env-var CPU/memory caps.
- **Done when:** Out-of-plan specs rejected; caps enforced.
- **Priority:** Must · Confidence: 🟢

### T-CMP-03 — Implement `getDefaultSpecification`
- **Behavior source:** `Modules/Compute/Base.php`
- **Description:** Return lowest allowed spec in plan, or highest available without plan; `preferFallback` → `APP_COMPUTE_SPECIFICATION_DEFAULT`.
- **Done when:** All branches behave per legacy rules.
- **Priority:** Must · Confidence: 🟢

### T-CMP-04 — Implement `resolveTemplateRef`
- **Behavior source:** `Modules/Compute/Base.php`
- **Description:** Resolve wildcard tag refs via GitHub VCS adapter to highest concrete tag; raw fallback on error.
- **Done when:** `0.3.*` resolves; adapter failure returns raw ref.
- **Priority:** Must · Confidence: 🟢

### T-CMP-05 — Implement `redeployVcsFunction`
- **Behavior source:** `Modules/Compute/Base.php`
- **Description:** Create deployments doc (resourceType=functions); template → BuildMessage; else → Deployments::createFromUrl.
- **Done when:** Both paths produce correct deployment docs and queue messages.
- **Priority:** Must · Confidence: 🟢

### T-CMP-06 — Implement `redeployVcsSite` + rule minting
- **Behavior source:** `Modules/Compute/Base.php`
- **Description:** Create deployments doc (resourceType=sites); mint deployment/commit/branch-preview rules; backfill empty manual rules.
- **Done when:** Deployment + 3 rules created; backfill applied.
- **Priority:** Must · Confidence: 🟢

### T-CMP-07 — Implement `activateBranchPreviewRule`
- **Behavior source:** `Modules/Compute/Base.php`
- **Description:** Upsert branch-preview rule; update manual rules pinned to branch.
- **Done when:** Preview rule repoints to new deployment.
- **Priority:** Must · Confidence: 🟢

### T-CMP-08 — Legacy rule ID compatibility
- **Behavior source:** `Modules/Compute/Base.php`
- **Description:** Honor `_APP_RULES_FORMAT=md5` for backward-compatible MD5 rule IDs.
- **Done when:** md5 format produces legacy IDs.
- **Priority:** Should · Confidence: 🟢

### T-CMP-09 — Unit tests
- **Behavior source:** `tests/unit/**` (compute spec validator, base helpers)
- **Description:** Unit tests for validator gating, spec resolution, template ref resolution.
- **Done when:** `docker compose exec appwrite test tests/unit/` green.
- **Priority:** Must · Confidence: 🟢
