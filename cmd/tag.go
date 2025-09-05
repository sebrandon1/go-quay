package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
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
	Use:   "info",
	Short: "Get detailed tag information",
	Long:  `Get detailed information about a specific tag including metadata, manifest digest, and size.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		tag, err := client.GetTag(namespace, repository, tagName)
		if err != nil {
			fmt.Printf("Error getting tag information: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Tag information for %s/%s:%s\n", namespace, repository, tagName)
		printJSON(tag)
	},
}

// Tag Update
var tagUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update tag metadata",
	Long:  `Update tag metadata such as expiration date.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		tag, err := client.UpdateTag(namespace, repository, tagName, tagExpiration)
		if err != nil {
			fmt.Printf("Error updating tag: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated tag %s/%s:%s\n", namespace, repository, tagName)
		printJSON(tag)
	},
}

// Tag Delete
var tagDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a tag",
	Long:  `Delete a specific tag from the repository. This action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmTagDeletion {
			fmt.Printf("Are you sure you want to delete tag %s/%s:%s? This action cannot be undone.\n", namespace, repository, tagName)
			fmt.Print("Use --confirm to proceed with deletion.\n")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.DeleteTag(namespace, repository, tagName)
		if err != nil {
			fmt.Printf("Error deleting tag: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted tag %s/%s:%s\n", namespace, repository, tagName)
	},
}

// Tag History
var tagHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Get tag history",
	Long:  `Get the history of changes for a specific tag, including previous versions and modifications.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		history, err := client.GetTagHistory(namespace, repository, tagName)
		if err != nil {
			fmt.Printf("Error getting tag history: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("History for tag %s/%s:%s\n", namespace, repository, tagName)
		printJSON(history)
	},
}

// Tag Revert
var tagRevertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Revert tag to a previous state",
	Long:  `Revert a tag to a previous state using its manifest digest.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		tag, err := client.RevertTag(namespace, repository, tagName, manifestDigest)
		if err != nil {
			fmt.Printf("Error reverting tag: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully reverted tag %s/%s:%s to manifest %s\n", namespace, repository, tagName, manifestDigest)
		printJSON(tag)
	},
}

func init() {
	// Add subcommands to tag command
	tagCmd.AddCommand(tagInfoCmd)
	tagCmd.AddCommand(tagUpdateCmd)
	tagCmd.AddCommand(tagDeleteCmd)
	tagCmd.AddCommand(tagHistoryCmd)
	tagCmd.AddCommand(tagRevertCmd)

	// Global tag flags (repository context)
	tagCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	tagCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	tagCmd.PersistentFlags().StringVarP(&tagName, "tag", "T", "", "Name of the tag")
	tagCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")

	// Mark global flags as required
	if err := tagCmd.MarkPersistentFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := tagCmd.MarkPersistentFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := tagCmd.MarkPersistentFlagRequired("tag"); err != nil {
		fmt.Printf("Error marking tag flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := tagCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}

	// Update command specific flags
	tagUpdateCmd.Flags().StringVarP(&tagExpiration, "expiration", "e", "", "Tag expiration date (ISO format)")

	// Delete command specific flags
	tagDeleteCmd.Flags().BoolVar(&confirmTagDeletion, "confirm", false, "Confirm tag deletion")

	// Revert command specific flags
	tagRevertCmd.Flags().StringVarP(&manifestDigest, "manifest", "m", "", "Manifest digest to revert to")
	if err := tagRevertCmd.MarkFlagRequired("manifest"); err != nil {
		fmt.Printf("Error marking manifest flag as required: %v\n", err)
		os.Exit(1)
	}
}
