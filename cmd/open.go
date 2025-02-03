package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(openCmd())
}

func openCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open project in Godot editor",
		Run:   runOpen,
	}
}

func runOpen(cobraCmd *cobra.Command, args []string) {
    godotPath := filepath.Join("dependencies", "godot.exe")
    
    if _, err := os.Stat(godotPath); os.IsNotExist(err) {
        fmt.Printf("Godot executable not found at %s\n", godotPath)
        fmt.Println("Run 'gdcli install' to install the required version")
        return
    }

    // Initialize project if not exists
    if _, err := os.Stat("project.godot"); os.IsNotExist(err) {
        fmt.Println("Initializing new Godot project...")
        initCmd := exec.Command(godotPath, "--path", ".", "--editor")
        initCmd.Stdout = os.Stdout
        initCmd.Stderr = os.Stderr
        if err := initCmd.Run(); err != nil {
            fmt.Printf("Failed to initialize project: %v\n", err)
            return
        }
    }

    godotCmd := exec.Command(godotPath, "--path", ".", "--editor")
    godotCmd.Stdout = nil
    godotCmd.Stderr = nil
    godotCmd.Stdin = nil

    if err := godotCmd.Start(); err != nil {
        fmt.Printf("Error launching Godot: %v\n", err)
    } else {
        fmt.Println("Godot editor launched successfully and detached from terminal.")
    }
}
