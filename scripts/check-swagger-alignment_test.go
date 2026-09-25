package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestScanSourceEndpointsHTTPMethods(t *testing.T) {
	const sourceFile = "fixture.go"

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
	if err := os.WriteFile(filepath.Join(libPath, sourceFile), []byte(source), 0600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	got, err := scanSourceEndpoints(libPath, "BaseURL")
	if err != nil {
		t.Fatalf("scan source endpoints: %v", err)
	}

	want := []ImplementedEndpoint{
		{Method: "GET", Path: "/get", SourceFile: sourceFile, Function: "Get"},
		{Method: "POST", Path: "/create", SourceFile: sourceFile, Function: "Create"},
		{Method: "DELETE", Path: "/delete", SourceFile: sourceFile, Function: "Delete"},
		{Method: "GET", Path: "/default", SourceFile: sourceFile, Function: "DefaultMethod"},
		{Method: "PATCH", Path: "/patch", SourceFile: sourceFile, Function: "StandardLibraryRequest"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scanSourceEndpoints() = %#v, want %#v", got, want)
	}
}
