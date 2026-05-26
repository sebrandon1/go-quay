package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMessages(t *testing.T) {
	mockResponse := Messages{
		Messages: []Message{
			{
				UUID:      "msg-uuid-123",
				Content:   "Scheduled maintenance on Sunday",
				Severity:  "info",
				MediaType: testMediaTypePlain,
			},
			{
				UUID:      "msg-uuid-456",
				Content:   "New feature available",
				Severity:  "info",
				MediaType: testMediaTypePlain,
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
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

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
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

func TestCreateMessage(t *testing.T) {
	mockMessage := Message{
		UUID:      "msg-uuid-new",
		Content:   "System maintenance scheduled",
		Severity:  "warning",
		MediaType: testMediaTypePlain,
	}
	mockResponseJSON, _ := json.Marshal(mockMessage)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/messages"
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

	message, err := client.CreateMessage("System maintenance scheduled", "warning", testMediaTypePlain)
	if err != nil {
		t.Fatalf("CreateMessage returned error: %v", err)
	}

	if message.UUID != "msg-uuid-new" {
		t.Errorf("Expected message UUID 'msg-uuid-new', got %s", message.UUID)
	}
	if message.Severity != "warning" {
		t.Errorf("Expected severity 'warning', got %s", message.Severity)
	}
}

func TestGetMessagesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetMessages()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
