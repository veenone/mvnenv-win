package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var appVersion string

var rootCmd = &cobra.Command{
	Use:   "mvnenv",
	Short: "Maven version manager for Windows",
	Long: `mvnenv-win is a command-line tool for managing multiple Apache Maven installations on Windows.
It allows you to easily switch between different Maven versions for various projects.`,
	SilenceUsage: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersion sets the application version
func SetVersion(v string) {
	appVersion = v
	rootCmd.Version = v
}

func init() {
	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// formatError formats error messages for consistent output
func formatError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%v", err)
}

// getMvnenvRoot returns the mvnenv installation root directory
func getMvnenvRoot() string {
	if root := os.Getenv("MVNENV_ROOT"); root != "" {
		return root
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user home directory: %v\n", err)
		os.Exit(1)
	}
	return home + "\\.mvnenv"
}
