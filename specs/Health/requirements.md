# Requirements: Health

> Identifier: `HLT-DIAGNOSTICS`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Health module is a comprehensive diagnostic API for server administrators. It exposes read-only endpoints probing each infrastructure component — HTTP layer, database pools, cache (Redis), pub/sub, NTP time sync, SSL certificates, antivirus (ClamAV), geo service, storage devices, and all background-job queues — returning pass/fail status with latency. Queue endpoints enforce configurable size thresholds to detect worker stalls.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Health` | Module purpose | 🟢 |
| `specs/Health/legacy-mapping.md` | 25 endpoints, business rules, DI | 🟢 |
| `specs/architecture.md#Infrastructure` | Component topology | 🟢 |
| `specs/security.md` | Scope gating (root, health.read) | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| Server admin | Check overall health | GET /v1/health |
| Server admin | Detect worker stall | Queue size endpoint with threshold |
| Server admin | Validate TLS | GET /v1/health/certificate for domain |
| SRE | Verify version | GET /v1/health/version (public) |

## 4. Business rules

1. **BR-01:** All queue-check endpoints share `threshold` param (default 5000); `size >= threshold` → `HEALTH_QUEUE_SIZE_EXCEEDED`. 🟢
2. **BR-02:** `/health/db` throws when any pool fails (current code `!empty($failures)`). 🟢
3. **BR-03:** `/health/time` connects to `time.google.com:123` via raw UDP, 5s receive timeout; rejects if drift > ±60 s. 🟢
4. **BR-04:** `/health/certificate` rejects expired cert (`$sslExpiration < time()`) → `HEALTH_CERTIFICATE_EXPIRED`. 🟢
5. **BR-05:** `/health/anti-virus` reads `_APP_STORAGE_ANTIVIRUS`; `disabled` → status=disabled without contacting ClamAV. 🟢
6. **BR-06:** `/health/geo` reads `_APP_GEO_ENDPOINT`; empty → `GENERAL_SERVER_ERROR`. 🟢
7. **BR-07:** `/health/stats` requires `root` scope; hidden from SDK docs (`docs: false`). 🟢
8. **BR-08:** `/health/version` uses `public` scope — no auth. 🟢
9. **BR-09:** `/health/queue/failed/:name` validates name via WhiteList from env vars (fallback `Event::*_QUEUE_NAME` constants). 🟢

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | HTTP-layer health (trivial pass) | Must | 200 pass with latency | 🟢 |
| RF-02 | Version (public) | Must | Stable version returned | 🟢 |
| RF-03 | DB pool ping | Must | All Console+Project pools pinged; failure → error | 🟢 |
| RF-04 | Cache ping | Must | Redis ping latency | 🟢 |
| RF-05 | PubSub ping | Must | All PubSub pools pinged | 🟢 |
| RF-06 | NTP time sync check | Must | UDP to time.google.com; ±60 s drift check | 🟢 |
| RF-07 | Certificate validation | Must | Fetch + validate TLS for public domain | 🟢 |
| RF-08 | Storage device test | Must | Write/read/delete on 4 devices | 🟢 |
| RF-09 | Local storage check | Must | Readability/writability of local volumes | 🟢 |
| RF-10 | AntiVirus check | Must | ClamAV ping; disabled flag respected | 🟢 |
| RF-11 | Geo service probe | Must | HTTP probe to `_APP_GEO_ENDPOINT` | 🟢 |
| RF-12 | Stats (root-only) | Must | Storage usage + Redis info | 🟢 |
| RF-13 | Queue size checks (13 queues) | Must | Threshold enforcement per queue | 🟢 |
| RF-14 | Failed-jobs count | Must | Named queue via WhiteList | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Observability | Latency in ms on all checks | Response schema | 🟢 |
| Reliability | Worker stall detection | Threshold + HEALTH_QUEUE_SIZE_EXCEEDED | 🟢 |
| Security | root scope on stats | label docs:false, scope root | 🟢 |
| Availability | Public version endpoint | scope public | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: Queue over threshold
  Given a queue with size >= threshold
  When GET /v1/health/queue/:name is called
  Then HEALTH_QUEUE_SIZE_EXCEEDED error is returned

Scenario: NTP drift
  Given server clock drifting more than 60s
  When GET /v1/health/time is called
  Then the check fails with drift error

Scenario: Certificate expired
  Given a domain with an expired TLS certificate
  When GET /v1/health/certificate is called
  Then HEALTH_CERTIFICATE_EXPIRED error is returned

Scenario: Antivirus disabled
  Given _APP_STORAGE_ANTIVIRUS=disabled
  When GET /v1/health/anti-virus is called
  Then status=disabled is returned without contacting ClamAV
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..06 (core probes) | Must | Baseline diagnostics |
| RF-07 (certificate) | Must | TLS validation |
| RF-08..09 (storage) | Must | Storage diagnostics |
| RF-10..11 (AV, geo) | Should | Optional infra |
| RF-12 (stats) | Should | Root-only |
| RF-13..14 (queues) | Must | Worker stall detection |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🔴 [GAP] Exact latency measurement method per check not fully verified.
- 🟡 [INFERRED] Stats response schema (storage usage + Redis info fields).

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
