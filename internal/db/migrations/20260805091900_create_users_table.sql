-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    ID       SERIAL PRIMARY KEY,
    login    VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
