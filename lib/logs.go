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
