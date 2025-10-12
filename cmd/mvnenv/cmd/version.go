package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current Maven version and origin",
	Long:  `Display the currently active Maven version and where it was set from (shell, local, or global).`,
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()
	resolver := version.NewVersionResolver(mvnenvRoot)

	resolved, err := resolver.ResolveVersion()
	if err != nil {
		if version.IsNoVersionSetError(err) {
			fmt.Println("No Maven version is set.")
			fmt.Println("Set a version with: mvnenv global <version>")
			return nil
		}
		if version.IsVersionNotInstalledError(err) {
			ver := version.ExtractVersionFromError(err)
			fmt.Printf("Maven version '%s' is set but not installed.\n", ver)
			fmt.Printf("Install it with: mvnenv install %s\n", ver)
			return nil
		}
		return formatError(err)
	}

	fmt.Printf("%s (set by %s)\n", resolved.Version, resolved.Source)
	return nil
}
