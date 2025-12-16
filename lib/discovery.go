/*
Package lib provides Quay.io API client functionality.

This file covers DISCOVERY operations:

API Discovery:
  - GET /api/v1/discovery - GetDiscovery()

The Discovery API provides information about available API endpoints and versions.
*/
package lib

import (
	"fmt"
)

// GetDiscovery retrieves API discovery information
func (c *Client) GetDiscovery() (*Discovery, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/discovery", QuayURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery request: %w", err)
	}

	var discovery Discovery
	if err := c.get(req, &discovery); err != nil {
		return nil, fmt.Errorf("failed to get discovery: %w", err)
	}

	return &discovery, nil
}
