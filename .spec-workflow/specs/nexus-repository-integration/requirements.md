# Requirements Document

## Introduction

The Nexus Repository Integration feature provides the capability to download Maven distributions from private Nexus Repository Manager instances instead of (or in addition to) the public Apache Maven archives. This enables enterprise teams to use mvnenv-win with their corporate infrastructure, supporting authentication, SSL/TLS, multiple repositories with priority ordering, and version discovery from Nexus metadata.

This spec focuses on Nexus repository client functionality—configuration, authentication, version discovery, and distribution downloads. It integrates with the core-version-management spec (which consumes the download functionality) and the cli-commands spec (which exposes repository management commands).

## Alignment with Product Vision

**Product Principle: Enterprise Ready**
The requirement for Nexus integration with authentication, SSL/TLS, and Windows Credential Manager storage directly supports the principle of being "built for corporate environments from day one."

**Business Objective: Corporate Compliance**
Repository configuration and authenticated downloads from private Nexus instances enable organizations to enforce security policies requiring Maven distributions to come from approved internal sources rather than public internet.

**Key Feature: Nexus Integration**
This spec directly implements the core product feature "Native support for downloading Maven distributions from private Nexus repositories with authentication."

**Target User: Enterprise Teams**
Organizations using private Maven repositories via Nexus are a primary target user group, and this spec enables their workflow.

## Requirements

### Requirement 1: Repository Configuration Management

**User Story:** As a developer, I want to configure multiple Nexus repositories, so that I can download Maven from my organization's approved sources.

#### Acceptance Criteria

1. WHEN user adds repository THEN system SHALL validate repository name is alphanumeric with hyphens/underscores
2. WHEN adding repository THEN system SHALL validate URL is valid HTTP/HTTPS endpoint
3. WHEN adding repository THEN system SHALL save to `%USERPROFILE%\.mvnenv\config\repositories.yaml`
4. IF repository name already exists THEN system SHALL return error "repository 'X' already exists"
5. WHEN adding repository THEN system SHALL support optional priority parameter (integer 1-100, default 50)
6. WHEN listing repositories THEN system SHALL display name, URL, priority, enabled status, and auth status
7. WHEN listing repositories THEN system SHALL sort by priority ascending (lower number = higher priority)
8. WHEN removing repository THEN system SHALL verify repository exists before removal
9. IF repository not found THEN system SHALL return error "repository 'X' not found"
10. WHEN removing repository THEN system SHALL preserve other repositories in configuration
11. WHEN repository configuration invalid THEN system SHALL provide specific validation error

### Requirement 2: Repository Authentication

**User Story:** As a developer, I want to configure authentication for Nexus repositories, so that I can access private distributions that require credentials.

#### Acceptance Criteria

1. WHEN configuring auth THEN system SHALL support basic authentication (username/password)
2. WHEN configuring auth THEN system SHALL support token authentication (bearer token)
3. WHEN user provides credentials THEN system SHALL store in Windows Credential Manager with target prefix `mvnenv:repo:`
4. WHEN storing credentials THEN system SHALL use repository name as credential identifier
5. IF Windows Credential Manager unavailable THEN system SHALL fall back to environment variable reference in config
6. WHEN reading config THEN system SHALL support `${ENV_VAR}` syntax for environment variable substitution
7. WHEN authenticating request THEN system SHALL first check Windows Credential Manager, then environment variables
8. IF credentials invalid during download THEN system SHALL return error with HTTP status code
9. WHEN removing repository THEN system SHALL also remove associated credentials from Credential Manager
10. WHEN listing repositories THEN system SHALL indicate auth status without exposing credentials

### Requirement 3: Version Discovery from Nexus

**User Story:** As a developer, I want to see which Maven versions are available in my Nexus repositories, so that I can install approved versions.

#### Acceptance Criteria

1. WHEN discovering versions THEN system SHALL query all enabled repositories in priority order
2. WHEN querying repository THEN system SHALL request `maven-metadata.xml` from Maven artifact path
3. WHEN parsing metadata THEN system SHALL extract `<version>` elements from `<versioning><versions>` section
4. WHEN combining results THEN system SHALL deduplicate versions across repositories
5. IF repository request fails THEN system SHALL log warning and continue with next repository
6. IF all repositories fail THEN system SHALL return error "no repositories available"
7. WHEN repository returns 401/403 THEN system SHALL return error "authentication failed for repository 'X'"
8. WHEN repository returns 404 THEN system SHALL log warning "Maven metadata not found at 'URL'"
9. WHEN discovery succeeds THEN system SHALL cache results with repository source information
10. WHEN version requested THEN system SHALL know which repository provides that version

