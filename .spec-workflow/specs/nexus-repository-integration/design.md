# Design Document

## Architecture Overview

The Nexus Repository Integration feature implements a flexible repository abstraction layer that supports multiple Maven distribution sources (Nexus Repository Manager, Apache Maven archives) with authentication, priority-based selection, and automatic fallback. The architecture uses the Repository pattern to abstract different source types and provides a unified interface for version discovery and distribution downloads.

### Core Design Principles

1. **Repository Abstraction**: Common `Repository` interface allows pluggable implementations for Nexus, Apache archives, and future sources
2. **Priority-Based Selection**: Repositories queried in priority order with automatic fallback on failure
3. **Secure by Default**: TLS certificate validation mandatory, credentials stored in Windows Credential Manager
4. **Graceful Degradation**: Single repository failure doesn't prevent using other repositories
5. **Efficient Caching**: Metadata cached with TTL to minimize network queries

## Component Architecture

```
internal/nexus/
├── repository.go         # Repository interface and types
├── manager.go           # RepositoryManager orchestrates multiple repositories
├── nexus_client.go      # NexusClient implements Repository for Nexus
├── apache_client.go     # ApacheArchiveClient implements Repository for Apache
├── credentials.go       # CredentialManager for Windows Credential Manager
├── metadata.go          # MetadataCache for version caching
├── config.go            # RepositoryConfig for repositories.yaml
└── http.go              # HTTPClient wrapper with retry and timeout logic
```

### Component Relationships

```
CLI Commands (repo add/list/remove)
         ↓
  RepositoryManager ←─────────────┐
         ↓                         │
    Repository (interface)         │
      ↙        ↘                  │
NexusClient  ApacheArchiveClient  │
      ↓            ↓              │
  HTTPClient   HTTPClient         │
      ↓            ↓              │
  CredentialManager               │
         ↓                        │
  MetadataCache ──────────────────┘
```

## Detailed Component Design

### 1. Repository Interface

The core abstraction for all Maven distribution sources.

```go
// Repository represents a source for Maven distributions
type Repository interface {
    // Name returns the repository identifier
    Name() string

    // Priority returns the repository priority (lower = higher priority)
    Priority() int

    // IsEnabled returns whether the repository is enabled
    IsEnabled() bool

    // DiscoverVersions queries the repository for available Maven versions
    // Returns list of versions or error if repository unavailable
    DiscoverVersions(ctx context.Context) ([]string, error)

    // DownloadDistribution downloads Maven distribution for specified version
    // Returns path to downloaded file or error
    DownloadDistribution(ctx context.Context, version string, destDir string, progress ProgressCallback) (string, error)

    // VerifyChecksum downloads and verifies checksum for distribution
    // Returns true if checksum valid, false otherwise
    VerifyChecksum(ctx context.Context, archivePath string, version string) (bool, error)

    // HasVersion checks if repository provides specific version (uses cache if available)
    HasVersion(version string) (bool, error)
}

// ProgressCallback is called during downloads to report progress
type ProgressCallback func(downloaded int64, total int64)

// RepositoryType identifies the type of repository
type RepositoryType string

const (
    RepositoryTypeNexus  RepositoryType = "nexus"
    RepositoryTypeApache RepositoryType = "apache"
)
```

### 2. RepositoryManager

Orchestrates multiple repositories with priority-based selection and fallback.

