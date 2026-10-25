package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const GET_VERSION_MANFIEST_BY_PLATFORM_QUERY = `
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

const GET_VERSION_FOR_PLATFORM_QUERY = "SELECT version FROM version WHERE platform = ? LIMIT 1"

type VersionManfiest struct {
	ID                    int64   `db:"id"`
	Platform              string  `db:"platform"`
	Version               string  `db:"version"`
	Crc32                 int     `db:"crc32"`
	UpdateScript          *string `db:"update_script"`
	UpdateScriptSignature *string `db:"update_script_signature"`
	UpdateScriptTimestamp *int64  `db:"update_script_timestamp"`
}

type VersionManfiestRepository struct {
	db *sqlx.DB
}

func NewVersionManfiestRepository(db *sqlx.DB) *VersionManfiestRepository {
	return &VersionManfiestRepository{db: db}
}

func (h *VersionManfiestRepository) GetVersionManfiestByPlatform(platform string) (*VersionManfiest, error) {
	var version VersionManfiest

	if err := h.db.Get(&version, GET_VERSION_MANFIEST_BY_PLATFORM_QUERY, platform); err != nil {
		return nil, err
	}

	return &version, nil
}

func (h *VersionManfiestRepository) GetVersionForPlatform(platform string) (*string, error) {
	var version string

	if err := h.db.Get(&version, GET_VERSION_FOR_PLATFORM_QUERY, platform); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("unknown error: %s", err.Error())
		}

		version = ""
	}

	return &version, nil
}
