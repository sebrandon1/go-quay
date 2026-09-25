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

// resetTagFlags resets all tag-related flags to their zero values
// so that Cobra's persistent flag state does not leak between tests.
func resetTagFlags(t *testing.T) {
	t.Helper()
	tagPage = 0
	tagLimit = 0
	for _, name := range []string{"page", "limit"} {
		if flag := tagListCmd.Flags().Lookup(name); flag != nil {
			flag.Changed = false
			if err := flag.Value.Set("0"); err != nil {
				t.Fatalf("reset %s flag: %v", name, err)
			}
		}
	}
	t.Cleanup(func() {
		namespace = ""
		repository = ""
		token = ""
		quayURL = ""
		tagName = ""
		tagExpiration = ""
		tagPage = 0
		tagLimit = 0
		for _, name := range []string{"page", "limit"} {
			if flag := tagListCmd.Flags().Lookup(name); flag != nil {
				flag.Changed = false
				_ = flag.Value.Set("0")
			}
		}
		manifestDigest = ""
		confirmTagDeletion = false

		rootCmd.SetArgs([]string{})
	})
}

func TestTagListCmdPassesPagination(t *testing.T) {
	tests := []struct {
		name      string
		flags     []string
		wantPage  string
		wantLimit string
	}{
		{name: "explicit page and limit", flags: []string{"--page", "2", "--limit", "25"}, wantPage: "2", wantLimit: "25"},
		{name: "Quay API defaults", wantPage: "", wantLimit: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetTagFlags(t)
			resetRootFlags(t)

			var gotPage, gotLimit string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPage = r.URL.Query().Get("page")
				gotLimit = r.URL.Query().Get("limit")
				w.Header().Set("Content-Type", "application/json")
				writeResponse(t, w, []byte(`{"tags": []}`))
			}))
			defer server.Close()

			args := []string{
				cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
				cmdTag, subcmdList, "-n", testNamespace, "-r", testRepository,
			}
			rootCmd.SetArgs(append(args, tt.flags...))
			var executeErr error
			captureStdout(t, func() { executeErr = rootCmd.Execute() })
			if executeErr != nil {
				t.Fatalf("tag list: %v", executeErr)
			}
			if gotPage != tt.wantPage || gotLimit != tt.wantLimit {
				t.Fatalf("pagination query = page %q, limit %q; want page %q, limit %q", gotPage, gotLimit, tt.wantPage, tt.wantLimit)
			}
		})
	}
}

func TestRequireTagName(t *testing.T) {
	original := tagName
	t.Cleanup(func() { tagName = original })
	tagName = ""
	if err := requireTagName(nil, nil); err == nil {
		t.Fatal("expected missing tag name to be rejected")
	}
	tagName = testLatestTagName
	if err := requireTagName(nil, nil); err != nil {
		t.Fatalf("expected tag name to be accepted: %v", err)
	}
}

func TestTagInfoCmd(t *testing.T) {
	resetTagFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		expectedPath := "/repository/" + testNamespace + "/" + testRepository + "/tag/v1.0"
		if !strings.HasSuffix(r.URL.Path, expectedPath) {
			t.Errorf("expected path ending with %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		writeResponse(t, w, []byte(`{"name": "v1.0", "manifest_digest": "sha256:deadbeef", "size": 12345}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdTag, subcmdInfo, "-n", testNamespace, "-r", testRepository, "-T", "v1.0",
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

	if !strings.Contains(output, `"name": "v1.0"`) {
		t.Errorf("expected tag name in output, got: %s", output)
	}
	if !strings.Contains(output, `"manifest_digest": "sha256:deadbeef"`) {
		t.Errorf("expected manifest_digest in output, got: %s", output)
	}
}

func TestTagChangeCmd(t *testing.T) {
	resetTagFlags(t)

	var receivedMethod string
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Discard stdout (change command prints to stderr, not stdout).
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdTag, "change", "-n", testNamespace, "-r", testRepository, "-T", testLatestTagName,
		"--manifest", "sha256:abc123",
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if receivedMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", receivedMethod)
	}
	expectedPath := "/repository/" + testNamespace + "/" + testRepository + "/tag/latest"
	if !strings.HasSuffix(receivedPath, expectedPath) {
		t.Errorf("expected path ending with %s, got %s", expectedPath, receivedPath)
	}
}

func TestTagDeleteWithoutConfirmErrors(t *testing.T) {
	resetTagFlags(t)

	// Discard stdout.
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue,
		cmdTag, subcmdDelete, "-n", testNamespace, "-r", testRepository, "-T", "v1.0",
		// --confirm intentionally omitted
	})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err == nil {
		t.Fatal("expected error when --confirm is not passed, got nil")
	}
	if !strings.Contains(err.Error(), "confirm") {
		t.Errorf("expected error message about --confirm, got: %v", err)
	}
}
