package cmd

import (
	"fmt"

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
	Use:   subcmdInfo,
	Short: "Get organization information",
	Long:  `Get detailed information about an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		org, err := client.GetOrganization(orgName)
		if err != nil {
			return fmt.Errorf("getting organization: %w", err)
		}
		return printJSON(org)
	},
}

// Organization Members
var orgMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Get organization members",
	Long:  `Get list of all members in an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		members, err := client.GetOrganizationMembers(orgName)
		if err != nil {
			return fmt.Errorf("getting organization members: %w", err)
		}
		return printJSON(members)
	},
}

// Organization Teams
var orgTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Get organization teams",
	Long:  `Get list of all teams in an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		teams, err := client.GetTeams(orgName)
		if err != nil {
			return fmt.Errorf("getting organization teams: %w", err)
		}
		return printJSON(teams)
	},
}

// Team Info
var teamInfoCmd = &cobra.Command{
	Use:   "team",
	Short: "Get team information",
	Long:  `Get detailed information about a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		team, err := client.GetTeam(orgName, teamName)
		if err != nil {
			return fmt.Errorf("getting team: %w", err)
		}
		return printJSON(team)
	},
}

// Team Members
var teamMembersCmd = &cobra.Command{
	Use:   "team-members",
	Short: "Get team members",
	Long:  `Get list of all members in a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		members, err := client.GetTeamMembers(orgName, teamName)
		if err != nil {
			return fmt.Errorf("getting team members: %w", err)
		}
		return printJSON(members)
	},
}

// Organization Robots
var orgRobotsCmd = &cobra.Command{
	Use:   "robots",
	Short: "Get organization robots",
	Long:  `Get list of all robot accounts in an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		robots, err := client.GetRobotAccounts(orgName)
		if err != nil {
			return fmt.Errorf("getting organization robots: %w", err)
		}
		return printJSON(robots)
	},
}

// Organization Quota
var orgQuotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Get organization quota",
	Long:  `Get quota information for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		quota, err := client.GetQuota(orgName)
		if err != nil {
			return fmt.Errorf("getting organization quota: %w", err)
		}
		return printJSON(quota)
	},
}

// Auto-prune Policies
var autoPruneCmd = &cobra.Command{
	Use:   "auto-prune",
	Short: "Get auto-prune policies",
	Long:  `Get list of auto-prune policies for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		policies, err := client.GetAutoPrunePolicies(orgName)
		if err != nil {
			return fmt.Errorf("getting auto-prune policies: %w", err)
		}
		return printJSON(policies)
	},
}

// Organization Applications
var orgApplicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Get organization applications",
	Long:  `Get list of OAuth applications for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		applications, err := client.GetApplications(orgName)
		if err != nil {
			return fmt.Errorf("getting organization applications: %w", err)
		}
		return printJSON(applications)
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
		org, err := client.CreateOrganization(orgName, email)
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
		org, err := client.UpdateOrganization(orgName, email)
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
		err = client.DeleteOrganization(orgName)
		if err != nil {
			return fmt.Errorf("deleting organization: %w", err)
		}
		fmt.Println("Organization deleted successfully")
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
		err = client.AddOrganizationMember(orgName, memberName)
		if err != nil {
			return fmt.Errorf("adding organization member: %w", err)
		}
		fmt.Println("Member added successfully")
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
		err = client.RemoveOrganizationMember(orgName, memberName)
		if err != nil {
			return fmt.Errorf("removing organization member: %w", err)
		}
		fmt.Println("Member removed successfully")
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
		member, err := client.GetOrganizationMember(orgName, memberName)
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
		collaborators, err := client.GetOrganizationCollaborators(orgName)
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
		repos, err := client.GetOrganizationRepositories(orgName)
		if err != nil {
			return fmt.Errorf("getting organization repositories: %w", err)
		}
		return printJSON(repos)
	},
}

// Get Proxy Cache Config
var proxyCacheCmd = &cobra.Command{
	Use:   "proxy-cache",
	Short: "Get proxy cache configuration",
	Long:  `Get proxy cache configuration for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		config, err := client.GetProxyCacheConfig(orgName)
		if err != nil {
			return fmt.Errorf("getting proxy cache config: %w", err)
		}
		return printJSON(config)
	},
}

