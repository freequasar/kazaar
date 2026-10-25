package repositories

import "github.com/jmoiron/sqlx"

type Device struct {
	ID       int    `db:"id"`
	DeviceID string `db:"device_id"`
	Platform string `db:"platform"`
}

type DeviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}
