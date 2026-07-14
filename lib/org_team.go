/*
Package lib provides Quay.io API client functionality.

This file covers ORGANIZATION TEAM MANAGEMENT endpoints:

Team Management:
  - GET    /api/v1/organization/{orgname}/teams                     - GetTeams()
  - PUT    /api/v1/organization/{orgname}/team/{teamname}           - CreateTeam()
  - GET    /api/v1/organization/{orgname}/team/{teamname}           - GetTeam()
  - PUT    /api/v1/organization/{orgname}/team/{teamname}           - UpdateTeam()
  - DELETE /api/v1/organization/{orgname}/team/{teamname}           - DeleteTeam()

Team Members:
  - GET    /api/v1/organization/{orgname}/team/{teamname}/members   - GetTeamMembers()
  - PUT    /api/v1/organization/{orgname}/team/{teamname}/members/{membername} - AddTeamMember()
  - DELETE /api/v1/organization/{orgname}/team/{teamname}/members/{membername} - RemoveTeamMember()

Team Permissions:
  - GET    /api/v1/organization/{orgname}/team/{teamname}/permissions - GetTeamPermissions()
  - PUT    /api/v1/organization/{orgname}/team/{teamname}/permissions/{repository} - SetTeamRepositoryPermission()
  - DELETE /api/v1/organization/{orgname}/team/{teamname}/permissions/{repository} - RemoveTeamRepositoryPermission()

Team Invitations:
  - PUT    /api/v1/organization/{orgname}/team/{teamname}/invite/{email} - InviteTeamMember()
  - DELETE /api/v1/organization/{orgname}/team/{teamname}/invite/{email} - DeleteTeamInvite()
*/
package lib

import (
	"fmt"
)

const (
	fieldRole = "role"
)

// Team Management

// GetTeams retrieves all teams for an organization
func (c *Client) GetTeams(orgname string) ([]Team, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/teams", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get teams request: %w", err)
	}

	var response struct {
		Teams []Team `json:"teams"`
	}
	if err := c.get(req, &response); err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}

	return response.Teams, nil
}

// CreateTeam creates a new team in an organization
func (c *Client) CreateTeam(orgname, teamname, description, role string) (*Team, error) {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s/team/%s", orgname, teamname), CreateTeamRequest{
		Name:        teamname,
		Description: description,
		Role:        role,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create team request: %w", err)
	}

	var team Team
	if err := c.put(req, &team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return &team, nil
}

// GetTeam retrieves details for a specific team
func (c *Client) GetTeam(orgname, teamname string) (*Team, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/team/%s", orgname, teamname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get team request: %w", err)
	}

	var team Team
	if err := c.get(req, &team); err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &team, nil
}

// DeleteTeam deletes a team from an organization
func (c *Client) DeleteTeam(orgname, teamname string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/team/%s", orgname, teamname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete team request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}

// UpdateTeam updates team settings
func (c *Client) UpdateTeam(orgname, teamname, description, role string) (*Team, error) {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s/team/%s", orgname, teamname), map[string]interface{}{
		"description": description,
		fieldRole:     role,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create update team request: %w", err)
	}

	var team Team
	if err := c.put(req, &team); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return &team, nil
}

// Team Members Management

// GetTeamMembers retrieves members of a team
func (c *Client) GetTeamMembers(orgname, teamname string) (*TeamMembers, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/team/%s/members", orgname, teamname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get team members request: %w", err)
	}

	var members TeamMembers
	if err := c.get(req, &members); err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	return &members, nil
}

// AddTeamMember adds a member to a team
func (c *Client) AddTeamMember(orgname, teamname, membername string) error {
	req, err := newRequest("PUT", c.buildURL("/organization/%s/team/%s/members/%s", orgname, teamname, membername), nil)
	if err != nil {
		return fmt.Errorf("failed to create add team member request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	return nil
}

// RemoveTeamMember removes a member from a team
func (c *Client) RemoveTeamMember(orgname, teamname, membername string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/team/%s/members/%s", orgname, teamname, membername), nil)
	if err != nil {
		return fmt.Errorf("failed to create remove team member request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}

	return nil
}

// Team Permissions Management

// GetTeamPermissions retrieves repository permissions for a team
func (c *Client) GetTeamPermissions(orgname, teamname string) (*TeamPermissions, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/team/%s/permissions", orgname, teamname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get team permissions request: %w", err)
	}

	var perms TeamPermissions
	if err := c.get(req, &perms); err != nil {
		return nil, fmt.Errorf("failed to get team permissions: %w", err)
	}

	return &perms, nil
}

// SetTeamRepositoryPermission sets repository permission for a team
func (c *Client) SetTeamRepositoryPermission(orgname, teamname, repository, role string) error {
	req, err := newRequestWithBody("PUT", c.buildURL("/organization/%s/team/%s/permissions/%s", orgname, teamname, repository), map[string]interface{}{
		fieldRole: role,
	})
	if err != nil {
		return fmt.Errorf("failed to create set team repository permission request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to set team repository permission: %w", err)
	}

	return nil
}

// RemoveTeamRepositoryPermission removes repository permission for a team
func (c *Client) RemoveTeamRepositoryPermission(orgname, teamname, repository string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/team/%s/permissions/%s", orgname, teamname, repository), nil)
	if err != nil {
		return fmt.Errorf("failed to create remove team repository permission request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to remove team repository permission: %w", err)
	}

	return nil
}

// Team Invitations

// InviteTeamMember invites a member to a team via email
func (c *Client) InviteTeamMember(orgname, teamname, email string) error {
	req, err := newRequest("PUT", c.buildURL("/organization/%s/team/%s/invite/%s", orgname, teamname, email), nil)
	if err != nil {
		return fmt.Errorf("failed to create invite team member request: %w", err)
	}

	if err := c.put(req, nil); err != nil {
		return fmt.Errorf("failed to invite team member: %w", err)
	}

	return nil
}

// DeleteTeamInvite deletes a pending team invitation
func (c *Client) DeleteTeamInvite(orgname, teamname, email string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/team/%s/invite/%s", orgname, teamname, email), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete team invite request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete team invite: %w", err)
	}

	return nil
}
