package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testManifestDigest = "sha256:def456"
)

func TestGetTag(t *testing.T) {
	mockTag := Tag{
		Name:           "v1.0.0",
		ManifestDigest: "sha256:abc123def456",
		Size:           1024000,
		LastModified:   "2024-01-15T10:30:00Z",
		IsManifestList: false,
		DockerImageID:  "abc123",
		ImageID:        "def456",
		V1Metadata: map[string]string{
			"architecture": "amd64",
			"os":           "linux",
		},
	}

	mockResponseJSON, _ := json.Marshal(mockTag)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/tag/v1.0.0" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/tag/v1.0.0, got %s", r.URL.Path)
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

	tag, err := client.GetTag("testorg", "testrepo", "v1.0.0")
	if err != nil {
		t.Fatalf("GetTag failed: %v", err)
	}

	if tag.Name != "v1.0.0" {
		t.Errorf("Expected tag name 'v1.0.0', got '%s'", tag.Name)
	}
	if tag.ManifestDigest != "sha256:abc123def456" {
		t.Errorf("Expected manifest digest 'sha256:abc123def456', got '%s'", tag.ManifestDigest)
	}
	if tag.Size != 1024000 {
		t.Errorf("Expected size 1024000, got %d", tag.Size)
	}
	if tag.V1Metadata["architecture"] != "amd64" {
		t.Errorf("Expected architecture 'amd64', got '%s'", tag.V1Metadata["architecture"])
	}
}

func TestUpdateTag(t *testing.T) {
	mockTag := Tag{
		Name:         "v1.0.0",
		Expiration:   "2024-12-31T23:59:59Z",
		LastModified: "2024-01-15T10:30:00Z",
	}

	mockResponseJSON, _ := json.Marshal(mockTag)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/tag/v1.0.0" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/tag/v1.0.0, got %s", r.URL.Path)
		}

		// Verify request body
		var req UpdateTagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Expiration != "2024-12-31T23:59:59Z" {
			t.Errorf("Expected expiration '2024-12-31T23:59:59Z', got '%s'", req.Expiration)
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

	tag, err := client.UpdateTag("testorg", "testrepo", "v1.0.0", "2024-12-31T23:59:59Z")
	if err != nil {
		t.Fatalf("UpdateTag failed: %v", err)
	}

	if tag.Expiration != "2024-12-31T23:59:59Z" {
		t.Errorf("Expected expiration '2024-12-31T23:59:59Z', got '%s'", tag.Expiration)
	}
}

func TestDeleteTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/tag/old-version" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/tag/old-version, got %s", r.URL.Path)
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

	err = client.DeleteTag("testorg", "testrepo", "old-version")
	if err != nil {
		t.Fatalf("DeleteTag failed: %v", err)
	}
}

func TestGetTagHistory(t *testing.T) {
	mockHistory := TagHistory{
		Tags: []Tag{
			{
				Name:           "latest",
				ManifestDigest: "sha256:abc123",
				LastModified:   "2024-01-15T10:30:00Z",
				Size:           1024000,
			},
			{
				Name:           "latest",
				ManifestDigest: testManifestDigest,
				LastModified:   "2024-01-10T08:15:00Z",
				Size:           1020000,
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockHistory)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/tag/latest/history" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/tag/latest/history, got %s", r.URL.Path)
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

	history, err := client.GetTagHistory("testorg", "testrepo", "latest")
	if err != nil {
		t.Fatalf("GetTagHistory failed: %v", err)
	}

	if len(history.Tags) != 2 {
		t.Errorf("Expected 2 tags in history, got %d", len(history.Tags))
	}

	// Check first tag (most recent)
	if history.Tags[0].ManifestDigest != "sha256:abc123" {
		t.Errorf("Expected first tag manifest 'sha256:abc123', got '%s'", history.Tags[0].ManifestDigest)
	}

	// Check second tag (older)
	if history.Tags[1].ManifestDigest != testManifestDigest {
		t.Errorf("Expected second tag manifest 'sha256:def456', got '%s'", history.Tags[1].ManifestDigest)
	}
}

func TestRevertTag(t *testing.T) {
	mockTag := Tag{
		Name:           "latest",
		ManifestDigest: testManifestDigest,
		LastModified:   "2024-01-15T11:00:00Z",
	}

	mockResponseJSON, _ := json.Marshal(mockTag)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/tag/latest/revert" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/tag/latest/revert, got %s", r.URL.Path)
		}

		// Verify request body
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req["manifest_digest"] != testManifestDigest {
			t.Errorf("Expected manifest_digest 'sha256:def456', got '%v'", req["manifest_digest"])
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

	tag, err := client.RevertTag("testorg", "testrepo", "latest", testManifestDigest)
	if err != nil {
		t.Fatalf("RevertTag failed: %v", err)
	}

	if tag.ManifestDigest != testManifestDigest {
		t.Errorf("Expected reverted tag manifest 'sha256:def456', got '%s'", tag.ManifestDigest)
	}
}

func TestTagErrorHandling(t *testing.T) {
	// Test 404 error for non-existent tag
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Tag not found"}`))
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetTag error
	_, err = client.GetTag("testorg", "testrepo", "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent tag, got nil")
	}

	// Test UpdateTag error
	_, err = client.UpdateTag("testorg", "testrepo", "nonexistent", "2024-12-31T23:59:59Z")
	if err == nil {
		t.Error("Expected error for non-existent tag, got nil")
	}

	// Test DeleteTag error
	err = client.DeleteTag("testorg", "testrepo", "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent tag, got nil")
	}

	// Test GetTagHistory error
	_, err = client.GetTagHistory("testorg", "testrepo", "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent tag, got nil")
	}

	// Test RevertTag error
	_, err = client.RevertTag("testorg", "testrepo", "nonexistent", "sha256:abc123")
	if err == nil {
		t.Error("Expected error for non-existent tag, got nil")
	}
}
