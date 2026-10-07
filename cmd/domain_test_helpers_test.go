package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/spf13/cobra"
)

const (
	testHTTPGet               = "GET"
	testHTTPPost              = "POST"
	testHTTPPut               = "PUT"
	testHTTPDelete            = "DELETE"
	testAPIErrorCase          = "API error"
	testAPIErrorResponse      = `{"error":"broken"}`
	testLimitParam            = "limit"
	testLogsEndTimeParam      = "endtime"
	testLogsStartTimeParam    = "starttime"
	testLogsStartDate         = "2026-01-01"
	testLogsEndDate           = "2026-01-31"
	testLogsPushKind          = "push"
	testManifestLabelValue    = "platform"
	testMirrorExternalRef     = "docker.io/library/nginx"
	testRobotPath             = "/user/robots/buildbot"
	testRobotFederationPath   = testRobotPath + "/federation"
	testBuildID               = "build-1"
	testBuildResponse         = `{"id":"build-1"}`
	testTriggerResponse       = `{"id":"trigger-1"}`
	testNotificationResponse  = `{"uuid":"notice-1"}`
	testNotificationEventKey  = "event"
	testNotificationMethodKey = "method"
	testNotificationUUID      = "notice-1"
	testTriggerID             = "trigger-1"
	testLogsNextPage          = "cursor-1"
	testLogsNextPageResult    = "cursor-2"
	testLogsNextPageParam     = "next_page"
)

type commandRequestExpectation struct {
	method string
	path   string
	query  map[string]string
	body   interface{}
	status int
}

type commandDomainCase struct {
	name       string
	run        func(*cobra.Command, []string) error
	request    commandRequestExpectation
	response   string
	wantOutput string
	wantError  string
	setup      func(*testing.T)
}

func runCommandDomainCases(t *testing.T, tests []commandDomainCase) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setup != nil {
				test.setup(t)
			}
			output, err := runCommandAgainstHTTP(t, test.run, test.request, test.response)
			checkCommandTestResult(t, output, err, test.wantOutput, test.wantError)
		})
	}
}

func setCommandTestValue[T any](t *testing.T, target *T, value T) {
	t.Helper()
	previous := *target
	*target = value
	t.Cleanup(func() { *target = previous })
}

func runCommandAgainstHTTP(t *testing.T, run func(*cobra.Command, []string) error, want commandRequestExpectation, response string) (string, error) {
	t.Helper()
	setCommandTestValue(t, &token, testTokenValue)
	setCommandTestValue(t, &outputFormat, outputJSON)
	setCommandTestValue(t, &dryRun, false)
	setCommandTestValue(t, &verbose, false)
	setCommandTestValue(t, &maxRetries, 0)

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		path := strings.TrimPrefix(r.URL.Path, "/api/v1")
		if r.Method != want.method || path != want.path {
			t.Errorf("request = %s %s, want %s %s", r.Method, path, want.method, want.path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+testTokenValue {
			t.Errorf("Authorization = %q, want bearer token", got)
		}
		for key, expected := range want.query {
			if got := r.URL.Query().Get(key); got != expected {
				t.Errorf("query %s = %q, want %q", key, got, expected)
			}
		}
		if want.body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
			} else if !jsonMatches(body, want.body) {
				t.Errorf("request body = %s, want JSON equivalent to %v", bytes.TrimSpace(body), want.body)
			}
		}

		status := want.status
		if status == 0 {
			status = http.StatusOK
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if response != "" {
			_, _ = io.WriteString(w, response)
		}
	}))
	defer server.Close()
	setCommandTestValue(t, &quayURL, server.URL+"/api/v1")

	command := &cobra.Command{}
	command.SetContext(context.Background())
	var runErr error
	output := captureStdout(t, func() { runErr = run(command, nil) })
	if got := requestCount.Load(); got != 1 {
		t.Errorf("HTTP request count = %d, want 1", got)
	}
	return output, runErr
}

func jsonMatches(got []byte, want interface{}) bool {
	var actualValue interface{}
	if err := json.Unmarshal(got, &actualValue); err != nil {
		return false
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		return false
	}
	var expectedValue interface{}
	if err := json.Unmarshal(wantJSON, &expectedValue); err != nil {
		return false
	}
	return reflect.DeepEqual(actualValue, expectedValue)
}

func checkCommandTestResult(t *testing.T, output string, err error, wantOutput, wantErr string) {
	t.Helper()
	if wantErr == "" {
		if err != nil {
			t.Fatalf("command returned error: %v", err)
		}
	} else if err == nil || !strings.Contains(err.Error(), wantErr) {
		t.Fatalf("command error = %v, want it to contain %q", err, wantErr)
	}
	if wantOutput != "" && !strings.Contains(output, wantOutput) {
		t.Errorf("command output %q does not contain %q", output, wantOutput)
	}
}
