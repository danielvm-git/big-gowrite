---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: target_business_rules
producedBy: curator
hash: "sha256:8f4c5a6549628bb95fb145762170d8a1c9dcb4f0768ef2f40ea12465f31b50b4"
---

# Target Business Rules

> Catálogo das regras de negócio do legado com decisão de migração: MIGRAR, DESCARTAR ou DECISÃO HUMANA.
> Fonte: specs consolidadas (`code-analysis.md`, `state-machines.md`, `permissions.md`, `specs/database/`) + `paradigm_decision.md` (transformacional, Go-native).
> Nota de granularidade: classificação por **tema/módulo** a partir das specs consolidadas (o Writer produziu units por módulo; regras duplicadas foram consolidadas).

## Resumo
- Regras/temas analisados: **28**
- MIGRAR: **18**
- DESCARTAR: **6** (todas vinculadas a paradigma — detalhe em `discard_log.md`)
- DECISÃO HUMANA: **4**

## Regras MIGRAR

> O caso de uso de negócio permanece; em paradigma transformacional muda **onde/como** é expresso (erro explícito, interface, goroutine, struct tipada), nunca o *que*.

### BR-MIGRAR-001 — Ciclo de vida MFA/TOTP (Account)
- **Origem**: `code-analysis.md §3.1`, `state-machines.md` (challenges), Account unit.
- **Confiança**: 🟢
- **Descrição**: secret TOTP → confirmação OTP → persiste authenticator; recálculo/dedupe de `factors` na sessão; rate-limits por desafio.
- **Compat. paradigma**: validação vira função pura na borda; erros (`409` recreate bloqueado) viram erro explícito Go.

### BR-MIGRAR-002 — Máquina de estados de Deployment (`waiting→processing→building→ready|failed|canceled`)
- **Origem**: `state-machines.md §1` (revisado). 🟢
- **Compat.**: estados literais preservados; transições finais gravadas por worker Go; guard de cancel = erro explícito.

### BR-MIGRAR-003 — Máquina de estados de Execution (`waiting/processing/scheduled→completed|failed|canceled`)
- **Origem**: `state-machines.md §2` (revisado). 🟢
- **Compat.**: status derivado do HTTP code (≥500→failed); sync-timeout→failed (não persiste `timeout`).

### BR-MIGRAR-004 — Transações com lock otimista (Databases)
- **Origem**: `code-analysis.md §3.6`, ADR-005. 🟢
- **Descrição**: staging multi-coleção, checagem de versão, commit atômico, rollback.
- **Compat.**: regra de negócio (atomicidade + versionamento) **permanece**; mecanismo vira transação Go idiomática. Não é descarte.

### BR-MIGRAR-005 — Upload em chunks com segurança concorrente (Storage/Functions/Sites)
- **Origem**: `code-analysis.md §3.16, §3.17`, ADR-004. 🟢
- **Descrição**: `Content-Range`, finalize com compressão/criptografia/antivírus, `LockContention→429 Retry-After`.
- **Compat.**: a **necessidade** de exclusão concorrente por recurso persiste entre múltiplos processos → mecanismo distribuído equivalente em Go (ver DECISÃO HUMANA BR-HUMANA-002 sobre lock store).

### BR-MIGRAR-006 — Criptografia em repouso (25 campos)
- **Origem**: `specs/database/business-rules.md`. 🟢
- **Descrição**: filtro `encrypt` em password, secrets, tokens OAuth, credenciais (25 campos).
- **Compat.**: vira camada de encrypt/decrypt tipada em Go (AES-GCM); **regra de segurança inegociável**.

### BR-MIGRAR-007 — Modelo de permissões (ACL de documento + scopes por rota)
- **Origem**: `permissions.md §2, §3`. 🟢
- **Compat.**: `Role::any/users/guests/user/team` e scopes `{resource}.{read|write}` preservados como tipos Go.

### BR-MIGRAR-008 — Enum de papéis de console (guests/users/admin/developer/owner/keys)
- **Origem**: `permissions.md §4.2` (revisado, `app/config/roles.php`). 🟢
- **Compat.**: `admin≡developer` (scopes idênticos), `owner=member∪admins` preservados como conjuntos de scope tipados.

### BR-MIGRAR-009 — Rate limits por rota (abuse-limit/time/key)
- **Origem**: `permissions.md §8` (revisado). 🟢
- **Compat.**: valores (10/15/50/60/100/120, janela 3600s/60s) e chaves (`ip`/`userId`/`url`/`email`) preservados; algoritmo sliding-window reimplementado em Go (era `utopia-php/abuse`).

### BR-MIGRAR-010 — Máquina de estados de Rule/Proxy (domínio + TLS)
- **Origem**: `state-machines.md §3`. 🟢

### BR-MIGRAR-011 — Máquina de estados de Migration (`pending→processing→completed|failed`, retry só de failed)
- **Origem**: `state-machines.md §4`. 🟢

### BR-MIGRAR-012 — Membership (invited→confirmed, ≥1 owner)
- **Origem**: `state-machines.md §5`, `code-analysis.md §3.18`. 🟢

### BR-MIGRAR-013 — Message lifecycle (`draft→scheduled→processing→sent|failed`, terminal idempotente)
- **Origem**: `state-machines.md §7b` (revisado). 🟢

