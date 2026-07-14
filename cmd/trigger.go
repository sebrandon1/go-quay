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
	"fmt"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	triggerNamespace  string
	triggerRepository string
	triggerUUID       string
	triggerCommitSHA  string
	triggerPullRobot  string
	triggerBuildLimit int
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
	Use:   subcmdList,
	Short: "List all build triggers for a repository",
	Long:  `List all build triggers configured for a repository.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		triggers, err := client.GetTriggers(triggerNamespace, triggerRepository)
		if err != nil {
			return fmt.Errorf("getting triggers: %w", err)
		}

		return printJSON(triggers)
	},
}

// triggerInfoCmd represents the trigger info command
var triggerInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get details of a specific build trigger",
	Long:  `Get detailed information about a specific build trigger by its UUID.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if triggerUUID == "" {
			return fmt.Errorf("--uuid is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		trigger, err := client.GetTrigger(triggerNamespace, triggerRepository, triggerUUID)
		if err != nil {
			return fmt.Errorf("getting trigger: %w", err)
		}

		return printJSON(trigger)
	},
}

// triggerDeleteCmd represents the trigger delete command
var triggerDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a build trigger",
	Long:  `Delete a build trigger from a repository.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if triggerUUID == "" {
			return fmt.Errorf("--uuid is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteTrigger(triggerNamespace, triggerRepository, triggerUUID)
		if err != nil {
			return fmt.Errorf("deleting trigger: %w", err)
		}

		fmt.Printf("Trigger %s deleted successfully\n", triggerUUID)
		return nil
	},
}

// triggerEnableCmd represents the trigger enable command
var triggerEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a build trigger",
	Long:  `Enable a build trigger for a repository.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if triggerUUID == "" {
			return fmt.Errorf("--uuid is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		trigger, err := client.UpdateTrigger(triggerNamespace, triggerRepository, triggerUUID, true)
		if err != nil {
			return fmt.Errorf("enabling trigger: %w", err)
		}

		return printJSON(trigger)
	},
}

// triggerDisableCmd represents the trigger disable command
var triggerDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a build trigger",
	Long:  `Disable a build trigger for a repository.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if triggerUUID == "" {
			return fmt.Errorf("--uuid is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		trigger, err := client.UpdateTrigger(triggerNamespace, triggerRepository, triggerUUID, false)
		if err != nil {
			return fmt.Errorf("disabling trigger: %w", err)
		}

		return printJSON(trigger)
	},
}

// triggerStartCmd represents the trigger start command
var triggerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Manually start a build from a trigger",
	Long:  `Manually start a build using a configured trigger.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if triggerUUID == "" {
			return fmt.Errorf("--uuid is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		var triggerReq *lib.ManualTriggerRequest
		if triggerCommitSHA != "" {
			triggerReq = &lib.ManualTriggerRequest{
				CommitSHA: triggerCommitSHA,
			}
		}

		build, err := client.StartTriggerBuild(triggerNamespace, triggerRepository, triggerUUID, triggerReq)
		if err != nil {
			return fmt.Errorf("starting trigger build: %w", err)
		}

		return printJSON(build)
	},
}

// triggerActivateCmd represents the trigger activate command
var triggerActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate a build trigger with configuration",
	Long:  `Activate a build trigger with the specified configuration.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if triggerUUID == "" {
			return fmt.Errorf("--uuid is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		activateReq := &lib.ActivateTriggerRequest{}
		if triggerPullRobot != "" {
			activateReq.PullRobot = triggerPullRobot
		}

		trigger, err := client.ActivateTrigger(triggerNamespace, triggerRepository, triggerUUID, activateReq)
		if err != nil {
			return fmt.Errorf("activating trigger: %w", err)
		}

		return printJSON(trigger)
	},
}

var triggerBuildsCmd = &cobra.Command{
	Use:   "builds",
	Short: "List builds for a trigger",
	Long:  `List builds that were started by a specific trigger.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if triggerUUID == "" {
			return fmt.Errorf("--uuid is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		builds, err := client.GetTriggerBuilds(triggerNamespace, triggerRepository, triggerUUID, triggerBuildLimit)
		if err != nil {
			return fmt.Errorf("getting trigger builds: %w", err)
		}

		return printJSON(builds)
	},
}

func setupTriggerFlags() {
	// Common flags for all subcommands
	for _, cmd := range []*cobra.Command{triggerListCmd, triggerInfoCmd, triggerDeleteCmd, triggerEnableCmd, triggerDisableCmd, triggerStartCmd, triggerActivateCmd, triggerBuildsCmd} {
		cmd.Flags().StringVarP(&triggerNamespace, "namespace", "n", "", "Namespace/organization")
		cmd.Flags().StringVarP(&triggerRepository, "repository", "r", "", "Repository name")
	}

	// UUID flags for subcommands that need it
	for _, cmd := range []*cobra.Command{triggerInfoCmd, triggerDeleteCmd, triggerEnableCmd, triggerDisableCmd, triggerStartCmd, triggerActivateCmd, triggerBuildsCmd} {
		cmd.Flags().StringVar(&triggerUUID, "uuid", "", "UUID of the build trigger")
	}

	// Start command specific flags
	triggerStartCmd.Flags().StringVar(&triggerCommitSHA, "commit-sha", "", "Specific commit SHA to build (optional)")

	// Activate command specific flags
	triggerActivateCmd.Flags().StringVar(&triggerPullRobot, "pull-robot", "", "Robot account to use for pulling base images (optional)")

	// Builds command specific flags
	triggerBuildsCmd.Flags().IntVarP(&triggerBuildLimit, "limit", "l", 10, "Maximum number of builds to return")
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
	triggerCmd.AddCommand(triggerBuildsCmd)

	// Setup flags
	setupTriggerFlags()
}
