---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: topology_decision
producedBy: designer
hash: "sha256:24963bd8f36cbcaeb754e408781b649b9b2283565507891dd11515b6675795e1"
---

# Topology Decision

> Como organizar o sistema Go novo: preservar a topologia do legado, adotar uma moderna, ou híbrido.
> Leitura obrigatória do Designer (Fase 2) e do agente de codificação (árvore de pastas).

## Topologia do legado detectada
- **Padrão organizacional**: **Modular Monolith, package-by-feature (módulos por domínio)** + workers assíncronos.
- **Confiança**: 🟢 CONFIRMADO
- **Evidências**:
  - `architecture.md §1`: "Modular Monolith with async workers" — processo único Swoole + 12 workers via filas Redis. 🟢
  - ADR-001 (`adrs/001-estrutura-modulos.md`): feature modules em `Modules/{Name}/` com Action pattern. 🟢
  - `AGENTS.md` (Module structure): cada módulo com `Http/`, `Workers/`, `Tasks/`, `Services/`. 🟢
- **Mapa da árvore legada** (resumido):
  ```
  src/Appwrite/Platform/
    Modules/{Account,Databases,Functions,Storage,...22}/
      Http/  Workers/  Tasks/  Services/  Module.php
    Workers/            # 8 workers cross-cutting
  app/  (init, config, controllers, http.php/worker.php/realtime.php)
  ```

## Diagnóstico estrutural
- **Acoplamento**: médio — fronteiras de módulo limpas, mas forte acoplamento de runtime ao framework compartilhado `utopia-php` (~40 pacotes) e a recursos globais via DI.
- **Coesão por módulo**: alta — cada módulo agrupa HTTP + workers + tasks do seu domínio (evidência: AGENTS.md, ADR-001). 🟢
- **Módulos órfãos / mortos**: nenhum detectado.
- **Camadas redundantes**: nenhuma clara.
- **Violações de fronteira**: acoplamento de dados entre plano-plataforma e plano-projeto (coleções compartilhadas), mas mediado por `dbForPlatform`/`dbForProject`.
- **Mistura de estilos**: núcleo OO-com-DI + bordas event-driven (workers) — coeso, não caótico.
- **Avaliação geral**: **parcialmente problemática** — o mapa de features é saudável e limpo; o débito está no acoplamento pesado ao framework PHP e no processo único.

## Topologia moderna proposta
- **Padrão**: **Hexagonal (ports & adapters) por feature, em monorepo Go com vertical slices** — um binário leve + binários de worker opcionais.
- **Justificativa**: alinha-se ao paradigma transformacional (interfaces/composição, erros explícitos nas portas, goroutines nos adapters), à estratégia Strangler Fig (portas explícitas por módulo permitem flip incremental e Parallel Run isolável) e ao objetivo de footprint (um binário Go estático, sem framework pesado).
- **Ganhos concretos**:
  - Testabilidade por feature (domínio puro testável sem I/O; adapters mockáveis via interface).
  - Isolamento para Strangler/Parallel Run: cada módulo tem porta HTTP e porta de dados próprias → roteável e comparável isoladamente.
  - Footprint: composição explícita, sem container de DI em runtime; stdlib-first (princípios bigpowers).
- **Custo / risco**:
  - Curva de aprendizado de hexagonal idiomático em Go (evitar over-engineering de portas).
  - Esforço de definir os contratos de porta por módulo antes de codar.
- **Esboço da árvore proposta**:
  ```
  cmd/
    server/        # binário API (HTTP)
    worker/        # binário(s) de workers
  internal/
    {account,databases,functions,storage,...}/   # 1 pacote por feature (bounded context)
      domain/      # entidades, invariantes, regras puras (sem I/O)
      ports/       # interfaces (repos, publishers, gateways)
      adapters/    # postgres, http handlers, queue consumers
      service.go   # orquestração do caso de uso
  platform/        # cross-cutting: auth, permissions, crypto, errors, ratelimit, db, queue
  pkg/             # utilitários públicos
  ```

