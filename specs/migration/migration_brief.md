---
schemaVersion: 1
generatedAt: 2026-08-05T00:00:00Z
reversa:
  version: "1.2.57"
kind: migration_brief
producedBy: orchestrator
hash: "sha256:8e8ca5c2b438b4c5abda36bb257703a005cd9ff212220da2960ea3a7edf8c30a"
---

# Migration Brief

> Documento de critério de migração coletado em entrevista no início do `/reversa-migrate`.
> Consumido pelos seis agentes do Time de Migração. Não pergunta paradigma (responsabilidade do Paradigm Advisor) nem apetite (derivado em `paradigm_decision.md`).

## Objetivo da migração
Produzir uma versão do Appwrite que consuma menos recursos (CPU/RAM) ao rodar em uma VPS modesta, mantendo o conjunto de features que hoje atraem o usuário para hospedar sites (Sites, Functions, Databases, Storage, Auth, Teams, etc.), mas sobre um runtime mais leve que a pilha PHP/Swoole atual.

## Métricas de sucesso
- **Paridade funcional**: todas as features hoje disponíveis no Appwrite relevantes para hosting de sites devem existir na versão migrada.
- **Paridade de UI**: o console deve ter a mesma estrutura, look and feel, opções, botões e features do Appwrite original — sem regressão perceptível para o usuário.
- **Testes de UI**: 100% dos testes de paridade de UI devem passar antes de considerar a migração concluída.
- **Consumo de recursos**: pegada de memória/CPU menor que a stack PHP/Swoole atual na mesma VPS (redução qualitativa esperada; sem alvo numérico fixado pelo usuário).

## Restrições
- **Prazo**: nenhum prazo fixo — ritmo determinado pela cota diária de uso de IA do desenvolvedor (único fator limitante declarado).
- **Orçamento**: solo-dev, sem orçamento de equipe ou contratação externa.
- **Técnicas**: nenhuma API externa travada foi declarada; os contratos hoje existentes em `specs/openapi/` devem ser tratados como referência de paridade.
- **Operacionais**: nenhuma janela de manutenção formal — ambiente pessoal, VPS única.

## Fatores de risco conhecidos
- Nenhum risco conhecido foi apontado pelo usuário no momento da entrevista.

## Stakeholders
| Nome / papel | Responsabilidade na migração |
|---|---|
| danielvm (solo-dev) | Todas as decisões, implementação e validação da migração |

## Stack alvo
- **Linguagem**: Go (Golang) — decisão já tomada pelo usuário.
- **Framework**: a definir pelo Paradigm Advisor / Designer, avaliado sob os princípios e convenções do bigpowers.
- **Banco**: a definir — priorizar opção mais leve que ainda cubra os recursos hoje usados pelo Appwrite (relacional/documento, cache, filas).
- **Mensageria** (se houver): a definir, equivalente funcional ao Redis pub/sub e às filas usadas hoje.
- **Infra**: VPS única, self-hosted, foco explícito em baixo consumo de recursos.
- **Outros componentes relevantes**: cache, observabilidade e demais componentes atuais (Redis, Traefik, etc.) devem ser reavaliados por variantes mais leves quando aplicável.

## Escopo declarado
- **Incluído**: migração completa de todas as features do Appwrite (full migration) — todos os 22 módulos mapeados em `specs/` (Account, Advisor, Avatars, Compute, Console, Databases, Functions, Health, Migrations, Notifications, Organization, Presences, Project, Projects, Proxy, Sites, Storage, Teams, Tokens, Users, VCS, Webhooks).
- **Excluído**: nenhum módulo excluído — escopo é migração total.

## Notas livres
- O usuário pediu explicitamente que o time de migração use os skills, principles e conventions do **bigpowers** (catálogo de skills instalado no ambiente) para orientar as decisões de stack e arquitetura ao longo da migração, além do Go já definido como linguagem-alvo.
- Prioridade explícita e repetida: preservar a experiência de UI/UX e o conjunto de features do Appwrite — especialmente as voltadas a hospedagem de sites — enquanto se reduz o custo de infraestrutura.
