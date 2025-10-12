# Requirements Document

## Introduction

The Core Version Management feature implements the business logic for managing Maven installations on Windows. This includes downloading and installing Maven distributions, uninstalling versions, listing available and installed versions, and resolving which version should be active based on global, local, and shell settings. This feature provides the core functionality that the CLI commands (from cli-commands spec) will invoke.

This spec focuses purely on version management operations—downloading, installing, uninstalling, listing, and version resolution. Repository management (Nexus integration) is handled in the nexus-repository-integration spec, and shim operations are in the shim-system-implementation spec.

## Alignment with Product Vision

**Product Principle: Fast Performance**
The requirement that "version switching must be imperceptibly fast (<100ms)" directly informs the design of version resolution, which must use efficient file lookups and in-memory caching rather than expensive operations.

**Product Principle: Fail Safe, Not Fast**
Atomic installation operations with rollback capabilities ensure that partial or failed installations don't corrupt the system state, aligning with the principle of prioritizing reliability over speed when conflicts arise.

**Business Objective: Build Reproducibility**
Version resolution logic (global/local/shell hierarchy) ensures consistent, reproducible build environments across development, CI/CD, and production by providing deterministic version selection.

**Key Feature: Multi-Version Management**
This spec directly implements the core product feature of installing, maintaining, and switching between multiple Maven versions simultaneously without conflicts.

## Requirements

### Requirement 1: Maven Version Installation

**User Story:** As a developer, I want to install specific Maven versions from configured repositories, so that I can use different Maven versions for different projects.

#### Acceptance Criteria

