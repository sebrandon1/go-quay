package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	orgName  string
	teamName string
)

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "Organization management commands",
	Long: `Commands for managing organizations, teams, members, robots, and other organization-related operations.

Available commands:
  info         - Get organization information
  members      - Get organization members
  teams        - Get organization teams
  team         - Get specific team information
  team-members - Get team members
  robots       - Get organization robots
  quota        - Get organization quota
  auto-prune   - Get auto-prune policies
  applications - Get organization applications`,
}

// Organization Info
var orgInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get organization information",
	Long:  `Get detailed information about an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		org, err := client.GetOrganization(orgName)
		if err != nil {
			fmt.Printf("Error getting organization: %v\n", err)
			os.Exit(1)
		}
		printJSON(org)
	},
}

// Organization Members
var orgMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Get organization members",
	Long:  `Get list of all members in an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		members, err := client.GetOrganizationMembers(orgName)
		if err != nil {
			fmt.Printf("Error getting organization members: %v\n", err)
			os.Exit(1)
		}
		printJSON(members)
	},
}

// Organization Teams
var orgTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Get organization teams",
	Long:  `Get list of all teams in an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		teams, err := client.GetTeams(orgName)
		if err != nil {
			fmt.Printf("Error getting organization teams: %v\n", err)
			os.Exit(1)
		}
		printJSON(teams)
	},
}

// Team Info
var teamInfoCmd = &cobra.Command{
	Use:   "team",
	Short: "Get team information",
	Long:  `Get detailed information about a specific team.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		team, err := client.GetTeam(orgName, teamName)
		if err != nil {
			fmt.Printf("Error getting team: %v\n", err)
			os.Exit(1)
		}
		printJSON(team)
	},
}

// Team Members
var teamMembersCmd = &cobra.Command{
	Use:   "team-members",
	Short: "Get team members",
	Long:  `Get list of all members in a specific team.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		members, err := client.GetTeamMembers(orgName, teamName)
		if err != nil {
			fmt.Printf("Error getting team members: %v\n", err)
			os.Exit(1)
		}
		printJSON(members)
	},
}

// Organization Robots
var orgRobotsCmd = &cobra.Command{
	Use:   "robots",
	Short: "Get organization robots",
	Long:  `Get list of all robot accounts in an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		robots, err := client.GetRobotAccounts(orgName)
		if err != nil {
			fmt.Printf("Error getting organization robots: %v\n", err)
			os.Exit(1)
		}
		printJSON(robots)
	},
}

// Organization Quota
var orgQuotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Get organization quota",
	Long:  `Get quota information for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		quota, err := client.GetQuota(orgName)
		if err != nil {
			fmt.Printf("Error getting organization quota: %v\n", err)
			os.Exit(1)
		}
		printJSON(quota)
	},
}

// Auto-prune Policies
var autoPruneCmd = &cobra.Command{
	Use:   "auto-prune",
	Short: "Get auto-prune policies",
	Long:  `Get list of auto-prune policies for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		policies, err := client.GetAutoPrunePolicies(orgName)
		if err != nil {
			fmt.Printf("Error getting auto-prune policies: %v\n", err)
			os.Exit(1)
		}
		printJSON(policies)
	},
}

// Organization Applications
var orgApplicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Get organization applications",
	Long:  `Get list of OAuth applications for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		applications, err := client.GetApplications(orgName)
		if err != nil {
			fmt.Printf("Error getting organization applications: %v\n", err)
			os.Exit(1)
		}
		printJSON(applications)
	},
}

func init() {
	// Add organization command to get command
	organizationCmd.AddCommand(orgInfoCmd)
	organizationCmd.AddCommand(orgMembersCmd)
	organizationCmd.AddCommand(orgTeamsCmd)
	organizationCmd.AddCommand(teamInfoCmd)
	organizationCmd.AddCommand(teamMembersCmd)
	organizationCmd.AddCommand(orgRobotsCmd)
	organizationCmd.AddCommand(orgQuotaCmd)
	organizationCmd.AddCommand(autoPruneCmd)
	organizationCmd.AddCommand(orgApplicationsCmd)

	// Add persistent flags for organization name
	organizationCmd.PersistentFlags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	organizationCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Authentication token")

	// Mark flags as required
	if err := organizationCmd.MarkPersistentFlagRequired("organization"); err != nil {
		fmt.Printf("Error marking organization flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := organizationCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}

	// Add specific flags for team commands
	teamInfoCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	if err := teamInfoCmd.MarkFlagRequired("team"); err != nil {
		fmt.Printf("Error marking team flag as required: %v\n", err)
		os.Exit(1)
	}

	teamMembersCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	if err := teamMembersCmd.MarkFlagRequired("team"); err != nil {
		fmt.Printf("Error marking team flag as required: %v\n", err)
		os.Exit(1)
	}
}

// Helper function to print JSON output
func printJSON(data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
