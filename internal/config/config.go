package config

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	PlaylistsAvailable []string `json:"playlists_available"`
	SelectedPlaylist   string   `json:"selected_playlist"`
	Shuffle            bool     `json:"shuffle"`
	Loop               bool     `json:"loop"`
	VlcPath            string   `json:"vlc_path"`

	vlc *exec.Cmd
}

func NewConfig() *Config {
	return &Config{
		PlaylistsAvailable: make([]string, 0),
		SelectedPlaylist:   "",
		Shuffle:            false,
		Loop:               false,
		VlcPath:            "C:\\Program Files\\VideoLAN\\VLC\\vlc.exe", // this is correct for a windows dev machine using WSL
	}
}

func (c *Config) ToString() string {
	var s string
	s += fmt.Sprintf("playlists_available=%s\n", strings.Join(c.PlaylistsAvailable, ","))
	s += fmt.Sprintf("selected_playlist=%s\n", c.SelectedPlaylist)
	s += fmt.Sprintf("shuffle=%t\n", c.Shuffle)
	s += fmt.Sprintf("loop=%t\n", c.Loop)
	return s
}

func GetConfigDir() string {
	u, err := user.Current()
	if err != nil {
		slog.Error("Failed to load current user", "error", err)
		os.Exit(1)
	}
	return filepath.Join(u.HomeDir, ".tvcfg")
}

func GetPlaylistPath(playlistName string) string {
	return filepath.Join(GetConfigDir(), "playlists", playlistName)
}

var singleton *Config

func LoadConfig() *Config {
	if singleton != nil {
		return singleton
	}
	// lazy-init: create the empty config if none exists
	var cfg *Config
	configBytes, err := os.ReadFile(filepath.Join(GetConfigDir(), "config.json"))
	if errors.Is(err, fs.ErrNotExist) || len(configBytes) == 0 {
		slog.Debug("Initializing default config")
		cfg = NewConfig()
		SaveConfig(cfg)
	} else if err != nil {
		slog.Error("Failed to load config.json", "error", err)
		os.Exit(1)
	} else {
		cfg = NewConfig()
		err = json.Unmarshal(configBytes, cfg)
		if err != nil {
			slog.Error("Failed to parse config.json", "error", err)
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

	singleton = cfg
	return cfg
}

func SaveConfig(cfg *Config) {
	configBytes, err := json.Marshal(cfg)
	if err != nil {
		slog.Error("Failed to marshal config.json", "error", err)
	}
	err = os.WriteFile(filepath.Join(GetConfigDir(), "config.json"), configBytes, fs.ModePerm)
	if err != nil {
		slog.Error("Failed to save config.json", "error", err)
	}
}

// Value implements the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (c *Config) Value() (driver.Value, error) {
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

func (c *Config) StartVlc() error {
	args := []string{
		fmt.Sprintf("--http-host=%s", "127.0.0.1"),
		fmt.Sprintf("--http-port=%s", "8090"),
		fmt.Sprintf("--http-password=%s", "bedroomtv123"),
		"--extraintf=http",
		"--fullscreen",
	}

	if c.Shuffle {
		args = append(args, "--random")
	}
	if c.Loop {
		args = append(args, "--loop")
	}

	if len(c.SelectedPlaylist) > 0 {
		args = append(args, GetPlaylistPath(c.SelectedPlaylist))
	}

	c.vlc = exec.Command(c.VlcPath, args...)
	slog.Debug("Launching VLC", "cmd", c.vlc.String())
	err := c.vlc.Start()
	if err != nil {
		slog.Error("Failed to start VLC", "error", err)
		return fmt.Errorf("failed to start VLC: %w", err)
	}

	// Success!
	return nil
}

func (c *Config) StopVlc() error {
	if c.vlc == nil {
		slog.Error("VLC already stopped")
		return errors.New("VLC already stopped")
	}
	if c.vlc.Process == nil {
		slog.Error("no handle to VLC process")
		return errors.New("no handle to VLC process")
	}
	return c.vlc.Process.Kill()
}
