package cmd

import (
	"testing"

	"github.com/sebrandon1/go-quay/lib"
)

func TestTokenFlagExistsOnGetCmd(t *testing.T) {
	flag := getCmd.PersistentFlags().Lookup("token")
	if flag == nil {
		t.Fatal("Expected --token flag on getCmd, not found")
	}

	if flag.Annotations != nil {
		if _, found := flag.Annotations["cobra_annotation_bash_completion_one_required_flag"]; found {
			t.Error("Token flag should not be marked as required — it defaults to $QUAY_TOKEN")
		}
	}
}

func TestTokenFlagOverridesEnvVar(t *testing.T) {
	flag := getCmd.PersistentFlags().Lookup("token")
	if flag == nil {
		t.Fatal("Expected --token flag on getCmd, not found")
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

func TestQuayURLFlagExists(t *testing.T) {
	flag := getCmd.PersistentFlags().Lookup("quay-url")
	if flag == nil {
		t.Fatal("Expected --quay-url flag on getCmd, not found")
	}
	if flag.DefValue != lib.DefaultQuayURL {
		// Default may come from QUAY_URL env var at init time
		if flag.Annotations != nil {
			if _, found := flag.Annotations["cobra_annotation_bash_completion_one_required_flag"]; found {
				t.Error("quay-url flag should not be marked as required")
			}
		}
	}
}
