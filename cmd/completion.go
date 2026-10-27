package cmd

import (
	"github.com/sebrandon1/go-toolbox/cobrautil"
)

var completionCmd = cobrautil.NewCompletionCommand(rootCmd, cliName)
