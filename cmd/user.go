package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

// userCmd represents the user command group
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User account management commands",
	Long: `Commands for managing user account information and starred repositories.

Available commands:
  info    - Get current user information
  starred - List starred repositories
  star    - Star a repository
  unstar  - Unstar a repository`,
}

// User Info
var userInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get current user information",
	Long:  `Get detailed information about the currently authenticated user account.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		user, err := client.GetUser()
		if err != nil {
			fmt.Printf("Error getting user information: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("User information for %s:\n", user.Username)
		printJSON(user)
	},
}

// User Starred Repositories
var userStarredCmd = &cobra.Command{
	Use:   "starred",
	Short: "List starred repositories",
	Long:  `List all repositories starred by the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		starred, err := client.GetStarredRepositories()
		if err != nil {
			fmt.Printf("Error getting starred repositories: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Starred repositories:")
		printJSON(starred)
	},
}

// Star Repository
var starRepoCmd = &cobra.Command{
	Use:   "star",
	Short: "Star a repository",
	Long:  `Add a repository to your starred list for easy discovery.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.StarRepository(namespace, repository)
		if err != nil {
			fmt.Printf("Error starring repository: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully starred repository %s/%s\n", namespace, repository)
	},
}

// Unstar Repository
var unstarRepoCmd = &cobra.Command{
	Use:   "unstar",
	Short: "Unstar a repository",
	Long:  `Remove a repository from your starred list.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.UnstarRepository(namespace, repository)
		if err != nil {
			fmt.Printf("Error unstarring repository: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully unstarred repository %s/%s\n", namespace, repository)
	},
}

var lookupUsername string

var userLookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Look up a user by username",
	Long:  `Get information about a specific user by their username.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		user, err := client.GetUserByUsername(lookupUsername)
		if err != nil {
			fmt.Printf("Error looking up user: %v\n", err)
			os.Exit(1)
		}

		printJSON(user)
	},
}

var userMarketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Get user marketplace information",
	Long:  `Get marketplace subscription information for the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		marketplace, err := client.GetUserMarketplace()
		if err != nil {
			fmt.Printf("Error getting user marketplace info: %v\n", err)
			os.Exit(1)
		}

		printJSON(marketplace)
	},
}

var starUserRepoCmd = &cobra.Command{
	Use:   "star-user",
	Short: "Star a repository (user endpoint)",
	Long:  `Add a repository to your starred list using the user star endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.StarUserRepository(namespace, repository)
		if err != nil {
			fmt.Printf("Error starring repository: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully starred repository %s/%s\n", namespace, repository)
	},
}

var unstarUserRepoCmd = &cobra.Command{
	Use:   "unstar-user",
	Short: "Unstar a repository (user endpoint)",
	Long:  `Remove a repository from your starred list using the user star endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.UnstarUserRepository(namespace, repository)
		if err != nil {
			fmt.Printf("Error unstarring repository: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully unstarred repository %s/%s\n", namespace, repository)
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
	userCmd.AddCommand(starUserRepoCmd)
	userCmd.AddCommand(unstarUserRepoCmd)

	// Global user flags
	userCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")

	// Mark token as required for all user commands
	if err := userCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}

	// Star command specific flags (requires repository context)
	starRepoCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	starRepoCmd.Flags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	if err := starRepoCmd.MarkFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := starRepoCmd.MarkFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}

	// Unstar command specific flags (requires repository context)
	unstarRepoCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	unstarRepoCmd.Flags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	if err := unstarRepoCmd.MarkFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := unstarRepoCmd.MarkFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}

	// Lookup command flags
	userLookupCmd.Flags().StringVarP(&lookupUsername, "username", "u", "", "Username to look up")
	if err := userLookupCmd.MarkFlagRequired("username"); err != nil {
		fmt.Printf("Error marking username flag as required: %v\n", err)
		os.Exit(1)
	}

	// Star-user command flags
	starUserRepoCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	starUserRepoCmd.Flags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	if err := starUserRepoCmd.MarkFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := starUserRepoCmd.MarkFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}

	// Unstar-user command flags
	unstarUserRepoCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	unstarUserRepoCmd.Flags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	if err := unstarUserRepoCmd.MarkFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := unstarUserRepoCmd.MarkFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}
}
