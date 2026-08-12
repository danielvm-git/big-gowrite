// Package jobs is the async dispatch backbone: enqueue work and consume it
// exactly once with SELECT … FOR UPDATE SKIP LOCKED. It replaces Swoole
// coroutines + Redis queues + the outbox pattern with a single Postgres table
// (BR-MIGRAR-014, BR-DESCARTAR-004; specs/migration/target_architecture.md).
package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Status values for a job's terminal-state lifecycle (ADR-008).
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusCanceled   = "canceled"
)

// Job is a unit of async work.
type Job struct {
	ID       int64
	Kind     string
	Payload  json.RawMessage
	Attempts int
}

// Queue enqueues and manages jobs against a pgx pool.
type Queue struct {
	pool *pgxpool.Pool
}

// New builds a Queue over the given pool.
func New(pool *pgxpool.Pool) *Queue { return &Queue{pool: pool} }

// EnqueueOptions configures an enqueue.
type EnqueueOptions struct {
	// DedupKey, if non-empty, makes the enqueue idempotent: a second enqueue
	// with the same key is a no-op and returns the existing job's id.
	DedupKey string
}

// Enqueue inserts a job. With a DedupKey, a duplicate is a no-op (idempotent).
// Returns the job id and whether a new row was inserted.
func (q *Queue) Enqueue(ctx context.Context, kind string, payload any, opts EnqueueOptions) (id int64, inserted bool, err error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return 0, false, fmt.Errorf("jobs: marshal payload: %w", err)
	}

	if opts.DedupKey != "" {
		// ON CONFLICT DO NOTHING makes the insert idempotent on dedup_key.
		err = q.pool.QueryRow(ctx, `
			INSERT INTO jobs (kind, payload, dedup_key)
			VALUES ($1, $2::jsonb, $3)
			ON CONFLICT (dedup_key) DO NOTHING
			RETURNING id`, kind, string(raw), opts.DedupKey).Scan(&id)
		if errors.Is(err, pgx.ErrNoRows) {
			// Conflict: fetch the existing job's id.
			if e := q.pool.QueryRow(ctx,
				`SELECT id FROM jobs WHERE dedup_key = $1`, opts.DedupKey).Scan(&id); e != nil {
				return 0, false, fmt.Errorf("jobs: fetch deduped id: %w", e)
			}
			return id, false, nil
		}
		if err != nil {
			return 0, false, fmt.Errorf("jobs: enqueue (dedup): %w", err)
		}
		return id, true, nil
	}

	err = q.pool.QueryRow(ctx, `
		INSERT INTO jobs (kind, payload)
		VALUES ($1, $2::jsonb)
		RETURNING id`, kind, string(raw)).Scan(&id)
	if err != nil {
		return 0, false, fmt.Errorf("jobs: enqueue: %w", err)
	}
	return id, true, nil
}

// Claim atomically claims up to one due pending job, marking it processing.
// Returns (nil, nil) when no job is available. Uses FOR UPDATE SKIP LOCKED so
// concurrent workers never claim the same row.
func (q *Queue) Claim(ctx context.Context, kinds ...string) (*Job, error) {
	// Build an optional kind filter.
	query := `
		UPDATE jobs SET status = 'processing', attempts = attempts + 1, updated_at = now()
		WHERE id = (
			SELECT id FROM jobs
			WHERE status = 'pending' AND run_after <= now()`
	args := []any{}
	if len(kinds) > 0 {
		query += ` AND kind = ANY($1)`
		args = append(args, kinds)
	}
	query += `
			ORDER BY id
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, kind, payload, attempts`

	var j Job
	err := q.pool.QueryRow(ctx, query, args...).Scan(&j.ID, &j.Kind, &j.Payload, &j.Attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("jobs: claim: %w", err)
	}
	return &j, nil
}

// Complete marks a job completed (a terminal state).
func (q *Queue) Complete(ctx context.Context, id int64) error {
	return q.setTerminal(ctx, id, StatusCompleted, "")
}

// Fail marks a job failed (a terminal state) with an error message.
func (q *Queue) Fail(ctx context.Context, id int64, reason string) error {
	return q.setTerminal(ctx, id, StatusFailed, reason)
}

// Cancel marks a job canceled (a terminal state).
func (q *Queue) Cancel(ctx context.Context, id int64) error {
	return q.setTerminal(ctx, id, StatusCanceled, "")
}

func (q *Queue) setTerminal(ctx context.Context, id int64, status, reason string) error {
	ct, err := q.pool.Exec(ctx, `
		UPDATE jobs SET status = $2, last_error = NULLIF($3,''), updated_at = now()
		WHERE id = $1`, id, status, reason)
	if err != nil {
		return fmt.Errorf("jobs: set %s: %w", status, err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("jobs: job %d not found", id)
	}
	return nil
}

// Status returns a job's current status (test/introspection helper).
func (q *Queue) Status(ctx context.Context, id int64) (string, error) {
	var s string
	err := q.pool.QueryRow(ctx, `SELECT status FROM jobs WHERE id = $1`, id).Scan(&s)
	return s, err
}
