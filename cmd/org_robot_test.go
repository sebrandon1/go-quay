package cmd

import (
	"net/http"
	"testing"

	"github.com/spf13/cobra"
)

const (
	testOrgRobotsPath    = "/organization/" + testNamespace + "/robots"
	testOrgRobotName     = "buildbot"
	testOrgRobotPath     = testOrgRobotsPath + "/" + testOrgRobotName
	testOrgRobotFullName = "testns+buildbot"
	testOrgRobotJSON     = `{"name":"` + testOrgRobotFullName + `","description":"ci robot"}`
	testOrgRobotIssuer   = "https://issuer.example.com"
	testOrgRobotSubject  = "system:serviceaccount:ci:builder"
	testOrgRobotRole     = "write"
)

func setOrgRobotVars(t *testing.T) {
	setCommandTestValue(t, &orgName, testNamespace)
	setCommandTestValue(t, &robotShortname, testOrgRobotName)
	setCommandTestValue(t, &description, "ci robot")
	setCommandTestValue(t, &repository, testRepository)
	setCommandTestValue(t, &role, testOrgRobotRole)
	setCommandTestValue(t, &federationIssuer, testOrgRobotIssuer)
	setCommandTestValue(t, &federationSubject, testOrgRobotSubject)
	setCommandTestValue(t, &confirm, true)
}

func TestOrgRobotCommands(t *testing.T) {
	permPath := testOrgRobotPath + "/permissions"
	repoPermPath := permPath + "/" + testRepository
	federationPath := testOrgRobotPath + "/federation"
	tests := []commandDomainCase{
		{
			name: "org-robot list", run: orgRobotsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgRobotsPath},
			response: `{"robots":[` + testOrgRobotJSON + `]}`, wantOutput: testOrgRobotFullName, setup: setOrgRobotVars,
		},
		{
			name: "org-robot list API error", run: orgRobotsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgRobotsPath, status: http.StatusInternalServerError},
			response: testAPIErrorResponse, wantError: "getting organization robots", setup: setOrgRobotVars,
		},
		{
			name: "org-robot info", run: orgRobotCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgRobotPath},
			response: testOrgRobotJSON, wantOutput: testOrgRobotFullName, setup: setOrgRobotVars,
		},
		{
			name: "org-robot info API error", run: orgRobotCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgRobotPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "getting robot account", setup: setOrgRobotVars,
		},
		{
			name: "org-robot create", run: createRobotCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testOrgRobotPath, body: map[string]interface{}{testFlagDescription: "ci robot"}},
			response: testOrgRobotJSON, wantOutput: testOrgRobotFullName, setup: setOrgRobotVars,
		},
		{
			name: "org-robot create API error", run: createRobotCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testOrgRobotPath, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "creating robot account", setup: setOrgRobotVars,
		},
		{
			name: "org-robot delete", run: deleteRobotCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testOrgRobotPath, status: http.StatusNoContent},
			setup:   setOrgRobotVars,
		},
		{
			name: "org-robot delete API error", run: deleteRobotCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPDelete, path: testOrgRobotPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "deleting robot account", setup: setOrgRobotVars,
		},
		{
			name: "org-robot regenerate", run: regenerateRobotCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testOrgRobotPath + "/regenerate"},
			response: testOrgRobotJSON, wantOutput: testOrgRobotFullName, setup: setOrgRobotVars,
		},
		{
			name: "org-robot regenerate API error", run: regenerateRobotCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: testOrgRobotPath + "/regenerate", status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "regenerating robot token", setup: setOrgRobotVars,
		},
		{
			name: "org-robot permissions", run: orgRobotPermissionsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: permPath},
			response: `{"permissions":[{"repository":{"name":"testrepo"},"role":"write"}]}`, wantOutput: "testrepo", setup: setOrgRobotVars,
		},
		{
			name: "org-robot permissions API error", run: orgRobotPermissionsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: permPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "getting robot permissions", setup: setOrgRobotVars,
		},
		{
			name: "org-robot set permission", run: setRobotPermissionCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: repoPermPath, body: map[string]interface{}{testFlagRole: testOrgRobotRole}, status: http.StatusNoContent},
			setup:   setOrgRobotVars,
		},
		{
			name: "org-robot set permission API error", run: setRobotPermissionCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: repoPermPath, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "setting robot permission", setup: setOrgRobotVars,
		},
		{
			name: "org-robot remove permission", run: removeRobotPermissionCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: repoPermPath, status: http.StatusNoContent},
			setup:   setOrgRobotVars,
		},
		{
			name: "org-robot remove permission API error", run: removeRobotPermissionCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPDelete, path: repoPermPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "removing robot permission", setup: setOrgRobotVars,
		},
		{
			name: "org-robot federation get", run: orgRobotFederationGetCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: federationPath},
			response: `{"federation":[{"issuer":"` + testOrgRobotIssuer + `","subject":"` + testOrgRobotSubject + `"}]}`, wantOutput: testOrgRobotIssuer, setup: setOrgRobotVars,
		},
		{
			name: "org-robot federation get API error", run: orgRobotFederationGetCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: federationPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "getting robot federation", setup: setOrgRobotVars,
		},
		{
			name: "org-robot federation create", run: orgRobotFederationCreateCmd.RunE,
			request: commandRequestExpectation{
				method: testHTTPPost, path: federationPath,
				body: []map[string]interface{}{{"issuer": testOrgRobotIssuer, "subject": testOrgRobotSubject}},
			},
			setup: setOrgRobotVars,
		},
		{
			name: "org-robot federation create API error", run: orgRobotFederationCreateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: federationPath, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "creating robot federation", setup: setOrgRobotVars,
		},
		{
			name: "org-robot federation delete", run: orgRobotFederationDeleteCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: federationPath, status: http.StatusNoContent},
			setup:   setOrgRobotVars,
		},
		{
			name: "org-robot federation delete API error", run: orgRobotFederationDeleteCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPDelete, path: federationPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "deleting robot federation", setup: setOrgRobotVars,
		},
	}
	runCommandDomainCases(t, tests)
}