```go
// RepositoryManager manages multiple Maven distribution repositories
type RepositoryManager struct {
    repositories []Repository
    config       *RepositoryConfig
    cache        *MetadataCache
    credentials  *CredentialManager
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(configPath string) (*RepositoryManager, error)

// AddRepository adds a new repository to configuration
func (m *RepositoryManager) AddRepository(name string, url string, priority int, repoType RepositoryType) error

// RemoveRepository removes a repository from configuration
func (m *RepositoryManager) RemoveRepository(name string) error

// ListRepositories returns all configured repositories
func (m *RepositoryManager) ListRepositories() []RepositoryInfo

// SetAuthentication configures authentication for a repository
func (m *RepositoryManager) SetAuthentication(repoName string, authType AuthType, credentials interface{}) error

// DiscoverVersions queries all repositories and returns combined version list
// Queries repositories in priority order, deduplicates results
func (m *RepositoryManager) DiscoverVersions(ctx context.Context, forceRefresh bool) ([]VersionInfo, error)

// DownloadVersion downloads Maven distribution from appropriate repository
// Uses cached metadata to identify which repository has the version
// Falls back to other repositories if primary fails
func (m *RepositoryManager) DownloadVersion(ctx context.Context, version string, destDir string, progress ProgressCallback) (string, error)

// GetRepository returns repository by name for direct access
func (m *RepositoryManager) GetRepository(name string) (Repository, error)

// RepositoryInfo contains display information about a repository
type RepositoryInfo struct {
    Name        string
    URL         string
    Type        RepositoryType
    Priority    int
    Enabled     bool
    HasAuth     bool
    LastChecked time.Time
    Available   bool
}

// VersionInfo contains version and its repository source
type VersionInfo struct {
    Version    string
    Source     string // repository name
    SourceType RepositoryType
}
```

**Key Algorithms:**

**DiscoverVersions Algorithm:**
```
1. Check if cache is fresh (< TTL) and forceRefresh is false
   - If yes: return cached versions
2. Sort repositories by priority (ascending)
3. For each enabled repository in priority order:
   a. Query repository.DiscoverVersions() with timeout
   b. If success: add versions to result set with source info
   c. If failure: log warning, continue to next repository
4. If no repositories succeeded: return error
5. Deduplicate versions (keep first occurrence = highest priority source)
6. Update cache with results and timestamp
7. Return combined version list
```

**DownloadVersion Algorithm:**
```
1. Query cache to find which repository provides version
2. If found: attempt download from that repository
   a. If success: return path
   b. If failure: log error, continue to fallback
3. Query all repositories for version (in priority order)
4. For each repository that has version:
   a. Attempt download with progress callback
   b. If success: verify checksum if available, return path
   c. If failure: log error, try next repository
5. If all repositories failed: return error "version not found"
```

### 3. NexusClient

Implements Repository interface for Nexus Repository Manager.

```go
// NexusClient implements Repository for Nexus Repository Manager 2.x/3.x
type NexusClient struct {
    name         string
    baseURL      string
    priority     int
    enabled      bool
    authType     AuthType
    httpClient   *HTTPClient
    credentials  *CredentialManager
    tlsConfig    *TLSConfig
}

// NewNexusClient creates a new Nexus repository client
func NewNexusClient(name string, baseURL string, priority int, credentials *CredentialManager, tlsConfig *TLSConfig) *NexusClient

// DiscoverVersions queries Nexus for available Maven versions
// Requests: {baseURL}/org/apache/maven/apache-maven/maven-metadata.xml
// Parses: <metadata><versioning><versions><version>3.9.4</version>...
func (c *NexusClient) DiscoverVersions(ctx context.Context) ([]string, error)

// DownloadDistribution downloads Maven distribution from Nexus
// URL: {baseURL}/org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip
func (c *NexusClient) DownloadDistribution(ctx context.Context, version string, destDir string, progress ProgressCallback) (string, error)

// VerifyChecksum downloads .sha256 file and verifies archive
// URL: {baseURL}/org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip.sha256
func (c *NexusClient) VerifyChecksum(ctx context.Context, archivePath string, version string) (bool, error)

// AuthType defines authentication method
type AuthType string

const (
    AuthTypeNone   AuthType = "none"
    AuthTypeBasic  AuthType = "basic"
    AuthTypeToken  AuthType = "token"
)

// TLSConfig contains TLS/SSL settings
type TLSConfig struct {
    Insecure       bool   // Skip certificate verification (must be explicit)
    CAFile         string // Path to custom CA certificate
    ClientCertFile string // Path to client certificate (future)
    ClientKeyFile  string // Path to client key (future)
}
```

**Key Operations:**

