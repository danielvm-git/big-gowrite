---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: screen_modernization_decision
producedBy: screen-translator
mode: literal
hash: "sha256:5560092ac17b494ceb6fe48cbfbd4c10f55d42b548a41084d4b2c84936a7aea4"
---

# Screen Modernization Decision (revised post red-team)

> Decisão sobre o modo de tradução das telas do console. Resolve BR-HUMANA-004 (console-in-Go).
> **Mode = literal/preservar** — validado por 2 red-teams. Mas a *definição de paridade* e o *escopo do que o Go precisa garantir* foram corrigidos: o red-team provou que "paridade ≡ contrato OpenAPI" era **falso**.

## Plataforma de origem detectada
- **Console SPA**: SvelteKit 2 + Svelte 5, **build estático** (`@sveltejs/adapter-static`, `fallback: index.html`, `ssr=false` no root layout — zero rotas/hooks server-side). SDK `@appwrite.io/console` + design system `pink`. 🟢

## Plataforma alvo confirmada
- **A mesma SPA**, servida pelo binário Go via `go:embed` + `net/http`. **Sem runtime Node** (build é 100% estático). Resolve BR-HUMANA-004 = **"Go serves existing SPA"**. Verdict de hosting: **FEASIBLE-WITH-CHANGES**.

## Decisão de modo
- **Modo**: **literal** (preservar a SPA; não reescrever front-end). Rewrite (WASM/SSR) viola paridade de UI + footprint sem ganho. **Sem specs de tela geradas.**
- Princípios bigpowers aplicados: **evidence-over-claims** (626 usos de SDK), **YAGNI/minimal-change** (não reescrever front-end que já satisfaz o requisito), **research-first** (investigou-se o acoplamento real).
- **Decidido em**: 2026-08-11.

## ⚠️ Correção crítica do oráculo de paridade (red-team CRITICAL-0)
A versão anterior afirmava **"paridade de UI ≡ paridade de contrato OpenAPI"**. **FALSO e verificado como falso:**
- `specs/openapi/appwrite.yaml` tem **0 blocos `responses:`** e **0 `requestBody:`** — define **rotas/params/security apenas, não payloads**. Um backend Go pode "honrar o OpenAPI" e retornar qualquer JSON.
- O console consome shape de resposta massivamente: `.total` lido **273×**, além de `$id`/`$createdAt` e todos os campos de modelo.
- **Logo**: conformidade de contrato ≠ conformidade de payload. **O oráculo correto é paridade de SHAPE de resposta + protocolo**, não o OpenAPI atual.

## O que o Go DEVE garantir para paridade de UI (além de rotas)
| DEV | Dependência | Fora do OpenAPI? | O que o Go deve garantir |
|-----|-------------|------------------|--------------------------|
| DEV-001 | **Schemas de resposta** (`.total`, campos de modelo, `$id`/`$createdAt`, listas) | sim (0 responses no spec) | shape byte-compatível por operação — derivar de `specs/database/` + respostas capturadas (golden) |
| DEV-002 | **Realtime WebSocket** (`/v1/realtime`) | sim | upgrade WS, grammar `project`+`channels[]`, auth-over-WS, envelope `{events,channels,timestamp,payload}`, strings de evento (`databases.*.collections.*`). **529 arquivos** tocam realtime |
| DEV-003 | **Assistant streaming** | spec diz "SSE" — **ERRADO** | `ai.svelte` usa `streamProtocol:'text'` (ai-sdk) = **stream de texto cru concatenado, NÃO SSE**. Seguir o spec literalmente **quebra** a UI. Go deve emitir texto cru, não `data:` frames |
| DEV-004 | **Error `type`/`code`/`message`** controla navegação | sim | catálogo exato de `type` strings (`user_email_not_verified`, `console_account_verification_required`, `billing_coupon_already_used`…), `code` numérico e substrings de `message` (`emailVerification.ts:7`, `hooks.client.ts:21`) |
| DEV-005 | **Cookie-session auth** | sim | `Set-Cookie a_session_console` (name/domain/SameSite) + padrão `X-Fallback-Cookies` (localStorage) para contexto sem 3rd-party cookies; `credentials:'include'` |
| DEV-006 | **Headers console-only** | sim | `X-Appwrite-Mode: admin` + `X-Appwrite-Impersonate-*` (impersonation), `X-Appwrite-Console-Fingerprint`, `X-Appwrite-Organization` — cada um com comportamento server-side |
| DEV-007 | Header `Date` (sync de relógio) | sim | `+layout.ts:36` lê `response.headers.get('Date')` → `syncServerTime()`; Go deve emitir `Date` correto |
| DEV-008 | `sdk/sources.ts:36` `client.call` raw (analytics) | sim | validar rota/contrato exato |
| DEV-009 | Topologia multi-região (`<region>.<host>/v1`) | sim | console deriva WS/endpoint da topologia de subdomínio; decisão de deployment |

## Requisitos de hosting (go:embed — red-team FEASIBLE-WITH-CHANGES)
Um `http.FileServer` ingênuo **quebra**. O handler Go deve:
1. Servir sob **base path `/console`** (compilado no build, `svelte.config.js` `paths.base`) + `301 / → /console`.
2. **SPA fallback** estilo `try_files $uri /console/index.html` (HTTP 200, não 404) para deep-links.
3. Negociar **`Accept-Encoding` br/gz** (build usa `precompress:true`) — ou dropar precompress e gzip on-the-fly.
4. Replicar as **3 tiers de cache-control** (immutable 1-ano para `/_app/immutable/`, `no-store` para HTML).
5. **Correção de env**: `PUBLIC_APPWRITE_ENDPOINT` é **build-time**, NÃO Go-runtime. Deixar **unset** → o SDK deriva o endpoint de `window.location` (`sdk.ts:43`) — ideal para same-origin Go. Demais `PUBLIC_*` (CONSOLE_MODE, feature flags) ficam congelados no build do console (aceitável self-hosted).
- **Same-origin** (Go serve `/console` + `/v1`) **simplifica** cookies/CORS (first-party, sem preflight).

## Implicações para o Inspector (oráculo corrigido)
Paridade de UI **NÃO** é medida contra o `appwrite.yaml` atual (é oco). Medir contra:
1. **Golden de respostas**: capturar respostas reais do Appwrite legado por operação; comparar shape (o Data Master já mapeou 48 coleções/606 atributos → base dos modelos).
2. **Protocolo WS realtime** (DEV-002) — testes dedicados fora de REST.
3. **Framing do assistant** (DEV-003, texto cru).
4. **Catálogo de error `type`/`code`** (DEV-004).
5. **Cookie-session + headers console-only** (DEV-005/006).
6. Suíte E2E do próprio console apontada ao Go API (same-origin) como teste de integração final.

> **Gap de discovery detectado**: `specs/openapi/appwrite.yaml` está incompleto (sem response schemas). Anotado para o handoff — não bloqueia, pois `specs/database/` supre os modelos, mas o Inspector deve gerar goldens de resposta.
