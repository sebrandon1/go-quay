package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
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
	delegateType     string
	delegateName     string
	prototypeID      string
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

// Create Organization
var createOrgCmd = &cobra.Command{
	Use:   "create-org",
	Short: "Create an organization",
	Long:  `Create a new organization with the specified name and email.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		org, err := client.CreateOrganization(orgName, email)
		if err != nil {
			fmt.Printf("Error creating organization: %v\n", err)
			os.Exit(1)
		}
		printJSON(org)
	},
}

// Update Organization
var updateOrgCmd = &cobra.Command{
	Use:   "update-org",
	Short: "Update an organization",
	Long:  `Update an organization's email.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		org, err := client.UpdateOrganization(orgName, email)
		if err != nil {
			fmt.Printf("Error updating organization: %v\n", err)
			os.Exit(1)
		}
		printJSON(org)
	},
}

// Delete Organization
var deleteOrgCmd = &cobra.Command{
	Use:   "delete-org",
	Short: "Delete an organization",
	Long:  `Delete an organization. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to delete an organization")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteOrganization(orgName)
		if err != nil {
			fmt.Printf("Error deleting organization: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Organization deleted successfully")
	},
}

// Add Organization Member
var addMemberCmd = &cobra.Command{
	Use:   "add-member",
	Short: "Add a member to an organization",
	Long:  `Add a member to an organization by member name.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.AddOrganizationMember(orgName, memberName)
		if err != nil {
			fmt.Printf("Error adding organization member: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Member added successfully")
	},
}

// Remove Organization Member
var removeMemberCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "Remove a member from an organization",
	Long:  `Remove a member from an organization. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to remove a member")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.RemoveOrganizationMember(orgName, memberName)
		if err != nil {
			fmt.Printf("Error removing organization member: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Member removed successfully")
	},
}

// Get Organization Member
var getMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Get organization member information",
	Long:  `Get detailed information about a specific organization member.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		member, err := client.GetOrganizationMember(orgName, memberName)
		if err != nil {
			fmt.Printf("Error getting organization member: %v\n", err)
			os.Exit(1)
		}
		printJSON(member)
	},
}

// Get Organization Collaborators
var collaboratorsCmd = &cobra.Command{
	Use:   "collaborators",
	Short: "Get organization collaborators",
	Long:  `Get list of collaborators for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		collaborators, err := client.GetOrganizationCollaborators(orgName)
		if err != nil {
			fmt.Printf("Error getting organization collaborators: %v\n", err)
			os.Exit(1)
		}
		printJSON(collaborators)
	},
}

// Get Organization Repositories
var orgRepositoriesCmd = &cobra.Command{
	Use:   "repositories",
	Short: "Get organization repositories",
	Long:  `Get list of repositories for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		repos, err := client.GetOrganizationRepositories(orgName)
		if err != nil {
			fmt.Printf("Error getting organization repositories: %v\n", err)
			os.Exit(1)
		}
		printJSON(repos)
	},
}

// Get Default Permissions
var defaultPermissionsCmd = &cobra.Command{
	Use:   "default-permissions",
	Short: "Get default permissions",
	Long:  `Get default permissions for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		perms, err := client.GetDefaultPermissions(orgName)
		if err != nil {
			fmt.Printf("Error getting default permissions: %v\n", err)
			os.Exit(1)
		}
		printJSON(perms)
	},
}

// Create Default Permission
var createDefaultPermissionCmd = &cobra.Command{
	Use:   "create-default-permission",
	Short: "Create a default permission",
	Long:  `Create a default permission for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		perm, err := client.CreateDefaultPermission(orgName, role, delegateType, delegateName)
		if err != nil {
			fmt.Printf("Error creating default permission: %v\n", err)
			os.Exit(1)
		}
		printJSON(perm)
	},
}

// Delete Default Permission
var deleteDefaultPermissionCmd = &cobra.Command{
	Use:   "delete-default-permission",
	Short: "Delete a default permission",
	Long:  `Delete a default permission for an organization. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to delete a default permission")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteDefaultPermission(orgName, prototypeID)
		if err != nil {
			fmt.Printf("Error deleting default permission: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Default permission deleted successfully")
	},
}

// Get Proxy Cache Config
var proxyCacheCmd = &cobra.Command{
	Use:   "proxy-cache",
	Short: "Get proxy cache configuration",
	Long:  `Get proxy cache configuration for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		config, err := client.GetProxyCacheConfig(orgName)
		if err != nil {
			fmt.Printf("Error getting proxy cache config: %v\n", err)
			os.Exit(1)
		}
		printJSON(config)
	},
}

