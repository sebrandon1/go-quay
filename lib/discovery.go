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
	"net/http"
)

// GetDiscovery retrieves API discovery information
func (c *Client) GetDiscovery() (*Discovery, error) {
	req, err := newRequest(http.MethodGet, c.buildURL("/discovery"), nil)
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
	req, err := newRequest(http.MethodGet, c.buildURL("/registry/capabilities"), nil)
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
	if clientID == "" {
		return nil, fmt.Errorf("clientID is required")
	}

	req, err := newRequest(http.MethodGet, c.buildURL("/app/%s", clientID), nil)
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
	if prefix == "" {
		return nil, fmt.Errorf("prefix is required")
	}

	req, err := newRequest(http.MethodGet, c.buildURL("/entities/%s", prefix), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get entities request: %w", err)
	}

	params := map[string]string{}
	if includeOrgs {
		params["includeOrgs"] = queryValueTrue
	}
	if includeTeams {
		params["includeTeams"] = queryValueTrue
	}
	addQueryParams(req, params)

	var entities Entities
	if err := c.get(req, &entities); err != nil {
		return nil, fmt.Errorf("failed to get entities: %w", err)
	}

	return &entities, nil
}
