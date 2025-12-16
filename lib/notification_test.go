package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetNotification    = "GET"
	httpPostNotification   = "POST"
	httpDeleteNotification = "DELETE"

	testNotificationNamespace  = "testorg"
	testNotificationRepository = "testrepo"
	testNotificationUUID       = "notification-uuid-123"
	testNotificationEvent      = "repo_push"
	testNotificationMethod     = "webhook"
)

func TestGetNotifications(t *testing.T) {
	mockResponse := RepositoryNotifications{
		Notifications: []RepositoryNotification{
			{UUID: testNotificationUUID, Event: testNotificationEvent, Method: testNotificationMethod, Title: "Push Webhook"},
			{UUID: "notification-uuid-456", Event: "build_success", Method: "slack", Title: "Build Slack"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetNotification {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNotificationNamespace + "/" + testNotificationRepository + "/notification/"
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

	notifications, err := client.GetNotifications(testNotificationNamespace, testNotificationRepository)
	if err != nil {
		t.Fatalf("GetNotifications returned error: %v", err)
	}

	if len(notifications.Notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(notifications.Notifications))
	}
	if notifications.Notifications[0].UUID != testNotificationUUID {
		t.Errorf("Expected first notification UUID %s, got %s", testNotificationUUID, notifications.Notifications[0].UUID)
	}
}

func TestGetNotification(t *testing.T) {
	mockResponse := RepositoryNotification{
		UUID:   testNotificationUUID,
		Event:  testNotificationEvent,
		Method: testNotificationMethod,
		Title:  "Push Webhook",
		Config: map[string]interface{}{"url": "https://example.com/webhook"},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetNotification {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNotificationNamespace + "/" + testNotificationRepository + "/notification/" + testNotificationUUID
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

	notification, err := client.GetNotification(testNotificationNamespace, testNotificationRepository, testNotificationUUID)
	if err != nil {
		t.Fatalf("GetNotification returned error: %v", err)
	}

	if notification.UUID != testNotificationUUID {
		t.Errorf("Expected notification UUID %s, got %s", testNotificationUUID, notification.UUID)
	}
	if notification.Event != testNotificationEvent {
		t.Errorf("Expected event %s, got %s", testNotificationEvent, notification.Event)
	}
}

func TestCreateNotification(t *testing.T) {
	mockResponse := RepositoryNotification{
		UUID:   testNotificationUUID,
		Event:  testNotificationEvent,
		Method: testNotificationMethod,
		Title:  "New Webhook",
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPostNotification {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNotificationNamespace + "/" + testNotificationRepository + "/notification/"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
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

	notificationReq := &CreateNotificationRequest{
		Event:  testNotificationEvent,
		Method: testNotificationMethod,
		Config: map[string]interface{}{"url": "https://example.com/webhook"},
		Title:  "New Webhook",
	}

	notification, err := client.CreateNotification(testNotificationNamespace, testNotificationRepository, notificationReq)
	if err != nil {
		t.Fatalf("CreateNotification returned error: %v", err)
	}

	if notification.UUID != testNotificationUUID {
		t.Errorf("Expected notification UUID %s, got %s", testNotificationUUID, notification.UUID)
	}
}

func TestDeleteNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeleteNotification {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNotificationNamespace + "/" + testNotificationRepository + "/notification/" + testNotificationUUID
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

	err = client.DeleteNotification(testNotificationNamespace, testNotificationRepository, testNotificationUUID)
	if err != nil {
		t.Fatalf("DeleteNotification returned error: %v", err)
	}
}

func TestTestNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPostNotification {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNotificationNamespace + "/" + testNotificationRepository + "/notification/" + testNotificationUUID + "/test"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.TestNotification(testNotificationNamespace, testNotificationRepository, testNotificationUUID)
	if err != nil {
		t.Fatalf("TestNotification returned error: %v", err)
	}
}

func TestResetNotification(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPostNotification {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/" + testNotificationNamespace + "/" + testNotificationRepository + "/notification/" + testNotificationUUID + "/reset"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.ResetNotification(testNotificationNamespace, testNotificationRepository, testNotificationUUID)
	if err != nil {
		t.Fatalf("ResetNotification returned error: %v", err)
	}
}

func TestGetNotificationsError(t *testing.T) {
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

	_, err = client.GetNotifications(testNotificationNamespace, testNotificationRepository)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetNotificationError(t *testing.T) {
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

	_, err = client.GetNotification(testNotificationNamespace, testNotificationRepository, "nonexistent-uuid")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
