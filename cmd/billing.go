package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	billingOrgName string
)

var billingCmd = &cobra.Command{
	Use:   "billing",
	Short: "Billing-related commands",
	Long:  `Commands for managing and viewing billing information, subscriptions, invoices, and usage statistics.`,
}

var orgBillingCmd = &cobra.Command{
	Use:   "org-info",
	Short: "Get billing information for an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		billing, err := client.GetOrganizationBilling(billingOrgName)
		if err != nil {
			return fmt.Errorf("getting organization billing: %w", err)
		}

		return printJSON(billing)
	},
}

var userBillingCmd = &cobra.Command{
	Use:   "user-info",
	Short: "Get billing information for the current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		billing, err := client.GetUserBilling()
		if err != nil {
			return fmt.Errorf("getting user billing: %w", err)
		}

		return printJSON(billing)
	},
}

var orgSubscriptionCmd = &cobra.Command{
	Use:   "org-subscription",
	Short: "Get subscription details for an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		subscription, err := client.GetOrganizationSubscription(billingOrgName)
		if err != nil {
			return fmt.Errorf("getting organization subscription: %w", err)
		}

		return printJSON(subscription)
	},
}

var userSubscriptionCmd = &cobra.Command{
	Use:   "user-subscription",
	Short: "Get subscription details for the current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		subscription, err := client.GetUserSubscription()
		if err != nil {
			return fmt.Errorf("getting user subscription: %w", err)
		}

		return printJSON(subscription)
	},
}

var orgInvoicesCmd = &cobra.Command{
	Use:   "org-invoices",
	Short: "Get invoices for an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		invoices, err := client.GetOrganizationInvoices(billingOrgName)
		if err != nil {
			return fmt.Errorf("getting organization invoices: %w", err)
		}

		return printJSON(invoices)
	},
}

var userInvoicesCmd = &cobra.Command{
	Use:   "user-invoices",
	Short: "Get invoices for the current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		invoices, err := client.GetUserInvoices()
		if err != nil {
			return fmt.Errorf("getting user invoices: %w", err)
		}

		return printJSON(invoices)
	},
}

var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Get available subscription plans",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		plans, err := client.GetAvailablePlans()
		if err != nil {
			return fmt.Errorf("getting available plans: %w", err)
		}

		return printJSON(plans)
	},
}

func init() {
	orgBillingCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	_ = orgBillingCmd.MarkPersistentFlagRequired("organization")

	orgSubscriptionCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	_ = orgSubscriptionCmd.MarkPersistentFlagRequired("organization")

	orgInvoicesCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	_ = orgInvoicesCmd.MarkPersistentFlagRequired("organization")

	billingCmd.AddCommand(orgBillingCmd)
	billingCmd.AddCommand(userBillingCmd)
	billingCmd.AddCommand(orgSubscriptionCmd)
	billingCmd.AddCommand(userSubscriptionCmd)
	billingCmd.AddCommand(orgInvoicesCmd)
	billingCmd.AddCommand(userInvoicesCmd)
	billingCmd.AddCommand(plansCmd)
}
