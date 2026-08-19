package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	orgName          string
	teamName         string
	email            string
	memberName       string
	confirm          bool
	description      string
	role             string
	upstreamRegistry string
	insecure         bool
	expiration       int
	clientID         string
	appName          string
	applicationURI   string
	redirectURI      string
	sku              string
	quantity         int
	subscriptionID   string
	subscriptionIDs  []string
	limitBytes       int64
	policyUUID       string
	method           string
	pruneValue       int
	tagPattern       string
)

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:   cmdOrganization,
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
	Use:   subcmdInfo,
	Short: "Get organization information",
	Long:  `Get detailed information about an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		org, err := client.GetOrganization(cmd.Context(), orgName)
		if err != nil {
			return fmt.Errorf("getting organization: %w", err)
		}
		return printJSON(org)
	},
}

// Organization Members
var orgMembersCmd = &cobra.Command{
	Use:   cmdMembers,
	Short: "Get organization members",
	Long:  `Get list of all members in an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		members, err := client.GetOrganizationMembers(cmd.Context(), orgName)
		if err != nil {
			return fmt.Errorf("getting organization members: %w", err)
		}
		return printJSON(members)
	},
}

// Create Organization
var createOrgCmd = &cobra.Command{
	Use:   "create-org",
	Short: "Create an organization",
	Long:  `Create a new organization with the specified name and email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		org, err := client.CreateOrganization(cmd.Context(), orgName, email)
		if err != nil {
			return fmt.Errorf("creating organization: %w", err)
		}
		return printJSON(org)
	},
}

// Update Organization
var updateOrgCmd = &cobra.Command{
	Use:   "update-org",
	Short: "Update an organization",
	Long:  `Update an organization's email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		org, err := client.UpdateOrganization(cmd.Context(), orgName, email)
		if err != nil {
			return fmt.Errorf("updating organization: %w", err)
		}
		return printJSON(org)
	},
}

// Delete Organization
var deleteOrgCmd = &cobra.Command{
	Use:   "delete-org",
	Short: "Delete an organization",
	Long:  `Delete an organization. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete an organization")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteOrganization(cmd.Context(), orgName)
		if err != nil {
			return fmt.Errorf("deleting organization: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Organization deleted successfully")
		return nil
	},
}

// Add Organization Member
var addMemberCmd = &cobra.Command{
	Use:   "add-member",
	Short: "Add a member to an organization",
	Long:  `Add a member to an organization by member name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.AddOrganizationMember(cmd.Context(), orgName, memberName)
		if err != nil {
			return fmt.Errorf("adding organization member: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Member added successfully")
		return nil
	},
}

// Remove Organization Member
var removeMemberCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "Remove a member from an organization",
	Long:  `Remove a member from an organization. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to remove a member")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.RemoveOrganizationMember(cmd.Context(), orgName, memberName)
		if err != nil {
			return fmt.Errorf("removing organization member: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Member removed successfully")
		return nil
	},
}

// Get Organization Member
var getMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Get organization member information",
	Long:  `Get detailed information about a specific organization member.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		member, err := client.GetOrganizationMember(cmd.Context(), orgName, memberName)
		if err != nil {
			return fmt.Errorf("getting organization member: %w", err)
		}
		return printJSON(member)
	},
}

// Get Organization Collaborators
var collaboratorsCmd = &cobra.Command{
	Use:   "collaborators",
	Short: "Get organization collaborators",
	Long:  `Get list of collaborators for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		collaborators, err := client.GetOrganizationCollaborators(cmd.Context(), orgName)
		if err != nil {
			return fmt.Errorf("getting organization collaborators: %w", err)
		}
		return printJSON(collaborators)
	},
}

// Get Organization Repositories
var orgRepositoriesCmd = &cobra.Command{
	Use:   "repositories",
	Short: "Get organization repositories",
	Long:  `Get list of repositories for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		repos, err := client.GetOrganizationRepositories(cmd.Context(), orgName)
		if err != nil {
			return fmt.Errorf("getting organization repositories: %w", err)
		}
		return printJSON(repos)
	},
}

func initOrgPersistentFlags() {
	organizationCmd.PersistentFlags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	_ = organizationCmd.MarkPersistentFlagRequired("organization")
}

func initOrgMemberFlags() {
	createOrgCmd.Flags().StringVar(&email, "email", "", "Email address")
	_ = createOrgCmd.MarkFlagRequired("email")

	updateOrgCmd.Flags().StringVar(&email, "email", "", "Email address")
	_ = updateOrgCmd.MarkFlagRequired("email")

	deleteOrgCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	addMemberCmd.Flags().StringVar(&memberName, "member", "", "Member name")
	_ = addMemberCmd.MarkFlagRequired("member")

	removeMemberCmd.Flags().StringVar(&memberName, "member", "", "Member name")
	_ = removeMemberCmd.MarkFlagRequired("member")
	removeMemberCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm removal")

	getMemberCmd.Flags().StringVar(&memberName, "member", "", "Member name")
	_ = getMemberCmd.MarkFlagRequired("member")
}

func init() {
	organizationCmd.AddCommand(orgInfoCmd)
	organizationCmd.AddCommand(orgMembersCmd)
	organizationCmd.AddCommand(createOrgCmd)
	organizationCmd.AddCommand(updateOrgCmd)
	organizationCmd.AddCommand(deleteOrgCmd)
	organizationCmd.AddCommand(addMemberCmd)
	organizationCmd.AddCommand(removeMemberCmd)
	organizationCmd.AddCommand(getMemberCmd)
	organizationCmd.AddCommand(collaboratorsCmd)
	organizationCmd.AddCommand(orgRepositoriesCmd)

	initOrgPersistentFlags()
	initOrgMemberFlags()
}
