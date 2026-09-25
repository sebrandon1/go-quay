package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompletionCommandGeneratesScriptsWithoutToken(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		t.Run(shell, func(t *testing.T) {
			resetRootFlags(t)
			var stdout, stderr bytes.Buffer
			rootCmd.SetOut(&stdout)
			rootCmd.SetErr(&stderr)
			rootCmd.SetArgs([]string{cmdCompletion, shell})

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("completion %s: %v (stderr: %s)", shell, err, stderr.String())
			}
			if output := stdout.String(); output == "" || !strings.Contains(output, cliName) {
				t.Fatalf("completion %s output should contain %q, got %q", shell, cliName, output)
			}
		})
	}
}

func TestCompletionCommandRejectsUnsupportedShell(t *testing.T) {
	resetRootFlags(t)
	rootCmd.SetArgs([]string{cmdCompletion, "tcsh"})

	err := rootCmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "unsupported shell") {
		t.Fatalf("completion tcsh error = %v, want unsupported shell error", err)
	}
}

func TestCompletionHelpListsSupportedShells(t *testing.T) {
	resetRootFlags(t)
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetArgs([]string{cmdCompletion, flagHelp})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("completion --help: %v", err)
	}

	for _, shell := range []string{"bash", "zsh", "fish", "powershell"} {
		if !strings.Contains(stdout.String(), shell) {
			t.Errorf("completion help should list %q; got %q", shell, stdout.String())
		}
	}
}
