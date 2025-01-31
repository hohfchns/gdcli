package cmd

import (
	"fmt"

	"github.com/IgorBayerl/gdcli/internal/config"
	"github.com/IgorBayerl/gdcli/internal/core"
	"github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(installCmd())
}

func installCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "install [version]",
        Short: "Install Godot engine version",
        Run:   runInstall,
    }
}

func runInstall(cmd *cobra.Command, args []string) {
    var version string
    var dotnet bool

    if len(args) > 0 {
        version = args[0]
    } else {
        cfg, err := config.LoadConfig()
        if err != nil {
            fmt.Println("No version specified and no config found")
            return
        }
        version = cfg.EngineVersion
        dotnet = cfg.IsDotNet
    }

    if err := core.InstallGodotVersion(version, dotnet); err != nil {
        fmt.Printf("Installation failed: %v\n", err)
        return
    }
    
    fmt.Println("Godot installed successfully")
}