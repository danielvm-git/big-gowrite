# ADR-006 — Deprecação das Dev Keys (DEV_KEY_GONE)

- **Status:** Aceito (deprecated)
- **Data:** Retrospectiva 2026-08-05
- **Confiança:** 🟢 CONFIRMADO

## Contexto
Dev keys eram credenciais de projeto que "bypassavam rate limits e melhoravam logs de erro durante o desenvolvimento". Tornaram-se vetor de risco e perderam o propósito com o amadurecimento do sistema de keys.

## Decisão
- **Deprecar a criação** de dev keys: rota de create lança `DEV_KEY_GONE`.
- Keys legadas viram **tombstones** (`devKeys` com flag de tombstone) — read/list ainda respondem com o marcador para compatibilidade de clientes.
- Testes do fluxo de criação removidos (`chore-deprecate-devkeys`).

## Alternativas consideradas
- Remoção total imediata (rejeitado: quebraria clientes existentes que ainda listam).
- Manter funcional (rejeitado: risco de abuso de rate limits).

## Consequências
- Compatibilidade: GET/LIST continuam retornando dados legados com marcador.
- Documentação/SDK atualizada (create/update/delete não funcionais).
- Escopos `devKeys.read/write` ainda presentes nos labels (legado).

## Evidências
- Commits: `chore-deprecate-devkeys`, `Deprecate dev keys creation + remove tests`.
- Módulo Projects — `Http/DevKeys/*`; `domain.md` glossário (Dev Key).
