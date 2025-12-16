/*
Package cmd provides the command-line interface for go-quay.

This file contains the commands for the Discovery API:
  - go-quay get discovery - Get API discovery information
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

// discoveryCmd represents the discovery command
var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Get API discovery information",
	Long: `Get API discovery information from Quay.io.

This endpoint returns information about available API endpoints and versions.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		discovery, err := client.GetDiscovery()
		if err != nil {
			fmt.Println("Error getting discovery:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(discovery, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

func init() {
	discoveryCmd.Flags().StringVarP(&token, "token", "t", "", "Quay.io API token")
}
