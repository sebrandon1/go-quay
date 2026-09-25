package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/sebrandon1/go-quay/lib"
	"gopkg.in/yaml.v3"
)

// outputFormat holds the selected output format (json, yaml, or table).
// Set via the --output/-O persistent flag on rootCmd.
var outputFormat string

const (
	tableHeaderName      = "NAME"
	tableSeverityUnknown = "UNKNOWN"
)

var tableCellReplacer = strings.NewReplacer("\t", " ", "\n", " ", "\r", " ")

// getClient creates a Quay client with the configured token and URL.
func getClient() (*lib.Client, error) {
	client, err := lib.NewClientWithURL(token, quayURL)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("unexpected nil client")
	}
	client.Version = rootCmd.Version
	if dryRun {
		client.HTTPClient.Transport = dryRunTransport{}
		client.Retry = nil
	}
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

// printTable writes a tab-separated table with aligned columns to stdout.
func printTable(headers []string, rows [][]string) error {
	for i, row := range rows {
		if len(row) != len(headers) {
			return fmt.Errorf("table row %d has %d columns, want %d", i+1, len(row), len(headers))
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, joinTableCells(headers)); err != nil {
		return fmt.Errorf("writing table header: %w", err)
	}
	for _, row := range rows {
		if _, err := fmt.Fprintln(w, joinTableCells(row)); err != nil {
			return fmt.Errorf("writing table row: %w", err)
		}
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("flushing table output: %w", err)
	}
	return nil
}

func joinTableCells(cells []string) string {
	clean := make([]string, len(cells))
	for i, cell := range cells {
		clean[i] = tableCellReplacer.Replace(cell)
	}
	return strings.Join(clean, "\t")
}
