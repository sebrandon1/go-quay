package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrganizationBilling(t *testing.T) {
	mockResponse := BillingInfo{
		Plan:     testBillingPlanFree,
		PlanType: testBillingPlanTypeFree,
		IsActive: true,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/v1/organization/" + testNamespace + "/plan"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+testTokenValue {
			t.Errorf("Expected Authorization header 'Bearer %s', got '%s'", testTokenValue, authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	billing, err := client.GetOrganizationBilling(testNamespace)
	if err != nil {
		t.Fatalf("GetOrganizationBilling returned error: %v", err)
	}

	if billing.Plan != testBillingPlanFree {
		t.Errorf("Expected plan '%s', got '%s'", testBillingPlanFree, billing.Plan)
	}

	if billing.PlanType != testBillingPlanTypeFree {
		t.Errorf("Expected plan type '%s', got '%s'", testBillingPlanTypeFree, billing.PlanType)
	}

	if !billing.IsActive {
		t.Error("Expected billing to be active")
	}
}

func TestGetUserBilling(t *testing.T) {
	mockResponse := BillingInfo{
		Plan:        testBillingPlanFree,
		PlanType:    testBillingPlanTypeFree,
		IsActive:    true,
		UsedPrivate: 2,
		PlanPrivate: 5,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/v1/user/plan"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+testTokenValue {
			t.Errorf("Expected Authorization header 'Bearer %s', got '%s'", testTokenValue, authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	billing, err := client.GetUserBilling()
	if err != nil {
		t.Fatalf("GetUserBilling returned error: %v", err)
	}

	if billing.Plan != testBillingPlanFree {
		t.Errorf("Expected plan '%s', got '%s'", testBillingPlanFree, billing.Plan)
	}

	if billing.UsedPrivate != 2 {
		t.Errorf("Expected used private 2, got %d", billing.UsedPrivate)
	}

	if billing.PlanPrivate != 5 {
		t.Errorf("Expected plan private 5, got %d", billing.PlanPrivate)
	}
}

func TestGetOrganizationSubscription(t *testing.T) {
	mockResponse := Subscription{
		ID:       testBillingSubscriptionID,
		Name:     "Business Plan",
		Price:    5000,
		Currency: testBillingCurrencyUSD,
		Period:   testBillingPeriodMonthly,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/v1/organization/" + testNamespace + "/plan"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	subscription, err := client.GetOrganizationSubscription(testNamespace)
	if err != nil {
		t.Fatalf("GetOrganizationSubscription returned error: %v", err)
	}

	if subscription.ID != testBillingSubscriptionID {
		t.Errorf("Expected subscription ID '%s', got '%s'", testBillingSubscriptionID, subscription.ID)
	}

	if subscription.Price != 5000 {
		t.Errorf("Expected price 5000, got %d", subscription.Price)
	}

	if subscription.Currency != testBillingCurrencyUSD {
		t.Errorf("Expected currency '%s', got '%s'", testBillingCurrencyUSD, subscription.Currency)
	}
}

func TestGetUserSubscription(t *testing.T) {
	mockResponse := Subscription{
		ID:       testBillingSubscriptionID,
		Name:     "Personal Plan",
		Price:    1500,
		Currency: testBillingCurrencyUSD,
		Period:   testBillingPeriodMonthly,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/v1/user/plan"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	subscription, err := client.GetUserSubscription()
	if err != nil {
		t.Fatalf("GetUserSubscription returned error: %v", err)
	}

	if subscription.ID != testBillingSubscriptionID {
		t.Errorf("Expected subscription ID '%s', got '%s'", testBillingSubscriptionID, subscription.ID)
	}

	if subscription.Price != 1500 {
		t.Errorf("Expected price 1500, got %d", subscription.Price)
	}
}

func TestGetOrganizationInvoices(t *testing.T) {
	mockResponse := struct {
		Invoices []Invoice `json:"invoices"`
	}{
		Invoices: []Invoice{
			{
				ID:       testBillingInvoiceID,
				Number:   "INV-001",
				Status:   testBillingInvoiceStatusPaid,
				Amount:   5000,
				Currency: testBillingCurrencyUSD,
			},
			{
				ID:       "inv-456",
				Number:   "INV-002",
				Status:   "pending",
				Amount:   5000,
				Currency: testBillingCurrencyUSD,
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/v1/organization/" + testNamespace + "/invoices"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	invoices, err := client.GetOrganizationInvoices(testNamespace)
	if err != nil {
		t.Fatalf("GetOrganizationInvoices returned error: %v", err)
	}

	if len(invoices) != 2 {
		t.Fatalf("Expected 2 invoices, got %d", len(invoices))
	}

	if invoices[0].ID != testBillingInvoiceID {
		t.Errorf("Expected first invoice ID '%s', got '%s'", testBillingInvoiceID, invoices[0].ID)
	}

	if invoices[0].Status != testBillingInvoiceStatusPaid {
		t.Errorf("Expected first invoice status '%s', got '%s'", testBillingInvoiceStatusPaid, invoices[0].Status)
	}

	if invoices[0].Amount != 5000 {
		t.Errorf("Expected first invoice amount 5000, got %d", invoices[0].Amount)
	}
}

func TestGetUserInvoices(t *testing.T) {
	client, err := NewClient(testTokenValue)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	invoices, err := client.GetUserInvoices()
	if err == nil {
		t.Error("Expected error for unsupported endpoint, got nil")
	}

	if invoices != nil {
		t.Error("Expected nil invoices for unsupported endpoint")
	}
}

func TestGetAvailablePlans(t *testing.T) {
	mockResponse := struct {
		Plans []Subscription `json:"plans"`
	}{
		Plans: []Subscription{
			{
				ID:       "plan-free",
				Name:     "Free",
				Price:    0,
				Currency: testBillingCurrencyUSD,
				Period:   testBillingPeriodMonthly,
			},
			{
				ID:       "plan-business",
				Name:     "Business",
				Price:    5000,
				Currency: testBillingCurrencyUSD,
				Period:   testBillingPeriodMonthly,
				Features: []string{"private repos", "team management"},
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/v1/plans"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	plans, err := client.GetAvailablePlans()
	if err != nil {
		t.Fatalf("GetAvailablePlans returned error: %v", err)
	}

	if len(plans) != 2 {
		t.Fatalf("Expected 2 plans, got %d", len(plans))
	}

	if plans[0].Name != "Free" {
		t.Errorf("Expected first plan name 'Free', got '%s'", plans[0].Name)
	}

	if plans[0].Price != 0 {
		t.Errorf("Expected first plan price 0, got %d", plans[0].Price)
	}

	if plans[1].Name != "Business" {
		t.Errorf("Expected second plan name 'Business', got '%s'", plans[1].Name)
	}

	if len(plans[1].Features) != 2 {
		t.Errorf("Expected 2 features for Business plan, got %d", len(plans[1].Features))
	}
}

func TestGetOrganizationBillingError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	billing, err := client.GetOrganizationBilling(testNamespace)
	if err == nil {
		t.Error("Expected error for 500 response, got nil")
	}

	if billing != nil {
		t.Error("Expected nil billing on error")
	}
}

func TestGetUserBillingError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	billing, err := client.GetUserBilling()
	if err == nil {
		t.Error("Expected error for 401 response, got nil")
	}

	if billing != nil {
		t.Error("Expected nil billing on error")
	}
}

func TestGetAvailablePlansError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL(testTokenValue, server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	plans, err := client.GetAvailablePlans()
	if err == nil {
		t.Error("Expected error for 404 response, got nil")
	}

	if plans != nil {
		t.Error("Expected nil plans on error")
	}
}
