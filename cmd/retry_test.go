package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestRetryFlagsControlCLIRequests(t *testing.T) {
	tests := []struct {
		name         string
		flags        []string
		wantAttempts int32
		wantError    bool
	}{
		{
			name:         "default disables retries",
			wantAttempts: 1,
			wantError:    true,
		},
		{
			name: "configured retries recover from rate limit",
			flags: []string{
				"--max-retries", "2",
				"--retry-backoff", "1ms",
				"--retry-max-backoff", "2ms",
			},
			wantAttempts: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRootFlags(t)
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if attempts.Add(1) == 1 {
					w.Header().Set("Retry-After", "0")
					w.WriteHeader(http.StatusTooManyRequests)
					_, _ = w.Write([]byte(`{"error":"rate limited"}`))
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"username":"retry-user"}`))
			}))
			defer server.Close()

			var stderr bytes.Buffer
			rootCmd.SetErr(&stderr)
			args := []string{cmdInfo, cmdUser, testTokenFlag, "retry-test-token", testQuayURLFlag, server.URL}
			args = append(args, tt.flags...)
			rootCmd.SetArgs(args)

			var runErr error
			stdout := captureStdout(t, func() { runErr = rootCmd.Execute() })
			if tt.wantError {
				if runErr == nil {
					t.Fatal("expected the rate-limited request to fail without retries")
				}
			} else {
				if runErr != nil {
					t.Fatalf("execute with retries: %v", runErr)
				}
				if !strings.Contains(stdout, `"username": "retry-user"`) {
					t.Errorf("stdout missing user JSON: %q", stdout)
				}
			}
			if got := attempts.Load(); got != tt.wantAttempts {
				t.Errorf("request attempts = %d, want %d", got, tt.wantAttempts)
			}
		})
	}
}

func TestRetryFlagsRejectNegativeValues(t *testing.T) {
	tests := []struct {
		name  string
		flag  string
		value string
	}{
		{name: "maximum retries", flag: "--max-retries", value: "-1"},
		{name: "initial backoff", flag: "--retry-backoff", value: "-1ms"},
		{name: "maximum backoff", flag: "--retry-max-backoff", value: "-1ms"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRootFlags(t)
			var stderr bytes.Buffer
			rootCmd.SetErr(&stderr)
			rootCmd.SetArgs([]string{cmdInfo, cmdUser, testTokenFlag, "retry-test-token", testQuayURLFlag, "http://127.0.0.1:1", tt.flag, tt.value})
			err := rootCmd.Execute()
			if err == nil || !strings.Contains(err.Error(), "zero or greater") {
				t.Fatalf("error = %v, want validation error for %s", err, tt.flag)
			}
		})
	}
}
