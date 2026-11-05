package cmd

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sebrandon1/go-quay/lib"
)

func TestResolvePullTrafficDate(t *testing.T) {
	now := time.Date(2026, time.November, 5, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "explicit date", input: pullTrafficTestStart, want: pullTrafficTestStart},
		{name: "relative days", input: "30d", want: pullTrafficTestStart},
		{name: "relative today", input: "0d", want: "11/05/2026"},
		{name: "empty", input: "", want: ""},
		{name: "invalid relative", input: "xd", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolveReportDate(test.input, now)
			if test.wantErr {
				if err == nil {
					t.Fatalf("resolveReportDate(%q) returned no error", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveReportDate(%q) returned error: %v", test.input, err)
			}
			if got != test.want {
				t.Errorf("resolveReportDate(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestPrintPullTrafficCSV(t *testing.T) {
	setCommandTestValue(t, &outputFormat, outputCSV)
	report := &pullTrafficReport{
		Namespace:    testNamespace,
		StartDate:    pullTrafficTestStart,
		EndDate:      pullTrafficTestStart,
		TotalPulls:   10,
		Repositories: []pullTrafficRepo{{Repository: testPullTrafficAlpha, Pulls: 10, SharePercent: 100}},
	}
	output := captureStdout(t, func() {
		if err := printPullTraffic(report, pullTrafficOptions{}); err != nil {
			t.Fatalf("printPullTraffic returned error: %v", err)
		}
	})
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		t.Fatalf("csv lines = %d (%q), want a header and one row", len(lines), output)
	}
	if lines[0] != "REPOSITORY,PULLS,PUSHES,SHARE" {
		t.Errorf("csv header = %q", lines[0])
	}
	if lines[1] != testPullTrafficAlpha+",10,0,100.0%" {
		t.Errorf("csv row = %q", lines[1])
	}
}

func TestCollectReportConsumersMergesRepositories(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == pullTrafficRepositoryPath:
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{
				{Name: testPullTrafficAlpha}, {Name: testPullTrafficBeta},
			}})
		case strings.HasSuffix(r.URL.Path, "/aggregatelogs"):
			writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 5}}})
		case strings.Contains(r.URL.Path, "/"+testPullTrafficAlpha+"/logs"):
			writeJSON(t, w, lib.Logs{Logs: []lib.LogEntry{
				{Kind: pullTrafficKindPull, IP: testConsumerIP, Performer: lib.Performer{Name: testConsumerPuller}},
				{Kind: pullTrafficKindPull, IP: testConsumerIP, Performer: lib.Performer{Name: testConsumerPuller}},
				{Kind: pullTrafficKindPull, IP: "2.2.2.2", Performer: lib.Performer{Name: "robot-y"}},
			}})
		default:
			writeJSON(t, w, lib.Logs{Logs: []lib.LogEntry{
				{Kind: pullTrafficKindPull, IP: testConsumerIP, Performer: lib.Performer{Name: testConsumerPuller}},
			}})
		}
	}
	client := newPullTrafficTestClient(t, handler)

	report, err := collectReportConsumers(context.Background(), client, testNamespace,
		pullTrafficTestStart, pullTrafficTestStart, 10, 5, 10)
	if err != nil {
		t.Fatalf("collectReportConsumers returned error: %v", err)
	}
	if report.RepositoriesSampled != 2 {
		t.Errorf("repositories sampled = %d, want 2", report.RepositoriesSampled)
	}
	if len(report.Consumers) != 2 {
		t.Fatalf("consumers = %+v, want 2 identities", report.Consumers)
	}
	top := report.Consumers[0]
	if top.Puller != testConsumerPuller || top.Pulls != 3 || top.Repositories != 2 {
		t.Errorf("top consumer = %+v, want robot-a with 3 pulls across 2 repositories", top)
	}
	if report.Consumers[1].Repositories != 1 {
		t.Errorf("second consumer = %+v, want a single repository", report.Consumers[1])
	}
}

