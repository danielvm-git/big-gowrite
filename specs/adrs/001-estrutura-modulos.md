# ADR-001 — Estrutura de módulos (`Modules/{Name}`)

- **Status:** Aceito (em andamento — migração incremental)
- **Data:** Retrospectiva 2026-08-05
- **Confiança:** 🟢 CONFIRMADO

## Contexto
O servidor evoluiu de um monólito organizado por serviço (API + workers no mesmo namespace) para features cada vez maiores. Commits recentes migram serviços inteiros para a estrutura de módulos (`refactor(users): migrate Users API to module structure`). A migração é incremental e convive com código legado fora de módulos.

## Decisão
Organizar cada feature em `src/Appwrite/Platform/Modules/{Name}/` com:
- `Module.php` — registra serviços (Http/Workers/Tasks)
- `Services/Http.php` — ações HTTP via padrão `Action` (getName, setHttpMethod, setHttpPath, label scope/event, param, inject)
- `Services/Workers.php` / `Services/Tasks.php` — workers e tarefas CLI
- `Http/{Service}/{Action}.php` — endpoints aninhados espelhando a URL

## Alternativas consideradas
- Manter o monólito por serviço (rejeitado: acoplamento crescente).
- Quebrar em microsserviços de rede (rejeitado: custo operacional vs ganho — plataforma já roda workers de fila).

## Consequências
- Convivência temporária com código legado; migração por módulo.
- Mudanças no padrão Action/DI afetam `appwrite/cloud` (cross-repo).
- `resourceType` sempre plural; ações CRUD-only; operações não-CRUD viram propriedades (`Memberships/Status/Update.php`).

## Evidências
- `git log`: `refactor(users): migrate Users API to module structure` (HEAD~1).
- 22 módulos em `src/Appwrite/Platform/Modules/`; estrutura verificada via `find`.
