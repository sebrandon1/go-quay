# lib

This folder contains Go source files for Quay.io API endpoints found in: 

https://docs.quay.io/api/swagger/

## Available Endpoints

### Organization Management

#### Core Organization Operations
- `CreateOrganization(name, email string) (*Organization, error)` - Create a new organization
- `GetOrganization(orgname string) (*Organization, error)` - Get organization details
- `UpdateOrganization(orgname, email string) (*Organization, error)` - Update organization settings
- `DeleteOrganization(orgname string) error` - Delete an organization

#### Organization Members
- `GetOrganizationMembers(orgname string) (*OrganizationMembers, error)` - Get organization members
- `AddOrganizationMember(orgname, membername string) error` - Add member to organization
- `RemoveOrganizationMember(orgname, membername string) error` - Remove member from organization

#### Organization Repositories
- `GetOrganizationRepositories(orgname string) (*OrganizationRepositories, error)` - Get organization repositories

#### Organization Billing
- `GetOrganizationBilling(orgname string) (*BillingInfo, error)` - Get organization billing information
- `GetOrganizationSubscription(orgname string) (*Subscription, error)` - Get organization subscription details
- `GetOrganizationInvoices(orgname string) ([]Invoice, error)` - Get organization invoices

#### Organization Logs
- `GetOrganizationLogs(namespace, nextPage string) (*OrganizationLogs, error)` - Get organization logs

### Team Management

#### Team Operations
- `GetTeams(orgname string) ([]Team, error)` - Get all teams in organization
- `CreateTeam(orgname, teamname, description, role string) (*Team, error)` - Create new team
- `GetTeam(orgname, teamname string) (*Team, error)` - Get team details
- `UpdateTeam(orgname, teamname, description, role string) (*Team, error)` - Update team settings
- `DeleteTeam(orgname, teamname string) error` - Delete team

#### Team Members
- `GetTeamMembers(orgname, teamname string) (*TeamMembers, error)` - Get team members
- `AddTeamMember(orgname, teamname, membername string) error` - Add member to team
- `RemoveTeamMember(orgname, teamname, membername string) error` - Remove member from team

#### Team Permissions
- `GetTeamPermissions(orgname, teamname string) (*TeamPermissions, error)` - Get team repository permissions
- `SetTeamRepositoryPermission(orgname, teamname, repository, role string) error` - Set repository permission for team
- `RemoveTeamRepositoryPermission(orgname, teamname, repository string) error` - Remove repository permission from team

### Robot Account Management

#### Robot Account Operations
- `GetRobotAccounts(orgname string) (*RobotAccounts, error)` - Get all robot accounts
- `CreateRobotAccount(orgname, robotShortname, description string, unstructured map[string]interface{}) (*RobotAccount, error)` - Create robot account
- `GetRobotAccount(orgname, robotShortname string) (*RobotAccount, error)` - Get robot account details
- `DeleteRobotAccount(orgname, robotShortname string) error` - Delete robot account
- `RegenerateRobotToken(orgname, robotShortname string) (*RobotAccount, error)` - Regenerate robot token

#### Robot Permissions
- `GetRobotPermissions(orgname, robotShortname string) (*RobotPermissions, error)` - Get robot repository permissions
- `SetRobotRepositoryPermission(orgname, robotShortname, repository, role string) error` - Set repository permission for robot
- `RemoveRobotRepositoryPermission(orgname, robotShortname, repository string) error` - Remove repository permission from robot

### Quota Management

- `GetQuota(orgname string) (*Quota, error)` - Get organization quota information
- `CreateQuota(orgname string, limitBytes int64) (*Quota, error)` - Create quota for organization
- `UpdateQuota(orgname string, limitBytes int64) (*Quota, error)` - Update quota limits
- `DeleteQuota(orgname string) error` - Delete organization quota

### Auto-Prune Policy Management

- `GetAutoPrunePolicies(orgname string) (*AutoPrunePolicies, error)` - Get auto-prune policies
- `CreateAutoPrunePolicy(orgname, method string, value int, tagPattern string) (*AutoPrunePolicy, error)` - Create auto-prune policy
- `GetAutoPrunePolicy(orgname, policyUUID string) (*AutoPrunePolicy, error)` - Get specific auto-prune policy
- `UpdateAutoPrunePolicy(orgname, policyUUID, method string, value int, tagPattern string) (*AutoPrunePolicy, error)` - Update auto-prune policy
- `DeleteAutoPrunePolicy(orgname, policyUUID string) error` - Delete auto-prune policy

