package cmd

import (
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

const (
	testPermBasePath = "/repository/" + testNamespace + "/" + testRepository + "/permissions"
	testPermUser     = "alice"
	testPermTeam     = "devs"
	testPermRole     = "write"
)

func setPermissionVars(t *testing.T) {
	setCommandTestValue(t, &namespace, testNamespace)
	setCommandTestValue(t, &repository, testRepository)
	setCommandTestValue(t, &permissionUser, testPermUser)
	setCommandTestValue(t, &permissionTeamName, testPermTeam)
	setCommandTestValue(t, &permissionRole, testPermRole)
	setCommandTestValue(t, &confirmPermDeletion, true)
}

func TestPermissionCommands(t *testing.T) {
	userPath := testPermBasePath + "/user/" + testPermUser
	teamPath := testPermBasePath + "/team/" + testPermTeam
	rolePermJSON := `{"role":"write","name":"alice"}`
	roleBody := map[string]interface{}{testFlagRole: testPermRole}
	tests := []commandDomainCase{
		{
			name: "permissions list", run: permListCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testPermBasePath},
			response: `{"permissions":[` + rolePermJSON + `]}`, wantOutput: testPermUser, setup: setPermissionVars,
		},
		{
			name: "permissions list API error", run: permListCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testPermBasePath, status: http.StatusForbidden},
			response: testAPIErrorResponse, wantError: "getting repository permissions", setup: setPermissionVars,
		},
		{
			name: "permissions set", run: permSetCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: testPermBasePath + "/" + testPermUser, body: roleBody},
			setup:   setPermissionVars,
		},
		{
			name: "permissions set API error", run: permSetCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testPermBasePath + "/" + testPermUser, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "setting repository permission", setup: setPermissionVars,
		},
		{
			name: "permissions remove", run: permRemoveCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testPermBasePath + "/" + testPermUser, status: http.StatusNoContent},
			setup:   setPermissionVars,
		},
		{
			name: "permissions remove API error", run: permRemoveCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPDelete, path: testPermBasePath + "/" + testPermUser, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "removing repository permission", setup: setPermissionVars,
		},
		{
			name: "permissions user permissions", run: permUserPermissionsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testPermBasePath + "/user/"},
			response: `{"permissions":[` + rolePermJSON + `]}`, wantOutput: testPermUser, setup: setPermissionVars,
		},
		{
			name: "permissions user permissions API error", run: permUserPermissionsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testPermBasePath + "/user/", status: http.StatusForbidden},
			response: testAPIErrorResponse, wantError: "getting user permissions", setup: setPermissionVars,
		},
		{
			name: "permissions user permission", run: permUserPermissionCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: userPath},
			response: rolePermJSON, wantOutput: testPermUser, setup: setPermissionVars,
		},
		{
			name: "permissions user permission API error", run: permUserPermissionCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: userPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "getting user permission", setup: setPermissionVars,
		},
		{
			name: "permissions set user permission", run: permSetUserPermCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: userPath, body: roleBody},
			setup:   setPermissionVars,
		},
		{
			name: "permissions set user permission API error", run: permSetUserPermCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: userPath, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "setting user permission", setup: setPermissionVars,
		},
		{
			name: "permissions delete user permission", run: permDeleteUserPermCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: userPath, status: http.StatusNoContent},
			setup:   setPermissionVars,
		},
		{
			name: "permissions delete user permission API error", run: permDeleteUserPermCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPDelete, path: userPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "deleting user permission", setup: setPermissionVars,
		},
		{
			name: "permissions user transitive permission", run: permUserTransitiveCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: userPath + "/transitive"},
			response: rolePermJSON, wantOutput: testPermUser, setup: setPermissionVars,
		},
		{
			name: "permissions user transitive permission API error", run: permUserTransitiveCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: userPath + "/transitive", status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "getting user transitive permission", setup: setPermissionVars,
		},
		{
			name: "permissions team permissions", run: permTeamPermissionsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testPermBasePath + "/team/"},
			response: `{"permissions":[{"role":"write","name":"devs"}]}`, wantOutput: testPermTeam, setup: setPermissionVars,
		},
		{
			name: "permissions team permissions API error", run: permTeamPermissionsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testPermBasePath + "/team/", status: http.StatusForbidden},
			response: testAPIErrorResponse, wantError: "getting team permissions", setup: setPermissionVars,
		},
		{
			name: "permissions team permission", run: permTeamPermissionCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: teamPath},
			response: `{"role":"write","name":"devs"}`, wantOutput: testPermTeam, setup: setPermissionVars,
		},
		{
			name: "permissions team permission API error", run: permTeamPermissionCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: teamPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "getting team permission", setup: setPermissionVars,
		},
		{
			name: "permissions set team permission", run: permSetTeamPermCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: teamPath, body: roleBody},
			setup:   setPermissionVars,
		},
		{
			name: "permissions set team permission API error", run: permSetTeamPermCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: teamPath, status: http.StatusBadRequest},
			response: testAPIErrorResponse, wantError: "setting team permission", setup: setPermissionVars,
		},
		{
			name: "permissions delete team permission", run: permDeleteTeamPermCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: teamPath, status: http.StatusNoContent},
			setup:   setPermissionVars,
		},
		{
			name: "permissions delete team permission API error", run: permDeleteTeamPermCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPDelete, path: teamPath, status: http.StatusNotFound},
			response: testAPIErrorResponse, wantError: "deleting team permission", setup: setPermissionVars,
		},
	}
	runCommandDomainCases(t, tests)
}

func TestPermissionDeleteRequiresConfirm(t *testing.T) {
	tests := []struct {
		name string
		run  func(*cobra.Command, []string) error
		want string
	}{
		{"permissions user", permDeleteUserPermCmd.RunE, "are you sure you want to remove permission for user 'alice'"},
		{"permissions team", permDeleteTeamPermCmd.RunE, "are you sure you want to remove permission for team 'devs'"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setPermissionVars(t)
			setCommandTestValue(t, &confirmPermDeletion, false)
			err := runCommandExpectingNoRequest(t, test.run)
			if err == nil || !strings.Contains(err.Error(), test.want) || !strings.Contains(err.Error(), "--confirm") {
				t.Fatalf("error = %v, want confirm prompt containing %q", err, test.want)
			}
		})
	}
}

func TestPermissionRequiredFlags(t *testing.T) {
	perm := func(sub string) []string { return []string{cmdGet, subcmdPermissions, sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "permissions set", args: perm("set"), want: []string{testRequiredFlagError, testFlagUser, testFlagRole}},
		{name: "permissions remove", args: perm("remove"), want: []string{testRequiredFlagError, testFlagUser}},
		{name: "permissions user permission", args: perm("user-permission"), want: []string{testRequiredFlagError, testFlagUser}},
		{name: "permissions set user permission", args: perm("set-user-permission"), want: []string{testRequiredFlagError, testFlagUser, testFlagRole}},
		{name: "permissions delete user permission", args: perm("delete-user-permission"), want: []string{testRequiredFlagError, testFlagUser}},
		{name: "permissions user transitive", args: perm("user-transitive-permission"), want: []string{testRequiredFlagError, testFlagUser}},
		{name: "permissions team permission", args: perm("team-permission"), want: []string{testRequiredFlagError, testFlagTeam}},
		{name: "permissions set team permission", args: perm("set-team-permission"), want: []string{testRequiredFlagError, testFlagTeam, testFlagRole}},
		{name: "permissions delete team permission", args: perm("delete-team-permission"), want: []string{testRequiredFlagError, testFlagTeam}},
	})
}
