# go-quay Library Guide

This guide covers how to use go-quay as a Go library in your applications.

## Installation

```bash
go get github.com/sebrandon1/go-quay
```

## Client Initialization

```go
package main

import (
    "log"
    "os"

    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    // Get token from environment (recommended)
    token := os.Getenv("QUAY_TOKEN")

    // Create client
    client, err := lib.NewClient(token)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Use client...
}
```

## API Categories

### User Operations

```go
// Get current user
user, err := client.GetUser()

// Get user by username
user, err := client.GetUserByUsername("johndoe")

// Get starred repositories
starred, err := client.GetStarredRepositories()

// Star/unstar repositories
err := client.StarRepository("namespace", "repo")
err := client.UnstarRepository("namespace", "repo")

// User robots
robots, err := client.GetUserRobots()
robot, err := client.CreateUserRobot("name", "description")
robot, err := client.GetUserRobot("name")
robot, err := client.RegenerateUserRobotToken("name")
err := client.DeleteUserRobot("name")
```

### Repository Operations

```go
// List repositories
repos, err := client.ListRepositories(namespace, public, starred, page, limit)

// CRUD operations
repo, err := client.CreateRepository(namespace, name, visibility, description)
repo, err := client.GetRepository(namespace, name)
repo, err := client.UpdateRepository(namespace, name, description, visibility)
err := client.DeleteRepository(namespace, name)

// Change visibility
err := client.ChangeRepositoryVisibility(namespace, name, "public")
```

### Tag Operations

```go
// Get tag info
tag, err := client.GetTag(namespace, repo, tagName)

// Update tag (set expiration)
tag, err := client.UpdateTag(namespace, repo, tagName, expiration)

// Delete tag
err := client.DeleteTag(namespace, repo, tagName)

// Tag history
history, err := client.GetTagHistory(namespace, repo, tagName)

// Restore tag
err := client.RestoreTag(namespace, repo, tagName, manifestDigest)

// Revert tag
tag, err := client.RevertTag(namespace, repo, tagName, manifestDigest)
```

### Manifest Operations

```go
// Get manifest
manifest, err := client.GetManifest(namespace, repo, digest)

// Delete manifest
err := client.DeleteManifest(namespace, repo, digest)

// Labels
labels, err := client.GetManifestLabels(namespace, repo, digest)
label, err := client.GetManifestLabel(namespace, repo, digest, labelID)
err := client.AddManifestLabel(namespace, repo, digest, key, value, mediaType)
err := client.DeleteManifestLabel(namespace, repo, digest, labelID)
```

### Security Scanning

```go
// Get security scan results
security, err := client.GetManifestSecurity(namespace, repo, digest, includeVulns)

// Access vulnerability data
if security.Status == "scanned" && security.Data != nil {
    for _, feature := range security.Data.Layer.Features {
        for _, vuln := range feature.Vulnerabilities {
            // Process vulnerability
        }
    }
}
```

### Organization Operations

```go
// CRUD
org, err := client.CreateOrganization(name, email)
org, err := client.GetOrganization(name)
org, err := client.UpdateOrganization(name, email)
err := client.DeleteOrganization(name)

// Members
members, err := client.GetOrganizationMembers(orgname)
member, err := client.GetOrganizationMember(orgname, membername)

// Collaborators
collaborators, err := client.GetOrganizationCollaborators(orgname)

// Robots
robots, err := client.GetOrganizationRobots(orgname)
robot, err := client.CreateOrganizationRobot(orgname, name, description, metadata)
robot, err := client.GetOrganizationRobot(orgname, name)
robot, err := client.RegenerateOrganizationRobotToken(orgname, name)
err := client.DeleteOrganizationRobot(orgname, name)
```

### Team Operations

```go
// CRUD
teams, err := client.GetOrganizationTeams(orgname)
team, err := client.CreateTeam(orgname, teamname, description, role)
team, err := client.GetTeam(orgname, teamname)
team, err := client.UpdateTeam(orgname, teamname, description, role)
err := client.DeleteTeam(orgname, teamname)

// Members
members, err := client.GetTeamMembers(orgname, teamname)
err := client.AddTeamMember(orgname, teamname, membername)
err := client.RemoveTeamMember(orgname, teamname, membername)

// Permissions
perms, err := client.GetTeamPermissions(orgname, teamname)
err := client.SetTeamRepositoryPermission(orgname, teamname, repo, role)
err := client.RemoveTeamRepositoryPermission(orgname, teamname, repo)
```

### Permission Operations

```go
// User permissions
perms, err := client.ListUserPermissions(namespace, repo)
perm, err := client.GetUserPermission(namespace, repo, username)
err := client.SetUserPermission(namespace, repo, username, role)
err := client.DeleteUserPermission(namespace, repo, username)

// Team permissions
perms, err := client.ListTeamPermissions(namespace, repo)
perm, err := client.GetTeamPermission(namespace, repo, teamname)
err := client.SetTeamPermission(namespace, repo, teamname, role)
err := client.DeleteTeamPermission(namespace, repo, teamname)
```

### Build Operations

```go
// List builds
builds, err := client.GetRepositoryBuilds(namespace, repo, limit, page)

// Get build
build, err := client.GetBuild(namespace, repo, buildUUID)

// Get build logs
logs, err := client.GetBuildLogs(namespace, repo, buildUUID)

// Get build status
status, err := client.GetBuildStatus(namespace, repo, buildUUID)

// Request new build
build, err := client.RequestBuild(namespace, repo, request)

// Cancel build
err := client.CancelBuild(namespace, repo, buildUUID)
```

### Build Trigger Operations

