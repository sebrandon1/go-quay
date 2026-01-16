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

Default Permissions Management:
  - GET    /api/v1/organization/{orgname}/prototypes                - GetDefaultPermissions()
  - POST   /api/v1/organization/{orgname}/prototypes                - CreateDefaultPermission()
  - DELETE /api/v1/organization/{orgname}/prototypes/{prototypeid}  - DeleteDefaultPermission()

Proxy Cache Configuration:
  - GET    /api/v1/organization/{orgname}/proxycache                - GetProxyCacheConfig()
  - POST   /api/v1/organization/{orgname}/proxycache                - CreateProxyCacheConfig()
  - DELETE /api/v1/organization/{orgname}/proxycache                - DeleteProxyCacheConfig()
*/
package lib

import (
	"fmt"
)

// Organization Management

// CreateOrganization creates a new organization
func (c *Client) CreateOrganization(name, email string) (*Organization, error) {
	req, err := newRequestWithBody("POST", QuayURL+"/organization/", CreateOrganizationRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	var org Organization
	if err := c.post(req, &org); err != nil {
		return nil, err
	}

	return &org, nil
}

// GetOrganization retrieves organization details
func (c *Client) GetOrganization(orgname string) (*Organization, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var org Organization
	if err := c.get(req, &org); err != nil {
		return nil, err
	}

	return &org, nil
}

// DeleteOrganization deletes an organization
func (c *Client) DeleteOrganization(orgname string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s", QuayURL, orgname), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// UpdateOrganization updates organization settings
func (c *Client) UpdateOrganization(orgname, email string) (*Organization, error) {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s", QuayURL, orgname), map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return nil, err
	}

	var org Organization
	if err := c.put(req, &org); err != nil {
		return nil, err
	}

	return &org, nil
}

// Organization Members Management

// GetOrganizationMembers retrieves organization members
func (c *Client) GetOrganizationMembers(orgname string) (*OrganizationMembers, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/members", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var members OrganizationMembers
	if err := c.get(req, &members); err != nil {
		return nil, err
	}

	return &members, nil
}

// AddOrganizationMember adds a member to an organization
func (c *Client) AddOrganizationMember(orgname, membername string) error {
	req, err := newRequest("PUT", fmt.Sprintf("%s/organization/%s/members/%s", QuayURL, orgname, membername), nil)
	if err != nil {
		return err
	}

	return c.put(req, nil)
}

// RemoveOrganizationMember removes a member from an organization
func (c *Client) RemoveOrganizationMember(orgname, membername string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/members/%s", QuayURL, orgname, membername), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// Organization Repositories

// GetOrganizationRepositories retrieves repositories for an organization
func (c *Client) GetOrganizationRepositories(orgname string) (*OrganizationRepositories, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/repositories", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var repos OrganizationRepositories
	if err := c.get(req, &repos); err != nil {
		return nil, err
	}

	return &repos, nil
}

// Default Permissions (Prototypes)

// GetDefaultPermissions retrieves default permissions for an organization
func (c *Client) GetDefaultPermissions(orgname string) (*DefaultPermissions, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/prototypes", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var perms DefaultPermissions
	if err := c.get(req, &perms); err != nil {
		return nil, err
	}

	return &perms, nil
}

// CreateDefaultPermission creates a default permission for an organization
func (c *Client) CreateDefaultPermission(orgname, role, delegateType, delegateName string) (*DefaultPermission, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/organization/%s/prototypes", QuayURL, orgname), map[string]interface{}{
		"role": role,
		"delegate": map[string]interface{}{
			"kind": delegateType,
			"name": delegateName,
		},
	})
	if err != nil {
		return nil, err
	}

	var perm DefaultPermission
	if err := c.post(req, &perm); err != nil {
		return nil, err
	}

	return &perm, nil
}

