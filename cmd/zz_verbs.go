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
// violating Cobra's single-parent rule. Flags are copied; RunE is shared.
func verbLeaf(use string, src *cobra.Command) *cobra.Command {
	dst := &cobra.Command{
		Use:   use,
		Short: src.Short,
		Long:  src.Long,
		Args:  src.Args,
		RunE:  src.RunE,
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

type verbSpec struct {
	use string
	src *cobra.Command
}

func addVerbs(parent *cobra.Command, deprecateSrc bool, leaves ...verbSpec) {
	for _, leaf := range leaves {
		parent.AddCommand(verbLeaf(leaf.use, leaf.src))
		if deprecateSrc {
			deprecate(leaf.src, parent.Use, leaf.use)
		}
	}
}

func registerVerbCommands() {
	addVerbs(createCmd, true,
		verbSpec{cmdRepository, repoCreateCmd},
		verbSpec{cmdOrganization, createOrgCmd},
		verbSpec{cmdQuota, createQuotaCmd},
		verbSpec{cmdProxyCache, createProxyCacheCmd},
		verbSpec{cmdOrgRobot, createRobotCmd},
		verbSpec{cmdApplication, createApplicationCmd},
		verbSpec{cmdAutoPrune, createAutoPruneCmd},
		verbSpec{"marketplace-subscription", createMarketplaceSubscriptionCmd},
		verbSpec{cmdRobot, robotCreateCmd},
		verbSpec{cmdTeam, teamCreateCmd},
		verbSpec{cmdNotification, notificationCreateCmd},
		verbSpec{"message", messagesCreateCmd},
		verbSpec{cmdPrototype, prototypeCreateCmd},
		verbSpec{cmdRepoToken, repotokenCreateCmd},
		verbSpec{cmdMirror, mirrorCreateCmd},
		verbSpec{cmdOrgMember, addMemberCmd},
		verbSpec{"team-member", teamAddMemberCmd},
		verbSpec{"manifest-label", manifestAddLabelCmd},
		verbSpec{cmdStar, starRepoCmd},
		verbSpec{cmdBuild, buildRequestCmd},
		verbSpec{"trigger-build", triggerStartCmd},
		verbSpec{cmdTrigger, triggerActivateCmd},
		verbSpec{"robot-federation", robotFederationCreateCmd},
		verbSpec{"org-robot-federation", orgRobotFederationCreateCmd},
		verbSpec{"invite", inviteMemberCmd},
		verbSpec{"org-logs-export", exportOrgLogsCmd},
		verbSpec{"user-logs-export", exportUserLogsCmd},
		verbSpec{"repo-logs-export", exportRepoLogsCmd},
	)

	addVerbs(deleteCmd, true,
		verbSpec{cmdRepository, repoDeleteCmd},
		verbSpec{cmdOrganization, deleteOrgCmd},
		verbSpec{cmdQuota, deleteQuotaCmd},
		verbSpec{cmdProxyCache, deleteProxyCacheCmd},
		verbSpec{cmdOrgRobot, deleteRobotCmd},
		verbSpec{cmdApplication, deleteApplicationCmd},
		verbSpec{cmdAutoPrune, deleteAutoPruneCmd},
		verbSpec{"marketplace-subscription", deleteMarketplaceSubscriptionCmd},
		verbSpec{"subscriptions", batchRemoveSubscriptionsCmd},
		verbSpec{cmdRobot, robotDeleteCmd},
		verbSpec{cmdTeam, teamDeleteCmd},
		verbSpec{cmdNotification, notificationDeleteCmd},
		verbSpec{cmdPrototype, prototypeDeleteCmd},
		verbSpec{cmdRepoToken, repotokenDeleteCmd},
		verbSpec{cmdTag, tagDeleteCmd},
		verbSpec{cmdManifest, manifestDeleteCmd},
		verbSpec{"permission", permRemoveCmd},
		verbSpec{cmdUserPermission, permDeleteUserPermCmd},
		verbSpec{cmdTeamPermission, permDeleteTeamPermCmd},
		verbSpec{cmdOrgMember, removeMemberCmd},
		verbSpec{"team-member", teamRemoveMemberCmd},
		verbSpec{"manifest-label", manifestRemoveLabelCmd},
		verbSpec{cmdStar, unstarRepoCmd},
		verbSpec{cmdBuild, buildCancelCmd},
		verbSpec{cmdTrigger, triggerDeleteCmd},
		verbSpec{"robot-federation", robotFederationDeleteCmd},
		verbSpec{"org-robot-federation", orgRobotFederationDeleteCmd},
		verbSpec{"invite", cancelInviteCmd},
		verbSpec{"team-repo-permission", teamRemovePermissionCmd},
		verbSpec{"org-robot-permission", removeRobotPermissionCmd},
	)

	addVerbs(updateCmd, true,
		verbSpec{cmdRepository, repoUpdateCmd},
		verbSpec{"visibility", repoChangeVisibilityCmd},
		verbSpec{cmdOrganization, updateOrgCmd},
		verbSpec{cmdQuota, updateQuotaCmd},
		verbSpec{cmdApplication, updateApplicationCmd},
		verbSpec{"application-secret", resetApplicationSecretCmd},
		verbSpec{cmdAutoPrune, updateAutoPruneCmd},
		verbSpec{cmdOrgRobot, regenerateRobotCmd},
		verbSpec{cmdRobot, robotRegenerateCmd},
		verbSpec{cmdTeam, teamUpdateCmd},
		verbSpec{cmdNotification, notificationUpdateCmd},
		verbSpec{"notification-reset", notificationResetCmd},
		verbSpec{"notification-test", notificationTestCmd},
		verbSpec{cmdPrototype, prototypeUpdateCmd},
		verbSpec{cmdRepoToken, repotokenUpdateCmd},
		verbSpec{cmdMirror, mirrorUpdateCmd},
		verbSpec{cmdTag, tagUpdateCmd},
		verbSpec{"tag-revert", tagRevertCmd},
		verbSpec{"tag-restore", tagRestoreCmd},
		verbSpec{"tag-change", tagChangeCmd},
		verbSpec{"permission", permSetCmd},
		verbSpec{cmdUserPermission, permSetUserPermCmd},
		verbSpec{cmdTeamPermission, permSetTeamPermCmd},
		verbSpec{"team-repo-permission", teamSetPermissionCmd},
		verbSpec{"org-robot-permission", setRobotPermissionCmd},
		verbSpec{"trigger-enable", triggerEnableCmd},
		verbSpec{"trigger-disable", triggerDisableCmd},
	)

	addVerbs(listCmd, false,
		verbSpec{cmdRepository, repoListCmd},
		verbSpec{"org-members", orgMembersCmd},
		verbSpec{"org-repositories", orgRepositoriesCmd},
		verbSpec{"org-collaborators", collaboratorsCmd},
		verbSpec{"org-robots", orgRobotsCmd},
		verbSpec{"org-teams", orgTeamsCmd},
		verbSpec{"org-applications", orgApplicationsCmd},
		verbSpec{cmdAutoPrune, autoPruneCmd},
		verbSpec{"prototypes", prototypeListCmd},
		verbSpec{"repotokens", repotokenListCmd},
		verbSpec{"robots", robotListCmd},
		verbSpec{"teams", teamListCmd},
		verbSpec{"notifications", notificationListCmd},
		verbSpec{"messages", messagesListCmd},
		verbSpec{"builds", buildListCmd},
		verbSpec{"triggers", triggerListCmd},
		verbSpec{"permissions", permListCmd},
		verbSpec{"user-permissions", permUserPermissionsCmd},
		verbSpec{"team-permissions", permTeamPermissionsCmd},
		verbSpec{"starred", userStarredCmd},
	)

	addVerbs(infoCmd, false,
		verbSpec{cmdRepository, repoInfoCmd},
		verbSpec{cmdOrganization, orgInfoCmd},
		verbSpec{cmdQuota, orgQuotaCmd},
		verbSpec{cmdProxyCache, proxyCacheCmd},
		verbSpec{cmdOrgRobot, orgRobotCmd},
		verbSpec{cmdApplication, applicationCmd},
		verbSpec{cmdAutoPrune, autoPrunePolicyCmd},
		verbSpec{cmdMarketplace, marketplaceCmd},
		verbSpec{cmdOrgMember, getMemberCmd},
		verbSpec{cmdTeam, teamCmdInfoCmd},
		verbSpec{"org-team", teamInfoCmd},
		verbSpec{cmdTag, tagInfoCmd},
		verbSpec{"user", userInfoCmd},
		verbSpec{cmdManifest, manifestInfoCmd},
		verbSpec{"secscan", secscanInfoCmd},
		verbSpec{cmdRobot, robotInfoCmd},
		verbSpec{cmdNotification, notificationInfoCmd},
		verbSpec{cmdPrototype, prototypeInfoCmd},
		verbSpec{cmdRepoToken, repotokenInfoCmd},
		verbSpec{cmdMirror, mirrorInfoCmd},
		verbSpec{cmdBuild, buildInfoCmd},
		verbSpec{cmdTrigger, triggerInfoCmd},
		verbSpec{"error", errorTypeCmd},
	)
}
