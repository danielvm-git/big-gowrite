# ADR-008 — Estado terminal obrigatório para jobs de fila

- **Status:** Aceito
- **Data:** Retrospectiva 2026-08-05
- **Confiança:** 🟢 CONFIRMADO

## Contexto
Mensagens de fila processadas por workers podem falhar no meio do processamento e ficar **presas em `processing`** para sempre, sem estado terminal — corrompendo contadores, retries e o estado de recursos dependentes.

## Decisão
- Todo worker/mensagem deve terminar em **estado terminal** (completed/failed/cancelled), inclusive quando o cleanup falha.
- Fix aplicado ao messaging: "stop messages getting stranded in processing" + "keep a terminal status when cleanup fails".
- Dedup de callbacks por **CloudEvent id** no padrão join (Builds/Jobs) para que múltiplos eventos não dupliquem transições.

## Alternativas consideradas
- Retry infinito sem estado terminal (rejeitado: stranded messages).
- Reaproveitar estados intermediários como finais (rejeitado: quebra contadores/retries).

## Consequências
- Resiliência: filas nunca ficam com mensagens órfãs em processamento.
- Padrão replicado nos workers (Builds/Jobs, Databases, Migrations, Messaging).

## Evidências
- Commits: `fix(messaging): stop messages getting stranded in processing`, `fix(messaging): keep a terminal status when cleanup fails`.
- `state-machines.md` §8 (regra transversal de estado).