// Create Proxy Cache Config
var createProxyCacheCmd = &cobra.Command{
	Use:   "create-proxy-cache",
	Short: "Create proxy cache configuration",
	Long:  `Create proxy cache configuration for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		config, err := client.CreateProxyCacheConfig(orgName, upstreamRegistry, insecure, expiration)
		if err != nil {
			return fmt.Errorf("creating proxy cache config: %w", err)
		}
		return printJSON(config)
	},
}

// Delete Proxy Cache Config
var deleteProxyCacheCmd = &cobra.Command{
	Use:   "delete-proxy-cache",
	Short: "Delete proxy cache configuration",
	Long:  `Delete proxy cache configuration for an organization. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete proxy cache config")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteProxyCacheConfig(orgName)
		if err != nil {
			return fmt.Errorf("deleting proxy cache config: %w", err)
		}
		fmt.Println("Proxy cache config deleted successfully")
		return nil
	},
}

// Get Robot Account
var orgRobotCmd = &cobra.Command{
	Use:   "robot",
	Short: "Get robot account information",
	Long:  `Get detailed information about a specific robot account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		robot, err := client.GetRobotAccount(orgName, robotShortname)
		if err != nil {
			return fmt.Errorf("getting robot account: %w", err)
		}
		return printJSON(robot)
	},
}

// Create Robot Account
var createRobotCmd = &cobra.Command{
	Use:   "create-robot",
	Short: "Create a robot account",
	Long:  `Create a new robot account in an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		robot, err := client.CreateRobotAccount(orgName, robotShortname, description, nil)
		if err != nil {
			return fmt.Errorf("creating robot account: %w", err)
		}
		return printJSON(robot)
	},
}

// Delete Robot Account
var deleteRobotCmd = &cobra.Command{
	Use:   "delete-robot",
	Short: "Delete a robot account",
	Long:  `Delete a robot account from an organization. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete a robot account")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteRobotAccount(orgName, robotShortname)
		if err != nil {
			return fmt.Errorf("deleting robot account: %w", err)
		}
		fmt.Println("Robot account deleted successfully")
		return nil
	},
}

// Regenerate Robot Token
var regenerateRobotCmd = &cobra.Command{
	Use:   "regenerate-robot",
	Short: "Regenerate a robot account token",
	Long:  `Regenerate the token for a robot account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		robot, err := client.RegenerateRobotToken(orgName, robotShortname)
		if err != nil {
			return fmt.Errorf("regenerating robot token: %w", err)
		}
		return printJSON(robot)
	},
}

// Get Robot Permissions
var orgRobotPermissionsCmd = &cobra.Command{
	Use:   "robot-permissions",
	Short: "Get robot account permissions",
	Long:  `Get permissions for a specific robot account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		perms, err := client.GetRobotPermissions(orgName, robotShortname)
		if err != nil {
			return fmt.Errorf("getting robot permissions: %w", err)
		}
		return printJSON(perms)
	},
}

// Set Robot Repository Permission
var setRobotPermissionCmd = &cobra.Command{
	Use:   "set-robot-permission",
	Short: "Set robot repository permission",
	Long:  `Set a robot account's permission on a repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.SetRobotRepositoryPermission(orgName, robotShortname, repository, role)
		if err != nil {
			return fmt.Errorf("setting robot permission: %w", err)
		}
		fmt.Println("Robot permission set successfully")
		return nil
	},
}

