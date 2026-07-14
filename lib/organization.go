/*
Package lib provides Quay.io API client functionality.

This file covers core ORGANIZATION MANAGEMENT endpoints:

Organization CRUD:
  - POST   /api/v1/organization/                                    - CreateOrganization()
  - GET    /api/v1/organization/{orgname}                           - GetOrganization()
  - PUT    /api/v1/organization/{orgname}                           - UpdateOrganization()
  - DELETE /api/v1/organization/{orgname}                           - DeleteOrganization()

Organization Members:
  - GET    /api/v1/organization/{orgname}/members                   - GetOrganizationMembers()
  - GET    /api/v1/organization/{orgname}/members/{membername}      - GetOrganizationMember()
  - PUT    /api/v1/organization/{orgname}/members/{membername}      - AddOrganizationMember()
  - DELETE /api/v1/organization/{orgname}/members/{membername}      - RemoveOrganizationMember()

Organization Repositories:
  - GET    /api/v1/organization/{orgname}/repositories              - GetOrganizationRepositories()

Organization Collaborators:
  - GET    /api/v1/organization/{orgname}/collaborators             - GetOrganizationCollaborators()
*/
package lib

import (
	"fmt"
)

// Organization Management

// CreateOrganization creates a new organization
func (c *Client) CreateOrganization(name, email string) (*Organization, error) {
	req, err := newRequestWithBody("POST", c.BaseURL+"/organization/", CreateOrganizationRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create organization request: %w", err)
	}

	var org Organization
	if err := c.post(req, &org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return &org, nil
}

// GetOrganization retrieves organization details
func (c *Client) GetOrganization(orgname string) (*Organization, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization request: %w", err)
	}

	var org Organization
	if err := c.get(req, &org); err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return &org, nil
}

// DeleteOrganization deletes an organization
func (c *Client) DeleteOrganization(orgname string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s", orgname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete organization request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}

// UpdateOrganization updates organization settings
func (c *Client) UpdateOrganization(orgname, email string) (*Organization, error) {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s", orgname), map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create update organization request: %w", err)
	}

	var org Organization
	if err := c.put(req, &org); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return &org, nil
}

// Organization Members Management

// GetOrganizationMembers retrieves organization members
func (c *Client) GetOrganizationMembers(orgname string) (*OrganizationMembers, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/members", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization members request: %w", err)
	}

	var members OrganizationMembers
	if err := c.get(req, &members); err != nil {
		return nil, fmt.Errorf("failed to get organization members: %w", err)
	}

	return &members, nil
}

// GetOrganizationMember gets information about a specific organization member
func (c *Client) GetOrganizationMember(orgname, membername string) (*OrganizationMember, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/members/%s", orgname, membername), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization member request: %w", err)
	}

	var member OrganizationMember
	if err := c.get(req, &member); err != nil {
		return nil, fmt.Errorf("failed to get organization member: %w", err)
	}

	return &member, nil
}

// AddOrganizationMember adds a member to an organization
func (c *Client) AddOrganizationMember(orgname, membername string) error {
	req, err := newRequest("PUT", c.buildURL("/organization/%s/members/%s", orgname, membername), nil)
	if err != nil {
		return fmt.Errorf("failed to create add organization member request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to add organization member: %w", err)
	}

	return nil
}

// RemoveOrganizationMember removes a member from an organization
func (c *Client) RemoveOrganizationMember(orgname, membername string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/members/%s", orgname, membername), nil)
	if err != nil {
		return fmt.Errorf("failed to create remove organization member request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to remove organization member: %w", err)
	}

	return nil
}

// Organization Repositories

// GetOrganizationRepositories retrieves repositories for an organization
func (c *Client) GetOrganizationRepositories(orgname string) (*OrganizationRepositories, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/repositories", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization repositories request: %w", err)
	}

	var repos OrganizationRepositories
	if err := c.get(req, &repos); err != nil {
		return nil, fmt.Errorf("failed to get organization repositories: %w", err)
	}

	return &repos, nil
}

// Organization Collaborators

// GetOrganizationCollaborators gets the list of collaborators for an organization
func (c *Client) GetOrganizationCollaborators(orgname string) (*Collaborators, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/collaborators", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization collaborators request: %w", err)
	}

	var collaborators Collaborators
	if err := c.get(req, &collaborators); err != nil {
		return nil, fmt.Errorf("failed to get organization collaborators: %w", err)
	}

	return &collaborators, nil
}
