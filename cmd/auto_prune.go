package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Auto-prune Policies
var autoPruneCmd = &cobra.Command{
	Use:   "auto-prune",
	Short: "Get auto-prune policies",
	Long:  `Get list of auto-prune policies for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		policies, err := client.GetAutoPrunePolicies(orgName)
		if err != nil {
			return fmt.Errorf("getting auto-prune policies: %w", err)
		}
		return printJSON(policies)
	},
}

// Get Auto-Prune Policy
var autoPrunePolicyCmd = &cobra.Command{
	Use:   "auto-prune-policy",
	Short: "Get a specific auto-prune policy",
	Long:  `Get detailed information about a specific auto-prune policy.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		policy, err := client.GetAutoPrunePolicy(orgName, policyUUID)
		if err != nil {
			return fmt.Errorf("getting auto-prune policy: %w", err)
		}
		return printJSON(policy)
	},
}

// Create Auto-Prune Policy
var createAutoPruneCmd = &cobra.Command{
	Use:   "create-auto-prune",
	Short: "Create an auto-prune policy",
	Long:  `Create an auto-prune policy for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		policy, err := client.CreateAutoPrunePolicy(orgName, method, pruneValue, tagPattern)
		if err != nil {
			return fmt.Errorf("creating auto-prune policy: %w", err)
		}
		return printJSON(policy)
	},
}

// Update Auto-Prune Policy
var updateAutoPruneCmd = &cobra.Command{
	Use:   "update-auto-prune",
	Short: "Update an auto-prune policy",
	Long:  `Update an existing auto-prune policy.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		policy, err := client.UpdateAutoPrunePolicy(orgName, policyUUID, method, pruneValue, tagPattern)
		if err != nil {
			return fmt.Errorf("updating auto-prune policy: %w", err)
		}
		return printJSON(policy)
	},
}

// Delete Auto-Prune Policy
var deleteAutoPruneCmd = &cobra.Command{
	Use:   "delete-auto-prune",
	Short: "Delete an auto-prune policy",
	Long:  `Delete an auto-prune policy. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete an auto-prune policy")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteAutoPrunePolicy(orgName, policyUUID)
		if err != nil {
			return fmt.Errorf("deleting auto-prune policy: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Auto-prune policy deleted successfully")
		return nil
	},
}

func initOrgAutoPruneFlags() {
	autoPrunePolicyCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	_ = autoPrunePolicyCmd.MarkFlagRequired("policy-uuid")

	createAutoPruneCmd.Flags().StringVar(&method, "method", "", "Prune method")
	_ = createAutoPruneCmd.MarkFlagRequired("method")
	createAutoPruneCmd.Flags().IntVar(&pruneValue, "value", 0, "Prune value")
	_ = createAutoPruneCmd.MarkFlagRequired("value")
	createAutoPruneCmd.Flags().StringVar(&tagPattern, "tag-pattern", "", "Tag pattern to match")

	updateAutoPruneCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	_ = updateAutoPruneCmd.MarkFlagRequired("policy-uuid")
	updateAutoPruneCmd.Flags().StringVar(&method, "method", "", "Prune method")
	_ = updateAutoPruneCmd.MarkFlagRequired("method")
	updateAutoPruneCmd.Flags().IntVar(&pruneValue, "value", 0, "Prune value")
	_ = updateAutoPruneCmd.MarkFlagRequired("value")
	updateAutoPruneCmd.Flags().StringVar(&tagPattern, "tag-pattern", "", "Tag pattern to match")

	deleteAutoPruneCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	_ = deleteAutoPruneCmd.MarkFlagRequired("policy-uuid")
	deleteAutoPruneCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func init() {
	organizationCmd.AddCommand(autoPruneCmd)
	organizationCmd.AddCommand(autoPrunePolicyCmd)
	organizationCmd.AddCommand(createAutoPruneCmd)
	organizationCmd.AddCommand(updateAutoPruneCmd)
	organizationCmd.AddCommand(deleteAutoPruneCmd)

	initOrgAutoPruneFlags()
}
