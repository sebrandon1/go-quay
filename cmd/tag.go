package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	tagName            string
	tagExpiration      string
	tagPage            int
	tagLimit           int
	manifestDigest     string
	confirmTagDeletion bool
)

// tagCmd represents the tag command group
var tagCmd = &cobra.Command{
	Use:   cmdTag,
	Short: "Repository tag management commands",
	Long: `Commands for managing repository tags including detailed information, updates, deletion, and history.

Available commands:
  list     - List repository tags
  info     - Get detailed tag information
  update   - Update tag metadata
  delete   - Delete a tag
  history  - Get tag history
  revert   - Revert tag to a previous state`,
	Example: `  go-quay get tag info --namespace myorg --repository myapp --tag latest --token "$QUAY_TOKEN"
  go-quay get tag history --namespace myorg --repository myapp --tag latest --token "$QUAY_TOKEN"`,
}

var tagListCmd = &cobra.Command{
	Use:   subcmdList,
	Short: "List repository tags",
	Long: `List tags in a repository. Use --page and --limit to select a page; leave them unset to use Quay's default.

Table columns: TAG, DIGEST, SIZE, LAST MODIFIED, EXPIRATION.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		tags, err := client.ListTagsPage(cmd.Context(), namespace, repository, tagLimit, tagPage, false)
		if err != nil {
			return fmt.Errorf("listing tags: %w", err)
		}
		if outputFormat != outputTable {
			return printJSON(tags)
		}

		var rows [][]string
		if tags != nil {
			rows = make([][]string, 0, len(tags.Tags))
			for _, tag := range tags.Tags {
				rows = append(rows, []string{
					tag.Name,
					tag.ManifestDigest,
					fmt.Sprintf("%d", tag.Size),
					tag.LastModified,
					tag.Expiration,
				})
			}
		}
		return printTable([]string{"TAG", "DIGEST", "SIZE", "LAST MODIFIED", "EXPIRATION"}, rows)
	},
}

// Tag Info
var tagInfoCmd = &cobra.Command{
	Use:     subcmdInfo,
	Short:   "Get detailed tag information",
	Long:    `Get detailed information about a specific tag including metadata, manifest digest, and size.`,
	PreRunE: requireTagName,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		tag, err := client.GetTag(cmd.Context(), namespace, repository, tagName)
		if err != nil {
			return fmt.Errorf("getting tag information: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Tag information for %s/%s:%s\n", namespace, repository, tagName)
		return printJSON(tag)
	},
}

// Tag Update
var tagUpdateCmd = &cobra.Command{
	Use:     subcmdUpdate,
	Short:   "Update tag metadata",
	Long:    `Update tag metadata such as expiration date.`,
	PreRunE: requireTagName,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		tag, err := client.UpdateTag(cmd.Context(), namespace, repository, tagName, tagExpiration)
		if err != nil {
			return fmt.Errorf("updating tag: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully updated tag %s/%s:%s\n", namespace, repository, tagName)
		return printJSON(tag)
	},
}

// Tag Delete
var tagDeleteCmd = &cobra.Command{
	Use:     subcmdDelete,
	Short:   "Delete a tag",
	Long:    `Delete a specific tag from the repository. This action cannot be undone.`,
	PreRunE: requireTagName,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmTagDeletion {
			return fmt.Errorf("are you sure you want to delete tag %s/%s:%s? This action cannot be undone.\nUse --confirm to proceed with deletion", namespace, repository, tagName)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteTag(cmd.Context(), namespace, repository, tagName)
		if err != nil {
			return fmt.Errorf("deleting tag: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully deleted tag %s/%s:%s\n", namespace, repository, tagName)
		return nil
	},
}

// Tag History
var tagHistoryCmd = &cobra.Command{
	Use:     "history",
	Short:   "Get tag history",
	Long:    `Get the history of changes for a specific tag, including previous versions and modifications.`,
	PreRunE: requireTagName,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		history, err := client.GetTagHistory(cmd.Context(), namespace, repository, tagName)
		if err != nil {
			return fmt.Errorf("getting tag history: %w", err)
		}

		fmt.Fprintf(os.Stderr, "History for tag %s/%s:%s\n", namespace, repository, tagName)
		return printJSON(history)
	},
}

// Tag Revert
var tagRevertCmd = &cobra.Command{
	Use:     "revert",
	Short:   "Revert tag to a previous state",
	Long:    `Revert a tag to a previous state using its manifest digest.`,
	PreRunE: requireTagName,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		tag, err := client.RevertTag(cmd.Context(), namespace, repository, tagName, manifestDigest)
		if err != nil {
			return fmt.Errorf("reverting tag: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully reverted tag %s/%s:%s to manifest %s\n", namespace, repository, tagName, manifestDigest)
		return printJSON(tag)
	},
}

// Tag Change (create/move)
var tagChangeCmd = &cobra.Command{
	Use:     "change",
	Short:   "Create or move a tag to a manifest",
	Long:    `Create a new tag or move an existing tag to point at a specific manifest digest.`,
	PreRunE: requireTagName,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.ChangeTag(cmd.Context(), namespace, repository, tagName, manifestDigest)
		if err != nil {
			return fmt.Errorf("changing tag: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully changed tag %s/%s:%s to manifest %s\n", namespace, repository, tagName, manifestDigest)
		return nil
	},
}

var tagRestoreCmd = &cobra.Command{
	Use:     "restore",
	Short:   "Restore a tag from a previous state",
	Long:    `Restore a previously deleted or modified tag using its manifest digest.`,
	PreRunE: requireTagName,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.RestoreTag(cmd.Context(), namespace, repository, tagName, manifestDigest)
		if err != nil {
			return fmt.Errorf("restoring tag: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully restored tag %s/%s:%s from manifest %s\n", namespace, repository, tagName, manifestDigest)
		return nil
	},
}

func init() {
	// Add subcommands to tag command
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagInfoCmd)
	tagCmd.AddCommand(tagUpdateCmd)
	tagCmd.AddCommand(tagDeleteCmd)
	tagCmd.AddCommand(tagHistoryCmd)
	tagCmd.AddCommand(tagRevertCmd)
	tagCmd.AddCommand(tagRestoreCmd)
	tagCmd.AddCommand(tagChangeCmd)

	// Global tag flags (repository context)
	tagCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", appCfg.Namespace, "Name of the namespace (default: config file)")
	tagCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	tagCmd.PersistentFlags().StringVarP(&tagName, "tag", "T", "", "Name of the tag")

	// Mark namespace as required only when no config default is available
	if appCfg.Namespace == "" {
		_ = tagCmd.MarkPersistentFlagRequired("namespace")
	}
	_ = tagCmd.MarkPersistentFlagRequired("repository")

	// Update command specific flags
	tagUpdateCmd.Flags().StringVarP(&tagExpiration, "expiration", "e", "", "Tag expiration date (ISO format)")

	// Tag list pagination flags (0 leaves the API default in effect).
	tagListCmd.Flags().IntVar(&tagPage, "page", 0, "Page number (default: Quay API default)")
	tagListCmd.Flags().IntVar(&tagLimit, "limit", 0, "Maximum results per page (default: Quay API default)")

	// Delete command specific flags
	tagDeleteCmd.Flags().BoolVar(&confirmTagDeletion, "confirm", false, "Confirm tag deletion")

	// Revert command specific flags
	tagRevertCmd.Flags().StringVarP(&manifestDigest, "manifest", "m", "", "Manifest digest to revert to")
	_ = tagRevertCmd.MarkFlagRequired("manifest")

	// Restore command specific flags
	tagRestoreCmd.Flags().StringVarP(&manifestDigest, "manifest", "m", "", "Manifest digest to restore")
	_ = tagRestoreCmd.MarkFlagRequired("manifest")

	// Change command specific flags
	tagChangeCmd.Flags().StringVarP(&manifestDigest, "manifest", "m", "", "Manifest digest to assign to the tag")
	_ = tagChangeCmd.MarkFlagRequired("manifest")
}

func requireTagName(_ *cobra.Command, _ []string) error {
	if tagName == "" {
		return fmt.Errorf("tag is required; provide --tag")
	}
	return nil
}