### Requirement 4: Maven Distribution Download from Nexus

**User Story:** As a developer, I want to download Maven distributions from Nexus, so that installations use my organization's approved sources.

#### Acceptance Criteria

1. WHEN downloading from Nexus THEN system SHALL construct artifact URL: `{repo-url}/org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip`
2. WHEN downloading THEN system SHALL include authentication headers if configured for repository
3. WHEN download starts THEN system SHALL show progress with bytes downloaded and total size
4. IF download fails with 401/403 THEN system SHALL return error "authentication failed"
5. IF download fails with 404 THEN system SHALL try next repository in priority order
6. IF all repositories fail THEN system SHALL return error "version 'X' not found in any repository"
7. WHEN download completes THEN system SHALL verify Content-Type is `application/zip` or `application/octet-stream`
8. WHEN Nexus provides checksum THEN system SHALL download and verify SHA-256 checksum from `{artifact-url}.sha256`
9. IF checksum file exists and verification fails THEN system SHALL return error "checksum mismatch"
10. WHEN download succeeds THEN system SHALL cache archive with repository source metadata

### Requirement 5: SSL/TLS Certificate Handling

**User Story:** As a developer, I want to securely connect to Nexus over HTTPS with proper certificate validation, so that my downloads are protected from tampering.

#### Acceptance Criteria

1. WHEN connecting to HTTPS repository THEN system SHALL verify TLS certificate by default
2. WHEN certificate invalid THEN system SHALL return error with certificate details
3. WHEN using self-signed certificate THEN system SHALL support certificate import via configuration
4. WHEN custom CA certificate configured THEN system SHALL load from file path in config
5. IF CA certificate file not found THEN system SHALL return error "CA certificate not found at path"
6. WHEN certificate verification disabled THEN system SHALL require explicit `insecure: true` flag in config
7. WHEN insecure mode enabled THEN system SHALL log warning on every connection
8. WHEN TLS handshake fails THEN system SHALL return error with TLS version and cipher details
9. WHEN repository uses HTTP (not HTTPS) THEN system SHALL log warning about unencrypted connection

### Requirement 6: Repository Priority and Fallback

**User Story:** As a developer, I want repositories to be tried in priority order, so that preferred sources are used first with automatic fallback.

#### Acceptance Criteria

1. WHEN multiple repositories configured THEN system SHALL try in priority order (lower number first)
2. WHEN repositories have same priority THEN system SHALL use alphabetical order by name
3. WHEN repository fails THEN system SHALL automatically try next repository without user intervention
4. WHEN repository disabled THEN system SHALL skip that repository
5. WHEN repository responds slowly THEN system SHALL timeout after configured duration (default 30s)
6. IF all repositories timeout THEN system SHALL return error "all repositories timed out"
7. WHEN version found in higher-priority repository THEN system SHALL not query lower-priority repositories
8. WHEN downloading THEN system SHALL use the repository that provided the version during discovery
9. IF that repository fails THEN system SHALL try other repositories that also have that version

### Requirement 7: Metadata Caching

**User Story:** As a developer, I want available versions to be cached locally, so that I don't have to wait for repository queries every time.

#### Acceptance Criteria

1. WHEN versions discovered THEN system SHALL cache to `%USERPROFILE%\.mvnenv\cache\repo-metadata.json`
2. WHEN writing cache THEN system SHALL include timestamp and repository source per version
3. WHEN reading cache THEN system SHALL check if cache is fresh (default TTL: 24 hours)
4. IF cache fresh THEN system SHALL use cached versions without querying repositories
5. IF cache stale THEN system SHALL refresh from repositories and update cache
6. WHEN cache file corrupt THEN system SHALL delete and rebuild cache
7. WHEN user runs `mvnenv update` THEN system SHALL force cache refresh regardless of TTL
8. WHEN repository configuration changes THEN system SHALL invalidate cache
9. WHEN listing versions THEN system SHALL indicate if results are from cache with timestamp

### Requirement 8: Apache Maven Fallback

**User Story:** As a developer, I want mvnenv-win to fall back to Apache Maven archives, so that I can install versions even if Nexus repositories are unavailable.

#### Acceptance Criteria

