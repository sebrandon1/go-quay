package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testSecScanManifestRef = "sha256:abc123def456789"
	httpGetSecScan         = "GET"
)

func TestGetManifestSecurity(t *testing.T) {
	mockSecurityScan := SecurityScan{
		Status: "scanned",
		Data: &SecurityData{
			Layer: &SecurityLayer{
				Name:             "sha256:abc123",
				NamespaceName:    "centos:8",
				IndexedByVersion: 4,
				Features: []SecurityFeature{
					{
						Name:          "openssl",
						Version:       "1.1.1k",
						VersionFormat: "rpm",
						NamespaceName: "centos:8",
						AddedBy:       "sha256:layer1",
						Vulnerabilities: []SecurityVulnerability{
							{
								Name:          "CVE-2021-3712",
								NamespaceName: "centos:8",
								Description:   "OpenSSL vulnerability in ASN.1 parsing",
								Link:          "https://access.redhat.com/security/cve/CVE-2021-3712",
								Severity:      "Medium",
								FixedBy:       "1.1.1k-5",
							},
						},
					},
					{
						Name:          "curl",
						Version:       "7.61.1",
						VersionFormat: "rpm",
						NamespaceName: "centos:8",
						AddedBy:       "sha256:layer2",
					},
				},
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockSecurityScan)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetSecScan {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/repository/testorg/testrepo/manifest/" + testSecScanManifestRef + "/security"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Check for vulnerabilities query parameter
		if r.URL.Query().Get("vulnerabilities") != "true" {
			t.Errorf("Expected vulnerabilities=true query parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

	securityScan, err := client.GetManifestSecurity("testorg", "testrepo", testSecScanManifestRef, true)
	if err != nil {
		t.Fatalf("GetManifestSecurity failed: %v", err)
	}

	if securityScan.Status != "scanned" {
		t.Errorf("Expected status 'scanned', got '%s'", securityScan.Status)
	}
	if securityScan.Data == nil {
		t.Fatal("Expected Data to be non-nil")
	}
	if securityScan.Data.Layer == nil {
		t.Fatal("Expected Layer to be non-nil")
	}
	if len(securityScan.Data.Layer.Features) != 2 {
		t.Errorf("Expected 2 features, got %d", len(securityScan.Data.Layer.Features))
	}

	// Check first feature with vulnerability
	feature := securityScan.Data.Layer.Features[0]
	if feature.Name != "openssl" {
		t.Errorf("Expected feature name 'openssl', got '%s'", feature.Name)
	}
	if len(feature.Vulnerabilities) != 1 {
		t.Errorf("Expected 1 vulnerability, got %d", len(feature.Vulnerabilities))
	}
	if feature.Vulnerabilities[0].Name != "CVE-2021-3712" {
		t.Errorf("Expected vulnerability 'CVE-2021-3712', got '%s'", feature.Vulnerabilities[0].Name)
	}
	if feature.Vulnerabilities[0].Severity != "Medium" {
		t.Errorf("Expected severity 'Medium', got '%s'", feature.Vulnerabilities[0].Severity)
	}
}

func TestGetManifestSecurityWithoutVulnerabilities(t *testing.T) {
	mockSecurityScan := SecurityScan{
		Status: "scanned",
		Data: &SecurityData{
			Layer: &SecurityLayer{
				Name:             "sha256:abc123",
				IndexedByVersion: 4,
			},
		},
	}

	mockResponseJSON, _ := json.Marshal(mockSecurityScan)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetSecScan {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check that vulnerabilities query parameter is NOT present
		if r.URL.Query().Get("vulnerabilities") == "true" {
			t.Errorf("Did not expect vulnerabilities=true query parameter")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

	securityScan, err := client.GetManifestSecurity("testorg", "testrepo", testSecScanManifestRef, false)
	if err != nil {
		t.Fatalf("GetManifestSecurity failed: %v", err)
	}

	if securityScan.Status != "scanned" {
		t.Errorf("Expected status 'scanned', got '%s'", securityScan.Status)
	}
}

func TestGetManifestSecurityQueued(t *testing.T) {
	// Test when scan is still queued/in progress
	mockSecurityScan := SecurityScan{
		Status: "queued",
	}

	mockResponseJSON, _ := json.Marshal(mockSecurityScan)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

	securityScan, err := client.GetManifestSecurity("testorg", "testrepo", testSecScanManifestRef, true)
	if err != nil {
		t.Fatalf("GetManifestSecurity failed: %v", err)
	}

	if securityScan.Status != "queued" {
		t.Errorf("Expected status 'queued', got '%s'", securityScan.Status)
	}
	if securityScan.Data != nil {
		t.Error("Expected Data to be nil for queued scan")
	}
}

func TestGetManifestSecurityErrorHandling(t *testing.T) {
	// Test 404 error for non-existent manifest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Manifest not found"}`))
	}))
	defer server.Close()

	originalURL := QuayURL
	QuayURL = server.URL + "/api/v1"
	defer func() { QuayURL = originalURL }()

	client, err := NewClient("test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetManifestSecurity("testorg", "testrepo", "nonexistent", true)
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}
}

func TestGetManifestSecurityUnsupported(t *testing.T) {
	// Test when security scanning is not supported for the image
	mockSecurityScan := SecurityScan{
		Status: "unsupported",
	}

	mockResponseJSON, _ := json.Marshal(mockSecurityScan)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

	securityScan, err := client.GetManifestSecurity("testorg", "testrepo", testSecScanManifestRef, true)
	if err != nil {
		t.Fatalf("GetManifestSecurity failed: %v", err)
	}

	if securityScan.Status != "unsupported" {
		t.Errorf("Expected status 'unsupported', got '%s'", securityScan.Status)
	}
}
