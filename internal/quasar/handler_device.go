package quasar

import (
	"net/http"

	"github.com/freequasar/kazaar/internal/utilities"
)

func (quasarHandler *QuasarHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
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
	serverVersion, err := quasarHandler.versionManfiestRepository.GetVersionForPlatform(platform)
	if err != nil {
		utilities.ReplyJSON(w, http.StatusForbidden, EnvelopeResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if err := quasarHandler.deviceService.CreateDevice(deviceId, platform); err != nil {
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
