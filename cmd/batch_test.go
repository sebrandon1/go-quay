package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	testBatchFalse      = "false"
	testBatchRepository = "repo"
)

func resetBatchCommandState(t *testing.T) {
	t.Helper()
	resetRootFlags(t)
	batchFile = ""
	batchFailFast = false
	batchConfirm = false
	batchDryRun = false
	for name, value := range map[string]string{
		"file": "", "fail-fast": testBatchFalse, "confirm": testBatchFalse, "dry-run": testBatchFalse,
	} {
		if flag := batchApplyCmd.Flags().Lookup(name); flag != nil {
			flag.Changed = false
			if err := flag.Value.Set(value); err != nil {
				t.Fatalf("reset %s flag: %v", name, err)
			}
		}
	}
	rootCmd.SetIn(nil)
	t.Cleanup(func() {
		rootCmd.SetIn(nil)
		rootCmd.SetArgs(nil)
	})
}

func runBatchCLI(t *testing.T, input, filePath, quayURL string, flags ...string) (string, error) {
	t.Helper()
	resetBatchCommandState(t)
	oldSilenceUsage := rootCmd.SilenceUsage
	rootCmd.SilenceUsage = true
	t.Cleanup(func() { rootCmd.SilenceUsage = oldSilenceUsage })
	if filePath == "-" {
		rootCmd.SetIn(strings.NewReader(input))
	}
	var stderr bytes.Buffer
	rootCmd.SetErr(&stderr)
	args := []string{cmdBatch, subcmdApply, "--file", filePath}
	if quayURL != "" {
		args = append(args, "--quay-url", quayURL)
	}
	args = append(args, flags...)
	rootCmd.SetArgs(args)
	var runErr error
	stdout := captureStdout(t, func() { runErr = rootCmd.Execute() })
	return stdout, runErr
}

func newBatchCreateDeleteServer(t *testing.T, requests *[]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*requests = append(*requests, r.Method+" "+r.URL.Path)
		if r.Method == http.MethodPost && r.URL.Path == "/api/v1/repository" {
			assertBatchCreateRequest(t, r)
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"name":"new-repo"}`)
			return
		}
		if r.Method == http.MethodDelete && r.URL.Path == "/api/v1/repository/"+testNamespace+"/old-repo" {
			w.WriteHeader(http.StatusOK)
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	}))
}

func assertBatchCreateRequest(t *testing.T, r *http.Request) {
	t.Helper()
	var body struct {
		Namespace   string `json:"namespace"`
		Repository  string `json:"repository"`
		Visibility  string `json:"visibility"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Errorf("decode create request: %v", err)
	}
	if body.Namespace != testNamespace || body.Repository != "new-repo" || body.Visibility != "public" || body.Description != "batch-created" {
		t.Errorf("create request body = %+v", body)
	}
}

func TestBatchApplyCreatesAndDeletesRepositoriesFromFile(t *testing.T) {
	input := `apiVersion: go-quay/v1alpha1
items:
  - action: create
    namespace: testns
    repository: new-repo
    visibility: public
    description: batch-created
  - action: delete
    namespace: testns
    repository: old-repo
`
	filePath := t.TempDir() + "/repos.yaml"
	if err := os.WriteFile(filePath, []byte(input), 0o600); err != nil {
		t.Fatalf("write batch file: %v", err)
	}

	var requests []string
	server := newBatchCreateDeleteServer(t, &requests)
	defer server.Close()

	stdout, err := runBatchCLI(t, "", filePath, server.URL+"/api/v1", "--token", testTokenValue, "--confirm")
	if err != nil {
		t.Fatalf("batch apply: %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("requests = %v, want create and delete", requests)
	}
	var results []repositoryBatchResult
	if err := json.Unmarshal([]byte(stdout), &results); err != nil {
		t.Fatalf("decode batch results: %v\n%s", err, stdout)
	}
	if len(results) != 2 || results[0].Status != batchStatusSucceeded || results[1].Status != batchStatusSucceeded {
		t.Fatalf("batch results = %+v, want two successes", results)
	}
}

func TestBatchApplyAcceptsJSONFromStdinAndDryRunsWithoutToken(t *testing.T) {
	t.Setenv("QUAY_TOKEN", "")
	input := `{"apiVersion":"go-quay/v1alpha1","items":[{"action":"create","namespace":"testns","repository":"json-repo"}]}`
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { requests++ }))
	defer server.Close()

	stdout, err := runBatchCLI(t, input, "-", server.URL+"/api/v1", "--dry-run")
	if err != nil {
		t.Fatalf("batch dry-run: %v", err)
	}
	if requests != 0 {
		t.Fatalf("dry-run sent %d requests, want none", requests)
	}
	var results []repositoryBatchResult
	if err := json.Unmarshal([]byte(stdout), &results); err != nil {
		t.Fatalf("decode dry-run results: %v\n%s", err, stdout)
	}
	if len(results) != 1 || results[0].Status != batchStatusDryRun {
		t.Fatalf("dry-run results = %+v, want one dry_run item", results)
	}
}

