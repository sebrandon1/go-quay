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
	"context"
	"fmt"
	"net/http"
)

// SearchRepositories searches for repositories matching the query
func (c *Client) SearchRepositories(ctx context.Context, query string, page int) (*SearchRepositoryResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/find/repositories"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search repositories request: %w", err)
	}

	params := map[string]string{"query": query}
	if page > 0 {
		params["page"] = fmt.Sprintf("%d", page)
	}
	addQueryParams(req, params)

	var result SearchRepositoryResult
	if err := c.get(req, &result); err != nil {
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	return &result, nil
}

// SearchAll searches for all entity types (repositories, users, organizations, teams, robots)
func (c *Client) SearchAll(ctx context.Context, query string) (*SearchAllResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/find/all"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search all request: %w", err)
	}

	addQueryParams(req, map[string]string{"query": query})

	var result SearchAllResult
	if err := c.get(req, &result); err != nil {
		return nil, fmt.Errorf("failed to search all: %w", err)
	}

	return &result, nil
}
