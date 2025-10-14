package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/veenone/mvnenv-win/cmd/mvnenv/cmd"
	_ "github.com/veenone/mvnenv-win/cmd/mvnenv/plugins/mirror" // Import plugins
)

// Version is the application version, can be overridden at build time with ldflags
var Version string

func main() {
	// If Version not set via ldflags, read from VERSION file
	if Version == "" {
		Version = readVersionFile()
	}

	// Set version in root command
	cmd.SetVersion(Version)

	// Execute root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// readVersionFile reads version from VERSION file as fallback
func readVersionFile() string {
	exePath, err := os.Executable()
	if err != nil {
		return "dev"
	}

	exeDir := filepath.Dir(exePath)
	versionFile := filepath.Join(exeDir, "VERSION")

	data, err := os.ReadFile(versionFile)
	if err != nil {
		// Try relative to working directory
		data, err = os.ReadFile("VERSION")
		if err != nil {
			return "dev"
		}
	}

	return string(data)
}
