package main

import (
	"net/http"
	"os"
	"path"

	"github.com/freequasar/kazaar/internal/database"
	"github.com/freequasar/kazaar/internal/quasar"
	"github.com/go-chi/chi/v5"
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

	router := chi.NewRouter()

	quasarHandler, err := quasar.NewQuasarHandler(db, versionsDir)
	if err != nil {
		panic(err)
	}
	quasarHandler.RegisterRoutes(router)

	http.ListenAndServe(":3000", router)
}
