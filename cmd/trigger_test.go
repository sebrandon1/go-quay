package cmd

import "testing"

func TestTriggerCommands(t *testing.T) {
	const basePath = "/repository/testns/testrepo/trigger/"
	setupRepository := func(t *testing.T) {
		setCommandTestValue(t, &triggerNamespace, testNamespace)
		setCommandTestValue(t, &triggerRepository, testRepository)
	}
	setupUUID := func(t *testing.T) { setupRepository(t); setCommandTestValue(t, &triggerUUID, testTriggerID) }
	tests := []commandDomainCase{
		{name: subcmdList, run: triggerListCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: basePath}, response: `{"triggers":[` + testTriggerResponse + `]}`, wantOutput: testTriggerID, setup: setupRepository},
		{name: subcmdInfo, run: triggerInfoCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: basePath + testTriggerID}, response: testTriggerResponse, wantOutput: testTriggerID, setup: setupUUID},
		{name: subcmdDelete, run: triggerDeleteCmd.RunE, request: commandRequestExpectation{method: testHTTPDelete, path: basePath + testTriggerID}, setup: setupUUID},
		{name: "enable", run: triggerEnableCmd.RunE, request: commandRequestExpectation{method: testHTTPPut, path: basePath + testTriggerID, body: map[string]interface{}{"enabled": true}}, response: `{"id":"` + testTriggerID + `","enabled":true}`, wantOutput: testTriggerID, setup: setupUUID},
		{name: "disable", run: triggerDisableCmd.RunE, request: commandRequestExpectation{method: testHTTPPut, path: basePath + testTriggerID, body: map[string]interface{}{"enabled": false}}, response: `{"id":"` + testTriggerID + `","enabled":false}`, wantOutput: testTriggerID, setup: setupUUID},
		{name: "start with commit", run: triggerStartCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath + testTriggerID + "/start", body: map[string]interface{}{"commit_sha": "deadbeef"}}, response: testBuildResponse, wantOutput: testBuildID, setup: func(t *testing.T) { setupUUID(t); setCommandTestValue(t, &triggerCommitSHA, "deadbeef") }},
		{name: "start without commit", run: triggerStartCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath + testTriggerID + "/start"}, response: testBuildResponse, wantOutput: testBuildID, setup: setupUUID},
		{name: "activate with pull robot", run: triggerActivateCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath + testTriggerID + "/activate", body: map[string]interface{}{"pull_robot": "testns+builder"}}, response: `{"id":"` + testTriggerID + `"}`, wantOutput: testTriggerID, setup: func(t *testing.T) { setupUUID(t); setCommandTestValue(t, &triggerPullRobot, "testns+builder") }},
		{name: "activate without pull robot", run: triggerActivateCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath + testTriggerID + "/activate", body: map[string]interface{}{}}, response: `{"id":"` + testTriggerID + `"}`, setup: setupUUID},
		{name: "builds", run: triggerBuildsCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: basePath + "trigger-1/builds", query: map[string]string{testLimitParam: "3"}}, response: `{"builds":[{"id":"build-2"}]}`, wantOutput: "build-2", setup: func(t *testing.T) { setupUUID(t); setCommandTestValue(t, &triggerBuildLimit, 3) }},
		{name: testAPIErrorCase, run: triggerListCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: basePath, status: 500}, response: testAPIErrorResponse, wantError: "getting triggers", setup: setupRepository},
	}
	runCommandDomainCases(t, tests)
}
