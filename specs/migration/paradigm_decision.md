---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: paradigm_decision
producedBy: paradigm_advisor
hash: "sha256:2738974d5e0781f939988528d44ae36cbfdc8a35d946099fa3b787f87b09b0da"
---

# Paradigm Decision

> Decisão consciente sobre como tratar a mudança de paradigma entre o legado (PHP/Swoole) e a stack alvo (Go).
> Leitura obrigatória primeiro para todos os agentes posteriores e para o agente de codificação.

## Paradigma do legado detectado
- **Paradigma principal**: **híbrido — OO-com-DI (núcleo) + event-driven (bordas assíncronas)**
- **Confiança**: 🟢 CONFIRMADO
- **Evidências**:
  - OO-com-DI: padrão Action (classes estendem `Appwrite\Platform\Action`), injeção via `->inject('dbForProject')` / `->inject('queueForEvents')`, registro de serviços em `Module.php`, repositórios `dbForProject`/`dbForPlatform` — `AGENTS.md` (Action pattern), `code-analysis.md §2.2, §2.5`. 🟢
  - Event-driven: `queueForEvents`, 12 workers consumindo filas Redis (`BuildV1`, `DeleteV1`…), publishers, realtime pub/sub, schema assíncrono, máquinas de estado com estado terminal — `code-analysis.md §2.4, §5`, `state-machines.md`, ADR-003 (build pipeline async), ADR-008 (terminal state em filas). 🟢
- **Variações observadas** (híbrido):
  - Núcleo HTTP (22 módulos): OO-com-DI síncrono. 🟢
  - Workers / filas / realtime: event-driven assíncrono com idempotência e estados terminais. 🟢
  - Concorrência in-process: corrotinas Swoole (ex.: análise concorrente de repos VCS — `code-analysis.md §3.21`). 🟢

## Stack alvo declarada
- Linguagem: **Go** (decidida pelo usuário)
- Framework: a definir (Designer, sob princípios bigpowers)
- Infra: VPS única self-hosted, foco em **baixo consumo de recursos**

## Paradigma natural inferido
- **Paradigma**: **CSP / goroutines (event-driven leve) + composição por interfaces**
- **Justificativa**: Go não tem classes nem herança; comportamento compartilhado se dá por **composição + interfaces (structural typing)**; erros são **valores retornados explícitos** (sem exceções); concorrência idiomática é **goroutines + channels** (CSP). Ver catálogo `paradigm-catalog.md` (linha Go).
- **Alternativas viáveis**: procedural estruturado (mais simples ainda); OO "simulada" via interfaces + um container de DI de terceiros (`wire`/`fx`/`dig`) — viável, porém menos idiomático e mais pesado.

## Gap identificado
- **Severidade**: **médio-alto** — é uma mudança de modelo mental, não sintática.
- **Implicações concretas** (com exemplo do próprio legado):
  - **1 — Exceções → erros explícitos.** O legado lança `AppwriteException` em toda parte (ex.: `Executions/Create.php:488` captura `\Throwable`, seta status `failed` e relança; guards `BUILD_NOT_READY`, `MEMBERSHIP_DOWNGRADE_PROHIBITED`). Em Go, cada função retorna `(valor, error)`; não há try/catch. Cada um dos ~35 códigos de exceção vira um valor de erro explícito + wrapping — afeta **todos** os handlers.
  - **2 — Container DI → wiring explícito.** O legado resolve recursos via container Utopia (`->inject(...)`, `code-analysis.md §2.2`). Go não tem esse idioma: dependências viram campos de struct / argumentos passados no *composition root*. O padrão `inject()` desaparece.
  - **3 — Herança da base `Action` → composição + interfaces.** Todo endpoint estende `Action` (herança). Go não tem herança; comportamento compartilhado vira struct embutido + uma interface `Handler`/`Route`. A hierarquia de classes precisa ser redesenhada como contrato de interface.
  - **4 — Corrotinas Swoole + workers de fila Redis → goroutines/channels.** 12 workers consomem filas Redis; Swoole faz corrotinas in-process. Em Go a concorrência natural é goroutine+channel (CSP). O modelo fila-worker pode permanecer (Redis), mas os idiomas in-process mudam — e há oportunidade de trocar parte do fan-out Redis por channels para o objetivo de **footprint menor**. Ordem/idempotência (ADR-008) precisam de padrões Go equivalentes.
  - **5 — Config PHP dinâmica → tipos estáticos Go.** O schema vive em arrays PHP em runtime (`app/config/collections/*.php` — 48 coleções, filtros `encrypt`/`subQuery` aplicados dinamicamente por `utopia-php/database`; ver `specs/database/`). Go precisa de structs tipadas ou um schema-engine; a abordagem array-como-schema precisa virar definição tipada.

