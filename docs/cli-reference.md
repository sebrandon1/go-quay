# CLI Reference

Complete command reference for the `go-quay` CLI. All commands require a Quay.io authentication token via `--token` or `-t`.

## Billing API

The billing API provides access to subscription plans, billing information, and invoices.

### Get available subscription plans
```bash
go-quay get billing plans --token YOUR_TOKEN
```

### Get user billing information
```bash
go-quay get billing user-info --token YOUR_TOKEN
go-quay get billing user-subscription --token YOUR_TOKEN
```

### Get organization billing information
```bash
go-quay get billing org-info --organization ORG_NAME --token YOUR_TOKEN
go-quay get billing org-subscription --organization ORG_NAME --token YOUR_TOKEN
go-quay get billing org-invoices --organization ORG_NAME --token YOUR_TOKEN
```

## Build API

Manage automated image builds from Dockerfiles.

### List builds for a repository
```bash
go-quay get build list \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

### Get build details
```bash
go-quay get build info \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid BUILD_UUID \
  --token YOUR_TOKEN
```

### Get build logs
```bash
go-quay get build logs \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid BUILD_UUID \
  --token YOUR_TOKEN
```

### Request a new build
```bash
go-quay get build request \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --archive-url "https://example.com/archive.tar.gz" \
  --tag latest \
  --token YOUR_TOKEN
```

### Cancel a build
```bash
go-quay get build cancel \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid BUILD_UUID \
  --confirm \
  --token YOUR_TOKEN
```

**Build Phases:**
- `waiting`: Build is queued
- `starting`: Build is starting
- `building`: Build is in progress
- `pushing`: Pushing built image
- `complete`: Build completed successfully
- `error`: Build failed

## Trigger API

Manage build triggers for automated builds when code is pushed to connected source repositories.

### List triggers for a repository
```bash
go-quay get trigger list \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

### Get trigger details
```bash
go-quay get trigger info \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid TRIGGER_UUID \
  --token YOUR_TOKEN
```

### Delete a trigger
```bash
go-quay get trigger delete \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid TRIGGER_UUID \
  --token YOUR_TOKEN
```

### Enable a trigger
```bash
go-quay get trigger enable \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid TRIGGER_UUID \
  --token YOUR_TOKEN
```

### Disable a trigger
```bash
go-quay get trigger disable \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid TRIGGER_UUID \
  --token YOUR_TOKEN
```

### Manually start a build from trigger
```bash
go-quay get trigger start \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid TRIGGER_UUID \
  --token YOUR_TOKEN

# Optionally specify a commit SHA
go-quay get trigger start \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid TRIGGER_UUID \
  --commit-sha abc123def456 \
  --token YOUR_TOKEN
```

### Activate a trigger
```bash
go-quay get trigger activate \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid TRIGGER_UUID \
  --token YOUR_TOKEN
```

**Supported Trigger Services:**
- `github`: GitHub repository
- `gitlab`: GitLab repository
- `bitbucket`: Bitbucket repository
- `custom-git`: Custom git repository

## Repository Notification API

Manage webhooks for repository events.

### List notifications for a repository
```bash
go-quay get notification list \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

### Get notification details
```bash
go-quay get notification info \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid NOTIFICATION_UUID \
  --token YOUR_TOKEN
```

### Create a webhook notification
```bash
go-quay get notification create \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --event repo_push \
  --method webhook \
  --url "https://example.com/webhook" \
  --title "Push Webhook" \
  --token YOUR_TOKEN
```

### Create a Slack notification
```bash
go-quay get notification create \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --event build_success \
  --method slack \
  --url "https://hooks.slack.com/services/..." \
  --title "Build Success" \
  --token YOUR_TOKEN
```

### Test a notification
```bash
go-quay get notification test \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid NOTIFICATION_UUID \
  --token YOUR_TOKEN
```

### Reset notification failure count
```bash
go-quay get notification reset \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid NOTIFICATION_UUID \
  --token YOUR_TOKEN