func TestBatchApplyValidatesWholeInputBeforeSendingRequests(t *testing.T) {
	input := `apiVersion: go-quay/v1alpha1
items:
  - action: create
    namespace: testns
    repository: valid-repo
  - action: update
    namespace: testns
    repository: unsupported-repo
`
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { requests++ }))
	defer server.Close()

	_, err := runBatchCLI(t, input, "-", server.URL+"/api/v1", "--token", testTokenValue)
	if err == nil || !strings.Contains(err.Error(), "supported actions are create and delete") {
		t.Fatalf("error = %v, want unsupported action validation error", err)
	}
	if requests != 0 {
		t.Fatalf("invalid batch sent %d requests, want none", requests)
	}
}

func TestBatchApplyContinuesAfterItemFailureAndReturnsAggregateError(t *testing.T) {
	input := `apiVersion: go-quay/v1alpha1
items:
  - action: create
    namespace: testns
    repository: denied-repo
  - action: create
    namespace: testns
    repository: later-repo
`
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		var body struct {
			Repository string `json:"repository"`
		}
		if decodeErr := json.NewDecoder(r.Body).Decode(&body); decodeErr != nil {
			t.Errorf("decode create request: %v", decodeErr)
		}
		if body.Repository == "denied-repo" {
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, `{"error":"denied"}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"name":"`+body.Repository+`"}`)
	}))
	defer server.Close()

	stdout, err := runBatchCLI(t, input, "-", server.URL+"/api/v1", "--token", testTokenValue)
	if err == nil || !strings.Contains(err.Error(), "1 failed operation") {
		t.Fatalf("error = %v, want aggregate failure", err)
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want continuation after first failure", requests)
	}
	var results []repositoryBatchResult
	if decodeErr := json.Unmarshal([]byte(stdout), &results); decodeErr != nil {
		t.Fatalf("decode results: %v\n%s", decodeErr, stdout)
	}
	if len(results) != 2 || results[0].Status != batchStatusFailed || results[1].Status != batchStatusSucceeded {
		t.Fatalf("batch results = %+v, want failed then succeeded", results)
	}
}

func TestBatchApplyReportsDeleteFailure(t *testing.T) {
	input := `apiVersion: go-quay/v1alpha1
items:
  - action: delete
    namespace: testns
    repository: denied-repo
`
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodDelete {
			t.Errorf("request method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"error":"denied"}`)
	}))
	defer server.Close()

	stdout, err := runBatchCLI(t, input, "-", server.URL+"/api/v1", "--token", testTokenValue, "--confirm")
	if err == nil || !strings.Contains(err.Error(), "1 failed operation") {
		t.Fatalf("error = %v, want aggregate failure", err)
	}
	var results []repositoryBatchResult
	if decodeErr := json.Unmarshal([]byte(stdout), &results); decodeErr != nil {
		t.Fatalf("decode results: %v\n%s", decodeErr, stdout)
	}
	if requests != 1 || len(results) != 1 || results[0].Status != batchStatusFailed || !strings.Contains(results[0].Error, "deleting repository") {
		t.Fatalf("requests=%d results=%+v, want one failed delete", requests, results)
	}
}

func TestBatchApplyFailFastMarksRemainingItemsNotRun(t *testing.T) {
	input := `apiVersion: go-quay/v1alpha1
items:
  - action: create
    namespace: testns
    repository: denied-repo
  - action: create
    namespace: testns
    repository: not-sent-repo
`
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"error":"denied"}`)
	}))
	defer server.Close()

	stdout, err := runBatchCLI(t, input, "-", server.URL+"/api/v1", "--token", testTokenValue, "--fail-fast")
	if err == nil || !strings.Contains(err.Error(), "not run") {
		t.Fatalf("error = %v, want fail-fast aggregate error", err)
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want one request", requests)
	}
	var results []repositoryBatchResult
	if decodeErr := json.Unmarshal([]byte(stdout), &results); decodeErr != nil {
		t.Fatalf("decode results: %v\n%s", decodeErr, stdout)
	}
	if len(results) != 2 || results[0].Status != batchStatusFailed || results[1].Status != batchStatusNotRun {
		t.Fatalf("batch results = %+v, want failed then not_run", results)
	}
}

func TestBatchApplyRequiresConfirmForDelete(t *testing.T) {
	input := `apiVersion: go-quay/v1alpha1