**DiscoverVersions Implementation:**
```
1. Construct metadata URL: {baseURL}/org/apache/maven/apache-maven/maven-metadata.xml
2. Create HTTP GET request with context
3. Add authentication headers if configured
4. Execute request with timeout
5. If 401/403: return error "authentication failed"
6. If 404: return error "metadata not found"
7. Parse XML response:
   <metadata>
     <versioning>
       <versions>
         <version>3.6.3</version>
         <version>3.8.6</version>
         ...
       </versions>
     </versioning>
   </metadata>
8. Extract all <version> elements
9. Return version list
```

**DownloadDistribution Implementation:**
```
1. Construct artifact URL: {baseURL}/org/apache/maven/apache-maven/{version}/apache-maven-{version}-bin.zip
2. Create destination file: {destDir}/apache-maven-{version}-bin.zip
3. Create HTTP GET request with context
4. Add authentication headers if configured
5. Execute request and get response body reader
6. If 401/403: return error "authentication failed"
7. If 404: return error "version not found"
8. Stream response to file with progress callback:
   - Read chunk from response
   - Write chunk to file
   - Call progress(bytesWritten, totalSize)
9. Close file and response
10. Return file path
```

### 4. ApacheArchiveClient

Implements Repository interface for Apache Maven official archives.

```go
// ApacheArchiveClient implements Repository for Apache Maven archive
type ApacheArchiveClient struct {
    name       string
    baseURL    string // https://archive.apache.org/dist/maven/maven-3/
    priority   int
    enabled    bool
    httpClient *HTTPClient
}

// NewApacheArchiveClient creates Apache archive repository client
func NewApacheArchiveClient(name string, priority int) *ApacheArchiveClient

// DiscoverVersions scrapes Apache archive directory listing
// Parses HTML: <a href="3.9.4/">3.9.4/</a>
func (c *ApacheArchiveClient) DiscoverVersions(ctx context.Context) ([]string, error)

// DownloadDistribution downloads from Apache archive
// URL: {baseURL}/{version}/binaries/apache-maven-{version}-bin.zip
func (c *ApacheArchiveClient) DownloadDistribution(ctx context.Context, version string, destDir string, progress ProgressCallback) (string, error)

// VerifyChecksum downloads .sha512 file and verifies archive
// URL: {baseURL}/{version}/binaries/apache-maven-{version}-bin.zip.sha512
func (c *ApacheArchiveClient) VerifyChecksum(ctx context.Context, archivePath string, version string) (bool, error)
```

**Key Operations:**

**DiscoverVersions Implementation (Directory Scraping):**
```
1. HTTP GET: https://archive.apache.org/dist/maven/maven-3/
2. Parse HTML response for directory links:
   - Match pattern: <a href="(\d+\.\d+\.\d+)/">
   - Extract version numbers
3. Filter valid semantic versions
4. Return version list
```

### 5. CredentialManager

Manages secure credential storage using Windows Credential Manager.

```go
// CredentialManager handles secure credential storage
type CredentialManager struct {
    targetPrefix string // "mvnenv:repo:"
}

// NewCredentialManager creates a credential manager
func NewCredentialManager() *CredentialManager

// StoreCredentials stores credentials in Windows Credential Manager
// Target: mvnenv:repo:{repoName}
func (m *CredentialManager) StoreCredentials(repoName string, username string, password string) error

// RetrieveCredentials retrieves credentials from Windows Credential Manager
// Returns username, password, error
func (m *CredentialManager) RetrieveCredentials(repoName string) (string, string, error)

// DeleteCredentials removes credentials from Windows Credential Manager
func (m *CredentialManager) DeleteCredentials(repoName string) error

// HasCredentials checks if credentials exist for repository
func (m *CredentialManager) HasCredentials(repoName string) bool
```

**Implementation Details:**
- Use `github.com/danieljoos/wincred` or syscall to Windows Credential Manager APIs
- Target name format: `mvnenv:repo:{repoName}` for namespacing
- Store username as CredentialAttribute, password as CredentialBlob (encrypted by OS)
- Handle errors gracefully: if Credential Manager unavailable, return error suggesting environment variables

### 6. MetadataCache

Caches version metadata to reduce network queries.

