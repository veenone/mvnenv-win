# Tasks Document

This document lists all implementation tasks for the Nexus Repository Integration feature. Each task includes detailed prompts for AI-assisted development.

## Task List

### Foundation Tasks

#### Task 1: Create Repository Interface and Core Types
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/repository.go`

**Description:** Define the Repository interface and core types (RepositoryType, ProgressCallback, RepositoryInfo, VersionInfo) that provide the abstraction layer for all Maven distribution sources.

**Dependencies:** None

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are implementing the Repository interface for mvnenv-win's Nexus integration feature.

**Role:** Senior Go backend developer with expertise in repository pattern and API design.

**Task:** Create internal/nexus/repository.go with:

1. Repository interface with methods:
   - Name() string
   - Priority() int
   - IsEnabled() bool
   - DiscoverVersions(ctx context.Context) ([]string, error)
   - DownloadDistribution(ctx context.Context, version string, destDir string, progress ProgressCallback) (string, error)
   - VerifyChecksum(ctx context.Context, archivePath string, version string) (bool, error)
   - HasVersion(version string) (bool, error)

2. Core types:
   - RepositoryType (nexus, apache)
   - ProgressCallback func(downloaded int64, total int64)
   - RepositoryInfo struct (Name, URL, Type, Priority, Enabled, HasAuth, LastChecked, Available)
   - VersionInfo struct (Version, Source, SourceType)

3. Common error variables:
   - ErrRepositoryNotFound
   - ErrRepositoryExists
   - ErrAuthenticationFailed
   - ErrVersionNotFound
   - ErrNoRepositoriesAvailable
   - ErrInvalidConfiguration
   - ErrCertificateInvalid

4. RepositoryError struct wrapping errors with repository context

**Restrictions:**
- Use Go 1.21+ features
- Follow Go standard error handling patterns
- Include comprehensive godoc comments
- Use const for RepositoryType values

**Success Criteria:**
- Repository interface clearly defines contract for all implementations
- Types are exported and well-documented
- Error types provide helpful context
- Code passes go vet and staticcheck

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 2: Create HTTPClient with Retry Logic
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/http.go`

**Description:** Implement HTTPClient wrapper using go-resty with retry logic, timeout handling, and TLS configuration support.

**Dependencies:** Task 1

**Estimated Effort:** 3 hours

**_Prompt:**
```
You are implementing the HTTPClient wrapper for mvnenv-win's Nexus integration.

**Role:** Senior Go developer with expertise in HTTP clients and error handling.

**Task:** Create internal/nexus/http.go with:

1. HTTPClient struct:
   - client *resty.Client
   - timeout time.Duration
   - retryCount int
   - retryWait time.Duration

2. NewHTTPClient(timeout time.Duration, tlsConfig *TLSConfig) *HTTPClient
   - Initialize resty client with timeout
   - Configure retry: 3 attempts, exponential backoff (1s, 2s, 4s)
   - Retry on: connection errors, timeouts, 5xx (except 501)
   - Don't retry on: 4xx (except 429)

3. Get(ctx context.Context, url string, auth *AuthConfig) (*resty.Response, error)
   - Execute GET with context cancellation
   - Apply authentication headers if provided
   - Return response or wrapped error

4. Download(ctx context.Context, url string, destPath string, auth *AuthConfig, progress ProgressCallback) error
   - Stream download to file
   - Call progress callback with downloaded/total bytes
   - Handle context cancellation
   - Clean up partial file on error

5. SetTLSConfig(tlsConfig *TLSConfig)
   - Update client TLS configuration
   - Load custom CA certificates if specified
   - Set InsecureSkipVerify if configured

6. TLSConfig struct:
   - Insecure bool
   - CAFile string
   - ClientCertFile string (for future)
   - ClientKeyFile string (for future)

7. AuthConfig struct:
   - Type AuthType (none, basic, token)
   - Username string
   - Password string
   - Token string

**Restrictions:**
- Use github.com/go-resty/resty/v2
- Implement proper context cancellation
- Clean up resources on error
- Log retry attempts at DEBUG level

**Success Criteria:**
- Retry logic works for transient errors
- Progress callback invoked during downloads
- TLS configuration properly applied
- Context cancellation stops requests immediately

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 3: Create CredentialManager for Windows Credential Manager
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/credentials.go`
- `internal/nexus/credentials_windows.go`

**Description:** Implement secure credential storage using Windows Credential Manager APIs.

**Dependencies:** None

**Estimated Effort:** 3 hours

