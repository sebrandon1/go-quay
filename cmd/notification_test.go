package cmd

import (
	"strings"
	"testing"
)

func TestNotificationCommands(t *testing.T) {
	const basePath = "/repository/testns/testrepo/notification/"
	setupRepository := func(t *testing.T) {
		setCommandTestValue(t, &notificationNamespace, testNamespace)
		setCommandTestValue(t, &notificationRepository, testRepository)
	}
	setupUUID := func(t *testing.T) {
		setupRepository(t)
		setCommandTestValue(t, &notificationUUID, testNotificationUUID)
	}
	setupWebhook := func(t *testing.T) {
		setCommandTestValue(t, &notificationEvent, "repo_push")
		setCommandTestValue(t, &notificationMethod, "webhook")
		setCommandTestValue(t, &notificationURL, "https://hooks.example.test/quay")
		setCommandTestValue(t, &notificationTitle, "push alert")
	}
	requestBody := map[string]interface{}{
		testNotificationEventKey: "repo_push", testNotificationMethodKey: "webhook", testConfigRootCommand: map[string]interface{}{"url": "https://hooks.example.test/quay"}, "title": "push alert",
	}
	emailBody := map[string]interface{}{
		testNotificationEventKey: "build_failure", testNotificationMethodKey: "email", testConfigRootCommand: map[string]interface{}{"email": "ops@example.test"},
	}
	slackBody := map[string]interface{}{
		testNotificationEventKey: "build_success", testNotificationMethodKey: "slack", testConfigRootCommand: map[string]interface{}{"url": "https://hooks.slack.test/quay"},
	}
	tests := []commandDomainCase{
		{name: subcmdList, executeArgs: []string{subcmdList, "notifications", testConfigNamespaceFlag, testNamespace, testDiffRepositoryFlag, testRepository}, request: commandRequestExpectation{method: testHTTPGet, path: basePath}, response: `{"notifications":[` + testNotificationResponse + `]}`, wantOutput: testNotificationUUID},
		{name: subcmdInfo, run: notificationInfoCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: basePath + testNotificationUUID}, response: testNotificationResponse, wantOutput: testNotificationUUID, setup: setupUUID},
		{name: subcmdCreate, run: notificationCreateCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath, body: requestBody}, response: testNotificationResponse, wantOutput: testNotificationUUID, setup: func(t *testing.T) { setupRepository(t); setupWebhook(t) }},
		{name: "create email", run: notificationCreateCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath, body: emailBody}, response: `{"uuid":"notice-2"}`, wantOutput: "notice-2", setup: func(t *testing.T) {
			setupRepository(t)
			setCommandTestValue(t, &notificationEvent, "build_failure")
			setCommandTestValue(t, &notificationMethod, "email")
			setCommandTestValue(t, &notificationURL, "ops@example.test")
		}},
		{name: "create slack", run: notificationCreateCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath, body: slackBody}, response: `{"uuid":"notice-3"}`, wantOutput: "notice-3", setup: func(t *testing.T) {
			setupRepository(t)
			setCommandTestValue(t, &notificationEvent, "build_success")
			setCommandTestValue(t, &notificationMethod, "slack")
			setCommandTestValue(t, &notificationURL, "https://hooks.slack.test/quay")
		}},
		{name: subcmdDelete, run: notificationDeleteCmd.RunE, request: commandRequestExpectation{method: testHTTPDelete, path: basePath + testNotificationUUID}, setup: func(t *testing.T) { setupUUID(t); setCommandTestValue(t, &confirmNotificationDel, true) }},
		{name: "test", run: notificationTestCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath + "notice-1/test"}, setup: setupUUID},
		{name: "reset", run: notificationResetCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath + "notice-1/reset"}, setup: setupUUID},
		{name: "update", run: notificationUpdateCmd.RunE, request: commandRequestExpectation{method: testHTTPPost, path: basePath + testNotificationUUID, body: requestBody}, response: testNotificationResponse, wantOutput: testNotificationUUID, setup: func(t *testing.T) { setupUUID(t); setupWebhook(t) }},
		{name: testAPIErrorCase, run: notificationListCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: basePath, status: 500}, response: testAPIErrorResponse, wantError: "getting notifications", setup: setupRepository},
	}
	runCommandDomainCases(t, tests)
}

func TestNotificationDeleteRequiresConfirmation(t *testing.T) {
	setCommandTestValue(t, &notificationNamespace, testNamespace)
	setCommandTestValue(t, &notificationRepository, testRepository)
	setCommandTestValue(t, &notificationUUID, testNotificationUUID)
	setCommandTestValue(t, &confirmNotificationDel, false)
	setCommandTestValue(t, &token, testTokenValue)
	setCommandTestValue(t, &quayURL, testUnreachableURL)
	if err := notificationDeleteCmd.RunE(nil, nil); err == nil || !strings.Contains(err.Error(), "Use --confirm to proceed") {
		t.Fatalf("notification delete error = %v, want confirmation prompt", err)
	}
}