```go
// MetadataCache manages cached repository metadata
type MetadataCache struct {
    cachePath string // %USERPROFILE%\.mvnenv\cache\repo-metadata.json
    ttl       time.Duration
    mu        sync.RWMutex
}

// NewMetadataCache creates a metadata cache
func NewMetadataCache(cachePath string, ttl time.Duration) *MetadataCache

// IsFresh checks if cache is within TTL
func (c *MetadataCache) IsFresh() bool

// Load loads cached metadata from disk
func (c *MetadataCache) Load() (*CachedMetadata, error)

// Save saves metadata to disk with current timestamp
func (c *MetadataCache) Save(metadata *CachedMetadata) error

// Invalidate deletes cache file
func (c *MetadataCache) Invalidate() error

// CachedMetadata contains cached version information
type CachedMetadata struct {
    Timestamp time.Time               `json:"timestamp"`
    Versions  []CachedVersionInfo     `json:"versions"`
    Repos     map[string]RepoSnapshot `json:"repos"` // repository state snapshot
}

// CachedVersionInfo contains version with source repository
type CachedVersionInfo struct {
    Version    string         `json:"version"`
    Source     string         `json:"source"`      // repository name
    SourceType RepositoryType `json:"source_type"`
}

// RepoSnapshot captures repository state at cache time
type RepoSnapshot struct {
    URL       string    `json:"url"`
    Priority  int       `json:"priority"`
    Available bool      `json:"available"` // was repository reachable?
    Timestamp time.Time `json:"timestamp"`
}
```

**Cache Invalidation Rules:**
1. TTL expired (default 24 hours)
2. Repository configuration changed (add/remove/modify)
3. User runs `mvnenv update` command
4. Cache file corrupted (JSON parse error)

### 7. RepositoryConfig

Manages repositories.yaml configuration file.

```go
// RepositoryConfig manages repository configuration
type RepositoryConfig struct {
    configPath string
    mu         sync.RWMutex
}

// NewRepositoryConfig creates repository config manager
func NewRepositoryConfig(configPath string) *RepositoryConfig

// Load loads configuration from disk
func (c *RepositoryConfig) Load() (*Config, error)

// Save saves configuration to disk
func (c *RepositoryConfig) Save(config *Config) error

// AddRepository adds repository to configuration
func (c *RepositoryConfig) AddRepository(repo RepositoryEntry) error

// RemoveRepository removes repository from configuration
func (c *RepositoryConfig) RemoveRepository(name string) error

// Config represents repositories.yaml structure
type Config struct {
    Repositories []RepositoryEntry `yaml:"repositories"`
}

// RepositoryEntry represents a single repository configuration
type RepositoryEntry struct {
    Name     string         `yaml:"name"`
    Type     RepositoryType `yaml:"type"`
    URL      string         `yaml:"url"`
    Priority int            `yaml:"priority"`
    Enabled  bool           `yaml:"enabled"`
    Auth     *AuthConfig    `yaml:"auth,omitempty"`
    TLS      *TLSConfig     `yaml:"tls,omitempty"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
    Type     AuthType `yaml:"type"`
    Username string   `yaml:"username,omitempty"` // Can be ${ENV_VAR}
    Password string   `yaml:"password,omitempty"` // Can be ${ENV_VAR}
    Token    string   `yaml:"token,omitempty"`    // Can be ${ENV_VAR}
}
```

**Environment Variable Substitution:**
```go
// expandEnvVars replaces ${VAR_NAME} with environment variable value
func expandEnvVars(value string) string {
    re := regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)
    return re.ReplaceAllStringFunc(value, func(match string) string {
        varName := match[2 : len(match)-1] // Strip ${ and }
        return os.Getenv(varName)
    })
}
```

**Default Configuration (no repositories.yaml):**
```yaml
repositories:
  - name: "apache"
    type: "apache"
    url: "https://archive.apache.org/dist/maven/maven-3/"
    priority: 100
    enabled: true
```

### 8. HTTPClient

Wrapper for HTTP client with retry logic and timeout handling.

```go
// HTTPClient wraps HTTP client with retry and timeout logic
type HTTPClient struct {
    client      *resty.Client
    timeout     time.Duration
    retryCount  int
    retryWait   time.Duration
}

