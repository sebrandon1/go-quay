/*
Package lib provides Quay.io API client functionality.

This file covers ORGANIZATION QUOTA MANAGEMENT endpoints:

Quota Management:
  - GET    /api/v1/organization/{orgname}/quota                     - GetQuota()
  - POST   /api/v1/organization/{orgname}/quota                     - CreateQuota()
  - PUT    /api/v1/organization/{orgname}/quota                     - UpdateQuota()
  - DELETE /api/v1/organization/{orgname}/quota                     - DeleteQuota()
*/
package lib

import (
	"fmt"
)

// GetQuota retrieves quota information for an organization
func (c *Client) GetQuota(orgname string) (*Quota, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/quota", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get quota request: %w", err)
	}

	var quota Quota
	if err := c.get(req, &quota); err != nil {
		return nil, fmt.Errorf("failed to get quota: %w", err)
	}

	return &quota, nil
}

// CreateQuota creates a quota for an organization
func (c *Client) CreateQuota(orgname string, limitBytes int64) (*Quota, error) {
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/quota", orgname), CreateQuotaRequest{
		LimitBytes: limitBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create quota request: %w", err)
	}

	var quota Quota
	if err := c.post(req, &quota); err != nil {
		return nil, fmt.Errorf("failed to create quota: %w", err)
	}

	return &quota, nil
}

// UpdateQuota updates quota limits for an organization
func (c *Client) UpdateQuota(orgname string, limitBytes int64) (*Quota, error) {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s/quota", orgname), CreateQuotaRequest{
		LimitBytes: limitBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create update quota request: %w", err)
	}

	var quota Quota
	if err := c.put(req, &quota); err != nil {
		return nil, fmt.Errorf("failed to update quota: %w", err)
	}

	return &quota, nil
}

// DeleteQuota deletes quota for an organization
func (c *Client) DeleteQuota(orgname string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/quota", orgname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete quota request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete quota: %w", err)
	}

	return nil
}
