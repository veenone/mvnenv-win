package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/veenone/mvnenv-win/internal/shim"
	versionpkg "github.com/veenone/mvnenv-win/internal/version"
)

func main() {
	// Detect command name from executable name
	command := detectCommand()

	// Get mvnenv root
	mvnenvRoot := getMvnenvRoot()

	// Create resolver and executor
	resolver := versionpkg.NewVersionResolver(mvnenvRoot)
	executor := shim.NewShimExecutor(resolver)

	// Execute Maven command
	exitCode, err := executor.Execute(command, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "mvnenv: %v\n", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

// detectCommand determines which Maven command was invoked
func detectCommand() string {
	exePath, err := os.Executable()
	if err != nil {
		return "mvn"
	}

	base := filepath.Base(exePath)
	name := strings.TrimSuffix(base, ".exe")
	return name
}

// getMvnenvRoot returns the mvnenv installation root
func getMvnenvRoot() string {
	if root := os.Getenv("MVNENV_ROOT"); root != "" {
		return root
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".mvnenv")
}
