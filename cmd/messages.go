/*
Package cmd provides the command-line interface for go-quay.

This file contains the commands for the Messages API:
  - go-quay get messages - Get system messages
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

// messagesCmd represents the messages command
var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Get system messages",
	Long: `Get system-wide messages for the authenticated user.

This includes maintenance notifications, announcements, and other important messages.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		messages, err := client.GetMessages()
		if err != nil {
			fmt.Println("Error getting messages:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(messages, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

func init() {
	messagesCmd.Flags().StringVarP(&token, "token", "t", "", "Quay.io API token")
}
