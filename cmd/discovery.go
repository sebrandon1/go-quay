package cmd

import (
	"fmt"

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
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		discovery, err := client.GetDiscovery()
		if err != nil {
			return fmt.Errorf("getting discovery: %w", err)
		}

		return printJSON(discovery)
	},
}

var discoveryCapabilitiesCmd = &cobra.Command{
	Use:   "capabilities",
	Short: "Get registry capabilities",
	Long:  `Get registry capabilities including sparse manifest support and available mirror architectures.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		capabilities, err := client.GetRegistryCapabilities()
		if err != nil {
			return fmt.Errorf("getting registry capabilities: %w", err)
		}

		return printJSON(capabilities)
	},
}

var discoveryAppInfoCmd = &cobra.Command{
	Use:   "app-info",
	Short: "Get application information by client ID",
	Long:  `Get detailed information about an OAuth application by its client ID.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		app, err := client.GetAppInfo(discoveryClientID)
		if err != nil {
			return fmt.Errorf("getting app info: %w", err)
		}

		return printJSON(app)
	},
}

var discoveryEntitiesCmd = &cobra.Command{
	Use:   "entities",
	Short: "Search for entities by prefix",
	Long:  `Search for users, organizations, and teams by name prefix.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		entities, err := client.GetEntities(entityPrefix, includeOrgs, includeTeams)
		if err != nil {
			return fmt.Errorf("getting entities: %w", err)
		}

		return printJSON(entities)
	},
}

func init() {
	discoveryCmd.AddCommand(discoveryAPICmd)
	discoveryCmd.AddCommand(discoveryCapabilitiesCmd)
	discoveryCmd.AddCommand(discoveryAppInfoCmd)
	discoveryCmd.AddCommand(discoveryEntitiesCmd)

	discoveryAppInfoCmd.Flags().StringVar(&discoveryClientID, "client-id", "", "OAuth application client ID")
	_ = discoveryAppInfoCmd.MarkFlagRequired("client-id")

	discoveryEntitiesCmd.Flags().StringVar(&entityPrefix, "prefix", "", "Entity name prefix to search")
	discoveryEntitiesCmd.Flags().BoolVar(&includeOrgs, "include-orgs", false, "Include organizations in results")
	discoveryEntitiesCmd.Flags().BoolVar(&includeTeams, "include-teams", false, "Include teams in results")
	_ = discoveryEntitiesCmd.MarkFlagRequired("prefix")
}
