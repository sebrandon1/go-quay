package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
)

func markFlagRequired(err error) {
	if err != nil {
		fmt.Printf("Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}

// mustGetClient creates a Quay client with the configured token and URL.
// Exits with error message if client creation fails.
func mustGetClient() *lib.Client {
	client, err := lib.NewClientWithURL(token, quayURL)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}
	return client
}

// printJSON marshals and prints data as formatted JSON
func printJSON(data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
