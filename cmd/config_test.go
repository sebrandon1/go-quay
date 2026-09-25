package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/pflag"
)

const (
	testConfigNamespaceFlag = "--namespace"
	testConfigInitCommand   = "init"
	testConfigRootCommand   = "config"
	testConfigNamespace     = "demo"
	testConfigFlagNamespace = "flag-namespace"
	testConfigToken         = "setup-secret"
	testConfigShowCommand   = "show"
	testWindowsOS           = "windows"
)

func resetConfigCommandState(t *testing.T) {
	t.Helper()
	resetRootFlags(t)
	oldNamespace, oldNonInteractive, oldShowToken := configNamespace, configInitNonInteractive, configShowToken
	t.Cleanup(func() {
		configNamespace, configInitNonInteractive, configShowToken = oldNamespace, oldNonInteractive, oldShowToken
		rootCmd.SetArgs([]string{})
		rootCmd.SetIn(nil)
	})
	configNamespace = ""
	configInitNonInteractive = false
	configShowToken = false
	if flag := configInitCmd.Flags().Lookup("namespace"); flag != nil {
		flag.Changed = false
		_ = flag.Value.Set("")
	}
	if flag := configInitCmd.Flags().Lookup("non-interactive"); flag != nil {
		flag.Changed = false
		_ = flag.Value.Set("false")
	}
	if flag := configShowCmd.Flags().Lookup("show-token"); flag != nil {
		flag.Changed = false
		_ = flag.Value.Set("false")
	}
}

func setConfigTestDir(t *testing.T, path string) {
	t.Helper()
	old := userConfigDir
	userConfigDir = func() (string, error) { return path, nil }
	t.Cleanup(func() { userConfigDir = old })
}

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
	configDir := filepath.Join(tmpDir, cliName)
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
	if filepath.Base(filepath.Dir(path)) != cliName {
		t.Errorf("Expected config dir named %q, got %q", cliName, filepath.Base(filepath.Dir(path)))
	}
}

func TestConfigInitNonInteractive(t *testing.T) {
	resetConfigCommandState(t)
	setConfigTestDir(t, t.TempDir())
	rootCmd.SetArgs([]string{testConfigRootCommand, testConfigInitCommand, "--non-interactive", "--token", testConfigToken, testConfigNamespaceFlag, testConfigNamespace})
	var runErr error
	output := captureStdout(t, func() { runErr = rootCmd.Execute() })
	if runErr != nil {
		t.Fatalf("config init: %v", runErr)
	}
	if strings.Contains(output, testConfigToken) {
		t.Fatal("config init output must not include the token")
	}

	path := configFilePath()
	cfg, err := readConfig(path)
	if err != nil {
		t.Fatalf("readConfig: %v", err)
	}
	if cfg.Token != testConfigToken || cfg.Namespace != testConfigNamespace {
		t.Fatalf("saved config = %+v, want token and namespace from flags", cfg)
	}
	if cfg.QuayURL != lib.DefaultQuayURL {
		t.Errorf("QuayURL = %q, want default %q", cfg.QuayURL, lib.DefaultQuayURL)
	}
	if runtime.GOOS != testWindowsOS {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat config file: %v", err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Errorf("config file permissions = %04o, want 0600", got)
		}
		info, err = os.Stat(filepath.Dir(path))
		if err != nil {
			t.Fatalf("stat config directory: %v", err)
		}
		if got := info.Mode().Perm(); got != 0o700 {
			t.Errorf("config directory permissions = %04o, want 0700", got)
		}
	}
}

