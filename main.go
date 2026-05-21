package main

import "github.com/sebrandon1/go-quay/cmd"

var version = "dev"

func main() {
	cmd.SetVersion(version)
	_ = cmd.Execute()
}
