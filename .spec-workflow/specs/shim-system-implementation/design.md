# Design Document

## Architecture Overview

The Shim System provides transparent Maven command interception through lightweight proxy executables that resolve the active Maven version and forward commands to the appropriate installation. The architecture uses a single shim executable that determines which command was invoked based on its filename, resolves the version, and executes the target Maven command with full I/O pass-through.

### Core Design Principles

1. **Minimal Overhead**: Shim adds <50ms overhead through efficient version resolution and direct process spawning
2. **Transparent Operation**: Users interact with Maven commands normally; shims are invisible
3. **Single Binary**: One shim executable serves all Maven commands (mvn, mvnDebug, mvnyjp) via filename detection
4. **Zero Buffering**: Direct I/O forwarding without buffering ensures real-time output
5. **Fail Fast**: Clear error messages at resolution time prevent confusing Maven errors

## Component Architecture

```
cmd/shim/
└── main.go                # Shim executable entry point

internal/shim/
├── generator.go           # ShimGenerator creates shim files
├── executor.go            # ShimExecutor handles command execution
├── resolver.go            # Integrates with VersionResolver
└── path_windows.go        # Windows-specific PATH management
```

### Component Relationships

```
User types "mvn clean"
         ↓
    Shim Executable (mvn.exe)
         ↓
    Detect Command (from filename)
         ↓
    ShimExecutor.Execute()
         ↓
    VersionResolver.ResolveVersion() ←── from core-version-management
         ↓
    Construct Maven Path
         ↓
    Set MAVEN_HOME
         ↓
    os/exec.Command with I/O forwarding
         ↓
    Maven Process Execution
         ↓
    Exit with Maven's exit code
```

## Detailed Component Design

### 1. Shim Executable (cmd/shim/main.go)

The main entry point for the shim executable. This is the actual binary that gets invoked when users run `mvn`, `mvnDebug`, etc.

```go
// cmd/shim/main.go
package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/veenone/mvnenv-win/internal/shim"
    "github.com/veenone/mvnenv-win/internal/version"
)

func main() {
    // Detect command name from executable name
    command := detectCommand()

    // Create executor with version resolver
    resolver := version.NewVersionResolver()
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
// by examining the executable filename (mvn.exe -> "mvn")
func detectCommand() string {
    exePath, err := os.Executable()
    if err != nil {
        return "mvn" // default fallback
    }

    base := filepath.Base(exePath)
    // Remove .exe extension if present
    name := strings.TrimSuffix(base, ".exe")
    return name
}
```

**Key Design Decisions:**
- Single binary approach: reduces code duplication and simplifies updates
- Filename-based detection: allows same executable to serve multiple commands
- Minimal dependencies: only uses internal packages and standard library
- Direct exit: preserves Maven's exit code exactly

### 2. ShimGenerator

Generates shim executable files and manages the shims directory.

