package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/version"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "List all installed Maven versions",
	Long: `Display all Maven versions that are currently installed.

Lists all installed Maven versions with the currently active version
marked with an asterisk (*).`,
	Example: `  mvnenv versions`,
	RunE:    runVersions,
}

func init() {
	rootCmd.AddCommand(versionsCmd)
}

func runVersions(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()
	lister := version.NewVersionLister(mvnenvRoot)

	versions, err := lister.ListInstalled()
	if err != nil {
		return formatError(err)
	}

	if len(versions) == 0 {
		fmt.Println("No Maven versions installed.")
		fmt.Println("Install a version with: mvnenv install <version>")
		return nil
	}

	// Get current version
	currentVersion := lister.GetCurrentVersion()

	// Display versions
	for _, v := range versions {
		if v == currentVersion {
			fmt.Printf("* %s\n", v)
		} else {
			fmt.Printf("  %s\n", v)
		}
	}

	return nil
}
