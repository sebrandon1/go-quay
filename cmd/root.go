package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var quayURL string

var rootCmd = &cobra.Command{
	Use:   "quay",
	Short: "Quay CLI interacts with Quay.io API",
}

func SetVersion(v string) {
	rootCmd.Version = v
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get objects from Quay.io",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate token is provided
		if token == "" {
			return fmt.Errorf(`authentication token required

Set QUAY_TOKEN environment variable or use --token/-t flag.
Get your token at https://quay.io/organization/<org>?tab=applications`)
		}
		return nil
	},
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func init() {
	getCmd.PersistentFlags().StringVarP(&token, "token", "t", envOrDefault("QUAY_TOKEN", ""), "Quay.io API token (default: $QUAY_TOKEN)")
	getCmd.PersistentFlags().StringVar(&quayURL, "quay-url", envOrDefault("QUAY_URL", lib.DefaultQuayURL), "Quay API base URL (default: $QUAY_URL)")
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(repositoryCmd)
	getCmd.AddCommand(billingCmd)
	getCmd.AddCommand(organizationCmd)
	getCmd.AddCommand(permissionsCmd)
	getCmd.AddCommand(tagCmd)
	getCmd.AddCommand(userCmd)
	getCmd.AddCommand(manifestCmd)
	getCmd.AddCommand(secscanCmd)
	getCmd.AddCommand(robotCmd)
	getCmd.AddCommand(searchCmd)
	getCmd.AddCommand(teamCmd)
	getCmd.AddCommand(buildCmd)
	getCmd.AddCommand(notificationCmd)
	getCmd.AddCommand(triggerCmd)
	getCmd.AddCommand(discoveryCmd)
	getCmd.AddCommand(errorTypeCmd)
	getCmd.AddCommand(messagesCmd)
	getCmd.AddCommand(prototypeCmd)
	getCmd.AddCommand(repotokenCmd)
	getCmd.AddCommand(logsCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
