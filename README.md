# go-quay
![Quay API Verified Nightly](https://github.com/sebrandon1/go-quay/actions/workflows/nightly.yaml/badge.svg)

A Go wrapper around Quay APIs

## Table of API Coverage

The following APIs are covered by the repo:
| API                    | Cmd     | Lib     | Covered                                                                                                                                                                                                             |
| ---------------------- | ------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Billing                | No      | No      |                                                                                                                                                                                                                     |
| Build                  | No      | No      |                                                                                                                                                                                                                     |
| Discovery              | No      | No      |                                                                                                                                                                                                                     |
| Error                  | No      | No      |                                                                                                                                                                                                                     |
| Messages               | No      | No      |                                                                                                                                                                                                                     |
| Logs                   | Partial | Partial | /api/v1/repository/{repository}/aggregatelogs |
| Manifest               | No      | No      |                                                                                                                                                                                                                     |
| Organization           | No      | No      |                                                                                                                                                                                                                     |
| Permission             | No      | No      |                                                                                                                                                                                                                     |
| Prototype              | No      | No      |                                                                                                                                                                                                                     |
| Repository             | Partial | Partial | /api/v1/repository/{repository}                                                                                                 |
| RepositoryNotification | No      | No      |                                                                                                                                                                                                                     |
| RepoToken              | No      | No      |                                                                                                                                                                                                                     |
| Robot                  | No      | No      |                                                                                                                                                                                                                     |
| Search                 | No      | No      |                                                                                                                                                                                                                     |
| SecScan                | No      | No      |                                                                                                                                                                                                                     |
| Tag                    | Partial | Partial | /api/v1/repository/{repository}/tag                                                                                                                   |
| Team                   | No      | No      |                                                                                                                                                                                                                     |
| Trigger                | No      | No      |                                                                                                                                                                                                                     |
| User                   | No      | No      | 