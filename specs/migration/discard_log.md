---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: discard_log
producedBy: curator
hash: "sha256:5600a0a4fbd8a59ad11bc1947870355c869103b0cdd2dc882b40d61440a5d4eb"
---

# Discard Log

> Registro completo do que foi descartado da migração e por quê. **Todos** os itens aqui são mecanismos do paradigma PHP absorvidos por construção em Go — nenhuma regra de negócio é perdida.

## Itens descartados

### BR-DESCARTAR-001 — Tratamento de erro por exceção (try/catch + `AppwriteException`)
- **Origem**: `code-analysis.md §2.5`; ex.: `Executions/Create.php:483-490`.
- **Descrição**: fluxo de erro via `throw`/`catch \Throwable` e propagação ao controller.
- **Vinculado a paradigma**: **sim** (OO/PHP → Go).
  - Go absorve o caso com **erros como valores** `(v, error)` + `errors.Is/As` + wrapping. A *semântica* (qual código HTTP para qual falha) é preservada como BR-MIGRAR (mapeada pelo Inspector, impl. 1 do `paradigm_decision.md`).
- **Reposição**: tipo `AppError` Go com código + status HTTP; middleware traduz para resposta.
- **Risco**: baixo (mecânico, porém abrangente — toca todo handler).

### BR-DESCARTAR-002 — Container de injeção de dependências (`->inject(...)`)
- **Origem**: `code-analysis.md §2.2`.
- **Descrição**: resolução de recursos (`dbForProject`, `queueForEvents`…) via container Utopia.
- **Vinculado a paradigma**: **sim**.
  - Go absorve com **wiring explícito** no composition root (campos de struct / args). Sem perda: as dependências continuam existindo, mudam de forma de entrega.
- **Reposição**: construtores explícitos; opcionalmente `google/wire` (codegen, zero-runtime) se o volume justificar.
- **Risco**: baixo.

### BR-DESCARTAR-003 — Herança da classe base `Action`
- **Origem**: `AGENTS.md` (Action pattern).
- **Descrição**: todo endpoint estende `Appwrite\Platform\Action`, herdando registro de rota/hooks.
- **Vinculado a paradigma**: **sim**.
  - Go absorve com **interface `Handler` + composição** (struct embutido para comportamento comum). Contrato REST (create/get/list/update/delete) preservado.
- **Reposição**: interface `Route`/`Handler` + registrador de rotas.
- **Risco**: baixo-médio (redesenho estrutural, mas mapeável 1:1 por endpoint).

### BR-DESCARTAR-004 — Corrotinas Swoole (concorrência in-process)
- **Origem**: `code-analysis.md §3.21` (análise concorrente de repos VCS), runtime Swoole.
- **Descrição**: concorrência cooperativa via corrotinas Swoole.
- **Vinculado a paradigma**: **sim**.
  - Go absorve **nativamente** com goroutines + channels (CSP) — mais leve, alinhado ao objetivo de footprint.
- **Reposição**: goroutines + `errgroup`/channels.
- **Risco**: baixo (Go é mais forte aqui).

### BR-DESCARTAR-005 — Schema como array PHP avaliado em runtime
- **Origem**: `app/config/collections/*.php` (48 coleções; ver `specs/database/`).
- **Descrição**: coleções/atributos/índices/filtros definidos como arrays PHP carregados em runtime.
- **Vinculado a paradigma**: **sim** (dinâmico PHP → estático Go).
  - Go absorve com **definição tipada** (structs + tags) ou um schema-engine tipado. O **schema em si** migra (BR-MIGRAR-018); só a representação runtime-array é descartada.
- **Reposição**: ver BR-HUMANA-001 (engine) — structs tipadas + migrations, ou lib (`ent`/`sqlc`).
- **Risco**: médio (é a espinha dorsal de dados; a decisão está aberta em BR-HUMANA-001/002).

### BR-DESCARTAR-006 — Restrição de ecossistema "só utopia-php"
- **Origem**: `AGENTS.md` (Conventions: "avoid dependencies outside utopia-php").
- **Descrição**: convenção que restringe libs ao ecossistema Utopia.
- **Vinculado a paradigma**: **sim** (é uma regra do stack PHP, sem sentido em Go).
  - Go absorve adotando o ecossistema Go idiomático (stdlib-first, conforme princípios bigpowers).
- **Reposição**: convenções Go do próprio projeto novo (a definir em seed-conventions do repo Go).
- **Risco**: baixo.

## Itens descartados por mudança de paradigma (subseção dedicada)

> **Todos os 6 itens** acima têm `Vinculado a paradigma = sim`. Auditoria explícita para o agente de codificação:
>
> | ID | Mecanismo PHP | Como o Go absorve |
> |----|---------------|-------------------|
> | 001 | try/catch + exceção | erro como valor `(v, error)` + wrapping |
> | 002 | container DI (`inject`) | wiring explícito / `google/wire` |
> | 003 | herança `Action` | interface `Handler` + composição |
> | 004 | corrotina Swoole | goroutine + channel (CSP) |
> | 005 | schema array runtime | struct tipada + migrations |
> | 006 | "só utopia-php" | ecossistema Go idiomático |
>
> Nenhuma regra de negócio descartada. Consistente com o `paradigm_decision.md` (transformacional): o paradigma alvo absorve cada mecanismo por construção.
