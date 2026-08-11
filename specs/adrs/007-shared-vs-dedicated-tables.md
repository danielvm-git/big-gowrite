# ADR-007 — Tabelas compartilhadas vs dedicadas por projeto

- **Status:** Aceito (transição em andamento)
- **Data:** Retrospectiva 2026-08-05
- **Confiança:** 🟢 CONFIRMADO (transição) / 🟡 (detalhes de pooling)

## Contexto
Cada projeto provisiona bancos de dados (data-plane) para documentsdb/vectorsdb. Existem duas estratégias: **tabelas/coleções compartilhadas** (mesmo cluster, namespace por projeto) e **bancos dedicados** por projeto. TODOs no código indicam que a migração para shared tables ainda está em curso.

## Decisão
- Seleção de DSN na criação do projeto via **pool de env vars** (`_APP_DB_*`), com filtro de região, seleção aleatória ou override explícito.
- Suporte a **shared tables** (namespace por projeto) com fallback/transição para dedicadas — TODOs: `Temporary until all projects are using shared tables` (Projects/Organization create).
- Bootstrap de coleções do projeto na criação.

## Alternativas consideradas
- Apenas bancos dedicados (rejeitado: custo de provisionamento por projeto).
- Apenas shared tables (rejeitado: isolamento/limites de tenant).

## Consequências
- Transição gradual mantém compatibilidade com projetos existentes.
- DSN por projeto persistido no doc do projeto (`database` map).
- VectorsDB (PostgreSQL/embeddings) segue fora desse fluxo de shared tables.

## Evidências
- TODOs em `Projects/Http/Projects/Create.php` e `Organization/Http/Projects/Create.php`.
- Módulo Databases — `Pool.php` (resolução de DSN); `domain.md` glossário.