func TestCollectReportConsumersLimit(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == pullTrafficRepositoryPath:
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
		case strings.HasSuffix(r.URL.Path, "/aggregatelogs"):
			writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 3}}})
		default:
			writeJSON(t, w, lib.Logs{Logs: []lib.LogEntry{
				{Kind: pullTrafficKindPull, IP: testConsumerIP, Performer: lib.Performer{Name: testConsumerPuller}},
				{Kind: pullTrafficKindPull, IP: "2.2.2.2", Performer: lib.Performer{Name: "robot-y"}},
			}})
		}
	}
	client := newPullTrafficTestClient(t, handler)

	report, err := collectReportConsumers(context.Background(), client, testNamespace,
		pullTrafficTestStart, pullTrafficTestStart, 1, 5, 1)
	if err != nil {
		t.Fatalf("collectReportConsumers returned error: %v", err)
	}
	if len(report.Consumers) != 1 {
		t.Errorf("consumers = %d, want the list capped by --limit", len(report.Consumers))
	}
}

func TestPrintReportConsumersCSV(t *testing.T) {
	setCommandTestValue(t, &outputFormat, outputCSV)
	report := &reportConsumersReport{
		Namespace: testNamespace,
		Consumers: []reportConsumer{{Puller: testConsumerPuller, IP: testConsumerIP, Pulls: 7, Repositories: 2, LastSeen: "2026-11-05 20:00"}},
	}
	output := captureStdout(t, func() {
		if err := printReportConsumers(report); err != nil {
			t.Fatalf("printReportConsumers returned error: %v", err)
		}
	})
	want := "PULLER,IP,PULLS,REPOS,LAST SEEN\n" + testConsumerPuller + "," + testConsumerIP + ",7,2,2026-11-05 20:00\n"
	if output != want {
		t.Errorf("csv output = %q, want %q", output, want)
	}
}

func TestReportConsumersCommandDefaultsToTable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == pullTrafficRepositoryPath:
			writeJSON(t, w, lib.RepositoryList{Repositories: []lib.OrganizationRepository{{Name: testPullTrafficAlpha}}})
		case strings.HasSuffix(r.URL.Path, "/aggregatelogs"):
			writeJSON(t, w, lib.AggregatedLogs{Aggregated: []lib.AggregatedLogEntry{{Kind: pullTrafficKindPull, Count: 3}}})
		default:
			writeJSON(t, w, lib.Logs{Logs: []lib.LogEntry{
				{Kind: pullTrafficKindPull, IP: testConsumerIP, Performer: lib.Performer{Name: testConsumerPuller}},
			}})
		}
	}))
	defer server.Close()

	resetRootFlags(t)
	resetCommandFlagChanges(rootCmd)
	rootCmd.SetArgs([]string{
		"report", "consumers",
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
	for _, want := range []string{"PULLER", "IP", "REPOS", testConsumerPuller} {
		if !strings.Contains(output, want) {
			t.Errorf("output %q does not contain %q; the command should default to a table", output, want)
		}
	}
}

func TestCollectReportConsumersRequiresNamespace(t *testing.T) {
	client := newPullTrafficTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request %s %s without a namespace", r.Method, r.URL.Path)
	})
	if _, err := collectReportConsumers(context.Background(), client, "", pullTrafficTestStart, pullTrafficTestStart, 1, 1, 1); err == nil {
		t.Fatal("collectReportConsumers accepted an empty namespace")
	}
}

func TestPrintReportConsumersJSON(t *testing.T) {
	setCommandTestValue(t, &outputFormat, outputJSON)
	report := &reportConsumersReport{
		Namespace: testNamespace,
		Consumers: []reportConsumer{{Puller: testConsumerPuller, Pulls: 2, Repositories: 1}},
	}
	output := captureStdout(t, func() {
		if err := printReportConsumers(report); err != nil {
			t.Fatalf("printReportConsumers returned error: %v", err)
		}
	})
	for _, want := range []string{`"consumers"`, testConsumerPuller, `"repositories": 1`} {
		if !strings.Contains(output, want) {
			t.Errorf("json output %q does not contain %q", output, want)
		}
	}
}
