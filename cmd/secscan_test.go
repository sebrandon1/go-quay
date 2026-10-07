package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

func resetSecscanCommandState(t *testing.T) {
	t.Helper()
	resetRootFlags(t)
	oldSilenceUsage := rootCmd.SilenceUsage
	rootCmd.SilenceUsage = true
	t.Cleanup(func() { rootCmd.SilenceUsage = oldSilenceUsage })
	reset := func() {
		namespace = appCfg.Namespace
		repository = ""
		secScanManifestRef = ""
		includeVulnerabilities = true
		secScanWatch = false
		secScanInterval = defaultSecScanInterval
		secScanWatchTimeout = defaultSecScanWatchTimeout
		for name, value := range map[string]string{
			"vulnerabilities": strconv.FormatBool(includeVulnerabilities),
			"watch":           "false",
			"interval":        defaultSecScanInterval.String(),
			"watch-timeout":   defaultSecScanWatchTimeout.String(),
		} {
			if flag := secscanInfoCmd.Flags().Lookup(name); flag != nil {
				flag.Changed = false
				_ = flag.Value.Set(value)
			}
		}
		for name, value := range map[string]string{
			testNamespaceFlagName: appCfg.Namespace,
			cmdRepository:         "",
			"manifest":            "",
		} {
			if flag := secscanCmd.PersistentFlags().Lookup(name); flag != nil {
				flag.Changed = false
				_ = flag.Value.Set(value)
			}
		}
		rootCmd.SetArgs(nil)
		rootCmd.SetContext(context.Background())
		secscanInfoCmd.SetContext(context.Background())
	}
	reset()
	t.Cleanup(reset)
}

func runSecscanInfo(t *testing.T, quayURL string, flags ...string) (string, string, error) {
	return runSecscanInfoWithContext(t, context.Background(), quayURL, flags...)
}

func runSecscanInfoWithContext(t *testing.T, ctx context.Context, quayURL string, flags ...string) (string, string, error) {
	t.Helper()
	rootCmd.SetContext(ctx)
	secscanInfoCmd.SetContext(ctx)
	var stderr strings.Builder
	rootCmd.SetErr(&stderr)
	args := []string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, quayURL,
		"secscan", subcmdInfo, "-n", testNamespace, "-r", testRepository,
		testManifestFlag, "sha256:test-manifest",
	}
	rootCmd.SetArgs(append(args, flags...))
	var executeErr error
	stdout := captureStdout(t, func() { executeErr = rootCmd.ExecuteContext(ctx) })
	return stdout, stderr.String(), executeErr
}

func serveSecurityStatus(w http.ResponseWriter, status string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"status":%q}`, status)
}

type cancelOnWaitingWriter struct {
	cancel context.CancelFunc
	output strings.Builder
}

func (w *cancelOnWaitingWriter) Write(p []byte) (int, error) {
	n, err := w.output.Write(p)
	if strings.Contains(string(p), "waiting") {
		w.cancel()
	}
	return n, err
}

func TestSecscanInfoWithoutWatchMakesOneRequest(t *testing.T) {
	resetSecscanCommandState(t)
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		serveSecurityStatus(w, "scanned")
	}))
	defer server.Close()

	stdout, _, err := runSecscanInfo(t, server.URL)
	if err != nil {
		t.Fatalf("secscan info: %v", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
	if !strings.Contains(stdout, `"status": "scanned"`) {
		t.Fatalf("output does not contain the scan response: %s", stdout)
	}
}

func TestSecscanWatchPollsUntilScanned(t *testing.T) {
	resetSecscanCommandState(t)
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if requests.Add(1) == 1 {
			serveSecurityStatus(w, "scanning")
			return
		}
		serveSecurityStatus(w, "scanned")
	}))
	defer server.Close()

	stdout, stderr, err := runSecscanInfo(t, server.URL, testWatchFlag, "--interval", "1ms", "--watch-timeout", "1s")
	if err != nil {
		t.Fatalf("secscan watch: %v", err)
	}
	if got := requests.Load(); got != 2 {
		t.Fatalf("requests = %d, want 2", got)
	}
	if !strings.Contains(stdout, `"status": "scanned"`) {
		t.Fatalf("output does not contain final scan response: %s", stdout)
	}
	if !strings.Contains(stderr, `status "scanning"`) {
		t.Fatalf("stderr does not report the pending status: %s", stderr)
	}
}

func TestSecscanWatchFailureReturnsNonzeroAfterPrintingResult(t *testing.T) {
	resetSecscanCommandState(t)
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		serveSecurityStatus(w, "failed")
	}))
	defer server.Close()

	stdout, stderr, err := runSecscanInfo(t, server.URL, testWatchFlag, "--interval", "1ms", "--watch-timeout", "1s")
	if err == nil || !strings.Contains(err.Error(), "security scan failed") {
		t.Fatalf("error = %v, want security scan failure", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1", got)
	}
	if !strings.Contains(stdout, `"status": "failed"`) || !strings.Contains(stderr, "security scan failed") {
		t.Fatalf("missing final failure result or error; stdout=%q stderr=%q", stdout, stderr)
	}
}

func TestSecscanWatchUnsupportedIsTerminal(t *testing.T) {
	resetSecscanCommandState(t)
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		serveSecurityStatus(w, "unsupported")
	}))
	defer server.Close()

	stdout, _, err := runSecscanInfo(t, server.URL, testWatchFlag, "--interval", "1ms", "--watch-timeout", "1s")
	if err != nil {
		t.Fatalf("unsupported scan should be terminal: %v", err)
	}
	if got := requests.Load(); got != 1 || !strings.Contains(stdout, `"status": "unsupported"`) {
		t.Fatalf("requests=%d output=%q; want one unsupported result", got, stdout)
	}
}

func TestSecscanWatchTimeout(t *testing.T) {
	resetSecscanCommandState(t)
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		serveSecurityStatus(w, "scanning")
	}))
	defer server.Close()

	_, _, err := runSecscanInfo(t, server.URL, testWatchFlag, "--interval", "1ms", "--watch-timeout", "10ms")
	if !errors.Is(err, context.DeadlineExceeded) || !strings.Contains(err.Error(), "watch timed out") {
		t.Fatalf("error = %v, want bounded watch timeout", err)
	}
	if got := requests.Load(); got == 0 {
		t.Fatal("watch did not request the scan status before timing out")
	}
}

func TestSecscanWatchReturnsRequestError(t *testing.T) {
	resetSecscanCommandState(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"unavailable"}`, http.StatusServiceUnavailable)
	}))
	defer server.Close()

	_, _, err := runSecscanInfo(t, server.URL, testWatchFlag, "--interval", "1ms", "--watch-timeout", "1s")
	if err == nil || !strings.Contains(err.Error(), "getting security scan") {
		t.Fatalf("error = %v, want security scan request error", err)
	}
}

