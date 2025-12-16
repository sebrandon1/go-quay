package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetError  = "GET"
	testErrorType = "invalid_token"
)

func TestGetErrorType(t *testing.T) {
	mockResponse := ErrorType{
		Type:        testErrorType,
		Title:       "Invalid Token",
		Description: "The provided authentication token is invalid or expired",
		Status:      401,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetError {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/error/" + testErrorType
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

	errType, err := client.GetErrorType(testErrorType)
	if err != nil {
		t.Fatalf("GetErrorType returned error: %v", err)
	}

	if errType.Type != testErrorType {
		t.Errorf("Expected type '%s', got %s", testErrorType, errType.Type)
	}
	if errType.Status != 401 {
		t.Errorf("Expected status 401, got %d", errType.Status)
	}
}

func TestGetErrorTypeError(t *testing.T) {
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

	_, err = client.GetErrorType("nonexistent_error")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