// NewHTTPClient creates HTTP client with configuration
func NewHTTPClient(timeout time.Duration, tlsConfig *TLSConfig) *HTTPClient

// Get performs HTTP GET with retry logic
func (c *HTTPClient) Get(ctx context.Context, url string, auth *AuthConfig) (*resty.Response, error)

// Download downloads file with progress callback
func (c *HTTPClient) Download(ctx context.Context, url string, destPath string, auth *AuthConfig, progress ProgressCallback) error

// SetTLSConfig updates TLS configuration
func (c *HTTPClient) SetTLSConfig(tlsConfig *TLSConfig)
```

**Retry Logic:**
- Retry on: connection errors, timeouts, 5xx server errors
- Don't retry on: 4xx client errors (except 429 Too Many Requests)
- Exponential backoff: 1s, 2s, 4s
- Max retries: 3

## Integration Points

### Integration with core-version-management

The `VersionInstaller` component (from core-version-management spec) will use `RepositoryManager`:

```go
// In internal/version/installer.go
type VersionInstaller struct {
    repoManager *nexus.RepositoryManager
    // ... other fields
}

func (i *VersionInstaller) InstallVersion(version string) error {
    // 1. Download distribution
    archivePath, err := i.repoManager.DownloadVersion(ctx, version, tempDir, progressCallback)
    if err != nil {
        return fmt.Errorf("download failed: %w", err)
    }

    // 2. Extract and install
    // ... extraction logic ...
}
```

### Integration with cli-commands

Repository management commands (from cli-commands spec) will use `RepositoryManager`:

```go
// In cmd/mvnenv/cmd/repo_add.go
func runRepoAdd(cmd *cobra.Command, args []string) error {
    repoManager, err := nexus.NewRepositoryManager(configPath)
    // ...
    return repoManager.AddRepository(name, url, priority, repoType)
}
```

## Error Handling

### Error Types

```go
// Error types for repository operations
var (
    ErrRepositoryNotFound     = errors.New("repository not found")
    ErrRepositoryExists       = errors.New("repository already exists")
    ErrAuthenticationFailed   = errors.New("authentication failed")
    ErrVersionNotFound        = errors.New("version not found in any repository")
    ErrNoRepositoriesAvailable = errors.New("no repositories available")
    ErrInvalidConfiguration   = errors.New("invalid repository configuration")
    ErrCertificateInvalid     = errors.New("TLS certificate validation failed")
)

// RepositoryError wraps errors with repository context
type RepositoryError struct {
    Repo      string
    Operation string
    Cause     error
}

