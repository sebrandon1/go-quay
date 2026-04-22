package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	messageContent   string
	messageSeverity  string
	messageMediaType string
)

var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "System message management commands",
}

var messagesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get system messages",
	Long:  `Get system-wide messages including maintenance notifications and announcements.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		messages, err := client.GetMessages()
		if err != nil {
			fmt.Printf("Error getting messages: %v\n", err)
			os.Exit(1)
		}

		printJSON(messages)
	},
}

var messagesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a system message",
	Long:  `Create a new system-wide message.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		message, err := client.CreateMessage(messageContent, messageSeverity, messageMediaType)
		if err != nil {
			fmt.Printf("Error creating message: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Message created successfully!")
		printJSON(message)
	},
}

func init() {
	messagesCmd.AddCommand(messagesListCmd)
	messagesCmd.AddCommand(messagesCreateCmd)

	messagesCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Quay.io API token")

	messagesCreateCmd.Flags().StringVar(&messageContent, "content", "", "Message content")
	messagesCreateCmd.Flags().StringVar(&messageSeverity, "severity", "", "Message severity (info, warning, error)")
	messagesCreateCmd.Flags().StringVar(&messageMediaType, "media-type", "", "Media type for the message")
	if err := messagesCreateCmd.MarkFlagRequired("content"); err != nil {
		fmt.Printf("Error marking content flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := messagesCreateCmd.MarkFlagRequired("severity"); err != nil {
		fmt.Printf("Error marking severity flag as required: %v\n", err)
		os.Exit(1)
	}
}
