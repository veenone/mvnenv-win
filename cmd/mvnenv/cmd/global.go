package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/config"
	"github.com/veenone/mvnenv-win/internal/version"
)

var (
	globalUnset bool
)

var globalCmd = &cobra.Command{
	Use:   "global [version]",
	Short: "Set or show the global Maven version",
	Long: `Set or show the global default Maven version.

The global version is used when no shell-specific version (MVNENV_MAVEN_VERSION)
or local version (.maven-version file) is set. This provides a system-wide
default Maven version.

Version Resolution Hierarchy:
  1. Shell: MVNENV_MAVEN_VERSION environment variable (highest priority)
  2. Local: .maven-version file in current or parent directory
  3. Global: Set via this command (lowest priority)`,
	Example: `  # Show current global version
  mvnenv global

  # Set global version
  mvnenv global 3.9.4

  # Unset global version
  mvnenv global --unset`,
	RunE: runGlobal,
}

func init() {
	rootCmd.AddCommand(globalCmd)
	globalCmd.Flags().BoolVar(&globalUnset, "unset", false, "Remove the global version setting")
}

func runGlobal(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()
	configMgr := config.NewManager(mvnenvRoot)

	// Case 1: Unset global version
	if globalUnset {
		if err := configMgr.UnsetGlobalVersion(); err != nil {
			return formatError(fmt.Errorf("failed to unset global version: %w", err))
		}
		fmt.Println("Global Maven version unset")
		return nil
	}

	// Case 2: Display current global version
	if len(args) == 0 {
		globalVersion, err := configMgr.GetGlobalVersion()
		if err != nil {
			return formatError(fmt.Errorf("failed to read configuration: %w", err))
		}

		if globalVersion == "" {
			fmt.Println("No global Maven version set (use 'mvnenv global <version>')")
		} else {
			fmt.Println(globalVersion)
		}
		return nil
	}

	// Case 3: Set global version
	if len(args) != 1 {
		return formatError(fmt.Errorf("expected 0 or 1 argument, got %d", len(args)))
	}

	newVersion := args[0]

	// Validate version format
	if err := validateVersionFormat(newVersion); err != nil {
		return formatError(err)
	}

	// Validate version is installed
	resolver := version.NewVersionResolver(mvnenvRoot)
	if !resolver.IsVersionInstalled(newVersion) {
		return formatError(fmt.Errorf("Maven %s is not installed (use 'mvnenv install %s' first)", newVersion, newVersion))
	}

	// Set global version
	if err := configMgr.SetGlobalVersion(newVersion); err != nil {
		return formatError(fmt.Errorf("failed to set global version: %w", err))
	}

	fmt.Printf("Global Maven version set to %s\n", newVersion)
	return nil
}

// validateVersionFormat validates that version string is safe and valid
func validateVersionFormat(ver string) error {
	if ver == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Check for path traversal attempts
	if containsPathTraversal(ver) {
		return fmt.Errorf("invalid version format: '%s' (must contain only alphanumeric characters, dots, and hyphens)", ver)
	}

	// Validate characters (alphanumeric, dots, hyphens only)
	for _, ch := range ver {
		if !isValidVersionChar(ch) {
			return fmt.Errorf("invalid version format: '%s' (must contain only alphanumeric characters, dots, and hyphens)", ver)
		}
	}

	return nil
}

// containsPathTraversal checks for path traversal attempts
func containsPathTraversal(s string) bool {
	return len(s) >= 2 && (s[0] == '.' && (s[1] == '.' || s[1] == '/' || s[1] == '\\'))
}

// isValidVersionChar checks if character is valid in version string
func isValidVersionChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '.' ||
		ch == '-'
}
