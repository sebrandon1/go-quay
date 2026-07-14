/*
Package lib provides Quay.io API client functionality.

This file covers ORGANIZATION ROBOT ACCOUNT endpoints:

Robot Account Management:
  - GET    /api/v1/organization/{orgname}/robots                    - GetRobotAccounts()
  - PUT    /api/v1/organization/{orgname}/robots/{robot_shortname}  - CreateRobotAccount()
  - GET    /api/v1/organization/{orgname}/robots/{robot_shortname}  - GetRobotAccount()
  - DELETE /api/v1/organization/{orgname}/robots/{robot_shortname}  - DeleteRobotAccount()
  - POST   /api/v1/organization/{orgname}/robots/{robot_shortname}/regenerate - RegenerateRobotToken()

Robot Permissions:
  - GET    /api/v1/organization/{orgname}/robots/{robot_shortname}/permissions - GetRobotPermissions()
  - PUT    /api/v1/organization/{orgname}/robots/{robot_shortname}/permissions/{repository} - SetRobotRepositoryPermission()
  - DELETE /api/v1/organization/{orgname}/robots/{robot_shortname}/permissions/{repository} - RemoveRobotRepositoryPermission()

Robot Federation:
  - GET    /api/v1/organization/{orgname}/robots/{robot_shortname}/federation  - GetRobotFederation()
  - POST   /api/v1/organization/{orgname}/robots/{robot_shortname}/federation  - CreateRobotFederation()
  - DELETE /api/v1/organization/{orgname}/robots/{robot_shortname}/federation  - DeleteRobotFederation()
*/
package lib

import (
	"fmt"
)

// Robot Account Management

// GetRobotAccounts retrieves all robot accounts for an organization
func (c *Client) GetRobotAccounts(orgname string) (*RobotAccounts, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/robots", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get robot accounts request: %w", err)
	}

	var robots RobotAccounts
	if err := c.get(req, &robots); err != nil {
		return nil, fmt.Errorf("failed to get robot accounts: %w", err)
	}

	return &robots, nil
}

// CreateRobotAccount creates a new robot account in an organization
func (c *Client) CreateRobotAccount(orgname, robotShortname, description string, unstructured map[string]interface{}) (*RobotAccount, error) {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s/robots/%s", orgname, robotShortname), CreateRobotRequest{
		Description:  description,
		Unstructured: unstructured,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create robot account request: %w", err)
	}

	var robot RobotAccount
	if err := c.put(req, &robot); err != nil {
		return nil, fmt.Errorf("failed to create robot account: %w", err)
	}

	return &robot, nil
}

// GetRobotAccount retrieves details for a specific robot account
func (c *Client) GetRobotAccount(orgname, robotShortname string) (*RobotAccount, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/robots/%s", orgname, robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get robot account request: %w", err)
	}

	var robot RobotAccount
	if err := c.get(req, &robot); err != nil {
		return nil, fmt.Errorf("failed to get robot account: %w", err)
	}

	return &robot, nil
}

// DeleteRobotAccount deletes a robot account from an organization
func (c *Client) DeleteRobotAccount(orgname, robotShortname string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/robots/%s", orgname, robotShortname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete robot account request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete robot account: %w", err)
	}

	return nil
}

// RegenerateRobotToken regenerates the token for a robot account
func (c *Client) RegenerateRobotToken(orgname, robotShortname string) (*RobotAccount, error) {
	req, err := newRequest("POST", c.buildURL("/organization/%s/robots/%s/regenerate", orgname, robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create regenerate robot token request: %w", err)
	}

	var robot RobotAccount
	if err := c.post(req, &robot); err != nil {
		return nil, fmt.Errorf("failed to regenerate robot token: %w", err)
	}

	return &robot, nil
}

// Robot Permissions Management

// GetRobotPermissions retrieves repository permissions for a robot account
func (c *Client) GetRobotPermissions(orgname, robotShortname string) (*RobotPermissions, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/robots/%s/permissions", orgname, robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get robot permissions request: %w", err)
	}

	var perms RobotPermissions
	if err := c.get(req, &perms); err != nil {
		return nil, fmt.Errorf("failed to get robot permissions: %w", err)
	}

	return &perms, nil
}

// SetRobotRepositoryPermission sets repository permission for a robot account
func (c *Client) SetRobotRepositoryPermission(orgname, robotShortname, repository, role string) error {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s/robots/%s/permissions/%s", orgname, robotShortname, repository), map[string]interface{}{
		fieldRole: role,
	})
	if err != nil {
		return fmt.Errorf("failed to create set robot repository permission request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to set robot repository permission: %w", err)
	}

	return nil
}

// RemoveRobotRepositoryPermission removes repository permission for a robot account
func (c *Client) RemoveRobotRepositoryPermission(orgname, robotShortname, repository string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/robots/%s/permissions/%s", orgname, robotShortname, repository), nil)
	if err != nil {
		return fmt.Errorf("failed to create remove robot repository permission request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to remove robot repository permission: %w", err)
	}

	return nil
}

// Robot Federation

// GetRobotFederation retrieves the federation configuration for an organization's robot account
func (c *Client) GetRobotFederation(orgname, robotShortname string) (*RobotFederation, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/robots/%s/federation", orgname, robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get robot federation request: %w", err)
	}

	var federation RobotFederation
	if err := c.get(req, &federation); err != nil {
		return nil, fmt.Errorf("failed to get robot federation: %w", err)
	}

	return &federation, nil
}

// CreateRobotFederation creates or updates the federation configuration for an organization's robot account
func (c *Client) CreateRobotFederation(orgname, robotShortname string, configs []RobotFederationConfig) error {
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/robots/%s/federation", orgname, robotShortname), configs)
	if err != nil {
		return fmt.Errorf("failed to create robot federation request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to create robot federation: %w", err)
	}

	return nil
}

// DeleteRobotFederation deletes the federation configuration for an organization's robot account
func (c *Client) DeleteRobotFederation(orgname, robotShortname string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/robots/%s/federation", orgname, robotShortname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete robot federation request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete robot federation: %w", err)
	}

	return nil
}
