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

### Lint
```bash
make lint
```

## Architecture

- **`cmd/`** - CLI command implementations using Cobra
- **`lib/`** - Quay API client library with full API coverage
- **`scripts/`** - Helper scripts
- **`main.go`** - Application entry point

## API Coverage

| API | Commands |
|-----|----------|
| Billing | user/plan, organization plans, invoices |
| Build | repository builds, logs |
| Discovery | API discovery endpoint |
| Manifest | manifest info, labels |
| Organization | org details, members, teams, robots, quota |
| Permission | repository permissions |
| Repository | CRUD operations, tags |
| Robot | user robots, permissions |
| Search | repository search |
| SecScan | security scanning results |

## Configuration

Set your Quay.io API token:
```bash
export QUAY_API_TOKEN="your-token"
```

## Requirements

- Go 1.21+
- Quay.io API token with appropriate permissions

## Code Style

- Follow standard Go conventions
- Use `go fmt` before committing
- Run `golangci-lint` for linting (config in `.golangci.yml`)
