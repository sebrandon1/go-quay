# go-quay

[![Pre-Main Checks](https://github.com/sebrandon1/go-quay/actions/workflows/pre-main.yaml/badge.svg)](https://github.com/sebrandon1/go-quay/actions/workflows/pre-main.yaml)
[![Quay API Verified Nightly](https://github.com/sebrandon1/go-quay/actions/workflows/nightly.yaml/badge.svg)](https://github.com/sebrandon1/go-quay/actions/workflows/nightly.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sebrandon1/go-quay)](https://golang.org/)
[![License](https://img.shields.io/github/license/sebrandon1/go-quay)](https://github.com/sebrandon1/go-quay/blob/main/LICENSE)

A Go wrapper around Quay APIs

## Table of API Coverage

The following APIs are covered by the repo:
| API                    | Cmd     | Lib     | Covered                                                                                                                                                                                                             |
| ---------------------- | ------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Billing                | Yes     | Yes     | /api/v1/user/plan, /api/v1/organization/{orgname}/plan, /api/v1/organization/{orgname}/invoices, /api/v1/plans/                                                                                                   |
| Build                  | No      | No      |                                                                                                                                                                                                                     |
| Discovery              | No      | No      |                                                                                                                                                                                                                     |
| Error                  | No      | No      |                                                                                                                                                                                                                     |
| Messages               | No      | No      |                                                                                                                                                                                                                     |
| Logs                   | Partial | Partial | /api/v1/repository/{namespace}/{repository}/aggregatelogs, /api/v1/repository/{namespace}/{repository}/logs, /api/v1/organization/{orgname}/logs |
| Manifest               | No      | No      |                                                                                                                                                                                                                     |
| Organization           | Yes     | Yes     | /api/v1/organization/{orgname}, /api/v1/organization/{orgname}/members, /api/v1/organization/{orgname}/teams, /api/v1/organization/{orgname}/team/{teamname}, /api/v1/organization/{orgname}/robots, /api/v1/organization/{orgname}/quota, /api/v1/organization/{orgname}/autoprunepolicy, /api/v1/organization/{orgname}/applications |
| Permission             | No      | No      |                                                                                                                                                                                                                     |
| Prototype              | No      | No      |                                                                                                                                                                                                                     |
| Repository             | Partial | Partial | /api/v1/repository/{namespace}/{repository}, /api/v1/repository/{namespace}/{repository}/tag                                   |
| RepositoryNotification | No      | No      |                                                                                                                                                                                                                     |
| RepoToken              | No      | No      |                                                                                                                                                                                                                     |
| Robot                  | No      | No      |                                                                                                                                                                                                                     |
| Search                 | No      | No      |                                                                                                                                                                                                                     |
| SecScan                | No      | No      |                                                                                                                                                                                                                     |
| Tag                    | Partial | Partial | /api/v1/repository/{namespace}/{repository}/tag (included in Repository API)                                                                       |
| Team                   | No      | No      |                                                                                                                                                                                                                     |
| Trigger                | No      | No      |                                                                                                                                                                                                                     |
| User                   | No      | No      | 

## Authentication

All API commands require a Quay.io authentication token. You can obtain a token from your Quay.io account settings:

1. Go to [Quay.io](https://quay.io) and log in to your account
2. Navigate to **Account Settings** → **Robot Accounts** or **User Settings** → **CLI Password**
3. Generate a new token with appropriate permissions
4. Use the token with the `--token` or `-t` flag in commands

## Usage Examples

### Billing API

The billing API provides access to subscription plans, billing information, and invoices.

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

### Logs API

The logs API provides access to repository activity logs and aggregated statistics.

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

The repository API provides access to repository information including metadata and tags.

#### Get repository information with tags
```bash
./go-quay get repository \
  --namespace NAMESPACE \
  --repository REPOSITORY \
  --token YOUR_TOKEN
```

#### Example
```bash
./go-quay get repository \
  -n myorg \
  -r myapp \
  -t YOUR_TOKEN
```

This returns comprehensive repository information including:
- Repository metadata (description, visibility, permissions)
- All repository tags with details (size, last modified, expiration)
- Organization and user permissions

### Organization API

The organization API provides comprehensive management of organizations, teams, members, robots, and related settings.

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