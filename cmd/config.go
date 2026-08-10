package cmd

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// appConfig holds values loaded from the config file.
// Fields are applied as defaults — CLI flags and environment variables always take priority.
type appConfig struct {
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	QuayURL   string `yaml:"quay-url"`
}

// appCfg is initialized at package load time, before any init() functions run.
// This ensures config values are available when flags are registered.
var appCfg = loadConfig()

// loadConfig reads the config file from the user's config directory.
// Returns an empty config if the file doesn't exist or can't be parsed.
func loadConfig() appConfig {
	path := configFilePath()
	if path == "" {
		return appConfig{}
	}

	data, err := os.ReadFile(path) // #nosec G304 -- path is from os.UserConfigDir(), not user input
	if err != nil {
		return appConfig{}
	}

	var cfg appConfig
	if err := parseConfig(data, &cfg); err != nil {
		return appConfig{}
	}

	return cfg
}

// parseConfig unmarshals YAML config data into the given appConfig.
func parseConfig(data []byte, cfg *appConfig) error {
	return yaml.Unmarshal(data, cfg)
}

// configFilePath returns the path to the config file.
// Uses os.UserConfigDir() to determine the platform-appropriate config directory:
//   - Linux:   ~/.config/go-quay/config.yaml
//   - macOS:   ~/Library/Application Support/go-quay/config.yaml
//   - Windows: %AppData%/go-quay/config.yaml
func configFilePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, "go-quay", "config.yaml")
}
