package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestFetchSwaggerEndpoints(t *testing.T) {
	tests := []struct {
		name            string
		status          int
		response        string
		invalidURL      bool
		canceledContext bool
		want            []Endpoint
		wantErr         string
		wantCanceled    bool
	}{
		{
			name:     "success",
			response: `{"paths":{"/z":{"post":{"summary":"Create"}},"/a":{"get":{"summary":"List","tags":["repositories"]}}}}`,
			want: []Endpoint{
				{Method: http.MethodGet, Path: "/a", Tags: []string{"repositories"}, Summary: "List"},
				{Method: http.MethodPost, Path: "/z", Summary: "Create"},
			},
		},
		{
			name:     "http error",
			status:   http.StatusServiceUnavailable,
			response: "unavailable",
			wantErr:  "HTTP 503: unavailable",
		},
		{
			name:     "invalid JSON",
			response: "{",
			wantErr:  "failed to parse swagger spec",
		},
		{
			name:       "invalid URL",
			invalidURL: true,
			wantErr:    "invalid control character",
		},
		{
			name:            "canceled context",
			canceledContext: true,
			wantCanceled:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := ""
			if !tt.invalidURL {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodGet {
						t.Errorf("request method = %q, want GET", r.Method)
					}
					status := tt.status
					if status == 0 {
						status = http.StatusOK
					}
					w.WriteHeader(status)
					_, _ = w.Write([]byte(tt.response))
				}))
				t.Cleanup(server.Close)
				url = server.URL
			} else {
				url = "\x00"
			}

			ctx := context.Background()
			if tt.canceledContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			got, err := fetchSwaggerEndpoints(ctx, url)
			if tt.wantCanceled {
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("error = %v, want context.Canceled", err)
				}
				return
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want error containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("fetchSwaggerEndpoints() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("fetchSwaggerEndpoints() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

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
	_, _ = http.NewRequestWithContext(ctx, "PATCH", c.buildURL("/patch"), nil)
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
