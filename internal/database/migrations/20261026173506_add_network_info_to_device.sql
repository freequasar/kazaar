-- +goose Up
ALTER TABLE device
ADD COLUMN network_info TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE
DROP COLUMN network_info;
