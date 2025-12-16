/*
Package lib provides Quay.io API client functionality.

This file covers SEARCH operations:

Search:
  - GET /api/v1/find/repositories - SearchRepositories()
  - GET /api/v1/find/all          - SearchAll()

Search operations provide discovery capabilities for finding repositories,
users, organizations, teams, and robots across Quay.io.
*/
package lib

import (
	"fmt"
	"net/url"
)

// SearchRepositories searches for repositories matching the query
func (c *Client) SearchRepositories(query string, page int) (*SearchRepositoryResult, error) {
	params := url.Values{}
	params.Add("query", query)
	if page > 0 {
		params.Add("page", fmt.Sprintf("%d", page))
	}

	req, err := newRequest("GET", fmt.Sprintf("%s/find/repositories?%s", QuayURL, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search repositories request: %w", err)
	}

	var result SearchRepositoryResult
	if err := c.get(req, &result); err != nil {
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	return &result, nil
}

// SearchAll searches for all entity types (repositories, users, organizations, teams, robots)
func (c *Client) SearchAll(query string) (*SearchAllResult, error) {
	params := url.Values{}
	params.Add("query", query)

	req, err := newRequest("GET", fmt.Sprintf("%s/find/all?%s", QuayURL, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search all request: %w", err)
	}

	var result SearchAllResult
	if err := c.get(req, &result); err != nil {
		return nil, fmt.Errorf("failed to search all: %w", err)
	}

	return &result, nil
}
