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
make fmt                      # gofmt + goimports
make coverage                 # Unit tests with coverage report
make govulncheck              # Vulnerability scan
make ci                       # lint + vet + test + build
make check-swagger-alignment  # Verify lib functions match Quay's Swagger spec
make clean                    # Remove the built binary
go test ./lib/ -run TestCreateRepository -v  # Run a single test
```

Integration tests require `QUAY_TOKEN` and `QUAY_ORG` environment variables:
```bash
make integration-test
```

Optional for repo-scoped integration tests: `QUAY_NAMESPACE` and `QUAY_REPOSITORY`.

## Architecture

### Two-layer design: lib → cmd

Every API domain follows the same pattern:

1. **`lib/<domain>.go`** — API client methods on `*Client`. Each file's doc comment lists the HTTP endpoints it covers. All methods use the shared HTTP helpers in `client.go` (`get`, `post`, `put`, `delete`) and build URLs from the client's `BaseURL` field.

2. **`lib/structs.go`** — All request/response types for the entire API in one file, with JSON tags matching the Quay API.

3. **`cmd/<domain>.go`** — Cobra commands that parse flags, call `getClient()` (`lib.NewClientWithURL(token, quayURL)`), invoke the lib method, and print output via `printJSON()`. Commands are registered in `cmd/root.go` under the `get` parent command.

### Key patterns

- **`lib.DefaultQuayURL`** is a package-level `const` (`https://quay.io/api/v1`) in `lib/client.go`. Each `Client` has its own `BaseURL`.
- **`lib.NewClient(bearerToken)`** creates a client against `DefaultQuayURL`. **`lib.NewClientWithURL(bearerToken, baseURL)`** is used by the CLI and by unit tests (point `baseURL` at `httptest.NewServer`).
- All exported `Client` methods take `context.Context` as the first argument. CLI commands pass `cmd.Context()`.
- **`Client.Retry`** (`*RetryConfig`) is optional. When set, HTTP calls retry on 429 and 5xx. Retry backoff honors the request context. `NewClient` leaves it nil.
- API errors that include a Quay JSON body are returned as `*lib.QuayError` (implements `error`; use `errors.As` and `StatusCode()`).
- CLI authentication: `--token` / `-t`, else `$QUAY_TOKEN`, else the config file. API base URL: `--quay-url`, else `$QUAY_URL`, else config `quay-url`, else `DefaultQuayURL`. Output: `--output` / `-O` (`json`, `yaml`, or `table`).
- Config file (optional defaults): `token`, `namespace`, `quay-url` in platform config dir (`~/.config/go-quay/config.yaml` on Linux). Flags and env vars override the file.
- CLI commands use persistent flags (`--namespace`, `--repository`, `--token`) on parent commands so subcommands inherit them.
- The `cmd/helpers.go` file contains shared CLI utilities like `getClient()` and `printJSON()`.

### CLI command tree

All API commands hang off `go-quay get`:

`repository`, `billing`, `organization`, `permissions`, `tag`, `user`, `manifest`, `secscan`, `robot`, `search`, `team`, `build`, `notification`, `trigger`, `discovery`, `error`, `messages`, `prototype`, `repotoken`, `logs`, `mirror`

## Configuration

```bash
export QUAY_TOKEN="your-token"          # required for API calls
export QUAY_URL="https://quay.io/api/v1"  # optional; self-hosted Quay
```

CLI flag precedence: flags > environment variables > config file > built-in defaults.

## Adding a new API endpoint

1. Add request/response structs to `lib/structs.go`
2. Add client method(s) to `lib/<domain>.go` following the existing pattern
3. Add unit tests to `lib/<domain>_test.go` using `httptest.NewServer` and `NewClientWithURL`
4. Add Cobra command in `cmd/<domain>.go` and register it in `cmd/root.go`
5. Update `README.md` API coverage, `docs/cli-reference.md`, and `docs/library-guide.md`

## Requirements

- Go 1.26+
- golangci-lint (for `make lint`)
