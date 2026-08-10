package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Teams (org-level) ---

func TestGetTeams(t *testing.T) {
	mockResponse := struct {
		Teams []Team `json:"teams"`
	}{
		Teams: []Team{
			{Name: testTeamName, Description: testTeamDescDev, Role: roleMember, MemberCount: 5},
			{Name: "admins", Description: "Admin team", Role: roleAdmin, MemberCount: 2},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/teams"
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

	teams, err := client.GetTeams(testOrgName)
	if err != nil {
		t.Fatalf("GetTeams returned error: %v", err)
	}

	if len(teams) != 2 {
		t.Errorf("Expected 2 teams, got %d", len(teams))
	}
	if teams[0].Name != testTeamName {
		t.Errorf("Expected first team name %s, got %s", testTeamName, teams[0].Name)
	}
	if teams[1].Name != "admins" {
		t.Errorf("Expected second team name 'admins', got %s", teams[1].Name)
	}
}

func TestGetTeam(t *testing.T) {
	mockResponse := Team{
		Name:        testTeamName,
		Description: "Development team",
		Role:        roleMember,
		MemberCount: 5,
		RepoCount:   3,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName
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

	team, err := client.GetTeam(testOrgName, testTeamName)
	if err != nil {
		t.Fatalf("GetTeam returned error: %v", err)
	}

	if team.Name != testTeamName {
		t.Errorf("Expected team name %s, got %s", testTeamName, team.Name)
	}
	if team.Description != "Development team" {
		t.Errorf("Expected description 'Development team', got %s", team.Description)
	}
	if team.Role != roleMember {
		t.Errorf("Expected role %s, got %s", roleMember, team.Role)
	}
}

func TestCreateTeam(t *testing.T) {
	mockResponse := Team{
		Name:        testTeamName,
		Description: testTeamDescNew,
		Role:        roleMember,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName
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

	team, err := client.CreateTeam(testOrgName, testTeamName, testTeamDescNew, roleMember)
	if err != nil {
		t.Fatalf("CreateTeam returned error: %v", err)
	}

	if team.Name != testTeamName {
		t.Errorf("Expected team name %s, got %s", testTeamName, team.Name)
	}
	if team.Description != testTeamDescNew {
		t.Errorf("Expected description 'New team', got %s", team.Description)
	}
}

func TestUpdateTeam(t *testing.T) {
	mockResponse := Team{
		Name:        testTeamName,
		Description: updatedDescription,
		Role:        roleAdmin,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName
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

	team, err := client.UpdateTeam(testOrgName, testTeamName, updatedDescription, roleAdmin)
	if err != nil {
		t.Fatalf("UpdateTeam returned error: %v", err)
	}

	if team.Name != testTeamName {
		t.Errorf("Expected team name %s, got %s", testTeamName, team.Name)
	}
	if team.Description != updatedDescription {
		t.Errorf("Expected description '%s', got %s", updatedDescription, team.Description)
	}
	if team.Role != roleAdmin {
		t.Errorf("Expected role %s, got %s", roleAdmin, team.Role)
	}
}

func TestDeleteTeam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName
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

	err = client.DeleteTeam(testOrgName, testTeamName)
	if err != nil {
		t.Fatalf("DeleteTeam returned error: %v", err)
	}
}

// --- Team Members ---

func TestGetTeamMembers(t *testing.T) {
	mockResponse := TeamMembers{
		Members: []TeamMember{
			{Name: testUserName, Kind: testKindUser, IsRobot: false},
			{Name: "robot+builder", Kind: testKindRobot, IsRobot: true},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/members"
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

	members, err := client.GetTeamMembers(testOrgName, testTeamName)
	if err != nil {
		t.Fatalf("GetTeamMembers returned error: %v", err)
	}

	if len(members.Members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members.Members))
	}
	if members.Members[0].Name != testUserName {
		t.Errorf("Expected first member name %s, got %s", testUserName, members.Members[0].Name)
	}
	if members.Members[0].IsRobot {
		t.Errorf("Expected first member to not be a robot")
	}
	if !members.Members[1].IsRobot {
		t.Errorf("Expected second member to be a robot")
	}
}

func TestAddTeamMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/members/" + testUserName
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

	err = client.AddTeamMember(testOrgName, testTeamName, testUserName)
	if err != nil {
		t.Fatalf("AddTeamMember returned error: %v", err)
	}
}

func TestRemoveTeamMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/members/" + testUserName
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

	err = client.RemoveTeamMember(testOrgName, testTeamName, testUserName)
	if err != nil {
		t.Fatalf("RemoveTeamMember returned error: %v", err)
	}
}

// --- Team Permissions ---

func TestGetTeamPermissions(t *testing.T) {
	mockResponse := TeamPermissions{
		Permissions: []TeamPermission{
			{Repository: Repository{Name: testRepoName}, Role: testRoleRead},
			{Repository: Repository{Name: "another-repo"}, Role: testRoleWrite},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/permissions"
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

	perms, err := client.GetTeamPermissions(testOrgName, testTeamName)
	if err != nil {
		t.Fatalf("GetTeamPermissions returned error: %v", err)
	}

	if len(perms.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(perms.Permissions))
	}
	if perms.Permissions[0].Repository.Name != testRepoName {
		t.Errorf("Expected first repo name %s, got %s", testRepoName, perms.Permissions[0].Repository.Name)
	}
	if perms.Permissions[0].Role != testRoleRead {
		t.Errorf("Expected first role %s, got %s", testRoleRead, perms.Permissions[0].Role)
	}
}

func TestSetTeamRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/permissions/" + testRepoName
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

	err = client.SetTeamRepositoryPermission(testOrgName, testTeamName, testRepoName, testRoleWrite)
	if err != nil {
		t.Fatalf("SetTeamRepositoryPermission returned error: %v", err)
	}
}

func TestRemoveTeamRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/permissions/" + testRepoName
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

	err = client.RemoveTeamRepositoryPermission(testOrgName, testTeamName, testRepoName)
	if err != nil {
		t.Fatalf("RemoveTeamRepositoryPermission returned error: %v", err)
	}
}

// --- Team Invites ---

func TestInviteTeamMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/invite/" + testEmailAddress
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

	err = client.InviteTeamMember(testOrgName, testTeamName, testEmailAddress)
	if err != nil {
		t.Fatalf("InviteTeamMember returned error: %v", err)
	}
}

