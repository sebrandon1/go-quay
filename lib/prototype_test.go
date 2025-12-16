package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpGetPrototype    = "GET"
	httpPostPrototype   = "POST"
	httpPutPrototype    = "PUT"
	httpDeletePrototype = "DELETE"

	testPrototypeOrg  = "testorg"
	testPrototypeUUID = "proto-uuid-123"
	prototypeRoleRead = "read"
)

func TestGetPrototypes(t *testing.T) {
	mockResponse := Prototypes{
		Prototypes: []Prototype{
			{
				ID:   testPrototypeUUID,
				Role: prototypeRoleRead,
				Delegate: PrototypeDelegate{
					Name: "devteam",
					Kind: "team",
				},
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetPrototype {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testPrototypeOrg + "/prototypes"
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

	prototypes, err := client.GetPrototypes(testPrototypeOrg)
	if err != nil {
		t.Fatalf("GetPrototypes returned error: %v", err)
	}

	if len(prototypes.Prototypes) != 1 {
		t.Errorf("Expected 1 prototype, got %d", len(prototypes.Prototypes))
	}
	if prototypes.Prototypes[0].ID != testPrototypeUUID {
		t.Errorf("Expected prototype ID %s, got %s", testPrototypeUUID, prototypes.Prototypes[0].ID)
	}
}

func TestCreatePrototype(t *testing.T) {
	mockResponse := Prototype{
		ID:   testPrototypeUUID,
		Role: prototypeRoleRead,
		Delegate: PrototypeDelegate{
			Name: "devteam",
			Kind: "team",
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPostPrototype {
			t.Errorf("Expected POST request, got %s", r.Method)
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

	createReq := &CreatePrototypeRequest{
		Delegate: PrototypeDelegateRequest{
			Name: "devteam",
			Kind: "team",
		},
		Role: prototypeRoleRead,
	}

	prototype, err := client.CreatePrototype(testPrototypeOrg, createReq)
	if err != nil {
		t.Fatalf("CreatePrototype returned error: %v", err)
	}

	if prototype.ID != testPrototypeUUID {
		t.Errorf("Expected prototype ID %s, got %s", testPrototypeUUID, prototype.ID)
	}
}

func TestGetPrototype(t *testing.T) {
	mockResponse := Prototype{
		ID:   testPrototypeUUID,
		Role: prototypeRoleRead,
		Delegate: PrototypeDelegate{
			Name: "devteam",
			Kind: "team",
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpGetPrototype {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testPrototypeOrg + "/prototypes/" + testPrototypeUUID
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

	prototype, err := client.GetPrototype(testPrototypeOrg, testPrototypeUUID)
	if err != nil {
		t.Fatalf("GetPrototype returned error: %v", err)
	}

	if prototype.ID != testPrototypeUUID {
		t.Errorf("Expected prototype ID %s, got %s", testPrototypeUUID, prototype.ID)
	}
}

func TestUpdatePrototype(t *testing.T) {
	mockResponse := Prototype{
		ID:   testPrototypeUUID,
		Role: "write",
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpPutPrototype {
			t.Errorf("Expected PUT request, got %s", r.Method)
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

	updateReq := &UpdatePrototypeRequest{
		Role: "write",
	}

	prototype, err := client.UpdatePrototype(testPrototypeOrg, testPrototypeUUID, updateReq)
	if err != nil {
		t.Fatalf("UpdatePrototype returned error: %v", err)
	}

	if prototype.Role != "write" {
		t.Errorf("Expected role 'write', got %s", prototype.Role)
	}
}

func TestDeletePrototype(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpDeletePrototype {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testPrototypeOrg + "/prototypes/" + testPrototypeUUID
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

	err = client.DeletePrototype(testPrototypeOrg, testPrototypeUUID)
	if err != nil {
		t.Fatalf("DeletePrototype returned error: %v", err)
	}
}

func TestGetPrototypesError(t *testing.T) {
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

	_, err = client.GetPrototypes(testPrototypeOrg)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
