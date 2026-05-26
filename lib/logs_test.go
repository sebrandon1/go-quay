package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testStartDate    = "2026-05-01"
	testEndDate      = "2026-05-21"
	testEndDateMay   = "2026-05-31"
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
	mockResponseJSON, _ := json.Marshal(mockResponse)

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

	logs, err := client.GetAggregatedLogs(testNamespace, testRepository, testStartDate, testEndDate)
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
	mockResponseJSON, _ := json.Marshal(mockResponse)

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

	logs, err := client.GetLogs(testNamespace, testRepository, "", "", "")
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
	mockResponseJSON, _ := json.Marshal(Logs{})

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

	_, err = client.GetLogs(testNamespace, testRepository, "", testStartDate, testEndDate)
	if err != nil {
		t.Fatalf("GetLogs with date range returned error: %v", err)
	}
}

func TestGetLogsWithNextPage(t *testing.T) {
	mockResponseJSON, _ := json.Marshal(Logs{})

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

	_, err = client.GetLogs(testNamespace, testRepository, testNextPage, "", "")
	if err != nil {
		t.Fatalf("GetLogs with next page returned error: %v", err)
	}
}

func TestGetOrganizationLogsWithDateRange(t *testing.T) {
	mockResponseJSON, _ := json.Marshal(Logs{})

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

	_, err = client.GetOrganizationLogs(testNamespace, "", testStartDate, testEndDateMay)
	if err != nil {
		t.Fatalf("GetOrganizationLogs with date range returned error: %v", err)
	}
}

func TestGetUserLogsWithDateRange(t *testing.T) {
	mockResponseJSON, _ := json.Marshal(Logs{})

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

	_, err = client.GetUserLogs("", testStartDate, testEndDateMay)
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
	mockResponseJSON, _ := json.Marshal(mockResponse)

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

	logs, err := client.GetOrganizationAggregatedLogs(testNamespace, testStartDate, testEndDateMay)
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

	err = client.ExportOrganizationLogs(testNamespace, exportReq)
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
	mockResponseJSON, _ := json.Marshal(mockResponse)

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

	logs, err := client.GetUserAggregatedLogs(testStartDate, testEndDateMay)
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

	err = client.ExportUserLogs(exportReq)
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

	err = client.ExportRepositoryLogs(testNamespace, testRepository, exportReq)
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

	_, err = client.GetAggregatedLogs(testNamespace, testRepository, "", "")
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

	_, err = client.GetLogs(testNamespace, testRepository, "", "", "")
	if err == nil {
		t.Error("Expected error, got nil")
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
