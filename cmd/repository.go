package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var (
	repoVisibility  string
	repoDescription string
	confirmDeletion bool
	repoPublic      bool
	repoStarred     bool
	repoPopularity  bool
	repoTable       bool
	repoPage        int
	repoLimit       int
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
	Use:   subcmdInfo,
	Short: "Get repository information from Quay.io",
	RunE: func(cmd *cobra.Command, args []string) error {
		if repository == "" {
			return fmt.Errorf("--repository is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		repo, err := client.GetRepository(namespace, repository)
		if err != nil {
			return fmt.Errorf("getting repository: %w", err)
		}

		return printJSON(repo)
	},
}

// Repository Creation
var repoCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
	Short: "Create a new repository",
	Long:  `Create a new repository in the specified namespace with optional description and visibility settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if repository == "" {
			return fmt.Errorf("--repository is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		repo, err := client.CreateRepository(namespace, repository, repoVisibility, repoDescription)
		if err != nil {
			return fmt.Errorf("creating repository: %w", err)
		}

		fmt.Printf("Successfully created repository %s/%s\n", namespace, repository)
		return printJSON(repo)
	},
}

// Repository Update
var repoUpdateCmd = &cobra.Command{
	Use:   subcmdUpdate,
	Short: "Update repository settings",
	Long:  `Update repository description and/or visibility settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if repository == "" {
			return fmt.Errorf("--repository is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		repo, err := client.UpdateRepository(namespace, repository, repoDescription, repoVisibility)
		if err != nil {
			return fmt.Errorf("updating repository: %w", err)
		}

		fmt.Printf("Successfully updated repository %s/%s\n", namespace, repository)
		return printJSON(repo)
	},
}

// Repository Deletion
var repoDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a repository",
	Long:  `Delete a repository. This action is irreversible and will remove all images and tags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if repository == "" {
			return fmt.Errorf("--repository is required")
		}

		if !confirmDeletion {
			return fmt.Errorf("are you sure you want to delete repository %s/%s? This action cannot be undone.\nUse --confirm to proceed with deletion", namespace, repository)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeleteRepository(namespace, repository)
		if err != nil {
			return fmt.Errorf("deleting repository: %w", err)
		}

		fmt.Printf("Successfully deleted repository %s/%s\n", namespace, repository)
		return nil
	},
}

var repoListCmd = &cobra.Command{
	Use:   subcmdList,
	Short: "List repositories in a namespace",
	Long: `List all repositories in a namespace with optional filters.

Use --table with --popularity for an enriched dashboard sorted by pull count:

  go-quay get repository list -n myorg --public --popularity --table

  REPOSITORY                  PULLS  PUSHES (30d)  TAGS  LATEST TAG  LAST PUSH   MULTI-ARCH
  certsuite-sample-workload   44594  2             2     tag1        2026-06-15  yes
  certsuite-probe             14226  5             42    latest      2026-06-15  yes
  kube-rbac-proxy             4052   0             1     v0.13.1     2026-03-13  no
  certsuite                   774    6             51    unstable    2026-06-15  yes
  collector                   32     4             46    unstable    2026-06-15  no`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		repos, err := client.ListRepositories(namespace, repoPublic, repoStarred, repoPopularity, repoPage, repoLimit)
		if err != nil {
			return fmt.Errorf("listing repositories: %w", err)
		}

		if !repoTable {
			return printJSON(repos)
		}

		sort.Slice(repos.Repositories, func(i, j int) bool {
			return repos.Repositories[i].Popularity > repos.Repositories[j].Popularity
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "REPOSITORY\tPULLS\tPUSHES (30d)\tTAGS\tLATEST TAG\tLAST PUSH\tMULTI-ARCH")

		thirtyDaysAgo := time.Now().UTC().AddDate(0, 0, -30).Unix()

		for _, repo := range repos.Repositories {
			latestTag := "-"
			lastPush := "-"
			multiArch := "-"
			tagDisplay := "0"
			recentPushes := 0

			tags, err := client.ListTags(namespace, repo.Name, 100, true)
			if err == nil && len(tags.Tags) > 0 {
				for _, tag := range tags.Tags {
					if strings.HasPrefix(tag.Name, "sha256-") || strings.HasSuffix(tag.Name, ".sig") || strings.HasSuffix(tag.Name, ".att") || strings.HasSuffix(tag.Name, ".sbom") {
						continue
					}
					latestTag = tag.Name
					multiArch = "no"
					if tag.IsManifestList {
						multiArch = "yes"
					}
					if tag.StartTs > 0 {
						lastPush = time.Unix(tag.StartTs, 0).UTC().Format("2006-01-02")
					}
					break
				}
				if tags.HasAdditional {
					tagDisplay = fmt.Sprintf("%d+", len(tags.Tags))
				} else {
					tagDisplay = fmt.Sprintf("%d", len(tags.Tags))
				}
				for _, tag := range tags.Tags {
					if tag.StartTs >= thirtyDaysAgo {
						recentPushes++
					}
				}
			}

			fmt.Fprintf(w, "%s\t%.0f\t%d\t%s\t%s\t%s\t%s\n",
				repo.Name, repo.Popularity, recentPushes, tagDisplay, latestTag, lastPush, multiArch)
		}

		if err := w.Flush(); err != nil {
			return fmt.Errorf("flushing table output: %w", err)
		}
		return nil
	},
}

var repoChangeVisibilityCmd = &cobra.Command{
	Use:   "change-visibility",
	Short: "Change repository visibility",
	Long:  `Change a repository's visibility between public and private.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if repository == "" {
			return fmt.Errorf("--repository is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.ChangeRepositoryVisibility(namespace, repository, repoVisibility)
		if err != nil {
			return fmt.Errorf("changing repository visibility: %w", err)
		}

		fmt.Printf("Successfully changed visibility of %s/%s to %s\n", namespace, repository, repoVisibility)
		return nil
	},
}

func init() {
	// Add subcommands to repository command
	repositoryCmd.AddCommand(repoInfoCmd)
	repositoryCmd.AddCommand(repoCreateCmd)
	repositoryCmd.AddCommand(repoUpdateCmd)
	repositoryCmd.AddCommand(repoDeleteCmd)
	repositoryCmd.AddCommand(repoListCmd)
	repositoryCmd.AddCommand(repoChangeVisibilityCmd)

	// Global repository flags
	repositoryCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	repositoryCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")

	// Mark global flags as required
	_ = repositoryCmd.MarkPersistentFlagRequired("namespace")

	// Create command specific flags
	repoCreateCmd.Flags().StringVarP(&repoVisibility, "visibility", "v", "private", "Repository visibility (private/public)")
	repoCreateCmd.Flags().StringVarP(&repoDescription, "description", "d", "", "Repository description")

	// Update command specific flags
	repoUpdateCmd.Flags().StringVarP(&repoVisibility, "visibility", "v", "", "Repository visibility (private/public)")
	repoUpdateCmd.Flags().StringVarP(&repoDescription, "description", "d", "", "Repository description")

	// Delete command specific flags
	repoDeleteCmd.Flags().BoolVar(&confirmDeletion, "confirm", false, "Confirm repository deletion")

	// List command specific flags
	repoListCmd.Flags().BoolVar(&repoPublic, "public", true, "Include public repositories")
	repoListCmd.Flags().BoolVar(&repoStarred, "starred", false, "Only starred repositories")
	repoListCmd.Flags().BoolVar(&repoPopularity, "popularity", false, "Include pull popularity scores")
	repoListCmd.Flags().BoolVar(&repoTable, "table", false, "Display results as a formatted table with tag details")
	repoListCmd.Flags().IntVar(&repoPage, "page", 1, "Page number")
	repoListCmd.Flags().IntVar(&repoLimit, "limit", 10, "Maximum results per page")

	// Change-visibility command specific flags
	repoChangeVisibilityCmd.Flags().StringVarP(&repoVisibility, "visibility", "v", "", "New visibility (private/public)")
	_ = repoChangeVisibilityCmd.MarkFlagRequired("visibility")
}
