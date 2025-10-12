package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/cache"
	"github.com/veenone/mvnenv-win/internal/repository"
	versionpkg "github.com/veenone/mvnenv-win/internal/version"
	"github.com/veenone/mvnenv-win/pkg/maven"
)

var (
	installList  bool
	installQuiet bool
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a Maven version",
	Long: `Download and install a specific Maven version from Apache archive.

Downloads the specified Maven version, verifies its checksum, and installs
it to the mvnenv versions directory. Use the -l flag to list all available
versions.`,
	Example: `  mvnenv install 3.9.4
  mvnenv install -l
  mvnenv install -q 3.8.6`,
	RunE: runInstall,
}

func init() {
	installCmd.Flags().BoolVarP(&installList, "list", "l", false, "List available versions")
	installCmd.Flags().BoolVarP(&installQuiet, "quiet", "q", false, "Suppress output")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()

	// Set quiet mode
	quietMode = installQuiet

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
	mvnenvRoot := getMvnenvRoot()
	cacheManager := cache.NewManager(mvnenvRoot)

	// Try to load from cache first
	versions, err := cacheManager.LoadVersions()

	// If cache doesn't exist or is stale (>24 hours), fetch from Apache
	if err != nil || len(versions) == 0 || cacheManager.IsCacheStale(24*time.Hour) {
		fmt.Println("Fetching available versions from Apache archive...")

		archive := repository.NewApacheArchive()
		versions, err = archive.ListVersions()
		if err != nil {
			return fmt.Errorf("failed to list versions: %w", err)
		}

		// Sort versions (newest first)
		versions, err = maven.SortVersions(versions)
		if err != nil {
			// If sorting fails, use unsorted
			fmt.Printf("Warning: Failed to sort versions: %v\n", err)
		}

		// Save to cache
		if saveErr := cacheManager.SaveVersions(versions); saveErr != nil {
			fmt.Printf("Warning: Failed to save cache: %v\n", saveErr)
		}
	} else {
		age, _ := cacheManager.GetCacheAge()
		fmt.Printf("Using cached versions (updated %v ago)\n", age.Round(time.Minute))
		fmt.Println("Run 'mvnenv update' to refresh the cache")
	}

	if len(versions) == 0 {
		fmt.Println("No versions found")
		return nil
	}

	fmt.Printf("\nAvailable Maven versions:\n")
	for _, v := range versions {
		fmt.Printf("  %s\n", v)
	}

	return nil
}
