# Contributing to go-quay

Thank you for contributing! This guide covers setup, development workflow, and how to add new API endpoints.

## Prerequisites

- Go 1.26+
- golangci-lint (for `make lint`)
- Quay.io API token (`QUAY_TOKEN`) and organization (`QUAY_ORG`) for integration tests

## Setup

```bash
git clone https://github.com/sebrandon1/go-quay.git
cd go-quay
go mod download
make build
```

## Development

```bash
make test              # Unit tests
make lint              # golangci-lint
make vet               # go vet
make fmt               # gofmt + goimports
make coverage          # Unit tests with coverage report
make govulncheck       # Vulnerability scan
make ci                # lint + vet + test + build
make integration-test  # Live API tests (requires QUAY_TOKEN and QUAY_ORG)
make check-swagger-alignment  # Verify lib methods match Quay Swagger spec
```

Run a single test:

```bash
go test ./lib/ -run TestCreateRepository -v
```

## Code Style

- Follow standard Go conventions and run `go fmt` before committing
- All tests and lint checks must pass
- Add tests for new functionality

## Pull Request Process

1. Fork the repository and create a feature branch
2. Make your changes with tests
3. Run `make test` and `make lint`
4. Update documentation if you add features or CLI commands
5. Open a pull request against `main`

## Adding a New API Endpoint

Every API domain follows the same pattern:

1. Add request/response structs to `lib/structs.go`
2. Add client method(s) to `lib/<domain>.go`
3. Add unit tests to `lib/<domain>_test.go` using `httptest.NewServer` and `lib.NewClientWithURL(token, server.URL+"/api/v1")`
4. Add Cobra command in `cmd/<domain>.go` and register it in `cmd/root.go`
5. Update `README.md` (API coverage table), `docs/cli-reference.md`, and `docs/library-guide.md`

See [CLAUDE.md](CLAUDE.md) for architecture details.

## Project Structure

```
go-quay/
├── cmd/           # CLI commands (Cobra)
├── lib/           # Quay API client library
├── docs/          # Guides, CLI reference, and tutorials
├── examples/      # Runnable example programs
├── scripts/       # Helper scripts
├── main.go        # Application entry point
└── Makefile       # Build and development commands
```

## Questions?

Open an issue on GitHub.
