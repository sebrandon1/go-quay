package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
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
	Use:   "list",
	Short: "List all teams in an organization",
	Long:  `List all teams within the specified organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		teams, err := client.GetTeams(teamCmdOrgname)
		if err != nil {
			fmt.Printf("Error getting teams: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Teams in organization '%s':\n", teamCmdOrgname)
		printJSON(teams)
	},
}

// Team Info
var teamCmdInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get team details",
	Long:  `Get detailed information about a specific team.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		team, err := client.GetTeam(teamCmdOrgname, teamCmdName)
		if err != nil {
			fmt.Printf("Error getting team: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Team: %s/%s\n", teamCmdOrgname, teamCmdName)
		printJSON(team)
	},
}

// Team Create
var teamCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new team",
	Long: `Create a new team within an organization.

Roles:
  - member: Inherits default permissions
  - creator: Can create new repositories
  - admin: Full administrative access`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		team, err := client.CreateTeam(teamCmdOrgname, teamCmdName, teamCmdDescription, teamCmdRole)
		if err != nil {
			fmt.Printf("Error creating team: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created team: %s/%s\n", teamCmdOrgname, teamCmdName)
		printJSON(team)
	},
}

// Team Update
var teamUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update team settings",
	Long:  `Update team description and role.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		team, err := client.UpdateTeam(teamCmdOrgname, teamCmdName, teamCmdDescription, teamCmdRole)
		if err != nil {
			fmt.Printf("Error updating team: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated team: %s/%s\n", teamCmdOrgname, teamCmdName)
		printJSON(team)
	},
}

// Team Delete
var teamDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a team",
	Long:  `Delete a team from an organization. This action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmTeamDelete {
			fmt.Printf("Are you sure you want to delete team '%s/%s'? This action cannot be undone.\n", teamCmdOrgname, teamCmdName)
			fmt.Println("Use --confirm to proceed with deletion.")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.DeleteTeam(teamCmdOrgname, teamCmdName)
		if err != nil {
			fmt.Printf("Error deleting team: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted team: %s/%s\n", teamCmdOrgname, teamCmdName)
	},
}

// Team Members
var teamCmdMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "List team members",
	Long:  `List all members of a specific team.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		members, err := client.GetTeamMembers(teamCmdOrgname, teamCmdName)
		if err != nil {
			fmt.Printf("Error getting team members: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Members of team '%s/%s':\n", teamCmdOrgname, teamCmdName)
		printJSON(members)
	},
}

// Team Add Member
var teamAddMemberCmd = &cobra.Command{
	Use:   "add-member",
	Short: "Add a member to a team",
	Long:  `Add a user or robot account to a team.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.AddTeamMember(teamCmdOrgname, teamCmdName, teamCmdMemberName)
		if err != nil {
			fmt.Printf("Error adding team member: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully added '%s' to team '%s/%s'\n", teamCmdMemberName, teamCmdOrgname, teamCmdName)
	},
}

