package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testClientID       = "client-abc123"
	testPolicyUUID     = "policy-uuid-456"
	testAppName        = "Test Application"
	testAppURI         = "https://app.example.com"
	testRedirectURI    = "https://app.example.com/callback"
	testUpstreamReg    = "docker.io"
	testPrototypeID    = "proto-123"
	testSubscriptionID = "sub-789"
	testMemberName     = "johndoe"
)

// --- Organization CRUD ---

func TestGetOrganization(t *testing.T) {
	mockOrg := Organization{
		Name:          testOrgName,
		Email:         testEmailAddress,
		IsOrgAdmin:    true,
		CanCreateRepo: true,
	}
	mockResponseJSON, _ := json.Marshal(mockOrg)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName
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

	org, err := client.GetOrganization(testOrgName)
	if err != nil {
		t.Fatalf("GetOrganization returned error: %v", err)
	}

	if org.Name != testOrgName {
		t.Errorf("Expected org name %s, got %s", testOrgName, org.Name)
	}
	if org.Email != testEmailAddress {
		t.Errorf("Expected email %s, got %s", testEmailAddress, org.Email)
	}
	if !org.IsOrgAdmin {
		t.Errorf("Expected IsOrgAdmin to be true")
	}
}

func TestGetOrganizationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetOrganization(testOrgName)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCreateOrganization(t *testing.T) {
	mockOrg := Organization{
		Name:  testOrgName,
		Email: testEmailAddress,
	}
	mockResponseJSON, _ := json.Marshal(mockOrg)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/organization/" {
			t.Errorf("Expected path /api/v1/organization/, got %s", r.URL.Path)
		}

		var req CreateOrganizationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.Name != testOrgName {
			t.Errorf("Expected org name %s, got %s", testOrgName, req.Name)
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

	org, err := client.CreateOrganization(testOrgName, testEmailAddress)
	if err != nil {
		t.Fatalf("CreateOrganization returned error: %v", err)
	}

	if org.Name != testOrgName {
		t.Errorf("Expected org name %s, got %s", testOrgName, org.Name)
	}
}

func TestUpdateOrganization(t *testing.T) {
	updatedEmail := "updated@example.com"
	mockOrg := Organization{
		Name:  testOrgName,
		Email: updatedEmail,
	}
	mockResponseJSON, _ := json.Marshal(mockOrg)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName
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

	org, err := client.UpdateOrganization(testOrgName, updatedEmail)
	if err != nil {
		t.Fatalf("UpdateOrganization returned error: %v", err)
	}

	if org.Email != updatedEmail {
		t.Errorf("Expected email %s, got %s", updatedEmail, org.Email)
	}
}

func TestDeleteOrganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName
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

	err = client.DeleteOrganization(testOrgName)
	if err != nil {
		t.Fatalf("DeleteOrganization returned error: %v", err)
	}
}

// --- Organization Members ---

func TestGetOrganizationMembers(t *testing.T) {
	mockMembers := OrganizationMembers{
		Members: []OrganizationMember{
			{Name: testMemberName, Kind: testKindUser},
			{Name: "janedoe", Kind: testKindUser},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockMembers)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/members"
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

	members, err := client.GetOrganizationMembers(testOrgName)
	if err != nil {
		t.Fatalf("GetOrganizationMembers returned error: %v", err)
	}

	if len(members.Members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members.Members))
	}
	if members.Members[0].Name != testMemberName {
		t.Errorf("Expected member name %s, got %s", testMemberName, members.Members[0].Name)
	}
}

func TestGetOrganizationMembersError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetOrganizationMembers(testOrgName)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddOrganizationMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/members/" + testMemberName
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

	err = client.AddOrganizationMember(testOrgName, testMemberName)
	if err != nil {
		t.Fatalf("AddOrganizationMember returned error: %v", err)
	}
}

func TestRemoveOrganizationMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/members/" + testMemberName
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

	err = client.RemoveOrganizationMember(testOrgName, testMemberName)
	if err != nil {
		t.Fatalf("RemoveOrganizationMember returned error: %v", err)
	}
}

