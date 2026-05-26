package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testUpdatedItem = "updated-item"
	testFieldName   = "name"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(testTokenValue)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if client.BearerToken != testTokenValue {
		t.Errorf("Expected bearer token %s, got %s", testTokenValue, client.BearerToken)
	}

	if client.BaseURL != DefaultQuayURL {
		t.Errorf("Expected base URL %s, got %s", DefaultQuayURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTP client to be set, got nil")
	}
}

func TestNewClientEmptyToken(t *testing.T) {
	client, err := NewClient("")
	if err == nil {
		t.Error("Expected error for empty token, got nil")
	}

	if client != nil {
		t.Error("Expected nil client for empty token")
	}
}

func TestNewClientWithURL(t *testing.T) {
	customURL := "https://custom-quay.example.com/api/v1"

	client, err := NewClientWithURL(testTokenValue, customURL)
	if err != nil {
		t.Fatalf("NewClientWithURL returned error: %v", err)
	}

	if client.BearerToken != testTokenValue {
		t.Errorf("Expected bearer token %s, got %s", testTokenValue, client.BearerToken)
	}

	if client.BaseURL != customURL {
		t.Errorf("Expected base URL %s, got %s", customURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTP client to be set, got nil")
	}
}

func TestNewClientWithURLEmptyToken(t *testing.T) {
	client, err := NewClientWithURL("", "https://quay.example.com/api/v1")
	if err == nil {
		t.Error("Expected error for empty token, got nil")
	}

	if client != nil {
		t.Error("Expected nil client for empty token")
	}
}

func TestGetRequest(t *testing.T) {
	type testResponse struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	mockResponse := testResponse{Name: "test-item", Value: 42}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+testTokenValue {
			t.Errorf("Expected Authorization header 'Bearer %s', got '%s'", testTokenValue, authHeader)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequest(httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result testResponse
	err = client.get(req, &result)
	if err != nil {
		t.Fatalf("get returned error: %v", err)
	}

	if result.Name != "test-item" {
		t.Errorf("Expected name 'test-item', got '%s'", result.Name)
	}

	if result.Value != 42 {
		t.Errorf("Expected value 42, got %d", result.Value)
	}
}

func TestGetRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequest(httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if err == nil {
		t.Error("Expected error for 500 response, got nil")
	}
}

func TestPostRequest(t *testing.T) {
	type postBody struct {
		Name string `json:"name"`
	}

	type postResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	mockResponse := postResponse{ID: "new-id-123", Name: "created-item"}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+testTokenValue {
			t.Errorf("Expected Authorization header 'Bearer %s', got '%s'", testTokenValue, authHeader)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		var reqBody postBody
		if err := json.Unmarshal(body, &reqBody); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}

		if reqBody.Name != "test-create" {
			t.Errorf("Expected request body name 'test-create', got '%s'", reqBody.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	reqBody := postBody{Name: "test-create"}
	req, err := newRequestWithBody(httpMethodPost, server.URL+"/api/v1/test", reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result postResponse
	err = client.post(req, &result)
	if err != nil {
		t.Fatalf("post returned error: %v", err)
	}

	if result.ID != "new-id-123" {
		t.Errorf("Expected ID 'new-id-123', got '%s'", result.ID)
	}

	if result.Name != "created-item" {
		t.Errorf("Expected name 'created-item', got '%s'", result.Name)
	}
}

func TestPostRequestNoResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequestWithBody(httpMethodPost, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	err = client.post(req, nil)
	if err != nil {
		t.Fatalf("post with nil response target returned error: %v", err)
	}
}

func TestPutRequest(t *testing.T) {
	type putResponse struct {
		Updated bool   `json:"updated"`
		Name    string `json:"name"`
	}

	mockResponse := putResponse{Updated: true, Name: testUpdatedItem}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+testTokenValue {
			t.Errorf("Expected Authorization header 'Bearer %s', got '%s'", testTokenValue, authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequestWithBody(httpMethodPut, server.URL+"/api/v1/test", map[string]string{testFieldName: testUpdatedItem})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result putResponse
	err = client.put(req, &result)
	if err != nil {
		t.Fatalf("put returned error: %v", err)
	}

	if !result.Updated {
		t.Error("Expected updated to be true")
	}

	if result.Name != testUpdatedItem {
		t.Errorf("Expected name '%s', got '%s'", testUpdatedItem, result.Name)
	}
}

func TestPutRequestNoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequestWithBody(httpMethodPut, server.URL+"/api/v1/test", map[string]string{testFieldName: "test"})
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	err = client.put(req, nil)
	if err != nil {
		t.Fatalf("put with 204 No Content returned error: %v", err)
	}
}

func TestDeleteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+testTokenValue {
			t.Errorf("Expected Authorization header 'Bearer %s', got '%s'", testTokenValue, authHeader)
		}

		expectedPath := "/api/v1/test/item-123"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequest(httpMethodDelete, server.URL+"/api/v1/test/item-123", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	err = client.delete(req)
	if err != nil {
		t.Fatalf("delete returned error: %v", err)
	}
}

func TestDeleteRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequest(httpMethodDelete, server.URL+"/api/v1/test/item-123", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	err = client.delete(req)
	if err == nil {
		t.Error("Expected error for 500 response, got nil")
	}
}
