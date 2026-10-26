package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	createCmd = &cobra.Command{Use: cmdCreate, Short: "Create Quay.io resources"}
	deleteCmd = &cobra.Command{Use: cmdDelete, Short: "Delete Quay.io resources"}
	updateCmd = &cobra.Command{Use: cmdUpdate, Short: "Update Quay.io resources"}
	listCmd   = &cobra.Command{Use: cmdList, Short: "List Quay.io resources"}
	infoCmd   = &cobra.Command{Use: cmdInfo, Short: "Show a single Quay.io resource"}
)

func init() {
	rootCmd.AddCommand(createCmd, deleteCmd, updateCmd, listCmd, infoCmd)
	registerVerbCommands()
}

// verbLeaf clones src so the same action can live under a verb parent without
// violating Cobra's single-parent rule. Flags, hooks, and RunE are copied, so
// local validation behaves the same on both paths. Example is passed in rather
// than inherited: the legacy examples spell out "go-quay get ...", which the
// verb-first tree must not advertise.
func verbLeaf(use, example string, src *cobra.Command) *cobra.Command {
	dst := &cobra.Command{
		Use:               use,
		Short:             src.Short,
		Long:              src.Long,
		Example:           example,
		Args:              src.Args,
		RunE:              src.RunE,
		PreRun:            src.PreRun,
		PostRun:           src.PostRun,
		PreRunE:           src.PreRunE,
		PostRunE:          src.PostRunE,
		ValidArgsFunction: src.ValidArgsFunction,
	}
	copyFlagSet(dst.Flags(), src.LocalNonPersistentFlags())
	copyInheritedResourceFlags(dst, src)
	return dst
}

func copyInheritedResourceFlags(dst, src *cobra.Command) {
	src.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if rootCmd.PersistentFlags().Lookup(f.Name) != nil {
			return
		}
		copyOneFlag(dst.Flags(), f)
	})
}

func copyFlagSet(dst, src *pflag.FlagSet) {
	if src == nil {
		return
	}
	src.VisitAll(func(f *pflag.Flag) {
		copyOneFlag(dst, f)
	})
}

func copyOneFlag(dst *pflag.FlagSet, f *pflag.Flag) {
	if f.Name == "help" || dst.Lookup(f.Name) != nil {
		return
	}
	cf := *f
	cf.Changed = false
	if f.Annotations != nil {
		cf.Annotations = make(map[string][]string, len(f.Annotations))
		for k, v := range f.Annotations {
			cf.Annotations[k] = append([]string(nil), v...)
		}
	}
	dst.AddFlag(&cf)
}

func deprecate(cmd *cobra.Command, verb, resource string) {
	cmd.Deprecated = fmt.Sprintf("use %q instead", cliName+" "+verb+" "+resource)
}

const exampleRepoFlags = `--namespace myorg --repository myapp --token "$QUAY_TOKEN"`

// verbSpec describes one verb-first alias. exampleArgs is only the flags part of
// the example; addVerbs prefixes the command path so it cannot drift from use.
type verbSpec struct {
	use         string
	src         *cobra.Command
	exampleArgs string
}

func addVerbs(parent *cobra.Command, deprecateSrc bool, leaves ...verbSpec) {
	for _, leaf := range leaves {
		example := ""
		if leaf.exampleArgs != "" {
			example = "  " + cliName + " " + parent.Use + " " + leaf.use + " " + leaf.exampleArgs
		}
		parent.AddCommand(verbLeaf(leaf.use, example, leaf.src))
		if deprecateSrc {
			deprecate(leaf.src, parent.Use, leaf.use)
		}
	}
}

