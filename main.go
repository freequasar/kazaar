package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var db *sqlx.DB

type Device struct {
	ID       int    `db:"id"`
	DeviceID string `db:"device_id"`
	Platform string `db:"platform"`
}

func main() {
	// connect to the database
	var err error
	db, err = sqlx.Connect("sqlite", "kazaar.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// execute the migrations
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {
		panic(err)
	}

	// setup a router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	r.Post("/register", register)

	// listen and serve
	http.ListenAndServe(":3000", r)
}

func envelope(w http.ResponseWriter, statusCode int, status string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(map[string]string{
		"status":  status,
		"message": message,
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !r.URL.Query().Has("device_id") || !r.URL.Query().Has("platform") {
		envelope(w, http.StatusForbidden, "error", "no device_id / platform")
		return
	}

	deviceId := r.URL.Query().Get("device_id")
	platform := r.URL.Query().Get("platform")

	_, err := db.Exec(
		"INSERT INTO device (device_id, platform) VALUES (?, ?)",
		deviceId,
		platform,
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				envelope(w, http.StatusForbidden, "error", "device must be unique")
				return
			}
		}

		envelope(w, http.StatusForbidden, "error", "internal server error")
		return
	}

	envelope(w, http.StatusOK, "error", "registered")
}
