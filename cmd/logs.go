package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	namespace     string
	repository    string
	startdate     string
	enddate       string
	token         string
	nextPage      string
	starttime     string
	endtime       string
	callbackURL   string
	callbackEmail string
)

var aggregatedLogsCmd = &cobra.Command{
	Use:   "aggregatedlogs",
	Short: "Get aggregated logs from Quay.io",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			return
		}

		logs, err := client.GetAggregatedLogs(namespace, repository, startdate, enddate)
		if err != nil {
			fmt.Println("Error getting aggregated logs:", err)
			return
		}

		jsonOutput, err := json.Marshal(logs)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		fmt.Println(string(jsonOutput))
	},
}

func init() {
	aggregatedLogsCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&startdate, "startdate", "s", "", "Start date for the logs")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&enddate, "enddate", "e", "", "End date for the logs")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")
}

// logsCmd is the parent command group for log management
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Log management commands",
}

var repoLogsCmd = &cobra.Command{
	Use:   "repo-logs",
	Short: "Get repository logs",
	Long:  `Get action logs for a specific repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		logs, err := client.GetLogs(namespace, repository, nextPage)
		if err != nil {
			fmt.Printf("Error getting repository logs: %v\n", err)
			os.Exit(1)
		}

		printJSON(logs)
	},
}

var orgLogsCmd = &cobra.Command{
	Use:   "org-logs",
	Short: "Get organization logs",
	Long:  `Get action logs for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		logs, err := client.GetOrganizationLogs(orgName, nextPage)
		if err != nil {
			fmt.Printf("Error getting organization logs: %v\n", err)
			os.Exit(1)
		}

		printJSON(logs)
	},
}

var orgAggregatedLogsCmd = &cobra.Command{
	Use:   "org-aggregated-logs",
	Short: "Get aggregated organization logs",
	Long:  `Get aggregated action logs for an organization.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		logs, err := client.GetOrganizationAggregatedLogs(orgName, startdate, enddate)
		if err != nil {
			fmt.Printf("Error getting organization aggregated logs: %v\n", err)
			os.Exit(1)
		}

		printJSON(logs)
	},
}

var userLogsCmd = &cobra.Command{
	Use:   "user-logs",
	Short: "Get user logs",
	Long:  `Get action logs for the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		logs, err := client.GetUserLogs(nextPage)
		if err != nil {
			fmt.Printf("Error getting user logs: %v\n", err)
			os.Exit(1)
		}

		printJSON(logs)
	},
}

var userAggregatedLogsCmd = &cobra.Command{
	Use:   "user-aggregated-logs",
	Short: "Get aggregated user logs",
	Long:  `Get aggregated action logs for the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		logs, err := client.GetUserAggregatedLogs(startdate, enddate)
		if err != nil {
			fmt.Printf("Error getting user aggregated logs: %v\n", err)
			os.Exit(1)
		}

		printJSON(logs)
	},
}

var exportOrgLogsCmd = &cobra.Command{
	Use:   "export-org-logs",
	Short: "Export organization logs",
	Long:  `Export action logs for an organization via callback URL or email.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.ExportOrganizationLogs(orgName, &lib.ExportLogsRequest{
			StartTime:   starttime,
			EndTime:     endtime,
			CallbackURL: callbackURL,
			Email:       callbackEmail,
		})
		if err != nil {
			fmt.Printf("Error exporting organization logs: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Organization logs export initiated successfully.")
	},
}

var exportUserLogsCmd = &cobra.Command{
	Use:   "export-user-logs",
	Short: "Export user logs",
	Long:  `Export action logs for the current user via callback URL or email.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.ExportUserLogs(&lib.ExportLogsRequest{
			StartTime:   starttime,
			EndTime:     endtime,
			CallbackURL: callbackURL,
			Email:       callbackEmail,
		})
		if err != nil {
			fmt.Printf("Error exporting user logs: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("User logs export initiated successfully.")
	},
}

var exportRepoLogsCmd = &cobra.Command{
	Use:   "export-repo-logs",
	Short: "Export repository logs",
	Long:  `Export action logs for a repository via callback URL or email.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		err = client.ExportRepositoryLogs(namespace, repository, &lib.ExportLogsRequest{
			StartTime:   starttime,
			EndTime:     endtime,
			CallbackURL: callbackURL,
			Email:       callbackEmail,
		})
		if err != nil {
			fmt.Printf("Error exporting repository logs: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Repository logs export initiated successfully.")
	},
}

func init() {
	logsCmd.AddCommand(repoLogsCmd)
	logsCmd.AddCommand(orgLogsCmd)
	logsCmd.AddCommand(orgAggregatedLogsCmd)
	logsCmd.AddCommand(userLogsCmd)
	logsCmd.AddCommand(userAggregatedLogsCmd)
	logsCmd.AddCommand(exportOrgLogsCmd)
	logsCmd.AddCommand(exportUserLogsCmd)
	logsCmd.AddCommand(exportRepoLogsCmd)

	logsCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")

	// repo-logs flags
	repoLogsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Repository namespace")
	repoLogsCmd.Flags().StringVarP(&repository, "repository", "r", "", "Repository name")
	repoLogsCmd.Flags().StringVar(&nextPage, "next-page", "", "Next page token for pagination")

	// org-logs flags
	orgLogsCmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	orgLogsCmd.Flags().StringVar(&nextPage, "next-page", "", "Next page token for pagination")

	// org-aggregated-logs flags
	orgAggregatedLogsCmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	orgAggregatedLogsCmd.Flags().StringVarP(&startdate, "startdate", "s", "", "Start date")
	orgAggregatedLogsCmd.Flags().StringVarP(&enddate, "enddate", "e", "", "End date")

	// user-logs flags
	userLogsCmd.Flags().StringVar(&nextPage, "next-page", "", "Next page token for pagination")

	// user-aggregated-logs flags
	userAggregatedLogsCmd.Flags().StringVarP(&startdate, "startdate", "s", "", "Start date")
	userAggregatedLogsCmd.Flags().StringVarP(&enddate, "enddate", "e", "", "End date")

	// export-org-logs flags
	exportOrgLogsCmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	exportOrgLogsCmd.Flags().StringVar(&starttime, "start-time", "", "Start time for export")
	exportOrgLogsCmd.Flags().StringVar(&endtime, "end-time", "", "End time for export")
	exportOrgLogsCmd.Flags().StringVar(&callbackURL, "callback-url", "", "Callback URL for export results")
	exportOrgLogsCmd.Flags().StringVar(&callbackEmail, "callback-email", "", "Email for export results")

	// export-user-logs flags
	exportUserLogsCmd.Flags().StringVar(&starttime, "start-time", "", "Start time for export")
	exportUserLogsCmd.Flags().StringVar(&endtime, "end-time", "", "End time for export")
	exportUserLogsCmd.Flags().StringVar(&callbackURL, "callback-url", "", "Callback URL for export results")
	exportUserLogsCmd.Flags().StringVar(&callbackEmail, "callback-email", "", "Email for export results")

	// export-repo-logs flags
	exportRepoLogsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Repository namespace")
	exportRepoLogsCmd.Flags().StringVarP(&repository, "repository", "r", "", "Repository name")
	exportRepoLogsCmd.Flags().StringVar(&starttime, "start-time", "", "Start time for export")
	exportRepoLogsCmd.Flags().StringVar(&endtime, "end-time", "", "End time for export")
	exportRepoLogsCmd.Flags().StringVar(&callbackURL, "callback-url", "", "Callback URL for export results")
	exportRepoLogsCmd.Flags().StringVar(&callbackEmail, "callback-email", "", "Email for export results")
}