// Remove Robot Repository Permission
var removeRobotPermissionCmd = &cobra.Command{
	Use:   "remove-robot-permission",
	Short: "Remove robot repository permission",
	Long:  `Remove a robot account's permission on a repository. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to remove a robot permission")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.RemoveRobotRepositoryPermission(orgName, robotShortname, repository)
		if err != nil {
			return fmt.Errorf("removing robot permission: %w", err)
		}
		fmt.Println("Robot permission removed successfully")
		return nil
	},
}

// Org Robot Federation Get
var orgRobotFederationGetCmd = &cobra.Command{
	Use:   "robot-federation-get",
	Short: "Get organization robot federation configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		federation, err := client.GetRobotFederation(orgName, robotShortname)
		if err != nil {
			return fmt.Errorf("getting robot federation: %w", err)
		}

		return printJSON(federation)
	},
}

// Org Robot Federation Create
var orgRobotFederationCreateCmd = &cobra.Command{
	Use:   "robot-federation-create",
	Short: "Create or update organization robot federation configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		configs := []lib.RobotFederationConfig{
			{Issuer: federationIssuer, Subject: federationSubject},
		}

		err = client.CreateRobotFederation(orgName, robotShortname, configs)
		if err != nil {
			return fmt.Errorf("creating robot federation: %w", err)
		}

		fmt.Printf("Successfully configured federation for robot %s in org %s\n", robotShortname, orgName)
		return nil
	},
}

// Org Robot Federation Delete
var orgRobotFederationDeleteCmd = &cobra.Command{
	Use:   "robot-federation-delete",
	Short: "Delete organization robot federation configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteRobotFederation(orgName, robotShortname)
		if err != nil {
			return fmt.Errorf("deleting robot federation: %w", err)
		}

		fmt.Printf("Successfully deleted federation for robot %s in org %s\n", robotShortname, orgName)
		return nil
	},
}

// Get Application
var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Get application information",
	Long:  `Get detailed information about a specific OAuth application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		app, err := client.GetApplication(orgName, clientID)
		if err != nil {
			return fmt.Errorf("getting application: %w", err)
		}
		return printJSON(app)
	},
}

// Create Application
var createApplicationCmd = &cobra.Command{
	Use:   "create-application",
	Short: "Create an OAuth application",
	Long:  `Create a new OAuth application for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		app, err := client.CreateApplication(orgName, appName, description, applicationURI, redirectURI)
		if err != nil {
			return fmt.Errorf("creating application: %w", err)
		}
		return printJSON(app)
	},
}

// Update Application
var updateApplicationCmd = &cobra.Command{
	Use:   "update-application",
	Short: "Update an OAuth application",
	Long:  `Update an existing OAuth application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		app, err := client.UpdateApplication(orgName, clientID, appName, description, applicationURI, redirectURI)
		if err != nil {
			return fmt.Errorf("updating application: %w", err)
		}
		return printJSON(app)
	},
}

// Delete Application
var deleteApplicationCmd = &cobra.Command{
	Use:   "delete-application",
	Short: "Delete an OAuth application",
	Long:  `Delete an OAuth application. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete an application")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteApplication(orgName, clientID)
		if err != nil {
			return fmt.Errorf("deleting application: %w", err)
		}
		fmt.Println("Application deleted successfully")
		return nil
	},
}

// Reset Application Client Secret
var resetApplicationSecretCmd = &cobra.Command{
	Use:   "reset-application-secret",
	Short: "Reset application client secret",
	Long:  `Reset the client secret for an OAuth application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		app, err := client.ResetApplicationClientSecret(orgName, clientID)
		if err != nil {
			return fmt.Errorf("resetting application secret: %w", err)
		}
		return printJSON(app)
	},
}

// Get Organization Marketplace
var marketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Get organization marketplace information",
	Long:  `Get marketplace information for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		marketplace, err := client.GetOrganizationMarketplace(orgName)
		if err != nil {
			return fmt.Errorf("getting marketplace info: %w", err)
		}
		return printJSON(marketplace)
	},
}

// Create Marketplace Subscription
var createMarketplaceSubscriptionCmd = &cobra.Command{
	Use:   "create-marketplace-subscription",
	Short: "Create a marketplace subscription",
	Long:  `Create a marketplace subscription for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.CreateOrganizationMarketplaceSubscription(orgName, &lib.MarketplaceSubscriptionRequest{SKU: sku, Quantity: quantity})
		if err != nil {
			return fmt.Errorf("creating marketplace subscription: %w", err)
		}
		fmt.Println("Marketplace subscription created successfully")
		return nil
	},
}

