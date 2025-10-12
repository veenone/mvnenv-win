# Project Structure

## Directory Organization

```
mvnenv-win/
├── cmd/                          # Command-line executables (Go convention)
│   ├── mvnenv/                   # Main CLI application
│   │   └── main.go              # Entry point for mvnenv command
│   └── shim/                     # Shim executable
│       └── main.go              # Entry point for Maven shim
│
├── internal/                     # Private application packages (Go convention)
│   ├── config/                   # Configuration management
│   │   ├── config.go            # Configuration loading and validation
│   │   ├── repository.go        # Repository configuration structures
│   │   └── paths.go             # Path resolution and defaults
│   ├── download/                 # Download manager
│   │   ├── downloader.go        # HTTP download with progress tracking
│   │   ├── checksum.go          # SHA-256 verification
│   │   └── cache.go             # Download cache management
│   ├── environment/              # Environment variable handling
│   │   ├── path.go              # PATH manipulation (Windows-specific)
│   │   ├── maven_home.go        # MAVEN_HOME management
│   │   └── registry.go          # Windows registry operations
│   ├── nexus/                    # Nexus repository client
│   │   ├── client.go            # Nexus REST API client
│   │   ├── auth.go              # Authentication handlers
│   │   └── metadata.go          # Version metadata parsing
│   ├── shim/                     # Shim generation and management
│   │   ├── generator.go         # Shim executable generation
│   │   ├── resolver.go          # Version resolution logic
│   │   └── executor.go          # Maven command execution
│   └── version/                  # Version management
│       ├── manager.go           # Install/uninstall operations
│       ├── resolver.go          # Version resolution (global/local/shell)
│       ├── parser.go            # Version string parsing
│       └── compare.go           # Version comparison logic
│
├── pkg/                          # Public reusable packages (Go convention)
│   └── maven/                    # Maven-specific utilities
│       ├── version.go           # Maven version structure and parsing
│       ├── paths.go             # Maven installation paths
│       └── metadata.go          # Maven metadata structures
│
├── test/                         # Test suites and test fixtures
│   ├── fixtures/                # Test data and fixtures
│   │   ├── maven-metadata.xml  # Sample Nexus metadata
│   │   └── config.yaml         # Sample configuration files
│   ├── integration/             # Integration tests
│   │   ├── install_test.go     # Maven installation tests
│   │   └── switch_test.go      # Version switching tests
│   └── testutil/                # Test utilities and helpers
│       ├── mock_nexus.go       # Mock Nexus server
│       └── filesystem.go       # Temporary filesystem helpers
│
├── scripts/                      # Build and development scripts
│   ├── install-mvnenv-win.ps1  # PowerShell installation script
│   ├── build.ps1               # Build script for releases
│   └── test.ps1                # Test runner script
│
├── docs/                         # Documentation
│   ├── installation.md          # Installation guide
│   ├── usage.md                # Command usage reference
│   ├── configuration.md        # Configuration guide
│   ├── nexus-integration.md    # Nexus setup guide
│   └── troubleshooting.md      # Common issues and solutions
│
├── .github/                      # GitHub-specific files
│   ├── workflows/               # GitHub Actions CI/CD
│   │   ├── test.yml            # Test workflow
│   │   ├── release.yml         # Release workflow
│   │   └── lint.yml            # Linting workflow
│   └── ISSUE_TEMPLATE/          # Issue templates
│
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── Makefile                      # Build automation
├── README.md                     # Project overview
├── LICENSE                       # License file
└── .gitignore                   # Git ignore patterns
```

## Installation Directory Structure

The application manages Maven versions in the user's home directory:

```
%USERPROFILE%\.mvnenv/
├── bin/                          # mvnenv executables
│   └── mvnenv.exe               # Main executable
│
├── shims/                        # Shim executables (added to PATH)
│   ├── mvn.exe                  # Maven shim (Windows executable)
│   └── mvn.cmd                  # Maven shim (batch wrapper)
│
├── versions/                     # Installed Maven versions
│   ├── 3.6.3/                   # Maven 3.6.3 installation
│   │   ├── bin/
│   │   ├── lib/
│   │   └── conf/
│   ├── 3.8.6/                   # Maven 3.8.6 installation
│   │   └── ...
│   └── 3.9.4/                   # Maven 3.9.4 installation
│       └── ...
│
├── cache/                        # Downloaded distribution archives
│   ├── apache-maven-3.6.3-bin.zip
│   ├── apache-maven-3.8.6-bin.zip
│   └── apache-maven-3.9.4-bin.zip
│
├── config/                       # Configuration files
│   ├── config.yaml              # Main configuration
│   └── repositories.yaml        # Repository configuration
│
└── logs/                         # Operation logs
    └── mvnenv.log               # Application log file
```

