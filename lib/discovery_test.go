package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetDiscovery = "GET"
)

func TestGetDiscovery(t *testing.T) {
	mockResponse := Discovery{
		Version: "v1",
		APIs: map[string]DiscoveryAPI{
			"repository": {
				Path:        "/api/v1/repository",
				Methods:     []string{"GET", "POST"},
				Description: "Repository operations",
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetDiscovery {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/discovery"
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

	discovery, err := client.GetDiscovery()
	if err != nil {
		t.Fatalf("GetDiscovery returned error: %v", err)
	}

	if discovery.Version != "v1" {
		t.Errorf("Expected version 'v1', got %s", discovery.Version)
	}
	if len(discovery.APIs) != 1 {
		t.Errorf("Expected 1 API, got %d", len(discovery.APIs))
	}
}

func TestGetDiscoveryError(t *testing.T) {
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

	_, err = client.GetDiscovery()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
