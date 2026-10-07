package cmd

import (
	"strings"
	"testing"
)

func TestRobotCommands(t *testing.T) {
	setupName := func(t *testing.T) { setCommandTestValue(t, &robotShortname, "buildbot") }
	setupFederation := func(t *testing.T) {
		setupName(t)
		setCommandTestValue(t, &federationIssuer, "https://issuer.example")
		setCommandTestValue(t, &federationSubject, "repo:team/project")
	}
	tests := []commandDomainCase{
		{name: subcmdList, run: robotListCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user/robots"}, response: `{"robots":[{"name":"buildbot"}]}`, wantOutput: "buildbot"},
		{name: subcmdInfo, run: robotInfoCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: testRobotPath}, response: `{"name":"buildbot"}`, wantOutput: "buildbot", setup: setupName},
		{name: subcmdCreate, run: robotCreateCmd.RunE, request: commandRequestExpectation{method: testHTTPPut, path: testRobotPath, body: map[string]interface{}{"description": "build worker"}}, response: `{"name":"buildbot","token":"robot-token"}`, wantOutput: "robot-token", setup: func(t *testing.T) { setupName(t); setCommandTestValue(t, &robotDescription, "build worker") }},
		{name: subcmdDelete, run: robotDeleteCmd.RunE, request: commandRequestExpectation{method: testHTTPDelete, path: testRobotPath}, setup: func(t *testing.T) { setupName(t); setCommandTestValue(t, &confirmRobotDelete, true) }},
		{name: "regenerate", run: robotRegenerateCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: "/user/robots/buildbot/regenerate"}, response: `{"name":"buildbot","token":"new-token"}`, wantOutput: "new-token", setup: setupName},
		{name: "permissions", run: robotPermissionsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user/robots/buildbot/permissions"}, response: `{"permissions":[]}`, setup: setupName},
		{name: "federation get", run: robotFederationGetCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: testRobotFederationPath}, response: `{"federation":[{"issuer":"https://issuer.example","subject":"repo:team/project"}]}`, wantOutput: "issuer.example", setup: setupName},
		{name: "federation create", run: robotFederationCreateCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: testRobotFederationPath, body: []map[string]string{{"issuer": "https://issuer.example", "subject": "repo:team/project"}}}, setup: setupFederation},
		{name: "federation delete", run: robotFederationDeleteCmd.RunE, request: commandRequestExpectation{method: testHTTPDelete, path: testRobotFederationPath}, setup: setupName},
		{name: testAPIErrorCase, run: robotListCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user/robots", status: 500}, response: testAPIErrorResponse, wantError: "getting robot accounts"},
	}
	runCommandDomainCases(t, tests)
}

func TestRobotDeleteRequiresConfirmation(t *testing.T) {
	setCommandTestValue(t, &robotShortname, "buildbot")
	setCommandTestValue(t, &confirmRobotDelete, false)
	setCommandTestValue(t, &token, testTokenValue)
	setCommandTestValue(t, &quayURL, testUnreachableURL)
	if err := robotDeleteCmd.RunE(nil, nil); err == nil || !strings.Contains(err.Error(), "Use --confirm to proceed") {
		t.Fatalf("robot delete error = %v, want confirmation prompt", err)
	}
}
