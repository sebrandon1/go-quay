package cmd

import (
	"net/http"
	"testing"
)

const (
	testRepoTokenName = "ci-token"
	testRepoTokenPath = "/repository/" + testNamespace + "/" + testRepository + "/tokens"
	testRepoTokenCode = "code-1"
	testRepoTokenJSON = `{"code":"code-1","friendlyName":"ci-token","role":"read"}`
)

func setRepoTokenVars(t *testing.T) {
	setCommandTestValue(t, &repoTokenNamespace, testNamespace)
	setCommandTestValue(t, &repoTokenRepository, testRepository)
	setCommandTestValue(t, &repoTokenCode, testRepoTokenCode)
	setCommandTestValue(t, &repoTokenName, testRepoTokenName)
	setCommandTestValue(t, &repoTokenRole, "read")
	setCommandTestValue(t, &confirmTokenDelete, true)
}

func TestRepoTokenCommands(t *testing.T) {
	codePath := testRepoTokenPath + "/" + testRepoTokenCode
	tests := []commandDomainCase{
		{
			name: "repotoken list", run: repotokenListCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testRepoTokenPath},
			response: `{"tokens":[` + testRepoTokenJSON + `]}`, wantOutput: testRepoTokenName, setup: setRepoTokenVars,
		},
		{
			name: "repotoken list API error", run: repotokenListCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testRepoTokenPath, status: http.StatusInternalServerError},
			response: testAPIErrorResponse, wantError: "getting tokens", setup: setRepoTokenVars,
		},
		{
			name: "repotoken info", run: repotokenInfoCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: codePath},
			response: testRepoTokenJSON, wantOutput: testRepoTokenName, setup: setRepoTokenVars,
		},
		{
			name: "repotoken info API error", run: repotokenInfoCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: codePath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "getting token", setup: setRepoTokenVars,
		},
		{
			name: "repotoken create", run: repotokenCreateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testRepoTokenPath, body: map[string]interface{}{"friendlyName": testRepoTokenName}},
			response: testRepoTokenJSON, wantOutput: testRepoTokenCode, setup: setRepoTokenVars,
		},
		{
			name: "repotoken create API error", run: repotokenCreateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testRepoTokenPath, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "creating token", setup: setRepoTokenVars,
		},
		{
			name: "repotoken update", run: repotokenUpdateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: codePath, body: map[string]interface{}{testFlagRole: "read"}},
			response: testRepoTokenJSON, wantOutput: testRepoTokenCode, setup: setRepoTokenVars,
		},
		{
			name: "repotoken update API error", run: repotokenUpdateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: codePath, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "updating token", setup: setRepoTokenVars,
		},
		{
			name: "repotoken delete", run: repotokenDeleteCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: codePath, status: http.StatusNoContent},
			setup:   setRepoTokenVars,
		},
		{
			name: "repotoken delete API error", run: repotokenDeleteCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPDelete, path: codePath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "deleting token", setup: setRepoTokenVars,
		},
	}
	runCommandDomainCases(t, tests)
}

func TestRepoTokenDeleteRequiresConfirm(t *testing.T) {
	setRepoTokenVars(t)
	setCommandTestValue(t, &confirmTokenDelete, false)
	err := runCommandExpectingNoRequest(t, repotokenDeleteCmd.RunE)
	if err == nil || err.Error() != "--confirm is required to delete a token" {
		t.Fatalf("error = %v, want confirm requirement", err)
	}
}

func TestRepoTokenRequiredFlags(t *testing.T) {
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "repotoken info", args: []string{cmdGet, cmdRepoToken, subcmdInfo}, want: []string{testRequiredFlagError, testFlagCode}},
		{name: "repotoken create", args: []string{cmdGet, cmdRepoToken, subcmdCreate}, want: []string{testRequiredFlagError, testFlagName}},
		{name: "repotoken update", args: []string{cmdGet, cmdRepoToken, subcmdUpdate}, want: []string{testRequiredFlagError, testFlagCode, testFlagRole}},
		{name: "repotoken delete", args: []string{cmdGet, cmdRepoToken, subcmdDelete}, want: []string{testRequiredFlagError, testFlagCode}},
	})
}
