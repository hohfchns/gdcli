// cmd/clean.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd())
}

func cleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Remove generated files and dependencies",
		Run:   runClean,
	}
}

func runClean(cmd *cobra.Command, args []string) {
	// Remove dependencies directory
	if err := os.RemoveAll("dependencies"); err != nil {
		fmt.Printf("Error cleaning dependencies: %v\n", err)
	} else {
		fmt.Println("Cleaned dependencies directory")
	}

	if err := os.RemoveAll(".godot"); err != nil {
		fmt.Printf("Error cleaning .godot: %v\n", err)
	} else {
		fmt.Println("Cleaned .godot directory")
	}

	// Add more cleanup tasks here if needed
	// Example: os.Remove("other-generated-file")
}