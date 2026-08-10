package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	testQuayURLFlag = "--quay-url"
	testTokenFlag   = "--token"
	testTokenValue  = "test-token"
	testNamespace   = "testns"
	testRepository  = "testrepo"
	testOrgName     = "testorg"
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
			_, _ = w.Write([]byte(`{"tags": [{"name": "latest", "manifest_digest": "sha256:abc123"}]}`))
		case strings.HasSuffix(r.URL.Path, "/repository/"+testNamespace+"/"+testRepository):
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"namespace": "` + testNamespace + `", "name": "` + testRepository + `", "is_public": true}`))
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
	_, _ = io.Copy(&buf, r)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"repositories": [{"name": "repo1", "namespace": "` + testNamespace + `"}, {"name": "repo2", "namespace": "` + testNamespace + `"}]}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, subcmdList, "-n", testNamespace,
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, `"name": "repo1"`) {
		t.Errorf("expected repo1 in output, got: %s", output)
	}
	if !strings.Contains(output, `"name": "repo2"`) {
		t.Errorf("expected repo2 in output, got: %s", output)
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
		_, _ = w.Write([]byte(`{"namespace": "` + testNamespace + `", "name": "newrepo", "is_public": false}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdRepository, subcmdCreate, "-n", testNamespace, "-r", "newrepo",
		"--visibility", "private", "--description", "a test repo",
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, `"name": "newrepo"`) {
		t.Errorf("expected newrepo in output, got: %s", output)
	}
	if !strings.Contains(output, `"namespace": "`+testNamespace+`"`) {
		t.Errorf("expected namespace in output, got: %s", output)
	}
}
