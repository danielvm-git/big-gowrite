---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: target_data_model
producedBy: designer
hash: "sha256:8e053d076cc9e5302b6d38af44964255b9fea2b2a699d2806cf7282b4298752b"
---

# Target Data Model — gowrite (PostgreSQL)

> Deriva das **48 coleções / 606 atributos / 233 índices** em `specs/database/`. Um único Postgres substitui Mongo+MariaDB+Postgres.
> Estratégia: tabelas relacionais para entidades com schema fixo; **JSONB** para documentos de usuário (coleções dinâmicas do contexto Data); **pgvector** para embeddings.

## Mapeamento de tipos (utopia-php/database → Postgres)
| Legado | Postgres |
|--------|----------|
| `string`/`text` | `text` (ou `varchar(n)` quando `size` pequeno) |
| `integer` | `bigint` |
| `double` | `double precision` |
| `boolean` | `boolean` |
| `datetime` | `timestamptz` |
| `object` | `jsonb` |
| `vector` | `vector(n)` (pgvector) |
| filtro `encrypt` | `bytea` (AES-GCM na camada `platform/crypto`) |
| filtro `json` | `jsonb` |
| `$id`/`$internalId` | `text id` (PK natural) + `bigint internal_id` (sequência) |

## Entidades por bounded context (amostra representativa — schema completo gerado do `specs/database/data-dictionary.md`)
| Tabela | Aggregate dono | PK | Contexto |
|--------|----------------|----|----------|
| `users` | User | id | Identity |
| `sessions`,`identities`,`targets`,`authenticators`,`challenges`,`tokens` | User | id | Identity |
| `teams`,`memberships` | Team | id | Organizations |
| `projects`,`keys`,`dev_keys`,`platforms`,`webhooks`,`schedules` | Project | id | Projects |
| `databases`,`collections`,`attributes`,`indexes`,`transactions`,`transaction_logs` | Database | id | Data |
| `buckets`,`files`,`resource_tokens` | Bucket | id | Storage |
| `functions`,`sites`,`deployments`,`executions`,`variables` | Function/Site | id | Compute |
| `installations`,`repositories`,`vcs_comments` | Installation | id | VCS |
| `rules`,`certificates` | Rule | id | Networking |
| `messages`,`providers`,`topics`,`subscribers` | Message | id | Messaging |
| `outbox` | — (cross-cutting) | id | platform |

## DDL (esboço — padrões-chave)
```sql
-- Multi-tenant scoping (invariante: projectInternalId em todo recurso)
CREATE TABLE users (
  id           text PRIMARY KEY,
  internal_id  bigserial UNIQUE,
  project_internal_id bigint NOT NULL,
  email        text NOT NULL,
  hash         bytea,               -- encrypt (AES-GCM em platform/crypto)
  status       boolean NOT NULL DEFAULT true,  -- (revisão: booleano, não 0/1/2)
  mfa          boolean NOT NULL DEFAULT false,
  prefs        jsonb  NOT NULL DEFAULT '{}',
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX uq_users_project_email ON users(project_internal_id, email);  -- unicidade por projeto

-- Documentos de usuário (contexto Data): coleção dinâmica → linha JSONB
CREATE TABLE documents (
  id text PRIMARY KEY,
  collection_internal_id bigint NOT NULL,
  project_internal_id bigint NOT NULL,
  data jsonb NOT NULL,
  permissions text[] NOT NULL DEFAULT '{}',
  version int NOT NULL DEFAULT 1,        -- lock otimista (ADR-005)
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX gin_documents_data ON documents USING gin (data jsonb_path_ops);

-- VectorsDB → pgvector
CREATE TABLE vector_documents (
  id text PRIMARY KEY,
  collection_internal_id bigint NOT NULL,
  embedding vector(1536),               -- dimensão configurável por coleção
  metadata jsonb
);
CREATE INDEX ann_vector_documents ON vector_documents USING hnsw (embedding vector_cosine_ops);

-- Jobs (fila única em Postgres — substitui Redis + outbox; red-team O3/O5)
CREATE TABLE jobs (
  id bigserial PRIMARY KEY,
  kind text NOT NULL,                    -- build|delete|message|certificate|migration|event
  payload jsonb NOT NULL,
  dedup_key text UNIQUE,                 -- idempotência por CloudEvent id (ADR-008)
  run_after timestamptz NOT NULL DEFAULT now(),
  locked_at timestamptz,
  attempts int NOT NULL DEFAULT 0,
  status text NOT NULL DEFAULT 'pending',-- pending|processing|done|failed (estado terminal obrigatório)
  created_at timestamptz NOT NULL DEFAULT now()
);
-- consumo: SELECT ... WHERE status='pending' AND run_after<=now()
--          FOR UPDATE SKIP LOCKED LIMIT n;
-- realtime cross-context: pg_notify() no commit; server escuta via LISTEN.
```

## Restrições preservadas
- **13 índices unique** e **18 fulltext** de `specs/database/business-rules.md` → unique constraints + `tsvector`/GIN em Postgres.
- **25 campos `encrypt`** → `bytea`, cifrados/deciphrados em `platform/crypto` (AES-GCM), nunca em texto plano.
- **Relacionamentos** (`*InternalId`) → FKs Postgres reais (ganho vs. legado app-side) onde a fronteira do contexto permite; cross-context via id lógico (sem FK).

## Considerações do paradigma (revisadas post red-team)
- **Fila única `jobs` em Postgres** (`SKIP LOCKED`) substitui Redis + outbox — uma só mecânica de dispatch (O5). Escrita do job na **mesma transação** do aggregate garante integridade sem dual-write.
- **Realtime** via `LISTEN/NOTIFY` (sem Redis pub/sub); **rate-limit** via token-bucket in-memory por instância (O3).
- **Acesso a dados via `sqlc`** — SQL tipado gerado, sem engine genérica (O4). Cripto explícita: `crypto.Encrypt(field)` nos ~poucos write-sites dos 25 campos, não um pipeline de filtro de coluna.
- **Lock otimista** nativo via coluna `version` + `UPDATE ... WHERE version = $expected`.
- Ownership de dado **por contexto** (não por coleção) durante Strangler; sem escrita cross-context para contexto ainda-em-PHP (S2).

## Origem no legado
Cada tabela rastreia para uma coleção em `specs/database/data-dictionary.md`. Renomeações: `snake_case` Postgres (ex.: `resourceTokens`→`resource_tokens`, `transactionLogs`→`transaction_logs`). Fusões: coleções `bucket_{seq}`/`database_{seq}` dinâmicas → `files`/`documents` particionadas por `*_internal_id`.
