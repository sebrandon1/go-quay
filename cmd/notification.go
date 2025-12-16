package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	notificationNamespace  string
	notificationRepository string
	notificationUUID       string
	notificationEvent      string
	notificationMethod     string
	notificationTitle      string
	notificationURL        string
	confirmNotificationDel bool
)

// notificationCmd represents the notification command group
var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Repository notification/webhook management commands",
	Long: `Commands for managing repository notifications (webhooks).

Repository notifications allow external services to be notified
of events such as image pushes, builds, and vulnerability scans.

Supported events:
  - repo_push: Image push to repository
  - build_queued: Build has been queued
  - build_start: Build has started
  - build_success: Build completed successfully
  - build_failure: Build failed
  - build_canceled: Build was canceled
  - vulnerability_found: New vulnerability discovered

Supported methods:
  - webhook: HTTP webhook
  - email: Email notification
  - slack: Slack notification

Available commands:
  list   - List notifications for a repository
  info   - Get notification details
  create - Create a new notification
  delete - Delete a notification
  test   - Test a notification
  reset  - Reset notification failure count`,
}

// Notification List
var notificationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List notifications for a repository",
	Long:  `List all notifications (webhooks) for the specified repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		notifications, err := client.GetNotifications(notificationNamespace, notificationRepository)
		if err != nil {
			fmt.Printf("Error getting notifications: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Notifications for %s/%s:\n", notificationNamespace, notificationRepository)
		printJSON(notifications)
	},
}

// Notification Info
var notificationInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get notification details",
	Long:  `Get detailed information about a specific notification.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		notification, err := client.GetNotification(notificationNamespace, notificationRepository, notificationUUID)
		if err != nil {
			fmt.Printf("Error getting notification: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Notification: %s\n", notificationUUID)
		printJSON(notification)
	},
}

// Notification Create
var notificationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new notification",
	Long: `Create a new notification (webhook) for the repository.

For webhook method, provide the --url flag with the webhook endpoint.
For email method, provide the --url flag with the email address.
For slack method, provide the --url flag with the Slack webhook URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		config := map[string]interface{}{}
		switch notificationMethod {
		case "webhook":
			config["url"] = notificationURL
		case "email":
			config["email"] = notificationURL
		case "slack":
			config["url"] = notificationURL
		}

		notificationReq := &lib.CreateNotificationRequest{
			Event:  notificationEvent,
			Method: notificationMethod,
			Config: config,
			Title:  notificationTitle,
		}

		notification, err := client.CreateNotification(notificationNamespace, notificationRepository, notificationReq)
		if err != nil {
			fmt.Printf("Error creating notification: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Notification created successfully!\n")
		printJSON(notification)
	},
}

// Notification Delete
var notificationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a notification",
	Long:  `Delete a notification from the repository. This action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmNotificationDel {
			fmt.Printf("Are you sure you want to delete notification '%s'? This action cannot be undone.\n", notificationUUID)
			fmt.Println("Use --confirm to proceed with deletion.")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.DeleteNotification(notificationNamespace, notificationRepository, notificationUUID)
		if err != nil {
			fmt.Printf("Error deleting notification: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Notification '%s' deleted successfully.\n", notificationUUID)
	},
}

// Notification Test
var notificationTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test a notification",
	Long:  `Send a test event to the notification endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.TestNotification(notificationNamespace, notificationRepository, notificationUUID)
		if err != nil {
			fmt.Printf("Error testing notification: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Test event sent to notification '%s'.\n", notificationUUID)
	},
}

// Notification Reset
var notificationResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset notification failure count",
	Long:  `Reset the failure count for a notification that has failed.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.ResetNotification(notificationNamespace, notificationRepository, notificationUUID)
		if err != nil {
			fmt.Printf("Error resetting notification: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Notification '%s' failure count reset.\n", notificationUUID)
	},
}

func init() {
	// Add subcommands to notification command
	notificationCmd.AddCommand(notificationListCmd)
	notificationCmd.AddCommand(notificationInfoCmd)
	notificationCmd.AddCommand(notificationCreateCmd)
	notificationCmd.AddCommand(notificationDeleteCmd)
	notificationCmd.AddCommand(notificationTestCmd)
	notificationCmd.AddCommand(notificationResetCmd)

	// Global notification flags
	notificationCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")
	notificationCmd.PersistentFlags().StringVarP(&notificationNamespace, "namespace", "n", "", "Repository namespace")
	notificationCmd.PersistentFlags().StringVarP(&notificationRepository, "repository", "r", "", "Repository name")
	markNotificationFlagRequired(notificationCmd.MarkPersistentFlagRequired("token"))
	markNotificationFlagRequired(notificationCmd.MarkPersistentFlagRequired("namespace"))
	markNotificationFlagRequired(notificationCmd.MarkPersistentFlagRequired("repository"))

	initNotificationInfoFlags()
	initNotificationCreateFlags()
	initNotificationDeleteFlags()
	initNotificationTestFlags()
	initNotificationResetFlags()
}

func initNotificationInfoFlags() {
	notificationInfoCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	markNotificationFlagRequired(notificationInfoCmd.MarkFlagRequired("uuid"))
}

func initNotificationCreateFlags() {
	notificationCreateCmd.Flags().StringVarP(&notificationEvent, "event", "e", "", "Event type (repo_push, build_success, etc.)")
	notificationCreateCmd.Flags().StringVarP(&notificationMethod, "method", "m", "", "Method (webhook, email, slack)")
	notificationCreateCmd.Flags().StringVar(&notificationURL, "url", "", "Webhook URL or email address")
	notificationCreateCmd.Flags().StringVar(&notificationTitle, "title", "", "Notification title")
	markNotificationFlagRequired(notificationCreateCmd.MarkFlagRequired("event"))
	markNotificationFlagRequired(notificationCreateCmd.MarkFlagRequired("method"))
	markNotificationFlagRequired(notificationCreateCmd.MarkFlagRequired("url"))
}

func initNotificationDeleteFlags() {
	notificationDeleteCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	notificationDeleteCmd.Flags().BoolVar(&confirmNotificationDel, "confirm", false, "Confirm notification deletion")
	markNotificationFlagRequired(notificationDeleteCmd.MarkFlagRequired("uuid"))
}

func initNotificationTestFlags() {
	notificationTestCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	markNotificationFlagRequired(notificationTestCmd.MarkFlagRequired("uuid"))
}

func initNotificationResetFlags() {
	notificationResetCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	markNotificationFlagRequired(notificationResetCmd.MarkFlagRequired("uuid"))
}

func markNotificationFlagRequired(err error) {
	if err != nil {
		fmt.Printf("Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
