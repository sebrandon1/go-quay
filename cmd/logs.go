package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/openshift/cluster-node-tuning-operator/test/e2e/performanceprofile/functests/utils/log"
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

func init() {
	aggregatedLogsCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "Name of the namespace", "")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "Name of the repository", "")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&startdate, "startdate", "s", "Start date for the logs", "")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&enddate, "enddate", "e", "End date for the logs", "")
	aggregatedLogsCmd.PersistentFlags().StringVarP(&token, "token", "t", "Bearer token", "")
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
			log.Infof("Error getting aggregated logs: %v", err)
			panic(err)
		}

		jsonOutput, err := json.Marshal(logs)
		if err != nil {
			log.Infof("Error marshaling logs to JSON: %v", err)
			panic(err)
		}

		fmt.Println(string(jsonOutput))
	},
}