## Opções apresentadas ao usuário
1. **Preservar topologia legada** (conservador) — replicar módulos-por-domínio em pacotes Go quase 1:1. *Consequências*: mapa mental idêntico, rastreabilidade trivial; mas carrega o acoplamento estrutural do PHP para Go e viola a regra "sem decomposição 1:1".
2. **Adotar topologia moderna proposta** (transformacional) — hexagonal por feature, do zero. *Consequências*: máximo ganho de testabilidade/isolamento/footprint; maior esforço de design de portas; distância do mapa legado exige rastreabilidade explícita.
3. **Híbrido** (equilibrado, **recomendado**) — **preservar o mapa macro de 22 bounded contexts** (essencial para roteamento Strangler por módulo e paridade de contrato de API/console), **modernizar o interior de cada módulo para hexagonal Go** (domain/ports/adapters). *Consequências*: mantém rastreabilidade e o gate de Parallel Run por módulo; moderniza onde há ganho real; evita 1:1 no micro sem perder o macro que serve à estratégia.

## Decisão do usuário
- **Escolha**: **3 — Híbrido**, com os bounded contexts **definidos via `model-domain`** (skill bigpowers, 2026-08-11).
- **Justificativa do usuário**: pediu explicitamente definir a topologia por modelagem de domínio (`/domain-model`) em vez de escolher uma opção genérica. Resultado: **9 bounded contexts** confirmados (Identity, Organizations, Projects, Data, Storage, Compute, VCS, Networking, Messaging) + camada `platform/` cross-cutting. Macro-mapa preservado para roteamento Strangler + paridade; interior de cada contexto modernizado para hexagonal Go.
- **Decidido em**: 2026-08-11T16:23:34Z
- **Fonte do domínio**: `specs/tech-architecture/tech-stack.md` (sessão model-domain).

## Mapeamento legado → novo (preliminar, sob opção 3)
| Módulo / pasta legada | Bounded context novo | Tipo | Observações |
|---|---|---|---|
| `Modules/{22 módulos}` | `internal/{22 features}` | preservado (macro) + modernizado (micro) | mesmo bounded context; interior vira domain/ports/adapters |
| `src/Appwrite/Platform/Workers/` + `Modules/*/Workers` | `cmd/worker` + `internal/*/adapters` (consumers) | fundido por responsabilidade | 12 workers viram consumers de fila por feature |
| `utopia-php/*` (framework) | `platform/` (auth, crypto, db, ratelimit, errors) | reescrito | subconjunto tipado nativo (BR-HUMANA-001) |
| `app/config/collections/*.php` | `platform/db` + `internal/*/adapters/postgres` | reescrito | schema tipado Postgres (BR-DESCARTAR-005) |

## ⚠️ Revisão pós red-team (2026-08-11)
Três agentes adversariais validaram esta topologia (`redteam_findings.md`). Correções **estruturais** aplicadas (não exigiram nova decisão humana):
- **9 → 11 contextos + BFF**: promovidos **Migrations** e **Advisor** a contextos reais (têm aggregate + state machine); **Console** reclassificado como **BFF** acima dos contextos (não dentro de Projects); **Webhooks** movido para Projects.
- **Ownership de `projects`** (B1): Projects é dono único do write-model; `/v1/organization/projects/*` vira view.
- **VCS↔Compute** (B2/S4): VCS publica `repository.pushed` via outbox; sem escritas diretas em `deployments`. Ordem: VCS antes/junto de Compute.
- **Unidade de flip = contexto** (S1) e **Identity+Organizations co-flipam** (S2). Ordem topológica em `migration_strategy.md`.
- **`platform/` deixa de ser god-package** (O2): só leaf packages (auth/permissions/crypto/errors/ratelimit) + serviços stateless (Health/Avatars/Presences).

> **Pendente de decisão humana** (pivôs de direção, não estruturais): O1 (hexagonal vs. flat), O3 (Redis vs. só Postgres), O4 (reimplementar engine vs. `sqlc`), O5 (outbox vs. `jobs`+`SKIP LOCKED`). Ver `redteam_findings.md §C`.

## Implicações para a Fase 2 (após aprovação)
- Bounded contexts = os 22 módulos (agrupamentos justificados onde fizer sentido, ex.: Functions+Sites+Compute compartilham build).
- `platform/` concentra os cross-cutting decididos: permissões, roles (`roles.php`), rate limits, criptografia (25 campos), erros.
- Data model alvo = Postgres (JSONB + pgvector), a partir das 48 coleções de `specs/database/`.
