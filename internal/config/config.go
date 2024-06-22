package config

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
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
	VlcPath            string   `json:"vlc_path"`
}

func NewConfig() *Config {
	return &Config{
		PlaylistsAvailable: make([]string, 0),
		SelectedPlaylist:   "",
		Shuffle:            false,
		Loop:               false,
		VlcPath:            "/mnt/c/Program Files/VideoLAN/VLC/vlc.exe", // this is correct for a windows dev machine using WSL
	}
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
	var cfg *Config
	err := db.QueryRow("SELECT data FROM config LIMIT 1").Scan(&cfg)
	if err != nil {
		var sqlite3Err sqlite3.Error
		if errors.As(err, &sqlite3Err) {
			slog.Debug("Sqlite3 error", "error", sqlite3Err, "code", sqlite3Err.Code)
			_, err := db.Exec("CREATE TABLE config (data JSONB)")
			if err != nil {
				slog.Error("Failed to initialize config table", "error", err)
				os.Exit(1)
			} else {
				slog.Info("Created config table")
				_, err := db.Exec("INSERT INTO config (data) VALUES (?)", NewConfig())
				if err != nil {
					slog.Error("Failed to write default config settings", "error", err)
				} else {
					slog.Info("Initialized config db with default settings")
				}
				cfg = NewConfig()
			}
		} else {
			slog.Error("Error querying config data", "error", err)
			os.Exit(1)
		}
	}

	// load the set of configured playlists
	playlistFiles, err := os.ReadDir(GetConfigDir() + "/playlists")
	if err != nil {
		slog.Error("Failed to list playlist files", "error", err.Error())
		os.Exit(1)
	}
	cfg.PlaylistsAvailable = make([]string, 0, len(playlistFiles))
	for i := range playlistFiles {
		cfg.PlaylistsAvailable = append(cfg.PlaylistsAvailable, playlistFiles[i].Name())
	}

	return cfg
}

func SaveConfig(cfg *Config, db *sql.DB) {
	_, err := db.Exec("UPDATE config SET data=?", cfg)
	if err != nil {
		slog.Error("Failed to write config update", "error", err)
		os.Exit(1)
	} else {
		slog.Debug("Updated config settings")
	}
}

// Value implements the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (c Config) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (c *Config) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}