// Create Proxy Cache Config
var createProxyCacheCmd = &cobra.Command{
	Use:   "create-proxy-cache",
	Short: "Create proxy cache configuration",
	Long:  `Create proxy cache configuration for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		config, err := client.CreateProxyCacheConfig(orgName, upstreamRegistry, insecure, expiration)
		if err != nil {
			fmt.Printf("Error creating proxy cache config: %v\n", err)
			os.Exit(1)
		}
		printJSON(config)
	},
}

// Delete Proxy Cache Config
var deleteProxyCacheCmd = &cobra.Command{
	Use:   "delete-proxy-cache",
	Short: "Delete proxy cache configuration",
	Long:  `Delete proxy cache configuration for an organization. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to delete proxy cache config")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteProxyCacheConfig(orgName)
		if err != nil {
			fmt.Printf("Error deleting proxy cache config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Proxy cache config deleted successfully")
	},
}

// Get Robot Account
var orgRobotCmd = &cobra.Command{
	Use:   "robot",
	Short: "Get robot account information",
	Long:  `Get detailed information about a specific robot account.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		robot, err := client.GetRobotAccount(orgName, robotShortname)
		if err != nil {
			fmt.Printf("Error getting robot account: %v\n", err)
			os.Exit(1)
		}
		printJSON(robot)
	},
}

// Create Robot Account
var createRobotCmd = &cobra.Command{
	Use:   "create-robot",
	Short: "Create a robot account",
	Long:  `Create a new robot account in an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		robot, err := client.CreateRobotAccount(orgName, robotShortname, description, nil)
		if err != nil {
			fmt.Printf("Error creating robot account: %v\n", err)
			os.Exit(1)
		}
		printJSON(robot)
	},
}

// Delete Robot Account
var deleteRobotCmd = &cobra.Command{
	Use:   "delete-robot",
	Short: "Delete a robot account",
	Long:  `Delete a robot account from an organization. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to delete a robot account")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteRobotAccount(orgName, robotShortname)
		if err != nil {
			fmt.Printf("Error deleting robot account: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Robot account deleted successfully")
	},
}

// Regenerate Robot Token
var regenerateRobotCmd = &cobra.Command{
	Use:   "regenerate-robot",
	Short: "Regenerate a robot account token",
	Long:  `Regenerate the token for a robot account.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		robot, err := client.RegenerateRobotToken(orgName, robotShortname)
		if err != nil {
			fmt.Printf("Error regenerating robot token: %v\n", err)
			os.Exit(1)
		}
		printJSON(robot)
	},
}

// Get Robot Permissions
var orgRobotPermissionsCmd = &cobra.Command{
	Use:   "robot-permissions",
	Short: "Get robot account permissions",
	Long:  `Get permissions for a specific robot account.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		perms, err := client.GetRobotPermissions(orgName, robotShortname)
		if err != nil {
			fmt.Printf("Error getting robot permissions: %v\n", err)
			os.Exit(1)
		}
		printJSON(perms)
	},
}

// Set Robot Repository Permission
var setRobotPermissionCmd = &cobra.Command{
	Use:   "set-robot-permission",
	Short: "Set robot repository permission",
	Long:  `Set a robot account's permission on a repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.SetRobotRepositoryPermission(orgName, robotShortname, repository, role)
		if err != nil {
			fmt.Printf("Error setting robot permission: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Robot permission set successfully")
	},
}

// Remove Robot Repository Permission
var removeRobotPermissionCmd = &cobra.Command{
	Use:   "remove-robot-permission",
	Short: "Remove robot repository permission",
	Long:  `Remove a robot account's permission on a repository. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to remove a robot permission")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.RemoveRobotRepositoryPermission(orgName, robotShortname, repository)
		if err != nil {
			fmt.Printf("Error removing robot permission: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Robot permission removed successfully")
	},
}

// Get Application
var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Get application information",
	Long:  `Get detailed information about a specific OAuth application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		app, err := client.GetApplication(orgName, clientID)
		if err != nil {
			fmt.Printf("Error getting application: %v\n", err)
			os.Exit(1)
		}
		printJSON(app)
	},
}

