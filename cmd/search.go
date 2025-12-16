package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	searchQuery string
	searchPage  int
)

// searchCmd represents the search command group
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for repositories and entities on Quay.io",
	Long: `Commands for searching repositories, users, organizations, and other entities.

Available commands:
  repositories - Search for repositories
  all          - Search all entity types`,
}

// Search Repositories
var searchReposCmd = &cobra.Command{
	Use:   "repositories",
	Short: "Search for repositories",
	Long:  `Search for repositories on Quay.io by name or description.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		result, err := client.SearchRepositories(searchQuery, searchPage)
		if err != nil {
			fmt.Printf("Error searching repositories: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Repository search results for: %s\n", searchQuery)
		printJSON(result)
	},
}

// Search All
var searchAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Search all entity types",
	Long: `Search for all entity types including repositories, users, organizations, teams, and robots.

Results include a 'kind' field indicating the entity type:
  - repository: Container image repository
  - user: User account
  - organization: Organization
  - team: Team within an organization
  - robot: Robot account`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		result, err := client.SearchAll(searchQuery)
		if err != nil {
			fmt.Printf("Error searching: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Search results for: %s\n", searchQuery)
		printJSON(result)
	},
}

func init() {
	// Add subcommands to search command
	searchCmd.AddCommand(searchReposCmd)
	searchCmd.AddCommand(searchAllCmd)

	// Global search flags
	searchCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")
	searchCmd.PersistentFlags().StringVarP(&searchQuery, "query", "q", "", "Search query")
	if err := searchCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := searchCmd.MarkPersistentFlagRequired("query"); err != nil {
		fmt.Printf("Error marking query flag as required: %v\n", err)
		os.Exit(1)
	}

	// Repositories command specific flags
	searchReposCmd.Flags().IntVarP(&searchPage, "page", "p", 0, "Page number for pagination (starts at 1)")
}
