package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

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
		fmt.Fprintln(os.Stderr, "Robot account deleted successfully")
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
		fmt.Fprintln(os.Stderr, "Robot permission set successfully")
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
		fmt.Fprintln(os.Stderr, "Robot permission removed successfully")
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

		fmt.Fprintf(os.Stderr, "Successfully configured federation for robot %s in org %s\n", robotShortname, orgName)
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

		fmt.Fprintf(os.Stderr, "Successfully deleted federation for robot %s in org %s\n", robotShortname, orgName)
		return nil
	},
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

func init() {
	organizationCmd.AddCommand(orgRobotsCmd)
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

	initOrgRobotFlags()
}
