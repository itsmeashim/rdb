package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ConnectionString string `json:"connection_string"`
	MaxConnections   int    `json:"max_connections"`
	DefaultProgram   string `json:"default_program"`
	DefaultPlatform  string `json:"default_platform"`
}

func DefaultConfig() *Config {
	return &Config{
		ConnectionString: "",
		MaxConnections:   10,
		DefaultProgram:   "default",
		DefaultPlatform:  "default",
	}
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "rdb", "config.json"), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Save(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
