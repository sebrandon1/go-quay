package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

func markFlagRequired(err error) {
	if err != nil {
		fmt.Printf("Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
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
