# Requirements Document

## Introduction

The CLI Commands feature provides the foundational command-line interface framework for mvnenv-win, directly inspired by pyenv-win's proven command structure. It establishes the command parsing, help system, and output formatting that all other features will build upon. Using the Cobra library, this feature creates an intuitive CLI experience that will feel immediately familiar to pyenv-win users.

This feature focuses on the CLI framework infrastructure—command registration, help generation, error handling, and flags—rather than the business logic of specific commands (version management, repository configuration, etc.), which will be implemented in separate specs.

**Design Philosophy from pyenv-win:**
- Simple, clean text output without decorative elements
- Consistent command naming and argument patterns
- Minimal but helpful error messages
- Fast command execution
- No emojis or fancy formatting in stdout

## Alignment with Product Vision

**Product Principle: Intuitive by Design**
By mirroring pyenv-win's command structure exactly, users familiar with pyenv-win will have zero learning curve. Commands like `mvnenv install -l`, `mvnenv global 3.9.4`, and `mvnenv versions` work exactly as expected.

**Product Principle: User Experience First**
The pyenv-win command structure has been battle-tested by thousands of users. Adopting this proven UX pattern ensures immediate usability and reduces friction.

**Business Objective: Developer Productivity**
Familiar commands mean developers can start using mvnenv-win immediately without reading documentation, directly supporting the goal of reducing time spent on version management.

**Appendix B: Migration Path from pyenv-win**
The PRD explicitly references making mvnenv-win commands intuitive for pyenv-win users, making command alignment a core requirement.

## Requirements

### Requirement 1: Command Discovery

**User Story:** As a developer, I want to discover all available commands, so that I can understand what mvnenv-win can do.

#### Acceptance Criteria

1. WHEN user runs `mvnenv commands` THEN system SHALL list all available commands (one per line)
2. WHEN user runs `mvnenv` without arguments THEN system SHALL display usage help with command list
3. WHEN user runs `mvnenv --help` THEN system SHALL display detailed help with all commands and descriptions
4. WHEN user runs `mvnenv -h` THEN system SHALL display same output as `mvnenv --help`

### Requirement 2: Version Display Commands

**User Story:** As a developer, I want to see my current Maven version and what versions are installed, so that I can verify my environment.

#### Acceptance Criteria

1. WHEN user runs `mvnenv version` THEN system SHALL display currently active Maven version
2. WHEN user runs `mvnenv version` and no version is set THEN system SHALL display "No Maven version set"
3. WHEN user runs `mvnenv versions` THEN system SHALL list all installed Maven versions (one per line)
4. WHEN user runs `mvnenv versions` THEN system SHALL indicate currently active version with `*` prefix
5. WHEN user runs `mvnenv --version` THEN system SHALL display mvnenv-win version in format "mvnenv-win X.Y.Z"

### Requirement 3: Version Selection Commands

**User Story:** As a developer, I want to set Maven versions at different scopes (global, local, shell), so that I can control which version is used in different contexts.

#### Acceptance Criteria

1. WHEN user runs `mvnenv global <version>` THEN system SHALL set global Maven version to specified version
2. WHEN user runs `mvnenv global` without arguments THEN system SHALL display current global version
3. WHEN user runs `mvnenv local <version>` THEN system SHALL set local Maven version (create `.maven-version` file)
4. WHEN user runs `mvnenv local` without arguments THEN system SHALL display current local version
5. WHEN user runs `mvnenv shell <version>` THEN system SHALL set shell-specific Maven version for current session
6. WHEN user runs `mvnenv shell` without arguments THEN system SHALL display current shell version
7. WHEN user runs `mvnenv shell --unset` THEN system SHALL clear shell-specific version

### Requirement 4: Installation and Listing Commands

**User Story:** As a developer, I want to install and uninstall Maven versions, so that I can manage what's available on my system.

#### Acceptance Criteria

1. WHEN user runs `mvnenv install -l` THEN system SHALL list all available Maven versions from configured sources
2. WHEN user runs `mvnenv install <version>` THEN system SHALL install specified Maven version
3. WHEN user runs `mvnenv install <version1> <version2>` THEN system SHALL install multiple versions
4. WHEN user runs `mvnenv install -q <version>` THEN system SHALL install quietly (suppress non-error output)
5. WHEN user runs `mvnenv uninstall <version>` THEN system SHALL uninstall specified Maven version
6. WHEN user runs `mvnenv uninstall <version1> <version2>` THEN system SHALL uninstall multiple versions

### Requirement 5: Utility Commands

**User Story:** As a developer, I want utility commands for cache updates and path resolution, so that I can maintain my mvnenv-win installation.