// DeleteDefaultPermission deletes a default permission
func (c *Client) DeleteDefaultPermission(orgname, prototypeid string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/prototypes/%s", QuayURL, orgname, prototypeid), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// Proxy Cache Configuration

// GetProxyCacheConfig retrieves proxy cache configuration for an organization
func (c *Client) GetProxyCacheConfig(orgname string) (*ProxyCacheConfig, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/proxycache", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var config ProxyCacheConfig
	if err := c.get(req, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateProxyCacheConfig creates proxy cache configuration for an organization
func (c *Client) CreateProxyCacheConfig(orgname, upstreamRegistry string, insecure bool, expiration int) (*ProxyCacheConfig, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/organization/%s/proxycache", QuayURL, orgname), map[string]interface{}{
		"upstream_registry": upstreamRegistry,
		"insecure":          insecure,
		"expiration":        expiration,
	})
	if err != nil {
		return nil, err
	}

	var config ProxyCacheConfig
	if err := c.post(req, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// DeleteProxyCacheConfig deletes proxy cache configuration for an organization
func (c *Client) DeleteProxyCacheConfig(orgname string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/proxycache", QuayURL, orgname), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// Team Management

// GetTeams retrieves all teams for an organization
func (c *Client) GetTeams(orgname string) ([]Team, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/teams", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Teams []Team `json:"teams"`
	}
	if err := c.get(req, &response); err != nil {
		return nil, err
	}

	return response.Teams, nil
}

// CreateTeam creates a new team in an organization
func (c *Client) CreateTeam(orgname, teamname, description, role string) (*Team, error) {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s/team/%s", QuayURL, orgname, teamname), CreateTeamRequest{
		Name:        teamname,
		Description: description,
		Role:        role,
	})
	if err != nil {
		return nil, err
	}

	var team Team
	if err := c.put(req, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// GetTeam retrieves details for a specific team
func (c *Client) GetTeam(orgname, teamname string) (*Team, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/team/%s", QuayURL, orgname, teamname), nil)
	if err != nil {
		return nil, err
	}

	var team Team
	if err := c.get(req, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// DeleteTeam deletes a team from an organization
func (c *Client) DeleteTeam(orgname, teamname string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/team/%s", QuayURL, orgname, teamname), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// UpdateTeam updates team settings
func (c *Client) UpdateTeam(orgname, teamname, description, role string) (*Team, error) {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s/team/%s", QuayURL, orgname, teamname), map[string]interface{}{
		"description": description,
		"role":        role,
	})
	if err != nil {
		return nil, err
	}

	var team Team
	if err := c.put(req, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// Team Members Management

// GetTeamMembers retrieves members of a team
func (c *Client) GetTeamMembers(orgname, teamname string) (*TeamMembers, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/team/%s/members", QuayURL, orgname, teamname), nil)
	if err != nil {
		return nil, err
	}

	var members TeamMembers
	if err := c.get(req, &members); err != nil {
		return nil, err
	}

	return &members, nil
}

// AddTeamMember adds a member to a team
func (c *Client) AddTeamMember(orgname, teamname, membername string) error {
	req, err := newRequest("PUT", fmt.Sprintf("%s/organization/%s/team/%s/members/%s", QuayURL, orgname, teamname, membername), nil)
	if err != nil {
		return err
	}

	return c.put(req, nil)
}

// RemoveTeamMember removes a member from a team
func (c *Client) RemoveTeamMember(orgname, teamname, membername string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/team/%s/members/%s", QuayURL, orgname, teamname, membername), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// Team Permissions Management

// GetTeamPermissions retrieves repository permissions for a team
func (c *Client) GetTeamPermissions(orgname, teamname string) (*TeamPermissions, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/team/%s/permissions", QuayURL, orgname, teamname), nil)
	if err != nil {
		return nil, err
	}

	var perms TeamPermissions
	if err := c.get(req, &perms); err != nil {
		return nil, err
	}

	return &perms, nil
}

// SetTeamRepositoryPermission sets repository permission for a team
func (c *Client) SetTeamRepositoryPermission(orgname, teamname, repository, role string) error {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s/team/%s/permissions/%s", QuayURL, orgname, teamname, repository), map[string]interface{}{
		"role": role,
	})
	if err != nil {
		return err
	}

	return c.put(req, nil)
}

// RemoveTeamRepositoryPermission removes repository permission for a team
func (c *Client) RemoveTeamRepositoryPermission(orgname, teamname, repository string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/team/%s/permissions/%s", QuayURL, orgname, teamname, repository), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// Robot Account Management

// GetRobotAccounts retrieves all robot accounts for an organization
func (c *Client) GetRobotAccounts(orgname string) (*RobotAccounts, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/robots", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var robots RobotAccounts
	if err := c.get(req, &robots); err != nil {
		return nil, err
	}

	return &robots, nil
}

// CreateRobotAccount creates a new robot account in an organization
func (c *Client) CreateRobotAccount(orgname, robotShortname, description string, unstructured map[string]interface{}) (*RobotAccount, error) {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s/robots/%s", QuayURL, orgname, robotShortname), CreateRobotRequest{
		Description:  description,
		Unstructured: unstructured,
	})
	if err != nil {
		return nil, err
	}

	var robot RobotAccount
	if err := c.put(req, &robot); err != nil {
		return nil, err
	}

	return &robot, nil
}

// GetRobotAccount retrieves details for a specific robot account
func (c *Client) GetRobotAccount(orgname, robotShortname string) (*RobotAccount, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/robots/%s", QuayURL, orgname, robotShortname), nil)
	if err != nil {
		return nil, err
	}

	var robot RobotAccount
	if err := c.get(req, &robot); err != nil {
		return nil, err
	}

	return &robot, nil
}

// DeleteRobotAccount deletes a robot account from an organization
func (c *Client) DeleteRobotAccount(orgname, robotShortname string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/robots/%s", QuayURL, orgname, robotShortname), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// RegenerateRobotToken regenerates the token for a robot account
func (c *Client) RegenerateRobotToken(orgname, robotShortname string) (*RobotAccount, error) {
	req, err := newRequest("POST", fmt.Sprintf("%s/organization/%s/robots/%s/regenerate", QuayURL, orgname, robotShortname), nil)
	if err != nil {
		return nil, err
	}

	var robot RobotAccount
	if err := c.post(req, &robot); err != nil {
		return nil, err
	}

	return &robot, nil
}

// GetRobotPermissions retrieves repository permissions for a robot account
func (c *Client) GetRobotPermissions(orgname, robotShortname string) (*RobotPermissions, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/robots/%s/permissions", QuayURL, orgname, robotShortname), nil)
	if err != nil {
		return nil, err
	}

	var perms RobotPermissions
	if err := c.get(req, &perms); err != nil {
		return nil, err
	}

	return &perms, nil
}

// SetRobotRepositoryPermission sets repository permission for a robot account
func (c *Client) SetRobotRepositoryPermission(orgname, robotShortname, repository, role string) error {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s/robots/%s/permissions/%s", QuayURL, orgname, robotShortname, repository), map[string]interface{}{
		"role": role,
	})
	if err != nil {
		return err
	}

	return c.put(req, nil)
}

// RemoveRobotRepositoryPermission removes repository permission for a robot account
func (c *Client) RemoveRobotRepositoryPermission(orgname, robotShortname, repository string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/robots/%s/permissions/%s", QuayURL, orgname, robotShortname, repository), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// Quota Management

// GetQuota retrieves quota information for an organization
func (c *Client) GetQuota(orgname string) (*Quota, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/quota", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var quota Quota
	if err := c.get(req, &quota); err != nil {
		return nil, err
	}

	return &quota, nil
}

// CreateQuota creates a quota for an organization
func (c *Client) CreateQuota(orgname string, limitBytes int64) (*Quota, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/organization/%s/quota", QuayURL, orgname), CreateQuotaRequest{
		LimitBytes: limitBytes,
	})
	if err != nil {
		return nil, err
	}

	var quota Quota
	if err := c.post(req, &quota); err != nil {
		return nil, err
	}

	return &quota, nil
}

// UpdateQuota updates quota limits for an organization
func (c *Client) UpdateQuota(orgname string, limitBytes int64) (*Quota, error) {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s/quota", QuayURL, orgname), CreateQuotaRequest{
		LimitBytes: limitBytes,
	})
	if err != nil {
		return nil, err
	}

	var quota Quota
	if err := c.put(req, &quota); err != nil {
		return nil, err
	}

	return &quota, nil
}

