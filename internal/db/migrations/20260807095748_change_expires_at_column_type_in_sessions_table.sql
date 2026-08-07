-- +goose Up
ALTER TABLE sessions
    ALTER COLUMN expires_at TYPE timestamptz USING expires_at AT TIME ZONE 'UTC';

-- +goose Down
ALTER TABLE sessions
    ALTER COLUMN expires_at TYPE timestamp USING expires_at AT TIME ZONE 'UTC';