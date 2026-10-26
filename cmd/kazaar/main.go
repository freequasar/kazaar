package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"path"

	"github.com/freequasar/kazaar/internal/database"
	"github.com/freequasar/kazaar/internal/emudhra"
	"github.com/freequasar/kazaar/internal/quasar"
	"github.com/freequasar/kazaar/internal/repositories"
	"github.com/freequasar/kazaar/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "modernc.org/sqlite"
)

var YANDEX_DOMAINS = []string{
	"quasar.yandex.net", "quasar.yandex.ru", "scbh.yandex.net",
	"clck.yandex.net", "mobileproxy.passport.yandex.net", "login.yandex.ru",
	"oauth.yandex.ru", "rpc.alice.yandex.net", "uniproxy.alice.yandex.net",
	"alice.yandex.net", "alice-ats.yandex.net", "alice-ats-test.yandex.net",
	"tel.yandex.net", "voiceservices.yandex.net", "ws.api.music.yandex.net",
	"startup.mobile.yandex.net", "startup.appmetrica.yandex.net", "report.appmetrica.yandex.net",
	"quasar.s3.yandex.net", "storage.mds.yandex.net", "strm.yandex.ru", "s3.yandex.net",
	"*.svc.cdn.yandex.net",
}

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
	quasarHandler.RegisterRoutes(router)

	tlsCert, err := emudhraCA.GenerateCert(YANDEX_DOMAINS)
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
