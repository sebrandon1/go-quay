package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testRepoPath       = "/api/v1/repository/testorg/testrepo"
	testTagPath        = "/api/v1/repository/testorg/testrepo/tag"
	httpGetRepo        = "GET"
	updatedDescription = "Updated description"
)

func TestCreateRepository(t *testing.T) {
	// Mock response
	mockRepo := Repository{
		Namespace:   "testorg",
		Name:        "testrepo",
		Description: "Test repository",
		IsPublic:    false,
		Kind:        "image",
	}

	mockResponseJSON, _ := json.Marshal(mockRepo)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository" {
			t.Errorf("Expected path /api/v1/repository, got %s", r.URL.Path)
		}

		// Verify request body
		var req CreateRepositoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		expectedReq := CreateRepositoryRequest{
			Repository:  "testrepo",
			Namespace:   "testorg",
			Visibility:  "private",
			Description: "Test repository",
		}

		if req != expectedReq {
			t.Errorf("Request body mismatch. Expected %+v, got %+v", expectedReq, req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	// Override QuayURL for testing
	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	// Create client and test
	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	repo, err := client.CreateRepository("testorg", "testrepo", "private", "Test repository")
	if err != nil {
		t.Fatalf("CreateRepository failed: %v", err)
	}

	// Verify response
	if repo.Name != "testrepo" {
		t.Errorf("Expected repository name 'testrepo', got '%s'", repo.Name)
	}
	if repo.Namespace != "testorg" {
		t.Errorf("Expected namespace 'testorg', got '%s'", repo.Namespace)
	}
}

func TestUpdateRepository(t *testing.T) {
	mockRepo := Repository{
		Namespace:   "testorg",
		Name:        "testrepo",
		Description: updatedDescription,
		IsPublic:    true,
	}

	mockResponseJSON, _ := json.Marshal(mockRepo)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != testRepoPath {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo, got %s", r.URL.Path)
		}

		var req UpdateRepositoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Description != updatedDescription {
			t.Errorf("Expected description 'Updated description', got '%s'", req.Description)
		}
		if req.Visibility != "public" {
			t.Errorf("Expected visibility 'public', got '%s'", req.Visibility)
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

	repo, err := client.UpdateRepository("testorg", "testrepo", updatedDescription, "public")
	if err != nil {
		t.Fatalf("UpdateRepository failed: %v", err)
	}

	if repo.Description != updatedDescription {
		t.Errorf("Expected description 'Updated description', got '%s'", repo.Description)
	}
}

func TestDeleteRepository(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != testRepoPath {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo, got %s", r.URL.Path)
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

	err = client.DeleteRepository("testorg", "testrepo")
	if err != nil {
		t.Fatalf("DeleteRepository failed: %v", err)
	}
}

func TestGetRepository(t *testing.T) {
	mockRepo := Repository{
		Namespace:   "testorg",
		Name:        "testrepo",
		Description: "Test repository",
		IsPublic:    false,
		Kind:        "image",
	}

	mockTags := RepositoryTags{
		Tags: []struct {
			Name           string `json:"name,omitempty"`
			Reversion      bool   `json:"reversion,omitempty"`
			StartTs        int    `json:"start_ts,omitempty"`
			ManifestDigest string `json:"manifest_digest,omitempty"`
			IsManifestList bool   `json:"is_manifest_list,omitempty"`
			Size           any    `json:"size,omitempty"`
			LastModified   string `json:"last_modified,omitempty"`
			EndTs          int    `json:"end_ts,omitempty"`
			Expiration     string `json:"expiration,omitempty"`
		}{
			{Name: "latest", ManifestDigest: "sha256:abc123"},
			{Name: "v1.0.0", ManifestDigest: "sha256:def456"},
		},
		Page:          1,
		HasAdditional: false,
	}

	mockRepoJSON, _ := json.Marshal(mockRepo)
	mockTagsJSON, _ := json.Marshal(mockTags)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetRepo {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Return different responses based on the path
		switch r.URL.Path {
		case testRepoPath:
			w.Write(mockRepoJSON)
		case testTagPath:
			w.Write(mockTagsJSON)
		default:
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	// Override base URL since GetRepository uses hardcoded URL
	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	repoWithTags, err := client.GetRepository("testorg", "testrepo")
	if err != nil {
		t.Fatalf("GetRepository failed: %v", err)
	}

	if repoWithTags.Name != "testrepo" {
		t.Errorf("Expected repository name 'testrepo', got '%s'", repoWithTags.Name)
	}
	if repoWithTags.Namespace != "testorg" {
		t.Errorf("Expected namespace 'testorg', got '%s'", repoWithTags.Namespace)
	}
	if len(repoWithTags.Tags.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(repoWithTags.Tags.Tags))
	}
}
