# go-quay

[![Pre-Main Checks](https://github.com/sebrandon1/go-quay/actions/workflows/pre-main.yaml/badge.svg)](https://github.com/sebrandon1/go-quay/actions/workflows/pre-main.yaml)
[![Quay API Verified Nightly](https://github.com/sebrandon1/go-quay/actions/workflows/nightly.yaml/badge.svg)](https://github.com/sebrandon1/go-quay/actions/workflows/nightly.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sebrandon1/go-quay)](https://golang.org/)
[![License](https://img.shields.io/github/license/sebrandon1/go-quay)](https://github.com/sebrandon1/go-quay/blob/main/LICENSE)

A Go wrapper around the [Quay.io REST API](https://docs.quay.io/api/swagger/).
Provides both a reusable library (`lib/`) and a CLI (`cmd/`) built with Cobra.

## Key Features

- **Full API Coverage** — Billing, builds, manifests, organizations, permissions, repos, robots, security scans, tags, teams, and more
- **Library + CLI** — Use as a Go package or as a standalone command-line tool
- **Container Image** — Available at `quay.io/bapalm/go-quay` for quick use without installation
- **Popularity Sorting** — List repositories with pull-count scores (`--popularity --table` sorts by that score)

## Quick Start

### Prebuilt binary

Download from [GitHub Releases](https://github.com/sebrandon1/go-quay/releases):

```bash
# Linux (amd64)
curl -sL https://github.com/sebrandon1/go-quay/releases/latest/download/go-quay_$(curl -sL https://api.github.com/repos/sebrandon1/go-quay/releases/latest | grep tag_name | cut -d '"' -f4)_linux_amd64.tar.gz | tar xz
sudo mv go-quay /usr/local/bin/
```

### Container image

```bash
podman run --rm quay.io/bapalm/go-quay info repository \
  -n myorg -r myapp -t "$QUAY_TOKEN"
```

### Go install

```bash
go install github.com/sebrandon1/go-quay@latest
```

### Build from source

```bash
git clone https://github.com/sebrandon1/go-quay.git
cd go-quay && make build
```

## Library Usage

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
    client, err := lib.NewClient(os.Getenv("QUAY_TOKEN"))
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    user, err := client.GetUser(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Logged in as: %s\n", user.Username)

    repos, err := client.ListRepositories(ctx, "my-namespace", false, false, false, 1, 10)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d repositories\n", len(repos.Repositories))
}
```

## CLI Usage

```bash
# Get repository info
go-quay info repository -n myorg -r myapp -t "$QUAY_TOKEN"

# List repositories sorted by pull count
go-quay list repository -n myorg --popularity --table -t "$QUAY_TOKEN"

# Security scan a manifest
go-quay info secscan -n myorg -r myapp -m sha256:abc123... -t "$QUAY_TOKEN"
```

Mutating commands use the same verb-first shape (`create` / `delete` / `update`). `go-quay get …` still works; mutating actions under `get` are deprecated.

## Guides

| Guide | Description |
|-------|-------------|
| [Library Guide](docs/library-guide.md) | Complete API reference for using go-quay as a Go package |
| [CLI Reference](docs/cli-reference.md) | Full command reference for every API endpoint |
| [Getting Started](docs/tutorials/01-getting-started.md) | Installation and first API calls |
| [Repository Management](docs/tutorials/02-repository-management.md) | Create, update, and manage repos |
| [Security Scanning](docs/tutorials/03-security-scanning.md) | Scan images for vulnerabilities |
| [CI/CD Automation](docs/tutorials/04-ci-cd-automation.md) | Robot accounts and webhooks |
| [Organization Admin](docs/tutorials/05-organization-admin.md) | Manage teams and quotas |
| [Examples](examples/) | Runnable example programs |

## API Coverage

Each API links to the corresponding [Quay.io Swagger documentation](https://docs.quay.io/api/swagger/):

| API | Cmd | Lib | Covered |
| --- | --- | --- | ------- |
| [Billing](https://docs.quay.io/api/swagger/#operation--api-v1-user-plan-get) | Yes | Yes | /api/v1/user/plan, /api/v1/organization/{orgname}/plan, /api/v1/organization/{orgname}/invoices, /api/v1/plans/ |
| [Build](https://docs.quay.io/api/swagger/#Build) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/build/, /api/v1/repository/{namespace}/{repository}/build/{build_uuid}, /api/v1/repository/{namespace}/{repository}/build/{build_uuid}/logs |
| [Discovery](https://docs.quay.io/api/swagger/#Discovery) | Yes | Yes | /api/v1/discovery |
| [Error](https://docs.quay.io/api/swagger/#Error) | Yes | Yes | /api/v1/error/{error_type} |
| [Messages](https://docs.quay.io/api/swagger/#Messages) | Yes | Yes | /api/v1/messages |
| [Logs](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--aggregatelogs-get) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/aggregatelogs, /api/v1/repository/{namespace}/{repository}/logs, /api/v1/organization/{orgname}/logs, /api/v1/organization/{orgname}/aggregatelogs, /api/v1/user/logs, /api/v1/user/aggregatelogs |
| [Manifest](https://docs.quay.io/api/swagger/#Manifest) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}, /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/labels, /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/labels/{labelid} |
| Mirror | Yes | Yes | /api/v1/repository/{namespace}/{repository}/mirror |
| [Organization](https://docs.quay.io/api/swagger/#operation--api-v1-organization--orgname--get) | Yes | Yes | /api/v1/organization/{orgname}, /api/v1/organization/{orgname}/members, /api/v1/organization/{orgname}/teams, /api/v1/organization/{orgname}/team/{teamname}, /api/v1/organization/{orgname}/robots, /api/v1/organization/{orgname}/quota, /api/v1/organization/{orgname}/autoprunepolicy, /api/v1/organization/{orgname}/applications |
| [Permission](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--permissions-get) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/permissions, /api/v1/repository/{namespace}/{repository}/permissions/{username} |
| [Prototype](https://docs.quay.io/api/swagger/#Prototype) | Yes | Yes | /api/v1/organization/{orgname}/prototypes, /api/v1/organization/{orgname}/prototypes/{uuid} |
| [Repository](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--get) | Yes | Yes | /api/v1/repository/{namespace}/{repository}, /api/v1/repository/{namespace}/{repository}/tag, /api/v1/repository, /api/v1/repository/{namespace}/{repository} (CRUD) |
| [RepositoryNotification](https://docs.quay.io/api/swagger/#RepositoryNotification) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/notification/, /api/v1/repository/{namespace}/{repository}/notification/{uuid}, /api/v1/repository/{namespace}/{repository}/notification/{uuid}/test |
| [RepoToken](https://docs.quay.io/api/swagger/#RepoToken) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/tokens, /api/v1/repository/{namespace}/{repository}/tokens/{code} (DEPRECATED) |
| [Robot](https://docs.quay.io/api/swagger/#Robot) | Yes | Yes | /api/v1/user/robots, /api/v1/user/robots/{robot_shortname}, /api/v1/user/robots/{robot_shortname}/regenerate, /api/v1/user/robots/{robot_shortname}/permissions |
| [Search](https://docs.quay.io/api/swagger/#Search) | Yes | Yes | /api/v1/find/repositories, /api/v1/find/all |
| [SecScan](https://docs.quay.io/api/swagger/#SecScan) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/security |
| [Tag](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--tag-get) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/tag, /api/v1/repository/{namespace}/{repository}/tag/{tag}, /api/v1/repository/{namespace}/{repository}/tag/{tag}/history |
| [Team](https://docs.quay.io/api/swagger/#Team) | Yes | Yes | /api/v1/organization/{orgname}/team/{teamname}, /api/v1/organization/{orgname}/team/{teamname}/members, /api/v1/organization/{orgname}/team/{teamname}/permissions |
| [Trigger](https://docs.quay.io/api/swagger/#Trigger) | Yes | Yes | /api/v1/repository/{namespace}/{repository}/trigger/, /api/v1/repository/{namespace}/{repository}/trigger/{trigger_uuid}, /api/v1/repository/{namespace}/{repository}/trigger/{trigger_uuid}/start, /api/v1/repository/{namespace}/{repository}/trigger/{trigger_uuid}/activate |
| [User](https://docs.quay.io/api/swagger/#operation--api-v1-user-get) | Yes | Yes | /api/v1/user, /api/v1/user/starred, /api/v1/repository/{namespace}/{repository}/star |

## Authentication

1. Go to [Quay.io](https://quay.io) and log in
2. Create an **application token** (Account Settings → Applications) or a **robot account** token. Encrypted/CLI passwords are for `docker login`, not this API.
3. Use the token with `--token` / `-t`, `QUAY_TOKEN`, or a config file (`token:` in `~/.config/go-quay/config.yaml` on Linux)

Optional: `--quay-url` / `QUAY_URL` for a self-hosted registry. Output format: `--output` / `-O` (`json`, `yaml`, `table`). See [CLI Reference](docs/cli-reference.md).

## Development

```bash
make build     # Build binary
make test      # Run unit tests
make lint      # Run linter
make vet       # Run go vet
make coverage  # Tests with coverage
make ci        # lint + vet + test + build
```

## Prerequisites

- Go 1.26+
- golangci-lint (for `make lint`)

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
