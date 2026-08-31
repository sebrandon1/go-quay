package lib

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
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

	if client.Version != "dev" {
		t.Errorf("Expected default version 'dev', got %s", client.Version)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTP client to be set, got nil")
	}
}

func TestNewClientEmptyToken(t *testing.T) {
	client, err := NewClient("")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client for empty token")
	}

	if client.BearerToken != "" {
		t.Errorf("Expected empty bearer token, got %s", client.BearerToken)
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
	if err != nil {
		t.Fatalf("NewClientWithURL returned unexpected error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client for empty token")
	}

	if client.BearerToken != "" {
		t.Errorf("Expected empty bearer token, got %s", client.BearerToken)
	}

	if client.BaseURL != "https://quay.example.com/api/v1" {
		t.Errorf("Expected custom base URL, got %s", client.BaseURL)
	}
}

func TestUserAgentHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if ua != "go-quay/1.2.3" {
			t.Errorf("Expected User-Agent 'go-quay/1.2.3', got '%s'", ua)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Version = "1.2.3"

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]any
	err = client.get(req, &result)
	if err != nil {
		t.Fatalf("get returned error: %v", err)
	}
}

func TestGetRequest(t *testing.T) {
	type testResponse struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	mockResponse := testResponse{Name: "test-item", Value: 42}
	mockResponseJSON := mustMarshal(t, mockResponse)

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

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
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

func TestQuayErrorParsing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":404,"error":"not_found","detail":"Repository not found","error_type":"not_found"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var quayErr *QuayError
	if !errors.As(err, &quayErr) {
		t.Fatalf("Expected QuayError, got %T: %v", err, err)
	}

	if quayErr.StatusCode() != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", quayErr.StatusCode())
	}

	if quayErr.Message != "not_found" {
		t.Errorf("Expected message 'not_found', got '%s'", quayErr.Message)
	}

	if quayErr.Detail != "Repository not found" {
		t.Errorf("Expected detail 'Repository not found', got '%s'", quayErr.Detail)
	}

	expectedErrStr := "quay API error (status 404): not_found — Repository not found"
	if quayErr.Error() != expectedErrStr {
		t.Errorf("Expected error string '%s', got '%s'", expectedErrStr, quayErr.Error())
	}
}

func TestQuayErrorFallbackToRawError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`not json at all`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	var quayErr *QuayError
	if errors.As(err, &quayErr) {
		t.Error("Expected raw error for non-JSON response, got QuayError")
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

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
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
	mockResponseJSON := mustMarshal(t, mockResponse)

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
	req, err := newRequestWithBody(context.Background(), httpMethodPost, server.URL+"/api/v1/test", reqBody)
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

	req, err := newRequestWithBody(context.Background(), httpMethodPost, server.URL+"/api/v1/test", nil)
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
	mockResponseJSON := mustMarshal(t, mockResponse)

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

	req, err := newRequestWithBody(context.Background(), httpMethodPut, server.URL+"/api/v1/test", map[string]string{testFieldName: testUpdatedItem})
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

	req, err := newRequestWithBody(context.Background(), httpMethodPut, server.URL+"/api/v1/test", map[string]string{testFieldName: testPlaceholder})
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

	req, err := newRequest(context.Background(), httpMethodDelete, server.URL+"/api/v1/test/item-123", nil)
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

	req, err := newRequest(context.Background(), httpMethodDelete, server.URL+"/api/v1/test/item-123", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	err = client.delete(req)
	if err == nil {
		t.Error("Expected error for 500 response, got nil")
	}
}

func TestRetryOn429(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := attempts.Add(1)
		if n <= 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate_limited"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"success"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Retry = &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
	}

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if err != nil {
		t.Fatalf("Expected success after retries, got error: %v", err)
	}

	if result["name"] != "success" {
		t.Errorf("Expected name 'success', got '%s'", result["name"])
	}

	if attempts.Load() != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts.Load())
	}
}

func TestRetryOn500(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := attempts.Add(1)
		if n == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"internal"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Retry = &RetryConfig{
		MaxRetries:     2,
		InitialBackoff: 1 * time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
	}

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]any
	err = client.get(req, &result)
	if err != nil {
		t.Fatalf("Expected success after retry, got error: %v", err)
	}

	if attempts.Load() != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts.Load())
	}
}

func TestRetryExhausted(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":"rate_limited"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Retry = &RetryConfig{
		MaxRetries:     2,
		InitialBackoff: 1 * time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
	}

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if err == nil {
		t.Fatal("Expected error after exhausted retries, got nil")
	}

	if attempts.Load() != 3 {
		t.Errorf("Expected 3 attempts (1 + 2 retries), got %d", attempts.Load())
	}
}

func TestNoRetryOn4xx(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not_found"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Retry = &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Millisecond,
	}

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if err == nil {
		t.Fatal("Expected error for 404, got nil")
	}

	if attempts.Load() != 1 {
		t.Errorf("Expected 1 attempt (no retry on 404), got %d", attempts.Load())
	}
}

func TestNoRetryWithoutConfig(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req, err := newRequest(context.Background(), httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if attempts.Load() != 1 {
		t.Errorf("Expected 1 attempt (no retry config), got %d", attempts.Load())
	}
}

func TestCanceledContextBeforeDo(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, err := newRequest(ctx, httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled, got %v", err)
	}
	if attempts.Load() != 0 {
		t.Errorf("Expected 0 HTTP attempts for canceled context, got %d", attempts.Load())
	}
}

func TestCanceledContextAbortsInFlightRequest(t *testing.T) {
	started := make(chan struct{})
	block := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-block
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()
	defer close(block)

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := newRequest(ctx, httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		var result map[string]string
		done <- client.get(req, &result)
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("server did not receive request")
	}
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("get did not return after cancel")
	}
}

func TestCanceledContextDuringBackoff(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		select {
		case <-started:
		default:
			close(started)
		}
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":"rate_limited"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Retry = &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 30 * time.Second,
		MaxBackoff:     30 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := newRequest(ctx, httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		var result map[string]string
		done <- client.get(req, &result)
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("server did not receive request")
	}
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Expected context.Canceled during backoff, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("get did not return after cancel during backoff")
	}
}

func TestDeadlineExceededDuringBackoff(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":"rate_limited"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Retry = &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 30 * time.Second,
		MaxBackoff:     30 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	req, err := newRequest(ctx, httpMethodGet, server.URL+"/api/v1/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.get(req, &result)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected context.DeadlineExceeded during backoff, got %v", err)
	}
}
