package cmd

import (
	"net/http"
	"testing"
)

const (
	testPrototypePath = "/organization/" + testNamespace + "/prototypes"
	testPrototypeUUID = "proto-1"
	testPrototypeJSON = `{"id":"proto-1","role":"write","delegate":{"name":"` + testPermUser + `","kind":"user"}}`
)

func setPrototypeVars(t *testing.T) {
	setCommandTestValue(t, &prototypeOrg, testNamespace)
	setCommandTestValue(t, &prototypeUUID, testPrototypeUUID)
	setCommandTestValue(t, &prototypeDelegateName, testPermUser)
	setCommandTestValue(t, &prototypeDelegateKind, "user")
	setCommandTestValue(t, &prototypeRole, "write")
	setCommandTestValue(t, &confirmProtoDelete, true)
}

func TestPrototypeCommands(t *testing.T) {
	uuidPath := testPrototypePath + "/" + testPrototypeUUID
	tests := []commandDomainCase{
		{
			name: "prototype list", run: prototypeListCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testPrototypePath},
			response: `{"prototypes":[` + testPrototypeJSON + `]}`, wantOutput: testPrototypeUUID, setup: setPrototypeVars,
		},
		apiErrorCase("prototype list", prototypeListCmd.RunE, testHTTPGet, testPrototypePath, http.StatusForbidden, "getting prototypes", setPrototypeVars),
		{
			name: "prototype info", run: prototypeInfoCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: uuidPath},
			response: testPrototypeJSON, wantOutput: testPrototypeUUID, setup: setPrototypeVars,
		},
		apiErrorCase("prototype info", prototypeInfoCmd.RunE, testHTTPGet, uuidPath, http.StatusNotFound, "getting prototype", setPrototypeVars),
		{
			name: "prototype create", run: prototypeCreateCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPost, path: testPrototypePath, body: map[string]interface{}{
				"delegate": map[string]interface{}{testFlagName: testPermUser, "kind": "user"}, testFlagRole: "write",
			}},
			response: testPrototypeJSON, wantOutput: testPrototypeUUID, setup: setPrototypeVars,
		},
		apiErrorCase("prototype create", prototypeCreateCmd.RunE, testHTTPPost, testPrototypePath, http.StatusBadRequest, "creating prototype", setPrototypeVars),
		{
			name: "prototype update", run: prototypeUpdateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: uuidPath, body: map[string]interface{}{testFlagRole: "write"}},
			response: testPrototypeJSON, wantOutput: testPrototypeUUID, setup: setPrototypeVars,
		},
		apiErrorCase("prototype update", prototypeUpdateCmd.RunE, testHTTPPut, uuidPath, http.StatusBadRequest, "updating prototype", setPrototypeVars),
		{
			name: "prototype delete", run: prototypeDeleteCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: uuidPath, status: http.StatusNoContent},
			setup:   setPrototypeVars,
		},
		apiErrorCase("prototype delete", prototypeDeleteCmd.RunE, testHTTPDelete, uuidPath, http.StatusNotFound, "deleting prototype", setPrototypeVars),
	}
	runCommandDomainCases(t, tests)
}

func TestPrototypeDeleteRequiresConfirm(t *testing.T) {
	setPrototypeVars(t)
	setCommandTestValue(t, &confirmProtoDelete, false)
	err := runCommandExpectingNoRequest(t, prototypeDeleteCmd.RunE)
	if err == nil || err.Error() != "--confirm is required to delete a prototype" {
		t.Fatalf("error = %v, want confirm requirement", err)
	}
}

func TestPrototypeRequiredFlags(t *testing.T) {
	proto := func(sub string) []string { return []string{cmdGet, "prototype", sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "prototype info", args: proto(subcmdInfo), want: []string{testRequiredFlagError, testFlagUUID}},
		{name: "prototype create", args: proto(subcmdCreate), want: []string{testRequiredFlagError, "delegate-name", "delegate-kind", testFlagRole}},
		{name: "prototype update", args: proto(subcmdUpdate), want: []string{testRequiredFlagError, testFlagUUID, testFlagRole}},
		{name: "prototype delete", args: proto(subcmdDelete), want: []string{testRequiredFlagError, testFlagUUID}},
	})
}
