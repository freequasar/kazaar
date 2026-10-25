package services

import (
	"database/sql"
	"errors"

	"github.com/freequasar/kazaar/internal/repositories"
	"github.com/jmoiron/sqlx"
)

const CHECK_FOR_MANFIEST_AVAILABILITY_QUERY = "SELECT 1 FROM version WHERE platform = ? AND version = ? LIMIT 1"

type VersionManfiestService struct {
	versionManfiestRepository *repositories.VersionManfiestRepository
	db                        *sqlx.DB
}

func NewVersionManfiestService(
	versionManfiestRepository *repositories.VersionManfiestRepository,
	db *sqlx.DB,
) *VersionManfiestService {
	return &VersionManfiestService{
		versionManfiestRepository: versionManfiestRepository,
		db:                        db,
	}
}

func (versionManfiestService *VersionManfiestService) CheckForManfiestAvailability(
	platform string, version string,
) (bool, error) {
	var exists int

	err := versionManfiestService.db.Get(
		&exists, CHECK_FOR_MANFIEST_AVAILABILITY_QUERY, platform, version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists == 1, nil
}
