package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sebrandon1/go-quay/lib"
)

const (
	pullTrafficTestStart        = "10/06/2026"
	pullTrafficRepositoryPath   = "/api/v1/repository"
	testPullTrafficAlpha        = "alpha"
	testPullTrafficBeta         = "beta"
	testPullTrafficGamma        = "gamma"
	testPullTrafficErrorKey     = "error"
	testKindPullTag             = "pull_tag"
	testPullTrafficEndDate      = "2026-11-05"
	testPullTrafficCoveredStart = "2026-10-06"
	testConsumerIP              = "1.1.1.1"
	testConsumerUserAgent       = "cri-o/1.35"
	testConsumerPuller          = "robot-a"
	testConsumerIPSecondary     = "2.2.2.2"
	testConsumerLastSeen        = "2026-11-05 20:00"
)

func newPullTrafficTestClient(t *testing.T, handler http.HandlerFunc) *lib.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := lib.NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}
	return client
}

func writeJSON(t *testing.T, w http.ResponseWriter, value interface{}) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Errorf("marshaling response: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

// pullTrafficHandler serves a repository list plus per-repository aggregates.
func pullTrafficHandler(t *testing.T, repositories []string, counts map[string]map[string]int) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pullTrafficRepositoryPath {
			list := lib.RepositoryList{}
			for _, name := range repositories {
				list.Repositories = append(list.Repositories, lib.OrganizationRepository{Name: name})
			}
			writeJSON(t, w, list)
			return
		}

		prefix := "/api/v1/repository/" + testNamespace + "/"
		if !strings.HasPrefix(r.URL.Path, prefix) || !strings.HasSuffix(r.URL.Path, "/aggregatelogs") {
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		name := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, prefix), "/aggregatelogs")
		aggregated := lib.AggregatedLogs{}
		for kind, count := range counts[name] {
			aggregated.Aggregated = append(aggregated.Aggregated, lib.AggregatedLogEntry{Kind: kind, Count: count})
		}
		writeJSON(t, w, aggregated)
	}
}

func TestCollectPullTrafficRanksRepositories(t *testing.T) {
	client := newPullTrafficTestClient(t, pullTrafficHandler(t,
		[]string{testPullTrafficAlpha, testPullTrafficBeta, testPullTrafficGamma},
		map[string]map[string]int{
			testPullTrafficAlpha: {pullTrafficKindPull: 100, pullTrafficKindPush: 5},
			testPullTrafficBeta:  {pullTrafficKindPull: 300},
			testPullTrafficGamma: {},
		}))

	report, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: pullTrafficTestStart,
		EndDate:   pullTrafficTestStart,
	})
	if err != nil {
		t.Fatalf("collectPullTraffic returned error: %v", err)
	}
	if report.TotalPulls != 400 || report.TotalPushes != 5 {
		t.Errorf("totals = %d pulls / %d pushes, want 400 / 5", report.TotalPulls, report.TotalPushes)
	}
	wantOrder := []string{testPullTrafficBeta, testPullTrafficAlpha, testPullTrafficGamma}
	for i, want := range wantOrder {
		if got := report.Repositories[i].Repository; got != want {
			t.Fatalf("row %d = %s, want %s (order %v)", i, got, want, wantOrder)
		}
	}
	if share := report.Repositories[0].SharePercent; share != 75 {
		t.Errorf("beta share = %v, want 75", share)
	}
	if share := report.Repositories[1].SharePercent; share != 25 {
		t.Errorf("alpha share = %v, want 25", share)
	}
}

func TestCollectPullTrafficTopLimitsRows(t *testing.T) {
	client := newPullTrafficTestClient(t, pullTrafficHandler(t,
		[]string{testPullTrafficAlpha, testPullTrafficBeta, testPullTrafficGamma},
		map[string]map[string]int{
			testPullTrafficAlpha: {pullTrafficKindPull: 100},
			testPullTrafficBeta:  {pullTrafficKindPull: 300},
			testPullTrafficGamma: {pullTrafficKindPull: 50},
		}))

	report, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: pullTrafficTestStart,
		EndDate:   pullTrafficTestStart,
		Top:       2,
	})
	if err != nil {
		t.Fatalf("collectPullTraffic returned error: %v", err)
	}
	if len(report.Repositories) != 2 {
		t.Fatalf("rows = %d, want 2", len(report.Repositories))
	}
	if report.TotalPulls != 450 {
		t.Errorf("total pulls = %d, want 450 from every repository, not just the top rows", report.TotalPulls)
	}
}

