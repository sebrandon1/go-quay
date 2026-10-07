package lib

import (
	"context"
	"testing"
)

const (
	validationNamespace     = "org"
	validationRepository    = "repo"
	validationUsername      = "user"
	validationTeam          = "team"
	validationRole          = "read"
	validationMember        = "member"
	validationEmail         = "user@example.com"
	namespaceRequiredError  = "namespace is required"
	repositoryRequiredError = "repository is required"
	usernameRequiredError   = "username is required"
	roleRequiredError       = "role is required"
	teamnameRequiredError   = "teamname is required"
	orgnameRequiredError    = "orgname is required"
	memberRequiredError     = "membername is required"
	emailRequiredError      = "email is required"
)

type validationErrorCase struct {
	name string
	want string
	call func() error
}

func TestPermissionRequiredArgumentValidation(t *testing.T) {
	client := &Client{}
	ctx := context.Background()
	tests := []validationErrorCase{
		{name: "get repository permissions namespace", want: namespaceRequiredError, call: func() error { _, err := client.GetRepositoryPermissions(ctx, "", validationRepository); return err }},
		{name: "get repository permissions repository", want: repositoryRequiredError, call: func() error { _, err := client.GetRepositoryPermissions(ctx, validationNamespace, ""); return err }},
		{name: "set repository permission namespace", want: namespaceRequiredError, call: func() error {
			return client.SetRepositoryPermission(ctx, "", validationRepository, validationUsername, validationRole)
		}},
		{name: "set repository permission repository", want: repositoryRequiredError, call: func() error {
			return client.SetRepositoryPermission(ctx, validationNamespace, "", validationUsername, validationRole)
		}},
		{name: "set repository permission username", want: usernameRequiredError, call: func() error {
			return client.SetRepositoryPermission(ctx, validationNamespace, validationRepository, "", validationRole)
		}},
		{name: "set repository permission role", want: roleRequiredError, call: func() error {
			return client.SetRepositoryPermission(ctx, validationNamespace, validationRepository, validationUsername, "")
		}},
		{name: "remove repository permission namespace", want: namespaceRequiredError, call: func() error {
			return client.RemoveRepositoryPermission(ctx, "", validationRepository, validationUsername)
		}},
		{name: "remove repository permission repository", want: repositoryRequiredError, call: func() error {
			return client.RemoveRepositoryPermission(ctx, validationNamespace, "", validationUsername)
		}},
		{name: "remove repository permission username", want: usernameRequiredError, call: func() error {
			return client.RemoveRepositoryPermission(ctx, validationNamespace, validationRepository, "")
		}},
		{name: "list user permissions namespace", want: namespaceRequiredError, call: func() error { _, err := client.ListUserPermissions(ctx, "", validationRepository); return err }},
		{name: "list user permissions repository", want: repositoryRequiredError, call: func() error { _, err := client.ListUserPermissions(ctx, validationNamespace, ""); return err }},
		{name: "get user permission namespace", want: namespaceRequiredError, call: func() error {
			_, err := client.GetUserPermission(ctx, "", validationRepository, validationUsername)
			return err
		}},
		{name: "get user permission repository", want: repositoryRequiredError, call: func() error {
			_, err := client.GetUserPermission(ctx, validationNamespace, "", validationUsername)
			return err
		}},
		{name: "get user permission username", want: usernameRequiredError, call: func() error {
			_, err := client.GetUserPermission(ctx, validationNamespace, validationRepository, "")
			return err
		}},
		{name: "set user permission namespace", want: namespaceRequiredError, call: func() error {
			return client.SetUserPermission(ctx, "", validationRepository, validationUsername, validationRole)
		}},
		{name: "set user permission repository", want: repositoryRequiredError, call: func() error {
			return client.SetUserPermission(ctx, validationNamespace, "", validationUsername, validationRole)
		}},
		{name: "set user permission username", want: usernameRequiredError, call: func() error {
			return client.SetUserPermission(ctx, validationNamespace, validationRepository, "", validationRole)
		}},
		{name: "set user permission role", want: roleRequiredError, call: func() error {
			return client.SetUserPermission(ctx, validationNamespace, validationRepository, validationUsername, "")
		}},
		{name: "delete user permission namespace", want: namespaceRequiredError, call: func() error { return client.DeleteUserPermission(ctx, "", validationRepository, validationUsername) }},
		{name: "delete user permission repository", want: repositoryRequiredError, call: func() error { return client.DeleteUserPermission(ctx, validationNamespace, "", validationUsername) }},
		{name: "delete user permission username", want: usernameRequiredError, call: func() error { return client.DeleteUserPermission(ctx, validationNamespace, validationRepository, "") }},
		{name: "get transitive permission namespace", want: namespaceRequiredError, call: func() error {
			_, err := client.GetUserTransitivePermission(ctx, "", validationRepository, validationUsername)
			return err
		}},
		{name: "get transitive permission repository", want: repositoryRequiredError, call: func() error {
			_, err := client.GetUserTransitivePermission(ctx, validationNamespace, "", validationUsername)
			return err
		}},
		{name: "get transitive permission username", want: usernameRequiredError, call: func() error {
			_, err := client.GetUserTransitivePermission(ctx, validationNamespace, validationRepository, "")
			return err
		}},
		{name: "list team permissions namespace", want: namespaceRequiredError, call: func() error { _, err := client.ListTeamPermissions(ctx, "", validationRepository); return err }},
		{name: "list team permissions repository", want: repositoryRequiredError, call: func() error { _, err := client.ListTeamPermissions(ctx, validationNamespace, ""); return err }},
		{name: "get team permission namespace", want: namespaceRequiredError, call: func() error {
			_, err := client.GetTeamPermission(ctx, "", validationRepository, validationTeam)
			return err
		}},
		{name: "get team permission repository", want: repositoryRequiredError, call: func() error {
			_, err := client.GetTeamPermission(ctx, validationNamespace, "", validationTeam)
			return err
		}},
		{name: "get team permission team", want: teamnameRequiredError, call: func() error {
			_, err := client.GetTeamPermission(ctx, validationNamespace, validationRepository, "")
			return err
		}},
		{name: "set team permission namespace", want: namespaceRequiredError, call: func() error {
			return client.SetTeamPermission(ctx, "", validationRepository, validationTeam, validationRole)
		}},
		{name: "set team permission repository", want: repositoryRequiredError, call: func() error {
			return client.SetTeamPermission(ctx, validationNamespace, "", validationTeam, validationRole)
		}},
		{name: "set team permission team", want: teamnameRequiredError, call: func() error {
			return client.SetTeamPermission(ctx, validationNamespace, validationRepository, "", validationRole)
		}},
		{name: "set team permission role", want: roleRequiredError, call: func() error {
			return client.SetTeamPermission(ctx, validationNamespace, validationRepository, validationTeam, "")
		}},
		{name: "delete team permission namespace", want: namespaceRequiredError, call: func() error { return client.DeleteTeamPermission(ctx, "", validationRepository, validationTeam) }},
		{name: "delete team permission repository", want: repositoryRequiredError, call: func() error { return client.DeleteTeamPermission(ctx, validationNamespace, "", validationTeam) }},
		{name: "delete team permission team", want: teamnameRequiredError, call: func() error { return client.DeleteTeamPermission(ctx, validationNamespace, validationRepository, "") }},
	}

	assertValidationErrors(t, tests)
}

