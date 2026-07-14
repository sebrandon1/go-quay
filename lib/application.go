/*
Package lib provides Quay.io API client functionality.

This file covers ORGANIZATION APPLICATION endpoints:

Applications Management:
  - GET    /api/v1/organization/{orgname}/applications              - GetApplications()
  - POST   /api/v1/organization/{orgname}/applications              - CreateApplication()
  - GET    /api/v1/organization/{orgname}/applications/{client_id}  - GetApplication()
  - PUT    /api/v1/organization/{orgname}/applications/{client_id}  - UpdateApplication()
  - DELETE /api/v1/organization/{orgname}/applications/{client_id}  - DeleteApplication()
  - POST   /api/v1/organization/{orgname}/applications/{client_id}/resetclientsecret - ResetApplicationClientSecret()
*/
package lib

import (
	"fmt"
)

// GetApplications retrieves applications for an organization
func (c *Client) GetApplications(orgname string) (*Applications, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/applications", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get applications request: %w", err)
	}

	var apps Applications
	if err := c.get(req, &apps); err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	return &apps, nil
}

// CreateApplication creates a new application for an organization
func (c *Client) CreateApplication(orgname, name, description, applicationURI, redirectURI string) (*Application, error) {
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/applications", orgname), CreateApplicationRequest{
		Name:           name,
		Description:    description,
		ApplicationURI: applicationURI,
		RedirectURI:    redirectURI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create application request: %w", err)
	}

	var app Application
	if err := c.post(req, &app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return &app, nil
}

// GetApplication retrieves details for a specific application
func (c *Client) GetApplication(orgname, clientID string) (*Application, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/applications/%s", orgname, clientID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get application request: %w", err)
	}

	var app Application
	if err := c.get(req, &app); err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return &app, nil
}

// UpdateApplication updates an application
func (c *Client) UpdateApplication(orgname, clientID, name, description, applicationURI, redirectURI string) (*Application, error) {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s/applications/%s", orgname, clientID), CreateApplicationRequest{
		Name:           name,
		Description:    description,
		ApplicationURI: applicationURI,
		RedirectURI:    redirectURI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create update application request: %w", err)
	}

	var app Application
	if err := c.put(req, &app); err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return &app, nil
}

// DeleteApplication deletes an application
func (c *Client) DeleteApplication(orgname, clientID string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/applications/%s", orgname, clientID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete application request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	return nil
}

// ResetApplicationClientSecret resets the client secret for an application
func (c *Client) ResetApplicationClientSecret(orgname, clientID string) (*Application, error) {
	req, err := newRequest("POST", c.buildURL("/organization/%s/applications/%s/resetclientsecret", orgname, clientID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create reset application client secret request: %w", err)
	}

	var app Application
	if err := c.post(req, &app); err != nil {
		return nil, fmt.Errorf("failed to reset application client secret: %w", err)
	}

	return &app, nil
}
