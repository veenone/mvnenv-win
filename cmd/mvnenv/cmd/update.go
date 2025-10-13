package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/cache"
	"github.com/veenone/mvnenv-win/internal/repository"
	"github.com/veenone/mvnenv-win/pkg/maven"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the cached list of available Maven versions",
	Long: `Fetches the latest list of available Maven versions from Apache archive and updates the local cache.

This command refreshes the cached version list used by the install -l command.
Run this periodically to see newly released Maven versions.`,
	Example: `  mvnenv update`,
	RunE:    runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()

	fmt.Println("Fetching available Maven versions from configured repositories...")

	// Fetch versions from all repositories (Apache + Nexus if configured)
	repoManager := repository.NewManager(mvnenvRoot)
	versions, err := repoManager.ListVersions()
	if err != nil {
		return formatError(err)
	}

	// Sort versions (newest first)
	sortedVersions, err := maven.SortVersions(versions)
	if err != nil {
		return formatError(fmt.Errorf("sort versions: %w", err))
	}

	// Save to cache
	cacheManager := cache.NewManager(mvnenvRoot)
	if err := cacheManager.SaveVersions(sortedVersions); err != nil {
		return formatError(fmt.Errorf("save cache: %w", err))
	}

	fmt.Printf("Successfully updated cache with %d Maven versions\n", len(sortedVersions))

	return nil
}
