package cmd

import (
	"fmt"
	"os"

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
	Use:   "list",
	Short: "List builds for a repository",
	Long:  `List all builds for the specified repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		builds, err := client.GetBuilds(buildNamespace, buildRepository, buildLimit)
		if err != nil {
			fmt.Printf("Error getting builds: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Builds for %s/%s:\n", buildNamespace, buildRepository)
		printJSON(builds)
	},
}

// Build Info
var buildInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get build details",
	Long:  `Get detailed information about a specific build.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		build, err := client.GetBuild(buildNamespace, buildRepository, buildUUID)
		if err != nil {
			fmt.Printf("Error getting build: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Build: %s\n", buildUUID)
		printJSON(build)
	},
}

// Build Logs
var buildLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get build logs",
	Long:  `Get the logs for a specific build.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		logs, err := client.GetBuildLogs(buildNamespace, buildRepository, buildUUID)
		if err != nil {
			fmt.Printf("Error getting build logs: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Logs for build %s:\n", buildUUID)
		printJSON(logs)
	},
}

// Build Request
var buildRequestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request a new build",
	Long: `Request a new build from an archive URL.

The archive should be a tar.gz file containing a Dockerfile and any necessary build context.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		buildReq := &lib.RequestBuildRequest{
			ArchiveURL:     buildArchiveURL,
			DockerfilePath: buildDockerfile,
			Subdirectory:   buildSubdirectory,
			Tags:           buildTags,
		}

		build, err := client.RequestBuild(buildNamespace, buildRepository, buildReq)
		if err != nil {
			fmt.Printf("Error requesting build: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Build requested successfully!\n")
		printJSON(build)
	},
}

// Build Cancel
var buildCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel an ongoing build",
	Long:  `Cancel an ongoing build. This action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmBuildCancel {
			fmt.Printf("Are you sure you want to cancel build '%s'? This action cannot be undone.\n", buildUUID)
			fmt.Println("Use --confirm to proceed with cancellation.")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.CancelBuild(buildNamespace, buildRepository, buildUUID)
		if err != nil {
			fmt.Printf("Error canceling build: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Build '%s' canceled successfully.\n", buildUUID)
	},
}

func init() {
	// Add subcommands to build command
	buildCmd.AddCommand(buildListCmd)
	buildCmd.AddCommand(buildInfoCmd)
	buildCmd.AddCommand(buildLogsCmd)
	buildCmd.AddCommand(buildRequestCmd)
	buildCmd.AddCommand(buildCancelCmd)

	// Global build flags
	buildCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")
	buildCmd.PersistentFlags().StringVarP(&buildNamespace, "namespace", "n", "", "Repository namespace")
	buildCmd.PersistentFlags().StringVarP(&buildRepository, "repository", "r", "", "Repository name")
	markBuildFlagRequired(buildCmd.MarkPersistentFlagRequired("token"))
	markBuildFlagRequired(buildCmd.MarkPersistentFlagRequired("namespace"))
	markBuildFlagRequired(buildCmd.MarkPersistentFlagRequired("repository"))

	initBuildListFlags()
	initBuildInfoFlags()
	initBuildRequestFlags()
	initBuildCancelFlags()
}

func initBuildListFlags() {
	buildListCmd.Flags().IntVarP(&buildLimit, "limit", "l", 10, "Maximum number of builds to return")
}

func initBuildInfoFlags() {
	buildInfoCmd.Flags().StringVarP(&buildUUID, "uuid", "u", "", "Build UUID")
	markBuildFlagRequired(buildInfoCmd.MarkFlagRequired("uuid"))
}

func initBuildRequestFlags() {
	buildRequestCmd.Flags().StringVarP(&buildArchiveURL, "archive-url", "a", "", "URL to archive containing Dockerfile")
	buildRequestCmd.Flags().StringVarP(&buildDockerfile, "dockerfile", "d", "", "Path to Dockerfile within archive")
	buildRequestCmd.Flags().StringVarP(&buildSubdirectory, "subdirectory", "s", "", "Subdirectory containing build context")
	buildRequestCmd.Flags().StringSliceVar(&buildTags, "tag", []string{"latest"}, "Tags for the built image")
	markBuildFlagRequired(buildRequestCmd.MarkFlagRequired("archive-url"))
}

func initBuildCancelFlags() {
	buildCancelCmd.Flags().StringVarP(&buildUUID, "uuid", "u", "", "Build UUID")
	buildCancelCmd.Flags().BoolVar(&confirmBuildCancel, "confirm", false, "Confirm build cancellation")
	markBuildFlagRequired(buildCancelCmd.MarkFlagRequired("uuid"))
}

func markBuildFlagRequired(err error) {
	if err != nil {
		fmt.Printf("Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
