// Example: CI/CD Integration with go-quay Library
//
// This example demonstrates how to use go-quay for CI/CD automation:
// - Creating and managing robot accounts
// - Setting up repository permissions for robots
// - Configuring webhook notifications for builds
// - Managing build triggers
//
// Usage:
//
//	export QUAY_TOKEN="your-quay-api-token"
//	go run main.go --namespace myorg --repository myapp
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sebrandon1/go-quay/lib"
)

func main() {
	// Parse command line arguments
	namespace := flag.String("namespace", "", "Organization namespace (required)")
	repository := flag.String("repository", "", "Repository name (required)")
	robotName := flag.String("robot", "ci-bot", "Name for the CI robot account")
	webhookURL := flag.String("webhook", "", "Webhook URL for notifications (optional)")
	flag.Parse()

	if *namespace == "" || *repository == "" {
		fmt.Println("Usage: go run main.go --namespace <org> --repository <repo> [--robot <name>] [--webhook <url>]")
		os.Exit(1)
	}

	// Get authentication token
	token := os.Getenv("QUAY_TOKEN")
	if token == "" {
		log.Fatal("QUAY_TOKEN environment variable is required")
	}

	// Initialize client
	client, err := lib.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("=== CI/CD Integration Setup ===")
	fmt.Printf("Organization: %s\n", *namespace)
	fmt.Printf("Repository: %s\n\n", *repository)

	// Step 1: Check organization access
	fmt.Println("1. Verifying organization access...")
	org, err := client.GetOrganization(*namespace)
	if err != nil {
		log.Fatalf("Failed to get organization (do you have access?): %v", err)
	}
	fmt.Printf("   Organization: %s\n", org.Name)
	fmt.Printf("   Email: %s\n\n", org.Email)

	// Step 2: List existing robot accounts
	fmt.Println("2. Checking existing robot accounts...")
	robots, err := client.GetOrganizationRobots(*namespace)
	if err != nil {
		log.Printf("   Could not list robots: %v\n", err)
	} else {
		fmt.Printf("   Found %d existing robot accounts:\n", len(robots.Robots))
		for _, robot := range robots.Robots {
			fmt.Printf("   - %s\n", robot.Name)
		}
	}
	fmt.Println()

	// Step 3: Create or get robot account
	fullRobotName := *namespace + "+" + *robotName
	fmt.Printf("3. Setting up robot account '%s'...\n", fullRobotName)

	// Check if robot exists
	existingRobot, err := client.GetOrganizationRobot(*namespace, *robotName)
	if err == nil {
		fmt.Println("   Robot already exists, using existing account")
		fmt.Printf("   Name: %s\n", existingRobot.Name)
		fmt.Println("   Note: Token is only shown at creation time")
	} else {
		// Create new robot
		fmt.Println("   Creating new robot account...")
		newRobot, err := client.CreateOrganizationRobot(*namespace, *robotName, "CI/CD automation robot", nil)
		if err != nil {
			log.Fatalf("Failed to create robot: %v", err)
		}
		fmt.Printf("   Created robot: %s\n", newRobot.Name)
		fmt.Printf("   IMPORTANT - Save this token (shown only once):\n")
		fmt.Printf("   Token: %s\n", newRobot.Token)
	}
	fmt.Println()

	// Step 4: Set repository permissions for the robot
	fmt.Printf("4. Setting repository permissions for '%s'...\n", fullRobotName)

	// Grant write permission to the robot
	err = client.SetUserPermission(*namespace, *repository, fullRobotName, "write")
	if err != nil {
		log.Printf("   Could not set permission: %v\n", err)
	} else {
		fmt.Printf("   Granted 'write' permission to %s on %s/%s\n", fullRobotName, *namespace, *repository)
	}

	// Verify permission was set
	perm, err := client.GetUserPermission(*namespace, *repository, fullRobotName)
	if err != nil {
		log.Printf("   Could not verify permission: %v\n", err)
	} else {
		fmt.Printf("   Verified: %s has '%s' role\n", perm.Name, perm.Role)
	}
	fmt.Println()

	// Step 5: Set up webhook notification (if URL provided)
	if *webhookURL != "" {
		fmt.Println("5. Setting up webhook notification...")

		notification, err := client.CreateNotification(*namespace, *repository, lib.CreateNotificationRequest{
			Event:  "repo_push",
			Method: "webhook",
			Title:  "CI/CD Push Notification",
			Config: lib.NotificationConfig{
				URL: *webhookURL,
			},
		})
		if err != nil {
			log.Printf("   Could not create notification: %v\n", err)
		} else {
			fmt.Printf("   Created notification: %s\n", notification.UUID)
			fmt.Printf("   Event: %s\n", notification.Event)
			fmt.Printf("   Method: %s\n", notification.Method)
			fmt.Printf("   URL: %s\n", *webhookURL)
		}
		fmt.Println()
	}

	// Step 6: List build triggers
	fmt.Println("6. Checking build triggers...")
	triggers, err := client.ListBuildTriggers(*namespace, *repository)
	if err != nil {
		log.Printf("   Could not list triggers: %v\n", err)
	} else if len(triggers.Triggers) == 0 {
		fmt.Println("   No build triggers configured")
		fmt.Println("   Tip: Configure build triggers in the Quay.io UI")
	} else {
		fmt.Printf("   Found %d build triggers:\n", len(triggers.Triggers))
		for _, trigger := range triggers.Triggers {
			status := "enabled"
			if !trigger.Enabled {
				status = "disabled"
			}
			fmt.Printf("   - %s (%s) [%s]\n", trigger.ID, trigger.Service, status)
		}
	}
	fmt.Println()

	// Step 7: List existing notifications
	fmt.Println("7. Checking existing notifications...")
	notifications, err := client.ListNotifications(*namespace, *repository)
	if err != nil {
		log.Printf("   Could not list notifications: %v\n", err)
	} else if len(notifications.Notifications) == 0 {
		fmt.Println("   No notifications configured")
	} else {
		fmt.Printf("   Found %d notifications:\n", len(notifications.Notifications))
		for _, n := range notifications.Notifications {
			fmt.Printf("   - %s: %s (%s)\n", n.UUID[:8], n.Event, n.Method)
		}
	}

	// Summary
	fmt.Println("\n=== CI/CD Setup Summary ===")
	fmt.Printf("Robot Account: %s\n", fullRobotName)
	fmt.Printf("Repository: %s/%s\n", *namespace, *repository)
	fmt.Println("\nDocker login command:")
	fmt.Printf("  docker login quay.io -u %s -p <robot-token>\n", fullRobotName)
	fmt.Println("\nPush command:")
	fmt.Printf("  docker push quay.io/%s/%s:tag\n", *namespace, *repository)
}
