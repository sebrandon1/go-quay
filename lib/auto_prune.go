/*
Package lib provides Quay.io API client functionality.

This file covers ORGANIZATION AUTO-PRUNE POLICY endpoints:

Auto-Prune Policy Management:
  - GET    /api/v1/organization/{orgname}/autoprunepolicy          - GetAutoPrunePolicies()
  - POST   /api/v1/organization/{orgname}/autoprunepolicy          - CreateAutoPrunePolicy()
  - GET    /api/v1/organization/{orgname}/autoprunepolicy/{policy_uuid} - GetAutoPrunePolicy()
  - PUT    /api/v1/organization/{orgname}/autoprunepolicy/{policy_uuid} - UpdateAutoPrunePolicy()
  - DELETE /api/v1/organization/{orgname}/autoprunepolicy/{policy_uuid} - DeleteAutoPrunePolicy()
*/
package lib

import (
	"fmt"
)

// GetAutoPrunePolicies retrieves auto-prune policies for an organization
func (c *Client) GetAutoPrunePolicies(orgname string) (*AutoPrunePolicies, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/autoprunepolicy", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get auto-prune policies request: %w", err)
	}

	var policies AutoPrunePolicies
	if err := c.get(req, &policies); err != nil {
		return nil, fmt.Errorf("failed to get auto-prune policies: %w", err)
	}

	return &policies, nil
}

// CreateAutoPrunePolicy creates an auto-prune policy for an organization
func (c *Client) CreateAutoPrunePolicy(orgname, method string, value int, tagPattern string) (*AutoPrunePolicy, error) {
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/autoprunepolicy", orgname), CreateAutoPruneRequest{
		Method:     method,
		Value:      value,
		TagPattern: tagPattern,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auto-prune policy request: %w", err)
	}

	var policy AutoPrunePolicy
	if err := c.post(req, &policy); err != nil {
		return nil, fmt.Errorf("failed to create auto-prune policy: %w", err)
	}

	return &policy, nil
}

// GetAutoPrunePolicy retrieves a specific auto-prune policy
func (c *Client) GetAutoPrunePolicy(orgname, policyUUID string) (*AutoPrunePolicy, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/autoprunepolicy/%s", orgname, policyUUID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get auto-prune policy request: %w", err)
	}

	var policy AutoPrunePolicy
	if err := c.get(req, &policy); err != nil {
		return nil, fmt.Errorf("failed to get auto-prune policy: %w", err)
	}

	return &policy, nil
}

// UpdateAutoPrunePolicy updates an auto-prune policy
func (c *Client) UpdateAutoPrunePolicy(orgname, policyUUID, method string, value int, tagPattern string) (*AutoPrunePolicy, error) {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s/autoprunepolicy/%s", orgname, policyUUID), CreateAutoPruneRequest{
		Method:     method,
		Value:      value,
		TagPattern: tagPattern,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create update auto-prune policy request: %w", err)
	}

	var policy AutoPrunePolicy
	if err := c.put(req, &policy); err != nil {
		return nil, fmt.Errorf("failed to update auto-prune policy: %w", err)
	}

	return &policy, nil
}

// DeleteAutoPrunePolicy deletes an auto-prune policy
func (c *Client) DeleteAutoPrunePolicy(orgname, policyUUID string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/autoprunepolicy/%s", orgname, policyUUID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete auto-prune policy request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete auto-prune policy: %w", err)
	}

	return nil
}
