package main

import (
	"database/sql"
	"fmt"
	"github.com/klaital/tv-controller/cmd/server/api"
	"github.com/klaital/tv-controller/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	dbfile := config.GetConfigDir() + "/config.db"
	slog.Info("Connecting to DB", "db", dbfile)
	os.Remove(dbfile)

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	slog.Info("Loading config")
	cfg := config.LoadConfig(db)
	fmt.Printf("%+v\n", cfg)

	srv := api.NewServer(db)
	r := http.NewServeMux()
	h := api.HandlerFromMux(srv, r)
	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	slog.Info("Listening for HTTP requests", "Addr", s.Addr)
	s.ListenAndServe()
}
