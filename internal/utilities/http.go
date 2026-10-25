package utilities

import (
	"encoding/json"
	"net/http"
)

func ReplyJSON(w http.ResponseWriter, statusCode int, d any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(d)
}
