/*
Package lib provides Quay.io API client functionality.

This file covers LOGGING endpoints:

Repository Logs:
  - GET /api/v1/repository/{namespace}/{repository}/aggregatelogs  - GetAggregatedLogs(), GetAggregatedLogsRange()
  - GET /api/v1/repository/{namespace}/{repository}/logs           - GetLogs(), GetAllLogs()
  - POST /api/v1/repository/{namespace}/{repository}/exportlogs    - ExportRepositoryLogs()

Organization Logs:
  - GET /api/v1/organization/{orgname}/logs                        - GetOrganizationLogs()
  - GET /api/v1/organization/{orgname}/aggregatelogs               - GetOrganizationAggregatedLogs()
  - POST /api/v1/organization/{orgname}/exportlogs                 - ExportOrganizationLogs()

User Logs:
  - GET /api/v1/user/logs                                          - GetUserLogs()
  - GET /api/v1/user/aggregatelogs                                 - GetUserAggregatedLogs()
  - POST /api/v1/user/exportlogs                                   - ExportUserLogs()

Log list endpoints support pagination via next_page. The Quay API rejects
aggregated-log ranges longer than 30 days and silently ignores dates that are
not MM/DD/YYYY, so callers may pass either MM/DD/YYYY or YYYY-MM-DD and use
GetAggregatedLogsRange for longer ranges.
*/
package lib

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	startTimeParam = "starttime"
	endTimeParam   = "endtime"

	// logDateLayout is the only date format the Quay log API honors.
	logDateLayout = "01/02/2006"
	// isoDateLayout is accepted from callers and converted to logDateLayout.
	isoDateLayout = "2006-01-02"

	// maxAggregatedLogWindowDays is the longest range the Quay API accepts for
	// aggregated logs; longer ranges fail with "Invalid date range for logs".
	maxAggregatedLogWindowDays = 30
)

// normalizeLogDate converts a caller-supplied date into the MM/DD/YYYY form the
// Quay log API requires. An unparseable value is an error rather than a silent
// fallback: Quay ignores bad dates and answers with a 48-hour default window,
// which looks like valid data.
func normalizeLogDate(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	for _, layout := range []string{logDateLayout, isoDateLayout, time.RFC3339} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.Format(logDateLayout), nil
		}
	}
	return "", fmt.Errorf("invalid log date %q: use MM/DD/YYYY or YYYY-MM-DD", value)
}

// ValidateLogDateRange normalizes a log date range and reports whether Quay
// will accept it, without sending a request. Empty values are allowed, matching
// the API's optional range.
func ValidateLogDateRange(startDate, endDate string) (string, string, error) {
	return normalizeLogRange(startDate, endDate)
}

// normalizeLogRange normalizes both ends of an optional date range.
func normalizeLogRange(startDate, endDate string) (string, string, error) {
	start, err := normalizeLogDate(startDate)
	if err != nil {
		return "", "", err
	}
	end, err := normalizeLogDate(endDate)
	if err != nil {
		return "", "", err
	}
	return start, end, nil
}

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

// GetAggregatedLogs returns the aggregated logs for a repository. The Quay API
// accepts at most 30 days per call; use GetAggregatedLogsRange for longer spans.
func (c *Client) GetAggregatedLogs(ctx context.Context, namespace, repository, startDate, endDate string) (*AggregatedLogs, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	start, end, err := normalizeLogRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get new request
	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/aggregatelogs", namespace, repository))
	if err != nil {
		return nil, fmt.Errorf("failed to create get repository aggregate logs request: %w", err)
	}

	addLogQueryParams(req, "", start, end)

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get repository aggregate logs: %w", err)
	}

	return &logs, nil
}

// GetLogs returns the logs for a repository
func (c *Client) GetLogs(ctx context.Context, namespace, repository, nextPage, startDate, endDate string) (*Logs, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	start, end, err := normalizeLogRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/logs", namespace, repository))
	if err != nil {
		return nil, fmt.Errorf("failed to create get repository logs request: %w", err)
	}

	addLogQueryParams(req, nextPage, start, end)

	var logs Logs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get repository logs: %w", err)
	}

	return &logs, nil
}

// GetOrganizationLogs returns the logs for an organization
func (c *Client) GetOrganizationLogs(ctx context.Context, orgname, nextPage, startDate, endDate string) (*Logs, error) {
	if orgname == "" {
		return nil, fmt.Errorf("orgname is required")
	}
	start, end, err := normalizeLogRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/organization/%s/logs", orgname))
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization logs request: %w", err)
	}

	addLogQueryParams(req, nextPage, start, end)

	var logs Logs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get organization logs: %w", err)
	}

	return &logs, nil
}

// GetOrganizationAggregatedLogs returns the aggregated logs for an organization
func (c *Client) GetOrganizationAggregatedLogs(ctx context.Context, orgname, startDate, endDate string) (*AggregatedLogs, error) {
	if orgname == "" {
		return nil, fmt.Errorf("orgname is required")
	}
	start, end, err := normalizeLogRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/organization/%s/aggregatelogs", orgname))
	if err != nil {
		return nil, fmt.Errorf("failed to create get organization aggregate logs request: %w", err)
	}

	addLogQueryParams(req, "", start, end)

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get organization aggregate logs: %w", err)
	}

	return &logs, nil
}

