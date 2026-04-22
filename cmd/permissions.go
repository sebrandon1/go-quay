package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
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
	Use:   "permissions",
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
	Use:   "list",
	Short: "List repository permissions",
	Long:  `List all users and robots that have permissions on this repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		permissions, err := client.GetRepositoryPermissions(namespace, repository)
		if err != nil {
			fmt.Printf("Error getting repository permissions: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Permissions for repository %s/%s:\n", namespace, repository)
		printJSON(permissions)
	},
}

// Permissions Set
var permSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set permission for a user/robot",
	Long:  `Set permission level for a user or robot account on this repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.SetRepositoryPermission(namespace, repository, permissionUser, permissionRole)
		if err != nil {
			fmt.Printf("Error setting repository permission: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set %s permission for %s on repository %s/%s\n",
			permissionRole, permissionUser, namespace, repository)
	},
}

// Permissions Remove
var permRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove permission for a user/robot",
	Long:  `Remove permission for a user or robot account from this repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.RemoveRepositoryPermission(namespace, repository, permissionUser)
		if err != nil {
			fmt.Printf("Error removing repository permission: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully removed permissions for %s from repository %s/%s\n",
			permissionUser, namespace, repository)
	},
}

var permUserPermissionsCmd = &cobra.Command{
	Use:   "user-permissions",
	Short: "List user permissions on repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		permissions, err := client.ListUserPermissions(namespace, repository)
		if err != nil {
			fmt.Printf("Error getting user permissions: %v\n", err)
			os.Exit(1)
		}

		printJSON(permissions)
	},
}

var permUserPermissionCmd = &cobra.Command{
	Use:   "user-permission",
	Short: "Get specific user permission on repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		permission, err := client.GetUserPermission(namespace, repository, permissionUser)
		if err != nil {
			fmt.Printf("Error getting user permission: %v\n", err)
			os.Exit(1)
		}

		printJSON(permission)
	},
}

var permSetUserPermCmd = &cobra.Command{
	Use:   "set-user-permission",
	Short: "Set user permission on repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.SetUserPermission(namespace, repository, permissionUser, permissionRole)
		if err != nil {
			fmt.Printf("Error setting user permission: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set %s permission for user %s on %s/%s\n",
			permissionRole, permissionUser, namespace, repository)
	},
}

var permDeleteUserPermCmd = &cobra.Command{
	Use:   "delete-user-permission",
	Short: "Delete user permission from repository",
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmPermDeletion {
			fmt.Printf("Are you sure you want to remove permission for user '%s' from %s/%s?\n",
				permissionUser, namespace, repository)
			fmt.Println("Use --confirm to proceed with removal.")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.DeleteUserPermission(namespace, repository, permissionUser)
		if err != nil {
			fmt.Printf("Error deleting user permission: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully removed permission for user %s from %s/%s\n",
			permissionUser, namespace, repository)
	},
}

var permUserTransitiveCmd = &cobra.Command{
	Use:   "user-transitive-permission",
	Short: "Get user transitive permission on repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		permission, err := client.GetUserTransitivePermission(namespace, repository, permissionUser)
		if err != nil {
			fmt.Printf("Error getting user transitive permission: %v\n", err)
			os.Exit(1)
		}

		printJSON(permission)
	},
}

var permTeamPermissionsCmd = &cobra.Command{
	Use:   "team-permissions",
	Short: "List team permissions on repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		permissions, err := client.ListTeamPermissions(namespace, repository)
		if err != nil {
			fmt.Printf("Error getting team permissions: %v\n", err)
			os.Exit(1)
		}

		printJSON(permissions)
	},
}

var permTeamPermissionCmd = &cobra.Command{
	Use:   "team-permission",
	Short: "Get specific team permission on repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		permission, err := client.GetTeamPermission(namespace, repository, permissionTeamName)
		if err != nil {
			fmt.Printf("Error getting team permission: %v\n", err)
			os.Exit(1)
		}

		printJSON(permission)
	},
}

var permSetTeamPermCmd = &cobra.Command{
	Use:   "set-team-permission",
	Short: "Set team permission on repository",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.SetTeamPermission(namespace, repository, permissionTeamName, permissionRole)
		if err != nil {
			fmt.Printf("Error setting team permission: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set %s permission for team %s on %s/%s\n",
			permissionRole, permissionTeamName, namespace, repository)
	},
}

var permDeleteTeamPermCmd = &cobra.Command{
	Use:   "delete-team-permission",
	Short: "Delete team permission from repository",
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmPermDeletion {
			fmt.Printf("Are you sure you want to remove permission for team '%s' from %s/%s?\n",
				permissionTeamName, namespace, repository)
			fmt.Println("Use --confirm to proceed with removal.")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.DeleteTeamPermission(namespace, repository, permissionTeamName)
		if err != nil {
			fmt.Printf("Error deleting team permission: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully removed permission for team %s from %s/%s\n",
			permissionTeamName, namespace, repository)
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
	permissionsCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")

	// Mark global flags as required
	if err := permissionsCmd.MarkPersistentFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := permissionsCmd.MarkPersistentFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := permissionsCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}

	// Set command specific flags
	permSetCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username or robot account")
	permSetCmd.Flags().StringVarP(&permissionRole, "role", "R", "", "Permission role (read/write/admin)")
	if err := permSetCmd.MarkFlagRequired("user"); err != nil {
		fmt.Printf("Error marking user flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := permSetCmd.MarkFlagRequired("role"); err != nil {
		fmt.Printf("Error marking role flag as required: %v\n", err)
		os.Exit(1)
	}

	// Remove command specific flags
	permRemoveCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username or robot account")
	if err := permRemoveCmd.MarkFlagRequired("user"); err != nil {
		fmt.Printf("Error marking user flag as required: %v\n", err)
		os.Exit(1)
	}

	initPermUserFlags()
	initPermTeamFlags()
}

func initPermUserFlags() {
	permUserPermissionCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username")
	markFlagRequired(permUserPermissionCmd.MarkFlagRequired("user"))

	permSetUserPermCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username")
	permSetUserPermCmd.Flags().StringVarP(&permissionRole, "role", "R", "", "Permission role (read/write/admin)")
	markFlagRequired(permSetUserPermCmd.MarkFlagRequired("user"))
	markFlagRequired(permSetUserPermCmd.MarkFlagRequired("role"))

	permDeleteUserPermCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username")
	permDeleteUserPermCmd.Flags().BoolVar(&confirmPermDeletion, "confirm", false, "Confirm permission removal")
	markFlagRequired(permDeleteUserPermCmd.MarkFlagRequired("user"))

	permUserTransitiveCmd.Flags().StringVarP(&permissionUser, "user", "u", "", "Username")
	markFlagRequired(permUserTransitiveCmd.MarkFlagRequired("user"))
}

func initPermTeamFlags() {
	permTeamPermissionCmd.Flags().StringVar(&permissionTeamName, "team", "", "Team name")
	markFlagRequired(permTeamPermissionCmd.MarkFlagRequired("team"))

	permSetTeamPermCmd.Flags().StringVar(&permissionTeamName, "team", "", "Team name")
	permSetTeamPermCmd.Flags().StringVarP(&permissionRole, "role", "R", "", "Permission role (read/write/admin)")
	markFlagRequired(permSetTeamPermCmd.MarkFlagRequired("team"))
	markFlagRequired(permSetTeamPermCmd.MarkFlagRequired("role"))

	permDeleteTeamPermCmd.Flags().StringVar(&permissionTeamName, "team", "", "Team name")
	permDeleteTeamPermCmd.Flags().BoolVar(&confirmPermDeletion, "confirm", false, "Confirm permission removal")
	markFlagRequired(permDeleteTeamPermCmd.MarkFlagRequired("team"))
}
