package quasar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/freequasar/kazaar/internal/repositories"
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

func (quasarHandler *QuasarHandler) UpdateDeviceState(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if !q.Has("device_id") || !q.Has("platform") || !q.Has("field") {
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
			Message: "invalid device id",
		})
		return
	}
	platform := q.Get("platform")
	if platform == "" {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "invalid platform",
		})
		return
	}

	device, err := quasarHandler.deviceRepository.GetDeviceByDeviceID(deviceId)
	if err != nil || device.Platform != platform {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: fmt.Sprintf("unknown device"),
		})
		return
	}

	field := q.Get("field")
	switch field {
	case "networkInfo":
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(strings.ToLower(ct), "application/json") {
			utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
				Status:  "error",
				Message: "body must be application/json",
			})
			return
		}

		var networkInfo repositories.NetworkInfo
		if err := json.NewDecoder(r.Body).Decode(&networkInfo); err != nil {
			utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
				Status:  "error",
				Message: "invalid body",
			})
			return
		}

		if _, err := quasarHandler.deviceRepository.UpdateNetworkInfo(device, networkInfo); err != nil {
			utilities.ReplyJSON(w, http.StatusInternalServerError, EnvelopeResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
	default:
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "unknown field",
		})
		return
	}

	utilities.ReplyJSON(w, http.StatusOK, EnvelopeResponse{
		Status:  "ok",
		Message: "",
	})
}
