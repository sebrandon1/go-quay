package cmd

import (
	"net/http"
	"testing"

	"github.com/spf13/cobra"
)

const (
	testOrgBodyJSON    = `{"name":"testns"}`
	testFlagMember     = "member"
	testOrgPath        = "/organization/" + testNamespace
	testOrgMember      = "bob"
	testOrgMemberPath  = testOrgPath + "/members/" + testOrgMember
	testOrgContactMail = "org@example.com"
)

func setOrgVars(t *testing.T) {
	setCommandTestValue(t, &orgName, testNamespace)
	setCommandTestValue(t, &email, testOrgContactMail)
	setCommandTestValue(t, &memberName, testOrgMember)
	setCommandTestValue(t, &confirm, true)
}

func TestOrganizationDomainCommands(t *testing.T) {
	tests := []commandDomainCase{
		{
			name: "org info", run: orgInfoCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgPath},
			response: testOrgBodyJSON, wantOutput: testNamespace, setup: setOrgVars,
		},
		apiErrorCase("org info", orgInfoCmd.RunE, testHTTPGet, testOrgPath, http.StatusNotFound, "getting organization", setOrgVars),
		{
			name: "org members", run: orgMembersCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgPath + "/members"},
			response: `{"members":[{"name":"bob","kind":"user"}]}`, wantOutput: testOrgMember, setup: setOrgVars,
		},
		apiErrorCase("org members", orgMembersCmd.RunE, testHTTPGet, testOrgPath+"/members", http.StatusForbidden, "getting organization members", setOrgVars),
		{
			name: "org create", run: createOrgCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPost, path: "/organization/", body: map[string]interface{}{testFlagName: testNamespace, testFlagEmail: testOrgContactMail}},
			response: testOrgBodyJSON, wantOutput: testNamespace, setup: setOrgVars,
		},
		apiErrorCase("org create", createOrgCmd.RunE, testHTTPPost, "/organization/", http.StatusBadRequest, "creating organization", setOrgVars),
		{
			name: "org update", run: updateOrgCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPPut, path: testOrgPath, body: map[string]interface{}{testFlagEmail: testOrgContactMail}},
			response: testOrgBodyJSON, wantOutput: testNamespace, setup: setOrgVars,
		},
		apiErrorCase("org update", updateOrgCmd.RunE, testHTTPPut, testOrgPath, http.StatusBadRequest, "updating organization", setOrgVars),
		{
			name: "org delete", run: deleteOrgCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testOrgPath, status: http.StatusNoContent},
			setup:   setOrgVars,
		},
		apiErrorCase("org delete", deleteOrgCmd.RunE, testHTTPDelete, testOrgPath, http.StatusNotFound, "deleting organization", setOrgVars),
		{
			name: "org add member", run: addMemberCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPut, path: testOrgMemberPath, status: http.StatusNoContent},
			setup:   setOrgVars,
		},
		apiErrorCase("org add member", addMemberCmd.RunE, testHTTPPut, testOrgMemberPath, http.StatusBadRequest, "adding organization member", setOrgVars),
		{
			name: "org remove member", run: removeMemberCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testOrgMemberPath, status: http.StatusNoContent},
			setup:   setOrgVars,
		},
		apiErrorCase("org remove member", removeMemberCmd.RunE, testHTTPDelete, testOrgMemberPath, http.StatusNotFound, "removing organization member", setOrgVars),
		{
			name: "org member", run: getMemberCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgMemberPath},
			response: `{"name":"bob","kind":"user"}`, wantOutput: testOrgMember, setup: setOrgVars,
		},
		apiErrorCase("org member", getMemberCmd.RunE, testHTTPGet, testOrgMemberPath, http.StatusNotFound, "getting organization member", setOrgVars),
		{
			name: "org collaborators", run: collaboratorsCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgPath + "/collaborators"},
			response: `{"collaborators":[{"name":"carol"}]}`, wantOutput: "carol", setup: setOrgVars,
		},
		apiErrorCase("org collaborators", collaboratorsCmd.RunE, testHTTPGet, testOrgPath+"/collaborators", http.StatusForbidden, "getting organization collaborators", setOrgVars),
		{
			name: "org repositories", run: orgRepositoriesCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testOrgPath + "/repositories"},
			response: `{"repositories":[{"name":"testrepo"}]}`, wantOutput: testRepository, setup: setOrgVars,
		},
		apiErrorCase("org repositories", orgRepositoriesCmd.RunE, testHTTPGet, testOrgPath+"/repositories", http.StatusForbidden, "getting organization repositories", setOrgVars),
	}
	runCommandDomainCases(t, tests)
}

func TestOrganizationConfirmGuards(t *testing.T) {
	tests := []struct {
		name string
		run  func(*cobra.Command, []string) error
		want string
	}{
		{"org delete guard", deleteOrgCmd.RunE, "must pass --confirm to delete an organization"},
		{"org remove member guard", removeMemberCmd.RunE, "must pass --confirm to remove a member"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setOrgVars(t)
			setCommandTestValue(t, &confirm, false)
			err := runCommandExpectingNoRequest(t, test.run)
			if err == nil || err.Error() != test.want {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestOrganizationRequiredFlags(t *testing.T) {
	org := func(sub string) []string { return []string{cmdGet, cmdOrganization, "-o", testNamespace, sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "org no organization", args: []string{cmdGet, cmdOrganization, subcmdInfo}, want: []string{testRequiredFlagError, cmdOrganization}},
		{name: "org create", args: org("create-org"), want: []string{testRequiredFlagError, testFlagEmail}},
		{name: "org update", args: org("update-org"), want: []string{testRequiredFlagError, testFlagEmail}},
		{name: "org add member", args: org("add-member"), want: []string{testRequiredFlagError, testFlagMember}},
		{name: "org remove member", args: org("remove-member"), want: []string{testRequiredFlagError, testFlagMember}},
		{name: "org member", args: org("member"), want: []string{testRequiredFlagError, testFlagMember}},
	})
}
