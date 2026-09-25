package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	quayURL         string
	dryRun          bool
	verbose         bool
	maxRetries      int
	retryBackoff    = 500 * time.Millisecond
	retryMaxBackoff = 5 * time.Second
)

var rootCmd = &cobra.Command{
	Use:   cliName,
	Short: "Quay CLI interacts with Quay.io API",
}

func SetVersion(v string) {
	rootCmd.Version = v
}

var getCmd = &cobra.Command{
	Use:   cmdGet,
	Short: "Get objects from Quay.io (legacy parent; prefer create/delete/update/list/info)",
}

func persistentPreRunE(cmd *cobra.Command, _ []string) error {
	token = resolveFlag(flagChanged(cmd, "token"), token, os.Getenv("QUAY_TOKEN"), appCfg.Token)
	quayURL = resolveFlag(flagChanged(cmd, "quay-url"), quayURL, os.Getenv("QUAY_URL"), appCfg.QuayURL, lib.DefaultQuayURL)
	if maxRetries < 0 {
		return fmt.Errorf("--max-retries must be zero or greater")
	}
	if retryBackoff < 0 {
		return fmt.Errorf("--retry-backoff must be zero or greater")
	}
	if retryMaxBackoff < 0 {
		return fmt.Errorf("--retry-max-backoff must be zero or greater")
	}

	if token == "" && !dryRun && !isAuthenticationExemptCommand(cmd) {
		return fmt.Errorf(`authentication token required

Set QUAY_TOKEN environment variable, use --token/-t flag, or add to config file (%s).
Get your token at https://quay.io/organization/<org>?tab=applications`, configFilePath())
	}

	switch outputFormat {
	case outputJSON, outputYAML, outputTable:
		// valid
	default:
		return fmt.Errorf("invalid output format %q: must be json, yaml, or table", outputFormat)
	}

	return nil
}

func isAuthenticationExemptCommand(cmd *cobra.Command) bool {
	for current := cmd; current != nil; current = current.Parent() {
		if current == configCmd || current == completionCmd {
			return true
		}
	}
	return false
}

func flagChanged(cmd *cobra.Command, name string) bool {
	f := cmd.Flag(name)
	return f != nil && f.Changed
}

// resolveFlag returns flagValue when the flag was explicitly set; otherwise the first non-empty fallback.
func resolveFlag(changed bool, flagValue string, fallbacks ...string) string {
	if changed {
		return flagValue
	}
	return firstNonEmpty(fallbacks...)
}

// firstNonEmpty returns the first non-empty string from the given values.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func init() {
	rootCmd.PersistentPreRunE = persistentPreRunE
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Quay.io API token ($QUAY_TOKEN or config file)")
	rootCmd.PersistentFlags().StringVar(&quayURL, "quay-url", lib.DefaultQuayURL, "Quay API base URL ($QUAY_URL or config file)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show API requests without sending them")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "O", "json", "Output format: json, yaml, or table")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Log HTTP requests and responses to stderr")
	rootCmd.PersistentFlags().IntVar(&maxRetries, "max-retries", 0, "Maximum retries after the initial request (0 disables retries)")
	rootCmd.PersistentFlags().DurationVar(&retryBackoff, "retry-backoff", 500*time.Millisecond, "Initial delay between retry attempts")
	rootCmd.PersistentFlags().DurationVar(&retryMaxBackoff, "retry-max-backoff", 5*time.Second, "Maximum delay between retry attempts")
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(diffCmd)
	getCmd.AddCommand(repositoryCmd)
	getCmd.AddCommand(billingCmd)
	getCmd.AddCommand(organizationCmd)
	getCmd.AddCommand(permissionsCmd)
	getCmd.AddCommand(tagCmd)
	getCmd.AddCommand(userCmd)
	getCmd.AddCommand(manifestCmd)
	getCmd.AddCommand(secscanCmd)
	getCmd.AddCommand(robotCmd)
	getCmd.AddCommand(searchCmd)
	getCmd.AddCommand(teamCmd)
	getCmd.AddCommand(buildCmd)
	getCmd.AddCommand(notificationCmd)
	getCmd.AddCommand(triggerCmd)
	getCmd.AddCommand(discoveryCmd)
	getCmd.AddCommand(errorTypeCmd)
	getCmd.AddCommand(messagesCmd)
	getCmd.AddCommand(prototypeCmd)
	getCmd.AddCommand(repotokenCmd)
	getCmd.AddCommand(logsCmd)
	getCmd.AddCommand(mirrorCmd)
}

// Execute executes the root command, canceling in-flight HTTP on interrupt.
func Execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return rootCmd.ExecuteContext(ctx)
}
