package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

func init() {
	repositoryCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "Name of the namespace", "")
	repositoryCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "Name of the repository", "")
	repositoryCmd.PersistentFlags().StringVarP(&token, "token", "t", "Bearer token", "")
}

var repositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Get repository information from Quay.io",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			panic(err)
		}

		repo, err := client.GetRepository(namespace, repository)
		if err != nil {
			panic(err)
		}

		jsonOutput, err := json.Marshal(repo)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(jsonOutput))
	},
}
