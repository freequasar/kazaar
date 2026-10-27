package quasar

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/freequasar/kazaar/internal/repositories"
	"github.com/freequasar/kazaar/internal/utilities"
	"github.com/go-chi/chi/v5"
)

func (quasarHandler *QuasarHandler) buildDownloadUrl(r *http.Request, manfiest *repositories.VersionManfiest) string {
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

func (quasarHandler *QuasarHandler) CheckUpdates(w http.ResponseWriter, r *http.Request) {
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
		platform = "yandexstation"
	}

	clientVersion := q.Get("version")
	manfiest, err := quasarHandler.versionManfiestRepository.GetVersionManfiestByPlatform(platform)
	if err != nil || manfiest.Version == clientVersion {
		utilities.ReplyJSON(w, http.StatusOK, CheckUpdatesResponse{HasUpdate: false})
		return
	}

	downloadUrl := quasarHandler.buildDownloadUrl(r, manfiest)
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

func (quasarHandler *QuasarHandler) DownloadUpdate(w http.ResponseWriter, r *http.Request) {
	if !slices.Contains([]string{http.MethodGet, http.MethodHead}, r.Method) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	platform := chi.URLParam(r, "platform")
	version := chi.URLParam(r, "version")

	manfiestExists, err := quasarHandler.versionManfiestService.CheckForManfiestAvailability(
		platform, version,
	)
	if err != nil || !manfiestExists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("update not found"))
		return
	}

	fileName := filepath.Base(version) + ".zip"
	filePath := path.Join(quasarHandler.versionsDir, fileName)

	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
