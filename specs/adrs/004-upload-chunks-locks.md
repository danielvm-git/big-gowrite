# ADR-004 — Upload multipart em chunks com distributed locks

- **Status:** Aceito
- **Data:** Retrospectiva 2026-08-05
- **Confiança:** 🟢 CONFIRMADO

## Contexto
Storage, Functions e Sites precisam receber arquivos grandes (até centenas de MB) sobre HTTP sem buffers gigantes, com **resumibilidade** e **concorrência controlada** no mesmo recurso.

## Decisão
- Upload via `Content-Range` em chunks (coleção de destino guarda `chunksTotal`/`chunksUploaded`).
- **Distributed lock** por recurso durante o upload — concorrência → `LockContention` (429 com `Retry-After` em Sites).
- Processamento de finalize: compressão (gz/zstd) + criptografia (AES-128-GCM) + scan ClamAV — aplicados apenas quando `sizeOriginal ≤ buffer` (upload otimizado).
- Sites: lock de 600s; Storage: download/view com suporte a HTTP Range.

## Alternativas consideradas
- Upload síncrono de arquivo inteiro (rejeitado: limite de memória e sem resumibilidade).
- Multipart único sem lock (rejeitado: uploads concorrentes corromperiam o mesmo documento).

## Consequências
- Resumibilidade e paralelismo de chunks.
- Complexidade de contrato (`Content-Range`, erros `LockContention`).
- Reuso do mesmo padrão em 3 módulos (Storage/Functions/Sites).

## Evidências
- Módulos Storage, Functions, Sites (upload em chunks).
- `domain.md` §2.4; `modules.json` (algorithms de Storage/Functions/Sites).
