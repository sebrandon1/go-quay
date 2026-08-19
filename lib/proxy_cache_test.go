package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Proxy Cache ---

func TestGetProxyCacheConfig(t *testing.T) {
	mockConfig := ProxyCacheConfig{
		UpstreamRegistry: testUpstreamReg,
		Insecure:         false,
		Expiration:       86400,
	}
	mockResponseJSON, _ := json.Marshal(mockConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/proxycache"
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

	config, err := client.GetProxyCacheConfig(context.Background(), testOrgName)
	if err != nil {
		t.Fatalf("GetProxyCacheConfig returned error: %v", err)
	}

	if config.UpstreamRegistry != testUpstreamReg {
		t.Errorf("Expected upstream registry %s, got %s", testUpstreamReg, config.UpstreamRegistry)
	}
	if config.Expiration != 86400 {
		t.Errorf("Expected expiration 86400, got %d", config.Expiration)
	}
}

func TestCreateProxyCacheConfig(t *testing.T) {
	mockConfig := ProxyCacheConfig{
		UpstreamRegistry: testUpstreamReg,
		Insecure:         true,
		Expiration:       3600,
	}
	mockResponseJSON, _ := json.Marshal(mockConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/proxycache"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		var req CreateProxyCacheConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.UpstreamRegistry != testUpstreamReg {
			t.Errorf("Expected upstream registry %s, got %s", testUpstreamReg, req.UpstreamRegistry)
		}
		if !req.Insecure {
			t.Errorf("Expected Insecure to be true")
		}
		if req.Expiration != 3600 {
			t.Errorf("Expected expiration 3600, got %d", req.Expiration)
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

	config, err := client.CreateProxyCacheConfig(context.Background(), testOrgName, testUpstreamReg, true, 3600)
	if err != nil {
		t.Fatalf("CreateProxyCacheConfig returned error: %v", err)
	}

	if !config.Insecure {
		t.Errorf("Expected Insecure to be true")
	}
}

func TestDeleteProxyCacheConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/proxycache"
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

	err = client.DeleteProxyCacheConfig(context.Background(), testOrgName)
	if err != nil {
		t.Fatalf("DeleteProxyCacheConfig returned error: %v", err)
	}
}

func TestProxyCacheHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetProxyCacheConfig(context.Background(), testOrgName)
	if err == nil {
		t.Error("Expected error from GetProxyCacheConfig, got nil")
	}

	_, err = client.CreateProxyCacheConfig(context.Background(), testOrgName, testUpstreamReg, false, 86400)
	if err == nil {
		t.Error("Expected error from CreateProxyCacheConfig, got nil")
	}

	err = client.DeleteProxyCacheConfig(context.Background(), testOrgName)
	if err == nil {
		t.Error("Expected error from DeleteProxyCacheConfig, got nil")
	}
}