## Opções apresentadas ao usuário
1. **Adotar o paradigma natural do Go (transformacional)** — Go idiomático: erros explícitos, composição por interfaces, goroutines/channels, schema tipado. *Consequências*: melhor para o objetivo de **menos recursos** (sem cerimônia OO, concorrência eficiente); porém é reescrita de modelo mental — não há port 1:1 de classes PHP; maior distância das specs de origem exige rastreabilidade cuidadosa.
2. **Forçar OO+DI estilo PHP em Go (conservador)** — replicar base `Action`, um container DI de terceiros, wrapping de erros estilo exceção. *Consequências*: facilita mapeamento 1:1 e paridade, mas **luta contra os idiomas do Go**, é mais pesado e aloca mais — **contraria o objetivo de reduzir recursos**; débito técnico idiomático permanente.
3. **Híbrido (equilibrado)** — Go idiomático nas camadas de concorrência/erro/dados (CSP, erros explícitos, schema tipado) para ganhar eficiência, **preservando** a fronteira módulo/serviço e a topologia fila-worker (que mapeia bem) para manter paridade e port tratáveis; DI leve (wiring explícito ou container pequeno) em vez de framework pesado. *Consequências*: melhor relação eficiência/paridade/risco para solo-dev; exige disciplina para decidir por-camada onde é idiomático vs. onde preserva estrutura.

## Decisão do usuário
- **Escolha**: **1 — Adotar o paradigma natural do Go (transformacional)**
- **Justificativa do usuário**: Objetivo dominante é reduzir footprint numa VPS modesta; Go idiomático (erros explícitos, composição por interfaces, goroutines/channels, schema tipado) maximiza eficiência e alinha-se aos princípios bigpowers. Aceita a reescrita de modelo mental e o custo de rastreabilidade em troca da menor pegada de recursos.
- **Decidido em**: 2026-08-11T16:23:34Z

## Apetite derivado
- `derived_appetite`: **transformational**

## Implicações pendentes para próximos agentes
| Agente | Implicação | Como honrar |
|---|---|---|
| Curator | Config PHP dinâmica (impl. 5) | Marcar o schema-as-code (48 coleções) como MIGRA; decidir por-coleção o que vira struct tipada vs. schema-engine |
| Strategist | Fila-worker vs. channels (impl. 4) | A estratégia de cutover precisa preservar semântica de idempotência/estado terminal (ADR-008) sob o novo modelo de concorrência |
| Designer | Herança `Action` + DI (impl. 2, 3) | Topologia/arquitetura alvo deve definir a interface `Handler` e o wiring explícito que substituem base-class + container |
| Inspector | Exceções → erros (impl. 1) | Specs de paridade devem mapear cada código de erro do legado a um erro Go observável na mesma resposta HTTP |

## Notas
- O usuário pediu explicitamente orientação pelos **skills/principles/conventions do bigpowers** nas decisões de stack/arquitetura — os agentes posteriores devem citar esses princípios ao propor framework, banco e topologia.
- Objetivo dominante: **reduzir footprint** preservando **paridade de features e de UI/UX do console**. Isso favorece idiomas Go eficientes nas camadas quentes.
