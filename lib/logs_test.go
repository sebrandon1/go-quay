package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	// The log API only honors MM/DD/YYYY; requests normalize to this form.
	testStartDate    = "05/01/2026"
	testEndDate      = "05/21/2026"
	testEndDateMay   = "05/31/2026"
	testDatetime     = "Mon, 19 May 2026 00:00:00 -0000"
	testNextPage     = "abc123"
	testKindPullRepo = "pull_repo"
)

func TestGetAggregatedLogs(t *testing.T) {
	mockResponse := AggregatedLogs{
		Aggregated: []AggregatedLogEntry{
			{Kind: testKindPullRepo, Count: 10, Datetime: testDatetime},
			{Kind: "push_repo", Count: 2, Datetime: testDatetime},
		},
	}
	mockResponseJSON := mustMarshal(t, mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/testorg/testrepo/aggregatelogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.URL.Query().Get("starttime") != testStartDate {
			t.Errorf("Expected starttime '%s', got '%s'", testStartDate, r.URL.Query().Get("starttime"))
		}
		if r.URL.Query().Get("endtime") != testEndDate {
			t.Errorf("Expected endtime '%s', got '%s'", testEndDate, r.URL.Query().Get("endtime"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	logs, err := client.GetAggregatedLogs(context.Background(), testNamespace, testRepository, testStartDate, testEndDate)
	if err != nil {
		t.Fatalf("GetAggregatedLogs returned error: %v", err)
	}

	if len(logs.Aggregated) != 2 {
		t.Errorf("Expected 2 aggregated entries, got %d", len(logs.Aggregated))
	}
	if logs.Aggregated[0].Kind != testKindPullRepo {
		t.Errorf("Expected kind 'pull_repo', got '%s'", logs.Aggregated[0].Kind)
	}
	if logs.Aggregated[0].Count != 10 {
		t.Errorf("Expected count 10, got %d", logs.Aggregated[0].Count)
	}
}

func TestGetLogs(t *testing.T) {
	mockResponse := Logs{
		StartTime: testDatetime,
		EndTime:   "Wed, 21 May 2026 00:00:00 -0000",
		Logs: []LogEntry{
			{
				Kind:     testKindPullRepo,
				Datetime: "Wed, 21 May 2026 10:57:43 -0000",
				Metadata: Metadata{Repo: testRepository, Tag: testTagNameLatest},
			},
		},
	}
	mockResponseJSON := mustMarshal(t, mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/testorg/testrepo/logs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	logs, err := client.GetLogs(context.Background(), testNamespace, testRepository, "", "", "")
	if err != nil {
		t.Fatalf("GetLogs returned error: %v", err)
	}

	if len(logs.Logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs.Logs))
	}
	if logs.Logs[0].Kind != testKindPullRepo {
		t.Errorf("Expected kind 'pull_repo', got '%s'", logs.Logs[0].Kind)
	}
	if logs.Logs[0].Metadata.Tag != testTagNameLatest {
		t.Errorf("Expected tag 'latest', got '%s'", logs.Logs[0].Metadata.Tag)
	}
}

func TestGetLogsWithDateRange(t *testing.T) {
	mockResponseJSON := mustMarshal(t, Logs{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("starttime") != testStartDate {
			t.Errorf("Expected starttime '%s', got '%s'", testStartDate, r.URL.Query().Get("starttime"))
		}
		if r.URL.Query().Get("endtime") != testEndDate {
			t.Errorf("Expected endtime '%s', got '%s'", testEndDate, r.URL.Query().Get("endtime"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetLogs(context.Background(), testNamespace, testRepository, "", testStartDate, testEndDate)
	if err != nil {
		t.Fatalf("GetLogs with date range returned error: %v", err)
	}
}

func TestGetLogsWithNextPage(t *testing.T) {
	mockResponseJSON := mustMarshal(t, Logs{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("next_page") != testNextPage {
			t.Errorf("Expected next_page 'abc123', got '%s'", r.URL.Query().Get("next_page"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetLogs(context.Background(), testNamespace, testRepository, testNextPage, "", "")
	if err != nil {
		t.Fatalf("GetLogs with next page returned error: %v", err)
	}
}

func TestGetOrganizationLogsWithDateRange(t *testing.T) {
	mockResponseJSON := mustMarshal(t, Logs{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/organization/testorg/logs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.URL.Query().Get("starttime") != testStartDate {
			t.Errorf("Expected starttime '%s', got '%s'", testStartDate, r.URL.Query().Get("starttime"))
		}
		if r.URL.Query().Get("endtime") != testEndDateMay {
			t.Errorf("Expected endtime '%s', got '%s'", testEndDateMay, r.URL.Query().Get("endtime"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetOrganizationLogs(context.Background(), testNamespace, "", testStartDate, testEndDateMay)
	if err != nil {
		t.Fatalf("GetOrganizationLogs with date range returned error: %v", err)
	}
}

func TestGetUserLogsWithDateRange(t *testing.T) {
	mockResponseJSON := mustMarshal(t, Logs{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/user/logs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.URL.Query().Get("starttime") != testStartDate {
			t.Errorf("Expected starttime '%s', got '%s'", testStartDate, r.URL.Query().Get("starttime"))
		}
		if r.URL.Query().Get("endtime") != testEndDateMay {
			t.Errorf("Expected endtime '%s', got '%s'", testEndDateMay, r.URL.Query().Get("endtime"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetUserLogs(context.Background(), "", testStartDate, testEndDateMay)
	if err != nil {
		t.Fatalf("GetUserLogs with date range returned error: %v", err)
	}
}

func TestGetOrganizationAggregatedLogs(t *testing.T) {
	mockResponse := AggregatedLogs{
		Aggregated: []AggregatedLogEntry{
			{Kind: testKindPullRepo, Count: 25, Datetime: testDatetime},
		},
	}
	mockResponseJSON := mustMarshal(t, mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/testorg/aggregatelogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.URL.Query().Get("starttime") != testStartDate {
			t.Errorf("Expected starttime '%s', got '%s'", testStartDate, r.URL.Query().Get("starttime"))
		}
		if r.URL.Query().Get("endtime") != testEndDateMay {
			t.Errorf("Expected endtime '%s', got '%s'", testEndDateMay, r.URL.Query().Get("endtime"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	logs, err := client.GetOrganizationAggregatedLogs(context.Background(), testNamespace, testStartDate, testEndDateMay)
	if err != nil {
		t.Fatalf("GetOrganizationAggregatedLogs returned error: %v", err)
	}

	if len(logs.Aggregated) != 1 {
		t.Errorf("Expected 1 aggregated entry, got %d", len(logs.Aggregated))
	}
	if logs.Aggregated[0].Count != 25 {
		t.Errorf("Expected count 25, got %d", logs.Aggregated[0].Count)
	}
}

func TestExportOrganizationLogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/testorg/exportlogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	exportReq := &ExportLogsRequest{
		StartTime:   testStartDate,
		EndTime:     testEndDateMay,
		CallbackURL: "https://example.com/callback",
	}

	err = client.ExportOrganizationLogs(context.Background(), testNamespace, exportReq)
	if err != nil {
		t.Fatalf("ExportOrganizationLogs returned error: %v", err)
	}
}

func TestGetUserAggregatedLogs(t *testing.T) {
	mockResponse := AggregatedLogs{
		Aggregated: []AggregatedLogEntry{
			{Kind: testKindPullRepo, Count: 15, Datetime: testDatetime},
		},
	}
	mockResponseJSON := mustMarshal(t, mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/user/aggregatelogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.URL.Query().Get("starttime") != testStartDate {
			t.Errorf("Expected starttime '%s', got '%s'", testStartDate, r.URL.Query().Get("starttime"))
		}
		if r.URL.Query().Get("endtime") != testEndDateMay {
			t.Errorf("Expected endtime '%s', got '%s'", testEndDateMay, r.URL.Query().Get("endtime"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	logs, err := client.GetUserAggregatedLogs(context.Background(), testStartDate, testEndDateMay)
	if err != nil {
		t.Fatalf("GetUserAggregatedLogs returned error: %v", err)
	}

	if len(logs.Aggregated) != 1 {
		t.Errorf("Expected 1 aggregated entry, got %d", len(logs.Aggregated))
	}
	if logs.Aggregated[0].Count != 15 {
		t.Errorf("Expected count 15, got %d", logs.Aggregated[0].Count)
	}
}

func TestExportUserLogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/user/exportlogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	exportReq := &ExportLogsRequest{
		StartTime: testStartDate,
		EndTime:   testEndDateMay,
		Email:     testEmailAddress,
	}

	err = client.ExportUserLogs(context.Background(), exportReq)
	if err != nil {
		t.Fatalf("ExportUserLogs returned error: %v", err)
	}
}

func TestExportRepositoryLogs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/testorg/testrepo/exportlogs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	exportReq := &ExportLogsRequest{
		StartTime:   testStartDate,
		EndTime:     testEndDateMay,
		CallbackURL: "https://example.com/callback",
	}

	err = client.ExportRepositoryLogs(context.Background(), testNamespace, testRepository, exportReq)
	if err != nil {
		t.Fatalf("ExportRepositoryLogs returned error: %v", err)
	}
}

func TestGetAggregatedLogsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetAggregatedLogs(context.Background(), testNamespace, testRepository, "", "")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetLogsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetLogs(context.Background(), testNamespace, testRepository, "", "", "")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestLogsHTTPErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetOrganizationLogs(context.Background(), testNamespace, "", "", "")
	if err == nil {
		t.Error("Expected error from GetOrganizationLogs, got nil")
	}

	_, err = client.GetOrganizationAggregatedLogs(context.Background(), testNamespace, testStartDate, testEndDateMay)
	if err == nil {
		t.Error("Expected error from GetOrganizationAggregatedLogs, got nil")
	}

	err = client.ExportOrganizationLogs(context.Background(), testNamespace, &ExportLogsRequest{StartTime: testStartDate, EndTime: testEndDateMay})
	if err == nil {
		t.Error("Expected error from ExportOrganizationLogs, got nil")
	}

	_, err = client.GetUserLogs(context.Background(), "", testStartDate, testEndDateMay)
	if err == nil {
		t.Error("Expected error from GetUserLogs, got nil")
	}

	_, err = client.GetUserAggregatedLogs(context.Background(), testStartDate, testEndDateMay)
	if err == nil {
		t.Error("Expected error from GetUserAggregatedLogs, got nil")
	}

	err = client.ExportUserLogs(context.Background(), &ExportLogsRequest{StartTime: testStartDate, EndTime: testEndDateMay})
	if err == nil {
		t.Error("Expected error from ExportUserLogs, got nil")
	}

	err = client.ExportRepositoryLogs(context.Background(), testNamespace, testRepository, &ExportLogsRequest{StartTime: testStartDate, EndTime: testEndDateMay})
	if err == nil {
		t.Error("Expected error from ExportRepositoryLogs, got nil")
	}
}

func TestOrganizationRepositoryPopularityZero(t *testing.T) {
	repo := OrganizationRepository{
		Name:       testRepository,
		Namespace:  testNamespace,
		Popularity: 0,
	}

	data, err := json.Marshal(repo)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, `"popularity":0`) && !strings.Contains(output, `"popularity": 0`) {
		t.Errorf("Expected popularity field in JSON output even when zero, got: %s", output)
	}
}

func TestNormalizeLogDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "api layout", input: testStartDate, want: testStartDate},
		{name: "iso layout", input: "2026-05-01", want: testStartDate},
		{name: "rfc3339", input: "2026-05-01T12:30:00Z", want: testStartDate},
		{name: "empty", input: "", want: ""},
		{name: "unparseable", input: "May 1, 2026", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := normalizeLogDate(test.input)
			if test.wantErr {
				if err == nil || !strings.Contains(err.Error(), "invalid log date") {
					t.Fatalf("normalizeLogDate(%q) error = %v, want an invalid-date error", test.input, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeLogDate(%q) returned error: %v", test.input, err)
			}
			if got != test.want {
				t.Errorf("normalizeLogDate(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestLogRequestsNormalizeISODates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("starttime"); got != testStartDate {
			t.Errorf("starttime = %q, want %q", got, testStartDate)
		}
		if got := r.URL.Query().Get("endtime"); got != testEndDate {
			t.Errorf("endtime = %q, want %q", got, testEndDate)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mustMarshal(t, AggregatedLogs{}))
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if _, err := client.GetAggregatedLogs(context.Background(), testNamespace, testRepository, "2026-05-01", "2026-05-21"); err != nil {
		t.Fatalf("GetAggregatedLogs with ISO dates returned error: %v", err)
	}
}

func TestAggregatedLogsOmitUnsetDates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.URL.Query()["starttime"]; ok {
			t.Errorf("starttime should be omitted when empty, got %q", r.URL.Query().Get("starttime"))
		}
		if _, ok := r.URL.Query()["endtime"]; ok {
			t.Errorf("endtime should be omitted when empty, got %q", r.URL.Query().Get("endtime"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mustMarshal(t, AggregatedLogs{}))
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if _, err := client.GetAggregatedLogs(context.Background(), testNamespace, testRepository, "", ""); err != nil {
		t.Fatalf("GetAggregatedLogs without dates returned error: %v", err)
	}
}

func TestLogsRejectInvalidDate(t *testing.T) {
	client, err := NewClient("")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if _, err := client.GetAggregatedLogs(context.Background(), testNamespace, testRepository, "not-a-date", ""); err == nil {
		t.Fatal("GetAggregatedLogs accepted an invalid date")
	}
	if _, err := client.GetLogs(context.Background(), testNamespace, testRepository, "", "", "not-a-date"); err == nil {
		t.Fatal("GetLogs accepted an invalid date")
	}
	if _, _, err := ValidateLogDateRange("2026-05-01", "2026-05-21"); err != nil {
		t.Errorf("ValidateLogDateRange rejected valid dates: %v", err)
	}
}

func TestGetAggregatedLogsRangeSplitsWindows(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		start, err := time.Parse("01/02/2006", r.URL.Query().Get("starttime"))
		if err != nil {
			t.Errorf("bad starttime %q: %v", r.URL.Query().Get("starttime"), err)
		}
		end, err := time.Parse("01/02/2006", r.URL.Query().Get("endtime"))
		if err != nil {
			t.Errorf("bad endtime %q: %v", r.URL.Query().Get("endtime"), err)
		}
		if days := end.Sub(start).Hours()/24 + 1; days > maxAggregatedLogWindowDays {
			t.Errorf("window %s..%s spans %.0f days, want <= %d", start, end, days, maxAggregatedLogWindowDays)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mustMarshal(t, AggregatedLogs{Aggregated: []AggregatedLogEntry{
			{Kind: testKindPullRepo, Count: 1, Datetime: testDatetime},
		}}))
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	logs, err := client.GetAggregatedLogsRange(context.Background(), testNamespace, testRepository, "2026-01-01", "2026-03-31")
	if err != nil {
		t.Fatalf("GetAggregatedLogsRange returned error: %v", err)
	}
	if requests != 3 {
		t.Errorf("requests = %d, want 3 windows for a 90-day range", requests)
	}
	if len(logs.Aggregated) != 3 {
		t.Errorf("entries = %d, want 3 merged windows", len(logs.Aggregated))
	}
}

func TestGetAggregatedLogsRangeRejectsReversedRange(t *testing.T) {
	client, err := NewClient("")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	_, err = client.GetAggregatedLogsRange(context.Background(), testNamespace, testRepository, "05/21/2026", "05/01/2026")
	if err == nil || !strings.Contains(err.Error(), "before start date") {
		t.Fatalf("error = %v, want a reversed-range error", err)
	}
}

func TestGetAllLogsFollowsNextPage(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		response := Logs{Logs: []LogEntry{{Kind: testKindPullRepo, Metadata: Metadata{Tag: "first"}}}}
		if r.URL.Query().Get("next_page") == "" {
			response.NextPage = testNextPage
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mustMarshal(t, response))
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	entries, truncated, err := client.GetAllLogs(context.Background(), testNamespace, testRepository, "", "", 5)
	if err != nil {
		t.Fatalf("GetAllLogs returned error: %v", err)
	}
	if requests != 2 || len(entries) != 2 {
		t.Errorf("requests = %d, entries = %d; want 2 and 2", requests, len(entries))
	}
	if truncated {
		t.Error("truncated = true, want false when the stream ends")
	}
}

func TestGetAllLogsStopsAtPageLimit(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
		w.Write(mustMarshal(t, Logs{
			Logs:     []LogEntry{{Kind: testKindPullRepo}},
			NextPage: testNextPage,
		}))
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	entries, truncated, err := client.GetAllLogs(context.Background(), testNamespace, testRepository, "", "", 2)
	if err != nil {
		t.Fatalf("GetAllLogs returned error: %v", err)
	}
	if requests != 2 || len(entries) != 2 {
		t.Errorf("requests = %d, entries = %d; want 2 and 2", requests, len(entries))
	}
	if !truncated {
		t.Error("truncated = false, want true when the page cap is hit")
	}
}

func TestGetAggregatedLogsRangeSingleRequestWhenOpenEnded(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
		w.Write(mustMarshal(t, AggregatedLogs{}))
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if _, err := client.GetAggregatedLogsRange(context.Background(), testNamespace, testRepository, "", ""); err != nil {
		t.Fatalf("GetAggregatedLogsRange returned error: %v", err)
	}
	if requests != 1 {
		t.Errorf("requests = %d, want 1 when the range is open-ended", requests)
	}
}

func TestGetAggregatedLogsRangeStopsOnWindowError(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Write(mustMarshal(t, map[string]string{"error": "denied"}))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mustMarshal(t, AggregatedLogs{}))
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if _, err := client.GetAggregatedLogsRange(context.Background(), testNamespace, testRepository, "2026-01-01", "2026-03-31"); err == nil {
		t.Fatal("GetAggregatedLogsRange swallowed a window error")
	}
	if requests != 2 {
		t.Errorf("requests = %d, want the range query to stop at the failing window", requests)
	}
}
