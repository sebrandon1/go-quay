package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sebrandon1/go-quay/lib"
)

const (
	testNamespaceFlagName       = "namespace"
	testDiffRepositoryFlagName  = "repository"
	testDiffRepositoryFlag      = "--repository"
	testDiffTagAFlagName        = "tag-a"
	testDiffTagBFlagName        = "tag-b"
	testDiffOldTagName          = "old"
	testDiffNewTagName          = "new"
	testDiffOldManifestDigest   = "sha256:old"
	testDiffNewManifestDigest   = "sha256:new"
	testDiffVulnerabilitiesTrue = "true"
	testDiffSecurityStatus      = "scanned"
)

func resetTagDiffFlags(t *testing.T) {
	t.Helper()
	diffNamespace = appCfg.Namespace
	diffRepository = ""
	diffTagA = ""
	diffTagB = ""
	for _, name := range []string{testNamespaceFlagName, testDiffRepositoryFlagName, testDiffTagAFlagName, testDiffTagBFlagName} {
		if f := diffTagCmd.Flags().Lookup(name); f != nil {
			f.Changed = false
			value := ""
			if name == testNamespaceFlagName {
				value = appCfg.Namespace
			}
			if err := f.Value.Set(value); err != nil {
				t.Fatalf("reset %s flag: %v", name, err)
			}
		}
	}
}

func TestDiffTagCommandShowsLabelAndVulnerabilityDeltas(t *testing.T) {
	resetRootFlags(t)
	resetTagDiffFlags(t)
	server, requestPaths := startTagDiffTestServer(t)
	defer server.Close()

	var stderr bytes.Buffer
	rootCmd.SetErr(&stderr)
	rootCmd.SetArgs([]string{
		diffCmd.Use, diffTagCmd.Use, testConfigNamespaceFlag, testNamespace, testDiffRepositoryFlag, testRepository,
		"--tag-a", testDiffOldTagName, "--tag-b", testDiffNewTagName, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL + "/api/v1",
	})
	var runErr error
	stdout := captureStdout(t, func() { runErr = rootCmd.Execute() })
	if runErr != nil {
		t.Fatalf("execute tag diff: %v", runErr)
	}
	if got := len(requestPaths()); got != 6 {
		t.Fatalf("got %d API requests, want 6: %v", got, requestPaths())
	}
	assertTagDiffJSON(t, stdout)
}

func assertTagDiffJSON(t *testing.T, stdout string) {
	t.Helper()
	var result tagDiffResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("decode diff JSON: %v\n%s", err, stdout)
	}
	if result.TagA.Name != testDiffOldTagName || result.TagA.ManifestDigest != testDiffOldManifestDigest || result.TagA.SecurityStatus != testDiffSecurityStatus {
		t.Errorf("tag A = %+v", result.TagA)
	}
	if result.TagB.Name != testDiffNewTagName || result.TagB.ManifestDigest != testDiffNewManifestDigest || result.TagB.SecurityStatus != testDiffSecurityStatus {
		t.Errorf("tag B = %+v", result.TagB)
	}
	assertTagDiffLabels(t, result)
	assertTagDiffVulnerabilities(t, result)
}

func assertTagDiffLabels(t *testing.T, result tagDiffResult) {
	t.Helper()
	if len(result.AddedLabels) != 2 || result.AddedLabels[0].Key != "build" || result.AddedLabels[1].Value != "2" {
		t.Errorf("added labels = %+v", result.AddedLabels)
	}
	if len(result.RemovedLabels) != 2 || result.RemovedLabels[0].Key != "source" || result.RemovedLabels[1].Value != "1" {
		t.Errorf("removed labels = %+v", result.RemovedLabels)
	}
}

func assertTagDiffVulnerabilities(t *testing.T, result tagDiffResult) {
	t.Helper()
	if len(result.AddedVulnerabilities) != 1 || result.AddedVulnerabilities[0] != "CVE-3" {
		t.Errorf("added vulnerabilities = %v", result.AddedVulnerabilities)
	}
	if len(result.RemovedVulnerabilities) != 1 || result.RemovedVulnerabilities[0] != "CVE-1" {
		t.Errorf("removed vulnerabilities = %v", result.RemovedVulnerabilities)
	}
}