// Create Application
var createApplicationCmd = &cobra.Command{
	Use:   "create-application",
	Short: "Create an OAuth application",
	Long:  `Create a new OAuth application for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		app, err := client.CreateApplication(orgName, appName, description, applicationURI, redirectURI)
		if err != nil {
			fmt.Printf("Error creating application: %v\n", err)
			os.Exit(1)
		}
		printJSON(app)
	},
}

// Update Application
var updateApplicationCmd = &cobra.Command{
	Use:   "update-application",
	Short: "Update an OAuth application",
	Long:  `Update an existing OAuth application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		app, err := client.UpdateApplication(orgName, clientID, appName, description, applicationURI, redirectURI)
		if err != nil {
			fmt.Printf("Error updating application: %v\n", err)
			os.Exit(1)
		}
		printJSON(app)
	},
}

// Delete Application
var deleteApplicationCmd = &cobra.Command{
	Use:   "delete-application",
	Short: "Delete an OAuth application",
	Long:  `Delete an OAuth application. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to delete an application")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteApplication(orgName, clientID)
		if err != nil {
			fmt.Printf("Error deleting application: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Application deleted successfully")
	},
}

// Reset Application Client Secret
var resetApplicationSecretCmd = &cobra.Command{
	Use:   "reset-application-secret",
	Short: "Reset application client secret",
	Long:  `Reset the client secret for an OAuth application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		app, err := client.ResetApplicationClientSecret(orgName, clientID)
		if err != nil {
			fmt.Printf("Error resetting application secret: %v\n", err)
			os.Exit(1)
		}
		printJSON(app)
	},
}

// Get Organization Marketplace
var marketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Get organization marketplace information",
	Long:  `Get marketplace information for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		marketplace, err := client.GetOrganizationMarketplace(orgName)
		if err != nil {
			fmt.Printf("Error getting marketplace info: %v\n", err)
			os.Exit(1)
		}
		printJSON(marketplace)
	},
}

// Create Marketplace Subscription
var createMarketplaceSubscriptionCmd = &cobra.Command{
	Use:   "create-marketplace-subscription",
	Short: "Create a marketplace subscription",
	Long:  `Create a marketplace subscription for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.CreateOrganizationMarketplaceSubscription(orgName, &lib.MarketplaceSubscriptionRequest{SKU: sku, Quantity: quantity})
		if err != nil {
			fmt.Printf("Error creating marketplace subscription: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Marketplace subscription created successfully")
	},
}

// Delete Marketplace Subscription
var deleteMarketplaceSubscriptionCmd = &cobra.Command{
	Use:   "delete-marketplace-subscription",
	Short: "Delete a marketplace subscription",
	Long:  `Delete a marketplace subscription. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to delete a marketplace subscription")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteOrganizationMarketplaceSubscription(orgName, subscriptionID)
		if err != nil {
			fmt.Printf("Error deleting marketplace subscription: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Marketplace subscription deleted successfully")
	},
}

// Batch Remove Marketplace Subscriptions
var batchRemoveSubscriptionsCmd = &cobra.Command{
	Use:   "batch-remove-subscriptions",
	Short: "Batch remove marketplace subscriptions",
	Long:  `Remove multiple marketplace subscriptions at once. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to batch remove subscriptions")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.BatchRemoveOrganizationMarketplaceSubscriptions(orgName, subscriptionIDs)
		if err != nil {
			fmt.Printf("Error batch removing subscriptions: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Marketplace subscriptions removed successfully")
	},
}

// Create Quota
var createQuotaCmd = &cobra.Command{
	Use:   "create-quota",
	Short: "Create organization quota",
	Long:  `Create a quota for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		quota, err := client.CreateQuota(orgName, limitBytes)
		if err != nil {
			fmt.Printf("Error creating quota: %v\n", err)
			os.Exit(1)
		}
		printJSON(quota)
	},
}

// Update Quota
var updateQuotaCmd = &cobra.Command{
	Use:   "update-quota",
	Short: "Update organization quota",
	Long:  `Update the quota for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		quota, err := client.UpdateQuota(orgName, limitBytes)
		if err != nil {
			fmt.Printf("Error updating quota: %v\n", err)
			os.Exit(1)
		}
		printJSON(quota)
	},
}

