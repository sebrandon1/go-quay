package lib

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	errBuildUUIDRequired   = "buildUUID is required"
	errManifestRefRequired = "manifestRef is required"
	errUUIDRequired        = "uuid is required"
	errCodeRequired        = "code is required"
	errTagRequired         = "tag is required"
	errDigestRequired      = "manifestDigest is required"
	errTriggerUUIDRequired = "triggerUUID is required"
	errLabelIDRequired     = "labelID is required"
)

// errorPathCase describes one exported method for the shared error-path tests.
// decodes is true when the method unmarshals a response body, so a malformed
// 200 response must surface as an error.
type errorPathCase struct {
	name    string
	decodes bool
	call    func(c *Client, namespace, repository string) error
}

func newErrorPathClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	client, err := NewClientWithURL("tok", serverURL+"/api/v1")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}
	return client
}

func errorPathCases() []errorPathCase {
	ctx := context.Background()
	return []errorPathCase{
		// build.go
		{"build.GetBuilds", true, func(c *Client, ns, r string) error { _, err := c.GetBuilds(ctx, ns, r, 5); return err }},
		{"build.GetBuild", true, func(c *Client, ns, r string) error { _, err := c.GetBuild(ctx, ns, r, "uuid"); return err }},
		{"build.GetBuildLogs", true, func(c *Client, ns, r string) error { _, err := c.GetBuildLogs(ctx, ns, r, "uuid"); return err }},
		{"build.RequestBuild", true, func(c *Client, ns, r string) error {
			_, err := c.RequestBuild(ctx, ns, r, &RequestBuildRequest{})
			return err
		}},
		{"build.CancelBuild", false, func(c *Client, ns, r string) error { return c.CancelBuild(ctx, ns, r, "uuid") }},
		{"build.GetBuildStatus", true, func(c *Client, ns, r string) error { _, err := c.GetBuildStatus(ctx, ns, r, "uuid"); return err }},

		// manifest.go
		{"manifest.GetManifest", true, func(c *Client, ns, r string) error { _, err := c.GetManifest(ctx, ns, r, "sha256:abc"); return err }},
		{"manifest.DeleteManifest", false, func(c *Client, ns, r string) error { return c.DeleteManifest(ctx, ns, r, "sha256:abc") }},
		{"manifest.GetManifestLabels", true, func(c *Client, ns, r string) error {
			_, err := c.GetManifestLabels(ctx, ns, r, "sha256:abc")
			return err
		}},
		{"manifest.AddManifestLabel", true, func(c *Client, ns, r string) error {
			_, err := c.AddManifestLabel(ctx, ns, r, "sha256:abc", "key", "value", "text/plain")
			return err
		}},
		{"manifest.GetManifestLabel", true, func(c *Client, ns, r string) error {
			_, err := c.GetManifestLabel(ctx, ns, r, "sha256:abc", "label")
			return err
		}},
		{"manifest.DeleteManifestLabel", false, func(c *Client, ns, r string) error {
			return c.DeleteManifestLabel(ctx, ns, r, "sha256:abc", "label")
		}},

		// notification.go
		{"notification.GetNotifications", true, func(c *Client, ns, r string) error { _, err := c.GetNotifications(ctx, ns, r); return err }},
		{"notification.GetNotification", true, func(c *Client, ns, r string) error { _, err := c.GetNotification(ctx, ns, r, "uuid"); return err }},
		{"notification.CreateNotification", true, func(c *Client, ns, r string) error {
			_, err := c.CreateNotification(ctx, ns, r, &CreateNotificationRequest{Event: "repo_push", Method: "email"})
			return err
		}},
		{"notification.DeleteNotification", false, func(c *Client, ns, r string) error { return c.DeleteNotification(ctx, ns, r, "uuid") }},
		{"notification.TestNotification", false, func(c *Client, ns, r string) error { return c.TestNotification(ctx, ns, r, "uuid") }},
		{"notification.ResetNotification", false, func(c *Client, ns, r string) error { return c.ResetNotification(ctx, ns, r, "uuid") }},
		{"notification.UpdateNotification", true, func(c *Client, ns, r string) error {
			_, err := c.UpdateNotification(ctx, ns, r, "uuid", &CreateNotificationRequest{Event: "repo_push", Method: "email"})
			return err
		}},

		// repotoken.go
		{"repotoken.GetRepoTokens", true, func(c *Client, ns, r string) error { _, err := c.GetRepoTokens(ctx, ns, r); return err }},
		{"repotoken.CreateRepoToken", true, func(c *Client, ns, r string) error {
			_, err := c.CreateRepoToken(ctx, ns, r, &CreateRepoTokenRequest{FriendlyName: "tok"})
			return err
		}},
		{"repotoken.GetRepoToken", true, func(c *Client, ns, r string) error { _, err := c.GetRepoToken(ctx, ns, r, "code"); return err }},
		{"repotoken.UpdateRepoToken", true, func(c *Client, ns, r string) error {
			_, err := c.UpdateRepoToken(ctx, ns, r, "code", &UpdateRepoTokenRequest{Role: "read"})
			return err
		}},
		{"repotoken.DeleteRepoToken", false, func(c *Client, ns, r string) error { return c.DeleteRepoToken(ctx, ns, r, "code") }},

		// tag.go
		{"tag.GetTag", true, func(c *Client, ns, r string) error { _, err := c.GetTag(ctx, ns, r, "v1"); return err }},
		{"tag.UpdateTag", true, func(c *Client, ns, r string) error { _, err := c.UpdateTag(ctx, ns, r, "v1", "never"); return err }},
		{"tag.DeleteTag", false, func(c *Client, ns, r string) error { return c.DeleteTag(ctx, ns, r, "v1") }},
		{"tag.GetTagHistory", true, func(c *Client, ns, r string) error { _, err := c.GetTagHistory(ctx, ns, r, "v1"); return err }},
		{"tag.RevertTag", true, func(c *Client, ns, r string) error { _, err := c.RevertTag(ctx, ns, r, "v1", "sha256:abc"); return err }},
		{"tag.ChangeTag", false, func(c *Client, ns, r string) error { return c.ChangeTag(ctx, ns, r, "v1", "sha256:abc") }},
		{"tag.RestoreTag", false, func(c *Client, ns, r string) error { return c.RestoreTag(ctx, ns, r, "v1", "sha256:abc") }},

		// trigger.go
		{"trigger.GetTriggers", true, func(c *Client, ns, r string) error { _, err := c.GetTriggers(ctx, ns, r); return err }},
		{"trigger.GetTrigger", true, func(c *Client, ns, r string) error { _, err := c.GetTrigger(ctx, ns, r, "uuid"); return err }},
		{"trigger.DeleteTrigger", false, func(c *Client, ns, r string) error { return c.DeleteTrigger(ctx, ns, r, "uuid") }},
		{"trigger.UpdateTrigger", true, func(c *Client, ns, r string) error { _, err := c.UpdateTrigger(ctx, ns, r, "uuid", true); return err }},
		{"trigger.StartTriggerBuild", true, func(c *Client, ns, r string) error {
			_, err := c.StartTriggerBuild(ctx, ns, r, "uuid", &ManualTriggerRequest{})
			return err
		}},
		{"trigger.ActivateTrigger", true, func(c *Client, ns, r string) error {
			_, err := c.ActivateTrigger(ctx, ns, r, "uuid", &ActivateTriggerRequest{})
			return err
		}},
		{"trigger.GetTriggerBuilds", true, func(c *Client, ns, r string) error { _, err := c.GetTriggerBuilds(ctx, ns, r, "uuid", 5); return err }},
	}
}