func TestCollectPullTrafficRejectsInvalidDate(t *testing.T) {
	client := newPullTrafficTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request %s %s for an invalid date", r.Method, r.URL.Path)
	})

	_, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: "not-a-date",
		EndDate:   pullTrafficTestStart,
	})
	if err == nil || !strings.Contains(err.Error(), "invalid log date") {
		t.Fatalf("error = %v, want an invalid-date error", err)
	}
}

func TestCollectPullTrafficContinuedOnRepoError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pullTrafficRepositoryPath {
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}, {Name: "broken"}}})
			return
		}
		if strings.Contains(r.URL.Path, "/broken/") {
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(t, w, map[string]string{testPullTrafficErrorKey: "broken"})
			return
		}
		writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 7}}})
	}
	client := newPullTrafficTestClient(t, handler)

	report, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: pullTrafficTestStart,
		EndDate:   pullTrafficTestStart,
	})
	if err != nil {
		t.Fatalf("collectPullTraffic returned error: %v", err)
	}
	if report.Repositories[0].Repository != testPullTrafficAlpha || report.Repositories[0].Pulls != 7 {
		t.Errorf("healthy repository = %+v, want alpha with 7 pulls", report.Repositories[0])
	}
	if report.Repositories[1].Error == "" {
		t.Errorf("broken repository has no recorded error: %+v", report.Repositories[1])
	}
}

func TestCollectPullTrafficFailsWhenEveryRepoFails(t *testing.T) {
	client := newPullTrafficTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pullTrafficRepositoryPath {
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
			return
		}
		w.WriteHeader(http.StatusForbidden)
		writeJSON(t, w, map[string]string{testPullTrafficErrorKey: "denied"})
	})

	_, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: pullTrafficTestStart,
		EndDate:   pullTrafficTestStart,
	})
	if err == nil || !strings.Contains(err.Error(), "no repository logs could be read") {
		t.Fatalf("error = %v, want a no-logs error", err)
	}
}

func TestCollectPullTrafficByTag(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == pullTrafficRepositoryPath:
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
		case strings.HasSuffix(r.URL.Path, "/aggregatelogs"):
			writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 9}}})
		case strings.HasSuffix(r.URL.Path, "/logs"):
			writeJSON(t, w, lib.Logs{Logs: []lib.LogEntry{
				{Kind: pullTrafficKindPull, Metadata: lib.Metadata{Tag: "v1"}},
				{Kind: pullTrafficKindPull, Metadata: lib.Metadata{Tag: "v1"}},
				{Kind: pullTrafficKindPull, Metadata: lib.Metadata{Tag: "v2"}},
				{Kind: pullTrafficKindPull},
				{Kind: pullTrafficKindPush, Metadata: lib.Metadata{Tag: "v3"}},
			}})
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	}
	client := newPullTrafficTestClient(t, handler)

	report, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: pullTrafficTestStart,
		EndDate:   pullTrafficTestStart,
		ByTag:     true,
		MaxPages:  5,
	})
	if err != nil {
		t.Fatalf("collectPullTraffic returned error: %v", err)
	}
	row := report.Repositories[0]
	if len(row.TopTags) != 2 || row.TopTags[0].Tag != "v1" || row.TopTags[0].Count != 2 {
		t.Fatalf("top tags = %+v, want v1 twice then v2", row.TopTags)
	}
	if row.UntaggedPulls != 1 {
		t.Errorf("untagged pulls = %d, want 1", row.UntaggedPulls)
	}
}

func TestRankPullTrafficTags(t *testing.T) {
	tags, untagged := rankPullTrafficTags([]lib.LogEntry{
		{Kind: testKindPullTag, Metadata: lib.Metadata{Tag: "b"}},
		{Kind: testKindPullTag, Metadata: lib.Metadata{Tag: "a"}},
		{Kind: testKindPullTag, Metadata: lib.Metadata{Tag: "a"}},
		{Kind: "pull_repo", Metadata: lib.Metadata{Tag: "a"}},
		{Kind: "pull_repo"},
	})
	if untagged != 1 {
		t.Errorf("untagged = %d, want 1", untagged)
	}
	if len(tags) != 2 || tags[0].Tag != "a" || tags[0].Count != 3 {
		t.Fatalf("tags = %+v, want a=3 then b=1", tags)
	}

	tied, _ := rankPullTrafficTags([]lib.LogEntry{
		{Kind: "pull_tag", Metadata: lib.Metadata{Tag: "b"}},
		{Kind: "pull_tag", Metadata: lib.Metadata{Tag: "a"}},
	})
	if tied[0].Tag != "a" {
		t.Errorf("tied tags = %+v, want the alphabetically first tag to win", tied)
	}
}

