package cmd

import (
	"fmt"

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
}

// SecScan Info
var secscanInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get security scan results for a manifest",
	Long: `Get security scan results for a specific manifest including vulnerability information.

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

		security, err := client.GetManifestSecurity(namespace, repository, secScanManifestRef, includeVulnerabilities)
		if err != nil {
			return fmt.Errorf("getting security scan: %w", err)
		}

		fmt.Printf("Security scan for %s/%s@%s\n", namespace, repository, secScanManifestRef)
		return printJSON(security)
	},
}

func init() {
	// Add subcommands to secscan command
	secscanCmd.AddCommand(secscanInfoCmd)

	// Global secscan flags (repository context)
	secscanCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	secscanCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	secscanCmd.PersistentFlags().StringVarP(&secScanManifestRef, "manifest", "m", "", "Manifest reference (digest like sha256:...)")

	// Mark global flags as required
	_ = secscanCmd.MarkPersistentFlagRequired("namespace")
	_ = secscanCmd.MarkPersistentFlagRequired("repository")
	_ = secscanCmd.MarkPersistentFlagRequired("manifest")

	// Info command specific flags
	secscanInfoCmd.Flags().BoolVarP(&includeVulnerabilities, "vulnerabilities", "V", true, "Include vulnerability details in the response")
}