// Team Remove Member
var teamRemoveMemberCmd = &cobra.Command{
	Use:   "remove-member",
	Short: "Remove a member from a team",
	Long:  `Remove a user or robot account from a team.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmTeamMemberDel {
			fmt.Printf("Are you sure you want to remove '%s' from team '%s/%s'?\n", teamCmdMemberName, teamCmdOrgname, teamCmdName)
			fmt.Println("Use --confirm to proceed with removal.")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.RemoveTeamMember(teamCmdOrgname, teamCmdName, teamCmdMemberName)
		if err != nil {
			fmt.Printf("Error removing team member: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully removed '%s' from team '%s/%s'\n", teamCmdMemberName, teamCmdOrgname, teamCmdName)
	},
}

// Team Permissions
var teamPermissionsCmd = &cobra.Command{
	Use:   "permissions",
	Short: "List team repository permissions",
	Long:  `List all repository permissions for a team.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		permissions, err := client.GetTeamPermissions(teamCmdOrgname, teamCmdName)
		if err != nil {
			fmt.Printf("Error getting team permissions: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Repository permissions for team '%s/%s':\n", teamCmdOrgname, teamCmdName)
		printJSON(permissions)
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
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.SetTeamRepositoryPermission(teamCmdOrgname, teamCmdName, teamCmdRepository, teamCmdPermissionRole)
		if err != nil {
			fmt.Printf("Error setting team permission: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set '%s' permission for team '%s/%s' on repository '%s'\n",
			teamCmdPermissionRole, teamCmdOrgname, teamCmdName, teamCmdRepository)
	},
}

// Team Remove Permission
var teamRemovePermissionCmd = &cobra.Command{
	Use:   "remove-permission",
	Short: "Remove repository permission from team",
	Long:  `Remove repository permission from a team.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmTeamPermDel {
			fmt.Printf("Are you sure you want to remove permissions for team '%s/%s' on repository '%s'?\n",
				teamCmdOrgname, teamCmdName, teamCmdRepository)
			fmt.Println("Use --confirm to proceed with removal.")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.RemoveTeamRepositoryPermission(teamCmdOrgname, teamCmdName, teamCmdRepository)
		if err != nil {
			fmt.Printf("Error removing team permission: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully removed permissions for team '%s/%s' on repository '%s'\n",
			teamCmdOrgname, teamCmdName, teamCmdRepository)
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
	teamCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")
	teamCmd.PersistentFlags().StringVarP(&teamCmdOrgname, "organization", "o", "", "Organization name")
	markTeamFlagRequired(teamCmd.MarkPersistentFlagRequired("token"))
	markTeamFlagRequired(teamCmd.MarkPersistentFlagRequired("organization"))
}

func initTeamCmdFlags() {
	// Info command flags
	teamCmdInfoCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	markTeamFlagRequired(teamCmdInfoCmd.MarkFlagRequired("name"))

	// Create command flags
	teamCreateCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamCreateCmd.Flags().StringVarP(&teamCmdDescription, "description", "d", "", "Team description")
	teamCreateCmd.Flags().StringVarP(&teamCmdRole, "role", "r", "member", "Team role (member, creator, admin)")
	markTeamFlagRequired(teamCreateCmd.MarkFlagRequired("name"))

	// Update command flags
	teamUpdateCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamUpdateCmd.Flags().StringVarP(&teamCmdDescription, "description", "d", "", "Team description")
	teamUpdateCmd.Flags().StringVarP(&teamCmdRole, "role", "r", "", "Team role (member, creator, admin)")
	markTeamFlagRequired(teamUpdateCmd.MarkFlagRequired("name"))

	// Delete command flags
	teamDeleteCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamDeleteCmd.Flags().BoolVar(&confirmTeamDelete, "confirm", false, "Confirm team deletion")
	markTeamFlagRequired(teamDeleteCmd.MarkFlagRequired("name"))
}

func initTeamMemberCmdFlags() {
	// Members command flags
	teamCmdMembersCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	markTeamFlagRequired(teamCmdMembersCmd.MarkFlagRequired("name"))

	// Add member command flags
	teamAddMemberCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamAddMemberCmd.Flags().StringVarP(&teamCmdMemberName, "member", "m", "", "Member name (username or robot)")
	markTeamFlagRequired(teamAddMemberCmd.MarkFlagRequired("name"))
	markTeamFlagRequired(teamAddMemberCmd.MarkFlagRequired("member"))

	// Remove member command flags
	teamRemoveMemberCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamRemoveMemberCmd.Flags().StringVarP(&teamCmdMemberName, "member", "m", "", "Member name (username or robot)")
	teamRemoveMemberCmd.Flags().BoolVar(&confirmTeamMemberDel, "confirm", false, "Confirm member removal")
	markTeamFlagRequired(teamRemoveMemberCmd.MarkFlagRequired("name"))
	markTeamFlagRequired(teamRemoveMemberCmd.MarkFlagRequired("member"))
}

func initTeamPermissionCmdFlags() {
	// Permissions command flags
	teamPermissionsCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	markTeamFlagRequired(teamPermissionsCmd.MarkFlagRequired("name"))

	// Set permission command flags
	teamSetPermissionCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamSetPermissionCmd.Flags().StringVarP(&teamCmdRepository, "repository", "R", "", "Repository name")
	teamSetPermissionCmd.Flags().StringVarP(&teamCmdPermissionRole, "role", "r", "", "Permission role (read, write, admin)")
	markTeamFlagRequired(teamSetPermissionCmd.MarkFlagRequired("name"))
	markTeamFlagRequired(teamSetPermissionCmd.MarkFlagRequired("repository"))
	markTeamFlagRequired(teamSetPermissionCmd.MarkFlagRequired("role"))

	// Remove permission command flags
	teamRemovePermissionCmd.Flags().StringVarP(&teamCmdName, "name", "n", "", "Team name")
	teamRemovePermissionCmd.Flags().StringVarP(&teamCmdRepository, "repository", "R", "", "Repository name")
	teamRemovePermissionCmd.Flags().BoolVar(&confirmTeamPermDel, "confirm", false, "Confirm permission removal")
	markTeamFlagRequired(teamRemovePermissionCmd.MarkFlagRequired("name"))
	markTeamFlagRequired(teamRemovePermissionCmd.MarkFlagRequired("repository"))
}

func markTeamFlagRequired(err error) {
	if err != nil {
		fmt.Printf("Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
