# spec-id: parity-02-team-invite
# origem: code-analysis.md §3.18, state-machines.md §5 (Membership), redteam S2
# target: target_domain_model.md (Organizations + Identity co-flip)
# paradigma: cross-context write numa requisição — Identity+Organizations flipam juntos
@paridade @critico
Funcionalidade: Convite de time cria usuário stub e sessão (cross-context)

  Cenário: Convite de e-mail não registrado cria stub e membership invited
    Dado um time com pelo menos um owner
    E um e-mail que não corresponde a nenhum usuário
    Quando um owner cria um convite para esse e-mail
    Então um usuário stub é criado no contexto Identity
    E uma membership é criada com status "invited"
    E a resposta tem o mesmo shape do legado (campos + total)

  @critico
  Cenário: Confirmação de convite valida secret e cria sessão
    Dado uma membership "invited" com secret válido
    Quando o convidado confirma com o proofForToken correto
    Então a membership passa a "confirmed"
    E uma sessão é criada para o convidado (cookie a_session_console)
    E a escrita atinge users+sessions (Identity) e memberships (Organizations) atomicamente

  Cenário: Downgrade do último owner é bloqueado
    Dado um time com exatamente um owner
    Quando se tenta remover ou rebaixar esse owner
    Então a operação falha com erro type "MEMBERSHIP_DOWNGRADE_PROHIBITED"
    E o status HTTP é idêntico ao legado
