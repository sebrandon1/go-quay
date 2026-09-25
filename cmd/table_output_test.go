package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runTableCommand(t *testing.T, response string, run func(*cobra.Command) error) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		writeResponse(t, w, []byte(response))
	}))
	t.Cleanup(server.Close)

	oldToken, oldURL, oldFormat := token, quayURL, outputFormat
	oldNamespace, oldRepository, oldOrgName := namespace, repository, orgName
	oldTeamOrg, oldManifest, oldVulnerabilities := teamCmdOrgname, secScanManifestRef, includeVulnerabilities
	token = testTokenValue
	quayURL = server.URL
	outputFormat = outputTable
	namespace = testNamespace
	repository = testRepository
	orgName = testOrgName
	teamCmdOrgname = testOrgName
	secScanManifestRef = "sha256:abc"
	includeVulnerabilities = true
	t.Cleanup(func() {
		token, quayURL, outputFormat = oldToken, oldURL, oldFormat
		namespace, repository, orgName = oldNamespace, oldRepository, oldOrgName
		teamCmdOrgname, secScanManifestRef, includeVulnerabilities = oldTeamOrg, oldManifest, oldVulnerabilities
	})

	command := &cobra.Command{}
	command.SetContext(context.Background())
	var runErr error
	output := captureStdout(t, func() { runErr = run(command) })
	if runErr != nil {
		t.Fatalf("command: %v", runErr)
	}
	return output
}

func requireTableOutput(t *testing.T, output string, headers ...string) {
	t.Helper()
	for _, header := range headers {
		if !strings.Contains(output, header) {
			t.Errorf("table output missing header %q: %s", header, output)
		}
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected a header and data row, got: %s", output)
	}
}

func TestTagListTableOutput(t *testing.T) {
	response := `{"tags":[{"name":"` + testLatestTagName + `","manifest_digest":"sha256:abc","size":123,"last_modified":"2026-09-20","expiration":"never"}]}`
	output := runTableCommand(t, response, func(cmd *cobra.Command) error {
		return tagListCmd.RunE(cmd, nil)
	})
	requireTableOutput(t, output, "TAG", "DIGEST", "SIZE", testLatestTagName, "sha256:abc")
}

func TestRobotListTableOutput(t *testing.T) {
	output := runTableCommand(t, `{"robots":[{"name":"buildbot","description":"CI robot","created":"2026-09-01","last_accessed":"2026-09-20","token":"secret-token"}]}`, func(cmd *cobra.Command) error {
		return robotListCmd.RunE(cmd, nil)
	})
	requireTableOutput(t, output, tableHeaderName, "DESCRIPTION", "CREATED", "LAST ACCESSED", "buildbot", "CI robot")
	if strings.Contains(output, "secret-token") {
		t.Fatal("table output must not include robot tokens")
	}
}

func TestOrganizationMembersTableOutput(t *testing.T) {
	output := runTableCommand(t, `{"members":[{"name":"alice","kind":"user","teams":[{"name":"owners"}],"repositories":["api"]}]}`, func(cmd *cobra.Command) error {
		return orgMembersCmd.RunE(cmd, nil)
	})
	requireTableOutput(t, output, tableHeaderName, "KIND", "TEAMS", "REPOSITORIES", "alice", "user")
}

func TestTeamListTableOutput(t *testing.T) {
	output := runTableCommand(t, `{"teams":[{"name":"operators","role":"admin","member_count":2,"repo_count":3}]}`, func(cmd *cobra.Command) error {
		return teamListCmd.RunE(cmd, nil)
	})
	requireTableOutput(t, output, tableHeaderName, "ROLE", "MEMBERS", "REPOSITORIES", "operators", "admin")
}

func TestSecurityScanTableOutput(t *testing.T) {
	response := `{"status":"scanned","data":{"Layer":{"Features":[{"Name":"openssl","Vulnerabilities":[{"Name":"CVE-1","Severity":"Critical"},{"Name":"CVE-2","Severity":"high"}]},{"Name":"zlib","Vulnerabilities":[{"Name":"CVE-3","Severity":"unlisted"}]}]}}}`
	output := runTableCommand(t, response, func(cmd *cobra.Command) error {
		return secscanInfoCmd.RunE(cmd, nil)
	})
	requireTableOutput(t, output, "STATUS", "FEATURES", "VULNERABILITIES", "CRITICAL", "HIGH", tableSeverityUnknown, "scanned", "3")
	if !strings.Contains(strings.Join(strings.Fields(output), " "), "scanned 2 3 1 1 0 0 0 1") {
		t.Errorf("unexpected security summary counts: %s", output)
	}
}

func TestSecuritySummaryWithoutScan(t *testing.T) {
	var printErr error
	output := captureStdout(t, func() { printErr = printSecuritySummary(nil) })
	if printErr != nil {
		t.Fatalf("printSecuritySummary: %v", printErr)
	}
	fields := strings.Fields(output)
	if len(fields) < 9 || strings.Join(fields[9:], " ") != "0 0 0 0 0 0 0 0" {
		t.Errorf("expected zero-valued summary for missing scan data, got: %s", output)
	}
}
