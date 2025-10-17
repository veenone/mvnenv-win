package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/cache"
	"github.com/veenone/mvnenv-win/internal/repository"
	versionpkg "github.com/veenone/mvnenv-win/internal/version"
	"github.com/veenone/mvnenv-win/pkg/maven"
)

var (
	installList         bool
	installQuiet        bool
	installForce        bool
	installSkipExisting bool
	installClear        bool
	installOffline      bool
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
	installCmd.Flags().BoolVarP(&installForce, "force", "f", false, "Force reinstall if version already exists")
	installCmd.Flags().BoolVarP(&installSkipExisting, "skip-existing", "s", false, "Skip installation if version already exists (no error)")
	installCmd.Flags().BoolVarP(&installClear, "clear", "c", false, "Clear cache before installing")
	installCmd.Flags().BoolVar(&installOffline, "offline", false, "Offline mode: only use Nexus (fail if unavailable)")
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
		return fmt.Errorf("version argument required\nUsage: mvnenv install <version> [version2 ...]\nList available versions with: mvnenv install -l")
	}

	// Clear cache if requested
	if installClear {
		if err := clearInstallCache(mvnenvRoot); err != nil {
			fmt.Printf("Warning: Failed to clear cache: %v\n", err)
		} else if !installQuiet {
			fmt.Println("Cache cleared")
		}
	}

	// Handle multiple version installation
	var successfulInstalls []string
	var failedInstalls []string

	for _, version := range args {
		// Handle "latest" keyword
		if version == "latest" {
			latestVersion, err := getLatestAvailableVersion(mvnenvRoot)
			if err != nil {
				failedInstalls = append(failedInstalls, fmt.Sprintf("%s (failed to determine latest: %v)", version, err))
				continue
			}
			version = latestVersion
			if !installQuiet {
				fmt.Printf("Installing latest Maven version: %s\n", version)
			}
		}

		// Install version with flags
		if err := installSingleVersion(mvnenvRoot, version); err != nil {
			failedInstalls = append(failedInstalls, fmt.Sprintf("%s (%v)", version, err))
		} else {
			successfulInstalls = append(successfulInstalls, version)
		}
	}

	// Report results
	if !installQuiet && len(args) > 1 {
		fmt.Println("\nInstallation Summary:")
		if len(successfulInstalls) > 0 {
			fmt.Printf("✓ Successfully installed: %v\n", successfulInstalls)
		}
		if len(failedInstalls) > 0 {
			fmt.Printf("✗ Failed: %v\n", failedInstalls)
		}
	}

	// Return error if any installations failed
	if len(failedInstalls) > 0 {
		if len(successfulInstalls) == 0 {
			return fmt.Errorf("all installations failed")
		}
		return fmt.Errorf("partial success: %d succeeded, %d failed", len(successfulInstalls), len(failedInstalls))
	}

	return nil
}

// installSingleVersion installs a single Maven version with flag handling
func installSingleVersion(mvnenvRoot, version string) error {
	installer := versionpkg.NewVersionInstaller(mvnenvRoot)

	// Configure installer based on flags
	installer.SetForce(installForce)
	installer.SetSkipExisting(installSkipExisting)
	installer.SetOffline(installOffline)
	installer.SetQuiet(installQuiet)

	// Install version
	if err := installer.InstallVersion(version); err != nil {
		return err
	}

	return nil
}

// clearInstallCache clears the download cache
func clearInstallCache(mvnenvRoot string) error {
	cacheDir := filepath.Join(mvnenvRoot, "cache")
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Cache doesn't exist, nothing to clear
		}
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		// Only delete .zip files, keep versions.json
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".zip" {
			path := filepath.Join(cacheDir, entry.Name())
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s: %w", entry.Name(), err)
			}
		}
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
