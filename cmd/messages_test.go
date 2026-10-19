package cmd

import (
	"net/http"
	"testing"
)

const (
	testMessagesPath = "/messages"
	testMessageText  = "maintenance tonight"
	testMessageJSON  = `{"uuid":"msg-1","content":"` + testMessageText + `","severity":"warning"}`
)

func setMessageVars(t *testing.T) {
	setCommandTestValue(t, &messageContent, testMessageText)
	setCommandTestValue(t, &messageSeverity, "warning")
	setCommandTestValue(t, &messageMediaType, "text/plain")
}

func TestMessagesCommands(t *testing.T) {
	tests := []commandDomainCase{
		{
			name: "messages list", run: messagesListCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testMessagesPath},
			response: `{"messages":[` + testMessageJSON + `]}`, wantOutput: testMessageText, setup: setMessageVars,
		},
		apiErrorCase("messages list", messagesListCmd.RunE, testHTTPGet, testMessagesPath, http.StatusInternalServerError, "getting messages", setMessageVars),
		{
			name: "messages create", run: messagesCreateCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPost, path: testMessagesPath, body: map[string]interface{}{
				"message": map[string]interface{}{"content": testMessageText, "severity": "warning", "media_type": "text/plain"},
			}},
			response: testMessageJSON, wantOutput: testMessageText, setup: setMessageVars,
		},
		apiErrorCase("messages create", messagesCreateCmd.RunE, testHTTPPost, testMessagesPath, http.StatusForbidden, "creating message", setMessageVars),
	}
	runCommandDomainCases(t, tests)
}

func TestMessagesRequiredFlags(t *testing.T) {
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "messages create", args: []string{cmdGet, "messages", subcmdCreate}, want: []string{testRequiredFlagError, "content", "severity"}},
	})
}
