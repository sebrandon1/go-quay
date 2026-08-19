package cmd

const (
	subcmdInfo        = "info"
	subcmdList        = "list"
	subcmdCreate      = "create"
	subcmdUpdate      = "update"
	subcmdDelete      = "delete"
	subcmdPermissions = "permissions"

	// Command name constants
	cliName         = "go-quay"
	cmdGet          = "get"
	cmdCreate       = "create"
	cmdDelete       = "delete"
	cmdUpdate       = "update"
	cmdList         = "list"
	cmdInfo         = "info"
	cmdRepository   = "repository"
	cmdOrganization = "organization"
	cmdTag          = "tag"
	cmdMembers      = "members"

	// Verb-first resource names (also used as Cobra Use strings in domain files).
	cmdQuota          = "quota"
	cmdProxyCache     = "proxy-cache"
	cmdOrgRobot       = "org-robot"
	cmdApplication    = "application"
	cmdAutoPrune      = "auto-prune"
	cmdRobot          = "robot"
	cmdTeam           = "team"
	cmdNotification   = "notification"
	cmdPrototype      = "prototype"
	cmdRepoToken      = "repotoken"
	cmdMirror         = "mirror"
	cmdOrgMember      = "org-member"
	cmdStar           = "star"
	cmdBuild          = "build"
	cmdTrigger        = "trigger"
	cmdManifest       = "manifest"
	cmdUserPermission = "user-permission"
	cmdTeamPermission = "team-permission"
	cmdMarketplace    = "marketplace"

	// Output format constants
	outputJSON  = "json"
	outputYAML  = "yaml"
	outputTable = "table"
)