func TestGetOrganizationMember(t *testing.T) {
	mockMember := OrganizationMember{
		Name: testMemberName,
		Kind: testKindUser,
		Teams: []Team{
			{Name: testTeamName, Role: roleMember},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockMember)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/members/" + testMemberName
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

	member, err := client.GetOrganizationMember(testOrgName, testMemberName)
	if err != nil {
		t.Fatalf("GetOrganizationMember returned error: %v", err)
	}

	if member.Name != testMemberName {
		t.Errorf("Expected member name %s, got %s", testMemberName, member.Name)
	}
	if len(member.Teams) != 1 {
		t.Errorf("Expected 1 team, got %d", len(member.Teams))
	}
}

// --- Organization Repositories ---

func TestGetOrganizationRepositories(t *testing.T) {
	mockRepos := OrganizationRepositories{
		Repositories: []OrganizationRepository{
			{Name: testRepoName, Namespace: testOrgName, IsPublic: true},
			{Name: "private-repo", Namespace: testOrgName, IsPublic: false},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockRepos)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/repositories"
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

	repos, err := client.GetOrganizationRepositories(testOrgName)
	if err != nil {
		t.Fatalf("GetOrganizationRepositories returned error: %v", err)
	}

	if len(repos.Repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(repos.Repositories))
	}
	if repos.Repositories[0].Name != testRepoName {
		t.Errorf("Expected repo name %s, got %s", testRepoName, repos.Repositories[0].Name)
	}
}

// --- Teams (org-level, calling organization.go functions) ---

func TestGetTeamsOrg(t *testing.T) {
	mockResponse := struct {
		Teams []Team `json:"teams"`
	}{
		Teams: []Team{
			{Name: testTeamName, Description: testTeamDescDev, Role: roleMember},
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

	if len(teams) != 1 {
		t.Errorf("Expected 1 team, got %d", len(teams))
	}
}

func TestCreateTeamOrg(t *testing.T) {
	mockTeam := Team{
		Name:        testTeamName,
		Description: testTeamDescNew,
		Role:        roleMember,
	}
	mockResponseJSON, _ := json.Marshal(mockTeam)

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
}

func TestGetTeamOrg(t *testing.T) {
	mockTeam := Team{
		Name:        testTeamName,
		Description: testTeamDescDev,
		Role:        roleMember,
		MemberCount: 5,
	}
	mockResponseJSON, _ := json.Marshal(mockTeam)

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

	if team.MemberCount != 5 {
		t.Errorf("Expected member count 5, got %d", team.MemberCount)
	}
}

func TestUpdateTeamOrg(t *testing.T) {
	mockTeam := Team{
		Name:        testTeamName,
		Description: "Updated desc",
		Role:        roleAdmin,
	}
	mockResponseJSON, _ := json.Marshal(mockTeam)

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

	team, err := client.UpdateTeam(testOrgName, testTeamName, "Updated desc", roleAdmin)
	if err != nil {
		t.Fatalf("UpdateTeam returned error: %v", err)
	}

	if team.Role != roleAdmin {
		t.Errorf("Expected role %s, got %s", roleAdmin, team.Role)
	}
}

func TestDeleteTeamOrg(t *testing.T) {
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

// --- Robot Accounts ---

func TestGetRobotAccounts(t *testing.T) {
	mockRobots := RobotAccounts{
		Robots: []RobotAccount{
			{Name: testRobotFullName, Description: testRobotDescValue, Created: testTimestamp},
			{Name: "test-org+cibot", Description: "CI robot", Created: testTimestamp},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockRobots)

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

	robots, err := client.GetRobotAccounts(testOrgName)
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
	mockResponseJSON, _ := json.Marshal(mockRobot)

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

	robot, err := client.CreateRobotAccount(testOrgName, testRobotShortname, testRobotDescValue, nil)
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
	mockResponseJSON, _ := json.Marshal(mockRobot)

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

	robot, err := client.GetRobotAccount(testOrgName, testRobotShortname)
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

	err = client.DeleteRobotAccount(testOrgName, testRobotShortname)
	if err != nil {
		t.Fatalf("DeleteRobotAccount returned error: %v", err)
	}
}

func TestRegenerateRobotToken(t *testing.T) {
	mockRobot := RobotAccount{
		Name:  testOrgName + "+" + testRobotShortname,
		Token: "regenerated-token",
	}
	mockResponseJSON, _ := json.Marshal(mockRobot)

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

	robot, err := client.RegenerateRobotToken(testOrgName, testRobotShortname)
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
	mockResponseJSON, _ := json.Marshal(mockPerms)

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

	perms, err := client.GetRobotPermissions(testOrgName, testRobotShortname)
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
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.SetRobotRepositoryPermission(testOrgName, testRobotShortname, testRepoName, testRoleWrite)
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

	err = client.RemoveRobotRepositoryPermission(testOrgName, testRobotShortname, testRepoName)
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
	mockResponseJSON, _ := json.Marshal(mockFederation)

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

	federation, err := client.GetRobotFederation(testOrgName, testRobotShortname)
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

	err = client.CreateRobotFederation(testOrgName, testRobotShortname, configs)
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

	err = client.DeleteRobotFederation(testOrgName, testRobotShortname)
	if err != nil {
		t.Fatalf("DeleteRobotFederation returned error: %v", err)
	}
}

// --- Quota ---

func TestGetQuota(t *testing.T) {
	mockQuota := Quota{
		ID:         "quota-1",
		LimitBytes: 1073741824,
	}
	mockResponseJSON, _ := json.Marshal(mockQuota)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/quota"
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

	quota, err := client.GetQuota(testOrgName)
	if err != nil {
		t.Fatalf("GetQuota returned error: %v", err)
	}

	if quota.LimitBytes != 1073741824 {
		t.Errorf("Expected limit bytes 1073741824, got %d", quota.LimitBytes)
	}
}

func TestCreateQuota(t *testing.T) {
	var limitBytes int64 = 2147483648
	mockQuota := Quota{
		ID:         "quota-2",
		LimitBytes: limitBytes,
	}
	mockResponseJSON, _ := json.Marshal(mockQuota)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/quota"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req CreateQuotaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.LimitBytes != limitBytes {
			t.Errorf("Expected limit bytes %d, got %d", limitBytes, req.LimitBytes)
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

	quota, err := client.CreateQuota(testOrgName, limitBytes)
	if err != nil {
		t.Fatalf("CreateQuota returned error: %v", err)
	}

	if quota.LimitBytes != limitBytes {
		t.Errorf("Expected limit bytes %d, got %d", limitBytes, quota.LimitBytes)
	}
}

func TestDeleteQuota(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/quota"
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

	err = client.DeleteQuota(testOrgName)
	if err != nil {
		t.Fatalf("DeleteQuota returned error: %v", err)
	}
}

// --- Auto-Prune ---

func TestGetAutoPrunePolicies(t *testing.T) {
	mockPolicies := AutoPrunePolicies{
		Policies: []AutoPrunePolicy{
			{UUID: testPolicyUUID, Method: testAutoPruneMethodNumberOfTags, Value: 10, TagPattern: "v*"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockPolicies)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy"
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

	policies, err := client.GetAutoPrunePolicies(testOrgName)
	if err != nil {
		t.Fatalf("GetAutoPrunePolicies returned error: %v", err)
	}

	if len(policies.Policies) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(policies.Policies))
	}
	if policies.Policies[0].UUID != testPolicyUUID {
		t.Errorf("Expected policy UUID %s, got %s", testPolicyUUID, policies.Policies[0].UUID)
	}
}

func TestCreateAutoPrunePolicy(t *testing.T) {
	mockPolicy := AutoPrunePolicy{
		UUID:       testPolicyUUID,
		Method:     testAutoPruneMethodNumberOfTags,
		Value:      20,
		TagPattern: testTagPatternRelease,
	}
	mockResponseJSON, _ := json.Marshal(mockPolicy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req CreateAutoPruneRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.Method != testAutoPruneMethodNumberOfTags {
			t.Errorf("Expected method 'number_of_tags', got %s", req.Method)
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

	policy, err := client.CreateAutoPrunePolicy(testOrgName, testAutoPruneMethodNumberOfTags, 20, testTagPatternRelease)
	if err != nil {
		t.Fatalf("CreateAutoPrunePolicy returned error: %v", err)
	}

	if policy.Value != 20 {
		t.Errorf("Expected value 20, got %d", policy.Value)
	}
}

func TestDeleteAutoPrunePolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy/" + testPolicyUUID
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

	err = client.DeleteAutoPrunePolicy(testOrgName, testPolicyUUID)
	if err != nil {
		t.Fatalf("DeleteAutoPrunePolicy returned error: %v", err)
	}
}

// --- Applications ---

func TestGetApplications(t *testing.T) {
	mockApps := Applications{
		Applications: []Application{
			{ClientID: testClientID, Name: testAppName, Description: "A test app"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockApps)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/applications"
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

	apps, err := client.GetApplications(testOrgName)
	if err != nil {
		t.Fatalf("GetApplications returned error: %v", err)
	}

	if len(apps.Applications) != 1 {
		t.Errorf("Expected 1 application, got %d", len(apps.Applications))
	}
	if apps.Applications[0].ClientID != testClientID {
		t.Errorf("Expected client ID %s, got %s", testClientID, apps.Applications[0].ClientID)
	}
}

func TestCreateApplication(t *testing.T) {
	mockApp := Application{
		ClientID:       testClientID,
		Name:           testAppName,
		Description:    "New application",
		ApplicationURI: testAppURI,
		RedirectURI:    testRedirectURI,
	}
	mockResponseJSON, _ := json.Marshal(mockApp)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/applications"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req CreateApplicationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.Name != testAppName {
			t.Errorf("Expected app name %s, got %s", testAppName, req.Name)
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

	app, err := client.CreateApplication(testOrgName, testAppName, "New application", testAppURI, testRedirectURI)
	if err != nil {
		t.Fatalf("CreateApplication returned error: %v", err)
	}

	if app.Name != testAppName {
		t.Errorf("Expected app name %s, got %s", testAppName, app.Name)
	}
	if app.RedirectURI != testRedirectURI {
		t.Errorf("Expected redirect URI %s, got %s", testRedirectURI, app.RedirectURI)
	}
}

// --- Proxy Cache ---

func TestGetProxyCacheConfig(t *testing.T) {
	mockConfig := ProxyCacheConfig{
		UpstreamRegistry: testUpstreamReg,
		Insecure:         false,
		Expiration:       86400,
	}
	mockResponseJSON, _ := json.Marshal(mockConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/proxycache"
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

	config, err := client.GetProxyCacheConfig(testOrgName)
	if err != nil {
		t.Fatalf("GetProxyCacheConfig returned error: %v", err)
	}

	if config.UpstreamRegistry != testUpstreamReg {
		t.Errorf("Expected upstream registry %s, got %s", testUpstreamReg, config.UpstreamRegistry)
	}
	if config.Expiration != 86400 {
		t.Errorf("Expected expiration 86400, got %d", config.Expiration)
	}
}

func TestCreateProxyCacheConfig(t *testing.T) {
	mockConfig := ProxyCacheConfig{
		UpstreamRegistry: testUpstreamReg,
		Insecure:         true,
		Expiration:       3600,
	}
	mockResponseJSON, _ := json.Marshal(mockConfig)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/proxycache"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
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

	config, err := client.CreateProxyCacheConfig(testOrgName, testUpstreamReg, true, 3600)
	if err != nil {
		t.Fatalf("CreateProxyCacheConfig returned error: %v", err)
	}

	if !config.Insecure {
		t.Errorf("Expected Insecure to be true")
	}
}

func TestDeleteProxyCacheConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/proxycache"
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

	err = client.DeleteProxyCacheConfig(testOrgName)
	if err != nil {
		t.Fatalf("DeleteProxyCacheConfig returned error: %v", err)
	}
}

// --- Default Permissions ---

func TestGetDefaultPermissions(t *testing.T) {
	mockPerms := DefaultPermissions{
		Prototypes: []DefaultPermission{
			{
				ID:   testPrototypeID,
				Role: testRoleRead,
				Delegate: User{
					Name: testMemberName,
					Kind: testKindUser,
				},
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockPerms)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/prototypes"
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

	perms, err := client.GetDefaultPermissions(testOrgName)
	if err != nil {
		t.Fatalf("GetDefaultPermissions returned error: %v", err)
	}

	if len(perms.Prototypes) != 1 {
		t.Errorf("Expected 1 prototype, got %d", len(perms.Prototypes))
	}
	if perms.Prototypes[0].Role != testRoleRead {
		t.Errorf("Expected role %s, got %s", testRoleRead, perms.Prototypes[0].Role)
	}
}

// --- Default Permissions (Create/Delete) ---

func TestCreateDefaultPermission(t *testing.T) {
	mockPerm := DefaultPermission{
		ID:   testPrototypeID,
		Role: testRoleRead,
		Delegate: User{
			Name: testMemberName,
			Kind: testKindUser,
		},
	}
	mockResponseJSON, _ := json.Marshal(mockPerm)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/prototypes"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
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

	perm, err := client.CreateDefaultPermission(testOrgName, testRoleRead, testKindUser, testMemberName)
	if err != nil {
		t.Fatalf("CreateDefaultPermission returned error: %v", err)
	}

	if perm.ID != testPrototypeID {
		t.Errorf("Expected prototype ID %s, got %s", testPrototypeID, perm.ID)
	}
	if perm.Role != testRoleRead {
		t.Errorf("Expected role %s, got %s", testRoleRead, perm.Role)
	}
}

func TestDeleteDefaultPermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/prototypes/" + testPrototypeID
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

	err = client.DeleteDefaultPermission(testOrgName, testPrototypeID)
	if err != nil {
		t.Fatalf("DeleteDefaultPermission returned error: %v", err)
	}
}

// --- Quota (Update) ---

func TestUpdateQuota(t *testing.T) {
	var limitBytes int64 = 4294967296
	mockQuota := Quota{
		ID:         "quota-1",
		LimitBytes: limitBytes,
	}
	mockResponseJSON, _ := json.Marshal(mockQuota)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/quota"
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

	quota, err := client.UpdateQuota(testOrgName, limitBytes)
	if err != nil {
		t.Fatalf("UpdateQuota returned error: %v", err)
	}

	if quota.LimitBytes != limitBytes {
		t.Errorf("Expected limit bytes %d, got %d", limitBytes, quota.LimitBytes)
	}
}

// --- Auto-Prune (Get/Update single policy) ---

func TestGetAutoPrunePolicy(t *testing.T) {
	mockPolicy := AutoPrunePolicy{
		UUID:       testPolicyUUID,
		Method:     testAutoPruneMethodNumberOfTags,
		Value:      10,
		TagPattern: "v*",
	}
	mockResponseJSON, _ := json.Marshal(mockPolicy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy/" + testPolicyUUID
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

	policy, err := client.GetAutoPrunePolicy(testOrgName, testPolicyUUID)
	if err != nil {
		t.Fatalf("GetAutoPrunePolicy returned error: %v", err)
	}

	if policy.UUID != testPolicyUUID {
		t.Errorf("Expected policy UUID %s, got %s", testPolicyUUID, policy.UUID)
	}
	if policy.Value != 10 {
		t.Errorf("Expected value 10, got %d", policy.Value)
	}
}

func TestUpdateAutoPrunePolicy(t *testing.T) {
	mockPolicy := AutoPrunePolicy{
		UUID:       testPolicyUUID,
		Method:     testAutoPruneMethodNumberOfTags,
		Value:      30,
		TagPattern: testTagPatternRelease,
	}
	mockResponseJSON, _ := json.Marshal(mockPolicy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy/" + testPolicyUUID
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

	policy, err := client.UpdateAutoPrunePolicy(testOrgName, testPolicyUUID, testAutoPruneMethodNumberOfTags, 30, testTagPatternRelease)
	if err != nil {
		t.Fatalf("UpdateAutoPrunePolicy returned error: %v", err)
	}

	if policy.Value != 30 {
		t.Errorf("Expected value 30, got %d", policy.Value)
	}
	if policy.TagPattern != testTagPatternRelease {
		t.Errorf("Expected tag pattern '%s', got %s", testTagPatternRelease, policy.TagPattern)
	}
}

// --- Applications (Get/Update/Delete/ResetSecret single app) ---

func TestGetApplication(t *testing.T) {
	mockApp := Application{
		ClientID:    testClientID,
		Name:        testAppName,
		Description: "A test app",
	}
	mockResponseJSON, _ := json.Marshal(mockApp)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/applications/" + testClientID
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

	app, err := client.GetApplication(testOrgName, testClientID)
	if err != nil {
		t.Fatalf("GetApplication returned error: %v", err)
	}

	if app.ClientID != testClientID {
		t.Errorf("Expected client ID %s, got %s", testClientID, app.ClientID)
	}
}

func TestUpdateApplication(t *testing.T) {
	mockApp := Application{
		ClientID:       testClientID,
		Name:           "Updated App",
		Description:    updatedDescription,
		ApplicationURI: testAppURI,
		RedirectURI:    testRedirectURI,
	}
	mockResponseJSON, _ := json.Marshal(mockApp)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/applications/" + testClientID
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

	app, err := client.UpdateApplication(testOrgName, testClientID, "Updated App", updatedDescription, testAppURI, testRedirectURI)
	if err != nil {
		t.Fatalf("UpdateApplication returned error: %v", err)
	}

	if app.Name != "Updated App" {
		t.Errorf("Expected app name 'Updated App', got %s", app.Name)
	}
}

func TestDeleteApplication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/applications/" + testClientID
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

	err = client.DeleteApplication(testOrgName, testClientID)
	if err != nil {
		t.Fatalf("DeleteApplication returned error: %v", err)
	}
}

func TestResetApplicationClientSecret(t *testing.T) {
	mockApp := Application{
		ClientID:     testClientID,
		ClientSecret: "new-secret-456",
		Name:         testAppName,
	}
	mockResponseJSON, _ := json.Marshal(mockApp)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/applications/" + testClientID + "/resetclientsecret"
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

	app, err := client.ResetApplicationClientSecret(testOrgName, testClientID)
	if err != nil {
		t.Fatalf("ResetApplicationClientSecret returned error: %v", err)
	}

	if app.ClientSecret != "new-secret-456" {
		t.Errorf("Expected client secret 'new-secret-456', got %s", app.ClientSecret)
	}
}

// --- Marketplace (Create/BatchRemove/Delete subscription) ---

func TestCreateOrganizationMarketplaceSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/marketplace"
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

	subReq := &MarketplaceSubscriptionRequest{
		SKU:      "premium-plan",
		Quantity: 1,
	}

	err = client.CreateOrganizationMarketplaceSubscription(testOrgName, subReq)
	if err != nil {
		t.Fatalf("CreateOrganizationMarketplaceSubscription returned error: %v", err)
	}
}

func TestBatchRemoveOrganizationMarketplaceSubscriptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/marketplace/batchremove"
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

	err = client.BatchRemoveOrganizationMarketplaceSubscriptions(testOrgName, []string{testSubscriptionID, "sub-456"})
	if err != nil {
		t.Fatalf("BatchRemoveOrganizationMarketplaceSubscriptions returned error: %v", err)
	}
}

func TestDeleteOrganizationMarketplaceSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/marketplace/" + testSubscriptionID
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

	err = client.DeleteOrganizationMarketplaceSubscription(testOrgName, testSubscriptionID)
	if err != nil {
		t.Fatalf("DeleteOrganizationMarketplaceSubscription returned error: %v", err)
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

// --- Collaborators ---

func TestGetOrganizationCollaborators(t *testing.T) {
	mockCollabs := Collaborators{
		Collaborators: []Collaborator{
			{Name: testMemberName, Kind: testKindUser},
			{Name: "external-user", Kind: testKindUser},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockCollabs)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/collaborators"
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

	collabs, err := client.GetOrganizationCollaborators(testOrgName)
	if err != nil {
		t.Fatalf("GetOrganizationCollaborators returned error: %v", err)
	}

	if len(collabs.Collaborators) != 2 {
		t.Errorf("Expected 2 collaborators, got %d", len(collabs.Collaborators))
	}
	if collabs.Collaborators[0].Name != testMemberName {
		t.Errorf("Expected collaborator name %s, got %s", testMemberName, collabs.Collaborators[0].Name)
	}
}

// --- Marketplace ---

func TestGetOrganizationMarketplace(t *testing.T) {
	mockMarketplace := MarketplaceInfo{
		HasPayment: true,
		Subscriptions: []MarketplaceSubscription{
			{ID: testSubscriptionID, SKU: "premium-plan", Status: "active"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockMarketplace)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/marketplace"
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

	marketplace, err := client.GetOrganizationMarketplace(testOrgName)
	if err != nil {
		t.Fatalf("GetOrganizationMarketplace returned error: %v", err)
	}

	if !marketplace.HasPayment {
		t.Errorf("Expected HasPayment to be true")
	}
	if len(marketplace.Subscriptions) != 1 {
		t.Errorf("Expected 1 subscription, got %d", len(marketplace.Subscriptions))
	}
	if marketplace.Subscriptions[0].ID != testSubscriptionID {
		t.Errorf("Expected subscription ID %s, got %s", testSubscriptionID, marketplace.Subscriptions[0].ID)
	}
}

func newOrgErrorClient(t *testing.T) *Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func TestOrganizationCRUDHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.CreateOrganization("testorg", testEmailAddress)
	if err == nil {
		t.Error("Expected error from CreateOrganization, got nil")
	}

	_, err = client.UpdateOrganization(testOrgName, testEmailAddress)
	if err == nil {
		t.Error("Expected error from UpdateOrganization, got nil")
	}

	err = client.DeleteOrganization(testOrgName)
	if err == nil {
		t.Error("Expected error from DeleteOrganization, got nil")
	}

	err = client.AddOrganizationMember(testOrgName, testMemberName)
	if err == nil {
		t.Error("Expected error from AddOrganizationMember, got nil")
	}

	err = client.RemoveOrganizationMember(testOrgName, testMemberName)
	if err == nil {
		t.Error("Expected error from RemoveOrganizationMember, got nil")
	}

	_, err = client.GetOrganizationMember(testOrgName, testMemberName)
	if err == nil {
		t.Error("Expected error from GetOrganizationMember, got nil")
	}

	_, err = client.GetOrganizationRepositories(testOrgName)
	if err == nil {
		t.Error("Expected error from GetOrganizationRepositories, got nil")
	}

	_, err = client.GetOrganizationCollaborators(testOrgName)
	if err == nil {
		t.Error("Expected error from GetOrganizationCollaborators, got nil")
	}

	_, err = client.GetDefaultPermissions(testOrgName)
	if err == nil {
		t.Error("Expected error from GetDefaultPermissions, got nil")
	}

	_, err = client.CreateDefaultPermission(testOrgName, testRoleRead, "user", testMemberName)
	if err == nil {
		t.Error("Expected error from CreateDefaultPermission, got nil")
	}

	err = client.DeleteDefaultPermission(testOrgName, testPrototypeID)
	if err == nil {
		t.Error("Expected error from DeleteDefaultPermission, got nil")
	}

	_, err = client.GetProxyCacheConfig(testOrgName)
	if err == nil {
		t.Error("Expected error from GetProxyCacheConfig, got nil")
	}

	_, err = client.CreateProxyCacheConfig(testOrgName, testUpstreamReg, false, 86400)
	if err == nil {
		t.Error("Expected error from CreateProxyCacheConfig, got nil")
	}

	err = client.DeleteProxyCacheConfig(testOrgName)
	if err == nil {
		t.Error("Expected error from DeleteProxyCacheConfig, got nil")
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

func TestOrganizationRobotHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetRobotAccounts(testOrgName)
	if err == nil {
		t.Error("Expected error from GetRobotAccounts, got nil")
	}

	_, err = client.CreateRobotAccount(testOrgName, "testbot", testRobotDescValue, nil)
	if err == nil {
		t.Error("Expected error from CreateRobotAccount, got nil")
	}

	_, err = client.GetRobotAccount(testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from GetRobotAccount, got nil")
	}

	err = client.DeleteRobotAccount(testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from DeleteRobotAccount, got nil")
	}

	_, err = client.RegenerateRobotToken(testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from RegenerateRobotToken, got nil")
	}

	_, err = client.GetRobotPermissions(testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from GetRobotPermissions, got nil")
	}

	err = client.SetRobotRepositoryPermission(testOrgName, "testbot", testRepository, testRoleRead)
	if err == nil {
		t.Error("Expected error from SetRobotRepositoryPermission, got nil")
	}

	err = client.RemoveRobotRepositoryPermission(testOrgName, "testbot", testRepository)
	if err == nil {
		t.Error("Expected error from RemoveRobotRepositoryPermission, got nil")
	}

	_, err = client.GetRobotFederation(testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from GetRobotFederation, got nil")
	}

	err = client.CreateRobotFederation(testOrgName, "testbot", []RobotFederationConfig{{Issuer: "https://example.com", Subject: testPlaceholder}})
	if err == nil {
		t.Error("Expected error from CreateRobotFederation, got nil")
	}

	err = client.DeleteRobotFederation(testOrgName, "testbot")
	if err == nil {
		t.Error("Expected error from DeleteRobotFederation, got nil")
	}
}

func TestOrganizationQuotaPolicyHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetQuota(testOrgName)
	if err == nil {
		t.Error("Expected error from GetQuota, got nil")
	}

	_, err = client.CreateQuota(testOrgName, 1073741824)
	if err == nil {
		t.Error("Expected error from CreateQuota, got nil")
	}

	_, err = client.UpdateQuota(testOrgName, 2147483648)
	if err == nil {
		t.Error("Expected error from UpdateQuota, got nil")
	}

	err = client.DeleteQuota(testOrgName)
	if err == nil {
		t.Error("Expected error from DeleteQuota, got nil")
	}

	_, err = client.GetAutoPrunePolicies(testOrgName)
	if err == nil {
		t.Error("Expected error from GetAutoPrunePolicies, got nil")
	}

	_, err = client.CreateAutoPrunePolicy(testOrgName, testAutoPruneMethodNumberOfTags, 10, "")
	if err == nil {
		t.Error("Expected error from CreateAutoPrunePolicy, got nil")
	}

	_, err = client.GetAutoPrunePolicy(testOrgName, testPolicyUUID)
	if err == nil {
		t.Error("Expected error from GetAutoPrunePolicy, got nil")
	}

	_, err = client.UpdateAutoPrunePolicy(testOrgName, testPolicyUUID, testAutoPruneMethodNumberOfTags, 20, "")
	if err == nil {
		t.Error("Expected error from UpdateAutoPrunePolicy, got nil")
	}

	err = client.DeleteAutoPrunePolicy(testOrgName, testPolicyUUID)
	if err == nil {
		t.Error("Expected error from DeleteAutoPrunePolicy, got nil")
	}
}

func TestOrganizationAppMarketplaceHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetApplications(testOrgName)
	if err == nil {
		t.Error("Expected error from GetApplications, got nil")
	}

	_, err = client.CreateApplication(testOrgName, testAppName, "desc", testAppURI, testRedirectURI)
	if err == nil {
		t.Error("Expected error from CreateApplication, got nil")
	}

	_, err = client.GetApplication(testOrgName, testClientID)
	if err == nil {
		t.Error("Expected error from GetApplication, got nil")
	}

	_, err = client.UpdateApplication(testOrgName, testClientID, testAppName, "desc", testAppURI, testRedirectURI)
	if err == nil {
		t.Error("Expected error from UpdateApplication, got nil")
	}

	err = client.DeleteApplication(testOrgName, testClientID)
	if err == nil {
		t.Error("Expected error from DeleteApplication, got nil")
	}

	_, err = client.ResetApplicationClientSecret(testOrgName, testClientID)
	if err == nil {
		t.Error("Expected error from ResetApplicationClientSecret, got nil")
	}

	_, err = client.GetOrganizationMarketplace(testOrgName)
	if err == nil {
		t.Error("Expected error from GetOrganizationMarketplace, got nil")
	}

	err = client.CreateOrganizationMarketplaceSubscription(testOrgName, &MarketplaceSubscriptionRequest{SKU: testPlaceholder})
	if err == nil {
		t.Error("Expected error from CreateOrganizationMarketplaceSubscription, got nil")
	}

	err = client.BatchRemoveOrganizationMarketplaceSubscriptions(testOrgName, []string{testSubscriptionID})
	if err == nil {
		t.Error("Expected error from BatchRemoveOrganizationMarketplaceSubscriptions, got nil")
	}

	err = client.DeleteOrganizationMarketplaceSubscription(testOrgName, testSubscriptionID)
	if err == nil {
		t.Error("Expected error from DeleteOrganizationMarketplaceSubscription, got nil")
	}
}
