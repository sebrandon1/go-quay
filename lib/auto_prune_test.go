package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Auto-Prune ---

func TestGetAutoPrunePolicies(t *testing.T) {
	mockPolicies := AutoPrunePolicies{
		Policies: []AutoPrunePolicy{
			{UUID: testPolicyUUID, Method: testAutoPruneMethodNumberOfTags, Value: 10, TagPattern: "v*"},
		},
	}
	mockResponseJSON := mustMarshal(t, mockPolicies)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy"
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

	policies, err := client.GetAutoPrunePolicies(context.Background(), testOrgName)
	if err != nil {
		t.Fatalf("GetAutoPrunePolicies returned error: %v", err)
	}

	if len(policies.Policies) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(policies.Policies))
	}
	if policies.Policies[0].UUID != testPolicyUUID {
		t.Errorf("Expected policy UUID %s, got %s", testPolicyUUID, policies.Policies[0].UUID)
	}
}

func TestCreateAutoPrunePolicy(t *testing.T) {
	mockPolicy := AutoPrunePolicy{
		UUID:       testPolicyUUID,
		Method:     testAutoPruneMethodNumberOfTags,
		Value:      20,
		TagPattern: testTagPatternRelease,
	}
	mockResponseJSON := mustMarshal(t, mockPolicy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var req CreateAutoPruneRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if req.Method != testAutoPruneMethodNumberOfTags {
			t.Errorf("Expected method 'number_of_tags', got %s", req.Method)
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

	policy, err := client.CreateAutoPrunePolicy(context.Background(), testOrgName, testAutoPruneMethodNumberOfTags, 20, testTagPatternRelease)
	if err != nil {
		t.Fatalf("CreateAutoPrunePolicy returned error: %v", err)
	}

	if policy.Value != 20 {
		t.Errorf("Expected value 20, got %d", policy.Value)
	}
}

func TestGetAutoPrunePolicy(t *testing.T) {
	mockPolicy := AutoPrunePolicy{
		UUID:       testPolicyUUID,
		Method:     testAutoPruneMethodNumberOfTags,
		Value:      10,
		TagPattern: "v*",
	}
	mockResponseJSON := mustMarshal(t, mockPolicy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy/" + testPolicyUUID
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

	policy, err := client.GetAutoPrunePolicy(context.Background(), testOrgName, testPolicyUUID)
	if err != nil {
		t.Fatalf("GetAutoPrunePolicy returned error: %v", err)
	}

	if policy.UUID != testPolicyUUID {
		t.Errorf("Expected policy UUID %s, got %s", testPolicyUUID, policy.UUID)
	}
	if policy.Value != 10 {
		t.Errorf("Expected value 10, got %d", policy.Value)
	}
}

func TestUpdateAutoPrunePolicy(t *testing.T) {
	mockPolicy := AutoPrunePolicy{
		UUID:       testPolicyUUID,
		Method:     testAutoPruneMethodNumberOfTags,
		Value:      30,
		TagPattern: testTagPatternRelease,
	}
	mockResponseJSON := mustMarshal(t, mockPolicy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy/" + testPolicyUUID
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

	policy, err := client.UpdateAutoPrunePolicy(context.Background(), testOrgName, testPolicyUUID, testAutoPruneMethodNumberOfTags, 30, testTagPatternRelease)
	if err != nil {
		t.Fatalf("UpdateAutoPrunePolicy returned error: %v", err)
	}

	if policy.Value != 30 {
		t.Errorf("Expected value 30, got %d", policy.Value)
	}
	if policy.TagPattern != testTagPatternRelease {
		t.Errorf("Expected tag pattern '%s', got %s", testTagPatternRelease, policy.TagPattern)
	}
}

func TestDeleteAutoPrunePolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/autoprunepolicy/" + testPolicyUUID
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

	err = client.DeleteAutoPrunePolicy(context.Background(), testOrgName, testPolicyUUID)
	if err != nil {
		t.Fatalf("DeleteAutoPrunePolicy returned error: %v", err)
	}
}

func TestAutoPruneHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetAutoPrunePolicies(context.Background(), testOrgName)
	if err == nil {
		t.Error("Expected error from GetAutoPrunePolicies, got nil")
	}

	_, err = client.CreateAutoPrunePolicy(context.Background(), testOrgName, testAutoPruneMethodNumberOfTags, 10, "")
	if err == nil {
		t.Error("Expected error from CreateAutoPrunePolicy, got nil")
	}

	_, err = client.GetAutoPrunePolicy(context.Background(), testOrgName, testPolicyUUID)
	if err == nil {
		t.Error("Expected error from GetAutoPrunePolicy, got nil")
	}

	_, err = client.UpdateAutoPrunePolicy(context.Background(), testOrgName, testPolicyUUID, testAutoPruneMethodNumberOfTags, 20, "")
	if err == nil {
		t.Error("Expected error from UpdateAutoPrunePolicy, got nil")
	}

	err = client.DeleteAutoPrunePolicy(context.Background(), testOrgName, testPolicyUUID)
	if err == nil {
		t.Error("Expected error from DeleteAutoPrunePolicy, got nil")
	}
}