// Delete Marketplace Subscription
var deleteMarketplaceSubscriptionCmd = &cobra.Command{
	Use:   "delete-marketplace-subscription",
	Short: "Delete a marketplace subscription",
	Long:  `Delete a marketplace subscription. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete a marketplace subscription")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteOrganizationMarketplaceSubscription(orgName, subscriptionID)
		if err != nil {
			return fmt.Errorf("deleting marketplace subscription: %w", err)
		}
		fmt.Println("Marketplace subscription deleted successfully")
		return nil
	},
}

// Batch Remove Marketplace Subscriptions
var batchRemoveSubscriptionsCmd = &cobra.Command{
	Use:   "batch-remove-subscriptions",
	Short: "Batch remove marketplace subscriptions",
	Long:  `Remove multiple marketplace subscriptions at once. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to batch remove subscriptions")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.BatchRemoveOrganizationMarketplaceSubscriptions(orgName, subscriptionIDs)
		if err != nil {
			return fmt.Errorf("batch removing subscriptions: %w", err)
		}
		fmt.Println("Marketplace subscriptions removed successfully")
		return nil
	},
}

// Create Quota
var createQuotaCmd = &cobra.Command{
	Use:   "create-quota",
	Short: "Create organization quota",
	Long:  `Create a quota for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		quota, err := client.CreateQuota(orgName, limitBytes)
		if err != nil {
			return fmt.Errorf("creating quota: %w", err)
		}
		return printJSON(quota)
	},
}

// Update Quota
var updateQuotaCmd = &cobra.Command{
	Use:   "update-quota",
	Short: "Update organization quota",
	Long:  `Update the quota for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		quota, err := client.UpdateQuota(orgName, limitBytes)
		if err != nil {
			return fmt.Errorf("updating quota: %w", err)
		}
		return printJSON(quota)
	},
}

// Delete Quota
var deleteQuotaCmd = &cobra.Command{
	Use:   "delete-quota",
	Short: "Delete organization quota",
	Long:  `Delete the quota for an organization. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete a quota")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteQuota(orgName)
		if err != nil {
			return fmt.Errorf("deleting quota: %w", err)
		}
		fmt.Println("Quota deleted successfully")
		return nil
	},
}

// Get Auto-Prune Policy
var autoPrunePolicyCmd = &cobra.Command{
	Use:   "auto-prune-policy",
	Short: "Get a specific auto-prune policy",
	Long:  `Get detailed information about a specific auto-prune policy.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		policy, err := client.GetAutoPrunePolicy(orgName, policyUUID)
		if err != nil {
			return fmt.Errorf("getting auto-prune policy: %w", err)
		}
		return printJSON(policy)
	},
}

// Create Auto-Prune Policy
var createAutoPruneCmd = &cobra.Command{
	Use:   "create-auto-prune",
	Short: "Create an auto-prune policy",
	Long:  `Create an auto-prune policy for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		policy, err := client.CreateAutoPrunePolicy(orgName, method, pruneValue, tagPattern)
		if err != nil {
			return fmt.Errorf("creating auto-prune policy: %w", err)
		}
		return printJSON(policy)
	},
}

// Update Auto-Prune Policy
var updateAutoPruneCmd = &cobra.Command{
	Use:   "update-auto-prune",
	Short: "Update an auto-prune policy",
	Long:  `Update an existing auto-prune policy.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		policy, err := client.UpdateAutoPrunePolicy(orgName, policyUUID, method, pruneValue, tagPattern)
		if err != nil {
			return fmt.Errorf("updating auto-prune policy: %w", err)
		}
		return printJSON(policy)
	},
}

// Delete Auto-Prune Policy
var deleteAutoPruneCmd = &cobra.Command{
	Use:   "delete-auto-prune",
	Short: "Delete an auto-prune policy",
	Long:  `Delete an auto-prune policy. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete an auto-prune policy")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteAutoPrunePolicy(orgName, policyUUID)
		if err != nil {
			return fmt.Errorf("deleting auto-prune policy: %w", err)
		}
		fmt.Println("Auto-prune policy deleted successfully")
		return nil
	},
}

// Invite Team Member
var inviteMemberCmd = &cobra.Command{
	Use:   "invite-member",
	Short: "Invite a member to a team",
	Long:  `Invite a member to a team by email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.InviteTeamMember(orgName, teamName, email)
		if err != nil {
			return fmt.Errorf("inviting team member: %w", err)
		}
		fmt.Println("Team member invited successfully")
		return nil
	},
}

