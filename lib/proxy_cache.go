/*
Package lib provides Quay.io API client functionality.

This file covers PROXY CACHE CONFIGURATION endpoints:

Proxy Cache Configuration:
  - GET    /api/v1/organization/{orgname}/proxycache                - GetProxyCacheConfig()
  - POST   /api/v1/organization/{orgname}/proxycache                - CreateProxyCacheConfig()
  - DELETE /api/v1/organization/{orgname}/proxycache                - DeleteProxyCacheConfig()
*/
package lib

import (
	"fmt"
	"net/http"
)

// GetProxyCacheConfig retrieves proxy cache configuration for an organization
func (c *Client) GetProxyCacheConfig(orgname string) (*ProxyCacheConfig, error) {
	if orgname == "" {
		return nil, fmt.Errorf("orgname is required")
	}

	req, err := newRequest(http.MethodGet, c.buildURL("/organization/%s/proxycache", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get proxy cache config request: %w", err)
	}

	var config ProxyCacheConfig
	if err := c.get(req, &config); err != nil {
		return nil, fmt.Errorf("failed to get proxy cache config: %w", err)
	}

	return &config, nil
}

// CreateProxyCacheConfig creates proxy cache configuration for an organization
func (c *Client) CreateProxyCacheConfig(orgname, upstreamRegistry string, insecure bool, expiration int) (*ProxyCacheConfig, error) {
	if orgname == "" {
		return nil, fmt.Errorf("orgname is required")
	}
	if upstreamRegistry == "" {
		return nil, fmt.Errorf("upstreamRegistry is required")
	}

	req, err := newRequestWithBody(http.MethodPost, c.buildURL("/organization/%s/proxycache", orgname), map[string]interface{}{
		"upstream_registry": upstreamRegistry,
		"insecure":          insecure,
		"expiration":        expiration,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy cache config request: %w", err)
	}

	var config ProxyCacheConfig
	if err := c.post(req, &config); err != nil {
		return nil, fmt.Errorf("failed to create proxy cache config: %w", err)
	}

	return &config, nil
}

// DeleteProxyCacheConfig deletes proxy cache configuration for an organization
func (c *Client) DeleteProxyCacheConfig(orgname string) error {
	if orgname == "" {
		return fmt.Errorf("orgname is required")
	}

	req, err := newRequest(http.MethodDelete, c.buildURL("/organization/%s/proxycache", orgname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete proxy cache config request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete proxy cache config: %w", err)
	}

	return nil
}
