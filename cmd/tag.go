package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	tagName            string
	tagExpiration      string
	manifestDigest     string
	confirmTagDeletion bool
)

// tagCmd represents the tag command group
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Repository tag management commands",
	Long: `Commands for managing repository tags including detailed information, updates, deletion, and history.

Available commands:
  info     - Get detailed tag information
  update   - Update tag metadata
  delete   - Delete a tag
  history  - Get tag history
  revert   - Revert tag to a previous state`,
}

// Tag Info
var tagInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get detailed tag information",
	Long:  `Get detailed information about a specific tag including metadata, manifest digest, and size.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		tag, err := client.GetTag(namespace, repository, tagName)
		if err != nil {
			return fmt.Errorf("getting tag information: %w", err)
		}

		fmt.Printf("Tag information for %s/%s:%s\n", namespace, repository, tagName)
		return printJSON(tag)
	},
}

// Tag Update
var tagUpdateCmd = &cobra.Command{
	Use:   subcmdUpdate,
	Short: "Update tag metadata",
	Long:  `Update tag metadata such as expiration date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		tag, err := client.UpdateTag(namespace, repository, tagName, tagExpiration)
		if err != nil {
			return fmt.Errorf("updating tag: %w", err)
		}

		fmt.Printf("Successfully updated tag %s/%s:%s\n", namespace, repository, tagName)
		return printJSON(tag)
	},
}

// Tag Delete
var tagDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a tag",
	Long:  `Delete a specific tag from the repository. This action cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmTagDeletion {
			return fmt.Errorf("are you sure you want to delete tag %s/%s:%s? This action cannot be undone.\nUse --confirm to proceed with deletion", namespace, repository, tagName)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteTag(namespace, repository, tagName)
		if err != nil {
			return fmt.Errorf("deleting tag: %w", err)
		}

		fmt.Printf("Successfully deleted tag %s/%s:%s\n", namespace, repository, tagName)
		return nil
	},
}

// Tag History
var tagHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Get tag history",
	Long:  `Get the history of changes for a specific tag, including previous versions and modifications.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		history, err := client.GetTagHistory(namespace, repository, tagName)
		if err != nil {
			return fmt.Errorf("getting tag history: %w", err)
		}

		fmt.Printf("History for tag %s/%s:%s\n", namespace, repository, tagName)
		return printJSON(history)
	},
}

// Tag Revert
var tagRevertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Revert tag to a previous state",
	Long:  `Revert a tag to a previous state using its manifest digest.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		tag, err := client.RevertTag(namespace, repository, tagName, manifestDigest)
		if err != nil {
			return fmt.Errorf("reverting tag: %w", err)
		}

		fmt.Printf("Successfully reverted tag %s/%s:%s to manifest %s\n", namespace, repository, tagName, manifestDigest)
		return printJSON(tag)
	},
}

var tagRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a tag from a previous state",
	Long:  `Restore a previously deleted or modified tag using its manifest digest.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.RestoreTag(namespace, repository, tagName, manifestDigest)
		if err != nil {
			return fmt.Errorf("restoring tag: %w", err)
		}

		fmt.Printf("Successfully restored tag %s/%s:%s from manifest %s\n", namespace, repository, tagName, manifestDigest)
		return nil
	},
}

func init() {
	// Add subcommands to tag command
	tagCmd.AddCommand(tagInfoCmd)
	tagCmd.AddCommand(tagUpdateCmd)
	tagCmd.AddCommand(tagDeleteCmd)
	tagCmd.AddCommand(tagHistoryCmd)
	tagCmd.AddCommand(tagRevertCmd)
	tagCmd.AddCommand(tagRestoreCmd)

	// Global tag flags (repository context)
	tagCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	tagCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	tagCmd.PersistentFlags().StringVarP(&tagName, "tag", "T", "", "Name of the tag")

	// Mark global flags as required
	_ = tagCmd.MarkPersistentFlagRequired("namespace")
	_ = tagCmd.MarkPersistentFlagRequired("repository")
	_ = tagCmd.MarkPersistentFlagRequired("tag")

	// Update command specific flags
	tagUpdateCmd.Flags().StringVarP(&tagExpiration, "expiration", "e", "", "Tag expiration date (ISO format)")

	// Delete command specific flags
	tagDeleteCmd.Flags().BoolVar(&confirmTagDeletion, "confirm", false, "Confirm tag deletion")

	// Revert command specific flags
	tagRevertCmd.Flags().StringVarP(&manifestDigest, "manifest", "m", "", "Manifest digest to revert to")
	_ = tagRevertCmd.MarkFlagRequired("manifest")

	// Restore command specific flags
	tagRestoreCmd.Flags().StringVarP(&manifestDigest, "manifest", "m", "", "Manifest digest to restore")
	_ = tagRestoreCmd.MarkFlagRequired("manifest")
}
