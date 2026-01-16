# Contributing to go-quay

Thank you for your interest in contributing to go-quay! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Valid Quay.io API token (for integration testing)
- golangci-lint (for linting)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/sebrandon1/go-quay.git
   cd go-quay
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   ```

## Development Workflow

### Running Tests

```bash
make test
```

### Running Linter

```bash
make lint
```

### Building

```bash
make build
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Ensure `golangci-lint` passes with no issues
- Add tests for new functionality

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Run tests and linting:
   ```bash
   make test
   make lint
   ```
5. Commit your changes with a descriptive commit message
6. Push to your fork
7. Open a Pull Request against the `main` branch

### PR Requirements

- All tests must pass
- Linting must pass with no issues
- Include tests for new functionality
- Update documentation if adding new features or CLI commands

## Adding New Quay API Endpoints

When adding support for a new Quay API endpoint:

1. Add data structures to `lib/structs.go`
2. Add the client method to `lib/client.go`
3. Add tests to `lib/client_test.go`
4. Add CLI command to appropriate file in `cmd/`
5. Add CLI tests if applicable
6. Update the README.md with usage documentation
7. Update the API coverage table in README.md

## Project Structure

```
go-quay/
├── cmd/           # CLI command implementations
│   ├── root.go    # Root command and subcommand registration
│   ├── billing.go # Billing API commands
│   ├── build.go   # Build API commands
│   ├── manifest.go# Manifest API commands
│   ├── organization.go # Organization API commands
│   ├── repository.go   # Repository API commands
│   └── ...        # Other API command files
├── lib/           # Quay API client library
│   ├── client.go  # HTTP client and API methods
│   └── structs.go # Data structures
├── scripts/       # Helper scripts
├── main.go        # Application entry point
└── Makefile       # Build and development commands
```

## Questions?

If you have questions, please open an issue on GitHub.
