# Getting Started with go-quay

This tutorial walks you through the basics of using the go-quay library to interact with the Quay.io container registry API.

## Prerequisites

- Go 1.26+
- A Quay.io account with API access
- An API token from Quay.io

## Step 1: Install the Library

Add go-quay to your Go project:

```bash
go get github.com/sebrandon1/go-quay
```

## Step 2: Obtain an API Token

1. Log in to [Quay.io](https://quay.io)
2. Open **Account Settings** → **Applications** (or create a **Robot Account** under your user or organization)
3. Create a token with the permissions you need:
   - **Read repositories** - for listing and pulling
   - **Write to repositories** - for pushing images
   - **Administer organization** - for managing teams and settings

Do not use **CLI Password** / **Generate Encrypted Password** — those are for `docker login` only.

Alternatively, create a **Robot Account** for automation:
1. Go to your organization or user settings
2. Click **Robot Accounts**
3. Create a new robot with the required permissions

## Step 3: Initialize the Client

Create a new Go file and initialize the client:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    // Get token from environment variable (recommended)
    token := os.Getenv("QUAY_TOKEN")
    if token == "" {
        log.Fatal("QUAY_TOKEN environment variable is required")
    }

    // Initialize the client
    client, err := lib.NewClient(token)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    fmt.Println("Client initialized successfully!")
}
```

Run your program:

```bash
export QUAY_TOKEN="your-token-here"
go run main.go
```

## Step 4: Make Your First API Call

Let's get information about the current user:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    token := os.Getenv("QUAY_TOKEN")
    if token == "" {
        log.Fatal("QUAY_TOKEN environment variable is required")
    }

    client, err := lib.NewClient(token)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    ctx := context.Background()

    // Get current user information
    user, err := client.GetUser(ctx)
    if err != nil {
        log.Fatalf("Failed to get user: %v", err)
    }

    fmt.Printf("Logged in as: %s\n", user.Username)
    fmt.Printf("Email: %s\n", user.Email)
    fmt.Printf("Verified: %v\n", user.Verified)
}
```

For a self-hosted registry, use `lib.NewClientWithURL(token, "https://quay.example.com/api/v1")`.

## Step 5: List Repositories

Now let's list repositories in a namespace:

```go
// List repositories in your namespace
repos, err := client.ListRepositories(ctx, "your-username", false, false, false, 1, 10)
if err != nil {
    log.Fatalf("Failed to list repositories: %v", err)
}

fmt.Printf("Found %d repositories:\n", len(repos.Repositories))
for _, repo := range repos.Repositories {
    visibility := "private"
    if repo.IsPublic {
        visibility = "public"
    }
    fmt.Printf("- %s/%s (%s)\n", repo.Namespace, repo.Name, visibility)
}
```

## Step 6: Get Repository Details

Get detailed information about a specific repository:

```go
// Get repository with tags
repo, err := client.GetRepository(ctx, "namespace", "repository-name")
if err != nil {
    log.Fatalf("Failed to get repository: %v", err)
}

fmt.Printf("Repository: %s/%s\n", repo.Namespace, repo.Name)
fmt.Printf("Description: %s\n", repo.Description)
fmt.Printf("Tags: %d\n", len(repo.Tags.Tags))

// List tags
for _, tag := range repo.Tags.Tags {
    fmt.Printf("  - %s (digest: %s)\n", tag.Name, tag.ManifestDigest[:12])
}
```

## Error Handling

The go-quay library returns `*lib.QuayError` when the API sends a JSON error body:

```go
import (
    "context"
    "errors"
    "fmt"

    "github.com/sebrandon1/go-quay/lib"
)

repo, err := client.GetRepository(ctx, "namespace", "nonexistent-repo")
if err != nil {
    var qerr *lib.QuayError
    if errors.As(err, &qerr) {
        switch qerr.StatusCode() {
        case 404:
            fmt.Println("Repository not found")
        case 401:
            fmt.Println("Authentication failed - check your token")
        case 403:
            fmt.Println("Access denied - insufficient permissions")
        default:
            fmt.Printf("API error: %v\n", err)
        }
        return
    }
    fmt.Printf("Unexpected error: %v\n", err)
    return
}
```

## Environment Variables

| Variable | Used by | Description |
|----------|---------|-------------|
| `QUAY_TOKEN` | library, CLI, examples | API token |
| `QUAY_URL` | CLI | Optional API base URL (self-hosted Quay) |
| `QUAY_NAMESPACE` | examples and integration tests only | Default namespace; the CLI uses `--namespace` or config `namespace` instead |

## Complete Example

See the [basic-usage example](../../examples/basic-usage/main.go) for a complete working program.

## Next Steps

- [Repository Management](./02-repository-management.md) - Create, update, and delete repositories
- [Security Scanning](./03-security-scanning.md) - Scan images for vulnerabilities
- [CI/CD Automation](./04-ci-cd-automation.md) - Set up robot accounts and webhooks
- [Organization Administration](./05-organization-admin.md) - Manage teams and quotas
