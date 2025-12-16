/*
Package cmd provides the command-line interface for go-quay.

This file contains the commands for the Trigger API:
  - go-quay get trigger list         - List all build triggers for a repository
  - go-quay get trigger info         - Get details of a specific build trigger
  - go-quay get trigger delete       - Delete a build trigger
  - go-quay get trigger enable       - Enable a build trigger
  - go-quay get trigger disable      - Disable a build trigger
  - go-quay get trigger start        - Manually start a build from a trigger
  - go-quay get trigger activate     - Activate a build trigger with configuration
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	triggerNamespace  string
	triggerRepository string
	triggerUUID       string
	triggerCommitSHA  string
	triggerPullRobot  string
)

// triggerCmd represents the trigger command
var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Manage repository build triggers",
	Long: `Manage build triggers for repositories.

Build triggers allow automated image builds when code is pushed to 
connected source repositories like GitHub, GitLab, or Bitbucket.`,
}

// triggerListCmd represents the trigger list command
var triggerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all build triggers for a repository",
	Long:  `List all build triggers configured for a repository.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		triggers, err := client.GetTriggers(triggerNamespace, triggerRepository)
		if err != nil {
			fmt.Println("Error getting triggers:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(triggers, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// triggerInfoCmd represents the trigger info command
var triggerInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get details of a specific build trigger",
	Long:  `Get detailed information about a specific build trigger by its UUID.`,
	Run: func(_ *cobra.Command, _ []string) {
		if triggerUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		trigger, err := client.GetTrigger(triggerNamespace, triggerRepository, triggerUUID)
		if err != nil {
			fmt.Println("Error getting trigger:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(trigger, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// triggerDeleteCmd represents the trigger delete command
var triggerDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a build trigger",
	Long:  `Delete a build trigger from a repository.`,
	Run: func(_ *cobra.Command, _ []string) {
		if triggerUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		err = client.DeleteTrigger(triggerNamespace, triggerRepository, triggerUUID)
		if err != nil {
			fmt.Println("Error deleting trigger:", err)
			os.Exit(1)
		}

		fmt.Printf("Trigger %s deleted successfully\n", triggerUUID)
	},
}

// triggerEnableCmd represents the trigger enable command
var triggerEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a build trigger",
	Long:  `Enable a build trigger for a repository.`,
	Run: func(_ *cobra.Command, _ []string) {
		if triggerUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		trigger, err := client.UpdateTrigger(triggerNamespace, triggerRepository, triggerUUID, true)
		if err != nil {
			fmt.Println("Error enabling trigger:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(trigger, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// triggerDisableCmd represents the trigger disable command
var triggerDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a build trigger",
	Long:  `Disable a build trigger for a repository.`,
	Run: func(_ *cobra.Command, _ []string) {
		if triggerUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		trigger, err := client.UpdateTrigger(triggerNamespace, triggerRepository, triggerUUID, false)
		if err != nil {
			fmt.Println("Error disabling trigger:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(trigger, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// triggerStartCmd represents the trigger start command
var triggerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Manually start a build from a trigger",
	Long:  `Manually start a build using a configured trigger.`,
	Run: func(_ *cobra.Command, _ []string) {
		if triggerUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		var triggerReq *lib.ManualTriggerRequest
		if triggerCommitSHA != "" {
			triggerReq = &lib.ManualTriggerRequest{
				CommitSHA: triggerCommitSHA,
			}
		}

		build, err := client.StartTriggerBuild(triggerNamespace, triggerRepository, triggerUUID, triggerReq)
		if err != nil {
			fmt.Println("Error starting trigger build:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(build, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// triggerActivateCmd represents the trigger activate command
var triggerActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate a build trigger with configuration",
	Long:  `Activate a build trigger with the specified configuration.`,
	Run: func(_ *cobra.Command, _ []string) {
		if triggerUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		activateReq := &lib.ActivateTriggerRequest{}
		if triggerPullRobot != "" {
			activateReq.PullRobot = triggerPullRobot
		}

		trigger, err := client.ActivateTrigger(triggerNamespace, triggerRepository, triggerUUID, activateReq)
		if err != nil {
			fmt.Println("Error activating trigger:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(trigger, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

func setupTriggerFlags() {
	// Common flags for all subcommands
	for _, cmd := range []*cobra.Command{triggerListCmd, triggerInfoCmd, triggerDeleteCmd, triggerEnableCmd, triggerDisableCmd, triggerStartCmd, triggerActivateCmd} {
		cmd.Flags().StringVarP(&triggerNamespace, "namespace", "n", "", "Namespace/organization")
		cmd.Flags().StringVarP(&triggerRepository, "repository", "r", "", "Repository name")
		cmd.Flags().StringVarP(&token, "token", "t", "", "Quay.io API token")
	}

	// UUID flags for subcommands that need it
	for _, cmd := range []*cobra.Command{triggerInfoCmd, triggerDeleteCmd, triggerEnableCmd, triggerDisableCmd, triggerStartCmd, triggerActivateCmd} {
		cmd.Flags().StringVar(&triggerUUID, "uuid", "", "UUID of the build trigger")
	}

	// Start command specific flags
	triggerStartCmd.Flags().StringVar(&triggerCommitSHA, "commit-sha", "", "Specific commit SHA to build (optional)")

	// Activate command specific flags
	triggerActivateCmd.Flags().StringVar(&triggerPullRobot, "pull-robot", "", "Robot account to use for pulling base images (optional)")
}

func init() {
	// Add subcommands to trigger command
	triggerCmd.AddCommand(triggerListCmd)
	triggerCmd.AddCommand(triggerInfoCmd)
	triggerCmd.AddCommand(triggerDeleteCmd)
	triggerCmd.AddCommand(triggerEnableCmd)
	triggerCmd.AddCommand(triggerDisableCmd)
	triggerCmd.AddCommand(triggerStartCmd)
	triggerCmd.AddCommand(triggerActivateCmd)

	// Setup flags
	setupTriggerFlags()
}
