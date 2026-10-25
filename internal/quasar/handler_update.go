package quasar

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/freequasar/kazaar/internal/utilities"
	"github.com/go-chi/chi/v5"
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

const CHECK_FOR_MANFIEST_AVAILABILITY_QUERY = "SELECT 1 FROM version WHERE platform = ? AND version = ? LIMIT 1"

type VersionManfiest struct {
	ID                    int64   `db:"id"`
	Platform              string  `db:"platform"`
	Version               string  `db:"version"`
	Crc32                 int     `db:"crc32"`
	UpdateScript          *string `db:"update_script"`
	UpdateScriptSignature *string `db:"update_script_signature"`
	UpdateScriptTimestamp *int64  `db:"update_script_timestamp"`
}

func (h *QuasarHandler) getVersionManfiestByPlatform(platform string) (*VersionManfiest, error) {
	var version VersionManfiest

	if err := h.db.Get(&version, GET_VERSION_MANFIEST_BY_PLATFORM_QUERY, platform); err != nil {
		return nil, err
	}

	return &version, nil
}

func (h *QuasarHandler) checkForManfiestAvailability(platform string, version string) (bool, error) {
	var exists int

	if err := h.db.Get(&exists, CHECK_FOR_MANFIEST_AVAILABILITY_QUERY, platform, version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return exists == 1, nil
}

func (h *QuasarHandler) buildDownloadUrl(r *http.Request, manfiest *VersionManfiest) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	host := r.Host
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	}

	platform := url.PathEscape(manfiest.Platform)
	version := url.PathEscape(manfiest.Version)

	return fmt.Sprintf("%s://%s/download_update/%s/%s", scheme, host, platform, version)
}

func (h *QuasarHandler) CheckUpdates(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if !q.Has("platform") {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "no platform",
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
	manfiest, err := h.getVersionManfiestByPlatform(platform)
	if err != nil || manfiest.Version == clientVersion {
		utilities.ReplyJSON(w, http.StatusOK, CheckUpdatesResponse{HasUpdate: false})
		return
	}

	downloadUrl := h.buildDownloadUrl(r, manfiest)
	utilities.ReplyJSON(w, http.StatusOK, CheckUpdatesResponse{
		HasUpdate:   true,
		Version:     manfiest.Version,
		DownloadUrl: downloadUrl,
		Crc32:       manfiest.Crc32,
		IsCritical:  true,
		FallbackUrl: downloadUrl,
		FallbackIp:  "127.0.0.1",
		UpdateId:    strconv.FormatInt(manfiest.ID, 10),
	})
}

func (h *QuasarHandler) DownloadUpdate(w http.ResponseWriter, r *http.Request) {
	if !slices.Contains([]string{http.MethodGet, http.MethodHead}, r.Method) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	platform := chi.URLParam(r, "platform")
	version := chi.URLParam(r, "version")

	manfiestExists, err := h.checkForManfiestAvailability(platform, version)
	if err != nil || !manfiestExists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("update not found"))
		return
	}

	fileName := filepath.Base(version) + ".zip"
	filePath := path.Join(h.versionsDir, fileName)

	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