// Delete Quota
var deleteQuotaCmd = &cobra.Command{
	Use:   "delete-quota",
	Short: "Delete organization quota",
	Long:  `Delete the quota for an organization. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to delete a quota")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteQuota(orgName)
		if err != nil {
			fmt.Printf("Error deleting quota: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Quota deleted successfully")
	},
}

// Get Auto-Prune Policy
var autoPrunePolicyCmd = &cobra.Command{
	Use:   "auto-prune-policy",
	Short: "Get a specific auto-prune policy",
	Long:  `Get detailed information about a specific auto-prune policy.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		policy, err := client.GetAutoPrunePolicy(orgName, policyUUID)
		if err != nil {
			fmt.Printf("Error getting auto-prune policy: %v\n", err)
			os.Exit(1)
		}
		printJSON(policy)
	},
}

// Create Auto-Prune Policy
var createAutoPruneCmd = &cobra.Command{
	Use:   "create-auto-prune",
	Short: "Create an auto-prune policy",
	Long:  `Create an auto-prune policy for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		policy, err := client.CreateAutoPrunePolicy(orgName, method, pruneValue, tagPattern)
		if err != nil {
			fmt.Printf("Error creating auto-prune policy: %v\n", err)
			os.Exit(1)
		}
		printJSON(policy)
	},
}

// Update Auto-Prune Policy
var updateAutoPruneCmd = &cobra.Command{
	Use:   "update-auto-prune",
	Short: "Update an auto-prune policy",
	Long:  `Update an existing auto-prune policy.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		policy, err := client.UpdateAutoPrunePolicy(orgName, policyUUID, method, pruneValue, tagPattern)
		if err != nil {
			fmt.Printf("Error updating auto-prune policy: %v\n", err)
			os.Exit(1)
		}
		printJSON(policy)
	},
}

// Delete Auto-Prune Policy
var deleteAutoPruneCmd = &cobra.Command{
	Use:   "delete-auto-prune",
	Short: "Delete an auto-prune policy",
	Long:  `Delete an auto-prune policy. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to delete an auto-prune policy")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteAutoPrunePolicy(orgName, policyUUID)
		if err != nil {
			fmt.Printf("Error deleting auto-prune policy: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Auto-prune policy deleted successfully")
	},
}

// Invite Team Member
var inviteMemberCmd = &cobra.Command{
	Use:   "invite-member",
	Short: "Invite a member to a team",
	Long:  `Invite a member to a team by email.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.InviteTeamMember(orgName, teamName, email)
		if err != nil {
			fmt.Printf("Error inviting team member: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Team member invited successfully")
	},
}

