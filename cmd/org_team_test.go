package cmd

import (
	"net/http"
	"testing"
)

const (
	testOrgTeamsPath = "/organization/" + testNamespace + "/teams"
	testOrgTeamName  = "devs"
	testOrgTeamPath  = "/organization/" + testNamespace + "/team/" + testOrgTeamName
	testOrgTeamEmail = "dev@example.com"
	testOrgTeamJSON  = `{"name":"devs","role":"member"}`
)

func setOrgTeamVars(t *testing.T) {
	setCommandTestValue(t, &orgName, testNamespace)
	setCommandTestValue(t, &teamName, testOrgTeamName)
	setCommandTestValue(t, &email, testOrgTeamEmail)
	setCommandTestValue(t, &confirm, true)
}

func TestOrgTeamCommands(t *testing.T) {
	invitePath := testOrgTeamPath + "/invite/" + testOrgTeamEmail
	tests := []commandDomainCase{
		{
			name: "team list", run: orgTeamsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgTeamsPath},
			response: `{"teams":[` + testOrgTeamJSON + `]}`, wantOutput: testOrgTeamName, setup: setOrgTeamVars,
		},
		apiErrorCase("team list", orgTeamsCmd.RunE, testHTTPGet, testOrgTeamsPath, http.StatusForbidden, "getting organization teams", setOrgTeamVars),
		{
			name: "team info", run: teamInfoCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgTeamPath},
			response: testOrgTeamJSON, wantOutput: testOrgTeamName, setup: setOrgTeamVars,
		},
		apiErrorCase("team info", teamInfoCmd.RunE, testHTTPGet, testOrgTeamPath, http.StatusNotFound, "getting team", setOrgTeamVars),
		{
			name: "team members", run: teamMembersCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgTeamPath + "/members"},
			response: `{"name":"devs","members":[{"name":"` + testPermUser + `"}]}`, wantOutput: testPermUser, setup: setOrgTeamVars,
		},
		apiErrorCase("team members", teamMembersCmd.RunE, testHTTPGet, testOrgTeamPath+"/members", http.StatusNotFound, "getting team members", setOrgTeamVars),
		{
			name: "team invite", run: inviteMemberCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: invitePath, status: http.StatusNoContent},
			setup:   setOrgTeamVars,
		},
		apiErrorCase("team invite", inviteMemberCmd.RunE, testHTTPPut, invitePath, http.StatusBadRequest, "inviting team member", setOrgTeamVars),
		{
			name: "team cancel invite", run: cancelInviteCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: invitePath, status: http.StatusNoContent},
			setup:   setOrgTeamVars,
		},
		apiErrorCase("team cancel invite", cancelInviteCmd.RunE, testHTTPDelete, invitePath, http.StatusNotFound, "canceling team invite", setOrgTeamVars),
	}
	runCommandDomainCases(t, tests)
}

func TestOrgTeamCancelInviteRequiresConfirm(t *testing.T) {
	setOrgTeamVars(t)
	setCommandTestValue(t, &confirm, false)
	err := runCommandExpectingNoRequest(t, cancelInviteCmd.RunE)
	if err == nil || err.Error() != "must pass --confirm to cancel an invite" {
		t.Fatalf("error = %v, want confirm requirement", err)
	}
}

func TestOrgTeamRequiredFlags(t *testing.T) {
	org := func(sub string) []string { return []string{cmdGet, cmdOrganization, sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "team info", args: org("team"), want: []string{testRequiredFlagError, testFlagTeam}},
		{name: "team members", args: org("team-members"), want: []string{testRequiredFlagError, testFlagTeam}},
		{name: "team invite", args: org("invite-member"), want: []string{testRequiredFlagError, testFlagTeam, testFlagEmail}},
		{name: "team cancel invite", args: org("cancel-invite"), want: []string{testRequiredFlagError, testFlagTeam, testFlagEmail}},
	})
}