func TestPullTrafficRowStatusDerivation(t *testing.T) {
	tests := []struct {
		name string
		row  pullTrafficRepo
		want reportRowStatus
	}{
		{name: "complete", row: pullTrafficRepo{}, want: rowStatusOK},
		{name: "error wins", row: pullTrafficRepo{Error: "boom", TagsTruncated: true}, want: rowStatusError},
		{name: "tags truncated", row: pullTrafficRepo{TagsTruncated: true}, want: rowStatusTruncated},
		{name: "consumers truncated", row: pullTrafficRepo{ConsumersTruncated: true}, want: rowStatusTruncated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.row.setStatus()
			if tt.row.Status != tt.want {
				t.Fatalf("status = %q, want %q", tt.row.Status, tt.want)
			}
		})
	}
}

func TestCollectPullTrafficRowStatuses(t *testing.T) {
	const errorRepo = "repo-12"
	repositories := make([]string, 0, 13)
	for i := 0; i < 13; i++ {
		repositories = append(repositories, fmt.Sprintf("repo-%02d", i))
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == pullTrafficRepositoryPath:
			list := lib.RepositoryList{}
			for _, name := range repositories {
				list.Repositories = append(list.Repositories, lib.OrganizationRepository{Name: name})
			}
			writeJSON(t, w, list)
		case strings.HasSuffix(r.URL.Path, "/aggregatelogs"):
			name := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/repository/"+testNamespace+"/"), "/aggregatelogs")
			if name == errorRepo {
				w.WriteHeader(http.StatusInternalServerError)
				writeJSON(t, w, map[string]string{testPullTrafficErrorKey: "unreadable"})
				return
			}
			writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 100}}})
		case strings.HasSuffix(r.URL.Path, "/logs"):
			// repo-01 always reports another page, so the sample cap truncates it.
			next := ""
			if strings.Contains(r.URL.Path, "/repo-01/") {
				next = "next-page"
			}
			writeJSON(t, w, lib.Logs{Logs: []lib.LogEntry{{Kind: pullTrafficKindPull}}, NextPage: next})
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	}
	client := newPullTrafficTestClient(t, handler)

	report, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: pullTrafficTestStart,
		EndDate:   pullTrafficTestStart,
		ByTag:     true,
		MaxPages:  2,
	})
	if err != nil {
		t.Fatalf("collectPullTraffic returned error: %v", err)
	}

	statuses := map[string]reportRowStatus{}
	for _, row := range report.Repositories {
		statuses[row.Repository] = row.Status
	}
	want := map[string]reportRowStatus{
		"repo-00": rowStatusOK,
		"repo-01": rowStatusTruncated,
		"repo-11": rowStatusPartial,
		errorRepo: rowStatusError,
	}
	for name, wantStatus := range want {
		if got := statuses[name]; got != wantStatus {
			t.Errorf("%s status = %q, want %q", name, got, wantStatus)
		}
	}
}

func TestPrintPullTrafficTable(t *testing.T) {
	setCommandTestValue(t, &outputFormat, outputTable)
	report := &pullTrafficReport{
		Namespace:   testNamespace,
		StartDate:   pullTrafficTestStart,
		EndDate:     pullTrafficTestStart,
		TotalPulls:  300,
		TotalPushes: 4,
		Repositories: []pullTrafficRepo{
			{Repository: testPullTrafficBeta, Pulls: 300, Pushes: 4, SharePercent: 100},
		},
	}
	output := captureStdout(t, func() {
		if err := printPullTraffic(report, pullTrafficOptions{}); err != nil {
			t.Fatalf("printPullTraffic returned error: %v", err)
		}
	})
	for _, want := range []string{"REPOSITORY", "PULLS", testPullTrafficBeta, "100.0%"} {
		if !strings.Contains(output, want) {
			t.Errorf("table output %q does not contain %q", output, want)
		}
	}
}

