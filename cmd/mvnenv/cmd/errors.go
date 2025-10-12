package cmd

import (
	"errors"
	"fmt"
	"os"
)

// Common error types
var (
	// ErrVersionNotInstalled indicates the requested version is not installed
	ErrVersionNotInstalled = errors.New("version not installed")

	// ErrVersionNotSet indicates no Maven version is currently set
	ErrVersionNotSet = errors.New("no Maven version is set")

	// ErrVersionAlreadyInstalled indicates the version is already installed
	ErrVersionAlreadyInstalled = errors.New("version already installed")

	// ErrInvalidVersion indicates the version string is invalid
	ErrInvalidVersion = errors.New("invalid version")

	// ErrNoVersionsInstalled indicates no Maven versions are installed
	ErrNoVersionsInstalled = errors.New("no Maven versions installed")

	// ErrCommandNotFound indicates the Maven command was not found
	ErrCommandNotFound = errors.New("command not found")

	// ErrNetworkFailure indicates a network operation failed
	ErrNetworkFailure = errors.New("network operation failed")

	// ErrCacheFailure indicates cache operation failed
	ErrCacheFailure = errors.New("cache operation failed")
)

// formatError formats error messages with consistent "Error: " prefix for stderr
func formatError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("Error: %v", err)
}

// printError prints an error message to stderr with consistent formatting
func printError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}

// exitWithError prints an error and exits with code 1
func exitWithError(format string, args ...interface{}) {
	printError(format, args...)
	os.Exit(1)
}
