package cmd

import (
	"net/http"
	"strings"
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

// TestOrgTeamShims exercises the deprecated `get organization team*` paths. They
// must keep working and route to the same shared implementation as cmd/team.go.
func TestOrgTeamShims(t *testing.T) {
	invitePath := testOrgTeamPath + "/invite/" + testOrgTeamEmail
	tests := []commandDomainCase{
		{
			name: "shim team list", run: orgTeamsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgTeamsPath},
			response: `{"teams":[` + testOrgTeamJSON + `]}`, wantOutput: testOrgTeamName, setup: setOrgTeamVars,
		},
		apiErrorCase("shim team list", orgTeamsCmd.RunE, testHTTPGet, testOrgTeamsPath, http.StatusForbidden, "getting teams", setOrgTeamVars),
		{
			name: "shim team info", run: teamInfoCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgTeamPath},
			response: testOrgTeamJSON, wantOutput: testOrgTeamName, setup: setOrgTeamVars,
		},
		apiErrorCase("shim team info", teamInfoCmd.RunE, testHTTPGet, testOrgTeamPath, http.StatusNotFound, "getting team", setOrgTeamVars),
		{
			name: "shim team members", run: teamMembersCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgTeamPath + "/members"},
			response: `{"name":"devs","members":[{"name":"` + testPermUser + `"}]}`, wantOutput: testPermUser, setup: setOrgTeamVars,
		},
		apiErrorCase("shim team members", teamMembersCmd.RunE, testHTTPGet, testOrgTeamPath+"/members", http.StatusNotFound, "getting team members", setOrgTeamVars),
		{
			name: "shim team invite", run: orgInviteMemberCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: invitePath, status: http.StatusNoContent},
			setup:   setOrgTeamVars,
		},
		apiErrorCase("shim team invite", orgInviteMemberCmd.RunE, testHTTPPut, invitePath, http.StatusBadRequest, "inviting team member", setOrgTeamVars),
		{
			name: "shim team cancel invite", run: orgCancelInviteCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: invitePath, status: http.StatusNoContent},
			setup:   setOrgTeamVars,
		},
		apiErrorCase("shim team cancel invite", orgCancelInviteCmd.RunE, testHTTPDelete, invitePath, http.StatusNotFound, "canceling team invite", setOrgTeamVars),
	}
	runCommandDomainCases(t, tests)
}

// TestOrgTeamShimDeprecationMessages guards the migration window: the legacy
// commands stay runnable but must announce the canonical verb-first replacement.
func TestOrgTeamShimDeprecationMessages(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"teams", orgTeamsCmd.Deprecated, cliName + " " + cmdList + " teams"},
		{cmdTeam, teamInfoCmd.Deprecated, cliName + " " + cmdInfo + " " + cmdTeam},
		{"team-members", teamMembersCmd.Deprecated, cliName + " " + cmdList + " team-members"},
		{"invite-member", orgInviteMemberCmd.Deprecated, cliName + " " + cmdCreate + " invite"},
		{"cancel-invite", orgCancelInviteCmd.Deprecated, cliName + " " + cmdDelete + " invite"},
	}
	for _, test := range tests {
		if !strings.Contains(test.got, test.want) {
			t.Errorf("%s deprecation = %q, want it to mention %q", test.name, test.got, test.want)
		}
	}
}

func TestOrgTeamCancelInviteShimRequiresConfirm(t *testing.T) {
	setOrgTeamVars(t)
	setCommandTestValue(t, &confirm, false)
	err := runCommandExpectingNoRequest(t, orgCancelInviteCmd.RunE)
	if err == nil || err.Error() != errConfirmCancelInvite.Error() {
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
