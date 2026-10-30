package lib

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildAPIErrorShapes(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		body        string
		wantTyped   bool
		wantMessage string
		wantCode    int
		check       func(error) bool
	}{
		{
			name:   "auth body uses detail and error_message",
			status: http.StatusUnauthorized,
			body: `{"detail":"Requires authentication","error_message":"Requires authentication",` +
				`"error_type":"invalid_token","status":401}`,
			wantTyped:   true,
			wantMessage: testAuthMessage,
			wantCode:    http.StatusUnauthorized,
			check:       IsUnauthorized,
		},
		{
			name:        "error body keeps using error",
			status:      http.StatusInternalServerError,
			body:        `{"error":"boom"}`,
			wantTyped:   true,
			wantMessage: "boom",
			wantCode:    http.StatusInternalServerError,
		},
		{
			name:        "detail only",
			status:      http.StatusNotFound,
			body:        `{"detail":"Not found"}`,
			wantTyped:   true,
			wantMessage: "Not found",
			wantCode:    http.StatusNotFound,
			check:       IsNotFound,
		},
		{
			name:        "error_type only",
			status:      http.StatusUnauthorized,
			body:        `{"error_type":"invalid_token"}`,
			wantTyped:   true,
			wantMessage: "invalid_token",
			wantCode:    http.StatusUnauthorized,
		},
		{
			name:        "empty object falls back to the status text",
			status:      http.StatusForbidden,
			body:        `{}`,
			wantTyped:   true,
			wantMessage: http.StatusText(http.StatusForbidden),
			wantCode:    http.StatusForbidden,
			check:       IsForbidden,
		},
		{
			name:      "non-json body stays a plain error",
			status:    http.StatusBadGateway,
			body:      `<html>bad gateway</html>`,
			wantTyped: false,
		},
		{
			name:      "json array stays a plain error",
			status:    http.StatusBadRequest,
			body:      `[]`,
			wantTyped: false,
		},
		{
			name:      "json scalar stays a plain error",
			status:    http.StatusBadRequest,
			body:      `"plain"`,
			wantTyped: false,
		},
		{
			name:        "json null becomes a typed error with the status text",
			status:      http.StatusBadGateway,
			body:        `null`,
			wantTyped:   true,
			wantMessage: http.StatusText(http.StatusBadGateway),
			wantCode:    http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := buildAPIError(tt.status, tt.body)

			var quayErr *QuayError
			typed := errors.As(err, &quayErr)
			if typed != tt.wantTyped {
				t.Fatalf("typed = %v, want %v (err: %v)", typed, tt.wantTyped, err)
			}
			if !tt.wantTyped {
				if _, ok := StatusCode(err); ok {
					t.Error("StatusCode() reported a code for a non-API error")
				}

				return
			}

			if quayErr.StatusCode() != tt.wantCode {
				t.Errorf("StatusCode() = %d, want %d", quayErr.StatusCode(), tt.wantCode)
			}
			if quayErr.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", quayErr.Message, tt.wantMessage)
			}
			if code, ok := StatusCode(err); !ok || code != tt.wantCode {
				t.Errorf("StatusCode(err) = %d/%v, want %d/true", code, ok, tt.wantCode)
			}
			if tt.check != nil && !tt.check(err) {
				t.Errorf("helper check failed for status %d", tt.wantCode)
			}

			assertHelpersRejectOtherStatuses(t, err, tt.wantCode)
		})
	}
}

// assertHelpersRejectOtherStatuses checks that a status helper does not match a
// status other than its own.
func assertHelpersRejectOtherStatuses(t *testing.T, err error, wantCode int) {
	t.Helper()

	if wantCode != http.StatusNotFound && IsNotFound(err) {
		t.Error("IsNotFound() = true for a different status")
	}
	if wantCode != http.StatusUnauthorized && IsUnauthorized(err) {
		t.Error("IsUnauthorized() = true for a different status")
	}
	if wantCode != http.StatusForbidden && IsForbidden(err) {
		t.Error("IsForbidden() = true for a different status")
	}
}

