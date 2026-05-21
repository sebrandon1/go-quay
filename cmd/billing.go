package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
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
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		billing, err := client.GetOrganizationBilling(billingOrgName)
		if err != nil {
			fmt.Printf("Error getting organization billing: %v\n", err)
			os.Exit(1)
		}

		printJSON(billing)
	},
}

var userBillingCmd = &cobra.Command{
	Use:   "user-info",
	Short: "Get billing information for the current user",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		billing, err := client.GetUserBilling()
		if err != nil {
			fmt.Printf("Error getting user billing: %v\n", err)
			os.Exit(1)
		}

		printJSON(billing)
	},
}

var orgSubscriptionCmd = &cobra.Command{
	Use:   "org-subscription",
	Short: "Get subscription details for an organization",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		subscription, err := client.GetOrganizationSubscription(billingOrgName)
		if err != nil {
			fmt.Printf("Error getting organization subscription: %v\n", err)
			os.Exit(1)
		}

		printJSON(subscription)
	},
}

var userSubscriptionCmd = &cobra.Command{
	Use:   "user-subscription",
	Short: "Get subscription details for the current user",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		subscription, err := client.GetUserSubscription()
		if err != nil {
			fmt.Printf("Error getting user subscription: %v\n", err)
			os.Exit(1)
		}

		printJSON(subscription)
	},
}

var orgInvoicesCmd = &cobra.Command{
	Use:   "org-invoices",
	Short: "Get invoices for an organization",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		invoices, err := client.GetOrganizationInvoices(billingOrgName)
		if err != nil {
			fmt.Printf("Error getting organization invoices: %v\n", err)
			os.Exit(1)
		}

		printJSON(invoices)
	},
}

var userInvoicesCmd = &cobra.Command{
	Use:   "user-invoices",
	Short: "Get invoices for the current user",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		invoices, err := client.GetUserInvoices()
		if err != nil {
			fmt.Printf("Error getting user invoices: %v\n", err)
			os.Exit(1)
		}

		printJSON(invoices)
	},
}

var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Get available subscription plans",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		plans, err := client.GetAvailablePlans()
		if err != nil {
			fmt.Printf("Error getting available plans: %v\n", err)
			os.Exit(1)
		}

		printJSON(plans)
	},
}

func init() {
	orgBillingCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	if err := orgBillingCmd.MarkPersistentFlagRequired("organization"); err != nil {
		fmt.Printf("Error marking organization flag required: %v\n", err)
		os.Exit(1)
	}

	orgSubscriptionCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	if err := orgSubscriptionCmd.MarkPersistentFlagRequired("organization"); err != nil {
		fmt.Printf("Error marking organization flag required: %v\n", err)
		os.Exit(1)
	}

	orgInvoicesCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	if err := orgInvoicesCmd.MarkPersistentFlagRequired("organization"); err != nil {
		fmt.Printf("Error marking organization flag required: %v\n", err)
		os.Exit(1)
	}

	billingCmd.AddCommand(orgBillingCmd)
	billingCmd.AddCommand(userBillingCmd)
	billingCmd.AddCommand(orgSubscriptionCmd)
	billingCmd.AddCommand(userSubscriptionCmd)
	billingCmd.AddCommand(orgInvoicesCmd)
	billingCmd.AddCommand(userInvoicesCmd)
	billingCmd.AddCommand(plansCmd)
}