func TestListPullTrafficRepositoriesPages(t *testing.T) {
	requests := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		requests++
		page := r.URL.Query().Get("page")
		list := lib.RepositoryList{}
		if page == "1" {
			list.Repositories = []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}
			list.HasAdditional = true
		} else {
			list.Repositories = []lib.OrganizationRepository{{Name: testPullTrafficBeta}}
		}
		writeJSON(t, w, list)
	}
	client := newPullTrafficTestClient(t, handler)

	names, err := listPullTrafficRepositories(context.Background(), client, testNamespace)
	if err != nil {
		t.Fatalf("listPullTrafficRepositories returned error: %v", err)
	}
	if fmt.Sprint(names) != "[alpha beta]" {
		t.Errorf("names = %v, want [alpha beta]", names)
	}
	if requests != 2 {
		t.Errorf("requests = %d, want 2 pages", requests)
	}
}

func TestCollectPullTrafficRequiresNamespace(t *testing.T) {
	client := newPullTrafficTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request %s %s without a namespace", r.Method, r.URL.Path)
	})

	_, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{StartDate: pullTrafficTestStart})
	if err == nil || !strings.Contains(err.Error(), "namespace is required") {
		t.Fatalf("error = %v, want a namespace error", err)
	}
}

func TestCollectPullTrafficTagErrorsRecorded(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == pullTrafficRepositoryPath:
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
		case strings.HasSuffix(r.URL.Path, "/aggregatelogs"):
			writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 4}}})
		default:
			w.WriteHeader(http.StatusForbidden)
			writeJSON(t, w, map[string]string{testPullTrafficErrorKey: "denied"})
		}
	}
	client := newPullTrafficTestClient(t, handler)

	report, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: pullTrafficTestStart,
		EndDate:   pullTrafficTestStart,
		ByTag:     true,
		MaxPages:  2,
	})
	if err != nil {
		t.Fatalf("collectPullTraffic returned error: %v", err)
	}
	row := report.Repositories[0]
	if row.Pulls != 4 {
		t.Errorf("pulls = %d, want the aggregate to survive a tag failure", row.Pulls)
	}
	if row.Error == "" {
		t.Error("tag sampling failure was not recorded on the row")
	}
}

func TestPrintPullTrafficJSON(t *testing.T) {
	setCommandTestValue(t, &outputFormat, outputJSON)
	report := &pullTrafficReport{
		Namespace:    testNamespace,
		Repositories: []pullTrafficRepo{{Repository: testPullTrafficAlpha, Pulls: 5, SharePercent: 100}},
	}
	output := captureStdout(t, func() {
		if err := printPullTraffic(report, pullTrafficOptions{}); err != nil {
			t.Fatalf("printPullTraffic returned error: %v", err)
		}
	})
	for _, want := range []string{`"repositories"`, testPullTrafficAlpha, `"share_percent": 100`} {
		if !strings.Contains(output, want) {
			t.Errorf("JSON output %q does not contain %q", output, want)
		}
	}
}

func TestListPullTrafficRepositoriesDetectsRejectedToken(t *testing.T) {
	client := newPullTrafficTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("public") == "true" {
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
			return
		}
		writeJSON(t, w, lib.RepositoryList{})
	})

	_, err := listPullTrafficRepositories(context.Background(), client, testNamespace)
	if err == nil || !strings.Contains(err.Error(), "token was not accepted") {
		t.Fatalf("error = %v, want an unauthenticated-token error", err)
	}
}

func TestListPullTrafficRepositoriesAllowsEmptyNamespace(t *testing.T) {
	client := newPullTrafficTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, lib.RepositoryList{})
	})

	names, err := listPullTrafficRepositories(context.Background(), client, testNamespace)
	if err != nil {
		t.Fatalf("empty namespace returned error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("names = %v, want none", names)
	}
}

func TestCollectPullTrafficReportsCoveredWindow(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pullTrafficRepositoryPath {
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
			return
		}
		writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{
			{Kind: pullTrafficKindPull, Count: 2, Datetime: "Tue, 06 Oct 2026 00:00:00 -0000"},
			{Kind: pullTrafficKindPull, Count: 3, Datetime: "Thu, 05 Nov 2026 00:00:00 -0000"},
		}})
	}
	client := newPullTrafficTestClient(t, handler)

	report, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: "2026-08-07",
		EndDate:   testPullTrafficEndDate,
	})
	if err != nil {
		t.Fatalf("collectPullTraffic returned error: %v", err)
	}
	if report.CoveredStart != testPullTrafficCoveredStart || report.CoveredEnd != testPullTrafficEndDate {
		t.Errorf("covered window = %s..%s, want 2026-10-06..2026-11-05", report.CoveredStart, report.CoveredEnd)
	}
}

