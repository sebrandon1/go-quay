/*
Package cmd provides the command-line interface for go-quay.

This file contains the commands for the RepoToken API:
  - go-quay get repotoken list     - List all repository tokens
  - go-quay get repotoken info     - Get a specific token
  - go-quay get repotoken create   - Create a new token
  - go-quay get repotoken update   - Update a token
  - go-quay get repotoken delete   - Delete a token

WARNING: Repository tokens are deprecated. Use robot accounts instead.
*/
package cmd

import (
	"fmt"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	repoTokenNamespace  string
	repoTokenRepository string
	repoTokenCode       string
	repoTokenName       string
	repoTokenRole       string
	confirmTokenDelete  bool
)

// repotokenCmd represents the repotoken command
var repotokenCmd = &cobra.Command{
	Use:   "repotoken",
	Short: "Manage repository tokens (DEPRECATED)",
	Long: `Manage repository tokens for authentication.

WARNING: Repository tokens are deprecated. Use robot accounts instead for
better security, auditing, and permission management.`,
}

// repotokenListCmd lists all tokens
var repotokenListCmd = &cobra.Command{
	Use:   subcmdList,
	Short: "List all repository tokens",
	Long:  `List all tokens for a repository. (DEPRECATED - use robot accounts)`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		tokens, err := client.GetRepoTokens(repoTokenNamespace, repoTokenRepository) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			return fmt.Errorf("getting tokens: %w", err)
		}

		return printJSON(tokens)
	},
}

// repotokenInfoCmd gets a specific token
var repotokenInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get a specific token",
	Long:  `Get detailed information about a specific repository token. (DEPRECATED - use robot accounts)`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if repoTokenCode == "" {
			return fmt.Errorf("--code is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		repoToken, err := client.GetRepoToken(repoTokenNamespace, repoTokenRepository, repoTokenCode) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			return fmt.Errorf("getting token: %w", err)
		}

		return printJSON(repoToken)
	},
}

// repotokenCreateCmd creates a new token
var repotokenCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
	Short: "Create a new token",
	Long:  `Create a new repository token. (DEPRECATED - use robot accounts)`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if repoTokenName == "" {
			return fmt.Errorf("--name is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		createReq := &lib.CreateRepoTokenRequest{
			FriendlyName: repoTokenName,
		}

		repoToken, err := client.CreateRepoToken(repoTokenNamespace, repoTokenRepository, createReq) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			return fmt.Errorf("creating token: %w", err)
		}

		return printJSON(repoToken)
	},
}

// repotokenUpdateCmd updates a token
var repotokenUpdateCmd = &cobra.Command{
	Use:   subcmdUpdate,
	Short: "Update a token",
	Long:  `Update a repository token's role. (DEPRECATED - use robot accounts)`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if repoTokenCode == "" {
			return fmt.Errorf("--code is required")
		}
		if repoTokenRole == "" {
			return fmt.Errorf("--role is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		updateReq := &lib.UpdateRepoTokenRequest{
			Role: repoTokenRole,
		}

		repoToken, err := client.UpdateRepoToken(repoTokenNamespace, repoTokenRepository, repoTokenCode, updateReq) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			return fmt.Errorf("updating token: %w", err)
		}

		return printJSON(repoToken)
	},
}

// repotokenDeleteCmd deletes a token
var repotokenDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a token",
	Long:  `Delete a repository token. (DEPRECATED - use robot accounts)`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if repoTokenCode == "" {
			return fmt.Errorf("--code is required")
		}
		if !confirmTokenDelete {
			return fmt.Errorf("--confirm is required to delete a token")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteRepoToken(repoTokenNamespace, repoTokenRepository, repoTokenCode) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			return fmt.Errorf("deleting token: %w", err)
		}

		fmt.Printf("Token %s deleted successfully\n", repoTokenCode)
		return nil
	},
}

func setupRepoTokenFlags() {
	// Common flags
	for _, cmd := range []*cobra.Command{repotokenListCmd, repotokenInfoCmd, repotokenCreateCmd, repotokenUpdateCmd, repotokenDeleteCmd} {
		cmd.Flags().StringVarP(&repoTokenNamespace, "namespace", "n", "", "Namespace/organization")
		cmd.Flags().StringVarP(&repoTokenRepository, "repository", "r", "", "Repository name")
	}

	// Code flags
	for _, cmd := range []*cobra.Command{repotokenInfoCmd, repotokenUpdateCmd, repotokenDeleteCmd} {
		cmd.Flags().StringVar(&repoTokenCode, "code", "", "Token code")
	}

	// Create flags
	repotokenCreateCmd.Flags().StringVar(&repoTokenName, "name", "", "Friendly name for the token")

	// Update flags
	repotokenUpdateCmd.Flags().StringVar(&repoTokenRole, "role", "", "New role (read, write, admin)")

	// Delete flags
	repotokenDeleteCmd.Flags().BoolVar(&confirmTokenDelete, "confirm", false, "Confirm deletion")
}

func init() {
	repotokenCmd.AddCommand(repotokenListCmd)
	repotokenCmd.AddCommand(repotokenInfoCmd)
	repotokenCmd.AddCommand(repotokenCreateCmd)
	repotokenCmd.AddCommand(repotokenUpdateCmd)
	repotokenCmd.AddCommand(repotokenDeleteCmd)

	setupRepoTokenFlags()
}
