package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	permissionUser string
	permissionRole string
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

func init() {
	// Add subcommands to permissions command
	permissionsCmd.AddCommand(permListCmd)
	permissionsCmd.AddCommand(permSetCmd)
	permissionsCmd.AddCommand(permRemoveCmd)

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
}
