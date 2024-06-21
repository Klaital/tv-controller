package config

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"os/user"
)

type Config struct {
	PlaylistsAvailable []string `json:"playlists_available"`
	SelectedPlaylist   string   `json:"selected_playlist"`
	Shuffle            bool     `json:"shuffle"`
	Loop               bool     `json:"loop"`
}

func GetConfigDir() string {
	u, err := user.Current()
	if err != nil {
		slog.Error("Failed to load current user", "error", err)
		os.Exit(1)
	}
	return fmt.Sprintf("%s/%s", u.HomeDir, ".tvcfg")
}

func LoadConfig(db *sql.DB) *Config {
	// lazy-init: create the empty config if none exists
	var cfg Config
	err := db.QueryRow("SELECT data FROM config LIMIT 1").Scan(&cfg)

	if err != nil {
		var sqlite3Err sqlite3.Error
		if errors.As(err, &sqlite3Err) {
			slog.Debug("Sqlite3 error", "error", sqlite3Err, "code", sqlite3Err.Code)
			_, err := db.Exec("CREATE TABLE config (data JSONB)")
			if err != nil {
				slog.Error("Failed to initialize config table", "error", err)
				os.Exit(1)
			}
		}
	}

	// load the set of configured playlists
	playlistFiles, err := os.ReadDir(GetConfigDir() + "/playlists")
	if err != nil {
		slog.Error("Failed to list playlist files", "error", err.Error())
		os.Exit(1)
	}
	for i := range playlistFiles {
		cfg.PlaylistsAvailable = append(cfg.PlaylistsAvailable, playlistFiles[i].Name())
	}

	return &cfg
}

func SaveConfig(cfg *Config, db *sql.DB) {
	_, err := db.Exec("UPDATE config SET data=?", cfg)
	if err != nil {
		slog.Error("Failed to write config update", "error", err)
		os.Exit(1)
	}
}
