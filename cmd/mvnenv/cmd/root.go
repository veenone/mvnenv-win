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

It allows you to easily switch between different Maven versions for various projects
without manual PATH updates or system-wide configuration changes.

Features:
  - Install and manage multiple Maven versions
  - Switch between versions using shell, local, or global settings
  - Automatic command interception via shims
  - Compatible with pyenv-win conventions`,
	Example: `  mvnenv install 3.9.4
  mvnenv global 3.9.4
  mvnenv local 3.8.6
  mvn --version`,
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
