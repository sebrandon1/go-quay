package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	mirrorExternalRef   string
	mirrorSyncInterval  int
	mirrorSyncStartDate string
	mirrorRobotUsername string
	mirrorTagRule       string
	mirrorTagRuleKind   string
	mirrorExtUser       string
	mirrorExtPassword   string
	mirrorEnabled       bool
)

var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Manage repository mirror configuration",
	Long:  `Commands for managing repository mirror configuration including viewing, creating, and updating mirror settings.`,
}

var mirrorInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get mirror configuration",
	Long:  `Get the mirror configuration for a repository.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		config, err := client.GetMirrorConfig(namespace, repository)
		if err != nil {
			return fmt.Errorf("getting mirror config: %w", err)
		}

		return printJSON(config)
	},
}

var mirrorCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
	Short: "Create mirror configuration",
	Long:  `Create mirror configuration for a repository to sync from an external registry.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		createReq := &lib.CreateMirrorConfigRequest{
			ExternalRef:              mirrorExternalRef,
			SyncInterval:             mirrorSyncInterval,
			SyncStartDate:            mirrorSyncStartDate,
			RobotUsername:            mirrorRobotUsername,
			ExternalRegistryUsername: mirrorExtUser,
			ExternalRegistryPassword: mirrorExtPassword,
		}
		createReq.RootRule.Rule = mirrorTagRule
		createReq.RootRule.RuleKind = mirrorTagRuleKind

		config, err := client.CreateMirrorConfig(namespace, repository, createReq)
		if err != nil {
			return fmt.Errorf("creating mirror config: %w", err)
		}

		fmt.Fprintln(os.Stderr, "Mirror configuration created successfully")
		return printJSON(config)
	},
}

var mirrorUpdateCmd = &cobra.Command{
	Use:   subcmdUpdate,
	Short: "Update mirror configuration",
	Long:  `Update mirror configuration for a repository.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		updateReq := &lib.UpdateMirrorConfigRequest{
			IsEnabled: &mirrorEnabled,
		}
		if mirrorExternalRef != "" {
			updateReq.ExternalRef = mirrorExternalRef
		}
		if mirrorSyncInterval > 0 {
			updateReq.SyncInterval = &mirrorSyncInterval
		}
		if mirrorSyncStartDate != "" {
			updateReq.SyncStartDate = mirrorSyncStartDate
		}
		if mirrorRobotUsername != "" {
			updateReq.RobotUsername = mirrorRobotUsername
		}

		config, err := client.UpdateMirrorConfig(namespace, repository, updateReq)
		if err != nil {
			return fmt.Errorf("updating mirror config: %w", err)
		}

		fmt.Fprintln(os.Stderr, "Mirror configuration updated successfully")
		return printJSON(config)
	},
}

func init() {
	mirrorCmd.AddCommand(mirrorInfoCmd)
	mirrorCmd.AddCommand(mirrorCreateCmd)
	mirrorCmd.AddCommand(mirrorUpdateCmd)

	mirrorCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", appCfg.Namespace, "Name of the namespace (default: config file)")
	mirrorCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	if appCfg.Namespace == "" {
		_ = mirrorCmd.MarkPersistentFlagRequired("namespace")
	}
	_ = mirrorCmd.MarkPersistentFlagRequired("repository")

	mirrorCreateCmd.Flags().StringVar(&mirrorExternalRef, "external-ref", "", "External registry reference (e.g. docker.io/library/nginx)")
	_ = mirrorCreateCmd.MarkFlagRequired("external-ref")
	mirrorCreateCmd.Flags().IntVar(&mirrorSyncInterval, "sync-interval", 86400, "Sync interval in seconds")
	mirrorCreateCmd.Flags().StringVar(&mirrorSyncStartDate, "sync-start-date", "", "Sync start date (ISO format)")
	mirrorCreateCmd.Flags().StringVar(&mirrorRobotUsername, "robot-username", "", "Robot account for pulling")
	_ = mirrorCreateCmd.MarkFlagRequired("robot-username")
	mirrorCreateCmd.Flags().StringVar(&mirrorTagRule, "tag-rule", ".*", "Tag filter rule (regex)")
	mirrorCreateCmd.Flags().StringVar(&mirrorTagRuleKind, "tag-rule-kind", "tag_glob_csv", "Tag rule kind")
	mirrorCreateCmd.Flags().StringVar(&mirrorExtUser, "ext-username", "", "External registry username")
	mirrorCreateCmd.Flags().StringVar(&mirrorExtPassword, "ext-password", "", "External registry password")

	mirrorUpdateCmd.Flags().BoolVar(&mirrorEnabled, "enabled", true, "Enable or disable mirroring")
	mirrorUpdateCmd.Flags().StringVar(&mirrorExternalRef, "external-ref", "", "External registry reference")
	mirrorUpdateCmd.Flags().IntVar(&mirrorSyncInterval, "sync-interval", 0, "Sync interval in seconds")
	mirrorUpdateCmd.Flags().StringVar(&mirrorSyncStartDate, "sync-start-date", "", "Sync start date")
	mirrorUpdateCmd.Flags().StringVar(&mirrorRobotUsername, "robot-username", "", "Robot account for pulling")
}
