package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Get Proxy Cache Config
var proxyCacheCmd = &cobra.Command{
	Use:   "proxy-cache",
	Short: "Get proxy cache configuration",
	Long:  `Get proxy cache configuration for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		config, err := client.GetProxyCacheConfig(orgName)
		if err != nil {
			return fmt.Errorf("getting proxy cache config: %w", err)
		}
		return printJSON(config)
	},
}

// Create Proxy Cache Config
var createProxyCacheCmd = &cobra.Command{
	Use:   "create-proxy-cache",
	Short: "Create proxy cache configuration",
	Long:  `Create proxy cache configuration for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		config, err := client.CreateProxyCacheConfig(orgName, upstreamRegistry, insecure, expiration)
		if err != nil {
			return fmt.Errorf("creating proxy cache config: %w", err)
		}
		return printJSON(config)
	},
}

// Delete Proxy Cache Config
var deleteProxyCacheCmd = &cobra.Command{
	Use:   "delete-proxy-cache",
	Short: "Delete proxy cache configuration",
	Long:  `Delete proxy cache configuration for an organization. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete proxy cache config")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteProxyCacheConfig(orgName)
		if err != nil {
			return fmt.Errorf("deleting proxy cache config: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Proxy cache config deleted successfully")
		return nil
	},
}

func initOrgProxyCacheFlags() {
	createProxyCacheCmd.Flags().StringVar(&upstreamRegistry, "upstream-registry", "", "Upstream registry URL")
	_ = createProxyCacheCmd.MarkFlagRequired("upstream-registry")
	createProxyCacheCmd.Flags().BoolVar(&insecure, "insecure", false, "Allow insecure connections")
	createProxyCacheCmd.Flags().IntVar(&expiration, "expiration", 0, "Cache expiration in seconds")

	deleteProxyCacheCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func init() {
	organizationCmd.AddCommand(proxyCacheCmd)
	organizationCmd.AddCommand(createProxyCacheCmd)
	organizationCmd.AddCommand(deleteProxyCacheCmd)

	initOrgProxyCacheFlags()
}
