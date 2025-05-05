package lib

// ResolvedIP represents resolved IP details.
type ResolvedIP struct {
	Provider       string `json:"provider,omitempty"`
	Service        string `json:"service,omitempty"`
	SyncToken      string `json:"sync_token,omitempty"`
	CountryIsoCode string `json:"country_iso_code,omitempty"`
	AwsRegion      any    `json:"aws_region,omitempty"`
	Continent      string `json:"continent,omitempty"`
}

// Avatar represents performer avatar details.
type Avatar struct {
	Name  string `json:"name,omitempty"`
	Hash  string `json:"hash,omitempty"`
	Color string `json:"color,omitempty"`
	Kind  string `json:"kind,omitempty"`
}

// Performer represents the performer of an action.
type Performer struct {
	Kind    string `json:"kind,omitempty"`
	Name    string `json:"name,omitempty"`
	IsRobot bool   `json:"is_robot,omitempty"`
	Avatar  Avatar `json:"avatar,omitempty"`
}

// Metadata represents metadata for logs.
type Metadata struct {
	Repo           string     `json:"repo,omitempty"`
	Namespace      string     `json:"namespace,omitempty"`
	UserAgent      string     `json:"user-agent,omitempty"`
	ManifestDigest string     `json:"manifest_digest,omitempty"`
	Username       string     `json:"username,omitempty"`
	IsRobot        bool       `json:"is_robot,omitempty"`
	ResolvedIP     ResolvedIP `json:"resolved_ip,omitempty"`
	Tag            string     `json:"tag,omitempty"`
}

// LogEntry represents a single log entry.
type LogEntry struct {
	Kind      string    `json:"kind,omitempty"`
	Metadata  Metadata  `json:"metadata,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Datetime  string    `json:"datetime,omitempty"`
	Performer Performer `json:"performer,omitempty"`
}

// AggregatedLogEntry represents a single aggregated log entry.
type AggregatedLogEntry struct {
	Kind     string `json:"kind"`
	Count    int    `json:"count"`
	Datetime string `json:"datetime"`
}

// AggregatedLogs represents aggregated logs for a repository.
type AggregatedLogs struct {
	Aggregated []AggregatedLogEntry `json:"aggregated"`
}

// Logs represents logs for a repository.
type Logs struct {
	StartTime string     `json:"start_time,omitempty"`
	EndTime   string     `json:"end_time,omitempty"`
	Logs      []LogEntry `json:"logs,omitempty"`
	NextPage  string     `json:"next_page,omitempty"`
}

// OrganizationLogs represents logs for an organization.
type OrganizationLogs struct {
	StartTime string     `json:"start_time,omitempty"`
	EndTime   string     `json:"end_time,omitempty"`
	Logs      []LogEntry `json:"logs,omitempty"`
	NextPage  string     `json:"next_page,omitempty"`
}