```go
// internal/shim/generator.go
package shim

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
)

// ShimGenerator creates and manages Maven command shims
type ShimGenerator struct {
    shimsDir     string // %USERPROFILE%\.mvnenv\shims
    shimBinary   string // Path to shim.exe
    versionsDir  string // %USERPROFILE%\.mvnenv\versions
}

// NewShimGenerator creates a shim generator
func NewShimGenerator(mvnenvRoot string) *ShimGenerator {
    return &ShimGenerator{
        shimsDir:    filepath.Join(mvnenvRoot, "shims"),
        shimBinary:  filepath.Join(mvnenvRoot, "bin", "shim.exe"),
        versionsDir: filepath.Join(mvnenvRoot, "versions"),
    }
}

// GenerateShims creates shim executables for all Maven commands
// Returns list of generated shim paths or error
func (g *ShimGenerator) GenerateShims() ([]string, error) {
    // Ensure shims directory exists
    if err := os.MkdirAll(g.shimsDir, 0755); err != nil {
        return nil, fmt.Errorf("create shims directory: %w", err)
    }

    // Core Maven commands to shim
    commands := []string{"mvn", "mvnDebug"}

    // Scan installed versions for additional commands (mvnyjp)
    additionalCmds, err := g.discoverAdditionalCommands()
    if err != nil {
        // Log warning but continue with core commands
        log.Warnf("Failed to discover additional commands: %v", err)
    } else {
        commands = append(commands, additionalCmds...)
    }

    var generatedPaths []string

    for _, cmd := range commands {
        // Generate .exe shim (primary for PATH resolution)
        exePath, err := g.generateShimFile(cmd, ".exe")
        if err != nil {
            return nil, fmt.Errorf("generate %s.exe: %w", cmd, err)
        }
        generatedPaths = append(generatedPaths, exePath)

        // Generate .cmd shim (compatibility for cmd.exe)
        cmdPath, err := g.generateBatchShim(cmd)
        if err != nil {
            return nil, fmt.Errorf("generate %s.cmd: %w", cmd, err)
        }
        generatedPaths = append(generatedPaths, cmdPath)
    }

    return generatedPaths, nil
}

// generateShimFile creates executable shim by copying shim.exe
func (g *ShimGenerator) generateShimFile(command string, ext string) (string, error) {
    destPath := filepath.Join(g.shimsDir, command+ext)

    // Copy shim.exe to destination with new name
    // This allows filename-based command detection
    if err := copyFile(g.shimBinary, destPath); err != nil {
        return "", fmt.Errorf("copy shim binary: %w", err)
    }

    // Verify executable
    if err := verifyExecutable(destPath); err != nil {
        return "", fmt.Errorf("verify executable: %w", err)
    }

    return destPath, nil
}

// generateBatchShim creates .cmd shim that calls .exe shim
func (g *ShimGenerator) generateBatchShim(command string) (string, error) {
    destPath := filepath.Join(g.shimsDir, command+".cmd")

    // Batch script that forwards to .exe shim
    script := fmt.Sprintf(`@echo off
"%~dp0%s.exe" %%*
exit /b %%ERRORLEVEL%%
`, command)

    if err := os.WriteFile(destPath, []byte(script), 0755); err != nil {
        return "", fmt.Errorf("write batch shim: %w", err)
    }

    return destPath, nil
}

// discoverAdditionalCommands scans installed versions for commands like mvnyjp
func (g *ShimGenerator) discoverAdditionalCommands() ([]string, error) {
    var additionalCmds []string

    entries, err := os.ReadDir(g.versionsDir)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil // No versions installed yet
        }
        return nil, err
    }

    cmdSet := make(map[string]bool)

    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }

        binDir := filepath.Join(g.versionsDir, entry.Name(), "bin")
        binEntries, err := os.ReadDir(binDir)
        if err != nil {
            continue
        }

        for _, binEntry := range binEntries {
            name := binEntry.Name()
            // Look for .cmd files that aren't mvn or mvnDebug
            if strings.HasSuffix(name, ".cmd") {
                cmd := strings.TrimSuffix(name, ".cmd")
                if cmd != "mvn" && cmd != "mvnDebug" && !cmdSet[cmd] {
                    cmdSet[cmd] = true
                    additionalCmds = append(additionalCmds, cmd)
                }
            }
        }
    }

    return additionalCmds, nil
}

// copyFile copies file from src to dst atomically
func copyFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    // Write to temp file first
    tmpDst := dst + ".tmp"
    destination, err := os.Create(tmpDst)
    if err != nil {
        return err
    }

    _, err = io.Copy(destination, source)
    destination.Close()
    if err != nil {
        os.Remove(tmpDst)
        return err
    }

    // Atomic rename
    if err := os.Rename(tmpDst, dst); err != nil {
        os.Remove(tmpDst)
        return err
    }

    return nil
}

// verifyExecutable checks if file is executable
func verifyExecutable(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        return err
    }

    if info.IsDir() {
        return fmt.Errorf("path is directory")
    }

    // On Windows, check .exe extension
    if !strings.HasSuffix(path, ".exe") && !strings.HasSuffix(path, ".cmd") {
        return fmt.Errorf("not an executable file")
    }

    return nil
}
```

