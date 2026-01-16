# Organization Administration with go-quay

This tutorial covers organization management tasks: teams, members, quotas, policies, and more.

## Prerequisites

- Completed [Getting Started](./01-getting-started.md)
- Admin access to a Quay.io organization

## Organization Overview

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    client, _ := lib.NewClient(os.Getenv("QUAY_TOKEN"))
    orgName := "my-org"

    org, err := client.GetOrganization(orgName)
    if err != nil {
        log.Fatalf("Failed to get org: %v", err)
    }

    fmt.Printf("Organization: %s\n", org.Name)
    fmt.Printf("Email: %s\n", org.Email)
    fmt.Printf("Is Admin: %v\n", org.IsAdmin)
    fmt.Printf("Is Member: %v\n", org.IsMember)
}
```

## Team Management

Teams group users together for easier permission management.

### Creating a Team

```go
team, err := client.CreateTeam(
    "my-org",
    "developers",           // team name
    "Development team",     // description
    "member",              // role: member, creator, or admin
)
if err != nil {
    log.Fatalf("Failed to create team: %v", err)
}

fmt.Printf("Created team: %s\n", team.Name)
```

### Team Roles

| Role | Capabilities |
|------|--------------|
| `member` | Access to assigned repositories |
| `creator` | Can create new repositories |
| `admin` | Full organization admin access |

### Listing Teams

```go
teams, err := client.GetOrganizationTeams("my-org")
if err != nil {
    log.Fatalf("Failed to list teams: %v", err)
}

for _, team := range teams.Teams {
    fmt.Printf("- %s (role: %s, members: %d)\n",
        team.Name, team.Role, team.MemberCount)
}
```

### Updating a Team

```go
updatedTeam, err := client.UpdateTeam(
    "my-org",
    "developers",
    "Updated description",
    "creator",  // promote to creator role
)
```

### Deleting a Team

```go
err := client.DeleteTeam("my-org", "developers")
if err != nil {
    log.Fatalf("Failed to delete team: %v", err)
}
```

## Member Management

### Adding Members to a Team

```go
err := client.AddTeamMember("my-org", "developers", "john.doe")
if err != nil {
    log.Fatalf("Failed to add member: %v", err)
}
fmt.Println("Member added!")
```

### Listing Team Members

```go
members, err := client.GetTeamMembers("my-org", "developers")
if err != nil {
    log.Fatalf("Failed to get members: %v", err)
}

for _, m := range members.Members {
    fmt.Printf("- %s (%s)\n", m.Name, m.Kind)
}
```

### Removing Members

```go
err := client.RemoveTeamMember("my-org", "developers", "john.doe")
if err != nil {
    log.Fatalf("Failed to remove member: %v", err)
}
```

### Organization Members

List all organization members:

```go
members, err := client.GetOrganizationMembers("my-org")
if err != nil {
    log.Fatalf("Failed to list members: %v", err)
}

fmt.Printf("Total members: %d\n", len(members.Members))
for _, m := range members.Members {
    fmt.Printf("- %s\n", m.Name)
}
```

## Team Repository Permissions

Grant teams access to repositories:

```go
// Set team permission for a repository
err := client.SetTeamRepositoryPermission(
    "my-org",
    "developers",   // team name
    "my-app",       // repository
    "write",        // role: read, write, admin
)
if err != nil {
    log.Fatalf("Failed to set permission: %v", err)
}

// List team permissions
perms, err := client.GetTeamPermissions("my-org", "developers")
if err != nil {
    log.Fatalf("Failed to get permissions: %v", err)
}

for _, p := range perms.Permissions {
    fmt.Printf("- %s: %s\n", p.Repository.Name, p.Role)
}

// Remove team permission
err = client.RemoveTeamRepositoryPermission("my-org", "developers", "my-app")
```

## Quota Management

Control storage usage in your organization:

### Getting Quota Information

```go
quotas, err := client.GetOrganizationQuota("my-org")
if err != nil {
    log.Fatalf("Failed to get quota: %v", err)
}

for _, q := range quotas.Quotas {
    usedGB := float64(q.UsedBytes) / (1024 * 1024 * 1024)
    limitGB := float64(q.LimitBytes) / (1024 * 1024 * 1024)
    percent := float64(q.UsedBytes) / float64(q.LimitBytes) * 100

    fmt.Printf("Used: %.2f GB / %.2f GB (%.1f%%)\n", usedGB, limitGB, percent)
}
```

### Setting Quota Limits

```go
// Set 10 GB limit (requires super user)
limitBytes := int64(10 * 1024 * 1024 * 1024)
quota, err := client.CreateOrganizationQuota("my-org", limitBytes)
if err != nil {
    log.Fatalf("Failed to create quota: %v", err)
}

fmt.Printf("Quota set: %d bytes\n", quota.LimitBytes)
```

### Updating Quota

```go
// Increase to 50 GB
newLimit := int64(50 * 1024 * 1024 * 1024)
quota, err := client.UpdateOrganizationQuota("my-org", quotaID, newLimit)
```

## Auto-Prune Policies

Automatically clean up old tags to save storage:

### Creating an Auto-Prune Policy

```go
// Keep only the last 10 tags
policy, err := client.CreateOrganizationAutoPrunePolicy(
    "my-org",
    "number_of_tags",   // method
    10,                 // value: keep 10 tags
    "",                 // tag pattern (optional)
)
if err != nil {
    log.Fatalf("Failed to create policy: %v", err)
}

