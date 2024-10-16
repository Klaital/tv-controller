package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/klaital/tv-controller/cmd/server/api"
	mqtt_server "github.com/klaital/tv-controller/cmd/server/mqtt"
	"github.com/klaital/tv-controller/internal/config"
	"github.com/klaital/tv-controller/vlcclient"
	"github.com/rs/cors"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("Loading config")
	cfg := config.LoadConfig()
	fmt.Printf("%+v\n", cfg)

	// Construct API client for controlling VLC
	vlcClient := &vlcclient.Client{
		Addr:         "http://localhost:8090",
		HttpPassword: "bedroomtv123",
	}

	// Launch VLC server
	err := cfg.StartVlc()
	if err != nil {
		slog.Error("Failed to launch VLC", "error", err)
		os.Exit(1)
	}

	// Set up the HTTP Server
	router := mux.NewRouter()
	spa := spaHandler{
		staticPath: "tv-controller-web/dist",
		indexPath:  "index.html",
	}
	router.PathPrefix("/web").Handler(spa)

	srv := api.NewServer(vlcClient)
	h := api.HandlerFromMux(&srv, router)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080", "http://bedroom-tv"},
		AllowedMethods: []string{"GET", "PUT"},
	})
	s := &http.Server{
		Handler: c.Handler(h),
		Addr:    "0.0.0.0:8080",
	}

	// Set up the MQTT Server
	mServer := mqtt_server.MqttServer{
		Host:                "klaital.com",
		Port:                1883,
		PublishConfigTopic:  "bedroom/tv/config",
		RequestConfigTopic:  "bedroom/tv/update",
		ChangePlaylistTopic: "bedroom/tv/playlist",

		Config: cfg,
		Vlc:    vlcClient,
	}

	go mServer.Start()

	slog.Info("Listening for HTTP requests", "Addr", s.Addr)
	if err = s.ListenAndServe(); err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
		os.Exit(1)
	}
}
