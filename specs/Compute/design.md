# Design: Compute

> Unit: `Compute`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

Compute is a **shared library module** (no HTTP routes, no workers). It sits under `Modules/Compute/` and is consumed by the Functions and Sites modules through class inheritance (`Base` abstract action) and direct service usage (`Specification`, `Validator/Specification`).

```
Functions module ──┐
                   ├──► Modules/Compute/Base.php (abstract Action)
Sites module ──────┘
                        ├── Specification.php (size constants)
                        └── Validator/Specification.php (plan gating)
```

## 2. Components

| Component | Type | Responsibility |
|-----------|------|----------------|
| `Base.php` | Abstract `Action` | VCS deployment helpers for functions/sites |
| `Specification.php` | Constants | Compute size identifiers (s-0.5vcpu-512mb … s-8vcpu-8gb) |
| `Validator/Specification.php` | Utopia `Validator` | Filters allowed specs by plan + env caps |

## 3. Data flow — key helpers

### 3.1 `redeployVcsFunction`

```
1. Resolve repo + branch/commit via VCS adapter
2. If template deployment:
   a. Create deployments document (resourceType='functions') directly
   b. Enqueue BuildMessage → publisherForBuilds
3. Else:
   a. Deployments::createFromUrl (presigned URL flow)
4. Return deployment
```

### 3.2 `redeployVcsSite`

```
1. Resolve repo + branch/commit via VCS adapter
2. Create deployments document (resourceType='sites')
3. Mint rules in dbForPlatform (authorization skip):
   a. unique deployment domain
   b. optional commit-hash-prefixed domain
   c. branch-preview domain
4. Backfill empty-deployment-id manual rules (updateEmptyManualRule)
5. Return deployment
```

### 3.3 `activateBranchPreviewRule`

```
1. Upsert branch-preview rule → new deploymentId
2. Iterate manual rules pinned to branch; update deploymentId
```

### 3.4 `resolveTemplateRef`

```
1. If ref contains wildcard (e.g. 0.3.*):
   a. Query GitHub VCS adapter for tags matching pattern
   b. Pick highest concrete tag
2. On error: return raw reference
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | Deployments collection |
| `dbForPlatform` | Rules collection |
| `publisherForBuilds` | Build queue (template deployments) |
| `authorization` | `skip()` for cross-project rule writes |
| `vcs` | Git adapter (tags, commits, repo URL) |
| `vcsFactory` | Provider-specific adapters for template resolution |
| `deployments` | Deployments service (presigned URL flow) |

## 5. Queue contract

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForBuilds` | Compute → worker-builds | `BuildMessage` (deploymentId, template info) — only for template deployments |

## 6. Security

- Rule minting uses `authorization->skip()` for controlled cross-project writes in dbForPlatform.
- Spec validation prevents resource over-allocation beyond billing plan.
