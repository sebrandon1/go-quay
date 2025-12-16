package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetTeam    = "GET"
	httpPutTeam    = "PUT"
	httpDeleteTeam = "DELETE"

	testOrgName    = "test-org"
	testTeamName   = "developers"
	testMemberName = "testuser"
	testRepoName   = "test-repo"
	roleAdmin      = "admin"
	roleMember     = "member"
	roleRead       = "read"
)

func TestGetTeams(t *testing.T) {
	mockResponse := struct {
		Teams []Team `json:"teams"`
	}{
		Teams: []Team{
			{Name: testTeamName, Description: "Dev team", Role: roleMember, MemberCount: 5},
			{Name: "admins", Description: "Admin team", Role: roleAdmin, MemberCount: 2},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetTeam {
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

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
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
		if r.Method != httpGetTeam {
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

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
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
		Description: "New team",
		Role:        roleMember,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPutTeam {
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

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	team, err := client.CreateTeam(testOrgName, testTeamName, "New team", roleMember)
	if err != nil {
		t.Fatalf("CreateTeam returned error: %v", err)
	}

	if team.Name != testTeamName {
		t.Errorf("Expected team name %s, got %s", testTeamName, team.Name)
	}
	if team.Description != "New team" {
		t.Errorf("Expected description 'New team', got %s", team.Description)
	}
}

func TestUpdateTeam(t *testing.T) {
	mockResponse := Team{
		Name:        testTeamName,
		Description: "Updated description",
		Role:        roleAdmin,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPutTeam {
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

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	team, err := client.UpdateTeam(testOrgName, testTeamName, "Updated description", roleAdmin)
	if err != nil {
		t.Fatalf("UpdateTeam returned error: %v", err)
	}

	if team.Name != testTeamName {
		t.Errorf("Expected team name %s, got %s", testTeamName, team.Name)
	}
	if team.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got %s", team.Description)
	}
	if team.Role != roleAdmin {
		t.Errorf("Expected role %s, got %s", roleAdmin, team.Role)
	}
}

func TestDeleteTeam(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteTeam {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName
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

	err = client.DeleteTeam(testOrgName, testTeamName)
	if err != nil {
		t.Fatalf("DeleteTeam returned error: %v", err)
	}
}

func TestGetTeamMembers(t *testing.T) {
	mockResponse := TeamMembers{
		Members: []TeamMember{
			{Name: testMemberName, Kind: "user", IsRobot: false},
			{Name: "robot+builder", Kind: "robot", IsRobot: true},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetTeam {
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

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
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
	if members.Members[0].Name != testMemberName {
		t.Errorf("Expected first member name %s, got %s", testMemberName, members.Members[0].Name)
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
		if r.Method != httpPutTeam {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/members/" + testMemberName
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
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

	err = client.AddTeamMember(testOrgName, testTeamName, testMemberName)
	if err != nil {
		t.Fatalf("AddTeamMember returned error: %v", err)
	}
}

func TestRemoveTeamMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteTeam {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/members/" + testMemberName
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

	err = client.RemoveTeamMember(testOrgName, testTeamName, testMemberName)
	if err != nil {
		t.Fatalf("RemoveTeamMember returned error: %v", err)
	}
}

func TestGetTeamPermissions(t *testing.T) {
	mockResponse := TeamPermissions{
		Permissions: []TeamPermission{
			{Repository: Repository{Name: testRepoName}, Role: roleRead},
			{Repository: Repository{Name: "another-repo"}, Role: "write"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetTeam {
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

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
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
	if perms.Permissions[0].Role != roleRead {
		t.Errorf("Expected first role %s, got %s", roleRead, perms.Permissions[0].Role)
	}
}

func TestSetTeamRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPutTeam {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/permissions/" + testRepoName
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
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

	err = client.SetTeamRepositoryPermission(testOrgName, testTeamName, testRepoName, "write")
	if err != nil {
		t.Fatalf("SetTeamRepositoryPermission returned error: %v", err)
	}
}

func TestRemoveTeamRepositoryPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteTeam {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/team/" + testTeamName + "/permissions/" + testRepoName
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

	err = client.RemoveTeamRepositoryPermission(testOrgName, testTeamName, testRepoName)
	if err != nil {
		t.Fatalf("RemoveTeamRepositoryPermission returned error: %v", err)
	}
}

func TestGetTeamsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
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

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetTeam(testOrgName, "nonexistent-team")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
