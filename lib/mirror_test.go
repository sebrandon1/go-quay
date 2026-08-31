package lib

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testExternalRef = "docker.io/library/nginx"

func TestGetMirrorConfig(t *testing.T) {
	mockConfig := MirrorConfig{
		IsEnabled:   true,
		MirrorType:  "PULL",
		ExternalRef: testExternalRef,
	}
	mockResponseJSON := mustMarshal(t, mockConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/mirror"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	config, err := client.GetMirrorConfig(context.Background(), testNamespace, testRepository)
	if err != nil {
		t.Fatalf("GetMirrorConfig returned error: %v", err)
	}

	if config.ExternalRef != testExternalRef {
		t.Errorf("Expected external ref 'docker.io/library/nginx', got '%s'", config.ExternalRef)
	}
}

func TestCreateMirrorConfig(t *testing.T) {
	mockConfig := MirrorConfig{
		IsEnabled:   true,
		ExternalRef: testExternalRef,
	}
	mockResponseJSON := mustMarshal(t, mockConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	createReq := &CreateMirrorConfigRequest{
		ExternalRef:   testExternalRef,
		SyncInterval:  86400,
		SyncStartDate: "2025-01-01T00:00:00Z",
		RobotUsername: "org+robot",
	}

	config, err := client.CreateMirrorConfig(context.Background(), testNamespace, testRepository, createReq)
	if err != nil {
		t.Fatalf("CreateMirrorConfig returned error: %v", err)
	}

	if config.ExternalRef != testExternalRef {
		t.Errorf("Expected external ref 'docker.io/library/nginx', got '%s'", config.ExternalRef)
	}
}

func TestUpdateMirrorConfig(t *testing.T) {
	enabled := false
	mockConfig := MirrorConfig{
		IsEnabled:   enabled,
		ExternalRef: testExternalRef,
	}
	mockResponseJSON := mustMarshal(t, mockConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	updateReq := &UpdateMirrorConfigRequest{
		IsEnabled: &enabled,
	}

	config, err := client.UpdateMirrorConfig(context.Background(), testNamespace, testRepository, updateReq)
	if err != nil {
		t.Fatalf("UpdateMirrorConfig returned error: %v", err)
	}

	if config.IsEnabled {
		t.Error("Expected mirror to be disabled")
	}
}

func TestMirrorConfigHTTPErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not_found"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetMirrorConfig(context.Background(), testNamespace, testRepository)
	if err == nil {
		t.Error("Expected error from GetMirrorConfig, got nil")
	}

	_, err = client.CreateMirrorConfig(context.Background(), testNamespace, testRepository, &CreateMirrorConfigRequest{})
	if err == nil {
		t.Error("Expected error from CreateMirrorConfig, got nil")
	}

	_, err = client.UpdateMirrorConfig(context.Background(), testNamespace, testRepository, &UpdateMirrorConfigRequest{})
	if err == nil {
		t.Error("Expected error from UpdateMirrorConfig, got nil")
	}
}

func TestMirrorConfigValidation(t *testing.T) {
	client, err := NewClient(testTokenValue)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.GetMirrorConfig(context.Background(), "", testRepository)
	if err == nil {
		t.Error("Expected error for empty namespace")
	}

	_, err = client.GetMirrorConfig(context.Background(), testNamespace, "")
	if err == nil {
		t.Error("Expected error for empty repository")
	}
}