1. WHEN no Nexus repositories configured THEN system SHALL use Apache Maven archive as default
2. WHEN Nexus repositories configured THEN system SHALL still include Apache archive as lowest-priority fallback
3. WHEN querying Apache archive THEN system SHALL use HTTP directory listing at `https://archive.apache.org/dist/maven/maven-3/`
4. WHEN parsing Apache archive THEN system SHALL extract version directories from HTML directory listing
5. WHEN downloading from Apache THEN system SHALL construct URL: `https://archive.apache.org/dist/maven/maven-3/{version}/binaries/apache-maven-{version}-bin.zip`
6. WHEN downloading from Apache THEN system SHALL verify SHA-512 checksum from `.sha512` file
7. IF Apache archive unavailable THEN system SHALL return error only after all sources exhausted
8. WHEN Apache archive used THEN system SHALL log "downloading from Apache Maven archive"

## Non-Functional Requirements

### Code Architecture and Modularity

- **Single Responsibility**: Nexus client logic in `internal/nexus/`, separate from download logic and version management
- **Repository Pattern**: Abstract `Repository` interface with implementations for Nexus and Apache archives
- **Configuration Isolation**: Repository configuration in separate `repositories.yaml` file
- **Credential Isolation**: Credential management in separate module with platform-specific implementations

### Performance

- **Version Discovery**: <500ms when cache is fresh, <5s when querying multiple repositories
- **Download Speed**: Limited by network speed, no artificial throttling
- **Timeout Handling**: 30s default timeout per repository, configurable
- **Connection Pooling**: Reuse HTTP connections for multiple requests to same repository
- **Parallel Queries**: Query multiple repositories in parallel (not sequentially) when cache invalid

### Security

- **TLS Verification**: Mandatory by default, explicit flag required to disable
- **Credential Storage**: Windows Credential Manager preferred, environment variables as fallback
- **No Plaintext Passwords**: Passwords never stored in plaintext in configuration files
- **Audit Logging**: Log all authentication attempts and repository access
- **Certificate Validation**: Full certificate chain validation including hostname verification

### Reliability

- **Graceful Degradation**: Single repository failure doesn't prevent using other repositories
- **Automatic Retry**: Transient network errors trigger automatic retry with exponential backoff
- **Circuit Breaker**: Temporarily skip repository after repeated failures
- **Cache Resilience**: Corrupted cache file doesn't break functionality
- **Configuration Validation**: Invalid repository configuration detected at add time, not runtime

### Usability

- **Clear Error Messages**: Authentication failures include repository name and HTTP status code
- **Helpful Suggestions**: "repository not found" errors suggest `mvnenv repo list` command
- **Progress Indication**: Download progress shown with percentage and speed
- **Automatic Directory Creation**: Create cache and config directories automatically
- **Default Configuration**: Apache Maven archive available by default without configuration

### Compatibility

- **Nexus Versions**: Support Nexus Repository Manager 2.x and 3.x
- **HTTP Standards**: Follow standard HTTP authentication and redirect behavior
- **Maven Metadata Format**: Parse standard Maven metadata.xml format
- **Windows Credential Manager**: Use native Windows credential APIs

### Maintainability

- **Logging**: Comprehensive DEBUG-level logging for repository interactions
- **Error Context**: Wrap errors with repository name and operation context
- **Test Coverage**: >90% coverage for Nexus client, authentication, and caching logic
- **Documentation**: Godoc comments for all exported types and functions

## Technical Constraints

- **Go Version**: Use Go 1.21+ standard library HTTP client
- **HTTP Library**: Use `github.com/go-resty/resty/v2` for HTTP requests with retries
- **Credential Storage**: Use Windows Credential Manager via `syscall` or `github.com/danieljoos/wincred`
- **Configuration Format**: YAML for repositories.yaml, JSON for cache files
- **Cache Location**: `%USERPROFILE%\.mvnenv\cache\` for metadata cache
- **Config Location**: `%USERPROFILE%\.mvnenv\config\` for repository configuration

## Out of Scope

The following are explicitly NOT part of this spec:
- Maven settings.xml management (future feature)
- Nexus repository creation or administration (external to mvnenv-win)
- Nexus authentication configuration on the server side
- Repository mirroring or proxying functionality
- Maven plugin or dependency downloads (only Maven distribution downloads)
- Integration with Maven central or other artifact repositories

## Success Criteria

1. Nexus repositories can be configured with name, URL, priority, and authentication
2. Credentials stored securely in Windows Credential Manager
3. Version discovery successfully queries Nexus metadata and returns available versions
4. Maven distributions download from Nexus with authentication
5. SSL/TLS certificates properly validated with support for custom CAs
6. Repository priority and fallback work correctly (higher priority tried first)
7. Metadata cached with configurable TTL to reduce network queries
8. Apache Maven archive used as default/fallback when Nexus unavailable
9. All operations have clear error messages with repository context
10. >90% test coverage for Nexus client and repository management
11. Integration with core-version-management for seamless version installation
12. CLI commands (repo add/list/remove/auth) fully functional
