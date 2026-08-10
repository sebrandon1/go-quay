package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigMissingFile(t *testing.T) {
	cfg := loadConfig()
	if cfg.Token != "" {
		t.Errorf("Expected empty token, got %q", cfg.Token)
	}
	if cfg.Namespace != "" {
		t.Errorf("Expected empty namespace, got %q", cfg.Namespace)
	}
	if cfg.QuayURL != "" {
		t.Errorf("Expected empty quay-url, got %q", cfg.QuayURL)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temp config dir
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "go-quay")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configContent := []byte(`token: "test-token-123"
namespace: "my-org"
quay-url: "https://custom.quay.io/api/v1"
`)
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, configContent, 0o600); err != nil {
		t.Fatal(err)
	}

	// Override XDG_CONFIG_HOME (Linux) or use platform-specific approach
	// Since loadConfig uses os.UserConfigDir, we test the parsing logic directly
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var cfg appConfig
	if err := parseConfig(data, &cfg); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if cfg.Token != "test-token-123" {
		t.Errorf("Expected token 'test-token-123', got %q", cfg.Token)
	}
	if cfg.Namespace != "my-org" {
		t.Errorf("Expected namespace 'my-org', got %q", cfg.Namespace)
	}
	if cfg.QuayURL != "https://custom.quay.io/api/v1" {
		t.Errorf("Expected quay-url 'https://custom.quay.io/api/v1', got %q", cfg.QuayURL)
	}
}

func TestLoadConfigPartialValues(t *testing.T) {
	configContent := []byte(`token: "only-token"
`)
	var cfg appConfig
	if err := parseConfig(configContent, &cfg); err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if cfg.Token != "only-token" {
		t.Errorf("Expected token 'only-token', got %q", cfg.Token)
	}
	if cfg.Namespace != "" {
		t.Errorf("Expected empty namespace, got %q", cfg.Namespace)
	}
	if cfg.QuayURL != "" {
		t.Errorf("Expected empty quay-url, got %q", cfg.QuayURL)
	}
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	configContent := []byte(`{invalid yaml: [`)
	var cfg appConfig
	err := parseConfig(configContent, &cfg)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestConfigFilePath(t *testing.T) {
	path := configFilePath()
	if path == "" {
		t.Skip("os.UserConfigDir() not available on this platform")
	}
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("Expected config file named 'config.yaml', got %q", filepath.Base(path))
	}
	if filepath.Base(filepath.Dir(path)) != "go-quay" {
		t.Errorf("Expected config dir named 'go-quay', got %q", filepath.Base(filepath.Dir(path)))
	}
}

func TestFirstNonEmpty(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   string
	}{
		{"all empty", []string{"", "", ""}, ""},
		{"first non-empty", []string{"a", "b", "c"}, "a"},
		{"middle non-empty", []string{"", "b", "c"}, "b"},
		{"last non-empty", []string{"", "", "c"}, "c"},
		{"no values", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstNonEmpty(tt.values...)
			if got != tt.want {
				t.Errorf("firstNonEmpty(%v) = %q, want %q", tt.values, got, tt.want)
			}
		})
	}
}
