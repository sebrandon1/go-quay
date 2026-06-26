package cmd

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		client := mustGetClient()

		robots, err := client.GetUserRobotAccounts()
		if err != nil {
			fmt.Printf("Error getting robot accounts: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("User robot accounts:")
		printJSON(robots)
	},
}

// Robot Info
var robotInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get robot account details",
	Long:  `Get detailed information about a specific robot account including its token.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := mustGetClient()

		robot, err := client.GetUserRobotAccount(robotShortname)
		if err != nil {
			fmt.Printf("Error getting robot account: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Robot account: %s\n", robotShortname)
		printJSON(robot)
	},
}

// Robot Create
var robotCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
	Short: "Create a new robot account",
	Long:  `Create a new robot account with the specified name and description.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := mustGetClient()

		robot, err := client.CreateUserRobotAccount(robotShortname, robotDescription, nil)
		if err != nil {
			fmt.Printf("Error creating robot account: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created robot account: %s\n", robotShortname)
		fmt.Println("IMPORTANT: Save the token below - it will not be shown again!")
		printJSON(robot)
	},
}

// Robot Delete
var robotDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a robot account",
	Long:  `Delete a robot account. This action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmRobotDelete {
			fmt.Printf("Are you sure you want to delete robot account '%s'? This action cannot be undone.\n", robotShortname)
			fmt.Println("Use --confirm to proceed with deletion.")
			os.Exit(1)
		}

		client := mustGetClient()

		err := client.DeleteUserRobotAccount(robotShortname)
		if err != nil {
			fmt.Printf("Error deleting robot account: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted robot account: %s\n", robotShortname)
	},
}

// Robot Regenerate
var robotRegenerateCmd = &cobra.Command{
	Use:   "regenerate",
	Short: "Regenerate robot token",
	Long:  `Regenerate the token for a robot account. The old token will be invalidated.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := mustGetClient()

		robot, err := client.RegenerateUserRobotToken(robotShortname)
		if err != nil {
			fmt.Printf("Error regenerating robot token: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully regenerated token for robot: %s\n", robotShortname)
		fmt.Println("IMPORTANT: Save the new token below - it will not be shown again!")
		printJSON(robot)
	},
}

// Robot Permissions
var robotPermissionsCmd = &cobra.Command{
	Use:   subcmdPermissions,
	Short: "Get robot repository permissions",
	Long:  `Get the repository permissions assigned to a robot account.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := mustGetClient()

		permissions, err := client.GetUserRobotPermissions(robotShortname)
		if err != nil {
			fmt.Printf("Error getting robot permissions: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Permissions for robot: %s\n", robotShortname)
		printJSON(permissions)
	},
}

// Robot Federation Get
var robotFederationGetCmd = &cobra.Command{
	Use:   "federation-get",
	Short: "Get robot federation configuration",
	Run: func(cmd *cobra.Command, args []string) {
		client := mustGetClient()

		federation, err := client.GetUserRobotFederation(robotShortname)
		if err != nil {
			fmt.Printf("Error getting robot federation: %v\n", err)
			os.Exit(1)
		}

		printJSON(federation)
	},
}

// Robot Federation Create
var robotFederationCreateCmd = &cobra.Command{
	Use:   "federation-create",
	Short: "Create or update robot federation configuration",
	Run: func(cmd *cobra.Command, args []string) {
		client := mustGetClient()

		configs := []lib.RobotFederationConfig{
			{Issuer: federationIssuer, Subject: federationSubject},
		}

		err := client.CreateUserRobotFederation(robotShortname, configs)
		if err != nil {
			fmt.Printf("Error creating robot federation: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully configured federation for robot: %s\n", robotShortname)
	},
}

// Robot Federation Delete
var robotFederationDeleteCmd = &cobra.Command{
	Use:   "federation-delete",
	Short: "Delete robot federation configuration",
	Run: func(cmd *cobra.Command, args []string) {
		client := mustGetClient()

		err := client.DeleteUserRobotFederation(robotShortname)
		if err != nil {
			fmt.Printf("Error deleting robot federation: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted federation for robot: %s\n", robotShortname)
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
	if err := robotInfoCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("Error marking name flag as required: %v\n", err)
		os.Exit(1)
	}

	// Create command flags
	robotCreateCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	robotCreateCmd.Flags().StringVarP(&robotDescription, "description", "d", "", "Robot description")
	if err := robotCreateCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("Error marking name flag as required: %v\n", err)
		os.Exit(1)
	}

	// Delete command flags
	robotDeleteCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	robotDeleteCmd.Flags().BoolVar(&confirmRobotDelete, "confirm", false, "Confirm robot deletion")
	if err := robotDeleteCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("Error marking name flag as required: %v\n", err)
		os.Exit(1)
	}

	// Regenerate command flags
	robotRegenerateCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	if err := robotRegenerateCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("Error marking name flag as required: %v\n", err)
		os.Exit(1)
	}

	// Permissions command flags
	robotPermissionsCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name (without username prefix)")
	if err := robotPermissionsCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("Error marking name flag as required: %v\n", err)
		os.Exit(1)
	}

	// Federation get command flags
	robotFederationGetCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name")
	markFlagRequired(robotFederationGetCmd.MarkFlagRequired("name"))

	// Federation create command flags
	robotFederationCreateCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name")
	robotFederationCreateCmd.Flags().StringVar(&federationIssuer, "issuer", "", "Federation token issuer")
	robotFederationCreateCmd.Flags().StringVar(&federationSubject, "subject", "", "Federation token subject")
	markFlagRequired(robotFederationCreateCmd.MarkFlagRequired("name"))
	markFlagRequired(robotFederationCreateCmd.MarkFlagRequired("issuer"))
	markFlagRequired(robotFederationCreateCmd.MarkFlagRequired("subject"))

	// Federation delete command flags
	robotFederationDeleteCmd.Flags().StringVarP(&robotShortname, "name", "n", "", "Robot short name")
	markFlagRequired(robotFederationDeleteCmd.MarkFlagRequired("name"))
}
