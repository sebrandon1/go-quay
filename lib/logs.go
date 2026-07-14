/*
Package lib provides Quay.io API client functionality.

This file covers LOGGING endpoints:

Repository Logs:
  - GET /api/v1/repository/{namespace}/{repository}/aggregatelogs  - GetAggregatedLogs()
  - GET /api/v1/repository/{namespace}/{repository}/logs           - GetLogs()

Organization Logs:
  - GET /api/v1/organization/{orgname}/logs                        - GetOrganizationLogs()

All log endpoints support pagination via next_page parameter.
*/
package lib

import (
	"fmt"
	"net/http"
)

const (
	startTimeParam = "starttime"
	endTimeParam   = "endtime"
)

// addLogQueryParams adds optional pagination and date range params to a log request.
func addLogQueryParams(req *http.Request, nextPage, startDate, endDate string) {
	params := map[string]string{}
	if nextPage != "" {
		params["next_page"] = nextPage
	}
	if startDate != "" {
		params[startTimeParam] = startDate
	}
	if endDate != "" {
		params[endTimeParam] = endDate
	}
	if len(params) > 0 {
		addQueryParams(req, params)
	}
}

// GetAggregatedLogs returns the aggregated logs for a repository
func (c *Client) GetAggregatedLogs(namespace, repository, startDate, endDate string) (*AggregatedLogs, error) {
	// Get new request
	req, err := newRequest("GET", c.buildURL("/repository/%s/%s/aggregatelogs", namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get repository aggregate logs request: %w", err)
	}

	// Set query parameters
	addQueryParams(req, map[string]string{
		startTimeParam: startDate,
		endTimeParam:   endDate,
	})

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get repository aggregate logs: %w", err)
	}

	return &logs, nil
}

// GetLogs returns the logs for a repository
func (c *Client) GetLogs(namespace, repository, nextPage, startDate, endDate string) (*Logs, error) {
	req, err := newRequest("GET", c.buildURL("/repository/%s/%s/logs", namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get repository logs request: %w", err)
	}

	addLogQueryParams(req, nextPage, startDate, endDate)

	var logs Logs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get repository logs: %w", err)
	}

	return &logs, nil
}

// GetOrganizationLogs returns the logs for an organization
func (c *Client) GetOrganizationLogs(orgname, nextPage, startDate, endDate string) (*Logs, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/logs", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization logs request: %w", err)
	}

	addLogQueryParams(req, nextPage, startDate, endDate)

	var logs Logs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get organization logs: %w", err)
	}

	return &logs, nil
}

// GetOrganizationAggregatedLogs returns the aggregated logs for an organization
func (c *Client) GetOrganizationAggregatedLogs(orgname, startDate, endDate string) (*AggregatedLogs, error) {
	req, err := newRequest("GET", c.buildURL("/organization/%s/aggregatelogs", orgname), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization aggregate logs request: %w", err)
	}

	addQueryParams(req, map[string]string{
		startTimeParam: startDate,
		endTimeParam:   endDate,
	})

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get organization aggregate logs: %w", err)
	}

	return &logs, nil
}

// ExportOrganizationLogs exports the logs for an organization
func (c *Client) ExportOrganizationLogs(orgname string, request *ExportLogsRequest) error {
	req, err := newRequestWithBody("POST", c.buildURL("/organization/%s/exportlogs", orgname), request)
	if err != nil {
		return fmt.Errorf("failed to create export organization logs request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to export organization logs: %w", err)
	}

	return nil
}

// GetUserLogs returns the logs for the current user
func (c *Client) GetUserLogs(nextPage, startDate, endDate string) (*Logs, error) {
	req, err := newRequest("GET", c.buildURL("/user/logs"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user logs request: %w", err)
	}

	addLogQueryParams(req, nextPage, startDate, endDate)

	var logs Logs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get user logs: %w", err)
	}

	return &logs, nil
}

// GetUserAggregatedLogs returns the aggregated logs for the current user
func (c *Client) GetUserAggregatedLogs(startDate, endDate string) (*AggregatedLogs, error) {
	req, err := newRequest("GET", c.buildURL("/user/aggregatelogs"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get user aggregate logs request: %w", err)
	}

	addQueryParams(req, map[string]string{
		startTimeParam: startDate,
		endTimeParam:   endDate,
	})

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get user aggregate logs: %w", err)
	}

	return &logs, nil
}

// ExportUserLogs exports the logs for the current user
func (c *Client) ExportUserLogs(request *ExportLogsRequest) error {
	req, err := newRequestWithBody("POST", c.buildURL("/user/exportlogs"), request)
	if err != nil {
		return fmt.Errorf("failed to create export user logs request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to export user logs: %w", err)
	}

	return nil
}

// ExportRepositoryLogs exports the logs for a repository
func (c *Client) ExportRepositoryLogs(namespace, repository string, request *ExportLogsRequest) error {
	req, err := newRequestWithBody("POST", c.buildURL("/repository/%s/%s/exportlogs", namespace, repository), request)
	if err != nil {
		return fmt.Errorf("failed to create export repository logs request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to export repository logs: %w", err)
	}

	return nil
}