```

### Delete a notification
```bash
go-quay get notification delete \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid NOTIFICATION_UUID \
  --confirm \
  --token YOUR_TOKEN
```

**Supported Events:**
- `repo_push`: Image push to repository
- `build_queued`: Build has been queued
- `build_start`: Build has started
- `build_success`: Build completed successfully
- `build_failure`: Build failed
- `build_canceled`: Build was canceled
- `vulnerability_found`: New vulnerability discovered

**Supported Methods:**
- `webhook`: HTTP POST to a URL
- `email`: Email notification
- `slack`: Slack webhook

## Logs API

Repository, organization, and user activity logs with optional date range filtering.

### Get aggregated repository logs
```bash
go-quay get logs repo-aggregated-logs \
  -n NAMESPACE \
  -r REPOSITORY \
  -s "2024-01-01" \
  -e "2024-01-07" \
  -t YOUR_TOKEN
```

### Get repository logs with date filtering
```bash
go-quay get logs repo-logs \
  -n NAMESPACE \
  -r REPOSITORY \
  --startdate "2024-05-01" \
  --enddate "2024-05-31" \
  -t YOUR_TOKEN
```

### Get organization and user logs
```bash
# Organization logs
go-quay get logs org-logs -o ORG_NAME -t YOUR_TOKEN

# Organization aggregated logs
go-quay get logs org-aggregated-logs -o ORG_NAME -s "2024-01-01" -e "2024-01-31" -t YOUR_TOKEN

# User logs
go-quay get logs user-logs -t YOUR_TOKEN

# User aggregated logs
go-quay get logs user-aggregated-logs -s "2024-01-01" -e "2024-01-31" -t YOUR_TOKEN
```

## Repository API

Full CRUD operations for repository management.

### Get repository information with tags
```bash
go-quay get repository info \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

### Create a new repository
```bash
go-quay get repository create \
  --namespace myorg \
  --repository mynewrepo \
  --visibility private \
  --description "My new application repository" \
  --token YOUR_TOKEN
```

### Update repository settings
```bash
go-quay get repository update \
  --namespace myorg \
  --repository myrepo \
  --description "Updated description" \
  --token YOUR_TOKEN
```

### Delete a repository
```bash
go-quay get repository delete \
  --namespace myorg \
  --repository oldrepo \
  --confirm \
  --token YOUR_TOKEN
```

## Repository Permissions API

Manage who can access your repositories and what level of access they have.

### List repository permissions
```bash
go-quay get permissions list \
  --namespace myorg \
  --repository myrepo \
  --token YOUR_TOKEN
```

### Set user/robot permissions
```bash
go-quay get permissions set \
  --namespace myorg \
  --repository myrepo \
  --user john.doe \
  --role write \
  --token YOUR_TOKEN
```

### Remove permissions
```bash
go-quay get permissions remove \
  --namespace myorg \
  --repository myrepo \
  --user john.doe \
  --token YOUR_TOKEN
```

**Permission Roles:**
- `read`: Pull images and view repository
- `write`: Push images, pull images, and view repository
- `admin`: Full access including permission management

## Tag API

Tag management with detailed metadata, history, and operations.

### Get detailed tag information
```bash
go-quay get tag info \
  --namespace myorg \
  --repository myrepo \
  --tag v1.0.0 \
  --token YOUR_TOKEN
```

### Update tag metadata
```bash
go-quay get tag update \
  --namespace myorg \
  --repository myrepo \
  --tag v1.0.0 \
  --expiration "2024-12-31T23:59:59Z" \
  --token YOUR_TOKEN
```

### Delete a specific tag
```bash
go-quay get tag delete \
  --namespace myorg \
  --repository myrepo \
  --tag old-version \
  --confirm \
  --token YOUR_TOKEN
```

### View tag history
```bash
go-quay get tag history \
  --namespace myorg \
  --repository myrepo \
  --tag latest \
  --token YOUR_TOKEN
```

### Revert tag to previous state
```bash
go-quay get tag revert \
  --namespace myorg \
  --repository myrepo \
  --tag latest \
  --manifest sha256:abc123... \
  --token YOUR_TOKEN
```

