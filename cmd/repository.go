package cmd

import (
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	repoVisibility  string
	repoDescription string
	confirmDeletion bool
)

// repositoryCmd represents the repository command group
var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Repository management commands",
	Long: `Commands for managing repositories including creation, updates, deletion, and information retrieval.

Available commands:
  info     - Get repository information (default)
  create   - Create a new repository
  update   - Update repository settings
  delete   - Delete a repository`,
}

// Repository Info (existing functionality)
var repoInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get repository information from Quay.io",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		repo, err := client.GetRepository(namespace, repository)
		if err != nil {
			fmt.Printf("Error getting repository: %v\n", err)
			os.Exit(1)
		}

		printJSON(repo)
	},
}

// Repository Creation
var repoCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new repository",
	Long:  `Create a new repository in the specified namespace with optional description and visibility settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		repo, err := client.CreateRepository(namespace, repository, repoVisibility, repoDescription)
		if err != nil {
			fmt.Printf("Error creating repository: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully created repository %s/%s\n", namespace, repository)
		printJSON(repo)
	},
}

// Repository Update
var repoUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update repository settings",
	Long:  `Update repository description and/or visibility settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		repo, err := client.UpdateRepository(namespace, repository, repoDescription, repoVisibility)
		if err != nil {
			fmt.Printf("Error updating repository: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated repository %s/%s\n", namespace, repository)
		printJSON(repo)
	},
}

// Repository Deletion
var repoDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a repository",
	Long:  `Delete a repository. This action is irreversible and will remove all images and tags.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmDeletion {
			fmt.Printf("Are you sure you want to delete repository %s/%s? This action cannot be undone.\n", namespace, repository)
			fmt.Print("Use --confirm to proceed with deletion.\n")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.DeleteRepository(namespace, repository)
		if err != nil {
			fmt.Printf("Error deleting repository: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully deleted repository %s/%s\n", namespace, repository)
	},
}

func init() {
	// Add subcommands to repository command
	repositoryCmd.AddCommand(repoInfoCmd)
	repositoryCmd.AddCommand(repoCreateCmd)
	repositoryCmd.AddCommand(repoUpdateCmd)
	repositoryCmd.AddCommand(repoDeleteCmd)

	// Global repository flags
	repositoryCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	repositoryCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	repositoryCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")

	// Mark global flags as required
	if err := repositoryCmd.MarkPersistentFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := repositoryCmd.MarkPersistentFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := repositoryCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}

	// Create command specific flags
	repoCreateCmd.Flags().StringVarP(&repoVisibility, "visibility", "v", "private", "Repository visibility (private/public)")
	repoCreateCmd.Flags().StringVarP(&repoDescription, "description", "d", "", "Repository description")

	// Update command specific flags
	repoUpdateCmd.Flags().StringVarP(&repoVisibility, "visibility", "v", "", "Repository visibility (private/public)")
	repoUpdateCmd.Flags().StringVarP(&repoDescription, "description", "d", "", "Repository description")

	// Delete command specific flags
	repoDeleteCmd.Flags().BoolVar(&confirmDeletion, "confirm", false, "Confirm repository deletion")

	// Set default behavior to info command when no subcommand is specified
	repositoryCmd.Run = repoInfoCmd.Run
}
