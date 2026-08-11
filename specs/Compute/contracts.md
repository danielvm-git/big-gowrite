# Contracts: Compute

> Unit: `Compute`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

None — Compute is a shared library module with no HTTP surface.

## Library API contract

### `Specification` constants
```php
Specification::S_0_5_VCPU_512MB   // "s-0.5vcpu-512mb"
Specification::S_1_VCPU_1GB       // "s-1vcpu-1gb"
... up to ...
Specification::S_8_VCPU_8GB       // "s-8vcpu-8gb"
```

### `Base` methods (public surface consumed by Functions/Sites)
| Method | Signature | Returns |
|--------|-----------|---------|
| `redeployVcsFunction` | `(Deployment $deployment, ...)` | Deployment document |
| `redeployVcsSite` | `(Deployment $deployment, ...)` | Deployment document |
| `activateBranchPreviewRule` | `(Deployment $deployment, ...)` | void |
| `updateEmptyManualRule` | `(...)` | void |
| `getDefaultSpecification` | `(bool $preferFallback = false)` | string (spec id) |
| `resolveTemplateRef` | `(string $ref, string $provider, string $repo)` | string (concrete tag) |

### `Validator/Specification` interface
- Input: candidate spec id (string)
- Context: billing plan `runtimeSpecifications`, env CPU/memory caps
- Output: `true` if allowed, error otherwise

## Queue contract

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForBuilds` | Compute (via Functions/Sites Base) → worker-builds | `BuildMessage`: deploymentId, resourceType (`functions`/`sites`), template ref, repo/branch/commit — only for template-based VCS deployments |

## Configuration contract

| Env var | Effect |
|---------|--------|
| `_APP_RULES_FORMAT=md5` | Legacy MD5 rule IDs for backward compatibility |
| `APP_COMPUTE_SPECIFICATION_DEFAULT` | Fallback default spec constant |

## Events

None emitted by Compute directly; events emitted by Functions/Sites modules extending Base.