**_Prompt:**
```
You are implementing secure credential storage for mvnenv-win.

**Role:** Senior Go developer with Windows platform expertise and security focus.

**Task:** Create credential management with Windows Credential Manager integration:

1. internal/nexus/credentials.go:
   - CredentialManager struct with targetPrefix string
   - NewCredentialManager() *CredentialManager
   - Interface methods (to be implemented in platform-specific file)

2. internal/nexus/credentials_windows.go:
   - StoreCredentials(repoName string, username string, password string) error
     * Target: "mvnenv:repo:{repoName}"
     * Use CRED_TYPE_GENERIC
     * Store username in TargetName or Attributes
     * Store password in CredentialBlob (encrypted by OS)

   - RetrieveCredentials(repoName string) (string, string, error)
     * Read from Windows Credential Manager
     * Return username, password
     * Return error if not found

   - DeleteCredentials(repoName string) error
     * Remove credential from Windows Credential Manager

   - HasCredentials(repoName string) bool
     * Check if credential exists without reading it

3. Use github.com/danieljoos/wincred OR syscall to Windows Credential Manager APIs:
   - CredWriteW for storing
   - CredReadW for retrieving
   - CredDeleteW for deletion

4. Error handling:
   - If Credential Manager unavailable, return descriptive error
   - Suggest environment variables as fallback in error message

**Restrictions:**
- Windows-specific code only in *_windows.go
- Handle credential not found gracefully
- Never log credentials
- Use proper Windows API error codes

**Success Criteria:**
- Credentials stored encrypted by OS
- Target name format: "mvnenv:repo:{repoName}"
- Retrieval works after restart
- Deletion removes credentials completely

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

### Core Components

#### Task 4: Create RepositoryConfig for repositories.yaml
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/config.go`

**Description:** Implement configuration management for repositories.yaml with environment variable substitution.

**Dependencies:** Task 1

**Estimated Effort:** 3 hours

**_Prompt:**
```
You are implementing repository configuration management for mvnenv-win.

**Role:** Senior Go developer with configuration management expertise.

**Task:** Create internal/nexus/config.go with:

1. RepositoryConfig struct:
   - configPath string
   - mu sync.RWMutex

2. Config struct (YAML structure):
   - Repositories []RepositoryEntry

3. RepositoryEntry struct:
   - Name string
   - Type RepositoryType
   - URL string
   - Priority int
   - Enabled bool
   - Auth *AuthConfig (optional)
   - TLS *TLSConfig (optional)

4. NewRepositoryConfig(configPath string) *RepositoryConfig
   - Initialize with config file path
   - Create directory if doesn't exist

5. Load() (*Config, error)
   - Read YAML file
   - Parse into Config struct
   - Apply environment variable substitution
   - Return default config if file doesn't exist (Apache archive only)

6. Save(config *Config) error
   - Marshal Config to YAML
   - Write atomically (temp file + rename)
   - Preserve file permissions

7. AddRepository(repo RepositoryEntry) error
   - Validate: name alphanumeric + hyphens/underscores
   - Validate: URL is valid HTTP/HTTPS
   - Check if name already exists
   - Add to repositories list
   - Save config

8. RemoveRepository(name string) error
   - Find repository by name
   - Remove from list
   - Save config
   - Return error if not found

9. expandEnvVars(value string) string
   - Replace ${VAR_NAME} with os.Getenv(VAR_NAME)
   - Use regex: \$\{([A-Z_][A-Z0-9_]*)\}

**Restrictions:**
- Use gopkg.in/yaml.v3 for YAML parsing
- Thread-safe with RWMutex
- Atomic file writes (temp + rename)
- Validate repository names (no path traversal)

**Success Criteria:**
- YAML round-trips correctly
- Environment variables substituted
- Default config (Apache archive) provided if no file
- Concurrent access is safe

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 5: Create MetadataCache for Version Caching
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/metadata.go`

**Description:** Implement metadata caching to reduce network queries with TTL-based freshness checks.

**Dependencies:** Task 1

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are implementing metadata caching for mvnenv-win's repository integration.

**Role:** Senior Go developer with caching and performance optimization expertise.

**Task:** Create internal/nexus/metadata.go with:

1. MetadataCache struct:
   - cachePath string
   - ttl time.Duration
   - mu sync.RWMutex

2. CachedMetadata struct (JSON structure):
   - Timestamp time.Time
   - Versions []CachedVersionInfo
   - Repos map[string]RepoSnapshot

3. CachedVersionInfo struct:
   - Version string
   - Source string (repository name)
   - SourceType RepositoryType

4. RepoSnapshot struct:
   - URL string
   - Priority int
   - Available bool
   - Timestamp time.Time

5. NewMetadataCache(cachePath string, ttl time.Duration) *MetadataCache
   - Initialize with cache file path and TTL
   - Default TTL: 24 hours