## Naming Conventions

### Files

#### Go Source Files
- **Main Packages**: `main.go` (standard Go convention for package main)
- **Implementation Files**: `snake_case.go` (e.g., `maven_home.go`, `version_manager.go`)
- **Test Files**: `[filename]_test.go` (Go standard, e.g., `downloader_test.go`)
- **Windows-specific**: `[filename]_windows.go` (Go build tags, e.g., `registry_windows.go`)

#### Configuration Files
- **YAML Configuration**: `snake_case.yaml` (e.g., `config.yaml`, `repositories.yaml`)
- **Documentation**: `kebab-case.md` (e.g., `nexus-integration.md`, `installation-guide.md`)

#### Executables
- **Windows Executables**: `kebab-case.exe` (e.g., `mvnenv.exe`)
- **Batch Scripts**: `kebab-case.cmd` or `.bat` (e.g., `mvn.cmd`)
- **PowerShell Scripts**: `PascalCase.ps1` or `kebab-case.ps1` (e.g., `install-mvnenv-win.ps1`)

### Code

#### Go Naming Conventions (following Go standards)
- **Packages**: `lowercase` single-word when possible (e.g., `config`, `version`, `download`)
- **Exported Types**: `PascalCase` (e.g., `VersionManager`, `NexusClient`, `ConfigLoader`)
- **Unexported Types**: `camelCase` (e.g., `versionResolver`, `downloadCache`)
- **Interfaces**: `PascalCase` with `-er` suffix (e.g., `Downloader`, `VersionResolver`, `Installer`)
- **Functions/Methods**: `PascalCase` for exported, `camelCase` for unexported
  - Exported: `InstallVersion()`, `GetCurrentVersion()`, `ResolveVersion()`
  - Unexported: `parseMetadata()`, `validateChecksum()`, `extractArchive()`
- **Constants**: `PascalCase` for exported, `camelCase` for unexported
  - Exported: `DefaultTimeout`, `MaxRetries`, `ConfigFileName`
  - Unexported: `defaultConfigPath`, `shimExecutableName`
- **Variables**: `camelCase` for locals, `PascalCase` for exported package-level
  - Local: `currentVersion`, `downloadPath`, `configFile`
  - Exported: `DefaultConfig`, `SupportedVersions`

#### Special Conventions
- **Error Variables**: Start with `Err` (e.g., `ErrVersionNotFound`, `ErrInvalidConfig`)
- **Error Returns**: Functions return `error` as last return value (Go convention)
- **Context Parameters**: Context should be first parameter when used

## Import Patterns

### Import Order (following Go conventions and goimports)

```go
package example

import (
    // 1. Standard library imports (alphabetical)
    "context"
    "fmt"
    "os"
    "path/filepath"

    // 2. Third-party imports (alphabetical)
    "github.com/go-resty/resty/v2"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "gopkg.in/yaml.v3"

    // 3. Internal imports - project packages (alphabetical)
    "github.com/veenone/mvnenv-win/internal/config"
    "github.com/veenone/mvnenv-win/internal/download"
    "github.com/veenone/mvnenv-win/pkg/maven"
)
```

### Module Organization
- **Module Path**: `github.com/veenone/mvnenv-win`
- **Import Style**: Absolute imports from module root
- **Internal Packages**: Can only be imported by code in parent tree (Go enforcement)
- **Pkg Packages**: Can be imported by external projects (public API)

## Code Structure Patterns

### File Organization

