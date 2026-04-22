package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	discoveryClientID string
	entityPrefix      string
	includeOrgs       bool
	includeTeams      bool
)

var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "API discovery and entity lookup commands",
}

var discoveryAPICmd = &cobra.Command{
	Use:   "api",
	Short: "Get API discovery information",
	Long:  `Get API discovery information from Quay.io including available endpoints and versions.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		discovery, err := client.GetDiscovery()
		if err != nil {
			fmt.Printf("Error getting discovery: %v\n", err)
			os.Exit(1)
		}

		printJSON(discovery)
	},
}

var discoveryAppInfoCmd = &cobra.Command{
	Use:   "app-info",
	Short: "Get application information by client ID",
	Long:  `Get detailed information about an OAuth application by its client ID.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		app, err := client.GetAppInfo(discoveryClientID)
		if err != nil {
			fmt.Printf("Error getting app info: %v\n", err)
			os.Exit(1)
		}

		printJSON(app)
	},
}

var discoveryEntitiesCmd = &cobra.Command{
	Use:   "entities",
	Short: "Search for entities by prefix",
	Long:  `Search for users, organizations, and teams by name prefix.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		entities, err := client.GetEntities(entityPrefix, includeOrgs, includeTeams)
		if err != nil {
			fmt.Printf("Error getting entities: %v\n", err)
			os.Exit(1)
		}

		printJSON(entities)
	},
}

func init() {
	discoveryCmd.AddCommand(discoveryAPICmd)
	discoveryCmd.AddCommand(discoveryAppInfoCmd)
	discoveryCmd.AddCommand(discoveryEntitiesCmd)

	discoveryCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Quay.io API token")

	discoveryAppInfoCmd.Flags().StringVar(&discoveryClientID, "client-id", "", "OAuth application client ID")
	if err := discoveryAppInfoCmd.MarkFlagRequired("client-id"); err != nil {
		fmt.Printf("Error marking client-id flag as required: %v\n", err)
		os.Exit(1)
	}

	discoveryEntitiesCmd.Flags().StringVar(&entityPrefix, "prefix", "", "Entity name prefix to search")
	discoveryEntitiesCmd.Flags().BoolVar(&includeOrgs, "include-orgs", false, "Include organizations in results")
	discoveryEntitiesCmd.Flags().BoolVar(&includeTeams, "include-teams", false, "Include teams in results")
	if err := discoveryEntitiesCmd.MarkFlagRequired("prefix"); err != nil {
		fmt.Printf("Error marking prefix flag as required: %v\n", err)
		os.Exit(1)
	}
}
