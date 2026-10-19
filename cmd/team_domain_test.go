package cmd

import (
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

const (
	testTeamBasePath  = "/organization/" + testNamespace + "/team/" + testOrgTeamName
	testTeamMember    = "dave"
	testTeamMemberURL = testTeamBasePath + "/members/" + testTeamMember
	testTeamPermURL   = testTeamBasePath + "/permissions/" + testRepository
)

func setTeamCmdVars(t *testing.T) {
	setCommandTestValue(t, &teamCmdOrgname, testNamespace)
	setCommandTestValue(t, &teamCmdName, testOrgTeamName)
	setCommandTestValue(t, &teamCmdDescription, "the devs")
	setCommandTestValue(t, &teamCmdRole, "member")
	setCommandTestValue(t, &teamCmdMemberName, testTeamMember)
	setCommandTestValue(t, &teamCmdRepository, testRepository)
	setCommandTestValue(t, &teamCmdPermissionRole, testPermRole)
	setCommandTestValue(t, &confirmTeamDelete, true)
	setCommandTestValue(t, &confirmTeamMemberDel, true)
	setCommandTestValue(t, &confirmTeamPermDel, true)
}

func TestTeamDomainCommands(t *testing.T) {
	teamBody := map[string]interface{}{testFlagDescription: "the devs", testFlagRole: testFlagMember}
	tests := []commandDomainCase{
		{
			name: "team cmd list", run: teamListCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgTeamsPath},
			response: `{"teams":[` + testOrgTeamJSON + `]}`, wantOutput: testOrgTeamName, setup: setTeamCmdVars,
		},
		apiErrorCase("team cmd list", teamListCmd.RunE, testHTTPGet, testOrgTeamsPath, http.StatusForbidden, "getting teams", setTeamCmdVars),
		{
			name: "team cmd info", run: teamCmdInfoCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testTeamBasePath},
			response: testOrgTeamJSON, wantOutput: testOrgTeamName, setup: setTeamCmdVars,
		},
		apiErrorCase("team cmd info", teamCmdInfoCmd.RunE, testHTTPGet, testTeamBasePath, http.StatusNotFound, "getting team", setTeamCmdVars),
		{
			name: "team cmd create", run: teamCreateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testTeamBasePath, body: mergeBody(teamBody, map[string]interface{}{testFlagName: testOrgTeamName})},
			response: testOrgTeamJSON, wantOutput: testOrgTeamName, setup: setTeamCmdVars,
		},
		apiErrorCase("team cmd create", teamCreateCmd.RunE, testHTTPPut, testTeamBasePath, http.StatusBadRequest, "creating team", setTeamCmdVars),
		{
			name: "team cmd update", run: teamUpdateCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testTeamBasePath, body: teamBody},
			response: testOrgTeamJSON, wantOutput: testOrgTeamName, setup: setTeamCmdVars,
		},
		apiErrorCase("team cmd update", teamUpdateCmd.RunE, testHTTPPut, testTeamBasePath, http.StatusBadRequest, "updating team", setTeamCmdVars),
		{
			name: "team cmd delete", run: teamDeleteCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testTeamBasePath, status: http.StatusNoContent},
			setup:   setTeamCmdVars,
		},
		apiErrorCase("team cmd delete", teamDeleteCmd.RunE, testHTTPDelete, testTeamBasePath, http.StatusNotFound, "deleting team", setTeamCmdVars),
		{
			name: "team cmd members", run: teamCmdMembersCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testTeamBasePath + "/members"},
			response: `{"name":"devs","members":[{"name":"dave"}]}`, wantOutput: testTeamMember, setup: setTeamCmdVars,
		},
		apiErrorCase("team cmd members", teamCmdMembersCmd.RunE, testHTTPGet, testTeamBasePath+"/members", http.StatusNotFound, "getting team members", setTeamCmdVars),
		{
			name: "team cmd add member", run: teamAddMemberCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: testTeamMemberURL, status: http.StatusNoContent},
			setup:   setTeamCmdVars,
		},
		apiErrorCase("team cmd add member", teamAddMemberCmd.RunE, testHTTPPut, testTeamMemberURL, http.StatusBadRequest, "adding team member", setTeamCmdVars),
		{
			name: "team cmd remove member", run: teamRemoveMemberCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testTeamMemberURL, status: http.StatusNoContent},
			setup:   setTeamCmdVars,
		},
		apiErrorCase("team cmd remove member", teamRemoveMemberCmd.RunE, testHTTPDelete, testTeamMemberURL, http.StatusNotFound, "removing team member", setTeamCmdVars),
		{
			name: "team cmd permissions", run: teamPermissionsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testTeamBasePath + "/permissions"},
			response: `{"permissions":[{"repository":{"name":"testrepo"},"role":"write"}]}`, wantOutput: testRepository, setup: setTeamCmdVars,
		},
		apiErrorCase("team cmd permissions", teamPermissionsCmd.RunE, testHTTPGet, testTeamBasePath+"/permissions", http.StatusNotFound, "getting team permissions", setTeamCmdVars),
		{
			name: "team cmd set permission", run: teamSetPermissionCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: testTeamPermURL, body: map[string]interface{}{testFlagRole: testPermRole}, status: http.StatusNoContent},
			setup:   setTeamCmdVars,
		},
		apiErrorCase("team cmd set permission", teamSetPermissionCmd.RunE, testHTTPPut, testTeamPermURL, http.StatusBadRequest, "setting team permission", setTeamCmdVars),
		{
			name: "team cmd remove permission", run: teamRemovePermissionCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testTeamPermURL, status: http.StatusNoContent},
			setup:   setTeamCmdVars,
		},
		apiErrorCase("team cmd remove permission", teamRemovePermissionCmd.RunE, testHTTPDelete, testTeamPermURL, http.StatusNotFound, "removing team permission", setTeamCmdVars),
	}
	runCommandDomainCases(t, tests)
}