// Cancel Team Invite
var cancelInviteCmd = &cobra.Command{
	Use:   "cancel-invite",
	Short: "Cancel a team member invitation",
	Long:  `Cancel a pending team member invitation. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to cancel an invite")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteTeamInvite(orgName, teamName, email)
		if err != nil {
			return fmt.Errorf("canceling team invite: %w", err)
		}
		fmt.Println("Team invite canceled successfully")
		return nil
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
	organizationCmd.AddCommand(orgRobotFederationGetCmd)
	organizationCmd.AddCommand(orgRobotFederationCreateCmd)
	organizationCmd.AddCommand(orgRobotFederationDeleteCmd)
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
	_ = organizationCmd.MarkPersistentFlagRequired("organization")
}

func initOrgMemberFlags() {
	teamInfoCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = teamInfoCmd.MarkFlagRequired("team")

	teamMembersCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = teamMembersCmd.MarkFlagRequired("team")

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

func initOrgProxyCacheFlags() {
	createProxyCacheCmd.Flags().StringVar(&upstreamRegistry, "upstream-registry", "", "Upstream registry URL")
	_ = createProxyCacheCmd.MarkFlagRequired("upstream-registry")
	createProxyCacheCmd.Flags().BoolVar(&insecure, "insecure", false, "Allow insecure connections")
	createProxyCacheCmd.Flags().IntVar(&expiration, "expiration", 0, "Cache expiration in seconds")

	deleteProxyCacheCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func initOrgRobotFlags() {
	orgRobotCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = orgRobotCmd.MarkFlagRequired("robot")

	createRobotCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = createRobotCmd.MarkFlagRequired("robot")
	createRobotCmd.Flags().StringVar(&description, "description", "", "Robot description")

	deleteRobotCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = deleteRobotCmd.MarkFlagRequired("robot")
	deleteRobotCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	regenerateRobotCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = regenerateRobotCmd.MarkFlagRequired("robot")

	orgRobotPermissionsCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = orgRobotPermissionsCmd.MarkFlagRequired("robot")

	setRobotPermissionCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = setRobotPermissionCmd.MarkFlagRequired("robot")
	setRobotPermissionCmd.Flags().StringVar(&repository, "repository", "", "Repository name")
	_ = setRobotPermissionCmd.MarkFlagRequired("repository")
	setRobotPermissionCmd.Flags().StringVar(&role, "role", "", "Permission role")
	_ = setRobotPermissionCmd.MarkFlagRequired("role")

	removeRobotPermissionCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = removeRobotPermissionCmd.MarkFlagRequired("robot")
	removeRobotPermissionCmd.Flags().StringVar(&repository, "repository", "", "Repository name")
	_ = removeRobotPermissionCmd.MarkFlagRequired("repository")
	removeRobotPermissionCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm removal")

	orgRobotFederationGetCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = orgRobotFederationGetCmd.MarkFlagRequired("robot")

	orgRobotFederationCreateCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	orgRobotFederationCreateCmd.Flags().StringVar(&federationIssuer, "issuer", "", "Federation token issuer")
	orgRobotFederationCreateCmd.Flags().StringVar(&federationSubject, "subject", "", "Federation token subject")
	_ = orgRobotFederationCreateCmd.MarkFlagRequired("robot")
	_ = orgRobotFederationCreateCmd.MarkFlagRequired("issuer")
	_ = orgRobotFederationCreateCmd.MarkFlagRequired("subject")

	orgRobotFederationDeleteCmd.Flags().StringVar(&robotShortname, "robot", "", "Robot short name")
	_ = orgRobotFederationDeleteCmd.MarkFlagRequired("robot")
}

func initOrgApplicationFlags() {
	applicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	_ = applicationCmd.MarkFlagRequired("client-id")

	createApplicationCmd.Flags().StringVar(&appName, "name", "", "Application name")
	_ = createApplicationCmd.MarkFlagRequired("name")
	createApplicationCmd.Flags().StringVar(&description, "description", "", "Application description")
	createApplicationCmd.Flags().StringVar(&applicationURI, "application-uri", "", "Application URI")
	createApplicationCmd.Flags().StringVar(&redirectURI, "redirect-uri", "", "Redirect URI")

	updateApplicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	_ = updateApplicationCmd.MarkFlagRequired("client-id")
	updateApplicationCmd.Flags().StringVar(&appName, "name", "", "Application name")
	_ = updateApplicationCmd.MarkFlagRequired("name")
	updateApplicationCmd.Flags().StringVar(&description, "description", "", "Application description")
	updateApplicationCmd.Flags().StringVar(&applicationURI, "application-uri", "", "Application URI")
	updateApplicationCmd.Flags().StringVar(&redirectURI, "redirect-uri", "", "Redirect URI")

	deleteApplicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	_ = deleteApplicationCmd.MarkFlagRequired("client-id")
	deleteApplicationCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	resetApplicationSecretCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	_ = resetApplicationSecretCmd.MarkFlagRequired("client-id")
}

func initOrgMarketplaceFlags() {
	createMarketplaceSubscriptionCmd.Flags().StringVar(&sku, "sku", "", "Subscription SKU")
	_ = createMarketplaceSubscriptionCmd.MarkFlagRequired("sku")
	createMarketplaceSubscriptionCmd.Flags().IntVar(&quantity, "quantity", 0, "Subscription quantity")

	deleteMarketplaceSubscriptionCmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "Subscription ID")
	_ = deleteMarketplaceSubscriptionCmd.MarkFlagRequired("subscription-id")
	deleteMarketplaceSubscriptionCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	batchRemoveSubscriptionsCmd.Flags().StringSliceVar(&subscriptionIDs, "subscription-ids", nil, "Comma-separated subscription IDs")
	_ = batchRemoveSubscriptionsCmd.MarkFlagRequired("subscription-ids")
	batchRemoveSubscriptionsCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm removal")
}

func initOrgQuotaFlags() {
	createQuotaCmd.Flags().Int64Var(&limitBytes, "limit-bytes", 0, "Quota limit in bytes")
	_ = createQuotaCmd.MarkFlagRequired("limit-bytes")

	updateQuotaCmd.Flags().Int64Var(&limitBytes, "limit-bytes", 0, "Quota limit in bytes")
	_ = updateQuotaCmd.MarkFlagRequired("limit-bytes")

	deleteQuotaCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func initOrgAutoPruneFlags() {
	autoPrunePolicyCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	_ = autoPrunePolicyCmd.MarkFlagRequired("policy-uuid")

	createAutoPruneCmd.Flags().StringVar(&method, "method", "", "Prune method")
	_ = createAutoPruneCmd.MarkFlagRequired("method")
	createAutoPruneCmd.Flags().IntVar(&pruneValue, "value", 0, "Prune value")
	_ = createAutoPruneCmd.MarkFlagRequired("value")
	createAutoPruneCmd.Flags().StringVar(&tagPattern, "tag-pattern", "", "Tag pattern to match")

	updateAutoPruneCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	_ = updateAutoPruneCmd.MarkFlagRequired("policy-uuid")
	updateAutoPruneCmd.Flags().StringVar(&method, "method", "", "Prune method")
	_ = updateAutoPruneCmd.MarkFlagRequired("method")
	updateAutoPruneCmd.Flags().IntVar(&pruneValue, "value", 0, "Prune value")
	_ = updateAutoPruneCmd.MarkFlagRequired("value")
	updateAutoPruneCmd.Flags().StringVar(&tagPattern, "tag-pattern", "", "Tag pattern to match")

	deleteAutoPruneCmd.Flags().StringVar(&policyUUID, "policy-uuid", "", "Policy UUID")
	_ = deleteAutoPruneCmd.MarkFlagRequired("policy-uuid")
	deleteAutoPruneCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")
}

func initOrgInviteFlags() {
	inviteMemberCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = inviteMemberCmd.MarkFlagRequired("team")
	inviteMemberCmd.Flags().StringVar(&email, "email", "", "Email address")
	_ = inviteMemberCmd.MarkFlagRequired("email")

	cancelInviteCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = cancelInviteCmd.MarkFlagRequired("team")
	cancelInviteCmd.Flags().StringVar(&email, "email", "", "Email address")
	_ = cancelInviteCmd.MarkFlagRequired("email")
	cancelInviteCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm cancellation")
}
