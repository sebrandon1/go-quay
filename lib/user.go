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

// StarUserRepository stars a repository using the /user/starred endpoint
func (c *Client) StarUserRepository(namespace, repository string) error {
	body := struct {
		Repository string `json:"repository"`
	}{
		Repository: namespace + "/" + repository,
	}
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/user/starred", QuayURL), body)
	if err != nil {
		return fmt.Errorf("failed to create star user repository request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to star repository: %w", err)
	}

	return nil
}

// UnstarUserRepository unstars a repository using the /user/starred endpoint
func (c *Client) UnstarUserRepository(namespace, repository string) error {
	// The spec uses {repository} as namespace/repo
	req, err := newRequest("DELETE", fmt.Sprintf("%s/user/starred/%s/%s", QuayURL, namespace, repository), nil)
	if err != nil {
		return fmt.Errorf("failed to create unstar user repository request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to unstar repository: %w", err)
	}

	return nil
}

// GetUserByUsername retrieves information about a specific user
func (c *Client) GetUserByUsername(username string) (*UserDetails, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/users/%s", QuayURL, username), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user by username request: %w", err)
	}

	var user UserDetails
	if err := c.get(req, &user); err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// GetUserMarketplace retrieves marketplace information for the current user
func (c *Client) GetUserMarketplace() (*MarketplaceInfo, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/user/marketplace", QuayURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user marketplace request: %w", err)
	}

	var marketplace MarketplaceInfo
	if err := c.get(req, &marketplace); err != nil {
		return nil, fmt.Errorf("failed to get user marketplace: %w", err)
	}

	return &marketplace, nil
}
