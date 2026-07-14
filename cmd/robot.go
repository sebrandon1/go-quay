package cmd

import (
	"fmt"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	robotShortname     string
	robotDescription   string
	confirmRobotDelete bool
	federationIssuer   string
	federationSubject  string
)

// robotCmd represents the robot command group
var robotCmd = &cobra.Command{
	Use:   "robot",
	Short: "User robot account management commands",
	Long: `Commands for managing user-level robot accounts.

Robot accounts provide automated access credentials for CI/CD pipelines
and other automated workflows. These are tied to your user account.

Available commands:
  list       - List all robot accounts
  info       - Get robot account details
  create     - Create a new robot account
  delete     - Delete a robot account
  regenerate - Regenerate robot token
  permissions - Get robot repository permissions`,
}

// Robot List
var robotListCmd = &cobra.Command{
	Use:   subcmdList,
	Short: "List all robot accounts",
	Long:  `List all robot accounts associated with your user account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		robots, err := client.GetUserRobotAccounts()
		if err != nil {
			return fmt.Errorf("getting robot accounts: %w", err)
		}

		fmt.Println("User robot accounts:")
		return printJSON(robots)
	},
}

// Robot Info
var robotInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get robot account details",
	Long:  `Get detailed information about a specific robot account including its token.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		robot, err := client.GetUserRobotAccount(robotShortname)
		if err != nil {
			return fmt.Errorf("getting robot account: %w", err)
		}

		fmt.Printf("Robot account: %s\n", robotShortname)
		return printJSON(robot)
	},
}

// Robot Create
var robotCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
	Short: "Create a new robot account",
	Long:  `Create a new robot account with the specified name and description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		robot, err := client.CreateUserRobotAccount(robotShortname, robotDescription, nil)
		if err != nil {
			return fmt.Errorf("creating robot account: %w", err)
		}

		fmt.Printf("Successfully created robot account: %s\n", robotShortname)
		fmt.Println("IMPORTANT: Save the token below - it will not be shown again!")
		return printJSON(robot)
	},
}

// Robot Delete
var robotDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a robot account",
	Long:  `Delete a robot account. This action cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmRobotDelete {
			return fmt.Errorf("are you sure you want to delete robot account '%s'? This action cannot be undone.\nUse --confirm to proceed with deletion", robotShortname)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteUserRobotAccount(robotShortname)
		if err != nil {
			return fmt.Errorf("deleting robot account: %w", err)
		}

		fmt.Printf("Successfully deleted robot account: %s\n", robotShortname)
		return nil
	},
}

// Robot Regenerate
var robotRegenerateCmd = &cobra.Command{
	Use:   "regenerate",
	Short: "Regenerate robot token",
	Long:  `Regenerate the token for a robot account. The old token will be invalidated.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		robot, err := client.RegenerateUserRobotToken(robotShortname)
		if err != nil {
			return fmt.Errorf("regenerating robot token: %w", err)
		}

		fmt.Printf("Successfully regenerated token for robot: %s\n", robotShortname)
		fmt.Println("IMPORTANT: Save the new token below - it will not be shown again!")
		return printJSON(robot)
	},
}

// Robot Permissions
var robotPermissionsCmd = &cobra.Command{
	Use:   subcmdPermissions,
	Short: "Get robot repository permissions",
	Long:  `Get the repository permissions assigned to a robot account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		permissions, err := client.GetUserRobotPermissions(robotShortname)
		if err != nil {
			return fmt.Errorf("getting robot permissions: %w", err)
		}

		fmt.Printf("Permissions for robot: %s\n", robotShortname)
		return printJSON(permissions)
	},
}

// Robot Federation Get
var robotFederationGetCmd = &cobra.Command{
	Use:   "federation-get",
	Short: "Get robot federation configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		federation, err := client.GetUserRobotFederation(robotShortname)
		if err != nil {
			return fmt.Errorf("getting robot federation: %w", err)
		}

		return printJSON(federation)
	},
}

// Robot Federation Create
var robotFederationCreateCmd = &cobra.Command{
	Use:   "federation-create",
	Short: "Create or update robot federation configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		configs := []lib.RobotFederationConfig{
			{Issuer: federationIssuer, Subject: federationSubject},
		}

		err = client.CreateUserRobotFederation(robotShortname, configs)
		if err != nil {
			return fmt.Errorf("creating robot federation: %w", err)
		}

		fmt.Printf("Successfully configured federation for robot: %s\n", robotShortname)
		return nil
	},
}

// Robot Federation Delete
var robotFederationDeleteCmd = &cobra.Command{
	Use:   "federation-delete",
	Short: "Delete robot federation configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteUserRobotFederation(robotShortname)
		if err != nil {
			return fmt.Errorf("deleting robot federation: %w", err)
		}

		fmt.Printf("Successfully deleted federation for robot: %s\n", robotShortname)
		return nil
	},
}

func init() {
	// Add subcommands to robot command
	robotCmd.AddCommand(robotListCmd)
	robotCmd.AddCommand(robotInfoCmd)
	robotCmd.AddCommand(robotCreateCmd)
	robotCmd.AddCommand(robotDeleteCmd)
	robotCmd.AddCommand(robotRegenerateCmd)
	robotCmd.AddCommand(robotPermissionsCmd)
	robotCmd.AddCommand(robotFederationGetCmd)
	robotCmd.AddCommand(robotFederationCreateCmd)
	robotCmd.AddCommand(robotFederationDeleteCmd)

	// Info command flags
	robotInfoCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	_ = robotInfoCmd.MarkFlagRequired("name")

	// Create command flags
	robotCreateCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	robotCreateCmd.Flags().StringVarP(&robotDescription, "description", "d", "", "Robot description")
	_ = robotCreateCmd.MarkFlagRequired("name")

	// Delete command flags
	robotDeleteCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	robotDeleteCmd.Flags().BoolVar(&confirmRobotDelete, "confirm", false, "Confirm robot deletion")
	_ = robotDeleteCmd.MarkFlagRequired("name")

	// Regenerate command flags
	robotRegenerateCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	_ = robotRegenerateCmd.MarkFlagRequired("name")

	// Permissions command flags
	robotPermissionsCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	_ = robotPermissionsCmd.MarkFlagRequired("name")

	// Federation get command flags
	robotFederationGetCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name")
	_ = robotFederationGetCmd.MarkFlagRequired("name")

	// Federation create command flags
	robotFederationCreateCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name")
	robotFederationCreateCmd.Flags().StringVar(&federationIssuer, "issuer", "", "Federation token issuer")
	robotFederationCreateCmd.Flags().StringVar(&federationSubject, "subject", "", "Federation token subject")
	_ = robotFederationCreateCmd.MarkFlagRequired("name")
	_ = robotFederationCreateCmd.MarkFlagRequired("issuer")
	_ = robotFederationCreateCmd.MarkFlagRequired("subject")

	// Federation delete command flags
	robotFederationDeleteCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name")
	_ = robotFederationDeleteCmd.MarkFlagRequired("name")
}
