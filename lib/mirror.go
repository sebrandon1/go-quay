/*
Package lib provides Quay.io API client functionality.

This file covers REPOSITORY MIRROR CONFIGURATION endpoints:

Repository Mirror Configuration:
  - GET  /api/v1/repository/{namespace}/{repository}/mirror   - GetMirrorConfig()
  - POST /api/v1/repository/{namespace}/{repository}/mirror   - CreateMirrorConfig()
  - PUT  /api/v1/repository/{namespace}/{repository}/mirror   - UpdateMirrorConfig()
*/
package lib

import (
	"context"
	"fmt"
	"net/http"
)

// GetMirrorConfig retrieves mirror configuration for a repository
func (c *Client) GetMirrorConfig(ctx context.Context, namespace, repository string) (*MirrorConfig, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/mirror", namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get mirror config request: %w", err)
	}

	var config MirrorConfig
	if err := c.get(req, &config); err != nil {
		return nil, fmt.Errorf("failed to get mirror config: %w", err)
	}

	return &config, nil
}

// CreateMirrorConfig creates mirror configuration for a repository
func (c *Client) CreateMirrorConfig(ctx context.Context, namespace, repository string, config *CreateMirrorConfigRequest) (*MirrorConfig, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/repository/%s/%s/mirror", namespace, repository), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create mirror config request: %w", err)
	}

	var result MirrorConfig
	if err := c.post(req, &result); err != nil {
		return nil, fmt.Errorf("failed to create mirror config: %w", err)
	}

	return &result, nil
}

// UpdateMirrorConfig updates mirror configuration for a repository
func (c *Client) UpdateMirrorConfig(ctx context.Context, namespace, repository string, config *UpdateMirrorConfigRequest) (*MirrorConfig, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPut, c.buildURL("/repository/%s/%s/mirror", namespace, repository), config)
	if err != nil {
		return nil, fmt.Errorf("failed to update mirror config request: %w", err)
	}

	var result MirrorConfig
	if err := c.put(req, &result); err != nil {
		return nil, fmt.Errorf("failed to update mirror config: %w", err)
	}

	return &result, nil
}
