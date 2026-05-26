package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRepositoryPermissions(t *testing.T) {
	mockPermissions := RepositoryPermissions{
		Permissions: []RepositoryPermission{
			{
				Name: testPermUserName,
				Kind: testKindUser,
				Role: testRoleWrite,
				Avatar: Avatar{
					Name: testPermUserName,
					Kind: testKindUser,
				},
				IsRobot:    false,
				IsOrgAdmin: false,
			},
			{
				Name: "testorg+deploybot",
				Kind: testKindRobot,
				Role: testRoleRead,
				Avatar: Avatar{
					Name: "deploybot",
					Kind: testKindRobot,
				},
				IsRobot:    true,
				IsOrgAdmin: false,
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockPermissions)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	permissions, err := client.GetRepositoryPermissions(testNamespace, testRepository)
	if err != nil {
		t.Fatalf("GetRepositoryPermissions failed: %v", err)
	}

	if len(permissions.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions.Permissions))
	}

	// Check first permission (user)
	userPerm := permissions.Permissions[0]
	if userPerm.Name != testPermUserName {
		t.Errorf("Expected user name 'john.doe', got '%s'", userPerm.Name)
	}
	if userPerm.Role != testRoleWrite {
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
	if robotPerm.Role != testRoleRead {
		t.Errorf("Expected role 'read', got '%s'", robotPerm.Role)
	}
	if robotPerm.IsRobot != true {
		t.Errorf("Expected IsRobot true for robot, got %v", robotPerm.IsRobot)
	}
}

func TestSetRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/permissions/"+testPermUserName {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/permissions/john.doe, got %s", r.URL.Path)
		}

		// Verify request body
		var req SetRepositoryPermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Role != testRoleWrite {
			t.Errorf("Expected role 'write', got '%s'", req.Role)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.SetRepositoryPermission(testNamespace, testRepository, testPermUserName, testRoleWrite)
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with valid roles
	validRoles := []string{testRoleRead, testRoleWrite, roleAdmin}
	for _, role := range validRoles {
		err = client.SetRepositoryPermission(testNamespace, testRepository, testKindUser, role)
		if err != nil {
			t.Errorf("SetRepositoryPermission failed for valid role '%s': %v", role, err)
		}
	}
}

func TestRemoveRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repository/testorg/testrepo/permissions/"+testPermUserName {
			t.Errorf("Expected path /api/v1/repository/testorg/testrepo/permissions/john.doe, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.RemoveRepositoryPermission(testNamespace, testRepository, testPermUserName)
	if err != nil {
		t.Fatalf("RemoveRepositoryPermission failed: %v", err)
	}
}

func TestListUserPermissions(t *testing.T) {
	mockPermissions := RepositoryPermissions{
		Permissions: []RepositoryPermission{
			{Name: testPermUserName, Kind: testKindUser, Role: testRoleWrite},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockPermissions)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/user/"
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

	perms, err := client.ListUserPermissions(testNamespace, testRepository)
	if err != nil {
		t.Fatalf("ListUserPermissions returned error: %v", err)
	}

	if len(perms.Permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(perms.Permissions))
	}
	if perms.Permissions[0].Name != testPermUserName {
		t.Errorf("Expected user name %s, got %s", testPermUserName, perms.Permissions[0].Name)
	}
}

func TestGetUserPermission(t *testing.T) {
	mockPermission := RepositoryPermission{
		Name: testPermUserName,
		Kind: testKindUser,
		Role: testRoleWrite,
	}
	mockResponseJSON, _ := json.Marshal(mockPermission)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/user/" + testPermUserName
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

	perm, err := client.GetUserPermission(testNamespace, testRepository, testPermUserName)
	if err != nil {
		t.Fatalf("GetUserPermission returned error: %v", err)
	}

	if perm.Role != testRoleWrite {
		t.Errorf("Expected role %s, got %s", testRoleWrite, perm.Role)
	}
}

func TestSetUserPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/user/" + testPermUserName
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.SetUserPermission(testNamespace, testRepository, testPermUserName, testRoleWrite)
	if err != nil {
		t.Fatalf("SetUserPermission returned error: %v", err)
	}
}

func TestDeleteUserPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/user/" + testPermUserName
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

	err = client.DeleteUserPermission(testNamespace, testRepository, testPermUserName)
	if err != nil {
		t.Fatalf("DeleteUserPermission returned error: %v", err)
	}
}

func TestGetUserTransitivePermission(t *testing.T) {
	mockPermission := RepositoryPermission{
		Name: testPermUserName,
		Kind: testKindUser,
		Role: roleAdmin,
	}
	mockResponseJSON, _ := json.Marshal(mockPermission)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/user/" + testPermUserName + "/transitive"
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

	perm, err := client.GetUserTransitivePermission(testNamespace, testRepository, testPermUserName)
	if err != nil {
		t.Fatalf("GetUserTransitivePermission returned error: %v", err)
	}

	if perm.Role != roleAdmin {
		t.Errorf("Expected role '%s', got %s", roleAdmin, perm.Role)
	}
}

