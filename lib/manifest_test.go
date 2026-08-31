package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testManifestRef = "sha256:abc123def456789"
	testLabelID     = "label-123"
)

func TestGetManifest(t *testing.T) {
	mockManifest := Manifest{
		Digest:               testManifestRef,
		IsManifestList:       false,
		LayersCompressedSize: 1024000,
		Layers: []ManifestLayer{
			{
				CompressedSize: 512000,
				BlobDigest:     "sha256:layer1digest",
				Index:          0,
				Command:        []string{"bazel build //foo"},
			},
			{
				CompressedSize: 512000,
				BlobDigest:     "sha256:layer2digest",
				Index:          1,
			},
		},
	}

	mockResponseJSON := mustMarshal(t, mockManifest)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	manifest, err := client.GetManifest(context.Background(), testNamespace, testRepository, testManifestRef)
	if err != nil {
		t.Fatalf("GetManifest failed: %v", err)
	}

	if manifest.Digest != testManifestRef {
		t.Errorf("Expected manifest digest '%s', got '%s'", testManifestRef, manifest.Digest)
	}
	if len(manifest.Layers) != 2 {
		t.Errorf("Expected 2 layers, got %d", len(manifest.Layers))
	}
	if manifest.LayersCompressedSize != 1024000 {
		t.Errorf("Expected layers_compressed_size 1024000, got %d", manifest.LayersCompressedSize)
	}
	if manifest.Layers[0].BlobDigest != "sha256:layer1digest" {
		t.Errorf("Expected layer 0 blob_digest 'sha256:layer1digest', got '%s'", manifest.Layers[0].BlobDigest)
	}
}

func TestDeleteManifest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/manifest/"+testManifestRef {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/manifest/%s, got %s", testManifestRef, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.DeleteManifest(context.Background(), testNamespace, testRepository, testManifestRef)
	if err != nil {
		t.Fatalf("DeleteManifest failed: %v", err)
	}
}

func TestGetManifestLabels(t *testing.T) {
	mockLabels := ManifestLabels{
		Labels: []ManifestLabel{
			{
				ID:         "label-1",
				Key:        testLabelKeyVersion,
				Value:      "1.0.0",
				SourceType: testLabelKeyAPI,
				MediaType:  testMediaTypePlain,
			},
			{
				ID:         "label-2",
				Key:        "maintainer",
				Value:      testEmailAddress,
				SourceType: testLabelKeyAPI,
				MediaType:  testMediaTypePlain,
			},
		},
	}

	mockResponseJSON := mustMarshal(t, mockLabels)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	labels, err := client.GetManifestLabels(context.Background(), testNamespace, testRepository, testManifestRef)
	if err != nil {
		t.Fatalf("GetManifestLabels failed: %v", err)
	}

	if len(labels.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(labels.Labels))
	}
	if labels.Labels[0].Key != testLabelKeyVersion {
		t.Errorf("Expected first label key 'version', got '%s'", labels.Labels[0].Key)
	}
	if labels.Labels[1].Value != testEmailAddress {
		t.Errorf("Expected second label value 'test@example.com', got '%s'", labels.Labels[1].Value)
	}
}

func TestAddManifestLabel(t *testing.T) {
	mockLabel := ManifestLabel{
		ID:         "new-label-123",
		Key:        testLabelKeyEnvironment,
		Value:      testLabelValProduction,
		SourceType: testLabelKeyAPI,
		MediaType:  testMediaTypePlain,
	}

	mockResponseJSON := mustMarshal(t, mockLabel)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
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

		if req.Key != testLabelKeyEnvironment {
			t.Errorf("Expected key 'environment', got '%s'", req.Key)
		}
		if req.Value != testLabelValProduction {
			t.Errorf("Expected value 'production', got '%s'", req.Value)
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

	label, err := client.AddManifestLabel(context.Background(), testNamespace, testRepository, testManifestRef, testLabelKeyEnvironment, testLabelValProduction, testMediaTypePlain)
	if err != nil {
		t.Fatalf("AddManifestLabel failed: %v", err)
	}

	if label.Key != testLabelKeyEnvironment {
		t.Errorf("Expected label key 'environment', got '%s'", label.Key)
	}
	if label.Value != testLabelValProduction {
		t.Errorf("Expected label value 'production', got '%s'", label.Value)
	}
	if label.ID != "new-label-123" {
		t.Errorf("Expected label ID 'new-label-123', got '%s'", label.ID)
	}
}

func TestGetManifestLabel(t *testing.T) {
	mockLabel := ManifestLabel{
		ID:         testLabelID,
		Key:        testLabelKeyVersion,
		Value:      "2.0.0",
		SourceType: testLabelKeyAPI,
		MediaType:  testMediaTypePlain,
	}

	mockResponseJSON := mustMarshal(t, mockLabel)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	label, err := client.GetManifestLabel(context.Background(), testNamespace, testRepository, testManifestRef, testLabelID)
	if err != nil {
		t.Fatalf("GetManifestLabel failed: %v", err)
	}

	if label.ID != testLabelID {
		t.Errorf("Expected label ID '%s', got '%s'", testLabelID, label.ID)
	}
	if label.Key != testLabelKeyVersion {
		t.Errorf("Expected label key 'version', got '%s'", label.Key)
	}
}

func TestDeleteManifestLabel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/testorg/testrepo/manifest/" + testManifestRef + "/labels/" + testLabelID
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

	err = client.DeleteManifestLabel(context.Background(), testNamespace, testRepository, testManifestRef, testLabelID)
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetManifest error
	_, err = client.GetManifest(context.Background(), testNamespace, testRepository, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}

	// Test DeleteManifest error
	err = client.DeleteManifest(context.Background(), testNamespace, testRepository, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}

	// Test GetManifestLabels error
	_, err = client.GetManifestLabels(context.Background(), testNamespace, testRepository, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}

	// Test AddManifestLabel error
	_, err = client.AddManifestLabel(context.Background(), testNamespace, testRepository, "nonexistent", "key", "value", "")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}

	// Test GetManifestLabel error
	_, err = client.GetManifestLabel(context.Background(), testNamespace, testRepository, "nonexistent", "labelid")
	if err == nil {
		t.Error("Expected error for non-existent manifest label, got nil")
	}

	// Test DeleteManifestLabel error
	err = client.DeleteManifestLabel(context.Background(), testNamespace, testRepository, "nonexistent", "labelid")
	if err == nil {
		t.Error("Expected error for non-existent manifest label, got nil")
	}
}
