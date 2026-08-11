# Design: Proxy

> Unit: `Proxy`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Proxy module is an **HTTP service managing custom-domain routing rules**. Rules are stored in `dbForPlatform` (`rules`, `certificates`). Creation validates DNS (CAA → CNAME/A/AAAA), enqueues certificate generation, and stores rule docs. Deletion enqueues cleanup.

```
Client
  │ scopes: rules.*
  ▼
app/http.php → Proxy\Http\** (8 actions) + Action.php base
  │
  ├── dbForPlatform (rules, certificates)
  ├── dbForProject (resource existence: functions, sites)
  ├── publisherForCertificates (TLS jobs)
  ├── publisherForDeletes (cleanup)
  └── queueForEvents (rule events)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/proxy/rules/api` | `createAPIRule` | rules.write |
| POST | `/v1/proxy/rules/function` | `createFunctionRule` | rules.write |
| POST | `/v1/proxy/rules/site` | `createSiteRule` | rules.write |
| POST | `/v1/proxy/rules/redirect` | `createRedirectRule` | rules.write |
| GET | `/v1/proxy/rules` | `listRules` | rules.read |
| GET | `/v1/proxy/rules/:ruleId` | `getRule` | rules.read |
| DELETE | `/v1/proxy/rules/:ruleId` | `deleteRule` | rules.write |
| PATCH | `/v1/proxy/rules/:ruleId/status` | `updateRuleVerification` | rules.write |

## 3. Data flows — key actions

### 3.1 Create rule (shared base)

```
1. Validate domain restrictions (hostname list, deny list, Appwrite-owned patterns)
2. Appwrite-owned → verified (owner=Appwrite), skip DNS
3. Else:
   a. CAA record check (_APP_DOMAIN_TARGET_CAA)
   b. CNAME/A/AAAA validation (target by type: sites → _APP_DOMAIN_SITES; functions/API → _APP_DOMAIN_TARGET_CNAME)
4. Duplicate handling: orphan rule cleanup + retry
5. Status → certificate_generating; enqueue Certificate publish
6. Persist rule (ID: md5(domain) legacy or ID::unique())
7. Return rule (created → unverified remap)
```

### 3.2 Get/List rule

```
1. Fetch rule(s) scoped to project
2. Hydrate certificate logs + renewAt from certificates doc
```

### 3.3 Delete rule

```
1. Fetch + verify ownership
2. Delete rule doc
3. Enqueue DELETE_TYPE_DOCUMENT (certificate cleanup)
4. Emit rules.[ruleId].delete
```

### 3.4 Retry verification

```
1. If verified or certificate_generating → no-op
2. Re-run DNS validation
3. On success: trigger TLS generation
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForPlatform` | rules, certificates |
| `dbForProject` | Resource existence (functions, sites) |
| `publisherForCertificates` | TLS generation |
| `publisherForDeletes` | Certificate cleanup |
| `queueForEvents` | Rule events |
| `project` | projectId, internalId, region |
| `authorization` | skip() for platform records |
| `platform` | Hostname deny-list |
| `log` | DNS timing diagnostics |

## 5. Events emitted

- `rules.[ruleId].create`
- `rules.[ruleId].update`
- `rules.[ruleId].delete`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| Domain denied | 400 | Hostname/deny-list match |
| Duplicate rule | 400 | Existing active rule |
| Resource missing | 400 | Function/site not found |
| DNS validation fail | 400 | CAA/CNAME checks fail |

## 7. Security

- Deny-list + platform hostname enforcement.
- CAA gating before issuance.
- Project scoping on all queries.
