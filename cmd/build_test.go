package cmd

import (
	"strings"
	"testing"
)

func TestBuildCommands(t *testing.T) {
	setupBuild := func(t *testing.T) {
		setCommandTestValue(t, &buildNamespace, testNamespace)
		setCommandTestValue(t, &buildRepository, testRepository)
	}
	setupUUID := func(t *testing.T) { setupBuild(t); setCommandTestValue(t, &buildUUID, testBuildID) }
	tests := []commandDomainCase{
		{name: subcmdList, run: buildListCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/repository/testns/testrepo/build/", query: map[string]string{testLimitParam: "4"}}, response: `{"builds":[` + testBuildResponse + `]}`, wantOutput: testBuildID, setup: func(t *testing.T) { setupBuild(t); setCommandTestValue(t, &buildLimit, 4) }},
		{name: subcmdInfo, run: buildInfoCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/repository/testns/testrepo/build/" + testBuildID}, response: testBuildResponse, wantOutput: testBuildID, setup: setupUUID},
		{name: "logs", run: buildLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/repository/testns/testrepo/build/" + testBuildID + "/logs"}, response: `{"logs":[{"message":"compile"}]}`, wantOutput: "compile", setup: setupUUID},
		{name: "request", run: buildRequestCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: "/repository/testns/testrepo/build/", body: map[string]interface{}{"archive_url": "https://example.test/source.tar.gz", "dockerfile_path": "Dockerfile", "subdirectory": "app", "tags": []string{"v1", "stable"}}}, response: `{"id":"build-2"}`, wantOutput: "build-2", setup: func(t *testing.T) {
			setupBuild(t)
			setCommandTestValue(t, &buildArchiveURL, "https://example.test/source.tar.gz")
			setCommandTestValue(t, &buildDockerfile, "Dockerfile")
			setCommandTestValue(t, &buildSubdirectory, "app")
			setCommandTestValue(t, &buildTags, []string{"v1", "stable"})
		}},
		{name: "cancel", run: buildCancelCmd.RunE, request: commandRequestExpectation{method: testHTTPDelete, path: "/repository/testns/testrepo/build/" + testBuildID}, setup: func(t *testing.T) { setupUUID(t); setCommandTestValue(t, &confirmBuildCancel, true) }},
		{name: "status", run: buildStatusCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/repository/testns/testrepo/build/" + testBuildID + "/status"}, response: `{"status":"complete"}`, wantOutput: "complete", setup: setupUUID},
		{name: testAPIErrorCase, run: buildStatusCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/repository/testns/testrepo/build/" + testBuildID + "/status", status: 500}, response: testAPIErrorResponse, wantError: "getting build status", setup: setupUUID},
	}
	runCommandDomainCases(t, tests)
}

func TestBuildCancelRequiresConfirmation(t *testing.T) {
	setCommandTestValue(t, &buildNamespace, testNamespace)
	setCommandTestValue(t, &buildRepository, testRepository)
	setCommandTestValue(t, &buildUUID, testBuildID)
	setCommandTestValue(t, &confirmBuildCancel, false)
	setCommandTestValue(t, &token, testTokenValue)
	setCommandTestValue(t, &quayURL, testUnreachableURL)
	if err := buildCancelCmd.RunE(nil, nil); err == nil || !strings.Contains(err.Error(), "Use --confirm to proceed") {
		t.Fatalf("build cancel error = %v, want confirmation prompt", err)
	}
}
