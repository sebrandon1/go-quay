package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	namespace  string
	repository string
	startdate  string
	enddate    string
	token      string
)

var rootCmd = &cobra.Command{
	Use:   "quay",
	Short: "Quay CLI interacts with Quay.io API",
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get objects from Quay.io",
}

var aggregatedLogsCmd = &cobra.Command{
	Use:   "aggregatedlogs",
	Short: "Get aggregated logs from Quay.io",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			panic(err)
		}

		logs, err := client.GetAggregatedLogs(namespace, repository, startdate, enddate)
		if err != nil {
			panic(err)
		}

		jsonOutput, err := json.Marshal(logs)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(jsonOutput))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(aggregatedLogsCmd)

	aggregatedLogsCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "Name of the namespace", "")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "Name of the repository", "")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&startdate, "startdate", "s", "Start date for the logs", "")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&enddate, "enddate", "e", "End date for the logs", "")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&token, "token", "t", "Bearer token", "")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