**Key Design Decisions:**
- Copy shim.exe approach: allows single binary to serve multiple commands
- Atomic file operations: temp file + rename prevents partial shims
- .exe + .cmd generation: ensures compatibility across Windows shells
- Additional command discovery: automatically includes mvnyjp if available
- Error handling: continues on non-critical errors (e.g., discovery failure)

### 3. ShimExecutor

Handles command execution with version resolution and I/O forwarding.

```go
// internal/shim/executor.go
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

    "github.com/veenone/mvnenv-win/internal/version"
)

// ShimExecutor executes Maven commands with version resolution
type ShimExecutor struct {
    resolver *version.VersionResolver
    debug    bool // Enable debug output
}

// NewShimExecutor creates a shim executor
func NewShimExecutor(resolver *version.VersionResolver) *ShimExecutor {
    debug := os.Getenv("MVNENV_DEBUG") == "1"
    return &ShimExecutor{
        resolver: resolver,
        debug:    debug,
    }
}

// Execute resolves Maven version and executes command
// Returns Maven's exit code and any resolution errors
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
    // versionPath: %USERPROFILE%\.mvnenv\versions\3.9.4
    // Result: %USERPROFILE%\.mvnenv\versions\3.9.4\bin\mvn.cmd
    return filepath.Join(versionPath, "bin", command+".cmd")
}

// executeMaven spawns Maven process with I/O forwarding
func (e *ShimExecutor) executeMaven(mavenPath string, args []string, resolved *version.ResolvedVersion) (int, error) {
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
        cancel() // Cancel context, which kills Maven process
    }()

    // Start Maven process
    if err := cmd.Start(); err != nil {
        return 1, fmt.Errorf("failed to start Maven: %w", err)
    }

    // Wait for Maven to complete
    if err := cmd.Wait(); err != nil {
        // Extract exit code from error
        if exitErr, ok := err.(*exec.ExitError); ok {
            if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
                return status.ExitStatus(), nil
            }
        }
        // Unknown error, return 1
        return 1, err
    }

    return 0, nil
}

// formatResolutionError creates user-friendly error messages
func (e *ShimExecutor) formatResolutionError(err error) error {
    switch {
    case version.IsVersionNotInstalledError(err):
        v := version.ExtractVersionFromError(err)
        return fmt.Errorf("Maven version '%s' is set but not installed.\nInstall it with: mvnenv install %s", v, v)

    case version.IsNoVersionSetError(err):
        return fmt.Errorf("No Maven version is set.\nSet a global version with: mvnenv global <version>\nOr see available versions with: mvnenv install -l")

    default:
        return fmt.Errorf("Failed to resolve Maven version: %w", err)
    }
}

// logDebug outputs diagnostic information to stderr
func (e *ShimExecutor) logDebug(command string, args []string, resolved *version.ResolvedVersion, mavenPath string, resolutionTime time.Duration) {
    fmt.Fprintf(os.Stderr, "[mvnenv] Debug Information:\n")
    fmt.Fprintf(os.Stderr, "[mvnenv]   Command: %s\n", command)
    fmt.Fprintf(os.Stderr, "[mvnenv]   Arguments: %v\n", args)
    fmt.Fprintf(os.Stderr, "[mvnenv]   Resolved version: %s\n", resolved.Version)
    fmt.Fprintf(os.Stderr, "[mvnenv]   Source: %s\n", resolved.Source)
    fmt.Fprintf(os.Stderr, "[mvnenv]   Maven path: %s\n", mavenPath)
    fmt.Fprintf(os.Stderr, "[mvnenv]   MAVEN_HOME: %s\n", resolved.Path)
    fmt.Fprintf(os.Stderr, "[mvnenv]   Resolution time: %v\n", resolutionTime)
}
```

**Key Design Decisions:**
- Context-based cancellation: proper signal handling
- Direct I/O forwarding: cmd.Stdin/Stdout/Stderr = os.Stdin/Stdout/Stderr
- Exit code extraction: preserves Maven's exact exit code
- MAVEN_HOME injection: set only for Maven subprocess
- Debug mode: controlled via MVNENV_DEBUG environment variable
- Error formatting: converts internal errors to user-friendly messages