items:
  - action: delete
    namespace: testns
    repository: old-repo
`
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { requests++ }))
	defer server.Close()

	_, err := runBatchCLI(t, input, "-", server.URL+"/api/v1", "--token", testTokenValue)
	if err == nil || !strings.Contains(err.Error(), "pass --confirm") {
		t.Fatalf("error = %v, want delete confirmation error", err)
	}
	if requests != 0 {
		t.Fatalf("unconfirmed batch sent %d requests, want none", requests)
	}
}

func TestValidateRepositoryBatchRejectsInvalidDocuments(t *testing.T) {
	tests := []struct {
		name  string
		batch *repositoryBatch
		want  string
	}{
		{
			name:  "missing api version",
			batch: &repositoryBatch{Items: []repositoryBatchItem{{Action: subcmdCreate, Namespace: testNamespace, Repository: testBatchRepository}}},
			want:  "apiVersion",
		},
		{
			name:  "empty items",
			batch: &repositoryBatch{APIVersion: batchAPIVersion},
			want:  "items must not be empty",
		},
		{
			name:  "missing namespace",
			batch: &repositoryBatch{APIVersion: batchAPIVersion, Items: []repositoryBatchItem{{Action: subcmdCreate, Repository: testBatchRepository}}},
			want:  "requires namespace and repository",
		},
		{
			name:  "missing repository",
			batch: &repositoryBatch{APIVersion: batchAPIVersion, Items: []repositoryBatchItem{{Action: subcmdCreate, Namespace: testNamespace}}},
			want:  "requires namespace and repository",
		},
		{
			name:  "invalid visibility",
			batch: &repositoryBatch{APIVersion: batchAPIVersion, Items: []repositoryBatchItem{{Action: subcmdCreate, Namespace: testNamespace, Repository: testBatchRepository, Visibility: "internal"}}},
			want:  "visibility must be public or private",
		},
		{
			name:  "delete with create-only field",
			batch: &repositoryBatch{APIVersion: batchAPIVersion, Items: []repositoryBatchItem{{Action: subcmdDelete, Namespace: testNamespace, Repository: testBatchRepository, Description: "not allowed"}}},
			want:  "do not accept visibility or description",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRepositoryBatch(tt.batch); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func TestValidateRepositoryBatchDefaultsCreateVisibility(t *testing.T) {
	batch := &repositoryBatch{
		APIVersion: batchAPIVersion,
		Items:      []repositoryBatchItem{{Action: subcmdCreate, Namespace: testNamespace, Repository: testBatchRepository}},
	}
	if err := validateRepositoryBatch(batch); err != nil {
		t.Fatalf("validate batch: %v", err)
	}
	if got := batch.Items[0].Visibility; got != repoVisibilityPrivate {
		t.Fatalf("visibility = %q, want private", got)
	}
}

func TestReadRepositoryBatchRejectsMalformedAndMultipleDocuments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "malformed YAML",
			input: "apiVersion: [\nitems:",
			want:  "decoding batch input",
		},
		{
			name:  "unknown field",
			input: "apiVersion: go-quay/v1alpha1\nunknown: true\nitems: []\n",
			want:  "field unknown not found",
		},
		{
			name:  "multiple documents",
			input: "apiVersion: go-quay/v1alpha1\nitems: []\n---\napiVersion: go-quay/v1alpha1\nitems: []\n",
			want:  "one YAML or JSON document",
		},
		{
			name:  "malformed trailing document",
			input: "apiVersion: go-quay/v1alpha1\nitems: []\n---\ninvalid: [\n",
			want:  "decoding batch input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := t.TempDir() + "/input.yaml"
			if err := os.WriteFile(filePath, []byte(tt.input), 0o600); err != nil {
				t.Fatalf("write batch file: %v", err)
			}
			if _, err := readRepositoryBatch(batchApplyCmd, filePath); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func TestReadRepositoryBatchReportsMissingFile(t *testing.T) {
	_, err := readRepositoryBatch(batchApplyCmd, t.TempDir()+"/missing.yaml")
	if err == nil || !strings.Contains(err.Error(), "opening batch file") {
		t.Fatalf("error = %v, want batch file open error", err)
	}
}
