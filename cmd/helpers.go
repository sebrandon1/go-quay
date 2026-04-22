package cmd

import (
	"fmt"
	"os"
)

func markFlagRequired(err error) {
	if err != nil {
		fmt.Printf("Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
