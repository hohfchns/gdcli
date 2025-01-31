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
    // Get current directory name
    wd, err := os.Getwd()
    if err != nil {
        fmt.Printf("Error getting current directory: %v\n", err)
        return
    }
    defaultProjectName := filepath.Base(wd)

    // Interactive prompts
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
            Prompt: &survey.Input{
                Message: "Godot version:",
                Default: "4.3.0",
            },
        },
        {
            Name: "dotnet",
            Prompt: &survey.Confirm{
                Message: "Use .NET version?",
                Default: false,
            },
        },
    }


    answers := struct {
        ProjectName string
        Version     string
        DotNet      bool
    }{}

    if err := survey.Ask(qs, &answers); err != nil {
        fmt.Printf("Error during survey: %v\n", err)
        return
    }

    if err := config.CreateConfig(answers.Version, answers.ProjectName, answers.DotNet); err != nil {
        fmt.Printf("Error creating config: %v\n", err)
        return
    }

    fmt.Println("Installing Godot version...")
    if err := core.InstallGodotVersion(answers.Version, answers.DotNet); err != nil {
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