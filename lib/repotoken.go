/*
Package lib provides Quay.io API client functionality.

This file covers REPOTOKEN operations:

Repository Token Management (DEPRECATED - use robot accounts instead):
  - GET    /api/v1/repository/{namespace}/{repository}/tokens        - GetRepoTokens()
  - POST   /api/v1/repository/{namespace}/{repository}/tokens        - CreateRepoToken()
  - GET    /api/v1/repository/{namespace}/{repository}/tokens/{code} - GetRepoToken()
  - PUT    /api/v1/repository/{namespace}/{repository}/tokens/{code} - UpdateRepoToken()
  - DELETE /api/v1/repository/{namespace}/{repository}/tokens/{code} - DeleteRepoToken()

WARNING: Repository tokens are deprecated. Use robot accounts for authentication instead.
Robot accounts provide better security, auditing, and permission management.
*/
package lib

import (
	"context"
	"fmt"
	"net/http"
)

// GetRepoTokens retrieves all tokens for a repository.
//
// Deprecated: Use robot accounts instead.
func (c *Client) GetRepoTokens(ctx context.Context, namespace, repository string) (*RepoTokens, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/tokens", namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get repo tokens request: %w", err)
	}

	var tokens RepoTokens
	if err := c.get(req, &tokens); err != nil {
		return nil, fmt.Errorf("failed to get repo tokens: %w", err)
	}

	return &tokens, nil
}

// CreateRepoToken creates a new repository token.
//
// Deprecated: Use robot accounts instead.
func (c *Client) CreateRepoToken(ctx context.Context, namespace, repository string, createReq *CreateRepoTokenRequest) (*RepoToken, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/repository/%s/%s/tokens", namespace, repository), createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create repo token request: %w", err)
	}

	var token RepoToken
	if err := c.post(req, &token); err != nil {
		return nil, fmt.Errorf("failed to create repo token: %w", err)
	}

	return &token, nil
}

// GetRepoToken retrieves a specific repository token.
//
// Deprecated: Use robot accounts instead.
func (c *Client) GetRepoToken(ctx context.Context, namespace, repository, code string) (*RepoToken, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if code == "" {
		return nil, fmt.Errorf("code is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/tokens/%s", namespace, repository, code), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get repo token request: %w", err)
	}

	var token RepoToken
	if err := c.get(req, &token); err != nil {
		return nil, fmt.Errorf("failed to get repo token: %w", err)
	}

	return &token, nil
}

// UpdateRepoToken updates a repository token.
//
// Deprecated: Use robot accounts instead.
func (c *Client) UpdateRepoToken(ctx context.Context, namespace, repository, code string, updateReq *UpdateRepoTokenRequest) (*RepoToken, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if code == "" {
		return nil, fmt.Errorf("code is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPut, c.buildURL("/repository/%s/%s/tokens/%s", namespace, repository, code), updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create update repo token request: %w", err)
	}

	var token RepoToken
	if err := c.put(req, &token); err != nil {
		return nil, fmt.Errorf("failed to update repo token: %w", err)
	}

	return &token, nil
}

// DeleteRepoToken deletes a repository token.
//
// Deprecated: Use robot accounts instead.
func (c *Client) DeleteRepoToken(ctx context.Context, namespace, repository, code string) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}
	if code == "" {
		return fmt.Errorf("code is required")
	}

	req, err := newRequest(ctx, http.MethodDelete, c.buildURL("/repository/%s/%s/tokens/%s", namespace, repository, code), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete repo token request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete repo token: %w", err)
	}

	return nil
}