### BR-MIGRAR-014 — Estado terminal obrigatório em jobs de fila (ADR-008)
- **Origem**: `state-machines.md §8`, ADR-008. 🟢
- **Compat.**: em Go/CSP, idempotência por CloudEvent id + estado terminal permanece obrigatória.

### BR-MIGRAR-015 — Fan-out de deploy VCS + preview rules
- **Origem**: `code-analysis.md §3.21, §3.15`. 🟢

### BR-MIGRAR-016 — Imports de hash de senha (Bcrypt/MD5/Argon2/SHA/PHPass/Scrypt/Scrypt-Modified)
- **Origem**: `code-analysis.md §3.20`. 🟢
- **Compat.**: essencial para paridade de importação de usuários; bibliotecas Go equivalentes por algoritmo.

### BR-MIGRAR-017 — Proteção SSRF em Avatars (allowlist, redirects limitados)
- **Origem**: `code-analysis.md §3.3`. 🟢

### BR-MIGRAR-018 — Schema de 48 coleções / 606 atributos / 233 índices (13 unique, 18 fulltext)
- **Origem**: `specs/database/`. 🟢
- **Compat.**: o **schema** migra; a representação PHP-array-em-runtime NÃO (ver BR-DESCARTAR-005). Vira definição tipada Go ou schema-engine (ver BR-HUMANA-001).

## Regras DESCARTAR (resumo)

| ID | Origem | Motivo curto | Vínculo a paradigma? |
|---|---|---|---|
| BR-DESCARTAR-001 | `code-analysis.md §2.5` (AppwriteException) | try/catch + exceções → erros explícitos Go | sim |
| BR-DESCARTAR-002 | `code-analysis.md §2.2` (`inject()`) | container DI Utopia → wiring explícito | sim |
| BR-DESCARTAR-003 | `AGENTS.md` (base `Action`) | herança de classe → interface + composição | sim |
| BR-DESCARTAR-004 | `code-analysis.md §3.21` (corrotinas Swoole) | corrotina Swoole → goroutine/channel | sim |
| BR-DESCARTAR-005 | `app/config/collections/*.php` | schema como array PHP runtime → tipos estáticos | sim |
| BR-DESCARTAR-006 | convenção "só utopia-php" (`AGENTS.md`) | restrição de ecossistema PHP → irrelevante em Go | sim |

> Detalhe completo em `discard_log.md`. **Nenhuma regra de negócio é perdida** — apenas mecanismos do paradigma PHP absorvidos por construção em Go.

## Regras DECISÃO HUMANA

### BR-HUMANA-001 — Reimplementar `utopia-php/database` (engine multi-DB + filtros + relationships) vs. adotar lib Go?
- **Confiança origem**: 🟢 (o legado depende fortemente dessa engine).
- **Contexto**: filtros dinâmicos (`encrypt`, `json`, `subQuery*`), relacionamentos espelhados, adapters multi-DB são providos por `utopia-php/database`. Em Go não existe equivalente 1:1.
- **Recomendação do Curator**: para footprint mínimo + paradigma transformacional, **reimplementar um subconjunto tipado** (só o que os 22 módulos usam) em vez de portar a engine genérica inteira. Alternativa: `ent`/`sqlc` + camada de filtros custom.
- **Impacto**: define o esforço central de toda a migração de dados.

### BR-HUMANA-002 — Estratégia de banco: manter multi-adapter (Mongo+MariaDB+Postgres) ou consolidar em um só?
- **Contexto**: brief pede "banco mais leve". Multi-adapter é peso morto se o alvo é uma VPS.
- **Recomendação**: consolidar em **PostgreSQL** (cobre relacional + JSONB para documentos + pgvector para VectorsDB, um só processo) — maior alavanca de footprint. SQLite como modo ultra-leve é opção.
- **Impacto**: Designer (data model) e Strategist (data migration) dependem desta escolha.

### BR-HUMANA-003 — Executor de Functions e subsistemas externos (executor, browser p/ screenshots, assistant IA)
- **Contexto**: Functions dependem de um executor externo (`_APP_EXECUTOR`), screenshots de um browser interno, assistant de um serviço IA. Não são código PHP portável — são serviços.
- **Recomendação**: manter como serviços externos na fase 1 (a migração Go troca o backend, não o executor); reavaliar footprint depois.
- **Impacto**: escopo de paridade do Inspector para Functions/Sites/Avatars.

### BR-HUMANA-004 — Console UI: preservar o SPA existente apontando p/ a API Go, ou rewrite?
- **Contexto**: brief exige **paridade de UI** ("mesmo look and feel, botões, features"). O console é um SPA separado (repo `appwrite-console`), consumido via API. A migração do backend Go **não** exige reescrever o console.
- **Recomendação (forte)**: **preservar o console SPA como está** e garantir paridade de **contrato de API** (`specs/openapi/appwrite.yaml`). Isso satisfaz "paridade de UI" sem reescrever o front-end. O Screen Translator deve operar em modo *literal/paridade de contrato*, não redesenho.
- **Impacto**: muda radicalmente o escopo do Screen Translator e do Inspector (paridade = contrato de API + testes de UI contra o mesmo console).

## Rastreabilidade
Cada item MIGRAR/DESCARTAR aponta para a spec de origem; itens DECISÃO HUMANA replicados em `ambiguity_log.md` com status PENDENTE.
