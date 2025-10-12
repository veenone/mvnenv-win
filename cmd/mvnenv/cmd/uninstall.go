package cmd

import (
	"github.com/spf13/cobra"
	versionpkg "github.com/veenone/mvnenv-win/internal/version"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall a Maven version",
	Long:  `Remove an installed Maven version.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) error {
	ver := args[0]
	mvnenvRoot := getMvnenvRoot()

	// Uninstall version
	installer := versionpkg.NewVersionInstaller(mvnenvRoot)
	if err := installer.UninstallVersion(ver); err != nil {
		return formatError(err)
	}

	return nil
}
