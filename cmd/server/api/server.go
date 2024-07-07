package api

import (
	"database/sql"
	"encoding/json"
	"github.com/klaital/tv-controller/internal/config"
	"github.com/klaital/tv-controller/vlcclient"
	"io"
	"log/slog"
	"net/http"
)

var _ ServerInterface = (*Server)(nil)

type Server struct {
	Db        *sql.DB
	VlcClient *vlcclient.Client
}

func NewServer(db *sql.DB, vlc *vlcclient.Client) Server {
	return Server{Db: db, VlcClient: vlc}
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

	// Save updates to disk
	config.SaveConfig(&newCfg, s.Db)

	// echo the new config back to the client
	w.Write(b)
}

func (s Server) PausePlayback(w http.ResponseWriter, req *http.Request) {
	slog.Debug("Toggling pause")
	err := s.VlcClient.PlayPause()
	if err != nil {
		slog.Error("Failed to play/pause", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s Server) TrackAhead(w http.ResponseWriter, req *http.Request) {
	slog.Debug("Skipping ahead in playlist")
	err := s.VlcClient.TrackAhead()
	if err != nil {
		slog.Error("Failed to trackahead", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s Server) TrackBack(w http.ResponseWriter, req *http.Request) {
	slog.Debug("Backtracking in playlist")
	err := s.VlcClient.TrackBack()
	if err != nil {
		slog.Error("Failed to trackback", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s Server) ToggleShuffle(w http.ResponseWriter, req *http.Request) {
	cfg := config.LoadConfig(s.Db)
	// Load the current state of the VLC player to determine what commands to send
	vlcStatus, err := s.VlcClient.GetStatus()
	if err != nil {
		slog.Error("Failed to fetch current status from VLC", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// toggle the VLC player's setting, then save that value to the config store
	cfg.Shuffle = !vlcStatus.Random
	err = s.VlcClient.Random()
	if err != nil {
		slog.Error("Failed to toggle random setting on VLC", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// save the setting change
	config.SaveConfig(cfg, s.Db)

	w.WriteHeader(http.StatusNoContent)
}

func (s Server) ToggleLoop(w http.ResponseWriter, req *http.Request) {
	cfg := config.LoadConfig(s.Db)
	// Load the current state of the VLC player to determine what commands to send
	vlcStatus, err := s.VlcClient.GetStatus()
	if err != nil {
		slog.Error("Failed to fetch current status from VLC", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// toggle the VLC player's setting, then save that value to the config store
	cfg.Loop = !vlcStatus.Loop
	err = s.VlcClient.Loop()
	if err != nil {
		slog.Error("Failed to toggle loop setting on VLC", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// save the setting change
	config.SaveConfig(cfg, s.Db)

	w.WriteHeader(http.StatusNoContent)
}

func (s Server) SelectPlaylist(w http.ResponseWriter, req *http.Request) {
	cfg := config.LoadConfig(s.Db)

	var playlistRequest SelectPlaylistRequest
	b, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(b, &playlistRequest)
	if err != nil {
		slog.Error("Failed to unmarshal request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate palylist selection
	found := false
	for _, playlist := range cfg.PlaylistsAvailable {
		if playlist == *playlistRequest.Playlist {
			found = true
			break
		}
	}

	if !found {
		slog.Debug("Requested playlist not found", "playlist", playlistRequest.Playlist)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Tell VLC to start playing the specified playlist
	cfg.SelectedPlaylist = *playlistRequest.Playlist
	err = cfg.StopVlc()
	if err != nil {
		slog.Error("Failed to stop existing VLC instance", "error", err)
	}
	err = cfg.StartVlc()
	if err != nil {
		slog.Error("Failed to start new VLC instance", "error", err)
	}

	// save the setting change
	config.SaveConfig(cfg, s.Db)

	w.WriteHeader(http.StatusNoContent)
}
