package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetSearch   = "GET"
	testSearchQuery = "quay"
)

func TestSearchRepositories(t *testing.T) {
	mockResult := SearchRepositoryResult{
		Results: []SearchRepository{
			{
				Namespace:   testSearchQuery,
				Name:        testSearchQuery,
				Description: "Quay container registry",
				IsPublic:    true,
				Kind:        "image",
				Popularity:  100.5,
				Score:       0.95,
			},
			{
				Namespace:   "redhat",
				Name:        "quay-operator",
				Description: "Quay Operator for OpenShift",
				IsPublic:    true,
				Kind:        "image",
				Popularity:  50.2,
				Score:       0.85,
			},
		},
		HasAdditional: true,
		Page:          1,
		StartIndex:    0,
	}

	mockResponseJSON, _ := json.Marshal(mockResult)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetSearch {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/find/repositories" {
			t.Errorf("Expected path /api/v1/find/repositories, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != testSearchQuery {
			t.Errorf("Expected query '%s', got '%s'", testSearchQuery, r.URL.Query().Get("query"))
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

	result, err := client.SearchRepositories(testSearchQuery, 0)
	if err != nil {
		t.Fatalf("SearchRepositories failed: %v", err)
	}

	if len(result.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result.Results))
	}
	if result.Results[0].Namespace != testSearchQuery {
		t.Errorf("Expected namespace '%s', got '%s'", testSearchQuery, result.Results[0].Namespace)
	}
	if result.Results[0].Name != testSearchQuery {
		t.Errorf("Expected name '%s', got '%s'", testSearchQuery, result.Results[0].Name)
	}
	if !result.HasAdditional {
		t.Error("Expected HasAdditional to be true")
	}
}

func TestSearchRepositoriesWithPage(t *testing.T) {
	mockResult := SearchRepositoryResult{
		Results:       []SearchRepository{},
		HasAdditional: false,
		Page:          2,
		StartIndex:    10,
	}

	mockResponseJSON, _ := json.Marshal(mockResult)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("Expected page '2', got '%s'", r.URL.Query().Get("page"))
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

	result, err := client.SearchRepositories("test", 2)
	if err != nil {
		t.Fatalf("SearchRepositories failed: %v", err)
	}

	if result.Page != 2 {
		t.Errorf("Expected page 2, got %d", result.Page)
	}
}

func TestSearchAll(t *testing.T) {
	mockResult := SearchAllResult{
		Results: []SearchEntity{
			{
				Name:        "quay",
				Description: "Quay container registry",
				Kind:        "repository",
				Score:       0.95,
				Avatar: Avatar{
					Name:  "quay",
					Hash:  "abc123",
					Kind:  "org",
					Color: "#ff0000",
				},
			},
			{
				Name:        "quayuser",
				Description: "Quay developer",
				Kind:        "user",
				Score:       0.75,
			},
			{
				Name:        "redhat",
				Description: "Red Hat organization",
				Kind:        "organization",
				Score:       0.85,
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockResult)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetSearch {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/find/all" {
			t.Errorf("Expected path /api/v1/find/all, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != testSearchQuery {
			t.Errorf("Expected query '%s', got '%s'", testSearchQuery, r.URL.Query().Get("query"))
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

	result, err := client.SearchAll(testSearchQuery)
	if err != nil {
		t.Fatalf("SearchAll failed: %v", err)
	}

	if len(result.Results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result.Results))
	}

	// Verify we got different entity types
	kinds := make(map[string]bool)
	for _, r := range result.Results {
		kinds[r.Kind] = true
	}
	if !kinds["repository"] {
		t.Error("Expected to find repository in results")
	}
	if !kinds["user"] {
		t.Error("Expected to find user in results")
	}
	if !kinds["organization"] {
		t.Error("Expected to find organization in results")
	}
}

func TestSearchErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test SearchRepositories error
	_, err = client.SearchRepositories("test", 0)
	if err == nil {
		t.Error("Expected error for failed request, got nil")
	}

	// Test SearchAll error
	_, err = client.SearchAll("test")
	if err == nil {
		t.Error("Expected error for failed request, got nil")
	}
}

func TestSearchEmptyResults(t *testing.T) {
	mockResult := SearchRepositoryResult{
		Results:       []SearchRepository{},
		HasAdditional: false,
		Page:          1,
		StartIndex:    0,
	}

	mockResponseJSON, _ := json.Marshal(mockResult)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	result, err := client.SearchRepositories("nonexistent-repo-xyz", 0)
	if err != nil {
		t.Fatalf("SearchRepositories failed: %v", err)
	}

	if len(result.Results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(result.Results))
	}
	if result.HasAdditional {
		t.Error("Expected HasAdditional to be false for empty results")
	}
}
