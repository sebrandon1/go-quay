package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGet    = "GET"
	httpPut    = "PUT"
	httpDelete = "DELETE"
	testOrg    = "testorg"
)

func TestGetUser(t *testing.T) {
	mockUser := UserDetails{
		Anonymous:      false,
		Username:       "testuser",
		Email:          "test@example.com",
		Verified:       true,
		CanCreateRepo:  true,
		PreferredUsers: false,
		TagExpirationS: 2592000, // 30 days
		Avatar: Avatar{
			Name: "testuser",
			Kind: "user",
		},
		Organizations: []User{
			{
				Name:     testOrg,
				Username: testOrg,
				Kind:     "organization",
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockUser)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user" {
			t.Errorf("Expected path /api/v1/user, got %s", r.URL.Path)
		}

		// Verify Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", authHeader)
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

	user, err := client.GetUser()
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
	if user.Anonymous != false {
		t.Errorf("Expected Anonymous false, got %v", user.Anonymous)
	}
	if user.Verified != true {
		t.Errorf("Expected Verified true, got %v", user.Verified)
	}
	if user.CanCreateRepo != true {
		t.Errorf("Expected CanCreateRepo true, got %v", user.CanCreateRepo)
	}
	if len(user.Organizations) != 1 {
		t.Errorf("Expected 1 organization, got %d", len(user.Organizations))
	}
	if user.Organizations[0].Name != testOrg {
		t.Errorf("Expected organization name 'testorg', got '%s'", user.Organizations[0].Name)
	}
}

func TestGetStarredRepositories(t *testing.T) {
	mockStarred := StarredRepositories{
		Repositories: []StarredRepository{
			{
				Namespace:    "quay",
				Name:         "quay",
				Description:  "Container registry and security scanner",
				IsPublic:     true,
				Kind:         "image",
				LastModified: "2024-01-15T10:30:00Z",
				Popularity:   95.5,
			},
			{
				Namespace:    testOrg,
				Name:         "myapp",
				Description:  "My test application",
				IsPublic:     false,
				Kind:         "image",
				LastModified: "2024-01-10T08:15:00Z",
				Popularity:   10.2,
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockStarred)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/starred" {
			t.Errorf("Expected path /api/v1/user/starred, got %s", r.URL.Path)
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

	starred, err := client.GetStarredRepositories()
	if err != nil {
		t.Fatalf("GetStarredRepositories failed: %v", err)
	}

	if len(starred.Repositories) != 2 {
		t.Errorf("Expected 2 starred repositories, got %d", len(starred.Repositories))
	}

	// Check first repository
	repo1 := starred.Repositories[0]
	if repo1.Namespace != "quay" || repo1.Name != "quay" {
		t.Errorf("Expected first repository 'quay/quay', got '%s/%s'", repo1.Namespace, repo1.Name)
	}
	if repo1.IsPublic != true {
		t.Errorf("Expected first repository to be public, got %v", repo1.IsPublic)
	}
	if repo1.Popularity != 95.5 {
		t.Errorf("Expected first repository popularity 95.5, got %f", repo1.Popularity)
	}

	// Check second repository
	repo2 := starred.Repositories[1]
	if repo2.Namespace != testOrg || repo2.Name != "myapp" {
		t.Errorf("Expected second repository 'testorg/myapp', got '%s/%s'", repo2.Namespace, repo2.Name)
	}
	if repo2.IsPublic != false {
		t.Errorf("Expected second repository to be private, got %v", repo2.IsPublic)
	}
}

func TestStarRepository(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/quay/quay/star" {
			t.Errorf("Expected path /api/v1/repository/quay/quay/star, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.StarRepository("quay", "quay")
	if err != nil {
		t.Fatalf("StarRepository failed: %v", err)
	}
}

func TestUnstarRepository(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/quay/quay/star" {
			t.Errorf("Expected path /api/v1/repository/quay/quay/star, got %s", r.URL.Path)
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

	err = client.UnstarRepository("quay", "quay")
	if err != nil {
		t.Fatalf("UnstarRepository failed: %v", err)
	}
}

func TestUserErrorHandling(t *testing.T) {
	// Test unauthorized error (401)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid token"}`))
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("invalid-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetUser error
	_, err = client.GetUser()
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// Test GetStarredRepositories error
	_, err = client.GetStarredRepositories()
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// Test StarRepository error
	err = client.StarRepository("quay", "quay")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// Test UnstarRepository error
	err = client.UnstarRepository("quay", "quay")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestUserStarRepositoryNotFound(t *testing.T) {
	// Test 404 error for non-existent repository
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Repository not found"}`))
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.StarRepository("nonexistent", "repository")
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}

	err = client.UnstarRepository("nonexistent", "repository")
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}
}
