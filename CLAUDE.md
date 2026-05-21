# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go wrapper around the Quay.io REST API (`https://quay.io/api/v1`). Provides both a reusable library (`lib/`) and a CLI (`cmd/`) built with Cobra.

## Common Commands

```bash
make build                    # Build the go-quay binary
make test                     # Run all unit tests (go test ./... -v)
make lint                     # Run golangci-lint
make vet                      # Run go vet
make check-swagger-alignment  # Verify lib functions match Quay's Swagger spec
go test ./lib/ -run TestCreateRepository -v  # Run a single test
```

Integration tests require `QUAY_TOKEN` and `QUAY_ORG` environment variables:
```bash
make integration-test
```

## Architecture

### Two-layer design: lib → cmd

Every API domain follows the same pattern:

1. **`lib/<domain>.go`** — API client methods on `*Client`. Each file's doc comment lists the HTTP endpoints it covers. All methods use the shared HTTP helpers in `client.go` (`get`, `post`, `put`, `delete`) and build URLs from the package-level `QuayURL` variable.

2. **`lib/structs.go`** — All request/response types for the entire API in one file, with JSON tags matching the Quay API.

3. **`cmd/<domain>.go`** — Cobra commands that parse flags, call `lib.NewClient(token)`, invoke the lib method, and print JSON output via `printJSON()`. Commands are registered in `cmd/root.go` under the `get` parent command.

### Key patterns

- **`lib.QuayURL`** is a package-level `var` (not `const`) set to `https://quay.io/api/v1` in `lib/logs.go`. Tests override it by pointing to `httptest.NewServer` and restoring with `defer`.
- **`lib.NewClient(bearerToken)`** creates an authenticated client. The token is passed via `--token` / `-t` flag in CLI commands.
- CLI commands use persistent flags (`--namespace`, `--repository`, `--token`) on parent commands so subcommands inherit them.
- The `cmd/helpers.go` file contains shared CLI utilities like `markFlagRequired()`.

### Adding a new API endpoint

1. Add request/response structs to `lib/structs.go`
2. Add client method(s) to `lib/<domain>.go` following the existing pattern
3. Add unit tests to `lib/<domain>_test.go` using `httptest.NewServer` with the `QuayURL` override pattern
4. Add Cobra command in `cmd/<domain>.go` and register it in `cmd/root.go`

## Configuration

Set your Quay.io API token:
```bash
export QUAY_TOKEN="your-token"
```

## Requirements

- Go 1.26+
- golangci-lint (for `make lint`)
