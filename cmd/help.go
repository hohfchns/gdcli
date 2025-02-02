package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	// Override default help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Get help about any command",
		Long: `Help provides help for any command in the application.
Simply type gdcli help [path to command] for full details.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Root().HelpFunc()(cmd, args)
		},
	})
}