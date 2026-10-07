package cmd

import "testing"

func TestUserCommands(t *testing.T) {
	setupRepository := func(t *testing.T) {
		setCommandTestValue(t, &namespace, testNamespace)
		setCommandTestValue(t, &repository, testRepository)
	}
	setupLookup := func(t *testing.T) { setCommandTestValue(t, &lookupUsername, "alice") }
	tests := []commandDomainCase{
		{name: subcmdInfo, executeArgs: []string{subcmdInfo, cmdUser}, request: commandRequestExpectation{method: testHTTPGet, path: "/user"}, response: `{"username":"alice"}`, wantOutput: "alice"},
		{name: "starred repositories", run: userStarredCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user/starred"}, response: `{"repositories":[{"name":"widget"}]}`, wantOutput: "widget"},
		{name: "star repository", run: starRepoCmd.RunE, request: commandRequestExpectation{method: testHTTPPut, path: "/repository/testns/testrepo/star"}, setup: setupRepository},
		{name: "unstar repository", run: unstarRepoCmd.RunE, request: commandRequestExpectation{method: testHTTPDelete, path: "/repository/testns/testrepo/star"}, setup: setupRepository},
		{name: "lookup", run: userLookupCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/users/alice"}, response: `{"username":"alice"}`, wantOutput: "alice", setup: setupLookup},
		{name: "marketplace", run: userMarketplaceCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user/marketplace"}, response: `{}`},
		{name: testAPIErrorCase, run: userInfoCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user", status: 500}, response: testAPIErrorResponse, wantError: "getting user information"},
	}
	runCommandDomainCases(t, tests)
}