// DeleteQuota deletes quota for an organization
func (c *Client) DeleteQuota(orgname string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/quota", QuayURL, orgname), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// Auto-Prune Policy Management

// GetAutoPrunePolicies retrieves auto-prune policies for an organization
func (c *Client) GetAutoPrunePolicies(orgname string) (*AutoPrunePolicies, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/autoprunepolicy", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var policies AutoPrunePolicies
	if err := c.get(req, &policies); err != nil {
		return nil, err
	}

	return &policies, nil
}

// CreateAutoPrunePolicy creates an auto-prune policy for an organization
func (c *Client) CreateAutoPrunePolicy(orgname, method string, value int, tagPattern string) (*AutoPrunePolicy, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/organization/%s/autoprunepolicy", QuayURL, orgname), CreateAutoPruneRequest{
		Method:     method,
		Value:      value,
		TagPattern: tagPattern,
	})
	if err != nil {
		return nil, err
	}

	var policy AutoPrunePolicy
	if err := c.post(req, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// GetAutoPrunePolicy retrieves a specific auto-prune policy
func (c *Client) GetAutoPrunePolicy(orgname, policyUUID string) (*AutoPrunePolicy, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/autoprunepolicy/%s", QuayURL, orgname, policyUUID), nil)
	if err != nil {
		return nil, err
	}

	var policy AutoPrunePolicy
	if err := c.get(req, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// UpdateAutoPrunePolicy updates an auto-prune policy
func (c *Client) UpdateAutoPrunePolicy(orgname, policyUUID, method string, value int, tagPattern string) (*AutoPrunePolicy, error) {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s/autoprunepolicy/%s", QuayURL, orgname, policyUUID), CreateAutoPruneRequest{
		Method:     method,
		Value:      value,
		TagPattern: tagPattern,
	})
	if err != nil {
		return nil, err
	}

	var policy AutoPrunePolicy
	if err := c.put(req, &policy); err != nil {
		return nil, err
	}

	return &policy, nil
}

// DeleteAutoPrunePolicy deletes an auto-prune policy
func (c *Client) DeleteAutoPrunePolicy(orgname, policyUUID string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/autoprunepolicy/%s", QuayURL, orgname, policyUUID), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// Applications Management

// GetApplications retrieves applications for an organization
func (c *Client) GetApplications(orgname string) (*Applications, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/applications", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var apps Applications
	if err := c.get(req, &apps); err != nil {
		return nil, err
	}

	return &apps, nil
}

// CreateApplication creates a new application for an organization
func (c *Client) CreateApplication(orgname, name, description, applicationURI, redirectURI string) (*Application, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/organization/%s/applications", QuayURL, orgname), CreateApplicationRequest{
		Name:           name,
		Description:    description,
		ApplicationURI: applicationURI,
		RedirectURI:    redirectURI,
	})
	if err != nil {
		return nil, err
	}

	var app Application
	if err := c.post(req, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

// GetApplication retrieves details for a specific application
func (c *Client) GetApplication(orgname, clientID string) (*Application, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/applications/%s", QuayURL, orgname, clientID), nil)
	if err != nil {
		return nil, err
	}

	var app Application
	if err := c.get(req, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

// UpdateApplication updates an application
func (c *Client) UpdateApplication(orgname, clientID, name, description, applicationURI, redirectURI string) (*Application, error) {
	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/organization/%s/applications/%s", QuayURL, orgname, clientID), CreateApplicationRequest{
		Name:           name,
		Description:    description,
		ApplicationURI: applicationURI,
		RedirectURI:    redirectURI,
	})
	if err != nil {
		return nil, err
	}

	var app Application
	if err := c.put(req, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

// DeleteApplication deletes an application
func (c *Client) DeleteApplication(orgname, clientID string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/applications/%s", QuayURL, orgname, clientID), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// ResetApplicationClientSecret resets the client secret for an application
func (c *Client) ResetApplicationClientSecret(orgname, clientID string) (*Application, error) {
	req, err := newRequest("POST", fmt.Sprintf("%s/organization/%s/applications/%s/resetclientsecret", QuayURL, orgname, clientID), nil)
	if err != nil {
		return nil, err
	}

	var app Application
	if err := c.post(req, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

// GetOrganizationCollaborators gets the list of collaborators for an organization
func (c *Client) GetOrganizationCollaborators(orgname string) (*Collaborators, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/collaborators", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var collaborators Collaborators
	if err := c.get(req, &collaborators); err != nil {
		return nil, err
	}

	return &collaborators, nil
}

// GetOrganizationMember gets information about a specific organization member
func (c *Client) GetOrganizationMember(orgname, membername string) (*OrganizationMember, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/members/%s", QuayURL, orgname, membername), nil)
	if err != nil {
		return nil, err
	}

	var member OrganizationMember
	if err := c.get(req, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// GetOrganizationMarketplace gets marketplace information for an organization
func (c *Client) GetOrganizationMarketplace(orgname string) (*MarketplaceInfo, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/marketplace", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	var marketplace MarketplaceInfo
	if err := c.get(req, &marketplace); err != nil {
		return nil, err
	}

	return &marketplace, nil
}

// CreateOrganizationMarketplaceSubscription creates a marketplace subscription
func (c *Client) CreateOrganizationMarketplaceSubscription(orgname string, subscription *MarketplaceSubscriptionRequest) error {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/organization/%s/marketplace", QuayURL, orgname), subscription)
	if err != nil {
		return err
	}

	return c.post(req, nil)
}

// BatchRemoveOrganizationMarketplaceSubscriptions removes multiple marketplace subscriptions
func (c *Client) BatchRemoveOrganizationMarketplaceSubscriptions(orgname string, subscriptionIDs []string) error {
	body := struct {
		SubscriptionIDs []string `json:"subscription_ids"`
	}{
		SubscriptionIDs: subscriptionIDs,
	}
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/organization/%s/marketplace/batchremove", QuayURL, orgname), body)
	if err != nil {
		return err
	}

	return c.post(req, nil)
}

// DeleteOrganizationMarketplaceSubscription removes a specific marketplace subscription
func (c *Client) DeleteOrganizationMarketplaceSubscription(orgname, subscriptionID string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/marketplace/%s", QuayURL, orgname, subscriptionID), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}

// InviteTeamMember invites a member to a team via email
func (c *Client) InviteTeamMember(orgname, teamname, email string) error {
	req, err := newRequest("PUT", fmt.Sprintf("%s/organization/%s/team/%s/invite/%s", QuayURL, orgname, teamname, email), nil)
	if err != nil {
		return err
	}

	return c.put(req, nil)
}

// DeleteTeamInvite deletes a pending team invitation
func (c *Client) DeleteTeamInvite(orgname, teamname, email string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/organization/%s/team/%s/invite/%s", QuayURL, orgname, teamname, email), nil)
	if err != nil {
		return err
	}

	return c.delete(req)
}
