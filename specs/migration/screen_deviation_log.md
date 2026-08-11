---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: screen_deviation_log
producedBy: screen-translator
hash: "sha256:9585b058438b0d7bcadeade849b159626481ae87e63d425197932a0407c81fee"
---

# Screen Deviation Log

> Mode = **literal/preservar**: nenhuma tela é reescrita, então não há deviations visuais. As entradas abaixo são **dependências não-REST que o Go deve garantir** para que a SPA preservada funcione — descobertas e confirmadas por red-team. Status `aprovado` = escopo aceito (não são questões abertas; são contratos que o Inspector cobre).

| DEV | Tela/área afetada | Tipo | Descrição | Status |
|-----|-------------------|------|-----------|--------|
| DEV-001 | todas as listagens | plataforma | schemas de resposta (`.total`×273, campos de modelo) não estão no OpenAPI | aprovado |
| DEV-002 | 529 arquivos c/ realtime | plataforma | protocolo WebSocket `/v1/realtime` fora do REST | aprovado |
| DEV-003 | command center / AI | correcao | assistant é **texto cru** (`streamProtocol:'text'`), spec dizia SSE (errado) | aprovado |
| DEV-004 | verify-email, billing, global | plataforma | error `type`/`code`/`message` dirige navegação | aprovado |
| DEV-005 | auth/sessão | plataforma | cookie `a_session_console` + `X-Fallback-Cookies` | aprovado |
| DEV-006 | impersonation, project layout | plataforma | headers `X-Appwrite-Mode/Impersonate/Console-Fingerprint/Organization` | aprovado |
| DEV-007 | console layout | tecnica | header `Date` p/ sync de relógio | aprovado |
| DEV-008 | analytics | tecnica | `sdk/sources.ts:36` raw `client.call` | aprovado |
| DEV-009 | multi-região | plataforma | topologia `<region>.<host>/v1` derivada client-side | aprovado |

> Nenhuma deviation **pendente** → não bloqueia o handoff ao Inspector. Todas propagam para `parity_specs.md § Exceções`.
