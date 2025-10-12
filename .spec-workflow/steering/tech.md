# Technology Stack

## Project Type
**Command-Line Tool** - A native Windows console application for managing multiple Apache Maven installations with version switching capabilities, shim-based command interception, and private repository integration.

## Core Technologies

### Primary Language(s)
- **Language**: Go (Golang) 1.21+
- **Runtime/Compiler**: Go compiler with CGO disabled for pure Go binaries
- **Language-specific tools**:
  - Go modules for dependency management
  - Go build toolchain for compilation and cross-compilation
  - Go test framework for unit and integration testing

**Rationale**: Go was chosen for its excellent cross-platform compilation, single-binary distribution, fast execution speed, and strong standard library support for file operations, HTTP clients, and command-line parsing.

### Key Dependencies/Libraries

#### CLI Framework & Configuration
- **github.com/spf13/cobra**: Command-line interface framework providing structured command parsing, help generation, and subcommand organization
- **github.com/spf13/viper**: Configuration management supporting YAML/TOML files, environment variables, and configuration merging

#### HTTP & Network
- **github.com/go-resty/resty/v2**: HTTP client for Nexus API communication with retry logic, authentication, and response handling
- **net/http**: Go standard library for HTTP operations and download management

#### Data Formats & Serialization
- **gopkg.in/yaml.v3**: YAML parsing for configuration files (config.yaml, repositories.yaml)
- **encoding/json**: Standard library JSON handling for API responses and metadata

#### File Operations & Compression
- **archive/zip**: Standard library for extracting Maven distribution archives
- **path/filepath**: Cross-platform file path manipulation
- **os**: File system operations, environment variables, and process management

#### Security & Cryptography
- **crypto/sha256**: Checksum verification for downloaded Maven distributions
- **crypto/tls**: TLS/SSL certificate handling for secure Nexus connections
- **Windows Credential Manager integration**: Secure credential storage (via syscall or third-party library)

### Application Architecture

**Modular Command-Line Architecture** with the following design patterns:

- **Command Pattern**: Each command (install, uninstall, global, local) is implemented as a discrete command module using Cobra
- **Shim Architecture**: Lightweight proxy executables intercept Maven commands and delegate to the appropriate version
- **Repository Pattern**: Abstraction layer for Maven distribution sources (Apache, Nexus) with pluggable implementations
- **Version Resolution Strategy**: Hierarchical version selection (shell > local > global) with explicit resolution logging
- **Configuration Singleton**: Centralized configuration management through Viper with lazy loading

#### Component Structure
```
mvnenv-win/
├── cmd/
│   ├── mvnenv/          # Main CLI entry point with command registration
│   └── shim/            # Shim executable for Maven command interception
├── internal/            # Private application code
│   ├── config/          # Configuration loading, validation, and management
│   ├── download/        # HTTP download with progress tracking and resume
│   ├── environment/     # PATH and MAVEN_HOME manipulation
│   ├── nexus/           # Nexus repository client with authentication
│   ├── shim/            # Shim generation and management
│   └── version/         # Version parsing, comparison, and resolution
├── pkg/                 # Public reusable packages
│   └── maven/           # Maven-specific utilities (version parsing, paths)
└── test/                # Test suites and fixtures
```

### Data Storage

