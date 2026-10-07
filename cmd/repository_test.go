package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const (
	testQuayURLFlag       = "--quay-url"
	testUnreachableURL    = "http://127.0.0.1:1"
	testTokenFlag         = "--token"
	testManifestFlag      = "--manifest"
	testWatchFlag         = "--watch"
	testTokenValue        = "test-token"
	testDescriptionFlag   = "--description"
	testVisibilityFlag    = "--visibility"
	testPrivateVisibility = "private"
	testNamespace         = "testns"
	testRepository        = "testrepo"
	testOrgName           = "testorg"
)

// resetRepositoryFlags resets all repository-related flags to their zero values
// so that Cobra's persistent flag state does not leak between tests.
func resetRepositoryFlags(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		namespace = ""
		repository = ""
		token = ""
		quayURL = ""
		repoVisibility = ""
		repoDescription = ""
		confirmDeletion = false
		repoPublic = false
		repoStarred = false
		repoPopularity = false
		repoTable = false
		repoPage = 0
		repoLimit = 0

		rootCmd.SetArgs([]string{})
	})
}

func TestRepoInfoCmd(t *testing.T) {
	resetRepositoryFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "/repository/"+testNamespace+"/"+testRepository+"/tag/"):
			// ListTags call made by GetRepository
			w.WriteHeader(http.StatusOK)
			writeResponse(t, w, []byte(`{"tags": [{"name": "latest", "manifest_digest": "sha256:abc123"}]}`))
		case strings.HasSuffix(r.URL.Path, "/repository/"+testNamespace+"/"+testRepository):
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			writeResponse(t, w, []byte(`{"namespace": "`+testNamespace+`", "name": "`+testRepository+`", "is_public": true}`))
		default:
			t.Errorf("unexpected request path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, subcmdInfo, "-n", testNamespace, "-r", testRepository,
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `"namespace": "`+testNamespace+`"`) {
		t.Errorf("expected namespace in output, got: %s", output)
	}
	if !strings.Contains(output, `"name": "`+testRepository+`"`) {
		t.Errorf("expected repo name in output, got: %s", output)
	}
	if !strings.Contains(output, `"is_public": true`) {
		t.Errorf("expected is_public in output, got: %s", output)
	}
}

func TestRepoListCmd(t *testing.T) {
	resetRepositoryFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/repository") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("namespace") != testNamespace {
			t.Errorf("expected namespace=%s, got %s", testNamespace, q.Get("namespace"))
		}
		if q.Get("page") != "2" || q.Get("limit") != "5" {
			t.Errorf("expected page=2 and limit=5, got page=%s limit=%s", q.Get("page"), q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, []byte(`{"repositories": [{"name": "repo1", "namespace": "`+testNamespace+`"}, {"name": "repo2", "namespace": "`+testNamespace+`"}]}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, subcmdList, "-n", testNamespace, "--page", "2", "--limit", "5",
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `"name": "repo1"`) {
		t.Errorf("expected repo1 in output, got: %s", output)
	}
	if !strings.Contains(output, `"name": "repo2"`) {
		t.Errorf("expected repo2 in output, got: %s", output)
	}
}

func TestRepoListTableParallelTagFetches(t *testing.T) {
	resetRepositoryFlags(t)

	repositories, repositoryJSON := repositoryListTestData(9)
	var maxActiveRequests atomic.Int32
	server := newRepositoryListTableTestServer(t, repositoryJSON, &maxActiveRequests)
	defer server.Close()

	output := executeRepositoryListTableCommand(t, server.URL)
	assertRepositoryListTableOutput(t, output, repositories)

	if got := maxActiveRequests.Load(); got < 2 || got > repositoryTagFetchConcurrency {
		t.Errorf("maximum concurrent tag requests = %d, want between 2 and %d", got, repositoryTagFetchConcurrency)
	}
}

func repositoryListTestData(count int) ([]string, string) {
	repositories := make([]string, count)
	entries := make([]string, count)
	for i := range repositories {
		name := "repo" + strconv.Itoa(i+1)
		repositories[i] = name
		// Descending popularity keeps the expected table order explicit.
		entries[i] = fmt.Sprintf(`{"name":%q,"popularity":%d}`, name, count-i)
	}
	return repositories, `{"repositories":[` + strings.Join(entries, ",") + `]}`
}

func newRepositoryListTableTestServer(t *testing.T, repositoryJSON string, maxActiveRequests *atomic.Int32) *httptest.Server {
	t.Helper()
	var activeRequests atomic.Int32
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repository" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(repositoryJSON))
			return
		}

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 2 || pathParts[len(pathParts)-1] != "tag" {
			t.Errorf("unexpected request path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		repoName := pathParts[len(pathParts)-2]
		active := activeRequests.Add(1)
		defer activeRequests.Add(-1)
		recordMaximumConcurrency(maxActiveRequests, active)

		// Keep requests in flight long enough to observe overlapping fetches.
		time.Sleep(20 * time.Millisecond)
		if repoName == "repo5" {
			http.Error(w, "tag lookup failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tags":[{"name":"v1","start_ts":1780000000,"is_manifest_list":true}]}`))
	}))
}

func recordMaximumConcurrency(maximum *atomic.Int32, active int32) {
	for previous := maximum.Load(); active > previous; previous = maximum.Load() {
		if maximum.CompareAndSwap(previous, active) {
			return
		}
	}
}

