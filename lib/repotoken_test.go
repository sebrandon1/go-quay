package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetRepoToken    = "GET"
	httpPostRepoToken   = "POST"
	httpPutRepoToken    = "PUT"
	httpDeleteRepoToken = "DELETE"

	testTokenNamespace  = "testorg"
	testTokenRepository = "testrepo"
	testTokenCode       = "ABCD1234"
	testTokenRole       = "read"
)

func TestGetRepoTokens(t *testing.T) {
	mockResponse := RepoTokens{
		Tokens: []RepoToken{
			{Code: testTokenCode, FriendlyName: "CI Token", Role: testTokenRole},
			{Code: "EFGH5678", FriendlyName: "Deploy Token", Role: "write"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetRepoToken {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testTokenNamespace + "/" + testTokenRepository + "/tokens"
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

	tokens, err := client.GetRepoTokens(testTokenNamespace, testTokenRepository)
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
		Role:         testTokenRole,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPostRepoToken {
			t.Errorf("Expected POST request, got %s", r.Method)
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

	createReq := &CreateRepoTokenRequest{
		FriendlyName: "New CI Token",
	}

	token, err := client.CreateRepoToken(testTokenNamespace, testTokenRepository, createReq)
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
		FriendlyName: "CI Token",
		Role:         testTokenRole,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetRepoToken {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testTokenNamespace + "/" + testTokenRepository + "/tokens/" + testTokenCode
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

	token, err := client.GetRepoToken(testTokenNamespace, testTokenRepository, testTokenCode)
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
		FriendlyName: "CI Token",
		Role:         "write",
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPutRepoToken {
			t.Errorf("Expected PUT request, got %s", r.Method)
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

	updateReq := &UpdateRepoTokenRequest{
		Role: "write",
	}

	token, err := client.UpdateRepoToken(testTokenNamespace, testTokenRepository, testTokenCode, updateReq)
	if err != nil {
		t.Fatalf("UpdateRepoToken returned error: %v", err)
	}

	if token.Role != "write" {
		t.Errorf("Expected role 'write', got %s", token.Role)
	}
}

func TestDeleteRepoToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteRepoToken {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testTokenNamespace + "/" + testTokenRepository + "/tokens/" + testTokenCode
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

	err = client.DeleteRepoToken(testTokenNamespace, testTokenRepository, testTokenCode)
	if err != nil {
		t.Fatalf("DeleteRepoToken returned error: %v", err)
	}
}

func TestGetRepoTokensError(t *testing.T) {
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

	_, err = client.GetRepoTokens(testTokenNamespace, testTokenRepository)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
