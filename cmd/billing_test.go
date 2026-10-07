package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestBillingCommands(t *testing.T) {
	setupOrg := func(t *testing.T) { setCommandTestValue(t, &billingOrgName, testOrgName) }
	tests := []commandDomainCase{
		{name: "organization info", run: orgBillingCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/organization/testorg/plan"}, response: `{"plan":"enterprise"}`, wantOutput: "enterprise", setup: setupOrg},
		{name: "user info", executeArgs: []string{cmdGet, billingCmd.Use, userBillingCmd.Use}, request: commandRequestExpectation{method: testHTTPGet, path: "/user/plan"}, response: `{"plan":"free"}`, wantOutput: "free"},
		{name: "organization subscription", run: orgSubscriptionCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/organization/testorg/plan"}, response: `{"name":"enterprise"}`, wantOutput: "enterprise", setup: setupOrg},
		{name: "user subscription", run: userSubscriptionCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/user/plan"}, response: `{"name":"free"}`, wantOutput: "free"},
		{name: "organization invoices", run: orgInvoicesCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/organization/testorg/invoices"}, response: `{"invoices":[{"id":"invoice-1"}]}`, wantOutput: "invoice-1", setup: setupOrg},
		{name: "available plans", run: plansCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/plans"}, response: `{"plans":[{"name":"starter"}]}`, wantOutput: "starter"},
		{name: testAPIErrorCase, run: plansCmd.RunE, request: commandRequestExpectation{method: testHTTPGet, path: "/plans", status: 500}, response: testAPIErrorResponse, wantError: "getting available plans"},
	}
	runCommandDomainCases(t, tests)
}

func TestUserInvoicesCommandReportsUnsupportedEndpoint(t *testing.T) {
	setCommandTestValue(t, &token, testTokenValue)
	setCommandTestValue(t, &quayURL, testUnreachableURL)
	command := &cobra.Command{}
	if err := userInvoicesCmd.RunE(command, nil); err == nil || !strings.Contains(err.Error(), "user invoices endpoint not available") {
		t.Fatalf("user invoices error = %v, want unsupported endpoint error", err)
	}
}
