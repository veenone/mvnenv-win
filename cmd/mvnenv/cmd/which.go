package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/version"
)

var whichCmd = &cobra.Command{
	Use:   "which <command>",
	Short: "Display the path to the Maven executable",
	Long:  `Show the full path to the Maven executable for the currently active version.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runWhich,
}

func init() {
	rootCmd.AddCommand(whichCmd)
}

func runWhich(cmd *cobra.Command, args []string) error {
	command := args[0]
	mvnenvRoot := getMvnenvRoot()

	resolver := version.NewVersionResolver(mvnenvRoot)
	resolved, err := resolver.ResolveVersion()
	if err != nil {
		if version.IsNoVersionSetError(err) {
			return fmt.Errorf("no Maven version is set")
		}
		if version.IsVersionNotInstalledError(err) {
			ver := version.ExtractVersionFromError(err)
			return fmt.Errorf("Maven version '%s' is set but not installed", ver)
		}
		return formatError(err)
	}

	// Construct path to command
	commandPath := filepath.Join(resolved.Path, "bin", command+".cmd")
	fmt.Println(commandPath)

	return nil
}
