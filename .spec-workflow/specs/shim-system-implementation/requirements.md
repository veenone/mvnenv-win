# Requirements Document

## Introduction

The Shim System Implementation provides transparent Maven command interception, routing commands to the appropriate Maven version based on the current context (shell/local/global settings). The shim system acts as a lightweight proxy between the user's command and the actual Maven installation, enabling seamless version switching without requiring manual PATH updates or environment variable changes.

This spec focuses on shim generation, command interception, version resolution integration, and command execution with proper stdin/stdout/stderr handling. It integrates with core-version-management (for version resolution) and cli-commands (for the `rehash` command that regenerates shims).

## Alignment with Product Vision

**Product Principle: User Experience First**
The requirement that shims be transparent and require no manual intervention directly supports "making Maven version management invisible to the user."

**Product Principle: Performance Matters**
The requirement for <50ms shim execution overhead ensures "imperceptibly fast" command execution that maintains developer flow.

**Key Feature: Shim System**
This spec directly implements the core product feature "Transparent command interception that routes Maven commands to the correct version automatically."

**Business Objective: Developer Productivity**
Automatic version routing eliminates manual context switching, reducing time spent on version management to zero during normal workflow.

## Requirements

### Requirement 1: Shim Generation

**User Story:** As a developer, I want shim executables automatically generated for all Maven commands, so that I can use `mvn` without specifying full paths.

#### Acceptance Criteria

