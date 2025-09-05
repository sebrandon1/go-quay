#!/bin/bash

# Test CLI commands structure for go-quay project
# This script verifies that all CLI commands and subcommands are properly structured

set -e

BINARY="./bin/go-quay"

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo "Error: Binary $BINARY not found. Please build the project first."
    exit 1
fi

echo "Testing CLI command structure..."

# Test main commands exist
echo "Testing main commands..."
$BINARY get --help

# Test repository commands
echo "Testing repository commands..."
$BINARY get repository --help
$BINARY get repository create --help
$BINARY get repository update --help
$BINARY get repository delete --help
$BINARY get repository info --help

# Test permissions commands
echo "Testing permissions commands..."
$BINARY get permissions --help
$BINARY get permissions list --help
$BINARY get permissions set --help
$BINARY get permissions remove --help

# Test tag commands
echo "Testing tag commands..."
$BINARY get tag --help
$BINARY get tag info --help
$BINARY get tag update --help
$BINARY get tag delete --help
$BINARY get tag history --help
$BINARY get tag revert --help

# Test user commands
echo "Testing user commands..."
$BINARY get user --help
$BINARY get user info --help
$BINARY get user starred --help
$BINARY get user star --help
$BINARY get user unstar --help

# Test existing commands still work
echo "Testing existing commands..."
$BINARY get billing --help
$BINARY get organization --help
$BINARY get aggregatedlogs --help

echo "✅ All CLI commands structure tests passed!"
