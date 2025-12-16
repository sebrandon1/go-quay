package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetBuild    = "GET"
	httpPostBuild   = "POST"
	httpDeleteBuild = "DELETE"

	testBuildNamespace  = "testorg"
	testBuildRepository = "testrepo"
	testBuildUUID       = "build-uuid-123"
	testBuildPhase      = "complete"
)

func TestGetBuilds(t *testing.T) {
	mockResponse := Builds{
		Builds: []Build{
			{ID: testBuildUUID, Phase: testBuildPhase, DisplayName: "Build 1"},
			{ID: "build-uuid-456", Phase: "building", DisplayName: "Build 2"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetBuild {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testBuildNamespace + "/" + testBuildRepository + "/build/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	builds, err := client.GetBuilds(testBuildNamespace, testBuildRepository, 0)
	if err != nil {
		t.Fatalf("GetBuilds returned error: %v", err)
	}

	if len(builds.Builds) != 2 {
		t.Errorf("Expected 2 builds, got %d", len(builds.Builds))
	}
	if builds.Builds[0].ID != testBuildUUID {
		t.Errorf("Expected first build ID %s, got %s", testBuildUUID, builds.Builds[0].ID)
	}
}

func TestGetBuild(t *testing.T) {
	mockResponse := Build{
		ID:          testBuildUUID,
		Phase:       testBuildPhase,
		DisplayName: "Test Build",
		Started:     "2024-01-15T10:30:00Z",
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetBuild {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testBuildNamespace + "/" + testBuildRepository + "/build/" + testBuildUUID
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	build, err := client.GetBuild(testBuildNamespace, testBuildRepository, testBuildUUID)
	if err != nil {
		t.Fatalf("GetBuild returned error: %v", err)
	}

	if build.ID != testBuildUUID {
		t.Errorf("Expected build ID %s, got %s", testBuildUUID, build.ID)
	}
	if build.Phase != testBuildPhase {
		t.Errorf("Expected phase %s, got %s", testBuildPhase, build.Phase)
	}
}

func TestGetBuildLogs(t *testing.T) {
	mockResponse := BuildLogs{
		Start: 0,
		Total: 2,
		Logs: []BuildLogEntry{
			{Type: "phase", Message: "Starting build"},
			{Type: "command", Message: "docker build ."},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetBuild {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testBuildNamespace + "/" + testBuildRepository + "/build/" + testBuildUUID + "/logs"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	logs, err := client.GetBuildLogs(testBuildNamespace, testBuildRepository, testBuildUUID)
	if err != nil {
		t.Fatalf("GetBuildLogs returned error: %v", err)
	}

	if logs.Total != 2 {
		t.Errorf("Expected 2 log entries, got %d", logs.Total)
	}
	if len(logs.Logs) != 2 {
		t.Errorf("Expected 2 log entries in array, got %d", len(logs.Logs))
	}
}

func TestRequestBuild(t *testing.T) {
	mockResponse := Build{
		ID:    testBuildUUID,
		Phase: "waiting",
		Tags:  []string{"latest"},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPostBuild {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testBuildNamespace + "/" + testBuildRepository + "/build/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	buildReq := &RequestBuildRequest{
		ArchiveURL: "https://example.com/archive.tar.gz",
		Tags:       []string{"latest"},
	}

	build, err := client.RequestBuild(testBuildNamespace, testBuildRepository, buildReq)
	if err != nil {
		t.Fatalf("RequestBuild returned error: %v", err)
	}

	if build.ID != testBuildUUID {
		t.Errorf("Expected build ID %s, got %s", testBuildUUID, build.ID)
	}
}

func TestCancelBuild(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteBuild {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testBuildNamespace + "/" + testBuildRepository + "/build/" + testBuildUUID
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.CancelBuild(testBuildNamespace, testBuildRepository, testBuildUUID)
	if err != nil {
		t.Fatalf("CancelBuild returned error: %v", err)
	}
}

func TestGetBuildsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetBuilds(testBuildNamespace, testBuildRepository, 0)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetBuildError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetBuild(testBuildNamespace, testBuildRepository, "nonexistent-uuid")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
