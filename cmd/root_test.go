package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/pflag"
)

func resetRootFlags(t *testing.T) {
	t.Helper()
	origCfg := appCfg
	origFormat := outputFormat
	resetTokenAndURLFlags()
	t.Cleanup(func() {
		resetTokenAndURLFlags()
		outputFormat = origFormat
		appCfg = origCfg
		rootCmd.SetArgs([]string{})
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
		getCmd.SetArgs([]string{})
		getCmd.SetOut(nil)
		getCmd.SetErr(nil)
	})
	outputFormat = outputJSON
	appCfg = appConfig{}
}

func resetTokenAndURLFlags() {
	token = ""
	quayURL = lib.DefaultQuayURL
	for _, f := range []*pflag.Flag{
		rootCmd.Flag("token"),
		rootCmd.PersistentFlags().Lookup("token"),
	} {
		if f != nil {
			f.Changed = false
			_ = f.Value.Set("")
		}
	}
	for _, f := range []*pflag.Flag{
		rootCmd.Flag("quay-url"),
		rootCmd.PersistentFlags().Lookup("quay-url"),
	} {
		if f != nil {
			f.Changed = false
			_ = f.Value.Set(lib.DefaultQuayURL)
		}
	}
}

func TestTokenFlagExistsOnGetCmd(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("token")
	if flag == nil {
		t.Fatal("Expected --token flag on rootCmd, not found")
	}

	if flag.Annotations != nil {
		if _, found := flag.Annotations["cobra_annotation_bash_completion_one_required_flag"]; found {
			t.Error("Token flag should not be marked as required — it defaults to $QUAY_TOKEN")
		}
	}
}

func TestTokenFlagDefValueEmpty(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("token")
	if flag == nil {
		t.Fatal("Expected --token flag on rootCmd, not found")
	}
	if flag.DefValue != "" {
		t.Errorf("token flag DefValue must be empty so --help does not leak secrets, got %q", flag.DefValue)
	}
}

func TestTokenFlagOverridesEnvVar(t *testing.T) {
	resetRootFlags(t)

	flag := rootCmd.PersistentFlags().Lookup("token")
	if flag == nil {
		t.Fatal("Expected --token flag on rootCmd, not found")
	}

	err := flag.Value.Set("explicit-token")
	if err != nil {
		t.Fatalf("Failed to set token flag: %v", err)
	}

	if token != "explicit-token" {
		t.Errorf("Expected token to be 'explicit-token', got '%s'", token)
	}
}

func TestSetVersion(t *testing.T) {
	SetVersion("v1.2.3")
	if rootCmd.Version != "v1.2.3" {
		t.Errorf("Expected version 'v1.2.3', got '%s'", rootCmd.Version)
	}
}

func TestRootCommandUse(t *testing.T) {
	if rootCmd.Use != cliName {
		t.Errorf("Expected root command Use %q, got %q", cliName, rootCmd.Use)
	}
}

func TestQuayURLFlagExists(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("quay-url")
	if flag == nil {
		t.Fatal("Expected --quay-url flag on rootCmd, not found")
	}
	if flag.DefValue != lib.DefaultQuayURL {
		t.Errorf("quay-url flag DefValue = %q, want %q", flag.DefValue, lib.DefaultQuayURL)
	}
}

func TestHelpDoesNotLeakToken(t *testing.T) {
	resetRootFlags(t)

	const envSecret = "env-secret-token-must-not-appear"
	const cfgSecret = "cfg-secret-token-must-not-appear"
	t.Setenv("QUAY_TOKEN", envSecret)
	appCfg.Token = cfgSecret

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{cmdGet, "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("help should not error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, envSecret) {
		t.Error("help output leaked QUAY_TOKEN")
	}
	if strings.Contains(out, cfgSecret) {
		t.Error("help output leaked config file token")
	}
}

