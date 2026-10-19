package cmd

import (
	"net/http"
	"testing"
)

const testErrorTypePath = "/error/invalid_token"

func setErrorTypeVars(t *testing.T) {
	setCommandTestValue(t, &errorTypeName, "invalid_token")
}

func TestErrorTypeCommand(t *testing.T) {
	tests := []commandDomainCase{
		{
			name: "error type info", run: errorTypeCmd.RunE,
			request:    commandRequestExpectation{method: testHTTPGet, path: testErrorTypePath},
			response:   `{"type":"https://quay.io/api/v1/error/invalid_token","title":"invalid_token","description":"bad token"}`,
			wantOutput: "bad token", setup: setErrorTypeVars,
		},
		apiErrorCase("error type info", errorTypeCmd.RunE, testHTTPGet, testErrorTypePath, http.StatusNotFound, "getting error type", setErrorTypeVars),
	}
	runCommandDomainCases(t, tests)
}

func TestErrorTypeRequiredFlags(t *testing.T) {
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "error type", args: []string{cmdGet, "error"}, want: []string{testRequiredFlagError, "type"}},
	})
}
