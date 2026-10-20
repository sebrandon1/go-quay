package cmd

import (
	"net/http"
	"testing"
)

const (
	testQuotaPath  = "/organization/" + testNamespace + "/quota"
	testQuotaJSON  = `{"id":"7","limit_bytes":1048576}`
	testQuotaBytes = 1048576
)

func setQuotaVars(t *testing.T) {
	setCommandTestValue(t, &orgName, testNamespace)
	setCommandTestValue(t, &limitBytes, int64(testQuotaBytes))
	setCommandTestValue(t, &confirm, true)
}

func TestQuotaCommands(t *testing.T) {
	limitBody := map[string]interface{}{"limit_bytes": testQuotaBytes}
	tests := []commandDomainCase{
		{
			name: "quota info", run: orgQuotaCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testQuotaPath},
			response: testQuotaJSON, wantOutput: "1048576", setup: setQuotaVars,
		},
		apiErrorCase("quota info", orgQuotaCmd.RunE, testHTTPGet, testQuotaPath, http.StatusNotFound, "getting organization quota", setQuotaVars),
		{
			name: "quota create", run: createQuotaCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testQuotaPath, body: limitBody},
			response: testQuotaJSON, wantOutput: "1048576", setup: setQuotaVars,
		},
		apiErrorCase("quota create", createQuotaCmd.RunE, testHTTPPost, testQuotaPath, http.StatusBadRequest, "creating quota", setQuotaVars),
		{
			name: "quota update", run: updateQuotaCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testQuotaPath, body: limitBody},
			response: testQuotaJSON, wantOutput: "1048576", setup: setQuotaVars,
		},
		apiErrorCase("quota update", updateQuotaCmd.RunE, testHTTPPut, testQuotaPath, http.StatusBadRequest, "updating quota", setQuotaVars),
		{
			name: "quota delete", run: deleteQuotaCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testQuotaPath, status: http.StatusNoContent},
			setup:   setQuotaVars,
		},
		apiErrorCase("quota delete", deleteQuotaCmd.RunE, testHTTPDelete, testQuotaPath, http.StatusNotFound, "deleting quota", setQuotaVars),
	}
	runCommandDomainCases(t, tests)
}

func TestQuotaDeleteRequiresConfirm(t *testing.T) {
	setQuotaVars(t)
	setCommandTestValue(t, &confirm, false)
	if err := runCommandExpectingNoRequest(t, deleteQuotaCmd.RunE); err == nil {
		t.Fatal("expected confirm requirement error, got nil")
	}
}
