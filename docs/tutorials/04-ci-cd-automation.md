# CI/CD Automation with go-quay

This tutorial covers how to use go-quay for CI/CD automation, including robot accounts, permissions, webhooks, and build triggers.

## Prerequisites

- Completed [Getting Started](./01-getting-started.md)
- Admin access to a Quay.io organization
- Understanding of CI/CD concepts

## Robot Accounts

Robot accounts are service accounts designed for automation. They provide:
- Non-expiring tokens
- Granular permissions
- Easy revocation
- Audit trails

### Creating a Robot Account

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    client, _ := lib.NewClient(os.Getenv("QUAY_TOKEN"))
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    orgName := "my-org"

    // Create a robot account
    robot, err := client.CreateRobotAccount(ctx, 
        orgName,
        "ci-deploy",                    // short name (becomes org+ci-deploy)
        "CI/CD deployment automation",  // description
        nil,                            // unstructured metadata
    )
    if err != nil {
        log.Fatalf("Failed to create robot: %v", err)
    }

    fmt.Printf("Robot created: %s\n", robot.Name)
    fmt.Printf("Token (save this!): %s\n", robot.Token)
    fmt.Println("\nIMPORTANT: The token is only shown once!")
}
```

### Listing Robot Accounts

```go
robots, err := client.GetRobotAccounts(ctx, "my-org")
if err != nil {
    log.Fatalf("Failed to list robots: %v", err)
}

fmt.Printf("Found %d robots:\n", len(robots.Robots))
for _, robot := range robots.Robots {
    fmt.Printf("- %s: %s\n", robot.Name, robot.Description)
}
```

### Regenerating Robot Token

If a token is compromised:

```go
newRobot, err := client.RegenerateRobotToken(ctx, "my-org", "ci-deploy")
if err != nil {
    log.Fatalf("Failed to regenerate: %v", err)
}

fmt.Printf("New token: %s\n", newRobot.Token)
fmt.Println("The old token is now invalid!")
```

### User-Level Robots

You can also create robots at the user level:

```go
// Create user robot
robot, err := client.CreateUserRobotAccount(ctx, "my-deploy-bot", "Personal deployment bot", nil)

// List user robots
robots, err := client.GetUserRobotAccounts(ctx)

// Regenerate user robot token
newRobot, err := client.RegenerateUserRobotToken(ctx, "my-deploy-bot")
```

## Setting Repository Permissions

Grant robots access to specific repositories:

```go
// Grant write access to a repository
err := client.SetUserPermission(ctx, 
    "my-org",                    // namespace
    "my-app",                    // repository
    "my-org+ci-deploy",          // robot username
    "write",                     // role: read, write, or admin
)
if err != nil {
    log.Fatalf("Failed to set permission: %v", err)
}

fmt.Println("Permission granted!")

// Verify the permission
perm, err := client.GetUserPermission(ctx, "my-org", "my-app", "my-org+ci-deploy")
if err != nil {
    log.Fatalf("Failed to verify: %v", err)
}

fmt.Printf("User %s has role: %s\n", perm.Name, perm.Role)
```

### Permission Roles

| Role | Capabilities |
|------|--------------|
| `read` | Pull images |
| `write` | Pull and push images |
| `admin` | Full access including settings |

### Team Permissions

You can also grant permissions to entire teams:

```go
// Set team permission
err := client.SetTeamPermission(ctx, 
    "my-org",
    "my-app",
    "developers",  // team name
    "write",
)

// List team permissions
perms, err := client.ListTeamPermissions(ctx, "my-org", "my-app")
```

## Webhook Notifications

Set up webhooks to notify external services when events occur:

### Creating a Webhook

```go
notification, err := client.CreateNotification(ctx, 
    "my-org",
    "my-app",
    &lib.CreateNotificationRequest{
        Event:  "repo_push",          // trigger on image push
        Method: "webhook",            // notification method
        Title:  "Image Push Alert",   // friendly name
        Config: map[string]interface{}{
            "url": "https://my-ci-server.com/webhook",
        },
    },
)
if err != nil {
    log.Fatalf("Failed to create notification: %v", err)
}

