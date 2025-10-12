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
	mvnenvRoot string
	repository *repository.ApacheArchive
	resolver   *VersionResolver
}

// NewVersionInstaller creates a new version installer
func NewVersionInstaller(mvnenvRoot string) *VersionInstaller {
	return &VersionInstaller{
		mvnenvRoot: mvnenvRoot,
		repository: repository.NewApacheArchive(),
		resolver:   NewVersionResolver(mvnenvRoot),
	}
}

// InstallVersion installs a Maven version
func (i *VersionInstaller) InstallVersion(version string) error {
	// Check if already installed
	if i.resolver.IsVersionInstalled(version) {
		return fmt.Errorf("version '%s' already installed", version)
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

	// Download to cache
	archivePath := filepath.Join(cacheDir, fmt.Sprintf("apache-maven-%s-bin.zip", version))

	progress := func(downloaded, total int64) {
		if total > 0 {
			percent := float64(downloaded) / float64(total) * 100
			fmt.Printf("\rDownloading: %.1f%%", percent)
		}
	}

	if err := i.repository.DownloadVersion(version, archivePath, progress); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	fmt.Println() // New line after progress

	// Extract to versions directory
	fmt.Printf("Installing Maven %s...\n", version)
	versionPath := filepath.Join(versionsDir, version)

	if err := i.extractZip(archivePath, versionsDir, version); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	// Verify installation
	mvnCmd := filepath.Join(versionPath, "bin", "mvn.cmd")
	if _, err := os.Stat(mvnCmd); err != nil {
		return fmt.Errorf("installation verification failed: mvn.cmd not found")
	}

	fmt.Printf("Maven %s installed successfully\n", version)
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
