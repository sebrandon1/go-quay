package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	mockUser := UserDetails{
		Anonymous:      false,
		Username:       testUserName,
		Email:          testEmailAddress,
		Verified:       true,
		CanCreateRepo:  true,
		PreferredUsers: false,
		TagExpirationS: 2592000, // 30 days
		Avatar: Avatar{
			Name: testUserName,
			Kind: testKindUser,
		},
		Organizations: []User{
			{
				Name:     testNamespace,
				Username: testNamespace,
				Kind:     "organization",
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockUser)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	user, err := client.GetUser()
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}

	if user.Username != testUserName {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
	if user.Email != testEmailAddress {
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
	if user.Organizations[0].Name != testNamespace {
		t.Errorf("Expected organization name 'testorg', got '%s'", user.Organizations[0].Name)
	}
}

func TestGetStarredRepositories(t *testing.T) {
	mockStarred := StarredRepositories{
		Repositories: []StarredRepository{
			{
				Namespace:    testSearchQueryValue,
				Name:         testSearchQueryValue,
				Description:  "Container registry and security scanner",
				IsPublic:     true,
				Kind:         testKindImage,
				LastModified: "2024-01-15T10:30:00Z",
				Popularity:   95.5,
			},
			{
				Namespace:    testNamespace,
				Name:         "myapp",
				Description:  "My test application",
				IsPublic:     false,
				Kind:         testKindImage,
				LastModified: "2024-01-10T08:15:00Z",
				Popularity:   10.2,
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockStarred)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
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
	if repo1.Namespace != testSearchQueryValue || repo1.Name != testSearchQueryValue {
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
	if repo2.Namespace != testNamespace || repo2.Name != "myapp" {
		t.Errorf("Expected second repository 'testorg/myapp', got '%s/%s'", repo2.Namespace, repo2.Name)
	}
	if repo2.IsPublic != false {
		t.Errorf("Expected second repository to be private, got %v", repo2.IsPublic)
	}
}

func TestStarRepository(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/quay/quay/star" {
			t.Errorf("Expected path /api/v1/repository/quay/quay/star, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.StarRepository(testSearchQueryValue, testSearchQueryValue)
	if err != nil {
		t.Fatalf("StarRepository failed: %v", err)
	}
}

func TestUnstarRepository(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/quay/quay/star" {
			t.Errorf("Expected path /api/v1/repository/quay/quay/star, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.UnstarRepository(testSearchQueryValue, testSearchQueryValue)
	if err != nil {
		t.Fatalf("UnstarRepository failed: %v", err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	mockUser := UserDetails{
		Username:      testUserName,
		Email:         testEmailAddress,
		Verified:      true,
		CanCreateRepo: true,
	}
	mockResponseJSON, _ := json.Marshal(mockUser)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/users/" + testUserName
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

	user, err := client.GetUserByUsername(testUserName)
	if err != nil {
		t.Fatalf("GetUserByUsername returned error: %v", err)
	}

	if user.Username != testUserName {
		t.Errorf("Expected username %s, got %s", testUserName, user.Username)
	}
	if user.Email != testEmailAddress {
		t.Errorf("Expected email %s, got %s", testEmailAddress, user.Email)
	}
}

func TestGetUserMarketplace(t *testing.T) {
	mockMarketplace := MarketplaceInfo{
		HasPayment: false,
		Subscriptions: []MarketplaceSubscription{
			{ID: "user-sub-123", SKU: "free-plan", Status: "active"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockMarketplace)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/user/marketplace"
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

	marketplace, err := client.GetUserMarketplace()
	if err != nil {
		t.Fatalf("GetUserMarketplace returned error: %v", err)
	}

	if marketplace.HasPayment {
		t.Errorf("Expected HasPayment to be false")
	}
	if len(marketplace.Subscriptions) != 1 {
		t.Errorf("Expected 1 subscription, got %d", len(marketplace.Subscriptions))
	}
	if marketplace.Subscriptions[0].ID != "user-sub-123" {
		t.Errorf("Expected subscription ID 'user-sub-123', got %s", marketplace.Subscriptions[0].ID)
	}
}

func TestUserErrorHandling(t *testing.T) {
	// Test unauthorized error (401)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid token"}`))
	}))
	defer server.Close()

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
	err = client.StarRepository(testSearchQueryValue, testSearchQueryValue)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// Test UnstarRepository error
	err = client.UnstarRepository(testSearchQueryValue, testSearchQueryValue)
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
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

func TestUserHTTPErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetUserByUsername(testUserName)
	if err == nil {
		t.Error("Expected error from GetUserByUsername, got nil")
	}

	_, err = client.GetUserMarketplace()
	if err == nil {
		t.Error("Expected error from GetUserMarketplace, got nil")
	}
}
