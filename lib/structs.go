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

// BillingInfo represents billing information for an organization or user.
type BillingInfo struct {
	Plan        string            `json:"plan,omitempty"`
	PlanType    string            `json:"plan_type,omitempty"`
	IsActive    bool              `json:"is_active,omitempty"`
	ExpiresAt   string            `json:"expires_at,omitempty"`
	HasPayment  bool              `json:"has_payment,omitempty"`
	IsTrial     bool              `json:"is_trial,omitempty"`
	TrialUntil  string            `json:"trial_until,omitempty"`
	UsedPrivate int               `json:"used_private,omitempty"`
	PlanPrivate int               `json:"plan_private,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Subscription represents a subscription plan.
type Subscription struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Price       int               `json:"price,omitempty"`
	Currency    string            `json:"currency,omitempty"`
	Period      string            `json:"period,omitempty"`
	Features    []string          `json:"features,omitempty"`
	Limits      map[string]int    `json:"limits,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// PaymentMethod represents a payment method.
type PaymentMethod struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	LastFour    string `json:"last_four,omitempty"`
	Brand       string `json:"brand,omitempty"`
	ExpiryMonth int    `json:"expiry_month,omitempty"`
	ExpiryYear  int    `json:"expiry_year,omitempty"`
	IsDefault   bool   `json:"is_default,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// Invoice represents a billing invoice.
type Invoice struct {
	ID          string `json:"id,omitempty"`
	Number      string `json:"number,omitempty"`
	Status      string `json:"status,omitempty"`
	Amount      int    `json:"amount,omitempty"`
	Currency    string `json:"currency,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	PaidAt      string `json:"paid_at,omitempty"`
	Description string `json:"description,omitempty"`
}

// UsageStats represents usage statistics for billing.
type UsageStats struct {
	PrivateRepos   int    `json:"private_repos,omitempty"`
	PublicRepos    int    `json:"public_repos,omitempty"`
	StorageBytes   int    `json:"storage_bytes,omitempty"`
	BandwidthBytes int    `json:"bandwidth_bytes,omitempty"`
	BuildMinutes   int    `json:"build_minutes,omitempty"`
	Period         string `json:"period,omitempty"`
}
