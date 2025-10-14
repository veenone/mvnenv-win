package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/cache"
	"github.com/veenone/mvnenv-win/internal/repository"
	"github.com/veenone/mvnenv-win/internal/version"
	"github.com/veenone/mvnenv-win/pkg/maven"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if current Maven version is the latest available",
	Long: `Check if the currently active Maven version is the latest available version.

This command compares your current Maven version with the latest version
available from configured repositories (Apache archive and Nexus if configured).`,
	Example: `  mvnenv status`,
	RunE:    runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()

	// Get current version
	resolver := version.NewVersionResolver(mvnenvRoot)
	resolved, err := resolver.ResolveVersion()
	if err != nil {
		if version.IsNoVersionSetError(err) {
			fmt.Println("No Maven version is currently set")
			fmt.Println("Set a version with: mvnenv global <version>")
			return nil
		}
		return formatError(err)
	}

	currentVersion := resolved.Version
	fmt.Printf("Current Maven version: %s (set by %s)\n", currentVersion, resolved.Source)

	// Get latest available version
	latestVersion, err := getLatestVersion(mvnenvRoot)
	if err != nil {
		return formatError(fmt.Errorf("failed to determine latest version: %w", err))
	}

	fmt.Printf("Latest available version: %s\n", latestVersion)

	// Compare versions
	if currentVersion == latestVersion {
		fmt.Println("\nYour Maven version is up to date")
	} else {
		fmt.Printf("\nA newer version is available: %s\n", latestVersion)
		fmt.Printf("Update with: mvnenv install %s && mvnenv global %s\n", latestVersion, latestVersion)
	}

	return nil
}

// getLatestVersion returns the latest available Maven version
func getLatestVersion(mvnenvRoot string) (string, error) {
	cacheManager := cache.NewManager(mvnenvRoot)

	// Try to load from cache first
	versions, err := cacheManager.LoadVersions()

	// If cache doesn't exist or is empty, fetch from repositories
	if err != nil || len(versions) == 0 {
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
