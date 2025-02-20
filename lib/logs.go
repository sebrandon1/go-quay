package lib

import (
	"net/url"
)

const (
	QuayURL = "https://quay.io/api/v1"
)

type AggregatedLogs struct {
	Aggregated []struct {
		Kind     string `json:"kind"`
		Count    int    `json:"count"`
		Datetime string `json:"datetime"`
	} `json:"aggregated"`
}

type Logs struct {
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Logs      []struct {
		Kind     string `json:"kind,omitempty"`
		Metadata struct {
			Repo           string `json:"repo,omitempty"`
			Namespace      string `json:"namespace,omitempty"`
			UserAgent      string `json:"user-agent,omitempty"`
			ManifestDigest string `json:"manifest_digest,omitempty"`
			Username       string `json:"username,omitempty"`
			IsRobot        bool   `json:"is_robot,omitempty"`
			ResolvedIP     struct {
				Provider       string `json:"provider,omitempty"`
				Service        string `json:"service,omitempty"`
				SyncToken      string `json:"sync_token,omitempty"`
				CountryIsoCode string `json:"country_iso_code,omitempty"`
				AwsRegion      any    `json:"aws_region,omitempty"`
				Continent      string `json:"continent,omitempty"`
			} `json:"resolved_ip,omitempty"`
		} `json:"metadata,omitempty"`
		IP        string `json:"ip,omitempty"`
		Datetime  string `json:"datetime,omitempty"`
		Performer struct {
			Kind    string `json:"kind,omitempty"`
			Name    string `json:"name,omitempty"`
			IsRobot bool   `json:"is_robot,omitempty"`
			Avatar  struct {
				Name  string `json:"name,omitempty"`
				Hash  string `json:"hash,omitempty"`
				Color string `json:"color,omitempty"`
				Kind  string `json:"kind,omitempty"`
			} `json:"avatar,omitempty"`
		} `json:"performer,omitempty"`
	} `json:"logs,omitempty"`
	NextPage string `json:"next_page,omitempty"`
}

type OrganizationLogs struct {
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Logs      []struct {
		Kind     string `json:"kind,omitempty"`
		Metadata struct {
			Repo       string `json:"repo,omitempty"`
			Namespace  string `json:"namespace,omitempty"`
			UserAgent  string `json:"user-agent,omitempty"`
			Tag        string `json:"tag,omitempty"`
			Username   string `json:"username,omitempty"`
			IsRobot    bool   `json:"is_robot,omitempty"`
			ResolvedIP struct {
				Provider       string `json:"provider,omitempty"`
				Service        string `json:"service,omitempty"`
				SyncToken      string `json:"sync_token,omitempty"`
				CountryIsoCode string `json:"country_iso_code,omitempty"`
				AwsRegion      any    `json:"aws_region,omitempty"`
				Continent      string `json:"continent,omitempty"`
			} `json:"resolved_ip,omitempty"`
		} `json:"metadata,omitempty"`
		IP        string `json:"ip,omitempty"`
		Datetime  string `json:"datetime,omitempty"`
		Performer struct {
			Kind    string `json:"kind,omitempty"`
			Name    string `json:"name,omitempty"`
			IsRobot bool   `json:"is_robot,omitempty"`
			Avatar  struct {
				Name  string `json:"name,omitempty"`
				Hash  string `json:"hash,omitempty"`
				Color string `json:"color,omitempty"`
				Kind  string `json:"kind,omitempty"`
			} `json:"avatar,omitempty"`
		} `json:"performer,omitempty"`
	} `json:"logs,omitempty"`
	NextPage string `json:"next_page,omitempty"`
}

// GetAggregatedLogs returns the aggregated logs for a repository
func (c *Client) GetAggregatedLogs(namespace, repository, startDate, endDate string) (*AggregatedLogs, error) {
	// Get new request
	req, err := newRequest("GET", QuayURL+"/repository/"+namespace+"/"+repository+"/aggregatelogs", nil)
	if err != nil {
		return nil, err
	}

	// Set the bearer token
	req.Header.Add("Authorization", "Bearer "+c.BearerToken)

	// set the query parameters for starttime and endtime
	q := req.URL.Query()
	q.Add("starttime", startDate)
	q.Add("endtime", endDate)

	decoded, err := url.QueryUnescape(q.Encode())
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = decoded

	var logs AggregatedLogs
	err = c.get(req, &logs)
	if err != nil {
		return nil, err
	}
	return &logs, nil
}

// GetLogs returns the logs for a repository
func (c *Client) GetLogs(namespace, repository, next_page string) (*Logs, error) {
	// Get new request
	req, err := newRequest("GET", QuayURL+"/repository/"+namespace+"/"+repository+"/logs", nil)
	if err != nil {
		return nil, err
	}

	// Set the bearer token
	req.Header.Add("Authorization", "Bearer "+c.BearerToken)

	// set the query parameters for next_page
	if next_page != "" {
		q := req.URL.Query()
		q.Add("next_page", next_page)
		req.URL.RawQuery = q.Encode()
	}

	var logs Logs
	err = c.get(req, &logs)
	if err != nil {
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

	// set the query parameters for next_page
	if next_page != "" {
		q := req.URL.Query()
		q.Add("next_page", next_page)
		req.URL.RawQuery = q.Encode()
	}

	var logs OrganizationLogs
	err = c.get(req, &logs)
	if err != nil {
		return nil, err
	}
	return &logs, nil
}
