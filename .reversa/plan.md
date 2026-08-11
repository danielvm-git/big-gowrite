# Plano de Exploração — gowrite

> Criado pelo Reversa em 2026-08-05
> Marque cada tarefa com ✅ quando concluída.
> Você pode editar este plano antes de iniciar: adicione, remova ou reordene tarefas conforme necessário.

---

## Fase 1: Reconhecimento 🔍

- [x] **Scout** — Mapeamento de estrutura de pastas e tecnologias
- [x] **Scout** — Análise de dependências e gerenciadores de pacotes
- [x] **Scout** — Identificação de entry points, CI/CD e configurações

## Decisão de organização das specs 🗂️

> Entre o Scout e o Arqueólogo, o Reversa pergunta como você quer organizar as specs (por módulo, caso de uso, endpoint, híbrida, por features ou customizada). A escolha fica persistida em `.reversa/config.toml` na seção `[specs]` e não será reperguntada em execuções futuras. Para reapresentar o menu, remova manualmente a seção.

## Fase 2: Escavação 🏗️

> O Reversa preenche esta seção com os módulos reais após o Scout concluir o reconhecimento.

- [x] **Arqueólogo** — Análise do módulo `Account`
- [x] **Arqueólogo** — Análise do módulo `Advisor`
- [x] **Arqueólogo** — Análise do módulo `Avatars`
- [x] **Arqueólogo** — Análise do módulo `Compute`
- [x] **Arqueólogo** — Análise do módulo `Console`
- [x] **Arqueólogo** — Análise do módulo `Databases`
- [x] **Arqueólogo** — Análise do módulo `Functions`
- [x] **Arqueólogo** — Análise do módulo `Health`
- [x] **Arqueólogo** — Análise do módulo `Migrations`
- [x] **Arqueólogo** — Análise do módulo `Notifications`
- [x] **Arqueólogo** — Análise do módulo `Organization`
- [x] **Arqueólogo** — Análise do módulo `Presences`
- [x] **Arqueólogo** — Análise do módulo `Project`
- [x] **Arqueólogo** — Análise do módulo `Projects`
- [x] **Arqueólogo** — Análise do módulo `Proxy`
- [x] **Arqueólogo** — Análise do módulo `Sites`
- [x] **Arqueólogo** — Análise do módulo `Storage`
- [x] **Arqueólogo** — Análise do módulo `Teams`
- [x] **Arqueólogo** — Análise do módulo `Tokens`
- [x] **Arqueólogo** — Análise do módulo `Users`
- [x] **Arqueólogo** — Análise do módulo `VCS`
- [x] **Arqueólogo** — Análise do módulo `Webhooks`

## Fase 3: Interpretação 🧠

- [x] **Detetive** — Arqueologia Git e ADRs retroativos
- [x] **Detetive** — Regras de negócio implícitas e máquinas de estado
- [x] **Detetive** — Matriz de permissões (RBAC/ACL)
- [x] **Arquiteto** — Diagramas C4 (Contexto, Containers, Componentes)
- [x] **Arquiteto** — ERD completo e integrações externas
- [x] **Arquiteto** — Spec Impact Matrix

## Fase 4: Geração 📝

- [x] **Redator** — Specs SDD por componente (22 modules: requirements/design/tasks/contracts)
- [x] **Redator** — OpenAPI (se aplicável) → `specs/openapi/appwrite.yaml`
- [x] **Redator** — User Stories (se aplicável) → `specs/user-stories/` (10 fluxos)
- [x] **Redator** — Code/Spec Matrix → `specs/traceability/code-spec-matrix.md`

## Fase 5: Revisão ✅

- [x] **Revisor** — Revisão cruzada de specs (Codex indisponível → auto-revisão + verificação no código; reclassificações + correções in-place)
- [x] **Revisor** — Resolução de lacunas com o usuário → `specs/questions.md` (3/3 resolvidas via opensrc + código)
- [x] **Revisor** — Relatório de confiança final → `specs/confidence-report.md` + `specs/gaps.md` (~93%)
- [x] **Revisor (2º passe)** — opensrc + Workers/: papéis de console, rate limits, máquina de estados de messaging, contagem de workers (14→12)

---

## Agentes Independentes

> Execute estes agentes quando os recursos estiverem disponíveis — podem rodar em qualquer fase.

- [ ] **Visor** — Análise de interface via screenshots
- [x] **Data Master** — Análise completa do banco de dados → `specs/database/` (48 coleções, 606 atributos, 233 índices; 25 campos criptografados). Fecha G5 e G6.
- [ ] **Design System** — Extração de tokens de design
- [ ] **Tracer** — Análise dinâmica (requer sistema acessível)

---

## Próximo passo

Após o Time de Descoberta concluir e o `_reversa_sdd/` estar populado, você pode disparar um dos fluxos seguintes:

- `/reversa-migrate`: orquestrador do **Time de Migração** (Paradigm Advisor → Curator → Strategist → Designer → Screen Translator → Inspector). Gera as specs do sistema novo. Saída em `_reversa_sdd/migration/` e `_reversa_sdd/screens/`.
- `/reversa-reconstructor`: gera plano bottom-up para reimplementar o software a partir das specs do legado (uma tarefa por sessão).
