-- +goose Up
CREATE TABLE IF NOT EXISTS version (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    platform TEXT NOT NULL UNIQUE,
    version TEXT NOT NULL UNIQUE,
    crc32 INTEGER NOT NULL,
    update_script TEXT DEFAULT NULL,
    update_script_signature TEXT DEFAULT NULL,
    update_script_timestamp INTEGER DEFAULT NULL
);

-- +goose Down
DROP TABLE IF EXISTS version;