### Applications Management

- `GetApplications(orgname string) (*Applications, error)` - Get organization applications
- `CreateApplication(orgname, name, description, applicationURI, redirectURI string) (*Application, error)` - Create OAuth application
- `GetApplication(orgname, clientID string) (*Application, error)` - Get application details
- `UpdateApplication(orgname, clientID, name, description, applicationURI, redirectURI string) (*Application, error)` - Update application
- `DeleteApplication(orgname, clientID string) error` - Delete application
- `ResetApplicationClientSecret(orgname, clientID string) (*Application, error)` - Reset application client secret

### Default Permissions Management

- `GetDefaultPermissions(orgname string) (*DefaultPermissions, error)` - Get default repository permissions
- `CreateDefaultPermission(orgname, role, delegateType, delegateName string) (*DefaultPermission, error)` - Create default permission
- `DeleteDefaultPermission(orgname, prototypeid string) error` - Delete default permission

### Proxy Cache Configuration

- `GetProxyCacheConfig(orgname string) (*ProxyCacheConfig, error)` - Get proxy cache configuration
- `CreateProxyCacheConfig(orgname, upstreamRegistry string, insecure bool, expiration int) (*ProxyCacheConfig, error)` - Create proxy cache config
- `DeleteProxyCacheConfig(orgname string) error` - Delete proxy cache configuration

### Repository Management

- `GetRepository(namespace, repository string) (RepositoryWithTags, error)` - Get repository with tags information

### Logs Management

- `GetAggregatedLogs(namespace, repository, startDate, endDate string) (*AggregatedLogs, error)` - Get aggregated repository logs
- `GetLogs(namespace, repository, nextPage string) (*Logs, error)` - Get repository logs

### Billing Management

- `GetUserBilling() (*BillingInfo, error)` - Get user billing information
- `GetUserSubscription() (*Subscription, error)` - Get user subscription details
- `GetUserInvoices() ([]Invoice, error)` - Get user invoices (Not available in Quay API)
- `GetAvailablePlans() ([]Subscription, error)` - Get available subscription plans

## Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    // Initialize client with bearer token
    client, err := lib.NewClient("your-bearer-token")
    if err != nil {
        log.Fatal(err)
    }

    // Organization Management
    org, err := client.CreateOrganization("my-org", "admin@example.com")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created organization: %s\n", org.Name)

    // Team Management
    team, err := client.CreateTeam("my-org", "developers", "Development team", "member")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created team: %s\n", team.Name)

    // Robot Account Management
    robot, err := client.CreateRobotAccount("my-org", "ci-bot", "CI/CD robot", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created robot: %s\n", robot.Name)

    // Set permissions
    err = client.SetTeamRepositoryPermission("my-org", "developers", "my-repo", "write")
    if err != nil {
        log.Fatal(err)
    }

    // Quota Management
    quota, err := client.CreateQuota("my-org", 1073741824) // 1GB
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created quota: %d bytes\n", quota.LimitBytes)
}
```

## Authentication

All endpoints require authentication via bearer token. You can obtain a token from:

1. **OAuth 2.0 Access Token**: Created via Quay.io UI (10-year lifespan)
2. **Robot Account Token**: Persistent, non-expiring by default
3. **OCI Referrers OAuth Access Token**: For OCI operations

## Rate Limiting

Quay.io implements rate limiting (few requests per second per IP) with bursting capabilities. If you exceed limits, endpoints will return HTTP 429 "Too many requests".

## Error Handling

All functions return appropriate errors for:
- Authentication failures
- Rate limiting
- Network issues
- API errors (400, 404, 500 series)

## Roles and Permissions

### Organization Roles
- `admin`: Full administrative access
- `member`: Inherits team permissions
- `creator`: Member permissions + repository creation

### Repository Roles
- `admin`: Full repository access + administrative tasks
- `write`: Read and write access
- `read`: Read-only access
- `none`: No access

### Auto-Prune Methods
- `number_of_tags`: Keep specified number of tags
- `creation_date`: Keep tags newer than specified days
- `tag_pattern`: Keep tags matching pattern

## Data Structures

All response structures are defined in `structs.go` and include comprehensive field mappings for JSON serialization/deserialization.
