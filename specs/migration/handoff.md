---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: handoff
producedBy: migration-orchestrator
hash: "sha256:21f218474a40fa83fe7ed89ac49d3d79743aeb7a86fc085068381cd694bbe075"
---

# Migration Handoff — gowrite (Appwrite PHP → Go)

> Entregável final do Time de Migração. O agente de codificação implementa a partir daqui.
> **Leitura obrigatória primeiro**: `paradigm_decision.md` (como pensar) → `topology_decision.md` (como organizar). Depois `redteam_findings.md` (o que quase deu errado e foi corrigido).

## Decisões-âncora (todas validadas por red-team)
| Dimensão | Decisão |
|----------|---------|
| **Paradigma** | Transformacional / **Go-native** (erros explícitos, composição, goroutines, schema tipado) |
| **Banco** | **PostgreSQL único** (relacional + JSONB + pgvector); **sem Redis** |
| **Estrutura** | Pacotes **flat** idiomáticos (sem hexagonal ports/adapters); `platform/` = leaf packages |
| **Dados** | **`sqlc`** + SQL à mão + `crypto.Encrypt()` explícito (não reimplementar engine utopia-php) |
| **Dispatch** | Tabela **`jobs`** + `SELECT … FOR UPDATE SKIP LOCKED`; realtime via `LISTEN/NOTIFY` |
| **Topologia** | **11 bounded contexts + Console BFF** (`tech-stack.md`) |
| **Estratégia** | **Strangler Fig + Parallel Run**, flip por **contexto**, ordem **topológica** |
| **Console** | **Preservar a SPA** SvelteKit; Go serve via `go:embed`; paridade = shape+protocolo |

## Ordem de implementação (bottom-up, topológica)
1. `platform/` (errors, auth, permissions, crypto, ratelimit) + schema Postgres base + tabela `jobs`.
2. **Identity + Organizations** (co-flip atômico).
3. **Projects** (inclui Webhooks).
4. **Data**, **Storage**.
5. **VCS** → **Compute** (VCS publica evento via `jobs`; Compute consome).
6. **Networking**, **Messaging**, **Migrations**, **Advisor**.
7. **Console BFF** + servir a SPA (ver requisitos de hosting).
8. Workers cross-cutting (Deletes) — por último.

## Artefatos produzidos
**Discovery (`specs/`)**: code-analysis, domain, state-machines (8+2 SMs), permissions, architecture, C4, ERD, `database/` (48 coleções/606 atributos), adrs/ (8), traceability, user-stories, confidence-report.
**Migração (`specs/migration/`)**:
- `migration_brief.md`, `paradigm_decision.md`, `target_business_rules.md`, `discard_log.md`
- `migration_strategy.md`, `risk_register.md`, `cutover_plan.md`
- `topology_decision.md`, `target_architecture.md`, `target_domain_model.md`, `target_data_model.md`, `data_migration_plan.md`
- `screen_modernization_decision.md`, `target_screens.md`, `screen_deviation_log.md`
- `parity_specs.md`, `parity_tests/*.feature`
- `redteam_findings.md`, `ambiguity_log.md`
**Domínio (`specs/tech-architecture/tech-stack.md`)**: mapa de 11 contextos + invariantes.

## Itens REFERIDOS À CODIFICAÇÃO
- **Requisitos de hosting da SPA** (`screen_modernization_decision.md` §hosting): base `/console`, `try_files` fallback, br/gz, cache tiers, endpoint via `window.location`.
- **9 deviations não-REST** (`screen_deviation_log.md` DEV-001..009): shape de resposta, WS realtime, assistant texto-cru, catálogo de erro, cookie-session, headers console-only, `Date`, raw source call, multi-região.
- **Re-cripto dos 25 campos** na migração de dados (R-D3).
- **BR-HUMANA-003**: executor/browser/assistant permanecem serviços externos na fase 1.
- **Gap de discovery**: `specs/openapi/appwrite.yaml` sem response schemas → gerar goldens de resposta a partir de `specs/database/` + captura do legado.

## Próximos passos para o agente de codificação
1. `git init` do repo Go novo; aplicar **seed-conventions** do bigpowers (CLAUDE.md/CONVENTIONS.md Go).
2. Provisionar Postgres (pgvector); rodar `platform/` + migrations base + tabela `jobs`.
3. Implementar bottom-up (ordem acima), um contexto por vez, **flat + sqlc**.
4. Por contexto: Parallel Run → gate de paridade (`parity_specs.md`) → flip no Traefik → canary → decomissão do PHP.
5. Validar paridade contra **goldens de resposta + protocolo WS + framing + catálogo de erro**, não contra o OpenAPI oco.
6. Servir a SPA preservada com os 5 requisitos de hosting.

## Estado
Pipeline de migração **concluído** (6 agentes + 2 rodadas de red-team). Ambiguidades pendentes: **0**.
