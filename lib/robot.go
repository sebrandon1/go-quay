/*
Package lib provides Quay.io API client functionality.

This file covers USER ROBOT ACCOUNT operations:

Robot Account Management:
  - GET    /api/v1/user/robots                                  - GetUserRobotAccounts()
  - GET    /api/v1/user/robots/{robot_shortname}                - GetUserRobotAccount()
  - PUT    /api/v1/user/robots/{robot_shortname}                - CreateUserRobotAccount()
  - DELETE /api/v1/user/robots/{robot_shortname}                - DeleteUserRobotAccount()
  - POST   /api/v1/user/robots/{robot_shortname}/regenerate     - RegenerateUserRobotToken()
  - GET    /api/v1/user/robots/{robot_shortname}/permissions    - GetUserRobotPermissions()

User robot accounts provide automated access credentials for CI/CD pipelines
and other automated workflows tied to a user account rather than an organization.
*/
package lib

import (
	"fmt"
)

// GetUserRobotAccounts retrieves all robot accounts for the authenticated user
func (c *Client) GetUserRobotAccounts() (*RobotAccounts, error) {
	req, err := newRequest("GET", QuayURL+"/user/robots", nil)
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
	req, err := newRequest("GET", fmt.Sprintf("%s/user/robots/%s", QuayURL, robotShortname), nil)
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

	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/user/robots/%s", QuayURL, robotShortname), createReq)
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
	req, err := newRequest("DELETE", fmt.Sprintf("%s/user/robots/%s", QuayURL, robotShortname), nil)
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
	req, err := newRequest("POST", fmt.Sprintf("%s/user/robots/%s/regenerate", QuayURL, robotShortname), nil)
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
	req, err := newRequest("GET", fmt.Sprintf("%s/user/robots/%s/permissions", QuayURL, robotShortname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user robot permissions request: %w", err)
	}

	var permissions RobotPermissions
	if err := c.get(req, &permissions); err != nil {
		return nil, fmt.Errorf("failed to get user robot permissions: %w", err)
	}

	return &permissions, nil
}
