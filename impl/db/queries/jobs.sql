-- name: GetJob :one
SELECT id, kind, payload, status, attempts FROM jobs WHERE id = $1;

-- name: CountJobsByStatus :one
SELECT count(*) FROM jobs WHERE status = $1;
