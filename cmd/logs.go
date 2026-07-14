package cmd

import (
	"fmt"

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

// logsCmd is the parent command group for log management
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Log management commands",
}

var repoLogsCmd = &cobra.Command{
	Use:   "repo-logs",
	Short: "Get repository logs",
	Long:  `Get action logs for a specific repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		logs, err := client.GetLogs(namespace, repository, nextPage, startdate, enddate)
		if err != nil {
			return fmt.Errorf("getting repository logs: %w", err)
		}

		return printJSON(logs)
	},
}

var repoAggregatedLogsCmd = &cobra.Command{
	Use:   "repo-aggregated-logs",
	Short: "Get aggregated repository logs",
	Long:  `Get aggregated action logs for a specific repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		logs, err := client.GetAggregatedLogs(namespace, repository, startdate, enddate)
		if err != nil {
			return fmt.Errorf("getting repository aggregated logs: %w", err)
		}

		return printJSON(logs)
	},
}

var orgLogsCmd = &cobra.Command{
	Use:   "org-logs",
	Short: "Get organization logs",
	Long:  `Get action logs for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		logs, err := client.GetOrganizationLogs(orgName, nextPage, startdate, enddate)
		if err != nil {
			return fmt.Errorf("getting organization logs: %w", err)
		}

		return printJSON(logs)
	},
}

var orgAggregatedLogsCmd = &cobra.Command{
	Use:   "org-aggregated-logs",
	Short: "Get aggregated organization logs",
	Long:  `Get aggregated action logs for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		logs, err := client.GetOrganizationAggregatedLogs(orgName, startdate, enddate)
		if err != nil {
			return fmt.Errorf("getting organization aggregated logs: %w", err)
		}

		return printJSON(logs)
	},
}

var userLogsCmd = &cobra.Command{
	Use:   "user-logs",
	Short: "Get user logs",
	Long:  `Get action logs for the current user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		logs, err := client.GetUserLogs(nextPage, startdate, enddate)
		if err != nil {
			return fmt.Errorf("getting user logs: %w", err)
		}

		return printJSON(logs)
	},
}

var userAggregatedLogsCmd = &cobra.Command{
	Use:   "user-aggregated-logs",
	Short: "Get aggregated user logs",
	Long:  `Get aggregated action logs for the current user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		logs, err := client.GetUserAggregatedLogs(startdate, enddate)
		if err != nil {
			return fmt.Errorf("getting user aggregated logs: %w", err)
		}

		return printJSON(logs)
	},
}

var exportOrgLogsCmd = &cobra.Command{
	Use:   "export-org-logs",
	Short: "Export organization logs",
	Long:  `Export action logs for an organization via callback URL or email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.ExportOrganizationLogs(orgName, &lib.ExportLogsRequest{
			StartTime:   starttime,
			EndTime:     endtime,
			CallbackURL: callbackURL,
			Email:       callbackEmail,
		})
		if err != nil {
			return fmt.Errorf("exporting organization logs: %w", err)
		}

		fmt.Println("Organization logs export initiated successfully.")
		return nil
	},
}

var exportUserLogsCmd = &cobra.Command{
	Use:   "export-user-logs",
	Short: "Export user logs",
	Long:  `Export action logs for the current user via callback URL or email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.ExportUserLogs(&lib.ExportLogsRequest{
			StartTime:   starttime,
			EndTime:     endtime,
			CallbackURL: callbackURL,
			Email:       callbackEmail,
		})
		if err != nil {
			return fmt.Errorf("exporting user logs: %w", err)
		}

		fmt.Println("User logs export initiated successfully.")
		return nil
	},
}

var exportRepoLogsCmd = &cobra.Command{
	Use:   "export-repo-logs",
	Short: "Export repository logs",
	Long:  `Export action logs for a repository via callback URL or email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.ExportRepositoryLogs(namespace, repository, &lib.ExportLogsRequest{
			StartTime:   starttime,
			EndTime:     endtime,
			CallbackURL: callbackURL,
			Email:       callbackEmail,
		})
		if err != nil {
			return fmt.Errorf("exporting repository logs: %w", err)
		}

		fmt.Println("Repository logs export initiated successfully.")
		return nil
	},
}

func init() {
	logsCmd.AddCommand(repoLogsCmd)
	logsCmd.AddCommand(repoAggregatedLogsCmd)
	logsCmd.AddCommand(orgLogsCmd)
	logsCmd.AddCommand(orgAggregatedLogsCmd)
	logsCmd.AddCommand(userLogsCmd)
	logsCmd.AddCommand(userAggregatedLogsCmd)
	logsCmd.AddCommand(exportOrgLogsCmd)
	logsCmd.AddCommand(exportUserLogsCmd)
	logsCmd.AddCommand(exportRepoLogsCmd)

	// repo-logs flags
	repoLogsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Repository namespace")
	repoLogsCmd.Flags().StringVarP(&repository, "repository", "r", "", "Repository name")
	repoLogsCmd.Flags().StringVar(&nextPage, "next-page", "", "Next page token for pagination")
	repoLogsCmd.Flags().StringVarP(&startdate, "startdate", "s", "", "Start date for the logs")
	repoLogsCmd.Flags().StringVarP(&enddate, "enddate", "e", "", "End date for the logs")

	// repo-aggregated-logs flags
	repoAggregatedLogsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Repository namespace")
	repoAggregatedLogsCmd.Flags().StringVarP(&repository, "repository", "r", "", "Repository name")
	repoAggregatedLogsCmd.Flags().StringVarP(&startdate, "startdate", "s", "", "Start date")
	repoAggregatedLogsCmd.Flags().StringVarP(&enddate, "enddate", "e", "", "End date")

	// org-logs flags
	orgLogsCmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	orgLogsCmd.Flags().StringVar(&nextPage, "next-page", "", "Next page token for pagination")
	orgLogsCmd.Flags().StringVarP(&startdate, "startdate", "s", "", "Start date for the logs")
	orgLogsCmd.Flags().StringVarP(&enddate, "enddate", "e", "", "End date for the logs")

	// org-aggregated-logs flags
	orgAggregatedLogsCmd.Flags().StringVarP(&orgName, "organization", "o", "", "Organization name")
	orgAggregatedLogsCmd.Flags().StringVarP(&startdate, "startdate", "s", "", "Start date")
	orgAggregatedLogsCmd.Flags().StringVarP(&enddate, "enddate", "e", "", "End date")

	// user-logs flags
	userLogsCmd.Flags().StringVar(&nextPage, "next-page", "", "Next page token for pagination")
	userLogsCmd.Flags().StringVarP(&startdate, "startdate", "s", "", "Start date for the logs")
	userLogsCmd.Flags().StringVarP(&enddate, "enddate", "e", "", "End date for the logs")

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
