package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/sebrandon1/go-toolbox/cliout"
	"github.com/sebrandon1/go-toolbox/yamlx"
)

// outputFormat holds the selected output format (json, yaml, or table).
// Set via the --output/-O persistent flag on rootCmd.
var outputFormat string

const (
	tableHeaderName      = "NAME"
	tableSeverityUnknown = "UNKNOWN"
)

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
	if verbose {
		client.RequestLogger = func(entry lib.RequestLogEntry) {
			fmt.Fprintf(rootCmd.ErrOrStderr(), "HTTP %s %s status=%d duration=%s attempt=%d\n",
				entry.Method, entry.URL, entry.StatusCode, entry.Duration.Round(time.Millisecond), entry.Attempt)
		}
	}
	if maxRetries > 0 && !dryRun {
		client.Retry = &lib.RetryConfig{
			MaxRetries:     maxRetries,
			InitialBackoff: retryBackoff,
			MaxBackoff:     retryMaxBackoff,
		}
	}
	return client, nil
}

// printJSON marshals and prints data in the format selected by --output.
// Supported formats: json (default), yaml, table (falls back to json).
func printJSON(data interface{}) error {
	switch outputFormat {
	case outputYAML:
		return yamlx.Print(os.Stdout, data)
	case outputTable:
		// Table output requires command-specific formatting.
		// Commands that support table mode check outputFormat (or --table)
		// and render their own table before calling printJSON.
		// For commands without custom table support, fall back to JSON — but say
		// so, since a caller piping into awk/cut wants columns and would otherwise
		// silently parse a different shape.
		fmt.Fprintf(rootCmd.ErrOrStderr(), "note: table output is not available for %s; printing json\n", currentCommandPath())
		return printAsJSON(data)
	default:
		return printAsJSON(data)
	}
}

// printAsJSON marshals and prints data as indented JSON.
func printAsJSON(data interface{}) error {
	return cliout.PrintJSON(os.Stdout, data)
}

// printTable writes a tab-separated table with aligned columns to stdout.
func printTable(headers []string, rows [][]string) error {
	return cliout.PrintTable(os.Stdout, headers, rows)
}
