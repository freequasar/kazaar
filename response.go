package main

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RegisterResponse struct {
	Status            string `json:"status"`
	HasCriticalUpdate bool   `json:"hasCriticalUpdate"`
}

type TimezoneResponse struct {
	Status          string `json:"status"`
	Timezone        string `json:"timezone"`
	OffsetInSeconds int    `json:"offset_sec"`
}

func reply(w http.ResponseWriter, statusCode int, d any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(d)
}