fmt.Printf("Notification created: %s\n", notification.UUID)
```

### Supported Events

| Event | Description |
|-------|-------------|
| `repo_push` | Image pushed to repository |
| `build_queued` | Build has been queued |
| `build_start` | Build has started |
| `build_success` | Build completed successfully |
| `build_failure` | Build failed |
| `build_canceled` | Build was canceled |
| `vulnerability_found` | New vulnerability discovered |

### Supported Methods

| Method | Description |
|--------|-------------|
| `webhook` | HTTP POST to a URL |
| `slack` | Slack webhook |
| `email` | Email notification |
| `hipchat` | HipChat integration |
| `flowdock` | Flowdock integration |

### Slack Integration

```go
notification, err := client.CreateNotification(ctx, 
    "my-org",
    "my-app",
    &lib.CreateNotificationRequest{
        Event:  "vulnerability_found",
        Method: "slack",
        Title:  "Security Alert",
        Config: map[string]interface{}{
            "url": "https://hooks.slack.com/services/XXX/YYY/ZZZ",
        },
    },
)
```

### Testing Notifications

```go
// Test a notification
err := client.TestNotification(ctx, "my-org", "my-app", notificationUUID)
if err != nil {
    log.Printf("Test failed: %v\n", err)
} else {
    fmt.Println("Test notification sent!")
}
```

### Managing Notifications

```go
// List all notifications
notifications, err := client.GetNotifications(ctx, "my-org", "my-app")
for _, n := range notifications.Notifications {
    fmt.Printf("- %s: %s (%s)\n", n.UUID[:8], n.Event, n.Method)
}

// Reset failure count
err = client.ResetNotification(ctx, "my-org", "my-app", notificationUUID)

// Delete a notification
err = client.DeleteNotification(ctx, "my-org", "my-app", notificationUUID)
```

## Build Triggers

Build triggers automatically build images when code is pushed to a repository.

### Listing Triggers

```go
triggers, err := client.GetTriggers(ctx, "my-org", "my-app")
if err != nil {
    log.Fatalf("Failed to list triggers: %v", err)
}

for _, t := range triggers.Triggers {
    status := "enabled"
    if !t.Enabled {
        status = "disabled"
    }
    fmt.Printf("- %s: %s [%s]\n", t.ID, t.Service, status)
}
```

### Getting Trigger Details

```go
trigger, err := client.GetTrigger(ctx, "my-org", "my-app", triggerUUID)
if err != nil {
    log.Fatalf("Failed to get trigger: %v", err)
}

fmt.Printf("Trigger: %s\n", trigger.ID)
fmt.Printf("Service: %s\n", trigger.Service)
fmt.Printf("Enabled: %v\n", trigger.Enabled)
```

### Starting a Build Manually

```go
// Start a build from a trigger
build, err := client.StartTriggerBuild(ctx, "my-org", "my-app", triggerUUID, nil)
if err != nil {
    log.Fatalf("Failed to start build: %v", err)
}

fmt.Printf("Build started: %s\n", build.ID)
fmt.Printf("Status: %s\n", build.Phase)
```

### Activating a Trigger

```go
// Activate a trigger with configuration
trigger, err := client.ActivateTrigger(ctx, 
    "my-org",
    "my-app",
    triggerUUID,
    &lib.ActivateTriggerRequest{
        Config: map[string]interface{}{
            "build_source":    "master",
            "dockerfile_path": "/Dockerfile",
        },
    },
)
if err != nil {
    log.Fatalf("Failed to activate trigger: %v", err)
}
fmt.Printf("Trigger activated: %s\n", trigger.ID)
```

## Complete CI/CD Setup

Typical order: create the repo, create a robot, `SetUserPermission` for `org+robot`, then `CreateNotification` for `repo_push` / `vulnerability_found`. The runnable program is [ci-cd-integration](../../examples/ci-cd-integration/main.go).

```go
robot, err := client.CreateRobotAccount(ctx, org, robotName, "CI/CD builder", nil)
err = client.SetUserPermission(ctx, org, repo, org+"+"+robotName, "write")
notification, err := client.CreateNotification(ctx, org, repo, &lib.CreateNotificationRequest{
    Event:  "repo_push",
    Method: "webhook",
    Title:  "CI Build Trigger",
    Config: map[string]interface{}{"url": webhookURL},
})
```

## Best Practices

1. **Use Robot Accounts** - Never use personal tokens in CI/CD
2. **Minimum Permissions** - Grant only the permissions needed
3. **Rotate Tokens** - Regenerate tokens periodically
4. **Monitor Notifications** - Check for failed webhooks
5. **Use Private Repos** - Keep images private by default
6. **Tag Strategically** - Use semantic versioning

## Complete Example

See the [ci-cd-integration example](../../examples/ci-cd-integration/main.go) for a complete working program.

## Next Steps

- [Organization Administration](./05-organization-admin.md) - Manage teams and quotas
