package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetMessages = "GET"
)

func TestGetMessages(t *testing.T) {
	mockResponse := Messages{
		Messages: []Message{
			{
				UUID:      "msg-uuid-123",
				Content:   "Scheduled maintenance on Sunday",
				Severity:  "info",
				MediaType: "text/plain",
			},
			{
				UUID:      "msg-uuid-456",
				Content:   "New feature available",
				Severity:  "info",
				MediaType: "text/plain",
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetMessages {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/messages"
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

	messages, err := client.GetMessages()
	if err != nil {
		t.Fatalf("GetMessages returned error: %v", err)
	}

	if len(messages.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages.Messages))
	}
	if messages.Messages[0].UUID != "msg-uuid-123" {
		t.Errorf("Expected first message UUID 'msg-uuid-123', got %s", messages.Messages[0].UUID)
	}
}

func TestGetMessagesEmpty(t *testing.T) {
	mockResponse := Messages{
		Messages: []Message{},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	messages, err := client.GetMessages()
	if err != nil {
		t.Fatalf("GetMessages returned error: %v", err)
	}

	if len(messages.Messages) != 0 {
		t.Errorf("Expected 0 messages, got %d", len(messages.Messages))
	}
}

func TestGetMessagesError(t *testing.T) {
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

	_, err = client.GetMessages()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
