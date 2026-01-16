# Repository Management with go-quay

This tutorial covers how to manage container image repositories using the go-quay library.

## Prerequisites

- Completed [Getting Started](./01-getting-started.md)
- Appropriate permissions in your namespace/organization

## Creating a Repository

Create a new container repository:

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

    // Create a private repository
    repo, err := client.CreateRepository(
        "my-namespace",          // namespace (org or username)
        "my-new-app",           // repository name
        "private",              // visibility: "private" or "public"
        "My application image", // description
    )
    if err != nil {
        log.Fatalf("Failed to create repository: %v", err)
    }

    fmt.Printf("Created repository: %s/%s\n", repo.Namespace, repo.Name)
}
```

## Listing Repositories

List all repositories you have access to:

```go
// List repositories with filters
repos, err := client.ListRepositories(
    "my-namespace", // namespace (empty for all)
    false,          // include public repos
    false,          // only starred repos
    1,              // page number
    25,             // results per page
)
if err != nil {
    log.Fatalf("Failed to list: %v", err)
}

fmt.Printf("Found %d repositories\n", len(repos.Repositories))
for _, r := range repos.Repositories {
    fmt.Printf("- %s/%s\n", r.Namespace, r.Name)
}
```

## Getting Repository Details

Fetch complete repository information including tags:

```go
repo, err := client.GetRepository("my-namespace", "my-app")
if err != nil {
    log.Fatalf("Failed to get repository: %v", err)
}

// Repository metadata
fmt.Printf("Name: %s\n", repo.Name)
fmt.Printf("Namespace: %s\n", repo.Namespace)
fmt.Printf("Description: %s\n", repo.Description)
fmt.Printf("Public: %v\n", repo.IsPublic)
fmt.Printf("State: %s\n", repo.State)

// Access permissions
fmt.Printf("Can Write: %v\n", repo.CanWrite)
fmt.Printf("Can Admin: %v\n", repo.CanAdmin)

// Tags
fmt.Printf("\nTags (%d):\n", len(repo.Tags.Tags))
for _, tag := range repo.Tags.Tags {
    fmt.Printf("  %s\n", tag.Name)
    fmt.Printf("    Digest: %s\n", tag.ManifestDigest[:16]+"...")
    fmt.Printf("    Modified: %s\n", tag.LastModified)
    if tag.Expiration != "" {
        fmt.Printf("    Expires: %s\n", tag.Expiration)
    }
}
```

## Updating a Repository

Modify repository settings:

```go
// Update description and visibility
updatedRepo, err := client.UpdateRepository(
    "my-namespace",
    "my-app",
    "Updated description for my app", // new description
    "public",                          // new visibility
)
if err != nil {
    log.Fatalf("Failed to update: %v", err)
}

fmt.Printf("Updated: %s/%s\n", updatedRepo.Namespace, updatedRepo.Name)
```

You can also change just the visibility:

```go
err := client.ChangeRepositoryVisibility("my-namespace", "my-app", "private")
if err != nil {
    log.Fatalf("Failed to change visibility: %v", err)
}
fmt.Println("Repository is now private")
```

## Managing Tags

### Get Tag Details

```go
tag, err := client.GetTag("my-namespace", "my-app", "v1.0.0")
if err != nil {
    log.Fatalf("Failed to get tag: %v", err)
}

fmt.Printf("Tag: %s\n", tag.Name)
fmt.Printf("Digest: %s\n", tag.ManifestDigest)
```

### Update Tag (Set Expiration)

```go
// Set tag to expire in 30 days
updatedTag, err := client.UpdateTag(
    "my-namespace",
    "my-app",
    "latest",
    "2024-12-31T23:59:59Z", // expiration timestamp
)
if err != nil {
    log.Fatalf("Failed to update tag: %v", err)
}

fmt.Printf("Tag %s will expire at %s\n", updatedTag.Name, updatedTag.Expiration)
```

### View Tag History

```go
history, err := client.GetTagHistory("my-namespace", "my-app", "latest")
if err != nil {
    log.Fatalf("Failed to get history: %v", err)
}

fmt.Printf("History for tag 'latest':\n")
for _, entry := range history.Tags {
    fmt.Printf("  %s -> %s\n", entry.StartTs, entry.ManifestDigest[:16])
}
```

### Delete a Tag

```go
err := client.DeleteTag("my-namespace", "my-app", "old-version")
if err != nil {
    log.Fatalf("Failed to delete tag: %v", err)
}
fmt.Println("Tag deleted")
```

### Restore a Tag

Restore a tag to a previous image:

```go
err := client.RestoreTag(
    "my-namespace",
    "my-app",
    "latest",
    "sha256:abc123...", // manifest digest to restore to
)
if err != nil {
    log.Fatalf("Failed to restore tag: %v", err)
}
fmt.Println("Tag restored")
```

## Deleting a Repository

**Warning:** This operation is irreversible!

```go
err := client.DeleteRepository("my-namespace", "my-app")
if err != nil {
    log.Fatalf("Failed to delete: %v", err)
}
fmt.Println("Repository deleted")
```

## Starring Repositories

Star a repository for easy access:

```go
// Star a repository
err := client.StarRepository("quay", "quay")
if err != nil {
    log.Fatalf("Failed to star: %v", err)
}

// Get starred repositories
starred, err := client.GetStarredRepositories()
if err != nil {
    log.Fatalf("Failed to get starred: %v", err)
}

fmt.Printf("Starred repositories: %d\n", len(starred.Repositories))
for _, r := range starred.Repositories {
    fmt.Printf("- %s/%s\n", r.Namespace, r.Name)
}

// Unstar a repository
err = client.UnstarRepository("quay", "quay")
```

## Complete Workflow Example

Here's a complete example managing a repository lifecycle:

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
    namespace := "my-org"
    repoName := "demo-app"

    // 1. Create repository
    fmt.Println("Creating repository...")
    repo, err := client.CreateRepository(namespace, repoName, "private", "Demo application")
    if err != nil {
        log.Fatalf("Create failed: %v", err)
    }
    fmt.Printf("Created: %s/%s\n\n", repo.Namespace, repo.Name)

    // 2. Update description
    fmt.Println("Updating repository...")
    _, err = client.UpdateRepository(namespace, repoName, "Production demo application", "")
    if err != nil {
        log.Fatalf("Update failed: %v", err)
    }
    fmt.Println("Updated description\n")

    // 3. Star the repository
    fmt.Println("Starring repository...")
    err = client.StarRepository(namespace, repoName)
    if err != nil {
        log.Printf("Star failed (might not be supported): %v\n", err)
    } else {
        fmt.Println("Repository starred\n")
    }

    // 4. Get final state
    fmt.Println("Final repository state:")
    finalRepo, _ := client.GetRepository(namespace, repoName)
    fmt.Printf("  Name: %s\n", finalRepo.Name)
    fmt.Printf("  Description: %s\n", finalRepo.Description)
    fmt.Printf("  Starred: %v\n", finalRepo.IsStarred)
}
```

## Next Steps

- [Security Scanning](./03-security-scanning.md) - Scan images for vulnerabilities
- [CI/CD Automation](./04-ci-cd-automation.md) - Automate with robot accounts
