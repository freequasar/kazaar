package quasar

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/freequasar/kazaar/internal/utilities"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

const INSERT_DEVICE_QUERY = "INSERT INTO device (device_id, platform) VALUES (?, ?)"

const GET_VERSION_FOR_PLATFORM_QUERY = "SELECT version FROM version WHERE platform = ? LIMIT 1"

func (h *QuasarHandler) registerDevice(deviceId string, platform string) error {
	_, err := h.db.Exec(INSERT_DEVICE_QUERY, deviceId, platform)

	if err != nil {
		var error *sqlite.Error

		if errors.As(err, &error) {
			if error.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return errors.New("device must be unique")
			}

			return fmt.Errorf("sqlite3 error: %s", error.Error())
		}

		return fmt.Errorf("unknown error: %s", err.Error())
	}

	return nil
}

func (h *QuasarHandler) getVersionForPlatform(platform string) (*string, error) {
	var version string

	if err := h.db.Get(&version, GET_VERSION_FOR_PLATFORM_QUERY, platform); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("unknown error: %s", err.Error())
		}

		version = ""
	}

	return &version, nil
}

func (h *QuasarHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if !q.Has("device_id") || !q.Has("platform") {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "invalid request",
		})
		return
	}

	deviceId := q.Get("device_id")
	if deviceId == "" {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "device_id is empty",
		})
		return
	}
	platform := q.Get("platform")
	if platform == "" {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "platform is empty",
		})
		return
	}

	clientVersion := q.Get("version")
	serverVersion, err := h.getVersionForPlatform(platform)
	if err != nil {
		utilities.ReplyJSON(w, http.StatusForbidden, EnvelopeResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if err := h.registerDevice(deviceId, platform); err != nil {
		utilities.ReplyJSON(w, http.StatusForbidden, EnvelopeResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	utilities.ReplyJSON(w, http.StatusOK, RegisterResponse{
		Status:            "ok",
		HasCriticalUpdate: clientVersion == *serverVersion,
	})
}
