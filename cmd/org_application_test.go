package cmd

import (
	"net/http"
	"testing"
)

const (
	testApplicationsPath = "/organization/" + testNamespace + "/applications"
	testApplicationID    = "client-1"
	testApplicationJSON  = `{"client_id":"client-1","name":"ci-app"}`
	testApplicationName  = "ci-app"
	testFlagClientID     = "client-id"
)

func setApplicationVars(t *testing.T) {
	setCommandTestValue(t, &orgName, testNamespace)
	setCommandTestValue(t, &clientID, testApplicationID)
	setCommandTestValue(t, &appName, testApplicationName)
	setCommandTestValue(t, &description, "builds things")
	setCommandTestValue(t, &applicationURI, "https://app.example.com")
	setCommandTestValue(t, &redirectURI, "https://app.example.com/cb")
	setCommandTestValue(t, &confirm, true)
}

func TestOrgApplicationCommands(t *testing.T) {
	idPath := testApplicationsPath + "/" + testApplicationID
	appBody := map[string]interface{}{
		testFlagName: testApplicationName, testFlagDescription: "builds things",
		"application_uri": "https://app.example.com", "redirect_uri": "https://app.example.com/cb",
	}
	tests := []commandDomainCase{
		{
			name: "application list", run: orgApplicationsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testApplicationsPath},
			response: `{"applications":[` + testApplicationJSON + `]}`, wantOutput: testApplicationID, setup: setApplicationVars,
		},
		apiErrorCase("application list", orgApplicationsCmd.RunE, testHTTPGet, testApplicationsPath, http.StatusForbidden, "getting organization applications", setApplicationVars),
		{
			name: "application info", run: applicationCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: idPath},
			response: testApplicationJSON, wantOutput: testApplicationID, setup: setApplicationVars,
		},
		apiErrorCase("application info", applicationCmd.RunE, testHTTPGet, idPath, http.StatusNotFound, "getting application", setApplicationVars),
		{
			name: "application create", run: createApplicationCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testApplicationsPath, body: appBody},
			response: testApplicationJSON, wantOutput: testApplicationID, setup: setApplicationVars,
		},
		apiErrorCase("application create", createApplicationCmd.RunE, testHTTPPost, testApplicationsPath, http.StatusBadRequest, "creating application", setApplicationVars),
		{
			name: "application update", run: updateApplicationCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: idPath, body: appBody},
			response: testApplicationJSON, wantOutput: testApplicationID, setup: setApplicationVars,
		},
		apiErrorCase("application update", updateApplicationCmd.RunE, testHTTPPut, idPath, http.StatusBadRequest, "updating application", setApplicationVars),
		{
			name: "application delete", run: deleteApplicationCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: idPath, status: http.StatusNoContent},
			setup:   setApplicationVars,
		},
		apiErrorCase("application delete", deleteApplicationCmd.RunE, testHTTPDelete, idPath, http.StatusNotFound, "deleting application", setApplicationVars),
		{
			name: "application reset secret", run: resetApplicationSecretCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: idPath + "/resetclientsecret"},
			response: testApplicationJSON, wantOutput: testApplicationID, setup: setApplicationVars,
		},
		apiErrorCase("application reset secret", resetApplicationSecretCmd.RunE, testHTTPPost, idPath+"/resetclientsecret", http.StatusNotFound, "resetting application secret", setApplicationVars),
	}
	runCommandDomainCases(t, tests)
}

func TestOrgApplicationDeleteRequiresConfirm(t *testing.T) {
	setApplicationVars(t)
	setCommandTestValue(t, &confirm, false)
	err := runCommandExpectingNoRequest(t, deleteApplicationCmd.RunE)
	if err == nil || err.Error() != "must pass --confirm to delete an application" {
		t.Fatalf("error = %v, want confirm requirement", err)
	}
}

func TestOrgApplicationRequiredFlags(t *testing.T) {
	org := func(sub string) []string { return []string{cmdGet, cmdOrganization, sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "application info", args: org("application"), want: []string{testRequiredFlagError, testFlagClientID}},
		{name: "application create", args: org("create-application"), want: []string{testRequiredFlagError, testFlagName}},
		{name: "application update", args: org("update-application"), want: []string{testRequiredFlagError, testFlagClientID, testFlagName}},
		{name: "application delete", args: org("delete-application"), want: []string{testRequiredFlagError, testFlagClientID}},
		{name: "application reset secret", args: org("reset-application-secret"), want: []string{testRequiredFlagError, testFlagClientID}},
	})
}