#### Standard Go File Structure
```go
// 1. Package declaration and documentation
// Package version provides Maven version management functionality.
package version

// 2. Imports (grouped as described above)
import (
    "fmt"
    "os"
    // ...
)

// 3. Constants (exported first, then unexported)
const (
    DefaultVersion = "3.9.4"
    MaxVersions    = 100
)

const (
    versionFileName = ".maven-version"
    globalVersionKey = "global_version"
)

// 4. Package-level variables (if needed, minimize these)
var (
    SupportedVersions = []string{"3.6.3", "3.8.6", "3.9.4"}
)

// 5. Type definitions (exported first, then unexported)
type Manager struct {
    config     *config.Config
    downloader Downloader
    logger     Logger
}

// 6. Constructor functions (New* pattern)
func NewManager(cfg *config.Config, opts ...Option) (*Manager, error) {
    // Implementation
}

// 7. Public methods (grouped by receiver type)
func (m *Manager) InstallVersion(ctx context.Context, version string) error {
    // Implementation
}

func (m *Manager) UninstallVersion(version string) error {
    // Implementation
}

// 8. Private methods (grouped by receiver type)
func (m *Manager) resolveVersionPath(version string) string {
    // Implementation
}

// 9. Package-level helper functions (exported first, then unexported)
func ParseVersion(input string) (string, error) {
    // Implementation
}

func validateVersionString(input string) bool {
    // Implementation
}
```

### Function/Method Organization

```go
func (m *Manager) InstallVersion(ctx context.Context, version string) error {
    // 1. Input validation
    if version == "" {
        return ErrEmptyVersion
    }
    if !isValidVersion(version) {
        return ErrInvalidVersion
    }

    // 2. Precondition checks
    if m.IsInstalled(version) {
        return ErrAlreadyInstalled
    }

    // 3. Core logic
    downloadPath, err := m.downloadVersion(ctx, version)
    if err != nil {
        return fmt.Errorf("download failed: %w", err)
    }

    if err := m.extractVersion(downloadPath, version); err != nil {
        m.cleanup(downloadPath) // Cleanup on error
        return fmt.Errorf("extraction failed: %w", err)
    }

    // 4. Post-processing
    if err := m.verifyInstallation(version); err != nil {
        return fmt.Errorf("verification failed: %w", err)
    }

    // 5. Success return
    return nil
}
```

## Code Organization Principles

1. **Single Responsibility**: Each package has one clear purpose
   - `internal/config`: Configuration only
   - `internal/download`: Download operations only
   - `internal/version`: Version management only

2. **Modularity**: Packages are self-contained with clear interfaces
   - Dependencies injected via constructors
   - Interfaces defined in consuming packages (Go best practice)

3. **Testability**: Code structured for easy testing
   - Dependencies are interfaces
   - File system operations abstracted for testing
   - Network operations mockable

4. **Consistency**: Follow established Go conventions
   - Standard project layout (cmd, internal, pkg)
   - Go standard library patterns
   - Effective Go guidelines

## Module Boundaries

### Internal vs Public API

**Internal Packages** (`internal/`):
- Cannot be imported by external projects (Go enforced)
- Implementation details of mvnenv-win
- Can change without breaking compatibility
- Examples: `internal/config`, `internal/download`, `internal/nexus`

**Public Packages** (`pkg/`):
- Can be imported by external projects
- Stable API for reuse (e.g., plugins, extensions)
- Breaking changes require major version bump
- Examples: `pkg/maven` (Maven version utilities)

### Dependency Direction

```
┌──────────────────────────────────────┐
│  cmd/mvnenv (CLI Entry Point)        │
└──────────────┬───────────────────────┘
               │ depends on
               ▼
┌──────────────────────────────────────┐
│  internal/* (Business Logic)         │
│  - config, download, version, etc.   │
└──────────────┬───────────────────────┘
               │ depends on
               ▼
┌──────────────────────────────────────┐
│  pkg/maven (Public Utilities)        │
└──────────────────────────────────────┘

Rule: Dependencies flow downward only
- cmd depends on internal (but not vice versa)
- internal can use pkg (but pkg cannot use internal)
- internal packages can depend on each other (minimize)
```

### Platform-Specific Code

**Windows-Specific Implementation**:
- Use Go build tags: `//go:build windows`
- Filename convention: `*_windows.go`
- Examples:
  - `internal/environment/registry_windows.go` (Windows registry)
  - `internal/environment/path_windows.go` (Windows PATH manipulation)

