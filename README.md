# go-quay

[![Pre-Main Checks](https://github.com/sebrandon1/go-quay/actions/workflows/pre-main.yaml/badge.svg)](https://github.com/sebrandon1/go-quay/actions/workflows/pre-main.yaml)
[![Quay API Verified Nightly](https://github.com/sebrandon1/go-quay/actions/workflows/nightly.yaml/badge.svg)](https://github.com/sebrandon1/go-quay/actions/workflows/nightly.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sebrandon1/go-quay)](https://golang.org/)
[![License](https://img.shields.io/github/license/sebrandon1/go-quay)](https://github.com/sebrandon1/go-quay/blob/main/LICENSE)

A Go wrapper around Quay APIs

## Table of API Coverage

The following APIs are covered by the repo. Each API links to the corresponding section in the [Quay.io Swagger documentation](https://docs.quay.io/api/swagger/):
| API                    | Cmd     | Lib     | Covered                                                                                                                                                                                                             |
| ---------------------- | ------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Billing](https://docs.quay.io/api/swagger/#operation--api-v1-user-plan-get)                | Yes     | Yes     | /api/v1/user/plan, /api/v1/organization/{orgname}/plan, /api/v1/organization/{orgname}/invoices, /api/v1/plans/                                                                                                   |
| [Build](https://docs.quay.io/api/swagger/#Build)                  | Yes     | Yes     | /api/v1/repository/{namespace}/{repository}/build/, /api/v1/repository/{namespace}/{repository}/build/{build_uuid}, /api/v1/repository/{namespace}/{repository}/build/{build_uuid}/logs |
| [Discovery](https://docs.quay.io/api/swagger/#Discovery)              | No      | No      |                                                                                                                                                                                                                     |
| [Error](https://docs.quay.io/api/swagger/#Error)                  | No      | No      |                                                                                                                                                                                                                     |
| [Messages](https://docs.quay.io/api/swagger/#Messages)               | No      | No      |                                                                                                                                                                                                                     |
| [Logs](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--aggregatelogs-get)                   | Partial | Partial | /api/v1/repository/{namespace}/{repository}/aggregatelogs, /api/v1/repository/{namespace}/{repository}/logs, /api/v1/organization/{orgname}/logs |
| [Manifest](https://docs.quay.io/api/swagger/#Manifest)               | Yes     | Yes     | /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}, /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/labels, /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/labels/{labelid} |
| [Organization](https://docs.quay.io/api/swagger/#operation--api-v1-organization--orgname--get)           | Yes     | Yes     | /api/v1/organization/{orgname}, /api/v1/organization/{orgname}/members, /api/v1/organization/{orgname}/teams, /api/v1/organization/{orgname}/team/{teamname}, /api/v1/organization/{orgname}/robots, /api/v1/organization/{orgname}/quota, /api/v1/organization/{orgname}/autoprunepolicy, /api/v1/organization/{orgname}/applications |
| [Permission](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--permissions-get)             | Yes     | Yes     | /api/v1/repository/{namespace}/{repository}/permissions, /api/v1/repository/{namespace}/{repository}/permissions/{username} |
| [Prototype](https://docs.quay.io/api/swagger/#Prototype)              | No      | No      |                                                                                                                                                                                                                     |
| [Repository](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--get)             | Yes     | Yes     | /api/v1/repository/{namespace}/{repository}, /api/v1/repository/{namespace}/{repository}/tag, /api/v1/repository, /api/v1/repository/{namespace}/{repository} (CRUD) |
| [RepositoryNotification](https://docs.quay.io/api/swagger/#RepositoryNotification) | Yes     | Yes     | /api/v1/repository/{namespace}/{repository}/notification/, /api/v1/repository/{namespace}/{repository}/notification/{uuid}, /api/v1/repository/{namespace}/{repository}/notification/{uuid}/test |
| [RepoToken](https://docs.quay.io/api/swagger/#RepoToken)              | No      | No      |                                                                                                                                                                                                                     |
| [Robot](https://docs.quay.io/api/swagger/#Robot)                  | Yes     | Yes     | /api/v1/user/robots, /api/v1/user/robots/{robot_shortname}, /api/v1/user/robots/{robot_shortname}/regenerate, /api/v1/user/robots/{robot_shortname}/permissions |
| [Search](https://docs.quay.io/api/swagger/#Search)                 | Yes     | Yes     | /api/v1/find/repositories, /api/v1/find/all |
| [SecScan](https://docs.quay.io/api/swagger/#SecScan)                | Yes     | Yes     | /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/security |
| [Tag](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--tag-get)                    | Yes     | Yes     | /api/v1/repository/{namespace}/{repository}/tag, /api/v1/repository/{namespace}/{repository}/tag/{tag}, /api/v1/repository/{namespace}/{repository}/tag/{tag}/history |
| [Team](https://docs.quay.io/api/swagger/#Team)                   | Yes     | Yes     | /api/v1/organization/{orgname}/team/{teamname}, /api/v1/organization/{orgname}/team/{teamname}/members, /api/v1/organization/{orgname}/team/{teamname}/permissions |
| [Trigger](https://docs.quay.io/api/swagger/#Trigger)                | No      | No      |                                                                                                                                                                                                                     |
| [User](https://docs.quay.io/api/swagger/#operation--api-v1-user-get)                   | Yes     | Yes     | /api/v1/user, /api/v1/user/starred, /api/v1/repository/{namespace}/{repository}/star | 

## Authentication

All API commands require a Quay.io authentication token. You can obtain a token from your Quay.io account settings:

1. Go to [Quay.io](https://quay.io) and log in to your account
2. Navigate to **Account Settings** → **Robot Accounts** or **User Settings** → **CLI Password**
3. Generate a new token with appropriate permissions
4. Use the token with the `--token` or `-t` flag in commands

## Usage Examples

### Billing API

The billing API provides access to subscription plans, billing information, and invoices.

📖 **API Reference:** [Billing endpoints in Swagger](https://docs.quay.io/api/swagger/#operation--api-v1-user-plan-get)

#### Get available subscription plans
```bash
./go-quay get billing plans --token YOUR_TOKEN
```

#### Get user billing information
```bash
./go-quay get billing user-info --token YOUR_TOKEN
./go-quay get billing user-subscription --token YOUR_TOKEN
```

#### Get organization billing information
```bash
./go-quay get billing org-info --organization ORG_NAME --token YOUR_TOKEN
./go-quay get billing org-subscription --organization ORG_NAME --token YOUR_TOKEN
./go-quay get billing org-invoices --organization ORG_NAME --token YOUR_TOKEN
```

### Build API

The build API allows you to manage automated image builds from Dockerfiles.

📖 **API Reference:** [Build endpoints in Swagger](https://docs.quay.io/api/swagger/#Build)

#### List builds for a repository
```bash
./go-quay get build list \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

#### Get build details
```bash
./go-quay get build info \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid BUILD_UUID \
  --token YOUR_TOKEN
```

#### Get build logs
```bash
./go-quay get build logs \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid BUILD_UUID \
  --token YOUR_TOKEN
```

#### Request a new build
```bash
./go-quay get build request \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --archive-url "https://example.com/archive.tar.gz" \
  --tag latest \
  --token YOUR_TOKEN
```

#### Cancel a build
```bash
./go-quay get build cancel \
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

### Repository Notification API

The notification API allows you to manage webhooks for repository events.

📖 **API Reference:** [RepositoryNotification endpoints in Swagger](https://docs.quay.io/api/swagger/#RepositoryNotification)

#### List notifications for a repository
```bash
./go-quay get notification list \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

#### Get notification details
```bash
./go-quay get notification info \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid NOTIFICATION_UUID \
  --token YOUR_TOKEN
```

#### Create a webhook notification
```bash
./go-quay get notification create \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --event repo_push \
  --method webhook \
  --url "https://example.com/webhook" \
  --title "Push Webhook" \
  --token YOUR_TOKEN
```

#### Create a Slack notification
```bash
./go-quay get notification create \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --event build_success \
  --method slack \
  --url "https://hooks.slack.com/services/..." \
  --title "Build Success" \
  --token YOUR_TOKEN
```

#### Test a notification
```bash
./go-quay get notification test \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid NOTIFICATION_UUID \
  --token YOUR_TOKEN
```

#### Reset notification failure count
```bash
./go-quay get notification reset \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --uuid NOTIFICATION_UUID \
  --token YOUR_TOKEN
```

#### Delete a notification
```bash
./go-quay get notification delete \
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

### Logs API

The logs API provides access to repository activity logs and aggregated statistics.

📖 **API Reference:** [Logs endpoints in Swagger](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--aggregatelogs-get)

#### Get aggregated repository logs
```bash
# Get logs for the last 7 days
./go-quay get aggregatedlogs \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --startdate "01/01/2024" \
  --enddate "01/07/2024" \
  --token YOUR_TOKEN
```

#### Example with specific dates
```bash
./go-quay get aggregatedlogs \
  -n quay \
  -r my-app \
  -s "12/01/2023" \
  -e "12/31/2023" \
  -t YOUR_TOKEN
```

### Repository API

The repository API provides full CRUD (Create, Read, Update, Delete) operations for repository management.

📖 **API Reference:** [Repository endpoints in Swagger](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--get)

#### Get repository information with tags
```bash
./go-quay get repository info \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

#### Create a new repository
```bash
# Create a private repository with description
./go-quay get repository create \
  --namespace myorg \
  --repository mynewrepo \
  --visibility private \
  --description "My new application repository" \
  --token YOUR_TOKEN

# Create a public repository
./go-quay get repository create \
  -n myorg \
  -r publicrepo \
  -v public \
  -d "Public demo repository" \
  -t YOUR_TOKEN
```

#### Update repository settings
```bash
# Update repository description
./go-quay get repository update \
  --namespace myorg \
  --repository myrepo \
  --description "Updated description" \
  --token YOUR_TOKEN

# Change repository visibility to public
./go-quay get repository update \
  -n myorg \
  -r myrepo \
  -v public \
  -t YOUR_TOKEN
```

#### Delete a repository
```bash
# Delete repository (requires confirmation)
./go-quay get repository delete \
  --namespace myorg \
  --repository oldrepo \
  --confirm \
  --token YOUR_TOKEN
```

**Note:** Repository deletion is irreversible and will remove all images and tags.

### Repository Permissions API

Manage who can access your repositories and what level of access they have.

📖 **API Reference:** [Permission endpoints in Swagger](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--permissions-get)

#### List repository permissions
```bash
./go-quay get permissions list \
  --namespace myorg \
  --repository myrepo \
  --token YOUR_TOKEN
```

#### Set user/robot permissions
```bash
# Give a user write access
./go-quay get permissions set \
  --namespace myorg \
  --repository myrepo \
  --user john.doe \
  --role write \
  --token YOUR_TOKEN

# Give a robot account read access
./go-quay get permissions set \
  -n myorg \
  -r myrepo \
  -u myorg+deploybot \
  -R read \
  -t YOUR_TOKEN

# Grant admin access
./go-quay get permissions set \
  -n myorg \
  -r myrepo \
  -u jane.smith \
  -R admin \
  -t YOUR_TOKEN
```

#### Remove permissions
```bash
./go-quay get permissions remove \
  --namespace myorg \
  --repository myrepo \
  --user john.doe \
  --token YOUR_TOKEN
```

**Permission Roles:**
- `read`: Pull images and view repository
- `write`: Push images, pull images, and view repository  
- `admin`: Full access including permission management

### Enhanced Tag API

Comprehensive tag management with detailed metadata, history, and operations.

📖 **API Reference:** [Tag endpoints in Swagger](https://docs.quay.io/api/swagger/#operation--api-v1-repository--namespace---repository--tag-get)

#### Get detailed tag information
```bash
./go-quay get tag info \
  --namespace myorg \
  --repository myrepo \
  --tag v1.0.0 \
  --token YOUR_TOKEN
```

#### Update tag metadata
```bash
# Set tag expiration
./go-quay get tag update \
  --namespace myorg \
  --repository myrepo \
  --tag v1.0.0 \
  --expiration "2024-12-31T23:59:59Z" \
  --token YOUR_TOKEN
```

#### Delete a specific tag
```bash
./go-quay get tag delete \
  --namespace myorg \
  --repository myrepo \
  --tag old-version \
  --confirm \
  --token YOUR_TOKEN
```

#### View tag history
```bash
./go-quay get tag history \
  --namespace myorg \
  --repository myrepo \
  --tag latest \
  --token YOUR_TOKEN
```

#### Revert tag to previous state
```bash
./go-quay get tag revert \
  --namespace myorg \
  --repository myrepo \
  --tag latest \
  --manifest sha256:abc123... \
  --token YOUR_TOKEN
```

### Manifest API

Inspect and manage container image manifests, including layers, configuration, and labels.

📖 **API Reference:** [Manifest endpoints in Swagger](https://docs.quay.io/api/swagger/#Manifest)

#### Get manifest information
```bash
# Get detailed manifest info by digest
./go-quay get manifest info \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --token YOUR_TOKEN
```

#### Delete a manifest
```bash
# Delete manifest (requires confirmation)
./go-quay get manifest delete \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --confirm \
  --token YOUR_TOKEN
```

**Note:** Manifest deletion is irreversible and will remove all tags pointing to this manifest.

#### List manifest labels
```bash
./go-quay get manifest labels \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --token YOUR_TOKEN
```

#### Get a specific label
```bash
./go-quay get manifest label \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --label-id label-123 \
  --token YOUR_TOKEN
```

#### Add a label to a manifest
```bash
# Add a simple label
./go-quay get manifest add-label \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --key "environment" \
  --value "production" \
  --token YOUR_TOKEN

# Add a label with custom media type
./go-quay get manifest add-label \
  -n myorg \
  -r myrepo \
  -m sha256:abc123def456... \
  -k "config" \
  -v '{"setting": "value"}' \
  --media-type "application/json" \
  -t YOUR_TOKEN
```

#### Remove a label from a manifest
```bash
./go-quay get manifest remove-label \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --label-id label-123 \
  --token YOUR_TOKEN
```

### Security Scan (SecScan) API

Retrieve security vulnerability information for container images.

📖 **API Reference:** [SecScan endpoints in Swagger](https://docs.quay.io/api/swagger/#SecScan)

#### Get security scan results for a manifest
```bash
# Get security scan with vulnerability details
./go-quay get secscan info \
  --namespace myorg \
  --repository myrepo \
  --manifest sha256:abc123def456... \
  --token YOUR_TOKEN

# Get security scan without vulnerability details (faster)
./go-quay get secscan info \
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

### Robot API

Manage user-level robot accounts for CI/CD automation and automated workflows.

📖 **API Reference:** [Robot endpoints in Swagger](https://docs.quay.io/api/swagger/#Robot)

#### List all robot accounts
```bash
./go-quay get robot list \
  --token YOUR_TOKEN
```

#### Get robot account details
```bash
./go-quay get robot info \
  --name deploybot \
  --token YOUR_TOKEN
```

#### Create a new robot account
```bash
./go-quay get robot create \
  --name mybot \
  --description "CI/CD automation robot" \
  --token YOUR_TOKEN
```

**Note:** Save the token from the response - it will not be shown again!

#### Delete a robot account
```bash
./go-quay get robot delete \
  --name oldbot \
  --confirm \
  --token YOUR_TOKEN
```

#### Regenerate robot token
```bash
./go-quay get robot regenerate \
  --name deploybot \
  --token YOUR_TOKEN
```

**Note:** The old token will be invalidated immediately.

#### Get robot permissions
```bash
./go-quay get robot permissions \
  --name deploybot \
  --token YOUR_TOKEN
```

### Search API

Search for repositories, users, organizations, and other entities on Quay.io.

📖 **API Reference:** [Search endpoints in Swagger](https://docs.quay.io/api/swagger/#Search)

#### Search for repositories
```bash
# Basic repository search
./go-quay get search repositories \
  --query "nginx" \
  --token YOUR_TOKEN

# Search with pagination
./go-quay get search repositories \
  -q "python" \
  --page 2 \
  -t YOUR_TOKEN
```

#### Search all entity types
```bash
# Search for repositories, users, organizations, teams, and robots
./go-quay get search all \
  --query "redhat" \
  --token YOUR_TOKEN
```

**Entity Types in Results:**
- `repository`: Container image repository
- `user`: User account
- `organization`: Organization
- `team`: Team within an organization
- `robot`: Robot account

### Team API

Manage teams within an organization, including team members and repository permissions.

📖 **API Reference:** [Team endpoints in Swagger](https://docs.quay.io/api/swagger/#Team)

#### List teams in an organization
```bash
./go-quay get team list \
  --organization myorg \
  --token YOUR_TOKEN
```

#### Get team information
```bash
./go-quay get team info \
  --organization myorg \
  --name developers \
  --token YOUR_TOKEN
```

#### Create a new team
```bash
# Create a team with member role
./go-quay get team create \
  --organization myorg \
  --name developers \
  --description "Development team" \
  --role member \
  --token YOUR_TOKEN
```

**Team Roles:**
- `member`: Inherits default permissions
- `creator`: Can create new repositories
- `admin`: Full administrative access

#### Update team settings
```bash
./go-quay get team update \
  --organization myorg \
  --name developers \
  --description "Updated description" \
  --role creator \
  --token YOUR_TOKEN
```

#### Delete a team
```bash
./go-quay get team delete \
  --organization myorg \
  --name developers \
  --confirm \
  --token YOUR_TOKEN
```

#### Manage team members
```bash
# List team members
./go-quay get team members \
  --organization myorg \
  --name developers \
  --token YOUR_TOKEN

# Add a member to a team
./go-quay get team add-member \
  --organization myorg \
  --name developers \
  --member username \
  --token YOUR_TOKEN

# Remove a member from a team
./go-quay get team remove-member \
  --organization myorg \
  --name developers \
  --member username \
  --confirm \
  --token YOUR_TOKEN
```

#### Manage team repository permissions
```bash
# List team repository permissions
./go-quay get team permissions \
  --organization myorg \
  --name developers \
  --token YOUR_TOKEN

# Set repository permission for a team
./go-quay get team set-permission \
  --organization myorg \
  --name developers \
  --repository myrepo \
  --role write \
  --token YOUR_TOKEN

# Remove repository permission from a team
./go-quay get team remove-permission \
  --organization myorg \
  --name developers \
  --repository myrepo \
  --confirm \
  --token YOUR_TOKEN
```

**Permission Roles:**
- `read`: Pull images
- `write`: Pull and push images
- `admin`: Full administrative access

### User API

Manage your user account information and starred repositories.

📖 **API Reference:** [User endpoints in Swagger](https://docs.quay.io/api/swagger/#operation--api-v1-user-get)

#### Get current user information
```bash
./go-quay get user info \
  --token YOUR_TOKEN
```

#### List starred repositories
```bash
./go-quay get user starred \
  --token YOUR_TOKEN
```

#### Star a repository
```bash
./go-quay get user star \
  --namespace quay \
  --repository quay \
  --token YOUR_TOKEN
```

#### Unstar a repository
```bash
./go-quay get user unstar \
  --namespace quay \
  --repository quay \
  --token YOUR_TOKEN
```

### Organization API

The organization API provides comprehensive management of organizations, teams, members, robots, and related settings.

📖 **API Reference:** [Organization endpoints in Swagger](https://docs.quay.io/api/swagger/#operation--api-v1-organization--orgname--get)

#### Get organization information
```bash
./go-quay get organization info \
  --organization ORG_NAME \
  --token YOUR_TOKEN
```

#### Manage organization members
```bash
# List all members
./go-quay get organization members \
  -o myorg \
  -t YOUR_TOKEN
```

#### Manage teams
```bash
# List all teams in an organization
./go-quay get organization teams \
  -o myorg \
  -t YOUR_TOKEN

# Get specific team information
./go-quay get organization team \
  -o myorg \
  --team TEAM_NAME \
  -t YOUR_TOKEN

# Get team members
./go-quay get organization team-members \
  -o myorg \
  --team TEAM_NAME \
  -t YOUR_TOKEN
```

#### Manage robot accounts
```bash
# List organization robot accounts
./go-quay get organization robots \
  -o myorg \
  -t YOUR_TOKEN
```

#### Get organization quota and policies
```bash
# Get organization quota information
./go-quay get organization quota \
  -o myorg \
  -t YOUR_TOKEN

# Get auto-prune policies
./go-quay get organization auto-prune \
  -o myorg \
  -t YOUR_TOKEN

# Get OAuth applications
./go-quay get organization applications \
  -o myorg \
  -t YOUR_TOKEN
```