### 4. Version Resolution Integration

Integration with VersionResolver from core-version-management spec.

```go
// internal/shim/resolver.go
package shim

import (
    "github.com/veenone/mvnenv-win/internal/version"
)

// ResolveActiveVersion is a convenience wrapper around VersionResolver
// Used by ShimExecutor
func ResolveActiveVersion(mvnenvRoot string) (*version.ResolvedVersion, error) {
    resolver := version.NewVersionResolver(mvnenvRoot)
    return resolver.ResolveVersion()
}
```

**Integration Points:**
1. Shim calls `version.VersionResolver.ResolveVersion()`
2. VersionResolver follows shell > local > global hierarchy
3. Returns `ResolvedVersion` with version string, path, and source
4. Shim uses path to construct Maven binary location

### 5. PATH Management (Windows-specific)

```go
// internal/shim/path_windows.go
//go:build windows

package shim

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "golang.org/x/sys/windows/registry"
)

// PathManager manages Windows PATH environment variable
type PathManager struct {
    shimsDir string
}

// NewPathManager creates a PATH manager
func NewPathManager(mvnenvRoot string) *PathManager {
    return &PathManager{
        shimsDir: filepath.Join(mvnenvRoot, "shims"),
    }
}

// AddShimsDirToPath adds shims directory to user PATH
// Returns true if PATH was modified, false if already present
func (m *PathManager) AddShimsDirToPath() (bool, error) {
    // Open user environment key
    key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
    if err != nil {
        return false, fmt.Errorf("open registry key: %w", err)
    }
    defer key.Close()

    // Read current PATH
    currentPath, _, err := key.GetStringValue("Path")
    if err != nil && err != registry.ErrNotExist {
        return false, fmt.Errorf("read PATH: %w", err)
    }

    // Check if shims directory already in PATH (case-insensitive)
    if m.isInPath(currentPath, m.shimsDir) {
        return false, nil // Already present
    }

    // Prepend shims directory to PATH
    var newPath string
    if currentPath == "" {
        newPath = m.shimsDir
    } else {
        newPath = m.shimsDir + ";" + currentPath
    }

    // Write new PATH to registry
    if err := key.SetStringValue("Path", newPath); err != nil {
        return false, fmt.Errorf("write PATH: %w", err)
    }

    // Broadcast environment change to Windows
    m.broadcastEnvironmentChange()

    return true, nil
}

// isInPath checks if directory is in PATH (case-insensitive)
func (m *PathManager) isInPath(path string, dir string) bool {
    paths := strings.Split(path, ";")
    dirLower := strings.ToLower(dir)

    for _, p := range paths {
        if strings.ToLower(strings.TrimSpace(p)) == dirLower {
            return true
        }
    }

    return false
}

// broadcastEnvironmentChange notifies Windows of environment change
func (m *PathManager) broadcastEnvironmentChange() {
    // Use SendMessageTimeout to broadcast WM_SETTINGCHANGE
    // This makes the change visible to new processes without logout
    // Implementation uses syscall to SendMessageTimeoutW
}

// RemoveShimsDirFromPath removes shims directory from user PATH
func (m *PathManager) RemoveShimsDirFromPath() (bool, error) {
    key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
    if err != nil {
        return false, fmt.Errorf("open registry key: %w", err)
    }
    defer key.Close()

    currentPath, _, err := key.GetStringValue("Path")
    if err != nil {
        return false, fmt.Errorf("read PATH: %w", err)
    }

    // Remove shims directory from PATH
    newPath := m.removeFromPath(currentPath, m.shimsDir)

    if newPath == currentPath {
        return false, nil // Not present
    }

    if err := key.SetStringValue("Path", newPath); err != nil {
        return false, fmt.Errorf("write PATH: %w", err)
    }

    m.broadcastEnvironmentChange()

    return true, nil
}

// removeFromPath removes directory from PATH string
func (m *PathManager) removeFromPath(path string, dir string) string {
    paths := strings.Split(path, ";")
    dirLower := strings.ToLower(dir)

    var newPaths []string
    for _, p := range paths {
        trimmed := strings.TrimSpace(p)
        if trimmed != "" && strings.ToLower(trimmed) != dirLower {
            newPaths = append(newPaths, trimmed)
        }
    }

    return strings.Join(newPaths, ";")
}
```

