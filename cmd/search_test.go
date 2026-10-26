package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	testSearchQuery      = "testquery"
	testSearchRepoName   = "search-result-repo"
	testSearchEntityName = "search-result-user"
	testSearchQueryParam = "query"
	testSearchPageParam  = "page"
)

// preserveSearchFlags restores the package-level search flags after a case that
// sets them through flag execution: resetRootFlags does not cover them, so
// without this a --query from one case would leak into the next.
func preserveSearchFlags(t *testing.T) {
	t.Helper()
	setCommandTestValue(t, &searchQuery, searchQuery)
	setCommandTestValue(t, &searchPage, searchPage)
}

func TestSearchCommands(t *testing.T) {
	repoResponse := `{"results":[{"namespace":"` + testNamespace + `","name":"` + testSearchRepoName + `","kind":"repository"}],"page":2}`
	allResponse := `{"results":[{"name":"` + testSearchEntityName + `","kind":"user"}]}`
	setupQuery := func(t *testing.T) { setCommandTestValue(t, &searchQuery, testSearchQuery) }
	searchReposPath := "/find/repositories"
	searchAllPath := "/find/all"

	tests := []commandDomainCase{
		{
			name:        "search repositories",
			executeArgs: []string{cmdGet, searchCmd.Use, searchReposCmd.Use, "-q", testSearchQuery, "-p", "2"},
			request:     commandRequestExpectation{method: testHTTPGet, path: searchReposPath, query: map[string]string{testSearchQueryParam: testSearchQuery, testSearchPageParam: "2"}},
			response:    repoResponse,
			wantOutput:  testSearchRepoName,
			setup:       preserveSearchFlags,
		},
		{
			name:        "search all",
			executeArgs: []string{cmdGet, searchCmd.Use, searchAllCmd.Use, "-q", testSearchQuery},
			request:     commandRequestExpectation{method: testHTTPGet, path: searchAllPath, query: map[string]string{testSearchQueryParam: testSearchQuery}},
			response:    allResponse,
			wantOutput:  testSearchEntityName,
			setup:       preserveSearchFlags,
		},
		apiErrorCase("search repositories", searchReposCmd.RunE, testHTTPGet, searchReposPath, 500, "searching repositories", setupQuery),
		apiErrorCase("search all", searchAllCmd.RunE, testHTTPGet, searchAllPath, 500, "searching:", setupQuery),
	}
	runCommandDomainCases(t, tests)
}

// TestSearchRepositoriesDefaultPageAndBanner covers what the shared harness
// cannot: runCommandAgainstHTTP captures stdout only, so it can neither see the
// stderr banner nor assert that a query parameter is absent.
func TestSearchRepositoriesDefaultPageAndBanner(t *testing.T) {
	pageFlag := searchReposCmd.Flags().Lookup(testSearchPageParam)
	if pageFlag == nil || pageFlag.DefValue != "0" {
		t.Fatalf("--%s default = %v, want 0 so the page parameter is omitted", testSearchPageParam, pageFlag)
	}

	var gotQuery url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		writeResponse(t, w, []byte(`{"results":[]}`))
	}))
	defer server.Close()

	resetRootFlags(t)
	resetCommandFlagChanges(rootCmd)
	setCommandTestValue(t, &searchPage, 0)

	var stdout string
	stderr := captureStderr(t, func() {
		stdout = captureStdout(t, func() {
			rootCmd.SetArgs([]string{
				testTokenFlag, testTokenValue, testQuayURLFlag, server.URL + "/api/v1",
				cmdGet, searchCmd.Use, searchReposCmd.Use, "-q", testSearchQuery,
			})
			rootCmd.SetErr(io.Discard)
			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("search repositories: %v", err)
			}
		})
	})

	if gotQuery.Has(testSearchPageParam) {
		t.Errorf("request query = %q, want no %s parameter when the flag is unset", gotQuery.Encode(), testSearchPageParam)
	}
	if got := gotQuery.Get(testSearchQueryParam); got != testSearchQuery {
		t.Errorf("request query %s = %q, want %q", testSearchQueryParam, got, testSearchQuery)
	}
	if want := "Repository search results for: " + testSearchQuery; !strings.Contains(stderr, want) {
		t.Errorf("stderr = %q, want it to contain %q", stderr, want)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(stdout), &decoded); err != nil {
		t.Errorf("stdout = %q, want pure JSON: %v", stdout, err)
	}
}
