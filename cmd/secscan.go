package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	secScanManifestRef     string
	includeVulnerabilities bool
	secScanWatch           bool
	secScanInterval        time.Duration
	secScanWatchTimeout    time.Duration
)

const (
	defaultSecScanInterval     = 3 * time.Second
	defaultSecScanWatchTimeout = 5 * time.Minute
)

// secscanCmd represents the secscan command group
var secscanCmd = &cobra.Command{
	Use:   "secscan",
	Short: "Security scanning commands for container images",
	Long: `Commands for retrieving security scan information for container images.

Security scanning provides vulnerability information including CVE details,
severity levels, affected packages, and available fixes.

Available commands:
  info - Get security scan results for a manifest`,
	Example: `  go-quay get secscan info --namespace myorg --repository myapp --manifest sha256:abc123 --token "$QUAY_TOKEN"`,
}

// SecScan Info
var secscanInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get security scan results for a manifest",
	Long: `Get security scan results for a specific manifest including vulnerability information.

Table columns: STATUS, FEATURES, VULNERABILITIES, CRITICAL, HIGH, MEDIUM, LOW, NEGLIGIBLE, UNKNOWN.

The scan status can be:
  - scanned: Scan completed successfully
  - queued: Scan is queued and pending
  - scanning: Scan is currently in progress
  - unsupported: Image type is not supported for scanning
  - failed: Scan failed

Use --watch to poll until the scan reaches a terminal state.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		var security *lib.SecurityScan
		if secScanWatch {
			security, err = waitForSecurityScan(cmd, client)
		} else {
			security, err = client.GetManifestSecurity(cmd.Context(), namespace, repository, secScanManifestRef, includeVulnerabilities)
		}
		if err != nil {
			return fmt.Errorf("getting security scan: %w", err)
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "Security scan for %s/%s@%s\n", namespace, repository, secScanManifestRef)
		if outputFormat == outputTable {
			err = printSecuritySummary(security)
		} else {
			err = printJSON(security)
		}
		if err != nil {
			return err
		}
		if secScanWatch && strings.EqualFold(strings.TrimSpace(security.Status), "failed") {
			return fmt.Errorf("security scan failed with status %q", security.Status)
		}
		return nil
	},
}

func waitForSecurityScan(cmd *cobra.Command, client *lib.Client) (*lib.SecurityScan, error) {
	if secScanInterval <= 0 {
		return nil, fmt.Errorf("--interval must be greater than 0")
	}
	if secScanWatchTimeout <= 0 {
		return nil, fmt.Errorf("--watch-timeout must be greater than 0")
	}

	parentCtx := cmd.Context()
	ctx, cancel := context.WithTimeout(parentCtx, secScanWatchTimeout)
	defer cancel()

	for poll := 1; ; poll++ {
		scan, err := client.GetManifestSecurity(ctx, namespace, repository, secScanManifestRef, includeVulnerabilities)
		if err != nil {
			if parentCtx.Err() != nil {
				return nil, parentCtx.Err()
			}
			if ctx.Err() != nil {
				return nil, fmt.Errorf("security scan watch timed out after %s: %w", secScanWatchTimeout, ctx.Err())
			}
			return nil, err
		}

		status := ""
		if scan != nil {
			status = strings.ToLower(strings.TrimSpace(scan.Status))
		}
		switch status {
		case "scanned", "failed", "unsupported":
			return scan, nil
		}

		fmt.Fprintf(cmd.ErrOrStderr(), "Security scan status %q (poll %d); waiting\n", status, poll)
		timer := time.NewTimer(secScanInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			if parentCtx.Err() != nil {
				return nil, parentCtx.Err()
			}
			return nil, fmt.Errorf("security scan watch timed out after %s: %w", secScanWatchTimeout, ctx.Err())
		case <-timer.C:
		}
	}
}

func printSecuritySummary(scan *lib.SecurityScan) error {
	counts := map[string]int{
		"CRITICAL":           0,
		"HIGH":               0,
		"MEDIUM":             0,
		"LOW":                0,
		"NEGLIGIBLE":         0,
		tableSeverityUnknown: 0,
	}
	features := 0
	if scan != nil && scan.Data != nil && scan.Data.Layer != nil {
		features = len(scan.Data.Layer.Features)
		for _, feature := range scan.Data.Layer.Features {
			for _, vulnerability := range feature.Vulnerabilities {
				severity := strings.ToUpper(strings.TrimSpace(vulnerability.Severity))
				if _, found := counts[severity]; !found {
					severity = tableSeverityUnknown
				}
				counts[severity]++
			}
		}
	}
	status := ""
	if scan != nil {
		status = scan.Status
	}
	total := 0
	for _, count := range counts {
		total += count
	}
	rows := [][]string{{
		status,
		fmt.Sprintf("%d", features),
		fmt.Sprintf("%d", total),
		fmt.Sprintf("%d", counts["CRITICAL"]),
		fmt.Sprintf("%d", counts["HIGH"]),
		fmt.Sprintf("%d", counts["MEDIUM"]),
		fmt.Sprintf("%d", counts["LOW"]),
		fmt.Sprintf("%d", counts["NEGLIGIBLE"]),
		fmt.Sprintf("%d", counts[tableSeverityUnknown]),
	}}
	return printTable([]string{
		"STATUS", "FEATURES", "VULNERABILITIES", "CRITICAL", "HIGH", "MEDIUM", "LOW", "NEGLIGIBLE", tableSeverityUnknown,
	}, rows)
}

func init() {
	// Add subcommands to secscan command
	secscanCmd.AddCommand(secscanInfoCmd)

	// Global secscan flags (repository context)
	secscanCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", appCfg.Namespace, "Name of the namespace (default: config file)")
	secscanCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	secscanCmd.PersistentFlags().StringVarP(&secScanManifestRef, "manifest", "m", "", "Manifest reference (digest like sha256:...)")

	// Mark namespace as required only when no config default is available
	if appCfg.Namespace == "" {
		_ = secscanCmd.MarkPersistentFlagRequired("namespace")
	}
	_ = secscanCmd.MarkPersistentFlagRequired("repository")
	_ = secscanCmd.MarkPersistentFlagRequired("manifest")

	// Info command specific flags
	secscanInfoCmd.Flags().BoolVarP(&includeVulnerabilities, "vulnerabilities", "V", true, "Include vulnerability details in the response")
	secscanInfoCmd.Flags().BoolVar(&secScanWatch, "watch", false, "Poll until the security scan reaches a terminal state")
	secscanInfoCmd.Flags().DurationVar(&secScanInterval, "interval", defaultSecScanInterval, "Polling interval when --watch is enabled")
	secscanInfoCmd.Flags().DurationVar(&secScanWatchTimeout, "watch-timeout", defaultSecScanWatchTimeout, "Maximum time to wait when --watch is enabled")
}
