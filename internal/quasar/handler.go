package quasar

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/ringsaturn/tzf/v2"
)

type QuasarHandler struct {
	db          *sqlx.DB
	finder      tzf.F
	versionsDir string
}

func NewQuasarHandler(db *sqlx.DB, versionsDir string) (*QuasarHandler, error) {
	finder, err := tzf.NewEmbeddedFinder()

	if err != nil {
		return nil, err
	}

	return &QuasarHandler{
		db:          db,
		finder:      finder,
		versionsDir: versionsDir,
	}, nil
}

func (h *QuasarHandler) RegisterRoutes(r chi.Router) {
	// Device endpoints
	r.Post("/register", h.RegisterDevice)

	// Update endpoints
	r.Get("/check_updates", h.CheckUpdates)
	r.HandleFunc("/download_update/{platform}/{version}", h.DownloadUpdate)

	// Common endpoints
	r.Get("/generate_204", h.Generate204)
	r.Get("/get_timezone", h.GetTimezone)
}
