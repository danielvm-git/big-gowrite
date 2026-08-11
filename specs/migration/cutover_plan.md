---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: cutover_plan
producedBy: strategist
hash: "sha256:433559d92bb261afe62e91bc7f63a89a48a15ad30a327585393f5533339e8424"
---

# Cutover Plan

> Para a estratégia recomendada (Strangler Fig + Parallel Run). **Não há janela única** — o cutover é **incremental, por módulo**. Cada módulo repete o ciclo abaixo.

## Pré-requisitos globais (uma vez)
- Repositório Go novo criado (seed-conventions do bigpowers aplicado).
- PostgreSQL provisionado (relacional + JSONB + pgvector).
- Traefik configurado como roteador de caminhos PHP↔Go.
- ETL base validado em staging (schema das 48 coleções em Postgres).
- Chave-mestra AES migrada e decrypt validado em Go (R-D3).
- Suite de testes de paridade do Inspector (`parity_tests/*.feature`) executável em CI.

## Ciclo por módulo (repetir na ordem bottom-up do `migration_strategy.md`)

| # | Passo | Owner | Duração típica | Go/No-Go |
|---|-------|-------|----------------|----------|
| 1 | Implementar módulo em Go (bottom-up: dados→lógica→HTTP) | dev | 1–N sessões (cota de IA) | compila + unit tests verdes |
| 2 | Migrar/rechecar dados do módulo para Postgres (se tem coleção própria) | dev | curta | ETL do módulo valida contagem + amostra |
| 3 | **Parallel Run**: espelhar tráfego real ao Go, comparar respostas com PHP | dev | dias de observação | divergência ≤ limiar nos `parity_tests` do módulo |
| 4 | **Flip**: Traefik roteia o caminho do módulo para o Go | dev | minutos | testes de contrato (OpenAPI) verdes pós-flip |
| 5 | Observar em produção (canary) | dev | 1–3 dias | sem regressão; métricas de recurso ≤ PHP |
| 6 | Decomissionar o handler PHP do módulo | dev | curta | módulo 100% em Go |

## Critérios go/no-go por flip (passo 4)
> **Revisão red-team**: a unidade de flip é o **bounded context** (não o módulo). Identity+Organizations flipam juntos; VCS antes/junto de Compute.
- ✅ **Todos os contextos referenciados já migrados** (S4) — ex.: não flipar Compute antes de VCS (`installations`), nem Projects antes de Organizations (`teamInternalId`).
- ✅ Todos os `parity_tests` do contexto passam (100% p/ contrato de UI, conforme brief).
- ✅ Paridade de contrato de API vs. `specs/openapi/appwrite.yaml`.
- ✅ Dados do contexto reconciliados; **nenhuma escrita cross-context para coleção ainda em PHP** (S2) — se existir fan-out (ex.: convite de time escreve `users`+`sessions`), o contexto-alvo do fan-out já deve estar migrado.
- ❌ Qualquer divergência de segurança (auth/permissões/cripto) = no-go imediato.
- ❌ Worker cross-cutting (Deletes) que faça cascade para contexto ainda-em-PHP = no-go até esse contexto migrar (S3).

## Rollback (por módulo)
- **Trigger**: regressão detectada pós-flip.
- **Ação**: Traefik re-roteia o caminho do módulo de volta ao PHP (rota é a chave de reversão). Dado permanece consistente porque o ownership do dado não foi transferido até o passo 6.
- **Tempo de rollback**: minutos (mudança de rota).

## Marco final
- Todos os 22 módulos flipados → decomissionar o stack PHP/Swoole inteiro → medir footprint final vs. baseline (métrica de sucesso do brief).
