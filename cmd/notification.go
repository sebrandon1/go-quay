package cmd

import (
	"fmt"

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
	Use:   subcmdList,
	Short: "List notifications for a repository",
	Long:  `List all notifications (webhooks) for the specified repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		notifications, err := client.GetNotifications(notificationNamespace, notificationRepository)
		if err != nil {
			return fmt.Errorf("getting notifications: %w", err)
		}

		fmt.Printf("Notifications for %s/%s:\n", notificationNamespace, notificationRepository)
		return printJSON(notifications)
	},
}

// Notification Info
var notificationInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get notification details",
	Long:  `Get detailed information about a specific notification.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		notification, err := client.GetNotification(notificationNamespace, notificationRepository, notificationUUID)
		if err != nil {
			return fmt.Errorf("getting notification: %w", err)
		}

		fmt.Printf("Notification: %s\n", notificationUUID)
		return printJSON(notification)
	},
}

// Notification Create
var notificationCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
	Short: "Create a new notification",
	Long: `Create a new notification (webhook) for the repository.

For webhook method, provide the --url flag with the webhook endpoint.
For email method, provide the --url flag with the email address.
For slack method, provide the --url flag with the Slack webhook URL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
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
			return fmt.Errorf("creating notification: %w", err)
		}

		fmt.Printf("Notification created successfully!\n")
		return printJSON(notification)
	},
}

// Notification Delete
var notificationDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a notification",
	Long:  `Delete a notification from the repository. This action cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmNotificationDel {
			return fmt.Errorf("are you sure you want to delete notification '%s'? This action cannot be undone.\nUse --confirm to proceed with deletion", notificationUUID)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteNotification(notificationNamespace, notificationRepository, notificationUUID)
		if err != nil {
			return fmt.Errorf("deleting notification: %w", err)
		}

		fmt.Printf("Notification '%s' deleted successfully.\n", notificationUUID)
		return nil
	},
}

// Notification Test
var notificationTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test a notification",
	Long:  `Send a test event to the notification endpoint.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.TestNotification(notificationNamespace, notificationRepository, notificationUUID)
		if err != nil {
			return fmt.Errorf("testing notification: %w", err)
		}

		fmt.Printf("Test event sent to notification '%s'.\n", notificationUUID)
		return nil
	},
}

// Notification Reset
var notificationResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset notification failure count",
	Long:  `Reset the failure count for a notification that has failed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.ResetNotification(notificationNamespace, notificationRepository, notificationUUID)
		if err != nil {
			return fmt.Errorf("resetting notification: %w", err)
		}

		fmt.Printf("Notification '%s' failure count reset.\n", notificationUUID)
		return nil
	},
}

var notificationUpdateCmd = &cobra.Command{
	Use:   subcmdUpdate,
	Short: "Update a notification",
	Long:  `Update an existing notification (webhook) configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
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

		notification, err := client.UpdateNotification(notificationNamespace, notificationRepository, notificationUUID, notificationReq)
		if err != nil {
			return fmt.Errorf("updating notification: %w", err)
		}

		fmt.Printf("Notification '%s' updated successfully!\n", notificationUUID)
		return printJSON(notification)
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
	notificationCmd.AddCommand(notificationUpdateCmd)

	// Global notification flags
	notificationCmd.PersistentFlags().StringVarP(&notificationNamespace, "namespace", "n", "", "Repository namespace")
	notificationCmd.PersistentFlags().StringVarP(&notificationRepository, "repository", "r", "", "Repository name")
	_ = notificationCmd.MarkPersistentFlagRequired("namespace")
	_ = notificationCmd.MarkPersistentFlagRequired("repository")

	initNotificationInfoFlags()
	initNotificationCreateFlags()
	initNotificationDeleteFlags()
	initNotificationTestFlags()
	initNotificationResetFlags()
	initNotificationUpdateFlags()
}

func initNotificationUpdateFlags() {
	notificationUpdateCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	notificationUpdateCmd.Flags().StringVarP(&notificationEvent, "event", "e", "", "Event type (repo_push, build_success, etc.)")
	notificationUpdateCmd.Flags().StringVarP(&notificationMethod, "method", "m", "", "Method (webhook, email, slack)")
	notificationUpdateCmd.Flags().StringVar(&notificationURL, "url", "", "Webhook URL or email address")
	notificationUpdateCmd.Flags().StringVar(&notificationTitle, "title", "", "Notification title")
	_ = notificationUpdateCmd.MarkFlagRequired("uuid")
}

func initNotificationInfoFlags() {
	notificationInfoCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	_ = notificationInfoCmd.MarkFlagRequired("uuid")
}

func initNotificationCreateFlags() {
	notificationCreateCmd.Flags().StringVarP(&notificationEvent, "event", "e", "", "Event type (repo_push, build_success, etc.)")
	notificationCreateCmd.Flags().StringVarP(&notificationMethod, "method", "m", "", "Method (webhook, email, slack)")
	notificationCreateCmd.Flags().StringVar(&notificationURL, "url", "", "Webhook URL or email address")
	notificationCreateCmd.Flags().StringVar(&notificationTitle, "title", "", "Notification title")
	_ = notificationCreateCmd.MarkFlagRequired("event")
	_ = notificationCreateCmd.MarkFlagRequired("method")
	_ = notificationCreateCmd.MarkFlagRequired("url")
}

func initNotificationDeleteFlags() {
	notificationDeleteCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	notificationDeleteCmd.Flags().BoolVar(&confirmNotificationDel, "confirm", false, "Confirm notification deletion")
	_ = notificationDeleteCmd.MarkFlagRequired("uuid")
}

func initNotificationTestFlags() {
	notificationTestCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	_ = notificationTestCmd.MarkFlagRequired("uuid")
}

func initNotificationResetFlags() {
	notificationResetCmd.Flags().StringVarP(&notificationUUID, "uuid", "u", "", "Notification UUID")
	_ = notificationResetCmd.MarkFlagRequired("uuid")
}