func mergeBody(base, extra map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(base)+len(extra))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range extra {
		merged[k] = v
	}
	return merged
}

func TestTeamDomainConfirmGuards(t *testing.T) {
	tests := []struct {
		name  string
		run   func(*cobra.Command, []string) error
		unset *bool
		want  string
	}{
		{"team cmd delete guard", teamDeleteCmd.RunE, &confirmTeamDelete, "are you sure you want to delete team"},
		{"team cmd remove member guard", teamRemoveMemberCmd.RunE, &confirmTeamMemberDel, "are you sure you want to remove"},
		{"team cmd remove permission guard", teamRemovePermissionCmd.RunE, &confirmTeamPermDel, "are you sure you want to remove permissions"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setTeamCmdVars(t)
			setCommandTestValue(t, test.unset, false)
			err := runCommandExpectingNoRequest(t, test.run)
			if err == nil || !strings.Contains(err.Error(), test.want) || !strings.Contains(err.Error(), "--confirm") {
				t.Fatalf("error = %v, want confirm prompt containing %q", err, test.want)
			}
		})
	}
}

func TestTeamDomainRequiredFlags(t *testing.T) {
	team := func(sub string) []string { return []string{cmdGet, cmdTeam, "-o", testNamespace, sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "team cmd no organization", args: []string{cmdGet, cmdTeam, subcmdList}, want: []string{testRequiredFlagError, cmdOrganization}},
		{name: "team cmd info", args: team(subcmdInfo), want: []string{testRequiredFlagError, testFlagName}},
		{name: "team cmd create", args: team(subcmdCreate), want: []string{testRequiredFlagError, testFlagName}},
		{name: "team cmd add member", args: team("add-member"), want: []string{testRequiredFlagError, testFlagName, testFlagMember}},
		{name: "team cmd set permission", args: team("set-permission"), want: []string{testRequiredFlagError, testFlagName}},
	})
}
