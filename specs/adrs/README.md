# ADRs — Registro de Decisões Arquiteturais (retroativas)

> Gerado pelo **Detective** em 2026-08-05 a partir da arqueologia Git (37.244 commits) e inspeção de código.
> Confiança: todos 🟢 CONFIRMADO (com exceções anotadas em cada ADR).

| # | ADR | Decisão central |
|---|-----|-----------------|
| 001 | [estrutura-modulos](001-estrutura-modulos.md) | Migração incremental para `Modules/{Name}` com padrão Action |
| 002 | [hashing-senha](002-hashing-senha.md) | Argon2id default + imports legados (Bcrypt/MD5/Argon2/SHA/PHPass/Scrypt) |
| 003 | [pipeline-builds-async](003-pipeline-builds-async.md) | Builds + Jobs workers com dedup por CloudEvent id |
| 004 | [upload-chunks-locks](004-upload-chunks-locks.md) | Upload multipart (Content-Range) + distributed locks |
| 005 | [transacoes-lock-otimista](005-transacoes-lock-otimista.md) | Transações com lock otimista multi-backend |
| 006 | [deprecacao-dev-keys](006-deprecacao-dev-keys.md) | Dev keys deprecated (tombstone `DEV_KEY_GONE`) |
| 007 | [shared-vs-dedicated-tables](007-shared-vs-dedicated-tables.md) | DSN pool com shared/dedicated tables (transição) |
| 008 | [estado-terminal-filas](008-estado-terminal-filas.md) | Estado terminal obrigatório para jobs de fila |
