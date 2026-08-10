package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestApplicationHTTPErrors(t *testing.T) {
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
}
