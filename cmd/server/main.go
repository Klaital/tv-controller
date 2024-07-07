package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/klaital/tv-controller/cmd/server/api"
	"github.com/klaital/tv-controller/internal/config"
	"github.com/klaital/tv-controller/vlcclient"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	dbfile := config.GetConfigDir() + "/config.db"
	slog.Info("Connecting to DB", "db", dbfile)

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	slog.Info("Loading config")
	cfg := config.LoadConfig(db)
	fmt.Printf("%+v\n", cfg)

	// Construct API client for controlling VLC
	vlcClient := &vlcclient.Client{
		Addr:         "http://localhost:8090",
		HttpPassword: "bedroomtv123",
	}

	// Launch VLC server
	err = cfg.StartVlc()
	if err != nil {
		slog.Error("Failed to launch VLC", "error", err)
		os.Exit(1)
	}

	router := mux.NewRouter()
	spa := spaHandler{
		staticPath: "tv-controller-web/dist",
		indexPath:  "index.html",
	}
	router.PathPrefix("/web").Handler(spa)

	srv := api.NewServer(db, vlcClient)
	h := api.HandlerFromMux(srv, router)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080", "http://bedroom-tv"},
		AllowedMethods: []string{"GET", "PUT"},
	})
	s := &http.Server{
		Handler: c.Handler(h),
		Addr:    "0.0.0.0:8080",
	}

	slog.Info("Listening for HTTP requests", "Addr", s.Addr)
	s.ListenAndServe()
}
