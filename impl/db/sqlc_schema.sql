CREATE TABLE users (
    id                   TEXT        PRIMARY KEY,
    name                 TEXT        NOT NULL DEFAULT '',
    email                TEXT,
    phone                TEXT,
    status               BOOLEAN     NOT NULL DEFAULT true,
    labels               TEXT[]      NOT NULL DEFAULT '{}',
    password             TEXT,
    hash                 TEXT        NOT NULL DEFAULT 'argon2',
    hash_options         JSONB       NOT NULL DEFAULT '{}',
    password_update      TIMESTAMPTZ,
    prefs                JSONB       NOT NULL DEFAULT '{}',
    registration         TIMESTAMPTZ NOT NULL DEFAULT now(),
    email_verification   BOOLEAN     NOT NULL DEFAULT false,
    phone_verification   BOOLEAN     NOT NULL DEFAULT false,
    reset                BOOLEAN     NOT NULL DEFAULT false,
    mfa                  BOOLEAN     NOT NULL DEFAULT false,
    accessed_at          TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id                           TEXT        PRIMARY KEY,
    user_id                      TEXT        NOT NULL REFERENCES users(id),
    provider                     TEXT        NOT NULL DEFAULT '',
    provider_uid                 TEXT        NOT NULL DEFAULT '',
    provider_access_token        TEXT,
    provider_access_token_expiry TIMESTAMPTZ,
    provider_refresh_token       TEXT,
    secret                       TEXT,
    user_agent                   TEXT        NOT NULL DEFAULT '',
    ip                           TEXT        NOT NULL DEFAULT '',
    country_code                 TEXT        NOT NULL DEFAULT '',
    factors                      TEXT[]      NOT NULL DEFAULT '{}',
    expire                       TIMESTAMPTZ NOT NULL,
    mfa_updated_at               TIMESTAMPTZ,
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE authenticators (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL REFERENCES users(id),
    type        TEXT        NOT NULL,
    verified    BOOLEAN     NOT NULL DEFAULT false,
    data        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE challenges (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL REFERENCES users(id),
    type        TEXT        NOT NULL,
    token       TEXT,
    code        TEXT,
    expire      TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tokens (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL REFERENCES users(id),
    type        INT         NOT NULL,
    secret      TEXT,
    expire      TIMESTAMPTZ NOT NULL,
    user_agent  TEXT        NOT NULL DEFAULT '',
    ip          TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE jobs (
    id          BIGSERIAL   PRIMARY KEY,
    kind        TEXT        NOT NULL,
    payload     JSONB       NOT NULL DEFAULT '{}'::jsonb,
    dedup_key   TEXT        UNIQUE,
    status      TEXT        NOT NULL DEFAULT 'pending',
    attempts    INT         NOT NULL DEFAULT 0,
    max_attempts INT        NOT NULL DEFAULT 5,
    last_error  TEXT,
    run_after   TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
