---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: target_screens
producedBy: screen-translator
mode: literal
hash: "sha256:8dd0c44162560a8ad5f5c7053338a01899be5a95482cfaa8a50765550a5a0f52"
---

# Target Screens

> **Mode = literal/preservar.** Nenhuma tela é traduzida ou reescrita — o console SPA (SvelteKit 2 / Svelte 5) é **preservado como está** e servido pelo binário Go (`go:embed`). Ver `screen_modernization_decision.md`.

## Por que não há specs de tela
- O front-end não muda: mesma UI, mesmo look-and-feel, mesmos componentes (`pink`), mesmos fluxos.
- O trabalho de migração é **backend** (Go) + **hosting** (Go serve a SPA) + **paridade de contrato de resposta/protocolo** (não de tela).
- Reescrever telas violaria o requisito de paridade de UI do brief e o objetivo de footprint.

## O que substitui as specs de tela
- **`screen_modernization_decision.md`** — decisão + oráculo de paridade corrigido + requisitos de hosting.
- **`screen_deviation_log.md`** — as 9 dependências não-REST que o Go deve garantir (DEV-001..009).
- **Inspector (`parity_specs.md`)** — testes de paridade contra respostas golden + protocolo WS + framing do assistant + catálogo de erros + suíte E2E do console apontada ao Go.

## Fonte da UI preservada
`/Users/danielvm/Developer/appwrite-app/appwrite-console` (repo separado; read-only nesta migração).
