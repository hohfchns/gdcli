package main

import (
	"github.com/IgorBayerl/gdcli/cmd"
)

// Set via ldflags during build
var (
	version   = "dev"    // Should match the git tag
	commit    = "none"   // Git commit hash
	buildTime = "unknown"// Build timestamp
)

func main() {
	cmd.SetVersionInfo(version, commit, buildTime)
	cmd.Execute()
}