#!/bin/bash

# Test CLI commands structure for go-quay project
# This script verifies that all CLI commands and subcommands are properly structured

set -e

BINARY="./go-quay"
if [ ! -f "$BINARY" ]; then
	BINARY="./bin/go-quay"
fi

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo "Error: Binary not found at ./bin/go-quay or ./go-quay. Please build the project first (make build)."
    exit 1
fi

echo "Testing CLI command structure..."

echo "Testing main commands..."
$BINARY get --help
$BINARY create --help
$BINARY delete --help
$BINARY update --help
$BINARY list --help
$BINARY info --help

echo "Testing verb-first repository commands..."
$BINARY create repository --help
$BINARY delete repository --help
$BINARY update repository --help
$BINARY list repository --help
$BINARY info repository --help

echo "Testing repository commands..."
$BINARY get repository --help
$BINARY get repository create --help
$BINARY get repository update --help
$BINARY get repository delete --help
$BINARY get repository info --help
$BINARY get repository list --help

echo "Testing permissions commands..."
$BINARY get permissions --help
$BINARY get permissions list --help
$BINARY get permissions set --help
$BINARY get permissions remove --help

echo "Testing tag commands..."
$BINARY get tag --help
$BINARY get tag info --help
$BINARY get tag update --help
$BINARY get tag delete --help
$BINARY get tag history --help
$BINARY get tag revert --help
$BINARY get tag change --help
$BINARY get tag restore --help

echo "Testing user commands..."
$BINARY get user --help
$BINARY get user info --help
$BINARY get user starred --help
$BINARY get user star --help
$BINARY get user unstar --help

echo "Testing other command groups..."
$BINARY get billing --help
$BINARY get organization --help
$BINARY get logs --help
$BINARY get logs repo-aggregated-logs --help
$BINARY get mirror --help
$BINARY get mirror info --help

echo "✅ All CLI commands structure tests passed!"
