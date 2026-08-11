---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: migration_strategy
producedBy: strategist
hash: "sha256:2d77c80312b2262e1c2bf3a436f8a98cd0bff542c3ab38e2791c7e4e6a64d482"
---

# Migration Strategy

> Estratégias avaliadas para migrar Appwrite (PHP/Swoole) → Go, paradigma transformacional, banco consolidado em PostgreSQL.

## Contexto sintetizado
- **Tamanho do legado**: grande — 22 módulos, 12 workers, 48 coleções/606 atributos, ~1.259 arquivos PHP. Integrações vivas: GitHub/GitLab VCS, OAuth providers, SMTP, executor de functions, browser (screenshots), assistant IA, ClamAV.
- **Apetite derivado**: **transformational** (`paradigm_decision.md`).
- **Gap de paradigma**: **médio-alto**.
- **Restrições**: solo-dev; sem prazo fixo (ritmo limitado pela cota diária de IA); sem orçamento; **não é produção crítica sob carga** (VPS pessoal). Sem integrações regulatórias.
- **Alavancas já decididas**: DB → PostgreSQL; console → provável preservação (gate no Screen Translator).

## Estratégias avaliadas

### ✅ RECOMENDADA — Strangler Fig (bordas profundas) + Parallel Run nos módulos críticos
- **Como**: subir o servidor Go **ao lado** do PHP; um roteador (Traefik, já presente) direciona caminhos/módulos migrados para o Go, o resto continua no PHP. Migração **bottom-up, módulo a módulo**. Cada módulo migrado passa por **Parallel Run** (Go e PHP recebem a mesma requisição espelhada; saídas comparadas) antes de flipar o tráfego real — obrigatório para os BR-MIGRAR críticos (auth/permissões, criptografia, transações, máquinas de estado).
- **Adequação ao apetite**: alta — permite reescrita idiomática Go por módulo sem parar tudo.
- **Adequação ao gap**: alta — o Parallel Run prova paridade justamente onde o paradigma muda (erros, concorrência).
- **Adequação ao solo-dev/cota de IA**: **ótima** — trabalho fatiável em incrementos pequenos; cada sessão avança 1 módulo; sem pressão de janela única.
- **Custo**: médio · **Risco**: baixo · **Tempo**: longo (aceitável — sem prazo).
- **Prós**: rollback por módulo (rota volta ao PHP); paridade contínua; footprint reduzido incrementalmente. **Contras**: exige rodar dois stacks em paralelo por um período; consistência de dados durante a coexistência (ver risco D-2).

### Alternativa — Big Bang
- Reescrever todo o Go e cortar de uma vez.
- **Por que NÃO recomendada**: o sistema é grande (22 módulos); lead time enorme antes de qualquer entrega; paridade difícil de validar holisticamente; risco alto. O catálogo reserva Big Bang para sistemas **pequenos** — não é o caso.

### Alternativa — Branch by Abstraction
- **Por que NÃO recomendada**: pressupõe migração *in-process* (mesma runtime trocando implementação atrás de uma abstração). PHP→Go são runtimes distintos; não há abstração in-process compartilhável. Inaplicável aqui.

## Recomendação
**Strangler Fig com bordas profundas, Parallel Run como gate de paridade.**

### ⚠️ Revisão pós red-team (2026-08-11) — unidade de flip e ordem topológica
O red-team (`redteam_findings.md`) mostrou que **a unidade de flip é o bounded context, não o módulo** (S1), e que a ordem deve ser **topológica** por dependência de dados/escrita, não por curva de aprendizado (S5). Correções:
- **Flip por contexto**, não por módulo: Account+Users flipam juntos (Identity); Functions+Sites+Compute juntos (Compute).
- **Identity + Organizations co-flipam** como uma unidade atômica — o fluxo de convite de time escreve `users`+`sessions`+`teams`+`memberships` numa só requisição (S2, `Teams/Http/Memberships/Status/Update.php:129,139`).
- **VCS antes ou junto de Compute** — Compute lê `installations` de forma síncrona no create/update de function/site (S4).
- **Workers cross-cutting (ex.: `Deletes.php`) migram por último**, só depois de todos os contextos que eles tocam estarem em Postgres (S3) — evita órfãos no cascade delete.

### Ordem topológica revisada
1. **Fundação** (`platform/`): auth, permissões (`roles.php`), criptografia (25 campos), erros, ratelimit + schema Postgres base.
2. **Identity + Organizations** (unidade atômica) — base de tenancy/auth.
3. **Projects** (inclui Webhooks) — depende de Organizations (teamInternalId).
4. **Data**, **Storage** — dependem de Projects (projectInternalId).
5. **VCS** → **Compute** (nesta ordem; VCS publica `repository.pushed` via outbox, Compute consome).
6. **Networking**, **Messaging**, **Migrations**, **Advisor**.
7. **Console (BFF)** — depois que Identity/Orgs/Projects estão em Go.
8. **Workers cross-cutting** (Deletes, Stats) — por último.
9. **Decomissão** do PHP.

> Gate adicional por flip: **"todos os contextos referenciados já migrados?"** (S4). Decisão final da estratégia é humana.
