package version

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/veenone/mvnenv-win/pkg/maven"
)

// VersionLister lists installed Maven versions
type VersionLister struct {
	mvnenvRoot string
	resolver   *VersionResolver
}

// NewVersionLister creates a new version lister
func NewVersionLister(mvnenvRoot string) *VersionLister {
	return &VersionLister{
		mvnenvRoot: mvnenvRoot,
		resolver:   NewVersionResolver(mvnenvRoot),
	}
}

// ListInstalled returns a list of installed Maven versions
func (l *VersionLister) ListInstalled() ([]string, error) {
	versionsDir := filepath.Join(l.mvnenvRoot, "versions")

	// Check if versions directory exists
	if _, err := os.Stat(versionsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read versions directory
	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, fmt.Errorf("read versions directory: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		version := entry.Name()
		// Verify it's a valid Maven installation
		mvnCmd := filepath.Join(versionsDir, version, "bin", "mvn.cmd")
		if _, err := os.Stat(mvnCmd); err == nil {
			versions = append(versions, version)
		}
	}

	// Sort versions (newest first)
	if len(versions) > 0 {
		sorted, err := maven.SortVersions(versions)
		if err != nil {
			// If sorting fails, return unsorted
			return versions, nil
		}
		versions = sorted
	}

	return versions, nil
}

// GetCurrentVersion returns the currently active version or empty string
func (l *VersionLister) GetCurrentVersion() string {
	resolved, err := l.resolver.ResolveVersion()
	if err != nil {
		return ""
	}
	return resolved.Version
}
