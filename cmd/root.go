package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var quayURL string

var rootCmd = &cobra.Command{
	Use:   cliName,
	Short: "Quay CLI interacts with Quay.io API",
}

func SetVersion(v string) {
	rootCmd.Version = v
}

var getCmd = &cobra.Command{
	Use:   cmdGet,
	Short: "Get objects from Quay.io",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate token is provided
		if token == "" {
			return fmt.Errorf(`authentication token required

Set QUAY_TOKEN environment variable, use --token/-t flag, or add to config file (%s).
Get your token at https://quay.io/organization/<org>?tab=applications`, configFilePath())
		}

		// Validate output format
		switch outputFormat {
		case outputJSON, outputYAML, outputTable:
			// valid
		default:
			return fmt.Errorf("invalid output format %q: must be json, yaml, or table", outputFormat)
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

// firstNonEmpty returns the first non-empty string from the given values.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func init() {
	getCmd.PersistentFlags().StringVarP(&token, "token", "t", envOrDefault("QUAY_TOKEN", appCfg.Token), "Quay.io API token (default: $QUAY_TOKEN or config file)")
	getCmd.PersistentFlags().StringVar(&quayURL, "quay-url", envOrDefault("QUAY_URL", firstNonEmpty(appCfg.QuayURL, lib.DefaultQuayURL)), "Quay API base URL (default: $QUAY_URL or config file)")
	getCmd.PersistentFlags().StringVarP(&outputFormat, "output", "O", "json", "Output format: json, yaml, or table")
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
	getCmd.AddCommand(mirrorCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
