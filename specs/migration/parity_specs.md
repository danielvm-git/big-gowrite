---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: parity_specs
producedBy: inspector
hash: "sha256:45fc208d913d2b97c2f62aa49f3a282f03e775674c89b60ab2cf724d65624661"
---

# Parity Specs — provar equivalência gowrite ≡ Appwrite legado

> Specs de paridade (não testes executáveis). O agente de codificação traduz para o framework Go/Gherkin.
> Paradigma: transformacional (Go-native) com bordas event-driven → cobertura adaptada obrigatória.

## Estratégia de validação (modos selecionados)
- ✅ **Parallel Run (shadow online)** — Go e PHP recebem a mesma requisição espelhada por contexto; comparação de resposta. É o gate de flip (`migration_strategy.md`).
- ✅ **Characterization / response goldens** — capturar respostas reais do legado por operação e comparar **shape** (corrige o oráculo oco: `appwrite.yaml` não tem response schemas). Base de modelos: `specs/database/` (48 coleções/606 atributos).
- ✅ **Contract tests** — WS realtime, framing do assistant, catálogo de erros (fora do REST).
- ✅ **Data parity** — snapshots + checksums na migração Mongo/MariaDB → Postgres.

## Critérios de "paridade aceita"
- **Métrica primária**: índice de divergência funcional **< 0,01%** das requisições espelhadas em janela de **30 dias** por contexto.
- **Janela de observação**: mínimo 3 dias de canary pós-flip por contexto (`cutover_plan.md`).
- **Bloqueio de cutover**: qualquer divergência de **segurança** (auth/permissões/cripto) ou de **shape de resposta** em campo consumido pelo console = no-go.
- **Campos de divergência aceitável** (Parallel Run): timestamps de geração, ordering não-determinístico de listas sem `orderBy` explícito, ids internos (`internalId`) — normalizados antes da comparação.

## Cobertura adaptada ao paradigma (transformacional + event-driven)
| Transição | Dimensões obrigatórias |
|---|---|
| OO-com-DI síncrono → Go-native | equivalência funcional + mapeamento erro→status HTTP (exceções viraram `AppError`) |
| filas Redis → `jobs` Postgres (SKIP LOCKED) | **ordem**, **idempotência** (por `dedup_key`/CloudEvent id), consistência eventual, comportamento sob falha de worker, **estado terminal obrigatório** (ADR-008) |
| Mongo/MariaDB → Postgres | data parity (contagens, checksums, integridade referencial, re-cripto dos 25 campos) |

Cada fluxo com efeito assíncrono gera ≥3 cenários: `@paridade`, `@idempotencia`, `@ordem`.

## § Exceções (deviations do Screen Translator propagadas)
Paridade de UI = paridade de **shape de resposta + protocolo**, não de tela (mode literal/preservar). As deviations `DEV-001..009` de `screen_deviation_log.md` são exceções que exigem testes dedicados:
- DEV-001 shape de resposta (`.total`×273) · DEV-002 WS realtime · DEV-003 assistant texto-cru (não SSE) · DEV-004 catálogo de erro `type`/`code` · DEV-005 cookie-session + `X-Fallback-Cookies` · DEV-006 headers console-only · DEV-007 header `Date` · DEV-008 raw source call · DEV-009 topologia multi-região.

## Fluxos críticos cobertos (`parity_tests/`)
| # | Flow | Tags | BR / origem |
|---|------|------|-------------|
| 01 | Auth + sessão + MFA | @paridade @critico | BR-MIGRAR-001, DEV-005 |
| 02 | Convite de time (cross-context write) | @paridade @critico | S2, BR-MIGRAR-012 |
| 03 | Ciclo de build/deployment | @paridade @idempotencia @ordem | BR-MIGRAR-002, ADR-008 |
| 04 | Realtime subscription (WS) | @paridade @critico | DEV-002 |
| 05 | Transação com lock otimista | @paridade @critico | BR-MIGRAR-004, ADR-005 |
| 06 | Shape de resposta de listagem | @paridade @critico | DEV-001 |
| 07 | Catálogo de erro dirigindo navegação | @paridade | DEV-004 |

## Lacunas 🔴
- `characterization_specs/` não existe → cenários derivados de `code-analysis.md` + `state-machines.md`. Recomenda-se capturar goldens de resposta do legado antes do primeiro flip.
- OpenAPI de discovery incompleto (sem response schemas) — o Inspector supre via `specs/database/` + captura.
