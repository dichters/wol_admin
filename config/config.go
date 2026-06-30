// Package config reads and holds the global configuration from config.json.
// It is a singleton — call Load once at startup, then read via Cfg anywhere.
package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Cfg is the global configuration singleton, populated by Load.
var Cfg *Config

// Config represents all configurable parameters of the application.
type Config struct {
	ServerPort     string   `json:"server_port"`
	StdoutLogLevel string   `json:"stdout_log_level"`
	FileLogLevel   string   `json:"file_log_level"`
	EnableAntiShake bool    `json:"enable_anti_shake"`
	Redis          RedisCfg `json:"redis"`
	NasIP          string   `json:"nas_ip"`
	NasUser        string   `json:"nas_user"`
	NasMAC         string   `json:"nas_mac"`
}

// RedisCfg holds Redis connection parameters.
type RedisCfg struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

// Load reads config.json from dir and parses it into Cfg.
// It exits the process on any error (missing file, bad JSON, invalid log level).
func Load(dir string) {
	path := dir + "/config.json"

	data, err := os.ReadFile(path)
	if err != nil {
		// Use fmt before slog is ready
		fmt.Fprintf(os.Stderr, "fatal: cannot read config file %s: %v\n", path, err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: cannot parse config file %s: %v\n", path, err)
		os.Exit(1)
	}

	// Validate log levels
	if !isValidLevel(cfg.StdoutLogLevel) {
		fmt.Fprintf(os.Stderr, "fatal: invalid stdout_log_level %q, must be Debug/Info/Warn/Error\n", cfg.StdoutLogLevel)
		os.Exit(1)
	}
	if !isValidLevel(cfg.FileLogLevel) {
		fmt.Fprintf(os.Stderr, "fatal: invalid file_log_level %q, must be Debug/Info/Warn/Error\n", cfg.FileLogLevel)
		os.Exit(1)
	}

	// Default server port
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	Cfg = &cfg
}

// ParseLevel converts a config log-level string to slog.Level.
func ParseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func isValidLevel(s string) bool {
	switch strings.ToLower(s) {
	case "debug", "info", "warn", "error":
		return true
	default:
		return false
	}
}
