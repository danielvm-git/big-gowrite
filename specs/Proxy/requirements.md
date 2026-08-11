# Requirements: Proxy

> Identifier: `PRX-DOMAINS`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Proxy module manages custom domain routing rules. Projects serve the Appwrite API, a Function execution endpoint, a Site deployment, or a URL redirect under a user-controlled domain. The module handles rule creation with DNS verification, TLS certificate provisioning (via Certificate publisher), and rule lifecycle (list, get, delete, retry verification). Used by Functions and Sites to expose deployments on branded domains.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Proxy` | Module purpose | 🟢 |
| `specs/Proxy/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/security.md#DNS` | DNS verification, CAA | 🟢 |
| `specs/workers.md#Certificates` | TLS generation worker | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App developer | Branded API domain | Create API rule for custom domain |
| App developer | Serve site on custom domain | Create site rule |
| App developer | Redirect domains | Redirect rule with status code |
| Platform ops | Domain safety | Deny-list, Appwrite-owned domains |

## 4. Business rules

1. **BR-01:** Rule IDs: `md5(strtolower(domain))` when `_APP_RULES_FORMAT=md5`; else `ID::unique()`. 🟢
2. **BR-02:** Appwrite-owned domain patterns (`_APP_DOMAIN_FUNCTIONS`, `_APP_DOMAIN_SITES`, `_APP_DOMAIN_FUNCTIONS_FALLBACK`) auto-verified, skip DNS checks (status `verified`, owner `'Appwrite'`). 🟢
3. **BR-03:** Creation fires domain validation; DNS pass → `certificate_generating` + Certificate publish enqueued. 🟢
4. **BR-04:** Internal `created` status remapped to `unverified` in responses. 🟢
5. **BR-05:** Platform hostnames, `_APP_DOMAIN_SITES`, `_APP_DOMAIN_FUNCTIONS`, `_APP_CUSTOM_DOMAIN_DENY_LIST` domains denied. 🟢
6. **BR-06:** Duplicate rules: if conflicting rule's project gone, orphan deleted and creation retried. 🟢
7. **BR-07:** CAA record check (`_APP_DOMAIN_TARGET_CAA`) gates issuance before CNAME/A/AAAA validation. 🟢
8. **BR-08:** CNAME targets: Sites → `_APP_DOMAIN_SITES`; Functions/API → `_APP_DOMAIN_TARGET_CNAME`. 🟢
9. **BR-09:** Redirect rules require parent resource (`site` or `function`). 🟢
10. **BR-10:** list/get hydrate certificate logs + `renewAt`. 🟢
11. **BR-11:** Rule delete enqueues `DELETE_TYPE_DOCUMENT` for certificates. 🟢
12. **BR-12:** `updateRuleVerification` no-op if already `verified` or `certificate_generating`. 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Create API rule | Must | DNS validation; cert enqueue | 🟢 |
| RF-02 | Create Function rule | Must | Resource existence validated | 🟢 |
| RF-03 | Create Site rule | Must | Sites CNAME target | 🟢 |
| RF-04 | Create Redirect rule | Must | Parent resource + status code (301/302/307/308) | 🟢 |
| RF-05 | Get rule | Must | Certificate logs + renewAt hydrated | 🟢 |
| RF-06 | List rules | Must | Project-scoped | 🟢 |
| RF-07 | Delete rule | Must | Cleanup enqueue | 🟢 |
| RF-08 | Retry verification | Must | DNS retry → cert generation; no-op on terminal states | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Domain deny-list enforcement | platform hostnames + deny list | 🟢 |
| Security | CAA gating before issuance | `_APP_DOMAIN_TARGET_CAA` | 🟢 |
| Reliability | Orphan rule recovery | Duplicate handling | 🟢 |
| Compatibility | Legacy MD5 rule IDs | `_APP_RULES_FORMAT=md5` | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Appwrite-owned domain auto-verified
  Given a domain matching _APP_DOMAIN_SITES pattern
  When a rule is created for it
  Then the rule is verified without DNS checks and owner is Appwrite

Scenario: Denied domain
  Given a domain in _APP_CUSTOM_DOMAIN_DENY_LIST
  When a rule creation is attempted
  Then the request is rejected

Scenario: Retry on terminal state
  Given a rule already verified
  When PATCH .../status is called
  Then the call is a no-op

Scenario: Redirect without parent
  Given a redirect rule without site/function parent
  When the rule is created
  Then the request is rejected
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..04 (create) | Must | Core routing |
| RF-05..07 (lifecycle) | Must | Management |
| RF-08 (retry) | Must | DNS recovery |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🔴 [GAP] Exact DNS verification record set (TXT/CNAME) not captured in module.
- 🟡 [INFERRED] Certificate document schema.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
