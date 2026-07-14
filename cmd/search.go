package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		result, err := client.SearchRepositories(searchQuery, searchPage)
		if err != nil {
			return fmt.Errorf("searching repositories: %w", err)
		}

		fmt.Printf("Repository search results for: %s\n", searchQuery)
		return printJSON(result)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		result, err := client.SearchAll(searchQuery)
		if err != nil {
			return fmt.Errorf("searching: %w", err)
		}

		fmt.Printf("Search results for: %s\n", searchQuery)
		return printJSON(result)
	},
}

func init() {
	// Add subcommands to search command
	searchCmd.AddCommand(searchReposCmd)
	searchCmd.AddCommand(searchAllCmd)

	// Global search flags
	searchCmd.PersistentFlags().StringVarP(&searchQuery, "query", "q", "", "Search query")
	_ = searchCmd.MarkPersistentFlagRequired("query")

	// Repositories command specific flags
	searchReposCmd.Flags().IntVarP(&searchPage, "page", "p", 0, "Page number for pagination (starts at 1)")
}
