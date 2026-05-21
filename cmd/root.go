package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

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
}

func init() {
	getCmd.PersistentFlags().StringVarP(&token, "token", "t", os.Getenv("QUAY_TOKEN"), "Quay.io API token (default: $QUAY_TOKEN)")
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
