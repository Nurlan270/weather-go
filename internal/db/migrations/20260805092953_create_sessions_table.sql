-- +goose Up
CREATE TABLE IF NOT EXISTS sessions
(
    ID         VARCHAR PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS sessions;
