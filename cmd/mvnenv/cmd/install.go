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
versions.

Use "latest" as the version to install the newest available Maven version.`,
	Example: `  mvnenv install 3.9.4
  mvnenv install latest
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

	// Handle "latest" keyword
	if version == "latest" {
		latestVersion, err := getLatestAvailableVersion(mvnenvRoot)
		if err != nil {
			return formatError(fmt.Errorf("failed to determine latest version: %w", err))
		}
		version = latestVersion
		fmt.Printf("Installing latest Maven version: %s\n", version)
	}

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

	// If cache doesn't exist or is stale (>24 hours), fetch from repositories
	if err != nil || len(versions) == 0 || cacheManager.IsCacheStale(24*time.Hour) {
		fmt.Println("Fetching available versions from configured repositories...")

		repoManager := repository.NewManager(mvnenvRoot)
		versions, err = repoManager.ListVersions()
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

// getLatestAvailableVersion returns the latest available Maven version
func getLatestAvailableVersion(mvnenvRoot string) (string, error) {
	cacheManager := cache.NewManager(mvnenvRoot)

	// Try to load from cache first
	versions, err := cacheManager.LoadVersions()

	// If cache doesn't exist or is stale (>24 hours), fetch from repositories
	if err != nil || len(versions) == 0 || cacheManager.IsCacheStale(24*time.Hour) {
		repoManager := repository.NewManager(mvnenvRoot)
		versions, err = repoManager.ListVersions()
		if err != nil {
			return "", fmt.Errorf("failed to list versions: %w", err)
		}

		// Sort versions (newest first)
		versions, err = maven.SortVersions(versions)
		if err != nil {
			return "", fmt.Errorf("failed to sort versions: %w", err)
		}

		// Save to cache
		if saveErr := cacheManager.SaveVersions(versions); saveErr != nil {
			// Just warn, don't fail
			fmt.Printf("Warning: Failed to save cache: %v\n", saveErr)
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions available")
	}

	// Return the first version (newest)
	return versions[0], nil
}
