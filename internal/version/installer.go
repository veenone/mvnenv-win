package version

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/veenone/mvnenv-win/internal/repository"
)

// VersionInstaller handles Maven version installation
type VersionInstaller struct {
	mvnenvRoot    string
	repoManager   *repository.Manager
	resolver      *VersionResolver
	autoRehash    bool
	force         bool
	skipExisting  bool
	offline       bool
	quiet         bool
}

// NewVersionInstaller creates a new version installer
func NewVersionInstaller(mvnenvRoot string) *VersionInstaller {
	return &VersionInstaller{
		mvnenvRoot:  mvnenvRoot,
		repoManager: repository.NewManager(mvnenvRoot),
		resolver:    NewVersionResolver(mvnenvRoot),
		autoRehash:  true, // Enable automatic shim regeneration
		force:       false,
		skipExisting: false,
		offline:     false,
		quiet:       false,
	}
}

// SetForce sets the force flag (reinstall even if exists)
func (i *VersionInstaller) SetForce(force bool) {
	i.force = force
}

// SetSkipExisting sets the skip-existing flag
func (i *VersionInstaller) SetSkipExisting(skip bool) {
	i.skipExisting = skip
}

// SetOffline sets the offline flag
func (i *VersionInstaller) SetOffline(offline bool) {
	i.offline = offline
}

// SetQuiet sets the quiet flag
func (i *VersionInstaller) SetQuiet(quiet bool) {
	i.quiet = quiet
}

// InstallVersion installs a Maven version
func (i *VersionInstaller) InstallVersion(version string) error {
	// Check if already installed
	if i.resolver.IsVersionInstalled(version) {
		if i.skipExisting {
			if !i.quiet {
				fmt.Printf("Maven %s is already installed (skipped)\n", version)
			}
			return nil
		}
		if !i.force {
			return fmt.Errorf("Maven %s is already installed (use --force to reinstall)", version)
		}
		// Force reinstall: remove existing version first
		if !i.quiet {
			fmt.Printf("Maven %s already installed, reinstalling...\n", version)
		}
		if err := i.UninstallVersion(version); err != nil {
			return fmt.Errorf("failed to remove existing version: %w", err)
		}
	}

	// Create directories
	cacheDir := filepath.Join(i.mvnenvRoot, "cache")
	versionsDir := filepath.Join(i.mvnenvRoot, "versions")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("create cache directory: %w", err)
	}
	if err := os.MkdirAll(versionsDir, 0755); err != nil {
		return fmt.Errorf("create versions directory: %w", err)
	}

	// Check disk space (require at least 100MB for safety)
	requiredSpace := int64(100 * 1024 * 1024) // 100MB
	availableSpace, err := i.getAvailableDiskSpace(i.mvnenvRoot)
	if err != nil {
		// Warn but don't fail if we can't check disk space
		if !i.quiet {
			fmt.Printf("Warning: Could not check disk space: %v\n", err)
		}
	} else if availableSpace < requiredSpace {
		return fmt.Errorf("insufficient disk space: required %d MB, available %d MB",
			requiredSpace/(1024*1024), availableSpace/(1024*1024))
	}

	// Download to cache
	archivePath := filepath.Join(cacheDir, fmt.Sprintf("apache-maven-%s-bin.zip", version))

	var progress func(int64, int64)
	if !i.quiet {
		progress = func(downloaded, total int64) {
			if total > 0 {
				percent := float64(downloaded) / float64(total) * 100
				fmt.Printf("\rDownloading: %.1f%%", percent)
			}
		}
	}

	// Configure repository manager for offline mode
	if i.offline {
		i.repoManager.SetOfflineMode(true)
	}

	if err := i.repoManager.DownloadVersion(version, archivePath, progress); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	if !i.quiet && progress != nil {
		fmt.Println() // New line after progress
	}

	// Extract to versions directory
	if !i.quiet {
		fmt.Printf("Installing Maven %s...\n", version)
	}
	versionPath := filepath.Join(versionsDir, version)

	if err := i.extractZip(archivePath, versionsDir, version); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// Verify installation
	mvnCmd := filepath.Join(versionPath, "bin", "mvn.cmd")
	if _, err := os.Stat(mvnCmd); err != nil {
		return fmt.Errorf("installation verification failed: mvn.cmd not found")
	}

	if !i.quiet {
		fmt.Printf("Maven %s installed successfully\n", version)
	}

	// Automatically regenerate shims
	if i.autoRehash {
		if err := i.regenerateShims(); err != nil {
			fmt.Printf("Warning: Failed to regenerate shims: %v\n", err)
			fmt.Println("Run 'mvnenv rehash' manually to update shims")
		}
	}

	return nil
}

// UninstallVersion removes a Maven version
func (i *VersionInstaller) UninstallVersion(version string) error {
	// Check if installed
	if !i.resolver.IsVersionInstalled(version) {
		return fmt.Errorf("version '%s' not installed", version)
	}

	// Get version path
	versionPath := i.resolver.GetVersionPath(version)

	// Remove directory
	if err := os.RemoveAll(versionPath); err != nil {
		return fmt.Errorf("remove version directory: %w", err)
	}

	fmt.Printf("Maven %s uninstalled successfully\n", version)

	// Automatically regenerate shims
	if i.autoRehash {
		if err := i.regenerateShims(); err != nil {
			fmt.Printf("Warning: Failed to regenerate shims: %v\n", err)
			fmt.Println("Run 'mvnenv rehash' manually to update shims")
		}
	}

	return nil
}

// extractZip extracts a ZIP archive to a destination directory
func (i *VersionInstaller) extractZip(archivePath string, destDir string, version string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer r.Close()

	// Maven archives have a root directory like "apache-maven-3.9.4/"
	// We want to extract to "versions/3.9.4/" so we need to strip the root directory
	var rootPrefix string

	for _, f := range r.File {
		// Detect root prefix from first file
		if rootPrefix == "" {
			parts := strings.SplitN(f.Name, "/", 2)
			if len(parts) > 0 {
				rootPrefix = parts[0] + "/"
			}
		}

		// Skip if doesn't start with root prefix
		if !strings.HasPrefix(f.Name, rootPrefix) {
			continue
		}

		// Strip root prefix
		relativePath := strings.TrimPrefix(f.Name, rootPrefix)
		if relativePath == "" {
			continue
		}

		// Construct destination path
		destPath := filepath.Join(destDir, version, relativePath)

		// Create directory or extract file
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.Mode()); err != nil {
				return fmt.Errorf("create directory: %w", err)
			}
		} else {
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("create parent directory: %w", err)
			}

			// Extract file
			if err := i.extractFile(f, destPath); err != nil {
				return fmt.Errorf("extract file %s: %w", f.Name, err)
			}
		}
	}

	return nil
}

// extractFile extracts a single file from ZIP archive
func (i *VersionInstaller) extractFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

// regenerateShims regenerates shim executables
func (i *VersionInstaller) regenerateShims() error {
	// Import shim package dynamically to avoid circular dependency
	// We'll just silently skip if shim binary doesn't exist yet
	shimBinary := filepath.Join(i.mvnenvRoot, "bin", "shim.exe")
	if _, err := os.Stat(shimBinary); os.IsNotExist(err) {
		// Shim binary doesn't exist yet, skip regeneration
		return nil
	}

	// Note: This will be implemented when shim package is integrated
	// For now, we just check if shim.exe exists
	return nil
}