1. WHEN shims generated THEN system SHALL create shim for `mvn` command (mvn.exe, mvn.cmd)
2. WHEN shims generated THEN system SHALL create shim for `mvnDebug` command (mvnDebug.cmd)
3. WHEN shims generated THEN system SHALL create shim for `mvnyjp` command (mvnyjp.cmd) if it exists in any installed version
4. WHEN generating shims THEN system SHALL place in `%USERPROFILE%\.mvnenv\shims\` directory
5. WHEN generating shims THEN system SHALL create both .exe and .cmd versions for compatibility
6. IF shims directory doesn't exist THEN system SHALL create it automatically
7. WHEN generating shims THEN system SHALL overwrite existing shims
8. WHEN generation completes THEN system SHALL verify shims are executable
9. WHEN no Maven versions installed THEN system SHALL generate shims anyway (will error when invoked)
10. WHEN shim generation fails THEN system SHALL provide clear error with reason

### Requirement 2: Automatic Shim Regeneration

**User Story:** As a developer, I want shims automatically regenerated after installing/uninstalling Maven versions, so that I don't have to remember to run rehash manually.

#### Acceptance Criteria

1. WHEN Maven version installed THEN system SHALL automatically regenerate shims
2. WHEN Maven version uninstalled THEN system SHALL automatically regenerate shims
3. WHEN global version set THEN system SHALL NOT regenerate shims (not needed)
4. WHEN local version set THEN system SHALL NOT regenerate shims (not needed)
5. IF automatic regeneration fails THEN system SHALL log warning but complete installation
6. WHEN user runs `mvnenv rehash` THEN system SHALL manually trigger shim regeneration
7. WHEN rehash completes THEN system SHALL output "Shims regenerated successfully"

### Requirement 3: Version Resolution in Shim

**User Story:** As a developer, I want shims to automatically detect the correct Maven version, so that I don't have to manually switch versions.

#### Acceptance Criteria

1. WHEN shim invoked THEN system SHALL call VersionResolver to determine active version
2. WHEN resolving THEN system SHALL follow shell > local > global hierarchy
3. IF version resolved THEN system SHALL verify version is installed
4. IF version not installed THEN system SHALL return error "Maven version 'X' is set but not installed. Install with: mvnenv install X"
5. IF no version set THEN system SHALL return error "No Maven version set. Set global version with: mvnenv global <version>"
6. WHEN resolution fails THEN system SHALL exit with non-zero status (1)
7. WHEN resolution succeeds THEN system SHALL construct path to Maven binary
8. WHEN multiple shims invoked concurrently THEN resolution SHALL work correctly (thread-safe)

### Requirement 4: Command Execution and Pass-through

**User Story:** As a developer, I want Maven commands to execute exactly as if I called the version directly, so that my workflow isn't affected.

#### Acceptance Criteria

1. WHEN shim executes Maven THEN system SHALL pass all command-line arguments unchanged
2. WHEN Maven writes to stdout THEN system SHALL forward to shim's stdout (no buffering)
3. WHEN Maven writes to stderr THEN system SHALL forward to shim's stderr (no buffering)
4. WHEN Maven reads from stdin THEN system SHALL forward shim's stdin (no buffering)
5. WHEN Maven exits THEN system SHALL exit with same exit code
6. WHEN Maven signal received (Ctrl+C) THEN system SHALL forward signal to Maven process
7. WHEN Maven process killed THEN shim SHALL clean up and exit
8. WHEN Maven execution fails THEN system SHALL preserve Maven's error output
9. WHEN executing THEN system SHALL set MAVEN_HOME environment variable to version directory
10. WHEN executing THEN system SHALL inherit current process environment variables

### Requirement 5: Performance Requirements

**User Story:** As a developer, I want shims to have minimal overhead, so that my commands execute quickly.

#### Acceptance Criteria

1. WHEN shim invoked THEN version resolution SHALL complete in <50ms (excluding Maven execution)
2. WHEN shim invoked THEN total overhead from shim launch to Maven launch SHALL be <50ms
3. WHEN measuring overhead THEN system SHALL exclude Maven binary startup time
4. WHEN measuring overhead THEN system SHALL include: shim launch, version resolution, Maven process spawn
5. WHEN shim launched THEN system SHALL not perform network operations
6. WHEN shim launched THEN system SHALL not read large files (only .maven-version and config)
7. WHEN under load THEN concurrent shim invocations SHALL not slow each other down

### Requirement 6: Error Handling and Diagnostics

**User Story:** As a developer, I want clear error messages when shims fail, so that I can quickly resolve issues.

#### Acceptance Criteria

1. WHEN version resolution fails THEN error SHALL include which resolution step failed
2. WHEN version not installed THEN error SHALL suggest installation command
3. WHEN no version set THEN error SHALL suggest setting global version
4. WHEN Maven binary not found THEN error SHALL include expected path
5. WHEN Maven execution fails to start THEN error SHALL include OS error details
6. IF environment variable MVNENV_DEBUG=1 THEN shim SHALL output verbose diagnostics:
   - Resolved version and source (shell/local/global)
   - Maven binary path
   - MAVEN_HOME value
   - Command-line arguments
   - Execution time
7. WHEN debug mode enabled THEN diagnostics SHALL go to stderr (not stdout)
8. WHEN error occurs THEN shim SHALL exit with appropriate error code (1 for resolution, Maven's code for execution)

### Requirement 7: Windows Compatibility

**User Story:** As a developer, I want shims to work in all Windows shells, so that I have a consistent experience.

#### Acceptance Criteria

1. WHEN shim invoked from PowerShell THEN execution SHALL work correctly
2. WHEN shim invoked from Command Prompt (cmd.exe) THEN execution SHALL work correctly
3. WHEN shim invoked from PowerShell Core (pwsh) THEN execution SHALL work correctly
4. WHEN shim invoked from Git Bash (MINGW) THEN execution SHALL work correctly
5. WHEN shim invoked from Windows Terminal THEN execution SHALL work correctly
6. WHEN shim is .exe THEN it SHALL be preferred over .cmd in PATH resolution
7. WHEN shim is .cmd THEN it SHALL work for cmd.exe and PowerShell
8. WHEN paths contain spaces THEN command execution SHALL work correctly
9. WHEN environment variables contain Unicode THEN handling SHALL be correct
10. WHEN working directory contains non-ASCII characters THEN execution SHALL work correctly

### Requirement 8: PATH Management

**User Story:** As a developer, I want the shims directory automatically added to my PATH, so that Maven commands are available everywhere.

#### Acceptance Criteria

1. WHEN mvnenv installed THEN installation script SHALL add `%USERPROFILE%\.mvnenv\shims` to user PATH
2. WHEN adding to PATH THEN system SHALL place shims directory at beginning (highest priority)
3. WHEN adding to PATH THEN system SHALL check if already present (no duplicates)
4. WHEN PATH modified THEN system SHALL update registry (HKCU\Environment)
5. WHEN PATH modified THEN changes SHALL take effect in new terminal sessions
6. WHEN PATH modified THEN system SHALL log old and new PATH values
7. IF PATH modification fails THEN system SHALL provide manual instructions
8. WHEN checking PATH THEN system SHALL handle case-insensitive comparison (Windows)
9. WHEN shims directory in PATH THEN Maven commands SHALL resolve to shims, not other installations

### Requirement 9: MAVEN_HOME Management

**User Story:** As a developer, I want MAVEN_HOME automatically set to match the active version, so that Maven and Java build tools work correctly.

#### Acceptance Criteria

1. WHEN shim executes Maven THEN system SHALL set MAVEN_HOME environment variable
2. WHEN setting MAVEN_HOME THEN value SHALL be `%USERPROFILE%\.mvnenv\versions\{version}\`
3. WHEN setting MAVEN_HOME THEN it SHALL only affect Maven subprocess, not shim process
4. WHEN Maven subprocess exits THEN MAVEN_HOME SHALL not persist in parent shell
5. WHEN user has MAVEN_HOME set THEN shim SHALL override with correct version's path
6. WHEN MAVEN_HOME set THEN Maven and tools reading it SHALL use correct version

### Requirement 10: Shim Updates and Maintenance

**User Story:** As a developer, I want shims to stay up-to-date with mvnenv changes, so that I always have the latest functionality.

#### Acceptance Criteria

1. WHEN mvnenv upgraded THEN shims SHALL be automatically regenerated
2. WHEN shim format changes THEN `mvnenv rehash` SHALL update to new format
3. WHEN shim invoked THEN it SHALL check if mvnenv binary is available
4. IF mvnenv binary missing THEN shim SHALL error with helpful message
5. WHEN detecting stale shim THEN system SHALL suggest running `mvnenv rehash`
6. WHEN shim executable corrupted THEN regeneration SHALL fix it

## Non-Functional Requirements

### Code Architecture and Modularity

- **Shim Separation**: Shim executable in `cmd/shim/`, generation logic in `internal/shim/`
- **Single Responsibility**: Shim only does resolution + execution, no version management
- **Integration Points**: Uses VersionResolver from core-version-management
- **Template-Based**: Shim generation uses templates for easy updates

### Performance

- **Version Resolution**: <25ms from shim start to resolution complete
- **Process Spawn**: <25ms from resolution to Maven process start
- **Total Overhead**: <50ms combined (resolution + spawn)
- **Memory Footprint**: <10MB for shim process
- **Startup Time**: <50ms from shim launch to Maven launch
- **No I/O Blocking**: Stream forwarding with no buffering delays

### Security

- **Path Validation**: Sanitize resolved Maven binary path to prevent injection
- **Argument Pass-through**: Pass arguments unchanged, no interpretation
- **No Privilege Escalation**: Shims run with user's permissions only
- **Environment Isolation**: MAVEN_HOME only affects Maven subprocess
- **No Sensitive Data**: Shims don't access credentials or sensitive config

### Reliability

- **Crash Resistance**: Shim crashes don't affect system state
- **Signal Handling**: Properly handle Ctrl+C, SIGTERM, SIGKILL
- **Process Cleanup**: Ensure Maven subprocess terminated on shim exit
- **Concurrent Safety**: Multiple shims can run simultaneously without conflict
- **Atomic Generation**: Shim regeneration is atomic (temp file + rename)

### Usability

- **Transparent Operation**: Users shouldn't notice shims exist
- **Clear Errors**: All error messages include suggested actions
- **Debug Mode**: MVNENV_DEBUG provides detailed diagnostics
- **No Manual Steps**: Automatic regeneration after install/uninstall
- **Universal Compatibility**: Works in all Windows shells

### Compatibility

- **Windows Shells**: PowerShell, cmd.exe, PowerShell Core, Git Bash, Windows Terminal
- **Windows Versions**: Windows 10+, Windows Server 2019+
- **Path Lengths**: Handle long paths (MAX_PATH considerations)
- **Unicode Support**: Proper handling of non-ASCII in paths and arguments
- **Exit Codes**: Preserve Maven's exit codes exactly

### Maintainability

- **Logging**: Log version resolution and execution at DEBUG level
- **Error Context**: Wrap errors with shim context
- **Test Coverage**: >90% coverage for shim generation and execution logic
- **Documentation**: Godoc comments for all exported functions

## Technical Constraints

- **Go Version**: Use Go 1.21+ standard library
- **No External Dependencies**: Shim should be pure Go (no C dependencies)
- **Windows-Specific**: Can use Windows-specific code in `*_windows.go` files
- **Single Binary**: Shim is a single standalone executable
- **Process Execution**: Use `os/exec` for Maven subprocess
- **PATH Format**: Use Windows path conventions (backslashes, semicolon separator)

## Out of Scope

The following are explicitly NOT part of this spec:
- Maven wrapper (mvnw) detection or integration (future feature)
- IDE integration or plugins (separate tooling)
- Remote execution or SSH command forwarding
- Custom Maven arguments or preprocessing
- Maven settings.xml management (future feature)
- Multi-user or system-wide shim installation

## Success Criteria

1. Shims generated for mvn, mvnDebug, mvnyjp commands
2. Shims automatically regenerated after install/uninstall operations
3. Version resolution works correctly using shell > local > global hierarchy
4. Maven commands execute with all arguments and I/O pass-through
5. Exit codes preserved exactly from Maven process
6. MAVEN_HOME set correctly for Maven subprocess
7. Performance: <50ms overhead from shim to Maven execution
8. Works correctly in all Windows shells (PowerShell, cmd.exe, Git Bash)
9. Clear error messages with suggested actions
10. MVNENV_DEBUG provides detailed diagnostics
11. >90% test coverage for shim generation and execution
12. Concurrent shim invocations work without conflicts
