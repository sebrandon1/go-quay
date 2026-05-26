package lib

// Shared test constants used across all test files in the lib package.
// Organized by category to avoid duplication and satisfy goconst linter.

const (
	// HTTP methods
	httpMethodGet    = "GET"
	httpMethodPost   = "POST"
	httpMethodPut    = "PUT"
	httpMethodDelete = "DELETE"

	// Common test namespaces and repositories
	testNamespace  = "testorg"
	testRepository = "testrepo"

	// Common test roles
	testRoleRead  = "read"
	testRoleWrite = "write"

	// Common test tag names
	testTagNameLatest = "latest"
	testTagNameV1     = "v1.0.0"

	// Common test timestamps
	testTimestamp      = "2024-01-15T10:30:00Z"
	testExpirationTime = "2024-12-31T23:59:59Z"

	// Common test digests
	testDigestSHA256 = "sha256:abc123"

	// Common test media types
	testMediaTypePlain = "text/plain"

	// Common test kinds
	testKindImage = "image"
	testKindUser  = "user"
	testKindRobot = "robot"
	testKindTeam  = "team"

	// Common test entity names
	testEmailAddress     = "test@example.com"
	testUserName         = "testuser"
	testSearchQueryValue = "quay"

	// Manifest label keys and values
	testLabelKeyVersion     = "version"
	testLabelKeyAPI         = "api"
	testLabelKeyEnvironment = "environment"
	testLabelValProduction  = "production"

	// Permissions test values
	testPermUserName = "john.doe"

	// Prototype test values
	testPrototypeTeamName = "devteam"

	// Repository test values
	testRepoDescription = "Test repository"

	// Repo token test values
	testRepoTokenName = "CI Token"

	// Robot test values
	testRobotFullName  = "testuser+deploybot"
	testRobotDescValue = "Deployment robot"

	// Security scan test values
	testSecScanStatus    = "scanned"
	testSecScanImageName = "centos:8"

	// Architecture test values
	testArchAmd64 = "amd64"

	// Client test values
	testTokenValue = "test-token"

	// Billing test values
	testBillingPlanFree          = "free"
	testBillingPlanTypeFree      = "free"
	testBillingSubscriptionID    = "sub-123"
	testBillingInvoiceID         = "inv-123"
	testBillingInvoiceStatusPaid = "paid"
	testBillingCurrencyUSD       = "USD"
	testBillingPeriodMonthly     = "monthly"

	// Team description test values (shared across team_test.go and organization_test.go)
	testTeamDescDev = "Dev team"
	testTeamDescNew = "New team"

	// Auto-prune method test values
	testAutoPruneMethodNumberOfTags = "number_of_tags"
)
