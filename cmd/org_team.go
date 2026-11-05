package cmd

import (
	"github.com/spf13/cobra"
)

// The commands in this file are deprecated shims kept for one release. The
// canonical implementation lives in cmd/team.go and uses --organization/--name
// everywhere; these legacy paths keep their original --team/--email flags so
// existing scripts continue to work while pointing users at the new forms.

// Organization Teams
var orgTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Get organization teams (deprecated)",
	Long:  `Get list of all teams in an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTeamList(cmd, orgName)
	},
}

// Team Info
var teamInfoCmd = &cobra.Command{
	Use:   "team",
	Short: "Get team information (deprecated)",
	Long:  `Get detailed information about a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTeamInfo(cmd, orgName, teamName)
	},
}

// Team Members
var teamMembersCmd = &cobra.Command{
	Use:   "team-members",
	Short: "Get team members (deprecated)",
	Long:  `Get list of all members in a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTeamMembers(cmd, orgName, teamName)
	},
}

// Invite Team Member
var orgInviteMemberCmd = &cobra.Command{
	Use:   "invite-member",
	Short: "Invite a member to a team (deprecated)",
	Long:  `Invite a member to a team by email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTeamInvite(cmd, orgName, teamName, email)
	},
}

// Cancel Team Invite
var orgCancelInviteCmd = &cobra.Command{
	Use:   "cancel-invite",
	Short: "Cancel a team member invitation (deprecated)",
	Long:  `Cancel a pending team member invitation. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return errConfirmCancelInvite
		}
		return runTeamCancelInvite(cmd, orgName, teamName, email)
	},
}

func initOrgTeamFlags() {
	teamInfoCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = teamInfoCmd.MarkFlagRequired("team")

	teamMembersCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = teamMembersCmd.MarkFlagRequired("team")
}

func initOrgInviteFlags() {
	orgInviteMemberCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = orgInviteMemberCmd.MarkFlagRequired("team")
	orgInviteMemberCmd.Flags().StringVar(&email, "email", "", "Email address")
	_ = orgInviteMemberCmd.MarkFlagRequired("email")

	orgCancelInviteCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = orgCancelInviteCmd.MarkFlagRequired("team")
	orgCancelInviteCmd.Flags().StringVar(&email, "email", "", "Email address")
	_ = orgCancelInviteCmd.MarkFlagRequired("email")
	orgCancelInviteCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm cancellation")
}

func init() {
	organizationCmd.AddCommand(orgTeamsCmd)
	organizationCmd.AddCommand(teamInfoCmd)
	organizationCmd.AddCommand(teamMembersCmd)
	organizationCmd.AddCommand(orgInviteMemberCmd)
	organizationCmd.AddCommand(orgCancelInviteCmd)

	deprecate(orgTeamsCmd, cmdList, "teams")
	deprecate(teamInfoCmd, cmdInfo, cmdTeam)
	deprecate(teamMembersCmd, cmdList, "team-members")
	deprecate(orgInviteMemberCmd, cmdCreate, "invite")
	deprecate(orgCancelInviteCmd, cmdDelete, "invite")

	initOrgTeamFlags()
	initOrgInviteFlags()
}
