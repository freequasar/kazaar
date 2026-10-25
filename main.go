package main

import (
	"embed"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	tzf "github.com/ringsaturn/tzf/v2"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var wd string
var db *sqlx.DB
var f tzf.F

func main() {
	var err error

	// get the working directory
	wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	// initialize the timezone finder
	f, err = tzf.NewEmbeddedFinder()
	if err != nil {
		panic(err)
	}

	// connect to the database
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
	r.Get("/generate_204", func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now().UTC().UnixMicro()
		w.Header().Set("X-Start-Time", strconv.FormatInt(ts, 10))
		w.WriteHeader(http.StatusNoContent)
	})
	r.Get("/check_updates", checkUpdates)
	r.Get("/download_update/{platform}/{version}", downloadUpdate)
	r.Head("/download_update/{platform}/{version}", downloadUpdate)
	r.Get("/get_timezone", getTimezone)
	r.Post("/register", register)

	// listen and serve
	http.ListenAndServe(":3000", r)
}

func getTimezone(w http.ResponseWriter, r *http.Request) {
	// check for latitude and longitude
	q := r.URL.Query()
	if !q.Has("lat") || !q.Has("lon") {
		reply(w, http.StatusNotFound, Envelope{
			Status:  "error",
			Message: "no lat / lon",
		})
		return
	}

	// parse the latitude and longitude
	var err error
	latitude, err := strconv.ParseFloat(q.Get("lat"), 64)
	if err != nil {
		reply(w, http.StatusBadRequest, Envelope{
			Status:  "error",
			Message: "invalid lat",
		})
		return
	}
	longitude, err := strconv.ParseFloat(q.Get("lon"), 64)
	if err != nil {
		reply(w, http.StatusBadRequest, Envelope{
			Status:  "error",
			Message: "invalid lon",
		})
		return
	}

	// reply with timezone and offset
	tz := f.GetTimezoneName(longitude, latitude)
	loc, _ := time.LoadLocation(tz)
	_, offsetInSeconds := time.Now().In(loc).Zone()
	reply(w, http.StatusOK, TimezoneResponse{
		Status:          "ok",
		Timezone:        tz,
		OffsetInSeconds: offsetInSeconds,
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	// get deviceId and platform
	q := r.URL.Query()
	if !q.Has("device_id") || !q.Has("platform") {
		reply(w, http.StatusForbidden, Envelope{
			Status:  "error",
			Message: "no device id / platform",
		})
		return
	}
	deviceId := q.Get("device_id")
	platform := q.Get("platform")

	// save the device in DB
	if _, err := createDevice(deviceId, platform); err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				reply(w, http.StatusForbidden, Envelope{
					Status:  "error",
					Message: "device must be unique",
				})
				return
			}
		}

		reply(w, http.StatusForbidden, Envelope{
			Status:  "error",
			Message: "internal server error",
		})
		return
	}

	// reply with success
	reply(w, http.StatusOK, RegisterResponse{
		Status:            "ok",
		HasCriticalUpdate: false,
	})
}
