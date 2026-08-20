/*
Package lib provides Quay.io API client functionality.

This file covers BUILD operations:

Build Management:
  - GET    /api/v1/repository/{namespace}/{repository}/build/                    - GetBuilds()
  - GET    /api/v1/repository/{namespace}/{repository}/build/{build_uuid}        - GetBuild()
  - POST   /api/v1/repository/{namespace}/{repository}/build/                    - RequestBuild()
  - DELETE /api/v1/repository/{namespace}/{repository}/build/{build_uuid}        - CancelBuild()
  - GET    /api/v1/repository/{namespace}/{repository}/build/{build_uuid}/logs   - GetBuildLogs()

Builds allow automated image creation from Dockerfiles stored in git repositories
or uploaded archives.
*/
package lib

import (
	"context"
	"fmt"
	"net/http"
)

// GetBuilds retrieves a list of builds for a repository
func (c *Client) GetBuilds(ctx context.Context, namespace, repository string, limit int) (*Builds, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	url := c.buildURL("/repository/%s/%s/build/", namespace, repository)
	if limit > 0 {
		url = fmt.Sprintf("%s?limit=%d", url, limit)
	}

	req, err := newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get builds request: %w", err)
	}

	var builds Builds
	if err := c.get(req, &builds); err != nil {
		return nil, fmt.Errorf("failed to get builds: %w", err)
	}

	return &builds, nil
}

// GetBuild retrieves a specific build by UUID
func (c *Client) GetBuild(ctx context.Context, namespace, repository, buildUUID string) (*Build, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if buildUUID == "" {
		return nil, fmt.Errorf("buildUUID is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/build/%s", namespace, repository, buildUUID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get build request: %w", err)
	}

	var build Build
	if err := c.get(req, &build); err != nil {
		return nil, fmt.Errorf("failed to get build: %w", err)
	}

	return &build, nil
}

// GetBuildLogs retrieves the logs for a specific build
func (c *Client) GetBuildLogs(ctx context.Context, namespace, repository, buildUUID string) (*BuildLogs, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if buildUUID == "" {
		return nil, fmt.Errorf("buildUUID is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/build/%s/logs", namespace, repository, buildUUID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get build logs request: %w", err)
	}

	var logs BuildLogs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get build logs: %w", err)
	}

	return &logs, nil
}

// RequestBuild triggers a new build for a repository
func (c *Client) RequestBuild(ctx context.Context, namespace, repository string, buildRequest *RequestBuildRequest) (*Build, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/repository/%s/%s/build/", namespace, repository), buildRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create request build request: %w", err)
	}

	var build Build
	if err := c.post(req, &build); err != nil {
		return nil, fmt.Errorf("failed to request build: %w", err)
	}

	return &build, nil
}

// CancelBuild cancels an ongoing build
func (c *Client) CancelBuild(ctx context.Context, namespace, repository, buildUUID string) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}
	if buildUUID == "" {
		return fmt.Errorf("buildUUID is required")
	}

	req, err := newRequest(ctx, http.MethodDelete, c.buildURL("/repository/%s/%s/build/%s", namespace, repository, buildUUID), nil)
	if err != nil {
		return fmt.Errorf("failed to create cancel build request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to cancel build: %w", err)
	}

	return nil
}

// GetBuildStatus gets the status of a build
func (c *Client) GetBuildStatus(ctx context.Context, namespace, repository, buildUUID string) (*BuildStatus, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if buildUUID == "" {
		return nil, fmt.Errorf("buildUUID is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/build/%s/status", namespace, repository, buildUUID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get build status request: %w", err)
	}

	var status BuildStatus
	if err := c.get(req, &status); err != nil {
		return nil, fmt.Errorf("failed to get build status: %w", err)
	}

	return &status, nil
}
