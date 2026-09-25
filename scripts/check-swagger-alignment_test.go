package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestScanSourceEndpointsHTTPMethods(t *testing.T) {
	libPath := t.TempDir()
	source := `package fixture

func Get(ctx context.Context) {
	_, _ = newRequest(ctx, http.MethodGet, c.buildURL("/get"), nil)
}

func Create(ctx context.Context) {
	_, _ = newRequestWithBody(ctx, http.MethodPost,
		c.buildURL("/create"), body)
}

func Delete(ctx context.Context) {
	_, _ = newRequest(ctx, http.MethodDelete, c.buildURL("/delete"), nil)
}

func DefaultMethod(ctx context.Context) {
	url := c.buildURL("/default")
	_ = url
}

func StandardLibraryRequest(ctx context.Context) {
	_, _ = http.NewRequest("PATCH", c.buildURL("/patch"), nil)
}
`
	if err := os.WriteFile(filepath.Join(libPath, "fixture.go"), []byte(source), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	got, err := scanSourceEndpoints(libPath, "BaseURL")
	if err != nil {
		t.Fatalf("scan source endpoints: %v", err)
	}

	want := []ImplementedEndpoint{
		{Method: "GET", Path: "/get", SourceFile: "fixture.go", Function: "Get"},
		{Method: "POST", Path: "/create", SourceFile: "fixture.go", Function: "Create"},
		{Method: "DELETE", Path: "/delete", SourceFile: "fixture.go", Function: "Delete"},
		{Method: "GET", Path: "/default", SourceFile: "fixture.go", Function: "DefaultMethod"},
		{Method: "PATCH", Path: "/patch", SourceFile: "fixture.go", Function: "StandardLibraryRequest"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scanSourceEndpoints() = %#v, want %#v", got, want)
	}
}
