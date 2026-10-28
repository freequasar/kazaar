package repositories

import (
	"encoding/json"
	"errors"

	"github.com/jmoiron/sqlx"
)

const GET_DEVICE_BY_DEVICE_ID_QUERY = `
SELECT
    id, device_id, platform, name, locale, jingle, spotter
FROM
    device
WHERE
    device_id = ?
LIMIT
    1
`

const UPDATE_NETWORK_INFO_QUERY = `
UPDATE device
SET network_info = ?
WHERE id = ?
`

type Device struct {
	ID          int    `db:"id"`
	DeviceID    string `db:"device_id"`
	Platform    string `db:"platform"`
	NetworkInfo string `db:"network_info"`
	Name        string `db:"name"`
	Locale      string `db:"locale"`
	Jingle      bool   `db:"jingle"`
	Spotter     string `db:"spotter"`
}

type NetworkInfo struct {
	ExternalPort int      `json:"external_port"`
	IPAddresses  []string `json:"ip_addresses"`
	MacAddresses []string `json:"mac_addresses"`
	Timestamp    int      `json:"ts"`
	WiFiSSID     string   `json:"wifi_ssid"`
}

type DeviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (deviceRepository *DeviceRepository) GetDeviceByDeviceID(deviceID string) (*Device, error) {
	var device Device

	if err := deviceRepository.db.Get(&device, GET_DEVICE_BY_DEVICE_ID_QUERY, deviceID); err != nil {
		return nil, err
	}

	return &device, nil
}

func (DeviceRepository *DeviceRepository) UpdateNetworkInfo(
	device *Device,
	networkInfo NetworkInfo,
) (*Device, error) {
	networkInfoBytes, err := json.Marshal(networkInfo)
	if err != nil {
		return nil, err
	}
	networkInfoString := string(networkInfoBytes)

	result, err := DeviceRepository.db.Exec(
		UPDATE_NETWORK_INFO_QUERY,
		networkInfoString, device.ID,
	)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("no device found")
	}

	device.NetworkInfo = networkInfoString
	return device, nil
}