## Manifest API

Inspect and manage container image manifests, including layers, configuration, and labels.

### Get manifest information
```bash
go-quay get manifest info \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --token YOUR_TOKEN
```

### Delete a manifest
```bash
go-quay get manifest delete \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --confirm \
  --token YOUR_TOKEN
```

### List manifest labels
```bash
go-quay get manifest labels \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --token YOUR_TOKEN
```

### Get a specific label
```bash
go-quay get manifest label \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --label-id label-123 \
  --token YOUR_TOKEN
```

### Add a label to a manifest
```bash
go-quay get manifest add-label \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --key "environment" \
  --value "production" \
  --token YOUR_TOKEN
```

### Remove a label from a manifest
```bash
go-quay get manifest remove-label \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --label-id label-123 \
  --token YOUR_TOKEN
```

## Security Scan (SecScan) API

Retrieve security vulnerability information for container images.

### Get security scan results for a manifest
```bash
go-quay get secscan info \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --token YOUR_TOKEN

# Without vulnerability details (faster)
go-quay get secscan info \
  -n myorg \
  -r myrepo \
  -m sha256:abc123def456... \
  --vulnerabilities=false \
  -t YOUR_TOKEN
```

**Scan Status Values:**
- `scanned`: Scan completed successfully, results available
- `queued`: Scan is queued and pending
- `scanning`: Scan is currently in progress
- `unsupported`: Image type is not supported for scanning
- `failed`: Scan failed

**Vulnerability Severity Levels:**
- `Critical`: Severe vulnerabilities requiring immediate attention
- `High`: Important vulnerabilities to address soon
- `Medium`: Moderate risk vulnerabilities
- `Low`: Minor vulnerabilities
- `Negligible`: Minimal impact vulnerabilities
- `Unknown`: Severity not determined

## Robot API

Manage user-level robot accounts for CI/CD automation.

### List all robot accounts
```bash
go-quay get robot list --token YOUR_TOKEN
```

### Get robot account details
```bash
go-quay get robot info --name deploybot --token YOUR_TOKEN
```

### Create a new robot account
```bash
go-quay get robot create \
  --name mybot \
  --description "CI/CD automation robot" \
  --token YOUR_TOKEN
```

### Delete a robot account
```bash
go-quay get robot delete --name oldbot --confirm --token YOUR_TOKEN
```

### Regenerate robot token
```bash
go-quay get robot regenerate --name deploybot --token YOUR_TOKEN
```

### Get robot permissions
```bash
go-quay get robot permissions --name deploybot --token YOUR_TOKEN
```

## Search API

Search for repositories, users, organizations, and other entities.

### Search for repositories
```bash
go-quay get search repositories --query "nginx" --token YOUR_TOKEN
```

### Search all entity types
```bash
go-quay get search all --query "redhat" --token YOUR_TOKEN
```

## Team API

Manage teams within an organization.

### List teams in an organization
```bash
go-quay get team list --organization myorg --token YOUR_TOKEN
```

### Get team information
```bash
go-quay get team info --organization myorg --name developers --token YOUR_TOKEN
```

### Create a new team
```bash
go-quay get team create \
  --organization myorg \
  --name developers \
  --description "Development team" \
  --role member \
  --token YOUR_TOKEN
```

### Update team settings
```bash
go-quay get team update \
  --organization myorg \
  --name developers \
  --description "Updated description" \
  --role creator \
  --token YOUR_TOKEN
```

### Delete a team
```bash
go-quay get team delete --organization myorg --name developers --confirm --token YOUR_TOKEN
```

### Manage team members
```bash
# List team members
go-quay get team members --organization myorg --name developers --token YOUR_TOKEN

# Add a member
go-quay get team add-member --organization myorg --name developers --member username --token YOUR_TOKEN

# Remove a member
go-quay get team remove-member --organization myorg --name developers --member username --confirm --token YOUR_TOKEN
```