func TestOrgRobotConfirmGuards(t *testing.T) {
	tests := []struct {
		name string
		run  func(*cobra.Command, []string) error
		want string
	}{
		{"org-robot delete robot", deleteRobotCmd.RunE, "must pass --confirm to delete a robot account"},
		{"org-robot remove permission guard", removeRobotPermissionCmd.RunE, "must pass --confirm to remove a robot permission"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setOrgRobotVars(t)
			setCommandTestValue(t, &confirm, false)
			err := runCommandExpectingNoRequest(t, test.run)
			if err == nil || err.Error() != test.want {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestOrgRobotRequiredFlags(t *testing.T) {
	org := func(sub string) []string { return []string{cmdGet, cmdOrganization, sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "org-robot info", args: org("robot"), want: []string{testRequiredFlagError, testFlagRobot}},
		{name: "org-robot create", args: org("create-robot"), want: []string{testRequiredFlagError, testFlagRobot}},
		{name: "org-robot delete", args: org("delete-robot"), want: []string{testRequiredFlagError, testFlagRobot}},
		{name: "org-robot regenerate", args: org("regenerate-robot"), want: []string{testRequiredFlagError, testFlagRobot}},
		{name: "org-robot permissions", args: org("robot-permissions"), want: []string{testRequiredFlagError, testFlagRobot}},
		{name: "org-robot set permission", args: org("set-robot-permission"), want: []string{testRequiredFlagError, testFlagRobot, testFlagRepository, testFlagRole}},
		{name: "org-robot remove permission", args: org("remove-robot-permission"), want: []string{testRequiredFlagError, testFlagRobot, testFlagRepository}},
		{name: "org-robot federation get", args: org("robot-federation-get"), want: []string{testRequiredFlagError, testFlagRobot}},
		{name: "org-robot federation create", args: org("robot-federation-create"), want: []string{testRequiredFlagError, testFlagRobot, testFlagIssuer, testFlagSubject}},
		{name: "org-robot federation delete", args: org("robot-federation-delete"), want: []string{testRequiredFlagError, testFlagRobot}},
	})
}
