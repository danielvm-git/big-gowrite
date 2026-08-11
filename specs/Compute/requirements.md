# Requirements: Compute

> Identifier: `CMP-SHARED-LIB`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Compute module is a **shared infrastructure library** — not a standalone API surface. It provides reusable base classes and helpers used by the Functions and Sites modules: (1) `Base` — a base `Action` class with methods for deploying VCS-backed functions/sites (deployment documents, preview-domain rules, template refs), (2) `Specification` — enum-style constants for compute size tiers (vCPU/memory), and (3) `Validator/Specification` — a validator that filters allowed compute specs against billing plan limits and environment-variable caps. It exposes no HTTP routes.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Compute` | Shared library purpose | 🟢 |
| `specs/Compute/legacy-mapping.md` | Components, business rules, DI | 🟢 |
| `specs/domain.md#Compute` | Compute sizes, VCS deployment model | 🟢 |
| `specs/architecture.md#Workers` | Builds worker consumes BuildMessage | 🟢 |
| `specs/adrs/ADR-*.md` | Preview domains, rules format | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Functions module | Deploy VCS-backed functions | Reuse `redeployVcsFunction` from Base |
| Sites module | Deploy VCS-backed sites with preview domains | Reuse `redeployVcsSite` + rule minting |
| Platform billing | Enforce compute tier limits | Validate spec against plan via Validator |

## 4. Business rules

1. **BR-01:** `getDefaultSpecification` returns the lowest allowed spec in the billing plan, or highest available when no plan set. With `preferFallback=true` and no `planKey`, constant `APP_COMPUTE_SPECIFICATION_DEFAULT` is returned. 🟢
2. **BR-02:** `resolveTemplateRef` resolves wildcard tag refs (e.g. `0.3.*`) to the highest matching concrete tag via the GitHub VCS adapter; falls back to the raw reference on error. 🟢
3. **BR-03:** `redeployVcsFunction` creates a `deployments` document with `resourceType='functions'`. Non-template VCS deployments use `Deployments::createFromUrl`; template deployments create the document directly and enqueue a `BuildMessage` to the builds publisher. 🟢
4. **BR-04:** `redeployVcsSite` creates a `deployments` document with `resourceType='sites'` and mints three `rules` docs in `dbForPlatform`: unique deployment domain, optional commit-hash-prefixed domain, branch-preview domain. Empty-deployment-id manual rules backfilled via `updateEmptyManualRule`. 🟢
5. **BR-05:** `activateBranchPreviewRule` upserts the branch-preview rule to the new deployment and updates manual rules pinned to the same branch. 🟢
6. **BR-06:** Template-based VCS deployments handled by the Builds worker; plain VCS deployments go through the Deployments service directly. 🟡
7. **BR-07:** `_APP_RULES_FORMAT=md5` triggers legacy MD5 rule IDs for backward compatibility. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Provide `getDefaultSpecification` helper | Must | Returns lowest allowed spec; fallback behavior per plan | 🟢 |
| RF-02 | Provide `resolveTemplateRef` helper | Must | Wildcard → highest concrete tag; raw fallback on error | 🟢 |
| RF-03 | Provide `redeployVcsFunction` | Must | Creates deployment doc (resourceType=functions); enqueues BuildMessage for templates | 🟢 |
| RF-04 | Provide `redeployVcsSite` | Must | Creates deployment doc (resourceType=sites) + 3 rule docs | 🟢 |
| RF-05 | Provide `activateBranchPreviewRule` | Must | Upserts branch-preview rule; updates same-branch manual rules | 🟢 |
| RF-06 | Expose compute Specification constants | Must | All size identifiers from s-0.5vcpu-512mb to s-8vcpu-8gb | 🟢 |
| RF-07 | Validate spec against plan + env caps | Must | Reject specs above plan `runtimeSpecifications` or env CPU/mem caps | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Compatibility | Legacy MD5 rule IDs supported via env flag | `_APP_RULES_FORMAT=md5` | 🟢 |
| Reliability | VCS adapter failures fall back gracefully (raw ref) | `resolveTemplateRef` try/catch | 🟢 |
| Maintainability | Single source for compute tiers | `Specification` constants reused by Functions/Sites | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Default spec resolution with plan
  Given a billing plan with allowed runtimeSpecifications
  When getDefaultSpecification is called
  Then the lowest allowed specification is returned

Scenario: Wildcard template ref resolution
  Given a tag reference "0.3.*"
  When resolveTemplateRef is called with the GitHub adapter
  Then the highest matching concrete tag is returned

Scenario: VCS site deployment mints rules
  Given a site VCS deployment
  When redeployVcsSite is called
  Then a deployments document and 3 rules documents are created
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..05 (Base helpers) | Must | Shared by Functions/Sites |
| RF-06 (Spec constants) | Must | Canonical size list |
| RF-07 (Validator) | Must | Billing enforcement |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] Builds worker template-deployment flow details (BuildMessage payload) not fully verified in this module.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
