package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

var (
	// quietMode suppresses non-error output when true
	quietMode bool
)

// shouldPrint checks if output should be printed based on quiet flag
func shouldPrint() bool {
	return !quietMode
}

// printVersion prints a version string to stdout if not in quiet mode
func printVersion(version string) {
	if shouldPrint() {
		fmt.Println(version)
	}
}

// printVersionList prints a list of versions with optional current marker
func printVersionList(versions []string, currentVersion string) {
	if !shouldPrint() {
		return
	}

	for _, v := range versions {
		if v == currentVersion {
			fmt.Printf("* %s\n", v)
		} else {
			fmt.Printf("  %s\n", v)
		}
	}
}

// printMessage prints a message to stdout if not in quiet mode
func printMessage(format string, args ...interface{}) {
	if shouldPrint() {
		fmt.Printf(format+"\n", args...)
	}
}

// printMessageNoNewline prints a message without newline if not in quiet mode
func printMessageNoNewline(format string, args ...interface{}) {
	if shouldPrint() {
		fmt.Printf(format, args...)
	}
}

// toWindowsPath converts a path to Windows format with backslashes
func toWindowsPath(path string) string {
	// Convert forward slashes to backslashes
	path = strings.ReplaceAll(path, "/", "\\")

	// Clean the path to remove any redundant separators
	path = filepath.Clean(path)

	return path
}

// formatPath formats a path for display, ensuring Windows backslash format
func formatPath(path string) string {
	return toWindowsPath(path)
}

// printPath prints a path in Windows format
func printPath(path string) {
	fmt.Println(formatPath(path))
}