func TestDeleteTeamInvite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/invite/" + testEmailAddress
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

	err = client.DeleteTeamInvite(testOrgName, testTeamName, testEmailAddress)
	if err != nil {
		t.Fatalf("DeleteTeamInvite returned error: %v", err)
	}
}

// --- Error tests ---

func TestGetTeamsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetTeams(testOrgName)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetTeamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetTeam(testOrgName, "nonexistent-team")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestOrganizationTeamHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetTeams(testOrgName)
	if err == nil {
		t.Error("Expected error from GetTeams, got nil")
	}

	_, err = client.CreateTeam(testOrgName, testTeamName, testTeamDescDev, roleMember)
	if err == nil {
		t.Error("Expected error from CreateTeam, got nil")
	}

	_, err = client.GetTeam(testOrgName, testTeamName)
	if err == nil {
		t.Error("Expected error from GetTeam, got nil")
	}

	err = client.DeleteTeam(testOrgName, testTeamName)
	if err == nil {
		t.Error("Expected error from DeleteTeam, got nil")
	}

	_, err = client.UpdateTeam(testOrgName, testTeamName, testTeamDescDev, roleMember)
	if err == nil {
		t.Error("Expected error from UpdateTeam, got nil")
	}

	_, err = client.GetTeamMembers(testOrgName, testTeamName)
	if err == nil {
		t.Error("Expected error from GetTeamMembers, got nil")
	}

	err = client.AddTeamMember(testOrgName, testTeamName, testMemberName)
	if err == nil {
		t.Error("Expected error from AddTeamMember, got nil")
	}

	err = client.RemoveTeamMember(testOrgName, testTeamName, testMemberName)
	if err == nil {
		t.Error("Expected error from RemoveTeamMember, got nil")
	}

	_, err = client.GetTeamPermissions(testOrgName, testTeamName)
	if err == nil {
		t.Error("Expected error from GetTeamPermissions, got nil")
	}

	err = client.SetTeamRepositoryPermission(testOrgName, testTeamName, testRepository, testRoleRead)
	if err == nil {
		t.Error("Expected error from SetTeamRepositoryPermission, got nil")
	}

	err = client.RemoveTeamRepositoryPermission(testOrgName, testTeamName, testRepository)
	if err == nil {
		t.Error("Expected error from RemoveTeamRepositoryPermission, got nil")
	}

	err = client.InviteTeamMember(testOrgName, testTeamName, testEmailAddress)
	if err == nil {
		t.Error("Expected error from InviteTeamMember, got nil")
	}

	err = client.DeleteTeamInvite(testOrgName, testTeamName, testEmailAddress)
	if err == nil {
		t.Error("Expected error from DeleteTeamInvite, got nil")
	}
}
