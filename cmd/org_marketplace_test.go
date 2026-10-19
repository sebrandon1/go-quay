package cmd

import (
	"net/http"
	"testing"

	"github.com/spf13/cobra"
)

const (
	testMarketplacePath = "/organization/" + testNamespace + "/marketplace"
	testSubscriptionID  = "sub-1"
	testMarketplaceSKU  = "quay-pro"
)

func setMarketplaceVars(t *testing.T) {
	setCommandTestValue(t, &orgName, testNamespace)
	setCommandTestValue(t, &sku, testMarketplaceSKU)
	setCommandTestValue(t, &quantity, 2)
	setCommandTestValue(t, &subscriptionID, testSubscriptionID)
	setCommandTestValue(t, &subscriptionIDs, []string{testSubscriptionID, "sub-2"})
	setCommandTestValue(t, &confirm, true)
}

func TestMarketplaceCommands(t *testing.T) {
	tests := []commandDomainCase{
		{
			name: "marketplace info", run: marketplaceCmd.RunE,
			request:  commandRequestExpectation{method: testHTTPGet, path: testMarketplacePath},
			response: `{"subscriptions":[{"id":"sub-1","sku":"quay-pro"}],"has_payment":true}`, wantOutput: "has_payment", setup: setMarketplaceVars,
		},
		apiErrorCase("marketplace info", marketplaceCmd.RunE, testHTTPGet, testMarketplacePath, http.StatusNotFound, "getting marketplace info", setMarketplaceVars),
		{
			name: "marketplace create", run: createMarketplaceSubscriptionCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPost, path: testMarketplacePath, body: map[string]interface{}{"sku": testMarketplaceSKU, "quantity": 2}},
			setup:   setMarketplaceVars,
		},
		apiErrorCase("marketplace create", createMarketplaceSubscriptionCmd.RunE, testHTTPPost, testMarketplacePath, http.StatusBadRequest, "creating marketplace subscription", setMarketplaceVars),
		{
			name: "marketplace delete", run: deleteMarketplaceSubscriptionCmd.RunE,
			request: commandRequestExpectation{method: testHTTPDelete, path: testMarketplacePath + "/" + testSubscriptionID, status: http.StatusNoContent},
			setup:   setMarketplaceVars,
		},
		apiErrorCase("marketplace delete", deleteMarketplaceSubscriptionCmd.RunE, testHTTPDelete, testMarketplacePath+"/"+testSubscriptionID, http.StatusNotFound, "deleting marketplace subscription", setMarketplaceVars),
		{
			name: "marketplace batch remove", run: batchRemoveSubscriptionsCmd.RunE,
			request: commandRequestExpectation{method: testHTTPPost, path: testMarketplacePath + "/batchremove", body: map[string]interface{}{"subscription_ids": []string{testSubscriptionID, "sub-2"}}},
			setup:   setMarketplaceVars,
		},
		apiErrorCase("marketplace batch remove", batchRemoveSubscriptionsCmd.RunE, testHTTPPost, testMarketplacePath+"/batchremove", http.StatusBadRequest, "batch removing subscriptions", setMarketplaceVars),
	}
	runCommandDomainCases(t, tests)
}

func TestMarketplaceConfirmGuards(t *testing.T) {
	tests := []struct {
		name string
		run  func(*cobra.Command, []string) error
		want string
	}{
		{"marketplace delete guard", deleteMarketplaceSubscriptionCmd.RunE, "must pass --confirm to delete a marketplace subscription"},
		{"marketplace batch remove guard", batchRemoveSubscriptionsCmd.RunE, "must pass --confirm to batch remove subscriptions"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setMarketplaceVars(t)
			setCommandTestValue(t, &confirm, false)
			err := runCommandExpectingNoRequest(t, test.run)
			if err == nil || err.Error() != test.want {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestMarketplaceRequiredFlags(t *testing.T) {
	org := func(sub string) []string { return []string{cmdGet, cmdOrganization, sub} }
	runExecuteValidationCases(t, []executeValidationCase{
		{name: "marketplace create", args: org("create-marketplace-subscription"), want: []string{testRequiredFlagError, "sku"}},
		{name: "marketplace delete", args: org("delete-marketplace-subscription"), want: []string{testRequiredFlagError, "subscription-id"}},
		{name: "marketplace batch remove", args: org("batch-remove-subscriptions"), want: []string{testRequiredFlagError, "subscription-ids"}},
	})
}
