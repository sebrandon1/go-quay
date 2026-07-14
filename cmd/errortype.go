/*
Package cmd provides the command-line interface for go-quay.

This file contains the commands for the Error API:
  - go-quay get error - Get details about a specific error type
*/
package cmd

import (
	"fmt"
	"os"

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

		client := mustGetClient()

		errType, err := client.GetErrorType(errorTypeName)
		if err != nil {
			fmt.Println("Error getting error type:", err)
			os.Exit(1)
		}

		printJSON(errType)
	},
}

func init() {
	errorTypeCmd.Flags().StringVar(&errorTypeName, "type", "", "Error type to look up (e.g., 'invalid_token')")
}
