package cmd

import "testing"

func TestLogsCommands(t *testing.T) {
	setupRepository := func(t *testing.T) {
		setCommandTestValue(t, &namespace, testNamespace)
		setCommandTestValue(t, &repository, testRepository)
	}
	setupOrganization := func(t *testing.T) { setCommandTestValue(t, &orgName, testOrgName) }
	setupRange := func(t *testing.T) {
		setCommandTestValue(t, &startdate, testLogsStartDate)
		setCommandTestValue(t, &enddate, testLogsEndDate)
	}
	setupExport := func(t *testing.T) {
		setCommandTestValue(t, &starttime, testStartTime)
		setCommandTestValue(t, &endtime, "2026-01-31T23:59:59Z")
		setCommandTestValue(t, &callbackURL, "https://example.test/export")
		setCommandTestValue(t, &callbackEmail, "logs@example.test")
	}
	exportBody := map[string]interface{}{
		testLogsStartTimeParam: testStartTime, testLogsEndTimeParam: "2026-01-31T23:59:59Z",
		"callback_url": "https://example.test/export", "callback_email": "logs@example.test",
	}
	logResponse := `{"logs":[{"kind":"push"}],"next_page":"` + testLogsNextPageResult + `"}`
	aggregatedResponse := `{"aggregated":[{"kind":"push","count":2,"datetime":"2026-01-05"}]}`
	tests := []commandDomainCase{
		{name: "repository logs", executeArgs: []string{cmdGet, logsCmd.Use, repoLogsCmd.Use, testConfigNamespaceFlag, testNamespace, testDiffRepositoryFlag, testRepository, "--next-page", testLogsNextPage, "--startdate", testLogsStartDate, "--enddate", testLogsEndDate}, request: commandRequestExpectation{method: testHTTPGet, path: "/repository/testns/testrepo/logs", query: map[string]string{testLogsNextPageParam: testLogsNextPage, testLogsStartTimeParam: testLogsStartDate, testLogsEndTimeParam: testLogsEndDate}}, response: logResponse, wantOutput: testLogsNextPageResult},
		{name: "repository aggregated logs", run: repoAggregatedLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/repository/testns/testrepo/aggregatelogs", query: map[string]string{testLogsStartTimeParam: testLogsStartDate, testLogsEndTimeParam: testLogsEndDate}}, response: aggregatedResponse, wantOutput: testLogsPushKind, setup: func(t *testing.T) { setupRepository(t); setupRange(t) }},
		{name: "organization logs", run: orgLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/organization/testorg/logs", query: map[string]string{testLogsNextPageParam: testLogsNextPage, testLogsStartTimeParam: testLogsStartDate, testLogsEndTimeParam: testLogsEndDate}}, response: logResponse, wantOutput: testLogsNextPageResult, setup: func(t *testing.T) {
			setupOrganization(t)
			setupRange(t)
			setCommandTestValue(t, &nextPage, testLogsNextPage)
		}},
		{name: "organization aggregated logs", run: orgAggregatedLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/organization/testorg/aggregatelogs", query: map[string]string{testLogsStartTimeParam: testLogsStartDate, testLogsEndTimeParam: testLogsEndDate}}, response: aggregatedResponse, wantOutput: testLogsPushKind, setup: func(t *testing.T) { setupOrganization(t); setupRange(t) }},
		{name: "user logs", run: userLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user/logs", query: map[string]string{testLogsNextPageParam: testLogsNextPage, testLogsStartTimeParam: testLogsStartDate, testLogsEndTimeParam: testLogsEndDate}}, response: logResponse, wantOutput: testLogsNextPageResult, setup: func(t *testing.T) { setupRange(t); setCommandTestValue(t, &nextPage, testLogsNextPage) }},
		{name: "user aggregated logs", run: userAggregatedLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user/aggregatelogs", query: map[string]string{testLogsStartTimeParam: testLogsStartDate, testLogsEndTimeParam: testLogsEndDate}}, response: aggregatedResponse, wantOutput: testLogsPushKind, setup: func(t *testing.T) { setupRange(t) }},
		{name: "export organization logs", run: exportOrgLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: "/organization/testorg/exportlogs", body: exportBody}, setup: func(t *testing.T) { setupOrganization(t); setupExport(t) }},
		{name: "export user logs", run: exportUserLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: "/user/exportlogs", body: exportBody}, setup: setupExport},
		{name: "export repository logs", run: exportRepoLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: "/repository/testns/testrepo/exportlogs", body: exportBody}, setup: func(t *testing.T) { setupRepository(t); setupExport(t) }},
		{name: testAPIErrorCase, run: repoLogsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/repository/testns/testrepo/logs", status: 500}, response: testAPIErrorResponse, wantError: "getting repository logs", setup: setupRepository},
	}
	runCommandDomainCases(t, tests)
}
