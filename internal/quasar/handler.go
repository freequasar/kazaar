package quasar

import (
	"net/http"

	"github.com/freequasar/kazaar/internal/repositories"
	"github.com/freequasar/kazaar/internal/services"
	"github.com/freequasar/kazaar/internal/utilities"
	"github.com/go-chi/chi/v5"
	"github.com/ringsaturn/tzf/v2"
)

var domains = []string{
	"quasar.yandex.ru", "quasar.yandex.net",
	"quasar.yasmina.yango.com", "quasar.yasmina.funtechservices.com",
}

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

func (quasarHandler *QuasarHandler) Handler() http.Handler {
	quasarRouter := chi.NewRouter()

	// Device endpoints
	quasarRouter.Post("/register", quasarHandler.RegisterDevice)
	quasarRouter.Post(
		"/update_device_state",
		quasarHandler.UpdateDeviceState,
	)

	// Update endpoints
	quasarRouter.Get("/check_updates", quasarHandler.CheckUpdates)
	quasarRouter.HandleFunc(
		"/download_update/{platform}/{version}",
		quasarHandler.DownloadUpdate,
	)

	// Common endpoints
	quasarRouter.Get("/generate_204", quasarHandler.Generate204)
	quasarRouter.Get("/get_timezone", quasarHandler.GetTimezone)

	return utilities.DomainSubrouter(domains, quasarRouter)
}

func (quasarHandler *QuasarHandler) Domains() []string {
	return domains
}
