package cmd

import (
	"fmt"
	"os"

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
}

// SecScan Info
var secscanInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get security scan results for a manifest",
	Long: `Get security scan results for a specific manifest including vulnerability information.

The scan status can be:
  - scanned: Scan completed successfully
  - queued: Scan is queued and pending
  - scanning: Scan is currently in progress
  - unsupported: Image type is not supported for scanning
  - failed: Scan failed`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		security, err := client.GetManifestSecurity(namespace, repository, secScanManifestRef, includeVulnerabilities)
		if err != nil {
			fmt.Printf("Error getting security scan: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Security scan for %s/%s@%s\n", namespace, repository, secScanManifestRef)
		printJSON(security)
	},
}

func init() {
	// Add subcommands to secscan command
	secscanCmd.AddCommand(secscanInfoCmd)

	// Global secscan flags (repository context)
	secscanCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Name of the namespace")
	secscanCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "Name of the repository")
	secscanCmd.PersistentFlags().StringVarP(&secScanManifestRef, "manifest", "m", "", "Manifest reference (digest like sha256:...)")
	secscanCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Bearer token")

	// Mark global flags as required
	if err := secscanCmd.MarkPersistentFlagRequired("namespace"); err != nil {
		fmt.Printf("Error marking namespace flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := secscanCmd.MarkPersistentFlagRequired("repository"); err != nil {
		fmt.Printf("Error marking repository flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := secscanCmd.MarkPersistentFlagRequired("manifest"); err != nil {
		fmt.Printf("Error marking manifest flag as required: %v\n", err)
		os.Exit(1)
	}
	if err := secscanCmd.MarkPersistentFlagRequired("token"); err != nil {
		fmt.Printf("Error marking token flag as required: %v\n", err)
		os.Exit(1)
	}

	// Info command specific flags
	secscanInfoCmd.Flags().BoolVarP(&includeVulnerabilities, "vulnerabilities", "V", true, "Include vulnerability details in the response")
}
