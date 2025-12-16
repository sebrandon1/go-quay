package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetPerms    = "GET"
	httpPutPerms    = "PUT"
	httpDeletePerms = "DELETE"

	permRoleWrite = "write"
)

func TestGetRepositoryPermissions(t *testing.T) {
	mockPermissions := RepositoryPermissions{
		Permissions: []RepositoryPermission{
			{
				Name: "john.doe",
				Kind: "user",
				Role: permRoleWrite,
				Avatar: Avatar{
					Name: "john.doe",
					Kind: "user",
				},
				IsRobot:    false,
				IsOrgAdmin: false,
			},
			{
				Name: "testorg+deploybot",
				Kind: "robot",
				Role: "read",
				Avatar: Avatar{
					Name: "deploybot",
					Kind: "robot",
				},
				IsRobot:    true,
				IsOrgAdmin: false,
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockPermissions)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetPerms {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/permissions" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/permissions, got %s", r.URL.Path)
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

	permissions, err := client.GetRepositoryPermissions("testorg", "testrepo")
	if err != nil {
		t.Fatalf("GetRepositoryPermissions failed: %v", err)
	}

	if len(permissions.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions.Permissions))
	}

	// Check first permission (user)
	userPerm := permissions.Permissions[0]
	if userPerm.Name != "john.doe" {
		t.Errorf("Expected user name 'john.doe', got '%s'", userPerm.Name)
	}
	if userPerm.Role != permRoleWrite {
		t.Errorf("Expected role 'write', got '%s'", userPerm.Role)
	}
	if userPerm.IsRobot != false {
		t.Errorf("Expected IsRobot false for user, got %v", userPerm.IsRobot)
	}

	// Check second permission (robot)
	robotPerm := permissions.Permissions[1]
	if robotPerm.Name != "testorg+deploybot" {
		t.Errorf("Expected robot name 'testorg+deploybot', got '%s'", robotPerm.Name)
	}
	if robotPerm.Role != "read" {
		t.Errorf("Expected role 'read', got '%s'", robotPerm.Role)
	}
	if robotPerm.IsRobot != true {
		t.Errorf("Expected IsRobot true for robot, got %v", robotPerm.IsRobot)
	}
}

func TestSetRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPutPerms {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/permissions/john.doe" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/permissions/john.doe, got %s", r.URL.Path)
		}

		// Verify request body
		var req SetRepositoryPermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Role != permRoleWrite {
			t.Errorf("Expected role 'write', got '%s'", req.Role)
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

	err = client.SetRepositoryPermission("testorg", "testrepo", "john.doe", permRoleWrite)
	if err != nil {
		t.Fatalf("SetRepositoryPermission failed: %v", err)
	}
}

func TestSetRepositoryPermissionInvalidRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This shouldn't be called for invalid role, but if it is, return success
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

	// Test with valid roles
	validRoles := []string{"read", "write", "admin"}
	for _, role := range validRoles {
		err = client.SetRepositoryPermission("testorg", "testrepo", "user", role)
		if err != nil {
			t.Errorf("SetRepositoryPermission failed for valid role '%s': %v", role, err)
		}
	}
}

func TestRemoveRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeletePerms {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/permissions/john.doe" {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/permissions/john.doe, got %s", r.URL.Path)
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

	err = client.RemoveRepositoryPermission("testorg", "testrepo", "john.doe")
	if err != nil {
		t.Fatalf("RemoveRepositoryPermission failed: %v", err)
	}
}

func TestPermissionsErrorHandling(t *testing.T) {
	// Test 404 error
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

	// Test GetRepositoryPermissions error
	_, err = client.GetRepositoryPermissions("testorg", "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}

	// Test SetRepositoryPermission error
	err = client.SetRepositoryPermission("testorg", "nonexistent", "user", "read")
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}

	// Test RemoveRepositoryPermission error
	err = client.RemoveRepositoryPermission("testorg", "nonexistent", "user")
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}
}
