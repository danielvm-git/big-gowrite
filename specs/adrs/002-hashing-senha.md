# ADR-002 — Hashing de senha: Argon2id default + imports legados

- **Status:** Aceito
- **Data:** Retrospectiva 2026-08-05
- **Confiança:** 🟢 CONFIRMADO

## Contexto
O Appwrite recebe usuários de múltiplas fontes. Senhas precisam de hash moderno (Argon2id) por padrão, mas o admin deve conseguir importar bases legadas com hashes de outros algoritmos sem forçar reset de senha.

## Decisão
- Hash padrão: **Argon2id** (PHP `password_hash`).
- Import de hashes legados no módulo Users: **Bcrypt, MD5, Argon2, SHA (1/256/512), PHPass, Scrypt, Scrypt-Modified** — com opção `force`/`userId` para override de hash existente.
- Unicidade de e-mail entre `users` **e** `identities` (OAuth2) para evitar colisão de identidade.

## Alternativas consideradas
- Apenas Argon2id (rejeitado: bloqueia migração de bases legadas).
- Suporte apenas bcrypt (rejeitado: matriz de origem mais ampla no ecossistema).

## Consequências
- Segurança: hashes fracos legados persistem até a próxima troca (política `passwordHistory`/expiração mitigam).
- Complexidade de validação por algoritmo (verificador específico por tipo).

## Evidências
- Módulo `Users` (Base.php, Http) — imports de hash com tipos listados.
- `domain.md` §2.1 (regras de auth).
