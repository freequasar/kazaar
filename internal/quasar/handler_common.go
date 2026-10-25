package quasar

import (
	"net/http"
	"strconv"
	"time"

	"github.com/freequasar/kazaar/internal/utilities"
)

func (quasarHandler *QuasarHandler) Generate204(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (quasarHandler *QuasarHandler) GetTimezone(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if !q.Has("lat") || !q.Has("lon") {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "invalid request",
		})
	}

	var err error
	latitude, err := strconv.ParseFloat(q.Get("lat"), 64)
	if err != nil {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "invalid latitude",
		})
		return
	}
	longitude, err := strconv.ParseFloat(q.Get("lon"), 64)
	if err != nil {
		utilities.ReplyJSON(w, http.StatusBadRequest, EnvelopeResponse{
			Status:  "error",
			Message: "invalid longitude",
		})
		return
	}

	timezone := quasarHandler.timezoneFinder.GetTimezoneName(longitude, latitude)
	location, _ := time.LoadLocation(timezone)
	_, zone := time.Now().In(location).Zone()

	utilities.ReplyJSON(w, http.StatusOK, GetTimezoneResponse{
		Status:          "ok",
		Timezone:        timezone,
		OffsetInSeconds: zone,
	})
}
