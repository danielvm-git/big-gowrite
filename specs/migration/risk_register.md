---
schemaVersion: 1
generatedAt: 2026-08-11T16:23:34Z
reversa:
  version: "1.2.57"
kind: risk_register
producedBy: strategist
hash: "sha256:e778057f19799b989d78d170d67772fda9d72b3218af3ffc046357dc92d3777d"
---

# Risk Register

> Riscos da estratégia recomendada (Strangler Fig + Parallel Run), da mudança de paradigma, de dados, operacionais e organizacionais.
> P = probabilidade, I = impacto (baixo/médio/alto). Owner = papel (solo-dev acumula todos, explicitado por função).

| ID | Risco | P | I | Mitigação | Contingência | Owner |
|----|-------|---|---|-----------|--------------|-------|
| **PARADIGM** | | | | | | |
| R-P1 | Erros idiomáticos ao traduzir exceções → valores de erro (perda de semântica de status HTTP) | Média | Alto | Inspector mapeia cada código de erro legado→Go; testes de paridade por endpoint | Reverter rota do módulo ao PHP até corrigir | dev (backend) |
| R-P2 | Curva de aprendizado Go (goroutines/channels, sem herança) para solo-dev | Alta | Médio | Migrar bottom-up começando por módulos simples (Health, Avatars) para ganhar fluência; princípios bigpowers | Estender prazo (sem deadline) | dev |
| R-P3 | Concorrência: reimplementar idempotência/estado terminal (ADR-008) com goroutines diverge do comportamento Swoole | Média | Alto | Parallel Run nos workers; testes de idempotência por CloudEvent id | Manter worker PHP correspondente até paridade | dev (workers) |
| **DADOS** | | | | | | |
| R-D1 | Consolidação Mongo+MariaDB → PostgreSQL: perda/deformação de dados (documentos → JSONB, relationships espelhados) | Média | Alto | ETL com validação linha-a-linha; Data Master já mapeou 48 coleções/606 atributos como base | Re-executar ETL a partir de snapshot; Postgres em staging antes do flip | dev (dados) |
| R-D2 | Consistência durante coexistência PHP↔Go (escritas em dois stacks no período Strangler) | Alta | Alto | Migrar por módulo com **ownership exclusivo do dado** por módulo (um stack é dono da coleção por vez); sem escrita dupla | Congelar escrita do módulo durante o flip; reconciliação pós-flip | dev (dados) |
| R-D3 | Re-criptografia dos 25 campos `encrypt` (chaves, esquema AES) na troca de engine | Média | Alto | Migrar chave-mestra e algoritmo AES-GCM idênticos; testar decrypt de amostra antes do flip | Manter campos legíveis só pelo PHP até validar decrypt em Go | dev (segurança) |
| **OPERACIONAL** | | | | | | |
| R-O1 | Serviços externos (executor, browser, assistant, ClamAV) acoplados ao contrato PHP | Média | Médio | Manter serviços externos na fase 1 (BR-HUMANA-003); Go fala o mesmo protocolo | Adaptador de compatibilidade Go↔executor | dev (infra) |
| R-O2 | Rodar dois stacks aumenta footprint **durante** a migração (contrário ao objetivo, temporariamente) | Alta | Baixo | Aceitar sobrecusto temporário; decomissionar PHP por módulo assim que flipado | — | dev (infra) |
| R-O3 | Paridade de contrato de API para o console SPA preservado | Média | Alto | `specs/openapi/appwrite.yaml` como oráculo de contrato; testes de contrato no Parallel Run | Ajustar handler Go até o contrato bater | dev (backend) |
| **ORGANIZACIONAL** | | | | | | |
| R-G1 | Solo-dev é gargalo único; ritmo limitado pela cota de IA | Alta | Médio | Trabalho fatiado por módulo; estado persistido a cada módulo (Reversa checkpoints) | Pausar/retomar sem perda (`/reversa-migrate --resume`) | dev |
| R-G2 | Escopo total (22 módulos) pode exaurir motivação antes do fim | Média | Médio | Entregar valor incremental (módulos flipados já rodam em Go); marcos visíveis | Reduzir escopo para subconjunto de hosting-de-sites (core do brief) | dev (produto) |

## Riscos críticos (destaque)
- **R-D2** (consistência na coexistência) e **R-D1/R-D3** (ETL + re-cripto) são os de maior severidade — a consolidação de banco é a operação mais arriscada da migração.
- **R-P1/R-P3** (semântica de erro e concorrência) são os riscos diretos do paradigma transformacional — endereçados pelo Parallel Run e pelo Inspector.