**Key Design Decisions:**
- Registry-based PATH management: modifies HKCU\Environment
- Case-insensitive comparison: Windows paths are case-insensitive
- Prepend to PATH: ensures shims take priority
- Broadcast change: SendMessageTimeout makes change visible immediately
- User-level only: no system-wide PATH modification

## Error Handling

### Error Types and Messages

```go
// Common shim errors
var (
    ErrShimBinaryNotFound = errors.New("shim binary not found")
    ErrVersionNotResolved = errors.New("version resolution failed")
    ErrMavenBinaryNotFound = errors.New("Maven binary not found")
    ErrExecutionFailed = errors.New("Maven execution failed")
)

// ShimError wraps errors with shim context
type ShimError struct {
    Command string
    Step    string // "resolution", "execution", "generation"
    Cause   error
}

func (e *ShimError) Error() string {
    return fmt.Sprintf("mvnenv shim error in %s for command '%s': %v", e.Step, e.Command, e.Cause)
}
```

### Error Handling Strategy

1. **Resolution Errors**: Clear message with suggested action (install version, set global)
2. **Execution Errors**: Preserve Maven's error output, add shim context only if Maven didn't start
3. **Generation Errors**: Detailed error with which step failed
4. **PATH Errors**: Provide manual instructions if automatic modification fails

## Performance Optimization

### Optimization Techniques

1. **Fast Version Resolution**:
   - Read only necessary files (.maven-version, config.yaml)
   - No network operations
   - No expensive parsing

2. **Direct Process Execution**:
   - Use exec.CommandContext for minimal overhead
   - No intermediate shells or scripts
   - Direct I/O forwarding (no buffering)

3. **Lazy Loading**:
   - Don't load configuration unless needed
   - Don't discover additional commands during execution

4. **Binary Copying for Shim Generation**:
   - Copying shim.exe is faster than generating scripts
   - Enables filename-based command detection

### Performance Benchmarks

Target performance metrics:
- Version resolution: <25ms
- Process spawn: <25ms
- Total overhead: <50ms

Measurement approach:
```go
func BenchmarkShimOverhead(b *testing.B) {
    for i := 0; i < b.N; i++ {
        start := time.Now()
        // Resolve version
        // Spawn process (no-op Maven for test)
        overhead := time.Since(start)
        if overhead > 50*time.Millisecond {
            b.Fatalf("Overhead too high: %v", overhead)
        }
    }
}
```

## Integration Points

### Integration with core-version-management

```go
// ShimExecutor uses VersionResolver
type ShimExecutor struct {
    resolver *version.VersionResolver
}

func (e *ShimExecutor) Execute(command string, args []string) (int, error) {
    // Call VersionResolver from core-version-management
    resolved, err := e.resolver.ResolveVersion()
    if err != nil {
        return 1, e.formatResolutionError(err)
    }

    // Use resolved.Path to construct Maven binary path
    mavenPath := filepath.Join(resolved.Path, "bin", command+".cmd")

    // Execute Maven with MAVEN_HOME set to resolved.Path
    return e.executeMaven(mavenPath, args, resolved)
}
```

### Integration with cli-commands

```go
// cmd/mvnenv/cmd/rehash.go
func runRehash(cmd *cobra.Command, args []string) error {
    generator := shim.NewShimGenerator(mvnenvRoot)

    generatedPaths, err := generator.GenerateShims()
    if err != nil {
        return fmt.Errorf("failed to regenerate shims: %w", err)
    }

    fmt.Printf("Shims regenerated successfully (%d files)\n", len(generatedPaths))
    return nil
}
```

### Automatic Regeneration Hooks

```go
// In internal/version/installer.go
func (i *VersionInstaller) InstallVersion(version string) error {
    // ... installation logic ...

    // Automatically regenerate shims after installation
    if err := i.regenerateShims(); err != nil {
        log.Warnf("Failed to regenerate shims: %v", err)
        log.Info("Run 'mvnenv rehash' manually to update shims")
    }

    return nil
}
```

