package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const UPDATES_DIR = "updates/"

func checkUpdates(w http.ResponseWriter, r *http.Request) {
	// get platform and version
	q := r.URL.Query()
	if !q.Has("platform") {
		reply(
			w, http.StatusBadRequest,
			Envelope{Status: "error", Message: "no platform"},
		)
		return
	}
	platform := q.Get("platform")
	version := q.Get("version")

	// get version and reply with it
	v, err := getVersionByPlatform(platform)
	if err != nil || v.Version == version {
		reply(w, http.StatusOK, CheckUpdatesResponse{HasUpdate: false})
		return
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	downloadUrl := fmt.Sprintf(
		"%s://%s/download_update/%s/%s",
		scheme, r.Host, url.PathEscape(v.Platform), url.PathEscape(v.Version),
	)
	reply(w, http.StatusOK, CheckUpdatesResponse{
		HasUpdate:   true,
		Version:     v.Version,
		DownloadUrl: downloadUrl,
		Crc32:       v.Crc32,
		IsCritical:  true,
		FallbackUrl: downloadUrl,
		FallbackIp:  "127.0.0.1",
		UpdateId:    strconv.FormatInt(v.ID, 10),
	})
}

func downloadUpdate(w http.ResponseWriter, r *http.Request) {
	// get parameters
	platform := chi.URLParam(r, "platform")
	version := chi.URLParam(r, "version")

	// get a version
	v, err := getVersion(platform, version)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != http.MethodHead {
			w.Write([]byte("update not found"))
		}
		return
	}

	// serve the update
	filename := filepath.Base(v.Version) + ".zip"
	filepath := path.Join(wd, "./versions", filename)
	info, err := os.Stat(filepath)
	if err != nil || info.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != http.MethodHead {
			w.Write([]byte("update not found"))
		}
		return
	}
	http.ServeFile(w, r, filepath)
}
