package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/klaital/tv-controller/cmd/server/api"
	"github.com/klaital/tv-controller/internal/config"
	"github.com/klaital/tv-controller/vlcclient"
	"github.com/krynr/cec"
	"github.com/krynr/cec/device/raspberrypi"
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

	// Initialize the CEC client to listen for power on/off signals from the TV
	d := raspberrypi.Init(cec.TV, cec.DeviceTypeTV)
	x, err := cec.New(d, cec.Config{OSDName: "RPI"})
	if err != nil {
		slog.Error("Failed to connect to CEC bus", "error", err)
		os.Exit(1)
	}
	// Logger handler. Returning false causes the next handler to always trigger
	x.AddHandleFunc(func(x *cec.Cec, msg cec.Message) bool {
		slog.Debug("CEC message", "msg", msg)
		return false
	})
	// Fallback default handler
	x.AddHandler(cec.DefaultHandler{})
	go x.Run()
	
	// Start the REST server
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

	slog.Info("Listening for HTTP requests", "Addr", s.Addr)
	s.ListenAndServe()
}
