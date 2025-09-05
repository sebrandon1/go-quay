/*
Package lib provides Quay.io API client functionality.

This file covers REPOSITORY PERMISSIONS endpoints:

Repository Permissions:
  - GET    /api/v1/repository/{namespace}/{repository}/permissions                - GetRepositoryPermissions()
  - PUT    /api/v1/repository/{namespace}/{repository}/permissions/{username}     - SetRepositoryPermission()
  - DELETE /api/v1/repository/{namespace}/{repository}/permissions/{username}     - RemoveRepositoryPermission()

Repository permissions control who can read, write, or administer repositories.
Supported roles: read, write, admin
*/
package lib

import (
	"fmt"
)

// GetRepositoryPermissions retrieves permissions for a repository
func (c *Client) GetRepositoryPermissions(namespace, repository string) (*RepositoryPermissions, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/permissions", QuayURL, namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get repository permissions request: %w", err)
	}

	var permissions RepositoryPermissions
	if err := c.get(req, &permissions); err != nil {
		return nil, fmt.Errorf("failed to get repository permissions: %w", err)
	}

	return &permissions, nil
}

// SetRepositoryPermission sets permission for a user/robot on a repository
// role should be one of: "read", "write", "admin"
func (c *Client) SetRepositoryPermission(namespace, repository, username, role string) error {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/repository/%s/%s/permissions/%s", QuayURL, namespace, repository, username), SetRepositoryPermissionRequest{
		Role: role,
	})
	if err != nil {
		return fmt.Errorf("failed to create set repository permission request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to set repository permission: %w", err)
	}

	return nil
}

// RemoveRepositoryPermission removes permission for a user/robot from a repository
func (c *Client) RemoveRepositoryPermission(namespace, repository, username string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/permissions/%s", QuayURL, namespace, repository, username), nil)
	if err != nil {
		return fmt.Errorf("failed to create remove repository permission request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to remove repository permission: %w", err)
	}

	return nil
}