### Manage team repository permissions
```bash
# List team repository permissions
go-quay get team permissions --organization myorg --name developers --token YOUR_TOKEN

# Set repository permission for a team
go-quay get team set-permission \
  --organization myorg \
  --name developers \
  --repository myrepo \
  --role write \
  --token YOUR_TOKEN

# Remove repository permission from a team
go-quay get team remove-permission \
  --organization myorg \
  --name developers \
  --repository myrepo \
  --confirm \
  --token YOUR_TOKEN
```

**Team Roles:**
- `member`: Inherits default permissions
- `creator`: Can create new repositories
- `admin`: Full administrative access

## User API

Manage your user account information and starred repositories.

### Get current user information
```bash
go-quay get user info --token YOUR_TOKEN
```

### List starred repositories
```bash
go-quay get user starred --token YOUR_TOKEN
```

### Star / unstar a repository
```bash
go-quay get user star --namespace quay --repository quay --token YOUR_TOKEN
go-quay get user unstar --namespace quay --repository quay --token YOUR_TOKEN
```

## Organization API

Comprehensive management of organizations, teams, members, robots, and settings.

### Get organization information
```bash
go-quay get organization info --organization ORG_NAME --token YOUR_TOKEN
```

### List members
```bash
go-quay get organization members -o myorg -t YOUR_TOKEN
```

### List teams
```bash
go-quay get organization teams -o myorg -t YOUR_TOKEN
```

### Get specific team information
```bash
go-quay get organization team -o myorg --team TEAM_NAME -t YOUR_TOKEN
go-quay get organization team-members -o myorg --team TEAM_NAME -t YOUR_TOKEN
```

### List organization robot accounts
```bash
go-quay get organization robots -o myorg -t YOUR_TOKEN
```

### Get organization quota and policies
```bash
go-quay get organization quota -o myorg -t YOUR_TOKEN
go-quay get organization auto-prune -o myorg -t YOUR_TOKEN
go-quay get organization applications -o myorg -t YOUR_TOKEN
```

## Discovery API

Get information about available API endpoints and versions.

```bash
go-quay get discovery --token YOUR_TOKEN
```

## Error API

Get details about specific error types returned by the Quay.io API.

```bash
go-quay get error --type invalid_token --token YOUR_TOKEN
```

## Messages API

Get system-wide messages for the authenticated user.

```bash
go-quay get messages --token YOUR_TOKEN
```

## Prototype API

Manage default permission prototypes automatically applied to new repositories.

### List all prototypes
```bash
go-quay get prototype list --organization myorg --token YOUR_TOKEN
```

### Get prototype details
```bash
go-quay get prototype info --organization myorg --uuid PROTOTYPE_UUID --token YOUR_TOKEN
```

### Create a prototype
```bash
go-quay get prototype create \
  --organization myorg \
  --delegate-name devteam \
  --delegate-kind team \
  --role write \
  --token YOUR_TOKEN
```

### Update a prototype
```bash
go-quay get prototype update \
  --organization myorg \
  --uuid PROTOTYPE_UUID \
  --role admin \
  --token YOUR_TOKEN
```

### Delete a prototype
```bash
go-quay get prototype delete \
  --organization myorg \
  --uuid PROTOTYPE_UUID \
  --confirm \
  --token YOUR_TOKEN
```

**Delegate Kinds:**
- `user`: A specific user account
- `team`: A team within the organization
- `robot`: A robot account

## RepoToken API (DEPRECATED)

> **Warning:** Repository tokens are deprecated. Use robot accounts instead.

### List repository tokens
```bash
go-quay get repotoken list \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

### Get token details
```bash
go-quay get repotoken info \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --code TOKEN_CODE \
  --token YOUR_TOKEN
```

### Create a token
```bash
go-quay get repotoken create \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --name "CI Token" \
  --token YOUR_TOKEN
```

### Update a token
```bash
go-quay get repotoken update \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --code TOKEN_CODE \
  --role write \
  --token YOUR_TOKEN
```

### Delete a token
```bash
go-quay get repotoken delete \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --code TOKEN_CODE \
  --confirm \
  --token YOUR_TOKEN
```