func registerVerbCommands() {
	addVerbs(createCmd, true,
		verbSpec{use: cmdRepository, src: repoCreateCmd},
		verbSpec{use: cmdOrganization, src: createOrgCmd},
		verbSpec{use: cmdQuota, src: createQuotaCmd},
		verbSpec{use: cmdProxyCache, src: createProxyCacheCmd},
		verbSpec{use: cmdOrgRobot, src: createRobotCmd},
		verbSpec{use: cmdApplication, src: createApplicationCmd},
		verbSpec{use: cmdAutoPrune, src: createAutoPruneCmd},
		verbSpec{use: "marketplace-subscription", src: createMarketplaceSubscriptionCmd},
		verbSpec{use: cmdRobot, src: robotCreateCmd},
		verbSpec{use: cmdTeam, src: teamCreateCmd},
		verbSpec{use: cmdNotification, src: notificationCreateCmd},
		verbSpec{use: "message", src: messagesCreateCmd},
		verbSpec{use: cmdPrototype, src: prototypeCreateCmd},
		verbSpec{use: cmdRepoToken, src: repotokenCreateCmd},
		verbSpec{use: cmdMirror, src: mirrorCreateCmd},
		verbSpec{use: cmdOrgMember, src: addMemberCmd},
		verbSpec{use: "team-member", src: teamAddMemberCmd},
		verbSpec{use: "manifest-label", src: manifestAddLabelCmd},
		verbSpec{use: cmdStar, src: starRepoCmd},
		verbSpec{use: cmdBuild, src: buildRequestCmd},
		verbSpec{use: "trigger-build", src: triggerStartCmd},
		verbSpec{use: cmdTrigger, src: triggerActivateCmd},
		verbSpec{use: "robot-federation", src: robotFederationCreateCmd},
		verbSpec{use: "org-robot-federation", src: orgRobotFederationCreateCmd},
		verbSpec{use: "invite", src: inviteMemberCmd},
		verbSpec{use: "org-logs-export", src: exportOrgLogsCmd},
		verbSpec{use: "user-logs-export", src: exportUserLogsCmd},
		verbSpec{use: "repo-logs-export", src: exportRepoLogsCmd},
	)

	addVerbs(deleteCmd, true,
		verbSpec{use: cmdRepository, src: repoDeleteCmd},
		verbSpec{use: cmdOrganization, src: deleteOrgCmd},
		verbSpec{use: cmdQuota, src: deleteQuotaCmd},
		verbSpec{use: cmdProxyCache, src: deleteProxyCacheCmd},
		verbSpec{use: cmdOrgRobot, src: deleteRobotCmd},
		verbSpec{use: cmdApplication, src: deleteApplicationCmd},
		verbSpec{use: cmdAutoPrune, src: deleteAutoPruneCmd},
		verbSpec{use: "marketplace-subscription", src: deleteMarketplaceSubscriptionCmd},
		verbSpec{use: "subscriptions", src: batchRemoveSubscriptionsCmd},
		verbSpec{use: cmdRobot, src: robotDeleteCmd},
		verbSpec{use: cmdTeam, src: teamDeleteCmd},
		verbSpec{use: cmdNotification, src: notificationDeleteCmd},
		verbSpec{use: cmdPrototype, src: prototypeDeleteCmd},
		verbSpec{use: cmdRepoToken, src: repotokenDeleteCmd},
		verbSpec{use: cmdTag, src: tagDeleteCmd},
		verbSpec{use: cmdManifest, src: manifestDeleteCmd},
		verbSpec{use: "permission", src: permRemoveCmd},
		verbSpec{use: cmdUserPermission, src: permDeleteUserPermCmd},
		verbSpec{use: cmdTeamPermission, src: permDeleteTeamPermCmd},
		verbSpec{use: cmdOrgMember, src: removeMemberCmd},
		verbSpec{use: "team-member", src: teamRemoveMemberCmd},
		verbSpec{use: "manifest-label", src: manifestRemoveLabelCmd},
		verbSpec{use: cmdStar, src: unstarRepoCmd},
		verbSpec{use: cmdBuild, src: buildCancelCmd},
		verbSpec{use: cmdTrigger, src: triggerDeleteCmd},
		verbSpec{use: "robot-federation", src: robotFederationDeleteCmd},
		verbSpec{use: "org-robot-federation", src: orgRobotFederationDeleteCmd},
		verbSpec{use: "invite", src: cancelInviteCmd},
		verbSpec{use: "team-repo-permission", src: teamRemovePermissionCmd},
		verbSpec{use: "org-robot-permission", src: removeRobotPermissionCmd},
	)

	addVerbs(updateCmd, true,
		verbSpec{use: cmdRepository, src: repoUpdateCmd},
		verbSpec{use: "visibility", src: repoChangeVisibilityCmd},
		verbSpec{use: cmdOrganization, src: updateOrgCmd},
		verbSpec{use: cmdQuota, src: updateQuotaCmd},
		verbSpec{use: cmdApplication, src: updateApplicationCmd},
		verbSpec{use: "application-secret", src: resetApplicationSecretCmd},
		verbSpec{use: cmdAutoPrune, src: updateAutoPruneCmd},
		verbSpec{use: cmdOrgRobot, src: regenerateRobotCmd},
		verbSpec{use: cmdRobot, src: robotRegenerateCmd},
		verbSpec{use: cmdTeam, src: teamUpdateCmd},
		verbSpec{use: cmdNotification, src: notificationUpdateCmd},
		verbSpec{use: "notification-reset", src: notificationResetCmd},
		verbSpec{use: "notification-test", src: notificationTestCmd},
		verbSpec{use: cmdPrototype, src: prototypeUpdateCmd},
		verbSpec{use: cmdRepoToken, src: repotokenUpdateCmd},
		verbSpec{use: cmdMirror, src: mirrorUpdateCmd},
		verbSpec{use: cmdTag, src: tagUpdateCmd},
		verbSpec{use: "tag-revert", src: tagRevertCmd},
		verbSpec{use: "tag-restore", src: tagRestoreCmd},
		verbSpec{use: "tag-change", src: tagChangeCmd},
		verbSpec{use: "permission", src: permSetCmd},
		verbSpec{use: cmdUserPermission, src: permSetUserPermCmd},
		verbSpec{use: cmdTeamPermission, src: permSetTeamPermCmd},
		verbSpec{use: "team-repo-permission", src: teamSetPermissionCmd},
		verbSpec{use: "org-robot-permission", src: setRobotPermissionCmd},
		verbSpec{use: "trigger-enable", src: triggerEnableCmd},
		verbSpec{use: "trigger-disable", src: triggerDisableCmd},
	)

	addVerbs(listCmd, false,
		verbSpec{use: cmdRepository, src: repoListCmd},
		verbSpec{use: "org-members", src: orgMembersCmd},
		verbSpec{use: "org-repositories", src: orgRepositoriesCmd},
		verbSpec{use: "org-collaborators", src: collaboratorsCmd},
		verbSpec{use: "org-robots", src: orgRobotsCmd},
		verbSpec{use: "org-teams", src: orgTeamsCmd},
		verbSpec{use: "org-applications", src: orgApplicationsCmd},
		verbSpec{use: cmdAutoPrune, src: autoPruneCmd},
		verbSpec{use: "prototypes", src: prototypeListCmd},
		verbSpec{use: "repotokens", src: repotokenListCmd},
		verbSpec{use: "robots", src: robotListCmd},
		verbSpec{use: "teams", src: teamListCmd},
		verbSpec{use: "notifications", src: notificationListCmd},
		verbSpec{use: "messages", src: messagesListCmd},
		verbSpec{use: "builds", src: buildListCmd, exampleArgs: exampleRepoFlags},
		verbSpec{use: "triggers", src: triggerListCmd, exampleArgs: exampleRepoFlags},
		verbSpec{use: "permissions", src: permListCmd, exampleArgs: exampleRepoFlags},
		verbSpec{use: "user-permissions", src: permUserPermissionsCmd},
		verbSpec{use: "team-permissions", src: permTeamPermissionsCmd},
		verbSpec{use: "starred", src: userStarredCmd},
	)

	addVerbs(infoCmd, false,
		verbSpec{use: cmdRepository, src: repoInfoCmd, exampleArgs: exampleRepoFlags},
		verbSpec{use: cmdOrganization, src: orgInfoCmd, exampleArgs: `--organization myorg --token "$QUAY_TOKEN"`},
		verbSpec{use: cmdQuota, src: orgQuotaCmd},
		verbSpec{use: cmdProxyCache, src: proxyCacheCmd},
		verbSpec{use: cmdOrgRobot, src: orgRobotCmd},
		verbSpec{use: cmdApplication, src: applicationCmd},
		verbSpec{use: cmdAutoPrune, src: autoPrunePolicyCmd},
		verbSpec{use: cmdMarketplace, src: marketplaceCmd},
		verbSpec{use: cmdOrgMember, src: getMemberCmd},
		verbSpec{use: cmdTeam, src: teamCmdInfoCmd, exampleArgs: `--organization myorg --name developers --token "$QUAY_TOKEN"`},
		verbSpec{use: "org-team", src: teamInfoCmd},
		verbSpec{use: cmdTag, src: tagInfoCmd, exampleArgs: `--namespace myorg --repository myapp --tag latest --token "$QUAY_TOKEN"`},
		verbSpec{use: cmdUser, src: userInfoCmd},
		verbSpec{use: cmdManifest, src: manifestInfoCmd},
		verbSpec{use: "secscan", src: secscanInfoCmd, exampleArgs: `--namespace myorg --repository myapp --manifest sha256:abc123 --token "$QUAY_TOKEN"`},
		verbSpec{use: cmdRobot, src: robotInfoCmd, exampleArgs: `--name buildbot --token "$QUAY_TOKEN"`},
		verbSpec{use: cmdNotification, src: notificationInfoCmd},
		verbSpec{use: cmdPrototype, src: prototypeInfoCmd},
		verbSpec{use: cmdRepoToken, src: repotokenInfoCmd},
		verbSpec{use: cmdMirror, src: mirrorInfoCmd},
		verbSpec{use: cmdBuild, src: buildInfoCmd},
		verbSpec{use: cmdTrigger, src: triggerInfoCmd},
		verbSpec{use: "error", src: errorTypeCmd},
	)
}
