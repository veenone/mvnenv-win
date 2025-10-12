# Tasks Document

## Implementation Tasks for core-version-management Spec

- [x] 1. Create Maven version structures in pkg/maven
  - Files: `pkg/maven/version.go`
  - Define Version struct (Major, Minor, Patch, Qualifier)
  - Implement ParseVersion() function
  - Implement Version.Compare() method for semantic versioning
  - Purpose: Provide public version parsing and comparison utilities
  - _Leverage: Go standard library (strings, strconv)_
  - _Requirements: Req 9 (Get Latest Version - comparison logic)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in semantic versioning and data structures | Task: Create pkg/maven/version.go implementing Version struct and semantic version parsing/comparison following design.md Component 5, supporting Maven version formats (X.Y.Z and X.Y.Z-qualifier) | Restrictions: Must handle qualifiers (alpha, beta, RC), use semantic versioning rules, export all types and functions, no external dependencies beyond stdlib | Success: ParseVersion("3.9.4") works, Compare() correctly orders versions, qualifiers handled properly, unit tests pass | Instructions: After completing, edit tasks.md and change this task from [ ] to [x]_

- [x] 2. Create Maven path utilities in pkg/maven
  - Files: `pkg/maven/paths.go`
  - Implement GetMavenBinaryPath() for bin/mvn.cmd location
  - Implement GetMavenHome() for Maven installation root
  - Implement ValidateMavenInstallation() to verify valid Maven directory
  - Purpose: Provide Maven-specific path resolution utilities
  - _Leverage: filepath package for Windows paths_
  - _Requirements: Req 1 (Installation verification), Req 3 (List installed)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in file path handling and Windows conventions | Task: Create pkg/maven/paths.go with Maven installation path utilities following design.md, using Windows-style paths and validating Maven directory structure | Restrictions: Must use filepath package, handle Windows backslashes, check for bin/mvn.cmd existence, export all functions | Success: GetMavenBinaryPath() returns correct Windows path, ValidateMavenInstallation() detects valid/invalid installations, paths use backslashes | Instructions: After completing, edit tasks.md and change this task from [ ] to [x]_

- [x] 3. Create download manager in internal/download
  - Files: `internal/download/downloader.go`, `internal/download/checksum.go`
  - Implement Downloader interface with HTTP download and progress tracking
  - Implement ChecksumVerifier for SHA-256 verification
  - Add retry logic and timeout handling
  - Purpose: Provide reliable download with integrity verification (shared with nexus spec)
  - _Leverage: net/http, crypto/sha256, io packages_
  - _Requirements: Req 1 (Install with checksum verification)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in HTTP operations and file integrity | Task: Create internal/download package with Downloader and ChecksumVerifier following design.md, supporting HTTP downloads with SHA-256 verification, progress tracking, retries, and timeouts | Restrictions: Must verify all downloads with SHA-256, handle network errors gracefully, provide progress callbacks, no external HTTP libraries (use stdlib) | Success: Downloads work with progress, checksums verified correctly, failed downloads cleaned up, retries on transient errors | Instructions: After completing, edit tasks.md and change this task from [ ] to [x]_

- [x] 4. Create config package additions
  - Files: `internal/config/config.go`, `internal/config/paths.go`
  - Add GlobalVersion field to Config struct
  - Implement GetGlobalVersion() and SetGlobalVersion() methods
  - Add path helpers: VersionsDir(), CacheDir(), ConfigDir()
  - Purpose: Extend config package for version management needs
  - _Leverage: Existing config package from steering docs, gopkg.in/yaml.v3_
  - _Requirements: Req 6 (Set Global Version)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in configuration management and YAML | Task: Extend internal/config package with global version support and path helpers following design.md, adding GlobalVersion field and path resolution methods | Restrictions: Must preserve existing config fields, use YAML for persistence, handle missing config file gracefully, use %USERPROFILE%/.mvnenv paths | Success: Global version can be get/set, config persists to YAML, path helpers return correct Windows paths, config loading doesn't break existing code | Instructions: After completing, edit tasks.md and change this task from [ ] to [x]_

