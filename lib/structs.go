/*
Package lib provides Quay.io API client functionality.

This file contains DATA STRUCTURES for all API responses and requests:

Common Types:
  - ResolvedIP, Avatar, Performer, Metadata - Common embedded types

Repository Types:
  - Repository, RepositoryWithTags, RepositoryTags - Repository data structures

Log Types:
  - LogEntry, AggregatedLogEntry, Logs, AggregatedLogs, OrganizationLogs - Logging data

Billing Types:
  - BillingInfo, Subscription, PaymentMethod, Invoice - Billing and subscription data

Organization Types:
  - Organization, OrganizationMember, OrganizationMembers, OrganizationRepository, OrganizationRepositories

Team Types:
  - Team, TeamMember, TeamMembers, TeamPermission, TeamPermissions

Robot Account Types:
  - RobotAccount, RobotAccounts, RobotPermission, RobotPermissions

Application Types:
  - Application, Applications

Quota Types:
  - Quota, QuotaReport

Auto-Prune Types:
  - AutoPrunePolicy, AutoPrunePolicies

Permission Types:
  - DefaultPermission, DefaultPermissions

Proxy Cache Types:
  - ProxyCacheConfig

User Types:
  - User

Request Types:
  - CreateOrganizationRequest, CreateTeamRequest, CreateRobotRequest, CreateApplicationRequest, CreateQuotaRequest, CreateAutoPruneRequest

All structs include appropriate JSON tags for API serialization/deserialization.
*/
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

// Organization represents organization details
type Organization struct {
	Name            string           `json:"name,omitempty"`
	Email           string           `json:"email,omitempty"`
	Avatar          Avatar           `json:"avatar,omitempty"`
	IsOrgAdmin      bool             `json:"is_org_admin,omitempty"`
	CanCreateRepo   bool             `json:"can_create_repo,omitempty"`
	PreferredUsers  bool             `json:"preferred_users,omitempty"`
	Public          bool             `json:"public,omitempty"`
	AutoPrunePolicy *AutoPrunePolicy `json:"auto_prune_policy,omitempty"`
	QuotaReport     *QuotaReport     `json:"quota_report,omitempty"`
	TagExpiration   int              `json:"tag_expiration,omitempty"`
}

// OrganizationMember represents a member of an organization
type OrganizationMember struct {
	Name         string   `json:"name,omitempty"`
	Kind         string   `json:"kind,omitempty"`
	Avatar       Avatar   `json:"avatar,omitempty"`
	Teams        []Team   `json:"teams,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
}

// OrganizationMembers represents the response for organization members
type OrganizationMembers struct {
	Members []OrganizationMember `json:"members,omitempty"`
}

// Team represents a team within an organization
type Team struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Role        string `json:"role,omitempty"`
	Avatar      Avatar `json:"avatar,omitempty"`
	CanView     bool   `json:"can_view,omitempty"`
	RepoCount   int    `json:"repo_count,omitempty"`
	MemberCount int    `json:"member_count,omitempty"`
	IsSynced    bool   `json:"is_synced,omitempty"`
}

// TeamMember represents a member of a team
type TeamMember struct {
	Name    string `json:"name,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Avatar  Avatar `json:"avatar,omitempty"`
	IsRobot bool   `json:"is_robot,omitempty"`
	Invited bool   `json:"invited,omitempty"`
}

// TeamMembers represents the response for team members
type TeamMembers struct {
	Members []TeamMember `json:"members,omitempty"`
}

// TeamPermission represents repository permissions for a team
type TeamPermission struct {
	Repository Repository `json:"repository,omitempty"`
	Role       string     `json:"role,omitempty"`
}

// TeamPermissions represents the response for team permissions
type TeamPermissions struct {
	Permissions []TeamPermission `json:"permissions,omitempty"`
}

// RobotAccount represents a robot account
type RobotAccount struct {
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Token        string                 `json:"token,omitempty"`
	Created      string                 `json:"created,omitempty"`
	LastAccessed string                 `json:"last_accessed,omitempty"`
	Teams        []Team                 `json:"teams,omitempty"`
	Repositories []Repository           `json:"repositories,omitempty"`
	Unstructured map[string]interface{} `json:"unstructured_metadata,omitempty"`
}

// RobotAccounts represents the response for robot accounts
type RobotAccounts struct {
	Robots []RobotAccount `json:"robots,omitempty"`
}

// RobotPermission represents repository permissions for a robot
type RobotPermission struct {
	Repository Repository `json:"repository,omitempty"`
	Role       string     `json:"role,omitempty"`
}

// RobotPermissions represents the response for robot permissions
type RobotPermissions struct {
	Permissions []RobotPermission `json:"permissions,omitempty"`
}

