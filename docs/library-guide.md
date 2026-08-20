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

    // Create client (quay.io)
    client, err := lib.NewClient(token)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Self-hosted Quay:
    // client, err := lib.NewClientWithURL(token, "https://quay.example.com/api/v1")

    // Optional retries on 429 and 5xx:
    // client.Retry = &lib.RetryConfig{
    //     MaxRetries:     3,
    //     InitialBackoff: time.Second,
    //     MaxBackoff:     30 * time.Second,
    // }

    // Use client...
}
```

Every API method takes `context.Context` as its first argument. Use `context.Background()`, `context.WithTimeout`, or a request-scoped context from your caller. Canceling the context aborts in-flight HTTP requests and retry backoff.

The snippets below use `ctx` for that argument.

## API Categories

### User Operations

```go
// Get current user
user, err := client.GetUser(ctx)

// Get user by username
user, err := client.GetUserByUsername(ctx, "johndoe")

// Get starred repositories
starred, err := client.GetStarredRepositories(ctx)

// Star/unstar repositories
err := client.StarRepository(ctx, "namespace", "repo")
err := client.UnstarRepository(ctx, "namespace", "repo")

// User robots
robots, err := client.GetUserRobotAccounts(ctx)
robot, err := client.CreateUserRobotAccount(ctx, "name", "description", nil)
robot, err := client.GetUserRobotAccount(ctx, "name")
robot, err := client.RegenerateUserRobotToken(ctx, "name")
err := client.DeleteUserRobotAccount(ctx, "name")
perms, err := client.GetUserRobotPermissions(ctx, "name")

// User robot federation
federation, err := client.GetUserRobotFederation(ctx, "name")
err := client.CreateUserRobotFederation(ctx, "name", []lib.RobotFederationConfig{{Issuer: "issuer", Subject: "subject"}})
err := client.DeleteUserRobotFederation(ctx, "name")

// User marketplace
marketplace, err := client.GetUserMarketplace(ctx)
```

### Repository Operations

```go
// List repositories
repos, err := client.ListRepositories(ctx, namespace, public, starred, popularity, page, limit)

// CRUD operations
repo, err := client.CreateRepository(ctx, namespace, name, visibility, description)
repo, err := client.GetRepository(ctx, namespace, name)
repo, err := client.UpdateRepository(ctx, namespace, name, description, visibility)
err := client.DeleteRepository(ctx, namespace, name)

// Change visibility
err := client.ChangeRepositoryVisibility(ctx, namespace, name, "public")

// Paginate automatically
allRepos, err := client.ListAllRepositories(ctx, namespace, public, starred, popularity)
allTags, err := client.ListAllTags(ctx, namespace, name, onlyActive)
page, err := client.ListTagsPage(ctx, namespace, name, limit, pageNum, onlyActive)
```

### Tag Operations

```go
// Get tag info
tag, err := client.GetTag(ctx, namespace, repo, tagName)

// Update tag (set expiration)
tag, err := client.UpdateTag(ctx, namespace, repo, tagName, expiration)

// Delete tag
err := client.DeleteTag(ctx, namespace, repo, tagName)

// Tag history
history, err := client.GetTagHistory(ctx, namespace, repo, tagName)

// Restore tag
err := client.RestoreTag(ctx, namespace, repo, tagName, manifestDigest)

// Revert tag
tag, err := client.RevertTag(ctx, namespace, repo, tagName, manifestDigest)

// Create or move a tag to a digest
err := client.ChangeTag(ctx, namespace, repo, tagName, manifestDigest)
```

### Manifest Operations

```go
// Get manifest
manifest, err := client.GetManifest(ctx, namespace, repo, digest)

// Delete manifest
err := client.DeleteManifest(ctx, namespace, repo, digest)

// Labels
labels, err := client.GetManifestLabels(ctx, namespace, repo, digest)
label, err := client.GetManifestLabel(ctx, namespace, repo, digest, labelID)
label, err := client.AddManifestLabel(ctx, namespace, repo, digest, key, value, mediaType)
err := client.DeleteManifestLabel(ctx, namespace, repo, digest, labelID)
```

### Security Scanning

```go
// Get security scan results
security, err := client.GetManifestSecurity(ctx, namespace, repo, digest, includeVulns)

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
org, err := client.CreateOrganization(ctx, name, email)
org, err := client.GetOrganization(ctx, name)
org, err := client.UpdateOrganization(ctx, name, email)
err := client.DeleteOrganization(ctx, name)

