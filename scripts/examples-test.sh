#!/bin/bash

# Examples integration testing for go-quay
#
# Environment variables:
#   QUAY_TOKEN - Quay.io API token (optional; enables live API smoke tests)
#
# Example: QUAY_TOKEN=your_token ./scripts/examples-test.sh

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "Running examples tests..."

run_expect_fail() {
	local description="$1"
	local expected="$2"
	shift 2

	echo -n "Testing $description... "
	local output
	local status
	output="$("$@" 2>&1)" && status=0 || status=$?

	if [ "$status" -eq 0 ]; then
		echo "✗ (expected failure, got success)"
		exit 1
	fi
	if ! echo "$output" | grep -q "$expected"; then
		echo "✗ (expected output containing '$expected')"
		echo "$output"
		exit 1
	fi
	echo "✓"
}

echo "=== Building examples ==="
for example in basic-usage security-scan ci-cd-integration organization-management; do
	echo -n "Building examples/$example... "
	if go build -o /dev/null "./examples/$example"; then
		echo "✓"
	else
		echo "✗"
		exit 1
	fi
done

echo "=== Offline smoke tests ==="

run_expect_fail "basic-usage missing QUAY_TOKEN" \
	"QUAY_TOKEN environment variable is required" \
	env -u QUAY_TOKEN go run "./examples/basic-usage/main.go"

run_expect_fail "security-scan missing required flags" \
	"Usage:" \
	env -u QUAY_TOKEN go run "./examples/security-scan/main.go"

run_expect_fail "ci-cd-integration missing required flags" \
	"Usage:" \
	env -u QUAY_TOKEN go run "./examples/ci-cd-integration/main.go"

run_expect_fail "organization-management missing required flags" \
	"Usage:" \
	env -u QUAY_TOKEN go run "./examples/organization-management/main.go"

run_expect_fail "security-scan missing QUAY_TOKEN" \
	"QUAY_TOKEN environment variable is required" \
	env -u QUAY_TOKEN go run "./examples/security-scan/main.go" \
		--namespace testns --repository testrepo

run_expect_fail "ci-cd-integration missing QUAY_TOKEN" \
	"QUAY_TOKEN environment variable is required" \
	env -u QUAY_TOKEN go run "./examples/ci-cd-integration/main.go" \
		--namespace testns --repository testrepo

run_expect_fail "organization-management missing QUAY_TOKEN" \
	"QUAY_TOKEN environment variable is required" \
	env -u QUAY_TOKEN go run "./examples/organization-management/main.go" \
		--organization testorg

if [ -n "${QUAY_TOKEN:-}" ]; then
	echo "=== Live API smoke tests ==="

	printf "Testing basic-usage with QUAY_TOKEN... "
	output="$(go run "./examples/basic-usage/main.go" 2>&1)" && status=0 || status=$?
	if [ "$status" -eq 0 ] && echo "$output" | grep -q "Example Complete"; then
		echo "✓"
	else
		echo "⚠ (may not have valid token or API access)"
	fi
else
	echo "No QUAY_TOKEN available, skipping live API example tests"
fi

echo "✓ All examples tests completed successfully"
