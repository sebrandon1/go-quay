package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

// --- HTTP Error Tests ---

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
}
