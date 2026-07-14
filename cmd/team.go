package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	teamCmdOrgname        string
	teamCmdName           string
	teamCmdDescription    string
	teamCmdRole           string
	teamCmdMemberName     string
	teamCmdRepository     string
	teamCmdPermissionRole string
	confirmTeamDelete     bool
	confirmTeamMemberDel  bool
	confirmTeamPermDel    bool
)

// teamCmd represents the team command group
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Organization team management commands",
	Long: `Commands for managing teams within an organization.

Teams allow you to organize members and manage repository permissions at scale.

Available commands:
  list         - List all teams in an organization
  info         - Get team details
  create       - Create a new team
  update       - Update team settings
  delete       - Delete a team
  members      - List team members
  add-member   - Add a member to a team
  remove-member - Remove a member from a team
  permissions  - List team repository permissions
  set-permission - Set repository permission for team
  remove-permission - Remove repository permission from team`,
}

// Team List
var teamListCmd = &cobra.Command{
	Use:   subcmdList,
	Short: "List all teams in an organization",
	Long:  `List all teams within the specified organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		teams, err := client.GetTeams(teamCmdOrgname)
		if err != nil {
			return fmt.Errorf("getting teams: %w", err)
		}

		fmt.Printf("Teams in organization '%s':\n", teamCmdOrgname)
		return printJSON(teams)
	},
}

// Team Info
var teamCmdInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get team details",
	Long:  `Get detailed information about a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		team, err := client.GetTeam(teamCmdOrgname, teamCmdName)
		if err != nil {
			return fmt.Errorf("getting team: %w", err)
		}

		fmt.Printf("Team: %s/%s\n", teamCmdOrgname, teamCmdName)
		return printJSON(team)
	},
}

// Team Create
var teamCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
	Short: "Create a new team",
	Long: `Create a new team within an organization.

Roles:
  - member: Inherits default permissions
  - creator: Can create new repositories
  - admin: Full administrative access`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		team, err := client.CreateTeam(teamCmdOrgname, teamCmdName, teamCmdDescription, teamCmdRole)
		if err != nil {
			return fmt.Errorf("creating team: %w", err)
		}

		fmt.Printf("Successfully created team: %s/%s\n", teamCmdOrgname, teamCmdName)
		return printJSON(team)
	},
}

// Team Update
var teamUpdateCmd = &cobra.Command{
	Use:   subcmdUpdate,
	Short: "Update team settings",
	Long:  `Update team description and role.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		team, err := client.UpdateTeam(teamCmdOrgname, teamCmdName, teamCmdDescription, teamCmdRole)
		if err != nil {
			return fmt.Errorf("updating team: %w", err)
		}

		fmt.Printf("Successfully updated team: %s/%s\n", teamCmdOrgname, teamCmdName)
		return printJSON(team)
	},
}

// Team Delete
var teamDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a team",
	Long:  `Delete a team from an organization. This action cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmTeamDelete {
			return fmt.Errorf("are you sure you want to delete team '%s/%s'? This action cannot be undone.\nUse --confirm to proceed with deletion", teamCmdOrgname, teamCmdName)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteTeam(teamCmdOrgname, teamCmdName)
		if err != nil {
			return fmt.Errorf("deleting team: %w", err)
		}

		fmt.Printf("Successfully deleted team: %s/%s\n", teamCmdOrgname, teamCmdName)
		return nil
	},
}

// Team Members
var teamCmdMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "List team members",
	Long:  `List all members of a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		members, err := client.GetTeamMembers(teamCmdOrgname, teamCmdName)
		if err != nil {
			return fmt.Errorf("getting team members: %w", err)
		}

		fmt.Printf("Members of team '%s/%s':\n", teamCmdOrgname, teamCmdName)
		return printJSON(members)
	},
}

// Team Add Member
var teamAddMemberCmd = &cobra.Command{
	Use:   "add-member",
	Short: "Add a member to a team",
	Long:  `Add a user or robot account to a team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.AddTeamMember(teamCmdOrgname, teamCmdName, teamCmdMemberName)
		if err != nil {
			return fmt.Errorf("adding team member: %w", err)
		}

		fmt.Printf("Successfully added '%s' to team '%s/%s'\n", teamCmdMemberName, teamCmdOrgname, teamCmdName)
		return nil
	},
}

// Team Remove Member
var teamRemoveMemberCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "Remove a member from a team",
	Long:  `Remove a user or robot account from a team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmTeamMemberDel {
			return fmt.Errorf("are you sure you want to remove '%s' from team '%s/%s'?\nUse --confirm to proceed with removal", teamCmdMemberName, teamCmdOrgname, teamCmdName)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.RemoveTeamMember(teamCmdOrgname, teamCmdName, teamCmdMemberName)
		if err != nil {
			return fmt.Errorf("removing team member: %w", err)
		}

		fmt.Printf("Successfully removed '%s' from team '%s/%s'\n", teamCmdMemberName, teamCmdOrgname, teamCmdName)
		return nil
	},
}

// Team Permissions
var teamPermissionsCmd = &cobra.Command{
	Use:   subcmdPermissions,
	Short: "List team repository permissions",
	Long:  `List all repository permissions for a team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		permissions, err := client.GetTeamPermissions(teamCmdOrgname, teamCmdName)
		if err != nil {
			return fmt.Errorf("getting team permissions: %w", err)
		}

		fmt.Printf("Repository permissions for team '%s/%s':\n", teamCmdOrgname, teamCmdName)
		return printJSON(permissions)
	},
}

