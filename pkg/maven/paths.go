package maven

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetMavenBinaryPath returns the path to the mvn.cmd executable for a Maven installation
func GetMavenBinaryPath(mavenHome string) string {
	return filepath.Join(mavenHome, "bin", "mvn.cmd")
}

// GetMavenDebugBinaryPath returns the path to the mvnDebug.cmd executable
func GetMavenDebugBinaryPath(mavenHome string) string {
	return filepath.Join(mavenHome, "bin", "mvnDebug.cmd")
}

// GetMavenHome returns the MAVEN_HOME path for a version installation
func GetMavenHome(versionsDir, version string) string {
	return filepath.Join(versionsDir, version)
}

// ValidateMavenInstallation checks if a directory contains a valid Maven installation
func ValidateMavenInstallation(mavenHome string) error {
	// Check if directory exists
	info, err := os.Stat(mavenHome)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Maven installation directory does not exist: %s", mavenHome)
		}
		return fmt.Errorf("cannot access Maven installation directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("Maven installation path is not a directory: %s", mavenHome)
	}

	// Check for bin/mvn.cmd
	mvnCmd := GetMavenBinaryPath(mavenHome)
	if _, err := os.Stat(mvnCmd); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Maven installation is invalid: bin/mvn.cmd not found in %s", mavenHome)
		}
		return fmt.Errorf("cannot access Maven binary: %w", err)
	}

	return nil
}

// GetBinDirectory returns the path to the bin directory in a Maven installation
func GetBinDirectory(mavenHome string) string {
	return filepath.Join(mavenHome, "bin")
}

// GetLibDirectory returns the path to the lib directory in a Maven installation
func GetLibDirectory(mavenHome string) string {
	return filepath.Join(mavenHome, "lib")
}

// GetConfDirectory returns the path to the conf directory in a Maven installation
func GetConfDirectory(mavenHome string) string {
	return filepath.Join(mavenHome, "conf")
}
