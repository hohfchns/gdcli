// Package cmd implements the command-line interface for Godot Project CLI.
// This file contains the implementation of the 'init' command which initializes
// a new Godot project structure and configuration.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/IgorBayerl/gdcli/internal/config"
	"github.com/IgorBayerl/gdcli/internal/core"
	"github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(initCmd())
}

func initCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "init",
        Short: "Initialize new Godot project",
        Run:   runInit,
    }
}

func runInit(cmd *cobra.Command, args []string) {
    // Check if config already exists
    if _, err := os.Stat("gdproj.json"); err == nil {
        fmt.Println("Project already initialized. Run 'gdcli install' to install dependencies.")
        return
    } else if !os.IsNotExist(err) {
        fmt.Printf("Error checking for existing config: %v\n", err)
        return
    }

    // Get current directory name
    wd, err := os.Getwd()
    if err != nil {
        fmt.Printf("Error getting current directory: %v\n", err)
        return
    }
    defaultProjectName := filepath.Base(wd)

    var versionOptions []string
    currentOS := runtime.GOOS
    for _, v := range core.VersionManifest {
        if v.OS == currentOS {
            versionOptions = append(versionOptions, v.DisplayName)
        }
    }

    qs := []*survey.Question{
        {
            Name: "projectName",
            Prompt: &survey.Input{
                Message: "Project name:",
                Default: defaultProjectName,
            },
        },
        {
            Name: "version",
            Prompt: &survey.Select{
                Message: "Select Godot version:",
                Options: versionOptions,
                Default: versionOptions[0],
            },
        },
    }

    answers := struct {
        ProjectName string
        Version     string
    }{}

    if err := survey.Ask(qs, &answers); err != nil {
        fmt.Printf("Error during survey: %v\n", err)
        return
    }

    selected, err := core.GetVersionByIdentifier(answers.Version)
    if err != nil {
        fmt.Printf("Version selection error: %v\n", err)
        return
    }

    if err := config.CreateConfig(selected.Version, answers.ProjectName, selected.DotNet); err != nil {
        fmt.Printf("Error creating config: %v\n", err)
        return
    }

    fmt.Printf("Installing Godot %s...\n", selected.DisplayName)
    if err := core.InstallGodotVersion(selected); err != nil {
        fmt.Printf("Installation failed: %v\n", err)
        return
    }

    if err := createGodotProjectFile(answers.ProjectName); err != nil {
        fmt.Printf("Error creating project file: %v\n", err)
        return
    }

    updateGitignore()
    runOpen(cmd, args)
}

func createGodotProjectFile(projectName string) error {
    config := fmt.Sprintf(`[application]

config/name="%s"
run/main_scene=""
config/icon=""

[rendering]

environment/default_environment=""
`, projectName)
    return os.WriteFile("project.godot", []byte(config), 0644)
}

func updateGitignore() {
	linesToAdd := []string{
		"# Godot CLI",
		"dependencies/*",
		"",
		"# Godot 4+ specific ignores",
		".godot/",
		"",
		"# Godot-specific ignores",
		".import/",
		"export.cfg",
		"export_presets.cfg",
		"",
		"# Imported translations (automatically generated from CSV files)",
		"*.translation",
		"",
		"# Mono-specific ignores",
		".mono/",
		"data_*/",
		"mono_crash.*.json",
	}

	gitignorePath := ".gitignore"
	var existingLines map[string]bool = make(map[string]bool)

	// Check if .gitignore exists
	if _, err := os.Stat(gitignorePath); err == nil {
		// Read existing lines to avoid duplicates
		file, err := os.Open(gitignorePath)
		if err != nil {
			fmt.Printf("Error opening .gitignore: %v\n", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			existingLines[strings.TrimSpace(scanner.Text())] = true
		}
	}

	// Open .gitignore for appending
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening .gitignore: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range linesToAdd {
		if !existingLines[line] {
			if _, err := writer.WriteString(line + "\n"); err != nil {
				fmt.Printf("Error writing to .gitignore: %v\n", err)
				return
			}
		}
	}

	writer.Flush()
}
