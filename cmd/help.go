package cmd

import (
	"github.com/spf13/cobra"
)

const helpTemplate = `{{.Short}} (version: {{.Version}})

Usage:
  {{.UseLine}}

Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Options:
  {{rpad "--help" 10}} Display this help message
  {{rpad "--version" 10}} Display version information

Examples:
  gdcli init      Initialize new project
  gdcli install   Install Godot engine
  gdcli open      Launch Godot editor
`

func init() {
	rootCmd.SetHelpTemplate(helpTemplate)

	helpCmd := &cobra.Command{
		Use:   "help [command]",
		Short: "Get help about any command",
		Long: `Help provides help for any command in the application.
Simply type gdcli help [command] for full details.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				subCmd, _, err := rootCmd.Find(args)
				if err != nil || subCmd == nil {
					cmd.Printf("Unknown help topic: %s\n", args[0])
					if err := cmd.Help(); err != nil {
						cmd.Printf("Error: %v\n", err)
					}
					return
				}
				if err := subCmd.Help(); err != nil {
					cmd.Printf("Error: %v\n", err)
				}
			} else {
				if err := cmd.Root().Help(); err != nil {
					cmd.Printf("Error: %v\n", err)
				}
			}
		},
	}

	rootCmd.SetHelpCommand(helpCmd)
}
