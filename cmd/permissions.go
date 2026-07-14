package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	permissionUser      string
	permissionRole      string
	permissionTeamName  string
	confirmPermDeletion bool
)

// permissionsCmd represents the permissions command group
var permissionsCmd = &cobra.Command{
	Use:   subcmdPermissions,
	Short: "Repository permissions management commands",
	Long: `Commands for managing repository permissions including listing, setting, and removing user/robot access.

Available commands:
  list   - List repository permissions
  set    - Set permission for a user/robot
  remove - Remove permission for a user/robot

Supported roles: read, write, admin`,
}

// Permissions List
var permListCmd = &cobra.Command{
	Use:   subcmdList,
	Short: "List repository permissions",
	Long:  `List all users and robots that have permissions on this repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		permissions, err := client.GetRepositoryPermissions(namespace, repository)
		if err != nil {
			return fmt.Errorf("getting repository permissions: %w", err)
		}

		fmt.Printf("Permissions for repository %s/%s:\n", namespace, repository)
		return printJSON(permissions)
	},
}

// Permissions Set
var permSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set permission for a user/robot",
	Long:  `Set permission level for a user or robot account on this repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.SetRepositoryPermission(namespace, repository, permissionUser, permissionRole)
		if err != nil {
			return fmt.Errorf("setting repository permission: %w", err)
		}

		fmt.Printf("Successfully set %s permission for %s on repository %s/%s\n",
			permissionRole, permissionUser, namespace, repository)
		return nil
	},
}

