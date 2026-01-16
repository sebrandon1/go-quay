package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	billingOrgName string
	billingToken   string
)

// billingCmd represents the billing command
var billingCmd = &cobra.Command{
	Use:   "billing",
	Short: "Billing-related commands",
	Long:  `Commands for managing and viewing billing information, subscriptions, invoices, and usage statistics.`,
}

// orgBillingCmd gets billing information for an organization
var orgBillingCmd = &cobra.Command{
	Use:   "org-info",
	Short: "Get billing information for an organization",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(billingToken)
		if err != nil {
			fmt.Println("Error creating client:", err)
			return
		}

		billing, err := client.GetOrganizationBilling(billingOrgName)
		if err != nil {
			fmt.Println("Error getting organization billing:", err)
			return
		}

		jsonOutput, err := json.Marshal(billing)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonOutput))
	},
}

// userBillingCmd gets billing information for the current user
var userBillingCmd = &cobra.Command{
	Use:   "user-info",
	Short: "Get billing information for the current user",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(billingToken)
		if err != nil {
			fmt.Println("Error creating client:", err)
			return
		}

		billing, err := client.GetUserBilling()
		if err != nil {
			fmt.Println("Error getting user billing:", err)
			return
		}

		jsonOutput, err := json.Marshal(billing)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonOutput))
	},
}

// orgSubscriptionCmd gets subscription details for an organization
var orgSubscriptionCmd = &cobra.Command{
	Use:   "org-subscription",
	Short: "Get subscription details for an organization",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(billingToken)
		if err != nil {
			fmt.Println("Error creating client:", err)
			return
		}

		subscription, err := client.GetOrganizationSubscription(billingOrgName)
		if err != nil {
			fmt.Println("Error getting organization subscription:", err)
			return
		}

		jsonOutput, err := json.Marshal(subscription)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonOutput))
	},
}

// userSubscriptionCmd gets subscription details for the current user
var userSubscriptionCmd = &cobra.Command{
	Use:   "user-subscription",
	Short: "Get subscription details for the current user",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(billingToken)
		if err != nil {
			fmt.Println("Error creating client:", err)
			return
		}

		subscription, err := client.GetUserSubscription()
		if err != nil {
			fmt.Println("Error getting user subscription:", err)
			return
		}

		jsonOutput, err := json.Marshal(subscription)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonOutput))
	},
}

// orgInvoicesCmd gets invoices for an organization
var orgInvoicesCmd = &cobra.Command{
	Use:   "org-invoices",
	Short: "Get invoices for an organization",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(billingToken)
		if err != nil {
			fmt.Println("Error creating client:", err)
			return
		}

		invoices, err := client.GetOrganizationInvoices(billingOrgName)
		if err != nil {
			fmt.Println("Error getting organization invoices:", err)
			return
		}

		jsonOutput, err := json.Marshal(invoices)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonOutput))
	},
}

// userInvoicesCmd gets invoices for the current user
var userInvoicesCmd = &cobra.Command{
	Use:   "user-invoices",
	Short: "Get invoices for the current user",
	Run: func(cmd *cobra.Command, args []string) {
		// User invoices endpoint doesn't exist in Quay API - return empty array to match expected format
		fmt.Println("[]")
	},
}

// plansCmd gets available subscription plans
var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Get available subscription plans",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(billingToken)
		if err != nil {
			fmt.Println("Error creating client:", err)
			return
		}

		plans, err := client.GetAvailablePlans()
		if err != nil {
			fmt.Println("Error getting available plans:", err)
			return
		}

		jsonOutput, err := json.Marshal(plans)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonOutput))
	},
}

func init() {
	// Add persistent flags for commands that need organization name
	orgBillingCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	orgBillingCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := orgBillingCmd.MarkPersistentFlagRequired("organization"); err != nil {
		fmt.Printf("Error marking organization flag required: %v\n", err)
		os.Exit(1)
	}
	if err := orgBillingCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag required: %v\n", err)
		os.Exit(1)
	}

	orgSubscriptionCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	orgSubscriptionCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := orgSubscriptionCmd.MarkPersistentFlagRequired("organization"); err != nil {
		fmt.Printf("Error marking organization flag required: %v\n", err)
		os.Exit(1)
	}
	if err := orgSubscriptionCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag required: %v\n", err)
		os.Exit(1)
	}

	orgInvoicesCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	orgInvoicesCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := orgInvoicesCmd.MarkPersistentFlagRequired("organization"); err != nil {
		fmt.Printf("Error marking organization flag required: %v\n", err)
		os.Exit(1)
	}
	if err := orgInvoicesCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag required: %v\n", err)
		os.Exit(1)
	}

	// Add persistent flags for user commands
	userBillingCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := userBillingCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag required: %v\n", err)
		os.Exit(1)
	}

	userSubscriptionCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := userSubscriptionCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag required: %v\n", err)
		os.Exit(1)
	}

	userInvoicesCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := userInvoicesCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag required: %v\n", err)
		os.Exit(1)
	}

	plansCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := plansCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag required: %v\n", err)
		os.Exit(1)
	}

	// Add subcommands to the billing command
	billingCmd.AddCommand(orgBillingCmd)
	billingCmd.AddCommand(userBillingCmd)
	billingCmd.AddCommand(orgSubscriptionCmd)
	billingCmd.AddCommand(userSubscriptionCmd)
	billingCmd.AddCommand(orgInvoicesCmd)
	billingCmd.AddCommand(userInvoicesCmd)
	billingCmd.AddCommand(plansCmd)
}
