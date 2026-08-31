# Examples

Runnable programs demonstrating go-quay library usage. Each example requires `QUAY_TOKEN`.

| Example | Description |
|---------|-------------|
| [basic-usage](./basic-usage/) | Client setup, user info, and repository listing |
| [security-scan](./security-scan/) | Vulnerability scanning and reporting |
| [ci-cd-integration](./ci-cd-integration/) | Robot accounts, permissions, and webhooks |
| [organization-management](./organization-management/) | Teams, quotas, and org administration |

## Running an Example

```bash
export QUAY_TOKEN="your-token"
cd basic-usage
go run main.go
```

`basic-usage` also reads `QUAY_NAMESPACE` (optional; defaults to the authenticated username). That variable is not a CLI flag — the CLI uses `--namespace` or config `namespace`.

Some examples accept flags — run with `-h` for usage.

Run offline smoke tests with `make examples-test` from the repository root.