**Future Cross-Platform Support**:
- Linux/macOS implementations: `*_unix.go` or `*_linux.go`, `*_darwin.go`
- Shared interfaces in non-suffixed files
- Example: `environment.go` defines interface, `environment_windows.go` implements

### Stable vs Experimental

**Stable Components** (v1.0.0):
- Core version management: `internal/version`
- Configuration: `internal/config`
- Shim system: `internal/shim`

**Future Experimental** (post-v1.0.0):
- IDE plugins: Would go in `internal/plugins` or separate repo
- Web dashboard: Would go in `internal/web` or separate repo
- Team collaboration features: TBD architecture

## Code Size Guidelines

### File Size
- **Target**: <500 lines per file (excluding generated code)
- **Maximum**: 1000 lines (consider splitting if exceeded)
- **Rationale**: Maintainability and ease of navigation

### Function/Method Size
- **Target**: <50 lines per function
- **Maximum**: 100 lines (consider refactoring if exceeded)
- **Extract helpers** for logic blocks that can be named meaningfully

### Function Complexity
- **Cyclomatic Complexity**: Aim for <10 per function
- **Nesting Depth**: Maximum 3-4 levels (use early returns)
- **Parameters**: Maximum 5 parameters (use config structs for more)

### Example of Refactoring Deep Nesting

**Before** (deep nesting):
```go
func processVersion(version string) error {
    if version != "" {
        if isValid(version) {
            if !isInstalled(version) {
                if err := download(version); err == nil {
                    if err := install(version); err == nil {
                        return nil
                    }
                }
            }
        }
    }
    return errors.New("failed")
}
```

**After** (early returns):
```go
func processVersion(version string) error {
    if version == "" {
        return ErrEmptyVersion
    }
    if !isValid(version) {
        return ErrInvalidVersion
    }
    if isInstalled(version) {
        return ErrAlreadyInstalled
    }

    if err := download(version); err != nil {
        return fmt.Errorf("download failed: %w", err)
    }

    if err := install(version); err != nil {
        return fmt.Errorf("install failed: %w", err)
    }

    return nil
}
```

## Documentation Standards

### Package Documentation
- Every package must have package-level documentation
- First sentence should be concise package description
- Example:
  ```go
  // Package version provides Maven version installation, management,
  // and resolution capabilities for mvnenv-win.
  package version
  ```

### Exported API Documentation
- All exported types, functions, constants must have godoc comments
- Start with the name of the item being documented
- Examples:
  ```go
  // Manager handles Maven version installation and management.
  type Manager struct { ... }

  // InstallVersion downloads and installs the specified Maven version.
  // It returns an error if the version is invalid or already installed.
  func (m *Manager) InstallVersion(ctx context.Context, version string) error { ... }

  // DefaultTimeout is the default timeout for download operations.
  const DefaultTimeout = 300 * time.Second
  ```

### Complex Logic Comments
- Use inline comments for non-obvious logic
- Explain "why" not "what" (code shows what)
- Example:
  ```go
  // Use a temporary directory for extraction to ensure atomic installation.
  // If extraction fails, we can simply delete the temp directory without
  // corrupting the versions directory.
  tempDir := filepath.Join(m.config.CacheDir, "tmp-"+version)
  ```

### README Files
- Root `README.md`: Project overview, installation, quick start
- Major modules can have `README.md` if needed
- Keep READMEs concise, link to detailed docs

### Examples and Tutorials
- `docs/` directory for comprehensive guides
- Code examples in `examples/` directory (if creating public API)
- Tutorials in `docs/tutorials/` for common workflows

## Error Handling Patterns

### Error Definition
```go
// Package-level errors
var (
    ErrVersionNotFound   = errors.New("version not found")
    ErrInvalidConfig     = errors.New("invalid configuration")
    ErrAlreadyInstalled  = errors.New("version already installed")
)
```

### Error Wrapping
```go
// Always wrap errors with context using %w
if err := downloadFile(url, dest); err != nil {
    return fmt.Errorf("failed to download from %s: %w", url, err)
}
```

### Error Checking
```go
// Check errors immediately after operation
file, err := os.Open(path)
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
defer file.Close()
```