func executeRepositoryListTableCommand(t *testing.T, serverURL string) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = w
	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, serverURL,
		cmdRepository, subcmdList, "-n", testNamespace, "--popularity", "--table",
	})
	execErr := rootCmd.Execute()
	if closeErr := w.Close(); closeErr != nil {
		t.Errorf("close stdout pipe: %v", closeErr)
	}
	os.Stdout = oldStdout
	if execErr != nil {
		t.Fatalf("expected no error, got: %v", execErr)
	}

	var output bytes.Buffer
	if _, err := io.Copy(&output, r); err != nil {
		t.Fatalf("read captured stdout: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Errorf("close stdout reader: %v", err)
	}
	return output.String()
}

func assertRepositoryListTableOutput(t *testing.T, output string, repositories []string) {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(repositories)+1 {
		t.Fatalf("expected header and %d repository rows, got %d lines: %s", len(repositories), len(lines)-1, output)
	}
	for i, expectedName := range repositories {
		fields := strings.Fields(lines[i+1])
		if len(fields) == 0 || fields[0] != expectedName {
			t.Errorf("row %d repository = %q, want %q", i+1, lines[i+1], expectedName)
		}
		if expectedName == "repo5" {
			if len(fields) < 7 || strings.Join(fields[1:], " ") != "5 0 0 - - -" {
				t.Errorf("failed tag lookup row should retain placeholder values, got %q", lines[i+1])
			}
		} else if len(fields) < 7 || fields[3] != "1" || fields[4] != "v1" || fields[6] != "yes" {
			t.Errorf("successful tag lookup row has unexpected values: %q", lines[i+1])
		}
	}
}

func TestRepoInfoMissingRepositoryFlag(t *testing.T) {
	resetRepositoryFlags(t)

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue,
		cmdRepository, subcmdInfo, "-n", testNamespace,
		// --repository flag intentionally omitted
	})

	// Discard stdout to keep test output clean.
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err == nil {
		t.Fatal("expected error for missing --repository flag, got nil")
	}
	if !strings.Contains(err.Error(), "repository") {
		t.Errorf("expected error about repository flag, got: %v", err)
	}
}

func TestRepoCreateCmd(t *testing.T) {
	resetRepositoryFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/repository") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		writeResponse(t, w, []byte(`{"namespace": "`+testNamespace+`", "name": "newrepo", "is_public": false}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, subcmdCreate, "-n", testNamespace, "-r", "newrepo",
		testVisibilityFlag, testPrivateVisibility, testDescriptionFlag, "a test repo",
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `"name": "newrepo"`) {
		t.Errorf("expected newrepo in output, got: %s", output)
	}
	if !strings.Contains(output, `"namespace": "`+testNamespace+`"`) {
		t.Errorf("expected namespace in output, got: %s", output)
	}
}

func TestVerbRepoCreateCmd(t *testing.T) {
	resetRepositoryFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/repository") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		writeResponse(t, w, []byte(`{"namespace": "`+testNamespace+`", "name": "newrepo", "is_public": false}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdCreate, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, "-n", testNamespace, "-r", "newrepo",
		testVisibilityFlag, testPrivateVisibility, testDescriptionFlag, "a test repo",
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `"name": "newrepo"`) {
		t.Errorf("expected newrepo in output, got: %s", output)
	}
	if !strings.Contains(output, `"namespace": "`+testNamespace+`"`) {
		t.Errorf("expected namespace in output, got: %s", output)
	}
}

func TestVerbRepoListCmd(t *testing.T) {
	resetRepositoryFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/repository") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, []byte(`{"repositories": [{"name": "repo1", "namespace": "`+testNamespace+`"}]}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdList, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, "-n", testNamespace,
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `"name": "repo1"`) {
		t.Errorf("expected repo1 in output, got: %s", output)
	}
}

func TestVerbRepoInfoCmd(t *testing.T) {
	resetRepositoryFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "/repository/"+testNamespace+"/"+testRepository+"/tag/"):
			w.WriteHeader(http.StatusOK)
			writeResponse(t, w, []byte(`{"tags": [{"name": "latest", "manifest_digest": "sha256:abc123"}]}`))
		case strings.HasSuffix(r.URL.Path, "/repository/"+testNamespace+"/"+testRepository):
			w.WriteHeader(http.StatusOK)
			writeResponse(t, w, []byte(`{"namespace": "`+testNamespace+`", "name": "`+testRepository+`", "is_public": true}`))
		default:
			t.Errorf("unexpected request path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdInfo, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, "-n", testNamespace, "-r", testRepository,
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, `"name": "`+testRepository+`"`) {
		t.Errorf("expected repo name in output, got: %s", output)
	}
}

func TestVerbRepoDeleteRequiresConfirm(t *testing.T) {
	resetRepositoryFlags(t)

	rootCmd.SetArgs([]string{
		cmdDelete, testTokenFlag, testTokenValue, testQuayURLFlag, testUnreachableURL,
		cmdRepository, "-n", testNamespace, "-r", testRepository,
	})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error without --confirm")
	}
	if !strings.Contains(err.Error(), "confirm") {
		t.Errorf("expected confirm error, got: %v", err)
	}
}

func TestVerbRepoDeleteCmd(t *testing.T) {
	resetRepositoryFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		wantPath := "/repository/" + testNamespace + "/" + testRepository
		if !strings.HasSuffix(r.URL.Path, wantPath) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	rootCmd.SetArgs([]string{
		cmdDelete, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, "-n", testNamespace, "-r", testRepository, "--confirm",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
