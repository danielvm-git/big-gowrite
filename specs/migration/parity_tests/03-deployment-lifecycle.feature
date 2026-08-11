# spec-id: parity-03-deployment-lifecycle
# origem: state-machines.md §1 (Deployment), ADR-003, ADR-008
# target: target_architecture.md (jobs+SKIP LOCKED), target_domain_model.md (Compute)
# paradigma: filas Redis → jobs Postgres — exige idempotência + ordem + estado terminal
@paridade @idempotencia @ordem @critico
Funcionalidade: Ciclo de build/deployment sob a fila jobs

  Cenário: Deployment percorre a máquina de estados
    Dado um upload de deployment
    Quando o worker processa o build com sucesso
    Então o status transiciona waiting → processing → building → ready
    E "ready" é pré-condição para execução (senão BUILD_NOT_READY)

  @idempotencia
  Cenário: Callback de build duplicado não duplica efeito
    Dado um build cujo callback é entregue duas vezes (mesmo CloudEvent id)
    Quando o worker consome ambos via jobs (dedup_key)
    Então o estado final é aplicado exatamente uma vez
    E não há dupla transição nem duplo evento de deployment.ready

  @ordem
  Cenário: Cancelamento respeita estado terminal
    Dado um deployment em "building"
    Quando um cancel chega
    Então o status vira "canceled" (persistido)
    Mas se o deployment já está "ready" ou "failed"
    Então o cancel falha com "BUILD_ALREADY_COMPLETED"

  Cenário: Job nunca fica preso em processing (ADR-008)
    Dado um worker que falha no meio do processamento
    Quando a limpeza roda
    Então a mensagem termina em um estado terminal (done|failed), nunca preso em processing