func (e *RepositoryError) Error() string {
    return fmt.Sprintf("repository '%s' %s: %v", e.Repo, e.Operation, e.Cause)
}
```

### Error Handling Strategy

1. **Network Errors**: Retry with exponential backoff, fall back to next repository
2. **Authentication Errors**: Return immediately with clear error (don't retry other repos with same creds)
3. **Not Found Errors**: Try next repository in priority order
4. **Certificate Errors**: Return immediately with certificate details
5. **Configuration Errors**: Return immediately (don't attempt operation)

## Testing Strategy

### Unit Tests

```go
// Test files
internal/nexus/manager_test.go        // RepositoryManager logic
internal/nexus/nexus_client_test.go   // NexusClient operations
internal/nexus/apache_client_test.go  // ApacheArchiveClient operations
internal/nexus/credentials_test.go    // Credential management
internal/nexus/metadata_test.go       // Cache operations
internal/nexus/config_test.go         // Configuration management
```

**Key Test Cases:**
- Repository priority ordering (multiple repos with different priorities)
- Fallback logic (primary repo fails, secondary succeeds)
- Authentication header injection (basic and token auth)
- Environment variable substitution in config
- Cache freshness calculation and TTL expiration
- Checksum verification (valid and invalid checksums)
- TLS certificate validation (valid, invalid, custom CA)
- Concurrent repository queries

### Integration Tests

```go
// Test files
test/integration/nexus_integration_test.go
```

**Key Test Scenarios:**
1. Add/remove repositories via RepositoryManager
2. Discover versions from real Nexus instance (if available) or mock server
3. Download distribution with progress tracking
4. Fallback: primary Nexus fails, Apache archive succeeds
5. Authentication: valid credentials succeed, invalid fail with clear error
6. Cache: first query hits network, second query uses cache
7. Configuration file persistence across manager restarts

### Mock Server for Testing

```go
// test/mocks/nexus_mock.go
type MockNexusServer struct {
    *httptest.Server
    metadata    string // XML metadata response
    versions    map[string][]byte // version -> archive content
    authEnabled bool
    authToken   string
}
```

## Security Considerations

### Credential Security

1. **Windows Credential Manager**: Primary storage, OS-managed encryption
2. **Environment Variables**: Fallback for automation, clearly documented as less secure
3. **No Plaintext in Config**: Configuration file never contains plaintext passwords
4. **Audit Logging**: All authentication attempts logged (success/failure)

### TLS Security

1. **Certificate Validation**: Mandatory by default, explicit flag to disable
2. **Custom CA Support**: Load custom CA certificates from file
3. **Warning Logging**: Log warning on every insecure connection
4. **Certificate Chain Validation**: Full chain including intermediates

### Input Validation

1. **Repository Name**: Alphanumeric, hyphens, underscores only (prevent path traversal)
2. **URL Validation**: Must be valid HTTP/HTTPS URL with hostname
3. **Version String**: Semantic version format validation
4. **Path Sanitization**: All file paths sanitized to prevent directory traversal

## Performance Optimization

### Parallel Repository Queries

```go
// Query multiple repositories in parallel
func (m *RepositoryManager) DiscoverVersions(ctx context.Context, forceRefresh bool) ([]VersionInfo, error) {
    var wg sync.WaitGroup
    results := make(chan []VersionInfo, len(m.repositories))

    for _, repo := range m.repositories {
        if !repo.IsEnabled() {
            continue
        }

        wg.Add(1)
        go func(r Repository) {
            defer wg.Done()
            versions, err := r.DiscoverVersions(ctx)
            if err != nil {
                log.Warnf("repository %s failed: %v", r.Name(), err)
                return
            }
            // Convert to VersionInfo and send to channel
        }(repo)
    }

    wg.Wait()
    close(results)

    // Combine and deduplicate results
}
```

### Connection Pooling

- Reuse HTTP client across requests to same repository
- Configure connection pool size: 10 connections per repository
- Keep-alive timeout: 90 seconds

### Download Resume

- Support HTTP Range requests for partial download resume (future enhancement)
- Current implementation: restart download on failure

## Configuration Examples

### repositories.yaml with Multiple Sources

```yaml
repositories:
  - name: "corporate-nexus"
    type: "nexus"
    url: "https://nexus.company.com/repository/maven-releases"
    priority: 1
    enabled: true
    auth:
      type: "basic"
      username: "${NEXUS_USER}"
      password: "${NEXUS_PASS}"
    tls:
      insecure: false
      ca_file: "C:\\certs\\company-ca.crt"

  - name: "backup-nexus"
    type: "nexus"
    url: "https://nexus-backup.company.com/repository/maven"
    priority: 2
    enabled: true
    auth:
      type: "token"
      token: "${NEXUS_TOKEN}"

  - name: "apache"
    type: "apache"
    url: "https://archive.apache.org/dist/maven/maven-3/"
    priority: 100
    enabled: true
```

### Minimal Configuration (Apache Only)

```yaml
repositories:
  - name: "apache"
    type: "apache"
    url: "https://archive.apache.org/dist/maven/maven-3/"
    priority: 1
    enabled: true
```

## Future Enhancements

### Phase 1 (v1.1.0)
- Maven Central repository support (in addition to Apache archives)
- Proxy server configuration for HTTP requests
- Download resume support (HTTP Range requests)

### Phase 2 (v1.2.0)
- Repository health monitoring and circuit breaker pattern
- Automatic repository discovery from Maven settings.xml
- Repository mirroring configuration
- Webhook notifications for new version availability

### Phase 3 (v2.0.0)
- GraphQL API for Nexus 3.x (in addition to REST)
- Artifact signing verification (GPG signatures)
- Bandwidth throttling for downloads
- Multi-threaded chunk downloads for large files