func TestResolveFlag(t *testing.T) {
	const (
		fromFlag   = "from-flag"
		fromEnv    = "from-env"
		fromConfig = "from-config"
	)

	tests := []struct {
		name      string
		changed   bool
		flagValue string
		fallbacks []string
		want      string
	}{
		{name: "flag wins when changed", changed: true, flagValue: fromFlag, fallbacks: []string{fromEnv, fromConfig}, want: fromFlag},
		{name: "empty flag wins when changed", changed: true, flagValue: "", fallbacks: []string{fromEnv}, want: ""},
		{name: "env when flag unchanged", changed: false, flagValue: "ignored-default", fallbacks: []string{fromEnv, fromConfig}, want: fromEnv},
		{name: "config when env empty", changed: false, flagValue: "", fallbacks: []string{"", fromConfig, "fallback"}, want: fromConfig},
		{name: "fallback when others empty", changed: false, flagValue: "", fallbacks: []string{"", "", lib.DefaultQuayURL}, want: lib.DefaultQuayURL},
		{name: "all empty", changed: false, flagValue: "", fallbacks: []string{"", ""}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveFlag(tt.changed, tt.flagValue, tt.fallbacks...)
			if got != tt.want {
				t.Errorf("resolveFlag() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPersistentPreRunResolvesTokenFromEnv(t *testing.T) {
	resetRootFlags(t)
	t.Setenv("QUAY_TOKEN", "from-env")

	if err := persistentPreRunE(rootCmd, nil); err != nil {
		t.Fatalf("PersistentPreRunE: %v", err)
	}
	if token != "from-env" {
		t.Errorf("token = %q, want %q", token, "from-env")
	}
}

func TestPersistentPreRunResolvesTokenFromConfig(t *testing.T) {
	resetRootFlags(t)
	t.Setenv("QUAY_TOKEN", "")
	appCfg.Token = "from-config"

	if err := persistentPreRunE(rootCmd, nil); err != nil {
		t.Fatalf("PersistentPreRunE: %v", err)
	}
	if token != "from-config" {
		t.Errorf("token = %q, want %q", token, "from-config")
	}
}

func TestPersistentPreRunFlagOverridesEnv(t *testing.T) {
	resetRootFlags(t)
	t.Setenv("QUAY_TOKEN", "from-env")
	token = "from-flag"
	if f := rootCmd.Flag("token"); f != nil {
		f.Changed = true
	}

	if err := persistentPreRunE(rootCmd, nil); err != nil {
		t.Fatalf("PersistentPreRunE: %v", err)
	}
	if token != "from-flag" {
		t.Errorf("token = %q, want %q", token, "from-flag")
	}
}

func TestPersistentPreRunResolvesQuayURLFromEnv(t *testing.T) {
	resetRootFlags(t)
	t.Setenv("QUAY_TOKEN", "from-env")
	t.Setenv("QUAY_URL", "https://custom.example/api/v1")

	if err := persistentPreRunE(rootCmd, nil); err != nil {
		t.Fatalf("PersistentPreRunE: %v", err)
	}
	if quayURL != "https://custom.example/api/v1" {
		t.Errorf("quayURL = %q, want custom env URL", quayURL)
	}
}

func TestPersistentPreRunRequiresToken(t *testing.T) {
	resetRootFlags(t)
	t.Setenv("QUAY_TOKEN", "")

	err := persistentPreRunE(rootCmd, nil)
	if err == nil {
		t.Fatal("expected error when token is missing")
	}
	if !strings.Contains(err.Error(), "authentication token required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerbPathRequiresToken(t *testing.T) {
	resetRootFlags(t)
	t.Setenv("QUAY_TOKEN", "")

	rootCmd.SetArgs([]string{cmdCreate, cmdRepository, "-n", "ns", "-r", "repo"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when token is missing on create repository")
	}
	if !strings.Contains(err.Error(), "authentication token required") {
		t.Errorf("unexpected error: %v", err)
	}
}