1. WHEN user requests install of version X THEN system SHALL download Maven distribution from configured repository
2. WHEN downloading distribution THEN system SHALL verify file integrity using SHA-256 checksum
3. WHEN download completes THEN system SHALL extract archive to `%USERPROFILE%\.mvnenv\versions\{version}\`
4. WHEN extraction completes THEN system SHALL verify Maven binary exists at expected path
5. IF version already installed THEN system SHALL return error "version 'X' already installed"
6. IF download fails THEN system SHALL clean up partial files and return error with network details
7. IF checksum verification fails THEN system SHALL delete downloaded file and return error "checksum mismatch"
8. IF extraction fails THEN system SHALL remove version directory and return error with details
9. WHEN installation succeeds THEN system SHALL cache downloaded archive in `%USERPROFILE%\.mvnenv\cache\`
10. WHEN installing multiple versions THEN system SHALL install each sequentially, continuing on individual failures

### Requirement 2: Maven Version Uninstallation

**User Story:** As a developer, I want to uninstall Maven versions I no longer need, so that I can free up disk space and keep my system organized.

#### Acceptance Criteria

1. WHEN user requests uninstall of version X THEN system SHALL verify version is installed
2. IF version not installed THEN system SHALL return error "version 'X' not installed"
3. WHEN uninstalling THEN system SHALL check if version is currently active (global/local/shell)
4. IF version is currently active THEN system SHALL warn user and require confirmation
5. WHEN confirmed THEN system SHALL remove version directory `%USERPROFILE%\.mvnenv\versions\{version}\`
6. WHEN removing directory THEN system SHALL handle locked files gracefully with clear error
7. IF version is global version THEN system SHALL clear global version setting after removal
8. WHEN uninstall completes THEN system SHALL keep cached archive (for faster reinstall)
9. WHEN uninstalling multiple versions THEN system SHALL process each sequentially, continuing on failures

### Requirement 3: List Installed Versions

**User Story:** As a developer, I want to see which Maven versions are installed, so that I know what versions are available to use.

#### Acceptance Criteria

1. WHEN user requests installed versions list THEN system SHALL scan `%USERPROFILE%\.mvnenv\versions\` directory
2. WHEN scanning THEN system SHALL identify valid Maven installations by checking for `bin\mvn.cmd`
3. WHEN listing THEN system SHALL return versions sorted by semantic version (newest first)
4. WHEN displaying list THEN system SHALL indicate current active version with marker
5. IF no versions installed THEN system SHALL return empty list (not an error)
6. WHEN version directory exists but is invalid THEN system SHALL exclude from list and log warning

### Requirement 4: List Available Versions

**User Story:** As a developer, I want to see which Maven versions are available to install, so that I can choose the right version for my needs.

#### Acceptance Criteria

1. WHEN user requests available versions THEN system SHALL query all configured repositories
2. WHEN querying THEN system SHALL combine results from Apache Maven and Nexus repositories
3. WHEN listing THEN system SHALL return unique versions sorted by semantic version (newest first)
4. WHEN repository unavailable THEN system SHALL continue with other repositories and log warning
5. IF all repositories unavailable THEN system SHALL return error "no repositories available"
6. WHEN cache exists and is fresh THEN system SHALL use cached version list without querying
7. WHEN cache stale (>24 hours default) THEN system SHALL refresh from repositories
8. WHEN filtering by prefix THEN system SHALL return only versions matching prefix (e.g., "3.8")

### Requirement 5: Version Resolution (Global/Local/Shell)

**User Story:** As a developer, I want Maven version to be automatically selected based on my project and environment, so that I don't have to manually switch versions constantly.

#### Acceptance Criteria

1. WHEN resolving version THEN system SHALL check in order: shell > local > global
2. WHEN checking shell THEN system SHALL read environment variable `MVNENV_MAVEN_VERSION`
3. WHEN checking local THEN system SHALL search for `.maven-version` file in current and parent directories
4. WHEN checking global THEN system SHALL read from `%USERPROFILE%\.mvnenv\config\config.yaml`
5. IF shell version set THEN system SHALL return that version without checking local/global
6. IF local version file found THEN system SHALL return that version without checking global
7. IF only global version set THEN system SHALL return global version
8. IF no version set at any level THEN system SHALL return error "no Maven version set"
9. WHEN reading `.maven-version` file THEN system SHALL trim whitespace from version string
10. WHEN version resolved THEN system SHALL verify version is installed before returning
11. IF resolved version not installed THEN system SHALL return error "version 'X' not installed"

### Requirement 6: Set Global Version

**User Story:** As a developer, I want to set a default Maven version system-wide, so that all projects use the same version unless overridden.

#### Acceptance Criteria

1. WHEN user sets global version X THEN system SHALL verify version X is installed
2. IF version not installed THEN system SHALL return error "version 'X' not installed"
3. WHEN setting global THEN system SHALL write version to `%USERPROFILE%\.mvnenv\config\config.yaml`
4. WHEN writing config THEN system SHALL preserve other configuration settings
5. WHEN config file doesn't exist THEN system SHALL create it with default structure
6. WHEN setting completes THEN system SHALL verify version can be resolved

### Requirement 7: Set Local Version

**User Story:** As a developer, I want to set a Maven version for a specific project, so that the project always uses the correct version regardless of global settings.

#### Acceptance Criteria

1. WHEN user sets local version X THEN system SHALL verify version X is installed
2. IF version not installed THEN system SHALL return error "version 'X' not installed"
3. WHEN setting local THEN system SHALL create `.maven-version` file in current directory
4. WHEN creating file THEN system SHALL write only version string (no prefix, single line)
5. IF `.maven-version` exists THEN system SHALL overwrite with new version
6. WHEN file created THEN system SHALL verify version can be resolved from current directory

### Requirement 8: Set Shell Version

**User Story:** As a developer, I want to temporarily override Maven version for my current terminal session, so that I can test different versions without affecting other projects.

#### Acceptance Criteria

1. WHEN user sets shell version X THEN system SHALL verify version X is installed
2. IF version not installed THEN system SHALL return error "version 'X' not installed"
3. WHEN setting shell THEN system SHALL set environment variable `MVNENV_MAVEN_VERSION=X`
4. WHEN environment variable set THEN system SHALL provide instructions to user for persisting in session
5. WHEN shell version set THEN system SHALL take precedence over local and global versions
6. WHEN user unsets shell version THEN system SHALL remove `MVNENV_MAVEN_VERSION` variable
7. IF no shell variable THEN unset operation SHALL complete silently (not an error)

### Requirement 9: Get Latest Version

**User Story:** As a developer, I want to find the latest available or installed Maven version, so that I can stay up-to-date.

#### Acceptance Criteria

1. WHEN user requests latest without prefix THEN system SHALL return newest installed version
2. WHEN user requests latest with prefix (e.g., "3.8") THEN system SHALL return newest version matching prefix
3. WHEN comparing versions THEN system SHALL use semantic versioning rules (3.9.0 > 3.8.10)
4. IF no versions match criteria THEN system SHALL return error "no version found"
5. WHEN no versions installed THEN system SHALL return error "no versions installed"
6. WHEN prefix matches multiple versions THEN system SHALL return newest in that prefix range

### Requirement 10: Update Version Cache

**User Story:** As a developer, I want to refresh the cached list of available versions, so that I can see newly released Maven versions.

#### Acceptance Criteria

1. WHEN user requests cache update THEN system SHALL query all configured repositories
2. WHEN querying THEN system SHALL fetch latest version metadata from each repository
3. WHEN fetch completes THEN system SHALL update cache file with timestamp
4. WHEN repository unavailable THEN system SHALL continue with others and report warnings
5. IF all repositories fail THEN system SHALL return error "failed to update cache"
6. WHEN update completes THEN system SHALL display count of available versions found

## Non-Functional Requirements

### Code Architecture and Modularity

- **Single Responsibility Principle**: Version management logic in `internal/version/`, download logic in `internal/download/`, config in `internal/config/`
- **Modular Design**: Manager, Resolver, Installer, and Downloader as separate components with clear interfaces
- **Dependency Management**: Version management depends on download and config packages, not vice versa
- **Clear Interfaces**: Public methods for install, uninstall, list, resolve operations with well-defined contracts

### Performance

- **Version Resolution**: <100ms from call to resolved version path
- **List Installed**: <100ms to scan and return installed versions
- **List Available**: <500ms when cache is fresh, acceptable up to 5s when querying repositories
- **Installation**: Limited by network speed (accept 30-60s for typical version)
- **Uninstallation**: <5s for typical version removal
- **Cache Lookup**: <10ms for version cache checks

### Security

- **Checksum Verification**: SHA-256 verification mandatory for all downloads, no override option
- **Path Validation**: All file paths must be sanitized to prevent directory traversal
- **Archive Extraction**: Use secure extraction preventing zip-slip vulnerabilities
- **File Permissions**: Version directories readable/executable by user only (Windows ACLs)
- **No Credential Storage**: This spec doesn't handle credentials (nexus-repository-integration spec handles that)

### Reliability

- **Atomic Operations**: Installations are atomic—either complete or fully rolled back
- **Graceful Degradation**: Single repository failure doesn't prevent installation from other repositories
- **Corrupt Installation Detection**: Verify Maven binary exists and is executable after installation
- **Transaction Log**: Log all operations for troubleshooting (install start, extract, verify, complete)
- **Idempotency**: Installing already-installed version returns clear error, doesn't corrupt
- **Concurrent Safety**: Prevent concurrent installations of same version using file locks

### Usability

- **Clear Error Messages**: All errors include actionable information (which step failed, how to resolve)
- **Progress Indication**: Long operations (download, extract) show progress
- **Helpful Suggestions**: "version not installed" errors suggest `mvnenv install -l` to see available
- **Automatic Directory Creation**: Create required directories automatically (versions/, cache/, config/)
- **Clean State Management**: Operations leave system in consistent state even on failure

### Compatibility

- **Windows Paths**: Use Windows path conventions (backslashes, case-insensitive)
- **Long Paths**: Handle Windows MAX_PATH limitations for deep Maven directory structures
- **File Locking**: Respect Windows file locking, provide clear errors when files are in use
- **Archive Formats**: Support ZIP format for Maven distributions (standard Maven packaging)
- **Version Formats**: Support semantic versioning (X.Y.Z) and Maven version schemes (X.Y.Z-qualifier)

### Maintainability

- **Logging**: Comprehensive logging at DEBUG level for all operations
- **Error Context**: Wrap errors with context at each layer for debugging
- **Test Coverage**: >90% coverage for version resolution, comparison, and installation logic
- **Documentation**: Godoc comments for all exported functions and types

## Technical Constraints

- **Go Version**: Use Go 1.21+ standard library features
- **No External Executables**: Don't shell out to unzip, curl, etc.—use Go libraries
- **Windows-Specific**: Can use Windows-specific code (registry, ACLs) in `*_windows.go` files
- **Configuration Format**: YAML for config.yaml, plain text for .maven-version
- **Cache Location**: `%USERPROFILE%\.mvnenv\cache\` for downloaded archives and version lists
- **Installation Root**: `%USERPROFILE%\.mvnenv\versions\` for extracted Maven installations

## Out of Scope

The following are explicitly NOT part of this spec (covered in other specs):
- Repository configuration and Nexus authentication (nexus-repository-integration spec)
- Shim generation and execution (shim-system-implementation spec)
- CLI command implementation (cli-commands spec)
- Downloading logic (implemented in internal/download package, shared)
- Configuration file parsing (internal/config package, shared)

## Success Criteria

1. Maven versions can be installed from any configured repository
2. All downloads are verified with SHA-256 checksums
3. Installations are atomic (complete or fully rolled back)
4. Installed versions are listed correctly with current version indicated
5. Available versions are listed from all repositories
6. Version resolution follows shell > local > global hierarchy correctly
7. Global, local, and shell versions can be set and retrieved
8. Latest version determination works with and without prefixes
9. Version cache updates successfully from repositories
10. All operations have <100ms overhead (excluding network I/O)
11. >90% test coverage for all version management logic
12. Clear, actionable error messages for all failure scenarios