```go
// List triggers
triggers, err := client.ListBuildTriggers(namespace, repo)

// Get trigger
trigger, err := client.GetBuildTrigger(namespace, repo, triggerUUID)

// Activate trigger
err := client.ActivateBuildTrigger(namespace, repo, triggerUUID, config)

// Start build from trigger
build, err := client.StartBuildTrigger(namespace, repo, triggerUUID, commitSHA)

// Enable/disable trigger
err := client.EnableBuildTrigger(namespace, repo, triggerUUID)
err := client.DisableBuildTrigger(namespace, repo, triggerUUID)

// Delete trigger
err := client.DeleteBuildTrigger(namespace, repo, triggerUUID)
```

### Notification Operations

```go
// List notifications
notifications, err := client.ListNotifications(namespace, repo)

// Get notification
notification, err := client.GetNotification(namespace, repo, uuid)

// Create notification
notification, err := client.CreateNotification(namespace, repo, lib.CreateNotificationRequest{
    Event:  "repo_push",
    Method: "webhook",
    Title:  "My Notification",
    Config: lib.NotificationConfig{URL: "https://example.com/webhook"},
})

// Update notification
notification, err := client.UpdateNotification(namespace, repo, uuid, request)

// Test notification
err := client.TestNotification(namespace, repo, uuid)

// Reset failure count
err := client.ResetNotification(namespace, repo, uuid)

// Delete notification
err := client.DeleteNotification(namespace, repo, uuid)
```

### Quota Operations

```go
// Get quota
quotas, err := client.GetOrganizationQuota(orgname)

// Create quota
quota, err := client.CreateOrganizationQuota(orgname, limitBytes)

// Update quota
quota, err := client.UpdateOrganizationQuota(orgname, quotaID, limitBytes)

// Delete quota
err := client.DeleteOrganizationQuota(orgname, quotaID)
```

### Auto-Prune Operations

```go
// Get policies
policies, err := client.GetOrganizationAutoPrunePolicies(orgname)

// Create policy
policy, err := client.CreateOrganizationAutoPrunePolicy(orgname, method, value, tagPattern)

// Get specific policy
policy, err := client.GetOrganizationAutoPrunePolicy(orgname, policyUUID)

// Update policy
policy, err := client.UpdateOrganizationAutoPrunePolicy(orgname, policyUUID, method, value, tagPattern)

// Delete policy
err := client.DeleteOrganizationAutoPrunePolicy(orgname, policyUUID)
```

### Logs Operations

```go
// Repository logs
logs, err := client.GetRepositoryLogs(namespace, repo, nextPage)
aggLogs, err := client.GetAggregatedLogs(namespace, repo, startDate, endDate)

// Organization logs
logs, err := client.GetOrganizationLogs(orgname, nextPage)
aggLogs, err := client.GetOrganizationAggregatedLogs(orgname, startDate, endDate)

// User logs
logs, err := client.GetUserLogs(nextPage)
aggLogs, err := client.GetUserAggregatedLogs(startDate, endDate)

// Export logs
err := client.ExportRepositoryLogs(namespace, repo, startDate, endDate, email)
err := client.ExportOrganizationLogs(orgname, startDate, endDate, email)
err := client.ExportUserLogs(startDate, endDate, email)
```

### Search Operations

```go
// Search repositories
results, err := client.SearchRepositories(query, page)

// Search all (repos, users, teams, etc.)
results, err := client.SearchAll(query, page)
```

### Billing Operations

```go
// Plans
plans, err := client.GetAvailablePlans()

// User billing
billing, err := client.GetUserBilling()
subscription, err := client.GetUserSubscription()

// Organization billing
billing, err := client.GetOrganizationBilling(orgname)
subscription, err := client.GetOrganizationSubscription(orgname)
invoices, err := client.GetOrganizationInvoices(orgname)
```

## Error Handling

```go
result, err := client.GetRepository("namespace", "repo")
if err != nil {
    // Check for specific errors
    errStr := err.Error()

    switch {
    case strings.Contains(errStr, "404"):
        // Not found
    case strings.Contains(errStr, "401"):
        // Authentication failed
    case strings.Contains(errStr, "403"):
        // Permission denied
    case strings.Contains(errStr, "429"):
        // Rate limited - implement backoff
    default:
        // Other error
    }
}
```

## Rate Limiting

Quay.io implements rate limiting. Implement exponential backoff:

```go
func withRetry(fn func() error, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = fn()
        if err == nil {
            return nil
        }
        if strings.Contains(err.Error(), "429") {
            time.Sleep(time.Duration(1<<i) * time.Second)
            continue
        }
        return err
    }
    return err
}
```

## Best Practices

1. **Environment Variables** - Store tokens in environment variables, not code
2. **Error Handling** - Always check errors and handle appropriately
3. **Rate Limiting** - Implement backoff for rate limit errors
4. **Pagination** - Use pagination for large result sets
5. **Minimal Permissions** - Use tokens with minimal required permissions

## Examples

See the [examples directory](../examples/) for complete working programs:

- [basic-usage](../examples/basic-usage/) - Getting started
- [security-scan](../examples/security-scan/) - Vulnerability scanning
- [ci-cd-integration](../examples/ci-cd-integration/) - CI/CD automation
- [organization-management](../examples/organization-management/) - Org admin

## Tutorials

Step-by-step guides in the [tutorials directory](./tutorials/):

1. [Getting Started](./tutorials/01-getting-started.md)
2. [Repository Management](./tutorials/02-repository-management.md)
3. [Security Scanning](./tutorials/03-security-scanning.md)
4. [CI/CD Automation](./tutorials/04-ci-cd-automation.md)
5. [Organization Administration](./tutorials/05-organization-admin.md)