// Members
members, err := client.GetOrganizationMembers(ctx, orgname)
member, err := client.GetOrganizationMember(ctx, orgname, membername)
err := client.AddOrganizationMember(ctx, orgname, membername)
err := client.RemoveOrganizationMember(ctx, orgname, membername)

// Collaborators
collaborators, err := client.GetOrganizationCollaborators(ctx, orgname)

// Repositories
repos, err := client.GetOrganizationRepositories(ctx, orgname)

// Robots
robots, err := client.GetRobotAccounts(ctx, orgname)
robot, err := client.CreateRobotAccount(ctx, orgname, name, description, nil)
robot, err := client.GetRobotAccount(ctx, orgname, name)
robot, err := client.RegenerateRobotToken(ctx, orgname, name)
err := client.DeleteRobotAccount(ctx, orgname, name)
perms, err := client.GetRobotPermissions(ctx, orgname, name)
err := client.SetRobotRepositoryPermission(ctx, orgname, name, repository, role)
err := client.RemoveRobotRepositoryPermission(ctx, orgname, name, repository)

// Robot federation
federation, err := client.GetRobotFederation(ctx, orgname, name)
err := client.CreateRobotFederation(ctx, orgname, name, []lib.RobotFederationConfig{{Issuer: "issuer", Subject: "subject"}})
err := client.DeleteRobotFederation(ctx, orgname, name)
```

### Team Operations

```go
// CRUD
teams, err := client.GetTeams(ctx, orgname)
team, err := client.CreateTeam(ctx, orgname, teamname, description, role)
team, err := client.GetTeam(ctx, orgname, teamname)
team, err := client.UpdateTeam(ctx, orgname, teamname, description, role)
err := client.DeleteTeam(ctx, orgname, teamname)

// Members
members, err := client.GetTeamMembers(ctx, orgname, teamname)
err := client.AddTeamMember(ctx, orgname, teamname, membername)
err := client.RemoveTeamMember(ctx, orgname, teamname, membername)

// Invitations
err := client.InviteTeamMember(ctx, orgname, teamname, email)
err := client.DeleteTeamInvite(ctx, orgname, teamname, email)

// Permissions
perms, err := client.GetTeamPermissions(ctx, orgname, teamname)
err := client.SetTeamRepositoryPermission(ctx, orgname, teamname, repo, role)
err := client.RemoveTeamRepositoryPermission(ctx, orgname, teamname, repo)
```

### Permission Operations

```go
// Combined user/robot permissions (used by CLI `permissions list/set/remove`)
perms, err := client.GetRepositoryPermissions(ctx, namespace, repo)
err := client.SetRepositoryPermission(ctx, namespace, repo, username, role)
err := client.RemoveRepositoryPermission(ctx, namespace, repo, username)

// User permissions
perms, err := client.ListUserPermissions(ctx, namespace, repo)
perm, err := client.GetUserPermission(ctx, namespace, repo, username)
err := client.SetUserPermission(ctx, namespace, repo, username, role)
err := client.DeleteUserPermission(ctx, namespace, repo, username)
perm, err := client.GetUserTransitivePermission(ctx, namespace, repo, username)

// Team permissions
perms, err := client.ListTeamPermissions(ctx, namespace, repo)
perm, err := client.GetTeamPermission(ctx, namespace, repo, teamname)
err := client.SetTeamPermission(ctx, namespace, repo, teamname, role)
err := client.DeleteTeamPermission(ctx, namespace, repo, teamname)
```

### Build Operations

```go
// List builds
builds, err := client.GetBuilds(ctx, namespace, repo, limit)

// Get build
build, err := client.GetBuild(ctx, namespace, repo, buildUUID)

// Get build logs
logs, err := client.GetBuildLogs(ctx, namespace, repo, buildUUID)

// Get build status
status, err := client.GetBuildStatus(ctx, namespace, repo, buildUUID)

// Request new build
build, err := client.RequestBuild(ctx, namespace, repo, &lib.RequestBuildRequest{...})

// Cancel build
err := client.CancelBuild(ctx, namespace, repo, buildUUID)
```

### Build Trigger Operations

```go
// List triggers
triggers, err := client.GetTriggers(ctx, namespace, repo)

// Get trigger
trigger, err := client.GetTrigger(ctx, namespace, repo, triggerUUID)

// Activate trigger
trigger, err := client.ActivateTrigger(ctx, namespace, repo, triggerUUID, &lib.ActivateTriggerRequest{...})

// Start build from trigger
build, err := client.StartTriggerBuild(ctx, namespace, repo, triggerUUID, &lib.ManualTriggerRequest{...})