## Testing Strategy

### Unit Tests

```go
// internal/shim/generator_test.go
func TestShimGenerator_GenerateShims(t *testing.T)
func TestShimGenerator_DiscoverAdditionalCommands(t *testing.T)
func TestShimGenerator_AtomicGeneration(t *testing.T)

// internal/shim/executor_test.go
func TestShimExecutor_Execute_Success(t *testing.T)
func TestShimExecutor_Execute_VersionNotFound(t *testing.T)
func TestShimExecutor_Execute_ExitCodePreserved(t *testing.T)
func TestShimExecutor_Execute_SignalHandling(t *testing.T)

// internal/shim/path_windows_test.go
func TestPathManager_AddShimsDirToPath(t *testing.T)
func TestPathManager_IsInPath_CaseInsensitive(t *testing.T)
func TestPathManager_RemoveShimsDirFromPath(t *testing.T)
```

### Integration Tests

```go
// test/integration/shim_integration_test.go

// Test full workflow: generate shim, invoke, verify Maven executed
func TestIntegration_ShimExecution(t *testing.T)

// Test shim with different version sources (shell, local, global)
func TestIntegration_ShimVersionResolution(t *testing.T)

// Test I/O forwarding (stdin, stdout, stderr)
func TestIntegration_ShimIOForwarding(t *testing.T)

// Test exit code preservation
func TestIntegration_ShimExitCode(t *testing.T)

// Test concurrent shim invocations
func TestIntegration_ConcurrentShims(t *testing.T)
```

### Performance Tests

```go
// test/benchmarks/shim_bench_test.go

// Measure shim overhead
func BenchmarkShimOverhead(b *testing.B)

// Measure version resolution time
func BenchmarkVersionResolution(b *testing.B)

// Measure process spawn time
func BenchmarkProcessSpawn(b *testing.B)
```

## Security Considerations

### Path Injection Prevention

```go
// Validate resolved Maven path before execution
func (e *ShimExecutor) validateMavenPath(path string, versionPath string) error {
    // Ensure path is within version directory
    absPath, err := filepath.Abs(path)
    if err != nil {
        return err
    }

    absVersionPath, err := filepath.Abs(versionPath)
    if err != nil {
        return err
    }

    // Check path is within version directory
    if !strings.HasPrefix(absPath, absVersionPath) {
        return fmt.Errorf("Maven path outside version directory: %s", absPath)
    }

    return nil
}
```

### Argument Safety

- Arguments passed unchanged to Maven (no interpretation)
- No shell expansion or command substitution
- Direct process execution (not via shell)

### Environment Isolation

- MAVEN_HOME only set for Maven subprocess
- Parent shell environment unaffected
- No persistent environment changes

## Configuration

### Debug Mode

Enable via environment variable:
```bash
set MVNENV_DEBUG=1
mvn clean
```

Debug output (to stderr):
```
[mvnenv] Debug Information:
[mvnenv]   Command: mvn
[mvnenv]   Arguments: [clean]
[mvnenv]   Resolved version: 3.9.4
[mvnenv]   Source: local (.maven-version)
[mvnenv]   Maven path: C:\Users\user\.mvnenv\versions\3.9.4\bin\mvn.cmd
[mvnenv]   MAVEN_HOME: C:\Users\user\.mvnenv\versions\3.9.4
[mvnenv]   Resolution time: 12ms
[mvnenv] Total execution time: 2.345s
```

## Future Enhancements

### Phase 1 (v1.1.0)
- Maven wrapper (mvnw) detection and preference
- Custom shim hooks for pre/post execution
- Performance telemetry (opt-in)

### Phase 2 (v1.2.0)
- IDE integration hints (environment file generation)
- Shim auto-update on mvnenv upgrade
- Plugin system for custom shim behavior

### Phase 3 (v2.0.0)
- Cross-platform shim support (Linux, macOS)
- Remote execution support (SSH forwarding)
- Container-aware shim behavior
