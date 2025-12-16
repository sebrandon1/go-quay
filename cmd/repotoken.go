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
	"encoding/json"
	"fmt"
	"os"

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
	Use:   "list",
	Short: "List all repository tokens",
	Long:  `List all tokens for a repository. (DEPRECATED - use robot accounts)`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		tokens, err := client.GetRepoTokens(repoTokenNamespace, repoTokenRepository) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			fmt.Println("Error getting tokens:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(tokens, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// repotokenInfoCmd gets a specific token
var repotokenInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get a specific token",
	Long:  `Get detailed information about a specific repository token. (DEPRECATED - use robot accounts)`,
	Run: func(_ *cobra.Command, _ []string) {
		if repoTokenCode == "" {
			fmt.Println("Error: --code is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		repoToken, err := client.GetRepoToken(repoTokenNamespace, repoTokenRepository, repoTokenCode) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			fmt.Println("Error getting token:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(repoToken, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// repotokenCreateCmd creates a new token
var repotokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new token",
	Long:  `Create a new repository token. (DEPRECATED - use robot accounts)`,
	Run: func(_ *cobra.Command, _ []string) {
		if repoTokenName == "" {
			fmt.Println("Error: --name is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		createReq := &lib.CreateRepoTokenRequest{
			FriendlyName: repoTokenName,
		}

		repoToken, err := client.CreateRepoToken(repoTokenNamespace, repoTokenRepository, createReq) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			fmt.Println("Error creating token:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(repoToken, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// repotokenUpdateCmd updates a token
var repotokenUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a token",
	Long:  `Update a repository token's role. (DEPRECATED - use robot accounts)`,
	Run: func(_ *cobra.Command, _ []string) {
		if repoTokenCode == "" {
			fmt.Println("Error: --code is required")
			os.Exit(1)
		}
		if repoTokenRole == "" {
			fmt.Println("Error: --role is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		updateReq := &lib.UpdateRepoTokenRequest{
			Role: repoTokenRole,
		}

		repoToken, err := client.UpdateRepoToken(repoTokenNamespace, repoTokenRepository, repoTokenCode, updateReq) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			fmt.Println("Error updating token:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(repoToken, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// repotokenDeleteCmd deletes a token
var repotokenDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a token",
	Long:  `Delete a repository token. (DEPRECATED - use robot accounts)`,
	Run: func(_ *cobra.Command, _ []string) {
		if repoTokenCode == "" {
			fmt.Println("Error: --code is required")
			os.Exit(1)
		}
		if !confirmTokenDelete {
			fmt.Println("Error: --confirm is required to delete a token")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		err = client.DeleteRepoToken(repoTokenNamespace, repoTokenRepository, repoTokenCode) //nolint:staticcheck // Intentionally using deprecated API
		if err != nil {
			fmt.Println("Error deleting token:", err)
			os.Exit(1)
		}

		fmt.Printf("Token %s deleted successfully\n", repoTokenCode)
	},
}

func setupRepoTokenFlags() {
	// Common flags
	for _, cmd := range []*cobra.Command{repotokenListCmd, repotokenInfoCmd, repotokenCreateCmd, repotokenUpdateCmd, repotokenDeleteCmd} {
		cmd.Flags().StringVarP(&repoTokenNamespace, "namespace", "n", "", "Namespace/organization")
		cmd.Flags().StringVarP(&repoTokenRepository, "repository", "r", "", "Repository name")
		cmd.Flags().StringVarP(&token, "token", "t", "", "Quay.io API token")
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