6. IsFresh() bool
   - Load cache file
   - Check if timestamp within TTL
   - Return false if file doesn't exist or parse error

7. Load() (*CachedMetadata, error)
   - Read JSON cache file
   - Parse into CachedMetadata
   - Return error if file doesn't exist or corrupt

8. Save(metadata *CachedMetadata) error
   - Set current timestamp
   - Marshal to JSON with indentation
   - Write atomically (temp file + rename)

9. Invalidate() error
   - Delete cache file
   - Return nil if file doesn't exist

**Restrictions:**
- Use encoding/json for JSON parsing
- Thread-safe with RWMutex
- Atomic file writes
- Handle corrupted cache gracefully

**Success Criteria:**
- Freshness check works with TTL
- Cache survives process restarts
- Corrupted cache doesn't break functionality
- Concurrent access is safe

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 6: Create NexusClient Implementation
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/nexus_client.go`

**Description:** Implement Repository interface for Nexus Repository Manager 2.x and 3.x.

**Dependencies:** Tasks 1, 2, 3

**Estimated Effort:** 4 hours

**_Prompt:**
```
You are implementing NexusClient for mvnenv-win's Nexus integration.

**Role:** Senior Go developer with REST API and Maven repository expertise.

**Task:** Create internal/nexus/nexus_client.go implementing Repository interface:

1. NexusClient struct:
   - name string
   - baseURL string
   - priority int
   - enabled bool
   - authType AuthType
   - httpClient *HTTPClient
   - credentials *CredentialManager
   - tlsConfig *TLSConfig

2. NewNexusClient(...) *NexusClient
   - Initialize all fields
   - Create HTTPClient with TLS config

3. DiscoverVersions(ctx context.Context) ([]string, error)
   - URL: {baseURL}/org/apache/maven/apache-maven/maven-metadata.xml
   - GET request with authentication
   - Parse XML: <metadata><versioning><versions><version>3.9.4</version>...
   - Extract all <version> elements
   - Handle errors: 401/403 (auth failed), 404 (not found), network errors

4. DownloadDistribution(ctx context.Context, version string, destDir string, progress ProgressCallback) (string, error)
   - URL: {baseURL}/org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip
   - Destination: {destDir}/apache-maven-{version}-bin.zip
   - Use httpClient.Download with authentication and progress callback
   - Verify Content-Type is application/zip or application/octet-stream
   - Return file path on success

5. VerifyChecksum(ctx context.Context, archivePath string, version string) (bool, error)
   - URL: {baseURL}/org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip.sha256
   - Download checksum file
   - Calculate SHA-256 of archive
   - Compare checksums
   - Return true if match, false if mismatch
   - Return error if checksum file not available

6. HasVersion(version string) (bool, error)
   - Quick check without full discovery (future optimization)
   - For now: return false, nil (always query)

7. Name(), Priority(), IsEnabled() - simple getters

8. Authentication handling:
   - Retrieve credentials from CredentialManager
   - Fall back to environment variables from config
   - Apply to httpClient requests

**Restrictions:**
- Use encoding/xml for XML parsing
- Handle context cancellation in all network operations
- Log all network operations at DEBUG level
- Return wrapped errors with repository name

**Success Criteria:**
- Successfully parses Nexus metadata XML
- Downloads distributions with authentication
- Checksum verification works
- Clear error messages for auth failures

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 7: Create ApacheArchiveClient Implementation
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/apache_client.go`

**Description:** Implement Repository interface for Apache Maven official archives.

**Dependencies:** Tasks 1, 2

**Estimated Effort:** 3 hours

**_Prompt:**
```
You are implementing ApacheArchiveClient for mvnenv-win.

**Role:** Senior Go developer with web scraping and HTTP client expertise.

**Task:** Create internal/nexus/apache_client.go implementing Repository interface:

1. ApacheArchiveClient struct:
   - name string
   - baseURL string (https://archive.apache.org/dist/maven/maven-3/)
   - priority int
   - enabled bool
   - httpClient *HTTPClient

2. NewApacheArchiveClient(name string, priority int) *ApacheArchiveClient
   - Initialize with Apache archive URL
   - Create HTTPClient (no auth needed)

3. DiscoverVersions(ctx context.Context) ([]string, error)
   - GET: https://archive.apache.org/dist/maven/maven-3/
   - Parse HTML directory listing
   - Regex pattern: <a href="(\d+\.\d+\.\d+)/">
   - Extract version numbers
   - Filter to valid semantic versions
   - Return sorted list

4. DownloadDistribution(ctx context.Context, version string, destDir string, progress ProgressCallback) (string, error)
   - URL: {baseURL}/{version}/binaries/apache-maven-{version}-bin.zip
   - Destination: {destDir}/apache-maven-{version}-bin.zip
   - Use httpClient.Download with progress callback
   - Return file path on success

5. VerifyChecksum(ctx context.Context, archivePath string, version string) (bool, error)
   - URL: {baseURL}/{version}/binaries/apache-maven-{version}-bin.zip.sha512
   - Download checksum file
   - Calculate SHA-512 of archive
   - Compare checksums (handle both formats: "checksum filename" and just "checksum")
   - Return true if match, false if mismatch

6. HasVersion(version string) (bool, error)
   - Return false, nil (always query)

7. Name(), Priority(), IsEnabled() - simple getters

**Restrictions:**
- Use regex for HTML parsing (simple directory listing)
- No authentication needed
- Handle missing checksum files gracefully
- SHA-512 for Apache (not SHA-256)

**Success Criteria:**
- Successfully scrapes Apache directory listing
- Downloads distributions correctly
- SHA-512 checksum verification works
- No authentication headers sent

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 8: Create RepositoryManager
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/manager.go`

**Description:** Implement RepositoryManager that orchestrates multiple repositories with priority-based selection and fallback.

**Dependencies:** Tasks 1, 4, 5, 6, 7

**Estimated Effort:** 5 hours

**_Prompt:**
```
You are implementing RepositoryManager for mvnenv-win's Nexus integration.

**Role:** Senior Go backend developer with distributed systems expertise.

**Task:** Create internal/nexus/manager.go with:

1. RepositoryManager struct:
   - repositories []Repository
   - config *RepositoryConfig
   - cache *MetadataCache
   - credentials *CredentialManager

2. NewRepositoryManager(configPath string) (*RepositoryManager, error)
   - Load config from configPath
   - Create CredentialManager
   - Create MetadataCache with 24h TTL
   - Instantiate Repository implementations from config:
     * NexusClient for type "nexus"
     * ApacheArchiveClient for type "apache"
   - Add default Apache repository if no config
   - Sort repositories by priority

3. DiscoverVersions(ctx context.Context, forceRefresh bool) ([]VersionInfo, error)
   - If cache fresh and !forceRefresh: return cached versions
   - Query repositories in parallel using goroutines:
     * One goroutine per enabled repository
     * Collect results in channel
     * Wait for all with WaitGroup
     * Continue on individual failures (log warning)
   - If all repositories failed: return ErrNoRepositoriesAvailable
   - Deduplicate versions (keep first = highest priority)
   - Convert to VersionInfo with source repository
   - Update cache with results
   - Return version list

4. DownloadVersion(ctx context.Context, version string, destDir string, progress ProgressCallback) (string, error)
   - Query cache for repository that provides version
   - If found: try that repository first
   - If failed or not in cache: try all repositories in priority order
   - For each repository:
     * Call DownloadDistribution
     * If success: verify checksum (if available), return path
     * If failure: log error, continue to next
   - If all failed: return ErrVersionNotFound

5. AddRepository(name string, url string, priority int, repoType RepositoryType) error
   - Validate name (alphanumeric + hyphens/underscores)
   - Validate URL format
   - Check if name exists
   - Create RepositoryEntry
   - Add via config.AddRepository
   - Invalidate cache
   - Reload repositories

6. RemoveRepository(name string) error
   - Remove via config.RemoveRepository
   - Delete credentials via credentials.DeleteCredentials
   - Invalidate cache
   - Reload repositories

7. SetAuthentication(repoName string, authType AuthType, credentials interface{}) error
   - Find repository by name
   - Store credentials in CredentialManager OR update config with env vars
   - Update repository auth config
   - Save config

8. ListRepositories() []RepositoryInfo
   - Return info for all repositories
   - Include auth status (without exposing credentials)

9. GetRepository(name string) (Repository, error)
   - Find and return repository by name

**Restrictions:**
- Parallel queries with timeout per repository (30s default)
- Thread-safe operations
- Proper error wrapping with context
- Cache invalidation on config changes

**Success Criteria:**
- Parallel queries improve performance
- Fallback works when primary repository fails
- Cache reduces network queries
- Priority ordering respected
- Configuration changes take effect immediately

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

### Integration Tasks

#### Task 9: Integrate with core-version-management
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Modify:**
- `internal/version/installer.go`
- `internal/download/downloader.go`

**Description:** Update VersionInstaller to use RepositoryManager for downloads instead of direct HTTP client.

**Dependencies:** Task 8, core-version-management spec completed

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are integrating Nexus repository support into mvnenv-win's version installer.

**Role:** Senior Go developer with integration and refactoring expertise.

**Task:** Update version installation to use RepositoryManager:

1. Modify internal/version/installer.go:
   - Add repoManager *nexus.RepositoryManager field to VersionInstaller struct
   - Update NewVersionInstaller to accept RepositoryManager
   - Update InstallVersion method:
     * Replace direct download logic with repoManager.DownloadVersion
     * Pass progress callback to DownloadVersion
     * Handle repository errors appropriately

2. Update internal/download/downloader.go (if separate):
   - Refactor to use RepositoryManager as download source
   - Maintain progress tracking and error handling

3. Error handling:
   - Wrap repository errors with installation context
   - Provide helpful messages: "version not found in any repository"
   - Suggest checking `mvnenv install -l` for available versions

**Restrictions:**
- Maintain backward compatibility with existing installer tests
- Preserve progress tracking functionality
- Don't break atomic installation logic

**Success Criteria:**
- Version installation works with Nexus repositories
- Downloads fall back to Apache archive if Nexus fails
- Progress tracking still works
- All existing installer tests pass

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 10: Integrate with CLI Commands (repo add/list/remove/auth)
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Modify:**
- `cmd/mvnenv/cmd/repo_add.go`
- `cmd/mvnenv/cmd/repo_list.go`
- `cmd/mvnenv/cmd/repo_remove.go`
- `cmd/mvnenv/cmd/repo_auth.go`

**Description:** Wire up CLI commands to use RepositoryManager for repository management operations.

**Dependencies:** Task 8, cli-commands spec completed

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are implementing CLI command handlers for repository management in mvnenv-win.

**Role:** Senior Go developer with CLI application expertise.

**Task:** Wire up repository CLI commands to RepositoryManager:

1. cmd/mvnenv/cmd/repo_add.go:
   - Parse: name, url, --priority, --type flags
   - Create RepositoryManager
   - Call repoManager.AddRepository
   - Output: "Repository '{name}' added successfully"
   - Handle errors with clear messages

2. cmd/mvnenv/cmd/repo_list.go:
   - Create RepositoryManager
   - Call repoManager.ListRepositories
   - Format output as table:
     NAME            TYPE     URL                              PRIORITY  ENABLED  AUTH
     corporate-nexus nexus    https://nexus.company.com/...    1         yes      yes
     apache          apache   https://archive.apache.org/...   100       yes      no
   - Indicate current/active repository with marker

3. cmd/mvnenv/cmd/repo_remove.go:
   - Parse: name argument
   - Create RepositoryManager
   - Call repoManager.RemoveRepository
   - Output: "Repository '{name}' removed successfully"
   - Also note: "Credentials removed from Windows Credential Manager"

4. cmd/mvnenv/cmd/repo_auth.go:
   - Parse: name, --username, --password, --token flags
   - Prompt for password if not provided (hidden input)
   - Create RepositoryManager
   - Determine AuthType from flags
   - Call repoManager.SetAuthentication
   - Output: "Authentication configured for repository '{name}'"
   - Never log credentials

**Restrictions:**
- Follow pyenv-win-style plain text output (no emojis)
- Use tabwriter for aligned columns
- Hide password input (use terminal.ReadPassword)
- Clear error messages matching requirements.md

**Success Criteria:**
- Commands match pyenv-win style
- Table output is properly aligned
- Password input is hidden
- All repository operations work end-to-end

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

### Testing Tasks

#### Task 11: Unit Tests for NexusClient
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/nexus_client_test.go`

**Description:** Comprehensive unit tests for NexusClient including metadata parsing and download operations.

**Dependencies:** Task 6

**Estimated Effort:** 3 hours

**_Prompt:**
```
You are writing unit tests for NexusClient in mvnenv-win.

**Role:** Senior Go developer with testing expertise and TDD mindset.

**Task:** Create internal/nexus/nexus_client_test.go with comprehensive test coverage:

1. TestNexusClient_DiscoverVersions:
   - Mock HTTP server returning maven-metadata.xml
   - Verify versions correctly parsed
   - Test: valid XML with multiple versions
   - Test: empty versions list
   - Test: malformed XML (error handling)
   - Test: 404 response (not found)
   - Test: 401/403 response (auth failed)

2. TestNexusClient_DownloadDistribution:
   - Mock HTTP server serving ZIP file
   - Verify file downloaded to correct path
   - Test: successful download with progress callback
   - Test: context cancellation mid-download
   - Test: 404 response (version not found)
   - Test: authentication header included

3. TestNexusClient_VerifyChecksum:
   - Mock server with .sha256 file
   - Test: valid checksum (match)
   - Test: invalid checksum (mismatch)
   - Test: checksum file not available (graceful handling)

4. TestNexusClient_Authentication:
   - Test: basic auth headers correctly set
   - Test: token auth (Bearer token)
   - Test: credential retrieval from CredentialManager
   - Test: fallback to environment variables

**Restrictions:**
- Use httptest.Server for mock HTTP server
- Use testify/assert for assertions
- Test all error paths
- Achieve >90% code coverage

**Success Criteria:**
- All tests pass
- Code coverage >90%
- Tests are deterministic and fast (<1s total)
- Mock server properly simulates Nexus responses

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 12: Unit Tests for ApacheArchiveClient
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/apache_client_test.go`

**Description:** Unit tests for ApacheArchiveClient including directory scraping and SHA-512 verification.

**Dependencies:** Task 7

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are writing unit tests for ApacheArchiveClient in mvnenv-win.

**Role:** Senior Go developer with testing expertise.

**Task:** Create internal/nexus/apache_client_test.go with:

1. TestApacheArchiveClient_DiscoverVersions:
   - Mock Apache archive HTML directory listing
   - Test: parse versions from HTML links
   - Test: filter invalid version formats
   - Test: handle network errors gracefully

2. TestApacheArchiveClient_DownloadDistribution:
   - Mock server serving ZIP file
   - Verify download path construction
   - Test: successful download
   - Test: 404 response (version not available)

3. TestApacheArchiveClient_VerifyChecksum:
   - Mock server with .sha512 file
   - Test: SHA-512 verification (match)
   - Test: SHA-512 verification (mismatch)
   - Test: both checksum formats ("checksum filename" and "checksum")

**Restrictions:**
- Use httptest.Server
- Test HTML parsing with various directory formats
- SHA-512 (not SHA-256) for Apache

**Success Criteria:**
- All tests pass
- Directory scraping handles format variations
- SHA-512 verification correct
- Code coverage >90%

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 13: Unit Tests for RepositoryManager
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/manager_test.go`

**Description:** Unit tests for RepositoryManager including priority ordering, fallback, and caching.

**Dependencies:** Task 8

**Estimated Effort:** 4 hours

**_Prompt:**
```
You are writing unit tests for RepositoryManager in mvnenv-win.

**Role:** Senior Go developer with expertise in testing complex orchestration logic.

**Task:** Create internal/nexus/manager_test.go with:

1. TestRepositoryManager_DiscoverVersions_Parallel:
   - Mock multiple repositories
   - Verify parallel queries (use timing to confirm)
   - Test: combine results from multiple repos
   - Test: deduplicate versions (keep highest priority)

2. TestRepositoryManager_DiscoverVersions_Cache:
   - Test: fresh cache returns cached versions (no network)
   - Test: stale cache triggers refresh
   - Test: forceRefresh bypasses cache

3. TestRepositoryManager_DiscoverVersions_Fallback:
   - Mock: first repo fails, second succeeds
   - Verify: result contains versions from second repo
   - Test: all repos fail -> ErrNoRepositoriesAvailable

4. TestRepositoryManager_DownloadVersion_Priority:
   - Mock: version available in multiple repos
   - Verify: highest priority repo tried first
   - Test: primary fails, fallback to secondary

5. TestRepositoryManager_DownloadVersion_CacheLookup:
   - Populate cache with version sources
   - Verify: cache consulted first
   - Test: cached source used preferentially

6. TestRepositoryManager_AddRepository:
   - Test: valid repository added successfully
   - Test: duplicate name rejected
   - Test: invalid name rejected (path traversal)
   - Test: invalid URL rejected

7. TestRepositoryManager_RemoveRepository:
   - Test: existing repository removed
   - Test: credentials also deleted
   - Test: non-existent repository -> error

8. TestRepositoryManager_SetAuthentication:
   - Test: credentials stored in CredentialManager
   - Test: environment variable reference stored in config

**Restrictions:**
- Use mock Repository implementations
- Test concurrency correctness (no race conditions)
- Verify cache invalidation on config changes
- Use testify/mock for mocking

**Success Criteria:**
- All tests pass
- Priority ordering verified
- Parallel queries confirmed
- Cache behavior correct
- Code coverage >90%

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 14: Unit Tests for CredentialManager
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/credentials_test.go`

**Description:** Unit tests for Windows Credential Manager integration.

**Dependencies:** Task 3

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are writing unit tests for CredentialManager in mvnenv-win.

**Role:** Senior Go developer with Windows platform testing expertise.

**Task:** Create internal/nexus/credentials_test.go:

1. TestCredentialManager_StoreAndRetrieve:
   - Store credentials with StoreCredentials
   - Retrieve with RetrieveCredentials
   - Verify: username and password match

2. TestCredentialManager_TargetNameFormat:
   - Verify target: "mvnenv:repo:{repoName}"
   - Test with various repository names

3. TestCredentialManager_Delete:
   - Store credentials
   - Delete with DeleteCredentials
   - Verify: retrieval fails after deletion

4. TestCredentialManager_HasCredentials:
   - Test: returns true after store
   - Test: returns false after delete
   - Test: returns false for non-existent

5. TestCredentialManager_NotFound:
   - Retrieve non-existent credentials
   - Verify: appropriate error returned

**Restrictions:**
- Tests must run on Windows only (build tag: //go:build windows)
- Clean up all test credentials after each test
- Use unique repository names per test (avoid conflicts)
- May require admin setup for CI (document if needed)

**Success Criteria:**
- All tests pass on Windows
- Credentials stored and retrieved correctly
- Cleanup prevents test pollution
- Target name format correct

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 15: Unit Tests for MetadataCache
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/metadata_test.go`

**Description:** Unit tests for metadata caching with TTL and invalidation.

**Dependencies:** Task 5

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are writing unit tests for MetadataCache in mvnenv-win.

**Role:** Senior Go developer with caching expertise.

**Task:** Create internal/nexus/metadata_test.go:

1. TestMetadataCache_SaveAndLoad:
   - Create CachedMetadata
   - Save with Save()
   - Load with Load()
   - Verify: data matches

2. TestMetadataCache_IsFresh_TTL:
   - Save metadata
   - Test: IsFresh() returns true immediately
   - Test: IsFresh() returns false after TTL expires (use short TTL like 1ms)

3. TestMetadataCache_IsFresh_NoCache:
   - Test: IsFresh() returns false when cache doesn't exist

4. TestMetadataCache_Invalidate:
   - Save metadata
   - Invalidate()
   - Verify: cache file deleted
   - Verify: IsFresh() returns false

5. TestMetadataCache_CorruptedCache:
   - Write invalid JSON to cache file
   - Test: Load() returns error
   - Test: IsFresh() handles gracefully

6. TestMetadataCache_ConcurrentAccess:
   - Spawn multiple goroutines
   - Concurrent Save() and Load() calls
   - Verify: no race conditions

**Restrictions:**
- Use temp directory for test cache files
- Clean up cache files after tests
- Test with short TTL for fast tests
- Use t.Parallel() where appropriate

**Success Criteria:**
- All tests pass
- TTL expiration works correctly
- Corrupted cache handled gracefully
- No race conditions under concurrent access

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 16: Integration Tests
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `test/integration/nexus_integration_test.go`

**Description:** End-to-end integration tests with mock Nexus server or real Apache archive.

**Dependencies:** Tasks 8, 9, 10

**Estimated Effort:** 4 hours

**_Prompt:**
```
You are writing integration tests for Nexus repository integration in mvnenv-win.

**Role:** Senior Go developer with integration testing expertise.

**Task:** Create test/integration/nexus_integration_test.go:

1. TestIntegration_AddAndListRepositories:
   - Start with clean config
   - Add Nexus repository via CLI command
   - List repositories via CLI command
   - Verify: repository appears in list

2. TestIntegration_ConfigureAuthentication:
   - Add repository
   - Configure authentication via CLI
   - Verify: credentials stored in Windows Credential Manager
   - List repositories: verify auth status shown

3. TestIntegration_DiscoverVersions_ApacheArchive:
   - Use real Apache archive (or mock server)
   - Call DiscoverVersions
   - Verify: returns list of Maven versions
   - Verify: versions are valid semantic versions

4. TestIntegration_DownloadAndInstall:
   - Add test repository (mock server with small archive)
   - Install Maven version via VersionInstaller
   - Verify: distribution downloaded
   - Verify: checksum verified
   - Verify: version installed correctly

5. TestIntegration_Fallback:
   - Configure: high-priority repo that fails, low-priority Apache
   - Install Maven version
   - Verify: falls back to Apache archive
   - Verify: installation succeeds

6. TestIntegration_CacheRefresh:
   - Call install -l (list available)
   - Verify: cache created
   - Call install -l again
   - Verify: cache used (no network query)
   - Call update command
   - Verify: cache refreshed

7. TestIntegration_RemoveRepository:
   - Add repository with credentials
   - Remove repository via CLI
   - Verify: repository gone from config
   - Verify: credentials removed from Windows Credential Manager

**Restrictions:**
- Use temporary directories for test installations
- Clean up test data after each test
- Real Apache archive access acceptable (read-only)
- Mock Nexus server for write operations
- Tests should be skippable in offline mode

**Success Criteria:**
- Full workflow tested end-to-end
- Real network interactions work (Apache archive)
- Mock server simulates Nexus correctly
- Cleanup prevents test pollution

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

### Documentation and Finalization

#### Task 17: Create Package Documentation
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Create:**
- `internal/nexus/doc.go`

**Description:** Package-level documentation with overview and usage examples.

**Dependencies:** All implementation tasks completed

**Estimated Effort:** 1 hour

**_Prompt:**
```
You are writing package documentation for mvnenv-win's Nexus integration.

**Role:** Technical writer with Go documentation expertise.

**Task:** Create internal/nexus/doc.go with:

1. Package overview:
   - Purpose: Maven repository integration (Nexus, Apache archive)
   - Key capabilities: discovery, download, authentication, caching
   - Architecture: Repository pattern with pluggable implementations

2. Usage examples:
   - Creating RepositoryManager
   - Adding Nexus repository with authentication
   - Discovering versions
   - Downloading distributions
   - Configuring repositories

3. Code example:
```go
// Example: Configure and use repository manager
manager, err := nexus.NewRepositoryManager(configPath)
if err != nil {
    log.Fatal(err)
}

// Add Nexus repository
err = manager.AddRepository("corporate", "https://nexus.company.com/repo", 1, nexus.RepositoryTypeNexus)

// Configure authentication
err = manager.SetAuthentication("corporate", nexus.AuthTypeBasic, map[string]string{
    "username": "user",
    "password": "pass",
})

// Discover available versions
versions, err := manager.DiscoverVersions(context.Background(), false)

// Download specific version
path, err := manager.DownloadVersion(context.Background(), "3.9.4", "/tmp", nil)
```

4. Security considerations:
   - Credential storage in Windows Credential Manager
   - TLS certificate validation
   - Environment variable substitution

**Restrictions:**
- Follow godoc conventions
- Keep examples concise and runnable
- Include links to key types and functions

**Success Criteria:**
- Package documentation renders correctly in godoc
- Examples are accurate and helpful
- Security considerations clearly stated

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

#### Task 18: Final Review and Error Message Audit
**Status:** [ ] Not Started | [ ] In Progress | [ ] Completed

**Files to Review:**
- All implementation files in `internal/nexus/`
- All CLI command files referencing repositories

**Description:** Comprehensive review ensuring all error messages are clear, actionable, and match requirements.

**Dependencies:** All previous tasks completed

**Estimated Effort:** 2 hours

**_Prompt:**
```
You are conducting a final review of the Nexus integration implementation for mvnenv-win.

**Role:** Senior Go developer and technical reviewer.

**Task:** Review all implementation for completeness and quality:

1. Error Message Audit:
   - Review all error messages in internal/nexus/
   - Verify: errors include repository name context
   - Verify: errors are actionable (suggest next steps)
   - Examples:
     * "authentication failed for repository 'corporate-nexus': 401 Unauthorized. Check credentials with 'mvnenv repo auth corporate-nexus'"
     * "version '3.9.999' not found in any repository. List available versions with 'mvnenv install -l'"
     * "repository 'myrepo' already exists. Use 'mvnenv repo remove myrepo' to remove it first"

2. Requirements Verification:
   - Check requirements.md against implementation
   - Verify all acceptance criteria met
   - Document any deviations or limitations

3. Code Quality:
   - Run go vet on all nexus package files
   - Run staticcheck
   - Run golangci-lint
   - Fix any issues found

4. Test Coverage:
   - Run: go test -cover ./internal/nexus/...
   - Verify: coverage >90%
   - Add tests for any uncovered critical paths

5. Documentation:
   - Verify all exported functions have godoc comments
   - Check that package doc.go is complete
   - Ensure README or usage docs updated

6. Integration Points:
   - Verify core-version-management integration works
   - Verify CLI commands work end-to-end
   - Test manual workflow: add repo, configure auth, install version

**Success Criteria:**
- All requirements acceptance criteria met
- Error messages clear and actionable
- Code quality checks pass (vet, staticcheck, lint)
- Test coverage >90%
- Documentation complete

**Instructions:**
- When starting this task, mark it as "In Progress": [x] In Progress
- When completed, mark it as "Completed": [x] Completed
```

---

## Summary

**Total Tasks:** 18
**Estimated Total Effort:** 48 hours

**Task Dependencies Flow:**
```
Foundation (Tasks 1-3)
    ↓
Core Components (Tasks 4-8)
    ↓
Integration (Tasks 9-10)
    ↓
Testing (Tasks 11-16)
    ↓
Documentation (Tasks 17-18)
```

**Critical Path:**
Task 1 → Task 2 → Task 6 → Task 8 → Task 9 → Task 16 → Task 18

**Parallel Work Opportunities:**
- Tasks 4, 5 can be done in parallel with Task 2
- Task 7 can be done in parallel with Task 6
- Tasks 11, 12, 14, 15 (unit tests) can be done in parallel after respective implementations
- Tasks 9, 10 (integrations) can be done in parallel
