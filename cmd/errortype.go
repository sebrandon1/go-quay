/*
Package cmd provides the command-line interface for go-quay.

This file contains the commands for the Error API:
  - go-quay get error - Get details about a specific error type
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	errorTypeName string
)

// errorTypeCmd represents the error command
var errorTypeCmd = &cobra.Command{
	Use:   "error",
	Short: "Get error type information",
	Long: `Get detailed information about a specific error type.

This endpoint provides details about error types that can be returned by the Quay.io API.`,
	Run: func(_ *cobra.Command, _ []string) {
		if errorTypeName == "" {
			fmt.Println("Error: --type is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		errType, err := client.GetErrorType(errorTypeName)
		if err != nil {
			fmt.Println("Error getting error type:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(errType, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

func init() {
	errorTypeCmd.Flags().StringVarP(&token, "token", "t", "", "Quay.io API token")
	errorTypeCmd.Flags().StringVar(&errorTypeName, "type", "", "Error type to look up (e.g., 'invalid_token')")
}
