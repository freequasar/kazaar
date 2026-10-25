package main

import (
	"net/http"
	"os"
	"path"

	"github.com/freequasar/kazaar/internal/database"
	"github.com/freequasar/kazaar/internal/quasar"
	"github.com/freequasar/kazaar/internal/repositories"
	"github.com/freequasar/kazaar/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "modernc.org/sqlite"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dbPath := path.Join(workingDirectory, "./kazaar.db")
	versionsDir := path.Join(workingDirectory, "./versions")

	db, err := database.InitDB(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	deviceRepository := repositories.NewDeviceRepository(db)
	versionManfiestRepository := repositories.NewVersionManfiestRepository(db)

	deviceService := services.NewDeviceService(deviceRepository, db)
	versionManfiestService := services.NewVersionManfiestService(versionManfiestRepository, db)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	quasarHandler, err := quasar.NewQuasarHandler(
		deviceRepository, deviceService,
		versionManfiestRepository, versionManfiestService,
		versionsDir,
	)
	if err != nil {
		panic(err)
	}
	quasarHandler.RegisterRoutes(router)

	http.ListenAndServe(":3000", router)
}