func startTagDiffTestServer(t *testing.T) (*httptest.Server, func() []string) {
	t.Helper()
	var requestPaths []string
	tagPath := "/api/v1/repository/testns/testrepo/tag/"
	manifestPath := "/api/v1/repository/testns/testrepo/manifest/"
	responses := map[string]any{
		tagPath + testDiffOldTagName: lib.Tag{Name: testDiffOldTagName, ManifestDigest: testDiffOldManifestDigest},
		tagPath + testDiffNewTagName: lib.Tag{Name: testDiffNewTagName, ManifestDigest: testDiffNewManifestDigest},
		manifestPath + testDiffOldManifestDigest + "/labels": lib.ManifestLabels{Labels: []lib.ManifestLabel{
			{ID: "old-version", Key: "version", Value: "1"},
			{ID: "old-team", Key: "team", Value: "security"},
			{ID: "old-source", Key: "source", Value: "legacy"},
		}},
		manifestPath + testDiffNewManifestDigest + "/labels": lib.ManifestLabels{Labels: []lib.ManifestLabel{
			{ID: "new-version", Key: "version", Value: "2"},
			{ID: "new-team", Key: "team", Value: "security"},
			{ID: "new-build", Key: "build", Value: "prod"},
		}},
		manifestPath + testDiffOldManifestDigest + "/security": testSecurityScan("CVE-1", "CVE-2", "CVE-2"),
		manifestPath + testDiffNewManifestDigest + "/security": testSecurityScan("CVE-2", "CVE-3"),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPaths = append(requestPaths, r.URL.Path)
		if r.Method != http.MethodGet {
			t.Errorf("request method = %s, want GET", r.Method)
		}
		if strings.HasSuffix(r.URL.Path, "/security") && r.URL.Query().Get("vulnerabilities") != testDiffVulnerabilitiesTrue {
			t.Errorf("security scan request must include vulnerabilities=true")
		}
		response, found := responses[r.URL.Path]
		if !found {
			t.Errorf("unexpected request path %q", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
	return server, func() []string { return requestPaths }
}

func TestDiffTagCommandRequiresTagFlags(t *testing.T) {
	resetRootFlags(t)
	resetTagDiffFlags(t)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{
		diffCmd.Use, diffTagCmd.Use, testConfigNamespaceFlag, testNamespace, testDiffRepositoryFlag, testRepository,
		testTokenFlag, testTokenValue, testQuayURLFlag, testUnreachableURL,
	})
	if err := rootCmd.Execute(); err == nil || !strings.Contains(err.Error(), testDiffTagAFlagName) || !strings.Contains(err.Error(), testDiffTagBFlagName) {
		t.Fatalf("error = %v, want missing tag flag error", err)
	}
}

func TestDiffTagCommandReportsMissingTag(t *testing.T) {
	resetRootFlags(t)
	resetTagDiffFlags(t)
	var requestCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	}))
	defer server.Close()
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{
		diffCmd.Use, diffTagCmd.Use, testConfigNamespaceFlag, testNamespace, testDiffRepositoryFlag, testRepository,
		"--tag-a", "missing", "--tag-b", testDiffNewTagName, testTokenFlag, testTokenValue, testQuayURLFlag, server.URL + "/api/v1",
	})
	if err := rootCmd.Execute(); err == nil || !strings.Contains(err.Error(), `getting tag "missing"`) {
		t.Fatalf("error = %v, want clear missing tag error", err)
	}
	if requestCount != 1 {
		t.Errorf("request count = %d, want 1", requestCount)
	}
}

func TestCompareTagsReportsFailuresByStage(t *testing.T) {
	tagPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/tag/"
	manifestPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/manifest/"
	tests := []struct {
		name            string
		failPath        string
		emptyDigestPath string
		wantError       string
	}{
		{name: "second tag request", failPath: tagPath + testDiffNewTagName, wantError: `getting tag "new"`},
		{name: "first tag missing digest", emptyDigestPath: tagPath + testDiffOldTagName, wantError: `tag "old" did not return a manifest digest`},
		{name: "second tag missing digest", emptyDigestPath: tagPath + testDiffNewTagName, wantError: `tag "new" did not return a manifest digest`},
		{name: "first tag labels", failPath: manifestPath + testDiffOldManifestDigest + "/labels", wantError: `getting labels for tag "old"`},
		{name: "first tag security", failPath: manifestPath + testDiffOldManifestDigest + "/security", wantError: `getting security scan for tag "old"`},
		{name: "second tag labels", failPath: manifestPath + testDiffNewManifestDigest + "/labels", wantError: `getting labels for tag "new"`},
		{name: "second tag security", failPath: manifestPath + testDiffNewManifestDigest + "/security", wantError: `getting security scan for tag "new"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responses := map[string]any{
				tagPath + testDiffOldTagName:                           lib.Tag{Name: testDiffOldTagName, ManifestDigest: testDiffOldManifestDigest},
				tagPath + testDiffNewTagName:                           lib.Tag{Name: testDiffNewTagName, ManifestDigest: testDiffNewManifestDigest},
				manifestPath + testDiffOldManifestDigest + "/labels":   lib.ManifestLabels{},
				manifestPath + testDiffNewManifestDigest + "/labels":   lib.ManifestLabels{},
				manifestPath + testDiffOldManifestDigest + "/security": testSecurityScan(),
				manifestPath + testDiffNewManifestDigest + "/security": testSecurityScan(),
			}
			if tt.emptyDigestPath != "" {
				responses[tt.emptyDigestPath] = lib.Tag{}
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == tt.failPath {
					http.Error(w, `{"error":"unavailable"}`, http.StatusServiceUnavailable)
					return
				}
				response, found := responses[r.URL.Path]
				if !found {
					http.NotFound(w, r)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Errorf("encode response: %v", err)
				}
			}))
			defer server.Close()

			client, err := lib.NewClientWithURL(testTokenValue, server.URL+"/api/v1")
			if err != nil {
				t.Fatalf("create client: %v", err)
			}
			_, err = compareTags(context.Background(), client, testNamespace, testRepository, testDiffOldTagName, testDiffNewTagName)
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error = %v, want error containing %q", err, tt.wantError)
			}
		})
	}
}

func testSecurityScan(vulnerabilityIDs ...string) lib.SecurityScan {
	vulnerabilities := make([]lib.SecurityVulnerability, 0, len(vulnerabilityIDs))
	for _, id := range vulnerabilityIDs {
		vulnerabilities = append(vulnerabilities, lib.SecurityVulnerability{Name: id})
	}
	return lib.SecurityScan{
		Status: testDiffSecurityStatus,
		Data:   &lib.SecurityData{Layer: &lib.SecurityLayer{Features: []lib.SecurityFeature{{Vulnerabilities: vulnerabilities}}}},
	}
}
