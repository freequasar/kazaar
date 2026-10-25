package quasar

import (
	"github.com/freequasar/kazaar/internal/repositories"
	"github.com/freequasar/kazaar/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/ringsaturn/tzf/v2"
)

type QuasarHandler struct {
	deviceRepository          *repositories.DeviceRepository
	deviceService             *services.DeviceService
	versionManfiestRepository *repositories.VersionManfiestRepository
	versionManfiestService    *services.VersionManfiestService
	timezoneFinder            tzf.F
	versionsDir               string
}

func NewQuasarHandler(
	deviceRepository *repositories.DeviceRepository,
	deviceService *services.DeviceService,
	versionManfiestRepository *repositories.VersionManfiestRepository,
	versionManfiestService *services.VersionManfiestService,
	versionsDir string,
) (*QuasarHandler, error) {
	timezoneFinder, err := tzf.NewEmbeddedFinder()

	if err != nil {
		return nil, err
	}

	return &QuasarHandler{
		deviceRepository:          deviceRepository,
		deviceService:             deviceService,
		versionManfiestRepository: versionManfiestRepository,
		versionManfiestService:    versionManfiestService,
		timezoneFinder:            timezoneFinder,
		versionsDir:               versionsDir,
	}, nil
}

func (quasarHandler *QuasarHandler) RegisterRoutes(r chi.Router) {
	// Device endpoints
	r.Post("/register", quasarHandler.RegisterDevice)

	// Update endpoints
	r.Get("/check_updates", quasarHandler.CheckUpdates)
	r.HandleFunc("/download_update/{platform}/{version}", quasarHandler.DownloadUpdate)

	// Common endpoints
	r.Get("/generate_204", quasarHandler.Generate204)
	r.Get("/get_timezone", quasarHandler.GetTimezone)
}
