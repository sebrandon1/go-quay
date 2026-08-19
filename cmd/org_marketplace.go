package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

// Get Organization Marketplace
var marketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Get organization marketplace information",
	Long:  `Get marketplace information for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		marketplace, err := client.GetOrganizationMarketplace(orgName)
		if err != nil {
			return fmt.Errorf("getting marketplace info: %w", err)
		}
		return printJSON(marketplace)
	},
}

// Create Marketplace Subscription
var createMarketplaceSubscriptionCmd = &cobra.Command{
	Use:   "create-marketplace-subscription",
	Short: "Create a marketplace subscription",
	Long:  `Create a marketplace subscription for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.CreateOrganizationMarketplaceSubscription(orgName, &lib.MarketplaceSubscriptionRequest{SKU: sku, Quantity: quantity})
		if err != nil {
			return fmt.Errorf("creating marketplace subscription: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Marketplace subscription created successfully")
		return nil
	},
}

// Delete Marketplace Subscription
var deleteMarketplaceSubscriptionCmd = &cobra.Command{
	Use:   "delete-marketplace-subscription",
	Short: "Delete a marketplace subscription",
	Long:  `Delete a marketplace subscription. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete a marketplace subscription")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteOrganizationMarketplaceSubscription(orgName, subscriptionID)
		if err != nil {
			return fmt.Errorf("deleting marketplace subscription: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Marketplace subscription deleted successfully")
		return nil
	},
}

// Batch Remove Marketplace Subscriptions
var batchRemoveSubscriptionsCmd = &cobra.Command{
	Use:   "batch-remove-subscriptions",
	Short: "Batch remove marketplace subscriptions",
	Long:  `Remove multiple marketplace subscriptions at once. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to batch remove subscriptions")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.BatchRemoveOrganizationMarketplaceSubscriptions(orgName, subscriptionIDs)
		if err != nil {
			return fmt.Errorf("batch removing subscriptions: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Marketplace subscriptions removed successfully")
		return nil
	},
}

func initOrgMarketplaceFlags() {
	createMarketplaceSubscriptionCmd.Flags().StringVar(&sku, "sku", "", "Subscription SKU")
	_ = createMarketplaceSubscriptionCmd.MarkFlagRequired("sku")
	createMarketplaceSubscriptionCmd.Flags().IntVar(&quantity, "quantity", 0, "Subscription quantity")

	deleteMarketplaceSubscriptionCmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "Subscription ID")
	_ = deleteMarketplaceSubscriptionCmd.MarkFlagRequired("subscription-id")
	deleteMarketplaceSubscriptionCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	batchRemoveSubscriptionsCmd.Flags().StringSliceVar(&subscriptionIDs, "subscription-ids", nil, "Comma-separated subscription IDs")
	_ = batchRemoveSubscriptionsCmd.MarkFlagRequired("subscription-ids")
	batchRemoveSubscriptionsCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm removal")
}

func init() {
	organizationCmd.AddCommand(marketplaceCmd)
	organizationCmd.AddCommand(createMarketplaceSubscriptionCmd)
	organizationCmd.AddCommand(deleteMarketplaceSubscriptionCmd)
	organizationCmd.AddCommand(batchRemoveSubscriptionsCmd)

	initOrgMarketplaceFlags()
}