- [x] 5. Create VersionResolver for version resolution
  - Files: `internal/version/resolver.go`
  - Implement VersionResolver struct with ResolveVersion() method
  - Implement shell > local > global hierarchy
  - Add GetShellVersion(), GetLocalVersion(), GetGlobalVersion() methods
  - Add IsVersionInstalled() verification
  - Purpose: Determine which Maven version should be active
  - _Leverage: internal/config, os package for environment variables, filepath for .maven-version search_
  - _Requirements: Req 5 (Version Resolution)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in file system operations and environment variables | Task: Create internal/version/resolver.go implementing VersionResolver following design.md Component 2, with shell > local > global hierarchy, searching parent directories for .maven-version files | Restrictions: Must follow exact hierarchy (shell first, then local, then global), search up directory tree for .maven-version, verify version is installed before returning, handle missing config gracefully | Success: ResolveVersion() returns correct version based on hierarchy, shell env takes precedence, .maven-version found in parents, global from config works, uninstalled versions return error | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [x] 6. Create VersionInstaller for install/uninstall
  - Files: `internal/version/installer.go`
  - Implement VersionInstaller struct with Install() and Uninstall() methods
  - Add atomic installation with temp directory and rollback
  - Implement VerifyInstallation() for post-install checks
  - Add IsInstalled() check
  - Purpose: Handle Maven version lifecycle (install/uninstall)
  - _Leverage: internal/download, pkg/maven paths, archive/zip for extraction_
  - _Requirements: Req 1 (Install), Req 2 (Uninstall)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in file operations and atomic transactions | Task: Create internal/version/installer.go implementing VersionInstaller following design.md Component 3 and installation flow diagram, with atomic install (extract to temp, verify, move), checksum verification, and complete rollback on failure | Restrictions: Must be atomic (all-or-nothing), use temp directory during install, rollback completely on any error, verify Maven binary exists after extraction, handle Windows file locking | Success: Installation atomic with rollback, checksums verified, temp files cleaned up on failure, uninstall removes directory cleanly, already-installed detected | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [x] 7. Create VersionLister for listing versions
  - Files: `internal/version/lister.go`
  - Implement VersionLister struct with ListInstalled() and ListAvailable() methods
  - Add GetLatest() for finding latest version with optional prefix
  - Implement version sorting using pkg/maven comparison
  - Purpose: List installed and available Maven versions
  - _Leverage: pkg/maven version comparison, internal/version/cache_
  - _Requirements: Req 3 (List Installed), Req 4 (List Available), Req 9 (Get Latest)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer specializing in directory scanning and sorting algorithms | Task: Create internal/version/lister.go implementing VersionLister following design.md Component 4, scanning versions directory for installed versions, sorting by semantic version, finding latest with prefix support | Restrictions: Must validate each version directory (check for Maven binary), sort descending (newest first), handle empty versions directory, filter by prefix for GetLatest | Success: ListInstalled() scans directory correctly, versions sorted newest first, invalid directories skipped, GetLatest("3.8") returns newest 3.8.x version | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [x] 8. Create VersionCache for caching available versions
  - Files: `internal/version/cache.go`
  - Implement VersionCache struct with Get() and Set() methods
  - Add timestamp-based TTL (24 hours default)
  - Store cache in JSON format
  - Purpose: Cache available versions to avoid repeated repository queries
  - _Leverage: encoding/json, time package_
  - _Requirements: Req 4 (List Available with caching), Req 10 (Update Cache)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in caching and time-based expiration | Task: Create internal/version/cache.go implementing VersionCache following design.md data models, storing available versions with timestamp in JSON, checking TTL (24h default) before returning cached data | Restrictions: Must use JSON format, include timestamp in cache file, check staleness before returning, create cache directory if missing, handle corrupt cache gracefully | Success: Cache stores and retrieves versions, TTL respected (stale cache returns nil), JSON format valid, cache file in correct location | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 9. Create VersionManager orchestrator
  - Files: `internal/version/manager.go`
  - Implement VersionManager struct composing Resolver, Installer, and Lister
  - Add high-level methods: Install(), Uninstall(), GetCurrentVersion(), SetGlobalVersion(), SetLocalVersion(), etc.
  - Purpose: Provide unified API for all version management operations
  - _Leverage: All internal/version components (resolver, installer, lister, cache)_
  - _Requirements: All requirements (orchestration layer)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Software Architect specializing in API design and component orchestration | Task: Create internal/version/manager.go implementing VersionManager following design.md Component 1, composing Resolver, Installer, Lister, and Cache into unified API for CLI commands | Restrictions: Must delegate to specialized components (don't duplicate logic), provide clean public API, handle errors from components gracefully, initialize all components in NewVersionManager | Success: VersionManager provides all required methods, delegates correctly to components, CLI commands can use simple API, error handling consistent | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [x] 10. Integrate VersionManager with CLI commands
  - Files: `cmd/mvnenv/version.go`, `cmd/mvnenv/versions.go`, `cmd/mvnenv/global.go`, `cmd/mvnenv/local.go`, `cmd/mvnenv/shell.go`, `cmd/mvnenv/install.go`, `cmd/mvnenv/uninstall.go`, `cmd/mvnenv/latest.go`, `cmd/mvnenv/update.go`
  - Replace placeholder logic in CLI commands with VersionManager calls
  - Initialize VersionManager in main.go for command access
  - Handle errors and display user-friendly messages
  - Purpose: Connect version management business logic to CLI
  - _Leverage: cli-commands spec implementation, internal/version/manager_
  - _Requirements: All requirements (integration)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in CLI integration and error handling | Task: Integrate VersionManager with CLI commands from cli-commands spec, replacing placeholder logic with real VersionManager method calls, handling errors and formatting output following pyenv-win conventions | Restrictions: Must maintain CLI command signatures, use consistent error handling, output plain text (no colors/emojis), silent on success for setters, format errors with "Error: " prefix | Success: All CLI commands work with real version management, install/uninstall functional, version resolution works, errors displayed clearly, follows pyenv-win output style | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 11. Add unit tests for pkg/maven/version.go
  - Files: `pkg/maven/version_test.go`
  - Test ParseVersion() with various formats (3.9.4, 3.9.0-beta-1, etc.)
  - Test Version.Compare() with all comparison cases
  - Test semantic versioning edge cases (qualifiers, leading zeros)
  - Purpose: Ensure version parsing and comparison correctness
  - _Leverage: Go testing package, table-driven tests_
  - _Requirements: Req 9 (version comparison for GetLatest)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in unit testing and edge case coverage | Task: Create comprehensive unit tests for pkg/maven/version.go testing ParseVersion and Compare methods with table-driven tests covering standard versions, qualifiers, and edge cases | Restrictions: Must use table-driven tests, test all version formats (X.Y.Z and X.Y.Z-qualifier), verify correct semantic ordering, test error cases (invalid formats), achieve 100% coverage | Success: All version parsing tested, comparison correctness verified, qualifiers handled properly, invalid formats return errors, 100% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 12. Add unit tests for VersionResolver
  - Files: `internal/version/resolver_test.go`
  - Test shell > local > global hierarchy
  - Test .maven-version file discovery in parent directories
  - Test version installed verification
  - Mock FileSystem and ConfigProvider
  - Purpose: Ensure version resolution logic correctness
  - _Leverage: Go testing, testify/mock for mocking_
  - _Requirements: Req 5 (Version Resolution)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer specializing in behavior testing and mocking | Task: Create unit tests for VersionResolver testing resolution hierarchy (shell > local > global), .maven-version discovery in parent dirs, and version verification using mocked filesystem and config | Restrictions: Must mock filesystem and config, test each hierarchy level independently, verify parent directory traversal for .maven-version, test version not installed error, >90% coverage | Success: Hierarchy precedence tested (shell wins), .maven-version found in parents, global from config works, uninstalled versions error, >90% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 13. Add unit tests for VersionInstaller
  - Files: `internal/version/installer_test.go`
  - Test successful installation flow
  - Test rollback on checksum failure
  - Test rollback on extraction failure
  - Test already-installed detection
  - Mock Downloader and filesystem
  - Purpose: Ensure installation atomicity and error handling
  - _Leverage: Go testing, mock filesystem for temp directory tests_
  - _Requirements: Req 1 (Install), Req 2 (Uninstall)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in testing atomic operations and rollback scenarios | Task: Create unit tests for VersionInstaller testing installation success, checksum failure rollback, extraction failure rollback, and already-installed detection using mocked downloader and filesystem | Restrictions: Must test rollback scenarios (checksum fail, extraction fail), verify temp directory cleanup, test already-installed early exit, mock all external dependencies, >90% coverage | Success: Installation flow tested end-to-end, rollbacks verified (no partial installs), already-installed detected, temp files cleaned up, >90% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 14. Add unit tests for VersionLister
  - Files: `internal/version/lister_test.go`
  - Test ListInstalled() with version sorting
  - Test GetLatest() with and without prefix
  - Test empty versions directory handling
  - Test invalid version directory exclusion
  - Purpose: Ensure listing and latest version logic correctness
  - _Leverage: Go testing, mock filesystem_
  - _Requirements: Req 3 (List Installed), Req 9 (Get Latest)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer specializing in sorting algorithms and filtering logic | Task: Create unit tests for VersionLister testing ListInstalled sorting (newest first), GetLatest with prefix filtering, empty directory handling, and invalid version exclusion using mocked filesystem | Restrictions: Must verify descending sort order (newest first), test prefix matching for GetLatest, verify invalid directories skipped, test empty directory returns empty list, >90% coverage | Success: Sorting verified (3.9.4 before 3.8.6), GetLatest("3.8") returns newest 3.8.x, invalid dirs excluded, empty dir handled, >90% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 15. Add unit tests for VersionCache
  - Files: `internal/version/cache_test.go`
  - Test cache Get() with fresh and stale data
  - Test cache Set() with timestamp
  - Test corrupt cache handling
  - Test missing cache directory creation
  - Purpose: Ensure cache reliability and TTL correctness
  - _Leverage: Go testing, time mocking_
  - _Requirements: Req 4 (caching), Req 10 (Update Cache)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Engineer with expertise in caching logic and time-based testing | Task: Create unit tests for VersionCache testing cache Get/Set with TTL checking (24h), stale cache handling, corrupt cache recovery, and directory creation using temp directories and time manipulation | Restrictions: Must test fresh cache returns data, stale cache returns nil, corrupt JSON handled gracefully, missing directory created automatically, 100% coverage | Success: Fresh cache retrieved correctly, stale (>24h) returns nil, corrupt cache doesn't crash, directory auto-created, 100% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 16. Add unit tests for VersionManager
  - Files: `internal/version/manager_test.go`
  - Test all VersionManager public methods
  - Test component orchestration
  - Test error propagation from components
  - Mock all dependencies (Resolver, Installer, Lister)
  - Purpose: Ensure VersionManager orchestration correctness
  - _Leverage: Go testing, mocks for components_
  - _Requirements: All requirements (orchestration testing)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Integration Tester specializing in component orchestration and API testing | Task: Create unit tests for VersionManager testing all public methods (Install, Uninstall, GetCurrentVersion, SetGlobal/Local, etc.) with mocked components verifying correct delegation and error handling | Restrictions: Must mock Resolver, Installer, Lister, Cache; verify method calls to components; test error propagation; test all public API methods; >90% coverage | Success: All VersionManager methods tested, delegation to components verified, errors propagated correctly, API contract validated, >90% coverage | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 17. Add integration tests for version resolution
  - Files: `internal/version/integration_test.go`
  - Test resolution with real filesystem (temp directories)
  - Test .maven-version in nested project directories
  - Test environment variable override
  - Test global config persistence
  - Purpose: Validate resolution in realistic scenarios
  - _Leverage: Go testing, temp directories_
  - _Requirements: Req 5 (Version Resolution), Req 6-8 (Set versions)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Integration Engineer with expertise in filesystem testing and environment setup | Task: Create integration tests for version resolution using real temp filesystem, testing shell env override, .maven-version in nested dirs, and global config with actual file operations | Restrictions: Must use real filesystem (temp dirs), test complete resolution flow, set environment variables for shell test, create nested directories for .maven-version test, cleanup temp files | Success: Integration tests pass with real files, shell env override works, .maven-version found in parents, global config persists, tests isolated and repeatable | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 18. Add integration tests for installation
  - Files: `internal/version/install_integration_test.go`
  - Test full installation flow with mock download
  - Test uninstallation cleanup
  - Test atomic installation rollback scenarios
  - Use real filesystem with temp directories
  - Purpose: Validate installation/uninstallation in realistic scenarios
  - _Leverage: Go testing, temp directories, mock HTTP server for downloads_
  - _Requirements: Req 1 (Install), Req 2 (Uninstall)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Integration Engineer specializing in file operations and atomic transaction testing | Task: Create integration tests for installation using real temp filesystem and mock HTTP server, testing complete install/uninstall flow, atomic operations, and rollback scenarios | Restrictions: Must use real filesystem (temp dirs), mock HTTP server for downloads, test atomicity (rollback on failure), verify no partial installs, cleanup all temp files | Success: Integration tests pass with real files, installation atomic, rollback works correctly, uninstall cleans up completely, no file leaks | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 19. Add performance benchmarks
  - Files: `internal/version/benchmark_test.go`
  - Benchmark ResolveVersion() targeting <100ms
  - Benchmark ListInstalled() targeting <100ms
  - Benchmark version comparison operations
  - Purpose: Ensure performance requirements are met
  - _Leverage: Go benchmark testing_
  - _Requirements: Performance non-functional requirements_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Performance Engineer with expertise in Go benchmarking and optimization | Task: Create performance benchmarks for VersionResolver.ResolveVersion() and VersionLister.ListInstalled() verifying <100ms targets, and version comparison <1ms | Restrictions: Must benchmark critical path operations only, measure realistic scenarios (with files present), document results, identify optimization opportunities if targets not met | Success: Benchmarks run successfully, ResolveVersion <100ms, ListInstalled <100ms, comparison <1ms, or gaps documented with optimization plan | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 20. Create end-to-end test scenarios
  - Files: `test/e2e/version_management_test.go`
  - Test complete user workflows (install -> set global -> use)
  - Test project-specific version (local) workflow
  - Test version switching scenarios
  - Use real mvnenv binary with temp environment
  - Purpose: Validate complete version management workflows
  - _Leverage: Go testing, exec package to run mvnenv CLI_
  - _Requirements: All requirements (E2E validation)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: QA Automation Engineer specializing in end-to-end testing and user workflow validation | Task: Create E2E tests running complete version management workflows (install version, set global, verify active, switch to local, verify switch) using built mvnenv binary and temp environment | Restrictions: Must test realistic user scenarios, use actual CLI binary, temp isolated environment (%USERPROFILE% override), verify version switching works end-to-end, cleanup after tests | Success: E2E tests pass with real binary, install workflow works, global/local version setting verified, version switching tested, user experience validated | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [x] 21. Create error handling documentation
  - Files: `internal/version/errors.go`
  - Define all error types (ErrVersionNotInstalled, ErrAlreadyInstalled, etc.)
  - Add error wrapping helpers for context
  - Document error messages and user actions
  - Purpose: Provide consistent error handling across version management
  - _Leverage: errors package, fmt.Errorf with %w_
  - _Requirements: Req 8 (Error Handling NFR)_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Developer with expertise in error handling patterns and user experience | Task: Create internal/version/errors.go defining all error types from design.md error handling section, providing exported error variables and helper functions for wrapping errors with context | Restrictions: Must define exported error variables for all scenarios, use errors.New for base errors, provide helpers for wrapping with context (fmt.Errorf with %w), document each error with suggested user action | Success: All error types defined, error variables exported, wrapping helpers provided, error messages clear and actionable, documentation complete | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

- [ ] 22. Final integration and code review
  - Files: All files in `internal/version/`, `pkg/maven/`
  - Review all component implementations for consistency
  - Verify error handling follows patterns
  - Ensure performance benchmarks meet targets
  - Verify >90% test coverage across all packages
  - Check CLI integration works end-to-end
  - Purpose: Ensure specification is fully implemented and ready
  - _Leverage: All previous tasks, requirements.md Success Criteria_
  - _Requirements: All requirements_
  - _Prompt: Implement the task for spec core-version-management, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Senior Go Developer and Code Reviewer with expertise in version management and code quality | Task: Perform final integration review of core version management implementation, verifying all Success Criteria from requirements.md are met, test coverage >90%, performance targets met, error handling consistent, CLI integration complete | Restrictions: Must verify all Success Criteria from requirements.md, ensure atomic installations work, test resolution hierarchy, confirm performance targets (<100ms), validate >90% coverage | Success: All Success Criteria from requirements.md verified and met, atomic installations confirmed, resolution works correctly, performance targets met, >90% coverage, CLI integration functional, ready for production use | Instructions: Before starting, edit tasks.md and change this task from [ ] to [-]. After completing, change from [-] to [x]_

## Task Execution Notes

### Sequential Dependencies

**Foundation (Tasks 1-4):**
- Task 1 (pkg/maven/version.go) must complete first (used by many components)
- Task 2 (pkg/maven/paths.go) must complete first (used by installer, lister)
- Task 3 (internal/download) must complete first (used by installer)
- Task 4 (internal/config) must complete first (used by resolver, manager)

**Core Components (Tasks 5-9):**
- Tasks 1-4 must complete before starting tasks 5-9
- Task 5 (VersionResolver) can start after task 4
- Task 6 (VersionInstaller) can start after tasks 2-3
- Task 7 (VersionLister) can start after tasks 1-2
- Task 8 (VersionCache) can start after task 4
- Task 9 (VersionManager) must wait for tasks 5-8 to complete

**Integration (Task 10):**
- Task 10 requires task 9 and all cli-commands spec tasks to be complete

**Testing (Tasks 11-20):**
- Unit tests (11-16) can start after corresponding component completes
- Integration tests (17-18) require components they test
- Benchmarks (19) can start after components complete
- E2E tests (20) require task 10 (CLI integration) to be complete

**Finalization (Tasks 21-22):**
- Task 21 (errors) can be done anytime (preferably early)
- Task 22 must be last (final review)

### Parallel Execution Opportunities

**Group 1: Foundation (Tasks 1-4):**
- Tasks 1 and 2 can be done in parallel (both in pkg/maven)
- Task 3 (download) can be done in parallel with 1-2
- Task 4 (config) can be done in parallel with 1-3

**Group 2: Core Components (Tasks 5-8):**
- After tasks 1-4 complete, tasks 5-8 can be done in parallel
- Each component is independent

**Group 3: Unit Tests (Tasks 11-16):**
- Can be done in parallel after corresponding components complete
- Each test file is independent

**Group 4: Integration Tests (Tasks 17-19):**
- Can be done in parallel after components they test are complete

### Integration with Other Specs

**cli-commands spec:**
- Task 10 integrates with cli-commands implementations
- CLI command placeholders replaced with VersionManager calls

**nexus-repository-integration spec:**
- VersionLister.ListAvailable() will use Nexus repository provider
- Installer.getDownloadURL() will query Nexus for URLs

**shim-system-implementation spec:**
- Shim will call VersionResolver.ResolveVersion() to determine which Maven to execute

### Performance Targets

**Version Resolution (<100ms):**
- Optimize file system checks (cache directory scans)
- Use efficient .maven-version search (stop at first match)

**List Installed (<100ms):**
- Single directory scan with validation
- Sort in memory (not slow for typical version count)

**Installation:**
- Network-bound, 30-60s acceptable
- Focus on reliability over speed (atomic operations)
