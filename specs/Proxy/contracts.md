# Contracts: Proxy

> Unit: `Proxy`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/proxy/rules`. Scopes: `rules.read` / `rules.write`.

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `/rules/api` | `domain` | 201 rule | 400 denied/invalid DNS |
| POST | `/rules/function` | `domain`, `functionId`, `deploymentId` | 201 rule | 400 resource |
| POST | `/rules/site` | `domain`, `siteId`, `deploymentId` | 201 rule | 400 resource |
| POST | `/rules/redirect` | `domain`, `siteId`/`functionId`, `statusCode` | 201 rule | 400 parent |
| GET | `/rules` | `cursor`, `limit` | 200 ruleList | 401 |
| GET | `/rules/:ruleId` | — | 200 rule (+ cert) | 404 |
| DELETE | `/rules/:ruleId` | — | 204 | 404 |
| PATCH | `/rules/:ruleId/status` | — | 200 rule | 404 |

## Rule status model

`created` (internal) → `unverified` (external remap) → `verified` → `certificate_generating` → active/ready.

## DNS verification contract

| Check | Order | Config |
|-------|-------|--------|
| CAA record | First | `_APP_DOMAIN_TARGET_CAA` |
| CNAME/A/AAAA | Second | Sites → `_APP_DOMAIN_SITES`; Functions/API → `_APP_DOMAIN_TARGET_CNAME` |

Appwrite-owned patterns (`_APP_DOMAIN_FUNCTIONS`, `_APP_DOMAIN_SITES`, `_APP_DOMAIN_FUNCTIONS_FALLBACK`) skip DNS → verified (owner=Appwrite).

## Queue contracts

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForCertificates` | Rule create/retry → worker-certificates | Certificate generation job (domain, project, ruleId) |
| `publisherForDeletes` | Rule delete → worker-deletes | `DELETE_TYPE_DOCUMENT` (certificates cleanup) |

## Event contract

- `rules.[ruleId].create`
- `rules.[ruleId].update`
- `rules.[ruleId].delete`

## Configuration contract

| Env var | Effect |
|---------|--------|
| `_APP_RULES_FORMAT=md5` | Legacy rule IDs: `md5(strtolower(domain))` |
| `_APP_DOMAIN_FUNCTIONS`, `_APP_DOMAIN_SITES`, `_APP_DOMAIN_FUNCTIONS_FALLBACK` | Appwrite-owned patterns |
| `_APP_CUSTOM_DOMAIN_DENY_LIST` | Denied domains |
| `_APP_DOMAIN_TARGET_CAA` | CAA record value |
| `_APP_DOMAIN_TARGET_CNAME` | CNAME target for functions/API |
