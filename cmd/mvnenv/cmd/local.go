package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/version"
)

var localCmd = &cobra.Command{
	Use:   "local <version>",
	Short: "Set the local Maven version",
	Long:  `Set the Maven version for the current directory by creating a .maven-version file.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runLocal,
}

func init() {
	rootCmd.AddCommand(localCmd)
}

func runLocal(cmd *cobra.Command, args []string) error {
	ver := args[0]
	mvnenvRoot := getMvnenvRoot()

	// Verify version is installed
	resolver := version.NewVersionResolver(mvnenvRoot)
	if !resolver.IsVersionInstalled(ver) {
		return fmt.Errorf("version '%s' not installed", ver)
	}

	// Write .maven-version file in current directory
	versionFile := ".maven-version"
	if err := os.WriteFile(versionFile, []byte(ver), 0644); err != nil {
		return fmt.Errorf("failed to write .maven-version file: %w", err)
	}

	fmt.Printf("%s\n", ver)
	return nil
}