func TestListTeamPermissions(t *testing.T) {
	mockPermissions := RepositoryPermissions{
		Permissions: []RepositoryPermission{
			{Name: testPrototypeTeamName, Kind: testKindTeam, Role: testRoleRead},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockPermissions)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/team/"
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

	perms, err := client.ListTeamPermissions(testNamespace, testRepository)
	if err != nil {
		t.Fatalf("ListTeamPermissions returned error: %v", err)
	}

	if len(perms.Permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(perms.Permissions))
	}
	if perms.Permissions[0].Name != testPrototypeTeamName {
		t.Errorf("Expected team name %s, got %s", testPrototypeTeamName, perms.Permissions[0].Name)
	}
}

func TestGetTeamPermission(t *testing.T) {
	mockPermission := RepositoryPermission{
		Name: testPrototypeTeamName,
		Kind: testKindTeam,
		Role: testRoleRead,
	}
	mockResponseJSON, _ := json.Marshal(mockPermission)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/team/" + testPrototypeTeamName
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

	perm, err := client.GetTeamPermission(testNamespace, testRepository, testPrototypeTeamName)
	if err != nil {
		t.Fatalf("GetTeamPermission returned error: %v", err)
	}

	if perm.Role != testRoleRead {
		t.Errorf("Expected role %s, got %s", testRoleRead, perm.Role)
	}
}

func TestSetTeamPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/team/" + testPrototypeTeamName
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.SetTeamPermission(testNamespace, testRepository, testPrototypeTeamName, testRoleWrite)
	if err != nil {
		t.Fatalf("SetTeamPermission returned error: %v", err)
	}
}

func TestDeleteTeamPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/permissions/team/" + testPrototypeTeamName
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

	err = client.DeleteTeamPermission(testNamespace, testRepository, testPrototypeTeamName)
	if err != nil {
		t.Fatalf("DeleteTeamPermission returned error: %v", err)
	}
}

func TestPermissionsErrorHandling(t *testing.T) {
	// Test 404 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Repository not found"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetRepositoryPermissions error
	_, err = client.GetRepositoryPermissions(testNamespace, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}

	// Test SetRepositoryPermission error
	err = client.SetRepositoryPermission(testNamespace, "nonexistent", testKindUser, testRoleRead)
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}

	// Test RemoveRepositoryPermission error
	err = client.RemoveRepositoryPermission(testNamespace, "nonexistent", testKindUser)
	if err == nil {
		t.Error("Expected error for non-existent repository, got nil")
	}

	// Test ListUserPermissions error
	_, err = client.ListUserPermissions(testNamespace, "nonexistent")
	if err == nil {
		t.Error("Expected error from ListUserPermissions, got nil")
	}

	// Test GetUserPermission error
	_, err = client.GetUserPermission(testNamespace, "nonexistent", testPermUserName)
	if err == nil {
		t.Error("Expected error from GetUserPermission, got nil")
	}

	// Test SetUserPermission error
	err = client.SetUserPermission(testNamespace, "nonexistent", testPermUserName, testRoleRead)
	if err == nil {
		t.Error("Expected error from SetUserPermission, got nil")
	}

	// Test DeleteUserPermission error
	err = client.DeleteUserPermission(testNamespace, "nonexistent", testPermUserName)
	if err == nil {
		t.Error("Expected error from DeleteUserPermission, got nil")
	}

	// Test GetUserTransitivePermission error
	_, err = client.GetUserTransitivePermission(testNamespace, "nonexistent", testPermUserName)
	if err == nil {
		t.Error("Expected error from GetUserTransitivePermission, got nil")
	}

	// Test ListTeamPermissions error
	_, err = client.ListTeamPermissions(testNamespace, "nonexistent")
	if err == nil {
		t.Error("Expected error from ListTeamPermissions, got nil")
	}

	// Test GetTeamPermission error
	_, err = client.GetTeamPermission(testNamespace, "nonexistent", testTeamName)
	if err == nil {
		t.Error("Expected error from GetTeamPermission, got nil")
	}

	// Test SetTeamPermission error
	err = client.SetTeamPermission(testNamespace, "nonexistent", testTeamName, testRoleRead)
	if err == nil {
		t.Error("Expected error from SetTeamPermission, got nil")
	}

	// Test DeleteTeamPermission error
	err = client.DeleteTeamPermission(testNamespace, "nonexistent", testTeamName)
	if err == nil {
		t.Error("Expected error from DeleteTeamPermission, got nil")
	}
}
