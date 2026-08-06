-- +goose Up
CREATE TABLE IF NOT EXISTS locations
(
    ID        SERIAL PRIMARY KEY,
    name      VARCHAR NOT NULL,
    user_id   INTEGER REFERENCES users (id) ON DELETE CASCADE,
    latitude  DECIMAL NOT NULL,
    longitude DECIMAL NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS locations;