func TestConfigInitDefaultsUseFlagEnvConfigPrecedence(t *testing.T) {
	resetConfigCommandState(t)
	setConfigTestDir(t, t.TempDir())
	if err := writeConfig(configFilePath(), appConfig{
		Token:     "file-token",
		Namespace: "file-namespace",
		QuayURL:   "https://file.example/api/v1",
	}); err != nil {
		t.Fatalf("writeConfig: %v", err)
	}
	t.Setenv("QUAY_TOKEN", "env-token")
	t.Setenv("QUAY_URL", "https://env.example/api/v1")

	cfg, err := configInitDefaults(configInitCmd)
	if err != nil {
		t.Fatalf("configInitDefaults: %v", err)
	}
	if cfg.Token != "env-token" || cfg.Namespace != "file-namespace" || cfg.QuayURL != "https://env.example/api/v1" {
		t.Errorf("environment defaults = %+v, want env token/URL and file namespace", cfg)
	}

	tokenFlag := rootCmd.PersistentFlags().Lookup("token")
	urlFlag := rootCmd.PersistentFlags().Lookup("quay-url")
	namespaceFlag := configInitCmd.Flags().Lookup("namespace")
	for flag, value := range map[*pflag.Flag]string{
		tokenFlag:     "flag-token",
		urlFlag:       "https://flag.example/api/v1",
		namespaceFlag: testConfigFlagNamespace,
	} {
		if err := flag.Value.Set(value); err != nil {
			t.Fatalf("set %s flag: %v", flag.Name, err)
		}
		flag.Changed = true
	}
	configNamespace = testConfigFlagNamespace
	cfg, err = configInitDefaults(configInitCmd)
	if err != nil {
		t.Fatalf("configInitDefaults with flags: %v", err)
	}
	if cfg.Token != "flag-token" || cfg.Namespace != testConfigFlagNamespace || cfg.QuayURL != "https://flag.example/api/v1" {
		t.Errorf("flag defaults = %+v, want flag values", cfg)
	}
}

func TestConfigPathDoesNotRequireToken(t *testing.T) {
	resetConfigCommandState(t)
	configRoot := t.TempDir()
	setConfigTestDir(t, configRoot)
	rootCmd.SetArgs([]string{testConfigRootCommand, "path"})
	var runErr error
	output := captureStdout(t, func() { runErr = rootCmd.Execute() })
	if runErr != nil {
		t.Fatalf("config path: %v", runErr)
	}
	if got, want := strings.TrimSpace(output), filepath.Join(configRoot, cliName, "config.yaml"); got != want {
		t.Errorf("config path = %q, want %q", got, want)
	}
}

func TestConfigShowRedactsTokenByDefault(t *testing.T) {
	resetConfigCommandState(t)
	setConfigTestDir(t, t.TempDir())
	if err := writeConfig(configFilePath(), appConfig{Token: "show-secret", Namespace: testConfigNamespace, QuayURL: "https://quay.example/api/v1"}); err != nil {
		t.Fatalf("writeConfig: %v", err)
	}

	rootCmd.SetArgs([]string{testConfigRootCommand, testConfigShowCommand})
	var runErr error
	output := captureStdout(t, func() { runErr = rootCmd.Execute() })
	if runErr != nil {
		t.Fatalf("config show: %v", runErr)
	}
	if strings.Contains(output, "show-secret") || !strings.Contains(output, "********") {
		t.Errorf("config show should redact the token, got: %s", output)
	}

	rootCmd.SetArgs([]string{testConfigRootCommand, testConfigShowCommand, "--show-token"})
	output = captureStdout(t, func() { runErr = rootCmd.Execute() })
	if runErr != nil {
		t.Fatalf("config show --show-token: %v", runErr)
	}
	if !strings.Contains(output, "show-secret") {
		t.Errorf("config show --show-token should include the token, got: %s", output)
	}
}

func TestConfigInitNonInteractiveRequiresFlags(t *testing.T) {
	resetConfigCommandState(t)
	setConfigTestDir(t, t.TempDir())
	rootCmd.SetArgs([]string{testConfigRootCommand, testConfigInitCommand, "--non-interactive", testConfigNamespaceFlag, testConfigNamespace})
	if err := rootCmd.Execute(); err == nil || !strings.Contains(err.Error(), "requires --token") {
		t.Fatalf("config init error = %v, want missing-token error", err)
	}
}

