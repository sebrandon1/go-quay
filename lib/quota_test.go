package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Quota ---

func TestGetQuota(t *testing.T) {
	mockQuota := Quota{
		ID:         "quota-1",
		LimitBytes: 1073741824,
	}
	mockResponseJSON := mustMarshal(t, mockQuota)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/quota"
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

	quota, err := client.GetQuota(context.Background(), testOrgName)
	if err != nil {
		t.Fatalf("GetQuota returned error: %v", err)
	}

	if quota.LimitBytes != 1073741824 {
		t.Errorf("Expected limit bytes 1073741824, got %d", quota.LimitBytes)
	}
}

func TestCreateQuota(t *testing.T) {
	var limitBytes int64 = 2147483648
	mockQuota := Quota{
		ID:         "quota-2",
		LimitBytes: limitBytes,
	}
	mockResponseJSON := mustMarshal(t, mockQuota)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/quota"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req CreateQuotaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.LimitBytes != limitBytes {
			t.Errorf("Expected limit bytes %d, got %d", limitBytes, req.LimitBytes)
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

	quota, err := client.CreateQuota(context.Background(), testOrgName, limitBytes)
	if err != nil {
		t.Fatalf("CreateQuota returned error: %v", err)
	}

	if quota.LimitBytes != limitBytes {
		t.Errorf("Expected limit bytes %d, got %d", limitBytes, quota.LimitBytes)
	}
}

func TestUpdateQuota(t *testing.T) {
	var limitBytes int64 = 4294967296
	mockQuota := Quota{
		ID:         "quota-1",
		LimitBytes: limitBytes,
	}
	mockResponseJSON := mustMarshal(t, mockQuota)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/quota"
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

	quota, err := client.UpdateQuota(context.Background(), testOrgName, limitBytes)
	if err != nil {
		t.Fatalf("UpdateQuota returned error: %v", err)
	}

	if quota.LimitBytes != limitBytes {
		t.Errorf("Expected limit bytes %d, got %d", limitBytes, quota.LimitBytes)
	}
}

func TestDeleteQuota(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/quota"
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

	err = client.DeleteQuota(context.Background(), testOrgName)
	if err != nil {
		t.Fatalf("DeleteQuota returned error: %v", err)
	}
}

func TestQuotaHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetQuota(context.Background(), testOrgName)
	if err == nil {
		t.Error("Expected error from GetQuota, got nil")
	}

	_, err = client.CreateQuota(context.Background(), testOrgName, 1073741824)
	if err == nil {
		t.Error("Expected error from CreateQuota, got nil")
	}

	_, err = client.UpdateQuota(context.Background(), testOrgName, 2147483648)
	if err == nil {
		t.Error("Expected error from UpdateQuota, got nil")
	}

	err = client.DeleteQuota(context.Background(), testOrgName)
	if err == nil {
		t.Error("Expected error from DeleteQuota, got nil")
	}
}
