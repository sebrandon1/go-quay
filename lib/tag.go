/*
Package lib provides Quay.io API client functionality.

This file covers ENHANCED TAG operations:

Tag Management:
  - GET    /api/v1/repository/{namespace}/{repository}/tag/{tag}                 - GetTag()
  - PUT    /api/v1/repository/{namespace}/{repository}/tag/{tag}                 - UpdateTag()
  - PUT    /api/v1/repository/{namespace}/{repository}/tag/{tag}                 - ChangeTag()
  - DELETE /api/v1/repository/{namespace}/{repository}/tag/{tag}                 - DeleteTag()
  - GET    /api/v1/repository/{namespace}/{repository}/tag/{tag}/history         - GetTagHistory()

Enhanced tag operations provide detailed metadata, history tracking, and
expiration management for individual tags.
*/
package lib

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// TagExpiringWindow is how close to expiration a tag must be to be reported as
// expiring rather than active.
const TagExpiringWindow = 24 * time.Hour

// TagLifecycle is where a tag sits relative to its expiration.
type TagLifecycle string

const (
	TagActive   TagLifecycle = "active"
	TagExpiring TagLifecycle = "expiring"
	TagExpired  TagLifecycle = "expired"
)

// Lifecycle reports the tag's lifecycle at now. A tag with no usable
// expiration is active: the client does not guess that it has expired.
func (t Tag) Lifecycle(now time.Time) TagLifecycle {
	expires, ok := t.expiry()
	if !ok {
		return TagActive
	}
	if !now.Before(expires) {
		return TagExpired
	}
	if expires.Sub(now) <= TagExpiringWindow {
		return TagExpiring
	}
	return TagActive
}

// Expired reports whether the tag has reached its expiration at now.
func (t Tag) Expired(now time.Time) bool {
	return t.Lifecycle(now) == TagExpired
}

// expiry returns the tag's expiration, preferring the dated Expiration field
// and falling back to the EndTs epoch.
func (t Tag) expiry() (time.Time, bool) {
	if t.Expiration != "" {
		for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02"} {
			if parsed, err := time.Parse(layout, t.Expiration); err == nil {
				return parsed, true
			}
		}
	}
	if t.EndTs > 0 {
		return time.Unix(t.EndTs, 0), true
	}
	return time.Time{}, false
}

// GetTag retrieves detailed information about a specific tag
func (c *Client) GetTag(ctx context.Context, namespace, repository, tag string) (*Tag, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/tag/%s", namespace, repository, tag))
	if err != nil {
		return nil, fmt.Errorf("failed to create get tag request: %w", err)
	}

	var tagInfo Tag
	if err := c.get(req, &tagInfo); err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return &tagInfo, nil
}

// UpdateTag updates tag metadata (currently supports expiration)
func (c *Client) UpdateTag(ctx context.Context, namespace, repository, tag, expiration string) (*Tag, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag is required")
	}

	updateReq := UpdateTagRequest{}

	if expiration != "" {
		updateReq.Expiration = expiration
	}

	req, err := newRequestWithBody(ctx, http.MethodPut, c.buildURL("/repository/%s/%s/tag/%s", namespace, repository, tag), updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create update tag request: %w", err)
	}

	var tagInfo Tag
	if err := c.put(req, &tagInfo); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return &tagInfo, nil
}

// DeleteTag deletes a specific tag from a repository
func (c *Client) DeleteTag(ctx context.Context, namespace, repository, tag string) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}
	if tag == "" {
		return fmt.Errorf("tag is required")
	}

	req, err := newRequest(ctx, http.MethodDelete, c.buildURL("/repository/%s/%s/tag/%s", namespace, repository, tag))
	if err != nil {
		return fmt.Errorf("failed to create delete tag request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// GetTagHistory retrieves the history of changes for a specific tag
func (c *Client) GetTagHistory(ctx context.Context, namespace, repository, tag string) (*TagHistory, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/tag/%s/history", namespace, repository, tag))
	if err != nil {
		return nil, fmt.Errorf("failed to create get tag history request: %w", err)
	}

	var history TagHistory
	if err := c.get(req, &history); err != nil {
		return nil, fmt.Errorf("failed to get tag history: %w", err)
	}

	return &history, nil
}

// RevertTag reverts a tag to a previous state by manifest digest
func (c *Client) RevertTag(ctx context.Context, namespace, repository, tag, manifestDigest string) (*Tag, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag is required")
	}
	if manifestDigest == "" {
		return nil, fmt.Errorf("manifestDigest is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/repository/%s/%s/tag/%s/revert", namespace, repository, tag), RevertTagRequest{
		ManifestDigest: manifestDigest,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create revert tag request: %w", err)
	}

	var tagInfo Tag
	if err := c.post(req, &tagInfo); err != nil {
		return nil, fmt.Errorf("failed to revert tag: %w", err)
	}

	return &tagInfo, nil
}

// ChangeTag creates or moves a tag to point at a specific manifest digest
func (c *Client) ChangeTag(ctx context.Context, namespace, repository, tag, manifestDigest string) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}
	if tag == "" {
		return fmt.Errorf("tag is required")
	}
	if manifestDigest == "" {
		return fmt.Errorf("manifestDigest is required")
	}

	body := struct {
		ManifestDigest string `json:"manifest_digest"`
	}{
		ManifestDigest: manifestDigest,
	}
	req, err := newRequestWithBody(ctx, http.MethodPut, c.buildURL("/repository/%s/%s/tag/%s", namespace, repository, tag), body)
	if err != nil {
		return fmt.Errorf("failed to create change tag request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to change tag: %w", err)
	}

	return nil
}

// RestoreTag restores a tag to a previous image
func (c *Client) RestoreTag(ctx context.Context, namespace, repository, tag, manifestDigest string) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}
	if tag == "" {
		return fmt.Errorf("tag is required")
	}
	if manifestDigest == "" {
		return fmt.Errorf("manifestDigest is required")
	}

	body := struct {
		ManifestDigest string `json:"manifest_digest"`
	}{
		ManifestDigest: manifestDigest,
	}
	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/repository/%s/%s/tag/%s/restore", namespace, repository, tag), body)
	if err != nil {
		return fmt.Errorf("failed to create restore tag request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to restore tag: %w", err)
	}

	return nil
}
