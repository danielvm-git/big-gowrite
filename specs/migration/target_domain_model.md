---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: target_domain_model
producedBy: designer
hash: "sha256:2fe6f57946f00c782f71da0d513c1437532d9a85fb1aadecf151bbc17e9e96a8"
---

# Target Domain Model — gowrite

> Aggregates por bounded context, honrando o paradigma Go-native. Invariantes e state machines vêm das fontes autoritativas (`state-machines.md`, `domain.md`).

## Aggregates por contexto

### Identity
- **User** (root): id, email, hash, mfa, status(bool), prefs, labels. Invariante: email único entre `users`+`identities`.
  - Entidades: Session (factors), Identity (OAuth), Target, Authenticator, Challenge, Token.
  - Comandos: Register, CreateSession, EnableMfa, CreateChallenge, VerifyChallenge, ImportHash.
  - Eventos (outbox): `user.created`, `session.created`, `user.verification.*`.
- **State machines**: Challenge (rate-limited, expira); ver `state-machines.md`.

### Organizations
- **Team** (root): id, name, prefs. Invariante: **≥1 owner** (`MEMBERSHIP_DOWNGRADE_PROHIBITED`).
  - Entidades: Membership (status derivado de `invited`/`joined`).
  - Comandos: CreateTeam, Invite, ConfirmMembership, UpdateRole.
  - Eventos: `team.created`, `membership.created`, `membership.confirmed`.

### Projects
- **Project** (root): id, teamInternalId, keys, platforms, oauth2, smtp, auths, onboarding, schedules. Invariante: tenant root; `projectInternalId` propaga.
  - Comandos: CreateProject, CreateKey (persistida/efêmera), UpdateAuthPolicy, TransferTeam.
  - Eventos: `project.created`, `project.transferred`.

### Data
- **Database** (root) → Collection → Attribute/Index/Document. **Transaction** (aggregate): staging multi-coleção + versão.
  - Invariante: optimistic-lock (checa versão, commit atômico, rollback). ADR-005.
  - Comandos: CreateCollection, CreateAttribute(async), RunTransaction, VectorSearch.
  - Eventos: `database.*`, `attribute.ready` (worker).

### Storage
- **Bucket** (root) → File. **ResourceToken** (aggregate): `resourceId=bucketId:fileId`, secret cifrado.
  - Invariante: chunked-upload lock; `fileSecurity` vs bucket perms.
  - Comandos: CreateBucket, UploadChunk, Finalize(compress+encrypt+scan), CreateToken.

### Compute
- **Function**/**Site** (roots) → Deployment → Execution; Variable.
  - Invariantes: Deployment SM (`waiting→processing→building→ready|failed|canceled`); Execution SM (`waiting/processing/scheduled→completed|failed|canceled`); build dedup por CloudEvent id; `ready` exigido p/ execução.
  - Comandos: CreateDeployment, CancelBuild, CreateExecution, ActivateDeployment.
  - Eventos: `deployment.ready`, `execution.completed`.

### VCS
- **Installation** (root) → Repository → (deploy trigger). Invariante: assinatura de webhook (403 em mismatch); state HMAC.
  - Comandos: CreateInstallation, HandleWebhook, TriggerDeploy.

### Networking
- **Rule** (root) + Certificate. Invariante: SM (`unverified→verifying→verified→certificate_generating`); CAA gate.
  - Comandos: CreateRule, UpdateVerification, IssueCertificate(worker).

### Messaging
- **Message** (root) → Target/Provider/Topic/Subscriber. Invariante: SM (`draft→scheduled→processing→sent|failed`), terminal idempotente.
  - Comandos: CreateMessage, Schedule, Send(worker).

## Regras de domínio → local no domínio novo
| BR-MIGRAR | Onde vive |
|-----------|-----------|
| 001 MFA/TOTP | Identity/domain (Challenge, Authenticator) |
| 002/003 Deployment/Execution SM | Compute/domain |
| 004 Transações lock otimista | Data/domain (Transaction aggregate) |
| 005 Upload chunks | Storage/domain + adapters (lock) |
| 006 Criptografia 25 campos | platform/crypto (aplicado nos adapters postgres) |
| 007/008 Permissões + roles | platform/permissions |
| 009 Rate limits | platform/ratelimit |
| 010–013 State machines | contexto respectivo/domain |
| 014 Estado terminal (ADR-008) | tabela `jobs` + `cmd/worker` |
| 016 Import de hash | Identity/domain (strategy por algoritmo) |

## Rastreabilidade para o legado
| Legado | Novo | Tipo |
|--------|------|------|
| Account + Users | Identity | fundido (mesmo aggregate User) |
| Teams + Organization | Organizations | fundido |
| Project + Projects + Webhooks | Projects | fundido (Console → BFF, ver abaixo) |
| Databases + VectorsDB | Data | fundido |
| Storage + Tokens | Storage | fundido |
| Functions + Sites + Compute | Compute | fundido |
| VCS | VCS | 1-para-1 (contexto) |
| Proxy | Networking | renomeado |
| Notifications + messaging worker | Messaging | fundido |
| Migrations | Migrations | contexto próprio (post red-team) |
| Advisor | Advisor | contexto próprio (post red-team) |
| Console | Console | **BFF** acima dos contextos |
| Health, Avatars, Presences | platform/ops | stateless (só estes) |

## Contextos adicionais (post red-team)
- **Migrations** (root: Migration): imports de provider (Firebase/Supabase/NHost) + CSV/JSON; SM `pending→processing→completed|failed`, retry só de `failed`; política `onDuplicate`. Toca Data no import.
- **Advisor** (root: Report → Insight): resultados de análise (Lighthouse/db-analyzer); Insight SM `active→dismissed` + cascade delete async.
- **Console** = BFF (não aggregate): compõe Identity/Organizations/Projects.

## Eventos de domínio (via tabela `jobs` — revisado)
Todo efeito cross-context (webhooks, function triggers, realtime, deletes em cascata) escreve um registro na tabela **`jobs`** (ver `target_data_model.md`) na **mesma transação** do aggregate; o `cmd/worker` consome via `SKIP LOCKED`, at-least-once, idempotente por `dedup_key`. Realtime imediato via `pg_notify`. Substitui o event-bus Redis legado (`queueForEvents`) por uma única mecânica Postgres — sem outbox dedicado nem Redis (red-team O3/O5). Ex.: VCS grava job `event:repository.pushed`; Compute consome e cria o deployment (remove a escrita direta VCS→`deployments`, B2).
