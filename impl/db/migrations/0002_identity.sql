-- 0002_identity.sql: users, sessions, authenticators, challenges, tokens

CREATE TABLE IF NOT EXISTS users (
    id                   TEXT        PRIMARY KEY,
    name                 TEXT        NOT NULL DEFAULT '',
    email                TEXT,
    phone                TEXT,
    status               BOOLEAN     NOT NULL DEFAULT true,
    labels               TEXT[]      NOT NULL DEFAULT '{}',
    password             TEXT,                          -- AES-128-GCM envelope
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
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_email_uniq UNIQUE (email),
    CONSTRAINT users_phone_uniq UNIQUE (phone)
);

CREATE INDEX IF NOT EXISTS users_status_idx ON users (status);
CREATE INDEX IF NOT EXISTS users_registration_idx ON users (registration);
CREATE INDEX IF NOT EXISTS users_email_verification_idx ON users (email_verification);

CREATE TABLE IF NOT EXISTS sessions (
    id                           TEXT        PRIMARY KEY,
    user_id                      TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider                     TEXT        NOT NULL DEFAULT '',
    provider_uid                 TEXT        NOT NULL DEFAULT '',
    provider_access_token        TEXT,                  -- AES-128-GCM envelope
    provider_access_token_expiry TIMESTAMPTZ,
    provider_refresh_token       TEXT,                  -- AES-128-GCM envelope
    secret                       TEXT,                  -- AES-128-GCM envelope
    user_agent                   TEXT        NOT NULL DEFAULT '',
    ip                           TEXT        NOT NULL DEFAULT '',
    country_code                 TEXT        NOT NULL DEFAULT '',
    factors                      TEXT[]      NOT NULL DEFAULT '{}',
    expire                       TIMESTAMPTZ NOT NULL,
    mfa_updated_at               TIMESTAMPTZ,
    created_at                   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS sessions_user_id_idx ON sessions (user_id);
CREATE INDEX IF NOT EXISTS sessions_expire_idx ON sessions (expire);

CREATE TABLE IF NOT EXISTS authenticators (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        TEXT        NOT NULL,
    verified    BOOLEAN     NOT NULL DEFAULT false,
    data        TEXT,                           -- AES-128-GCM envelope (JSON)
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT authenticators_user_type_uniq UNIQUE (user_id, type)
);

CREATE INDEX IF NOT EXISTS authenticators_user_id_idx ON authenticators (user_id);

CREATE TABLE IF NOT EXISTS challenges (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        TEXT        NOT NULL,
    token       TEXT,                           -- AES-128-GCM envelope
    code        TEXT,                           -- AES-128-GCM envelope
    expire      TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS challenges_user_id_idx ON challenges (user_id);
CREATE INDEX IF NOT EXISTS challenges_expire_idx ON challenges (expire);

CREATE TABLE IF NOT EXISTS tokens (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        INT         NOT NULL,
    secret      TEXT,                           -- AES-128-GCM envelope
    expire      TIMESTAMPTZ NOT NULL,
    user_agent  TEXT        NOT NULL DEFAULT '',
    ip          TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS tokens_user_id_idx ON tokens (user_id);
CREATE INDEX IF NOT EXISTS tokens_type_expire_idx ON tokens (type, expire);
