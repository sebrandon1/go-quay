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

// User Permission Endpoints (per Quay API spec)

// ListUserPermissions lists all user permissions for a repository
func (c *Client) ListUserPermissions(namespace, repository string) (*RepositoryPermissions, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/permissions/user/", QuayURL, namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list user permissions request: %w", err)
	}

	var permissions RepositoryPermissions
	if err := c.get(req, &permissions); err != nil {
		return nil, fmt.Errorf("failed to list user permissions: %w", err)
	}

	return &permissions, nil
}

// GetUserPermission gets permission for a specific user on a repository
func (c *Client) GetUserPermission(namespace, repository, username string) (*RepositoryPermission, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/permissions/user/%s", QuayURL, namespace, repository, username), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user permission request: %w", err)
	}

	var permission RepositoryPermission
	if err := c.get(req, &permission); err != nil {
		return nil, fmt.Errorf("failed to get user permission: %w", err)
	}

	return &permission, nil
}

// SetUserPermission sets permission for a user on a repository
// role should be one of: "read", "write", "admin"
func (c *Client) SetUserPermission(namespace, repository, username, role string) error {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/repository/%s/%s/permissions/user/%s", QuayURL, namespace, repository, username), SetRepositoryPermissionRequest{
		Role: role,
	})
	if err != nil {
		return fmt.Errorf("failed to create set user permission request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to set user permission: %w", err)
	}

	return nil
}

// DeleteUserPermission removes permission for a user from a repository
func (c *Client) DeleteUserPermission(namespace, repository, username string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/permissions/user/%s", QuayURL, namespace, repository, username), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete user permission request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete user permission: %w", err)
	}

	return nil
}

// GetUserTransitivePermission gets the transitive permission for a user on a repository
func (c *Client) GetUserTransitivePermission(namespace, repository, username string) (*RepositoryPermission, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/permissions/user/%s/transitive", QuayURL, namespace, repository, username), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user transitive permission request: %w", err)
	}

	var permission RepositoryPermission
	if err := c.get(req, &permission); err != nil {
		return nil, fmt.Errorf("failed to get user transitive permission: %w", err)
	}

	return &permission, nil
}

// Team Permission Endpoints (per Quay API spec)

// ListTeamPermissions lists all team permissions for a repository
func (c *Client) ListTeamPermissions(namespace, repository string) (*RepositoryPermissions, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/permissions/team/", QuayURL, namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list team permissions request: %w", err)
	}

	var permissions RepositoryPermissions
	if err := c.get(req, &permissions); err != nil {
		return nil, fmt.Errorf("failed to list team permissions: %w", err)
	}

	return &permissions, nil
}

// GetTeamPermission gets permission for a specific team on a repository
func (c *Client) GetTeamPermission(namespace, repository, teamname string) (*RepositoryPermission, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/permissions/team/%s", QuayURL, namespace, repository, teamname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get team permission request: %w", err)
	}

	var permission RepositoryPermission
	if err := c.get(req, &permission); err != nil {
		return nil, fmt.Errorf("failed to get team permission: %w", err)
	}

	return &permission, nil
}

// SetTeamPermission sets permission for a team on a repository
// role should be one of: "read", "write", "admin"
func (c *Client) SetTeamPermission(namespace, repository, teamname, role string) error {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/repository/%s/%s/permissions/team/%s", QuayURL, namespace, repository, teamname), SetRepositoryPermissionRequest{
		Role: role,
	})
	if err != nil {
		return fmt.Errorf("failed to create set team permission request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to set team permission: %w", err)
	}

	return nil
}

// DeleteTeamPermission removes permission for a team from a repository
func (c *Client) DeleteTeamPermission(namespace, repository, teamname string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/permissions/team/%s", QuayURL, namespace, repository, teamname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete team permission request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete team permission: %w", err)
	}

	return nil
}
