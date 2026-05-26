package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDiscovery(t *testing.T) {
	mockResponse := Discovery{
		Version: "v1",
		APIs: map[string]DiscoveryAPI{
			"repository": {
				Path:        testAPIPathRepository,
				Methods:     []string{"GET", "POST"},
				Description: "Repository operations",
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
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

func TestGetRegistryCapabilities(t *testing.T) {
	mockResponse := RegistryCapabilities{
		SparseManifests: SparseManifests{
			Supported:                    false,
			RequiredArchitectures:        []string{},
			OptionalArchitecturesAllowed: false,
		},
		MirrorArchitectures: []string{testArchAmd64, "arm64", "ppc64le", "s390x"},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/registry/capabilities"
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

	capabilities, err := client.GetRegistryCapabilities()
	if err != nil {
		t.Fatalf("GetRegistryCapabilities returned error: %v", err)
	}

	if len(capabilities.MirrorArchitectures) != 4 {
		t.Errorf("Expected 4 mirror architectures, got %d", len(capabilities.MirrorArchitectures))
	}
	if capabilities.SparseManifests.Supported {
		t.Error("Expected sparse manifests to be unsupported")
	}
}

func TestGetAppInfo(t *testing.T) {
	mockApp := Application{
		ClientID:       testClientID,
		Name:           "My App",
		Description:    "An OAuth application",
		ApplicationURI: "https://app.example.com",
		RedirectURI:    "https://app.example.com/callback",
	}
	mockResponseJSON, _ := json.Marshal(mockApp)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/app/" + testClientID
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

	app, err := client.GetAppInfo(testClientID)
	if err != nil {
		t.Fatalf("GetAppInfo returned error: %v", err)
	}

	if app.ClientID != testClientID {
		t.Errorf("Expected client ID %s, got %s", testClientID, app.ClientID)
	}
	if app.Name != "My App" {
		t.Errorf("Expected name 'My App', got %s", app.Name)
	}
}

func TestGetEntities(t *testing.T) {
	mockEntities := Entities{
		Results: []Entity{
			{Name: testUserName, Kind: testKindUser, IsRobot: false},
			{Name: "testorg+bot", Kind: testKindRobot, IsRobot: true},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockEntities)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/entities/test"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if r.URL.Query().Get("includeOrgs") != testQueryValueTrue {
			t.Errorf("Expected includeOrgs=true query parameter")
		}
		if r.URL.Query().Get("includeTeams") != testQueryValueTrue {
			t.Errorf("Expected includeTeams=true query parameter")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	entities, err := client.GetEntities("test", true, true)
	if err != nil {
		t.Fatalf("GetEntities returned error: %v", err)
	}

	if len(entities.Results) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(entities.Results))
	}
	if entities.Results[0].Name != testUserName {
		t.Errorf("Expected first entity name %s, got %s", testUserName, entities.Results[0].Name)
	}
}

func TestGetDiscoveryError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetDiscovery()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
