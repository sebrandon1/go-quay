/*
Package lib provides Quay.io API client functionality.

This file covers ORGANIZATION MANAGEMENT endpoints:

Organization Management:
  - POST   /api/v1/organization/                                    - CreateOrganization()
  - GET    /api/v1/organization/{orgname}                           - GetOrganization()
  - PUT    /api/v1/organization/{orgname}                           - UpdateOrganization()
  - DELETE /api/v1/organization/{orgname}                           - DeleteOrganization()

Organization Members:
  - GET    /api/v1/organization/{orgname}/members                   - GetOrganizationMembers()
  - PUT    /api/v1/organization/{orgname}/members/{membername}      - AddOrganizationMember()
  - DELETE /api/v1/organization/{orgname}/members/{membername}      - RemoveOrganizationMember()

Organization Repositories:
  - GET    /api/v1/organization/{orgname}/repositories              - GetOrganizationRepositories()

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

Quota Management:
  - GET    /api/v1/organization/{orgname}/quota                     - GetQuota()
  - POST   /api/v1/organization/{orgname}/quota                     - CreateQuota()
  - PUT    /api/v1/organization/{orgname}/quota                     - UpdateQuota()
  - DELETE /api/v1/organization/{orgname}/quota                     - DeleteQuota()

Auto-Prune Policy Management:
  - GET    /api/v1/organization/{orgname}/autoprunepolicy          - GetAutoPrunePolicies()
  - POST   /api/v1/organization/{orgname}/autoprunepolicy          - CreateAutoPrunePolicy()
  - GET    /api/v1/organization/{orgname}/autoprunepolicy/{policy_uuid} - GetAutoPrunePolicy()
  - PUT    /api/v1/organization/{orgname}/autoprunepolicy/{policy_uuid} - UpdateAutoPrunePolicy()
  - DELETE /api/v1/organization/{orgname}/autoprunepolicy/{policy_uuid} - DeleteAutoPrunePolicy()

Applications Management:
  - GET    /api/v1/organization/{orgname}/applications              - GetApplications()
  - POST   /api/v1/organization/{orgname}/applications              - CreateApplication()
  - GET    /api/v1/organization/{orgname}/applications/{client_id}  - GetApplication()
  - PUT    /api/v1/organization/{orgname}/applications/{client_id}  - UpdateApplication()
  - DELETE /api/v1/organization/{orgname}/applications/{client_id}  - DeleteApplication()
  - POST   /api/v1/organization/{orgname}/applications/{client_id}/resetclientsecret - ResetApplicationClientSecret()

Proxy Cache Configuration:
  - GET    /api/v1/organization/{orgname}/proxycache                - GetProxyCacheConfig()
  - POST   /api/v1/organization/{orgname}/proxycache                - CreateProxyCacheConfig()
  - DELETE /api/v1/organization/{orgname}/proxycache                - DeleteProxyCacheConfig()
*/
package lib

import (
	"fmt"
)

const (
	fieldRole = "role"
	fieldName = "name"
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

// Proxy Cache Configuration

// GetProxyCacheConfig retrieves proxy cache configuration for an organization
func (c *Client) GetProxyCacheConfig(orgname string) (*ProxyCacheConfig, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/proxycache", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get proxy cache config request: %w", err)
	}

	var config ProxyCacheConfig
	if err := c.get(req, &config); err != nil {
		return nil, fmt.Errorf("failed to get proxy cache config: %w", err)
	}

	return &config, nil
}

// CreateProxyCacheConfig creates proxy cache configuration for an organization
func (c *Client) CreateProxyCacheConfig(orgname, upstreamRegistry string, insecure bool, expiration int) (*ProxyCacheConfig, error) {
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/proxycache", orgname), map[string]interface{}{
		"upstream_registry": upstreamRegistry,
		"insecure":          insecure,
		"expiration":        expiration,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy cache config request: %w", err)
	}

	var config ProxyCacheConfig
	if err := c.post(req, &config); err != nil {
		return nil, fmt.Errorf("failed to create proxy cache config: %w", err)
	}

	return &config, nil
}

// DeleteProxyCacheConfig deletes proxy cache configuration for an organization
func (c *Client) DeleteProxyCacheConfig(orgname string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/proxycache", orgname), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete proxy cache config request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete proxy cache config: %w", err)
	}

	return nil
}

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

// Quota Management

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

// Auto-Prune Policy Management

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

// Applications Management

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

// GetOrganizationMarketplace gets marketplace information for an organization
func (c *Client) GetOrganizationMarketplace(orgname string) (*MarketplaceInfo, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/marketplace", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization marketplace request: %w", err)
	}

	var marketplace MarketplaceInfo
	if err := c.get(req, &marketplace); err != nil {
		return nil, fmt.Errorf("failed to get organization marketplace: %w", err)
	}

	return &marketplace, nil
}

// CreateOrganizationMarketplaceSubscription creates a marketplace subscription
func (c *Client) CreateOrganizationMarketplaceSubscription(orgname string, subscription *MarketplaceSubscriptionRequest) error {
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/marketplace", orgname), subscription)
	if err != nil {
		return fmt.Errorf("failed to create marketplace subscription request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to create marketplace subscription: %w", err)
	}

	return nil
}

// BatchRemoveOrganizationMarketplaceSubscriptions removes multiple marketplace subscriptions
func (c *Client) BatchRemoveOrganizationMarketplaceSubscriptions(orgname string, subscriptionIDs []string) error {
	body := struct {
		SubscriptionIDs []string `json:"subscription_ids"`
	}{
		SubscriptionIDs: subscriptionIDs,
	}
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/marketplace/batchremove", orgname), body)
	if err != nil {
		return fmt.Errorf("failed to create batch remove marketplace subscriptions request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to batch remove marketplace subscriptions: %w", err)
	}

	return nil
}

// DeleteOrganizationMarketplaceSubscription removes a specific marketplace subscription
func (c *Client) DeleteOrganizationMarketplaceSubscription(orgname, subscriptionID string) error {
	req, err := newRequest("DELETE", c.buildURL("/organization/%s/marketplace/%s", orgname, subscriptionID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete marketplace subscription request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete marketplace subscription: %w", err)
	}

	return nil
}

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
