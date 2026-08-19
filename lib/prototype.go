/*
Package lib provides Quay.io API client functionality.

This file covers PROTOTYPE (Default Permission) operations:

Prototype Management:
  - GET    /api/v1/organization/{orgname}/prototypes             - GetPrototypes()
  - POST   /api/v1/organization/{orgname}/prototypes             - CreatePrototype()
  - GET    /api/v1/organization/{orgname}/prototypes/{uuid}      - GetPrototype()
  - PUT    /api/v1/organization/{orgname}/prototypes/{uuid}      - UpdatePrototype()
  - DELETE /api/v1/organization/{orgname}/prototypes/{uuid}      - DeletePrototype()

Prototypes define default permissions that are automatically applied to new
repositories created within an organization. They allow setting up permission
templates for users, teams, or robot accounts.

Delegate kinds:
  - user: A specific user account
  - team: A team within the organization
  - robot: A robot account
*/
package lib

import (
	"context"
	"fmt"
	"net/http"
)

// GetPrototypes retrieves all permission prototypes for an organization
func (c *Client) GetPrototypes(ctx context.Context, orgname string) (*Prototypes, error) {
	if orgname == "" {
		return nil, fmt.Errorf("orgname is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/organization/%s/prototypes", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get prototypes request: %w", err)
	}

	var prototypes Prototypes
	if err := c.get(req, &prototypes); err != nil {
		return nil, fmt.Errorf("failed to get prototypes: %w", err)
	}

	return &prototypes, nil
}

// CreatePrototype creates a new permission prototype for an organization
func (c *Client) CreatePrototype(ctx context.Context, orgname string, createReq *CreatePrototypeRequest) (*Prototype, error) {
	if orgname == "" {
		return nil, fmt.Errorf("orgname is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/organization/%s/prototypes", orgname), createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create prototype request: %w", err)
	}

	var prototype Prototype
	if err := c.post(req, &prototype); err != nil {
		return nil, fmt.Errorf("failed to create prototype: %w", err)
	}

	return &prototype, nil
}

// GetPrototype retrieves a specific prototype by UUID
func (c *Client) GetPrototype(ctx context.Context, orgname, prototypeUUID string) (*Prototype, error) {
	if orgname == "" {
		return nil, fmt.Errorf("orgname is required")
	}
	if prototypeUUID == "" {
		return nil, fmt.Errorf("prototypeUUID is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/organization/%s/prototypes/%s", orgname, prototypeUUID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get prototype request: %w", err)
	}

	var prototype Prototype
	if err := c.get(req, &prototype); err != nil {
		return nil, fmt.Errorf("failed to get prototype: %w", err)
	}

	return &prototype, nil
}

// UpdatePrototype updates an existing prototype
func (c *Client) UpdatePrototype(ctx context.Context, orgname, prototypeUUID string, updateReq *UpdatePrototypeRequest) (*Prototype, error) {
	if orgname == "" {
		return nil, fmt.Errorf("orgname is required")
	}
	if prototypeUUID == "" {
		return nil, fmt.Errorf("prototypeUUID is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPut, c.buildURL("/organization/%s/prototypes/%s", orgname, prototypeUUID), updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create update prototype request: %w", err)
	}

	var prototype Prototype
	if err := c.put(req, &prototype); err != nil {
		return nil, fmt.Errorf("failed to update prototype: %w", err)
	}

	return &prototype, nil
}

// DeletePrototype deletes a prototype
func (c *Client) DeletePrototype(ctx context.Context, orgname, prototypeUUID string) error {
	if orgname == "" {
		return fmt.Errorf("orgname is required")
	}
	if prototypeUUID == "" {
		return fmt.Errorf("prototypeUUID is required")
	}

	req, err := newRequest(ctx, http.MethodDelete, c.buildURL("/organization/%s/prototypes/%s", orgname, prototypeUUID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete prototype request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete prototype: %w", err)
	}

	return nil
}