func TestConfigInitTreatsNonTerminalInputAsNonInteractive(t *testing.T) {
	resetConfigCommandState(t)
	setConfigTestDir(t, t.TempDir())
	rootCmd.SetIn(strings.NewReader(""))
	rootCmd.SetArgs([]string{testConfigRootCommand, testConfigInitCommand, "--token", testConfigToken, testConfigNamespaceFlag, testConfigNamespace})
	var runErr error
	output := captureStdout(t, func() { runErr = rootCmd.Execute() })
	if runErr != nil {
		t.Fatalf("config init with piped stdin: %v", runErr)
	}
	if strings.Contains(output, testConfigToken) {
		t.Fatal("config init output must not include the token")
	}
	if cfg := loadConfig(); cfg.Token != testConfigToken || cfg.Namespace != testConfigNamespace {
		t.Errorf("saved config = %+v, want token and namespace from flags", cfg)
	}
}

func TestPromptForConfigUsesDefaultsWithoutPrintingToken(t *testing.T) {
	cfg := appConfig{Token: "hidden-secret", Namespace: "old-org", QuayURL: "https://old.example/api/v1"}
	var prompt bytes.Buffer
	updated, err := promptForConfig(strings.NewReader("\nnew-org\n\n"), &prompt, cfg)
	if err != nil {
		t.Fatalf("promptForConfig: %v", err)
	}
	if updated.Token != cfg.Token || updated.Namespace != "new-org" || updated.QuayURL != cfg.QuayURL {
		t.Errorf("prompt config = %+v, unexpected values", updated)
	}
	if strings.Contains(prompt.String(), cfg.Token) {
		t.Errorf("token prompt must not display the configured token: %s", prompt.String())
	}
}

func TestIsTerminalRejectsPipeInput(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer reader.Close()
	defer writer.Close()
	if isTerminal(reader) {
		t.Fatal("pipe input must not be treated as a terminal")
	}
}

func TestConfigShowRejectsInvalidYAML(t *testing.T) {
	resetConfigCommandState(t)
	setConfigTestDir(t, t.TempDir())
	if err := os.MkdirAll(filepath.Dir(configFilePath()), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFilePath(), []byte("{invalid"), 0o600); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{testConfigRootCommand, testConfigShowCommand})
	if err := rootCmd.Execute(); err == nil || !strings.Contains(err.Error(), "parsing") {
		t.Fatalf("config show error = %v, want YAML parsing error", err)
	}
}

func TestWriteConfigRejectsInvalidDirectory(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(filePath, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeConfig(filepath.Join(filePath, "config.yaml"), appConfig{Token: "token"}); err == nil {
		t.Fatal("expected writeConfig to reject a file as the parent directory")
	}
}

func TestLoadConfigFromUserConfigDir(t *testing.T) {
	configContent := []byte(`token: "disk-token"
namespace: "disk-ns"
quay-url: "https://disk.example/api/v1"
`)

	var configDir string
	switch runtime.GOOS {
	case testWindowsOS:
		base := t.TempDir()
		t.Setenv("APPDATA", base)
		configDir = filepath.Join(base, cliName)
	case "darwin":
		home := t.TempDir()
		t.Setenv("HOME", home)
		configDir = filepath.Join(home, "Library", "Application Support", cliName)
	default:
		base := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", base)
		configDir = filepath.Join(base, cliName)
	}

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, configContent, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig()
	if cfg.Token != "disk-token" {
		t.Errorf("Token = %q, want disk-token", cfg.Token)
	}
	if cfg.Namespace != "disk-ns" {
		t.Errorf("Namespace = %q, want disk-ns", cfg.Namespace)
	}
	if cfg.QuayURL != "https://disk.example/api/v1" {
		t.Errorf("QuayURL = %q, want disk URL", cfg.QuayURL)
	}
}

func TestLoadConfigInvalidFileReturnsEmpty(t *testing.T) {
	var configDir string
	switch runtime.GOOS {
	case testWindowsOS:
		base := t.TempDir()
		t.Setenv("APPDATA", base)
		configDir = filepath.Join(base, cliName)
	case "darwin":
		home := t.TempDir()
		t.Setenv("HOME", home)
		configDir = filepath.Join(home, "Library", "Application Support", cliName)
	default:
		base := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", base)
		configDir = filepath.Join(base, cliName)
	}

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(`{invalid`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig()
	if cfg.Token != "" || cfg.Namespace != "" || cfg.QuayURL != "" {
		t.Errorf("expected empty config for invalid YAML, got %+v", cfg)
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