func TestPrintPullTrafficTableShowsCoveredWindow(t *testing.T) {
	setCommandTestValue(t, &outputFormat, outputTable)
	report := &pullTrafficReport{
		Namespace:    testNamespace,
		StartDate:    "2026-08-07",
		EndDate:      testPullTrafficEndDate,
		CoveredStart: testPullTrafficCoveredStart,
		CoveredEnd:   testPullTrafficEndDate,
		TotalPulls:   5,
		Repositories: []pullTrafficRepo{{Repository: testPullTrafficAlpha, Pulls: 5, SharePercent: 100}},
	}
	stderr := captureStderr(t, func() {
		captureStdout(t, func() {
			if err := printPullTraffic(report, pullTrafficOptions{}); err != nil {
				t.Fatalf("printPullTraffic returned error: %v", err)
			}
		})
	})
	for _, want := range []string{"2026-10-06 to 2026-11-05", "requested 2026-08-07 to 2026-11-05"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("summary %q does not contain %q", stderr, want)
		}
	}
}

func TestDefaultPullTrafficWindow(t *testing.T) {
	now := time.Date(2026, time.November, 5, 12, 0, 0, 0, time.UTC)
	start, end := defaultReportWindow(now)
	if start != "10/07/2026" || end != "11/05/2026" {
		t.Errorf("window = %s..%s, want 10/07/2026..11/05/2026", start, end)
	}
}

func TestRankPullTrafficConsumers(t *testing.T) {
	consumers := rankPullTrafficConsumers([]lib.LogEntry{
		{Kind: pullTrafficKindPull, IP: testConsumerIP, Datetime: "Thu, 05 Nov 2026 19:00:00 -0000",
			Performer: lib.Performer{Name: testConsumerPuller}, Metadata: lib.Metadata{UserAgent: testConsumerUserAgent}},
		{Kind: pullTrafficKindPull, IP: testConsumerIP, Datetime: "Thu, 05 Nov 2026 20:00:00 -0000",
			Performer: lib.Performer{Name: testConsumerPuller}, Metadata: lib.Metadata{UserAgent: testConsumerUserAgent}},
		{Kind: pullTrafficKindPull, IP: testConsumerIPSecondary, Performer: lib.Performer{Name: "robot-b"}},
		{Kind: pullTrafficKindPull, IP: "3.3.3.3"},
		{Kind: pullTrafficKindPush, IP: testConsumerIP, Performer: lib.Performer{Name: testConsumerPuller}},
	})

	if len(consumers) != 3 {
		t.Fatalf("consumers = %+v, want 3 identities", consumers)
	}
	if consumers[0].Puller != testConsumerPuller || consumers[0].Pulls != 2 || consumers[0].IP != testConsumerIP {
		t.Errorf("top consumer = %+v, want robot-a from 1.1.1.1 with 2 pulls", consumers[0])
	}
	if consumers[0].LastSeen != testConsumerLastSeen {
		t.Errorf("last seen = %q, want 2026-11-05 20:00", consumers[0].LastSeen)
	}
	if consumers[0].UserAgent != testConsumerUserAgent {
		t.Errorf("user agent = %q, want cri-o/1.35", consumers[0].UserAgent)
	}
	if consumers[1].Puller != pullTrafficAnonymousPuller {
		t.Errorf("second consumer = %q, want the anonymous puller", consumers[1].Puller)
	}
}

func TestCollectPullTrafficConsumers(t *testing.T) {
	logRequests := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == pullTrafficRepositoryPath:
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
		case strings.HasSuffix(r.URL.Path, "/aggregatelogs"):
			writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 5}}})
		default:
			logRequests++
			writeJSON(t, w, lib.Logs{Logs: []lib.LogEntry{
				{Kind: pullTrafficKindPull, IP: testConsumerIP, Performer: lib.Performer{Name: testConsumerPuller}},
			}})
		}
	}
	client := newPullTrafficTestClient(t, handler)

	report, err := collectPullTraffic(context.Background(), client, pullTrafficOptions{
		Namespace: testNamespace,
		StartDate: pullTrafficTestStart,
		EndDate:   pullTrafficTestStart,
		ByTag:     true,
		Consumers: true,
		MaxPages:  2,
	})
	if err != nil {
		t.Fatalf("collectPullTraffic returned error: %v", err)
	}
	row := report.Repositories[0]
	if len(row.Consumers) != 1 || row.Consumers[0].Puller != testConsumerPuller {
		t.Errorf("consumers = %+v, want robot-a", row.Consumers)
	}
	if logRequests != 1 {
		t.Errorf("log requests = %d, want one shared sample for tags and consumers", logRequests)
	}
}

