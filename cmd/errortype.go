/*
Package cmd provides the command-line interface for go-quay.

This file contains the commands for the Error API:
  - go-quay get error - Get details about a specific error type
*/
package cmd

import (
	"fmt"

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
	RunE: func(_ *cobra.Command, _ []string) error {
		if errorTypeName == "" {
			return fmt.Errorf("--type is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		errType, err := client.GetErrorType(errorTypeName)
		if err != nil {
			return fmt.Errorf("getting error type: %w", err)
		}

		return printJSON(errType)
	},
}

func init() {
	errorTypeCmd.Flags().StringVar(&errorTypeName, "type", "", "Error type to look up (e.g., 'invalid_token')")
}
