package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sebrandon1/go-quay/lib"
	"gopkg.in/yaml.v3"
)

// outputFormat holds the selected output format (json, yaml, or table).
// Set via the --output/-O persistent flag on getCmd.
var outputFormat string

// getClient creates a Quay client with the configured token and URL.
func getClient() (*lib.Client, error) {
	client, err := lib.NewClientWithURL(token, quayURL)
	if err != nil {
		return nil, err
	}
	client.Version = rootCmd.Version
	return client, nil
}

// printJSON marshals and prints data in the format selected by --output.
// Supported formats: json (default), yaml, table (falls back to json).
func printJSON(data interface{}) error {
	switch outputFormat {
	case outputYAML:
		output, err := yaml.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshaling YAML: %w", err)
		}
		fmt.Print(string(output))
		return nil
	case outputTable:
		// Table output requires command-specific formatting.
		// Commands that support table mode check outputFormat (or --table)
		// and render their own table before calling printJSON.
		// For commands without custom table support, fall back to JSON.
		return printAsJSON(data)
	default:
		return printAsJSON(data)
	}
}

// printAsJSON marshals and prints data as indented JSON.
func printAsJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	fmt.Println(string(output))
	return nil
}
