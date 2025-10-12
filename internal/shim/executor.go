package shim

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	versionpkg "github.com/veenone/mvnenv-win/internal/version"
)

// ShimExecutor executes Maven commands with version resolution
type ShimExecutor struct {
	resolver *versionpkg.VersionResolver
	debug    bool
}

// NewShimExecutor creates a shim executor
func NewShimExecutor(resolver *versionpkg.VersionResolver) *ShimExecutor {
	debug := os.Getenv("MVNENV_DEBUG") == "1"
	return &ShimExecutor{
		resolver: resolver,
		debug:    debug,
	}
}

// Execute resolves Maven version and executes command
func (e *ShimExecutor) Execute(command string, args []string) (int, error) {
	startTime := time.Now()

	// Resolve active Maven version
	resolved, err := e.resolver.ResolveVersion()
	if err != nil {
		return 1, e.formatResolutionError(err)
	}

	resolutionTime := time.Since(startTime)

	// Construct path to Maven command
	mavenPath := e.constructMavenPath(resolved.Path, command)

	// Verify Maven binary exists
	if _, err := os.Stat(mavenPath); err != nil {
		return 1, fmt.Errorf("Maven binary not found at %s\nVersion %s may be corrupted. Try reinstalling with: mvnenv install %s",
			mavenPath, resolved.Version, resolved.Version)
	}

	if e.debug {
		e.logDebug(command, args, resolved, mavenPath, resolutionTime)
	}

	// Execute Maven with I/O forwarding
	exitCode, err := e.executeMaven(mavenPath, args, resolved)

	if e.debug {
		executionTime := time.Since(startTime)
		fmt.Fprintf(os.Stderr, "[mvnenv] Total execution time: %v\n", executionTime)
	}

	return exitCode, err
}

// constructMavenPath builds path to Maven command binary
func (e *ShimExecutor) constructMavenPath(versionPath string, command string) string {
	return filepath.Join(versionPath, "bin", command+".cmd")
}

// executeMaven spawns Maven process with I/O forwarding
func (e *ShimExecutor) executeMaven(mavenPath string, args []string, resolved *versionpkg.ResolvedVersion) (int, error) {
	// Create command with context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, mavenPath, args...)

	// Set MAVEN_HOME environment variable
	cmd.Env = append(os.Environ(), fmt.Sprintf("MAVEN_HOME=%s", resolved.Path))

	// Forward stdin/stdout/stderr (no buffering)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set working directory to current directory
	cmd.Dir, _ = os.Getwd()

	// Handle signals (Ctrl+C, etc.)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	// Start Maven process
	if err := cmd.Start(); err != nil {
		return 1, fmt.Errorf("failed to start Maven: %w", err)
	}

	// Wait for Maven to complete
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), nil
			}
		}
		return 1, err
	}

	return 0, nil
}

// formatResolutionError creates user-friendly error messages
func (e *ShimExecutor) formatResolutionError(err error) error {
	switch {
	case versionpkg.IsVersionNotInstalledError(err):
		ver := versionpkg.ExtractVersionFromError(err)
		return fmt.Errorf("Maven version '%s' is set but not installed.\nInstall it with: mvnenv install %s", ver, ver)

	case versionpkg.IsNoVersionSetError(err):
		return fmt.Errorf("No Maven version is set.\nSet a global version with: mvnenv global <version>\nOr see available versions with: mvnenv install -l")

	default:
		return fmt.Errorf("Failed to resolve Maven version: %w", err)
	}
}

// logDebug outputs diagnostic information to stderr
func (e *ShimExecutor) logDebug(command string, args []string, resolved *versionpkg.ResolvedVersion, mavenPath string, resolutionTime time.Duration) {
	fmt.Fprintf(os.Stderr, "[mvnenv] Debug Information:\n")
	fmt.Fprintf(os.Stderr, "[mvnenv]   Command: %s\n", command)
	fmt.Fprintf(os.Stderr, "[mvnenv]   Arguments: %v\n", args)
	fmt.Fprintf(os.Stderr, "[mvnenv]   Resolved version: %s\n", resolved.Version)
	fmt.Fprintf(os.Stderr, "[mvnenv]   Source: %s\n", resolved.Source)
	fmt.Fprintf(os.Stderr, "[mvnenv]   Maven path: %s\n", mavenPath)
	fmt.Fprintf(os.Stderr, "[mvnenv]   MAVEN_HOME: %s\n", resolved.Path)
	fmt.Fprintf(os.Stderr, "[mvnenv]   Resolution time: %v\n", resolutionTime)
}
