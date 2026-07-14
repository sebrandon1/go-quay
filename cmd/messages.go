package cmd

import (
	"fmt"

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
	Use:   subcmdList,
	Short: "Get system messages",
	Long:  `Get system-wide messages including maintenance notifications and announcements.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		messages, err := client.GetMessages()
		if err != nil {
			return fmt.Errorf("getting messages: %w", err)
		}

		return printJSON(messages)
	},
}

var messagesCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
	Short: "Create a system message",
	Long:  `Create a new system-wide message.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		message, err := client.CreateMessage(messageContent, messageSeverity, messageMediaType)
		if err != nil {
			return fmt.Errorf("creating message: %w", err)
		}

		fmt.Println("Message created successfully!")
		return printJSON(message)
	},
}

func init() {
	messagesCmd.AddCommand(messagesListCmd)
	messagesCmd.AddCommand(messagesCreateCmd)

	messagesCreateCmd.Flags().StringVar(&messageContent, "content", "", "Message content")
	messagesCreateCmd.Flags().StringVar(&messageSeverity, "severity", "", "Message severity (info, warning, error)")
	messagesCreateCmd.Flags().StringVar(&messageMediaType, "media-type", "", "Media type for the message")
	_ = messagesCreateCmd.MarkFlagRequired("content")
	_ = messagesCreateCmd.MarkFlagRequired("severity")
}
