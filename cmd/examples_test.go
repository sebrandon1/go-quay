package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestMajorResourceCommandsShowExamplesInHelp(t *testing.T) {
	commands := []*cobra.Command{
		repositoryCmd,
		tagCmd,
		organizationCmd,
		teamCmd,
		robotCmd,
		permissionsCmd,
		secscanCmd,
		buildCmd,
		triggerCmd,
	}
	for _, command := range commands {
		t.Run(command.CommandPath(), func(t *testing.T) {
			if strings.TrimSpace(command.Example) == "" {
				t.Fatal("expected command example")
			}
			if !strings.Contains(command.UsageString(), command.Example) {
				t.Errorf("help output does not include the command examples:\n%s", command.UsageString())
			}
		})
	}
}
