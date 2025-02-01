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
        Long: `Install a specific Godot version or use the version from config.
Examples:
  gdcli install 4.3.0-mono    # Install specific version
  gdcli install               # Use version from gdproj.json`,
        Run: runInstall,
    }
}

func runInstall(cmd *cobra.Command, args []string) {
    var version core.GodotVersion
    var err error

    if len(args) > 0 {
        // Install specified version
        version, err = core.GetVersionByIdentifier(args[0])
        if err != nil {
            fmt.Printf("❌ Version error: %v\n", err)
            fmt.Println("💡 Available versions:")
            for _, v := range core.VersionManifest {
                fmt.Printf("  - %s\n", v.DisplayName)
            }
            return
        }
    } else {
        // Try to use config version
        cfg, err := config.LoadConfig()
        if err != nil {
            fmt.Println("❌ No version specified and no config found")
            fmt.Println("💡 First create a project with: gdcli init")
            fmt.Println("   Or specify a version: gdcli install [version]")
            return
        }

        // Find config version in manifest
        found := false
        for _, v := range core.VersionManifest {
            if v.Version == cfg.EngineVersion && v.DotNet == cfg.IsDotNet {
                version = v
                found = true
                break
            }
        }

        if !found {
            fmt.Printf("❌ Configured version %s (%s) not found\n",
                cfg.EngineVersion,
                map[bool]string{true: "Mono", false: "Standard"}[cfg.IsDotNet],
            )
            fmt.Println("💡 Update your config or install manually:")
            fmt.Println("   gdcli install [version]")
            return
        }
    }

    fmt.Printf("🚀 Installing %s...\n", version.DisplayName)
    if err := core.InstallGodotVersion(version); err != nil {
        fmt.Printf("❌ Installation failed: %v\n", err)
        return
    }
    
    fmt.Printf("✅ Successfully installed %s\n", version.DisplayName)
    fmt.Println("🎮 Run your project with: gdcli open")
}