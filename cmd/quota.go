package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Organization Quota
var orgQuotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Get organization quota",
	Long:  `Get quota information for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		quota, err := client.GetQuota(orgName)
		if err != nil {
			return fmt.Errorf("getting organization quota: %w", err)
		}
		return printJSON(quota)
	},
}

// Create Quota
var createQuotaCmd = &cobra.Command{
	Use:   "create-quota",
	Short: "Create organization quota",
	Long:  `Create a quota for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		quota, err := client.CreateQuota(orgName, limitBytes)
		if err != nil {
			return fmt.Errorf("creating quota: %w", err)
		}
		return printJSON(quota)
	},
}

// Update Quota
var updateQuotaCmd = &cobra.Command{
	Use:   "update-quota",
	Short: "Update organization quota",
	Long:  `Update the quota for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		quota, err := client.UpdateQuota(orgName, limitBytes)
		if err != nil {
			return fmt.Errorf("updating quota: %w", err)
		}
		return printJSON(quota)
	},
}

// Delete Quota
var deleteQuotaCmd = &cobra.Command{
	Use:   "delete-quota",
	Short: "Delete organization quota",
	Long:  `Delete the quota for an organization. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete a quota")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteQuota(orgName)
		if err != nil {
			return fmt.Errorf("deleting quota: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Quota deleted successfully")
		return nil
	},
}

func initOrgQuotaFlags() {
	createQuotaCmd.Flags().Int64Var(&limitBytes, "limit-bytes", 0, "Quota limit in bytes")
	_ = createQuotaCmd.MarkFlagRequired("limit-bytes")

	updateQuotaCmd.Flags().Int64Var(&limitBytes, "limit-bytes", 0, "Quota limit in bytes")
	_ = updateQuotaCmd.MarkFlagRequired("limit-bytes")

	deleteQuotaCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func init() {
	organizationCmd.AddCommand(orgQuotaCmd)
	organizationCmd.AddCommand(createQuotaCmd)
	organizationCmd.AddCommand(updateQuotaCmd)
	organizationCmd.AddCommand(deleteQuotaCmd)

	initOrgQuotaFlags()
}
