package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/internal/cache"
	"github.com/veenone/mvnenv-win/internal/repository"
	versionpkg "github.com/veenone/mvnenv-win/internal/version"
	"github.com/veenone/mvnenv-win/pkg/maven"
)

var (
	latestRemote bool
)

var latestCmd = &cobra.Command{
	Use:   "latest [prefix]",
	Short: "Show the latest Maven version matching optional prefix",
	Long: `Show the latest installed Maven version, or with --remote flag, the latest available version.

Examples:
  mvnenv latest          # Latest installed version
  mvnenv latest 3.9      # Latest installed 3.9.x version
  mvnenv latest --remote # Latest available version from Apache archive
  mvnenv latest --remote 3.8 # Latest available 3.8.x version`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLatest,
}

func init() {
	latestCmd.Flags().BoolVarP(&latestRemote, "remote", "r", false, "Check remote versions from Apache archive")
	rootCmd.AddCommand(latestCmd)
}

func runLatest(cmd *cobra.Command, args []string) error {
	mvnenvRoot := getMvnenvRoot()
	var prefix string
	if len(args) > 0 {
		prefix = args[0]
	}

	var versions []string
	var err error

	if latestRemote {
		// Get remote versions
		cacheManager := cache.NewManager(mvnenvRoot)

		// Try to load from cache first
		versions, err = cacheManager.LoadVersions()
		if err != nil || len(versions) == 0 {
			// Cache doesn't exist or is empty, fetch from Apache
			fmt.Println("Fetching versions from Apache archive...")
			archive := repository.NewApacheArchive()
			versions, err = archive.ListVersions()
			if err != nil {
				return formatError(fmt.Errorf("fetch versions: %w", err))
			}

			// Save to cache for future use
			if saveErr := cacheManager.SaveVersions(versions); saveErr != nil {
				fmt.Printf("Warning: Failed to save cache: %v\n", saveErr)
			}
		}
	} else {
		// Get installed versions
		lister := versionpkg.NewVersionLister(mvnenvRoot)
		versions, err = lister.ListInstalled()
		if err != nil {
			return formatError(err)
		}
	}

	if len(versions) == 0 {
		if latestRemote {
			return formatError(fmt.Errorf("no versions available from Apache archive"))
		}
		return formatError(fmt.Errorf("no versions installed"))
	}

	// Sort versions (newest first)
	sortedVersions, err := maven.SortVersions(versions)
	if err != nil {
		return formatError(fmt.Errorf("sort versions: %w", err))
	}

	// Filter by prefix if provided
	var filtered []string
	if prefix != "" {
		for _, v := range sortedVersions {
			if strings.HasPrefix(v, prefix) {
				filtered = append(filtered, v)
			}
		}

		if len(filtered) == 0 {
			if latestRemote {
				return formatError(fmt.Errorf("no remote versions match prefix '%s'", prefix))
			}
			return formatError(fmt.Errorf("no installed versions match prefix '%s'", prefix))
		}

		sortedVersions = filtered
	}

	// Return the latest (first after sorting newest-first)
	latest := sortedVersions[0]
	fmt.Println(latest)

	return nil
}
