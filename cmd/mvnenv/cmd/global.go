package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/config"
	"github.com/veenone/mvnenv-win/internal/version"
)

var globalCmd = &cobra.Command{
	Use:   "global <version>",
	Short: "Set the global Maven version",
	Long: `Set the global Maven version that will be used by default in all directories.

The global version is stored in the mvnenv configuration file and serves as
the fallback when no local or shell version is set.`,
	Example: `  mvnenv global 3.9.4
  mvnenv global 3.8.6`,
	Args: cobra.ExactArgs(1),
	RunE: runGlobal,
}

func init() {
	rootCmd.AddCommand(globalCmd)
}

func runGlobal(cmd *cobra.Command, args []string) error {
	ver := args[0]
	mvnenvRoot := getMvnenvRoot()

	// Verify version is installed
	resolver := version.NewVersionResolver(mvnenvRoot)
	if !resolver.IsVersionInstalled(ver) {
		return fmt.Errorf("version '%s' not installed", ver)
	}

	// Set global version
	configMgr := config.NewManager(mvnenvRoot)
	if err := configMgr.SetGlobalVersion(ver); err != nil {
		return fmt.Errorf("failed to set global version: %w", err)
	}

	fmt.Printf("%s\n", ver)
	return nil
}
