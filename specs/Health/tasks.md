# Tasks: Health

> Unit: `Health`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## Task list

### T-HLT-01 — Module bootstrap + route registration
- **Behavior source:** `Modules/Health/Module.php`, `Services/Http.php`
- **Description:** Register HTTP service with 25 endpoints and scopes (health.read, public, root).
- **Done when:** Route table matches mapping.
- **Priority:** Must · Confidence: 🟢

### T-HLT-02 — Core probes (HTTP, version, db, cache, pubsub)
- **Behavior source:** `Http/Health/{Get,Version/Get,DB/Get,Cache/Get,PubSub/Get}.php`
- **Description:** Implement trivial pass, public version, DB pool ping, Redis ping, PubSub ping with latency.
- **Done when:** All probes respond with status + ping.
- **Priority:** Must · Confidence: 🟢

### T-HLT-03 — Time sync (NTP)
- **Behavior source:** `Http/Health/Time/Get.php`
- **Description:** UDP to time.google.com:123, 5s timeout, ±60 s drift assertion.
- **Done when:** Drift beyond bounds → error.
- **Priority:** Must · Confidence: 🟢

### T-HLT-04 — Certificate validation
- **Behavior source:** `Http/Health/Certificate/Get.php`
- **Description:** Fetch + validate TLS cert; expired → HEALTH_CERTIFICATE_EXPIRED.
- **Done when:** Expired cert detected.
- **Priority:** Must · Confidence: 🟢

### T-HLT-05 — Storage checks (devices + local)
- **Behavior source:** `Http/Health/Storage/{Get,Local/Get}.php`
- **Description:** Write/read/delete probe on 4 devices; local volume readability/writability.
- **Done when:** All storage probes work.
- **Priority:** Must · Confidence: 🟢

### T-HLT-06 — AntiVirus + Geo
- **Behavior source:** `Http/Health/{AntiVirus,Geo}/Get.php`
- **Description:** ClamAV ping with disabled flag; geo probe with empty-endpoint error.
- **Done when:** disabled status honored; geo errors on missing endpoint.
- **Priority:** Should · Confidence: 🟢

### T-HLT-07 — Stats (root scope)
- **Behavior source:** `Http/Health/Stats/Get.php`
- **Description:** Storage usage + Redis info; docs:false; scope root.
- **Done when:** Stats returned for root; hidden from SDK docs.
- **Priority:** Should · Confidence: 🟢

### T-HLT-08 — Queue checks (13) + failed jobs
- **Behavior source:** `Http/Health/Queue/**`
- **Description:** Threshold enforcement via base helper; failed jobs by named queue WhiteList.
- **Done when:** Thresholds enforced; failed count returned.
- **Priority:** Must · Confidence: 🟢

### T-HLT-09 — Contract / OpenAPI surface
- **Behavior source:** All `Modules/Health/Http/**`
- **Description:** OpenAPI 3.1 for health endpoints.
- **Done when:** `openapi/health.yaml` matches actions.
- **Priority:** Should · Confidence: 🟡

### T-HLT-10 — E2E test suite
- **Behavior source:** `tests/e2e/Services/Health/**`
- **Description:** Cover all probes + queue thresholds + failed jobs + version public access.
- **Done when:** E2E suite passes.
- **Priority:** Must · Confidence: 🟢
