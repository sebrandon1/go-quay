package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Robot Accounts ---

func TestGetRobotAccounts(t *testing.T) {
	mockRobots := RobotAccounts{
		Robots: []RobotAccount{
			{Name: testRobotFullName, Description: testRobotDescValue, Created: testTimestamp},
			{Name: "test-org+cibot", Description: "CI robot", Created: testTimestamp},
		},
	}
	mockResponseJSON := mustMarshal(t, mockRobots)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots"
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

	robots, err := client.GetRobotAccounts(context.Background(), testOrgName)
	if err != nil {
		t.Fatalf("GetRobotAccounts returned error: %v", err)
	}

	if len(robots.Robots) != 2 {
		t.Errorf("Expected 2 robots, got %d", len(robots.Robots))
	}
	if robots.Robots[0].Name != testRobotFullName {
		t.Errorf("Expected robot name %s, got %s", testRobotFullName, robots.Robots[0].Name)
	}
}

func TestCreateRobotAccount(t *testing.T) {
	mockRobot := RobotAccount{
		Name:        testOrgName + "+" + testRobotShortname,
		Description: testRobotDescValue,
		Token:       "new-robot-token",
		Created:     testTimestamp,
	}
	mockResponseJSON := mustMarshal(t, mockRobot)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req CreateRobotRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.Description != testRobotDescValue {
			t.Errorf("Expected description %s, got %s", testRobotDescValue, req.Description)
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

	robot, err := client.CreateRobotAccount(context.Background(), testOrgName, testRobotShortname, testRobotDescValue, nil)
	if err != nil {
		t.Fatalf("CreateRobotAccount returned error: %v", err)
	}

	if robot.Token != "new-robot-token" {
		t.Errorf("Expected token 'new-robot-token', got %s", robot.Token)
	}
}

func TestGetRobotAccount(t *testing.T) {
	mockRobot := RobotAccount{
		Name:        testOrgName + "+" + testRobotShortname,
		Description: testRobotDescValue,
		Token:       "robot-token",
		Created:     testTimestamp,
	}
	mockResponseJSON := mustMarshal(t, mockRobot)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname
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

	robot, err := client.GetRobotAccount(context.Background(), testOrgName, testRobotShortname)
	if err != nil {
		t.Fatalf("GetRobotAccount returned error: %v", err)
	}

	if robot.Description != testRobotDescValue {
		t.Errorf("Expected description %s, got %s", testRobotDescValue, robot.Description)
	}
}

func TestDeleteRobotAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname
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

	err = client.DeleteRobotAccount(context.Background(), testOrgName, testRobotShortname)
	if err != nil {
		t.Fatalf("DeleteRobotAccount returned error: %v", err)
	}
}

func TestRegenerateRobotToken(t *testing.T) {
	mockRobot := RobotAccount{
		Name:  testOrgName + "+" + testRobotShortname,
		Token: "regenerated-token",
	}
	mockResponseJSON := mustMarshal(t, mockRobot)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname + "/regenerate"
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

	robot, err := client.RegenerateRobotToken(context.Background(), testOrgName, testRobotShortname)
	if err != nil {
		t.Fatalf("RegenerateRobotToken returned error: %v", err)
	}

	if robot.Token != "regenerated-token" {
		t.Errorf("Expected token 'regenerated-token', got %s", robot.Token)
	}
}

// --- Robot Permissions ---

func TestGetRobotPermissions(t *testing.T) {
	mockPerms := RobotPermissions{
		Permissions: []RobotPermission{
			{Repository: Repository{Name: testRepoName}, Role: testRoleWrite},
			{Repository: Repository{Name: "other-repo"}, Role: testRoleRead},
		},
	}
	mockResponseJSON := mustMarshal(t, mockPerms)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname + "/permissions"
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

	perms, err := client.GetRobotPermissions(context.Background(), testOrgName, testRobotShortname)
	if err != nil {
		t.Fatalf("GetRobotPermissions returned error: %v", err)
	}

	if len(perms.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(perms.Permissions))
	}
	if perms.Permissions[0].Role != testRoleWrite {
		t.Errorf("Expected role %s, got %s", testRoleWrite, perms.Permissions[0].Role)
	}
}

func TestSetRobotRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname + "/permissions/" + testRepoName
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		var req SetRepositoryPermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.Role != testRoleWrite {
			t.Errorf("Expected role %s, got %s", testRoleWrite, req.Role)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.SetRobotRepositoryPermission(context.Background(), testOrgName, testRobotShortname, testRepoName, testRoleWrite)
	if err != nil {
		t.Fatalf("SetRobotRepositoryPermission returned error: %v", err)
	}
}

func TestRemoveRobotRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname + "/permissions/" + testRepoName
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

	err = client.RemoveRobotRepositoryPermission(context.Background(), testOrgName, testRobotShortname, testRepoName)
	if err != nil {
		t.Fatalf("RemoveRobotRepositoryPermission returned error: %v", err)
	}
}

// --- Robot Federation ---

func TestGetRobotFederation(t *testing.T) {
	mockFederation := RobotFederation{
		Federation: []RobotFederationConfig{
			{Issuer: testFederationIssuer, Subject: testFederationSubject},
		},
	}
	mockResponseJSON := mustMarshal(t, mockFederation)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname + "/federation"
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

	federation, err := client.GetRobotFederation(context.Background(), testOrgName, testRobotShortname)
	if err != nil {
		t.Fatalf("GetRobotFederation returned error: %v", err)
	}

	if len(federation.Federation) != 1 {
		t.Fatalf("Expected 1 federation config, got %d", len(federation.Federation))
	}
	if federation.Federation[0].Issuer != testFederationIssuer {
		t.Errorf("Expected issuer %s, got %s", testFederationIssuer, federation.Federation[0].Issuer)
	}
}

func TestCreateRobotFederation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname + "/federation"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	configs := []RobotFederationConfig{
		{Issuer: testFederationIssuer, Subject: testFederationSubject},
	}

	err = client.CreateRobotFederation(context.Background(), testOrgName, testRobotShortname, configs)
	if err != nil {
		t.Fatalf("CreateRobotFederation returned error: %v", err)
	}
}

func TestDeleteRobotFederation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/robots/" + testRobotShortname + "/federation"
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

	err = client.DeleteRobotFederation(context.Background(), testOrgName, testRobotShortname)
	if err != nil {
		t.Fatalf("DeleteRobotFederation returned error: %v", err)
	}
}

func TestOrganizationRobotHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetRobotAccounts(context.Background(), testOrgName)
	if err == nil {
		t.Error("Expected error from GetRobotAccounts, got nil")
	}

	_, err = client.CreateRobotAccount(context.Background(), testOrgName, "testbot", testRobotDescValue, nil)
	if err == nil {
		t.Error("Expected error from CreateRobotAccount, got nil")
	}

	_, err = client.GetRobotAccount(context.Background(), testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from GetRobotAccount, got nil")
	}

	err = client.DeleteRobotAccount(context.Background(), testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from DeleteRobotAccount, got nil")
	}

	_, err = client.RegenerateRobotToken(context.Background(), testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from RegenerateRobotToken, got nil")
	}

	_, err = client.GetRobotPermissions(context.Background(), testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from GetRobotPermissions, got nil")
	}

	err = client.SetRobotRepositoryPermission(context.Background(), testOrgName, "testbot", testRepository, testRoleRead)
	if err == nil {
		t.Error("Expected error from SetRobotRepositoryPermission, got nil")
	}

	err = client.RemoveRobotRepositoryPermission(context.Background(), testOrgName, "testbot", testRepository)
	if err == nil {
		t.Error("Expected error from RemoveRobotRepositoryPermission, got nil")
	}

	_, err = client.GetRobotFederation(context.Background(), testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from GetRobotFederation, got nil")
	}

	err = client.CreateRobotFederation(context.Background(), testOrgName, "testbot", []RobotFederationConfig{{Issuer: "https://example.com", Subject: testPlaceholder}})
	if err == nil {
		t.Error("Expected error from CreateRobotFederation, got nil")
	}

	err = client.DeleteRobotFederation(context.Background(), testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from DeleteRobotFederation, got nil")
	}
}