func TestGetRepositoryReportsTruncatedTags(t *testing.T) {
	repoJSON := mustMarshal(t, Repository{Namespace: testNamespace, Name: testRepository})
	tagsJSON := mustMarshal(t, RepositoryTags{
		Tags:          []Tag{{Name: testTagNameLatest}},
		Page:          1,
		HasAdditional: true,
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case testRepoPath:
			_, _ = w.Write(repoJSON)
		case testTagPath:
			_, _ = w.Write(tagsJSON)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := NewClientWithURL("", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("NewClientWithURL() error = %v", err)
	}

	repo, err := client.GetRepository(context.Background(), testNamespace, testRepository)
	if err != nil {
		t.Fatalf("GetRepository() error = %v", err)
	}

	if !repo.TagsTruncated() {
		t.Error("TagsTruncated() = false, want true when has_additional is set")
	}
	if len(repo.Tags.Tags) != 1 {
		t.Errorf("len(Tags) = %d, want 1", len(repo.Tags.Tags))
	}
}

func TestAPIErrorHelpersRejectForeignErrors(t *testing.T) {
	plain := errors.New("dial tcp: connection refused")

	for name, check := range map[string]func(error) bool{
		"IsNotFound":     IsNotFound,
		"IsUnauthorized": IsUnauthorized,
		"IsForbidden":    IsForbidden,
	} {
		if check(plain) {
			t.Errorf("%s = true for a non-Quay error", name)
		}
		if check(nil) {
			t.Errorf("%s = true for a nil error", name)
		}
	}

	if _, ok := StatusCode(plain); ok {
		t.Error("StatusCode() = true for a non-Quay error")
	}
}

func TestQuayErrorStringDoesNotRepeatDetail(t *testing.T) {
	repeated := &QuayError{Status: 401, Message: testAuthMessage, Detail: testAuthMessage}
	if got, want := repeated.Error(), "quay API error (status 401): Requires authentication"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	distinct := &QuayError{Status: 404, Message: testNotFoundErrorCode, Detail: "Repository not found"}
	want := "quay API error (status 404): not_found — Repository not found"
	if got := distinct.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestClientReturnsTypedErrorForAuthShape(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"Requires authentication","error_message":"Requires authentication",` +
			`"error_type":"invalid_token"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL("", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("NewClientWithURL() error = %v", err)
	}

	_, err = client.GetRepository(context.Background(), "redhat", "ubi9")
	if err == nil {
		t.Fatal("GetRepository() error = nil, want an API error")
	}

	if !IsUnauthorized(err) {
		t.Errorf("IsUnauthorized() = false for %v", err)
	}
	if code, ok := StatusCode(err); !ok || code != http.StatusUnauthorized {
		t.Errorf("StatusCode() = %d/%v, want 401/true", code, ok)
	}
}

func TestStatusCodeSurvivesRetryableWrapper(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"try later"}`))
	}))
	defer server.Close()

	client, err := NewClientWithURL("", server.URL+"/api/v1")
	if err != nil {
		t.Fatalf("NewClientWithURL() error = %v", err)
	}

	_, err = client.GetRepository(context.Background(), "redhat", "ubi9")
	if err == nil {
		t.Fatal("GetRepository() error = nil, want an API error")
	}

	if code, ok := StatusCode(err); !ok || code != http.StatusServiceUnavailable {
		t.Errorf("StatusCode() = %d/%v, want 503/true", code, ok)
	}
}

func TestRepositoryWithTagsTagsTruncated(t *testing.T) {
	partial := &RepositoryWithTags{Tags: RepositoryTags{HasAdditional: true}}
	if !partial.TagsTruncated() {
		t.Error("TagsTruncated() = false, want true when the API reports more tags")
	}

	complete := &RepositoryWithTags{Tags: RepositoryTags{HasAdditional: false}}
	if complete.TagsTruncated() {
		t.Error("TagsTruncated() = true, want false for a complete list")
	}
}
