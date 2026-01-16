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

var (
	QuayURL = "https://quay.io/api/v1"
)

// addQueryParams adds query parameters to a request URL.
func addQueryParams(req *http.Request, params map[string]string) {
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
}

// GetAggregatedLogs returns the aggregated logs for a repository
func (c *Client) GetAggregatedLogs(namespace, repository, startDate, endDate string) (*AggregatedLogs, error) {
	// Get new request
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/aggregatelogs", QuayURL, namespace, repository), nil)
	if err != nil {
		return nil, err
	}

	// Set query parameters
	addQueryParams(req, map[string]string{
		"starttime": startDate,
		"endtime":   endDate,
	})

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, err
	}

	return &logs, nil
}

// GetLogs returns the logs for a repository
func (c *Client) GetLogs(namespace, repository, nextPage string) (*Logs, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/logs", QuayURL, namespace, repository), nil)
	if err != nil {
		return nil, err
	}

	// Set query parameters
	if nextPage != "" {
		addQueryParams(req, map[string]string{"next_page": nextPage})
	}

	var logs Logs
	if err := c.get(req, &logs); err != nil {
		return nil, err
	}

	return &logs, nil
}

// GetOrganizationLogs returns the logs for an organization
func (c *Client) GetOrganizationLogs(orgname, nextPage string) (*OrganizationLogs, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/logs", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.BearerToken)

	if nextPage != "" {
		addQueryParams(req, map[string]string{"next_page": nextPage})
	}

	var logs OrganizationLogs
	if err := c.get(req, &logs); err != nil {
		return nil, err
	}

	return &logs, nil
}

// GetOrganizationAggregatedLogs returns the aggregated logs for an organization
func (c *Client) GetOrganizationAggregatedLogs(orgname, startDate, endDate string) (*AggregatedLogs, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/organization/%s/aggregatelogs", QuayURL, orgname), nil)
	if err != nil {
		return nil, err
	}

	addQueryParams(req, map[string]string{
		"starttime": startDate,
		"endtime":   endDate,
	})

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, err
	}

	return &logs, nil
}

// ExportOrganizationLogs exports the logs for an organization
func (c *Client) ExportOrganizationLogs(orgname string, request *ExportLogsRequest) error {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/organization/%s/exportlogs", QuayURL, orgname), request)
	if err != nil {
		return err
	}

	if err := c.post(req, nil); err != nil {
		return err
	}

	return nil
}

// GetUserLogs returns the logs for the current user
func (c *Client) GetUserLogs(nextPage string) (*Logs, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/user/logs", QuayURL), nil)
	if err != nil {
		return nil, err
	}

	if nextPage != "" {
		addQueryParams(req, map[string]string{"next_page": nextPage})
	}

	var logs Logs
	if err := c.get(req, &logs); err != nil {
		return nil, err
	}

	return &logs, nil
}

// GetUserAggregatedLogs returns the aggregated logs for the current user
func (c *Client) GetUserAggregatedLogs(startDate, endDate string) (*AggregatedLogs, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/user/aggregatelogs", QuayURL), nil)
	if err != nil {
		return nil, err
	}

	addQueryParams(req, map[string]string{
		"starttime": startDate,
		"endtime":   endDate,
	})

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, err
	}

	return &logs, nil
}

// ExportUserLogs exports the logs for the current user
func (c *Client) ExportUserLogs(request *ExportLogsRequest) error {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/user/exportlogs", QuayURL), request)
	if err != nil {
		return err
	}

	if err := c.post(req, nil); err != nil {
		return err
	}

	return nil
}

// ExportRepositoryLogs exports the logs for a repository
func (c *Client) ExportRepositoryLogs(namespace, repository string, request *ExportLogsRequest) error {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/repository/%s/%s/exportlogs", QuayURL, namespace, repository), request)
	if err != nil {
		return err
	}

	if err := c.post(req, nil); err != nil {
		return err
	}

	return nil
}
