# spec-id: parity-04-realtime-and-shape
# origem: redteam contract-parity CRITICAL-0/CRITICAL-1, screen_deviation_log DEV-001/002
# target: screen_modernization_decision.md (oráculo corrigido)
# nota: WS realtime e shape de resposta estão FORA do OpenAPI — testes dedicados
@paridade @critico
Funcionalidade: Paridade de protocolo realtime e shape de resposta

  Cenário: Shape de resposta de listagem bate com o legado
    Dado uma listagem (ex. GET /v1/users) com resultados
    Quando o Go responde
    Então a resposta contém o campo "total" com a mesma semântica do legado
    E cada item tem "$id" e "$createdAt" e os campos de modelo do legado
    E o console consegue renderizar sem campos ausentes

  @critico
  Cenário: Subscription realtime entrega o envelope esperado
    Dado um cliente que abre WS em /v1/realtime com project e channels[]
    Quando um recurso subscrito muda
    Então o cliente recebe um envelope {events, channels, timestamp, payload}
    E as strings de evento seguem o padrão do legado (ex. databases.*.collections.*)

  Cenário: Stream do assistant é texto cru (não SSE)
    Dado o console chama /console/assistant com streamProtocol text
    Quando o Go responde
    Então o corpo é um stream de texto concatenado
    E NÃO usa framing SSE "data:" (que quebraria a UI)
