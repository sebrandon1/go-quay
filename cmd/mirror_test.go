package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sebrandon1/go-quay/lib"
)

func resetMirrorCommandState(t *testing.T) {
	t.Helper()
	namespace = ""
	repository = ""
	for _, name := range []string{testNamespaceFlagName, cmdRepository} {
		if flag := mirrorCmd.PersistentFlags().Lookup(name); flag != nil {
			flag.Changed = false
			if err := flag.Value.Set(""); err != nil {
				t.Fatalf("reset mirror %s flag: %v", name, err)
			}
		}
	}
	t.Cleanup(func() {
		namespace = ""
		repository = ""
		rootCmd.SetArgs(nil)
	})
}

func TestMirrorInfoTableShowsScheduleConfiguration(t *testing.T) {
	resetRootFlags(t)
	resetMirrorCommandState(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/repository/"+testNamespace+"/"+testRepository+"/mirror" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		writeResponse(t, w, []byte(`{"is_enabled":true,"mirror_type":"PULL","external_reference":"`+testMirrorExternalRef+`","sync_interval":86400,"sync_start_date":"2025-01-01T00:00:00Z","robot_username":"myorg+mirrorbot","root_rule":{"rule":"`+testMirrorTagRule+`","rule_kind":"tag_glob_csv"}}`))
	}))
	defer server.Close()

	rootCmd.SetArgs([]string{
		cmdInfo, cmdMirror, "--namespace", testNamespace, "--repository", testRepository,
		testTokenFlag, testTokenValue, testQuayURLFlag, server.URL + "/api/v1", "--output", outputTable,
	})
	var runErr error
	stdout := captureStdout(t, func() { runErr = rootCmd.Execute() })
	if runErr != nil {
		t.Fatalf("mirror info: %v", runErr)
	}
	for _, expected := range []string{"ENABLED", "MIRROR TYPE", "EXTERNAL REF", "SYNC INTERVAL (SEC)", "SYNC START DATE", testTrueValue, "PULL", testMirrorExternalRef, "86400", "2025-01-01T00:00:00Z", "myorg+mirrorbot", testMirrorTagRule} {
		if !strings.Contains(stdout, expected) {
			t.Errorf("mirror table output %q is missing %q", stdout, expected)
		}
	}
}

func TestMirrorCreateAndUpdateCommands(t *testing.T) {
	const mirrorPath = "/repository/testns/testrepo/mirror"
	setupRepository := func(t *testing.T) {
		setCommandTestValue(t, &namespace, testNamespace)
		setCommandTestValue(t, &repository, testRepository)
	}
	tests := []commandDomainCase{
		{
			name:        subcmdCreate,
			executeArgs: []string{subcmdCreate, cmdMirror, testConfigNamespaceFlag, testNamespace, testDiffRepositoryFlag, testRepository, "--external-ref", testMirrorExternalRef, "--robot-username", testMirrorRobotUsername, "--sync-interval", "3600", "--sync-start-date", testStartTime, "--ext-username", "mirror-user", "--ext-password", "mirror-password", "--tag-rule", testMirrorTagRule, "--tag-rule-kind", "tag_glob_csv"},
			request: commandRequestExpectation{
				method: testHTTPPost, path: mirrorPath,
				body: map[string]interface{}{
					"external_reference": testMirrorExternalRef, "sync_interval": 3600,
					"sync_start_date": testStartTime, "robot_username": testMirrorRobotUsername,
					"external_registry_username": "mirror-user", "external_registry_password": "mirror-password",
					"root_rule": map[string]interface{}{"rule": "stable-*", "rule_kind": "tag_glob_csv"},
				},
			},
			response: `{"external_reference":"` + testMirrorExternalRef + `"}`, wantOutput: testMirrorExternalRef,
			setup: func(t *testing.T) {
				setupRepository(t)
				setCommandTestValue(t, &mirrorExternalRef, "")
				setCommandTestValue(t, &mirrorSyncInterval, 86400)
				setCommandTestValue(t, &mirrorSyncStartDate, "")
				setCommandTestValue(t, &mirrorRobotUsername, "")
				setCommandTestValue(t, &mirrorExtUser, "")
				setCommandTestValue(t, &mirrorExtPassword, "")
				setCommandTestValue(t, &mirrorTagRule, ".*")
				setCommandTestValue(t, &mirrorTagRuleKind, "tag_glob_csv")
			},
		},
		{
			name: "update configured fields",
			run:  mirrorUpdateCmd.RunE,
			request: commandRequestExpectation{
				method: testHTTPPut, path: mirrorPath,
				body: map[string]interface{}{
					"is_enabled": true, "external_reference": testMirrorExternalRef, "sync_interval": 7200,
					"sync_start_date": "2026-02-01T00:00:00Z", "robot_username": testMirrorRobotUsername,
				},
			},
			response: `{"is_enabled":true,"external_reference":"` + testMirrorExternalRef + `"}`, wantOutput: testMirrorExternalRef,
			setup: func(t *testing.T) {
				setupRepository(t)
				setCommandTestValue(t, &mirrorEnabled, true)
				setCommandTestValue(t, &mirrorExternalRef, testMirrorExternalRef)
				setCommandTestValue(t, &mirrorSyncInterval, 7200)
				setCommandTestValue(t, &mirrorSyncStartDate, "2026-02-01T00:00:00Z")
				setCommandTestValue(t, &mirrorRobotUsername, testMirrorRobotUsername)
			},
		},
		{
			name: "update enabled only",
			run:  mirrorUpdateCmd.RunE,
			request: commandRequestExpectation{
				method: testHTTPPut, path: mirrorPath,
				body: map[string]interface{}{"is_enabled": false},
			},
			response: `{}`,
			setup: func(t *testing.T) {
				setupRepository(t)
				setCommandTestValue(t, &mirrorEnabled, false)
			},
		},
		{
			name:     "create API error",
			run:      mirrorCreateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: mirrorPath, status: http.StatusInternalServerError},
			response: testAPIErrorResponse, wantError: "creating mirror config",
			setup: func(t *testing.T) {
				setupRepository(t)
				setCommandTestValue(t, &mirrorExternalRef, testMirrorExternalRef)
				setCommandTestValue(t, &mirrorRobotUsername, testMirrorRobotUsername)
			},
		},
	}
	runCommandDomainCases(t, tests)
}

func TestPrintMirrorSummaryRejectsNilConfig(t *testing.T) {
	if err := printMirrorSummary(nil); err == nil || !strings.Contains(err.Error(), "empty mirror config") {
		t.Fatalf("error = %v, want empty mirror config error", err)
	}
}

func TestPrintMirrorSummaryUsesDashesForUnsetFields(t *testing.T) {
	var runErr error
	stdout := captureStdout(t, func() { runErr = printMirrorSummary(&lib.MirrorConfig{}) })
	if runErr != nil {
		t.Fatalf("print empty mirror summary: %v", runErr)
	}
	if !strings.Contains(stdout, "false") || !strings.Contains(stdout, "-") {
		t.Fatalf("empty mirror summary = %q, want false and dash values", stdout)
	}
}