func TestPrintPullTrafficConsumersTable(t *testing.T) {
	setCommandTestValue(t, &outputFormat, outputTable)
	report := &pullTrafficReport{
		Namespace: testNamespace, StartDate: pullTrafficTestStart, EndDate: pullTrafficTestStart,
		Repositories: []pullTrafficRepo{{
			Repository: testPullTrafficAlpha,
			Pulls:      5,
			Consumers:  []pullTrafficConsumer{{Puller: testConsumerPuller, IP: testConsumerIP, Pulls: 5, LastSeen: testConsumerLastSeen}},
		}},
	}
	output := captureStdout(t, func() {
		if err := printPullTraffic(report, pullTrafficOptions{Consumers: true}); err != nil {
			t.Fatalf("printPullTraffic returned error: %v", err)
		}
	})
	for _, want := range []string{"PULLER", "IP", "LAST SEEN", testConsumerPuller, testConsumerIP} {
		if !strings.Contains(output, want) {
			t.Errorf("consumer table %q does not contain %q", output, want)
		}
	}
}

func TestPrintPullTrafficTagsTable(t *testing.T) {
	setCommandTestValue(t, &outputFormat, outputTable)
	report := &pullTrafficReport{
		Namespace: testNamespace, StartDate: pullTrafficTestStart, EndDate: pullTrafficTestStart,
		Repositories: []pullTrafficRepo{{
			Repository: testPullTrafficAlpha,
			Pulls:      5,
			TopTags:    []pullTrafficTag{{Tag: "v1", Count: 5}},
		}},
	}
	output := captureStdout(t, func() {
		if err := printPullTraffic(report, pullTrafficOptions{ByTag: true}); err != nil {
			t.Fatalf("printPullTraffic returned error: %v", err)
		}
	})
	for _, want := range []string{"TAGS", "v1 (5)"} {
		if !strings.Contains(output, want) {
			t.Errorf("tag table %q does not contain %q", output, want)
		}
	}
}

func TestPullTrafficCommandDefaultsToTableAndLast30Days(t *testing.T) {
	var startParam, endParam string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pullTrafficRepositoryPath {
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
			return
		}
		startParam = r.URL.Query().Get(testLogsStartTimeParam)
		endParam = r.URL.Query().Get(testLogsEndTimeParam)
		writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 7}}})
	}))
	defer server.Close()

	resetRootFlags(t)
	resetCommandFlagChanges(rootCmd)
	rootCmd.SetArgs([]string{
		"list", "pull-traffic",
		testConfigNamespaceFlag, testNamespace,
		testTokenFlag, testTokenValue,
		testQuayURLFlag, server.URL + "/api/v1",
	})
	rootCmd.SetErr(io.Discard)
	var runErr error
	output := captureStdout(t, func() { runErr = rootCmd.Execute() })
	if runErr != nil {
		t.Fatalf("command returned error: %v", runErr)
	}
	for _, want := range []string{"REPOSITORY", "PULLS", testPullTrafficAlpha} {
		if !strings.Contains(output, want) {
			t.Errorf("output %q does not contain %q; the command should default to a table", output, want)
		}
	}

	start, err := time.Parse(reportAPIDateLayout, startParam)
	if err != nil {
		t.Fatalf("starttime %q is not MM/DD/YYYY: %v", startParam, err)
	}
	end, err := time.Parse(reportAPIDateLayout, endParam)
	if err != nil {
		t.Fatalf("endtime %q is not MM/DD/YYYY: %v", endParam, err)
	}
	if days := end.Sub(start).Hours()/24 + 1; days != pullTrafficRetentionDays {
		t.Errorf("default window = %s..%s (%.0f days), want %d", startParam, endParam, days, pullTrafficRetentionDays)
	}
}

func TestRankPullTrafficConsumersCapsResults(t *testing.T) {
	entries := make([]lib.LogEntry, 0, pullTrafficTopConsumers+1)
	for i := 0; i <= pullTrafficTopConsumers; i++ {
		entries = append(entries, lib.LogEntry{
			Kind:      pullTrafficKindPull,
			Performer: lib.Performer{Name: fmt.Sprintf("robot-%d", i)},
		})
	}
	if got := len(rankPullTrafficConsumers(entries)); got != pullTrafficTopConsumers {
		t.Errorf("consumers = %d, want the list capped at %d", got, pullTrafficTopConsumers)
	}
}
