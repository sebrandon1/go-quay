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

// resetOrgFlags resets all organization-related flags to their zero values
// so that Cobra's persistent flag state does not leak between tests.
func resetOrgFlags(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		orgName = ""
		token = ""
		quayURL = ""

		rootCmd.SetArgs([]string{})
	})
}

func TestOrgInfoCmd(t *testing.T) {
	resetOrgFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		expectedPath := "/organization/" + testOrgName
		if !strings.HasSuffix(r.URL.Path, expectedPath) {
			t.Errorf("expected path ending with %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name": "` + testOrgName + `", "email": "admin@testorg.com", "is_org_admin": true}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdOrganization, subcmdInfo, "-o", testOrgName,
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

	if !strings.Contains(output, `"name": "`+testOrgName+`"`) {
		t.Errorf("expected org name in output, got: %s", output)
	}
	if !strings.Contains(output, `"email": "admin@testorg.com"`) {
		t.Errorf("expected email in output, got: %s", output)
	}
}

func TestOrgMembersCmd(t *testing.T) {
	resetOrgFlags(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		expectedPath := "/organization/" + testOrgName + "/members"
		if !strings.HasSuffix(r.URL.Path, expectedPath) {
			t.Errorf("expected path ending with %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"members": [{"name": "user1", "kind": "user"}, {"name": "user2", "kind": "user"}]}`))
	}))
	defer server.Close()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{
		cmdGet, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL,
		cmdOrganization, cmdMembers, "-o", testOrgName,
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

	if !strings.Contains(output, `"name": "user1"`) {
		t.Errorf("expected user1 in output, got: %s", output)
	}
	if !strings.Contains(output, `"name": "user2"`) {
		t.Errorf("expected user2 in output, got: %s", output)
	}
}
