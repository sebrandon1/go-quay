package cmd

import (
	"fmt"

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
	Use:   subcmdInfo,
	Short: "Get detailed manifest information",
	Long:  `Get detailed information about a specific manifest including layers, config, and metadata.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		manifest, err := client.GetManifest(namespace, repository, manifestRef)
		if err != nil {
			return fmt.Errorf("getting manifest information: %w", err)
		}

		fmt.Printf("Manifest information for %s/%s@%s\n", namespace, repository, manifestRef)
		return printJSON(manifest)
	},
}

// Manifest Delete
var manifestDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a manifest",
	Long:  `Delete a specific manifest from the repository. This action cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmManifestDeletion {
			return fmt.Errorf("are you sure you want to delete manifest %s/%s@%s? This action cannot be undone.\nUse --confirm to proceed with deletion", namespace, repository, manifestRef)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteManifest(namespace, repository, manifestRef)
		if err != nil {
			return fmt.Errorf("deleting manifest: %w", err)
		}

		fmt.Printf("Successfully deleted manifest %s/%s@%s\n", namespace, repository, manifestRef)
		return nil
	},
}

// Manifest Labels (list all)
var manifestLabelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "List manifest labels",
	Long:  `List all labels associated with a specific manifest.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		labels, err := client.GetManifestLabels(namespace, repository, manifestRef)
		if err != nil {
			return fmt.Errorf("getting manifest labels: %w", err)
		}

		fmt.Printf("Labels for manifest %s/%s@%s\n", namespace, repository, manifestRef)
		return printJSON(labels)
	},
}

// Manifest Label (get specific)
var manifestLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "Get a specific manifest label",
	Long:  `Get detailed information about a specific label on a manifest.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		label, err := client.GetManifestLabel(namespace, repository, manifestRef, labelID)
		if err != nil {
			return fmt.Errorf("getting manifest label: %w", err)
		}

		fmt.Printf("Label %s for manifest %s/%s@%s\n", labelID, namespace, repository, manifestRef)
		return printJSON(label)
	},
}

// Add Manifest Label
var manifestAddLabelCmd = &cobra.Command{
	Use:   "add-label",
	Short: "Add a label to a manifest",
	Long:  `Add a new label with a key-value pair to a specific manifest.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		label, err := client.AddManifestLabel(namespace, repository, manifestRef, labelKey, labelValue, labelMediaType)
		if err != nil {
			return fmt.Errorf("adding manifest label: %w", err)
		}

		fmt.Printf("Successfully added label to manifest %s/%s@%s\n", namespace, repository, manifestRef)
		return printJSON(label)
	},
}

// Remove Manifest Label
var manifestRemoveLabelCmd = &cobra.Command{
	Use:   "remove-label",
	Short: "Remove a label from a manifest",
	Long:  `Remove a specific label from a manifest by its label ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteManifestLabel(namespace, repository, manifestRef, labelID)
		if err != nil {
			return fmt.Errorf("removing manifest label: %w", err)
		}

		fmt.Printf("Successfully removed label %s from manifest %s/%s@%s\n", labelID, namespace, repository, manifestRef)
		return nil
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

	// Mark global flags as required
	_ = manifestCmd.MarkPersistentFlagRequired("namespace")
	_ = manifestCmd.MarkPersistentFlagRequired("repository")
	_ = manifestCmd.MarkPersistentFlagRequired("manifest")

	// Delete command specific flags
	manifestDeleteCmd.Flags().BoolVar(&confirmManifestDeletion, "confirm", false, "Confirm manifest deletion")

	// Label command specific flags (for getting a specific label)
	manifestLabelCmd.Flags().StringVarP(&labelID, "label-id", "l", "", "Label ID")
	_ = manifestLabelCmd.MarkFlagRequired("label-id")

	// Add label command specific flags
	manifestAddLabelCmd.Flags().StringVarP(&labelKey, "key", "k", "", "Label key")
	manifestAddLabelCmd.Flags().StringVarP(&labelValue, "value", "v", "", "Label value")
	manifestAddLabelCmd.Flags().StringVar(&labelMediaType, "media-type", "", "Label media type (optional, defaults to text/plain)")
	_ = manifestAddLabelCmd.MarkFlagRequired("key")
	_ = manifestAddLabelCmd.MarkFlagRequired("value")

	// Remove label command specific flags
	manifestRemoveLabelCmd.Flags().StringVarP(&labelID, "label-id", "l", "", "Label ID to remove")
	_ = manifestRemoveLabelCmd.MarkFlagRequired("label-id")
}