// Permissions Remove
var permRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove permission for a user/robot",
	Long:  `Remove permission for a user or robot account from this repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.RemoveRepositoryPermission(namespace, repository, permissionUser)
		if err != nil {
			return fmt.Errorf("removing repository permission: %w", err)
		}

		fmt.Printf("Successfully removed permissions for %s from repository %s/%s\n",
			permissionUser, namespace, repository)
		return nil
	},
}

var permUserPermissionsCmd = &cobra.Command{
	Use:   "user-permissions",
	Short: "List user permissions on repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		permissions, err := client.ListUserPermissions(namespace, repository)
		if err != nil {
			return fmt.Errorf("getting user permissions: %w", err)
		}

		return printJSON(permissions)
	},
}

var permUserPermissionCmd = &cobra.Command{
	Use:   "user-permission",
	Short: "Get specific user permission on repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		permission, err := client.GetUserPermission(namespace, repository, permissionUser)
		if err != nil {
			return fmt.Errorf("getting user permission: %w", err)
		}

		return printJSON(permission)
	},
}

var permSetUserPermCmd = &cobra.Command{
	Use:   "set-user-permission",
	Short: "Set user permission on repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.SetUserPermission(namespace, repository, permissionUser, permissionRole)
		if err != nil {
			return fmt.Errorf("setting user permission: %w", err)
		}

		fmt.Printf("Successfully set %s permission for user %s on %s/%s\n",
			permissionRole, permissionUser, namespace, repository)
		return nil
	},
}

var permDeleteUserPermCmd = &cobra.Command{
	Use:   "delete-user-permission",
	Short: "Delete user permission from repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmPermDeletion {
			return fmt.Errorf("are you sure you want to remove permission for user '%s' from %s/%s?\nUse --confirm to proceed with removal",
				permissionUser, namespace, repository)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteUserPermission(namespace, repository, permissionUser)
		if err != nil {
			return fmt.Errorf("deleting user permission: %w", err)
		}

		fmt.Printf("Successfully removed permission for user %s from %s/%s\n",
			permissionUser, namespace, repository)
		return nil
	},
}

var permUserTransitiveCmd = &cobra.Command{
	Use:   "user-transitive-permission",
	Short: "Get user transitive permission on repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		permission, err := client.GetUserTransitivePermission(namespace, repository, permissionUser)
		if err != nil {
			return fmt.Errorf("getting user transitive permission: %w", err)
		}

		return printJSON(permission)
	},
}

var permTeamPermissionsCmd = &cobra.Command{
	Use:   "team-permissions",
	Short: "List team permissions on repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		permissions, err := client.ListTeamPermissions(namespace, repository)
		if err != nil {
			return fmt.Errorf("getting team permissions: %w", err)
		}

		return printJSON(permissions)
	},
}

var permTeamPermissionCmd = &cobra.Command{
	Use:   "team-permission",
	Short: "Get specific team permission on repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		permission, err := client.GetTeamPermission(namespace, repository, permissionTeamName)
		if err != nil {
			return fmt.Errorf("getting team permission: %w", err)
		}

		return printJSON(permission)
	},
}

var permSetTeamPermCmd = &cobra.Command{
	Use:   "set-team-permission",
	Short: "Set team permission on repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.SetTeamPermission(namespace, repository, permissionTeamName, permissionRole)
		if err != nil {
			return fmt.Errorf("setting team permission: %w", err)
		}

		fmt.Printf("Successfully set %s permission for team %s on %s/%s\n",
			permissionRole, permissionTeamName, namespace, repository)
		return nil
	},
}

var permDeleteTeamPermCmd = &cobra.Command{
	Use:   "delete-team-permission",
	Short: "Delete team permission from repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmPermDeletion {
			return fmt.Errorf("are you sure you want to remove permission for team '%s' from %s/%s?\nUse --confirm to proceed with removal",
				permissionTeamName, namespace, repository)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteTeamPermission(namespace, repository, permissionTeamName)
		if err != nil {
			return fmt.Errorf("deleting team permission: %w", err)
		}

		fmt.Printf("Successfully removed permission for team %s from %s/%s\n",
			permissionTeamName, namespace, repository)
		return nil
	},
}

func init() {
	// Add subcommands to permissions command
	permissionsCmd.AddCommand(permListCmd)
	permissionsCmd.AddCommand(permSetCmd)
	permissionsCmd.AddCommand(permRemoveCmd)
	permissionsCmd.AddCommand(permUserPermissionsCmd)
	permissionsCmd.AddCommand(permUserPermissionCmd)
	permissionsCmd.AddCommand(permSetUserPermCmd)
	permissionsCmd.AddCommand(permDeleteUserPermCmd)
	permissionsCmd.AddCommand(permUserTransitiveCmd)
	permissionsCmd.AddCommand(permTeamPermissionsCmd)
	permissionsCmd.AddCommand(permTeamPermissionCmd)
	permissionsCmd.AddCommand(permSetTeamPermCmd)
	permissionsCmd.AddCommand(permDeleteTeamPermCmd)

	// Global permissions flags (repository context)
	permissionsCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	permissionsCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")

	// Mark global flags as required
	_ = permissionsCmd.MarkPersistentFlagRequired("namespace")
	_ = permissionsCmd.MarkPersistentFlagRequired("repository")

	// Set command specific flags
	permSetCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username or robot account")
	permSetCmd.Flags().StringVarP(&permissionRole, "role", "R", "", "Permission role (read/write/admin)")
	_ = permSetCmd.MarkFlagRequired("user")
	_ = permSetCmd.MarkFlagRequired("role")

	// Remove command specific flags
	permRemoveCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username or robot account")
	_ = permRemoveCmd.MarkFlagRequired("user")

	initPermUserFlags()
	initPermTeamFlags()
}

func initPermUserFlags() {
	permUserPermissionCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username")
	_ = permUserPermissionCmd.MarkFlagRequired("user")

	permSetUserPermCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username")
	permSetUserPermCmd.Flags().StringVarP(&permissionRole, "role", "R", "", "Permission role (read/write/admin)")
	_ = permSetUserPermCmd.MarkFlagRequired("user")
	_ = permSetUserPermCmd.MarkFlagRequired("role")

	permDeleteUserPermCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username")
	permDeleteUserPermCmd.Flags().BoolVar(&confirmPermDeletion, "confirm", false, "Confirm permission removal")
	_ = permDeleteUserPermCmd.MarkFlagRequired("user")

	permUserTransitiveCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username")
	_ = permUserTransitiveCmd.MarkFlagRequired("user")
}

func initPermTeamFlags() {
	permTeamPermissionCmd.Flags().StringVar(&permissionTeamName, "team", "", "Team name")
	_ = permTeamPermissionCmd.MarkFlagRequired("team")

	permSetTeamPermCmd.Flags().StringVar(&permissionTeamName, "team", "", "Team name")
	permSetTeamPermCmd.Flags().StringVarP(&permissionRole, "role", "R", "", "Permission role (read/write/admin)")
	_ = permSetTeamPermCmd.MarkFlagRequired("team")
	_ = permSetTeamPermCmd.MarkFlagRequired("role")

	permDeleteTeamPermCmd.Flags().StringVar(&permissionTeamName, "team", "", "Team name")
	permDeleteTeamPermCmd.Flags().BoolVar(&confirmPermDeletion, "confirm", false, "Confirm permission removal")
	_ = permDeleteTeamPermCmd.MarkFlagRequired("team")
}
