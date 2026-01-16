// Example: Organization Management with go-quay Library
//
// This example demonstrates organization administration tasks:
// - Getting organization information
// - Managing teams and members
// - Setting up quotas and auto-prune policies
// - Configuring default permissions
//
// Usage:
//
//	export QUAY_TOKEN="your-quay-api-token"
//	go run main.go --organization myorg
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
	orgName := flag.String("organization", "", "Organization name (required)")
	showDetails := flag.Bool("details", false, "Show detailed information")
	flag.Parse()

	if *orgName == "" {
		fmt.Println("Usage: go run main.go --organization <org-name> [--details]")
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

	fmt.Println("=== Organization Management Dashboard ===")
	fmt.Printf("Organization: %s\n\n", *orgName)

	// Step 1: Get organization info
	fmt.Println("1. Organization Information")
	fmt.Println("   " + repeat("-", 40))

	org, err := client.GetOrganization(*orgName)
	if err != nil {
		log.Fatalf("Failed to get organization: %v", err)
	}
	fmt.Printf("   Name: %s\n", org.Name)
	fmt.Printf("   Email: %s\n", org.Email)
	fmt.Printf("   Is Admin: %v\n", org.IsOrgAdmin)
	fmt.Printf("   Can Create Repo: %v\n", org.CanCreateRepo)
	fmt.Println()

	// Step 2: List organization members
	fmt.Println("2. Organization Members")
	fmt.Println("   " + repeat("-", 40))

	members, err := client.GetOrganizationMembers(*orgName)
	if err != nil {
		log.Printf("   Could not get members: %v\n", err)
	} else {
		fmt.Printf("   Total Members: %d\n", len(members.Members))
		if *showDetails {
			for _, member := range members.Members {
				fmt.Printf("   - %s (kind: %s)\n", member.Name, member.Kind)
			}
		}
	}
	fmt.Println()

	// Step 3: List teams
	fmt.Println("3. Teams")
	fmt.Println("   " + repeat("-", 40))

	teams, err := client.GetTeams(*orgName)
	if err != nil {
		log.Printf("   Could not get teams: %v\n", err)
	} else {
		fmt.Printf("   Total Teams: %d\n", len(teams))
		for _, team := range teams {
			fmt.Printf("   - %s (role: %s, members: %d)\n",
				team.Name, team.Role, team.MemberCount)

			if *showDetails {
				// Get team members
				teamMembers, err := client.GetTeamMembers(*orgName, team.Name)
				if err == nil {
					for _, m := range teamMembers.Members {
						fmt.Printf("     * %s\n", m.Name)
					}
				}
			}
		}
	}
	fmt.Println()

	// Step 4: List robot accounts
	fmt.Println("4. Robot Accounts")
	fmt.Println("   " + repeat("-", 40))

	robots, err := client.GetRobotAccounts(*orgName)
	if err != nil {
		log.Printf("   Could not get robots: %v\n", err)
	} else {
		fmt.Printf("   Total Robots: %d\n", len(robots.Robots))
		if *showDetails {
			for _, robot := range robots.Robots {
				fmt.Printf("   - %s\n", robot.Name)
				if robot.Description != "" {
					fmt.Printf("     Description: %s\n", robot.Description)
				}
			}
		}
	}
	fmt.Println()

	// Step 5: Check quota usage
	fmt.Println("5. Quota Information")
	fmt.Println("   " + repeat("-", 40))

	quota, err := client.GetQuota(*orgName)
	if err != nil {
		log.Printf("   Could not get quota: %v\n", err)
	} else if quota.LimitBytes == 0 {
		fmt.Println("   No quota limits configured")
	} else {
		fmt.Printf("   Limit: %s\n", formatBytes(quota.LimitBytes))
		if quota.DefaultLimitBytes > 0 {
			fmt.Printf("   Default Limit: %s\n", formatBytes(quota.DefaultLimitBytes))
		}
	}
	fmt.Println()

	// Step 6: Auto-prune policies
	fmt.Println("6. Auto-Prune Policies")
	fmt.Println("   " + repeat("-", 40))

	policies, err := client.GetAutoPrunePolicies(*orgName)
	if err != nil {
		log.Printf("   Could not get auto-prune policies: %v\n", err)
	} else if len(policies.Policies) == 0 {
		fmt.Println("   No auto-prune policies configured")
		fmt.Println("   Tip: Auto-prune helps manage storage by removing old tags")
	} else {
		fmt.Printf("   Total Policies: %d\n", len(policies.Policies))
		for _, policy := range policies.Policies {
			fmt.Printf("   - %s: %s = %v\n", policy.UUID[:8], policy.Method, policy.Value)
			if policy.TagPattern != "" {
				fmt.Printf("     Pattern: %s\n", policy.TagPattern)
			}
		}
	}
	fmt.Println()

	// Step 7: OAuth Applications
	fmt.Println("7. OAuth Applications")
	fmt.Println("   " + repeat("-", 40))

	apps, err := client.GetApplications(*orgName)
	if err != nil {
		log.Printf("   Could not get applications: %v\n", err)
	} else if len(apps.Applications) == 0 {
		fmt.Println("   No OAuth applications configured")
	} else {
		fmt.Printf("   Total Applications: %d\n", len(apps.Applications))
		for _, app := range apps.Applications {
			fmt.Printf("   - %s (client_id: %s...)\n", app.Name, app.ClientID[:8])
		}
	}
	fmt.Println()

	// Step 8: Default permissions (prototypes)
	fmt.Println("8. Default Permission Prototypes")
	fmt.Println("   " + repeat("-", 40))

	prototypes, err := client.GetPrototypes(*orgName)
	if err != nil {
		log.Printf("   Could not get prototypes: %v\n", err)
	} else if len(prototypes.Prototypes) == 0 {
		fmt.Println("   No default permission prototypes configured")
		fmt.Println("   Tip: Prototypes auto-apply permissions to new repositories")
	} else {
		fmt.Printf("   Total Prototypes: %d\n", len(prototypes.Prototypes))
		for _, p := range prototypes.Prototypes {
			fmt.Printf("   - %s: %s (%s) -> role: %s\n",
				p.ID[:8], p.Delegate.Name, p.Delegate.Kind, p.Role)
		}
	}

	// Summary
	fmt.Println("\n=== Organization Summary ===")
	if members != nil {
		fmt.Printf("Members: %d\n", len(members.Members))
	}
	if teams != nil {
		fmt.Printf("Teams: %d\n", len(teams))
	}
	if robots != nil {
		fmt.Printf("Robots: %d\n", len(robots.Robots))
	}
}

// repeat creates a string by repeating a character
func repeat(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
