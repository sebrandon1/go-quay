package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "quay",
	Short: "Quay CLI interacts with Quay.io API",
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get objects from Quay.io",
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(aggregatedLogsCmd)
	getCmd.AddCommand(repositoryCmd)
	getCmd.AddCommand(billingCmd)
	getCmd.AddCommand(organizationCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
