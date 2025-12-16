package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testRobotShortname = "deploybot"
	httpGetRobot       = "GET"
	httpPutRobot       = "PUT"
	httpPostRobot      = "POST"
	httpDeleteRobot    = "DELETE"
	roleWrite          = "write"
)

func TestGetUserRobotAccounts(t *testing.T) {
	mockRobots := RobotAccounts{
		Robots: []RobotAccount{
			{
				Name:        "testuser+deploybot",
				Description: "Deployment robot",
				Created:     "2024-01-15T10:30:00Z",
			},
			{
				Name:        "testuser+cibot",
				Description: "CI/CD robot",
				Created:     "2024-01-10T08:00:00Z",
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockRobots)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetRobot {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/robots" {
			t.Errorf("Expected path /api/v1/user/robots, got %s", r.URL.Path)
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

	robots, err := client.GetUserRobotAccounts()
	if err != nil {
		t.Fatalf("GetUserRobotAccounts failed: %v", err)
	}

	if len(robots.Robots) != 2 {
		t.Errorf("Expected 2 robots, got %d", len(robots.Robots))
	}
	if robots.Robots[0].Name != "testuser+deploybot" {
		t.Errorf("Expected robot name 'testuser+deploybot', got '%s'", robots.Robots[0].Name)
	}
}

func TestGetUserRobotAccount(t *testing.T) {
	mockRobot := RobotAccount{
		Name:        "testuser+deploybot",
		Description: "Deployment robot",
		Token:       "secret-token-123",
		Created:     "2024-01-15T10:30:00Z",
	}

	mockResponseJSON, _ := json.Marshal(mockRobot)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetRobot {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/user/robots/" + testRobotShortname
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

	robot, err := client.GetUserRobotAccount(testRobotShortname)
	if err != nil {
		t.Fatalf("GetUserRobotAccount failed: %v", err)
	}

	if robot.Name != "testuser+deploybot" {
		t.Errorf("Expected robot name 'testuser+deploybot', got '%s'", robot.Name)
	}
	if robot.Token != "secret-token-123" {
		t.Errorf("Expected token 'secret-token-123', got '%s'", robot.Token)
	}
}

func TestCreateUserRobotAccount(t *testing.T) {
	mockRobot := RobotAccount{
		Name:        "testuser+newbot",
		Description: "New automation robot",
		Token:       "new-secret-token",
		Created:     "2024-01-20T14:00:00Z",
	}

	mockResponseJSON, _ := json.Marshal(mockRobot)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPutRobot {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/robots/newbot" {
			t.Errorf("Expected path /api/v1/user/robots/newbot, got %s", r.URL.Path)
		}

		// Verify request body
		var req CreateRobotRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.Description != "New automation robot" {
			t.Errorf("Expected description 'New automation robot', got '%s'", req.Description)
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

	robot, err := client.CreateUserRobotAccount("newbot", "New automation robot", nil)
	if err != nil {
		t.Fatalf("CreateUserRobotAccount failed: %v", err)
	}

	if robot.Name != "testuser+newbot" {
		t.Errorf("Expected robot name 'testuser+newbot', got '%s'", robot.Name)
	}
	if robot.Token != "new-secret-token" {
		t.Errorf("Expected token 'new-secret-token', got '%s'", robot.Token)
	}
}

func TestDeleteUserRobotAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteRobot {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/user/robots/" + testRobotShortname
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

	err = client.DeleteUserRobotAccount(testRobotShortname)
	if err != nil {
		t.Fatalf("DeleteUserRobotAccount failed: %v", err)
	}
}

func TestRegenerateUserRobotToken(t *testing.T) {
	mockRobot := RobotAccount{
		Name:        "testuser+deploybot",
		Description: "Deployment robot",
		Token:       "regenerated-token-456",
		Created:     "2024-01-15T10:30:00Z",
	}

	mockResponseJSON, _ := json.Marshal(mockRobot)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPostRobot {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/user/robots/" + testRobotShortname + "/regenerate"
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

	robot, err := client.RegenerateUserRobotToken(testRobotShortname)
	if err != nil {
		t.Fatalf("RegenerateUserRobotToken failed: %v", err)
	}

	if robot.Token != "regenerated-token-456" {
		t.Errorf("Expected token 'regenerated-token-456', got '%s'", robot.Token)
	}
}

func TestGetUserRobotPermissions(t *testing.T) {
	mockPermissions := RobotPermissions{
		Permissions: []RobotPermission{
			{
				Repository: Repository{Name: "myrepo"},
				Role:       roleWrite,
			},
			{
				Repository: Repository{Name: "otherrepo"},
				Role:       "read",
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockPermissions)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetRobot {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/user/robots/" + testRobotShortname + "/permissions"
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

	permissions, err := client.GetUserRobotPermissions(testRobotShortname)
	if err != nil {
		t.Fatalf("GetUserRobotPermissions failed: %v", err)
	}

	if len(permissions.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions.Permissions))
	}
	if permissions.Permissions[0].Role != roleWrite {
		t.Errorf("Expected role 'write', got '%s'", permissions.Permissions[0].Role)
	}
}

func TestUserRobotErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Robot not found"}`))
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetUserRobotAccounts error
	_, err = client.GetUserRobotAccounts()
	if err == nil {
		t.Error("Expected error for failed request, got nil")
	}

	// Test GetUserRobotAccount error
	_, err = client.GetUserRobotAccount("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent robot, got nil")
	}

	// Test CreateUserRobotAccount error
	_, err = client.CreateUserRobotAccount("badbot", "desc", nil)
	if err == nil {
		t.Error("Expected error for failed create, got nil")
	}

	// Test DeleteUserRobotAccount error
	err = client.DeleteUserRobotAccount("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent robot, got nil")
	}

	// Test RegenerateUserRobotToken error
	_, err = client.RegenerateUserRobotToken("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent robot, got nil")
	}

	// Test GetUserRobotPermissions error
	_, err = client.GetUserRobotPermissions("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent robot, got nil")
	}
}