// Enable/disable trigger
trigger, err := client.UpdateTrigger(ctx, namespace, repo, triggerUUID, true)
trigger, err := client.UpdateTrigger(ctx, namespace, repo, triggerUUID, false)

// List builds from a specific trigger
builds, err := client.GetTriggerBuilds(ctx, namespace, repo, triggerUUID, limit)

// Delete trigger
err := client.DeleteTrigger(ctx, namespace, repo, triggerUUID)
```

### Notification Operations

```go
// List notifications
notifications, err := client.GetNotifications(ctx, namespace, repo)

// Get notification
notification, err := client.GetNotification(ctx, namespace, repo, uuid)

// Create notification
notification, err := client.CreateNotification(ctx, namespace, repo, &lib.CreateNotificationRequest{
    Event:  "repo_push",
    Method: "webhook",
    Title:  "My Notification",
    Config: map[string]interface{}{"url": "https://example.com/webhook"},
})

// Update notification
notification, err := client.UpdateNotification(ctx, namespace, repo, uuid, request)

// Test notification
err := client.TestNotification(ctx, namespace, repo, uuid)

// Reset failure count
err := client.ResetNotification(ctx, namespace, repo, uuid)

// Delete notification
err := client.DeleteNotification(ctx, namespace, repo, uuid)
```

### Quota Operations

```go
// Get quota
quota, err := client.GetQuota(ctx, orgname)

// Create quota
quota, err := client.CreateQuota(ctx, orgname, limitBytes)

// Update quota
quota, err := client.UpdateQuota(ctx, orgname, limitBytes)

// Delete quota
err := client.DeleteQuota(ctx, orgname)
```

### Auto-Prune Operations

```go
// Get policies
policies, err := client.GetAutoPrunePolicies(ctx, orgname)

// Create policy
policy, err := client.CreateAutoPrunePolicy(ctx, orgname, method, value, tagPattern)

// Get specific policy
policy, err := client.GetAutoPrunePolicy(ctx, orgname, policyUUID)

// Update policy
policy, err := client.UpdateAutoPrunePolicy(ctx, orgname, policyUUID, method, value, tagPattern)

// Delete policy
err := client.DeleteAutoPrunePolicy(ctx, orgname, policyUUID)
```

### Logs Operations

```go
// Repository logs
logs, err := client.GetLogs(ctx, namespace, repo, nextPage, startDate, endDate)
aggLogs, err := client.GetAggregatedLogs(ctx, namespace, repo, startDate, endDate)

// Organization logs
logs, err := client.GetOrganizationLogs(ctx, orgname, nextPage, startDate, endDate)
aggLogs, err := client.GetOrganizationAggregatedLogs(ctx, orgname, startDate, endDate)

// User logs
logs, err := client.GetUserLogs(ctx, nextPage, startDate, endDate)
aggLogs, err := client.GetUserAggregatedLogs(ctx, startDate, endDate)

// Export logs
err := client.ExportRepositoryLogs(ctx, namespace, repo, &lib.ExportLogsRequest{...})
err := client.ExportOrganizationLogs(ctx, orgname, &lib.ExportLogsRequest{...})
err := client.ExportUserLogs(ctx, &lib.ExportLogsRequest{...})
```

### Search Operations

```go
// Search repositories
results, err := client.SearchRepositories(ctx, query, page)

// Search all (repos, users, teams, etc.)
results, err := client.SearchAll(ctx, query)
```

### Billing Operations

```go
// Plans
plans, err := client.GetAvailablePlans(ctx)

// User billing
billing, err := client.GetUserBilling(ctx)
subscription, err := client.GetUserSubscription(ctx)

// User invoices are not available in the Quay API.
// GetUserInvoices() always returns an error.

// Organization billing
billing, err := client.GetOrganizationBilling(ctx, orgname)
subscription, err := client.GetOrganizationSubscription(ctx, orgname)
invoices, err := client.GetOrganizationInvoices(ctx, orgname)
```

### Discovery Operations

```go
// Get API discovery information
discovery, err := client.GetDiscovery(ctx)

// Get registry capabilities
capabilities, err := client.GetRegistryCapabilities(ctx)

// Get application info by client ID
app, err := client.GetAppInfo(ctx, clientID)

// Search entities by prefix
entities, err := client.GetEntities(ctx, prefix, includeOrgs, includeTeams)
```

### Messages Operations

```go
// Get system messages
messages, err := client.GetMessages(ctx)

// Create a message
message, err := client.CreateMessage(ctx, content, severity, mediaType)
```

### Prototype Operations

```go
// List default permission prototypes
prototypes, err := client.GetPrototypes(ctx, orgname)

