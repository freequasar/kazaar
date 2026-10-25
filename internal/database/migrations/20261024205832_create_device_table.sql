-- +goose Up
CREATE TABLE IF NOT EXISTS device (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT NOT NULL UNIQUE,
    platform TEXT NOT NULL
);

-- +goose Down
DROP TAEBLE IF EXISTS device;
