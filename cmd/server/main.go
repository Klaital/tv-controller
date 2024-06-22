package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/klaital/tv-controller/cmd/server/api"
	"github.com/klaital/tv-controller/internal/config"
	"github.com/klaital/tv-controller/vlcclient"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
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

	vlcStatus, err := vlcClient.GetStatus()
	//slog.Debug("GetStatus", "error", err, "status", vlcStatus)
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		slog.Info("Timeout calling VLC. Re-launching service", "error", err)
		vlcWait := sync.WaitGroup{}
		vlcWait.Add(1)
		// Launch VLC as a background process that will live as long as this server process is running.
		go func() {
			cmd := exec.Command("/mnt/c/Program Files/VideoLAN/VLC/vlc.exe", "--http-host=127.0.0.1", "--http-port=8090", "--extraintf=http", "--http-password=bedroomtv123")
			err = cmd.Run()
			if err != nil {
				slog.Error("Failed to launch vlc", "error", err)
				os.Exit(1)
			}
			vlcWait.Done()
		}()

		// Wait for VLC to start up before continuing on
		vlcWait.Wait()
		vlcStatus, err = vlcClient.GetStatus()
		if err != nil {
			slog.Error("Still cannot connect to VLC. Check settings", "error", err)
			os.Exit(1)
		}
	} else if err != nil {
		slog.Error("Failed to fetch initial status from VLC")
		os.Exit(1)
	}
	// synchronize settings. Set loop and shuffle settings if needed to make VLC match the user's saved preferences
	if vlcStatus.Random != cfg.Shuffle {
		if err = vlcClient.Random(); err != nil {
			slog.Error("Failed to set VLC random setting to match saved config", "error", err)
			os.Exit(1)
		}
	}
	if vlcStatus.Loop != cfg.Loop {
		if err = vlcClient.Loop(); err != nil {
			slog.Error("Failed to set VLC loop setting to match saved config", "error", err)
			os.Exit(1)
		}
	}

	srv := api.NewServer(db, vlcClient)
	r := http.NewServeMux()
	h := api.HandlerFromMux(srv, r)
	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	slog.Info("Listening for HTTP requests", "Addr", s.Addr)
	s.ListenAndServe()
}