// Team Set Permission
var teamSetPermissionCmd = &cobra.Command{
	Use:   "set-permission",
	Short: "Set repository permission for team",
	Long: `Set repository permission for a team.

Permission roles:
  - read: Pull images
  - write: Pull and push images
  - admin: Full administrative access`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.SetTeamRepositoryPermission(teamCmdOrgname, teamCmdName, teamCmdRepository, teamCmdPermissionRole)
		if err != nil {
			return fmt.Errorf("setting team permission: %w", err)
		}

		fmt.Printf("Successfully set '%s' permission for team '%s/%s' on repository '%s'\n",
			teamCmdPermissionRole, teamCmdOrgname, teamCmdName, teamCmdRepository)
		return nil
	},
}

// Team Remove Permission
var teamRemovePermissionCmd = &cobra.Command{
	Use:   "remove-permission",
	Short: "Remove repository permission from team",
	Long:  `Remove repository permission from a team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmTeamPermDel {
			return fmt.Errorf("are you sure you want to remove permissions for team '%s/%s' on repository '%s'?\nUse --confirm to proceed with removal",
				teamCmdOrgname, teamCmdName, teamCmdRepository)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.RemoveTeamRepositoryPermission(teamCmdOrgname, teamCmdName, teamCmdRepository)
		if err != nil {
			return fmt.Errorf("removing team permission: %w", err)
		}

		fmt.Printf("Successfully removed permissions for team '%s/%s' on repository '%s'\n",
			teamCmdOrgname, teamCmdName, teamCmdRepository)
		return nil
	},
}

func init() {
	// Add subcommands to team command
	teamCmd.AddCommand(teamListCmd)
	teamCmd.AddCommand(teamCmdInfoCmd)
	teamCmd.AddCommand(teamCreateCmd)
	teamCmd.AddCommand(teamUpdateCmd)
	teamCmd.AddCommand(teamDeleteCmd)
	teamCmd.AddCommand(teamCmdMembersCmd)
	teamCmd.AddCommand(teamAddMemberCmd)
	teamCmd.AddCommand(teamRemoveMemberCmd)
	teamCmd.AddCommand(teamPermissionsCmd)
	teamCmd.AddCommand(teamSetPermissionCmd)
	teamCmd.AddCommand(teamRemovePermissionCmd)

	initTeamGlobalFlags()
	initTeamCmdFlags()
	initTeamMemberCmdFlags()
	initTeamPermissionCmdFlags()
}

func initTeamGlobalFlags() {
	teamCmd.PersistentFlags().StringVarP(&teamCmdOrgname, "organization", "o", "", "Organization name")
	_ = teamCmd.MarkPersistentFlagRequired("organization")
}

func initTeamCmdFlags() {
	// Info command flags
	teamCmdInfoCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	_ = teamCmdInfoCmd.MarkFlagRequired("name")

	// Create command flags
	teamCreateCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamCreateCmd.Flags().StringVarP(&teamCmdDescription, "description", "d", "", "Team description")
	teamCreateCmd.Flags().StringVarP(&teamCmdRole, "role", "r", "member", "Team role (member, creator, admin)")
	_ = teamCreateCmd.MarkFlagRequired("name")

	// Update command flags
	teamUpdateCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamUpdateCmd.Flags().StringVarP(&teamCmdDescription, "description", "d", "", "Team description")
	teamUpdateCmd.Flags().StringVarP(&teamCmdRole, "role", "r", "", "Team role (member, creator, admin)")
	_ = teamUpdateCmd.MarkFlagRequired("name")

	// Delete command flags
	teamDeleteCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamDeleteCmd.Flags().BoolVar(&confirmTeamDelete, "confirm", false, "Confirm team deletion")
	_ = teamDeleteCmd.MarkFlagRequired("name")
}

func initTeamMemberCmdFlags() {
	// Members command flags
	teamCmdMembersCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	_ = teamCmdMembersCmd.MarkFlagRequired("name")

	// Add member command flags
	teamAddMemberCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamAddMemberCmd.Flags().StringVarP(&teamCmdMemberName, "member", "m", "", "Member name (username or robot)")
	_ = teamAddMemberCmd.MarkFlagRequired("name")
	_ = teamAddMemberCmd.MarkFlagRequired("member")

	// Remove member command flags
	teamRemoveMemberCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamRemoveMemberCmd.Flags().StringVarP(&teamCmdMemberName, "member", "m", "", "Member name (username or robot)")
	teamRemoveMemberCmd.Flags().BoolVar(&confirmTeamMemberDel, "confirm", false, "Confirm member removal")
	_ = teamRemoveMemberCmd.MarkFlagRequired("name")
	_ = teamRemoveMemberCmd.MarkFlagRequired("member")
}

func initTeamPermissionCmdFlags() {
	// Permissions command flags
	teamPermissionsCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	_ = teamPermissionsCmd.MarkFlagRequired("name")

	// Set permission command flags
	teamSetPermissionCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamSetPermissionCmd.Flags().StringVarP(&teamCmdRepository, "repository", "R", "", "Repository name")
	teamSetPermissionCmd.Flags().StringVarP(&teamCmdPermissionRole, "role", "r", "", "Permission role (read, write, admin)")
	_ = teamSetPermissionCmd.MarkFlagRequired("name")
	_ = teamSetPermissionCmd.MarkFlagRequired("repository")
	_ = teamSetPermissionCmd.MarkFlagRequired("role")

	// Remove permission command flags
	teamRemovePermissionCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamRemovePermissionCmd.Flags().StringVarP(&teamCmdRepository, "repository", "R", "", "Repository name")
	teamRemovePermissionCmd.Flags().BoolVar(&confirmTeamPermDel, "confirm", false, "Confirm permission removal")
	_ = teamRemovePermissionCmd.MarkFlagRequired("name")
	_ = teamRemovePermissionCmd.MarkFlagRequired("repository")
}