#### Primary Storage
- **File system**: All data stored in `%USERPROFILE%\.mvnenv\` directory structure
- **Configuration files**: YAML format for human readability and editability
- **Version files**: Plain text `.maven-version` files for project-specific versions

#### Storage Layout
```
%USERPROFILE%\.mvnenv\
├── versions/            # Extracted Maven installations (one directory per version)
├── cache/               # Downloaded distribution archives
├── config/              # config.yaml, repositories.yaml
├── logs/                # Operation logs
├── bin/                 # mvnenv executable
└── shims/               # Generated shim executables (mvn.exe, mvn.cmd)
```

#### Caching Strategy
- **Distribution cache**: Downloaded Maven archives cached indefinitely
- **Repository metadata**: Available versions list cached with TTL (configurable, default 24h)
- **Version resolution cache**: In-memory cache of resolved versions per directory (process lifetime)

#### Data Formats
- **YAML**: Configuration files for human readability
- **Plain text**: Version files (.maven-version)
- **ZIP**: Maven distribution format
- **JSON**: Nexus API responses and internal metadata

### External Integrations

#### APIs
- **Apache Maven Archive**: HTTP-based download from archive.apache.org
- **Sonatype Nexus Repository Manager**: REST API for version discovery and distribution download
  - Metadata queries: `/repository/{name}/org/apache/maven/apache-maven/maven-metadata.xml`
  - Artifact downloads: `/repository/{name}/org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip`

#### Protocols
- **HTTP/HTTPS**: Primary protocol for Maven distribution downloads
- **REST**: Nexus Repository Manager REST API
- **File system**: Local version management and configuration

#### Authentication
- **Basic Authentication**: Username/password for Nexus repositories
- **Token Authentication**: Bearer tokens for modern Nexus deployments
- **Environment Variables**: Credential injection via `${VAR}` syntax in configuration
- **Windows Credential Manager**: Secure credential storage (optional, recommended)

### Monitoring & Dashboard Technologies
Not applicable - mvnenv-win is a command-line tool without a persistent dashboard. Monitoring is achieved through:
- **CLI output**: Rich terminal output with color coding and formatting
- **Log files**: Structured logging to `%USERPROFILE%\.mvnenv\logs\`
- **Status commands**: `mvnenv version`, `mvnenv versions`, `mvnenv which` for state inspection

## Development Environment

### Build & Development Tools
- **Build System**: Go build system with Makefile for convenience commands
- **Package Management**: Go modules (go.mod, go.sum)
- **Development workflow**:
  - `go run` for rapid iteration during development
  - `go build` for local binary compilation
  - `make` for common development tasks (build, test, lint, clean)
- **Cross-compilation**: Go's built-in cross-compilation for Windows 32/64-bit

### Code Quality Tools
- **Static Analysis**:
  - `go vet` for code correctness
  - `staticcheck` for advanced static analysis
  - `golangci-lint` as meta-linter aggregator
- **Formatting**:
  - `gofmt` for standard Go formatting
  - `goimports` for import organization
- **Testing Framework**:
  - `go test` for unit tests
  - `testify/assert` for test assertions
  - `testify/mock` for mocking dependencies
  - Integration tests using temporary file systems
- **Documentation**:
  - `godoc` for package documentation
  - Markdown for user-facing documentation

### Version Control & Collaboration
- **VCS**: Git
- **Branching Strategy**: GitHub Flow (feature branches, main branch, pull requests)
- **Code Review Process**:
  - Pull request-based reviews on GitHub
  - Required CI checks (tests, linting) before merge
  - At least one approval required for merge

## Deployment & Distribution

### Target Platform(s)
- **Primary**: Windows 10+ (64-bit)
- **Architecture**: AMD64 (x86-64)
- **Future**: Windows 32-bit, Windows ARM64 (if demand exists)

### Distribution Method
1. **GitHub Releases**: Primary distribution channel with pre-built binaries
2. **PowerShell Install Script**: Automated installation via `install-mvnenv-win.ps1`
3. **Chocolatey**: Windows package manager for enterprise deployments
4. **Scoop**: Developer-focused Windows package manager
5. **Go Install**: Direct installation via `go install` for Go developers

### Installation Requirements
- **Operating System**: Windows 10 version 1809+ or Windows 11
- **Disk Space**: ~100MB for mvnenv-win + ~10MB per Maven version
- **Network**: Internet access for downloading Maven distributions from:
  - Apache Maven official archives (archive.apache.org)
  - Private Nexus repositories (if configured)
  - Custom repository mirrors
- **Permissions**: User-level permissions (no administrator required)
- **PATH modification**: User PATH updated during installation
- **Repository Access**:
  - For Nexus repositories: Valid credentials (username/password or token)
  - For Apache archives: No authentication required

### Update Mechanism
- **Self-update command**: `mvnenv self-update` checks GitHub releases and updates in-place
- **Package manager updates**: Chocolatey/Scoop update mechanisms
- **Manual update**: Download and replace executable

## Technical Requirements & Constraints

### Performance Requirements
- **Version switching**: <100ms for `mvnenv global/local/shell` commands
- **Shim execution overhead**: <50ms from command invocation to Maven execution
- **Maven installation**: Limited by network speed and disk I/O (~30-60 seconds for typical version)
- **Memory footprint**: <50MB resident memory during normal operation
- **Startup time**: <200ms for command parsing and initialization

### Compatibility Requirements

#### Platform Support
- **Operating Systems**: Windows 10 (1809+), Windows 11, Windows Server 2019+
- **Architectures**: x86-64 (AMD64) required
- **Shell Support**:
  - PowerShell 5.1+
  - PowerShell Core 7+
  - Command Prompt (cmd.exe)
  - Windows Terminal
  - ConEmu, cmder, and other terminal emulators

#### Dependency Versions
- **Go**: 1.21 minimum for development and compilation
- **Maven**: Supports Maven 3.0.0+ versions (primary focus on 3.6+ and 3.8+)
- **Nexus**: Compatible with Nexus Repository Manager 2.x and 3.x

#### Standards Compliance
- **Windows PATH**: Follows Windows path resolution conventions
- **Exit codes**: Standard Unix-style exit codes (0 = success, non-zero = error)
- **Semantic Versioning**: Maven version parsing and comparison follows SemVer principles

### Security & Compliance

#### Security Requirements
- **Checksum Verification**: SHA-256 checksum validation for all downloaded distributions
- **TLS/SSL Verification**: Certificate validation for HTTPS connections (no insecure modes)
- **Credential Storage**: Windows Credential Manager for secure password storage
- **Configuration Security**: File permissions restricting config access to user account
- **No Elevation**: Explicitly designed to run without administrator privileges

#### Compliance Standards
- **No PII Collection**: No personally identifiable information collected or transmitted
- **Opt-in Telemetry**: If telemetry added in future, must be opt-in with clear disclosure
- **Offline Mode**: Full functionality without internet access (except downloads)

#### Threat Model
- **Man-in-the-Middle**: Mitigated by TLS certificate validation
- **Malicious Distributions**: Mitigated by checksum verification
- **Credential Theft**: Mitigated by Windows Credential Manager usage
- **PATH Hijacking**: Mitigated by shim directory ordering and integrity checks

### Scalability & Reliability

#### Expected Load
- **Concurrent Versions**: 10-50 Maven versions per installation (realistic: 3-5)
- **Repository Count**: 1-10 configured Nexus repositories
- **Concurrent Operations**: Single-user, single-process model (no concurrency requirements)

#### Availability Requirements
- **Offline Operation**: Core functionality (version switching) must work offline
- **Graceful Degradation**: Repository unavailability should not break local operations
- **Atomic Operations**: Version installations are transactional (rollback on failure)

#### Growth Projections
- **Version Storage**: Linear growth based on number of installed versions (~10MB per version)
- **Cache Growth**: Bounded by user-initiated installations (no unbounded cache growth)

## Technical Decisions & Rationale

### Decision Log

#### 1. Go Language Selection
**Decision**: Use Go as the primary implementation language instead of Python, Batch, or PowerShell.

**Rationale**:
- Single-binary distribution eliminates runtime dependencies (no Python interpreter required)
- Fast execution speed meets performance requirements (<100ms version switching)
- Excellent cross-compilation support for future multi-platform expansion
- Strong standard library for file operations, HTTP, and ZIP handling
- Better maintainability than batch scripts or shell scripts

**Alternatives Considered**:
- Python: Requires Python runtime, slower execution, dependency management complexity
- PowerShell: Windows-only, slower startup, version compatibility issues
- Batch scripts: Limited functionality, poor error handling, difficult to maintain

#### 2. Shim-Based Architecture
**Decision**: Use shim executables to intercept Maven commands rather than PATH manipulation alone.

**Rationale**:
- Transparent to users (no need to remember to run special commands)
- Allows per-directory version resolution automatically
- Enables shell-specific overrides without modifying global state
- Provides consistent behavior across different shell environments

**Alternatives Considered**:
- Shell aliases: Not portable across shells, requires shell-specific configuration
- Wrapper scripts: Similar to shims but less transparent
- Direct PATH manipulation: Requires manual switching, error-prone

#### 3. YAML Configuration Format
**Decision**: Use YAML for configuration files instead of JSON, TOML, or INI.

**Rationale**:
- Human-readable and easily editable by hand
- Support for comments (important for configuration documentation)
- Widely adopted in DevOps and configuration management tools
- Better multi-line string support than JSON

**Alternatives Considered**:
- JSON: Less human-friendly, no comment support
- TOML: Less familiar to average users, more verbose for nested structures
- INI: Too simplistic for hierarchical configuration needs

#### 4. User-Level Installation
**Decision**: Install in `%USERPROFILE%` without requiring administrator elevation.

**Rationale**:
- Easier enterprise adoption (no IT approval for admin rights)
- Safer security model (limited blast radius)
- Aligns with modern development tools (pyenv, rbenv, nvm)
- Per-user isolation prevents conflicts in multi-user systems

**Trade-offs Accepted**:
- Cannot be installed system-wide easily
- Each user needs separate installation

#### 5. Windows Credential Manager for Secrets
**Decision**: Use Windows Credential Manager for storing repository credentials.

**Rationale**:
- Native Windows security infrastructure
- OS-level encryption and access control
- No custom encryption key management needed
- Integrates with Windows security policies

**Trade-offs Accepted**:
- Windows-specific (not portable to Linux/macOS)
- Requires additional syscall complexity

## Known Limitations

### 1. Windows-Only Support
**Limitation**: Currently only supports Windows platforms.

**Impact**: Cannot be used on Linux or macOS development environments.

**Future Solution**: Planned cross-platform support in v2.0.0 with unified codebase and platform-specific adapters.

### 2. Single-User Model
**Limitation**: Designed for single-user, local installation without multi-user or network share support.

**Impact**: Each user on a system must install mvnenv-win separately.

**Rationale**: Simplifies implementation and security model. Multi-user support adds complexity around permissions, shared state, and version conflicts.

### 3. No IDE Plugin Integration
**Limitation**: No direct integration with IDEs like IntelliJ IDEA or VS Code in v1.0.

**Impact**: IDE Maven settings must be manually updated to point to mvnenv-managed Maven installations.

**Future Solution**: IDE plugins planned for v1.2.0 to automatically detect and configure mvnenv versions.

### 4. Manual Nexus Repository Configuration
**Limitation**: Nexus repositories must be manually configured via commands or config file editing.

**Impact**: Requires manual steps for enterprise setup rather than auto-discovery.

**Future Solution**: Consider auto-discovery via Nexus API or import from Maven settings.xml in future versions.

### 5. No Maven Wrapper Integration
**Limitation**: Does not automatically detect or integrate with Maven Wrapper (mvnw) in v1.0.

**Impact**: Projects using mvnw may have version conflicts or unexpected behavior.

**Future Solution**: Planned for v2.0.0 with intelligent detection and preference for mvnw when present.
