package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testTokenCode = "ABCD1234"
)

func TestGetRepoTokens(t *testing.T) {
	mockResponse := RepoTokens{
		Tokens: []RepoToken{
			{Code: testTokenCode, FriendlyName: testRepoTokenName, Role: testRoleRead},
			{Code: "EFGH5678", FriendlyName: "Deploy Token", Role: testRoleWrite},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/tokens"
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

	tokens, err := client.GetRepoTokens(testNamespace, testRepository)
	if err != nil {
		t.Fatalf("GetRepoTokens returned error: %v", err)
	}

	if len(tokens.Tokens) != 2 {
		t.Errorf("Expected 2 tokens, got %d", len(tokens.Tokens))
	}
	if tokens.Tokens[0].Code != testTokenCode {
		t.Errorf("Expected first token code %s, got %s", testTokenCode, tokens.Tokens[0].Code)
	}
}

func TestCreateRepoToken(t *testing.T) {
	mockResponse := RepoToken{
		Code:         testTokenCode,
		FriendlyName: "New CI Token",
		Role:         testRoleRead,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
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

	createReq := &CreateRepoTokenRequest{
		FriendlyName: "New CI Token",
	}

	token, err := client.CreateRepoToken(testNamespace, testRepository, createReq)
	if err != nil {
		t.Fatalf("CreateRepoToken returned error: %v", err)
	}

	if token.Code != testTokenCode {
		t.Errorf("Expected token code %s, got %s", testTokenCode, token.Code)
	}
}

func TestGetRepoToken(t *testing.T) {
	mockResponse := RepoToken{
		Code:         testTokenCode,
		FriendlyName: testRepoTokenName,
		Role:         testRoleRead,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/tokens/" + testTokenCode
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

	token, err := client.GetRepoToken(testNamespace, testRepository, testTokenCode)
	if err != nil {
		t.Fatalf("GetRepoToken returned error: %v", err)
	}

	if token.Code != testTokenCode {
		t.Errorf("Expected token code %s, got %s", testTokenCode, token.Code)
	}
}

func TestUpdateRepoToken(t *testing.T) {
	mockResponse := RepoToken{
		Code:         testTokenCode,
		FriendlyName: testRepoTokenName,
		Role:         "write",
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	updateReq := &UpdateRepoTokenRequest{
		Role: testRoleWrite,
	}

	token, err := client.UpdateRepoToken(testNamespace, testRepository, testTokenCode, updateReq)
	if err != nil {
		t.Fatalf("UpdateRepoToken returned error: %v", err)
	}

	if token.Role != testRoleWrite {
		t.Errorf("Expected role 'write', got %s", token.Role)
	}
}

func TestDeleteRepoToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNamespace + "/" + testRepository + "/tokens/" + testTokenCode
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

	err = client.DeleteRepoToken(testNamespace, testRepository, testTokenCode)
	if err != nil {
		t.Fatalf("DeleteRepoToken returned error: %v", err)
	}
}

func TestGetRepoTokensError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetRepoTokens(testNamespace, testRepository)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRepoTokenHTTPErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.CreateRepoToken(testNamespace, testRepository, &CreateRepoTokenRequest{FriendlyName: testPlaceholder})
	if err == nil {
		t.Error("Expected error from CreateRepoToken, got nil")
	}

	_, err = client.GetRepoToken(testNamespace, testRepository, testTokenCode)
	if err == nil {
		t.Error("Expected error from GetRepoToken, got nil")
	}

	_, err = client.UpdateRepoToken(testNamespace, testRepository, testTokenCode, &UpdateRepoTokenRequest{Role: testRoleWrite})
	if err == nil {
		t.Error("Expected error from UpdateRepoToken, got nil")
	}

	err = client.DeleteRepoToken(testNamespace, testRepository, testTokenCode)
	if err == nil {
		t.Error("Expected error from DeleteRepoToken, got nil")
	}
}
