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
	"fmt"
)

// GetBuilds retrieves a list of builds for a repository
func (c *Client) GetBuilds(namespace, repository string, limit int) (*Builds, error) {
	url := fmt.Sprintf("%s/repository/%s/%s/build/", QuayURL, namespace, repository)
	if limit > 0 {
		url = fmt.Sprintf("%s?limit=%d", url, limit)
	}

	req, err := newRequest("GET", url, nil)
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
func (c *Client) GetBuild(namespace, repository, buildUUID string) (*Build, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/build/%s", QuayURL, namespace, repository, buildUUID), nil)
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
func (c *Client) GetBuildLogs(namespace, repository, buildUUID string) (*BuildLogs, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/build/%s/logs", QuayURL, namespace, repository, buildUUID), nil)
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
func (c *Client) RequestBuild(namespace, repository string, buildRequest *RequestBuildRequest) (*Build, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/repository/%s/%s/build/", QuayURL, namespace, repository), buildRequest)
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
func (c *Client) CancelBuild(namespace, repository, buildUUID string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/build/%s", QuayURL, namespace, repository, buildUUID), nil)
	if err != nil {
		return fmt.Errorf("failed to create cancel build request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to cancel build: %w", err)
	}

	return nil
}
