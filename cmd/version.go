package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   string
	Commit    string
	BuildTime string
)

func SetVersionInfo(version, commit, buildTime string) {
	Version = version
	Commit = commit
	BuildTime = buildTime
}

func init() {
	rootCmd.AddCommand(versionCmd())
}

func versionCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "version",
        Short: "Print version information",
        Run: func(cmd *cobra.Command, args []string) {
            full, _ := cmd.Flags().GetBool("full")
            if full {
                fmt.Printf("gdcli version %s\nCommit: %s\nBuild time: %s\n", Version, Commit, BuildTime)
            } else {
                fmt.Println(Version)
            }
        },
    }
    cmd.Flags().Bool("full", false, "Show detailed version information")
    return cmd
}

