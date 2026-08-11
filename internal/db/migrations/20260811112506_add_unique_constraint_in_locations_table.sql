-- +goose Up
ALTER TABLE locations
    ADD CONSTRAINT uq_locations_user_id_name UNIQUE (user_id, name);

-- +goose Down
ALTER TABLE locations
    DROP CONSTRAINT uq_locations_user_id_name;
