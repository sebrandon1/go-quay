package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	secScanManifestRef     string
	includeVulnerabilities bool
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
  - failed: Scan failed`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		security, err := client.GetManifestSecurity(cmd.Context(), namespace, repository, secScanManifestRef, includeVulnerabilities)
		if err != nil {
			return fmt.Errorf("getting security scan: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Security scan for %s/%s@%s\n", namespace, repository, secScanManifestRef)
		if outputFormat == outputTable {
			return printSecuritySummary(security)
		}
		return printJSON(security)
	},
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
}
