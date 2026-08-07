-- +goose Up
CREATE TABLE IF NOT EXISTS locations
(
    ID        SERIAL PRIMARY KEY,
    user_id   INTEGER REFERENCES users (id) ON DELETE CASCADE,
    name      VARCHAR NOT NULL,
    latitude  DECIMAL NOT NULL,
    longitude DECIMAL NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS locations;