#### Acceptance Criteria

1. WHEN user runs `mvnenv update` THEN system SHALL update cached version database from configured repositories
2. WHEN user runs `mvnenv rehash` THEN system SHALL rebuild shim executables
3. WHEN user runs `mvnenv which mvn` THEN system SHALL display full path to Maven executable being used
4. WHEN user runs `mvnenv which <command>` THEN system SHALL display path to specified command in current Maven
5. WHEN user runs `mvnenv latest` THEN system SHALL display latest installed Maven version
6. WHEN user runs `mvnenv latest <prefix>` THEN system SHALL display latest version matching prefix (e.g., `mvnenv latest 3.8`)

### Requirement 6: Help System

**User Story:** As a developer, I want help for each command, so that I can understand usage without leaving the terminal.

#### Acceptance Criteria

1. WHEN user runs `mvnenv help` THEN system SHALL display same output as `mvnenv --help`
2. WHEN user runs `mvnenv help <command>` THEN system SHALL display detailed help for that command
3. WHEN user runs `mvnenv <command> --help` THEN system SHALL display detailed help for that command
4. WHEN user runs `mvnenv <command> -h` THEN system SHALL display same help as `--help`
5. WHEN displaying command help THEN system SHALL include: usage syntax, description, available flags, and examples
6. WHEN user runs invalid command THEN system SHALL suggest similar valid command if one exists

### Requirement 7: Global Flags

**User Story:** As a developer, I want flags that control output behavior, so that I can use mvnenv-win in scripts and interactive sessions.

#### Acceptance Criteria

1. WHEN user adds `-q` or `--quiet` flag to any command THEN system SHALL suppress all non-error output
2. WHEN user runs `mvnenv <command> --help` THEN system SHALL display command-specific help
3. WHEN user runs `mvnenv --version` THEN system SHALL display mvnenv-win version and exit
4. WHEN command supports `-l` flag THEN it SHALL list available items (versions, commands, etc.)

### Requirement 8: Error Handling

**User Story:** As a developer, I want clear error messages, so that I can quickly understand and fix problems.

#### Acceptance Criteria

1. WHEN command fails THEN system SHALL display error message starting with "Error: "
2. WHEN error is recoverable THEN system SHALL include suggested action
3. WHEN command fails THEN system SHALL exit with non-zero status code
4. WHEN command succeeds THEN system SHALL exit with status code 0
5. WHEN user provides invalid version argument THEN system SHALL display "Error: version '<version>' not installed"
6. WHEN user provides invalid command THEN system SHALL display "mvnenv: no such command '<command>'"

### Requirement 9: Output Formatting

**User Story:** As a developer, I want simple, consistent output, so that I can easily parse results in scripts or read them interactively.

#### Acceptance Criteria

1. WHEN commands list items THEN system SHALL format output as plain text (one item per line)
2. WHEN displaying current version THEN system SHALL output version number only (e.g., "3.9.4")
3. WHEN displaying paths THEN system SHALL use Windows-style backslashes
4. WHEN `mvnenv versions` shows current version THEN system SHALL prefix with `* ` (asterisk and space)
5. WHEN output goes to stdout THEN system SHALL NOT include emojis, symbols (except `*` for current), or decorative elements
6. WHEN displaying errors THEN system SHALL write to stderr (not stdout)
7. WHEN command has no output (e.g., successful `mvnenv global 3.9.4`) THEN system SHALL output nothing and exit 0

### Requirement 10: Command Execution Performance

**User Story:** As a developer, I want fast command execution, so that CLI interactions don't interrupt my workflow.

#### Acceptance Criteria

1. WHEN user runs `mvnenv --version` THEN system SHALL respond in <50ms
2. WHEN user runs `mvnenv version` THEN system SHALL respond in <100ms
3. WHEN user runs `mvnenv versions` THEN system SHALL respond in <100ms
4. WHEN user runs `mvnenv --help` THEN system SHALL respond in <100ms

## Non-Functional Requirements

### Code Architecture and Modularity

- **Single Responsibility Principle**: Each command file contains only that command's implementation
- **Modular Design**: CLI framework isolated in `cmd/mvnenv/` with business logic in `internal/` packages
- **Dependency Management**: Commands shall not directly import each other; shared logic goes in internal packages
- **Clear Interfaces**: Commands interact with business logic through well-defined internal package interfaces
- **Cobra Integration**: Follow Cobra best practices; each command is a separate `cobra.Command` instance

### Performance

- **Command Startup**: CLI startup overhead must be <50ms from process launch to command execution
- **Help Generation**: Help text generation must not exceed 100ms
- **Memory Footprint**: CLI process memory usage must not exceed 20MB for info commands (version, help, versions)
- **Configuration Loading**: Lazy-load configuration only when needed by specific commands

