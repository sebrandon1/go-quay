package cmd

import (
	"fmt"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	buildNamespace     string
	buildRepository    string
	buildUUID          string
	buildLimit         int
	buildArchiveURL    string
	buildDockerfile    string
	buildSubdirectory  string
	buildTags          []string
	confirmBuildCancel bool
)

// buildCmd represents the build command group
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Repository build management commands",
	Long: `Commands for managing repository builds.

Builds allow automated image creation from Dockerfiles stored in git repositories
or uploaded archives.

Available commands:
  list    - List builds for a repository
  info    - Get build details
  logs    - Get build logs
  request - Request a new build
  cancel  - Cancel an ongoing build`,
}

// Build List
var buildListCmd = &cobra.Command{
	Use:   subcmdList,
	Short: "List builds for a repository",
	Long:  `List all builds for the specified repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		builds, err := client.GetBuilds(buildNamespace, buildRepository, buildLimit)
		if err != nil {
			return fmt.Errorf("getting builds: %w", err)
		}

		fmt.Printf("Builds for %s/%s:\n", buildNamespace, buildRepository)
		return printJSON(builds)
	},
}

// Build Info
var buildInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get build details",
	Long:  `Get detailed information about a specific build.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		build, err := client.GetBuild(buildNamespace, buildRepository, buildUUID)
		if err != nil {
			return fmt.Errorf("getting build: %w", err)
		}

		fmt.Printf("Build: %s\n", buildUUID)
		return printJSON(build)
	},
}

// Build Logs
var buildLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get build logs",
	Long:  `Get the logs for a specific build.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		logs, err := client.GetBuildLogs(buildNamespace, buildRepository, buildUUID)
		if err != nil {
			return fmt.Errorf("getting build logs: %w", err)
		}

		fmt.Printf("Logs for build %s:\n", buildUUID)
		return printJSON(logs)
	},
}

// Build Request
var buildRequestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request a new build",
	Long: `Request a new build from an archive URL.

The archive should be a tar.gz file containing a Dockerfile and any necessary build context.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		buildReq := &lib.RequestBuildRequest{
			ArchiveURL:     buildArchiveURL,
			DockerfilePath: buildDockerfile,
			Subdirectory:   buildSubdirectory,
			Tags:           buildTags,
		}

		build, err := client.RequestBuild(buildNamespace, buildRepository, buildReq)
		if err != nil {
			return fmt.Errorf("requesting build: %w", err)
		}

		fmt.Printf("Build requested successfully!\n")
		return printJSON(build)
	},
}

// Build Cancel
var buildCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel an ongoing build",
	Long:  `Cancel an ongoing build. This action cannot be undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirmBuildCancel {
			return fmt.Errorf("are you sure you want to cancel build '%s'? This action cannot be undone.\nUse --confirm to proceed with cancellation", buildUUID)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.CancelBuild(buildNamespace, buildRepository, buildUUID)
		if err != nil {
			return fmt.Errorf("canceling build: %w", err)
		}

		fmt.Printf("Build '%s' canceled successfully.\n", buildUUID)
		return nil
	},
}

var buildStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get build status",
	Long:  `Get the current status of a specific build.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		status, err := client.GetBuildStatus(buildNamespace, buildRepository, buildUUID)
		if err != nil {
			return fmt.Errorf("getting build status: %w", err)
		}

		fmt.Printf("Status for build %s:\n", buildUUID)
		return printJSON(status)
	},
}

func init() {
	// Add subcommands to build command
	buildCmd.AddCommand(buildListCmd)
	buildCmd.AddCommand(buildInfoCmd)
	buildCmd.AddCommand(buildLogsCmd)
	buildCmd.AddCommand(buildRequestCmd)
	buildCmd.AddCommand(buildCancelCmd)
	buildCmd.AddCommand(buildStatusCmd)

	// Global build flags
	buildCmd.PersistentFlags().StringVarP(&buildNamespace, "namespace", "n", "", "Repository namespace")
	buildCmd.PersistentFlags().StringVarP(&buildRepository, "repository", "r", "", "Repository name")
	_ = buildCmd.MarkPersistentFlagRequired("namespace")
	_ = buildCmd.MarkPersistentFlagRequired("repository")

	initBuildListFlags()
	initBuildInfoFlags()
	initBuildRequestFlags()
	initBuildCancelFlags()
	initBuildStatusFlags()
}

func initBuildStatusFlags() {
	buildStatusCmd.Flags().StringVarP(&buildUUID, "uuid", "u", "", "Build UUID")
	_ = buildStatusCmd.MarkFlagRequired("uuid")
}

func initBuildListFlags() {
	buildListCmd.Flags().IntVarP(&buildLimit, "limit", "l", 10, "Maximum number of builds to return")
}

func initBuildInfoFlags() {
	buildInfoCmd.Flags().StringVarP(&buildUUID, "uuid", "u", "", "Build UUID")
	_ = buildInfoCmd.MarkFlagRequired("uuid")
}

func initBuildRequestFlags() {
	buildRequestCmd.Flags().StringVarP(&buildArchiveURL, "archive-url", "a", "", "URL to archive containing Dockerfile")
	buildRequestCmd.Flags().StringVarP(&buildDockerfile, "dockerfile", "d", "", "Path to Dockerfile within archive")
	buildRequestCmd.Flags().StringVarP(&buildSubdirectory, "subdirectory", "s", "", "Subdirectory containing build context")
	buildRequestCmd.Flags().StringSliceVar(&buildTags, "tag", []string{"latest"}, "Tags for the built image")
	_ = buildRequestCmd.MarkFlagRequired("archive-url")
}

func initBuildCancelFlags() {
	buildCancelCmd.Flags().StringVarP(&buildUUID, "uuid", "u", "", "Build UUID")
	buildCancelCmd.Flags().BoolVar(&confirmBuildCancel, "confirm", false, "Confirm build cancellation")
	_ = buildCancelCmd.MarkFlagRequired("uuid")
}
