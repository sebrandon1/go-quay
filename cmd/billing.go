package cmd

import (
	"encoding/json"
	"fmt"

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
			panic(err)
		}

		billing, err := client.GetOrganizationBilling(billingOrgName)
		if err != nil {
			panic(err)
		}

		jsonOutput, err := json.Marshal(billing)
		if err != nil {
			panic(err)
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
			panic(err)
		}

		billing, err := client.GetUserBilling()
		if err != nil {
			panic(err)
		}

		jsonOutput, err := json.Marshal(billing)
		if err != nil {
			panic(err)
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
			panic(err)
		}

		subscription, err := client.GetOrganizationSubscription(billingOrgName)
		if err != nil {
			panic(err)
		}

		jsonOutput, err := json.Marshal(subscription)
		if err != nil {
			panic(err)
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
			panic(err)
		}

		subscription, err := client.GetUserSubscription()
		if err != nil {
			panic(err)
		}

		jsonOutput, err := json.Marshal(subscription)
		if err != nil {
			panic(err)
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
			panic(err)
		}

		invoices, err := client.GetOrganizationInvoices(billingOrgName)
		if err != nil {
			panic(err)
		}

		jsonOutput, err := json.Marshal(invoices)
		if err != nil {
			panic(err)
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

// NOTE: Usage commands are commented out because the usage endpoints don't exist in the Quay API
// orgUsageCmd gets usage statistics for an organization
// var orgUsageCmd = &cobra.Command{
// 	Use:   "org-usage",
// 	Short: "Get usage statistics for an organization",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Usage endpoints are not available in the Quay API")
// 	},
// }

// userUsageCmd gets usage statistics for the current user
// var userUsageCmd = &cobra.Command{
// 	Use:   "user-usage",
// 	Short: "Get usage statistics for the current user",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Usage endpoints are not available in the Quay API")
// 	},
// }

// plansCmd gets available subscription plans
var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Get available subscription plans",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(billingToken)
		if err != nil {
			panic(err)
		}

		plans, err := client.GetAvailablePlans()
		if err != nil {
			panic(err)
		}

		jsonOutput, err := json.Marshal(plans)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(jsonOutput))
	},
}

func init() {
	// Add persistent flags for commands that need organization name
	orgBillingCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	orgBillingCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := orgBillingCmd.MarkPersistentFlagRequired("organization"); err != nil {
		panic(err)
	}
	if err := orgBillingCmd.MarkPersistentFlagRequired("token"); err != nil {
		panic(err)
	}

	orgSubscriptionCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	orgSubscriptionCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := orgSubscriptionCmd.MarkPersistentFlagRequired("organization"); err != nil {
		panic(err)
	}
	if err := orgSubscriptionCmd.MarkPersistentFlagRequired("token"); err != nil {
		panic(err)
	}

	orgInvoicesCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	orgInvoicesCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := orgInvoicesCmd.MarkPersistentFlagRequired("organization"); err != nil {
		panic(err)
	}
	if err := orgInvoicesCmd.MarkPersistentFlagRequired("token"); err != nil {
		panic(err)
	}

	// orgUsageCmd flags commented out - endpoint doesn't exist
	// orgUsageCmd.PersistentFlags().StringVarP(&billingOrgName, "organization", "o", "", "Organization name")
	// orgUsageCmd.PersistentFlags().StringVarP(&billingPeriod, "period", "p", "", "Usage period (e.g., '30d', 'current_month')")
	// orgUsageCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	// if err := orgUsageCmd.MarkPersistentFlagRequired("organization"); err != nil {
	// 	panic(err)
	// }
	// if err := orgUsageCmd.MarkPersistentFlagRequired("token"); err != nil {
	// 	panic(err)
	// }

	// Add persistent flags for user commands
	userBillingCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := userBillingCmd.MarkPersistentFlagRequired("token"); err != nil {
		panic(err)
	}

	userSubscriptionCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := userSubscriptionCmd.MarkPersistentFlagRequired("token"); err != nil {
		panic(err)
	}

	userInvoicesCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := userInvoicesCmd.MarkPersistentFlagRequired("token"); err != nil {
		panic(err)
	}

	// userUsageCmd flags commented out - endpoint doesn't exist
	// userUsageCmd.PersistentFlags().StringVarP(&billingPeriod, "period", "p", "", "Usage period (e.g., '30d', 'current_month')")
	// userUsageCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	// if err := userUsageCmd.MarkPersistentFlagRequired("token"); err != nil {
	// 	panic(err)
	// }

	plansCmd.PersistentFlags().StringVarP(&billingToken, "token", "t", "", "Bearer token")
	if err := plansCmd.MarkPersistentFlagRequired("token"); err != nil {
		panic(err)
	}

	// Add subcommands to the billing command
	billingCmd.AddCommand(orgBillingCmd)
	billingCmd.AddCommand(userBillingCmd)
	billingCmd.AddCommand(orgSubscriptionCmd)
	billingCmd.AddCommand(userSubscriptionCmd)
	billingCmd.AddCommand(orgInvoicesCmd)
	billingCmd.AddCommand(userInvoicesCmd)
	// billingCmd.AddCommand(orgUsageCmd)  // Commented out - endpoint doesn't exist
	// billingCmd.AddCommand(userUsageCmd)  // Commented out - endpoint doesn't exist
	billingCmd.AddCommand(plansCmd)
}