// Create prototype
prototype, err := client.CreatePrototype(ctx, orgname, &lib.CreatePrototypeRequest{...})

// Get specific prototype
prototype, err := client.GetPrototype(ctx, orgname, prototypeUUID)

// Update prototype
prototype, err := client.UpdatePrototype(ctx, orgname, prototypeUUID, &lib.UpdatePrototypeRequest{...})

// Delete prototype
err := client.DeletePrototype(ctx, orgname, prototypeUUID)
```

### RepoToken Operations (Deprecated)

```go
// List repository tokens
tokens, err := client.GetRepoTokens(ctx, namespace, repo)

// Create token
token, err := client.CreateRepoToken(ctx, namespace, repo, &lib.CreateRepoTokenRequest{...})

// Get specific token
token, err := client.GetRepoToken(ctx, namespace, repo, code)

// Update token
token, err := client.UpdateRepoToken(ctx, namespace, repo, code, &lib.UpdateRepoTokenRequest{...})

// Delete token
err := client.DeleteRepoToken(ctx, namespace, repo, code)
```

### Application Operations

```go
// List applications
apps, err := client.GetApplications(ctx, orgname)

// Create application
app, err := client.CreateApplication(ctx, orgname, name, description, applicationURI, redirectURI)

// Get application
app, err := client.GetApplication(ctx, orgname, clientID)

// Update application
app, err := client.UpdateApplication(ctx, orgname, clientID, name, description, applicationURI, redirectURI)

// Delete application
err := client.DeleteApplication(ctx, orgname, clientID)

// Reset client secret
app, err := client.ResetApplicationClientSecret(ctx, orgname, clientID)
```

### Marketplace Operations

```go
// Get organization marketplace info
marketplace, err := client.GetOrganizationMarketplace(ctx, orgname)

// Create subscription
err := client.CreateOrganizationMarketplaceSubscription(ctx, orgname, &lib.MarketplaceSubscriptionRequest{...})

// Delete subscription
err := client.DeleteOrganizationMarketplaceSubscription(ctx, orgname, subscriptionID)

// Batch remove subscriptions
err := client.BatchRemoveOrganizationMarketplaceSubscriptions(ctx, orgname, subscriptionIDs)
```

### Proxy Cache Operations

```go
// Get proxy cache config
config, err := client.GetProxyCacheConfig(ctx, orgname)

// Create proxy cache config
config, err := client.CreateProxyCacheConfig(ctx, orgname, upstreamRegistry, insecure, expiration)

// Delete proxy cache config
err := client.DeleteProxyCacheConfig(ctx, orgname)
```

### Error Operations

```go
// Get error type details
details, err := client.GetErrorType(ctx, "invalid_token")
```

### Repository Tag Listing

```go
// List tags for a repository
tags, err := client.ListTags(ctx, namespace, repo, limit, onlyActive)
```

### Mirror Operations

```go
config, err := client.GetMirrorConfig(ctx, namespace, repo)
config, err := client.CreateMirrorConfig(ctx, namespace, repo, &lib.CreateMirrorConfigRequest{
    ExternalRef:   "docker.io/library/nginx",
    SyncInterval:  86400,
    RobotUsername: "myorg+mirrorbot",
})
config, err := client.UpdateMirrorConfig(ctx, namespace, repo, &lib.UpdateMirrorConfigRequest{
    ExternalRef: "docker.io/library/nginx",
})
```

## Error Handling

API errors that include a Quay JSON body are returned as `*lib.QuayError`:

```go
import "errors"

result, err := client.GetRepository(ctx, "namespace", "repo")
if err != nil {
    var qerr *lib.QuayError
    if errors.As(err, &qerr) {
        switch qerr.StatusCode() {
        case 404:
            // Not found
        case 401:
            // Authentication failed
        case 403:
            // Permission denied
        case 429:
            // Rate limited — set client.Retry or back off
        }
        return
    }
    // Network or other error
}
```

## Rate Limiting

Quay.io rate-limits requests. Enable built-in retries instead of wrapping every call:

```go
import "time"

client.Retry = &lib.RetryConfig{
    MaxRetries:     3,
    InitialBackoff: time.Second,
    MaxBackoff:     30 * time.Second,
}
```

Retries apply to HTTP 429 and 5xx. `NewClient` leaves `Retry` nil until you set it.

## Best Practices

1. **Environment Variables** - Store tokens in environment variables, not code
2. **Error Handling** - Always check errors and handle appropriately
3. **Rate Limiting** - Set `client.Retry` for 429/5xx backoff
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