func TestErrorPaths_APIStatusErrors(t *testing.T) {
	statuses := []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError}

	for _, status := range statuses {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"error":"boom"}`))
		}))

		client := newErrorPathClient(t, server.URL)
		for _, tc := range errorPathCases() {
			t.Run(tc.name+"/"+http.StatusText(status), func(t *testing.T) {
				err := tc.call(client, "testorg", "testrepo")
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var quayErr *QuayError
				if !errors.As(err, &quayErr) {
					t.Fatalf("expected *QuayError in chain, got %T: %v", err, err)
				}
				if quayErr.StatusCode() != status {
					t.Errorf("StatusCode() = %d, want %d", quayErr.StatusCode(), status)
				}
			})
		}
		server.Close()
	}
}

func TestErrorPaths_NonJSONErrorBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("<html>gateway down</html>"))
	}))
	defer server.Close()

	client := newErrorPathClient(t, server.URL)
	for _, tc := range errorPathCases() {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call(client, "testorg", "testrepo")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var quayErr *QuayError
			if errors.As(err, &quayErr) {
				t.Errorf("non-JSON body should produce a plain error, got *QuayError: %v", err)
			}
			if !strings.Contains(err.Error(), "500") {
				t.Errorf("error should mention status code 500: %v", err)
			}
		})
	}
}

func TestErrorPaths_MalformedSuccessBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{not json"))
	}))
	defer server.Close()

	client := newErrorPathClient(t, server.URL)
	for _, tc := range errorPathCases() {
		if !tc.decodes {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(client, "testorg", "testrepo"); err == nil {
				t.Fatal("expected decode error, got nil")
			}
		})
	}
}

func TestErrorPaths_RequiredArgumentGuards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("no request should be sent when a required argument is empty")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newErrorPathClient(t, server.URL)
	guards := []struct {
		name       string
		namespace  string
		repository string
		wantSubstr string
	}{
		{"empty namespace", "", "testrepo", "namespace is required"},
		{"empty repository", "testorg", "", "repository is required"},
	}

	for _, tc := range errorPathCases() {
		for _, g := range guards {
			t.Run(tc.name+"/"+g.name, func(t *testing.T) {
				err := tc.call(client, g.namespace, g.repository)
				if err == nil {
					t.Fatal("expected validation error, got nil")
				}
				if !strings.Contains(err.Error(), g.wantSubstr) {
					t.Errorf("error = %q, want substring %q", err.Error(), g.wantSubstr)
				}
			})
		}
	}
}

func TestErrorPaths_ExtraArgumentGuards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("no request should be sent when a required argument is empty")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newErrorPathClient(t, server.URL)
	ctx := context.Background()
	const ns, repo = "testorg", "testrepo"

	tests := []struct {
		name       string
		call       func() error
		wantSubstr string
	}{
		{"build.GetBuild", func() error { _, err := c.GetBuild(ctx, ns, repo, ""); return err }, errBuildUUIDRequired},
		{"build.GetBuildLogs", func() error { _, err := c.GetBuildLogs(ctx, ns, repo, ""); return err }, errBuildUUIDRequired},
		{"build.CancelBuild", func() error { return c.CancelBuild(ctx, ns, repo, "") }, errBuildUUIDRequired},
		{"build.GetBuildStatus", func() error { _, err := c.GetBuildStatus(ctx, ns, repo, ""); return err }, errBuildUUIDRequired},

		{"manifest.GetManifest", func() error { _, err := c.GetManifest(ctx, ns, repo, ""); return err }, errManifestRefRequired},
		{"manifest.DeleteManifest", func() error { return c.DeleteManifest(ctx, ns, repo, "") }, errManifestRefRequired},
		{"manifest.GetManifestLabels", func() error { _, err := c.GetManifestLabels(ctx, ns, repo, ""); return err }, errManifestRefRequired},
		{"manifest.AddManifestLabel/manifestRef", func() error { _, err := c.AddManifestLabel(ctx, ns, repo, "", "k", "v", ""); return err }, errManifestRefRequired},
		{"manifest.AddManifestLabel/key", func() error { _, err := c.AddManifestLabel(ctx, ns, repo, "sha256:abc", "", "v", ""); return err }, "key is required"},
		{"manifest.AddManifestLabel/value", func() error { _, err := c.AddManifestLabel(ctx, ns, repo, "sha256:abc", "k", "", ""); return err }, "value is required"},
		{"manifest.GetManifestLabel/manifestRef", func() error { _, err := c.GetManifestLabel(ctx, ns, repo, "", "id"); return err }, errManifestRefRequired},
		{"manifest.GetManifestLabel/labelID", func() error { _, err := c.GetManifestLabel(ctx, ns, repo, "sha256:abc", ""); return err }, errLabelIDRequired},
		{"manifest.DeleteManifestLabel/manifestRef", func() error { return c.DeleteManifestLabel(ctx, ns, repo, "", "id") }, errManifestRefRequired},
		{"manifest.DeleteManifestLabel/labelID", func() error { return c.DeleteManifestLabel(ctx, ns, repo, "sha256:abc", "") }, errLabelIDRequired},

		{"notification.GetNotification", func() error { _, err := c.GetNotification(ctx, ns, repo, ""); return err }, errUUIDRequired},
		{"notification.DeleteNotification", func() error { return c.DeleteNotification(ctx, ns, repo, "") }, errUUIDRequired},
		{"notification.TestNotification", func() error { return c.TestNotification(ctx, ns, repo, "") }, errUUIDRequired},
		{"notification.ResetNotification", func() error { return c.ResetNotification(ctx, ns, repo, "") }, errUUIDRequired},
		{"notification.UpdateNotification", func() error {
			_, err := c.UpdateNotification(ctx, ns, repo, "", &CreateNotificationRequest{})
			return err
		}, errUUIDRequired},

		{"repotoken.GetRepoToken", func() error { _, err := c.GetRepoToken(ctx, ns, repo, ""); return err }, errCodeRequired},
		{"repotoken.UpdateRepoToken", func() error { _, err := c.UpdateRepoToken(ctx, ns, repo, "", &UpdateRepoTokenRequest{}); return err }, errCodeRequired},
		{"repotoken.DeleteRepoToken", func() error { return c.DeleteRepoToken(ctx, ns, repo, "") }, errCodeRequired},

		{"tag.GetTag", func() error { _, err := c.GetTag(ctx, ns, repo, ""); return err }, errTagRequired},
		{"tag.UpdateTag", func() error { _, err := c.UpdateTag(ctx, ns, repo, "", "never"); return err }, errTagRequired},
		{"tag.DeleteTag", func() error { return c.DeleteTag(ctx, ns, repo, "") }, errTagRequired},
		{"tag.GetTagHistory", func() error { _, err := c.GetTagHistory(ctx, ns, repo, ""); return err }, errTagRequired},
		{"tag.RevertTag/tag", func() error { _, err := c.RevertTag(ctx, ns, repo, "", "sha256:abc"); return err }, errTagRequired},
		{"tag.RevertTag/digest", func() error { _, err := c.RevertTag(ctx, ns, repo, "v1", ""); return err }, errDigestRequired},
		{"tag.ChangeTag/tag", func() error { return c.ChangeTag(ctx, ns, repo, "", "sha256:abc") }, errTagRequired},
		{"tag.ChangeTag/digest", func() error { return c.ChangeTag(ctx, ns, repo, "v1", "") }, errDigestRequired},
		{"tag.RestoreTag/tag", func() error { return c.RestoreTag(ctx, ns, repo, "", "sha256:abc") }, errTagRequired},
		{"tag.RestoreTag/digest", func() error { return c.RestoreTag(ctx, ns, repo, "v1", "") }, errDigestRequired},

		{"trigger.GetTrigger", func() error { _, err := c.GetTrigger(ctx, ns, repo, ""); return err }, errTriggerUUIDRequired},
		{"trigger.DeleteTrigger", func() error { return c.DeleteTrigger(ctx, ns, repo, "") }, errTriggerUUIDRequired},
		{"trigger.UpdateTrigger", func() error { _, err := c.UpdateTrigger(ctx, ns, repo, "", true); return err }, errTriggerUUIDRequired},
		{"trigger.StartTriggerBuild", func() error {
			_, err := c.StartTriggerBuild(ctx, ns, repo, "", &ManualTriggerRequest{})
			return err
		}, errTriggerUUIDRequired},
		{"trigger.ActivateTrigger", func() error {
			_, err := c.ActivateTrigger(ctx, ns, repo, "", &ActivateTriggerRequest{})
			return err
		}, errTriggerUUIDRequired},
		{"trigger.GetTriggerBuilds", func() error { _, err := c.GetTriggerBuilds(ctx, ns, repo, "", 5); return err }, errTriggerUUIDRequired},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantSubstr) {
				t.Errorf("error = %q, want substring %q", err.Error(), tc.wantSubstr)
			}
		})
	}
}
