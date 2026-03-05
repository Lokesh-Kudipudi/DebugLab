package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents global CLI configuration
type Config struct {
	WorkspaceDir string `json:"workspace_dir"`
}

// configPath returns the path to the global config file
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home dir: %w", err)
	}
	dir := filepath.Join(home, ".dblab")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("could not create config dir: %w", err)
	}
	return filepath.Join(dir, "config.json"), nil
}

// SaveConfig writes the configuration to disk
func SaveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadConfig reads the configuration from disk
func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no workspace configured. Please run 'dblab init <directory>' first")
		}
		return nil, err
	}
	
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	if cfg.WorkspaceDir == "" {
		return nil, fmt.Errorf("no workspace configured. Please run 'dblab init <directory>' first")
	}
	return &cfg, nil
}
