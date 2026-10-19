package cmd

import (
	"net/http"
	"testing"
)

const (
	testAutoPrunePath   = "/organization/" + testNamespace + "/autoprunepolicy"
	testAutoPruneUUID   = "policy-1"
	testAutoPruneMethod = "number_of_tags"
	testAutoPruneJSON   = `{"uuid":"policy-1","method":"number_of_tags","value":10}`
)

func setAutoPruneVars(t *testing.T) {
	setCommandTestValue(t, &orgName, testNamespace)
	setCommandTestValue(t, &policyUUID, testAutoPruneUUID)
	setCommandTestValue(t, &method, testAutoPruneMethod)
	setCommandTestValue(t, &pruneValue, 10)
	setCommandTestValue(t, &tagPattern, "v*")
	setCommandTestValue(t, &confirm, true)
}

func TestAutoPruneCommands(t *testing.T) {
	policyBody := map[string]interface{}{"method": testAutoPruneMethod, "value": 10, "tag_pattern": "v*"}
	tests := []commandDomainCase{
		{
			name: "auto-prune list", run: autoPruneCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testAutoPrunePath},
			response: `{"policies":[` + testAutoPruneJSON + `]}`, wantOutput: testAutoPruneUUID, setup: setAutoPruneVars,
		},
		{
			name: "auto-prune list API error", run: autoPruneCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testAutoPrunePath, status: http.StatusInternalServerError},
			response: testAPIErrorResponse, wantError: "getting auto-prune policies", setup: setAutoPruneVars,
		},
		{
			name: "auto-prune info", run: autoPrunePolicyCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testAutoPrunePath + "/" + testAutoPruneUUID},
			response: testAutoPruneJSON, wantOutput: testAutoPruneMethod, setup: setAutoPruneVars,
		},
		{
			name: "auto-prune info API error", run: autoPrunePolicyCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testAutoPrunePath + "/" + testAutoPruneUUID, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "getting auto-prune policy", setup: setAutoPruneVars,
		},
		{
			name: "auto-prune create", run: createAutoPruneCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testAutoPrunePath, body: policyBody},
			response: testAutoPruneJSON, wantOutput: testAutoPruneUUID, setup: setAutoPruneVars,
		},
		{
			name: "auto-prune create API error", run: createAutoPruneCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testAutoPrunePath, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "creating auto-prune policy", setup: setAutoPruneVars,
		},
		{
			name: "auto-prune update", run: updateAutoPruneCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testAutoPrunePath + "/" + testAutoPruneUUID, body: policyBody},
			response: testAutoPruneJSON, wantOutput: testAutoPruneUUID, setup: setAutoPruneVars,
		},
		{
			name: "auto-prune update API error", run: updateAutoPruneCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testAutoPrunePath + "/" + testAutoPruneUUID, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "updating auto-prune policy", setup: setAutoPruneVars,
		},
		{
			name: "auto-prune delete", run: deleteAutoPruneCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testAutoPrunePath + "/" + testAutoPruneUUID, status: http.StatusNoContent},
			setup:   setAutoPruneVars,
		},
		{
			name: "auto-prune delete API error", run: deleteAutoPruneCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPDelete, path: testAutoPrunePath + "/" + testAutoPruneUUID, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "deleting auto-prune policy", setup: setAutoPruneVars,
		},
	}
	runCommandDomainCases(t, tests)
}

func TestDeleteAutoPruneRequiresConfirm(t *testing.T) {
	setAutoPruneVars(t)
	setCommandTestValue(t, &confirm, false)
	err := runCommandExpectingNoRequest(t, deleteAutoPruneCmd.RunE)
	if err == nil || err.Error() != "must pass --confirm to delete an auto-prune policy" {
		t.Fatalf("error = %v, want confirm requirement", err)
	}
}

func TestAutoPruneRequiredFlags(t *testing.T) {
	org := []string{cmdGet, cmdOrganization}
	args := func(sub string) []string { return append(append([]string{}, org...), sub) }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "auto-prune info", args: args("auto-prune-policy"), want: []string{testRequiredFlagError, testFlagPolicyUUID}},
		{name: "auto-prune create", args: args("create-auto-prune"), want: []string{testRequiredFlagError, testFlagMethod, testFlagValue}},
		{name: "auto-prune update", args: args("update-auto-prune"), want: []string{testRequiredFlagError, testFlagPolicyUUID}},
		{name: "auto-prune delete", args: args("delete-auto-prune"), want: []string{testRequiredFlagError, testFlagPolicyUUID}},
	})
}
