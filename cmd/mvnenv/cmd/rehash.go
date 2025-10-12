package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/shim"
)

var rehashCmd = &cobra.Command{
	Use:   "rehash",
	Short: "Regenerate shim executables",
	Long: `Regenerate shim executables for all Maven commands.

Run this command after installing or uninstalling Maven versions to rebuild
the shim files that intercept Maven commands. This is typically done automatically
but can be run manually if needed.`,
	Example: `  mvnenv rehash`,
	RunE:    runRehash,
}

func init() {
	rootCmd.AddCommand(rehashCmd)
}

func runRehash(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()

	fmt.Println("Regenerating shims...")

	generator := shim.NewShimGenerator(mvnenvRoot)
	generatedPaths, err := generator.GenerateShims()
	if err != nil {
		return formatError(err)
	}

	// Extract command names from paths
	commands := make(map[string]bool)
	for _, path := range generatedPaths {
		// Extract base name without extension
		name := path[len(path)-len(".exe"):]
		if name[len(name)-4:] == ".exe" {
			name = name[:len(name)-4]
		} else if name[len(name)-4:] == ".cmd" {
			name = name[:len(name)-4]
		}
		// Get just the filename
		for i := len(name) - 1; i >= 0; i-- {
			if name[i] == '\\' || name[i] == '/' {
				name = name[i+1:]
				break
			}
		}
		commands[name] = true
	}

	fmt.Printf("Shims regenerated successfully (%d files)\n", len(generatedPaths))
	for cmd := range commands {
		fmt.Printf("  %s\n", cmd)
	}

	return nil
}
