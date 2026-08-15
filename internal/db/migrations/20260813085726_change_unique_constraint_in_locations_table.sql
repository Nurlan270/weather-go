-- +goose Up
ALTER TABLE locations
    DROP CONSTRAINT uq_locations_user_id_name;

ALTER TABLE locations
    ADD CONSTRAINT uq_locations_user_id_lat_lon UNIQUE (user_id, latitude, longitude);

-- +goose Down
ALTER TABLE locations
    DROP CONSTRAINT uq_locations_user_id_lat_lon;

ALTER TABLE locations
    ADD CONSTRAINT uq_locations_user_id_name UNIQUE (user_id, name);
