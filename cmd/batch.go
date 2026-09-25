package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const batchAPIVersion = "go-quay/v1alpha1"

const (
	batchStatusSucceeded = "succeeded"
	batchStatusFailed    = "failed"
	batchStatusNotRun    = "not_run"
	batchStatusDryRun    = "dry_run"
)

var (
	batchFile     string
	batchFailFast bool
	batchConfirm  bool
	batchDryRun   bool
)

type repositoryBatch struct {
	APIVersion string                `yaml:"apiVersion"`
	Items      []repositoryBatchItem `yaml:"items"`
}

type repositoryBatchItem struct {
	Action      string `yaml:"action"`
	Namespace   string `yaml:"namespace"`
	Repository  string `yaml:"repository"`
	Visibility  string `yaml:"visibility,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type repositoryBatchResult struct {
	Action     string `json:"action"`
	Namespace  string `json:"namespace"`
	Repository string `json:"repository"`
	Status     string `json:"status"`
	Error      string `json:"error,omitempty"`
}

var batchCmd = &cobra.Command{
	Use:   cmdBatch,
	Short: "Run repository operations from a YAML or JSON batch file",
}

var batchApplyCmd = &cobra.Command{
	Use:   subcmdApply,
	Short: "Create or delete repositories from a batch file or stdin",
	Long: `Apply repository create and delete operations from a YAML or JSON document.

Each item specifies an action, namespace, and repository. The entire document is
validated before any API requests are sent. Operations run sequentially; failures
are reported per item and the command exits with an error if any item fails.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		batch, err := readRepositoryBatch(cmd, batchFile)
		if err != nil {
			return err
		}
		if err := validateRepositoryBatch(batch); err != nil {
			return err
		}

		containsDelete := false
		for _, item := range batch.Items {
			if item.Action == subcmdDelete {
				containsDelete = true
				break
			}
		}
		if containsDelete && !batchConfirm && !batchDryRun {
			return fmt.Errorf("batch contains delete operations; pass --confirm to proceed")
		}

		if batchDryRun {
			results := make([]repositoryBatchResult, 0, len(batch.Items))
			for _, item := range batch.Items {
				results = append(results, repositoryBatchResult{
					Action: item.Action, Namespace: item.Namespace, Repository: item.Repository, Status: batchStatusDryRun,
				})
			}
			return printJSON(results)
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		results := make([]repositoryBatchResult, 0, len(batch.Items))
		failed := 0
		skipped := 0
		for i, item := range batch.Items {
			result := repositoryBatchResult{
				Action: item.Action, Namespace: item.Namespace, Repository: item.Repository, Status: batchStatusSucceeded,
			}
			if err := applyRepositoryBatchItem(cmd, client, item); err != nil {
				result.Status = batchStatusFailed
				result.Error = err.Error()
				failed++
			}
			results = append(results, result)

			if batchFailFast && result.Status == batchStatusFailed {
				for _, remaining := range batch.Items[i+1:] {
					results = append(results, repositoryBatchResult{
						Action: remaining.Action, Namespace: remaining.Namespace, Repository: remaining.Repository,
						Status: batchStatusNotRun, Error: "stopped by --fail-fast",
					})
					skipped++
				}
				break
			}
		}

		if err := printJSON(results); err != nil {
			return err
		}
		if failed > 0 {
			if skipped > 0 {
				return fmt.Errorf("batch completed with %d failed operation(s) and %d not run", failed, skipped)
			}
			return fmt.Errorf("batch completed with %d failed operation(s)", failed)
		}
		return nil
	},
}

func readRepositoryBatch(cmd *cobra.Command, path string) (*repositoryBatch, error) {
	var input io.Reader
	if path == "-" {
		input = cmd.InOrStdin()
	} else {
		file, err := os.Open(path) // #nosec G304 -- path comes from the explicit user-provided --file flag
		if err != nil {
			return nil, fmt.Errorf("opening batch file %q: %w", path, err)
		}
		defer file.Close()
		input = file
	}

	decoder := yaml.NewDecoder(input)
	decoder.KnownFields(true)
	var batch repositoryBatch
	if err := decoder.Decode(&batch); err != nil {
		return nil, fmt.Errorf("decoding batch input: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("batch input must contain one YAML or JSON document")
		}
		return nil, fmt.Errorf("decoding batch input: %w", err)
	}
	return &batch, nil
}

func validateRepositoryBatch(batch *repositoryBatch) error {
	if batch.APIVersion != batchAPIVersion {
		return fmt.Errorf("apiVersion must be %q", batchAPIVersion)
	}
	if len(batch.Items) == 0 {
		return fmt.Errorf("batch items must not be empty")
	}

	for i := range batch.Items {
		item := &batch.Items[i]
		item.Action = strings.ToLower(strings.TrimSpace(item.Action))
		item.Namespace = strings.TrimSpace(item.Namespace)
		item.Repository = strings.TrimSpace(item.Repository)
		item.Visibility = strings.ToLower(strings.TrimSpace(item.Visibility))
		if item.Namespace == "" || item.Repository == "" {
			return fmt.Errorf("items[%d] requires namespace and repository", i)
		}
		switch item.Action {
		case subcmdCreate:
			if item.Visibility == "" {
				item.Visibility = repoVisibilityPrivate
			}
			if item.Visibility != "public" && item.Visibility != repoVisibilityPrivate {
				return fmt.Errorf("items[%d] visibility must be public or private", i)
			}
		case subcmdDelete:
			if item.Visibility != "" || item.Description != "" {
				return fmt.Errorf("items[%d] delete operations do not accept visibility or description", i)
			}
		default:
			return fmt.Errorf("items[%d] has unsupported action %q; supported actions are %s and %s", i, item.Action, subcmdCreate, subcmdDelete)
		}
	}
	return nil
}

func applyRepositoryBatchItem(cmd *cobra.Command, client *lib.Client, item repositoryBatchItem) error {
	switch item.Action {
	case subcmdCreate:
		_, err := client.CreateRepository(cmd.Context(), item.Namespace, item.Repository, item.Visibility, item.Description)
		if err != nil {
			return fmt.Errorf("creating repository %s/%s: %w", item.Namespace, item.Repository, err)
		}
	case subcmdDelete:
		if err := client.DeleteRepository(cmd.Context(), item.Namespace, item.Repository); err != nil {
			return fmt.Errorf("deleting repository %s/%s: %w", item.Namespace, item.Repository, err)
		}
	}
	return nil
}

func init() {
	batchCmd.AddCommand(batchApplyCmd)
	batchApplyCmd.Flags().StringVarP(&batchFile, "file", "f", "", "YAML or JSON input file, or - to read stdin")
	_ = batchApplyCmd.MarkFlagRequired("file")
	batchApplyCmd.Flags().BoolVar(&batchFailFast, "fail-fast", false, "Stop after the first failed operation")
	batchApplyCmd.Flags().BoolVar(&batchConfirm, "confirm", false, "Confirm batch delete operations")
	batchApplyCmd.Flags().BoolVar(&batchDryRun, "dry-run", false, "Validate and preview operations without sending API requests")
}
