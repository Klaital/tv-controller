package api

import (
	"database/sql"
	"encoding/json"
	"github.com/klaital/tv-controller/internal/config"
	"io"
	"log/slog"
	"net/http"
)

var _ ServerInterface = (*Server)(nil)

type Server struct {
	Db *sql.DB
}

func NewServer(db *sql.DB) Server {
	return Server{Db: db}
}

func pointerTo[T any](s T) *T {
	return &s
}

// GET /cfg
func (s Server) GetConfig(w http.ResponseWriter, req *http.Request) {
	cfg := config.LoadConfig(s.Db)
	b, err := json.Marshal(cfg)
	if err != nil {
		slog.Error("Failed to serialize config data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		slog.Error("Failed to write config data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// PUT /cfg
func (s Server) SetConfig(w http.ResponseWriter, req *http.Request) {
	// read the new config value
	var newCfg config.Config
	b, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(b, &newCfg)
	if err != nil {
		slog.Error("Failed to unmarshal request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oldCfg := config.LoadConfig(s.Db)
	// Detect changes and update VLC's settings to match
	if newCfg.SelectedPlaylist != oldCfg.SelectedPlaylist {
		slog.Debug("New playlist requested", "old", oldCfg.SelectedPlaylist, "new", newCfg.SelectedPlaylist)
		// TODO: command VLC to start the new playlist
	}
	if newCfg.Shuffle != oldCfg.Shuffle {
		slog.Debug("Shuffle toggled", "old", oldCfg.Shuffle, "new", newCfg.Shuffle)
		// TODO: command VLC to switch to/from shuffle mode
	}
	if newCfg.Loop != oldCfg.Loop {
		slog.Debug("Loop toggled", "old", oldCfg.Loop, "new", newCfg.Loop)
		// TODO: command VLC to switch to/from loop mode
	}
	// Save updates to disk
	config.SaveConfig(&newCfg, s.Db)

	// echo the new config back to the client
	w.Write(b)
}
