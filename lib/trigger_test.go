package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testTriggerUUID    = "trigger-uuid-123"
	testTriggerService = "github"
)

func TestGetTriggers(t *testing.T) {
	mockResponse := BuildTriggers{
		Triggers: []BuildTrigger{
			{ID: testTriggerUUID, Service: testTriggerService, IsActive: true, BuildSource: "org/repo"},
			{ID: "trigger-uuid-456", Service: "gitlab", IsActive: false, BuildSource: "group/project"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/trigger/"
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

	triggers, err := client.GetTriggers(testNamespace, testRepository)
	if err != nil {
		t.Fatalf("GetTriggers returned error: %v", err)
	}

	if len(triggers.Triggers) != 2 {
		t.Errorf("Expected 2 triggers, got %d", len(triggers.Triggers))
	}
	if triggers.Triggers[0].ID != testTriggerUUID {
		t.Errorf("Expected first trigger ID %s, got %s", testTriggerUUID, triggers.Triggers[0].ID)
	}
	if triggers.Triggers[0].Service != testTriggerService {
		t.Errorf("Expected service %s, got %s", testTriggerService, triggers.Triggers[0].Service)
	}
}

func TestGetTrigger(t *testing.T) {
	mockResponse := BuildTrigger{
		ID:            testTriggerUUID,
		Service:       testTriggerService,
		IsActive:      true,
		BuildSource:   "myorg/myrepo",
		RepositoryURL: "https://github.com/myorg/myrepo",
		CanInvoke:     true,
		Enabled:       true,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/trigger/" + testTriggerUUID
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

	trigger, err := client.GetTrigger(testNamespace, testRepository, testTriggerUUID)
	if err != nil {
		t.Fatalf("GetTrigger returned error: %v", err)
	}

	if trigger.ID != testTriggerUUID {
		t.Errorf("Expected trigger ID %s, got %s", testTriggerUUID, trigger.ID)
	}
	if trigger.Service != testTriggerService {
		t.Errorf("Expected service %s, got %s", testTriggerService, trigger.Service)
	}
	if !trigger.IsActive {
		t.Error("Expected trigger to be active")
	}
}

func TestDeleteTrigger(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/trigger/" + testTriggerUUID
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.DeleteTrigger(testNamespace, testRepository, testTriggerUUID)
	if err != nil {
		t.Fatalf("DeleteTrigger returned error: %v", err)
	}
}

func TestUpdateTrigger(t *testing.T) {
	mockResponse := BuildTrigger{
		ID:       testTriggerUUID,
		Service:  testTriggerService,
		IsActive: true,
		Enabled:  false,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/trigger/" + testTriggerUUID
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

	trigger, err := client.UpdateTrigger(testNamespace, testRepository, testTriggerUUID, false)
	if err != nil {
		t.Fatalf("UpdateTrigger returned error: %v", err)
	}

	if trigger.ID != testTriggerUUID {
		t.Errorf("Expected trigger ID %s, got %s", testTriggerUUID, trigger.ID)
	}
}

func TestStartTriggerBuild(t *testing.T) {
	mockResponse := Build{
		ID:    "build-uuid-789",
		Phase: "waiting",
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/trigger/" + testTriggerUUID + "/start"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	triggerReq := &ManualTriggerRequest{
		CommitSHA: "abc123def456",
	}

	build, err := client.StartTriggerBuild(testNamespace, testRepository, testTriggerUUID, triggerReq)
	if err != nil {
		t.Fatalf("StartTriggerBuild returned error: %v", err)
	}

	if build.ID != "build-uuid-789" {
		t.Errorf("Expected build ID 'build-uuid-789', got %s", build.ID)
	}
}

func TestActivateTrigger(t *testing.T) {
	mockResponse := BuildTrigger{
		ID:       testTriggerUUID,
		Service:  testTriggerService,
		IsActive: true,
		Enabled:  true,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/trigger/" + testTriggerUUID + "/activate"
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

	activateReq := &ActivateTriggerRequest{
		Config: map[string]interface{}{
			"build_source":    "myorg/myrepo",
			"dockerfile_path": "/Dockerfile",
		},
	}

	trigger, err := client.ActivateTrigger(testNamespace, testRepository, testTriggerUUID, activateReq)
	if err != nil {
		t.Fatalf("ActivateTrigger returned error: %v", err)
	}

	if trigger.ID != testTriggerUUID {
		t.Errorf("Expected trigger ID %s, got %s", testTriggerUUID, trigger.ID)
	}
	if !trigger.IsActive {
		t.Error("Expected trigger to be active")
	}
}

func TestGetTriggersError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetTriggers(testNamespace, testRepository)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetTriggerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetTrigger(testNamespace, testRepository, "nonexistent-uuid")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
