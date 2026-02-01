# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A comprehensive Go wrapper around Quay.io APIs. Provides both a CLI and a library for interacting with Quay container registries.

## Common Commands

### Build
```bash
make build
```

### Run
```bash
./go-quay [command]
```

### Test
```bash
go test ./...
make test
```

### Integration Test
```bash
# Requires QUAY_TOKEN and QUAY_ORG environment variables
make integration-test
```

### Lint
```bash
make lint
make vet
```

### Clean
```bash
make clean
```

### Check API Alignment
```bash
make check-swagger-alignment
```

## Architecture

- **`cmd/`** - CLI command implementations using Cobra
- **`lib/`** - Quay API client library with full API coverage
- **`docs/`** - Documentation including library guide and tutorials
- **`examples/`** - Runnable example programs
- **`scripts/`** - Helper scripts for testing and validation
- **`main.go`** - Application entry point

## API Coverage

| API | Commands |
|-----|----------|
| Billing | user/plan, organization plans, invoices |
| Build | repository builds, logs, request, cancel |
| Discovery | API discovery endpoint |
| Error | error type details |
| Logs | aggregated logs, repository logs, organization logs |
| Manifest | manifest info, labels (CRUD) |
| Messages | system messages |
| Organization | org details, members, teams, robots, quota, auto-prune, applications |
| Permission | repository permissions (CRUD) |
| Prototype | default permission prototypes |
| Repository | CRUD operations, tags |
| RepositoryNotification | webhooks for repository events |
| Robot | user robots, permissions, regenerate |
| Search | repository search, global search |
| SecScan | security scanning results |
| Tag | tag info, history, update, delete, revert |
| Team | team CRUD, members, permissions |
| Trigger | build triggers |
| User | user info, starred repositories |

## Configuration

Set your Quay.io API token:
```bash
export QUAY_TOKEN="your-token"
```

## Requirements

- Go 1.25+
- Quay.io API token with appropriate permissions

## Code Style

- Follow standard Go conventions
- Use `go fmt` before committing
- Run `golangci-lint` for linting (config in `.golangci.yml`)
