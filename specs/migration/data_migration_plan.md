---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: data_migration_plan
producedBy: designer
hash: "sha256:12cd548cbb08f11d9f572baecabf1c48da7b7c4b07fa13a897b751ac88c50b06"
---

# Data Migration Plan — legado → PostgreSQL

> Move dados de MongoDB (plataforma + projeto) + MariaDB para um Postgres consolidado. Executado **por módulo** (alinhado ao Strangler), com ownership de dado exclusivo por módulo durante a coexistência (R-D2).

## Mapeamento legado → novo
| Fonte legada | Destino Postgres | Tipo |
|--------------|------------------|------|
| Mongo `dbForPlatform.*` (projects, keys, rules, certificates, webhooks, installations…) | tabelas dos contextos Projects/Networking/VCS | renomeação snake_case |
| Mongo `dbForProject.*` (users, sessions, teams, buckets, functions, deployments…) | tabelas dos contextos correspondentes | renomeação |
| Coleções dinâmicas `database_{seq}` (documentos de usuário) | `documents` (JSONB, particionada por `collection_internal_id`) | fusão |
| VectorsDB (Postgres legado) | `vector_documents` (pgvector) | rehospedagem |
| Arquivos (`bucket_{seq}` metadados) | `files` + blobs no object storage | divisão (metadado vs. blob) |

## Transformações por coluna/tabela
| Regra | Tratamento | Inválidos |
|-------|-----------|-----------|
| `object`/filtro `json` → `jsonb` | parse + validação de schema | log + quarentena, não aborta |
| filtro `encrypt` (25 campos) → `bytea` | **decrypt com chave legada → re-encrypt AES-GCM** com a chave do novo sistema | falha de decrypt = parar o módulo (R-D3), nunca gravar texto plano |
| `datetime` string → `timestamptz` | normalizar TZ (UTC) | datas nulas/ inválidas → `NULL` + flag |
| `*InternalId` → FK `bigint` | reconstruir sequência interna; mapear id lógico→internal_id | id órfão → quarentena |
| status de user (bool) | mapear direto | — |
| relationships espelhados (Data) | reconstruir a partir dos metadados de atributo | ambíguos → LACUNA p/ validação |

## Estratégia de ETL
- **Ferramenta**: job Go dedicado (`cmd/etl`) — reusa os structs tipados do `target_data_model.md` (mesma validação do runtime).
- **Fluxo**: extrair (driver Mongo/MariaDB) → transformar (funções puras por tabela) → carregar (COPY em Postgres staging).
- **Idempotência**: upsert por PK natural (`id`); reprocessável sem duplicar.
- **Throughput**: batch + COPY; paralelizável por coleção (goroutines), respeitando ordem de FK.

## Backfill e captura de delta
- **Backfill**: snapshot inicial por módulo em staging.
- **Delta**: durante o Parallel Run, o módulo ainda dono no PHP continua escrevendo; capturar delta por `updated_at` até o flip. Como o ownership é exclusivo por módulo, **não há escrita dupla** — o delta é one-way (PHP→Postgres) até o flip, depois cessa.

## Cutover de dados (por módulo)
1. Backfill do módulo → staging Postgres.
2. Validar (contagens + checksums + amostra).
3. Aplicar delta final (freeze curto de escrita do módulo).
4. Flip de rota (Traefik) — Go passa a ser dono do dado.
5. Verificação pós-corte (contagem + integridade referencial).

## Validação de qualidade
- **Contagens**: `count(*)` por tabela legado vs. Postgres (tolerância 0 pós-freeze).
- **Checksums**: hash de amostra estratificada por tabela.
- **Integridade referencial**: toda FK resolve; zero órfãos fora de quarentena.
- **Cripto**: decrypt de amostra dos 25 campos em Go antes do flip (gate R-D3).

## Lacunas 🔴 (validar no agente de codificação)
- Relacionamentos bidirecionais do contexto Data espelhados por worker — regra exata de reconstrução precisa de validação contra dados reais.
- Blobs de arquivo: caminho/URL de origem (device storage) → destino object storage precisa de mapeamento por `_APP_STORAGE_DEVICE_*`.
