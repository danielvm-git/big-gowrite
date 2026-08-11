# ADR-003 — Pipeline assíncrono de builds (Builds + Jobs workers)

- **Status:** Aceito
- **Data:** Retrospectiva 2026-08-05
- **Confiança:** 🟢 CONFIRMADO

## Contexto
Builds de functions/sites (upload, VCS, template) são longos e assíncronos. É preciso orquestrar executor + build + ativação sem travar a API e com resiliência a falhas parciais.

## Decisão
- **Builds worker** executa o build (executor) e dispara callbacks (`exit(0)`, `complete(+manifest)`, etc.).
- **Jobs worker** sincroniza o estado: seta deployment `building` (estampa início), aplica resultado `ready`/`failed` de forma **deduplicada por CloudEvent id** (padrão join).
- Builds conect-time (VCS update) roteados por fila de jobs, não inline na API.
- Template-into-repo push fica no Builds worker e faz hand-off ao Jobs após.

## Alternativas consideradas
- Build síncrono na API (rejeitado: bloqueia request por minutos).
- Callback único sem dedup (rejeitado: builds geram múltiplos eventos → risco de dupla aplicação).

## Consequências
- Estados terminais sempre aplicados (inclusive em falha de cleanup — lição do fix de messaging).
- Erros internos de build **mascarados** nos logs visíveis ao usuário.
- Cancelamento de build bloqueado quando deployment já `ready`/`failed` (`BUILD_ALREADY_COMPLETED`).

## Evidências
- `Workers/Builds.php`, `Workers/Jobs.php` (Functions).
- Commits: `Route the VCS connect-time build (Functions/Update) through jobs`, `Keep template-into-repo push in the Builds worker; hand off to jobs after`, `fix(builds): mask internal errors`.
