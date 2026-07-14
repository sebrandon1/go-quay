package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command group
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User account management commands",
	Long: `Commands for managing user account information and starred repositories.

Available commands:
  info        - Get current user information
  starred     - List starred repositories
  star        - Star a repository
  unstar      - Unstar a repository
  lookup      - Look up a user by username
  marketplace - Get user marketplace information`,
}

// User Info
var userInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get current user information",
	Long:  `Get detailed information about the currently authenticated user account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		user, err := client.GetUser()
		if err != nil {
			return fmt.Errorf("getting user information: %w", err)
		}

		fmt.Printf("User information for %s:\n", user.Username)
		return printJSON(user)
	},
}

// User Starred Repositories
var userStarredCmd = &cobra.Command{
	Use:   "starred",
	Short: "List starred repositories",
	Long:  `List all repositories starred by the current user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		starred, err := client.GetStarredRepositories()
		if err != nil {
			return fmt.Errorf("getting starred repositories: %w", err)
		}

		fmt.Println("Starred repositories:")
		return printJSON(starred)
	},
}

// Star Repository
var starRepoCmd = &cobra.Command{
	Use:   "star",
	Short: "Star a repository",
	Long:  `Add a repository to your starred list for easy discovery.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.StarRepository(namespace, repository)
		if err != nil {
			return fmt.Errorf("starring repository: %w", err)
		}

		fmt.Printf("Successfully starred repository %s/%s\n", namespace, repository)
		return nil
	},
}

// Unstar Repository
var unstarRepoCmd = &cobra.Command{
	Use:   "unstar",
	Short: "Unstar a repository",
	Long:  `Remove a repository from your starred list.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.UnstarRepository(namespace, repository)
		if err != nil {
			return fmt.Errorf("unstarring repository: %w", err)
		}

		fmt.Printf("Successfully unstarred repository %s/%s\n", namespace, repository)
		return nil
	},
}

var lookupUsername string

var userLookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Look up a user by username",
	Long:  `Get information about a specific user by their username.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		user, err := client.GetUserByUsername(lookupUsername)
		if err != nil {
			return fmt.Errorf("looking up user: %w", err)
		}

		return printJSON(user)
	},
}

var userMarketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Get user marketplace information",
	Long:  `Get marketplace subscription information for the current user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		marketplace, err := client.GetUserMarketplace()
		if err != nil {
			return fmt.Errorf("getting user marketplace info: %w", err)
		}

		return printJSON(marketplace)
	},
}

func init() {
	// Add subcommands to user command
	userCmd.AddCommand(userInfoCmd)
	userCmd.AddCommand(userStarredCmd)
	userCmd.AddCommand(starRepoCmd)
	userCmd.AddCommand(unstarRepoCmd)
	userCmd.AddCommand(userLookupCmd)
	userCmd.AddCommand(userMarketplaceCmd)
	// Star command specific flags (requires repository context)
	starRepoCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	starRepoCmd.Flags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	_ = starRepoCmd.MarkFlagRequired("namespace")
	_ = starRepoCmd.MarkFlagRequired("repository")

	// Unstar command specific flags (requires repository context)
	unstarRepoCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	unstarRepoCmd.Flags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	_ = unstarRepoCmd.MarkFlagRequired("namespace")
	_ = unstarRepoCmd.MarkFlagRequired("repository")

	// Lookup command flags
	userLookupCmd.Flags().StringVarP(&lookupUsername, "username", "u", "", "Username to look up")
	_ = userLookupCmd.MarkFlagRequired("username")
}
