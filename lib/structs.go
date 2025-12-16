/*
Package lib provides Quay.io API client functionality.

This file contains DATA STRUCTURES for all API responses and requests:

Common Types:
  - ResolvedIP, Avatar, Performer, Metadata - Common embedded types

Repository Types:
  - Repository, RepositoryWithTags, RepositoryTags - Repository data structures
  - CreateRepositoryRequest, UpdateRepositoryRequest - Repository management requests

Repository Permission Types:
  - RepositoryPermission, RepositoryPermissions, SetRepositoryPermissionRequest - Repository permissions

Enhanced Tag Types:
  - Tag, TagHistory, UpdateTagRequest - Enhanced tag operations with metadata

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
  - User, UserDetails - Basic and detailed user information
  - StarredRepository, StarredRepositories - User starred repositories

Search Types:
  - SearchResult - Repository search results

Error Types:
  - QuayError - API error responses

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

// Repository Management Structures

// CreateRepositoryRequest represents the request to create a repository
type CreateRepositoryRequest struct {
	Repository  string `json:"repository"`
	Namespace   string `json:"namespace,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateRepositoryRequest represents the request to update a repository
type UpdateRepositoryRequest struct {
	Description string `json:"description,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

// Repository Permissions Structures

// RepositoryPermission represents a permission for a repository
type RepositoryPermission struct {
	Name       string `json:"name,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Avatar     Avatar `json:"avatar,omitempty"`
	Role       string `json:"role,omitempty"`
	IsRobot    bool   `json:"is_robot,omitempty"`
	IsOrgAdmin bool   `json:"is_org_admin,omitempty"`
}

// RepositoryPermissions represents the response for repository permissions
type RepositoryPermissions struct {
	Permissions []RepositoryPermission `json:"permissions,omitempty"`
}

// SetRepositoryPermissionRequest represents the request to set repository permission
type SetRepositoryPermissionRequest struct {
	Role string `json:"role"`
}

// Enhanced Tag Structures

// Tag represents a repository tag with enhanced metadata
type Tag struct {
	Name           string            `json:"name,omitempty"`
	Reversion      bool              `json:"reversion,omitempty"`
	StartTs        int64             `json:"start_ts,omitempty"`
	ManifestDigest string            `json:"manifest_digest,omitempty"`
	IsManifestList bool              `json:"is_manifest_list,omitempty"`
	Size           int64             `json:"size,omitempty"`
	LastModified   string            `json:"last_modified,omitempty"`
	EndTs          int64             `json:"end_ts,omitempty"`
	Expiration     string            `json:"expiration,omitempty"`
	DockerImageID  string            `json:"docker_image_id,omitempty"`
	ImageID        string            `json:"image_id,omitempty"`
	V1Metadata     map[string]string `json:"v1_metadata,omitempty"`
}

// TagHistory represents the history of a tag
type TagHistory struct {
	Tags []Tag `json:"tags,omitempty"`
}

// UpdateTagRequest represents the request to update a tag
type UpdateTagRequest struct {
	Expiration string `json:"expiration,omitempty"`
}

// User Account Structures

// UserDetails represents detailed user information
type UserDetails struct {
	Anonymous      bool   `json:"anonymous,omitempty"`
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	Verified       bool   `json:"verified,omitempty"`
	Avatar         Avatar `json:"avatar,omitempty"`
	Organizations  []User `json:"organizations,omitempty"`
	CanCreateRepo  bool   `json:"can_create_repo,omitempty"`
	PreferredUsers bool   `json:"preferred_users,omitempty"`
	TagExpirationS int    `json:"tag_expiration_s,omitempty"`
}

// StarredRepository represents a starred repository
type StarredRepository struct {
	Namespace    string  `json:"namespace,omitempty"`
	Name         string  `json:"name,omitempty"`
	Description  string  `json:"description,omitempty"`
	IsPublic     bool    `json:"is_public,omitempty"`
	Kind         string  `json:"kind,omitempty"`
	LastModified string  `json:"last_modified,omitempty"`
	Popularity   float64 `json:"popularity,omitempty"`
}

// StarredRepositories represents the response for starred repositories
type StarredRepositories struct {
	Repositories []StarredRepository `json:"repositories,omitempty"`
}

// Search and Discovery Structures

// SearchRepository represents a repository in search results
type SearchRepository struct {
	Namespace   string  `json:"namespace,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	IsPublic    bool    `json:"is_public,omitempty"`
	Kind        string  `json:"kind,omitempty"`
	Popularity  float64 `json:"popularity,omitempty"`
	Score       float64 `json:"score,omitempty"`
	Href        string  `json:"href,omitempty"`
}

// SearchRepositoryResult represents the response for repository search
type SearchRepositoryResult struct {
	Results       []SearchRepository `json:"results,omitempty"`
	HasAdditional bool               `json:"has_additional,omitempty"`
	Page          int                `json:"page,omitempty"`
	StartIndex    int                `json:"start_index,omitempty"`
}

// SearchEntity represents a generic entity in search results (user, org, team, robot)
type SearchEntity struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Kind        string  `json:"kind,omitempty"`
	Score       float64 `json:"score,omitempty"`
	Href        string  `json:"href,omitempty"`
	Avatar      Avatar  `json:"avatar,omitempty"`
}

// SearchAllResult represents the response for searching all entities
type SearchAllResult struct {
	Results []SearchEntity `json:"results,omitempty"`
}

// SearchResult represents a repository search result (legacy)
type SearchResult struct {
	Repositories []struct {
		Namespace   string `json:"namespace,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		IsPublic    bool   `json:"is_public,omitempty"`
		Popularity  int    `json:"popularity,omitempty"`
	} `json:"repositories,omitempty"`
}

// Manifest Structures

// Manifest represents a container image manifest
type Manifest struct {
	Digest         string          `json:"digest,omitempty"`
	SchemaVersion  int             `json:"schemaVersion,omitempty"`
	MediaType      string          `json:"mediaType,omitempty"`
	Size           int64           `json:"size,omitempty"`
	Layers         []ManifestLayer `json:"layers,omitempty"`
	Config         ManifestConfig  `json:"config,omitempty"`
	IsManifestList bool            `json:"is_manifest_list,omitempty"`
	ManifestData   string          `json:"manifest_data,omitempty"`
}

// ManifestLayer represents a layer in a manifest
type ManifestLayer struct {
	MediaType string `json:"mediaType,omitempty"`
	Size      int64  `json:"size,omitempty"`
	Digest    string `json:"digest,omitempty"`
	Index     int    `json:"index,omitempty"`
	Command   string `json:"command,omitempty"`
}

// ManifestConfig represents the config of a manifest
type ManifestConfig struct {
	MediaType string `json:"mediaType,omitempty"`
	Size      int64  `json:"size,omitempty"`
	Digest    string `json:"digest,omitempty"`
}

// ManifestLabel represents a label on a manifest
type ManifestLabel struct {
	ID         string `json:"id,omitempty"`
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
	SourceType string `json:"source_type,omitempty"`
	MediaType  string `json:"media_type,omitempty"`
}

// ManifestLabels represents the response for manifest labels
type ManifestLabels struct {
	Labels []ManifestLabel `json:"labels,omitempty"`
}

// AddManifestLabelRequest represents the request to add a label to a manifest
type AddManifestLabelRequest struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	MediaType string `json:"media_type,omitempty"`
}

// Security Scan Structures

// SecurityScan represents the security scan result for an image
type SecurityScan struct {
	Status string        `json:"status,omitempty"`
	Data   *SecurityData `json:"data,omitempty"`
}

// SecurityData contains the detailed security scan information
type SecurityData struct {
	Layer *SecurityLayer `json:"Layer,omitempty"`
}

// SecurityLayer represents the scanned layer information
type SecurityLayer struct {
	Name             string            `json:"Name,omitempty"`
	ParentName       string            `json:"ParentName,omitempty"`
	NamespaceName    string            `json:"NamespaceName,omitempty"`
	IndexedByVersion int               `json:"IndexedByVersion,omitempty"`
	Features         []SecurityFeature `json:"Features,omitempty"`
}

// SecurityFeature represents a package/feature found in the image
type SecurityFeature struct {
	Name            string                  `json:"Name,omitempty"`
	VersionFormat   string                  `json:"VersionFormat,omitempty"`
	NamespaceName   string                  `json:"NamespaceName,omitempty"`
	AddedBy         string                  `json:"AddedBy,omitempty"`
	Version         string                  `json:"Version,omitempty"`
	Vulnerabilities []SecurityVulnerability `json:"Vulnerabilities,omitempty"`
}

// SecurityVulnerability represents a vulnerability found in a feature
type SecurityVulnerability struct {
	Name          string                 `json:"Name,omitempty"`
	NamespaceName string                 `json:"NamespaceName,omitempty"`
	Description   string                 `json:"Description,omitempty"`
	Link          string                 `json:"Link,omitempty"`
	Severity      string                 `json:"Severity,omitempty"`
	Metadata      map[string]interface{} `json:"Metadata,omitempty"`
	FixedBy       string                 `json:"FixedBy,omitempty"`
}

// Build Structures

// Build represents a repository build
type Build struct {
	ID              string                 `json:"id,omitempty"`
	Phase           string                 `json:"phase,omitempty"`
	Started         string                 `json:"started,omitempty"`
	DisplayName     string                 `json:"display_name,omitempty"`
	Status          map[string]interface{} `json:"status,omitempty"`
	Subdirectory    string                 `json:"subdirectory,omitempty"`
	Dockerfile      string                 `json:"dockerfile_path,omitempty"`
	Context         string                 `json:"context,omitempty"`
	IsWriter        bool                   `json:"is_writer,omitempty"`
	Trigger         *BuildTrigger          `json:"trigger,omitempty"`
	TriggerMetadata map[string]interface{} `json:"trigger_metadata,omitempty"`
	ResourceKey     string                 `json:"resource_key,omitempty"`
	Pull            *BuildPullRobot        `json:"pull_robot,omitempty"`
	Repository      *BuildRepository       `json:"repository,omitempty"`
	Error           string                 `json:"error,omitempty"`
	ManualUser      string                 `json:"manual_user,omitempty"`
	Archive         string                 `json:"archive_url,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
}

// BuildTrigger represents a build trigger
type BuildTrigger struct {
	ID            string `json:"id,omitempty"`
	Service       string `json:"service,omitempty"`
	IsActive      bool   `json:"is_active,omitempty"`
	BuildSource   string `json:"build_source,omitempty"`
	RepositoryURL string `json:"repository_url,omitempty"`
}

// BuildPullRobot represents a robot account for pulling
type BuildPullRobot struct {
	Name    string `json:"name,omitempty"`
	IsRobot bool   `json:"is_robot,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Avatar  Avatar `json:"avatar,omitempty"`
}

// BuildRepository represents repository info in a build
type BuildRepository struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

// Builds represents a list of builds
type Builds struct {
	Builds []Build `json:"builds,omitempty"`
}

// BuildLogs represents build logs
type BuildLogs struct {
	Start int             `json:"start,omitempty"`
	Total int             `json:"total,omitempty"`
	Logs  []BuildLogEntry `json:"logs,omitempty"`
}

// BuildLogEntry represents a single log entry in build logs
type BuildLogEntry struct {
	Type    string                 `json:"type,omitempty"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// RequestBuildRequest represents the request to trigger a build
type RequestBuildRequest struct {
	FileID         string   `json:"file_id,omitempty"`
	ArchiveURL     string   `json:"archive_url,omitempty"`
	Subdirectory   string   `json:"subdirectory,omitempty"`
	DockerfilePath string   `json:"dockerfile_path,omitempty"`
	Context        string   `json:"context,omitempty"`
	PullRobot      string   `json:"pull_robot,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

// Repository Notification Structures

// RepositoryNotification represents a repository notification/webhook
type RepositoryNotification struct {
	UUID             string                 `json:"uuid,omitempty"`
	Title            string                 `json:"title,omitempty"`
	Event            string                 `json:"event,omitempty"`
	Method           string                 `json:"method,omitempty"`
	Config           map[string]interface{} `json:"config,omitempty"`
	EventConfig      map[string]interface{} `json:"event_config,omitempty"`
	NumberOfFailures int                    `json:"number_of_failures,omitempty"`
}

// RepositoryNotifications represents a list of notifications
type RepositoryNotifications struct {
	Notifications []RepositoryNotification `json:"notifications,omitempty"`
}

// CreateNotificationRequest represents the request to create a notification
type CreateNotificationRequest struct {
	Event       string                 `json:"event"`
	Method      string                 `json:"method"`
	Config      map[string]interface{} `json:"config"`
	EventConfig map[string]interface{} `json:"eventConfig,omitempty"`
	Title       string                 `json:"title,omitempty"`
}

// TestNotificationResponse represents the response from testing a notification
type TestNotificationResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// Error Response Structure

// QuayError represents a Quay API error response
type QuayError struct {
	Status      int                    `json:"status,omitempty"`
	Error       string                 `json:"error,omitempty"`
	ErrorType   string                 `json:"error_type,omitempty"`
	Detail      string                 `json:"detail,omitempty"`
	ErrorDetail map[string]interface{} `json:"error_detail,omitempty"`
}
