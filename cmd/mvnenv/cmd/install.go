package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/repository"
	versionpkg "github.com/veenone/mvnenv-win/internal/version"
	"github.com/veenone/mvnenv-win/pkg/maven"
)

var (
	installList bool
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a Maven version",
	Long:  `Download and install a specific Maven version.`,
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().BoolVarP(&installList, "list", "l", false, "List available versions")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()

	// Handle list flag
	if installList {
		return listAvailableVersions()
	}

	// Require version argument
	if len(args) == 0 {
		return fmt.Errorf("version argument required\nUsage: mvnenv install <version>\nList available versions with: mvnenv install -l")
	}

	version := args[0]

	// Install version
	installer := versionpkg.NewVersionInstaller(mvnenvRoot)
	if err := installer.InstallVersion(version); err != nil {
		return formatError(err)
	}

	return nil
}

func listAvailableVersions() error {
	fmt.Println("Fetching available versions from Apache archive...")

	archive := repository.NewApacheArchive()
	versions, err := archive.ListVersions()
	if err != nil {
		return fmt.Errorf("failed to list versions: %w", err)
	}

	if len(versions) == 0 {
		fmt.Println("No versions found")
		return nil
	}

	// Sort versions (newest first)
	sorted, err := maven.SortVersions(versions)
	if err != nil {
		// If sorting fails, use unsorted
		sorted = versions
	}

	fmt.Printf("\nAvailable Maven versions:\n")
	for _, v := range sorted {
		fmt.Printf("  %s\n", v)
	}

	return nil
}
