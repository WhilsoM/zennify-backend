-- +goose Up
CREATE TABLE "user" (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT user_username_unique UNIQUE (username)
);

CREATE INDEX user_created_at_idx ON "user" (created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS "user";
