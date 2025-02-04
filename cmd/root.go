package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gdcli",
	Short: "Godot Project CLI - Manage Godot projects efficiently",
	Long: `A comprehensive CLI tool for managing Godot projects with features for
version management, project initialization, and workflow automation.`,
	Version: Version, // Version is set from main.go
}

func init() {
	cobra.AddTemplateFunc("rpad", func(s string, padding int) string {
		return fmt.Sprintf("%-*s", padding, s)
	})

	rootCmd.SetVersionTemplate("gdcli version {{.Version}}\n")

	rootCmd.AddCommand(&cobra.Command{
		Use:    "completion",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
