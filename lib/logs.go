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
func (c *Client) GetOrganizationLogs(namespace, next_page string) (*OrganizationLogs, error) {
	// Get new request
	req, err := newRequest("GET", QuayURL+"/organization/"+namespace+"/logs", nil)
	if err != nil {
		return nil, err
	}

	// Set the bearer token
	req.Header.Add("Authorization", "Bearer "+c.BearerToken)

	// Set query parameters
	if next_page != "" {
		addQueryParams(req, map[string]string{"next_page": next_page})
	}

	var logs OrganizationLogs
	err = c.get(req, &logs)
	if err != nil {
		return nil, err
	}

	return &logs, nil
}