func TestOrganizationTeamRequiredArgumentValidation(t *testing.T) {
	client := &Client{}
	ctx := context.Background()
	tests := []validationErrorCase{
		{name: "get teams organization", want: orgnameRequiredError, call: func() error { _, err := client.GetTeams(ctx, ""); return err }},
		{name: "create team organization", want: orgnameRequiredError, call: func() error { _, err := client.CreateTeam(ctx, "", validationTeam, "", validationRole); return err }},
		{name: "create team team", want: teamnameRequiredError, call: func() error {
			_, err := client.CreateTeam(ctx, validationNamespace, "", "", validationRole)
			return err
		}},
		{name: "create team role", want: roleRequiredError, call: func() error {
			_, err := client.CreateTeam(ctx, validationNamespace, validationTeam, "", "")
			return err
		}},
		{name: "get team organization", want: orgnameRequiredError, call: func() error { _, err := client.GetTeam(ctx, "", validationTeam); return err }},
		{name: "get team team", want: teamnameRequiredError, call: func() error { _, err := client.GetTeam(ctx, validationNamespace, ""); return err }},
		{name: "delete team organization", want: orgnameRequiredError, call: func() error { return client.DeleteTeam(ctx, "", validationTeam) }},
		{name: "delete team team", want: teamnameRequiredError, call: func() error { return client.DeleteTeam(ctx, validationNamespace, "") }},
		{name: "update team organization", want: orgnameRequiredError, call: func() error { _, err := client.UpdateTeam(ctx, "", validationTeam, "", validationRole); return err }},
		{name: "update team team", want: teamnameRequiredError, call: func() error {
			_, err := client.UpdateTeam(ctx, validationNamespace, "", "", validationRole)
			return err
		}},
		{name: "update team role", want: roleRequiredError, call: func() error {
			_, err := client.UpdateTeam(ctx, validationNamespace, validationTeam, "", "")
			return err
		}},
		{name: "get team members organization", want: orgnameRequiredError, call: func() error { _, err := client.GetTeamMembers(ctx, "", validationTeam); return err }},
		{name: "get team members team", want: teamnameRequiredError, call: func() error { _, err := client.GetTeamMembers(ctx, validationNamespace, ""); return err }},
		{name: "add team member organization", want: orgnameRequiredError, call: func() error { return client.AddTeamMember(ctx, "", validationTeam, validationMember) }},
		{name: "add team member team", want: teamnameRequiredError, call: func() error { return client.AddTeamMember(ctx, validationNamespace, "", validationMember) }},
		{name: "add team member member", want: memberRequiredError, call: func() error { return client.AddTeamMember(ctx, validationNamespace, validationTeam, "") }},
		{name: "remove team member organization", want: orgnameRequiredError, call: func() error { return client.RemoveTeamMember(ctx, "", validationTeam, validationMember) }},
		{name: "remove team member team", want: teamnameRequiredError, call: func() error { return client.RemoveTeamMember(ctx, validationNamespace, "", validationMember) }},
		{name: "remove team member member", want: memberRequiredError, call: func() error { return client.RemoveTeamMember(ctx, validationNamespace, validationTeam, "") }},
		{name: "get team permissions organization", want: orgnameRequiredError, call: func() error { _, err := client.GetTeamPermissions(ctx, "", validationTeam); return err }},
		{name: "get team permissions team", want: teamnameRequiredError, call: func() error { _, err := client.GetTeamPermissions(ctx, validationNamespace, ""); return err }},
		{name: "set repository permission organization", want: orgnameRequiredError, call: func() error {
			return client.SetTeamRepositoryPermission(ctx, "", validationTeam, validationRepository, validationRole)
		}},
		{name: "set repository permission team", want: teamnameRequiredError, call: func() error {
			return client.SetTeamRepositoryPermission(ctx, validationNamespace, "", validationRepository, validationRole)
		}},
		{name: "set repository permission repository", want: repositoryRequiredError, call: func() error {
			return client.SetTeamRepositoryPermission(ctx, validationNamespace, validationTeam, "", validationRole)
		}},
		{name: "set repository permission role", want: roleRequiredError, call: func() error {
			return client.SetTeamRepositoryPermission(ctx, validationNamespace, validationTeam, validationRepository, "")
		}},
		{name: "remove repository permission organization", want: orgnameRequiredError, call: func() error {
			return client.RemoveTeamRepositoryPermission(ctx, "", validationTeam, validationRepository)
		}},
		{name: "remove repository permission team", want: teamnameRequiredError, call: func() error {
			return client.RemoveTeamRepositoryPermission(ctx, validationNamespace, "", validationRepository)
		}},
		{name: "remove repository permission repository", want: repositoryRequiredError, call: func() error {
			return client.RemoveTeamRepositoryPermission(ctx, validationNamespace, validationTeam, "")
		}},
		{name: "invite member organization", want: orgnameRequiredError, call: func() error { return client.InviteTeamMember(ctx, "", validationTeam, validationEmail) }},
		{name: "invite member team", want: teamnameRequiredError, call: func() error { return client.InviteTeamMember(ctx, validationNamespace, "", validationEmail) }},
		{name: "invite member email", want: emailRequiredError, call: func() error { return client.InviteTeamMember(ctx, validationNamespace, validationTeam, "") }},
		{name: "delete invitation organization", want: orgnameRequiredError, call: func() error { return client.DeleteTeamInvite(ctx, "", validationTeam, validationEmail) }},
		{name: "delete invitation team", want: teamnameRequiredError, call: func() error { return client.DeleteTeamInvite(ctx, validationNamespace, "", validationEmail) }},
		{name: "delete invitation email", want: emailRequiredError, call: func() error { return client.DeleteTeamInvite(ctx, validationNamespace, validationTeam, "") }},
	}

	assertValidationErrors(t, tests)
}

func assertValidationErrors(t *testing.T, tests []validationErrorCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil || err.Error() != tt.want {
				t.Errorf("error = %v, want %q", err, tt.want)
			}
		})
	}
}
