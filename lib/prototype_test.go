package lib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testPrototypeUUID = "proto-uuid-123"
)

func TestGetPrototypes(t *testing.T) {
	mockResponse := Prototypes{
		Prototypes: []Prototype{
			{
				ID:   testPrototypeUUID,
				Role: testRoleRead,
				Delegate: PrototypeDelegate{
					Name: testPrototypeTeamName,
					Kind: testKindTeam,
				},
			},
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testNamespace + "/prototypes"
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

	prototypes, err := client.GetPrototypes(testNamespace)
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
		Role: testRoleRead,
		Delegate: PrototypeDelegate{
			Name: testPrototypeTeamName,
			Kind: testKindTeam,
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
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

	createReq := &CreatePrototypeRequest{
		Delegate: PrototypeDelegateRequest{
			Name: testPrototypeTeamName,
			Kind: testKindTeam,
		},
		Role: testRoleRead,
	}

	prototype, err := client.CreatePrototype(testNamespace, createReq)
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
		Role: testRoleRead,
		Delegate: PrototypeDelegate{
			Name: testPrototypeTeamName,
			Kind: testKindTeam,
		},
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testNamespace + "/prototypes/" + testPrototypeUUID
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

	prototype, err := client.GetPrototype(testNamespace, testPrototypeUUID)
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
		Role: testRoleWrite,
	}
	mockResponseJSON, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResponseJSON)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	updateReq := &UpdatePrototypeRequest{
		Role: testRoleWrite,
	}

	prototype, err := client.UpdatePrototype(testNamespace, testPrototypeUUID, updateReq)
	if err != nil {
		t.Fatalf("UpdatePrototype returned error: %v", err)
	}

	if prototype.Role != testRoleWrite {
		t.Errorf("Expected role 'write', got %s", prototype.Role)
	}
}

func TestDeletePrototype(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/api/v1/organization/" + testNamespace + "/prototypes/" + testPrototypeUUID
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

	err = client.DeletePrototype(testNamespace, testPrototypeUUID)
	if err != nil {
		t.Fatalf("DeletePrototype returned error: %v", err)
	}
}

func TestGetPrototypesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := NewClientWithURL("test-token", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetPrototypes(testNamespace)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
