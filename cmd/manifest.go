package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	manifestRef             string
	labelID                 string
	labelKey                string
	labelValue              string
	labelMediaType          string
	confirmManifestDeletion bool
)

// manifestCmd represents the manifest command group
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Container image manifest management commands",
	Long: `Commands for managing container image manifests including inspection, labels, and deletion.

Available commands:
  info     - Get detailed manifest information
  delete   - Delete a manifest
  labels   - List manifest labels
  label    - Get a specific manifest label
  add-label    - Add a label to a manifest
  remove-label - Remove a label from a manifest`,
}

// Manifest Info
var manifestInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get detailed manifest information",
	Long:  `Get detailed information about a specific manifest including layers, config, and metadata.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		manifest, err := client.GetManifest(namespace, repository, manifestRef)
		if err != nil {
			fmt.Printf("Error getting manifest information: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Manifest information for %s/%s@%s\n", namespace, repository, manifestRef)
		printJSON(manifest)
	},
}

// Manifest Delete
var manifestDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a manifest",
	Long:  `Delete a specific manifest from the repository. This action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmManifestDeletion {
			fmt.Printf("Are you sure you want to delete manifest %s/%s@%s? This action cannot be undone.\n", namespace, repository, manifestRef)
			fmt.Print("Use --confirm to proceed with deletion.\n")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.DeleteManifest(namespace, repository, manifestRef)
		if err != nil {
			fmt.Printf("Error deleting manifest: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted manifest %s/%s@%s\n", namespace, repository, manifestRef)
	},
}

// Manifest Labels (list all)
var manifestLabelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "List manifest labels",
	Long:  `List all labels associated with a specific manifest.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		labels, err := client.GetManifestLabels(namespace, repository, manifestRef)
		if err != nil {
			fmt.Printf("Error getting manifest labels: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Labels for manifest %s/%s@%s\n", namespace, repository, manifestRef)
		printJSON(labels)
	},
}

// Manifest Label (get specific)
var manifestLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "Get a specific manifest label",
	Long:  `Get detailed information about a specific label on a manifest.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		label, err := client.GetManifestLabel(namespace, repository, manifestRef, labelID)
		if err != nil {
			fmt.Printf("Error getting manifest label: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Label %s for manifest %s/%s@%s\n", labelID, namespace, repository, manifestRef)
		printJSON(label)
	},
}

// Add Manifest Label
var manifestAddLabelCmd = &cobra.Command{
	Use:   "add-label",
	Short: "Add a label to a manifest",
	Long:  `Add a new label with a key-value pair to a specific manifest.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		label, err := client.AddManifestLabel(namespace, repository, manifestRef, labelKey, labelValue, labelMediaType)
		if err != nil {
			fmt.Printf("Error adding manifest label: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully added label to manifest %s/%s@%s\n", namespace, repository, manifestRef)
		printJSON(label)
	},
}

// Remove Manifest Label
var manifestRemoveLabelCmd = &cobra.Command{
	Use:   "remove-label",
	Short: "Remove a label from a manifest",
	Long:  `Remove a specific label from a manifest by its label ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.DeleteManifestLabel(namespace, repository, manifestRef, labelID)
		if err != nil {
			fmt.Printf("Error removing manifest label: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully removed label %s from manifest %s/%s@%s\n", labelID, namespace, repository, manifestRef)
	},
}

func init() {
	// Add subcommands to manifest command
	manifestCmd.AddCommand(manifestInfoCmd)
	manifestCmd.AddCommand(manifestDeleteCmd)
	manifestCmd.AddCommand(manifestLabelsCmd)
	manifestCmd.AddCommand(manifestLabelCmd)
	manifestCmd.AddCommand(manifestAddLabelCmd)
	manifestCmd.AddCommand(manifestRemoveLabelCmd)

	// Global manifest flags (repository context)
	manifestCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	manifestCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	manifestCmd.PersistentFlags().StringVarP(&manifestRef, "manifest", "m", "", "Manifest reference (digest like sha256:...)")
	manifestCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")

	// Mark global flags as required
	if err := manifestCmd.MarkPersistentFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := manifestCmd.MarkPersistentFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := manifestCmd.MarkPersistentFlagRequired("manifest"); err != nil {
		fmt.Printf("Error marking manifest flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := manifestCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}

	// Delete command specific flags
	manifestDeleteCmd.Flags().BoolVar(&confirmManifestDeletion, "confirm", false, "Confirm manifest deletion")

	// Label command specific flags (for getting a specific label)
	manifestLabelCmd.Flags().StringVarP(&labelID, "label-id", "l", "", "Label ID")
	if err := manifestLabelCmd.MarkFlagRequired("label-id"); err != nil {
		fmt.Printf("Error marking label-id flag as required: %v\n", err)
		os.Exit(1)
	}

	// Add label command specific flags
	manifestAddLabelCmd.Flags().StringVarP(&labelKey, "key", "k", "", "Label key")
	manifestAddLabelCmd.Flags().StringVarP(&labelValue, "value", "v", "", "Label value")
	manifestAddLabelCmd.Flags().StringVar(&labelMediaType, "media-type", "", "Label media type (optional, defaults to text/plain)")
	if err := manifestAddLabelCmd.MarkFlagRequired("key"); err != nil {
		fmt.Printf("Error marking key flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := manifestAddLabelCmd.MarkFlagRequired("value"); err != nil {
		fmt.Printf("Error marking value flag as required: %v\n", err)
		os.Exit(1)
	}

	// Remove label command specific flags
	manifestRemoveLabelCmd.Flags().StringVarP(&labelID, "label-id", "l", "", "Label ID to remove")
	if err := manifestRemoveLabelCmd.MarkFlagRequired("label-id"); err != nil {
		fmt.Printf("Error marking label-id flag as required: %v\n", err)
		os.Exit(1)
	}
}
