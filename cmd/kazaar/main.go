package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"path"
	"slices"

	"github.com/freequasar/kazaar/internal/database"
	"github.com/freequasar/kazaar/internal/emudhra"
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

	emudhraCA, err := emudhra.NewEmudhraCA()
	if err != nil {
		panic(err)
	}

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
	router.Mount("/", quasarHandler.Handler())

	domains := slices.Concat(
		quasarHandler.Domains(),
		// ...
	)

	tlsCert, err := emudhraCA.GenerateCert(domains)
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates:     []tls.Certificate{tlsCert},
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP384, tls.CurveP256},
	}

	server := &http.Server{
		Addr:      ":8443",
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}
