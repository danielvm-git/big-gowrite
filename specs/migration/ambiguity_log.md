---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: ambiguity_log
producedBy: migration-orchestrator
hash: "sha256:92cf9034b0ee6657b98a64bbca564c309e2e5f1d8c7c927dac1cd832a94dbcc9"
---

# Ambiguity Log

> Itens que precisam de decisão humana ou ficam referidos à codificação. Consolidado ao final (após Inspector) em três grupos.

## PENDENTES (decisão humana)

_Nenhum item pendente._ (BR-HUMANA-004 resolvido — ver abaixo.)

## RESOLVIDOS COM DECISÃO HUMANA

| ID | Tema | Decisão | Data |
|----|------|---------|------|
| — | Paradigma | **Transformacional / Go-native** | 2026-08-11 |
| BR-HUMANA-002 | Estratégia de banco | **Consolidar em PostgreSQL** (relacional + JSONB p/ documentos + pgvector p/ VectorsDB) | 2026-08-11 |
| BR-HUMANA-004 | Console | **Preservar a SPA** (mode literal); Go serve via `go:embed` (FEASIBLE-WITH-CHANGES); paridade = shape de resposta + protocolo | 2026-08-11 |

## RESOLVIDOS PELO RED-TEAM (2026-08-11)

| ID | Tema | Decisão | Fonte |
|----|------|---------|-------|
| BR-HUMANA-001 | Camada de dados | **`sqlc` + SQL à mão + `crypto.Encrypt()` explícito**; NÃO reimplementar engine `utopia-php/database` | red-team O4, confirmado pelo usuário |
| O1 | Estrutura de código | **Pacotes flat**, sem hexagonal ports/adapters | confirmado pelo usuário |
| O3/O5 | Infra/dispatch | **Postgres-only**: `jobs`+`SKIP LOCKED`, `LISTEN/NOTIFY`; sem Redis, sem outbox | confirmado pelo usuário |

## RESOLVIDOS COM RECOMENDAÇÃO (default provisório — confirmável no agente indicado)

| ID | Tema | Default aplicado | Confirmar em |
|----|------|------------------|--------------|
| BR-HUMANA-003 | Executor / browser / assistant | Manter como **serviços externos na fase 1** | Strategist (escopo) / Inspector (paridade) |

## REFERIDOS À CODIFICAÇÃO
- (a preencher pelos agentes posteriores)
