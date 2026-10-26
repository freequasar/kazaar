package utilities

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func DomainSubrouter(domains []string, r chi.Router) http.Handler {
	domainMap := make(map[string]bool)
	for _, d := range domains {
		domainMap[d] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}

		if domainMap[host] {
			r.ServeHTTP(w, req)
			return
		}

		http.NotFound(w, req)
	})
}

func ReplyJSON(w http.ResponseWriter, statusCode int, d any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(d)
}