// Cancel Team Invite
var cancelInviteCmd = &cobra.Command{
	Use:   "cancel-invite",
	Short: "Cancel a team member invitation",
	Long:  `Cancel a pending team member invitation. Requires --confirm flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm {
			fmt.Println("Error: must pass --confirm to cancel an invite")
			os.Exit(1)
		}
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}
		err = client.DeleteTeamInvite(orgName, teamName, email)
		if err != nil {
			fmt.Printf("Error canceling team invite: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Team invite canceled successfully")
	},
}

func init() {
	// Add all subcommands to organization command
	organizationCmd.AddCommand(orgInfoCmd)
	organizationCmd.AddCommand(orgMembersCmd)
	organizationCmd.AddCommand(orgTeamsCmd)
	organizationCmd.AddCommand(teamInfoCmd)
	organizationCmd.AddCommand(teamMembersCmd)
	organizationCmd.AddCommand(orgRobotsCmd)
	organizationCmd.AddCommand(orgQuotaCmd)
	organizationCmd.AddCommand(autoPruneCmd)
	organizationCmd.AddCommand(orgApplicationsCmd)
	organizationCmd.AddCommand(createOrgCmd)
	organizationCmd.AddCommand(updateOrgCmd)
	organizationCmd.AddCommand(deleteOrgCmd)
	organizationCmd.AddCommand(addMemberCmd)
	organizationCmd.AddCommand(removeMemberCmd)
	organizationCmd.AddCommand(getMemberCmd)
	organizationCmd.AddCommand(collaboratorsCmd)
	organizationCmd.AddCommand(orgRepositoriesCmd)
	organizationCmd.AddCommand(defaultPermissionsCmd)
	organizationCmd.AddCommand(createDefaultPermissionCmd)
	organizationCmd.AddCommand(deleteDefaultPermissionCmd)
	organizationCmd.AddCommand(proxyCacheCmd)
	organizationCmd.AddCommand(createProxyCacheCmd)
	organizationCmd.AddCommand(deleteProxyCacheCmd)
	organizationCmd.AddCommand(orgRobotCmd)
	organizationCmd.AddCommand(createRobotCmd)
	organizationCmd.AddCommand(deleteRobotCmd)
	organizationCmd.AddCommand(regenerateRobotCmd)
	organizationCmd.AddCommand(orgRobotPermissionsCmd)
	organizationCmd.AddCommand(setRobotPermissionCmd)
	organizationCmd.AddCommand(removeRobotPermissionCmd)
	organizationCmd.AddCommand(applicationCmd)
	organizationCmd.AddCommand(createApplicationCmd)
	organizationCmd.AddCommand(updateApplicationCmd)
	organizationCmd.AddCommand(deleteApplicationCmd)
	organizationCmd.AddCommand(resetApplicationSecretCmd)
	organizationCmd.AddCommand(marketplaceCmd)
	organizationCmd.AddCommand(createMarketplaceSubscriptionCmd)
	organizationCmd.AddCommand(deleteMarketplaceSubscriptionCmd)
	organizationCmd.AddCommand(batchRemoveSubscriptionsCmd)
	organizationCmd.AddCommand(createQuotaCmd)
	organizationCmd.AddCommand(updateQuotaCmd)
	organizationCmd.AddCommand(deleteQuotaCmd)
	organizationCmd.AddCommand(autoPrunePolicyCmd)
	organizationCmd.AddCommand(createAutoPruneCmd)
	organizationCmd.AddCommand(updateAutoPruneCmd)
	organizationCmd.AddCommand(deleteAutoPruneCmd)
	organizationCmd.AddCommand(inviteMemberCmd)
	organizationCmd.AddCommand(cancelInviteCmd)

	// Register persistent and command-specific flags
	initOrgPersistentFlags()
	initOrgMemberFlags()
	initOrgPermissionFlags()
	initOrgProxyCacheFlags()
	initOrgRobotFlags()
	initOrgApplicationFlags()
	initOrgMarketplaceFlags()
	initOrgQuotaFlags()
	initOrgAutoPruneFlags()
	initOrgInviteFlags()
}

func initOrgPersistentFlags() {
	organizationCmd.PersistentFlags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	organizationCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Authentication token")

	if err := organizationCmd.MarkPersistentFlagRequired("organization"); err != nil {
		fmt.Printf("Error marking organization flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := organizationCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}
}

func initOrgMemberFlags() {
	// team info flags
	teamInfoCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	if err := teamInfoCmd.MarkFlagRequired("team"); err != nil {
		fmt.Printf("Error marking team flag as required: %v\n", err)
		os.Exit(1)
	}

	// team members flags
	teamMembersCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	if err := teamMembersCmd.MarkFlagRequired("team"); err != nil {
		fmt.Printf("Error marking team flag as required: %v\n", err)
		os.Exit(1)
	}

	// create-org flags
	createOrgCmd.Flags().StringVar(&email, "email", "", "Email address")
	if err := createOrgCmd.MarkFlagRequired("email"); err != nil {
		fmt.Printf("Error marking email flag as required: %v\n", err)
		os.Exit(1)
	}

	// update-org flags
	updateOrgCmd.Flags().StringVar(&email, "email", "", "Email address")
	if err := updateOrgCmd.MarkFlagRequired("email"); err != nil {
		fmt.Printf("Error marking email flag as required: %v\n", err)
		os.Exit(1)
	}

	// delete-org flags
	deleteOrgCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	// add-member flags
	addMemberCmd.Flags().StringVar(&memberName, "member", "", "Member name")
	if err := addMemberCmd.MarkFlagRequired("member"); err != nil {
		fmt.Printf("Error marking member flag as required: %v\n", err)
		os.Exit(1)
	}

	// remove-member flags
	removeMemberCmd.Flags().StringVar(&memberName, "member", "", "Member name")
	if err := removeMemberCmd.MarkFlagRequired("member"); err != nil {
		fmt.Printf("Error marking member flag as required: %v\n", err)
		os.Exit(1)
	}
	removeMemberCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm removal")

	// member flags
	getMemberCmd.Flags().StringVar(&memberName, "member", "", "Member name")
	if err := getMemberCmd.MarkFlagRequired("member"); err != nil {
		fmt.Printf("Error marking member flag as required: %v\n", err)
		os.Exit(1)
	}
}

func initOrgPermissionFlags() {
	// create-default-permission flags
	createDefaultPermissionCmd.Flags().StringVar(&role, "role", "", "Permission role")
	if err := createDefaultPermissionCmd.MarkFlagRequired("role"); err != nil {
		fmt.Printf("Error marking role flag as required: %v\n", err)
		os.Exit(1)
	}
	createDefaultPermissionCmd.Flags().StringVar(&delegateType, "delegate-type", "", "Delegate type")
	if err := createDefaultPermissionCmd.MarkFlagRequired("delegate-type"); err != nil {
		fmt.Printf("Error marking delegate-type flag as required: %v\n", err)
		os.Exit(1)
	}
	createDefaultPermissionCmd.Flags().StringVar(&delegateName, "delegate-name", "", "Delegate name")
	if err := createDefaultPermissionCmd.MarkFlagRequired("delegate-name"); err != nil {
		fmt.Printf("Error marking delegate-name flag as required: %v\n", err)
		os.Exit(1)
	}

	// delete-default-permission flags
	deleteDefaultPermissionCmd.Flags().StringVar(&prototypeID, "prototype-id", "", "Prototype ID")
	if err := deleteDefaultPermissionCmd.MarkFlagRequired("prototype-id"); err != nil {
		fmt.Printf("Error marking prototype-id flag as required: %v\n", err)
		os.Exit(1)
	}
	deleteDefaultPermissionCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func initOrgProxyCacheFlags() {
	// create-proxy-cache flags
	createProxyCacheCmd.Flags().StringVar(&upstreamRegistry, "upstream-registry", "", "Upstream registry URL")
	if err := createProxyCacheCmd.MarkFlagRequired("upstream-registry"); err != nil {
		fmt.Printf("Error marking upstream-registry flag as required: %v\n", err)
		os.Exit(1)
	}
	createProxyCacheCmd.Flags().BoolVar(&insecure, "insecure", false, "Allow insecure connections")
	createProxyCacheCmd.Flags().IntVar(&expiration, "expiration", 0, "Cache expiration in seconds")

	// delete-proxy-cache flags
	deleteProxyCacheCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func initOrgRobotFlags() {
	// robot flags
	orgRobotCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	if err := orgRobotCmd.MarkFlagRequired("robot"); err != nil {
		fmt.Printf("Error marking robot flag as required: %v\n", err)
		os.Exit(1)
	}

	// create-robot flags
	createRobotCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	if err := createRobotCmd.MarkFlagRequired("robot"); err != nil {
		fmt.Printf("Error marking robot flag as required: %v\n", err)
		os.Exit(1)
	}
	createRobotCmd.Flags().StringVar(&description, "description", "", "Robot description")

	// delete-robot flags
	deleteRobotCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	if err := deleteRobotCmd.MarkFlagRequired("robot"); err != nil {
		fmt.Printf("Error marking robot flag as required: %v\n", err)
		os.Exit(1)
	}
	deleteRobotCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	// regenerate-robot flags
	regenerateRobotCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	if err := regenerateRobotCmd.MarkFlagRequired("robot"); err != nil {
		fmt.Printf("Error marking robot flag as required: %v\n", err)
		os.Exit(1)
	}

	// robot-permissions flags
	orgRobotPermissionsCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	if err := orgRobotPermissionsCmd.MarkFlagRequired("robot"); err != nil {
		fmt.Printf("Error marking robot flag as required: %v\n", err)
		os.Exit(1)
	}

	// set-robot-permission flags
	setRobotPermissionCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	if err := setRobotPermissionCmd.MarkFlagRequired("robot"); err != nil {
		fmt.Printf("Error marking robot flag as required: %v\n", err)
		os.Exit(1)
	}
	setRobotPermissionCmd.Flags().StringVar(&repository, "repository", "", "Repository name")
	if err := setRobotPermissionCmd.MarkFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}
	setRobotPermissionCmd.Flags().StringVar(&role, "role", "", "Permission role")
	if err := setRobotPermissionCmd.MarkFlagRequired("role"); err != nil {
		fmt.Printf("Error marking role flag as required: %v\n", err)
		os.Exit(1)
	}

	// remove-robot-permission flags
	removeRobotPermissionCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	if err := removeRobotPermissionCmd.MarkFlagRequired("robot"); err != nil {
		fmt.Printf("Error marking robot flag as required: %v\n", err)
		os.Exit(1)
	}
	removeRobotPermissionCmd.Flags().StringVar(&repository, "repository", "", "Repository name")
	if err := removeRobotPermissionCmd.MarkFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}
	removeRobotPermissionCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm removal")
}

func initOrgApplicationFlags() {
	// application flags
	applicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	if err := applicationCmd.MarkFlagRequired("client-id"); err != nil {
		fmt.Printf("Error marking client-id flag as required: %v\n", err)
		os.Exit(1)
	}

	// create-application flags
	createApplicationCmd.Flags().StringVar(&appName, "name", "", "Application name")
	if err := createApplicationCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("Error marking name flag as required: %v\n", err)
		os.Exit(1)
	}
	createApplicationCmd.Flags().StringVar(&description, "description", "", "Application description")
	createApplicationCmd.Flags().StringVar(&applicationURI, "application-uri", "", "Application URI")
	createApplicationCmd.Flags().StringVar(&redirectURI, "redirect-uri", "", "Redirect URI")

	// update-application flags
	updateApplicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	if err := updateApplicationCmd.MarkFlagRequired("client-id"); err != nil {
		fmt.Printf("Error marking client-id flag as required: %v\n", err)
		os.Exit(1)
	}
	updateApplicationCmd.Flags().StringVar(&appName, "name", "", "Application name")
	if err := updateApplicationCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("Error marking name flag as required: %v\n", err)
		os.Exit(1)
	}
	updateApplicationCmd.Flags().StringVar(&description, "description", "", "Application description")
	updateApplicationCmd.Flags().StringVar(&applicationURI, "application-uri", "", "Application URI")
	updateApplicationCmd.Flags().StringVar(&redirectURI, "redirect-uri", "", "Redirect URI")

	// delete-application flags
	deleteApplicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	if err := deleteApplicationCmd.MarkFlagRequired("client-id"); err != nil {
		fmt.Printf("Error marking client-id flag as required: %v\n", err)
		os.Exit(1)
	}
	deleteApplicationCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	// reset-application-secret flags
	resetApplicationSecretCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	if err := resetApplicationSecretCmd.MarkFlagRequired("client-id"); err != nil {
		fmt.Printf("Error marking client-id flag as required: %v\n", err)
		os.Exit(1)
	}
}

func initOrgMarketplaceFlags() {
	// create-marketplace-subscription flags
	createMarketplaceSubscriptionCmd.Flags().StringVar(&sku, "sku", "", "Subscription SKU")
	if err := createMarketplaceSubscriptionCmd.MarkFlagRequired("sku"); err != nil {
		fmt.Printf("Error marking sku flag as required: %v\n", err)
		os.Exit(1)
	}
	createMarketplaceSubscriptionCmd.Flags().IntVar(&quantity, "quantity", 0, "Subscription quantity")

	// delete-marketplace-subscription flags
	deleteMarketplaceSubscriptionCmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "Subscription ID")
	if err := deleteMarketplaceSubscriptionCmd.MarkFlagRequired("subscription-id"); err != nil {
		fmt.Printf("Error marking subscription-id flag as required: %v\n", err)
		os.Exit(1)
	}
	deleteMarketplaceSubscriptionCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	// batch-remove-subscriptions flags
	batchRemoveSubscriptionsCmd.Flags().StringSliceVar(&subscriptionIDs, "subscription-ids", nil, "Comma-separated subscription IDs")
	if err := batchRemoveSubscriptionsCmd.MarkFlagRequired("subscription-ids"); err != nil {
		fmt.Printf("Error marking subscription-ids flag as required: %v\n", err)
		os.Exit(1)
	}
	batchRemoveSubscriptionsCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm removal")
}

func initOrgQuotaFlags() {
	// create-quota flags
	createQuotaCmd.Flags().Int64Var(&limitBytes, "limit-bytes", 0, "Quota limit in bytes")
	if err := createQuotaCmd.MarkFlagRequired("limit-bytes"); err != nil {
		fmt.Printf("Error marking limit-bytes flag as required: %v\n", err)
		os.Exit(1)
	}

	// update-quota flags
	updateQuotaCmd.Flags().Int64Var(&limitBytes, "limit-bytes", 0, "Quota limit in bytes")
	if err := updateQuotaCmd.MarkFlagRequired("limit-bytes"); err != nil {
		fmt.Printf("Error marking limit-bytes flag as required: %v\n", err)
		os.Exit(1)
	}

	// delete-quota flags
	deleteQuotaCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func initOrgAutoPruneFlags() {
	// auto-prune-policy flags
	autoPrunePolicyCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	if err := autoPrunePolicyCmd.MarkFlagRequired("policy-uuid"); err != nil {
		fmt.Printf("Error marking policy-uuid flag as required: %v\n", err)
		os.Exit(1)
	}

	// create-auto-prune flags
	createAutoPruneCmd.Flags().StringVar(&method, "method", "", "Prune method")
	if err := createAutoPruneCmd.MarkFlagRequired("method"); err != nil {
		fmt.Printf("Error marking method flag as required: %v\n", err)
		os.Exit(1)
	}
	createAutoPruneCmd.Flags().IntVar(&pruneValue, "value", 0, "Prune value")
	if err := createAutoPruneCmd.MarkFlagRequired("value"); err != nil {
		fmt.Printf("Error marking value flag as required: %v\n", err)
		os.Exit(1)
	}
	createAutoPruneCmd.Flags().StringVar(&tagPattern, "tag-pattern", "", "Tag pattern to match")

	// update-auto-prune flags
	updateAutoPruneCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	if err := updateAutoPruneCmd.MarkFlagRequired("policy-uuid"); err != nil {
		fmt.Printf("Error marking policy-uuid flag as required: %v\n", err)
		os.Exit(1)
	}
	updateAutoPruneCmd.Flags().StringVar(&method, "method", "", "Prune method")
	if err := updateAutoPruneCmd.MarkFlagRequired("method"); err != nil {
		fmt.Printf("Error marking method flag as required: %v\n", err)
		os.Exit(1)
	}
	updateAutoPruneCmd.Flags().IntVar(&pruneValue, "value", 0, "Prune value")
	if err := updateAutoPruneCmd.MarkFlagRequired("value"); err != nil {
		fmt.Printf("Error marking value flag as required: %v\n", err)
		os.Exit(1)
	}
	updateAutoPruneCmd.Flags().StringVar(&tagPattern, "tag-pattern", "", "Tag pattern to match")

	// delete-auto-prune flags
	deleteAutoPruneCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	if err := deleteAutoPruneCmd.MarkFlagRequired("policy-uuid"); err != nil {
		fmt.Printf("Error marking policy-uuid flag as required: %v\n", err)
		os.Exit(1)
	}
	deleteAutoPruneCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func initOrgInviteFlags() {
	// invite-member flags
	inviteMemberCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	if err := inviteMemberCmd.MarkFlagRequired("team"); err != nil {
		fmt.Printf("Error marking team flag as required: %v\n", err)
		os.Exit(1)
	}
	inviteMemberCmd.Flags().StringVar(&email, "email", "", "Email address")
	if err := inviteMemberCmd.MarkFlagRequired("email"); err != nil {
		fmt.Printf("Error marking email flag as required: %v\n", err)
		os.Exit(1)
	}

	// cancel-invite flags
	cancelInviteCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	if err := cancelInviteCmd.MarkFlagRequired("team"); err != nil {
		fmt.Printf("Error marking team flag as required: %v\n", err)
		os.Exit(1)
	}
	cancelInviteCmd.Flags().StringVar(&email, "email", "", "Email address")
	if err := cancelInviteCmd.MarkFlagRequired("email"); err != nil {
		fmt.Printf("Error marking email flag as required: %v\n", err)
		os.Exit(1)
	}
	cancelInviteCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm cancellation")
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
