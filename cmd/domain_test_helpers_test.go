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
	"github.com/spf13/pflag"
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
	testStartTime             = "2026-01-01T00:00:00Z"
	testMirrorTagRule         = "stable-*"
	testBuildsCommand         = "builds"
	testMirrorRobotUsername   = "testns+mirrorbot"
	testRequiredFlagError     = `required flag(s)`
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
	name        string
	run         func(*cobra.Command, []string) error
	executeArgs []string
	request     commandRequestExpectation
	response    string
	wantOutput  string
	wantError   string
	setup       func(*testing.T)
}

func runCommandDomainCases(t *testing.T, tests []commandDomainCase) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.executeArgs != nil {
				resetRootFlags(t)
			}
			if test.setup != nil {
				test.setup(t)
			}
			output, err := runCommandAgainstHTTP(t, test.run, test.executeArgs, test.request, test.response)
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

func runCommandAgainstHTTP(t *testing.T, run func(*cobra.Command, []string) error, executeArgs []string, want commandRequestExpectation, response string) (string, error) {
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

	var runErr error
	output := captureStdout(t, func() {
		if executeArgs != nil {
			args := append([]string{testTokenFlag, testTokenValue, testQuayURLFlag, server.URL + "/api/v1"}, executeArgs...)
			rootCmd.SetArgs(args)
			rootCmd.SetErr(io.Discard)
			runErr = rootCmd.Execute()
			return
		}
		command := &cobra.Command{}
		command.SetContext(context.Background())
		runErr = run(command, nil)
	})
	if got := requestCount.Load(); got != 1 {
		t.Errorf("HTTP request count = %d, want 1", got)
	}
	return output, runErr
}

func TestCommandDomainExecuteValidation(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		want  []string
		setup func(*testing.T)
	}{
		{name: billingCmd.Use, args: []string{cmdGet, billingCmd.Use, orgBillingCmd.Use}, want: []string{testRequiredFlagError, "organization"}},
		{name: cmdBuild, args: []string{subcmdList, testBuildsCommand}, want: []string{testRequiredFlagError, testNamespaceFlagName, cmdRepository}},
		{name: logsCmd.Use, args: []string{cmdGet, logsCmd.Use, repoLogsCmd.Use}, want: []string{"namespace is required"}, setup: func(t *testing.T) {
			setCommandTestValue(t, &namespace, "")
			setCommandTestValue(t, &repository, "")
		}},
		{name: cmdManifest, args: []string{subcmdInfo, cmdManifest}, want: []string{testRequiredFlagError, cmdManifest, testNamespaceFlagName, cmdRepository}},
		{name: cmdMirror, args: []string{subcmdCreate, cmdMirror, testConfigNamespaceFlag, testNamespace, testDiffRepositoryFlag, testRepository}, want: []string{testRequiredFlagError, "external-ref", "robot-username"}},
		{name: cmdNotification, args: []string{subcmdList, "notifications"}, want: []string{testRequiredFlagError, testNamespaceFlagName, cmdRepository}},
		{name: cmdRobot, args: []string{subcmdInfo, cmdRobot}, want: []string{testRequiredFlagError}},
		{name: cmdTrigger, args: []string{subcmdInfo, cmdTrigger}, want: []string{testRequiredFlagError, "uuid"}},
		{name: cmdUser, args: []string{cmdGet, cmdUser, userLookupCmd.Use}, want: []string{testRequiredFlagError, "username"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resetRootFlags(t)
			resetCommandFlagChanges(rootCmd)
			if test.setup != nil {
				test.setup(t)
			}
			args := append([]string{testTokenFlag, testTokenValue, testQuayURLFlag, testUnreachableURL}, test.args...)
			rootCmd.SetArgs(args)
			rootCmd.SetErr(io.Discard)
			var runErr error
			captureStdout(t, func() { runErr = rootCmd.Execute() })
			if runErr == nil {
				t.Fatal("root command returned no validation error")
			}
			for _, want := range test.want {
				if !strings.Contains(runErr.Error(), want) {
					t.Errorf("root command error = %v, want it to contain %q", runErr, want)
				}
			}
		})
	}
}

func resetCommandFlagChanges(command *cobra.Command) {
	command.Flags().VisitAll(func(flag *pflag.Flag) { flag.Changed = false })
	command.PersistentFlags().VisitAll(func(flag *pflag.Flag) { flag.Changed = false })
	for _, child := range command.Commands() {
		resetCommandFlagChanges(child)
	}
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
