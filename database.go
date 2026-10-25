package main

const INSERT_DEVICE_QUERY = "INSERT INTO device (device_id, platform) VALUES (?, ?)"

const GET_VERSION_QUERY = `
SELECT
    id, platform, version,
    crc32, update_script,
    update_script_signature,
    update_script_timestamp
FROM
    version
WHERE
    platform = ? AND version = ?
LIMIT 1
`

const GET_VERSION_BY_PLATFORM_QUERY = `
SELECT
    id, platform, version,
    crc32, update_script,
    update_script_signature,
    update_script_timestamp
FROM
    version
WHERE
    platform = ?
LIMIT 1
`

type Device struct {
	ID       int64  `db:"id"`
	DeviceID string `db:"device_id"`
	Platform string `db:"platform"`
}

type Version struct {
	ID                    int64   `db:"id"`
	Platform              string  `db:"platform"`
	Version               string  `db:"version"`
	Crc32                 int     `db:"crc32"`
	UpdateScript          *string `db:"update_script"`
	UpdateScriptSignature *string `db:"update_script_signature"`
	UpdateScriptTimestamp *int64  `db:"update_script_timestamp"`
}

func createDevice(deviceId string, platform string) (*Device, error) {
	res, err := db.Exec(INSERT_DEVICE_QUERY, deviceId, platform)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Device{ID: id, DeviceID: deviceId, Platform: platform}, nil
}

func getVersion(platform string, version string) (*Version, error) {
	var v Version
	if err := db.Get(&v, GET_VERSION_QUERY, platform, version); err != nil {
		return nil, err
	}
	return &v, nil
}

func getVersionByPlatform(platform string) (*Version, error) {
	var v Version
	if err := db.Get(&v, GET_VERSION_BY_PLATFORM_QUERY, platform); err != nil {
		return nil, err
	}
	return &v, nil
}
