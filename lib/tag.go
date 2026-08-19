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
	"fmt"
	"net/http"
)

// GetTag retrieves detailed information about a specific tag
func (c *Client) GetTag(namespace, repository, tag string) (*Tag, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag is required")
	}

	req, err := newRequest(http.MethodGet, c.buildURL("/repository/%s/%s/tag/%s", namespace, repository, tag), nil)
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
func (c *Client) UpdateTag(namespace, repository, tag, expiration string) (*Tag, error) {
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

	req, err := newRequestWithBody(http.MethodPut, c.buildURL("/repository/%s/%s/tag/%s", namespace, repository, tag), updateReq)
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
func (c *Client) DeleteTag(namespace, repository, tag string) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}
	if tag == "" {
		return fmt.Errorf("tag is required")
	}

	req, err := newRequest(http.MethodDelete, c.buildURL("/repository/%s/%s/tag/%s", namespace, repository, tag), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete tag request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// GetTagHistory retrieves the history of changes for a specific tag
func (c *Client) GetTagHistory(namespace, repository, tag string) (*TagHistory, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag is required")
	}

	req, err := newRequest(http.MethodGet, c.buildURL("/repository/%s/%s/tag/%s/history", namespace, repository, tag), nil)
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
func (c *Client) RevertTag(namespace, repository, tag, manifestDigest string) (*Tag, error) {
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

	req, err := newRequestWithBody(http.MethodPost, c.buildURL("/repository/%s/%s/tag/%s/revert", namespace, repository, tag), RevertTagRequest{
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
func (c *Client) ChangeTag(namespace, repository, tag, manifestDigest string) error {
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
	req, err := newRequestWithBody(http.MethodPut, c.buildURL("/repository/%s/%s/tag/%s", namespace, repository, tag), body)
	if err != nil {
		return fmt.Errorf("failed to create change tag request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to change tag: %w", err)
	}

	return nil
}

// RestoreTag restores a tag to a previous image
func (c *Client) RestoreTag(namespace, repository, tag, manifestDigest string) error {
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
	req, err := newRequestWithBody(http.MethodPost, c.buildURL("/repository/%s/%s/tag/%s/restore", namespace, repository, tag), body)
	if err != nil {
		return fmt.Errorf("failed to create restore tag request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to restore tag: %w", err)
	}

	return nil
}
