package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testManifestRef    = "sha256:abc123def456789"
	testLabelID        = "label-123"
	httpGetManifest    = "GET"
	httpPostManifest   = "POST"
	httpDeleteManifest = "DELETE"
)

func TestGetManifest(t *testing.T) {
	mockManifest := Manifest{
		Digest:         testManifestRef,
		SchemaVersion:  2,
		MediaType:      "application/vnd.docker.distribution.manifest.v2+json",
		Size:           1024000,
		IsManifestList: false,
		Layers: []ManifestLayer{
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      512000,
				Digest:    "sha256:layer1digest",
				Index:     0,
			},
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      512000,
				Digest:    "sha256:layer2digest",
				Index:     1,
			},
		},
		Config: ManifestConfig{
			MediaType: "application/vnd.docker.container.image.v1+json",
			Size:      1500,
			Digest:    "sha256:configdigest",
		},
	}

	mockResponseJSON, _ := json.Marshal(mockManifest)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetManifest {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/manifest/"+testManifestRef {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/manifest/%s, got %s", testManifestRef, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

	manifest, err := client.GetManifest("testorg", "testrepo", testManifestRef)
	if err != nil {
		t.Fatalf("GetManifest failed: %v", err)
	}

	if manifest.Digest != testManifestRef {
		t.Errorf("Expected manifest digest '%s', got '%s'", testManifestRef, manifest.Digest)
	}
	if manifest.SchemaVersion != 2 {
		t.Errorf("Expected schema version 2, got %d", manifest.SchemaVersion)
	}
	if len(manifest.Layers) != 2 {
		t.Errorf("Expected 2 layers, got %d", len(manifest.Layers))
	}
	if manifest.Config.Digest != "sha256:configdigest" {
		t.Errorf("Expected config digest 'sha256:configdigest', got '%s'", manifest.Config.Digest)
	}
}

func TestDeleteManifest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteManifest {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/manifest/"+testManifestRef {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/manifest/%s, got %s", testManifestRef, r.URL.Path)
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

	err = client.DeleteManifest("testorg", "testrepo", testManifestRef)
	if err != nil {
		t.Fatalf("DeleteManifest failed: %v", err)
	}
}

func TestGetManifestLabels(t *testing.T) {
	mockLabels := ManifestLabels{
		Labels: []ManifestLabel{
			{
				ID:         "label-1",
				Key:        "version",
				Value:      "1.0.0",
				SourceType: "api",
				MediaType:  "text/plain",
			},
			{
				ID:         "label-2",
				Key:        "maintainer",
				Value:      "test@example.com",
				SourceType: "api",
				MediaType:  "text/plain",
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockLabels)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetManifest {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/manifest/"+testManifestRef+"/labels" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/manifest/%s/labels, got %s", testManifestRef, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

	labels, err := client.GetManifestLabels("testorg", "testrepo", testManifestRef)
	if err != nil {
		t.Fatalf("GetManifestLabels failed: %v", err)
	}

	if len(labels.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(labels.Labels))
	}
	if labels.Labels[0].Key != "version" {
		t.Errorf("Expected first label key 'version', got '%s'", labels.Labels[0].Key)
	}
	if labels.Labels[1].Value != "test@example.com" {
		t.Errorf("Expected second label value 'test@example.com', got '%s'", labels.Labels[1].Value)
	}
}

func TestAddManifestLabel(t *testing.T) {
	mockLabel := ManifestLabel{
		ID:         "new-label-123",
		Key:        "environment",
		Value:      "production",
		SourceType: "api",
		MediaType:  "text/plain",
	}

	mockResponseJSON, _ := json.Marshal(mockLabel)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPostManifest {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/manifest/"+testManifestRef+"/labels" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/manifest/%s/labels, got %s", testManifestRef, r.URL.Path)
		}

		// Verify request body
		var req AddManifestLabelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Key != "environment" {
			t.Errorf("Expected key 'environment', got '%s'", req.Key)
		}
		if req.Value != "production" {
			t.Errorf("Expected value 'production', got '%s'", req.Value)
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

	label, err := client.AddManifestLabel("testorg", "testrepo", testManifestRef, "environment", "production", "text/plain")
	if err != nil {
		t.Fatalf("AddManifestLabel failed: %v", err)
	}

	if label.Key != "environment" {
		t.Errorf("Expected label key 'environment', got '%s'", label.Key)
	}
	if label.Value != "production" {
		t.Errorf("Expected label value 'production', got '%s'", label.Value)
	}
	if label.ID != "new-label-123" {
		t.Errorf("Expected label ID 'new-label-123', got '%s'", label.ID)
	}
}

func TestGetManifestLabel(t *testing.T) {
	mockLabel := ManifestLabel{
		ID:         testLabelID,
		Key:        "version",
		Value:      "2.0.0",
		SourceType: "api",
		MediaType:  "text/plain",
	}

	mockResponseJSON, _ := json.Marshal(mockLabel)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetManifest {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/testorg/testrepo/manifest/" + testManifestRef + "/labels/" + testLabelID
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

	label, err := client.GetManifestLabel("testorg", "testrepo", testManifestRef, testLabelID)
	if err != nil {
		t.Fatalf("GetManifestLabel failed: %v", err)
	}

	if label.ID != testLabelID {
		t.Errorf("Expected label ID '%s', got '%s'", testLabelID, label.ID)
	}
	if label.Key != "version" {
		t.Errorf("Expected label key 'version', got '%s'", label.Key)
	}
}

func TestDeleteManifestLabel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteManifest {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/testorg/testrepo/manifest/" + testManifestRef + "/labels/" + testLabelID
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

	err = client.DeleteManifestLabel("testorg", "testrepo", testManifestRef, testLabelID)
	if err != nil {
		t.Fatalf("DeleteManifestLabel failed: %v", err)
	}
}

func TestManifestErrorHandling(t *testing.T) {
	// Test 404 error for non-existent manifest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Manifest not found"}`))
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetManifest error
	_, err = client.GetManifest("testorg", "testrepo", "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}

	// Test DeleteManifest error
	err = client.DeleteManifest("testorg", "testrepo", "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}

	// Test GetManifestLabels error
	_, err = client.GetManifestLabels("testorg", "testrepo", "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}

	// Test AddManifestLabel error
	_, err = client.AddManifestLabel("testorg", "testrepo", "nonexistent", "key", "value", "")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}

	// Test GetManifestLabel error
	_, err = client.GetManifestLabel("testorg", "testrepo", "nonexistent", "labelid")
	if err == nil {
		t.Error("Expected error for non-existent manifest label, got nil")
	}

	// Test DeleteManifestLabel error
	err = client.DeleteManifestLabel("testorg", "testrepo", "nonexistent", "labelid")
	if err == nil {
		t.Error("Expected error for non-existent manifest label, got nil")
	}
}
