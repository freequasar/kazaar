package services

import (
	"errors"

	"github.com/freequasar/kazaar/internal/repositories"
	"github.com/jmoiron/sqlx"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

const CREATE_DEVICE_QUERY = "INSERT INTO device (device_id, platform) VALUES (?, ?)"

type DeviceService struct {
	deviceRepository *repositories.DeviceRepository
	db               *sqlx.DB
}

func NewDeviceService(deviceRepository *repositories.DeviceRepository, db *sqlx.DB) *DeviceService {
	return &DeviceService{deviceRepository: deviceRepository, db: db}
}

func (deviceService *DeviceService) CreateDevice(deviceId string, platform string) error {
	_, err := deviceService.db.Exec(
		CREATE_DEVICE_QUERY, deviceId, platform,
	)

	if err != nil {
		var error *sqlite.Error

		if errors.As(err, &error) {
			if error.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return errors.New("device must be unique")
			}
		}

		return err
	}

	return nil
}