// ExportOrganizationLogs exports the logs for an organization
func (c *Client) ExportOrganizationLogs(ctx context.Context, orgname string, request *ExportLogsRequest) error {
	if orgname == "" {
		return fmt.Errorf("orgname is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/organization/%s/exportlogs", orgname), request)
	if err != nil {
		return fmt.Errorf("failed to create export organization logs request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to export organization logs: %w", err)
	}

	return nil
}

// GetUserLogs returns the logs for the current user
func (c *Client) GetUserLogs(ctx context.Context, nextPage, startDate, endDate string) (*Logs, error) {
	start, end, err := normalizeLogRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/user/logs"))
	if err != nil {
		return nil, fmt.Errorf("failed to create get user logs request: %w", err)
	}

	addLogQueryParams(req, nextPage, start, end)

	var logs Logs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get user logs: %w", err)
	}

	return &logs, nil
}

// GetUserAggregatedLogs returns the aggregated logs for the current user
func (c *Client) GetUserAggregatedLogs(ctx context.Context, startDate, endDate string) (*AggregatedLogs, error) {
	start, end, err := normalizeLogRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/user/aggregatelogs"))
	if err != nil {
		return nil, fmt.Errorf("failed to create get user aggregate logs request: %w", err)
	}

	addLogQueryParams(req, "", start, end)

	var logs AggregatedLogs
	if err := c.get(req, &logs); err != nil {
		return nil, fmt.Errorf("failed to get user aggregate logs: %w", err)
	}

	return &logs, nil
}

// ExportUserLogs exports the logs for the current user
func (c *Client) ExportUserLogs(ctx context.Context, request *ExportLogsRequest) error {
	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/user/exportlogs"), request)
	if err != nil {
		return fmt.Errorf("failed to create export user logs request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to export user logs: %w", err)
	}

	return nil
}

// ExportRepositoryLogs exports the logs for a repository
func (c *Client) ExportRepositoryLogs(ctx context.Context, namespace, repository string, request *ExportLogsRequest) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/repository/%s/%s/exportlogs", namespace, repository), request)
	if err != nil {
		return fmt.Errorf("failed to create export repository logs request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to export repository logs: %w", err)
	}

	return nil
}

// GetAggregatedLogsRange returns the aggregated logs for a repository across an
// arbitrary date range. The Quay API accepts at most 30 days per request, so
// longer ranges are split into consecutive windows and the entries are merged.
// Passing an empty start or end date makes a single request, the way
// GetAggregatedLogs does.
func (c *Client) GetAggregatedLogsRange(ctx context.Context, namespace, repository, startDate, endDate string) (*AggregatedLogs, error) {
	start, end, err := normalizeLogRange(startDate, endDate)
	if err != nil {
		return nil, err
	}
	if start == "" || end == "" {
		return c.GetAggregatedLogs(ctx, namespace, repository, start, end)
	}

	rangeStart, err := time.Parse(logDateLayout, start)
	if err != nil {
		return nil, fmt.Errorf("invalid start date %q: %w", start, err)
	}
	rangeEnd, err := time.Parse(logDateLayout, end)
	if err != nil {
		return nil, fmt.Errorf("invalid end date %q: %w", end, err)
	}
	if rangeEnd.Before(rangeStart) {
		return nil, fmt.Errorf("end date %s is before start date %s", end, start)
	}

	combined := &AggregatedLogs{}
	for windowStart := rangeStart; !windowStart.After(rangeEnd); {
		windowEnd := windowStart.AddDate(0, 0, maxAggregatedLogWindowDays-1)
		if windowEnd.After(rangeEnd) {
			windowEnd = rangeEnd
		}
		part, err := c.GetAggregatedLogs(ctx, namespace, repository,
			windowStart.Format(logDateLayout), windowEnd.Format(logDateLayout))
		if err != nil {
			return nil, err
		}
		combined.Aggregated = append(combined.Aggregated, part.Aggregated...)
		windowStart = windowEnd.AddDate(0, 0, 1)
	}
	return combined, nil
}

// GetAllLogs follows next_page for a repository until the log stream ends or
// maxPages pages have been read. It returns the entries collected so far and
// whether the result was truncated by the page cap.
func (c *Client) GetAllLogs(ctx context.Context, namespace, repository, startDate, endDate string, maxPages int) ([]LogEntry, bool, error) {
	entries := []LogEntry{}
	nextPage := ""
	for page := 0; maxPages <= 0 || page < maxPages; page++ {
		logs, err := c.GetLogs(ctx, namespace, repository, nextPage, startDate, endDate)
		if err != nil {
			return nil, false, err
		}
		entries = append(entries, logs.Logs...)
		nextPage = logs.NextPage
		if nextPage == "" {
			return entries, false, nil
		}
	}
	return entries, true, nil
}
