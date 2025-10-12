package version

import (
	"errors"
	"fmt"
)

// Sentinel errors for version management operations
var (
	// ErrVersionNotInstalled indicates the requested version is not installed
	ErrVersionNotInstalled = errors.New("version not installed")

	// ErrVersionAlreadyInstalled indicates the version is already installed
	ErrVersionAlreadyInstalled = errors.New("version already installed")

	// ErrVersionNotSet indicates no Maven version is currently set
	ErrVersionNotSet = errors.New("no Maven version is set")

	// ErrInvalidVersion indicates the version string is invalid
	ErrInvalidVersion = errors.New("invalid version format")

	// ErrNoVersionsInstalled indicates no Maven versions are installed
	ErrNoVersionsInstalled = errors.New("no Maven versions installed")

	// ErrVersionInUse indicates the version cannot be uninstalled because it is in use
	ErrVersionInUse = errors.New("version is currently in use")

	// ErrDownloadFailed indicates a download operation failed
	ErrDownloadFailed = errors.New("download failed")

	// ErrChecksumMismatch indicates checksum verification failed
	ErrChecksumMismatch = errors.New("checksum verification failed")

	// ErrExtractionFailed indicates archive extraction failed
	ErrExtractionFailed = errors.New("extraction failed")

	// ErrInvalidMavenInstallation indicates the Maven installation is invalid
	ErrInvalidMavenInstallation = errors.New("invalid Maven installation")
)

// VersionNotInstalledError wraps ErrVersionNotInstalled with version details
type VersionNotInstalledError struct {
	Version string
}

func (e *VersionNotInstalledError) Error() string {
	return fmt.Sprintf("version '%s' is not installed", e.Version)
}

func (e *VersionNotInstalledError) Unwrap() error {
	return ErrVersionNotInstalled
}

// VersionAlreadyInstalledError wraps ErrVersionAlreadyInstalled with version details
type VersionAlreadyInstalledError struct {
	Version string
}

func (e *VersionAlreadyInstalledError) Error() string {
	return fmt.Sprintf("version '%s' is already installed", e.Version)
}

func (e *VersionAlreadyInstalledError) Unwrap() error {
	return ErrVersionAlreadyInstalled
}

// NoVersionSetError wraps ErrVersionNotSet with context
type NoVersionSetError struct {
	Message string
}

func (e *NoVersionSetError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "no Maven version is set"
}

func (e *NoVersionSetError) Unwrap() error {
	return ErrVersionNotSet
}

// Helper functions for error checking

// IsVersionNotInstalledError checks if error is a version not installed error
func IsVersionNotInstalledError(err error) bool {
	return errors.Is(err, ErrVersionNotInstalled)
}

// IsVersionAlreadyInstalledError checks if error is a version already installed error
func IsVersionAlreadyInstalledError(err error) bool {
	return errors.Is(err, ErrVersionAlreadyInstalled)
}

// IsNoVersionSetError checks if error is a no version set error
func IsNoVersionSetError(err error) bool {
	return errors.Is(err, ErrVersionNotSet)
}

// IsInvalidVersionError checks if error is an invalid version error
func IsInvalidVersionError(err error) bool {
	return errors.Is(err, ErrInvalidVersion)
}

// ExtractVersionFromError extracts the version string from a VersionNotInstalledError
func ExtractVersionFromError(err error) string {
	var vErr *VersionNotInstalledError
	if errors.As(err, &vErr) {
		return vErr.Version
	}
	return ""
}

// WrapDownloadError wraps a download error with context
func WrapDownloadError(err error, url string) error {
	return fmt.Errorf("%w: failed to download from %s: %v", ErrDownloadFailed, url, err)
}

// WrapChecksumError wraps a checksum error with context
func WrapChecksumError(version string) error {
	return fmt.Errorf("%w: checksum verification failed for version %s", ErrChecksumMismatch, version)
}

// WrapExtractionError wraps an extraction error with context
func WrapExtractionError(err error, file string) error {
	return fmt.Errorf("%w: failed to extract %s: %v", ErrExtractionFailed, file, err)
}

// NewVersionNotInstalledError creates a new version not installed error
func NewVersionNotInstalledError(version string) error {
	return &VersionNotInstalledError{Version: version}
}

// NewVersionAlreadyInstalledError creates a new version already installed error
func NewVersionAlreadyInstalledError(version string) error {
	return &VersionAlreadyInstalledError{Version: version}
}

// NewNoVersionSetError creates a new no version set error
func NewNoVersionSetError(message string) error {
	return &NoVersionSetError{Message: message}
}
