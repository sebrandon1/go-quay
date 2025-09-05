/*
Package lib provides Quay.io API client functionality.

This file covers USER ACCOUNT operations:

User Management:
  - GET    /api/v1/user                                                 - GetUser()
  - GET    /api/v1/user/starred                                         - GetStarredRepositories()

Repository Stars:
  - PUT    /api/v1/repository/{namespace}/{repository}/star             - StarRepository()
  - DELETE /api/v1/repository/{namespace}/{repository}/star             - UnstarRepository()

User operations provide access to account information, starred repositories,
and the ability to star/unstar repositories for easy discovery.
*/
package lib

import (
	"fmt"
)

// GetUser retrieves information about the current authenticated user
func (c *Client) GetUser() (*UserDetails, error) {
	req, err := newRequest("GET", QuayURL+"/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user request: %w", err)
	}

	var user UserDetails
	if err := c.get(req, &user); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetStarredRepositories retrieves repositories starred by the current user
func (c *Client) GetStarredRepositories() (*StarredRepositories, error) {
	req, err := newRequest("GET", QuayURL+"/user/starred", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get starred repositories request: %w", err)
	}

	var starred StarredRepositories
	if err := c.get(req, &starred); err != nil {
		return nil, fmt.Errorf("failed to get starred repositories: %w", err)
	}

	return &starred, nil
}

// StarRepository adds a repository to the user's starred list
func (c *Client) StarRepository(namespace, repository string) error {
	req, err := newRequest("PUT", fmt.Sprintf("%s/repository/%s/%s/star", QuayURL, namespace, repository), nil)
	if err != nil {
		return fmt.Errorf("failed to create star repository request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to star repository: %w", err)
	}

	return nil
}

// UnstarRepository removes a repository from the user's starred list
func (c *Client) UnstarRepository(namespace, repository string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/star", QuayURL, namespace, repository), nil)
	if err != nil {
		return fmt.Errorf("failed to create unstar repository request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to unstar repository: %w", err)
	}

	return nil
}
