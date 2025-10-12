package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/veenone/mvnenv-win/internal/config"
)

// ResolvedVersion contains version resolution result
type ResolvedVersion struct {
	Version string
	Source  Source
	Path    string // Path to Maven installation
}

// Source indicates where the version was resolved from
type Source string

const (
	SourceShell  Source = "shell"
	SourceLocal  Source = "local"
	SourceGlobal Source = "global"
)

// VersionResolver resolves the active Maven version
type VersionResolver struct {
	mvnenvRoot    string
	configManager *config.Manager
}

// NewVersionResolver creates a new version resolver
func NewVersionResolver(mvnenvRoot string) *VersionResolver {
	return &VersionResolver{
		mvnenvRoot:    mvnenvRoot,
		configManager: config.NewManager(mvnenvRoot),
	}
}

// ResolveVersion resolves the active Maven version using shell > local > global hierarchy
func (r *VersionResolver) ResolveVersion() (*ResolvedVersion, error) {
	// 1. Check shell environment variable
	if version, ok := r.getShellVersion(); ok {
		if !r.isVersionInstalled(version) {
			return nil, &VersionError{
				Version: version,
				Source:  SourceShell,
				Err:     ErrVersionNotInstalled,
			}
		}
		return &ResolvedVersion{
			Version: version,
			Source:  SourceShell,
			Path:    r.getVersionPath(version),
		}, nil
	}

	// 2. Check .maven-version file (local)
	if version, ok := r.getLocalVersion(); ok {
		if !r.isVersionInstalled(version) {
			return nil, &VersionError{
				Version: version,
				Source:  SourceLocal,
				Err:     ErrVersionNotInstalled,
			}
		}
		return &ResolvedVersion{
			Version: version,
			Source:  SourceLocal,
			Path:    r.getVersionPath(version),
		}, nil
	}

	// 3. Check global configuration
	if version, ok := r.getGlobalVersion(); ok {
		if !r.isVersionInstalled(version) {
			return nil, &VersionError{
				Version: version,
				Source:  SourceGlobal,
				Err:     ErrVersionNotInstalled,
			}
		}
		return &ResolvedVersion{
			Version: version,
			Source:  SourceGlobal,
			Path:    r.getVersionPath(version),
		}, nil
	}

	return nil, NewNoVersionSetError("")
}

// getShellVersion reads version from MVNENV_MAVEN_VERSION environment variable
func (r *VersionResolver) getShellVersion() (string, bool) {
	version := strings.TrimSpace(os.Getenv("MVNENV_MAVEN_VERSION"))
	if version != "" {
		return version, true
	}
	return "", false
}

// getLocalVersion reads version from .maven-version file in current or parent directories
func (r *VersionResolver) getLocalVersion() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	for {
		versionFile := filepath.Join(dir, ".maven-version")
		if data, err := os.ReadFile(versionFile); err == nil {
			version := strings.TrimSpace(string(data))
			if version != "" {
				return version, true
			}
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}

	return "", false
}

// getGlobalVersion reads version from global configuration
func (r *VersionResolver) getGlobalVersion() (string, bool) {
	version, err := r.configManager.GetGlobalVersion()
	if err != nil || version == "" {
		return "", false
	}
	return version, true
}

// IsVersionInstalled checks if a Maven version is installed
func (r *VersionResolver) IsVersionInstalled(version string) bool {
	versionPath := r.GetVersionPath(version)
	mvnCmd := filepath.Join(versionPath, "bin", "mvn.cmd")
	_, err := os.Stat(mvnCmd)
	return err == nil
}

// GetVersionPath returns the installation path for a version
func (r *VersionResolver) GetVersionPath(version string) string {
	return filepath.Join(r.mvnenvRoot, "versions", version)
}

// isVersionInstalled is a private wrapper
func (r *VersionResolver) isVersionInstalled(version string) bool {
	return r.IsVersionInstalled(version)
}

// getVersionPath is a private wrapper
func (r *VersionResolver) getVersionPath(version string) string {
	return r.GetVersionPath(version)
}

// VersionError wraps version resolution errors with context
type VersionError struct {
	Version string
	Source  Source
	Err     error
}

func (e *VersionError) Error() string {
	return fmt.Sprintf("%s version '%s' %v", e.Source, e.Version, e.Err)
}

func (e *VersionError) Unwrap() error {
	return e.Err
}
