package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestDryRunBlocksAPIRequests(t *testing.T) {
	const (
		apiToken  = "dry-run-test-token"
		bodyToken = "dry-run-body-secret"
	)
	tests := []struct {
		name       string
		args       []string
		wantMethod string
		wantBody   string
	}{
		{
			name:       "create",
			args:       []string{cmdCreate, cmdRepository, "-n", testNamespace, "-r", testRepository, testVisibilityFlag, testPrivateVisibility, testDescriptionFlag, bodyToken},
			wantMethod: http.MethodPost,
			wantBody:   "request body omitted",
		},
		{
			name:       "delete",
			args:       []string{cmdDelete, cmdRepository, "-n", testNamespace, "-r", testRepository, "--confirm"},
			wantMethod: http.MethodDelete,
			wantBody:   "no request body",
		},
		{
			name:       "read-only",
			args:       []string{cmdInfo, cmdRepository, "-n", testNamespace, "-r", testRepository},
			wantMethod: http.MethodGet,
			wantBody:   "no request body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRootFlags(t)
			resetRepositoryFlags(t)

			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				requests.Add(1)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			var stderr strings.Builder
			rootCmd.SetErr(&stderr)
			args := []string{"--dry-run", testTokenFlag, apiToken, testQuayURLFlag, server.URL}
			rootCmd.SetArgs(append(args, tt.args...))
			err := rootCmd.Execute()
			if err == nil {
				t.Fatal("expected dry-run to report that the request was skipped")
			}

			message := err.Error() + stderr.String()
			if !strings.Contains(stderr.String(), tt.wantMethod) || !strings.Contains(stderr.String(), server.URL) {
				t.Errorf("CLI stderr %q should include the method and URL", stderr.String())
			}
			for _, want := range []string{"dry-run", tt.wantMethod, server.URL, tt.wantBody} {
				if !strings.Contains(message, want) {
					t.Errorf("dry-run message %q does not contain %q", message, want)
				}
			}
			for _, secret := range []string{apiToken, bodyToken} {
				if strings.Contains(message, secret) {
					t.Errorf("dry-run output leaked %q: %s", secret, message)
				}
			}
			if got := requests.Load(); got != 0 {
				t.Errorf("dry-run sent %d requests; want none", got)
			}
		})
	}
}

func TestDryRunDoesNotRequireToken(t *testing.T) {
	resetRootFlags(t)
	t.Setenv("QUAY_TOKEN", "")
	dryRun = true

	if err := persistentPreRunE(repoDeleteCmd, nil); err != nil {
		t.Fatalf("dry-run without a token: %v", err)
	}
}
