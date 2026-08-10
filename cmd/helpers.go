package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-quay/lib"
)

// getClient creates a Quay client with the configured token and URL.
func getClient() (*lib.Client, error) {
	client, err := lib.NewClientWithURL(token, quayURL)
	if err != nil {
		return nil, err
	}
	client.Version = rootCmd.Version
	return client, nil
}

// printJSON marshals and prints data as formatted JSON
func printJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	fmt.Println(string(output))
	return nil
}
