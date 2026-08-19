CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    avatar      VARCHAR(255),
    provider    VARCHAR(50)  NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT users_email_unique       UNIQUE (email),
    CONSTRAINT users_provider_id_unique UNIQUE (provider_id)
);
