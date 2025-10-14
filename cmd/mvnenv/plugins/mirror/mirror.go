// +build mirror

package mirror

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/veenone/mvnenv-win/cmd/mvnenv/plugins"
	"github.com/veenone/mvnenv-win/internal/config"
	"github.com/veenone/mvnenv-win/internal/download"
	"github.com/veenone/mvnenv-win/internal/nexus"
	"github.com/veenone/mvnenv-win/internal/repository"
	"github.com/veenone/mvnenv-win/pkg/maven"
)

func init() {
	plugins.Register(&MirrorPlugin{})
}

type MirrorPlugin struct{}

func (p *MirrorPlugin) Name() string {
	return "mirror"
}

func (p *MirrorPlugin) Description() string {
	return "Mirror Maven versions from Apache to Nexus repository"
}

func (p *MirrorPlugin) Command() *cobra.Command {
	var (
		dryRun      bool
		skipExisting bool
		maxVersions  int
	)

	cmd := &cobra.Command{
		Use:   "mirror",
		Short: "Mirror Maven versions from Apache to Nexus repository",
		Long: `Download all available Maven versions from Apache Maven archive and upload them
to a configured Nexus repository. This is useful for creating an internal mirror
of Maven distributions.

This command requires Nexus to be configured in the config file.`,
		Example: `  mvnenv mirror
  mvnenv mirror --dry-run
  mvnenv mirror --skip-existing
  mvnenv mirror --max 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMirror(dryRun, skipExisting, maxVersions)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be mirrored without uploading")
	cmd.Flags().BoolVar(&skipExisting, "skip-existing", true, "Skip versions that already exist in Nexus")
	cmd.Flags().IntVar(&maxVersions, "max", 0, "Maximum number of versions to mirror (0 = all)")

	return cmd
}

func runMirror(dryRun, skipExisting bool, maxVersions int) error {
	mvnenvRoot := getMvnenvRoot()

	// Load config
	configMgr := config.NewManager(mvnenvRoot)
	cfg, err := configMgr.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Verify Nexus is configured
	if cfg.Nexus == nil || !cfg.Nexus.Enabled {
		return fmt.Errorf("Nexus repository is not configured\nPlease configure Nexus in %s", filepath.Join(mvnenvRoot, "config", "config.yaml"))
	}

	fmt.Println("Maven Version Mirror")
	fmt.Println("===================")
	fmt.Printf("Source: Apache Maven Archive\n")
	fmt.Printf("Target: %s\n", cfg.Nexus.BaseURL)
	if dryRun {
		fmt.Println("Mode: DRY RUN (no uploads will be performed)")
	}
	fmt.Println()

	// Get list of available versions from Apache
	fmt.Println("Fetching available versions from Apache Maven archive...")
	apache := repository.NewApacheArchive()
	versions, err := apache.ListVersions()
	if err != nil {
		return fmt.Errorf("failed to list versions: %w", err)
	}

	// Sort versions (newest first)
	versions, err = maven.SortVersions(versions)
	if err != nil {
		return fmt.Errorf("failed to sort versions: %w", err)
	}

	// Limit versions if requested
	if maxVersions > 0 && len(versions) > maxVersions {
		versions = versions[:maxVersions]
		fmt.Printf("Limiting to %d most recent versions\n", maxVersions)
	}

	fmt.Printf("Found %d versions to mirror\n\n", len(versions))

	// Create Nexus client
	var nexusTLS *nexus.TLSConfig
	if cfg.Nexus.TLS != nil {
		nexusTLS = &nexus.TLSConfig{
			InsecureSkipVerify: cfg.Nexus.TLS.InsecureSkipVerify,
			CAFile:             cfg.Nexus.TLS.CAFile,
		}
	}
	nexusClient, err := nexus.NewClient(cfg.Nexus.BaseURL, cfg.Nexus.Username, cfg.Nexus.Password, nexusTLS)
	if err != nil {
		return fmt.Errorf("failed to create Nexus client: %w", err)
	}

	// Check which versions already exist in Nexus if skip-existing is enabled
	var existingVersions map[string]bool
	if skipExisting {
		fmt.Println("Checking existing versions in Nexus...")
		existing, err := nexusClient.ListVersions(context.Background())
		if err != nil {
			fmt.Printf("Warning: Could not check existing versions: %v\n", err)
			existingVersions = make(map[string]bool)
		} else {
			existingVersions = make(map[string]bool)
			for _, v := range existing {
				existingVersions[v] = true
			}
			fmt.Printf("Found %d existing versions in Nexus\n\n", len(existing))
		}
	}

	// Create temp directory for downloads
	tempDir := filepath.Join(mvnenvRoot, "cache", "mirror-temp")
	if !dryRun {
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			return fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(tempDir)
	}

	// Process each version
	successCount := 0
	skippedCount := 0
	failedCount := 0

	for i, version := range versions {
		fmt.Printf("[%d/%d] Processing Maven %s...\n", i+1, len(versions), version)

		// Skip if already exists
		if skipExisting && existingVersions[version] {
			fmt.Printf("  ✓ Already exists in Nexus, skipping\n\n")
			skippedCount++
			continue
		}

		if dryRun {
			fmt.Printf("  → Would download and upload to Nexus\n\n")
			successCount++
			continue
		}

		// Download from Apache
		archivePath := filepath.Join(tempDir, fmt.Sprintf("apache-maven-%s-bin.zip", version))
		fmt.Printf("  Downloading from Apache...")

		progressCallback := download.ProgressCallback(func(downloaded, total int64) {
			if total > 0 {
				percent := float64(downloaded) / float64(total) * 100
				fmt.Printf("\r  Downloading from Apache... %.1f%%", percent)
			}
		})

		err = apache.DownloadVersion(version, archivePath, progressCallback)
		if err != nil {
			fmt.Printf("\r  ✗ Download failed: %v\n\n", err)
			failedCount++
			continue
		}
		fmt.Printf("\r  ✓ Downloaded from Apache    \n")

		// Upload to Nexus
		fmt.Printf("  Uploading to Nexus...")

		uploadCallback := func(uploaded, total int64) {
			if total > 0 {
				percent := float64(uploaded) / float64(total) * 100
				fmt.Printf("\r  Uploading to Nexus... %.1f%%", percent)
			}
		}

		err = nexusClient.UploadVersion(context.Background(), version, archivePath, uploadCallback)
		if err != nil {
			fmt.Printf("\r  ✗ Upload failed: %v\n\n", err)
			failedCount++
			os.Remove(archivePath)
			continue
		}
		fmt.Printf("\r  ✓ Uploaded to Nexus      \n")

		// Clean up downloaded file
		os.Remove(archivePath)

		successCount++
		fmt.Printf("  ✓ Mirrored successfully\n\n")
	}

	// Print summary
	fmt.Println("Mirror Summary")
	fmt.Println("==============")
	if dryRun {
		fmt.Printf("Versions to mirror: %d\n", successCount)
		fmt.Printf("Already exists:     %d\n", skippedCount)
	} else {
		fmt.Printf("Successfully mirrored: %d\n", successCount)
		fmt.Printf("Skipped (existing):    %d\n", skippedCount)
		fmt.Printf("Failed:                %d\n", failedCount)
	}

	return nil
}

func getMvnenvRoot() string {
	if root := os.Getenv("MVNENV_ROOT"); root != "" {
		return root
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mvnenv")
}