fmt.Printf("Policy created: %s\n", policy.UUID)
```

### Prune Methods

| Method | Value | Description |
|--------|-------|-------------|
| `number_of_tags` | Integer | Keep N most recent tags |
| `creation_date` | Days | Keep tags newer than N days |

### With Tag Pattern

```go
// Only prune tags matching pattern
policy, err := client.CreateOrganizationAutoPrunePolicy(
    "my-org",
    "creation_date",
    30,                 // keep tags from last 30 days
    "^v[0-9]+\\.",      // only match version tags
)
```

### Listing Policies

```go
policies, err := client.GetOrganizationAutoPrunePolicies("my-org")
if err != nil {
    log.Fatalf("Failed to get policies: %v", err)
}

for _, p := range policies.Policies {
    fmt.Printf("- %s: %s = %v\n", p.UUID[:8], p.Method, p.Value)
    if p.TagPattern != "" {
        fmt.Printf("  Pattern: %s\n", p.TagPattern)
    }
}
```

### Updating a Policy

```go
updatedPolicy, err := client.UpdateOrganizationAutoPrunePolicy(
    "my-org",
    policyUUID,
    "number_of_tags",
    20,  // increase to keep 20 tags
    "",
)
```

### Deleting a Policy

```go
err := client.DeleteOrganizationAutoPrunePolicy("my-org", policyUUID)
```

## Default Permissions (Prototypes)

Automatically apply permissions to new repositories:

### Creating a Prototype

```go
prototype, err := client.CreateOrganizationPrototype(
    "my-org",
    "write",              // role
    "team",               // delegate kind: user, team, or robot
    "developers",         // delegate name
    nil,                  // activating user (optional)
)
if err != nil {
    log.Fatalf("Failed to create prototype: %v", err)
}

fmt.Printf("Prototype created: %s\n", prototype.ID)
```

### Listing Prototypes

```go
prototypes, err := client.GetOrganizationPrototypes("my-org")
if err != nil {
    log.Fatalf("Failed to get prototypes: %v", err)
}

for _, p := range prototypes.Prototypes {
    fmt.Printf("- %s: %s (%s) -> %s\n",
        p.ID[:8], p.Delegate.Name, p.Delegate.Kind, p.Role)
}
```

### Deleting a Prototype

```go
err := client.DeleteOrganizationPrototype("my-org", prototypeID)
```

## OAuth Applications

Manage OAuth applications for your organization:

### Creating an Application

```go
app, err := client.CreateOrganizationApplication(
    "my-org",
    "My CI Tool",                           // name
    "CI/CD integration app",                // description
    "https://ci.example.com",               // application URI
    "https://ci.example.com/oauth/callback", // redirect URI
)
if err != nil {
    log.Fatalf("Failed to create app: %v", err)
}

fmt.Printf("App created!\n")
fmt.Printf("Client ID: %s\n", app.ClientID)
fmt.Printf("Client Secret: %s\n", app.ClientSecret)
```

### Listing Applications

```go
apps, err := client.GetOrganizationApplications("my-org")
if err != nil {
    log.Fatalf("Failed to list apps: %v", err)
}

for _, app := range apps.Applications {
    fmt.Printf("- %s (ID: %s)\n", app.Name, app.ClientID)
}
```

### Resetting Client Secret

```go
app, err := client.ResetOrganizationApplicationClientSecret("my-org", clientID)
if err != nil {
    log.Fatalf("Failed to reset secret: %v", err)
}

fmt.Printf("New secret: %s\n", app.ClientSecret)
```

## Audit Logging

Access organization activity logs:

```go
logs, err := client.GetOrganizationLogs("my-org", "")
if err != nil {
    log.Fatalf("Failed to get logs: %v", err)
}

fmt.Printf("Log entries: %d\n", len(logs.Logs))
for _, entry := range logs.Logs {
    fmt.Printf("[%s] %s: %s\n",
        entry.Datetime, entry.Performer.Name, entry.Kind)
}

// Pagination
if logs.NextPage != "" {
    nextLogs, err := client.GetOrganizationLogs("my-org", logs.NextPage)
    // ... process next page
}
```

### Aggregated Logs

Get log statistics:

```go
aggLogs, err := client.GetOrganizationAggregatedLogs(
    "my-org",
    "01/01/2024",  // start date
    "01/31/2024",  // end date
)
if err != nil {
    log.Fatalf("Failed to get aggregated logs: %v", err)
}

for _, entry := range aggLogs.Aggregated {
    fmt.Printf("%s: %d\n", entry.Kind, entry.Count)
}
```

## Complete Example

See the [organization-management example](../../examples/organization-management/main.go) for a complete working program.

## Best Practices

1. **Use Teams** - Group users by function (dev, ops, security)
2. **Least Privilege** - Grant minimum necessary permissions
3. **Auto-Prune** - Set up policies to manage storage
4. **Audit Regularly** - Review logs and permissions
5. **Rotate Secrets** - Regenerate app secrets periodically
6. **Document Prototypes** - Document default permission settings
