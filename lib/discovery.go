/*
Package lib provides Quay.io API client functionality.

This file covers DISCOVERY and CAPABILITIES operations:

API Discovery:
  - GET /api/v1/discovery              - GetDiscovery()

Registry Capabilities:
  - GET /api/v1/registry/capabilities  - GetRegistryCapabilities()
*/
package lib

import (
	"fmt"
)

// GetDiscovery retrieves API discovery information
func (c *Client) GetDiscovery() (*Discovery, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/discovery", c.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery request: %w", err)
	}

	var discovery Discovery
	if err := c.get(req, &discovery); err != nil {
		return nil, fmt.Errorf("failed to get discovery: %w", err)
	}

	return &discovery, nil
}

// GetRegistryCapabilities retrieves the registry capabilities including sparse manifest support and mirror architectures
func (c *Client) GetRegistryCapabilities() (*RegistryCapabilities, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/registry/capabilities", c.BaseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry capabilities request: %w", err)
	}

	var capabilities RegistryCapabilities
	if err := c.get(req, &capabilities); err != nil {
		return nil, fmt.Errorf("failed to get registry capabilities: %w", err)
	}

	return &capabilities, nil
}

// GetAppInfo retrieves public information about an OAuth application by client ID
func (c *Client) GetAppInfo(clientID string) (*Application, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/app/%s", c.BaseURL, clientID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get app info request: %w", err)
	}

	var app Application
	if err := c.get(req, &app); err != nil {
		return nil, fmt.Errorf("failed to get app info: %w", err)
	}

	return &app, nil
}

// GetEntities searches for entities (users, robots, teams) matching a prefix
func (c *Client) GetEntities(prefix string, includeOrgs, includeTeams bool) (*Entities, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/entities/%s", c.BaseURL, prefix), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get entities request: %w", err)
	}

	q := req.URL.Query()
	if includeOrgs {
		q.Add("includeOrgs", "true")
	}
	if includeTeams {
		q.Add("includeTeams", "true")
	}
	req.URL.RawQuery = q.Encode()

	var entities Entities
	if err := c.get(req, &entities); err != nil {
		return nil, fmt.Errorf("failed to get entities: %w", err)
	}

	return &entities, nil
}