### Security

- **Input Validation**: All command arguments and flag values must be validated before use
- **Path Sanitization**: File paths from user input must be sanitized to prevent directory traversal
- **No Credential Exposure**: Error messages and help text must never display credentials or sensitive config values
- **Safe Flag Parsing**: Flag parsing must prevent injection attacks through malformed input

### Reliability

- **Graceful Degradation**: If config file is missing or invalid, info commands (help, version, commands) must still work
- **Error Recovery**: Invalid commands must not crash; display helpful error and exit cleanly
- **Consistent Exit Codes**: Use standard exit codes (0=success, 1=general error, 2=usage error)
- **Signal Handling**: Handle Ctrl+C gracefully without leaving terminal in corrupted state

### Usability

- **Intuitive Commands**: Command names must match pyenv-win exactly where applicable (install, uninstall, global, local, shell, versions, version, rehash, update, which, latest, commands)
- **Consistent Naming**: Use single-word commands where possible; multi-word commands rare
- **Short and Long Forms**: Provide short flags (-l, -q, -h) and long equivalents (--help) where appropriate
- **Helpful Suggestions**: When command is mistyped, suggest closest matching valid command using Levenshtein distance
- **Examples in Help**: Include practical examples in help text for each command
- **Minimal Output**: Follow pyenv-win pattern of minimal output; silence on success for setter commands

### Compatibility

- **PowerShell Support**: All commands must work identically in PowerShell 5.1, PowerShell 7+, and Command Prompt
- **Windows Terminal**: Output formatting must display correctly in Windows Terminal
- **Redirected Output**: When output is redirected (piped), continue plain text output (already no colors by default)
- **Unicode Support**: Handle UTF-8 paths correctly for international characters

### Maintainability

- **Command Registration**: New commands registrable by adding single file in `cmd/mvnenv/` directory
- **Centralized Configuration**: Global flags and config handling in single root command
- **Test Coverage**: All command parsing, help generation, and error handling must have unit tests (>90% coverage)
- **Documentation**: Each command file must include godoc comments explaining purpose and usage

## Technical Constraints

- **Cobra Version**: Use github.com/spf13/cobra v1.8+ for command framework
- **Go Version**: Compatible with Go 1.21+ standard library
- **Exit Codes**: Follow standard conventions (0=success, 1=error, 2=usage error)
- **Output Streams**: Errors to stderr, normal output to stdout
- **No Color Library**: Output is plain text only (no ANSI color codes), matching pyenv-win behavior
- **No Fancy Formatting**: No progress bars, spinners, or decorative elements in output

## Out of Scope

The following are explicitly NOT part of this spec (covered in other specs):
- Business logic for version management (install, uninstall, list versions) - in core-version-management spec
- Repository configuration commands (repo add, repo list) - in nexus-repository-integration spec
- Shim generation and execution logic - in shim-system-implementation spec
- Configuration file parsing and validation - in configuration-management spec (future)
- Maven version resolution algorithms - in core-version-management spec
- Download and installation mechanisms - in core-version-management spec

## Command Reference (Summary)

For clarity, here are all commands this spec defines (framework only; business logic in other specs):

```
mvnenv commands              # List all available commands
mvnenv --version             # Show mvnenv-win version
mvnenv version               # Show current Maven version
mvnenv versions              # List installed Maven versions
mvnenv global [<version>]    # Set/show global Maven version
mvnenv local [<version>]     # Set/show local Maven version
mvnenv shell [<version>]     # Set/show shell Maven version
mvnenv install -l            # List available Maven versions
mvnenv install <version>...  # Install Maven version(s)
mvnenv uninstall <version>... # Uninstall Maven version(s)
mvnenv update                # Update cached version database
mvnenv rehash                # Rebuild shim executables
mvnenv which <command>       # Show path to command executable
mvnenv latest [<prefix>]     # Show latest installed version
mvnenv help [<command>]      # Show help
```

## Success Criteria

1. All commands display help text with `--help` or `-h` flag
2. `mvnenv --version` displays correct version number
3. `mvnenv commands` lists all available commands
4. Invalid commands display error with "no such command" message and suggestion
5. All output is plain text without emojis or decorative elements
6. Error messages go to stderr, normal output to stdout
7. Commands follow pyenv-win naming and argument patterns exactly
8. All info commands (version, versions, commands, help) respond in <100ms
9. Exit codes are consistent (0=success, non-zero=error)
10. 100% test coverage for command parsing and help generation
11. Zero crashes or panics from invalid user input
