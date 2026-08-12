-- 0001_base.sql — gowrite base schema + async dispatch backbone.
-- PostgreSQL only. Replaces Mongo+MariaDB+Postgres+Redis (BR-HUMANA-002).
-- pgvector is added in a later migration (e05); not required here.

-- ── jobs: single async dispatch table (BR-MIGRAR-014, ADR-008) ──────────────
-- Consumed with SELECT … FOR UPDATE SKIP LOCKED. Every job drives to a terminal
-- state (completed|failed|canceled). Idempotency via dedup_key.
CREATE TABLE IF NOT EXISTS jobs (
    id          BIGSERIAL   PRIMARY KEY,
    kind        TEXT        NOT NULL,
    payload     JSONB       NOT NULL DEFAULT '{}'::jsonb,
    dedup_key   TEXT        UNIQUE,
    status      TEXT        NOT NULL DEFAULT 'pending'
                            CHECK (status IN ('pending','processing','completed','failed','canceled')),
    attempts    INT         NOT NULL DEFAULT 0,
    max_attempts INT        NOT NULL DEFAULT 5,
    last_error  TEXT,
    run_after   TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Partial index for the hot consumer query (pending & due).
CREATE INDEX IF NOT EXISTS jobs_claimable_idx
    ON jobs (run_after)
    WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS jobs_kind_status_idx ON jobs (kind, status);

-- ── realtime: LISTEN/NOTIFY fanout (replaces Redis pub/sub) ─────────────────
-- A generic trigger raises a NOTIFY on channel 'gowrite_realtime' whenever a
-- row is written to a realtime-enabled table. Contexts opt in by attaching it.
CREATE OR REPLACE FUNCTION notify_realtime() RETURNS trigger AS $$
DECLARE
    payload JSON;
BEGIN
    payload := json_build_object(
        'table', TG_TABLE_NAME,
        'op', TG_OP,
        'id', COALESCE(NEW.id, OLD.id)
    );
    PERFORM pg_notify('gowrite_realtime', payload::text);
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;
