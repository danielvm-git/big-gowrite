# ADR-005 — Transações com lock otimista (Databases)

- **Status:** Aceito
- **Data:** Retrospectiva 2026-08-05
- **Confiança:** 🟢 CONFIRMADO

## Contexto
O Appwrite precisa oferecer transações multi-documento sobre backends heterogêneos (MongoDB, PostgreSQL, MariaDB/MySQL) mantidos como camada de dados plana (`utopia-php/database`). Transações nativas têm semânticas e suporte diferentes por backend.

## Decisão
- Sistema de transações próprio com **lock otimista**: operações são encadeadas em staging (coleção `transactions`, campo `version`), validadas por cheque de versão e aplicadas em **commit único atômico**.
- Rollback via `publisherForDeletes` (worker de deletes) quando o commit falha.
- Trilha de execução em `transactionLogs`.

## Alternativas consideradas
- Transações nativas de cada backend (rejeitado: inconsistência de suporte entre adaptadores).
- Lock pessimista global (rejeitado: gargalo de concorrência).

## Consequências
- Consistência uniforme multi-backend com custo de re-tentativa por conflito de versão.
- Complexidade no cliente (versionamento de documentos).
- Integração com o worker `Databases` para schema assíncrono.

## Evidências
- Módulo Databases (`transactions`/`transactionLogs`), `domain.md` §2.5.
- `modules.json` — Databases algorithms: `Optimistic-lock transactions`.
