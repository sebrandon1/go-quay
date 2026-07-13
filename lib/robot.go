/*
Package lib provides Quay.io API client functionality.

This file covers USER ROBOT ACCOUNT operations:

Robot Account Management:
  - GET    /api/v1/user/robots                                      - GetUserRobotAccounts()
  - GET    /api/v1/user/robots/{robot_shortname}                    - GetUserRobotAccount()
  - PUT    /api/v1/user/robots/{robot_shortname}                    - CreateUserRobotAccount()
  - DELETE /api/v1/user/robots/{robot_shortname}                    - DeleteUserRobotAccount()
  - POST   /api/v1/user/robots/{robot_shortname}/regenerate         - RegenerateUserRobotToken()
  - GET    /api/v1/user/robots/{robot_shortname}/permissions        - GetUserRobotPermissions()
  - GET    /api/v1/user/robots/{robot_shortname}/federation         - GetUserRobotFederation()
  - POST   /api/v1/user/robots/{robot_shortname}/federation         - CreateUserRobotFederation()
  - DELETE /api/v1/user/robots/{robot_shortname}/federation         - DeleteUserRobotFederation()
*/
package lib

import (
	"fmt"
)

// GetUserRobotAccounts retrieves all robot accounts for the authenticated user
func (c *Client) GetUserRobotAccounts() (*RobotAccounts, error) {
	req, err := newRequest("GET", c.BaseURL+"/user/robots", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user robots request: %w", err)
	}

	var robots RobotAccounts
	if err := c.get(req, &robots); err != nil {
		return nil, fmt.Errorf("failed to get user robots: %w", err)
	}

	return &robots, nil
}

// GetUserRobotAccount retrieves a specific robot account for the authenticated user
func (c *Client) GetUserRobotAccount(robotShortname string) (*RobotAccount, error) {
	req, err := newRequest("GET", c.buildURL("/user/robots/%s", robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user robot request: %w", err)
	}

	var robot RobotAccount
	if err := c.get(req, &robot); err != nil {
		return nil, fmt.Errorf("failed to get user robot: %w", err)
	}

	return &robot, nil
}

// CreateUserRobotAccount creates a new robot account for the authenticated user
func (c *Client) CreateUserRobotAccount(robotShortname, description string, unstructured map[string]interface{}) (*RobotAccount, error) {
	createReq := CreateRobotRequest{
		Description:  description,
		Unstructured: unstructured,
	}

	req, err := newRequestWithBody("PUT", c.buildURL("/user/robots/%s", robotShortname), createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create user robot request: %w", err)
	}

	var robot RobotAccount
	if err := c.put(req, &robot); err != nil {
		return nil, fmt.Errorf("failed to create user robot: %w", err)
	}

	return &robot, nil
}

// DeleteUserRobotAccount deletes a robot account for the authenticated user
func (c *Client) DeleteUserRobotAccount(robotShortname string) error {
	req, err := newRequest("DELETE", c.buildURL("/user/robots/%s", robotShortname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete user robot request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete user robot: %w", err)
	}

	return nil
}

// RegenerateUserRobotToken regenerates the token for a user's robot account
func (c *Client) RegenerateUserRobotToken(robotShortname string) (*RobotAccount, error) {
	req, err := newRequest("POST", c.buildURL("/user/robots/%s/regenerate", robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create regenerate user robot token request: %w", err)
	}

	var robot RobotAccount
	if err := c.post(req, &robot); err != nil {
		return nil, fmt.Errorf("failed to regenerate user robot token: %w", err)
	}

	return &robot, nil
}

// GetUserRobotPermissions retrieves the repository permissions for a user's robot account
func (c *Client) GetUserRobotPermissions(robotShortname string) (*RobotPermissions, error) {
	req, err := newRequest("GET", c.buildURL("/user/robots/%s/permissions", robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user robot permissions request: %w", err)
	}

	var permissions RobotPermissions
	if err := c.get(req, &permissions); err != nil {
		return nil, fmt.Errorf("failed to get user robot permissions: %w", err)
	}

	return &permissions, nil
}

// GetUserRobotFederation retrieves the federation configuration for a user's robot account
func (c *Client) GetUserRobotFederation(robotShortname string) (*RobotFederation, error) {
	req, err := newRequest("GET", c.buildURL("/user/robots/%s/federation", robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user robot federation request: %w", err)
	}

	var federation RobotFederation
	if err := c.get(req, &federation); err != nil {
		return nil, fmt.Errorf("failed to get user robot federation: %w", err)
	}

	return &federation, nil
}

// CreateUserRobotFederation creates or updates the federation configuration for a user's robot account
func (c *Client) CreateUserRobotFederation(robotShortname string, configs []RobotFederationConfig) error {
	req, err := newRequestWithBody("POST", c.buildURL("/user/robots/%s/federation", robotShortname), configs)
	if err != nil {
		return fmt.Errorf("failed to create user robot federation request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to create user robot federation: %w", err)
	}

	return nil
}

// DeleteUserRobotFederation deletes the federation configuration for a user's robot account
func (c *Client) DeleteUserRobotFederation(robotShortname string) error {
	req, err := newRequest("DELETE", c.buildURL("/user/robots/%s/federation", robotShortname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete user robot federation request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete user robot federation: %w", err)
	}

	return nil
}