// Application represents an OAuth application
type Application struct {
	ClientID       string `json:"client_id,omitempty"`
	ClientSecret   string `json:"client_secret,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Organization   string `json:"organization,omitempty"`
	RedirectURI    string `json:"redirect_uri,omitempty"`
	ApplicationURI string `json:"application_uri,omitempty"`
	Avatar         Avatar `json:"avatar,omitempty"`
}

// Applications represents the response for applications
type Applications struct {
	Applications []Application `json:"applications,omitempty"`
}

// QuotaReport represents quota usage information
type QuotaReport struct {
	QuotaBytes      int64 `json:"quota_bytes,omitempty"`
	ConfiguredQuota int64 `json:"configured_quota,omitempty"`
	Quota           int64 `json:"quota,omitempty"`
	RunningTotal    int64 `json:"running_total,omitempty"`
}

// Quota represents quota configuration
type Quota struct {
	ID                string `json:"id,omitempty"`
	LimitBytes        int64  `json:"limit_bytes,omitempty"`
	DefaultLimit      int64  `json:"default_limit,omitempty"`
	DefaultLimitBytes int64  `json:"default_limit_bytes,omitempty"`
}

// AutoPrunePolicy represents auto-prune policy configuration
type AutoPrunePolicy struct {
	UUID              string   `json:"uuid,omitempty"`
	Method            string   `json:"method,omitempty"`
	Value             int      `json:"value,omitempty"`
	TagPattern        string   `json:"tag_pattern,omitempty"`
	TagPatternMatches []string `json:"tag_pattern_matches,omitempty"`
	CreationDate      string   `json:"creation_date,omitempty"`
	LastUpdated       string   `json:"last_updated,omitempty"`
}

// AutoPrunePolicies represents the response for auto-prune policies
type AutoPrunePolicies struct {
	Policies []AutoPrunePolicy `json:"policies,omitempty"`
}

// ProxyCacheConfig represents proxy cache configuration
type ProxyCacheConfig struct {
	UpstreamRegistry string `json:"upstream_registry,omitempty"`
	Insecure         bool   `json:"insecure,omitempty"`
	Expiration       int    `json:"expiration,omitempty"`
}

// DefaultPermission represents default repository permissions
type DefaultPermission struct {
	Role      string `json:"role,omitempty"`
	Delegate  User   `json:"delegate,omitempty"`
	AppliedTo User   `json:"applied_to,omitempty"`
	ID        string `json:"id,omitempty"`
}

// DefaultPermissions represents the response for default permissions
type DefaultPermissions struct {
	Prototypes []DefaultPermission `json:"prototypes,omitempty"`
}

// User represents a user account
type User struct {
	Name        string `json:"name,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Avatar      Avatar `json:"avatar,omitempty"`
	Kind        string `json:"kind,omitempty"`
	IsRobot     bool   `json:"is_robot,omitempty"`
	IsOrgMember bool   `json:"is_org_member,omitempty"`
}

// OrganizationRepository represents a repository within an organization
type OrganizationRepository struct {
	Name         string  `json:"name,omitempty"`
	Description  string  `json:"description,omitempty"`
	IsPublic     bool    `json:"is_public,omitempty"`
	Kind         string  `json:"kind,omitempty"`
	Namespace    string  `json:"namespace,omitempty"`
	LastModified string  `json:"last_modified,omitempty"`
	Popularity   float64 `json:"popularity,omitempty"`
	TagsCount    int     `json:"tags_count,omitempty"`
}

// OrganizationRepositories represents the response for organization repositories
type OrganizationRepositories struct {
	Repositories []OrganizationRepository `json:"repositories,omitempty"`
}

// CreateOrganizationRequest represents the request to create an organization
type CreateOrganizationRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateTeamRequest represents the request to create a team
type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Role        string `json:"role,omitempty"`
}

// CreateRobotRequest represents the request to create a robot account
type CreateRobotRequest struct {
	Description  string                 `json:"description,omitempty"`
	Unstructured map[string]interface{} `json:"unstructured_metadata,omitempty"`
}

// CreateApplicationRequest represents the request to create an application
type CreateApplicationRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	ApplicationURI string `json:"application_uri,omitempty"`
	RedirectURI    string `json:"redirect_uri,omitempty"`
}

// CreateQuotaRequest represents the request to create a quota
type CreateQuotaRequest struct {
	LimitBytes int64 `json:"limit_bytes"`
}

// CreateAutoPruneRequest represents the request to create an auto-prune policy
type CreateAutoPruneRequest struct {
	Method     string `json:"method"`
	Value      int    `json:"value"`
	TagPattern string `json:"tag_pattern,omitempty"`
}