func TestSecscanWatchHandlesContextErrorsDuringRequest(t *testing.T) {
	tests := []struct {
		name          string
		cancelParent  bool
		watchTimeout  string
		wantErr       error
		wantErrorText string
	}{
		{
			name:          "parent context canceled",
			cancelParent:  true,
			watchTimeout:  "1s",
			wantErr:       context.Canceled,
			wantErrorText: "context canceled",
		},
		{
			name:          "watch deadline expires",
			watchTimeout:  "10ms",
			wantErr:       context.DeadlineExceeded,
			wantErrorText: "watch timed out",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSecscanCommandState(t)
			requestStarted := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				close(requestStarted)
				<-r.Context().Done()
			}))
			defer server.Close()

			ctx := context.Background()
			if tt.cancelParent {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				defer cancel()
				go func() {
					<-requestStarted
					cancel()
				}()
			}

			_, _, err := runSecscanInfoWithContext(t, ctx, server.URL,
				testWatchFlag, "--interval", "1ms", "--watch-timeout", tt.watchTimeout)
			if !errors.Is(err, tt.wantErr) || !strings.Contains(err.Error(), tt.wantErrorText) {
				t.Fatalf("error = %v, want %v and %q", err, tt.wantErr, tt.wantErrorText)
			}
		})
	}
}

func TestSecscanWatchHonorsCanceledContext(t *testing.T) {
	resetSecscanCommandState(t)
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		serveSecurityStatus(w, "scanning")
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rootCmd.SetContext(ctx)
	secscanInfoCmd.SetContext(ctx)
	stderr := &cancelOnWaitingWriter{cancel: cancel}
	rootCmd.SetErr(stderr)
	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		"secscan", subcmdInfo, "-n", testNamespace, "-r", testRepository,
		testManifestFlag, "sha256:test-manifest", testWatchFlag, "--interval", "1h", "--watch-timeout", "1m",
	})
	var err error
	captureStdout(t, func() { err = rootCmd.ExecuteContext(ctx) })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("canceled watch sent %d requests, want exactly one before cancellation", got)
	}
	if !strings.Contains(stderr.output.String(), `status "scanning"`) {
		t.Fatalf("watch did not report pending status before cancellation: %s", stderr.output.String())
	}
}

func TestSecscanWatchRejectsNonpositiveTiming(t *testing.T) {
	tests := []struct {
		name  string
		flags []string
		want  string
	}{
		{name: "interval", flags: []string{testWatchFlag, "--interval", "0s"}, want: "--interval must be greater than 0"},
		{name: "watch timeout", flags: []string{testWatchFlag, "--watch-timeout", "0s"}, want: "--watch-timeout must be greater than 0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetSecscanCommandState(t)
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				requests.Add(1)
				serveSecurityStatus(w, "scanning")
			}))
			defer server.Close()

			_, _, err := runSecscanInfo(t, server.URL, tt.flags...)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
			if got := requests.Load(); got != 0 {
				t.Fatalf("invalid watch timing sent %d requests", got)
			}
		})
	}
}
