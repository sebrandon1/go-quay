# lib

Go client for the [Quay.io REST API](https://docs.quay.io/api/swagger/).

## Installation

```bash
go get github.com/sebrandon1/go-quay
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    client, err := lib.NewClient(os.Getenv("QUAY_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    user, err := client.GetUser()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Logged in as: %s\n", user.Username)
}
```

For a self-hosted registry, use `lib.NewClientWithURL(token, "https://quay.example.com/api/v1")`. Optional retries: set `client.Retry` to a `*lib.RetryConfig`.

## Documentation

- [Library Guide](../docs/library-guide.md) — complete API reference by domain
- [Tutorials](../docs/tutorials/01-getting-started.md) — step-by-step guides
- [Examples](../examples/) — runnable programs

## API Domains

Client methods are organized by domain in `lib/<domain>.go` (for example `repository.go`, `organization.go`, `mirror.go`). Each file's doc comment lists the HTTP endpoints it covers. Request and response types live in `structs.go`.

Coverage includes billing, builds, discovery, logs, manifests, notifications, organizations, permissions, prototypes, repositories, repository mirroring, robots, search, security scans, tags, teams, triggers, users, applications, marketplace, proxy cache, quota, auto-prune, messages, error types, and (deprecated) repository tokens.

## Authentication

All endpoints require a bearer token. Set `QUAY_TOKEN` or pass the token to `lib.NewClient(token)`. The default API base is `lib.DefaultQuayURL` (`https://quay.io/api/v1`).

## Testing

Unit tests create a client with `NewClientWithURL(token, server.URL+"/api/v1")` against `httptest.NewServer`. See any `lib/*_test.go` file for the pattern.
