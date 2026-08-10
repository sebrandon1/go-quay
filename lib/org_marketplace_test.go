package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Marketplace ---

func TestGetOrganizationMarketplace(t *testing.T) {
	mockMarketplace := MarketplaceInfo{
		HasPayment: true,
		Subscriptions: []MarketplaceSubscription{
			{ID: testSubscriptionID, SKU: "premium-plan", Status: "active"},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockMarketplace)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/marketplace"
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

	marketplace, err := client.GetOrganizationMarketplace(testOrgName)
	if err != nil {
		t.Fatalf("GetOrganizationMarketplace returned error: %v", err)
	}

	if !marketplace.HasPayment {
		t.Errorf("Expected HasPayment to be true")
	}
	if len(marketplace.Subscriptions) != 1 {
		t.Errorf("Expected 1 subscription, got %d", len(marketplace.Subscriptions))
	}
	if marketplace.Subscriptions[0].ID != testSubscriptionID {
		t.Errorf("Expected subscription ID %s, got %s", testSubscriptionID, marketplace.Subscriptions[0].ID)
	}
}

func TestCreateOrganizationMarketplaceSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/marketplace"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	subReq := &MarketplaceSubscriptionRequest{
		SKU:      "premium-plan",
		Quantity: 1,
	}

	err = client.CreateOrganizationMarketplaceSubscription(testOrgName, subReq)
	if err != nil {
		t.Fatalf("CreateOrganizationMarketplaceSubscription returned error: %v", err)
	}
}

func TestBatchRemoveOrganizationMarketplaceSubscriptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/marketplace/batchremove"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.BatchRemoveOrganizationMarketplaceSubscriptions(testOrgName, []string{testSubscriptionID, "sub-456"})
	if err != nil {
		t.Fatalf("BatchRemoveOrganizationMarketplaceSubscriptions returned error: %v", err)
	}
}

func TestDeleteOrganizationMarketplaceSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testOrgName + "/marketplace/" + testSubscriptionID
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

	err = client.DeleteOrganizationMarketplaceSubscription(testOrgName, testSubscriptionID)
	if err != nil {
		t.Fatalf("DeleteOrganizationMarketplaceSubscription returned error: %v", err)
	}
}

func TestMarketplaceHTTPErrors(t *testing.T) {
	client := newOrgErrorClient(t)

	_, err := client.GetOrganizationMarketplace(testOrgName)
	if err == nil {
		t.Error("Expected error from GetOrganizationMarketplace, got nil")
	}

	err = client.CreateOrganizationMarketplaceSubscription(testOrgName, &MarketplaceSubscriptionRequest{SKU: testPlaceholder})
	if err == nil {
		t.Error("Expected error from CreateOrganizationMarketplaceSubscription, got nil")
	}

	err = client.BatchRemoveOrganizationMarketplaceSubscriptions(testOrgName, []string{testSubscriptionID})
	if err == nil {
		t.Error("Expected error from BatchRemoveOrganizationMarketplaceSubscriptions, got nil")
	}

	err = client.DeleteOrganizationMarketplaceSubscription(testOrgName, testSubscriptionID)
	if err == nil {
		t.Error("Expected error from DeleteOrganizationMarketplaceSubscription, got nil")
	}
}
